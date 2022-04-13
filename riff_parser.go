package wavebin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/karupanerura/riffbin"
)

var (
	ErrUnexpectedFormType   = errors.New("unexpected form type")
	ErrUnexpectedChunkType  = errors.New("unexpected chunk type")
	ErrUnexpectedChunkSize  = errors.New("unexpected chunk size")
	ErrUnexpectedBlockAlign = errors.New("unexpected block align")
	ErrUnknownChunk         = errors.New("unknown chunk")
	ErrUnknownListType      = errors.New("unknown list type")
	ErrLackOfRequiredChunks = errors.New("lack of required chunks")
)

func ParseWaveRIFF(riffChunk *riffbin.RIFFChunk, ignoreUnknownChunk bool) (fmtChunk FormatChunk, infoChunk *InfoChunk, factChunk *FactChunk, sampleReader riffbin.SubChunk, err error) {
	if riffChunk.FormType != waveBytes {
		err = fmt.Errorf("%w: %s", ErrUnexpectedFormType, string(riffChunk.FormType[:]))
		return
	}

	for _, chunk := range riffChunk.Payload {
		if bytes.Equal(chunk.ChunkID(), junkBytes[:]) {
			continue
		}

		if bytes.Equal(chunk.ChunkID(), fmtBytes[:]) {
			fmtChunk, err = parseFormatChunk(chunk)
			if err != nil {
				return
			}
		} else if bytes.Equal(chunk.ChunkID(), dataBytes[:]) {
			sampleReader, err = parseDataChunk(chunk)
			if err != nil {
				return
			}
		} else if bytes.Equal(chunk.ChunkID(), listBytes[:]) {
			listChunk, ok := chunk.(*riffbin.ListChunk)
			if !ok {
				err = fmt.Errorf("RIFF[WAVE].LIST: %w", ErrUnexpectedChunkType)
				return
			}

			if listChunk.ListType == infoBytes {
				infoChunk, err = parseInfoChunk(listChunk)
				if err != nil {
					return
				}
			} else {
				// unknown chunk
				if ignoreUnknownChunk {
					continue
				}

				err = fmt.Errorf("RIFF[WAVE].LIST[%s]: %w", string(listChunk.ListType[:]), ErrUnknownListType)
				return
			}
		} else if bytes.Equal(chunk.ChunkID(), factBytes[:]) {
			factChunk, err = parseFactChunk(chunk)
			if err != nil {
				return
			}
		} else {
			// unknown chunk
			if ignoreUnknownChunk {
				continue
			}

			err = fmt.Errorf("RIFF[WAVE].%s: %w", string(chunk.ChunkID()), ErrUnknownChunk)
			return
		}
	}
	if fmtChunk == nil || sampleReader == nil {
		err = ErrLackOfRequiredChunks
		return
	}

	return
}

func parseFormatChunk(chunk riffbin.Chunk) (FormatChunk, error) {
	subChunk, ok := chunk.(riffbin.SubChunk)
	if !ok {
		return nil, fmt.Errorf("RIFF[WAVE].fmt: %w", ErrUnexpectedChunkType)
	}

	if subChunk.BodySize() == 14 {
		fmtChunk := &BasicFormatChunk{}
		_, err := fmtChunk.ReadFrom(subChunk)
		if err != nil {
			return nil, fmt.Errorf("RIFF[WAVE].fmt: %w", err)
		}

		return fmtChunk, nil
	} else {
		fmtChunk := &ExtendedFormatChunk{}
		_, err := fmtChunk.ReadFrom(subChunk)
		if err != nil {
			return nil, fmt.Errorf("RIFF[WAVE].fmt: %w", err)
		}

		return fmtChunk, nil
	}
}

func parseDataChunk(chunk riffbin.Chunk) (riffbin.SubChunk, error) {
	subChunk, ok := chunk.(riffbin.SubChunk)
	if !ok {
		return nil, fmt.Errorf("RIFF[WAVE].data: %w", ErrUnexpectedChunkType)
	}

	return subChunk, nil
}

func parseFactChunk(chunk riffbin.Chunk) (*FactChunk, error) {
	subChunk, ok := chunk.(riffbin.SubChunk)
	if !ok {
		return nil, fmt.Errorf("RIFF[WAVE].fact: %w", ErrUnexpectedChunkType)
	}
	if subChunk.BodySize() != 4 {
		return nil, fmt.Errorf("RIFF[WAVE].fact: %w", ErrUnexpectedChunkSize)
	}

	var b [4]byte
	_, err := io.ReadFull(subChunk, b[:])
	if err != nil {
		return nil, fmt.Errorf("RIFF[WAVE].fact: %w", err)
	}

	return &FactChunk{SampleLength: SampleLength(binary.LittleEndian.Uint32(b[:]))}, nil
}

func parseInfoChunk(listChunk *riffbin.ListChunk) (*InfoChunk, error) {
	infoChunk := &InfoChunk{Data: map[InfoKey]string{}}
	for _, chunk := range listChunk.Payload {
		subChunk, ok := chunk.(riffbin.SubChunk)
		if !ok {
			return nil, fmt.Errorf("RIFF[WAVE].INFO.%s: %w", string(chunk.ChunkID()), ErrUnexpectedChunkType)
		}

		var s strings.Builder
		_, err := io.Copy(&s, subChunk)
		if err != nil {
			return nil, fmt.Errorf("RIFF[WAVE].INFO.%s: %w", string(chunk.ChunkID()), err)
		}

		var key InfoKey
		copy(key[:], subChunk.ChunkID())
		infoChunk.Data[key] = s.String()
	}

	return infoChunk, nil
}
