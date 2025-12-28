package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const HeaderSize = 44

type Header struct {
	// Size is the number of bytes in the file minus 8 for the RIFF and size fields.
	Size uint32
	// Channels is the number of channels. Mono = 1, Stereo = 2, etc.
	Channels uint16
	// SampleRate is the sampling frequency (e.g., 44100 Hz).
	SampleRate uint32
	// ByteRate is SampleRate * NumChannels * BitsPerSample / 8.
	ByteRate uint32
	// BlockAlign is NumChannels * BitsPerSample / 8.
	BlockAlign uint16
	// BitsPerSample is the bit depth (e.g., 16, 24, or 32 bits).
	BitsPerSample uint16
	// Subchunk2Size is the number of bytes in the data section.
	Subchunk2Size uint32
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

	return &Header{
		Size:          binary.LittleEndian.Uint32(b[4:8]),
		Channels:      binary.LittleEndian.Uint16(b[22:24]),
		SampleRate:    binary.LittleEndian.Uint32(b[24:28]),
		ByteRate:      binary.LittleEndian.Uint32(b[28:32]),
		BlockAlign:    binary.LittleEndian.Uint16(b[32:34]),
		BitsPerSample: binary.LittleEndian.Uint16(b[34:36]),
		Subchunk2Size: binary.LittleEndian.Uint32(b[40:44]),
	}, nil
}
