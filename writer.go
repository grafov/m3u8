package m3u8

/*
 Part of M3U8 parser & generator library.

 Copyleft Alexander I.Grafov aka Axel <grafov@gmail.com>
 Library licensed under GPLv3

 ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"bytes"
	"errors"
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

func NewFixedPlaylist() *FixedPlaylist {
	p := new(FixedPlaylist)
	p.ver = minver
	p.TargetDuration = 0
	return p
}

func (p *FixedPlaylist) AddSegment(segment Segment) {
	p.segments = append(p.segments, segment)
	if segment.Key != nil { // due section 7 of HLS spec
		version(&p.ver, 5)
	}
	if p.TargetDuration < segment.Duration {
		p.TargetDuration = segment.Duration
	}
}

func (p *FixedPlaylist) Buffer() *bytes.Buffer {
	var buf bytes.Buffer

	buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
	buf.WriteString(strver(p.ver))
	buf.WriteRune('\n')
	buf.WriteString("#EXT-X-ALLOW-CACHE:YES\n")
	buf.WriteString("#EXT-X-TARGETDURATION:")
	buf.WriteString(strconv.FormatFloat(p.TargetDuration, 'f', 2, 64))
	buf.WriteRune('\n')
	buf.WriteString("#EXT-X-MEDIA-SEQUENCE:1\n")

	for _, s := range p.segments {
		if s.Key != nil {
			buf.WriteString("#EXT-X-KEY:")
			buf.WriteString("METHOD=")
			buf.WriteString(s.Key.Method)
			buf.WriteString(",URI=")
			buf.WriteString(s.Key.URI)
			if s.Key.IV != "" {
				buf.WriteString(",IV=")
				buf.WriteString(s.Key.IV)
			}
			buf.WriteRune('\n')
		}
		buf.WriteString("#EXTINF:")
		buf.WriteString(strconv.FormatFloat(s.Duration, 'f', 2, 32))
		buf.WriteString("\n")
		buf.WriteString(s.URI)
		if p.SID != "" {
			buf.WriteRune('?')
			buf.WriteString(p.SID)
		}
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
	p.variants = append(p.variants, variant)
}

func (p *VariantPlaylist) Buffer() *bytes.Buffer {
	var buf bytes.Buffer

	buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
	buf.WriteString(strver(p.ver))
	buf.WriteRune('\n')

	for _, pl := range p.variants {
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
		if p.SID != "" {
			buf.WriteRune('?')
			buf.WriteString(p.SID)
		}
		buf.WriteRune('\n')
	}

	return &buf
}

func NewSlidingPlaylist(winsize uint16, capacity uint16) (*SlidingPlaylist, error) {
	if capacity < winsize {
		return nil, errors.New("capacity must be greater then winsize")
	}
	p := new(SlidingPlaylist)
	p.ver = minver
	p.TargetDuration = 0
	p.SeqNo = 0
	p.winsize = winsize
	p.capacity = capacity
	p.segments = make(chan Segment, capacity)
	return p, nil
}

func (p *SlidingPlaylist) AddSegment(segment Segment) error {
	if uint16(len(p.segments)) >= p.winsize*2-1 {
		return errors.New("segments channel is full")
	}
	p.segments <- segment
	if segment.Key.Method != "" { // due section 7
		version(&p.ver, 5)
	}
	if p.TargetDuration < segment.Duration {
		p.TargetDuration = segment.Duration
	}
	p.buf.Reset()

	return nil
}

func (p *SlidingPlaylist) Buffer() *bytes.Buffer {
	var key *Key

	if len(p.segments) == 0 && p.buf.Len() > 0 {
		return p.buf
	}

	p.buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
	p.buf.WriteString(strver(p.ver))
	p.buf.WriteRune('\n')
	p.buf.WriteString("#EXT-X-ALLOW-CACHE:NO\n")
	p.buf.WriteString("#EXT-X-TARGETDURATION:")
	p.buf.WriteString(strconv.FormatFloat(p.TargetDuration, 'f', 2, 64))
	p.buf.WriteRune('\n')
	p.buf.WriteString("#EXT-X-MEDIA-SEQUENCE:")
	p.buf.WriteString(strconv.FormatUint(p.SeqNo, 10))
	p.buf.WriteRune('\n')
	p.SeqNo++

	for i := 0; i <= len(p.segments); i++ {
		select {
		case seg := <-p.segments:
			key = nil
			if seg.Key != nil {
				key = seg.Key
			} else {
				if p.key != nil {
					key = p.key
				}
			}
			if key != nil {
				p.buf.WriteString("#EXT-X-KEY:")
				p.buf.WriteString("METHOD=")
				p.buf.WriteString(key.Method)
				p.buf.WriteString(",URI=")
				p.buf.WriteString(key.URI)
				if key.IV != "" {
					p.buf.WriteString(",IV=")
					p.buf.WriteString(key.IV)
				}
				p.buf.WriteRune('\n')
			}
			if p.wv != nil {
				if p.wv.CypherVersion != "" {
					p.buf.WriteString("#WV-CYPHER-VERSION:")
					p.buf.WriteString(p.wv.CypherVersion)
					p.buf.WriteRune('\n')
				}
				if p.wv.ECM != "" {
					p.buf.WriteString("#WV-ECM:")
					p.buf.WriteString(p.wv.ECM)
					p.buf.WriteRune('\n')
				}
			}
			p.buf.WriteString("#EXTINF:")
			p.buf.WriteString(strconv.FormatFloat(seg.Duration, 'f', 2, 32))
			p.buf.WriteString("\n")
			p.buf.WriteString(seg.URI)
			if p.SID != "" {
				p.buf.WriteRune('?')
				p.buf.WriteString(p.SID)
			}
			p.buf.WriteString("\n")
			// TODO key
		default:
		}
	}
	return p.buf
}

func (p *SlidingPlaylist) Close() *bytes.Buffer {
	p.Buffer()
	p.buf.WriteString("#EXT-X-ENDLIST\n")

	return p.buf
}

func (p *SlidingPlaylist) SetKey(key *Key) {
	p.key = key
}

func (p *SlidingPlaylist) SetWV(wv *WV) {
	p.wv = wv
}
