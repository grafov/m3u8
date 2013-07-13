package m3u8

/*
 Part of M3U8 parser & generator library.

 Copyleft Alexander I.Grafov aka Axel <grafov@gmail.com>
 Library licensed under GPLv3

 ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
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
			p.variants = append(p.variants, variant)
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

func (p *MediaPlaylist) Decode(reader io.Reader) error {
	var eof, started, tagInf bool
	var seqid uint64

	buf := bufio.NewReader(reader)

	for !eof {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			eof = true
		} else if err != nil {
			break
		}
		// start tag first
		if line == "#EXTM3U" {
			started = true
		}
		// version tag
		if started && strings.HasPrefix(line, "#EXT-X-VERSION:") {
			_, err = fmt.Sscanf(line, "#EXT-X-VERSION:%d", p.ver)
			return err
		}
		if started && strings.HasPrefix(line, "#EXT-X-STREAM-INF:") {
			tagInf = true
			_, err = fmt.Sscanf(line, "", p.ver)

			seqid++
		}
		if started && tagInf {

		}
	}
	return nil
}
