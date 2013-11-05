package m3u8

/*
 Part of M3U8 parser & generator library.
 This file defines functions related to playlist generation.

 Copyleft 2013  Alexander I.Grafov aka Axel <grafov@gmail.com>

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

 ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"bytes"
	"errors"
	"fmt"
	"math"
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

// Create new empty master playlist.
// Master playlist consists of variants.
func NewMasterPlaylist() *MasterPlaylist {
	p := new(MasterPlaylist)
	p.ver = minver
	return p
}

// Add variant to master playlist
func (p *MasterPlaylist) Add(uri string, chunklist *MediaPlaylist, params VariantParams) {
	v := new(Variant)
	v.URI = uri
	v.Chunklist = chunklist
	v.VariantParams = params
	p.Variants = append(p.Variants, v)
}

func (p *MasterPlaylist) ResetCache() {
	p.buf.Reset()
}

// Generate output in M3U8 format.
func (p *MasterPlaylist) Encode() *bytes.Buffer {
	if p.buf.Len() > 0 {
		return &p.buf
	}

	p.buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
	p.buf.WriteString(strver(p.ver))
	p.buf.WriteRune('\n')

	for _, pl := range p.Variants {
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

	return &p.buf
}

// winsize defines how much items will displayed on playlist generation
// capacity is total size of a playlist
func NewMediaPlaylist(winsize uint, capacity uint) (*MediaPlaylist, error) {
	if capacity < winsize {
		return nil, errors.New("capacity must be greater then winsize or equal")
	}
	p := new(MediaPlaylist)
	p.ver = minver
	p.winsize = winsize
	p.capacity = capacity
	p.Segments = make([]*MediaSegment, capacity)
	return p, nil
}

// Get next segment from the media playlist. Until all segments done.
func (p *MediaPlaylist) next() (err error) {
	if p.count == 0 {
		return errors.New("playlist is empty")
	}
	p.head = (p.head + 1) % p.capacity
	p.count--
	return nil
}

// Add general chunk to media playlist
func (p *MediaPlaylist) Add(uri string, duration float64, title string) error {
	if p.head == p.tail && p.count > 0 {
		return errors.New("playlist is full")
	}
	seg := new(MediaSegment)
	seg.URI = uri
	seg.Duration = duration
	seg.Title = title
	p.Segments[p.tail] = seg
	p.tail = (p.tail + 1) % p.capacity
	p.count++
	if p.TargetDuration < duration {
		p.TargetDuration = duration
	}
	p.buf.Reset()
	return nil
}

func (p *MediaPlaylist) ResetCache() {
	p.buf.Reset()
}

// Generate output in M3U8 format. Marshal `winsize` elements from bottom of the `segments` queue.
func (p *MediaPlaylist) Encode() *bytes.Buffer {
	var err error
	var seg *MediaSegment

	if p.buf.Len() > 0 {
		return &p.buf
	}

	if p.Closed {
		p.SeqNo = 1
	} else {
		p.SeqNo++
	}
	p.buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
	p.buf.WriteString(strver(p.ver))
	p.buf.WriteRune('\n')
	p.buf.WriteString("#EXT-X-ALLOW-CACHE:NO\n")
	// default key
	if p.Key != nil {
		p.buf.WriteString("#EXT-X-KEY:")
		p.buf.WriteString("METHOD=")
		p.buf.WriteString(p.Key.Method)
		p.buf.WriteString(",URI=")
		p.buf.WriteString(p.Key.URI)
		if p.Key.IV != "" {
			p.buf.WriteString(",IV=")
			p.buf.WriteString(p.Key.IV)
		}
		p.buf.WriteRune('\n')
	}
	p.buf.WriteString("#EXT-X-MEDIA-SEQUENCE:")
	p.buf.WriteString(strconv.FormatUint(p.SeqNo, 10))
	p.buf.WriteRune('\n')
	p.buf.WriteString("#EXT-X-TARGETDURATION:")
	p.buf.WriteString(strconv.FormatInt(int64(math.Ceil(p.TargetDuration)), 10)) // due section 3.4.2 of M3U8 specs EXT-X-TARGETDURATION must be integer
	p.buf.WriteRune('\n')
	if p.Closed {
		p.buf.WriteString("#EXT-X-ENDLIST\n")
	}
	// Widevine tags
	if p.WV != nil {
		if p.WV.AudioChannels != 0 {
			p.buf.WriteString("#WV-AUDIO-CHANNELS ")
			p.buf.WriteString(strconv.FormatUint(uint64(p.WV.AudioChannels), 10))
			p.buf.WriteRune('\n')
		}
		if p.WV.AudioFormat != 0 {
			p.buf.WriteString("#WV-AUDIO-FORMAT ")
			p.buf.WriteString(strconv.FormatUint(uint64(p.WV.AudioFormat), 10))
			p.buf.WriteRune('\n')
		}
		if p.WV.AudioProfileIDC != 0 {
			p.buf.WriteString("#WV-AUDIO-PROFILE-IDC ")
			p.buf.WriteString(strconv.FormatUint(uint64(p.WV.AudioProfileIDC), 10))
			p.buf.WriteRune('\n')
		}
		if p.WV.AudioSampleSize != 0 {
			p.buf.WriteString("#WV-AUDIO-SAMPLE-SIZE ")
			p.buf.WriteString(strconv.FormatUint(uint64(p.WV.AudioSampleSize), 10))
			p.buf.WriteRune('\n')
		}
		if p.WV.AudioSamplingFrequency != 0 {
			p.buf.WriteString("#WV-AUDIO-SAMPLING-FREQUENCY ")
			p.buf.WriteString(strconv.FormatUint(uint64(p.WV.AudioSamplingFrequency), 10))
			p.buf.WriteRune('\n')
		}
		if p.WV.CypherVersion != "" {
			p.buf.WriteString("#WV-CYPHER-VERSION ")
			p.buf.WriteString(p.WV.CypherVersion)
			p.buf.WriteRune('\n')
		}
		if p.WV.ECM != "" {
			p.buf.WriteString("#WV-ECM ")
			p.buf.WriteString(p.WV.ECM)
			p.buf.WriteRune('\n')
		}
		if p.WV.VideoFormat != 0 {
			p.buf.WriteString("#WV-VIDEO-FORMAT ")
			p.buf.WriteString(strconv.FormatUint(uint64(p.WV.VideoFormat), 10))
			p.buf.WriteRune('\n')
		}
		if p.WV.VideoFrameRate != 0 {
			p.buf.WriteString("#WV-VIDEO-FRAME-RATE ")
			p.buf.WriteString(strconv.FormatUint(uint64(p.WV.VideoFrameRate), 10))
			p.buf.WriteRune('\n')
		}
		if p.WV.VideoLevelIDC != 0 {
			p.buf.WriteString("#WV-VIDEO-LEVEL-IDC")
			p.buf.WriteString(strconv.FormatUint(uint64(p.WV.VideoLevelIDC), 10))
			p.buf.WriteRune('\n')
		}
		if p.WV.VideoProfileIDC != 0 {
			p.buf.WriteString("#WV-VIDEO-PROFILE-IDC ")
			p.buf.WriteString(strconv.FormatUint(uint64(p.WV.VideoProfileIDC), 10))
			p.buf.WriteRune('\n')
		}
		if p.WV.VideoResolution != "" {
			p.buf.WriteString("#WV-VIDEO-RESOLUTION ")
			p.buf.WriteString(p.WV.VideoResolution)
			p.buf.WriteRune('\n')
		}
		if p.WV.VideoSAR != "" {
			p.buf.WriteString("#WV-VIDEO-SAR ")
			p.buf.WriteString(p.WV.VideoSAR)
			p.buf.WriteRune('\n')
		}
	}

	err = p.next()
	if err != nil {
		p.SeqNo--
		return &p.buf
	}
	head := p.head - 1
	count := p.count + 1

	for i := uint(0); i <= p.winsize && count > 0; count-- {
		seg = p.Segments[head]
		head = (head + 1) % p.capacity
		if seg == nil { // protection from badly filled chunklists
			continue
		}
		if p.winsize > 0 { // skip for VOD playlists, where winsize = 0
			i++
		}
		// key changed
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
		p.buf.WriteString("#EXTINF:")
		if p.durationAsInt {
			// Wowza Mediaserver and some others prefer floats.
			p.buf.WriteString(strconv.FormatFloat(seg.Duration, 'f', 3, 32))
		} else {
			// Old Android players has problems with non integer Duration.
			p.buf.WriteString(strconv.FormatInt(int64(math.Ceil(seg.Duration)), 10))
		}
		p.buf.WriteRune(',')
		p.buf.WriteString(seg.Title)
		p.buf.WriteString("\n")
		p.buf.WriteString(seg.URI)
		if p.SID != "" {
			p.buf.WriteRune('?')
			p.buf.WriteString(p.SID)
		}
		p.buf.WriteString("\n")
	}
	return &p.buf
}

// TargetDuration will be int on Encode
func (p *MediaPlaylist) DurationAsInt(yes bool) {
	if yes {
		// duration must be integers if protocol version is less than 3
		version(&p.ver, 3)
	}
	p.durationAsInt = yes
}

// Close sliding playlist and make them fixed.
func (p *MediaPlaylist) Close() {
	if p.buf.Len() > 0 {
		p.buf.WriteString("#EXT-X-ENDLIST\n")
	}
	p.Closed = true
}

// Set encryption key appeared once in header of the playlist (pointer to MediaPlaylist.Key). It useful when keys not changed during playback.
func (p *MediaPlaylist) SetDefaultKey(method, uri, iv, keyformat, keyformatversions string) {
	version(&p.ver, 5) // due section 7
	p.Key = &Key{method, uri, iv, keyformat, keyformatversions}
}

// Set encryption key for the current segment of media playlist (pointer to Segment.Key)
func (p *MediaPlaylist) SetKey(method, uri, iv, keyformat, keyformatversions string) error {
	if p.count == 0 {
		return errors.New("playlist is empty")
	}
	if p.head == p.tail && p.count > 0 {
		return errors.New("playlist is full")
	}
	version(&p.ver, 5) // due section 7
	p.Segments[(p.tail-1)%p.capacity].Key = &Key{method, uri, iv, keyformat, keyformatversions}
	return nil
}

// Helper. Dumper function for debug.
func dd(vars ...interface{}) {
	print("DEBUG: ")
	for _, msg := range vars {
		fmt.Printf("%v ", msg)
	}
	print("\n")
}
