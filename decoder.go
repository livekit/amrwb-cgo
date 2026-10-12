package amrwb

import (
	"errors"

	amrdec "github.com/livekit/amrwb-cgo/dec"
)

var (
	ErrInvalidBlock = errors.New("amrwb: invalid block")
)

// NewDecoder creates a new Decoder.
// Caller must call Close to avoid resource leak.
func NewDecoder(format Format) *Decoder {
	return &Decoder{
		Decoder: amrdec.New(),
		format:  format,
	}
}

// Decoder for AMR-WB audio codec.
type Decoder struct {
	*amrdec.Decoder
	buf    [FrameSizeMax + 1]byte
	format Format
}

// Decode PCM16 audio frame from src and return a number of bytes read.
// It returns ErrInvalidBlock if the block is empty or malformed.
func (d *Decoder) Decode(dst *PCMFrame, src []byte) (int, error) {
	if len(src) == 0 {
		return 0, ErrInvalidBlock
	}
	read := 0
	if d.format == RTPOctetAligned {
		// Skip CMR byte, the rest of the frame is in the storage format.
		// Interleaving is not supported, thus there's no ILL/ILP byte.
		src = src[1:]
		read = 1
		if len(src) == 0 {
			return 0, ErrInvalidBlock
		}
	}
	n := len(src)
	if n > FrameSizeMax {
		n = BlockSize(src)
		if n <= 0 {
			return 0, ErrInvalidBlock
		}
	}
	copy(d.buf[:], src[:n])
	read += n
	switch d.format {
	case RTPBandwidthEfficient:
		n = rtp2storage(d.buf[:n])
	case RTPOctetAligned:
		// Only a single frame per payload is supported, so F bit is ignored.
		// It occupies the padding bit in the storage format.
		d.buf[0] &^= tocFlagF
	}
	d.Decoder.Decode(dst, (*[61]byte)(d.buf[:FrameSizeMax]), false)
	return read, nil
}
