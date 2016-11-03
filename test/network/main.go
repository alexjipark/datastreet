package main

import (
	"github.com/tendermint/go-rpc/client"
	. "github.com/tendermint/go-common"
	"fmt"
	"github.com/alexjipark/datastreet/test"
	"github.com/alexjipark/datastreet/types"
	"github.com/tendermint/go-wire"
	"github.com/tendermint/go-rpc/types"

	tmsp "github.com/tendermint/tmsp/types"

	"github.com/gorilla/websocket"

	tdtypes "github.com/tendermint/tendermint/types"
	"bytes"
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
type ResultData struct {
	Result tmsp.Result `json:"result"`
}

func main() {
	ws := rpcclient.NewWSClient("35.160.145.128:46657", "/websocket")
	chainID := "chain-AMUKE0"

	_,err := ws.Start()
	if err != nil {
		Exit(err.Error())
	}
/*
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
*/
	// Get the root account
	root := test.PrivateAccountFromSecret("test")
	sequence := int(14)
	//Make a bunch of PrivateAccount
	destAccount := test.PrivateAccountFromSecret("test1")

	//====== Check Account
	fmt.Printf("Private Key : %X\n", root.PrivKey)
	fmt.Printf("Public Byte : %X\n", root.Account.PubKey.Bytes())
	fmt.Printf("Public Addr : %X\n", root.Account.PubKey.Address())

	// ======= Send coins to each account
	tx := &types.SendTx{
		Inputs: []types.TxInput {
			types.TxInput {
				Address: root.Account.PubKey.Address(),
				PubKey: root.Account.PubKey,
				Coins:	types.Coins {{"KRW", 10}},
				Sequence: sequence + 1,
			},
		},
		Outputs: []types.TxOutput {
			types.TxOutput {
				Address: destAccount.Account.PubKey.Address(),
				Coins: types.Coins{{"KRW", 10}},
			},
		},
	}

	//Sign request
	signBytes := tx.SignBytes(chainID)
	sig := root.PrivKey.Sign(signBytes)
	tx.Inputs[0].Signature = sig
	fmt.Println("tx: ", tx)

	txBytes := wire.BinaryBytes(struct{types.Tx}{tx})

	//==Event Subscribe..
	eid := tdtypes.EventStringTx(tdtypes.Tx(txBytes))
	if err = ws.Subscribe(eid); err != nil {
		fmt.Println("Subscribe Error .. ", err.Error())
	}
	defer func() {
		ws.Unsubscribe(eid)
	}()

	//Write request
	request := rpctypes.NewRPCRequest("fakeid", "broadcast_tx_async", Arr(txBytes))
	//request := rpctypes.NewRPCRequest("fakeid", "broadcast_tx_commit", Arr(txBytes))
	fmt.Println("request: ", request)


	reqBytes := wire.JSONBytes(request)
	err = ws.WriteMessage(websocket.TextMessage, reqBytes)
	if err != nil {
		Exit("writing websocket request: " + err.Error())
	}

	//waitForEvent(t, wsc, eid, true, func() {}, func(eid string, b interface{}) error {
	test.WaitForEvent(ws, eid, true, func(){}, func(eid string, b interface{}) error {
		evt, ok := b.(tdtypes.EventDataTx)
		if !ok {
			fmt.Println("Get Wrong Event Type ", b)
		}else {
			if bytes.Compare(evt.Tx, txBytes) != 0 {
				fmt.Println("Event Returned different tx")
			}
			fmt.Println("Event: ", evt)
			fmt.Println("Event Tx :", evt.Tx)
			fmt.Println("Event Tx Hash: ", evt.Tx.Hash())
			fmt.Println("Event Code :", evt.Code)
			fmt.Println("Event Log :", evt.Log)
			fmt.Println("Event Result :", evt.Result)
		}

		return nil
	})

	// Wait Forever
	TrapSignal(func() {
		ws.Stop()
	})
}

