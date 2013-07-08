package m3u8

/*
 M3U8 v3 playlists for HTTP Live Streaming. Generator and parser.
 Coded acordingly with http://tools.ietf.org/html/draft-pantos-http-live-streaming-11

 Copyleft Alexander I.Grafov aka Axel <grafov@gmail.com>
 Library licensed under GPLv3

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

/*
 Single bitrate playlist.
 It presents both simple media playlists and sliding window media playlists.
 All URI lines in the Playlist identify media segments.

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
	SeqNo          uint64
	segments       []*MediaSegment
	SID            string
	Iframe         bool // EXT-X-I-FRAMES-ONLY
	keyformat      int
	winsize        uint // max number of segments removed from queue on playlist generation
	capacity       uint // total capacity of slice used for the playlist
	head           uint // head of FIFO, we add segments to head
	tail           uint // tail of FIFO, we remove segments from tail
	count          uint // number of segments in the playlist
	buf            bytes.Buffer
	ver            uint8
}

/*
 Master playlist combines media playlists for multiple bitrates.
 All URI lines in the Playlist identify Media Playlists.
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
	variants []*Variant
	buf      bytes.Buffer
	ver      uint8
}

// Variants are items included in a master playlist. They linked to media playlists.
type Variant struct {
	ProgramId  uint8
	URI        string
	Bandwidth  uint32
	Codecs     string
	Resolution string
	Audio      string
	Video      string
	Subtitles  string
	Iframe     bool // EXT-X-I-FRAME-STREAM-INF
	AltMedia   []*AltMedia
	medialist  *MediaPlaylist
}

// Realizes EXT-X-MEDIA.
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

// Media segment included in a playlist.
type MediaSegment struct {
	SeqId    uint64
	URI      string
	Duration float64
	Key      *Key
	WV       *WV
}

// Information about stream encryption.
// Realizes EXT-X-KEY.
type Key struct {
	Method            string
	URI               string
	IV                string
	Keyformat         string
	Keyformatversions string
}

// Service information for Google Widevine playlists.
// This format not described in IETF draft but provied by Widevine packager as
// additional tags in the playlist.
type WV struct {
	AudioChannels        int
	AudioFormat          int
	AudioProfileIDC      int
	AudioSampleSize      int
	AudioSampleFrequency int
	CypherVersion        string
	ECM                  string
	VideoFormat          int
	VideoFrameRate       int
	VideoLevelIDC        int
	VideoProfileIDC      int
	VideoResolution      string
	VideoSAR             string
}
