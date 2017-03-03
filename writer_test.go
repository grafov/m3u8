/*
 Package m3u8. Playlist generation tests.

 Copyright 2013-2016 The Project Developers.
 See the AUTHORS and LICENSE files at the top-level directory of this distribution
 and at https://github.com/grafov/m3u8/

 ॐ तारे तुत्तारे तुरे स्व
*/
package m3u8

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// Check how master and media playlists implement common Playlist interface
func TestInterfaceImplemented(t *testing.T) {
	m := NewMasterPlaylist()
	CheckType(t, m)
	p, e := NewMediaPlaylist(1, 2)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	CheckType(t, p)
}

// Create new media playlist with wrong size (must be failed)
func TestCreateMediaPlaylistWithWrongSize(t *testing.T) {
	_, e := NewMediaPlaylist(2, 1) //wrong winsize
	if e == nil {
		t.Fatal("Create new media playlist must be failed, but it's don't")
	}
}

// Tests the last method on media playlist
func TestLastSegmentMediaPlaylist(t *testing.T) {
	p, _ := NewMediaPlaylist(5, 5)
	if p.last() != 4 {
		t.Errorf("last is %v, expected: 4", p.last())
	}
	for i := uint(0); i < 5; i++ {
		_ = p.Append("uri.ts", 4, "")
		if p.last() != i {
			t.Errorf("last is: %v, expected: %v", p.last(), i)
		}
	}
}

// Create new media playlist
// Add two segments to media playlist
func TestAddSegmentToMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(1, 2)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	e = p.Append("test01.ts", 10.0, "title")
	if e != nil {
		t.Errorf("Add 1st segment to a media playlist failed: %s", e)
	}
	if p.Segments[0].URI != "test01.ts" {
		t.Errorf("Expected: test01.ts, got: %v", p.Segments[0].URI)
	}
	if p.Segments[0].Duration != 10 {
		t.Errorf("Expected: 10, got: %v", p.Segments[0].Duration)
	}
	if p.Segments[0].Title != "title" {
		t.Errorf("Expected: title, got: %v", p.Segments[0].Title)
	}
}

func TestAppendSegmentToMediaPlaylist(t *testing.T) {
	p, _ := NewMediaPlaylist(2, 2)
	e := p.AppendSegment(&MediaSegment{Duration: 10})
	if e != nil {
		t.Errorf("Add 1st segment to a media playlist failed: %s", e)
	}
	if p.TargetDuration != 10 {
		t.Errorf("Failed to increase TargetDuration, expected: 10, got: %v", p.TargetDuration)
	}
	e = p.AppendSegment(&MediaSegment{Duration: 10})
	if e != nil {
		t.Errorf("Add 2nd segment to a media playlist failed: %s", e)
	}
	e = p.AppendSegment(&MediaSegment{Duration: 10})
	if e != ErrPlaylistFull {
		t.Errorf("Add 3rd expected full error, got: %s", e)
	}
}

// Create new media playlist
// Add three segments to media playlist
// Set discontinuity tag for the 2nd segment.
func TestDiscontinuityForMediaPlaylist(t *testing.T) {
	var e error
	p, e := NewMediaPlaylist(3, 4)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	p.Close()
	if e = p.Append("test01.ts", 5.0, ""); e != nil {
		t.Errorf("Add 1st segment to a media playlist failed: %s", e)
	}
	if e = p.Append("test02.ts", 6.0, ""); e != nil {
		t.Errorf("Add 2nd segment to a media playlist failed: %s", e)
	}
	if e = p.SetDiscontinuity(); e != nil {
		t.Error("Can't set discontinuity tag")
	}
	if e = p.Append("test03.ts", 6.0, ""); e != nil {
		t.Errorf("Add 3nd segment to a media playlist failed: %s", e)
	}
	//fmt.Println(p.Encode().String())
}

// Create new media playlist
// Add three segments to media playlist
// Set program date and time for 2nd segment.
// Set discontinuity tag for the 2nd segment.
func TestProgramDateTimeForMediaPlaylist(t *testing.T) {
	var e error
	p, e := NewMediaPlaylist(3, 4)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	p.Close()
	if e = p.Append("test01.ts", 5.0, ""); e != nil {
		t.Errorf("Add 1st segment to a media playlist failed: %s", e)
	}
	if e = p.Append("test02.ts", 6.0, ""); e != nil {
		t.Errorf("Add 2nd segment to a media playlist failed: %s", e)
	}
	loc, _ := time.LoadLocation("Europe/Moscow")
	if e = p.SetProgramDateTime(time.Date(2010, time.November, 30, 16, 25, 0, 125*1e6, loc)); e != nil {
		t.Error("Can't set program date and time")
	}
	if e = p.SetDiscontinuity(); e != nil {
		t.Error("Can't set discontinuity tag")
	}
	if e = p.Append("test03.ts", 6.0, ""); e != nil {
		t.Errorf("Add 3nd segment to a media playlist failed: %s", e)
	}
	//fmt.Println(p.Encode().String())
}

