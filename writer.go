package m3u8

/*
 Part of M3U8 parser & generator library.

 Copyleft Alexander I.Grafov aka Axel <grafov@gmail.com>
 Library licensed under GPLv3

 ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"bytes"
	"strconv"
)

func New() *M3U8 {
	return new(M3U8)
}

func (p *M3U8) AddPlaylist(ProgramId uint8, Bandwidth uint32, URI string) {
	p.Playlists = append(p.Playlists, Playlist{ProgramId, Bandwidth, URI})
}

func (p *M3U8) AddSegment(Duration float64, URI string) {
	p.Segments = append(p.Segments, Segment{Duration, URI})
}

func (p *M3U8) AddPlaylists ([]Playlist) {
}

func (p *M3U8) AddSegments ([]Segment) {
}

func (p *M3U8) Encrypt () {
}

func (p *M3U8) String() string {
	var buf bytes.Buffer

	buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
	buf.WriteString(M3U8Version)
	buf.WriteString("\n")

	for _, pl := range p.Playlists {
		buf.WriteString("#EXT-X-STREAM-INF:PROGRAM-ID=")
		buf.WriteString(strconv.FormatUint(uint64(pl.ProgramId), 10))
		buf.WriteString(",BANDWIDTH=")
		buf.WriteString(strconv.FormatUint(uint64(pl.Bandwidth), 10))
		buf.WriteString("\n")
		buf.WriteString(pl.URI)
		buf.WriteString("\n")
	}

	for _, s := range p.Segments {
		buf.WriteString("#EXTINF:")
		buf.WriteString(strconv.FormatFloat(s.Duration, 'f', 2, 32))
		buf.WriteString("\n")
		buf.WriteString(s.URI)
		buf.WriteString("\n")
	}

	return buf.String()
}

