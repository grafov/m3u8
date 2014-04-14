/*
 Playlist parsing tests.

 Copyleft 2013-2014 Alexander I.Grafov aka Axel <grafov@gmail.com>

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
	err = p.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		fmt.Println(err)
	}
	// check parsed values
	if p.ver != 3 {
		panic(fmt.Sprintf("Version of parsed playlist = %d (must = 3)", p.ver))
	}
	if len(p.Variants) != 5 {
		panic("Not all variants in master playlist parsed.")
	}
	// TODO check other values
	// fmt.Println(p.Encode().String())
}

func TestDecodeMasterPlaylistWithAlternatives(t *testing.T) {
	f, err := os.Open("sample-playlists/master-with-alternatives.m3u8")
	if err != nil {
		fmt.Println(err)
	}
	p := NewMasterPlaylist()
	err = p.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		fmt.Println(err)
	}
	// check parsed values
	if p.ver != 3 {
		panic(fmt.Sprintf("Version of parsed playlist = %d (must = 3)", p.ver))
	}
	// if len(p.Variants) != 5 {
	// 	panic("Not all variants in master playlist parsed.")
	// }
	// TODO check other values
	fmt.Println(p.Encode().String())
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
	err = p.DecodeFrom(bufio.NewReader(f), true)
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
	//fmt.Println(p.Encode().String()), stream.Name}
}

func TestDecodeMediaPlaylistWithWidevine(t *testing.T) {
	f, err := os.Open("sample-playlists/widevine-bitrate.m3u8")
	if err != nil {
		panic(err)
	}
	p, err := NewMediaPlaylist(5, 798)
	if err != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", err))
	}
	err = p.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Playlist object: %+v\n", p)
	// check parsed values
	if p.ver != 2 {
		panic(fmt.Sprintf("Version of parsed playlist = %d (must = 2)", p.ver))
	}
	if p.TargetDuration != 9 {
		panic(fmt.Sprintf("TargetDuration of parsed playlist = %f (must = 9.0)", p.TargetDuration))
	}
	// TODO check other values…
	//fmt.Printf("%+v\n", p.Key)
	fmt.Println(p.Encode().String())
}

func TestDecodeMasterPlaylistWithAutodetection(t *testing.T) {
	print("test")
	f, err := os.Open("sample-playlists/master.m3u8")
	if err != nil {
		panic(err)
	}
	m, listType, err := DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		panic(err)
	}
	if listType != MASTER {
		panic("Sample not recognized as master playlist.")
	}
	mp := m.(*MasterPlaylist)
	// fmt.Printf(">%+v\n", mp)
	// for _, v := range mp.Variants {
	// 	fmt.Printf(">>%+v +v\n", v)
	// }
	fmt.Println("Type below must be MasterPlaylist:")
	CheckType(mp)
}

func TestDecodeMediaPlaylistWithAutodetection(t *testing.T) {
	f, err := os.Open("sample-playlists/wowza-vod-chunklist.m3u8")
	if err != nil {
		panic(err)
	}
	p, listType, err := DecodeFrom(bufio.NewReader(f), true)
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
	// fmt.Println(pp.Encode().String())
}
