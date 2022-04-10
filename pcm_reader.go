package wavebin

import "io"

type PCMReader[T PCMSample] struct {
	r io.Reader
	p PCMSampleParser[T]
}

func NewPCMReader[T PCMSample](r io.Reader, p PCMSampleParser[T]) *PCMReader[T] {
	return &PCMReader[T]{r: r, p: p}
}

func (r *PCMReader[T]) ReadSample() (T, error) {
	return r.p.ParseFromReader(r.r)
}