// Create new media playlist
// Add two segments to media playlist with duration 9.0 and 9.1.
// Target duration must be set to nearest greater integer (= 10).
func TestTargetDurationForMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(1, 2)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	e = p.Append("test01.ts", 9.0, "")
	if e != nil {
		t.Errorf("Add 1st segment to a media playlist failed: %s", e)
	}
	e = p.Append("test02.ts", 9.1, "")
	if e != nil {
		t.Errorf("Add 2nd segment to a media playlist failed: %s", e)
	}
	if p.TargetDuration < 10.0 {
		t.Errorf("Target duration must = 10 (nearest greater integer to durations 9.0 and 9.1)")
	}
}

// Create new media playlist with capacity 10 elements
// Try to add 11 segments to media playlist (oversize error)
func TestOverAddSegmentsToMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(1, 10)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 0; i < 11; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			t.Logf("As expected new segment #%d not assigned to a media playlist: %s due oversize\n", i, e)
		}
	}
}

// Create new media playlist
// Add segment to media playlist
// Set SCTE
func TestSetSCTEForMediaPlaylist(t *testing.T) {
	tests := []struct {
		Cue      string
		ID       string
		Time     float64
		Expected string
	}{
		{"CueData1", "", 0, `#EXT-SCTE35:CUE="CueData1"` + "\n"},
		{"CueData2", "ID2", 0, `#EXT-SCTE35:CUE="CueData2",ID="ID2"` + "\n"},
		{"CueData3", "ID3", 3.141, `#EXT-SCTE35:CUE="CueData3",ID="ID3",TIME=3.141` + "\n"},
		{"CueData4", "", 3.1, `#EXT-SCTE35:CUE="CueData4",TIME=3.1` + "\n"},
		{"CueData5", "", 3.0, `#EXT-SCTE35:CUE="CueData5",TIME=3` + "\n"},
	}

	for _, test := range tests {
		p, e := NewMediaPlaylist(1, 1)
		if e != nil {
			t.Fatalf("Create media playlist failed: %s", e)
		}
		if e = p.Append("test01.ts", 5.0, ""); e != nil {
			t.Errorf("Add 1st segment to a media playlist failed: %s", e)
		}
		if e := p.SetSCTE(test.Cue, test.ID, test.Time); e != nil {
			t.Errorf("SetSCTE to a media playlist failed: %s", e)
		}
		if !strings.Contains(p.String(), test.Expected) {
			t.Errorf("Test %+v did not contain: %q, playlist: %v", test, test.Expected, p.String())
		}
	}
}

// Create new media playlist
// Add segment to media playlist
// Set encryption key
func TestSetKeyForMediaPlaylist(t *testing.T) {
	tests := []struct {
		KeyFormat         string
		KeyFormatVersions string
		ExpectVersion     uint8
	}{
		{"", "", 3},
		{"Format", "", 5},
		{"", "Version", 5},
		{"Format", "Version", 5},
	}

	for _, test := range tests {
		p, e := NewMediaPlaylist(3, 5)
		if e != nil {
			t.Fatalf("Create media playlist failed: %s", e)
		}
		if e = p.Append("test01.ts", 5.0, ""); e != nil {
			t.Errorf("Add 1st segment to a media playlist failed: %s", e)
		}
		if e := p.SetKey("AES-128", "https://example.com", "iv", test.KeyFormat, test.KeyFormatVersions); e != nil {
			t.Errorf("Set key to a media playlist failed: %s", e)
		}
		if p.ver != test.ExpectVersion {
			t.Errorf("Set key playlist version: %v, expected: %v", p.ver, test.ExpectVersion)
		}
	}
}

