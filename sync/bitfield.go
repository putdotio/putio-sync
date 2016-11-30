package sync

import (
	"bytes"
	"encoding/binary"
	"encoding/json"

	"github.com/cenk/bitfield"
)

// Bitfield is a gob wrapper for bitfield.Bitfield.
//
// The bitfield.Bitfield type has unexported fields, which gob.Encoder cannot
// access. We therefore write a BinaryMarshal/BinaryUnmarshal method pair to
// allow us to send and receive the type with the gob package.
type Bitfield struct {
	length uint32
	*bitfield.Bitfield
}

// MarshalBinary returns internal bytes representation of Bitfield.
func (b Bitfield) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, b.length)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.BigEndian, b.Bytes())
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary creates new Bitfield from the given data.
func (b *Bitfield) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, binary.BigEndian, &b.length)
	if err != nil {
		return err
	}
	b.Bitfield = bitfield.NewBytes(buf.Bytes(), b.length)
	return nil
}

func (b *Bitfield) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Length int `json:"bit_count_all"`
		Count  int `json:"bit_count_set"`
	}{
		Length: int(b.Len()),
		Count:  int(b.Count()),
	})
}
