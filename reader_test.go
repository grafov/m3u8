/*
 Playlist parsing tests.

 Copyright 2013-2017 The Project Developers.
 See the AUTHORS and LICENSE files at the top-level directory of this distribution
 and at https://github.com/grafov/m3u8/

 ॐ तारे तुत्तारे तुरे स्व
*/
package m3u8

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestDecodeMasterPlaylist(t *testing.T) {
	f, err := os.Open("sample-playlists/master.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := NewMasterPlaylist()
	err = p.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		t.Fatal(err)
	}
	// check parsed values
	if p.ver != 3 {
		t.Errorf("Version of parsed playlist = %d (must = 3)", p.ver)
	}
	if len(p.Variants) != 5 {
		t.Error("Not all variants in master playlist parsed.")
	}
	// TODO check other values
	// fmt.Println(p.Encode().String())
}

func TestDecodeMasterPlaylistWithMultipleCodecs(t *testing.T) {
	f, err := os.Open("sample-playlists/master-with-multiple-codecs.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := NewMasterPlaylist()
	err = p.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		t.Fatal(err)
	}
	// check parsed values
	if p.ver != 3 {
		t.Errorf("Version of parsed playlist = %d (must = 3)", p.ver)
	}
	if len(p.Variants) != 5 {
		t.Error("Not all variants in master playlist parsed.")
	}
	for _, v := range p.Variants {
		if v.Codecs != "avc1.42c015,mp4a.40.2" {
			t.Error("Codec string is wrong")
		}
	}
	// TODO check other values
	// fmt.Println(p.Encode().String())
}

func TestDecodeMasterPlaylistWithAlternatives(t *testing.T) {
	f, err := os.Open("sample-playlists/master-with-alternatives.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := NewMasterPlaylist()
	err = p.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		t.Fatal(err)
	}
	// check parsed values
	if p.ver != 3 {
		t.Errorf("Version of parsed playlist = %d (must = 3)", p.ver)
	}
	if len(p.Variants) != 4 {
		t.Fatal("not all variants in master playlist parsed")
	}
	// TODO check other values
	for i, v := range p.Variants {
		if i == 0 && len(v.Alternatives) != 3 {
			t.Fatalf("not all alternatives from #EXT-X-MEDIA parsed (has %d but should be 3", len(v.Alternatives))
		}
		if i == 1 && len(v.Alternatives) != 3 {
			t.Fatalf("not all alternatives from #EXT-X-MEDIA parsed (has %d but should be 3", len(v.Alternatives))
		}
		if i == 2 && len(v.Alternatives) != 3 {
			t.Fatalf("not all alternatives from #EXT-X-MEDIA parsed (has %d but should be 3", len(v.Alternatives))
		}
		if i == 3 && len(v.Alternatives) > 0 {
			t.Fatal("should not be alternatives for this variant")
		}
	}
	// fmt.Println(p.Encode().String())
}

func TestDecodeMasterPlaylistWithClosedCaptionEqNone(t *testing.T) {
	f, err := os.Open("sample-playlists/master-with-closed-captions-eq-none.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := NewMasterPlaylist()
	err = p.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		t.Fatal(err)
	}

	if len(p.Variants) != 3 {
		t.Fatal("not all variants in master playlist parsed")
	}
	for _, v := range p.Variants {
		if v.Captions != "NONE" {
			t.Fatal("variant field for CLOSED-CAPTIONS should be equal to NONE but it equals", v.Captions)
		}
	}
}

