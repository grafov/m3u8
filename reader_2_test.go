package m3u8

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func openPlaylist(t *testing.T, file string, listType ListType) (*MasterPlaylist, *MediaPlaylist) {
	f, err := os.Open(file)
	assert.NoError(t, err)
	m, typ, err := DecodeFrom(f, true)
	assert.NoError(t, err)
	assert.Equal(t, listType, typ)

	if listType == MASTER {
		return m.(*MasterPlaylist), nil
	}
	return nil, m.(*MediaPlaylist)
}

func TestDecode(t *testing.T) {
	master, _ := openPlaylist(t, "test-playlists/molotov-master.m3u8", MASTER)
	assert.EqualValues(t, 6, master.Version())

	assert.Len(t, master.Variants, 11)
	expectedVariants := []VariantParams{
		{Bandwidth: 485000, AverageBandwidth: 441000, Codecs: "mp4a.40.2,avc1.640015", Resolution: "426x240", FrameRate: 25, Audio: "audio-aacl-128", Subtitles: "textstream", Captions: "NONE"},
		{Bandwidth: 778000, AverageBandwidth: 707000, Codecs: "mp4a.40.2,avc1.64001E", Resolution: "640x360", FrameRate: 25, Audio: "audio-aacl-128", Subtitles: "textstream", Captions: "NONE"},
		{Bandwidth: 1221000, AverageBandwidth: 1110000, Codecs: "mp4a.40.2,avc1.64001F", Resolution: "854x480", FrameRate: 25, Audio: "audio-aacl-128", Subtitles: "textstream", Captions: "NONE"},
		{Bandwidth: 2294000, AverageBandwidth: 2085000, Codecs: "mp4a.40.2,avc1.64001F", Resolution: "1280x720", FrameRate: 25, Audio: "audio-aacl-128", Subtitles: "textstream", Captions: "NONE"},
		{Bandwidth: 4708000, AverageBandwidth: 4280000, Codecs: "mp4a.40.2,avc1.640028", Resolution: "1920x1080", FrameRate: 25, Audio: "audio-aacl-128", Subtitles: "textstream", Captions: "NONE"},
		{Bandwidth: 144000, AverageBandwidth: 131000, Codecs: "mp4a.40.2", Audio: "audio-aacl-128", Subtitles: "textstream"},
		{Iframe: true, Bandwidth: 43000, Codecs: "avc1.640015", Resolution: "426x240"},
		{Iframe: true, Bandwidth: 80000, Codecs: "avc1.64001E", Resolution: "640x360"},
		{Iframe: true, Bandwidth: 135000, Codecs: "avc1.64001F", Resolution: "854x480"},
		{Iframe: true, Bandwidth: 269000, Codecs: "avc1.64001F", Resolution: "1280x720"},
		{Iframe: true, Bandwidth: 571000, Codecs: "avc1.640028", Resolution: "1920x1080"},
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
		{GroupId: "audio-aacl-128", URI: "testvideo3-audio_eng=128000.m3u8", Type: "AUDIO", Language: "en", Name: "English", Default: true, Autoselect: "YES"},
		{GroupId: "audio-aacl-128", URI: "testvideo3-audio_fre=128000.m3u8", Type: "AUDIO", Language: "fr", Name: "French", Autoselect: "YES"},
		{GroupId: "textstream", URI: "testvideo3-textstream_fre=1000.m3u8", Type: "SUBTITLES", Language: "fr", Name: "French", Default: true, Autoselect: "YES"},
	}
	alternatives := master.Variants[0].Alternatives
	assert.Len(t, alternatives, 3)
	for i, alternative := range alternatives {
		assert.Equal(t, expectedAlternatives[i].GroupId, alternative.GroupId)
		assert.Equal(t, expectedAlternatives[i].URI, alternative.URI)
		assert.Equal(t, expectedAlternatives[i].Type, alternative.Type)
		assert.Equal(t, expectedAlternatives[i].Language, alternative.Language)
		assert.Equal(t, expectedAlternatives[i].Name, alternative.Name)
		assert.Equal(t, expectedAlternatives[i].Default, alternative.Default)
		assert.Equal(t, expectedAlternatives[i].Autoselect, alternative.Autoselect)
	}
}
