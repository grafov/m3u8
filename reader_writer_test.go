// Tests for checking how playlists "survive" Decode+Encode

package m3u8

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

// Test for bug where #EXT-X-KEY was duplicated
func TestDecodeEncodeKeySimple(t *testing.T) {
	encoded := encodeDecode(t, "sample-playlists/media-playlist-with-key.m3u8")
	count := strings.Count(encoded, "#EXT-X-KEY")
	if count != 1 {
		t.Errorf("Expected number of EXT-X-KEY: 1 actual: %d", count)
	}
}

// Test that #EXT-X-KEY remains when not adjacent to #EXTINF
func TestDecodeEncodeKeyComplex(t *testing.T) {
	encoded := encodeDecode(t, "sample-playlists/widevine-bitrate.m3u8")
	count := strings.Count(encoded, "#EXT-X-KEY")
	if count != 1 {
		t.Errorf("Expected number of EXT-X-KEY: 1 actual: %d", count)
	}
}

func encodeDecode(t *testing.T, fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	p, _, err := DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		t.Fatal(err)
	}
	pp := p.(*MediaPlaylist)
	return pp.Encode().String()
}
