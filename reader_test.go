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
	//fmt.Printf("Playlist object: %+v\n", p)
}

func TestDecodeMediaPlaylist(t *testing.T) {
	f, err := os.Open("sample-playlists/wowza-vod-chunklist.m3u8")
	if err != nil {
		fmt.Println(err)
	}
	p, err := NewMediaPlaylist(5, 512)
	if err != nil {
		panic(fmt.Sprintf("Create media playlist failed: %s", err))
	}
	err = p.Decode(bufio.NewReader(f), false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Playlist object: %+v\n", p)
}
