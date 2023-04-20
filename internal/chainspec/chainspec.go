package chainspec

import (
	"encoding/binary"
	"fmt"

	"github.com/OneOfOne/xxhash"
	"github.com/crypto2lab/ground/internal/primitives"
)

type Database interface {
	Put(key, value []byte) error
}

type ChainSpec struct {
	Genesis *Genesis `json:"genesis"`
}

func (c *ChainSpec) StoreGenesis(database Database) error {
	if c.Genesis != nil {
		return c.storeGenesis(database)
	}

	return nil
}

func (s *ChainSpec) storeGenesis(database Database) error {
	for _, acc := range s.Genesis.Accounts {
		encodedKey, err := acc.StoreKey()
		if err != nil {
			return fmt.Errorf("while generating storage key: %w", err)
		}

		encodedValue := acc.StoreValue()
		err = database.Put(encodedKey, encodedValue)
		if err != nil {
			return fmt.Errorf("while inserting account: %w", err)
		}
	}

	return nil
}

type Genesis struct {
	Runtime  string              `json:"runtime"`
	Accounts []*ChainSpecAccount `json:"accounts"`
}

type ChainSpecAccount struct {
	PublicAddress primitives.PublicAddress `json:"public_address"`
	Currency      primitives.Currency      `json:"currency"`
	Balance       primitives.Balance       `json:"balance"`
}

func (ca *ChainSpecAccount) StoreKey() ([]byte, error) {
	const accountModule = "ACCOUNT::"

	encPubkey, err := ca.PublicAddress.Encode()
	if err != nil {
		return nil, fmt.Errorf("while encoding public key: %w", err)
	}

	encCurrency, err := ca.Currency.Encode()
	if err != nil {
		return nil, fmt.Errorf("while encoding currency: %w", err)
	}

	//mod::module::publickey::currency
	storageKey := fmt.Sprintf("mod::%s::%s::%s",
		accountModule, encPubkey, encCurrency)

	hasher := xxhash.NewS64(0)
	hasher.WriteString(storageKey)

	result := hasher.Sum64()
	hash := make([]byte, 8)
	binary.LittleEndian.PutUint64(hash, result)

	return hash, nil
}

func (ca *ChainSpecAccount) StoreValue() []byte {
	enc, _ := ca.Balance.Encode()
	return enc
}
