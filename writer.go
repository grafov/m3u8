package m3u8

/*
 Part of M3U8 parser & generator library.
 This file defines functions related to playlist generation.

 Copyright 2013-2017 The Project Developers.
 See the AUTHORS and LICENSE files at the top-level directory of this distribution
 and at https://github.com/grafov/m3u8/

 ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	ErrPlaylistFull = errors.New("playlist is full")
)

// Set version of the playlist accordingly with section 7
func version(ver *uint8, newver uint8) {
	if *ver < newver {
		*ver = newver
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

// Append variant to master playlist.
// This operation does reset playlist cache.
func (p *MasterPlaylist) Append(uri string, chunklist *MediaPlaylist, params VariantParams) {
	v := new(Variant)
	v.URI = uri
	v.Chunklist = chunklist
	v.VariantParams = params
	p.Variants = append(p.Variants, v)
	if len(v.Alternatives) > 0 {
		// From section 7:
		// The EXT-X-MEDIA tag and the AUDIO, VIDEO and SUBTITLES attributes of
		// the EXT-X-STREAM-INF tag are backward compatible to protocol version
		// 1, but playback on older clients may not be desirable.  A server MAY
		// consider indicating a EXT-X-VERSION of 4 or higher in the Master
		// Playlist but is not required to do so.
		version(&p.ver, 4) // so it is optional and in theory may be set to ver.1
		// but more tests required
	}
	p.buf.Reset()
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

	var altsWritten map[string]bool = make(map[string]bool)

	for _, pl := range p.Variants {
		if pl.Alternatives != nil {
			for _, alt := range pl.Alternatives {
				// Make sure that we only write out an alternative once
				altKey := fmt.Sprintf("%s-%s-%s-%s", alt.Type, alt.GroupId, alt.Name, alt.Language)
				if altsWritten[altKey] {
					continue
				}
				altsWritten[altKey] = true

				p.buf.WriteString("#EXT-X-MEDIA:")
				if alt.Type != "" {
					p.buf.WriteString("TYPE=") // Type should not be quoted
					p.buf.WriteString(alt.Type)
				}
				if alt.GroupId != "" {
					p.buf.WriteString(",GROUP-ID=\"")
					p.buf.WriteString(alt.GroupId)
					p.buf.WriteRune('"')
				}
				if alt.Name != "" {
					p.buf.WriteString(",NAME=\"")
					p.buf.WriteString(alt.Name)
					p.buf.WriteRune('"')
				}
				p.buf.WriteString(",DEFAULT=")
				if alt.Default {
					p.buf.WriteString("YES")
				} else {
					p.buf.WriteString("NO")
				}
				if alt.Autoselect != "" {
					p.buf.WriteString(",AUTOSELECT=")
					p.buf.WriteString(alt.Autoselect)
				}
				if alt.Language != "" {
					p.buf.WriteString(",LANGUAGE=\"")
					p.buf.WriteString(alt.Language)
					p.buf.WriteRune('"')
				}
				if alt.Forced != "" {
					p.buf.WriteString(",FORCED=\"")
					p.buf.WriteString(alt.Forced)
					p.buf.WriteRune('"')
				}
				if alt.Characteristics != "" {
					p.buf.WriteString(",CHARACTERISTICS=\"")
					p.buf.WriteString(alt.Characteristics)
					p.buf.WriteRune('"')
				}
				if alt.Subtitles != "" {
					p.buf.WriteString(",SUBTITLES=\"")
					p.buf.WriteString(alt.Subtitles)
					p.buf.WriteRune('"')
				}
				if alt.URI != "" {
					p.buf.WriteString(",URI=\"")
					p.buf.WriteString(alt.URI)
					p.buf.WriteRune('"')
				}
				p.buf.WriteRune('\n')
			}
		}
		if pl.Iframe {
			p.buf.WriteString("#EXT-X-I-FRAME-STREAM-INF:PROGRAM-ID=")
			p.buf.WriteString(strconv.FormatUint(uint64(pl.ProgramId), 10))
			p.buf.WriteString(",BANDWIDTH=")
			p.buf.WriteString(strconv.FormatUint(uint64(pl.Bandwidth), 10))
			if pl.Codecs != "" {
				p.buf.WriteString(",CODECS=\"")
				p.buf.WriteString(pl.Codecs)
				p.buf.WriteRune('"')
			}
			if pl.Resolution != "" {
				p.buf.WriteString(",RESOLUTION=") // Resolution should not be quoted
				p.buf.WriteString(pl.Resolution)
			}
			if pl.Video != "" {
				p.buf.WriteString(",VIDEO=\"")
				p.buf.WriteString(pl.Video)
				p.buf.WriteRune('"')
			}
			if pl.URI != "" {
				p.buf.WriteString(",URI=\"")
				p.buf.WriteString(pl.URI)
				p.buf.WriteRune('"')
			}
			p.buf.WriteRune('\n')
		} else {
			p.buf.WriteString("#EXT-X-STREAM-INF:PROGRAM-ID=")
			p.buf.WriteString(strconv.FormatUint(uint64(pl.ProgramId), 10))
			p.buf.WriteString(",BANDWIDTH=")
			p.buf.WriteString(strconv.FormatUint(uint64(pl.Bandwidth), 10))
			if pl.Codecs != "" {
				p.buf.WriteString(",CODECS=\"")
				p.buf.WriteString(pl.Codecs)
				p.buf.WriteRune('"')
			}
			if pl.Resolution != "" {
				p.buf.WriteString(",RESOLUTION=") // Resolution should not be quoted
				p.buf.WriteString(pl.Resolution)
			}
			if pl.Audio != "" {
				p.buf.WriteString(",AUDIO=\"")
				p.buf.WriteString(pl.Audio)
				p.buf.WriteRune('"')
			}
			if pl.Video != "" {
				p.buf.WriteString(",VIDEO=\"")
				p.buf.WriteString(pl.Video)
				p.buf.WriteRune('"')
			}
			if pl.Captions != "" {
				p.buf.WriteString(",CLOSED-CAPTIONS=")
				if pl.Captions == "NONE" {
					p.buf.WriteString(pl.Captions) // CC should not be quoted when eq NONE
				} else {
					p.buf.WriteRune('"')
					p.buf.WriteString(pl.Captions)
					p.buf.WriteRune('"')
				}
			}
			if pl.Subtitles != "" {
				p.buf.WriteString(",SUBTITLES=\"")
				p.buf.WriteString(pl.Subtitles)
				p.buf.WriteRune('"')
			}
			if pl.Name != "" {
				p.buf.WriteString(",NAME=\"")
				p.buf.WriteString(pl.Name)
				p.buf.WriteRune('"')
			}
			p.buf.WriteRune('\n')
			p.buf.WriteString(pl.URI)
			if p.Args != "" {
				if strings.Contains(pl.URI, "?") {
					p.buf.WriteRune('&')
				} else {
					p.buf.WriteRune('?')
				}
				p.buf.WriteString(p.Args)
			}
			p.buf.WriteRune('\n')
		}
	}

	return &p.buf
}

// Version returns the current playlist version number
func (p *MasterPlaylist) Version() uint8 {
	return p.ver
}

// SetVersion sets the playlist version number, note the version maybe changed
// automatically by other Set methods.
func (p *MasterPlaylist) SetVersion(ver uint8) {
	p.ver = ver
}

// For compatibility with Stringer interface
// For example fmt.Printf("%s", sampleMediaList) will encode
// playist and print its string representation.
func (p *MasterPlaylist) String() string {
	return p.Encode().String()
}

// Creates new media playlist structure.
// Winsize defines how much items will displayed on playlist generation.
// Capacity is total size of a playlist.
func NewMediaPlaylist(winsize uint, capacity uint) (*MediaPlaylist, error) {
	p := new(MediaPlaylist)
	p.ver = minver
	p.capacity = capacity
	if err := p.SetWinSize(winsize); err != nil {
		return nil, err
	}
	p.Segments = make([]*MediaSegment, capacity)
	return p, nil
}

// last returns the previously written segment's index
func (p *MediaPlaylist) last() uint {
	if p.tail == 0 {
		return p.capacity - 1
	}
	return p.tail - 1
}

// Remove current segment from the head of chunk slice form a media playlist. Useful for sliding playlists.
// This operation does reset playlist cache.
func (p *MediaPlaylist) Remove() (err error) {
	if p.count == 0 {
		return errors.New("playlist is empty")
	}
	p.head = (p.head + 1) % p.capacity
	p.count--
	if !p.Closed {
		p.SeqNo++
	}
	p.buf.Reset()
	return nil
}

// Append general chunk to the tail of chunk slice for a media playlist.
// This operation does reset playlist cache.
func (p *MediaPlaylist) Append(uri string, duration float64, title string) error {
	seg := new(MediaSegment)
	seg.URI = uri
	seg.Duration = duration
	seg.Title = title
	return p.AppendSegment(seg)
}

// AppendSegment appends a MediaSegment to the tail of chunk slice for a media playlist.
// This operation does reset playlist cache.
func (p *MediaPlaylist) AppendSegment(seg *MediaSegment) error {
	if p.head == p.tail && p.count > 0 {
		return ErrPlaylistFull
	}
	p.Segments[p.tail] = seg
	p.tail = (p.tail + 1) % p.capacity
	p.count++
	if p.TargetDuration < seg.Duration {
		p.TargetDuration = math.Ceil(seg.Duration)
	}
	p.buf.Reset()
	return nil
}

// Combines two operations: firstly it removes one chunk from the head of chunk slice and move pointer to
// next chunk. Secondly it appends one chunk to the tail of chunk slice. Useful for sliding playlists.
// This operation does reset cache.
func (p *MediaPlaylist) Slide(uri string, duration float64, title string) {
	if !p.Closed && p.count >= p.winsize {
		p.Remove()
	} else if !p.Closed {
		p.SeqNo++
	}
	p.Append(uri, duration, title)
}

// Reset playlist cache. Next called Encode() will regenerate playlist from the chunk slice.
func (p *MediaPlaylist) ResetCache() {
	p.buf.Reset()
}

// Generate output in M3U8 format. Marshal `winsize` elements from bottom of the `segments` queue.
func (p *MediaPlaylist) Encode() *bytes.Buffer {
	if p.buf.Len() > 0 {
		return &p.buf
	}

	p.buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
	p.buf.WriteString(strver(p.ver))
	p.buf.WriteRune('\n')
	// default key (workaround for Widevine)
	if p.Key != nil {
		p.buf.WriteString("#EXT-X-KEY:")
		p.buf.WriteString("METHOD=")
		p.buf.WriteString(p.Key.Method)
		p.buf.WriteString(",URI=\"")
		p.buf.WriteString(p.Key.URI)
		p.buf.WriteRune('"')
		if p.Key.IV != "" {
			p.buf.WriteString(",IV=")
			p.buf.WriteString(p.Key.IV)
		}
		if p.Key.Keyformat != "" {
			p.buf.WriteString(",KEYFORMAT=\"")
			p.buf.WriteString(p.Key.Keyformat)
			p.buf.WriteRune('"')
		}
		if p.Key.Keyformatversions != "" {
			p.buf.WriteString(",KEYFORMATVERSIONS=\"")
			p.buf.WriteString(p.Key.Keyformatversions)
			p.buf.WriteRune('"')
		}
		p.buf.WriteRune('\n')
	}
	if p.Map != nil {
		p.buf.WriteString("#EXT-X-MAP:")
		p.buf.WriteString("URI=\"")
		p.buf.WriteString(p.Map.URI)
		p.buf.WriteRune('"')
		if p.Map.Limit > 0 {
			p.buf.WriteString(",BYTERANGE=")
			p.buf.WriteString(strconv.FormatInt(p.Map.Limit, 10))
			p.buf.WriteRune('@')
			p.buf.WriteString(strconv.FormatInt(p.Map.Offset, 10))
		}
		p.buf.WriteRune('\n')
	}
	if p.MediaType > 0 {
		p.buf.WriteString("#EXT-X-PLAYLIST-TYPE:")
		switch p.MediaType {
		case EVENT:
			p.buf.WriteString("EVENT\n")
			p.buf.WriteString("#EXT-X-ALLOW-CACHE:NO\n")
		case VOD:
			p.buf.WriteString("VOD\n")
		}
	}
	p.buf.WriteString("#EXT-X-MEDIA-SEQUENCE:")
	p.buf.WriteString(strconv.FormatUint(p.SeqNo, 10))
	p.buf.WriteRune('\n')
	p.buf.WriteString("#EXT-X-TARGETDURATION:")
	p.buf.WriteString(strconv.FormatInt(int64(math.Ceil(p.TargetDuration)), 10)) // due section 3.4.2 of M3U8 specs EXT-X-TARGETDURATION must be integer
	p.buf.WriteRune('\n')
	if p.Iframe {
		p.buf.WriteString("#EXT-X-I-FRAMES-ONLY\n")
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

	var (
		seg           *MediaSegment
		durationCache = make(map[float64]string)
	)

	head := p.head
	count := p.count
	for i := uint(0); (i < p.winsize || p.winsize == 0) && count > 0; count-- {
		seg = p.Segments[head]
		head = (head + 1) % p.capacity
		if seg == nil { // protection from badly filled chunklists
			continue
		}
		if p.winsize > 0 { // skip for VOD playlists, where winsize = 0
			i++
		}
		if seg.SCTE != nil {
			switch seg.SCTE.Syntax {
			case SCTE35_67_2014:
				p.buf.WriteString("#EXT-SCTE35:")
				p.buf.WriteString("CUE=\"")
				p.buf.WriteString(seg.SCTE.Cue)
				p.buf.WriteRune('"')
				if seg.SCTE.ID != "" {
					p.buf.WriteString(",ID=\"")
					p.buf.WriteString(seg.SCTE.ID)
					p.buf.WriteRune('"')
				}
				if seg.SCTE.Time != 0 {
					p.buf.WriteString(",TIME=")
					p.buf.WriteString(strconv.FormatFloat(seg.SCTE.Time, 'f', -1, 64))
				}
				p.buf.WriteRune('\n')
			case SCTE35_OATCLS:
				switch seg.SCTE.CueType {
				case SCTE35Cue_Start:
					p.buf.WriteString("#EXT-OATCLS-SCTE35:")
					p.buf.WriteString(seg.SCTE.Cue)
					p.buf.WriteRune('\n')
					p.buf.WriteString("#EXT-X-CUE-OUT:")
					p.buf.WriteString(strconv.FormatFloat(seg.SCTE.Time, 'f', -1, 64))
					p.buf.WriteRune('\n')
				case SCTE35Cue_Mid:
					p.buf.WriteString("#EXT-X-CUE-OUT-CONT:")
					p.buf.WriteString("ElapsedTime=")
					p.buf.WriteString(strconv.FormatFloat(seg.SCTE.Elapsed, 'f', -1, 64))
					p.buf.WriteString(",Duration=")
					p.buf.WriteString(strconv.FormatFloat(seg.SCTE.Time, 'f', -1, 64))
					p.buf.WriteString(",SCTE35=")
					p.buf.WriteString(seg.SCTE.Cue)
					p.buf.WriteRune('\n')
				case SCTE35Cue_End:
					p.buf.WriteString("#EXT-X-CUE-IN")
					p.buf.WriteRune('\n')
				}
			}
		}
		// check for key change
		if seg.Key != nil && p.Key != seg.Key {
			p.buf.WriteString("#EXT-X-KEY:")
			p.buf.WriteString("METHOD=")
			p.buf.WriteString(seg.Key.Method)
			p.buf.WriteString(",URI=\"")
			p.buf.WriteString(seg.Key.URI)
			p.buf.WriteRune('"')
			if seg.Key.IV != "" {
				p.buf.WriteString(",IV=")
				p.buf.WriteString(seg.Key.IV)
			}
			if seg.Key.Keyformat != "" {
				p.buf.WriteString(",KEYFORMAT=\"")
				p.buf.WriteString(seg.Key.Keyformat)
				p.buf.WriteRune('"')
			}
			if seg.Key.Keyformatversions != "" {
				p.buf.WriteString(",KEYFORMATVERSIONS=\"")
				p.buf.WriteString(seg.Key.Keyformatversions)
				p.buf.WriteRune('"')
			}
			p.buf.WriteRune('\n')
		}
		if seg.Discontinuity {
			p.buf.WriteString("#EXT-X-DISCONTINUITY\n")
		}
		if !seg.ProgramDateTime.IsZero() {
			p.buf.WriteString("#EXT-X-PROGRAM-DATE-TIME:")
			p.buf.WriteString(seg.ProgramDateTime.Format(DATETIME))
			p.buf.WriteRune('\n')
		}
		if seg.Limit > 0 {
			p.buf.WriteString("#EXT-X-BYTERANGE:")
			p.buf.WriteString(strconv.FormatInt(seg.Limit, 10))
			p.buf.WriteRune('@')
			p.buf.WriteString(strconv.FormatInt(seg.Offset, 10))
			p.buf.WriteRune('\n')
		}
		p.buf.WriteString("#EXTINF:")
		if str, ok := durationCache[seg.Duration]; ok {
			p.buf.WriteString(str)
		} else {
			if p.durationAsInt {
				// Old Android players has problems with non integer Duration.
				durationCache[seg.Duration] = strconv.FormatInt(int64(math.Ceil(seg.Duration)), 10)
			} else {
				// Wowza Mediaserver and some others prefer floats.
				durationCache[seg.Duration] = strconv.FormatFloat(seg.Duration, 'f', 3, 32)
			}
			p.buf.WriteString(durationCache[seg.Duration])
		}
		p.buf.WriteRune(',')
		p.buf.WriteString(seg.Title)
		p.buf.WriteRune('\n')
		p.buf.WriteString(seg.URI)
		if p.Args != "" {
			p.buf.WriteRune('?')
			p.buf.WriteString(p.Args)
		}
		p.buf.WriteRune('\n')
	}
	if p.Closed {
		p.buf.WriteString("#EXT-X-ENDLIST\n")
	}
	return &p.buf
}

// For compatibility with Stringer interface
// For example fmt.Printf("%s", sampleMediaList) will encode
// playist and print its string representation.
func (p *MediaPlaylist) String() string {
	return p.Encode().String()
}

// TargetDuration will be int on Encode
func (p *MediaPlaylist) DurationAsInt(yes bool) {
	if yes {
		// duration must be integers if protocol version is less than 3
		version(&p.ver, 3)
	}
	p.durationAsInt = yes
}

// Count tells us the number of items that are currently in the media playlist
func (p *MediaPlaylist) Count() uint {
	return p.count
}

// Close sliding playlist and make them fixed.
func (p *MediaPlaylist) Close() {
	if p.buf.Len() > 0 {
		p.buf.WriteString("#EXT-X-ENDLIST\n")
	}
	p.Closed = true
}

// Set encryption key appeared once in header of the playlist (pointer to MediaPlaylist.Key).
// It useful when keys not changed during playback.
// Set tag for the whole list.
func (p *MediaPlaylist) SetDefaultKey(method, uri, iv, keyformat, keyformatversions string) error {
	// A Media Playlist MUST indicate a EXT-X-VERSION of 5 or higher if it
	// contains:
	//   - The KEYFORMAT and KEYFORMATVERSIONS attributes of the EXT-X-KEY tag.
	if keyformat != "" || keyformatversions != "" {
		version(&p.ver, 5)
	}
	p.Key = &Key{method, uri, iv, keyformat, keyformatversions}

	return nil
}

// Set map appeared once in header of the playlist (pointer to MediaPlaylist.Key).
// It useful when map not changed during playback.
// Set tag for the whole list.
func (p *MediaPlaylist) SetDefaultMap(uri string, limit, offset int64) {
	version(&p.ver, 5) // due section 4
	p.Map = &Map{uri, limit, offset}
}

// Mark medialist as consists of only I-frames (Intra frames).
// Set tag for the whole list.
func (p *MediaPlaylist) SetIframeOnly() {
	version(&p.ver, 4) // due section 4.3.3
	p.Iframe = true
}

// Set encryption key for the current segment of media playlist (pointer to Segment.Key)
func (p *MediaPlaylist) SetKey(method, uri, iv, keyformat, keyformatversions string) error {
	if p.count == 0 {
		return errors.New("playlist is empty")
	}

	// A Media Playlist MUST indicate a EXT-X-VERSION of 5 or higher if it
	// contains:
	//   - The KEYFORMAT and KEYFORMATVERSIONS attributes of the EXT-X-KEY tag.
	if keyformat != "" || keyformatversions != "" {
		version(&p.ver, 5)
	}

	p.Segments[p.last()].Key = &Key{method, uri, iv, keyformat, keyformatversions}
	return nil
}

// Set encryption key for the current segment of media playlist (pointer to Segment.Key)
func (p *MediaPlaylist) SetMap(uri string, limit, offset int64) error {
	if p.count == 0 {
		return errors.New("playlist is empty")
	}
	version(&p.ver, 5) // due section 4
	p.Segments[p.last()].Map = &Map{uri, limit, offset}
	return nil
}

// Set limit and offset for the current media segment (EXT-X-BYTERANGE support for protocol version 4).
func (p *MediaPlaylist) SetRange(limit, offset int64) error {
	if p.count == 0 {
		return errors.New("playlist is empty")
	}
	version(&p.ver, 4) // due section 3.4.1
	p.Segments[p.last()].Limit = limit
	p.Segments[p.last()].Offset = offset
	return nil
}

// SetSCTE sets the SCTE cue format for the current media segment.
//
// Deprecated: Use SetSCTE35 instead.
func (p *MediaPlaylist) SetSCTE(cue string, id string, time float64) error {
	return p.SetSCTE35(&SCTE{Syntax: SCTE35_67_2014, Cue: cue, ID: id, Time: time})
}

// SetSCTE35 sets the SCTE cue format for the current media segment
func (p *MediaPlaylist) SetSCTE35(scte35 *SCTE) error {
	if p.count == 0 {
		return errors.New("playlist is empty")
	}
	p.Segments[p.last()].SCTE = scte35
	return nil
}

// Set discontinuity flag for the current media segment.
// EXT-X-DISCONTINUITY indicates an encoding discontinuity between the media segment
// that follows it and the one that preceded it (i.e. file format, number and type of tracks,
// encoding parameters, encoding sequence, timestamp sequence).
func (p *MediaPlaylist) SetDiscontinuity() error {
	if p.count == 0 {
		return errors.New("playlist is empty")
	}
	p.Segments[p.last()].Discontinuity = true
	return nil
}

// Set program date and time for the current media segment.
// EXT-X-PROGRAM-DATE-TIME tag associates the first sample of a
// media segment with an absolute date and/or time.  It applies only
// to the current media segment.
// Date/time format is YYYY-MM-DDThh:mm:ssZ (ISO8601) and includes time zone.
func (p *MediaPlaylist) SetProgramDateTime(value time.Time) error {
	if p.count == 0 {
		return errors.New("playlist is empty")
	}
	p.Segments[p.last()].ProgramDateTime = value
	return nil
}

// Version returns the current playlist version number
func (p *MediaPlaylist) Version() uint8 {
	return p.ver
}

// SetVersion sets the playlist version number, note the version maybe changed
// automatically by other Set methods.
func (p *MediaPlaylist) SetVersion(ver uint8) {
	p.ver = ver
}

// WinSize returns the playlist's window size.
func (p *MediaPlaylist) WinSize() uint {
	return p.winsize
}

// SetWinSize overwrites the playlist's window size.
func (p *MediaPlaylist) SetWinSize(winsize uint) error {
	if winsize > p.capacity {
		return errors.New("capacity must be greater than winsize or equal")
	}
	p.winsize = winsize
	return nil
}
