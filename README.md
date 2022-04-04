# github.com/karupanerura/wavebin ![](https://github.com/karupanerura/wavebin/workflows/test/badge.svg?branch=main) [![Go Reference](https://pkg.go.dev/badge/github.com/karupanerura/wavebin.svg)](https://pkg.go.dev/github.com/karupanerura/wavebin) [![codecov.io](https://codecov.io/github/karupanerura/wavebin/coverage.svg?branch=main)](https://codecov.io/github/karupanerura/wavebin?branch=main)

Go module implementing for WAVE binary format.

# Features

* Construct WAVE data structure
* Write WAVE data structure
  * Can write WAVE data from io.Reader
* Parse WAVE binary to data structure

# Motivation

There has never been a library that can be used without pre-determining the binary size in Go.

# Examples

## Example1: write WAVE file

```go
_, err := riffbin.NewCompletedChunkWriter(w).Write(
	wavebin.CreateCompletedRIFF(
		&wavebin.ExtendedFormatChunk{
			MetaFormat: wavebin.NewPCMMetaFormat(wavebin.MonoralChannels, 44100, 8),
		},
		[]byte{...},
	),
)
```

## Example2: write WAVE file from io.Reader

```go
w, err := riffbin.NewIncompleteChunkWriter(f)
if err != nil {
	panic(err)
}

_, err = w.Write(
	wavebin.CreateIncompleteRIFF(
		&wavebin.ExtendedFormatChunk{
			MetaFormat: wavebin.NewPCMMetaFormat(wavebin.MonoralChannels, 44100, 8),
		},
		r,
	),
)
```

## Example3: write WAVE file to io.WriteCloser

```go
w, err := wavebin.CreateSampleWriter(f, &wavebin.ExtendedFormatChunk{
		MetaFormat: wavebin.NewPCMMetaFormat(wavebin.MonoralChannels, 44100, 8),
})
if err != nil {
	panic(err)
}

// write samples
max := 2000 * int(math.Ceil(44100/440))
for i := 0; i < max; i++ {
	b := byte(math.Floor((float64(i) / float64(max)) * 192.0 * (1.0 + math.Sin(2.0*math.Pi*float64(i)/(44100/440))) / 2.0))
	_, err = w.Write([]byte{b})
	if err != nil {
		panic(err)
	}
}

err = w.Close()
if err != nil {
	panic(err)
}
```

## Example4: write WAVE file to io.WriteCloser with PCMWriter

```go
w, err := wavebin.CreateSampleWriter(f, &wavebin.ExtendedFormatChunk{
		MetaFormat: wavebin.NewPCMMetaFormat(wavebin.StereoChannels, 44100, 16),
})
if err != nil {
	panic(err)
}

pcmWriter := &wavebin.PCMWriter{W: bufio.NewWriter(w)} // bufio for performance

// write samples
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

err = w.Close()
if err != nil {
	panic(err)
}
```
