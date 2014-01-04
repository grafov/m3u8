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

// Read and parse master playlist from buffer.
// Call with strict=true will stop parsing on first format error.
func (p *MasterPlaylist) Decode(data bytes.Buffer, strict bool) error {
	return p.decode(&data, strict)
}

// Read and parse master playlist from Reader.
func (p *MasterPlaylist) DecodeFrom(reader io.Reader, strict bool) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return err
	}
	return p.decode(buf, strict)
}

func (p *MasterPlaylist) decode(buf *bytes.Buffer, strict bool) error {
	var eof, m3u, tagInf bool
	var variant *Variant
	var alt *Alternative
	var alternatives []*Alternative

	for !eof {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			eof = true
		} else if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		// start tag first
		if strings.HasPrefix(line, "#EXTM3U") {
			m3u = true
		}
		// version tag
		if strings.HasPrefix(line, "#EXT-X-VERSION:") {
			_, err = fmt.Sscanf(line, "#EXT-X-VERSION:%d", &p.ver)
			if strict && err != nil {
				return err
			}
		}
		if strings.HasPrefix(line, "#EXT-X-MEDIA:") {
			alt = new(Alternative)
			alternatives = append(alternatives, alt)
			for _, param := range strings.Split(line[13:], ",") {
				if strings.HasPrefix(param, "TYPE") {
					_, err = fmt.Sscanf(param, "TYPE=%s", &alt.Type)
					if strict && err != nil {
						return err
					}
					alt.Type = strings.Trim(alt.Type, "\"")
				}
				if strings.HasPrefix(param, "GROUP-ID") {
					_, err = fmt.Sscanf(param, "GROUP-ID=%s", &alt.GroupId)
					if strict && err != nil {
						return err
					}
					alt.GroupId = strings.Trim(alt.GroupId, "\"")
				}
				if strings.HasPrefix(param, "LANGUAGE") {
					_, err = fmt.Sscanf(param, "LANGUAGE=%s", &alt.Language)
					if strict && err != nil {
						return err
					}
					alt.Language = strings.Trim(alt.Language, "\"")
				}
				if strings.HasPrefix(param, "NAME") {
					_, err = fmt.Sscanf(param, "NAME=%s", &alt.Name)
					if strict && err != nil {
						return err
					}
					alt.Name = strings.Trim(alt.Name, "\"")
				}
				if strings.HasPrefix(param, "DEFAULT") {
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
				}
				if strings.HasPrefix(param, "AUTOSELECT") {
					_, err = fmt.Sscanf(param, "AUTOSELECT=%s", &alt.Autoselect)
					if strict && err != nil {
						return err
					}
					alt.Autoselect = strings.Trim(alt.Autoselect, "\"")
				}
				if strings.HasPrefix(param, "FORCED") {
					_, err = fmt.Sscanf(param, "FORCED=%s", &alt.Forced)
					if strict && err != nil {
						return err
					}
					alt.Forced = strings.Trim(alt.Forced, "\"")
				}
				if strings.HasPrefix(param, "CHARACTERISTICS") {
					_, err = fmt.Sscanf(param, "CHARACTERISTICS=%s", &alt.Characteristics)
					if strict && err != nil {
						return err
					}
					alt.Characteristics = strings.Trim(alt.Characteristics, "\"")
				}
				if strings.HasPrefix(param, "SUBTITLES") {
					_, err = fmt.Sscanf(param, "SUBTITLES=%s", &alt.Subtitles)
					if strict && err != nil {
						return err
					}
					alt.Subtitles = strings.Trim(alt.Subtitles, "\"")
				}

				if strings.HasPrefix(param, "URI") {
					_, err = fmt.Sscanf(param, "URI=%s", &alt.URI)
					if strict && err != nil {
						return err
					}
					alt.URI = strings.Trim(alt.URI, "\"")
				}
			}
			continue
		}
		if !tagInf && strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
			tagInf = true
			variant = new(Variant)
			if len(alternatives) > 0 {
				variant.Alternatives = alternatives
				alternatives = nil
			}
			p.Variants = append(p.Variants, variant)
			for _, param := range strings.Split(line[18:], ",") {
				if strings.HasPrefix(param, "PROGRAM-ID") {
					_, err = fmt.Sscanf(param, "PROGRAM-ID=%d", &variant.ProgramId)
					if strict && err != nil {
						return err
					}
				}
				if strings.HasPrefix(param, "BANDWIDTH") {
					_, err = fmt.Sscanf(param, "BANDWIDTH=%d", &variant.Bandwidth)
					if strict && err != nil {
						return err
					}
				}
				if strings.HasPrefix(param, "CODECS") {
					_, err = fmt.Sscanf(param, "CODECS=%s", &variant.Codecs)
					if strict && err != nil {
						return err
					}
					variant.Codecs = strings.Trim(variant.Codecs, "\"")
				}
				if strings.HasPrefix(param, "RESOLUTION") {
					_, err = fmt.Sscanf(param, "RESOLUTION=%s", &variant.Resolution)
					if strict && err != nil {
						return err
					}
					variant.Resolution = strings.Trim(variant.Resolution, "\"")
				}
				if strings.HasPrefix(param, "AUDIO") {
					_, err = fmt.Sscanf(param, "AUDIO=%s", &variant.Audio)
					if strict && err != nil {
						return err
					}
					variant.Audio = strings.Trim(variant.Audio, "\"")
				}
				if strings.HasPrefix(param, "VIDEO") {
					_, err = fmt.Sscanf(param, "VIDEO=%s", &variant.Video)
					if strict && err != nil {
						return err
					}
					variant.Video = strings.Trim(variant.Video, "\"")
				}
				if strings.HasPrefix(param, "SUBTITLES") {
					_, err = fmt.Sscanf(param, "SUBTITLES=%s", &variant.Subtitles)
					if strict && err != nil {
						return err
					}
					variant.Subtitles = strings.Trim(variant.Subtitles, "\"")
				}
			}
			continue
		}
		if tagInf {
			tagInf = false
			variant.URI = line
		}
	}
	if strict && !m3u {
		return errors.New("#EXT3MU absent")
	}
	return nil
}