// Decode a master playlist with Name tag in EXT-X-STREAM-INF
func TestDecodeMasterPlaylistWithStreamInfName(t *testing.T) {
	f, err := os.Open("sample-playlists/master-with-stream-inf-name.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := NewMasterPlaylist()
	err = p.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, variant := range p.Variants {
		if variant.Name == "" {
			t.Errorf("Empty name tag on variant URI: %s", variant.URI)
		}
	}
}

func TestDecodeMediaPlaylistByteRange(t *testing.T) {
	f, _ := os.Open("sample-playlists/media-playlist-with-byterange.m3u8")
	p, _ := NewMediaPlaylist(3, 3)
	_ = p.DecodeFrom(bufio.NewReader(f), true)
	expected := []*MediaSegment{
		{URI: "video.ts", Duration: 10, Limit: 75232},
		{URI: "video.ts", Duration: 10, Limit: 82112, Offset: 752321},
		{URI: "video.ts", Duration: 10, Limit: 69864},
	}
	for i, seg := range p.Segments {
		if *seg != *expected[i] {
			t.Errorf("exp: %+v\ngot: %+v", expected[i], seg)
		}
	}
}

// Decode a master playlist with i-frame-stream-inf
func TestDecodeMasterPlaylistWithIFrameStreamInf(t *testing.T) {
	f, err := os.Open("sample-playlists/master-with-i-frame-stream-inf.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := NewMasterPlaylist()
	err = p.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		t.Fatal(err)
	}
	expected := map[int]*Variant{
		86000:  {URI: "low/iframe.m3u8", VariantParams: VariantParams{Bandwidth: 86000, ProgramId: 1, Codecs: "c1", Resolution: "1x1", Video: "1", Iframe: true}},
		150000: {URI: "mid/iframe.m3u8", VariantParams: VariantParams{Bandwidth: 150000, ProgramId: 1, Codecs: "c2", Resolution: "2x2", Video: "2", Iframe: true}},
		550000: {URI: "hi/iframe.m3u8", VariantParams: VariantParams{Bandwidth: 550000, ProgramId: 1, Codecs: "c2", Resolution: "2x2", Video: "2", Iframe: true}},
	}
	for _, variant := range p.Variants {
		for k, expect := range expected {
			if reflect.DeepEqual(variant, expect) {
				delete(expected, k)
			}
		}
	}
	for _, expect := range expected {
		t.Errorf("not found:%+v", expect)
	}
}

func TestDecodeMasterPlaylistWithStreamInfAverageBandwidth(t *testing.T) {
	f, err := os.Open("sample-playlists/master-with-stream-inf-1.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := NewMasterPlaylist()
	err = p.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, variant := range p.Variants {
		if variant.AverageBandwidth == 0 {
			t.Errorf("Empty average bandwidth tag on variant URI: %s", variant.URI)
		}
	}
}

func TestDecodeMasterPlaylistWithStreamInfFrameRate(t *testing.T) {
	f, err := os.Open("sample-playlists/master-with-stream-inf-1.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := NewMasterPlaylist()
	err = p.DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, variant := range p.Variants {
		if variant.FrameRate == 0 {
			t.Errorf("Empty frame rate tag on variant URI: %s", variant.URI)
		}
	}
}

func TestDecodeMediaPlaylist(t *testing.T) {
	f, err := os.Open("sample-playlists/wowza-vod-chunklist.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p, err := NewMediaPlaylist(5, 798)
	if err != nil {
		t.Fatalf("Create media playlist failed: %s", err)
	}
	err = p.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("Playlist object: %+v\n", p)
	// check parsed values
	if p.ver != 3 {
		t.Errorf("Version of parsed playlist = %d (must = 3)", p.ver)
	}
	if p.TargetDuration != 12 {
		t.Errorf("TargetDuration of parsed playlist = %f (must = 12.0)", p.TargetDuration)
	}
	if !p.Closed {
		t.Error("This is a closed (VOD) playlist but Close field = false")
	}
	titles := []string{"Title 1", "Title 2", ""}
	for i, s := range p.Segments {
		if i > len(titles)-1 {
			break
		}
		if s.Title != titles[i] {
			t.Errorf("Segment %v's title = %v (must = %q)", i, s.Title, titles[i])
		}
	}
	// TODO check other values…
	//fmt.Println(p.Encode().String()), stream.Name}
}

func TestDecodeMediaPlaylistExtInfNonStrict2(t *testing.T) {
	header := `#EXTM3U
#EXT-X-TARGETDURATION:10
#EXT-X-VERSION:3
#EXT-X-MEDIA-SEQUENCE:0
%s
`

	tests := []struct {
		strict      bool
		extInf      string
		wantError   bool
		wantSegment *MediaSegment
	}{
		// strict mode on
		{true, "#EXTINF:10.000,", false, &MediaSegment{Duration: 10.0, Title: ""}},
		{true, "#EXTINF:10.000,Title", false, &MediaSegment{Duration: 10.0, Title: "Title"}},
		{true, "#EXTINF:10.000,Title,Track", false, &MediaSegment{Duration: 10.0, Title: "Title,Track"}},
		{true, "#EXTINF:invalid,", true, nil},
		{true, "#EXTINF:10.000", true, nil},

		// strict mode off
		{false, "#EXTINF:10.000,", false, &MediaSegment{Duration: 10.0, Title: ""}},
		{false, "#EXTINF:10.000,Title", false, &MediaSegment{Duration: 10.0, Title: "Title"}},
		{false, "#EXTINF:10.000,Title,Track", false, &MediaSegment{Duration: 10.0, Title: "Title,Track"}},
		{false, "#EXTINF:invalid,", false, &MediaSegment{Duration: 0.0, Title: ""}},
		{false, "#EXTINF:10.000", false, &MediaSegment{Duration: 10.0, Title: ""}},
	}

	for _, test := range tests {
		p, err := NewMediaPlaylist(1, 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		reader := bytes.NewBufferString(fmt.Sprintf(header, test.extInf))
		err = p.DecodeFrom(reader, test.strict)
		if test.wantError {
			if err == nil {
				t.Errorf("expected error but have: %v", err)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(p.Segments[0], test.wantSegment) {
			t.Errorf("\nhave: %+v\nwant: %+v", p.Segments[0], test.wantSegment)
		}
	}
}

func TestDecodeMediaPlaylistWithWidevine(t *testing.T) {
	f, err := os.Open("sample-playlists/widevine-bitrate.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p, err := NewMediaPlaylist(5, 798)
	if err != nil {
		t.Fatalf("Create media playlist failed: %s", err)
	}
	err = p.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("Playlist object: %+v\n", p)
	// check parsed values
	if p.ver != 2 {
		t.Errorf("Version of parsed playlist = %d (must = 2)", p.ver)
	}
	if p.TargetDuration != 9 {
		t.Errorf("TargetDuration of parsed playlist = %f (must = 9.0)", p.TargetDuration)
	}
	// TODO check other values…
	//fmt.Printf("%+v\n", p.Key)
	//fmt.Println(p.Encode().String())
}

func TestDecodeMasterPlaylistWithAutodetection(t *testing.T) {
	f, err := os.Open("sample-playlists/master.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	m, listType, err := DecodeFrom(bufio.NewReader(f), false)
	if err != nil {
		t.Fatal(err)
	}
	if listType != MASTER {
		t.Error("Sample not recognized as master playlist.")
	}
	mp := m.(*MasterPlaylist)
	// fmt.Printf(">%+v\n", mp)
	// for _, v := range mp.Variants {
	// 	fmt.Printf(">>%+v +v\n", v)
	// }
	//fmt.Println("Type below must be MasterPlaylist:")
	CheckType(t, mp)
}

func TestDecodeMediaPlaylistWithAutodetection(t *testing.T) {
	f, err := os.Open("sample-playlists/wowza-vod-chunklist.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p, listType, err := DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		t.Fatal(err)
	}
	pp := p.(*MediaPlaylist)
	CheckType(t, pp)
	if listType != MEDIA {
		t.Error("Sample not recognized as media playlist.")
	}
	// check parsed values
	if pp.TargetDuration != 12 {
		t.Errorf("TargetDuration of parsed playlist = %f (must = 12.0)", pp.TargetDuration)
	}

	if !pp.Closed {
		t.Error("This is a closed (VOD) playlist but Close field = false")
	}
	if pp.winsize != 0 {
		t.Errorf("Media window size %v != 0", pp.winsize)
	}
	// TODO check other values…
	// fmt.Println(pp.Encode().String())
}

// TestDecodeMediaPlaylistAutoDetectExtend tests a very large playlist auto
// extends to the appropriate size.
func TestDecodeMediaPlaylistAutoDetectExtend(t *testing.T) {
	f, err := os.Open("sample-playlists/media-playlist-large.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p, listType, err := DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		t.Fatal(err)
	}
	pp := p.(*MediaPlaylist)
	CheckType(t, pp)
	if listType != MEDIA {
		t.Error("Sample not recognized as media playlist.")
	}
	var exp uint = 40001
	if pp.Count() != exp {
		t.Errorf("Media segment count %v != %v", pp.Count(), exp)
	}
}

// Test for FullTimeParse of EXT-X-PROGRAM-DATE-TIME
// We testing ISO/IEC 8601:2004 where we can get time in UTC, UTC with Nanoseconds
// timeZone in formats '±00:00', '±0000', '±00'
// m3u8.FullTimeParse()
func TestFullTimeParse(t *testing.T) {
	var timestamps = []struct {
		name  string
		value string
	}{
		{"time_in_utc", "2006-01-02T15:04:05Z"},
		{"time_in_utc_nano", "2006-01-02T15:04:05.123456789Z"},
		{"time_with_positive_zone_and_colon", "2006-01-02T15:04:05+01:00"},
		{"time_with_positive_zone_no_colon", "2006-01-02T15:04:05+0100"},
		{"time_with_positive_zone_2digits", "2006-01-02T15:04:05+01"},
		{"time_with_negative_zone_and_colon", "2006-01-02T15:04:05-01:00"},
		{"time_with_negative_zone_no_colon", "2006-01-02T15:04:05-0100"},
		{"time_with_negative_zone_2digits", "2006-01-02T15:04:05-01"},
	}

	var err error
	for _, tstamp := range timestamps {
		_, err = FullTimeParse(tstamp.value)
		if err != nil {
			t.Errorf("FullTimeParse Error at %s [%s]: %s", tstamp.name, tstamp.value, err)
		}
	}
}

// Test for StrictTimeParse of EXT-X-PROGRAM-DATE-TIME
// We testing Strict format of RFC3339 where we can get time in UTC, UTC with Nanoseconds
// timeZone in formats '±00:00', '±0000', '±00'
// m3u8.StrictTimeParse()
func TestStrictTimeParse(t *testing.T) {
	var timestamps = []struct {
		name  string
		value string
	}{
		{"time_in_utc", "2006-01-02T15:04:05Z"},
		{"time_in_utc_nano", "2006-01-02T15:04:05.123456789Z"},
		{"time_with_positive_zone_and_colon", "2006-01-02T15:04:05+01:00"},
		{"time_with_negative_zone_and_colon", "2006-01-02T15:04:05-01:00"},
	}

	var err error
	for _, tstamp := range timestamps {
		_, err = StrictTimeParse(tstamp.value)
		if err != nil {
			t.Errorf("StrictTimeParse Error at %s [%s]: %s", tstamp.name, tstamp.value, err)
		}
	}
}

func TestMediaPlaylistWithOATCLSSCTE35Tag(t *testing.T) {
	f, err := os.Open("sample-playlists/media-playlist-with-oatcls-scte35.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p, _, err := DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		t.Fatal(err)
	}
	pp := p.(*MediaPlaylist)

	expect := map[int]*SCTE{
		0: {Syntax: SCTE35_OATCLS, CueType: SCTE35Cue_Start, Cue: "/DAlAAAAAAAAAP/wFAUAAAABf+/+ANgNkv4AFJlwAAEBAQAA5xULLA==", Time: 15},
		1: {Syntax: SCTE35_OATCLS, CueType: SCTE35Cue_Mid, Cue: "/DAlAAAAAAAAAP/wFAUAAAABf+/+ANgNkv4AFJlwAAEBAQAA5xULLA==", Time: 15, Elapsed: 8.844},
		2: {Syntax: SCTE35_OATCLS, CueType: SCTE35Cue_End},
	}
	for i := 0; i < int(pp.Count()); i++ {
		if !reflect.DeepEqual(pp.Segments[i].SCTE, expect[i]) {
			t.Errorf("OATCLS SCTE35 segment %v (uri: %v)\ngot: %#v\nexp: %#v",
				i, pp.Segments[i].URI, pp.Segments[i].SCTE, expect[i],
			)
		}
	}
}

func TestDecodeMediaPlaylistWithDiscontinuitySeq(t *testing.T) {
	f, err := os.Open("sample-playlists/media-playlist-with-discontinuity-seq.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p, listType, err := DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		t.Fatal(err)
	}
	pp := p.(*MediaPlaylist)
	CheckType(t, pp)
	if listType != MEDIA {
		t.Error("Sample not recognized as media playlist.")
	}
	if pp.DiscontinuitySeq == 0 {
		t.Error("Empty discontinuity sequenece tag")
	}
}

/***************************
 *  Code parsing examples  *
 ***************************/

// Example of parsing a playlist with EXT-X-DISCONTINIUTY tag
// and output it with integer segment durations.
func ExampleMediaPlaylist_DurationAsInt() {
	f, _ := os.Open("sample-playlists/media-playlist-with-discontinuity.m3u8")
	p, _, _ := DecodeFrom(bufio.NewReader(f), true)
	pp := p.(*MediaPlaylist)
	pp.DurationAsInt(true)
	fmt.Printf("%s", pp)
	// Output:
	// #EXTM3U
	// #EXT-X-VERSION:3
	// #EXT-X-MEDIA-SEQUENCE:0
	// #EXT-X-TARGETDURATION:10
	// #EXTINF:10,
	// ad0.ts
	// #EXTINF:8,
	// ad1.ts
	// #EXT-X-DISCONTINUITY
	// #EXTINF:10,
	// movieA.ts
	// #EXTINF:10,
	// movieB.ts
}

func TestMediaPlaylistWithSCTE35Tag(t *testing.T) {
	test_cases := []struct {
		playlistLocation  string
		expectedSCTEIndex int
		expectedSCTECue   string
		expectedSCTEID    string
		expectedSCTETime  float64
	}{
		{
			"sample-playlists/media-playlist-with-scte35.m3u8",
			2,
			"/DAIAAAAAAAAAAAQAAZ/I0VniQAQAgBDVUVJQAAAAH+cAAAAAA==",
			"123",
			123.12,
		},
		{
			"sample-playlists/media-playlist-with-scte35-1.m3u8",
			1,
			"/DAIAAAAAAAAAAAQAAZ/I0VniQAQAgBDVUVJQAA",
			"",
			0,
		},
	}
	for _, c := range test_cases {
		f, _ := os.Open(c.playlistLocation)
		playlist, _, _ := DecodeFrom(bufio.NewReader(f), true)
		mediaPlaylist := playlist.(*MediaPlaylist)
		for index, item := range mediaPlaylist.Segments {
			if item == nil {
				break
			}
			if index != c.expectedSCTEIndex && item.SCTE != nil {
				t.Error("Not expecting SCTE information on this segment")
			} else if index == c.expectedSCTEIndex && item.SCTE == nil {
				t.Error("Expecting SCTE information on this segment")
			} else if index == c.expectedSCTEIndex && item.SCTE != nil {
				if (*item.SCTE).Cue != c.expectedSCTECue {
					t.Error("Expected ", c.expectedSCTECue, " got ", (*item.SCTE).Cue)
				} else if (*item.SCTE).ID != c.expectedSCTEID {
					t.Error("Expected ", c.expectedSCTEID, " got ", (*item.SCTE).ID)
				} else if (*item.SCTE).Time != c.expectedSCTETime {
					t.Error("Expected ", c.expectedSCTETime, " got ", (*item.SCTE).Time)
				}
			}
		}
	}
}

/****************
 *  Benchmarks  *
 ****************/

func BenchmarkDecodeMasterPlaylist(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f, err := os.Open("sample-playlists/master.m3u8")
		if err != nil {
			b.Fatal(err)
		}
		p := NewMasterPlaylist()
		if err := p.DecodeFrom(bufio.NewReader(f), false); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeMediaPlaylist(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
	}
}
