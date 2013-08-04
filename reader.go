package m3u8

/*
 Part of M3U8 parser & generator library.
 This file defines functions related to playlist parsing.

 Copyleft Alexander I.Grafov aka Axel <grafov@gmail.com>
 Library licensed under GPLv3

 ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ListType uint

const (
	UNKNOWN ListType = iota
	MASTER
	MEDIA
)

// Read and parse master playlist.
// Call with strict=true will stop parsing on first format error.
func (p *MasterPlaylist) Decode(reader io.Reader, strict bool) error {
	var eof, m3u, tagInf bool
	var variant *Variant

	buf := bufio.NewReader(reader)

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

func (p *MediaPlaylist) Decode(reader io.Reader, strict bool) error {
	var eof, m3u, tagInf bool
	var title string
	var duration float64

	buf := bufio.NewReader(reader)

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
func Decode(reader io.Reader, strict bool) (interface{}, ListType, error) {
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
	buf := bufio.NewReader(reader)

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
	default:
		return nil, listType, errors.New("Can't detect playlist type.")
	}
}
