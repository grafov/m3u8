package template

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/grafov/m3u8"
)

// #CUSTOM-PLAYLIST-TAG:<number>

// CustomPlaylistTag implements both CustomTag and CustomDecoder
// interfaces.
type CustomPlaylistTag struct {
	Number int
}

// TagName should return the full indentifier including the leading
// '#' and trailing ':' if the tag also contains a value or attribute
// list.
func (tag *CustomPlaylistTag) TagName() string {
	return "#CUSTOM-PLAYLIST-TAG:"
}

// Decode decodes the input line. The line will be the entire matched
// line, including the identifier
func (tag *CustomPlaylistTag) Decode(line string) (m3u8.CustomTag, error) {
	_, err := fmt.Sscanf(line, "#CUSTOM-PLAYLIST-TAG:%d", &tag.Number)

	return tag, err
}

// SegmentTag is a playlist tag example.
func (tag *CustomPlaylistTag) SegmentTag() bool {
	return false
}

// Encode formats the structure to the text result.
func (tag *CustomPlaylistTag) Encode() *bytes.Buffer {
	buf := new(bytes.Buffer)

	buf.WriteString(tag.TagName())
	buf.WriteString(strconv.Itoa(tag.Number))

	return buf
}

// String implements Stringer interface.
func (tag *CustomPlaylistTag) String() string {
	return tag.Encode().String()
}
