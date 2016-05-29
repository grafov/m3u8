package m3u8

/*
 Part of M3U8 parser & generator library.
 This file defines functions related to playlist parsing.

 Copyleft 2013-2015 Alexander I.Grafov aka Axel <grafov@gmail.com>
 Copyleft 2013-2015 library authors (see AUTHORS file).

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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Parse master playlist from the buffer.
// If `strict` parameter is true then return first syntax error.
func (p *MasterPlaylist) Decode(buf bytes.Buffer, strict bool) error {
	var eof bool

	state := new(decodingState)

	for !eof {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			eof = true
		} else if err != nil {
			break
		}
		err = decodeLineOfMasterPlaylist(p, state, line, strict)
		if strict && err != nil {
			return err
		}
	}
	if strict && !state.m3u {
		return errors.New("#EXTM3U absent")
	}

	p.buf = *bytes.NewBuffer(make([]byte, 0, buf.Len()*2))
	return nil
}

// Parse master playlist from the io.Reader stream.
// If `strict` parameter is true then return first syntax error.
func (p *MasterPlaylist) DecodeFrom(reader io.Reader, strict bool) error {
	var eof bool

	state := new(decodingState)
	size := 0

	scanner := bufio.NewScanner(reader)

	for !eof {
		eof = !scanner.Scan()
		line := scanner.Text()
		size += len(line)

		if err := decodeLineOfMasterPlaylist(p, state, line, strict); strict && err != nil {
			return err
		}
	}
	if strict && !state.m3u {
		return errors.New("#EXT3MU absent")
	}

	p.buf = *bytes.NewBuffer(make([]byte, 0, size*2))
	return nil
}

// Parse media playlist from the buffer.
// If `strict` parameter is true then return first syntax error.
func (p *MediaPlaylist) Decode(buf bytes.Buffer, strict bool) error {
	var eof bool
	var line string
	var err error

	state := new(decodingState)
	wv := new(WV)

	for !eof {
		if line, err = buf.ReadString('\n'); err == io.EOF {
			eof = true
		} else if err != nil {
			break
		}

		if err = decodeLineOfMediaPlaylist(p, wv, state, line, strict); strict && err != nil {
			return err
		}

	}
	if state.tagWV {
		p.WV = wv
	}
	if strict && !state.m3u {
		return errors.New("#EXT3MU absent")
	}

	p.buf = *bytes.NewBuffer(make([]byte, 0, buf.Len()*2))
	return nil
}

// Parse media playlist from the io.Reader stream.
// If `strict` parameter is true then return first syntax error.
func (p *MediaPlaylist) DecodeFrom(reader io.Reader, strict bool) error {
	var eof bool
	var err error

	state := new(decodingState)
	wv := new(WV)
	size := 0

	scanner := bufio.NewScanner(reader)

	for !eof {
		eof = !scanner.Scan()
		line := scanner.Text()
		size += len(line)

		err = decodeLineOfMediaPlaylist(p, wv, state, line, strict)
		if strict && err != nil {
			return err
		}

	}
	if state.tagWV {
		p.WV = wv
	}
	if strict && !state.m3u {
		return errors.New("#EXTM3U absent")
	}

	p.buf = *bytes.NewBuffer(make([]byte, 0, size*2))
	return nil
}

// Detect playlist type and decode it from the buffer.
func Decode(data bytes.Buffer, strict bool) (Playlist, ListType, error) {
	return decode(&data, strict)
}

// Detect playlist type and decode it from input stream.
func DecodeFrom(reader io.Reader, strict bool) (Playlist, ListType, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil, 0, err
	}
	return decode(buf, strict)
}

// Detect playlist type and decode it. May be used as decoder for both master and media playlists.
func decode(buf *bytes.Buffer, strict bool) (Playlist, ListType, error) {
	var eof bool
	var line string
	var master *MasterPlaylist
	var media *MediaPlaylist
	var listType ListType
	var err error

	state := new(decodingState)
	wv := new(WV)

	master = NewMasterPlaylist()
	media, err = NewMediaPlaylist(8, 1024) // TODO make it autoextendable
	if err != nil {
		return nil, 0, fmt.Errorf("Create media playlist failed: %s", err)
	}

	for !eof {
		if line, err = buf.ReadString('\n'); err == io.EOF {
			eof = true
		} else if err != nil {
			break
		}

		// fixes the issues https://github.com/grafov/m3u8/issues/25
		// TODO: the same should be done in decode functions of both Master- and MediaPlaylists
		// so some DRYing would be needed.
		if len(line) < 1 || line == "\r" {
			continue
		}

		err = decodeLineOfMasterPlaylist(master, state, line, strict)
		if strict && err != nil {
			return master, state.listType, err
		}

		err = decodeLineOfMediaPlaylist(media, wv, state, line, strict)
		if strict && err != nil {
			return media, state.listType, err
		}

	}
	if state.listType == MEDIA && state.tagWV {
		media.WV = wv
	}

	if strict && !state.m3u {
		return nil, listType, errors.New("#EXTM3U absent")
	}

	switch state.listType {
	case MASTER:
		return master, MASTER, nil
	case MEDIA:
		return media, MEDIA, nil
	default:
		return nil, state.listType, errors.New("Can't detect playlist type")
	}
	return nil, state.listType, errors.New("This return is impossible. Saved for compatibility with go 1.0")
}

func decodeParamsLine(line string) map[string]string {
	out := make(map[string]string, 8)
	for {
		equalIndex := strings.Index(line, "=")
		if equalIndex == -1 {
			// error
			return out
		}
		key := strings.TrimLeftFunc(line[:equalIndex], unicode.IsSpace)
		if equalIndex == len(line)-1 {
			// error
			return out
		}
		line = line[equalIndex+1:]
		if line[0] == '"' {
			if len(line) < 3 {
				// error
				return out
			}
			line = line[1:]
			quoteIndex := strings.Index(line, `"`)
			if quoteIndex == -1 {
				// error
				return out
			}
			out[key] = line[:quoteIndex]
			if quoteIndex > len(line)-3 {
				return out
			}
			line = line[quoteIndex+2:]
		} else {
			commaIndex := strings.Index(line, ",")
			if commaIndex == -1 {
				commaIndex = len(line)
			}
			out[key] = line[:commaIndex]
			if commaIndex > len(line)-2 {
				return out
			}
			line = line[commaIndex+1:]
		}
	}
	return out
}

// Parse one line of master playlist.
func decodeLineOfMasterPlaylist(p *MasterPlaylist, state *decodingState, line string, strict bool) error {
	var err error

	line = strings.TrimSpace(line)

	switch {
	case line == "#EXTM3U": // start tag first
		state.m3u = true
	case strings.HasPrefix(line, "#EXT-X-VERSION:"): // version tag
		state.listType = MASTER
		num, err := strconv.ParseUint(line[len("#EXT-X-VERSION:"):], 10, 8)
		if strict && err != nil {
			return err
		}
		p.ver = uint8(num)
	case strings.HasPrefix(line, "#EXT-X-MEDIA:"):
		var alt Alternative
		state.listType = MASTER
		for k, v := range decodeParamsLine(line[13:]) {
			switch k {
			case "TYPE":
				alt.Type = v
			case "GROUP-ID":
				alt.GroupId = v
			case "LANGUAGE":
				alt.Language = v
			case "NAME":
				alt.Name = v
			case "DEFAULT":
				if strings.ToUpper(v) == "YES" {
					alt.Default = true
				} else if strings.ToUpper(v) == "NO" {
					alt.Default = false
				} else if strict {
					return errors.New("value must be YES or NO")
				}
			case "AUTOSELECT":
				alt.Autoselect = v
			case "FORCED":
				alt.Forced = v
			case "CHARACTERISTICS":
				alt.Characteristics = v
			case "SUBTITLES":
				alt.Subtitles = v
			case "URI":
				alt.URI = v
			}
		}
		state.alternatives = append(state.alternatives, &alt)
	case !state.tagStreamInf && strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
		state.tagStreamInf = true
		state.listType = MASTER
		state.variant = new(Variant)
		if len(state.alternatives) > 0 {
			state.variant.Alternatives = state.alternatives
			state.alternatives = nil
		}
		p.Variants = append(p.Variants, state.variant)
		for k, v := range decodeParamsLine(line[18:]) {
			switch k {
			case "PROGRAM-ID":
				var val int
				val, err = strconv.Atoi(v)
				if strict && err != nil {
					return err
				}
				state.variant.ProgramId = uint32(val)
			case "BANDWIDTH":
				var val int
				val, err = strconv.Atoi(v)
				if strict && err != nil {
					return err
				}
				state.variant.Bandwidth = uint32(val)
			case "CODECS":
				state.variant.Codecs = v
			case "RESOLUTION":
				state.variant.Resolution = v
			case "AUDIO":
				state.variant.Audio = v
			case "VIDEO":
				state.variant.Video = v
			case "SUBTITLES":
				state.variant.Subtitles = v
			case "CLOSED-CAPTIONS":
				state.variant.Captions = v
			case "NAME":
				state.variant.Name = v
			}
		}
	case state.tagStreamInf && !strings.HasPrefix(line, "#"):
		state.tagStreamInf = false
		state.variant.URI = line
	case strings.HasPrefix(line, "#EXT-X-I-FRAME-STREAM-INF:"):
		state.listType = MASTER
		state.variant = new(Variant)
		state.variant.Iframe = true
		if len(state.alternatives) > 0 {
			state.variant.Alternatives = state.alternatives
			state.alternatives = nil
		}
		p.Variants = append(p.Variants, state.variant)
		for k, v := range decodeParamsLine(line[26:]) {
			switch k {
			case "URI":
				state.variant.URI = v
			case "PROGRAM-ID":
				var val int
				val, err = strconv.Atoi(v)
				if strict && err != nil {
					return err
				}
				state.variant.ProgramId = uint32(val)
			case "BANDWIDTH":
				var val int
				val, err = strconv.Atoi(v)
				if strict && err != nil {
					return err
				}
				state.variant.Bandwidth = uint32(val)
			case "CODECS":
				state.variant.Codecs = v
			case "RESOLUTION":
				state.variant.Resolution = v
			case "AUDIO":
				state.variant.Audio = v
			case "VIDEO":
				state.variant.Video = v
			}
		}
	case strings.HasPrefix(line, "#"): // unknown tags treated as comments
		return err
	}
	return err
}

// Parse one line of media playlist.
func decodeLineOfMediaPlaylist(p *MediaPlaylist, wv *WV, state *decodingState, line string, strict bool) error {
	var title string
	var err error

	line = strings.TrimRightFunc(line, unicode.IsSpace)
	switch {
	case !state.tagInf && strings.HasPrefix(line, "#EXTINF:"):
		state.tagInf = true
		state.listType = MEDIA
		sepIndex := strings.Index(line, ",")
		if sepIndex == -1 {
			sepIndex = len(line)
		}
		state.durationString = line[8:sepIndex]
		if len(line) > sepIndex {
			title = line[sepIndex+1:]
		}
	case !strings.HasPrefix(line, "#"):
		if state.tagInf {
			p.AppendWithDurationString(line, state.durationString, title)
			state.tagInf = false
		}
		if state.tagRange {
			if err = p.SetRange(state.limit, state.offset); strict && err != nil {
				return err
			}
			state.tagRange = false
		}
		if state.tagSCTE35 {
			state.tagSCTE35 = false
			scte := *state.scte
			if err = p.SetSCTE(scte.Cue, scte.ID, scte.Time); strict && err != nil {
				return err
			}
		}
		if state.tagDiscontinuity {
			state.tagDiscontinuity = false
			if err = p.SetDiscontinuity(); strict && err != nil {
				return err
			}
		}
		if state.tagProgramDateTime {
			state.tagProgramDateTime = false
			if err = p.SetProgramDateTime(state.programDateTime); strict && err != nil {
				return err
			}
		}
		// If EXT-X-KEY appeared before reference to segment (EXTINF) then it linked to this segment
		if state.tagKey {
			p.Segments[p.last()].Key = &Key{state.xkey.Method, state.xkey.URI, state.xkey.IV, state.xkey.Keyformat, state.xkey.Keyformatversions}
			// First EXT-X-KEY may appeared in the header of the playlist and linked to first segment
			// but for convenient playlist generation it also linked as default playlist key
			if p.Key == nil {
				p.Key = state.xkey
			}
			state.tagKey = false
		}
		// If EXT-X-MAP appeared before reference to segment (EXTINF) then it linked to this segment
		if state.tagMap {
			p.Segments[p.last()].Map = &Map{state.xmap.URI, state.xmap.Limit, state.xmap.Offset}
			// First EXT-X-MAP may appeared in the header of the playlist and linked to first segment
			// but for convenient playlist generation it also linked as default playlist map
			if p.Map == nil {
				p.Map = state.xmap
			}
			state.tagMap = false
		}
	// start tag first
	case line == "#EXTM3U":
		state.m3u = true
	case line == "#EXT-X-ENDLIST":
		state.listType = MEDIA
		p.Closed = true
	case strings.HasPrefix(line, "#EXT-X-VERSION:"):
		state.listType = MEDIA
		num, err := strconv.ParseUint(line[len("#EXT-X-VERSION:"):], 10, 8)
		if strict && err != nil {
			return err
		}
		p.ver = uint8(num)
	case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
		state.listType = MEDIA
		num, err := strconv.ParseFloat(line[len("#EXT-X-TARGETDURATION:"):], 64)
		if strict && err != nil {
			return err
		}
		p.TargetDuration = num
	case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
		state.listType = MEDIA
		num, err := strconv.ParseUint(line[len("#EXT-X-MEDIA-SEQUENCE:"):], 10, 64)
		if strict && err != nil {
			return err
		}
		p.SeqNo = num
	case strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE:"):
		state.listType = MEDIA
		playlistType := line[len("#EXT-X-PLAYLIST-TYPE:"):]
		switch playlistType {
		case "EVENT":
			p.MediaType = EVENT
		case "VOD":
			p.MediaType = VOD
		}
	case strings.HasPrefix(line, "#EXT-X-KEY:"):
		state.listType = MEDIA
		state.xkey = new(Key)
		for k, v := range decodeParamsLine(line[len("#EXT-X-KEY:"):]) {
			switch k {
			case "METHOD":
				state.xkey.Method = v
			case "URI":
				state.xkey.URI = v
			case "IV":
				state.xkey.IV = v
			case "KEYFORMAT":
				state.xkey.Keyformat = v
			case "KEYFORMATVERSIONS":
				state.xkey.Keyformatversions = v
			}
		}
		state.tagKey = true
	case strings.HasPrefix(line, "#EXT-X-MAP:"):
		state.listType = MEDIA
		state.xmap = new(Map)
		for k, v := range decodeParamsLine(line[len("#EXT-X-MAP:"):]) {
			switch k {
			case "URI":
				state.xmap.URI = v
			case "BYTERANGE":
				if _, err = fmt.Sscanf(v, "%d@%d", &state.xmap.Limit, &state.xmap.Offset); strict && err != nil {
					return fmt.Errorf("Byterange sub-range length value parsing error: %s", err)
				}
			}
		}
		state.tagMap = true
	case !state.tagProgramDateTime && strings.HasPrefix(line, "#EXT-X-PROGRAM-DATE-TIME:"):
		state.tagProgramDateTime = true
		state.listType = MEDIA
		if state.programDateTime, err = time.Parse(DATETIME, line[len("#EXT-X-PROGRAM-DATE-TIME:"):]); strict && err != nil {
			return err
		}
	case !state.tagRange && strings.HasPrefix(line, "#EXT-X-BYTERANGE:"):
		state.tagRange = true
		state.listType = MEDIA
		state.offset = 0
		params := strings.SplitN(line[len("#EXT-X-BYTERANGE:"):], "@", 2)
		if state.limit, err = strconv.ParseInt(params[0], 10, 64); strict && err != nil {
			return fmt.Errorf("Byterange sub-range length value parsing error: %s", err)
		}
		if len(params) > 1 {
			if state.offset, err = strconv.ParseInt(params[1], 10, 64); strict && err != nil {
				return fmt.Errorf("Byterange sub-range offset value parsing error: %s", err)
			}
		}
	case !state.tagSCTE35 && strings.HasPrefix(line, "#EXT-SCTE35:"):
		state.tagSCTE35 = true
		state.listType = MEDIA
		state.scte = new(SCTE)
		for attribute, value := range decodeParamsLine(line[len("#EXT-SCTE35:"):]) {
			switch attribute {
			case "CUE":
				state.scte.Cue = value
			case "ID":
				state.scte.ID = value
			case "TIME":
				state.scte.Time, _ = strconv.ParseFloat(value, 64)
			}
		}
	case !state.tagDiscontinuity && strings.HasPrefix(line, "#EXT-X-DISCONTINUITY"):
		state.tagDiscontinuity = true
		state.listType = MEDIA
	case strings.HasPrefix(line, "#EXT-X-I-FRAMES-ONLY"):
		state.listType = MEDIA
		p.Iframe = true
	case strings.HasPrefix(line, "#WV-AUDIO-CHANNELS"):
		state.listType = MEDIA
		num, err := strconv.ParseUint(line[len("#WV-AUDIO-CHANNELS "):], 10, 0)
		if err == nil {
			wv.AudioChannels = uint(num)
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-FORMAT"):
		state.listType = MEDIA
		num, err := strconv.ParseUint(line[len("#WV-AUDIO-FORMAT "):], 10, 0)
		if err == nil {
			wv.AudioFormat = uint(num)
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-PROFILE-IDC"):
		state.listType = MEDIA
		num, err := strconv.ParseUint(line[len("#WV-AUDIO-PROFILE-IDC "):], 10, 0)
		if err == nil {
			wv.AudioProfileIDC = uint(num)
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-SAMPLE-SIZE"):
		state.listType = MEDIA
		num, err := strconv.ParseUint(line[len("#WV-AUDIO-SAMPLE-SIZE "):], 10, 0)
		if strict && err != nil {
			return err
		}
		if err == nil {
			wv.AudioSampleSize = uint(num)
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-SAMPLING-FREQUENCY"):
		state.listType = MEDIA
		num, err := strconv.ParseUint(line[len("#WV-AUDIO-SAMPLING-FREQUENCY "):], 10, 0)
		if strict && err != nil {
			return err
		}
		if err == nil {
			wv.AudioSamplingFrequency = uint(num)
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-CYPHER-VERSION"):
		state.listType = MEDIA
		wv.CypherVersion = line[len("#WV-CYPHER-VERSION "):]
		state.tagWV = true
	case strings.HasPrefix(line, "#WV-ECM"):
		state.listType = MEDIA
		wv.CypherVersion = line[len("#WV-ECM "):]
		state.tagWV = true
	case strings.HasPrefix(line, "#WV-VIDEO-FORMAT"):
		state.listType = MEDIA
		num, err := strconv.ParseUint(line[len("#WV-VIDEO-FORMAT "):], 10, 0)
		if strict && err != nil {
			return err
		}
		if err == nil {
			wv.VideoFormat = uint(num)
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-FRAME-RATE"):
		state.listType = MEDIA
		num, err := strconv.ParseUint(line[len("#WV-VIDEO-FRAME-RATE "):], 10, 0)
		if strict && err != nil {
			return err
		}
		if err == nil {
			wv.VideoFrameRate = uint(num)
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-LEVEL-IDC"):
		state.listType = MEDIA
		num, err := strconv.ParseUint(line[len("#WV-VIDEO-LEVEL-IDC "):], 10, 0)
		if strict && err != nil {
			return err
		}
		if err == nil {
			wv.VideoLevelIDC = uint(num)
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-PROFILE-IDC"):
		state.listType = MEDIA
		num, err := strconv.ParseUint(line[len("#WV-VIDEO-PROFILE-IDC "):], 10, 0)
		if strict && err != nil {
			return err
		}
		if err == nil {
			wv.VideoProfileIDC = uint(num)
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-RESOLUTION"):
		state.listType = MEDIA
		wv.VideoResolution = line[len("#WV-VIDEO-RESOLUTION "):]
		state.tagWV = true
	case strings.HasPrefix(line, "#WV-VIDEO-SAR"):
		state.listType = MEDIA
		wv.VideoResolution = line[len("#WV-VIDEO-SAR "):]
		state.tagWV = true
	case strings.HasPrefix(line, "#"): // unknown tags treated as comments
		return err
	}
	return err
}
