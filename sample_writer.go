package wavebin

import (
	"io"

	"github.com/karupanerura/riffbin"
)

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
