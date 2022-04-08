package wavebin

import (
	"encoding/binary"

	"github.com/karupanerura/riffbin"
)

type FormatChunk interface {
	ChunkProvider
	Bytes() []byte
}

type BasicFormatChunk struct {
	MetaFormat
}

func (f *BasicFormatChunk) Bytes() (b []byte) {
	b = make([]byte, 14)

	binary.LittleEndian.PutUint16(b[:2], f.CompressionCode())
	binary.LittleEndian.PutUint16(b[2:4], f.Channels())
	binary.LittleEndian.PutUint32(b[4:8], f.SamplesPerSecond())
	binary.LittleEndian.PutUint32(b[8:12], f.AverageBytesPerSecond())
	binary.LittleEndian.PutUint16(b[12:14], f.BlockAlign())

	return
}

func (f *BasicFormatChunk) Chunk() riffbin.Chunk {
	return &riffbin.OnMemorySubChunk{
		ID:      fmtBytes,
		Payload: f.Bytes(),
	}
}

type ExtendedFormatChunk struct {
	MetaFormat
}

func (f *ExtendedFormatChunk) Bytes() (b []byte) {
	s := 16
	ef := f.ExtraField()
	if len(ef) != 0 {
		s += 2 + len(ef)
	}
	b = make([]byte, s)

	binary.LittleEndian.PutUint16(b[:2], f.CompressionCode())
	binary.LittleEndian.PutUint16(b[2:4], f.Channels())
	binary.LittleEndian.PutUint32(b[4:8], f.SamplesPerSecond())
	binary.LittleEndian.PutUint32(b[8:12], f.AverageBytesPerSecond())
	binary.LittleEndian.PutUint16(b[12:14], f.BlockAlign())
	binary.LittleEndian.PutUint16(b[14:16], f.SignificantBitsPerSample())
	if len(ef) != 0 {
		binary.LittleEndian.PutUint16(b[16:18], uint16(len(ef)))
		copy(b[18:], ef)
	}

	return
}

func (f *ExtendedFormatChunk) Chunk() riffbin.Chunk {
	return &riffbin.OnMemorySubChunk{
		ID:      fmtBytes,
		Payload: f.Bytes(),
	}
}
