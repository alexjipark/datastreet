package types

import "github.com/tendermint/go-crypto"

type Account struct {
	PubKey 		crypto.PubKey	`json:"pub_key"`	// May be nil, if not known
	Sequence 	int 		`json:"sequence"`
	Balance		Coins		`json:"coins"`
}

type PrivateAccount struct {
	crypto.PrivKey
	Account
}
