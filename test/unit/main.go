package main

import (
	eyescli "github.com/tendermint/merkleeyes/client"
	"github.com/alexjipark/datastreet/app"
	"fmt"
	"github.com/alexjipark/datastreet/test"
	"github.com/alexjipark/datastreet/types"
	"github.com/tendermint/go-wire"
	. "github.com/tendermint/go-common"
)

func main() {
	testSendTx()
}

func testSendTx() {
	eyesCli := eyescli.NewLocalClient()
	chainID := "test_chain_id"
	datastApp := datastreet.NewDataStreet(eyesCli)

	datastApp.SetOption("base/chainID", chainID)
	fmt.Println(datastApp.Info())

	test1PrivateAcc := test.PrivateAccountFromSecret("test1")
	test2PrivateAcc := test.PrivateAccountFromSecret("test2")

	// Seed DataStreetApp with account
	test1Acc := test1PrivateAcc.Account
	test1Acc.Balance = types.Coins{{"",1000}}
	fmt.Println(datastApp.SetOption("base/account", string(wire.JSONBytes(test1Acc))))

	// Construct a SendTx signature
	tx := &types.SendTx{
		Fee: 0,
		Gas: 0,
		Inputs : []types.TxInput {
			types.TxInput{
				Address: test1PrivateAcc.Account.PubKey.Address(),
				PubKey:  test1PrivateAcc.Account.PubKey,
				Coins: types.Coins{{"",1}},
				Sequence: 1,
			},
		},
		Outputs: []types.TxOutput {
			types.TxOutput{
				Address: test2PrivateAcc.Account.PubKey.Address(),
				Coins: types.Coins{{"",1}},
			},
		},
	}

	// Sign Request
	signBytes := tx.SignBytes(chainID)
	fmt.Printf("Sign bytes: %X\n", signBytes)
	sig := test1PrivateAcc.PrivKey.Sign(signBytes)
	tx.Inputs[0].Signature = sig
	fmt.Printf("Signed TX bytes: %X\n", wire.BinaryBytes(struct{types.Tx}{tx}))

	// Write Request
	txBytes := wire.BinaryBytes(struct{types.Tx}{tx})
	res := datastApp.AppendTx(txBytes)
	fmt.Println(res)
	if res.IsErr() {
		Exit(Fmt("Failed :%v", res.Error()))
	}


}