package main

import (
	"github.com/tendermint/go-rpc/client"
	. "github.com/tendermint/go-common"
	"fmt"
	"github.com/alexjipark/datastreet/test"
	"github.com/alexjipark/datastreet/types"
	"github.com/tendermint/go-wire"
	"github.com/tendermint/go-rpc/types"
	"github.com/gorilla/websocket"

)
func testQuery(addr []byte){

	/*
	txBytes := wire.BinaryBytes(struct{types.Tx}{tx})
	request := rpctypes.NewRPCRequest("fakeid", "broadcast_tx_sync", Arr(txBytes))
	fmt.Println("request: ", request)
	reqBytes := wire.JSONBytes(request)

	err = ws.WriteMessage(websocket.TextMessage, reqBytes)
	if err != nil {
		Exit("writing websocket request: " + err.Error())
	}

	 */

}
func main() {
	ws := rpcclient.NewWSClient("35.161.51.6:46657", "/websocket")
	chainID := "chain-AMUKE0"

	_,err := ws.Start()
	if err != nil {
		Exit(err.Error())
	}
	var counter = 0;

	// Read a bunch of responses
	go func() {
		for {
			res, ok := <-ws.ResultsCh
			if !ok {
				fmt.Println("Not ok from rpcclient")
				break
			}
			fmt.Println(counter, "res:", Blue(string(res)))
		}
	}()

	// Get the root account
	root := test.PrivateAccountFromSecret("test")
	sequence := int(0)
	//Make a bunch of PrivateAccount
	destAccount := test.PrivateAccountFromSecret("test1")

	// ====== Query
	addr := root.Account.PubKey.Address()
	fmt.Printf("Addr: %X", addr)
	queryBytes := make([]byte, 1+ wire.ByteSliceSize(addr))
	buf := queryBytes
	buf[0] = 0x02
	buf = buf[1:]
	wire.PutByteSlice(buf, addr)
	fmt.Println("query: ", queryBytes)

	requestQ := rpctypes.NewRPCRequest("fakeid", "tmsp_query", Arr(queryBytes))
	fmt.Println("request: ", requestQ)
	reqBytesQ := wire.JSONBytes(requestQ)
	fmt.Println("reqBytes: ", reqBytesQ)

	err = ws.WriteMessage(websocket.TextMessage, reqBytesQ)
	if err != nil {
		Exit("writing websocket request: " + err.Error())
	}


	// ======= Send coins to each account
	tx := &types.SendTx{
		Inputs: []types.TxInput {
			types.TxInput {
				Address: root.Account.PubKey.Address(),
				PubKey: root.Account.PubKey,
				Coins:	types.Coins {{"", 10}},
				Sequence: sequence,
			},
		},
		Outputs: []types.TxOutput {
			types.TxOutput {
				Address: destAccount.Account.PubKey.Address(),
				Coins: types.Coins{{"", 9}},
			},
		},
	}
	sequence += 1

	//Sign request
	signBytes := tx.SignBytes(chainID)
	sig := root.PrivKey.Sign(signBytes)
	tx.Inputs[0].Signature = sig
	fmt.Println("tx: ", tx)

	//Write request
	txBytes := wire.BinaryBytes(struct{types.Tx}{tx})
	request := rpctypes.NewRPCRequest("fakeid", "broadcast_tx_sync", Arr(txBytes))
	fmt.Println("request: ", request)
	//reqBytes := wire.JSONBytes(request)

	//err = ws.WriteMessage(websocket.TextMessage, reqBytes)
	if err != nil {
		Exit("writing websocket request: " + err.Error())
	}

	// Wait Forever
	TrapSignal(func() {
		ws.Stop()
	})
}
