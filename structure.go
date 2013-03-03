package m3u8

/*
 M3U8 v3 playlists for HTTP Live Streaming. Generator and parser.
 Coded acordingly with http://tools.ietf.org/html/draft-pantos-http-live-streaming-10

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

// Simple playlist with fixed duration and with all segments 
// referenced from the single playlist file.
type FixedPlaylist struct {
	TargetDuration float64
	Segments       []Segment
	SID            string
	ver            uint8
}

type VariantPlaylist struct {
	ver      uint8
	Variants []Variant
	SID      string
}

// Playlist with sliding window
type SlidingPlaylist struct {
	TargetDuration float64
	SeqNo          uint64
	Segments       chan Segment
	SID            string
	key            *Key
	wv             *WV
	keyformat      int
	winsize        uint16
	cache          bytes.Buffer
	ver            uint8
}

// Variants included in a variant playlist
type Variant struct {
	ProgramId  uint8
	URI        string
	Bandwidth  uint32
	Codecs     string
	Resolution string
	Audio      string
	Video      string
	Subtitles  string
	Iframe     bool
	AltMedia   []AltMedia
}

// Realizes EXT-X-MEDIA
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

// Media segment included in a playlist
type Segment struct {
	SeqId    uint64
	URI      string
	Duration float64
	Key      *Key
	WV       *WV
}

// Information about stream encryption
type Key struct {
	Method            string
	URI               string
	IV                string
	Keyformat         string
	Keyformatversions string
}

// Additional information for Widevine
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

type Playlist interface {
	Buffer() *bytes.Buffer
}
