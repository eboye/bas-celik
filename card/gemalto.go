package card

import (
	"encoding/binary"
	"fmt"

	"github.com/ebfe/scard"
)

var GEMALTO_ATR_1 = []byte{
	0x3B, 0xFF, 0x94, 0x00, 0x00, 0x81, 0x31, 0x80,
	0x43, 0x80, 0x31, 0x80, 0x65, 0xB0, 0x85, 0x02,
	0x01, 0xF3, 0x12, 0x0F, 0xFF, 0x82, 0x90, 0x00,
	0x79,
}

var GEMALTO_ATR_2 = []byte{
	0x3B, 0xF9, 0x96, 0x00, 0x00, 0x80, 0x31, 0xFE,
	0x45, 0x53, 0x43, 0x45, 0x37, 0x20, 0x47, 0x43,
	0x4E, 0x33, 0x5E,
}

type Gemalto struct {
	smartCard *scard.Card
}

func connectGemalto(card *scard.Card) bool {
	data := []byte{0xF3, 0x81, 0x00, 0x00, 0x02, 0x53, 0x45, 0x52, 0x49, 0x44, 0x01}
	apu, _ := buildAPDU(0x00, 0xA4, 0x04, 0x00, data, 0)
	rsp, err := card.Transmit(apu)
	if err != nil || !responseOK(rsp) {
		return false
	}

	return true
}

func (card Gemalto) readFile(name []byte, trim bool) ([]byte, error) {
	output := make([]byte, 0)

	_, err := selectFile(card.smartCard, name, 4)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	data, err := read(card.smartCard, 0, 4)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	length := uint(binary.LittleEndian.Uint16(data[2:]))
	offset := uint(len(data))

	for length > 0 {
		data, err := read(card.smartCard, offset, length)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		output = append(output, data...)

		offset += uint(len(data))
		length -= uint(len(data))
	}

	if trim {
		return output[4:], nil
	}

	return output, nil
}
