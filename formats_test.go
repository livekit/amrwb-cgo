package amrwb

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

const frameMode2RTP = `f15411801dfff7817c56d9e3e33ddc01f6a0b2b45690b796b0113bb4d6b9913476`

func TestRTPToStorage(t *testing.T) {
	data, err := hex.DecodeString(frameMode2RTP)
	if err != nil {
		t.Fatal(err)
	}
	orig := data
	buf := make([]byte, len(data), len(data)+1)
	copy(buf, data)
	data = buf

	n := len(data)
	t.Logf("rtp: %x (last: %08b)", data[:n], data[n-1])
	n = rtp2storage(data[:n])
	data = data[:n]
	t.Logf("str: %x (last: %08b)", data[:n], data[n-1])
	n = storage2rtp(data[:n])
	data = data[:n]
	t.Logf("rtp: %x (last: %08b)", data[:n], data[n-1])
	if !bytes.Equal(orig, data[:n]) {
		t.Fatal("not equal")
	}
}

func TestStorageToRTP(t *testing.T) {
	frames := [9]string{
		0: "048f49019e7b2fe510f9464403957b9f9530",
		1: "0c9d191954c295069c79400c4f42905b9accb99ad8dbfc80",
		2: "149bd21929859ab506983f48147f1641469e9bbaaf68abdb8b9cef9a8f35bed8d8",
		3: "1c9bd21f49840ab5069c3a48045f11ec0ae5278fbbaa0deeeef1bbadf89d8b2a49cf93d998",
		4: "249bd21c7bf51ab506983f48046f101124c471bf8eecfb8f0beffdd9acda4f81f9d0ad59d89ea998f8",
		5: "2c9bd418b1d6b8b5069c3b48145f1d8293df20c8ef15e847fec87abd2d6dde65172daccfb2337d59d4216fd1d1a8a8",
		6: "349bd43d698400a506983d401c5f1628514bec36e79a69889baeec3adccd87d96815608cd9127ddcdd7f9e189d8e8f8dfd2d70",
		7: "3c9bd21de3651ab506983a58144f175dec30ecfe36aff20f85a8c37c077e5adeeb7acd5c3b767c53cf5b08b1f4fbaf5ebc1de962f6d64cdaf78558",
		8: "449bd21c79948ab50698dfee3a58047f17dde435b9ef2f7fc5c9d7aa8199a3fc54a2f7cefae93b44740ba5db8a39e0dfaf569f7a5aa93fb94cef3b5f18",
	}
	for mode, frame := range frames {
		t.Run(fmt.Sprintf("mode %d", mode), func(t *testing.T) {
			data, err := hex.DecodeString(frame)
			if err != nil {
				t.Fatal(err)
			}
			orig := data
			buf := make([]byte, len(data), len(data)+1)
			copy(buf, data)
			data = buf

			n := len(data)
			t.Logf("str: %x (last: %08b)", data[:n], data[n-1])
			n = storage2rtp(data[:n])
			data = data[:n]
			t.Logf("rtp: %x (last: %08b)", data[:n], data[n-1])
			n = rtp2storage(data[:n])
			data = data[:n]
			t.Logf("str: %x (last: %08b)", data[:n], data[n-1])
			if !bytes.Equal(orig, data[:n]) {
				t.Fatal("not equal")
			}
		})
	}
	data, err := hex.DecodeString(frameMode2RTP)
	if err != nil {
		t.Fatal(err)
	}
	orig := data
	buf := make([]byte, len(data), len(data)+1)
	copy(buf, data)
	data = buf

	n := len(data)
	t.Logf("rtp: %x (last: %08b)", data[:n], data[n-1])
	n = rtp2storage(data[:n])
	data = data[:n]
	t.Logf("str: %x (last: %08b)", data[:n], data[n-1])
	n = storage2rtp(data[:n])
	data = data[:n]
	t.Logf("rtp: %x (last: %08b)", data[:n], data[n-1])
	if !bytes.Equal(orig, data[:n]) {
		t.Fatal("not equal")
	}
}
