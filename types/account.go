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


//-------------------------------
type AccountGetter interface {
	GetAccount(addr []byte) *Account
}
type AccountSetter interface {
	SetAccount(addr []byte, acc *Account)
}

type AccountGetterSetter interface {
	GetAccount(addr []byte) *Account
	SetAccount(addr []byte, acc *Account)
}

