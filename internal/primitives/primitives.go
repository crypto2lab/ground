package primitives

import "encoding/binary"

type PublicAddress string

func (p PublicAddress) Encode() ([]byte, error) {
	return nil, nil
}

func (p *PublicAddress) Decode(d []byte) error {
	return nil
}

type Currency string

func (c Currency) Encode() ([]byte, error) {
	return nil, nil
}

func (c *Currency) Decode(d []byte) error {
	return nil
}

// balance encoded using LittleEndian
type Balance uint64

func (b Balance) Encode() ([]byte, error) {
	encoded := make([]byte, 8)
	binary.LittleEndian.PutUint64(encoded, uint64(b))
	return encoded, nil
}

func (b *Balance) Decode(d []byte) error {
	*b = Balance(binary.LittleEndian.Uint64(d))
	return nil
}
