package wavebin

import "encoding/binary"

type CompressionCode uint16

const (
	pcmCompressionCode        CompressionCode = 0x0001
	ieeeFloatCompressionCode  CompressionCode = 0x0003
	extensibleCompressionCode CompressionCode = 0xFFFE
)

type Channels uint16

const (
	MonoralChannels      Channels = 0x0001
	StereoChannels       Channels = 0x0002
	QuadraphonicChannels Channels = 0x0004
	SurroundChannels     Channels = 0x0006
)

type SamplesPerSecond uint32

type SignificantBitsPerSample uint16

type MetaFormat interface {
	CompressionCode() uint16
	Channels() uint16
	SamplesPerSecond() uint32
	SignificantBitsPerSample() uint16
	AverageBytesPerSecond() uint32
	BlockAlign() uint16
	ExtraField() []byte
}

type rawMetaFormat struct {
	channels                 Channels
	samplesPerSecond         SamplesPerSecond
	significantBitsPerSample SignificantBitsPerSample
}

func (f *rawMetaFormat) Channels() uint16 {
	return uint16(f.channels)
}

func (f *rawMetaFormat) SamplesPerSecond() uint32 {
	return uint32(f.samplesPerSecond)
}

func (f *rawMetaFormat) SignificantBitsPerSample() uint16 {
	return uint16(f.significantBitsPerSample)
}

func (f *rawMetaFormat) AverageBytesPerSecond() uint32 {
	return f.SamplesPerSecond() * uint32(f.BlockAlign())
}

func (f *rawMetaFormat) BlockAlign() uint16 {
	return f.Channels() * f.SignificantBitsPerSample() / 8
}

func (f *rawMetaFormat) ExtraField() []byte {
	return nil
}

type PCMMetaFormat struct {
	rawMetaFormat
}

func NewPCMMetaFormat(channels Channels, samplesPerSecond SamplesPerSecond, significantBitsPerSample SignificantBitsPerSample) *PCMMetaFormat {
	return &PCMMetaFormat{
		rawMetaFormat: rawMetaFormat{
			channels:                 channels,
			samplesPerSecond:         samplesPerSecond,
			significantBitsPerSample: significantBitsPerSample,
		},
	}
}

func (f *PCMMetaFormat) CompressionCode() uint16 {
	return uint16(pcmCompressionCode)
}

type IEEEFloatMetaFormat struct {
	rawMetaFormat
}

func NewIEEEFloatMetaFormat(channels Channels, samplesPerSecond SamplesPerSecond, significantBitsPerSample SignificantBitsPerSample) *IEEEFloatMetaFormat {
	return &IEEEFloatMetaFormat{
		rawMetaFormat: rawMetaFormat{
			channels:                 channels,
			samplesPerSecond:         samplesPerSecond,
			significantBitsPerSample: significantBitsPerSample,
		},
	}
}

func (f *IEEEFloatMetaFormat) CompressionCode() uint16 {
	return uint16(ieeeFloatCompressionCode)
}

type ValidBitsPerSample uint16

type ChannelMask uint32

const (
	ChannelMaskFrontLeft          ChannelMask = 0x00000001
	ChannelMaskFrontRight         ChannelMask = 0x00000002
	ChannelMaskFrontCenter        ChannelMask = 0x00000004
	ChannelMaskLowFrequency       ChannelMask = 0x00000008
	ChannelMaskBackLeft           ChannelMask = 0x00000010
	ChannelMaskBackRight          ChannelMask = 0x00000020
	ChannelMaskFrontLeftOfCenter  ChannelMask = 0x00000040
	ChannelMaskFrontRightOfCenter ChannelMask = 0x00000080
	ChannelMaskBackCenter         ChannelMask = 0x00000100
	ChannelMaskSideLeft           ChannelMask = 0x00000200
	ChannelMaskSideRight          ChannelMask = 0x00000400
	ChannelMaskTopCenter          ChannelMask = 0x00000800
	ChannelMaskTopFrontLeft       ChannelMask = 0x00001000
	ChannelMaskTopFrontCenter     ChannelMask = 0x00002000
	ChannelMaskTopFrontRight      ChannelMask = 0x00004000
	ChannelMaskTopBackLeft        ChannelMask = 0x00008000
	ChannelMaskTopBackCenter      ChannelMask = 0x00010000
	ChannelMaskTopBackRight       ChannelMask = 0x00020000
)

type ExtensibleMetaFormat struct {
	rawMetaFormat
	validBitsPerSample ValidBitsPerSample
	channelMask        ChannelMask
	subFormat          [16]byte
}

func NewExtensibleMetaFormat(channels Channels, samplesPerSecond SamplesPerSecond, significantBitsPerSample SignificantBitsPerSample, validBitsPerSample ValidBitsPerSample, channelMask ChannelMask, subFormat [16]byte) *ExtensibleMetaFormat {
	return &ExtensibleMetaFormat{
		rawMetaFormat: rawMetaFormat{
			channels:                 channels,
			samplesPerSecond:         samplesPerSecond,
			significantBitsPerSample: significantBitsPerSample,
		},
		validBitsPerSample: validBitsPerSample,
		channelMask:        channelMask,
		subFormat:          subFormat,
	}
}

func (f *ExtensibleMetaFormat) CompressionCode() uint16 {
	return uint16(extensibleCompressionCode)
}

func (f *ExtensibleMetaFormat) ExtraField() (b []byte) {
	b = make([]byte, 22)
	binary.LittleEndian.PutUint16(b[:2], uint16(f.validBitsPerSample))
	binary.LittleEndian.PutUint32(b[2:6], uint32(f.channelMask))
	copy(b[6:], f.subFormat[:])
	return
}
