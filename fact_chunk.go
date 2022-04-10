package wavebin

import (
	"encoding/binary"
	"io"

	"github.com/karupanerura/riffbin"
)

type SampleLength uint32

type FactChunk struct {
	SampleLength
}

func (f *FactChunk) Chunk() riffbin.Chunk {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(f.SampleLength))
	return &riffbin.OnMemorySubChunk{
		ID:      factBytes,
		Payload: b,
	}
}

func (f *FactChunk) ReadFrom(r io.Reader) (int64, error) {
	var b [4]byte

	n, err := io.ReadFull(r, b[:])
	if err != nil {
		return int64(n), err
	}

	f.SampleLength = SampleLength(binary.LittleEndian.Uint32(b[:]))
	return 4, nil
}
