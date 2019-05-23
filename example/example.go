package main

import (
	"bufio"
	"fmt"
	"os"
	"path"

	"github.com/grafov/m3u8"
	"github.com/grafov/m3u8/example/template"
)

func main() {
	GOPATH := os.Getenv("GOPATH")
	if GOPATH == "" {
		panic("$GOPATH is empty")
	}

	m3u8File := "github.com/grafov/m3u8/sample-playlists/media-playlist-with-custom-tags.m3u8"
	f, err := os.Open(path.Join(GOPATH, "src", m3u8File))
	if err != nil {
		panic(err)
	}

	customTags := []m3u8.CustomDecoder{
		&template.CustomPlaylistTag{},
		&template.CustomSegmentTag{},
	}

	p, listType, err := m3u8.DecodeWith(bufio.NewReader(f), true, customTags)
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
