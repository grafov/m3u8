/*
 * Playlist generation tests.
**/

package m3u8

import (
	"github.com/grafov/m3u8"
	"testing"
)

func TestOkNewFixedPlaylist(t *testing.T) {
	m3u8.NewFixedPlaylist()
}

func TestOkAddSegmentToFixedPlaylist(t *testing.T) {
	p := m3u8.NewFixedPlaylist()
	p.AddSegment(m3u8.Segment{0, "test02.ts", 5.0, nil, nil})
	print(p.Buffer().String())
}

func TestOkNewSlidingPlaylist(t *testing.T) err {
	_, e := m3u8.NewSlidingPlaylist(3, 5)
	return e
}

func TestBadNewSlidingPlaylist(t *testing.T) err {
	_, e := m3u8.NewSlidingPlaylist(5, 3)
	return e
}
