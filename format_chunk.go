package wavebin

import (
	"encoding/binary"
	"errors"
	"io"

	"github.com/karupanerura/riffbin"
)

type FormatChunk interface {
	ChunkProvider
	MetaFormat
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

func (f *BasicFormatChunk) ReadFrom(r io.Reader) (int64, error) {
	var b [14]byte

	n, err := io.ReadFull(r, b[:])
	if err != nil {
		return int64(n), err
	}

	mf := &rawMetaFormat{
		compressionCode:       binary.LittleEndian.Uint16(b[:2]),
		channels:              binary.LittleEndian.Uint16(b[2:4]),
		samplesPerSecond:      binary.LittleEndian.Uint32(b[4:8]),
		averageBytesPerSecond: binary.LittleEndian.Uint32(b[8:12]),
		blockAlign:            binary.LittleEndian.Uint16(b[12:14]),
	}
	mf.significantBitsPerSample = 8 * mf.blockAlign / mf.channels
	f.MetaFormat = mf

	return 14, nil
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

func (f *ExtendedFormatChunk) ReadFrom(r io.Reader) (int64, error) {
	var b [16]byte

	n, err := io.ReadFull(r, b[:])
	if err != nil {
		return int64(n), err
	}

	mf := &rawMetaFormat{
		compressionCode:          binary.LittleEndian.Uint16(b[:2]),
		channels:                 binary.LittleEndian.Uint16(b[2:4]),
		samplesPerSecond:         binary.LittleEndian.Uint32(b[4:8]),
		averageBytesPerSecond:    binary.LittleEndian.Uint32(b[8:12]),
		blockAlign:               binary.LittleEndian.Uint16(b[12:14]),
		significantBitsPerSample: binary.LittleEndian.Uint16(b[14:16]),
	}
	f.MetaFormat = mf

	n, err = io.ReadFull(r, b[:2])
	if errors.Is(err, io.EOF) {
		return int64(len(b)), nil
	} else if err != nil {
		return int64(len(b) + n), nil
	}

	extraFieldSize := binary.LittleEndian.Uint16(b[:2])
	mf.extraField = make([]byte, extraFieldSize)
	n, err = io.ReadFull(r, mf.extraField)
	if err != nil {
		return int64(len(b) + 2 + n), err
	}

	return int64(len(b) + 2 + n), nil
}

func (f *ExtendedFormatChunk) Chunk() riffbin.Chunk {
	return &riffbin.OnMemorySubChunk{
		ID:      fmtBytes,
		Payload: f.Bytes(),
	}
}