// Create new media playlist
// Add segment to media playlist
// Set encryption key
func TestSetDefaultKeyForMediaPlaylist(t *testing.T) {
	tests := []struct {
		KeyFormat         string
		KeyFormatVersions string
		ExpectVersion     uint8
	}{
		{"", "", 3},
		{"Format", "", 5},
		{"", "Version", 5},
		{"Format", "Version", 5},
	}

	for _, test := range tests {
		p, e := NewMediaPlaylist(3, 5)
		if e != nil {
			t.Fatalf("Create media playlist failed: %s", e)
		}
		if e := p.SetDefaultKey("AES-128", "https://example.com", "iv", test.KeyFormat, test.KeyFormatVersions); e != nil {
			t.Errorf("Set key to a media playlist failed: %s", e)
		}
		if p.ver != test.ExpectVersion {
			t.Errorf("Set key playlist version: %v, expected: %v", p.ver, test.ExpectVersion)
		}
	}
}

// Create new media playlist
// Add segment to media playlist
// Set map
func TestSetMapForMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	e = p.Append("test01.ts", 5.0, "")
	if e != nil {
		t.Errorf("Add 1st segment to a media playlist failed: %s", e)
	}
	e = p.SetMap("https://example.com", 1000*1024, 1024*1024)
	if e != nil {
		t.Errorf("Set map to a media playlist failed: %s", e)
	}
}

// Create new media playlist
// Add two segments to media playlist
// Encode structures to HLS
func TestEncodeMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	e = p.Append("test01.ts", 5.0, "")
	if e != nil {
		t.Errorf("Add 1st segment to a media playlist failed: %s", e)
	}
	p.DurationAsInt(true)
	//fmt.Println(p.Encode().String())
}

// Create new media playlist
// Add 10 segments to media playlist
// Test iterating over segments
func TestLoopSegmentsOfMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 0; i < 5; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			t.Errorf("Add segment #%d to a media playlist failed: %s", i, e)
		}
	}
	p.DurationAsInt(true)
	//fmt.Println(p.Encode().String())
}

// Create new media playlist with capacity 5
// Add 5 segments and 5 unique keys
// Test correct keys set on correct segments
func TestEncryptionKeysInMediaPlaylist(t *testing.T) {
	p, _ := NewMediaPlaylist(5, 5)
	// Add 5 segments and set custom encryption key
	for i := uint(0); i < 5; i++ {
		uri := fmt.Sprintf("uri-%d", i)
		expected := &Key{
			Method:            "AES-128",
			URI:               uri,
			IV:                fmt.Sprintf("%d", i),
			Keyformat:         "identity",
			Keyformatversions: "1",
		}
		_ = p.Append(uri+".ts", 4, "")
		_ = p.SetKey(expected.Method, expected.URI, expected.IV, expected.Keyformat, expected.Keyformatversions)

		if p.Segments[i].Key == nil {
			t.Fatalf("Key was not set on segment %v", i)
		}
		if *p.Segments[i].Key != *expected {
			t.Errorf("Key %+v does not match expected %+v", p.Segments[i].Key, expected)
		}
	}
}

// Create new media playlist
// Add 10 segments to media playlist
// Encode structure to HLS with integer target durations
func TestMediaPlaylistWithIntegerDurations(t *testing.T) {
	p, e := NewMediaPlaylist(3, 10)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 0; i < 9; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.6, "")
		if e != nil {
			t.Errorf("Add segment #%d to a media playlist failed: %s", i, e)
		}
	}
	p.DurationAsInt(false)
	//	fmt.Println(p.Encode().String())
}

// Create new media playlist
// Add 9 segments to media playlist
// 11 times encode structure to HLS with integer target durations
// Last playlist must be empty
func TestMediaPlaylistWithEmptyMedia(t *testing.T) {
	p, e := NewMediaPlaylist(3, 10)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 1; i < 10; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.6, "")
		if e != nil {
			t.Errorf("Add segment #%d to a media playlist failed: %s", i, e)
		}
	}
	for i := 1; i < 11; i++ {
		//fmt.Println(p.Encode().String())
		p.Remove()
	} // TODO add check for buffers equality
}

// Create new media playlist with winsize == capacity
func TestMediaPlaylistWinsize(t *testing.T) {
	p, e := NewMediaPlaylist(6, 6)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 1; i < 10; i++ {
		p.Slide(fmt.Sprintf("test%d.ts", i), 5.6, "")
		//fmt.Println(p.Encode().String()) // TODO check playlist sizes and mediasequence values
	}
}

// Create new media playlist as sliding playlist.
// Close it.
func TestClosedMediaPlaylist(t *testing.T) {
	p, e := NewMediaPlaylist(1, 10)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 0; i < 10; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			t.Errorf("Due oversize new segment #%d not assigned to a media playlist: %s\n", i, e)
		}
	}
	p.Close()
}

