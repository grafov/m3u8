/*
 * Playlist parsing tests.
**/

package m3u8

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestDecodeMasterPlaylist(t *testing.T) {
	f, err := os.Open("sample-playlists/master.m3u8")
	if err != nil {
		fmt.Println(err)
	}
	p := NewMasterPlaylist()
	err = p.Decode(bufio.NewReader(f), false)
	if err != nil {
		fmt.Println(err)
	}
	// check parsed values
	if p.ver != 3 {
		panic(fmt.Sprintf("Version of parsed playlist = %d (must = 3)", p.ver))
	}
	// TODO check other values
}

func TestDecodeMediaPlaylist(t *testing.T) {
	f, err := os.Open("sample-playlists/wowza-vod-chunklist.m3u8")
	if err != nil {
		panic(err)
	}
	p, err := NewMediaPlaylist(5, 798)
	if err != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", err))
	}
	err = p.Decode(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Playlist object: %+v\n", p)
	// check parsed values
	if p.ver != 3 {
		panic(fmt.Sprintf("Version of parsed playlist = %d (must = 3)", p.ver))
	}
	if p.TargetDuration != 12 {
		panic(fmt.Sprintf("TargetDuration of parsed playlist = %f (must = 12.0)", p.TargetDuration))
	}
	if !p.Closed {
		panic("This is a closed (VOD) playlist but Close field = false")
	}
	// TODO check other values…

	//fmt.Println(p.Encode(true).String())
}

func TestDecodeMasterPlaylistWithCommonDecode(t *testing.T) {
	f, err := os.Open("sample-playlists/master.m3u8")
	if err != nil {
		panic(err)
	}
	m, listType, err := Decode(bufio.NewReader(f), false)
	if err != nil {
		panic(err)
	}
	mp := m.(*MasterPlaylist)
	fmt.Println("Type below must be MasterPlaylist:")
	CheckType(mp)
	if listType != MASTER {
		panic("Sample not recognized as master playlist.")
	}
}

func TestDecodeMediaPlaylistWithCommonDecode(t *testing.T) {
	f, err := os.Open("sample-playlists/wowza-vod-chunklist.m3u8")
	if err != nil {
		panic(err)
	}
	p, listType, err := Decode(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	pp := p.(*MediaPlaylist)
	fmt.Println("Type below must be MediaPlaylist:")
	CheckType(pp)
	if listType != MEDIA {
		panic("Sample not recognized as media playlist.")
	}
	// check parsed values
	if pp.TargetDuration != 12 {
		panic(fmt.Sprintf("TargetDuration of parsed playlist = %f (must = 12.0)", pp.TargetDuration))
	}

	if !pp.Closed {
		panic("This is a closed (VOD) playlist but Close field = false")
	}
	// TODO check other values…
	//fmt.Println(pp.Encode().String())
}
