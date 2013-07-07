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
	"fmt"
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
func NewMediaPlaylist(winsize uint, capacity uint) (*MediaPlaylist, error) {
	if capacity < winsize {
		return nil, errors.New("capacity must be greater then winsize")
	}
	p := new(MediaPlaylist)
	p.ver = minver
	p.winsize = winsize
	p.capacity = capacity
	p.segments = make([]*MediaSegment, capacity)
	return p, nil
}

func (p *MediaPlaylist) Next() (seg *MediaSegment, err error) {
	if p.count == 0 || p.head == p.tail {
		return nil, errors.New("playlist is empty")
	}
	seg = p.segments[p.head]
	p.head = (p.head + 1) % p.capacity
	p.count--
	return seg, nil
}

//
func (p *MediaPlaylist) Add(uri string, duration float64) error {
	if p.head == p.tail && p.count > 0 {
		return errors.New("playlist is full")
	}
	seg := new(MediaSegment)
	seg.URI = uri
	seg.Duration = duration
	p.segments[p.tail] = seg
	p.tail = (p.tail + 1) % p.capacity
	p.count++
	if p.TargetDuration < duration {
		p.TargetDuration = duration
	}
	p.buf.Reset()
	return nil
}

// Generate output in HLS. Marshal `winsize` elements from bottom of the `segments` queue.
func (p *MediaPlaylist) Encode() *bytes.Buffer {
	var err error
	var seg *MediaSegment

	if p.buf.Len() > 0 {
		return &p.buf
	}
	p.SeqNo++
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

	for ; err == nil; seg, err = p.Next() {
		if seg == nil {
			continue
		}
		if seg.Key != nil {
			p.buf.WriteString("#EXT-X-KEY:")
			p.buf.WriteString("METHOD=")
			p.buf.WriteString(seg.Key.Method)
			p.buf.WriteString(",URI=")
			p.buf.WriteString(seg.Key.URI)
			if seg.Key.IV != "" {
				p.buf.WriteString(",IV=")
				p.buf.WriteString(seg.Key.IV)
			}
			p.buf.WriteRune('\n')
		}
		if seg.WV != nil {
			if seg.WV.CypherVersion != "" {
				p.buf.WriteString("#WV-CYPHER-VERSION:")
				p.buf.WriteString(seg.WV.CypherVersion)
				p.buf.WriteRune('\n')
			}
			if seg.WV.ECM != "" {
				p.buf.WriteString("#WV-ECM:")
				p.buf.WriteString(seg.WV.ECM)
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
	return &p.buf
}

func (p *MediaPlaylist) End() bytes.Buffer {
	p.buf.WriteString("#EXT-X-ENDLIST\n")
	return p.buf
}

func (p *MediaPlaylist) Key(method, uri, iv, keyformat, keyformatversions string) error {
	if p.count == 0 {
		return errors.New("playlist is empty")
	}
	if p.head == p.tail && p.count > 0 {
		return errors.New("playlist is full")
	}
	version(&p.ver, 5) // due section 7
	p.segments[(p.tail-1)%p.capacity].Key = &Key{method, uri, iv, keyformat, keyformatversions}
	return nil
}

func NewMasterPlaylist() *MasterPlaylist {
	p := new(MasterPlaylist)
	p.ver = minver
	return p
}

func (p *MasterPlaylist) Add(variant *Variant) error {
	p.variants = append(p.variants, variant)

	return nil
}

func (p *MasterPlaylist) Encode() bytes.Buffer {
	p.buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
	p.buf.WriteString(strver(p.ver))
	p.buf.WriteRune('\n')

	for _, pl := range p.variants {
		p.buf.WriteString("#EXT-X-STREAM-INF:PROGRAM-ID=")
		p.buf.WriteString(strconv.FormatUint(uint64(pl.ProgramId), 10))
		p.buf.WriteString(",BANDWIDTH=")
		p.buf.WriteString(strconv.FormatUint(uint64(pl.Bandwidth), 10))
		if pl.Codecs != "" {
			p.buf.WriteString(",CODECS=")
			p.buf.WriteString(pl.Codecs)
		}
		if pl.Resolution != "" {
			p.buf.WriteString(",RESOLUTION=\"")
			p.buf.WriteString(pl.Resolution)
			p.buf.WriteRune('"')
		}
		p.buf.WriteRune('\n')
		p.buf.WriteString(pl.URI)
		if p.SID != "" {
			p.buf.WriteRune('?')
			p.buf.WriteString(p.SID)
		}
		p.buf.WriteRune('\n')
	}

	return p.buf
}

func dd(vars ...interface{}) {
	print("DEBUG: ")
	for _, msg := range vars {
		fmt.Printf("%v ", msg)
	}
	print("\n")
}
