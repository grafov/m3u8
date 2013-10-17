package m3u8

/*
 Part of M3U8 parser & generator library.
 This file defines data structures related to package.

 Copyleft 2013  Alexander I.Grafov aka Axel <grafov@gmail.com>

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

import (
	"bytes"
)

const (
	/*
		Compatibility rules described in section 7:
		Clients and servers MUST implement protocol version 2 or higher to use:
		   o  The IV attribute of the EXT-X-KEY tag.
		   Clients and servers MUST implement protocol version 3 or higher to use:
		   o  Floating-point EXTINF duration values.
		   Clients and servers MUST implement protocol version 4 or higher to use:
		   o  The EXT-X-BYTERANGE tag.
		   o  The EXT-X-I-FRAME-STREAM-INF tag.
		   o  The EXT-X-I-FRAMES-ONLY tag.
		   o  The EXT-X-MEDIA tag.
		   o  The AUDIO and VIDEO attributes of the EXT-X-STREAM-INF tag.
	*/
	minver = uint8(3)
)

type ListType uint

const (
	UNKNOWN ListType = iota
	MASTER
	MEDIA
)

/*
 This structure represents a single bitrate playlist aka media playlist.
 It related to both a simple media playlists and a sliding window media playlists.
 URI lines in the Playlist point to media segments.

 Simple Media Playlist file sample:

   #EXTM3U
   #EXT-X-VERSION:3
   #EXT-X-TARGETDURATION:5220
   #EXTINF:5219.2,
   http://media.example.com/entire.ts
   #EXT-X-ENDLIST

 Sample of Sliding Window Media Playlist, using HTTPS:

   #EXTM3U
   #EXT-X-VERSION:3
   #EXT-X-TARGETDURATION:8
   #EXT-X-MEDIA-SEQUENCE:2680

   #EXTINF:7.975,
   https://priv.example.com/fileSequence2680.ts
   #EXTINF:7.941,
   https://priv.example.com/fileSequence2681.ts
   #EXTINF:7.975,
   https://priv.example.com/fileSequence2682.ts
*/
type MediaPlaylist struct {
	TargetDuration float64
	SeqNo          uint64 // EXT-X-MEDIA-SEQUENCE
	Segments       []*MediaSegment
	SID            string // optional session identifier (out of scope of HLS specs but useful in some cases)
	Iframe         bool   // EXT-X-I-FRAMES-ONLY
	Closed         bool   // is this VOD (closed) or Live (sliding) playlist?
	durationAsInt  bool   // output durations as integers of floats?
	keyformat      int
	winsize        uint // max number of segments displayed in an encoded playlist; need set to zero for VOD playlists
	capacity       uint // total capacity of slice used for the playlist
	head           uint // head of FIFO, we add segments to head
	tail           uint // tail of FIFO, we remove segments from tail
	count          uint // number of segments added to the playlist
	buf            bytes.Buffer
	ver            uint8
	Key            *Key // encryption key displayed before any segments
	WV             *WV  // Widevine related tags
}

/*
 This structure represents a master playlist which combines media playlists for multiple bitrates.
 URI lines in the playlist identify media playlists.
 Sample of Master Playlist file:

   #EXTM3U
   #EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1280000
   http://example.com/low.m3u8
   #EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=2560000
   http://example.com/mid.m3u8
   #EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=7680000
   http://example.com/hi.m3u8
   #EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=65000,CODECS="mp4a.40.5"
   http://example.com/audio-only.m3u8
*/
type MasterPlaylist struct {
	SID      string
	Variants []*Variant
	buf      bytes.Buffer
	ver      uint8
}

// This structure represents variants for master playlist.
// Variants included in a master playlist and point to media playlists.
type Variant struct {
	URI       string
	Chunklist *MediaPlaylist
	VariantParams
}

// This stucture represents additional parameters for a variant
type VariantParams struct {
	ProgramId  uint8
	Bandwidth  uint32
	Codecs     string
	Resolution string
	Audio      string
	Video      string
	Subtitles  string
	Iframe     bool // EXT-X-I-FRAME-STREAM-INF
	AltMedia   []*AltMedia
}

// This structure represents EXT-X-MEDIA tag in variants.
type AltMedia struct {
	GroupId         string
	URI             string
	Type            string
	Language        string
	Name            string
	Default         string
	Autoselect      string
	Forced          string
	Characteristics string
	Subtitles       string
}

// This structure represents a media segment included in a media playlist.
// Media segment may be encrypted.
// Widevine supports own tags for encryption metadata.
type MediaSegment struct {
	SeqId    uint64
	Title    string // optional second parameter for EXTINF tag
	URI      string
	Duration float64 // first parameter for EXTINF tag; duration must be integers if protocol version is less than 3 but we are always keep them float
	Key      *Key    // displayed before the segment and means changing of encryption key (in theory each segment may have own key)
}

// This structure represents information about stream encryption.
// Realizes EXT-X-KEY tag.
type Key struct {
	Method            string
	URI               string
	IV                string
	Keyformat         string
	Keyformatversions string
}

// This structure represents metadata  for Google Widevine playlists.
// This format not described in IETF draft but provied by Widevine packager as
// additional tags in the playlist.
type WV struct {
	AudioChannels          uint
	AudioFormat            uint
	AudioProfileIDC        uint
	AudioSampleSize        uint
	AudioSamplingFrequency uint
	CypherVersion          string
	ECM                    string
	VideoFormat            uint
	VideoFrameRate         uint
	VideoLevelIDC          uint
	VideoProfileIDC        uint
	VideoResolution        string
	VideoSAR               string
}

// Interface applied to various playlist types.
type Playlist interface {
	Encode() *bytes.Buffer
	Decode(bytes.Buffer, bool) error
}
