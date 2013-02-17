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

func version(ver *uint8, newver uint8) {
	if *ver < newver {
		ver = &newver
	}
}

func strver(ver uint8) string {
	return strconv.FormatUint(uint64(ver), 10)
}

func NewSimplePlaylist() *SimplePlaylist {
	p := new(SimplePlaylist)
	p.ver = minver
	p.TargetDuration = 0
	return p
}

func (p *SimplePlaylist) AddSegment(segment Segment) {
	p.Segments = append(p.Segments, segment)
	if segment.Key != nil { // due section 7
		version(&p.ver, 5)
	}
	if p.TargetDuration < segment.Duration {
		p.TargetDuration = segment.Duration
	}
}

func (p *SimplePlaylist) Buffer() *bytes.Buffer {
	var buf bytes.Buffer

	buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
	buf.WriteString(strver(p.ver))
	buf.WriteRune('\n')
	buf.WriteString("#EXT-X-ALLOW-CACHE:NO\n")
	buf.WriteString("#EXT-X-TARGET-DURATION:")
	buf.WriteString(strconv.FormatFloat(p.TargetDuration, 'f', 2, 64))
	buf.WriteRune('\n')
	//buf.WriteString("#EXT-X-MEDIA-SEQUENCE:0\n")

	for _, s := range p.Segments {
		buf.WriteString("#EXTINF:")
		buf.WriteString(strconv.FormatFloat(s.Duration, 'f', 2, 32))
		buf.WriteString("\n")
		buf.WriteString(s.URI)
		buf.WriteString("\n")
	}

	buf.WriteString("#EXT-X-ENDLIST\n")

	return &buf
}

func NewVariantPlaylist() *VariantPlaylist {
	p := new(VariantPlaylist)
	p.ver = minver
	return p
}

func (p *VariantPlaylist) AddVariant(variant Variant) {
	p.Variants = append(p.Variants, variant)
}

func (p *VariantPlaylist) Buffer() *bytes.Buffer {
	var buf bytes.Buffer

	buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
	buf.WriteString(strver(p.ver))
	buf.WriteRune('\n')

	for _, pl := range p.Variants {
		buf.WriteString("#EXT-X-STREAM-INF:PROGRAM-ID=")
		buf.WriteString(strconv.FormatUint(uint64(pl.ProgramId), 10))
		buf.WriteString(",BANDWIDTH=")
		buf.WriteString(strconv.FormatUint(uint64(pl.Bandwidth), 10))
		if pl.Codecs != "" {
			buf.WriteString(",CODECS=")
			buf.WriteString(pl.Codecs)
		}
		if pl.Resolution != "" {
			buf.WriteString(",RESOLUTION=\"")
			buf.WriteString(pl.Resolution)
			buf.WriteRune('"')
		}
		buf.WriteRune('\n')
		buf.WriteString(pl.URI)
		buf.WriteRune('\n')
	}

	return &buf
}

func NewSlidingPlaylist(winsize uint8) *SlidingPlaylist {
	p := new(SlidingPlaylist)
	p.ver = minver
	p.SeqNo = 0
	p.winsize = winsize
	return p
}

func (p *SlidingPlaylist) AddSegment(segment Segment) {
	p.Segments = append(p.Segments, segment)
	if segment.Key != nil { // due section 7
		version(&p.ver, 5)
	}
}

func NewKey(Method string, IV string, URI string) *Key {
	k := new(Key)
	k = &Key{Method, IV, URI}
	return k
}
