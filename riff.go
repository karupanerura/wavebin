package wavebin

import (
	"io"

	"github.com/karupanerura/riffbin"
)

var (
	waveBytes = [4]byte{'W', 'A', 'V', 'E'}
	fmtBytes  = [4]byte{'f', 'm', 't', ' '}
	factBytes = [4]byte{'f', 'a', 'c', 't'}
	infoBytes = [4]byte{'I', 'N', 'F', 'O'}
	dataBytes = [4]byte{'d', 'a', 't', 'a'}
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

type pipeWriter struct {
	pw    *io.PipeWriter
	errCh chan error
}

func (w *pipeWriter) Write(data []byte) (int, error) {
	return w.pw.Write(data)
}

func (w *pipeWriter) Close() error {
	_ = w.pw.Close() // always be nil
	return <-w.errCh
}

func CreateSampleWriter(w io.WriteSeeker, format FormatChunk, extras ...ChunkProvider) (io.WriteCloser, error) {
	cw, err := riffbin.NewIncompleteChunkWriter(w)
	if err != nil {
		return nil, err
	}

	pr, pw := io.Pipe()
	dataChunk := riffbin.NewIncompleteSubChunk(dataBytes, pr)
	riffChunk := createRIFF(format, dataChunk, extras...)

	errCh := make(chan error, 1)
	go func() {
		_, err := cw.Write(riffChunk)
		_ = pr.CloseWithError(err) // always be nil
		errCh <- err
		close(errCh)
	}()

	return &pipeWriter{pw: pw, errCh: errCh}, nil
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
