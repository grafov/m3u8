package m3u8

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	t.Run("parse molotov master playlist", func(t *testing.T) {
		master, _ := openPlaylist(t, "test-playlists/molotov-master.m3u8", MASTER)
		output := master.String()

		expected := []string{
			`#EXTM3U`,
			`#EXT-X-VERSION:6`,
			`#EXT-X-SESSION-KEY:METHOD=SAMPLE-AES,URI="skd://drmtoday?assetid=vod2live",KEYFORMAT="com.apple.streamingkeydelivery",KEYFORMATVERSIONS="1"`,
		}

		lines := strings.Split(output, "\n")
		for _, e := range expected {
			assert.Contains(t, lines, e, e)
		}
	})

	t.Run("ottera", func(t *testing.T) {
		master, _ := openPlaylist(t, "test-playlists/ottera-gusto-tv-master.m3u8", MASTER)
		output := master.String()

		expected := []string{
			`#EXTM3U`,
			`#EXT-X-VERSION:3`,
			`#EXT-X-MEDIA:TYPE=CLOSED-CAPTIONS,GROUP-ID="CC",LANGUAGE="eng",NAME="English",INSTREAM-ID="CC1"`,
			`#EXT-X-STREAM-INF:BANDWIDTH=3193687,AVERAGE-BANDWIDTH=3075600,CODECS="avc1.64001f,mp4a.40.2",RESOLUTION=1280x720,CLOSED-CAPTIONS="CC",SUBTITLES="subs",FRAME-RATE=29.970`,
		}

		lines := strings.Split(output, "\n")
		for _, e := range expected {
			require.Contains(t, lines, e, e)
		}
	})
}

func lineIsInOutput(lines []string, expectedLine string) bool {
	for _, line := range lines {
		if line == expectedLine {
			return true
		}
	}
	return false
}
