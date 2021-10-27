package m3u8

import (
	"github.com/stretchr/testify/assert"
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

		for _, e := range expected {
			assert.Contains(t, output, e, e)
		}
	})
}
