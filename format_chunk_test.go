package wavebin_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/karupanerura/wavebin"
)

func TestBasicFormatChunk(t *testing.T) {
	f := wavebin.BasicFormatChunk{
		MetaFormat: wavebin.NewPCMMetaFormat(wavebin.MonoralChannels, 44100, 8),
	}

	expected := []byte{
		0x01, 0x00, // Compression Code (Linear PCM)
		0x01, 0x00, // Number of channels (Monoral)
		0x44, 0xAC, 0x00, 0x00, // Sample rate (44.1Hz)
		0x44, 0xAC, 0x00, 0x00, // Average bytes per second (44.1Hz/Monoral)
		0x01, 0x00, // Block align (8bit/Monoral)
	}
	if !bytes.Equal(f.Bytes(), expected) {
		t.Errorf("Invalid Bytes: %v", f.Bytes())
		t.Log(hex.Dump(f.Bytes()))
		t.Log(hex.Dump(expected))
	}
}

func TestExtendedFormatChunk(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		f := wavebin.ExtendedFormatChunk{
			MetaFormat: wavebin.NewPCMMetaFormat(wavebin.MonoralChannels, 44100, 8),
		}

		expected := []byte{
			0x01, 0x00, // Compression Code (Linear PCM)
			0x01, 0x00, // Number of channels (Monoral)
			0x44, 0xAC, 0x00, 0x00, // Sample rate (44.1Hz)
			0x44, 0xAC, 0x00, 0x00, // Average bytes per second (44.1Hz/Monoral)
			0x01, 0x00, // Block align (8bit/Monoral)
			0x08, 0x00, // Significant bits per sample (8bit)
		}
		if !bytes.Equal(f.Bytes(), expected) {
			t.Errorf("Invalid Bytes: %v", f.Bytes())
			t.Log(hex.Dump(f.Bytes()))
			t.Log(hex.Dump(expected))
		}
	})
	t.Run("WithExtraField", func(t *testing.T) {
		f := wavebin.ExtendedFormatChunk{
			MetaFormat: wavebin.NewExtensibleMetaFormat(wavebin.MonoralChannels, 44100, 8, 8, wavebin.ChannelMaskFrontCenter, [16]byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
			}),
		}

		expected := []byte{
			0xFE, 0xFF, // Compression Code (Extensible)
			0x01, 0x00, // Number of channels (Monoral)
			0x44, 0xAC, 0x00, 0x00, // Sample rate (44.1Hz)
			0x44, 0xAC, 0x00, 0x00, // Average bytes per second (44.1Hz/Monoral)
			0x01, 0x00, // Block align (8bit/Monoral)
			0x08, 0x00, // Significant bits per sample (8bit)
			0x16, 0x00, // extra field size
			0x08, 0x00, // ValidBitsPerSample
			0x04, 0x00, 0x00, 0x00, // ChannelMask (FrontCenter)
			0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, // SubFormatChunk
			0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		}
		if !bytes.Equal(f.Bytes(), expected) {
			t.Errorf("Invalid Bytes: %v", f.Bytes())
			t.Log(hex.Dump(f.Bytes()))
			t.Log(hex.Dump(expected))
		}
	})
}
