package m3u8

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func openPlaylist(t *testing.T, file string, listType ListType) (*MasterPlaylist, *MediaPlaylist) {
	f, err := os.Open(file)
	require.NoError(t, err)
	m, typ, err := DecodeFrom(f, true)
	require.NoError(t, err)
	require.Equal(t, listType, typ)

	if listType == MASTER {
		return m.(*MasterPlaylist), nil
	}
	return nil, m.(*MediaPlaylist)
}

func TestDecode(t *testing.T) {
	verifyExpectedSegments := func(t *testing.T, media *MediaPlaylist, expectedSegments []MediaSegment) {
		for i := range expectedSegments {
			require.Equal(t, expectedSegments[i].URI, media.Segments[i].URI)
			require.Equal(t, expectedSegments[i].Duration, media.Segments[i].Duration)
			if expectedSegments[i].SCTE != nil {
				require.NotNil(t, media.Segments[i].SCTE)
				require.Equal(t, expectedSegments[i].SCTE.Elapsed, media.Segments[i].SCTE.Elapsed)
				require.Equal(t, expectedSegments[i].SCTE.Time, media.Segments[i].SCTE.Time)
			}
		}
	}

	t.Run("parse molotov master playlist", func(t *testing.T) {
		master, _ := openPlaylist(t, "test-playlists/molotov-master.m3u8", MASTER)
		assert.EqualValues(t, 6, master.Version())

		require.NotNil(t, master.Key)
		assert.Equal(t, "SAMPLE-AES", master.Key.Method)
		assert.Equal(t, "skd://drmtoday?assetid=vod2live", master.Key.URI)
		assert.Equal(t, "com.apple.streamingkeydelivery", master.Key.Keyformat)
		assert.Equal(t, "1", master.Key.Keyformatversions)

		assert.Len(t, master.Variants, 11)
		expectedVariants := []VariantParams{
			{Bandwidth: 483000, AverageBandwidth: 439000, Codecs: "mp4a.40.2,avc1.640015", Resolution: "426x240", FrameRate: 25, Audio: "audio-aacl-128", Captions: "NONE"},
			{Bandwidth: 788000, AverageBandwidth: 716000, Codecs: "mp4a.40.2,avc1.64001E", Resolution: "640x360", FrameRate: 25, Audio: "audio-aacl-128", Captions: "NONE"},
			{Bandwidth: 1220000, AverageBandwidth: 1109000, Codecs: "mp4a.40.2,avc1.64001F", Resolution: "854x480", FrameRate: 25, Audio: "audio-aacl-128", Captions: "NONE"},
			{Bandwidth: 2253000, AverageBandwidth: 2048000, Codecs: "mp4a.40.2,avc1.64001F", Resolution: "1280x720", FrameRate: 25, Audio: "audio-aacl-128", Captions: "NONE"},
			{Bandwidth: 4903000, AverageBandwidth: 4458000, Codecs: "mp4a.40.2,avc1.640028", Resolution: "1920x1080", FrameRate: 25, Audio: "audio-aacl-128", Captions: "NONE"},
			{Bandwidth: 143000, AverageBandwidth: 130000, Codecs: "mp4a.40.2", Audio: "audio-aacl-128"},
			{Iframe: true, Bandwidth: 43000, Codecs: "avc1.640015", Resolution: "426x240"},
			{Iframe: true, Bandwidth: 81000, Codecs: "avc1.64001E", Resolution: "640x360"},
			{Iframe: true, Bandwidth: 135000, Codecs: "avc1.64001F", Resolution: "854x480"},
			{Iframe: true, Bandwidth: 264000, Codecs: "avc1.64001F", Resolution: "1280x720"},
			{Iframe: true, Bandwidth: 596000, Codecs: "avc1.640028", Resolution: "1920x1080"},
		}
		for i, variant := range master.Variants {
			assert.Equal(t, expectedVariants[i].Iframe, variant.Iframe)
			assert.Equal(t, expectedVariants[i].Bandwidth, variant.Bandwidth)
			assert.Equal(t, expectedVariants[i].AverageBandwidth, variant.AverageBandwidth)
			assert.Equal(t, expectedVariants[i].Codecs, variant.Codecs)
			assert.Equal(t, expectedVariants[i].Resolution, variant.Resolution)
			assert.Equal(t, expectedVariants[i].FrameRate, variant.FrameRate)
			assert.Equal(t, expectedVariants[i].Audio, variant.Audio)
			assert.Equal(t, expectedVariants[i].Subtitles, variant.Subtitles)
			assert.Equal(t, expectedVariants[i].Captions, variant.Captions)
		}

		expectedAlternatives := []Alternative{
			{GroupId: "audio-aacl-128", URI: "publica-audio_fre=128000.m3u8", Type: "AUDIO", Language: "fr", Name: "French", Autoselect: "YES", Default: true},
		}
		alternatives := master.Variants[0].Alternatives
		assert.Len(t, alternatives, 1)
		for i, alternative := range alternatives {
			assert.Equal(t, expectedAlternatives[i].GroupId, alternative.GroupId)
			assert.Equal(t, expectedAlternatives[i].URI, alternative.URI)
			assert.Equal(t, expectedAlternatives[i].Type, alternative.Type)
			assert.Equal(t, expectedAlternatives[i].Language, alternative.Language)
			assert.Equal(t, expectedAlternatives[i].Name, alternative.Name)
			assert.Equal(t, expectedAlternatives[i].Default, alternative.Default)
			assert.Equal(t, expectedAlternatives[i].Autoselect, alternative.Autoselect)
		}
	})

	t.Run("parse molotov media playlist", func(t *testing.T) {
		t.Run("molotov-media-video", func(t *testing.T) {
			_, media := openPlaylist(t, "test-playlists/molotov-media-video.m3u8", MEDIA)
			assert.EqualValues(t, 6, media.Version())
			assert.EqualValues(t, 15, media.Count())
			expectedSegments := []MediaSegment{
				{URI: "testvideo3-video=3914000-601019.ts", Duration: 5.76},
				{URI: "testvideo3-video=3914000-601020.ts", Duration: 5.76},
			}

			for i := range expectedSegments {
				assert.Equal(t, expectedSegments[i].URI, media.Segments[i].URI)
				assert.Equal(t, expectedSegments[i].Duration, media.Segments[i].Duration)
			}

			out := media.String()
			assert.Contains(t, out, "#EXT-X-INDEPENDENT-SEGMENTS")
		})
	})

	t.Run("parse veset media playlist", func(t *testing.T) {
		_, media := openPlaylist(t, "test-playlists/veset-media.m3u8", MEDIA)
		assert.EqualValues(t, 6, media.Version())
		assert.EqualValues(t, 3, media.Count())
		expectedSegments := []MediaSegment{
			{URI: "segment_861.ts", Duration: 6.006, SCTE: &SCTE{Elapsed: 0, Time: 20.020}},
			{URI: "segment_862.ts", Duration: 6.006, SCTE: &SCTE{Elapsed: 6.006, Time: 20.020}},
			{URI: "segment_863.ts", Duration: 6.006, SCTE: &SCTE{Elapsed: 12.012, Time: 20.020}},
		}
		verifyExpectedSegments(t, media, expectedSegments)

		media.IndependentSegments = true
		out := media.String()
		assert.Contains(t, out, "#EXT-X-INDEPENDENT-SEGMENTS")
	})

	t.Run("parse ottera elec_en media playlist", func(t *testing.T) {
		_, media := openPlaylist(t, "test-playlists/ottera-elec-en-media.m3u8", MEDIA)
		assert.EqualValues(t, 6, media.Version())
		assert.EqualValues(t, 5, media.Count())
		expectedSegments := []MediaSegment{
			{URI: "content1.ts", Duration: 6},
			{URI: "content2.ts", Duration: 0.126},
			{URI: "https://ov-static.ottera.tv/scte_v3/elec/en/elec_en_ad_slate_1444_720_high/00000/elec_en_ad_slate_1444_720_high_00001.ts", Duration: 8, SCTE: &SCTE{Elapsed: 0, Time: 120}},
			{URI: "https://ov-static.ottera.tv/scte_v3/elec/en/elec_en_ad_slate_1444_720_high/00000/elec_en_ad_slate_1444_720_high_00002.ts", Duration: 8, SCTE: &SCTE{Elapsed: 10, Time: 120}},
			{URI: "https://ov-static.ottera.tv/scte_v3/elec/en/elec_en_ad_slate_1444_720_high/00000/elec_en_ad_slate_1444_720_high_00003.ts", Duration: 8, SCTE: &SCTE{Elapsed: 18, Time: 120}},
		}
		verifyExpectedSegments(t, media, expectedSegments)
	})

	t.Run("parse newsy media playlist", func(t *testing.T) {
		_, media := openPlaylist(t, "test-playlists/newsy-media.m3u8", MEDIA)
		assert.EqualValues(t, 3, media.Version())
		assert.EqualValues(t, 8, media.Count())
		expectedSegments := []MediaSegment{
			{URI: "re_1634161151_23587.ts", Duration: 4.096},
			{URI: "re_1634161151_23588.ts", Duration: 4.096},
			{URI: "re_1634161151_23589.ts", Duration: 4.096},
			{URI: "re_1634161151_23590.ts", Duration: 4.096, SCTE: &SCTE{Elapsed: 0,Time: 60}},
			{URI: "re_1634161151_23591.ts", Duration: 2.005},
			{URI: "re_1634161151_23592.ts", Duration: 4.096, SCTE: &SCTE{Elapsed: 0,Time: 60}},
			{URI: "re_1634161151_23593.ts", Duration: 4.096},
			{URI: "re_1634161151_23594.ts", Duration: 0.832},
		}
		verifyExpectedSegments(t, media, expectedSegments)
	})
}
