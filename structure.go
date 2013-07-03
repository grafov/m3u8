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

// Single bitrate playlist.
// It presents both simple media playlists and sliding window media playlists.
// All URI lines in the Playlist identify media segments.
type MediaPlaylist struct {
	TargetDuration float64
	SeqNo          uint64
	segments       []MediaSegment
	SID            string
	Iframe         bool // EXT-X-I-FRAMES-ONLY
	key            *Key
	wv             *WV
	keyformat      int
	winsize        uint16 // size of visible window
	capacity       uint16 // total capacity of slice used for the playlist
	buf            *bytes.Buffer
	ver            uint8
}

// Master playlist combines media playlists for multiple bitrates.
// All URI lines in the Playlist identify Media Playlists.
type MasterPlaylist struct {
	SID      string
	variants []Variant
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
	AltMedia   []AltMedia
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
