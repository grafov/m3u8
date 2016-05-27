package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/grafov/m3u8"
)

func main() {
	GOPATH := os.Getenv("GOPATH")
	if GOPATH == "" {
		panic("$GOPATH is empty")
	}
	this := "github.com/grafov/m3u8"
	f, err := os.Open(GOPATH + "/src/" + this + "/sample-playlists/media-playlist-with-byterange.m3u8")
	if err != nil {
		panic(err)
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		fmt.Printf("%+v\n", mediapl)
	case m3u8.MASTER:
		masterpl := p.(*m3u8.MasterPlaylist)
		fmt.Printf("%+v\n", masterpl)
	}
}