// Create new media playlist as sliding playlist.
func TestLargeMediaPlaylistWithParallel(t *testing.T) {
	testCount := 10
	expect, err := ioutil.ReadFile("sample-playlists/media-playlist-large.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < testCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			f, err := os.Open("sample-playlists/media-playlist-large.m3u8")
			if err != nil {
				t.Fatal(err)
			}
			p, err := NewMediaPlaylist(50000, 50000)
			if err != nil {
				t.Fatalf("Create media playlist failed: %s", err)
			}
			if err = p.DecodeFrom(bufio.NewReader(f), true); err != nil {
				t.Fatal(err)
			}

			actual := p.Encode().Bytes() // disregard output
			if bytes.Compare(expect, actual) != 0 {
				t.Fatal("not matched")
			}
		}()
		wg.Wait()
	}
}

func TestMediaVersion(t *testing.T) {
	m, _ := NewMediaPlaylist(3, 3)
	m.ver = 5
	if m.Version() != m.ver {
		t.Errorf("Expected version: %v, got: %v", m.ver, m.Version())
	}
}

func TestMediaSetVersion(t *testing.T) {
	m, _ := NewMediaPlaylist(3, 3)
	m.ver = 3
	m.SetVersion(5)
	if m.ver != 5 {
		t.Errorf("Expected version: %v, got: %v", 5, m.ver)
	}
}

func TestMediaWinSize(t *testing.T) {
	m, _ := NewMediaPlaylist(3, 3)
	if m.WinSize() != m.winsize {
		t.Errorf("Expected winsize: %v, got: %v", m.winsize, m.WinSize())
	}
}

func TestMediaSetWinSize(t *testing.T) {
	m, _ := NewMediaPlaylist(3, 5)
	err := m.SetWinSize(5)
	if err != nil {
		t.Fatal(err)
	}
	if m.winsize != 5 {
		t.Errorf("Expected winsize: %v, got: %v", 5, m.winsize)
	}
	// Check winsize cannot exceed capacity
	err = m.SetWinSize(99999)
	if err == nil {
		t.Error("Expected error, received: ", err)
	}
	// Ensure winsize didn't change
	if m.winsize != 5 {
		t.Errorf("Expected winsize: %v, got: %v", 5, m.winsize)
	}
}

// Create new master playlist without params
// Add media playlist
func TestNewMasterPlaylist(t *testing.T) {
	m := NewMasterPlaylist()
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 0; i < 5; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			t.Errorf("Add segment #%d to a media playlist failed: %s", i, e)
		}
	}
	m.Append("chunklist1.m3u8", p, VariantParams{})
}

// Create new master playlist without params
// Add media playlist with Alternatives
func TestNewMasterPlaylistWithAlternatives(t *testing.T) {
	m := NewMasterPlaylist()
	audioUri := fmt.Sprintf("%s/rendition.m3u8", "800")
	audioAlt := &Alternative{
		GroupId:    "audio",
		URI:        audioUri,
		Type:       "AUDIO",
		Name:       "main",
		Default:    true,
		Autoselect: "YES",
		Language:   "english",
	}
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 0; i < 5; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			t.Errorf("Add segment #%d to a media playlist failed: %s", i, e)
		}
	}
	m.Append("chunklist1.m3u8", p, VariantParams{Alternatives: []*Alternative{audioAlt}})

	if m.ver != 4 {
		t.Fatalf("Expected version 4, actual, %d", m.ver)
	}
	fmt.Printf("%v\n", m)
}

