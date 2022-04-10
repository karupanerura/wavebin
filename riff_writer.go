package wavebin

import (
	"io"

	"github.com/karupanerura/riffbin"
)

var (
	waveBytes = [4]byte{'W', 'A', 'V', 'E'}
	listBytes = [4]byte{'L', 'I', 'S', 'T'}
	fmtBytes  = [4]byte{'f', 'm', 't', ' '}
	factBytes = [4]byte{'f', 'a', 'c', 't'}
	infoBytes = [4]byte{'I', 'N', 'F', 'O'}
	dataBytes = [4]byte{'d', 'a', 't', 'a'}
	junkBytes = [4]byte{'j', 'u', 'n', 'k'}
)

type ChunkProvider interface {
	Chunk() riffbin.Chunk
}

func CreateCompletedRIFF(format FormatChunk, samples []byte, extras ...ChunkProvider) *riffbin.RIFFChunk {
	return createRIFF(format, &riffbin.OnMemorySubChunk{
		ID:      dataBytes,
		Payload: samples,
	}, extras...)
}

func CreateIncompleteRIFF(format FormatChunk, r io.Reader, extras ...ChunkProvider) *riffbin.RIFFChunk {
	dataChunk := riffbin.NewIncompleteSubChunk(dataBytes, r)
	return createRIFF(format, dataChunk, extras...)
}

func createRIFF(format FormatChunk, dataChunk riffbin.Chunk, extras ...ChunkProvider) *riffbin.RIFFChunk {
	chunks := make([]riffbin.Chunk, 2+len(extras))
	chunks[0] = format.Chunk()
	for i, extra := range extras {
		chunks[i+1] = extra.Chunk()
	}
	chunks[len(chunks)-1] = dataChunk
	return &riffbin.RIFFChunk{
		FormType: waveBytes,
		Payload:  chunks,
	}
}
