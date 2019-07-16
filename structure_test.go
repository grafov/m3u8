/*
 Playlist structures tests.

 Copyright 2013-2017 The Project Developers.
 See the AUTHORS and LICENSE files at the top-level directory of this distribution
 and at https://github.com/grafov/m3u8/

 ॐ तारे तुत्तारे तुरे स्व
*/
package m3u8

import (
	"bytes"
	"testing"
)

func CheckType(t *testing.T, p Playlist) {
	t.Logf("%T implements Playlist interface OK\n", p)
}

// Create new media playlist.
func TestNewMediaPlaylist(t *testing.T) {
	_, e := NewMediaPlaylist(1, 2)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
}

type MockCustomTag struct {
	name          string
	err           error
	segment       bool
	encodedString string
}

func (t *MockCustomTag) TagName() string {
	return t.name
}

func (t *MockCustomTag) Decode(line string) (CustomTag, error) {
	return t, t.err
}

func (t *MockCustomTag) Encode() *bytes.Buffer {
	if t.encodedString == "" {
		return nil
	}

	buf := new(bytes.Buffer)

	buf.WriteString(t.encodedString)

	return buf
}

func (t *MockCustomTag) String() string {
	return t.encodedString
}

func (t *MockCustomTag) SegmentTag() bool {
	return t.segment
}