func (p *MediaPlaylist) Decode(data bytes.Buffer, strict bool) error {
	return p.decode(&data, strict)
}

func (p *MediaPlaylist) DecodeFrom(reader io.Reader, strict bool) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return err
	}
	return p.decode(buf, strict)
}

func (p *MediaPlaylist) decode(buf *bytes.Buffer, strict bool) error {
	var eof, m3u, tagWV, tagInf, tagRange, tagDiscontinuity, tagProgramDateTime, tagKey bool
	var title, line string
	var duration float64
	var programDateTime time.Time
	var limit, offset int64
	var key *Key
	var err error

	wv := new(WV)
	for !eof {
		if line, err = buf.ReadString('\n'); err == io.EOF {
			eof = true
		} else if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		// start tag first
		if line == "#EXTM3U" {
			m3u = true
		}
		if line == "#EXT-X-ENDLIST" {
			p.Closed = true
		}
		if strings.HasPrefix(line, "#EXT-X-VERSION:") {
			if _, err = fmt.Sscanf(line, "#EXT-X-VERSION:%d", &p.ver); strict && err != nil {
				return err
			}
		}
		if strings.HasPrefix(line, "#EXT-X-TARGETDURATION:") {
			if _, err = fmt.Sscanf(line, "#EXT-X-TARGETDURATION:%f", &p.TargetDuration); strict && err != nil {
				return err
			}
		}
		if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:") {
			if _, err = fmt.Sscanf(line, "#EXT-X-MEDIA-SEQUENCE:%d", &p.SeqNo); strict && err != nil {
				return err
			}
		}
		if strings.HasPrefix(line, "#EXT-X-KEY:") {
			key = new(Key)
			for _, param := range strings.Split(line[11:], ",") {
				if strings.HasPrefix(param, "METHOD=") {
					if _, err = fmt.Sscanf(param, "METHOD=%s", &key.Method); strict && err != nil {
						return err
					}
				}
				if strings.HasPrefix(param, "URI=") {
					if _, err = fmt.Sscanf(param, "URI=%s", &key.URI); strict && err != nil {
						return err
					}
				}
				if strings.HasPrefix(param, "IV=") {
					if _, err = fmt.Sscanf(param, "IV=%s", &key.IV); strict && err != nil {
						return err
					}
				}
				if strings.HasPrefix(param, "KEYFORMAT=") {
					if _, err = fmt.Sscanf(param, "KEYFORMAT=%s", &key.Keyformat); strict && err != nil {
						return err
					}
				}
				if strings.HasPrefix(param, "KEYFORMATVERSIONS=") {
					if _, err = fmt.Sscanf(param, "KEYFORMATVERSIONS=%s", &key.Keyformatversions); strict && err != nil {
						return err
					}
				}
			}
			tagKey = true
		}

		if !tagProgramDateTime && strings.HasPrefix(line, "#EXT-X-PROGRAM-DATE-TIME:") {
			tagProgramDateTime = true
			if programDateTime, err = time.Parse(DATETIME, line[25:]); strict && err != nil {
				return err
			}
			continue
		}

		if !tagRange && strings.HasPrefix(line, "#EXT-X-BYTERANGE:") {
			tagRange = true
			params := strings.SplitN(line[17:], "@", 2)
			if limit, err = strconv.ParseInt(params[0], 10, 64); strict && err != nil {
				return errors.New(fmt.Sprintf("Byterange sub-range length value parsing error: %s", err))
			}
			if len(params) > 1 {
				if offset, err = strconv.ParseInt(params[1], 10, 64); strict && err != nil {
					return errors.New(fmt.Sprintf("Byterange sub-range offset value parsing error: %s", err))
				}
			}
			continue
		}

		if !tagInf && strings.HasPrefix(line, "#EXTINF:") {
			tagInf = true
			params := strings.SplitN(line[8:], ",", 2)
			if duration, err = strconv.ParseFloat(params[0], 64); strict && err != nil {
				return errors.New(fmt.Sprintf("Duration parsing error: %s", err))
			}
			title = params[1]
			continue
		}

		if !tagDiscontinuity && strings.HasPrefix(line, "#EXT-X-DISCONTINUITY") {
			tagDiscontinuity = true
			continue
		}

		if !strings.HasPrefix(line, "#") {
			if tagInf {
				p.Append(line, duration, title)
				tagInf = false
			} else if tagRange {
				if err = p.SetRange(limit, offset); strict && err != nil {
					return err
				}
				tagRange = false
			} else if tagDiscontinuity {
				tagDiscontinuity = false
				if err = p.SetDiscontinuity(); strict && err != nil {
					return err
				}
			} else if tagProgramDateTime {
				tagProgramDateTime = false
				if err = p.SetProgramDateTime(programDateTime); strict && err != nil {
					return err
				}
			}
			// if EXT-X-KEY appeared before reference to segment (EXTINF) then it linked to this segment
			if tagKey {
				p.SetKey(key.Method, key.URI, key.IV, key.Keyformat, key.Keyformatversions)
				tagKey = false
			}
		}
		// if EXT-X-KEY appeared before references to  it linked to whole playlist object
		if tagKey {
			tagKey = false
			p.Key = key
		}
		// There are a lot of Widevine tags follow:
		if strings.HasPrefix(line, "#WV-AUDIO-CHANNELS") {
			if _, err = fmt.Sscanf(line, "#WV-AUDIO-CHANNELS %d", &wv.AudioChannels); strict && err != nil {
				return err
			}
			if err == nil {
				tagWV = true
			}
		}
		if strings.HasPrefix(line, "#WV-AUDIO-FORMAT") {
			if _, err = fmt.Sscanf(line, "#WV-AUDIO-FORMAT %d", &wv.AudioFormat); strict && err != nil {
				return err
			}
			if err == nil {
				tagWV = true
			}
		}
		if strings.HasPrefix(line, "#WV-AUDIO-PROFILE-IDC") {
			if _, err = fmt.Sscanf(line, "#WV-AUDIO-PROFILE-IDC %d", &wv.AudioProfileIDC); strict && err != nil {
				return err
			}
			if err == nil {
				tagWV = true
			}
		}
		if strings.HasPrefix(line, "#WV-AUDIO-SAMPLE-SIZE") {
			if _, err = fmt.Sscanf(line, "#WV-AUDIO-SAMPLE-SIZE %d", &wv.AudioSampleSize); strict && err != nil {
				return err
			}
			if err == nil {
				tagWV = true
			}
		}
		if strings.HasPrefix(line, "#WV-AUDIO-SAMPLING-FREQUENCY") {
			if _, err = fmt.Sscanf(line, "#WV-AUDIO-SAMPLING-FREQUENCY %d", &wv.AudioSamplingFrequency); strict && err != nil {
				return err
			}
			if err == nil {
				tagWV = true
			}
		}
		if strings.HasPrefix(line, "#WV-CYPHER-VERSION") {
			wv.CypherVersion = line[19:]
			tagWV = true
		}
		if strings.HasPrefix(line, "#WV-ECM") {
			if _, err = fmt.Sscanf(line, "#WV-ECM %s", &wv.ECM); strict && err != nil {
				return err
			}
			if err == nil {
				tagWV = true
			}
		}
		if strings.HasPrefix(line, "#WV-VIDEO-FORMAT") {
			if _, err = fmt.Sscanf(line, "#WV-VIDEO-FORMAT %d", &wv.VideoFormat); strict && err != nil {
				return err
			}
			if err == nil {
				tagWV = true
			}
		}
		if strings.HasPrefix(line, "#WV-VIDEO-FRAME-RATE") {
			if _, err = fmt.Sscanf(line, "#WV-VIDEO-FRAME-RATE %d", &wv.VideoFrameRate); strict && err != nil {
				return err
			}
			if err == nil {
				tagWV = true
			}
		}
		if strings.HasPrefix(line, "#WV-VIDEO-LEVEL-IDC") {
			if _, err = fmt.Sscanf(line, "#WV-VIDEO-LEVEL-IDC %d", &wv.VideoLevelIDC); strict && err != nil {
				return err
			}
			if err == nil {
				tagWV = true
			}
		}
		if strings.HasPrefix(line, "#WV-VIDEO-PROFILE-IDC") {
			if _, err = fmt.Sscanf(line, "#WV-VIDEO-PROFILE-IDC %d", &wv.VideoProfileIDC); strict && err != nil {
				return err
			}
			if err == nil {
				tagWV = true
			}
		}
		if strings.HasPrefix(line, "#WV-VIDEO-RESOLUTION") {
			wv.VideoResolution = line[21:]
			tagWV = true
		}
		if strings.HasPrefix(line, "#WV-VIDEO-SAR") {
			if _, err = fmt.Sscanf(line, "#WV-VIDEO-SAR %s", &wv.VideoSAR); strict && err != nil {
				return err
			}
			if err == nil {
				tagWV = true
			}
		}
	}
	if tagWV {
		p.WV = wv
	}
	if strict && !m3u {
		return errors.New("#EXT3MU absent")
	}
	return nil
}

