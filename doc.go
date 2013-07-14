// Copyleft 2013 Alexander I.Grafov aka Axel <grafov@gmail.com>
// Library licensed under GPLv3
//
// ॐ तारे तुत्तारे तुरे स्व

/*

__This is only draft of the library. API may be changed!__

Library may be used for parsing and generation of M3U8 playlists. M3U8 format used in HTTP Live Streaming (Apple HLS) for internet video translations. Also the library may be useful for common M3U format parsing and generation.

Planned support of specific extensions such as Widevine or Verimatrix.

Library coded acordingly with http://tools.ietf.org/html/draft-pantos-http-live-streaming-11

Examples of usage may be found in *_test.go files of a package. Also see below some simple examples.

Create simple media playlist with sliding window of 3 segments and maximum of 50 segments.

        p, e := NewMediaPlaylist(3, 50)
        if e != nil {
          panic(fmt.Sprintf("Create media playlist failed: %s", e))
        }
        for i := 0; i < 5; i++ {
          e = p.Add(fmt.Sprintf("test%d.ts", i), 5.0)
          if e != nil {
            panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
          }
        }
        fmt.Println(p.Encode(true).String())

We add 5 testX.ts segments to playlist then encode it to M3U8 format and convert to string.

Next example shows parsing of master playlist:

        f, err := os.Open("sample-playlists/master.m3u8")
        if err != nil {
          fmt.Println(err)
        }
        p := NewMasterPlaylist()
        err = p.Decode(bufio.NewReader(f), false)
        if err != nil {
          fmt.Println(err)
        }

        fmt.Printf("Playlist object: %+v\n", p)

We are open playlist from the file and parse it as master playlist.

*/
package m3u8
