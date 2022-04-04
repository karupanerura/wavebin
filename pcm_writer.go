package wavebin

import (
	"encoding/binary"
	"io"
)

type PCMWriter struct {
	W io.Writer
}

func (w *PCMWriter) WriteSamples(samples ...PCMSample) (n int64, err error) {
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

type PCMSample interface {
	PutSamples([]byte)
	ByteSize() int
}

type PCM8BitMonoralSample uint8

func (s PCM8BitMonoralSample) PutSamples(p []byte) {
	p[0] = uint8(s)
}

func (s PCM8BitMonoralSample) ByteSize() int {
	return 1
}

type PCM8BitStereoSample struct{ L, R uint8 }

func (s PCM8BitStereoSample) PutSamples(p []byte) {
	_ = p[1] // early bounds check to guarantee safety of writes below
	p[0] = uint8(s.L)
	p[1] = uint8(s.R)
}

func (s PCM8BitStereoSample) ByteSize() int {
	return 2
}

type PCM16BitMonoralSample int16

func (s PCM16BitMonoralSample) PutSamples(p []byte) {
	_ = p[1] // early bounds check to guarantee safety of writes below
	binary.LittleEndian.PutUint16(p, encode16bitSignedInt(int16(s)))
}

func (s PCM16BitMonoralSample) ByteSize() int {
	return 2
}

type PCM16BitStereoSample struct{ L, R int16 }

func (s PCM16BitStereoSample) PutSamples(p []byte) {
	_ = p[3] // early bounds check to guarantee safety of writes below
	binary.LittleEndian.PutUint16(p[0:2], encode16bitSignedInt(int16(s.L)))
	binary.LittleEndian.PutUint16(p[2:4], encode16bitSignedInt(int16(s.R)))
}

func (s PCM16BitStereoSample) ByteSize() int {
	return 4
}

func encode16bitSignedInt(s16 int16) (u16 uint16) {
	negative := s16 < 0
	if negative {
		s16 &= 32767
	}

	u16 = uint16(s16)
	if negative {
		u16 |= 1 << 15
	}
	return
}
