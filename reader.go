package m3u8

/*
 Part of M3U8 parser & generator library.
 This file defines functions related to playlist parsing.

 Copyleft 2013-2014 Alexander I.Grafov aka Axel <grafov@gmail.com>

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
	"io"
	"strconv"
	"strings"
	"time"
)

// Parse master playlist from the buffer.
// If `strict` parameter is true then return first syntax error.
func (p *MasterPlaylist) Decode(data bytes.Buffer, strict bool) error {
	return p.decode(&data, strict)
}

// Parse master playlist from the io.Reader stream.
// If `strict` parameter is true then return first syntax error.
func (p *MasterPlaylist) DecodeFrom(reader io.Reader, strict bool) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return err
	}
	return p.decode(buf, strict)
}

// Parse master playlist. Internal function.
func (p *MasterPlaylist) decode(buf *bytes.Buffer, strict bool) error {
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
		return errors.New("#EXT3MU absent")
	}
	return nil
}

// Parse media playlist from the buffer.
// If `strict` parameter is true then return first syntax error.
func (p *MediaPlaylist) Decode(data bytes.Buffer, strict bool) error {
	return p.decode(&data, strict)
}

// Parse media playlist from the io.Reader stream.
// If `strict` parameter is true then return first syntax error.
func (p *MediaPlaylist) DecodeFrom(reader io.Reader, strict bool) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return err
	}
	return p.decode(buf, strict)
}

func (p *MediaPlaylist) decode(buf *bytes.Buffer, strict bool) error {
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

		err = decodeLineOfMediaPlaylist(p, wv, state, line, strict)
		if strict && err != nil {
			return err
		}

	}
	if state.tagWV {
		p.WV = wv
	}
	if strict && !state.m3u {
		return errors.New("#EXT3MU absent")
	}
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
		return nil, listType, errors.New("#EXT3MU absent")
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

// Parse one line of master playlist.
func decodeLineOfMasterPlaylist(p *MasterPlaylist, state *decodingState, line string, strict bool) error {
	var alt *Alternative
	var alternatives []*Alternative
	var err error

	line = strings.TrimSpace(line)

	switch {
	case line == "#EXTM3U": // start tag first
		state.m3u = true
	case strings.HasPrefix(line, "#EXT-X-VERSION:"): // version tag
		state.listType = MASTER
		_, err = fmt.Sscanf(line, "#EXT-X-VERSION:%d", &p.ver)
		if strict && err != nil {
			return err
		}
	case strings.HasPrefix(line, "#EXT-X-MEDIA:"):
		state.listType = MASTER
		alt = new(Alternative)
		alternatives = append(alternatives, alt)
		for _, param := range strings.Split(line[13:], ",") {
			switch {
			case strings.HasPrefix(param, "TYPE"):
				_, err = fmt.Sscanf(param, "TYPE=%s", &alt.Type)
				if strict && err != nil {
					return err
				}
				alt.Type = strings.Trim(alt.Type, "\"")
			case strings.HasPrefix(param, "GROUP-ID"):
				_, err = fmt.Sscanf(param, "GROUP-ID=%s", &alt.GroupId)
				if strict && err != nil {
					return err
				}
				alt.GroupId = strings.Trim(alt.GroupId, "\"")
			case strings.HasPrefix(param, "LANGUAGE"):
				_, err = fmt.Sscanf(param, "LANGUAGE=%s", &alt.Language)
				if strict && err != nil {
					return err
				}
				alt.Language = strings.Trim(alt.Language, "\"")
			case strings.HasPrefix(param, "NAME"):
				_, err = fmt.Sscanf(param, "NAME=%s", &alt.Name)
				if strict && err != nil {
					return err
				}
				alt.Name = strings.Trim(alt.Name, "\"")
			case strings.HasPrefix(param, "DEFAULT"):
				var val string
				_, err = fmt.Sscanf(param, "DEFAULT=%s", &val)
				if strict && err != nil {
					return err
				}
				val = strings.Trim(val, "\"")
				if strings.ToUpper(val) == "YES" {
					alt.Default = true
				} else if strings.ToUpper(val) == "NO" {
					alt.Default = false
				} else if strict {
					return errors.New("value must be YES or NO")
				}
			case strings.HasPrefix(param, "AUTOSELECT"):
				_, err = fmt.Sscanf(param, "AUTOSELECT=%s", &alt.Autoselect)
				if strict && err != nil {
					return err
				}
				alt.Autoselect = strings.Trim(alt.Autoselect, "\"")
			case strings.HasPrefix(param, "FORCED"):
				_, err = fmt.Sscanf(param, "FORCED=%s", &alt.Forced)
				if strict && err != nil {
					return err
				}
				alt.Forced = strings.Trim(alt.Forced, "\"")
			case strings.HasPrefix(param, "CHARACTERISTICS"):
				_, err = fmt.Sscanf(param, "CHARACTERISTICS=%s", &alt.Characteristics)
				if strict && err != nil {
					return err
				}
				alt.Characteristics = strings.Trim(alt.Characteristics, "\"")
			case strings.HasPrefix(param, "SUBTITLES"):
				_, err = fmt.Sscanf(param, "SUBTITLES=%s", &alt.Subtitles)
				if strict && err != nil {
					return err
				}
				alt.Subtitles = strings.Trim(alt.Subtitles, "\"")
			case strings.HasPrefix(param, "URI"):
				_, err = fmt.Sscanf(param, "URI=%s", &alt.URI)
				if strict && err != nil {
					return err
				}
				alt.URI = strings.Trim(alt.URI, "\"")
			}
		}
	case !state.tagStreamInf && strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
		state.tagStreamInf = true
		state.listType = MASTER
		state.variant = new(Variant)
		if len(alternatives) > 0 {
			state.variant.Alternatives = alternatives
			alternatives = nil
		}
		p.Variants = append(p.Variants, state.variant)
		for _, param := range strings.Split(line[18:], ",") {
			switch {
			case strings.HasPrefix(param, "PROGRAM-ID"):
				_, err = fmt.Sscanf(param, "PROGRAM-ID=%d", &state.variant.ProgramId)
				if strict && err != nil {
					return err
				}
			case strings.HasPrefix(param, "BANDWIDTH"):
				_, err = fmt.Sscanf(param, "BANDWIDTH=%d", &state.variant.Bandwidth)
				if strict && err != nil {
					return err
				}
			case strings.HasPrefix(param, "CODECS"):
				_, err = fmt.Sscanf(param, "CODECS=%s", &state.variant.Codecs)
				if strict && err != nil {
					return err
				}
				state.variant.Codecs = strings.Trim(state.variant.Codecs, "\"")
			case strings.HasPrefix(param, "RESOLUTION"):
				_, err = fmt.Sscanf(param, "RESOLUTION=%s", &state.variant.Resolution)
				if strict && err != nil {
					return err
				}
				state.variant.Resolution = strings.Trim(state.variant.Resolution, "\"")
			case strings.HasPrefix(param, "AUDIO"):
				_, err = fmt.Sscanf(param, "AUDIO=%s", &state.variant.Audio)
				if strict && err != nil {
					return err
				}
				state.variant.Audio = strings.Trim(state.variant.Audio, "\"")
			case strings.HasPrefix(param, "VIDEO"):
				_, err = fmt.Sscanf(param, "VIDEO=%s", &state.variant.Video)
				if strict && err != nil {
					return err
				}
				state.variant.Video = strings.Trim(state.variant.Video, "\"")
			case strings.HasPrefix(param, "SUBTITLES"):
				_, err = fmt.Sscanf(param, "SUBTITLES=%s", &state.variant.Subtitles)
				if strict && err != nil {
					return err
				}
				state.variant.Subtitles = strings.Trim(state.variant.Subtitles, "\"")
			}
		}
	case state.tagStreamInf && !strings.HasPrefix(line, "#"):
		state.tagStreamInf = false
		state.variant.URI = line
	case strings.HasPrefix(line, "#"): // unknown tags treated as comments
		return err
	}
	return err
}

// Parse one line of media playlist.
func decodeLineOfMediaPlaylist(p *MediaPlaylist, wv *WV, state *decodingState, line string, strict bool) error {
	var title string
	var err error

	line = strings.TrimSpace(line)
	switch {
	// start tag first
	case line == "#EXTM3U":
		state.m3u = true
	case line == "#EXT-X-ENDLIST":
		state.listType = MEDIA
		p.Closed = true
	case strings.HasPrefix(line, "#EXT-X-VERSION:"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#EXT-X-VERSION:%d", &p.ver); strict && err != nil {
			return err
		}
	case strings.HasPrefix(line, "#EXT-X-TARGETDURATION:"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#EXT-X-TARGETDURATION:%f", &p.TargetDuration); strict && err != nil {
			return err
		}
	case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#EXT-X-MEDIA-SEQUENCE:%d", &p.SeqNo); strict && err != nil {
			return err
		}
	case strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE:"):
		state.listType = MEDIA
		var playlistType string
		_, err = fmt.Sscanf(line, "#EXT-X-PLAYLIST-TYPE:%s", &playlistType)
		if err != nil {
			if strict {
				return err
			}
		} else {
			switch playlistType {
			case "EVENT":
				p.MediaType = EVENT
			case "VOD":
				p.MediaType = VOD
			}
		}
	case strings.HasPrefix(line, "#EXT-X-KEY:"):
		state.listType = MEDIA
		state.key = new(Key)
		for _, param := range strings.Split(line[11:], ",") {
			if strings.HasPrefix(param, "METHOD=") {
				if _, err = fmt.Sscanf(param, "METHOD=%s", &state.key.Method); strict && err != nil {
					return err
				}
			}
			if strings.HasPrefix(param, "URI=") {
				if _, err = fmt.Sscanf(param, "URI=%s", &state.key.URI); strict && err != nil {
					return err
				}
			}
			if strings.HasPrefix(param, "IV=") {
				if _, err = fmt.Sscanf(param, "IV=%s", &state.key.IV); strict && err != nil {
					return err
				}
			}
			if strings.HasPrefix(param, "KEYFORMAT=") {
				if _, err = fmt.Sscanf(param, "KEYFORMAT=%s", &state.key.Keyformat); strict && err != nil {
					return err
				}
			}
			if strings.HasPrefix(param, "KEYFORMATVERSIONS=") {
				if _, err = fmt.Sscanf(param, "KEYFORMATVERSIONS=%s", &state.key.Keyformatversions); strict && err != nil {
					return err
				}
			}
		}
		state.tagKey = true
	case !state.tagProgramDateTime && strings.HasPrefix(line, "#EXT-X-PROGRAM-DATE-TIME:"):
		state.tagProgramDateTime = true
		state.listType = MEDIA
		if state.programDateTime, err = time.Parse(DATETIME, line[25:]); strict && err != nil {
			return err
		}
	case !state.tagRange && strings.HasPrefix(line, "#EXT-X-BYTERANGE:"):
		state.tagRange = true
		state.listType = MEDIA
		params := strings.SplitN(line[17:], "@", 2)
		if state.limit, err = strconv.ParseInt(params[0], 10, 64); strict && err != nil {
			return fmt.Errorf("Byterange sub-range length value parsing error: %s", err)
		}
		if len(params) > 1 {
			if state.offset, err = strconv.ParseInt(params[1], 10, 64); strict && err != nil {
				return fmt.Errorf("Byterange sub-range offset value parsing error: %s", err)
			}
		}
	case !state.tagInf && strings.HasPrefix(line, "#EXTINF:"):
		state.tagInf = true
		state.listType = MEDIA
		params := strings.SplitN(line[8:], ",", 2)
		if state.duration, err = strconv.ParseFloat(params[0], 64); strict && err != nil {
			return fmt.Errorf("Duration parsing error: %s", err)
		}
		title = params[1]
	case !state.tagDiscontinuity && strings.HasPrefix(line, "#EXT-X-DISCONTINUITY"):
		state.tagDiscontinuity = true
		state.listType = MEDIA
	case !strings.HasPrefix(line, "#"):
		if state.tagInf {
			p.Append(line, state.duration, title)
			state.tagInf = false
		} else if state.tagRange {
			if err = p.SetRange(state.limit, state.offset); strict && err != nil {
				return err
			}
			state.tagRange = false
		} else if state.tagDiscontinuity {
			state.tagDiscontinuity = false
			if err = p.SetDiscontinuity(); strict && err != nil {
				return err
			}
		} else if state.tagProgramDateTime {
			state.tagProgramDateTime = false
			if err = p.SetProgramDateTime(state.programDateTime); strict && err != nil {
				return err
			}
		}
		// If EXT-X-KEY appeared before reference to segment (EXTINF) then it linked to this segment
		if state.tagKey {
			p.Segments[(p.tail-1)%p.capacity].Key = &Key{state.key.Method, state.key.URI, state.key.IV, state.key.Keyformat, state.key.Keyformatversions}
			// First EXT-X-KEY may appeared in the header of the playlist and linked to first segment
			// but for convenient playlist generation it also linked as default playlist key
			if p.Key == nil {
				p.Key = state.key
			}
			state.tagKey = false
		}
	case strings.HasPrefix(line, "#WV-AUDIO-CHANNELS"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-AUDIO-CHANNELS %d", &wv.AudioChannels); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-FORMAT"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-AUDIO-FORMAT %d", &wv.AudioFormat); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-PROFILE-IDC"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-AUDIO-PROFILE-IDC %d", &wv.AudioProfileIDC); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-SAMPLE-SIZE"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-AUDIO-SAMPLE-SIZE %d", &wv.AudioSampleSize); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-AUDIO-SAMPLING-FREQUENCY"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-AUDIO-SAMPLING-FREQUENCY %d", &wv.AudioSamplingFrequency); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-CYPHER-VERSION"):
		state.listType = MEDIA
		wv.CypherVersion = line[19:]
		state.tagWV = true
	case strings.HasPrefix(line, "#WV-ECM"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-ECM %s", &wv.ECM); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-FORMAT"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-VIDEO-FORMAT %d", &wv.VideoFormat); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-FRAME-RATE"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-VIDEO-FRAME-RATE %d", &wv.VideoFrameRate); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-LEVEL-IDC"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-VIDEO-LEVEL-IDC %d", &wv.VideoLevelIDC); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-PROFILE-IDC"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-VIDEO-PROFILE-IDC %d", &wv.VideoProfileIDC); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#WV-VIDEO-RESOLUTION"):
		state.listType = MEDIA
		wv.VideoResolution = line[21:]
		state.tagWV = true
	case strings.HasPrefix(line, "#WV-VIDEO-SAR"):
		state.listType = MEDIA
		if _, err = fmt.Sscanf(line, "#WV-VIDEO-SAR %s", &wv.VideoSAR); strict && err != nil {
			return err
		}
		if err == nil {
			state.tagWV = true
		}
	case strings.HasPrefix(line, "#"): // unknown tags treated as comments
		return err
	}
	return err
}
