package amrwb

import amrenc "github.com/livekit/amrwb-cgo/enc"

// Mode specifies the bit rate the codec supports.
type Mode = amrenc.Mode

const (
	Mode6Kb  = amrenc.Mode6Kb  // 6.60kbps
	Mode8Kb  = amrenc.Mode8Kb  // 8.85kbps
	Mode12Kb = amrenc.Mode12Kb // 12.65kbps
	Mode14Kb = amrenc.Mode14Kb // 14.25kbps
	Mode16Kb = amrenc.Mode16Kb // 15.85bps
	Mode18Kb = amrenc.Mode18Kb // 18.25bps
	Mode20Kb = amrenc.Mode20Kb // 19.85kbps
	Mode23Kb = amrenc.Mode23Kb // 23.05kbps
	Mode24Kb = amrenc.Mode24Kb // 23.85kbps

	Fastest = Mode6Kb  // fastest encoding
	Best    = Mode24Kb // best encoding
)

// NewEncoder creates a new Encoder.
// Caller must call Close to avoid resource leak.
func NewEncoder(mode Mode) *Encoder {
	return &Encoder{
		Encoder: amrenc.New(mode),
	}
}

// Encoder for AMR-WB audio codec.
type Encoder struct {
	*amrenc.Encoder
	buf [FrameSizeMax]byte
}

// Encode PCM16 audio frame and append it to dst.
func (e *Encoder) Encode(dst []byte, src *PCMFrame) []byte {
	n := e.Encoder.Encode(&e.buf, src)
	return append(dst, e.buf[:n]...)
}
