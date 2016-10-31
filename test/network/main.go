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

	tmsp "github.com/tendermint/tmsp/types"

	"encoding/json"
	"encoding/hex"
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
	ws := rpcclient.NewWSClient("35.161.26.62:46657", "/websocket")
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

			//==== Check the result
			//res - *json.RawMessage
			var result []interface{}
			err := json.Unmarshal([]byte(string(res)), &result)
			if err != nil {
				fmt.Println("Error in Unmarshalling with ", err.Error())
			}
			fmt.Printf("result num :%v\n" , len(result))
			fmt.Println(result[1])	// map
			fmt.Println(result[0])	// 112

			resData := result[1].(map[string]interface{})["result"].(map[string]interface{})["Data"]

			//fmt.Println([]byte(str))
			hexBytes, err := hex.DecodeString(resData.(string))
			fmt.Println(hexBytes)

			var acc *types.Account
			err = wire.ReadBinaryBytes(hexBytes, &acc)
			if err != nil {
				fmt.Printf("Error Reading Account %X error: %v",
					resData, err.Error())
			}
			fmt.Printf("Account : %X\n", acc.PubKey)
			fmt.Printf("Balance : %v\n", acc.Balance)
			fmt.Printf("Sequence : %v", acc.Sequence)

		}
	}()

	// Get the root account
	root := test.PrivateAccountFromSecret("test")
	sequence := int(0)
	//Make a bunch of PrivateAccount
	destAccount := test.PrivateAccountFromSecret("test1")

	//====== Check Account
	fmt.Printf("Private Key : %X\n", root.PrivKey)
	fmt.Printf("Public Byte : %X\n", root.Account.PubKey.Bytes())
	fmt.Printf("Public Addr : %X\n", root.Account.PubKey.Address())

	// ====== Query

	addrBytes := root.Account.PubKey.Address()
	fmt.Printf("Addr: %X\n", addrBytes)

	queryBytes := make([]byte, 1+ wire.ByteSliceSize(addrBytes))

	buf := queryBytes
	buf[0] = 0x01	//Get TypeByte
	buf = buf[1:]
	wire.PutByteSlice(buf, addrBytes)

	fmt.Printf("query: %X\n", queryBytes)

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
				Coins:	types.Coins {{"USD", 10}},
				Sequence: sequence + 1,
			},
		},
		Outputs: []types.TxOutput {
			types.TxOutput {
				Address: destAccount.Account.PubKey.Address(),
				Coins: types.Coins{{"USD", 10}},
			},
		},
	}

	//Sign request
	signBytes := tx.SignBytes(chainID)
	sig := root.PrivKey.Sign(signBytes)
	tx.Inputs[0].Signature = sig
	fmt.Println("tx: ", tx)

	//Write request
	txBytes := wire.BinaryBytes(struct{types.Tx}{tx})
	request := rpctypes.NewRPCRequest("fakeid", "broadcast_tx_sync", Arr(txBytes))
	fmt.Println("request: ", request)
	reqBytes := wire.JSONBytes(request)

	err = ws.WriteMessage(websocket.TextMessage, reqBytes)
	if err != nil {
		Exit("writing websocket request: " + err.Error())
	}

	// Wait Forever
	TrapSignal(func() {
		ws.Stop()
	})
}
