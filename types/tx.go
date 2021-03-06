package types

import (
	"github.com/tendermint/go-wire"
	"github.com/tendermint/go-crypto"
	. "github.com/tendermint/go-common"
	"bytes"
	tmsp "github.com/tendermint/tmsp/types"
)

/*
Tx (transaction) is an atomic operation on the ledger state

Account Types:
- SendTx	Send coins to Address
- AppTx		Send a msg to a contract that runs in the vm
 */

type Tx interface {
	AssertIsTx()
	SignBytes(chainID string) []byte
}

// Types of Tx implementation
const (
	// Account tracsactions
	TxTypeSend = byte(0x01)
	TxTypeApp  = byte(0x02)
)

func (_ *SendTx) AssertIsTx() {}
func (_ *AppTx)  AssertIsTx() {}

//go-wire : Go library for encoding/decoding structures into binary and JSON format.
var _= wire.RegisterInterface(
	struct{Tx}{},
	wire.ConcreteType{&SendTx{}, TxTypeSend},
	wire.ConcreteType{&AppTx{},  TxTypeApp},
)

type TxInput struct {
	Address 	[]byte 			`json:"address"`	// Hash of the PubKey
	Coins		Coins			`json:"coins"`		//
	Sequence	int			`json:"sequence"`	// Must be 1 greater than last commit
	Signature 	crypto.Signature	`json:"signature"`	// Depends on the PubKey type and the whole Tx
	PubKey		crypto.PubKey		`json:"pub_key"`	// Is present iff Sequence == 0

	// Unique Key For Data
	DataIndex	string			`json:"data_index"'`	// Data Index
}

func (txIn TxInput) String() string {
	return Fmt("TxInput{%X,%v,%v,%v,%v, %s}", txIn.Address,
			txIn.Coins, txIn.Sequence, txIn.Signature, txIn.PubKey, txIn.DataIndex)
}

type TxOutput struct {
	Address		[]byte			`json:"address"`	// Hash of the PubKey
	Coins		Coins			`json:"coins"`		//
}

func (txOut TxOutput) String() string {
	return Fmt("TxOutput{%X,%v}", txOut.Address, txOut.Coins)
}

func (txIn TxInput) ValidateBasic() tmsp.Result {
	// TBD - To Be Developed..
	return tmsp.OK
}

func (txOut TxOutput) ValidateBasic() tmsp.Result {
	// TBD - To Be Developed..
	return tmsp.OK
}


//----------------------------------

type SendTx struct {
	Fee 	int64		`json:"fee"`	// Fee
	Gas	int64		`json:"fee"`	// Gas
	Inputs	[]TxInput	`json:"inputs"`
	Outputs	[]TxOutput	`json:"outputs"`
}

func (tx *SendTx) SignBytes(chainID string) []byte {
	signBytes := wire.BinaryBytes(chainID)
	sigz := make( []crypto.Signature, len(tx.Inputs))
	for i, input := range tx.Inputs {
		sigz[i] = input.Signature
		tx.Inputs[i].Signature = nil
	}
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	for i := range tx.Inputs {
		tx.Inputs[i].Signature = sigz[i]
	}
	return signBytes
	/*

	signBytes := wire.BinaryBytes(chainID)
	sigz := make([]crypto.Signature, len(tx.Inputs))
	for i, input := range tx.Inputs {
		sigz[i] = input.Signature
		tx.Inputs[i].Signature = nil
	}
	signBytes = append(signBytes, wire.BinaryBytes(tx)...)
	for i := range tx.Inputs {
		tx.Inputs[i].Signature = sigz[i]
	}
	return signBytes
	 */
}

func (tx *SendTx) SetSignature(addr []byte, sig crypto.Signature) bool {
	for i, input := range tx.Inputs {
		if bytes.Equal(input.Address, addr) {
			tx.Inputs[i].Signature = sig
			return true	// AJ - why it returns?! only one matching to addrress in Inputs??
		}
	}
	return false
}

//----------------------------------

type AppTx struct {
	Fee 	int64	`json:"fee"`	// Fee
	Gas	int64	`json:"gas"`	// Gas
	Type	byte	`json:"type"`	// Which App
	Input	TxInput	`json:"input"`	// ..
	Data	[]byte	`json:"data"`
}

