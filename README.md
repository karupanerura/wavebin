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