// Create new master playlist supporting closed-caption=none
func TestNewMasterPlaylistWithClosedCaptionEqNone(t *testing.T) {
	m := NewMasterPlaylist()
	langs := []string{"fra", "eng"}
	audioBitrates := []string{"64", "128"}
	videoBitrates := []uint32{1170400, 3630000, 6380000}
	videoResolutions := []string{"320x240", "854x480", "1920x1080"}
	alts := []*Alternative{}
	vps := []*VariantParams{}

	for _, lang := range langs {
		for i, bitrate := range audioBitrates {
			a := &Alternative{
				GroupId:    fmt.Sprintf("audio%d", i),
				URI:        fmt.Sprintf("audio_%s_%s_rendition.m3u8", bitrate, lang),
				Type:       "AUDIO",
				Name:       fmt.Sprintf("%s", lang),
				Default:    lang == "fra",
				Autoselect: "YES",
				Language:   lang,
			}
			alts = append(alts, a)
		}

	}
	alts = append(alts, &Alternative{Type: "SUBTITLES", GroupId: "subtitles0", Name: "eng_subtitle", Default: false, Autoselect: "YES", Language: "eng", URI: "subtitle_eng_rendition.m3u8"})

	for i, videoBitrate := range videoBitrates {
		s := func(videoBitrate uint32) string {
			HQ := 0
			if videoBitrate >= 3630000 {
				HQ = 1
			}
			return fmt.Sprintf("audio%d", HQ)
		}
		vp := &VariantParams{
			ProgramId:    0,
			Bandwidth:    videoBitrate,
			Codecs:       "avc1",
			Resolution:   videoResolutions[i],
			Audio:        s(videoBitrate),
			Subtitles:    "subtitles0",
			Captions:     "NONE",
			Alternatives: alts,
		}
		vps = append(vps, vp)
	}

	p, err := NewMediaPlaylist(1, 1)
	if err != nil {
		t.Fatalf("Create media playlist failed: %s", err)
	}
	for i, vp := range vps {
		m.Append(fmt.Sprintf("%d_rendition.m3u8", i+1), p, *vp)
	}
	if m.ver != 4 {
		t.Fatalf("Expected version 4, actual, %d", m.ver)
	}
	// diff with master file using reader.
	// REQUIRES "github.com/pmezard/go-difflib/difflib"
	//f, err := os.Open("sample-playlists/master-with-closed-captions-eq-none.m3u8")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//nmp := NewMasterPlaylist()
	//err = nmp.DecodeFrom(bufio.NewReader(f), false)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	//diff := difflib.UnifiedDiff{
	//	A:        difflib.SplitLines(fmt.Sprintf("%s", m)),
	//	B:        difflib.SplitLines(fmt.Sprintf("%s", nmp)),
	//	FromFile: "TestData",
	//	ToFile:   "FileSource",
	//	Context:  1,
	//}
	//
	//text, _ := difflib.GetUnifiedDiffString(diff)
	//if text != "" {
	//	fmt.Println(text)
	//	t.Fatalf("Data not equal to File: sample-playlists/master-with-closed-captions-eq-none.m3u8")
	//}
	//fmt.Printf("%v\n", m)
}

// Create new master playlist with params
// Add media playlist
func TestNewMasterPlaylistWithParams(t *testing.T) {
	m := NewMasterPlaylist()
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 0; i < 5; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			t.Errorf("Add segment #%d to a media playlist failed: %s", i, e)
		}
	}
	m.Append("chunklist1.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
}

// Create new master playlist
// Add media playlist with existing query params in URI
// Append more query params and ensure it encodes correctly
func TestEncodeMasterPlaylistWithExistingQuery(t *testing.T) {
	m := NewMasterPlaylist()
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 0; i < 5; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			t.Errorf("Add segment #%d to a media playlist failed: %s", i, e)
		}
	}
	m.Append("chunklist1.m3u8?k1=v1&k2=v2", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
	m.Args = "k3=v3"
	if !strings.Contains(m.String(), "chunklist1.m3u8?k1=v1&k2=v2&k3=v3\n") {
		t.Errorf("Encode master with existing args failed")
	}
}

