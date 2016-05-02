/*
 Playlist parsing tests.

 Copyleft 2013-2015 Alexander I.Grafov aka Axel <grafov@gmail.com>
 Copyleft 2013-2015 library authors (see AUTHORS file).

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
package m3u8

import (
	"bufio"
	"fmt"
	"os"
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
	// TODO check other values…
	//fmt.Println(p.Encode().String()), stream.Name}
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
	print("test")
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
	// TODO check other values…
	// fmt.Println(pp.Encode().String())
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
