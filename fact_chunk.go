package wavebin

import (
	"encoding/binary"

	"github.com/karupanerura/riffbin"
)

type SampleLength uint32

type FactChunk struct {
	SampleLength
}

func (f *FactChunk) Chunk() riffbin.Chunk {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(f.SampleLength))
	return &riffbin.CompletedSubChunk{
		ID:      factBytes,
		Payload: b,
	}
}