// Tries to detect playlist type and returns playlist structure of appropriate type.
func Decode(data bytes.Buffer, strict bool) (Playlist, ListType, error) {
	return decode(&data, strict)
}

func DecodeFrom(reader io.Reader, strict bool) (Playlist, ListType, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil, UNKNOWN, err
	}
	return decode(buf, strict)
}

func decode(buf *bytes.Buffer, strict bool) (Playlist, ListType, error) {
	var eof, m3u, mediaExtinf, masterStreamInf bool
	var variant *Variant
	var title string
	var duration float64
	var ver uint8
	var listType ListType

	master := NewMasterPlaylist()
	media, err := NewMediaPlaylist(8, 1024) // TODO find better way instead of hardcoded values
	if err != nil {
		return nil, UNKNOWN, errors.New(fmt.Sprintf("Create media playlist failed: %s", err))
	}

	for !eof {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			eof = true
		} else if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		// start tag first
		if strings.HasPrefix(line, "#EXTM3U") {
			m3u = true
		}
		// version tag
		if strings.HasPrefix(line, "#EXT-X-VERSION:") {
			_, err = fmt.Sscanf(line, "#EXT-X-VERSION:%d", &ver)
			if strict && err != nil {
				return nil, listType, err
			}
		}

		// Master playlist parsing
		if listType != MEDIA && !masterStreamInf && strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
			listType = MASTER
			masterStreamInf = true
			variant = new(Variant)
			master.Variants = append(master.Variants, variant)
			for _, param := range strings.Split(line[18:], ",") {
				if strings.HasPrefix(param, "PROGRAM-ID") {
					_, err = fmt.Sscanf(param, "PROGRAM-ID=%d", &variant.ProgramId)
					if strict && err != nil {
						return nil, MASTER, err
					}
				}
				if strings.HasPrefix(param, "BANDWIDTH") {
					_, err = fmt.Sscanf(param, "BANDWIDTH=%d", &variant.Bandwidth)
					if strict && err != nil {
						return nil, MASTER, err
					}
				}
				if strings.HasPrefix(param, "CODECS") {
					_, err = fmt.Sscanf(param, "CODECS=%s", &variant.Codecs)
					if strict && err != nil {
						return nil, MASTER, err
					}
				}
				if strings.HasPrefix(param, "RESOLUTION") {
					_, err = fmt.Sscanf(param, "RESOLUTION=%s", &variant.Resolution)
					if strict && err != nil {
						return nil, MASTER, err
					}
				}
				if strings.HasPrefix(param, "AUDIO") {
					_, err = fmt.Sscanf(param, "AUDIO=%s", &variant.Audio)
					if strict && err != nil {
						return nil, MASTER, err
					}
				}
				if strings.HasPrefix(param, "VIDEO") {
					_, err = fmt.Sscanf(param, "VIDEO=%s", &variant.Video)
					if strict && err != nil {
						return nil, MASTER, err
					}
				}
				if strings.HasPrefix(param, "SUBTITLES") {
					_, err = fmt.Sscanf(param, "SUBTITLES=%s", &variant.Subtitles)
					if strict && err != nil {
						return nil, MASTER, err
					}
				}
			}
			continue
		}
		if listType == MEDIA && masterStreamInf {
			masterStreamInf = false
			variant.URI = line
		}

		if listType != MASTER && line == "#EXT-X-ENDLIST" {
			listType = MEDIA
			media.Closed = true
		}
		if listType != MASTER && strings.HasPrefix(line, "#EXT-X-TARGETDURATION:") {
			listType = MEDIA
			_, err = fmt.Sscanf(line, "#EXT-X-TARGETDURATION:%f", &media.TargetDuration)
			if strict && err != nil {
				return nil, MEDIA, err
			}
		}
		if listType != MASTER && strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:") {
			listType = MEDIA
			_, err = fmt.Sscanf(line, "#EXT-X-MEDIA-SEQUENCE:%d", &media.SeqNo)
			if strict && err != nil {
				return nil, MEDIA, err
			}
		}
		if listType != MASTER && !mediaExtinf && strings.HasPrefix(line, "#EXTINF:") {
			listType = MEDIA
			mediaExtinf = true
			params := strings.SplitN(line[8:], ",", 2)
			duration, err = strconv.ParseFloat(params[0], 64)
			if strict && err != nil {
				return nil, MEDIA, errors.New(fmt.Sprintf("Media playlist duration parsing error: %s", err))
			}
			title = params[1]
			continue
		}
		if listType == MEDIA && mediaExtinf {
			mediaExtinf = false
			media.Append(line, duration, title)
		}
	}

	if strict && !m3u {
		return nil, listType, errors.New("#EXT3MU absent")
	}

	switch listType {
	case MASTER:
		master.ver = ver
		return master, MASTER, nil
	case MEDIA:
		media.ver = ver
		return media, MEDIA, nil
	}
	return nil, listType, errors.New("Can't detect playlist type.")
}
