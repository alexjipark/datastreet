package test

import (
	"github.com/alexjipark/datastreet/types"
	"github.com/tendermint/go-crypto"
)

// Create a PrivAccount from secret
// The amount is not set
func PrivateAccountFromSecret(secret string) types.PrivateAccount {
	privateKey := crypto.GenPrivKeyEd25519FromSecret([]byte(secret))
	privateAccount := types.PrivateAccount{
		PrivKey : privateKey,
		Account : types.Account{
			PubKey: privateKey.PubKey(),
			Sequence: 0,
		},
	}
	return privateAccount
}