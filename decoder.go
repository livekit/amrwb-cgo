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
func NewDecoder() *Decoder {
	return &Decoder{
		Decoder: amrdec.New(),
	}
}

// Decoder for AMR-WB audio codec.
type Decoder struct {
	*amrdec.Decoder
	buf [FrameSizeMax]byte
}

// Decode PCM16 audio frame from src and return a number of bytes read.
// It returns ErrInvalidBlock if the block is empty or malformed.
func (d *Decoder) Decode(dst *PCMFrame, src []byte) (int, error) {
	if len(src) == 0 {
		return 0, ErrInvalidBlock
	}
	n := BlockSize(src)
	if n <= 0 {
		return 0, ErrInvalidBlock
	}
	copy(d.buf[:], src[:n])
	d.Decoder.Decode(dst, &d.buf, false)
	return n, nil
}
