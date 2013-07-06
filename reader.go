package m3u8

/*
 Part of M3U8 parser & generator library.

 Copyleft Alexander I.Grafov aka Axel <grafov@gmail.com>
 Library licensed under GPLv3

 ॐ तारे तुत्तारे तुरे स्व
*/

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func (p *MediaPlaylist) Decode(reader io.Reader) error {
	var eof, started, streamInf bool
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
			streamInf = true
			_, err = fmt.Sscanf(line, "", p.ver)

			seqid++
			//seg := Segment{SeqId: seqid
		}
		if started && streamInf {

		}
	}
	return nil
}

func (p *MasterPlaylist) Load(reader io.Reader) {
	// buf := bufio.NewReader(reader)

	//	line, err := buf.ReadString('\n')
}
