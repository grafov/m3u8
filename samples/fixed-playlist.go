package main

import (
	"github.com/grafov/m3u8"
)

func main() {
	p := m3u8.NewFixedPlaylist()
	p.AddSegment(m3u8.Segment{0, "test02.ts?111", 5.0, nil})
	p.AddSegment(m3u8.Segment{0, "test03.ts?111", 6.1, nil})
	print(p.Buffer().String());	print("\n")
}