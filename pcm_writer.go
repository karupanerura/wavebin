package wavebin

import (
	"io"
)

type PCMWriter[T PCMSample] struct {
	W io.Writer
}

func (w *PCMWriter[T]) WriteSamples(samples ...T) (n int64, err error) {
	var buf [4]byte
	var bufLen int
	for _, sample := range samples {
		if len(buf)-bufLen < sample.ByteSize() {
			// flush
			var nn int
			nn, err = w.W.Write(buf[:bufLen])
			n += int64(nn)
			if err != nil {
				return
			}

			bufLen = 0
		}

		sample.PutSamples(buf[bufLen:])
		bufLen += sample.ByteSize()
	}
	if bufLen > 0 {
		// flush
		var nn int
		nn, err = w.W.Write(buf[:bufLen])
		n += int64(nn)
	}

	return
}
