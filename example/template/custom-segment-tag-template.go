package template

import (
	"bytes"
	"errors"

	"github.com/grafov/m3u8"
)

// #CUSTOM-SEGMENT-TAG:<attribute-list>

// CustomSegmentTag implements both CustomTag and CustomDecoder
// interfaces.
type CustomSegmentTag struct {
	Name string
	Jedi bool
}

// TagName should return the full indentifier including the leading '#' and trailing ':'
// if the tag also contains a value or attribute list
func (tag *CustomSegmentTag) TagName() string {
	return "#CUSTOM-SEGMENT-TAG:"
}

// Decode decodes the input string to the internal structure. The line
// will be the entire matched line, including the identifier.
func (tag *CustomSegmentTag) Decode(line string) (m3u8.CustomTag, error) {
	var err error

	// Since this is a Segment tag, we want to create a new tag every time it is decoded
	// as there can be one for each segment with
	newTag := new(CustomSegmentTag)

	for k, v := range m3u8.DecodeAttributeList(line[20:]) {
		switch k {
		case "NAME":
			newTag.Name = v
		case "JEDI":
			if v == "YES" {
				newTag.Jedi = true
			} else if v == "NO" {
				newTag.Jedi = false
			} else {
				err = errors.New("Valid strings for JEDI attribute are YES and NO.")
			}
		}
	}

	return newTag, err
}

// SegmentTag is a playlist tag example.
func (tag *CustomSegmentTag) SegmentTag() bool {
	return true
}

// Encode encodes the structure to the text result.
func (tag *CustomSegmentTag) Encode() *bytes.Buffer {
	buf := new(bytes.Buffer)

	if tag.Name != "" {
		buf.WriteString(tag.TagName())
		buf.WriteString("NAME=\"")
		buf.WriteString(tag.Name)
		buf.WriteString("\",JEDI=")
		if tag.Jedi {
			buf.WriteString("YES")
		} else {
			buf.WriteString("NO")
		}
	}

	return buf
}

// String implements Stringer interface.
func (tag *CustomSegmentTag) String() string {
	return tag.Encode().String()
}
