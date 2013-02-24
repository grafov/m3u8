package main

import (
	"github.com/grafov/m3u8"
)

func main() {
	s := m3u8.NewSlidingPlaylist(4)
	for i := 0; i < 10; i++ {
		s.AddSegment(m3u8.Segment{0, "sample.ts", 5.0, nil})
		if i%5 == 1 {
			print(s.Buffer().String())
			print("\n")
		}
	}
	print(s.Buffer().String())
	print("\n")
	print(s.Buffer().String())
	print("\n")
}
