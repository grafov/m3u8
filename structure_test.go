/*
 * Playlist structures tests.
**/

package m3u8

import (
	"fmt"
	"testing"
)

func CheckType(p Playlist) {
	fmt.Printf("%T implements Playlist interface OK\n", p)
}

// Create new media playlist.
func TestNewMediaPlaylist(t *testing.T) {
	_, e := NewMediaPlaylist(1, 2)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
}
