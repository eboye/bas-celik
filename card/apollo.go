package card

import (
	"encoding/binary"
	"fmt"

	"github.com/ebfe/scard"
)

type Apollo struct {
	smartCard *scard.Card
}

var APOLLO_ATR = []byte{
	0x3B, 0xB9, 0x18, 0x00, 0x81, 0x31, 0xFE, 0x9E, 0x80,
	0x73, 0xFF, 0x61, 0x40, 0x83, 0x00, 0x00, 0x00, 0xDF,
}

func (card Apollo) readFile(name []byte, trim bool) ([]byte, error) {
	output := make([]byte, 0)

	_, err := selectFile(card.smartCard, name, 4)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	data, err := read(card.smartCard, 0, 6)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	length := uint(binary.LittleEndian.Uint16(data[4:]))
	offset := uint(6)

	if trim {
		length -= 4
		offset += 4
	}

	for length > 0 {
		data, err := read(card.smartCard, offset, length)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		output = append(output, data...)

		offset += uint(len(data))
		length -= uint(len(data))
	}

	return output, nil
}
