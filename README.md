<!--*- mode:markdown -*-->
M3U8
====

This is a most complete opensource library for parsing and generating of M3U8 playlists
used in HTTP Live Streaming (Apple HLS) for internet video translations.

Features are:

* Support HLS specs up to IETF draft 11 (version 3 of the protocol).
* Parsing and generation of master-playlists and media-playlists.
* Offers structures for keeping playlists metadata.
* Encryption keys support for usage with DRM systems like Verimatrix etc.
* Support for non standard [Google Widevine](http://www.widevine.com) tags.

Copyleft Alexander I.Grafov aka Axel <grafov@gmail.com>

Library licensed under GPLv3.

Install
-------

	go get github.com/grafov/m3u8

or get releases from https://github.com/grafov/m3u8/releases

Documentation
-------------

Package online documentation (examples included) available at:

* http://gowalker.org/github.com/grafov/m3u8
* http://godoc.org/github.com/grafov/m3u8

Supported by the HLS protocol tags and their library support explained in [M3U8 cheatsheet](M3U8.md).

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
* http://gonze.com/playlists/playlist-format-survey.html

Library usage samples in other software
---------------------------------------

* [Stream Surfer](http://streamsurfer.org) monitoring software.
* [gohls](https://github.com/kz26/gohls) — HLS downloader.

M3U8 parsing/generation in other languages
------------------------------------------

* https://github.com/globocom/m3u8 in Python
* https://github.com/zencoder/m3uzi in Ruby
* https://github.com/Jeanvf/M3U8Paser in Objective C
* https://github.com/tedconf/node-m3u8 in Javascript
* http://sourceforge.net/projects/m3u8parser/ in Java
* https://github.com/karlll/erlm3u8 in Erlang

Project status [![Is maintained?](http://stillmaintained.com/grafov/m3u8.png)](http://stillmaintained.com/grafov/m3u8)
---------------

In development.

[![Build Status](https://travis-ci.org/grafov/m3u8.png?branch=master)](https://travis-ci.org/grafov/m3u8) for `master` branch.

[![Build Status](https://drone.io/github.com/grafov/m3u8/status.png)](https://drone.io/github.com/grafov/m3u8/latest) for `master` branch.

Development rules:

* Changes firstly applied to `draft` branch then after minimal testing with unit tests it merged with `master` branch.
* Code in `draft` branch may be in inconsistent state.
* After complete testing and one week usage with my prober for HLS [Stream Surfer](http://streamsurfer.org) it may be released as new library version (minor version number increased).
* Each new API call or tag accompanied by at least one unit test until it released in a new library version.
* Major version number increased on API changes or new tags (until v1.0 this rule is not applied).
* I reserve the right to break the development rules if my life plans make me do it.

Also see project roadmap and planned features in [TODO](TODO.org).
