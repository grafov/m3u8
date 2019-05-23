package template

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/grafov/m3u8"
)

// #CUSTOM-PLAYLIST-TAG:<number>

// Implements both CustomTag and CustomDecoder interfaces
type CustomPlaylistTag struct {
	Number int
}

// TagName() should return the full indentifier including the leading '#' and trailing ':'
// if the tag also contains a value or attribute list
func (tag *CustomPlaylistTag) TagName() string {
	return "#CUSTOM-PLAYLIST-TAG:"
}

// line will be the entire matched line, including the identifier
func (tag *CustomPlaylistTag) Decode(line string) (m3u8.CustomTag, error) {
	_, err := fmt.Sscanf(line, "#CUSTOM-PLAYLIST-TAG:%d", &tag.Number)

	return tag, err
}

// This is a playlist tag example
func (tag *CustomPlaylistTag) Segment() bool {
	return false
}

func (tag *CustomPlaylistTag) Encode() *bytes.Buffer {
	buf := new(bytes.Buffer)

	buf.WriteString(tag.TagName())
	buf.WriteString(strconv.Itoa(tag.Number))

	return buf
}

func (tag *CustomPlaylistTag) String() string {
	return tag.Encode().String()
}
