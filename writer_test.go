/*
 * Playlist generation tests.
**/

package m3u8

import (
	"fmt"
	"testing"
)

// Create new media playlist.
func TestNewMediaPlaylist(t *testing.T) {
	_, e := NewMediaPlaylist(1, 2)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
}

// Create new media playlist
// Add two segments to media playlist
func TestAddSegmentToMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(1, 2)
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
}

// Create new media playlist with capacity 10 elements
// Try to add 11 segments to media playlist (oversize error)
func TestOverAddSegmentsToMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(1, 10)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	for i := 0; i < 11; i++ {
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0)
		if e != nil {
			fmt.Printf("Due oversize new segment #%d not assigned to a media playlist: %s\n", i, e)
		}
	}
}

// Create new media playlist
// Add segment to media playlist
// Set encryption key
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
}

// Create new media playlist
// Add two segments to media playlist
// Encode structures to HLS
func TestEncodeMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	e = p.Add("test01.ts", 5.0)
	if e != nil {
		panic(fmt.Sprintf("Add 1st segment to a media playlist failed: %s", e))
	}
	fmt.Println(p.Encode(true).String())
}

// Create new media playlist
// Add 10 segments to media playlist
// Test iterating over segments
func TestLoopSegmentsOfMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	for i := 0; i < 5; i++ {
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0)
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	for ; e == nil; _, e = p.Next() {
	}
}

// Create new master playlist without params
// Add media playlist
func TestNewMasterPlaylist(t *testing.T) {
	m := NewMasterPlaylist()
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	for i := 0; i < 5; i++ {
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0)
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	m.Add("chunklist1.m3u8", p, VariantParams{})
}

// Create new master playlist with params
// Add media playlist
func TestNewMasterPlaylistWithParams(t *testing.T) {
	m := NewMasterPlaylist()
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	for i := 0; i < 5; i++ {
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0)
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	m.Add("chunklist1.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
}

// Create new master playlist
// Add media playlist
// Encode structures to HLS
func TestEncodeMasterPlaylist(t *testing.T) {
	m := NewMasterPlaylist()
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	for i := 0; i < 5; i++ {
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0)
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	m.Add("chunklist1.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
	m.Add("chunklist2.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
	fmt.Println(m.Encode().String())
}
