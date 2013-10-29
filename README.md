M3U8
====

The library for parsing and generation of M3U8 playlists. M3U8 playlist format used in HTTP Live Streaming (Apple HLS) for internet video translations.
Also the library may be useful for common M3U format parsing and generation.

Features are:

* Structures for keeping playlists metadata.
* Parsing and generation of master-playlists and media-playlists.
* Encryption keys support for usage with DRM systems like Verimatrix etc.
* Support of non standard [Google Widevine](http://www.widevine.com) tags.

Copyleft Alexander I.Grafov aka Axel <grafov@gmail.com>

Library licensed under GPLv3.

Install
-------

	go get github.com/grafov/m3u8

Documentation
-------------

Package online documentation (examples included) available at http://godoc.org/github.com/grafov/m3u8

Examples
--------

Parse playlist:

```go
	f, err := os.Open("playlist.m3u8")
	if err != nil {
		panic(err)
	}
	p, listType, err := DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}
	switch listType {
	case MEDIA:
	    mediapl := p.(*MediaPlaylist)
		fmt.Printf("%+v\n", mediapl)
	case MASTER:
	    masterpl := p.(*MasterPlaylist)
		fmt.Printf("%+v\n", masterpl)
	}

```

Then you get filled with parsed data structures. For master playlists you get ``Master`` struct with slice consists of pointers to ``Variant`` structures (which represent playlists to each bitrate).
For media playlist parser returns ``MediaPlaylist`` structure with slice of ``Segments``. Each segment is of ``MediaSegment`` type.
See ``structure.go`` or full documentation (link below).

You may use API methods to fill structures or create them manually to generate playlists. Example of media playlist generation:

```go
	p, e := NewMediaPlaylist(3, 10) // with window of size 3 and capacity 10
	if e != nil {
		panic(fmt.Sprintf("Creating of media playlist failed: %s", e))
	}
	for i := 0; i < 5; i++ {
		e = p.Add(fmt.Sprintf("test%d.ts", i), 6.0, "")
		if e != nil {
			panic(fmt.Sprintf("Add segment #%d to a media playlist failed: %s", i, e))
		}
	}
	fmt.Println(p.Encode().String())
```

Library structure
-----------------

Library has compact code and bundled in three files:

* `structure.go` — declares all structures related to playlists and their properties
* `reader.go` — playlist parser methods
* `writer.go` — playlist generator methods

Each file has own test suite placed in `*_test.go` accordingly.

Related links
-------------

* http://en.wikipedia.org/wiki/M3U
* http://en.wikipedia.org/wiki/HTTP_Live_Streaming
* http://tools.ietf.org/html/draft-pantos-http-live-streaming-11
* http://gonze.com/playlists/playlist-format-survey.html

`m3u8` library used in [Stream Surfer](http://streamsurfer.org) monitoring software.

Project status
--------------

In development. API may be changed in a future.

[![Build Status](https://travis-ci.org/grafov/m3u8.png?branch=master)](https://travis-ci.org/grafov/m3u8) for `master` branch.

[![Is maintained?](http://stillmaintained.com/grafov/m3u8.png)](http://stillmaintained.com/grafov/m3u8)
