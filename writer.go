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

// winsize defines how much items will displayed on playlist generation
// capacity is total size of a playlist
func NewMediaPlaylist(winsize uint16, capacity uint16) (*MediaPlaylist, error) {
	if capacity < winsize {
		return nil, errors.New("capacity must be greater then winsize")
	}
	p := new(MediaPlaylist)
	p.ver = minver
	p.TargetDuration = 0
	p.SeqNo = 0
	p.winsize = winsize
	p.capacity = capacity
	p.segments = make([]*MediaSegment, winsize, capacity)
	return p, nil
}

//
func (p *MediaPlaylist) Add(segment *MediaSegment) error {
	if p.head == p.tail && len(p.segments) > 0 {
		return errors.New("playlist is full")
	}
	if segment.Key.Method != "" { // due section 7
		version(&p.ver, 5)
	}

	p.segments[p.tail] = segment
	p.tail = (p.tail + 1) % p.capacity

	if p.TargetDuration < segment.Duration {
		p.TargetDuration = segment.Duration
	}

	p.buf.Reset()
	return nil
}

// Generate output in HLS. Marshal `winsize` elements from bottom of the `segments` queue.
func (p *MediaPlaylist) Marshal() *bytes.Buffer {
	var key *Key
	var start, end uint16

	if p.buf.Len() > 0 {
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

	head := (p.head + p.winsize) % p.capacity
	if head < p.tail {
		start = head
		end = p.tail
	} else {
		start = p.tail
		end = head
	}
	p.tail = head

	for seg := range p.segments[start:end] {
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
	}
	return p.buf
}

func (p *MediaPlaylist) SetEndlist() *bytes.Buffer {
	p.buf.WriteString("#EXT-X-ENDLIST\n")

	return p.buf
}

func NewMasterPlaylist() *MasterPlaylist {
	p := new(MasterPlaylist)
	p.ver = minver
	return p
}

func (p *MasterPlaylist) Add(variant *Variant) {
	p.variants = append(p.variants, variant)
}

func (p *MasterPlaylist) Marshal() *bytes.Buffer {
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

func (p *MediaPlaylist) SetKey(key *Key) {
	p.key = key
}

func (p *MediaPlaylist) SetWV(wv *WV) {
	p.wv = wv
}