// Create new master playlist
// Add media playlist
// Encode structures to HLS
func TestEncodeMasterPlaylist(t *testing.T) {
	m := NewMasterPlaylist()
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 0; i < 5; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			t.Errorf("Add segment #%d to a media playlist failed: %s", i, e)
		}
	}
	m.Append("chunklist1.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
	m.Append("chunklist2.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
}

// Create new master playlist with Name tag in EXT-X-STREAM-INF
func TestEncodeMasterPlaylistWithStreamInfName(t *testing.T) {
	m := NewMasterPlaylist()
	p, e := NewMediaPlaylist(3, 5)
	if e != nil {
		t.Fatalf("Create media playlist failed: %s", e)
	}
	for i := 0; i < 5; i++ {
		e = p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
		if e != nil {
			t.Errorf("Add segment #%d to a media playlist failed: %s", i, e)
		}
	}
	m.Append("chunklist1.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 3000000, Resolution: "1152x960", Name: "HD 960p"})

	if m.Variants[0].Name != "HD 960p" {
		t.Fatalf("Create master with Name in EXT-X-STREAM-INF failed")
	}
	if !strings.Contains(m.String(), "NAME=\"HD 960p\"") {
		t.Fatalf("Encode master with Name in EXT-X-STREAM-INF failed")
	}
}

func TestMasterVersion(t *testing.T) {
	m := NewMasterPlaylist()
	m.ver = 5
	if m.Version() != m.ver {
		t.Errorf("Expected version: %v, got: %v", m.ver, m.Version())
	}
}

func TestMasterSetVersion(t *testing.T) {
	m := NewMasterPlaylist()
	m.ver = 3
	m.SetVersion(5)
	if m.ver != 5 {
		t.Errorf("Expected version: %v, got: %v", 5, m.ver)
	}
}

/******************************
 *  Code generation examples  *
 ******************************/

// Create new media playlist
// Add two segments to media playlist
// Print it
func ExampleMediaPlaylist_String() {
	p, _ := NewMediaPlaylist(1, 2)
	p.Append("test01.ts", 5.0, "")
	p.Append("test02.ts", 6.0, "")
	fmt.Printf("%s\n", p)
	// Output:
	// #EXTM3U
	// #EXT-X-VERSION:3
	// #EXT-X-MEDIA-SEQUENCE:0
	// #EXT-X-TARGETDURATION:6
	// #EXTINF:5.000,
	// test01.ts
}

// Create new media playlist
// Add two segments to media playlist
// Print it
func ExampleMediaPlaylist_String_Winsize0() {
	p, _ := NewMediaPlaylist(0, 2)
	p.Append("test01.ts", 5.0, "")
	p.Append("test02.ts", 6.0, "")
	fmt.Printf("%s\n", p)
	// Output:
	// #EXTM3U
	// #EXT-X-VERSION:3
	// #EXT-X-MEDIA-SEQUENCE:0
	// #EXT-X-TARGETDURATION:6
	// #EXTINF:5.000,
	// test01.ts
	// #EXTINF:6.000,
	// test02.ts
}

// Create new media playlist
// Add two segments to media playlist
// Print it
func ExampleMediaPlaylist_String_Winsize0_VOD() {
	p, _ := NewMediaPlaylist(0, 2)
	p.Append("test01.ts", 5.0, "")
	p.Append("test02.ts", 6.0, "")
	p.Close()
	fmt.Printf("%s\n", p)
	// Output:
	// #EXTM3U
	// #EXT-X-VERSION:3
	// #EXT-X-MEDIA-SEQUENCE:0
	// #EXT-X-TARGETDURATION:6
	// #EXTINF:5.000,
	// test01.ts
	// #EXTINF:6.000,
	// test02.ts
	// #EXT-X-ENDLIST
}

// Create new master playlist
// Add media playlist
// Encode structures to HLS
func ExampleMasterPlaylist_String() {
	m := NewMasterPlaylist()
	p, _ := NewMediaPlaylist(3, 5)
	for i := 0; i < 5; i++ {
		p.Append(fmt.Sprintf("test%d.ts", i), 5.0, "")
	}
	m.Append("chunklist1.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
	m.Append("chunklist2.m3u8", p, VariantParams{ProgramId: 123, Bandwidth: 1500000, Resolution: "576x480"})
	fmt.Printf("%s", m)
	// Output:
	// #EXTM3U
	// #EXT-X-VERSION:3
	// #EXT-X-STREAM-INF:PROGRAM-ID=123,BANDWIDTH=1500000,RESOLUTION=576x480
	// chunklist1.m3u8
	// #EXT-X-STREAM-INF:PROGRAM-ID=123,BANDWIDTH=1500000,RESOLUTION=576x480
	// chunklist2.m3u8
}

/****************
 *  Benchmarks  *
 ****************/

func BenchmarkEncodeMasterPlaylist(b *testing.B) {
	f, err := os.Open("sample-playlists/master.m3u8")
	if err != nil {
		b.Fatal(err)
	}
	p := NewMasterPlaylist()
	if err := p.DecodeFrom(bufio.NewReader(f), true); err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		p.ResetCache()
		_ = p.Encode() // disregard output
	}
}

func BenchmarkEncodeMediaPlaylist(b *testing.B) {
	f, err := os.Open("sample-playlists/media-playlist-large.m3u8")
	if err != nil {
		b.Fatal(err)
	}
	p, err := NewMediaPlaylist(50000, 50000)
	if err != nil {
		b.Fatalf("Create media playlist failed: %s", err)
	}
	if err = p.DecodeFrom(bufio.NewReader(f), true); err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		p.ResetCache()
		_ = p.Encode() // disregard output
	}
}
