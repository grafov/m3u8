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
	"time"
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

// Create new media playlist with wrong size (must be failed)
func TestCreateMediaPlaylistWithWrongSize(t *testing.T) {
	_, e := NewMediaPlaylist(2, 1) //wrong winsize
	if e == nil {
		panic("Create new media playlist must be failed, but it's don't")
	}
}

// Create new media playlist
// Add two segments to media playlist
func TestAddSegmentToMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(1, 2)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	e = p.Append("test01.ts", 5.0, "")
	if e != nil {
		panic(fmt.Sprintf("Add 1st segment to a media playlist failed: %s", e))
	}
	e = p.Append("test02.ts", 6.0, "")
	if e != nil {
		panic(fmt.Sprintf("Add 2nd segment to a media playlist failed: %s", e))
	}
}

// Create new media playlist
// Add three segments to media playlist
// Set discontinuity tag for the 2nd segment.
func TestDiscontinuityForMediaPlaylist(t *testing.T) {
	var e error
	p, e := NewMediaPlaylist(3, 4)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	p.Close()
	if e = p.Append("test01.ts", 5.0, ""); e != nil {
		panic(fmt.Sprintf("Add 1st segment to a media playlist failed: %s", e))
	}
	if e = p.Append("test02.ts", 6.0, ""); e != nil {
		panic(fmt.Sprintf("Add 2nd segment to a media playlist failed: %s", e))
	}
	if e = p.SetDiscontinuity(); e != nil {
		panic("Can't set discontinuity tag")
	}
	if e = p.Append("test03.ts", 6.0, ""); e != nil {
		panic(fmt.Sprintf("Add 3nd segment to a media playlist failed: %s", e))
	}
	fmt.Println(p.Encode().String())
}

// Create new media playlist
// Add three segments to media playlist
// Set program date and time for 2nd segment.
// Set discontinuity tag for the 2nd segment.
func TestProgramDateTimeForMediaPlaylist(t *testing.T) {
	var e error
	p, e := NewMediaPlaylist(3, 4)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	p.Close()
	if e = p.Append("test01.ts", 5.0, ""); e != nil {
		panic(fmt.Sprintf("Add 1st segment to a media playlist failed: %s", e))
	}
	if e = p.Append("test02.ts", 6.0, ""); e != nil {
		panic(fmt.Sprintf("Add 2nd segment to a media playlist failed: %s", e))
	}
	loc, _ := time.LoadLocation("Europe/Moscow")
	if e = p.SetProgramDateTime(time.Date(2010, time.November, 30, 16, 25, 0, 0, loc)); e != nil {
		panic("Can't set program date and time")
	}
	if e = p.SetDiscontinuity(); e != nil {
		panic("Can't set discontinuity tag")
	}
	if e = p.Append("test03.ts", 6.0, ""); e != nil {
		panic(fmt.Sprintf("Add 3nd segment to a media playlist failed: %s", e))
	}
	fmt.Println(p.Encode().String())
}

// Create new media playlist
// Add two segments to media playlist with duration 9.0 and 9.1.
// Target duration must be set to nearest greater integer (= 10).
func TestTargetDurationForMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(1, 2)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	e = p.Append("test01.ts", 9.0, "")
	if e != nil {
		panic(fmt.Sprintf("Add 1st segment to a media playlist failed: %s", e))
	}
	e = p.Append("test02.ts", 9.1, "")
	if e != nil {
		panic(fmt.Sprintf("Add 2nd segment to a media playlist failed: %s", e))
	}
	if p.TargetDuration < 10.0 {
		panic(fmt.Sprintf("Target duration must = 10 (nearest greater integer to durations 9.0 and 9.1)"))
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
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
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
	e = p.Append("test01.ts", 5.0, "")
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
	e = p.Append("test01.ts", 5.0, "")
	if e != nil {
		panic(fmt.Sprintf("Add 1st segment to a media playlist failed: %s", e))
	}
	p.DurationAsInt(true)
	fmt.Println(p.Encode().String())
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
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	p.DurationAsInt(true)
	fmt.Println(p.Encode().String())
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
		e = p.Append(fmt.Sprintf("test0-%d.ts", i), 6.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	// Add encryption key
	p.SetKey("AES-128", "https://example.com/", "0X00000000000000000000000000000000", "key-format1", "version x.x")
	// Add 10 segments to media playlist
	for i := 0; i < 5; i++ {
		e = p.Append(fmt.Sprintf("test1-%d.ts", i), 6.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	// Add encryption key
	p.SetKey("AES-128", "https://example.com/", "0X00000000000000000000000000000001", "key-format2", "version x.x")
	// Add 10 segments to media playlist
	for i := 0; i < 5; i++ {
		e = p.Append(fmt.Sprintf("test2-%d.ts", i), 6.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	for i := 0; i < 3; i++ {
		fmt.Printf("Iteration %d:\n%s\n", i, p.Encode().String())
		p.Remove()
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
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.6, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	p.DurationAsInt(false)
	fmt.Println(p.Encode().String())
}

// Create new media playlist
// Add 9 segments to media playlist
// 11 times encode structure to HLS with integer target durations
// Last playlist must be empty
func TestMediaPlaylistWithEmptyMedia(t *testing.T) {
	p, e := NewMediaPlaylist(3, 10)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	for i := 1; i < 10; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.6, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	for i := 1; i < 11; i++ {
		fmt.Println(p.Encode().String())
		p.Remove()
	} // TODO add check for buffers equality
}

// Create new media playlist with winsize == capacity
func TestMediaPlaylistWinsize(t *testing.T) {
	p, e := NewMediaPlaylist(6, 6)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	for i := 1; i < 10; i++ {
		p.Slide(fmt.Sprintf("test%d.ts", i), 5.6, "")
		fmt.Println(p.Encode().String()) // TODO check playlist sizes and mediasequence values
	}
}

// Create new media playlist as sliding playlist.
// Close it.
func TestClosedMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(1, 10)
	if e != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", e))
	}
	for i := 0; i < 10; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
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
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	m.Append("chunklist1.m3u8", p, VariantParams{})
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
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	m.Append("chunklist1.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
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
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	m.Append("chunklist1.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
	m.Append("chunklist2.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
	fmt.Println(m.Encode().String())
}
