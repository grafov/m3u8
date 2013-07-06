/*
 * Playlist generation tests.
**/

package m3u8

import (
	"fmt"
	"testing"
)

func TestNewMediaPlaylist(t *testing.T) {
	_, e := NewMediaPlaylist(3, 5)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
}

func TestAddSegmentToMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	e = p.Add("test01.ts", 5.0)
	if e != nil {
		panic(fmt.Sprintf("Add 1st segment to a media playlist failed: %s", e))
	}
	e = p.Add("test02.ts", 6.0)
	if e != nil {
		panic(fmt.Sprintf("Add 2nd segment to a media playlist failed: %s", e))
	}
	p.Encode().String()
}

func TestSetKeyForMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	e = p.Add("test01.ts", 5.0)
	if e != nil {
		panic(fmt.Sprintf("Add 1st segment to a media playlist failed: %s", e))
	}
	e = p.Key("AES", "example.com", "iv", "format", "vers")
	if e != nil {
		panic(fmt.Sprintf("Set key to a media playlist failed: %s", e))
	}
	p.Encode().String()
}

/*
func TestNewMasterPlaylist(t *testing.T) {
	NewMasterPlaylist()
}
*/
