package wavebin

import (
	"encoding/binary"
	"io"
	"math"
)

type PCMSample interface {
	PutSamples([]byte)
	ByteSize() int
}

type PCMSampleParser[T PCMSample] interface {
	ParseFromReader(r io.Reader) (T, error)
}

type PCM8BitMonoralSample uint8

func (s PCM8BitMonoralSample) PutSamples(p []byte) {
	p[0] = uint8(s)
}

func (s PCM8BitMonoralSample) ByteSize() int {
	return 1
}

type PCM8BitMonoralSampleParser struct{}

var _ PCMSampleParser[PCM8BitMonoralSample] = PCM8BitMonoralSampleParser{}

func (PCM8BitMonoralSampleParser) ParseFromReader(r io.Reader) (PCM8BitMonoralSample, error) {
	var b [1]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return 0, err
	}

	return PCM8BitMonoralSample(b[0]), nil
}

type PCM8BitStereoSample struct{ L, R uint8 }

func (s PCM8BitStereoSample) PutSamples(p []byte) {
	_ = p[1] // early bounds check to guarantee safety of writes below
	p[0] = s.L
	p[1] = s.R
}

func (s PCM8BitStereoSample) ByteSize() int {
	return 2
}

type PCM8BitStereoSampleParser struct{}

var _ PCMSampleParser[PCM8BitStereoSample] = PCM8BitStereoSampleParser{}

func (PCM8BitStereoSampleParser) ParseFromReader(r io.Reader) (PCM8BitStereoSample, error) {
	var b [2]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return PCM8BitStereoSample{}, err
	}

	return PCM8BitStereoSample{L: b[0], R: b[1]}, nil
}

type PCM16BitMonoralSample int16

func (s PCM16BitMonoralSample) PutSamples(p []byte) {
	_ = p[1] // early bounds check to guarantee safety of writes below
	binary.LittleEndian.PutUint16(p, encode16bitSignedInt(int16(s)))
}

func (s PCM16BitMonoralSample) ByteSize() int {
	return 2
}

type PCM16BitMonoralSampleParser struct{}

var _ PCMSampleParser[PCM16BitMonoralSample] = PCM16BitMonoralSampleParser{}

func (PCM16BitMonoralSampleParser) ParseFromReader(r io.Reader) (PCM16BitMonoralSample, error) {
	var b [2]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return 0, err
	}

	return PCM16BitMonoralSample(decode16bitSignedInt(binary.LittleEndian.Uint16(b[:]))), nil
}

type PCM16BitStereoSample struct{ L, R int16 }

func (s PCM16BitStereoSample) PutSamples(p []byte) {
	_ = p[3] // early bounds check to guarantee safety of writes below
	binary.LittleEndian.PutUint16(p[0:2], encode16bitSignedInt(s.L))
	binary.LittleEndian.PutUint16(p[2:4], encode16bitSignedInt(s.R))
}

func (s PCM16BitStereoSample) ByteSize() int {
	return 4
}

type PCM16BitStereoSampleParser struct{}

var _ PCMSampleParser[PCM16BitStereoSample] = PCM16BitStereoSampleParser{}

func (PCM16BitStereoSampleParser) ParseFromReader(r io.Reader) (PCM16BitStereoSample, error) {
	var b [4]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return PCM16BitStereoSample{}, err
	}

	return PCM16BitStereoSample{
		L: decode16bitSignedInt(binary.LittleEndian.Uint16(b[0:2])),
		R: decode16bitSignedInt(binary.LittleEndian.Uint16(b[2:4])),
	}, nil
}

func encode16bitSignedInt(s16 int16) (u16 uint16) {
	negative := s16 < 0
	if negative {
		s16 &= 32767 // ^uint16(1<<15)
	}

	u16 = uint16(s16)
	if negative {
		u16 |= 1 << 15
	}
	return
}

func decode16bitSignedInt(u16 uint16) (s16 int16) {
	if negative := (u16 & (1 << 15)) != 0; negative {
		s16 = math.MinInt16 // to set negative flag
	}
	s16 |= int16(u16 & 32767) // ^uint16(1<<15)
	return
}
