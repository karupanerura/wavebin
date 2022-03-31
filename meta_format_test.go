package wavebin_test

import (
	"testing"

	"github.com/karupanerura/wavebin"
)

func TestPCMMetaFormat(t *testing.T) {
	f := wavebin.NewPCMMetaFormat(wavebin.StereoChannels, 44100, 16)

	if f.AverageBytesPerSecond() != 176400 {
		t.Errorf("AverageBytesPerSecond should be 176400 but got: %d", f.AverageBytesPerSecond())
	}
	if f.BlockAlign() != 4 {
		t.Errorf("BlockAlign should be 4 but got: %d", f.BlockAlign())
	}
}
