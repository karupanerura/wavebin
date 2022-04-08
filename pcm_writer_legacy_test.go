//go:build !go1.18
// +build !go1.18

package wavebin_test

import (
	"bufio"
	"bytes"
	"math"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/karupanerura/wavebin"
)

func TestPCMWriter(t *testing.T) {
	for _, tt := range []struct {
		name     string
		samples  []wavebin.PCMSample
		expected []byte
	}{
		{"Empty", []wavebin.PCMSample{}, nil},
		{"8BitMonoral1Sample", []wavebin.PCMSample{wavebin.PCM8BitMonoralSample(0)}, []byte{0x00}},
		{"8BitMonoral2Samples", []wavebin.PCMSample{wavebin.PCM8BitMonoralSample(0), wavebin.PCM8BitMonoralSample(255)}, []byte{0x00, 0xFF}},
		{"8BitMonoral3Samples", []wavebin.PCMSample{wavebin.PCM8BitMonoralSample(0), wavebin.PCM8BitMonoralSample(255), wavebin.PCM8BitMonoralSample(255)}, []byte{0x00, 0xFF, 0xFF}},
		{"8BitMonoral4Samples", []wavebin.PCMSample{wavebin.PCM8BitMonoralSample(0), wavebin.PCM8BitMonoralSample(0), wavebin.PCM8BitMonoralSample(255), wavebin.PCM8BitMonoralSample(255)}, []byte{0x00, 0x00, 0xFF, 0xFF}},
		{"8BitMonoral5Samples", []wavebin.PCMSample{wavebin.PCM8BitMonoralSample(0), wavebin.PCM8BitMonoralSample(0), wavebin.PCM8BitMonoralSample(255), wavebin.PCM8BitMonoralSample(255), wavebin.PCM8BitMonoralSample(0)}, []byte{0x00, 0x00, 0xFF, 0xFF, 0x00}},
		{"8BitMonoral6Samples", []wavebin.PCMSample{wavebin.PCM8BitMonoralSample(0), wavebin.PCM8BitMonoralSample(0), wavebin.PCM8BitMonoralSample(255), wavebin.PCM8BitMonoralSample(255), wavebin.PCM8BitMonoralSample(0), wavebin.PCM8BitMonoralSample(0)}, []byte{0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00}},
		{"8BitStereo1Sample", []wavebin.PCMSample{wavebin.PCM8BitStereoSample{L: 0, R: 1}}, []byte{0x00, 0x01}},
		{"8BitStereo2Sample", []wavebin.PCMSample{wavebin.PCM8BitStereoSample{L: 0, R: 1}, wavebin.PCM8BitStereoSample{L: 1, R: 0}}, []byte{0x00, 0x01, 0x01, 0x00}},
		{"8BitStereo3Sample", []wavebin.PCMSample{wavebin.PCM8BitStereoSample{L: 0, R: 1}, wavebin.PCM8BitStereoSample{L: 0, R: 0}, wavebin.PCM8BitStereoSample{L: 1, R: 0}}, []byte{0x00, 0x01, 0x00, 0x00, 0x01, 0x00}},
		{"16BitMonoral1Sample", []wavebin.PCMSample{wavebin.PCM16BitMonoralSample(0)}, []byte{0x00, 0x00}},
		{"16BitMonoral2Samples", []wavebin.PCMSample{wavebin.PCM16BitMonoralSample(-1), wavebin.PCM16BitMonoralSample(1)}, []byte{0xff, 0xff, 0x01, 0x00}},
		{"16BitMonoral3Samples", []wavebin.PCMSample{wavebin.PCM16BitMonoralSample(-1), wavebin.PCM16BitMonoralSample(1), wavebin.PCM16BitMonoralSample(-1)}, []byte{0xff, 0xff, 0x01, 0x00, 0xff, 0xff}},
		{"16BitStereo1Sample", []wavebin.PCMSample{wavebin.PCM16BitStereoSample{L: -1, R: 1}}, []byte{0xff, 0xff, 0x01, 0x00}},
		{"16BitStereo2Sample", []wavebin.PCMSample{wavebin.PCM16BitStereoSample{L: -1, R: 1}, wavebin.PCM16BitStereoSample{L: 2, R: -2}}, []byte{0xff, 0xff, 0x01, 0x00, 0x02, 0x00, 0xfe, 0xff}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := &wavebin.PCMWriter{W: &buf}

			_, err := w.WriteSamples(tt.samples...)
			if err != nil {
				t.Fatal(err)
			}

			if df := cmp.Diff(tt.expected, buf.Bytes()); df != "" {
				t.Error(df)
			}
		})
	}
}

func ExamplePCMWriter() {
	f, err := os.CreateTemp("", "riffbin")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())

	w, err := wavebin.CreateSampleWriter(f, &wavebin.ExtendedFormatChunk{
		MetaFormat: wavebin.NewPCMMetaFormat(wavebin.StereoChannels, 44100, 16),
	})
	if err != nil {
		panic(err)
	}

	// write samples
	{
		bw := bufio.NewWriter(w) // bufio for performance
		pcmWriter := &wavebin.PCMWriter{W: bw}
		max := 2000 * int(math.Ceil(44100/440))
		for i := 0; i < max; i++ {
			sample := wavebin.PCM16BitStereoSample{
				L: int16(math.Floor((float64(i) / float64(max)) * 32767 * math.Sin(2.0*math.Pi*float64(i)/(44100/440)))),
				R: int16(math.Floor((float64(i) / float64(max)) * 32767 * math.Sin(2.0*math.Pi*float64(i)/(44100/110)))),
			}
			_, err = pcmWriter.WriteSamples(sample)
			if err != nil {
				panic(err)
			}
		}

		err := bw.Flush()
		if err != nil {
			panic(err)
		}
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}

	// preview on macOS
	if os.Getenv("DEBUG_TEST_PLAY") != "" && runtime.GOOS == "darwin" {
		err = exec.Command("afplay", f.Name()).Run()
		if err != nil {
			panic(err)
		}

		err = exec.Command("cp", f.Name(), "/Users/karupanerura/a.wav").Run()
		if err != nil {
			panic(err)
		}
	}

	// Output:
}
