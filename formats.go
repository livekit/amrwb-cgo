package amrwb

// Storage format
//
//  0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+
// |P|  FT   |Q|P|P|
// +-+-+-+-+-+-+-+-+

// Bandwidth-efficient RTP format
//
//  0 1 2 3   4 5 6 7 8 9
// +-+-+-+-+ +-+-+-+-+-+-+
// |  CMR  | |F|  FT   |Q|
// +-+-+-+-+ +-+-+-+-+-+-+

func storage2rtp(data []byte) int {
	// The old header is 8 bit, new one is 10 bit (4+6 bit).
	// We need an extra byte in the destination buffer.
	if cap(data) < len(data)+1 {
		panic("short buffer")
	}
	data = data[:len(data)+1]
	copy(data[2:], data[1:])

	toc := (data[0] >> 2) & 0b0_1111_1 // erase pad
	ft := (toc >> 1) & 0b1111          // aka coding mode
	q := (toc >> 0) & 0b1              // good quality
	_ = q

	bits := int(blockSizesBits[ft])

	readBits := bits + 8
	readBytes := readBits / 8
	if readBits%8 != 0 {
		readBytes++
	}

	writeBits := bits + 10
	writeBytes := writeBits / 8
	if writeBits%8 != 0 {
		writeBytes++
	}
	if bits == 0 {
		// unsupported frame, still works but may have extra padding
		readBytes = len(data)
		writeBytes = len(data)
	}

	const (
		cmr = byte(0b1111) // no codec mode request
		f   = byte(0b0)    // just one block in the frame
	)
	toc = (f << 5) | toc
	toc1 := (toc >> 2) & 0b1111
	toc2 := (toc >> 0) & 0b11

	data[0] = cmr<<4 | toc1
	data[1] = toc2 << 6
	cur := data[1]

	// 6 lower bits are empty now

	// [X X 7 6 5 4 3 2] <- [7 6 5 4 3 2 _ _]
	// [1 0 X X X X X X] <- [_ _ _ _ _ _ 1 0]

	i := 1
	// << 6
	for ; i < readBytes; i++ {
		next := data[i+1]
		data[i] = cur | ((next >> 2) & 0b11_1111)
		cur = (next & 0b11) << 6
	}
	data[i] = cur

	return writeBytes
}

func rtp2storage(data []byte) int {
	// The old header is 10 bit (4+6 bit), new one is 8 bit.
	// We might need to remove one byte from the output.

	cmr := (data[0] >> 4) & 0b1111
	toc1 := (data[0] >> 0) & 0b1111
	toc2 := (data[1] >> 6) & 0b11
	cur := ((data[1] >> 0) & 0b11_1111) << 2
	_ = cmr // ignored

	toc := (toc1 << 2) | toc2

	f := (toc >> 5) & 0b1     // has next block
	ft := (toc >> 1) & 0b1111 // aka coding mode
	q := (toc >> 0) & 0b1     // good quality
	_ = f                     // ignored, must be 0
	_ = q

	bits := int(blockSizesBits[ft])

	readBits := bits + 10
	readBytes := readBits / 8
	if readBits%8 != 0 {
		readBytes++
	}

	writeBits := bits + 8
	writeBytes := writeBits / 8
	if writeBits%8 != 0 {
		writeBytes++
	}
	if bits == 0 {
		// unsupported frame, still works but may have extra padding
		readBytes = len(data)
		writeBytes = len(data)
	}

	toc &= 0b0_1111_1 // erase f

	data[0] = toc << 2

	// writes aligned, but reading with offset of 6 bits
	// first 6 are already in cur

	// [_ _ _ _ _ _ 7 6] <- [7 6 X X X X X X]
	// [5 4 3 2 1 0 _ _] <- [X X 5 4 3 2 1 0]

	i := 1
	// << 6
	for ; i < readBytes-1; i++ {
		next := data[i+1]
		data[i] = cur | ((next >> 6) & 0b11)
		cur = ((next >> 0) & 0b11_1111) << 2
	}
	data[i] = cur

	return writeBytes
}
