package m3u8

/*
 M3U8 v3 playlists for HTTP Live Streaming. Generator and parser.
 Coded acordingly with http://tools.ietf.org/html/draft-pantos-http-live-streaming-10

 Copyleft Alexander I.Grafov aka Axel <grafov@gmail.com>
 Library licensed under GPLv3

 ॐ तारे तुत्तारे तुरे स्व
*/

const (
	M3U8Version = "3"
)

// General playlist structure
type M3U8 struct {
	Name string
	Key Key
	Playlists []Playlist
	Duration uint
	Segments []Segment
	SeqNo uint32
}

// Playlists included in a variant playlist
type Playlist struct {
	ProgramId uint8
	Bandwidth uint32
	URI string
}

// Media segments included in a playlist
type Segment struct {
	Duration float64
	URI string
}

// Information about stream encryption
type Key struct {
	Method string
	URI string
	IV string
}

