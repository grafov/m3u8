package m3u8

import (
	"fmt"
	"testing"
)

func TestLoadEncodeString(t *testing.T) {
	v, _ := NewMediaPlaylist(0, 10)
	v.SetVersion(6)
	v.SeqNo = 0
	v.Closed = true
	v.MediaType = 2

	for i := 0; i < 10; i++ {
		v.Append(fmt.Sprintf("https://host/identifier/quality/720p_segment_%d.ts", i), 4, "")
	}

	for i := 0; i < 100; i++ {
		go loadRender(v, t)
	}
}

const expectedPlaylist = `#EXTM3U
#EXT-X-VERSION:6
#EXT-X-PLAYLIST-TYPE:VOD
#EXT-X-MEDIA-SEQUENCE:0
#EXT-X-TARGETDURATION:4
#EXTINF:4.000,
https://host/identifier/quality/720p_segment_0.ts
#EXTINF:4.000,
https://host/identifier/quality/720p_segment_1.ts
#EXTINF:4.000,
https://host/identifier/quality/720p_segment_2.ts
#EXTINF:4.000,
https://host/identifier/quality/720p_segment_3.ts
#EXTINF:4.000,
https://host/identifier/quality/720p_segment_4.ts
#EXTINF:4.000,
https://host/identifier/quality/720p_segment_5.ts
#EXTINF:4.000,
https://host/identifier/quality/720p_segment_6.ts
#EXTINF:4.000,
https://host/identifier/quality/720p_segment_7.ts
#EXTINF:4.000,
https://host/identifier/quality/720p_segment_8.ts
#EXTINF:4.000,
https://host/identifier/quality/720p_segment_9.ts
#EXT-X-ENDLIST
`

// loadRender will render the playlist in a loop 500 times
// the intention is to launch this concurrently in goroutines to simulate high traffic
func loadRender(v *MediaPlaylist, t *testing.T) {
	for i := 0; i < 500; i++ {
		if v.String() != expectedPlaylist {
			t.Fatalf("expected:\n %s\n\n got:\n %s", expectedPlaylist, v.String())
		}
	}
}
