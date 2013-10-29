/*
 Playlist generation tests.

 Copyleft 2013  Alexander I.Grafov aka Axel <grafov@gmail.com>

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package m3u8

import (
	"fmt"
	"testing"
)

// Check how master and media playlists implement common Playlist interface
func TestInterfaceImplemented(t *testing.T) {
	m := NewMasterPlaylist()
	CheckType(m)
	p, e := NewMediaPlaylist(1, 2)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	CheckType(p)
}

// Create new media playlist
// Add two segments to media playlist
func TestAddSegmentToMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(1, 2)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	e = p.Add("test01.ts", 5.0, "")
	if e != nil {
		panic(fmt.Sprintf("Add 1st segment to a media playlist failed: %s", e))
	}
	e = p.Add("test02.ts", 6.0, "")
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
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0, "")
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
	e = p.Add("test01.ts", 5.0, "")
	if e != nil {
		panic(fmt.Sprintf("Add 1st segment to a media playlist failed: %s", e))
	}
	e = p.SetKey("AES-128", "https://example.com", "iv", "format", "vers")
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
	e = p.Add("test01.ts", 5.0, "")
	if e != nil {
		panic(fmt.Sprintf("Add 1st segment to a media playlist failed: %s", e))
	}
	p.DurationAsInt(true)
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
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	p.DurationAsInt(true)
	fmt.Println(p.Encode(true).String())
}

// Create new media playlist with capacity 30
// Add 10 segments to media playlist
// Add encryption key
// Add another 10 segments to media playlist
// Add new encryption key
// Add another 10 segments to media playlist
// Iterate over segments
func TestEncryptionKeysInMediaPlaylist(t *testing.T) {
	// Create new media playlist with capacity 30
	p, e := NewMediaPlaylist(5, 15)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	// Add 10 segments to media playlist
	for i := 0; i < 5; i++ {
		e = p.Add(fmt.Sprintf("test0-%d.ts", i), 6.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	// Add encryption key
	p.SetKey("AES-128", "https://example.com/", "0X00000000000000000000000000000000", "key-format1", "version x.x")
	// Add 10 segments to media playlist
	for i := 0; i < 5; i++ {
		e = p.Add(fmt.Sprintf("test1-%d.ts", i), 6.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	// Add encryption key
	p.SetKey("AES-128", "https://example.com/", "0X00000000000000000000000000000001", "key-format2", "version x.x")
	// Add 10 segments to media playlist
	for i := 0; i < 5; i++ {
		e = p.Add(fmt.Sprintf("test2-%d.ts", i), 6.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	for i := 0; i < 3; i++ {
		fmt.Printf("Iteration %d:\n%s\n", i, p.Encode(true).String())
	}
}

// Create new media playlist
// Add 10 segments to media playlist
// Encode structure to HLS with integer target durations
func TestMediaPlaylistWithIntegerDurations(t *testing.T) {
	p, e := NewMediaPlaylist(3, 10)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	for i := 0; i < 9; i++ {
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.6, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	p.DurationAsInt(false)
	//fmt.Println(p.Encode(true).String())
}

// Create new media playlist as sliding playlist.
// Close it.
func TestClosedMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(1, 10)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	for i := 0; i < 10; i++ {
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			fmt.Printf("Due oversize new segment #%d not assigned to a media playlist: %s\n", i, e)
		}
	}
	p.Close()
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
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0, "")
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
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0, "")
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
		e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	m.Add("chunklist1.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
	m.Add("chunklist2.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
	fmt.Println(m.Encode(true).String())
}
