package chainspec

import (
	"github.com/crypto2lab/ground/internal/primitives"
)

type Database interface {
	Put(key, value []byte) error
}

type ChainSpec struct {
	Genesis *Genesis `json:"genesis"`
}

type Genesis struct {
	Code     string              `json:"code"`
	Accounts []*ChainSpecAccount `json:"accounts"`
}

type ChainSpecAccount struct {
	PublicAddress primitives.PublicAddress `json:"public_address"`
	Currency      primitives.Currency      `json:"currency"`
	Balance       primitives.Balance       `json:"balance"`
}
