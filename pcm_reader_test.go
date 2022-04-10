package wavebin_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/karupanerura/wavebin"
)

func TestPCMReader(t *testing.T) {
	t.Run("8BitMonoral", func(t *testing.T) {
		for _, tt := range []struct {
			name     string
			rawBytes []byte
			expected []wavebin.PCM8BitMonoralSample
		}{
			{"Empty", []byte{}, nil},
			{"8BitMonoral1Sample", []byte{0x00}, []wavebin.PCM8BitMonoralSample{0}},
			{"8BitMonoral2Samples", []byte{0x00, 0xFF}, []wavebin.PCM8BitMonoralSample{0, 255}},
			{"8BitMonoral3Samples", []byte{0x00, 0xFF, 0xFF}, []wavebin.PCM8BitMonoralSample{0, 255, 255}},
			{"8BitMonoral4Samples", []byte{0x00, 0x00, 0xFF, 0xFF}, []wavebin.PCM8BitMonoralSample{0, 0, 255, 255}},
			{"8BitMonoral5Samples", []byte{0x00, 0x00, 0xFF, 0xFF, 0x00}, []wavebin.PCM8BitMonoralSample{0, 0, 255, 255, 0}},
			{"8BitMonoral6Samples", []byte{0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00}, []wavebin.PCM8BitMonoralSample{0, 0, 255, 255, 0, 0}},
		} {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				r := wavebin.NewPCMReader[wavebin.PCM8BitMonoralSample](bytes.NewReader(tt.rawBytes), wavebin.PCM8BitMonoralSampleParser{})

				var samples []wavebin.PCM8BitMonoralSample
				for {
					sample, err := r.ReadSample()
					if errors.Is(err, io.EOF) {
						break
					} else if err != nil {
						t.Fatal(err)
					}

					samples = append(samples, sample)
				}

				if df := cmp.Diff(tt.expected, samples); df != "" {
					t.Error(df)
				}
			})
		}
	})
	t.Run("8BitStereo", func(t *testing.T) {
		for _, tt := range []struct {
			name     string
			rawBytes []byte
			expected []wavebin.PCM8BitStereoSample
		}{
			{"Empty", []byte{}, nil},
			{"8BitStereo1Sample", []byte{0x00, 0x01}, []wavebin.PCM8BitStereoSample{{L: 0, R: 1}}},
			{"8BitStereo2Sample", []byte{0x00, 0x01, 0x01, 0x00}, []wavebin.PCM8BitStereoSample{{L: 0, R: 1}, {L: 1, R: 0}}},
			{"8BitStereo3Sample", []byte{0x00, 0x01, 0x00, 0x00, 0x01, 0x00}, []wavebin.PCM8BitStereoSample{{L: 0, R: 1}, {L: 0, R: 0}, {L: 1, R: 0}}},
		} {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				r := wavebin.NewPCMReader[wavebin.PCM8BitStereoSample](bytes.NewReader(tt.rawBytes), wavebin.PCM8BitStereoSampleParser{})

				var samples []wavebin.PCM8BitStereoSample
				for {
					sample, err := r.ReadSample()
					if errors.Is(err, io.EOF) {
						break
					} else if err != nil {
						t.Fatal(err)
					}

					samples = append(samples, sample)
				}

				if df := cmp.Diff(tt.expected, samples); df != "" {
					t.Error(df)
				}
			})
		}

		t.Run("InvalidBytes", func(t *testing.T) {
			r := wavebin.NewPCMReader[wavebin.PCM8BitStereoSample](bytes.NewReader([]byte{0x00}), wavebin.PCM8BitStereoSampleParser{})

			_, err := r.ReadSample()
			if !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	})
	t.Run("16BitMonoral", func(t *testing.T) {
		for _, tt := range []struct {
			name     string
			rawBytes []byte
			expected []wavebin.PCM16BitMonoralSample
		}{
			{"Empty", []byte{}, nil},
			{"16BitMonoral1Sample", []byte{0x00, 0x00}, []wavebin.PCM16BitMonoralSample{0}},
			{"16BitMonoral2Samples", []byte{0xff, 0xff, 0x01, 0x00}, []wavebin.PCM16BitMonoralSample{-1, 1}},
			{"16BitMonoral3Samples", []byte{0xff, 0xff, 0x01, 0x00, 0xff, 0xff}, []wavebin.PCM16BitMonoralSample{-1, 1, -1}},
		} {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				r := wavebin.NewPCMReader[wavebin.PCM16BitMonoralSample](bytes.NewReader(tt.rawBytes), wavebin.PCM16BitMonoralSampleParser{})

				var samples []wavebin.PCM16BitMonoralSample
				for {
					sample, err := r.ReadSample()
					if errors.Is(err, io.EOF) {
						break
					} else if err != nil {
						t.Fatal(err)
					}

					samples = append(samples, sample)
				}

				if df := cmp.Diff(tt.expected, samples); df != "" {
					t.Error(df)
				}
			})
		}

		t.Run("InvalidBytes", func(t *testing.T) {
			r := wavebin.NewPCMReader[wavebin.PCM16BitMonoralSample](bytes.NewReader([]byte{0x00}), wavebin.PCM16BitMonoralSampleParser{})

			_, err := r.ReadSample()
			if !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	})
	t.Run("16BitStereo", func(t *testing.T) {
		for _, tt := range []struct {
			name     string
			rawBytes []byte
			expected []wavebin.PCM16BitStereoSample
		}{
			{"Empty", []byte{}, nil},
			{"16BitStereo1Sample", []byte{0xff, 0xff, 0x01, 0x00}, []wavebin.PCM16BitStereoSample{{L: -1, R: 1}}},
			{"16BitStereo2Sample", []byte{0xff, 0xff, 0x01, 0x00, 0x02, 0x00, 0xfe, 0xff}, []wavebin.PCM16BitStereoSample{{L: -1, R: 1}, {L: 2, R: -2}}},
		} {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				r := wavebin.NewPCMReader[wavebin.PCM16BitStereoSample](bytes.NewReader(tt.rawBytes), wavebin.PCM16BitStereoSampleParser{})

				var samples []wavebin.PCM16BitStereoSample
				for {
					sample, err := r.ReadSample()
					if errors.Is(err, io.EOF) {
						break
					} else if err != nil {
						t.Fatal(err)
					}

					samples = append(samples, sample)
				}

				if df := cmp.Diff(tt.expected, samples); df != "" {
					t.Error(df)
				}
			})
		}

		t.Run("InvalidBytes", func(t *testing.T) {
			r := wavebin.NewPCMReader[wavebin.PCM16BitStereoSample](bytes.NewReader([]byte{0x00, 0x00, 0x00}), wavebin.PCM16BitStereoSampleParser{})

			_, err := r.ReadSample()
			if !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Errorf("unexpected err: %v", err)
			}
		})
	})
}
