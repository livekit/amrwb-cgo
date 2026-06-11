// Package amrwb implements AMR-WB audio codec encoding and decoding.
package amrwb

const (
	Magic        = "#!AMR-WB\n" // File magic for AMR-WB, as described in RFC 4867 #5.1.
	PCMFrameSize = 320          // PCM frame size at 16kHz
	FrameSizeMax = 61           // max encoded AMR-WB block size
)

// Format represent the variation of the encoding/decoding format for AMR-WB.
type Format int

const (
	// Storage format for AMR-WB. See RFC 4867 #5
	Storage = Format(iota)
	// RTPBandwidthEfficient is a bandwidth-efficient AMR-WB format for RTP payloads.
	RTPBandwidthEfficient
)

type PCMFrame = [PCMFrameSize]int16

var blockSizes = [16]byte{
	18,
	24,
	33,
	37,
	41,
	47,
	51,
	59,
	61,
	6,
	6,
	0,
	0,
	0,
	1,
	1,
}

// blockSizesBits is the number of bits for AMR-WB blocks, excluding TOC.
var blockSizesBits = [16]uint16{
	132, // 6.60k
	177, // 8.85k
	253, // 12.65k
	285, // 14.25k
	317, // 15.85k
	365, // 18.25k
	397, // 19.85k
	461, // 23.05k
	477, // 23.85k
}

// BlockSize returns size of the first AMR-WB block in file data, as described in RFC 4867 #5.3.
// This function is not suitable for use in parsing of AMR-WB RTP payload, as it uses a different format.
func BlockSize(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	mode := int((data[0] >> 3) & 0x0f)
	if mode >= len(blockSizes) {
		return -1
	}
	n := int(blockSizes[mode])
	if n > len(data) {
		return -1
	}
	return n
}

// Blocks counts the number of AMR-WB blocks in file data, as described in RFC 4867 #5.3.
// See BlockSize for details.
func Blocks(data []byte) (sz, cnt int) {
	for len(data) > 0 {
		n := BlockSize(data)
		if n <= 0 {
			break
		}
		sz += n
		cnt++
		data = data[sz:]
	}
	return
}
