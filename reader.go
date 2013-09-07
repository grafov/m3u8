package m3u8

/*
 Part of M3U8 parser & generator library.
 This file defines functions related to playlist parsing.

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
	"io"
	"strconv"
	"strings"
)

// Read and parse master playlist.
// Call with strict=true will stop parsing on first format error.
func (p *MasterPlaylist) Decode(data bytes.Buffer, strict bool) error {
	return p.decode(&data, strict)
}

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
		if !tagInf && strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
			tagInf = true
			variant = new(Variant)
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
				}
				if strings.HasPrefix(param, "RESOLUTION") {
					_, err = fmt.Sscanf(param, "RESOLUTION=%s", &variant.Resolution)
					if strict && err != nil {
						return err
					}
				}
				if strings.HasPrefix(param, "AUDIO") {
					_, err = fmt.Sscanf(param, "AUDIO=%s", &variant.Audio)
					if strict && err != nil {
						return err
					}
				}
				if strings.HasPrefix(param, "VIDEO") {
					_, err = fmt.Sscanf(param, "VIDEO=%s", &variant.Video)
					if strict && err != nil {
						return err
					}
				}
				if strings.HasPrefix(param, "SUBTITLES") {
					_, err = fmt.Sscanf(param, "SUBTITLES=%s", &variant.Subtitles)
					if strict && err != nil {
						return err
					}
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
	var eof, m3u, tagInf bool
	var title string
	var duration float64

	for !eof {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
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
			_, err = fmt.Sscanf(line, "#EXT-X-VERSION:%d", &p.ver)
			if strict && err != nil {
				return err
			}
		}
		if strings.HasPrefix(line, "#EXT-X-TARGETDURATION:") {
			_, err = fmt.Sscanf(line, "#EXT-X-TARGETDURATION:%f", &p.TargetDuration)
			if strict && err != nil {
				return err
			}
		}
		if strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE:") {
			_, err = fmt.Sscanf(line, "#EXT-X-MEDIA-SEQUENCE:%d", &p.SeqNo)
			if strict && err != nil {
				return err
			}
		}
		if !tagInf && strings.HasPrefix(line, "#EXTINF:") {
			tagInf = true
			params := strings.SplitN(line[8:], ",", 2)
			duration, err = strconv.ParseFloat(params[0], 64)
			if strict && err != nil {
				return errors.New(fmt.Sprintf("Duration parsing error: %s", err))
			}
			title = params[1]
			continue
		}
		if tagInf {
			tagInf = false
			p.Add(line, duration, title)
		}
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
			media.Add(line, duration, title)
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
