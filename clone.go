package m3u8

/*
 Part of M3U8 parser & generator library.
 This file defines data structures related to package.

 Copyright 2013-2017 The Project Developers.
 See the AUTHORS and LICENSE files at the top-level directory of this distribution
 and at https://github.com/grafov/m3u8/

 ॐ तारे तुत्तारे तुरे स्व
*/

func (p *MediaPlaylist) Clone() *MediaPlaylist {
	var newSegments []*MediaSegment
	for _, segment := range p.Segments {
		newSegments = append(newSegments, segment.Clone())
	}

	mp := &MediaPlaylist{
		TargetDuration:   p.TargetDuration,
		SeqNo:            p.SeqNo,
		Segments:         newSegments,
		Args:             p.Args,
		Iframe:           p.Iframe,
		Closed:           p.Closed,
		MediaType:        p.MediaType,
		DiscontinuitySeq: p.DiscontinuitySeq,
		StartTime:        p.StartTime,
		StartTimePrecise: p.StartTimePrecise,
		durationAsInt:    p.durationAsInt,
		keyformat:        p.keyformat,
		winsize:          p.winsize,
		capacity:         p.capacity,
		head:             p.head,
		tail:             p.tail,
		count:            p.count,
		buf:              p.buf,
		ver:              p.ver,
	}

	if p.Key != nil {
		mp.Key = &*p.Key
	}
	if p.Map != nil {
		mp.Map = &*p.Map
	}
	if p.WV != nil {
		mp.WV = &*p.WV
	}

	return mp
}

func (s *MediaSegment) Clone() *MediaSegment {
	ms := &MediaSegment{
		SeqId:           s.SeqId,
		Title:           s.Title,
		URI:             s.URI,
		Duration:        s.Duration,
		Attributes:      s.Attributes,
		Limit:           s.Limit,
		Offset:          s.Offset,
		Discontinuity:   s.Discontinuity,
		ProgramDateTime: s.ProgramDateTime,
	}

	if s.Key != nil {
		ms.Key = &*s.Key
	}
	if s.Map != nil {
		ms.Map = &*s.Map
	}
	if s.SCTE != nil {
		ms.SCTE = &*s.SCTE
	}

	return ms
}

func (a *Attribute) Clone() *Attribute {
	return &Attribute{
		Key:   a.Key,
		Value: a.Value,
	}
}
