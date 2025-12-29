package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const HeaderSize = 44

type Header struct {
	// Size is the number of bytes in the file minus 8 for the RIFF and size fields.
	Size uint32
	// Channels is the number of channels. Mono = 1, Stereo = 2, etc.
	Channels uint16
	// SampleRate is the sampling frequency (e.g., 44100 Hz).
	SampleRate uint32
	// ByteRate is SampleRate * NumChannels * BitsPerSample / 8. This value is the average number
	// of bytes per second at which the waveform data should be transferred.
	ByteRate uint32
	// BlockAlign is NumChannels * BitsPerSample / 8.
	BlockAlign uint16
	// BitsPerSample is the bit depth (e.g., 16, 24, or 32 bits).
	BitsPerSample uint16
	// DataSize is the number of bytes in the data section.
	DataSize uint32
}

func NewHeader(b []byte) (*Header, error) {
	if len(b) < HeaderSize {
		return nil, fmt.Errorf("need %d bytes", HeaderSize)
	}

	if string(b[0:4]) != "RIFF" {
		return nil, fmt.Errorf("not RIFF chunk")
	}

	if string(b[8:12]) != "WAVE" {
		return nil, fmt.Errorf("not WAVE format")
	}

	if string(b[12:16]) != "fmt " {
		return nil, fmt.Errorf("incorrect format chunk marker")
	}

	subChunk1Size := binary.LittleEndian.Uint32(b[16:20])
	if subChunk1Size != 16 {
		return nil, errors.New("wrong subchunk 1 size")
	}

	audioFormat := binary.LittleEndian.Uint16(b[20:22])
	if audioFormat != 1 {
		return nil, errors.New("wrong audio format")
	}

	if string(b[36:40]) != "data" {
		return nil, fmt.Errorf("expected 'data', got %s", string(b))
	}

	channels := binary.LittleEndian.Uint16(b[22:24])
	sampleRate := binary.LittleEndian.Uint32(b[24:28])
	byteRate := binary.LittleEndian.Uint32(b[28:32])
	blockAlign := binary.LittleEndian.Uint16(b[32:34])
	bitsPerSample := binary.LittleEndian.Uint16(b[34:36])

	if byteRate != sampleRate*uint32(channels)*uint32(bitsPerSample)/8 {
		return nil, errors.New("incorrect byte rate")
	}

	if blockAlign != channels*bitsPerSample/8 {
		return nil, errors.New("incorrect byte rate")
	}

	return &Header{
		Size:          binary.LittleEndian.Uint32(b[4:8]),
		Channels:      channels,
		SampleRate:    sampleRate,
		ByteRate:      byteRate,
		BlockAlign:    blockAlign,
		BitsPerSample: bitsPerSample,
		DataSize:      binary.LittleEndian.Uint32(b[40:44]),
	}, nil
}

func (h *Header) String() string {
	return fmt.Sprintf("Size (bytes): %d\n", h.Size) +
		fmt.Sprintf("Channels: %d\n", h.Channels) +
		fmt.Sprintf("SampleRate (Hz): %d\n", h.SampleRate) +
		fmt.Sprintf("ByteRate (bytes/sec): %d\n", h.ByteRate) +
		fmt.Sprintf("BlockAlign (bytes): %d\n", h.BlockAlign) +
		fmt.Sprintf("DataSize (bytes): %d\n", h.BitsPerSample) +
		fmt.Sprintf("Duration (s): %f", h.Duration().Seconds())
}

func (h *Header) MarshalJSON() ([]byte, error) {
	type jsonHeader struct {
		Size          uint32  `json:"size"`
		Channels      uint16  `json:"channels"`
		SampleRate    uint32  `json:"sampleRate"`
		ByteRate      uint32  `json:"byteRate"`
		BlockAlign    uint16  `json:"blockAlign"`
		BitsPerSample uint16  `json:"bitsPerSample"`
		DataSize      uint32  `json:"dataSize"`
		Duration      float64 `json:"duration"`
	}

	return json.Marshal(jsonHeader{
		Size:          h.Size,
		Channels:      h.Channels,
		SampleRate:    h.SampleRate,
		ByteRate:      h.ByteRate,
		BlockAlign:    h.BlockAlign,
		BitsPerSample: h.BitsPerSample,
		DataSize:      h.DataSize,
		Duration:      h.Duration().Seconds(),
	})
}

func (h *Header) Duration() time.Duration {
	nanos := float64(h.DataSize) / float64(h.ByteRate) * 1e9
	return time.Duration(nanos)
}
