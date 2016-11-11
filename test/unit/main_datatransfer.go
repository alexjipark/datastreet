package main

import (
	eyescli "github.com/tendermint/merkleeyes/client"

	"github.com/alexjipark/datastreet/app"
	"fmt"
	"encoding/json"
	"reflect"

	. "github.com/tendermint/go-common"
	"github.com/alexjipark/datastreet/test"
	"encoding/hex"
	"github.com/alexjipark/datastreet/types"
	"github.com/tendermint/go-wire"
)

func main() {

}

func testDataOwnership() {
	eyeCli := eyescli.NewLocalClient()
	chainID := "test_chain_id"
	datastApp := datastreet.NewDataStreet(eyeCli)
	dataToken := ""

	datastApp.SetOption("base/chainID", chainID)
	fmt.Println(datastApp.Info())

	kvz := loadGenesis("/Users/Park-jihun/Desktop/1_BlockChain/5_Tendermint/workspace/testnet_basecoin/bctest/app/genesis.json")
	for _, kv := range kvz {
		log := datastApp.SetOption(kv.Key, kv.Value)
		fmt.Println(Fmt("Set %v=%v. Log: %v"), kv.Key, kv.Value, log)
	}

	// Get the root account
	root := test.PrivateAccountFromSecret("test")
	// Make a bunch of PrivateAccount
	dstAccount := test.PrivateAccountFromSecret("test1")

	hex_root_addr := hex.EncodeToString(root.Account.PubKey.Address())
	hex_dst_addr  := hex.EncodeToString(dstAccount.Account.PubKey.Address())
	fmt.Printf("root's Address[%s].. dst's Address[%s]\n", hex_root_addr, hex_dst_addr)

	// ====== Send Data Ownership to Each Account?!
	tx := &types.SendTx{
		Inputs: []types.TxInput {
			types.TxInput{
				Address: root.Account.PubKey.Address(),
				PubKey:  root.Account.PubKey,
				Coins:   types.Coins{{dataToken,1}},
				Sequence: 0,
			},
		},
		Outputs: []types.TxOutput {
			types.TxOutput{
				Address:root.Account.PubKey.Address(),
				Coins:  types.Coins{{dataToken, 1}},
			},
		},
	}

	//Sign Request
	signBytes := tx.SignBytes(chainID)
	sig := root.PrivKey.Sign(signBytes)
	tx.Inputs[0].Signature = sig
	fmt.Println("tx: ", tx)

	//Write Request
	txBytes := wire.BinaryBytes(struct{types.Tx}{tx})
	datastApp.BeginBlock(0)
	res := datastApp.AppendTx(txBytes)
	fmt.Println("\nAppendTx: ", res)
	if res.IsErr() {
		Exit(Fmt("Failed :%v", res.Error()))
	} else {
		res = datastApp.AppendTx(txBytes)
		fmt.Println("\nAppendTx: ", res)
	}

	res = datastApp.Commit()
	if res.IsErr() {
		Exit (Fmt("Failed :%v", res.Error()))
	}
	fmt.Println("\nCommit: ", res)

	resEnd := datastApp.EndBlock(0)
	if len(resEnd) == 0 {
		fmt.Println("EndBlock.. no difference..\n")
	}
}

//============= Temp Testing =============//

type KeyValue struct {
	Key string `json:"key"`
	Value string `json:"value"`
}

func loadGenesis (filePath string) (kvz []KeyValue){
	kvz_ := []interface{}{}
	bytes, err := ReadFile(filePath)
	if err != nil {
		Exit("Loading Genesis File.." + err.Error())
	}
	err = json.Unmarshal(bytes, &kvz_)
	if err != nil {
		Exit ("Parsing Genesis File.." + err.Error())
	}
	if len(kvz_)%2 != 0 {
		Exit ("Genesis Cannot have an odd number of Items. Format = [key1, value1, key2, value2...]")
	}

	for i:=0; i <len(kvz_); i+=2 {
		keyIfc := kvz_[i]
		valueIfc := kvz_[i+1]
		var key, value string
		key, ok := keyIfc.(string)

		if !ok {
			Exit(Fmt("Genesis Had invalid key %v of type %v", keyIfc, reflect.TypeOf(keyIfc)))
		}
		if value_ , ok := valueIfc.(string); ok {
			value = value_
		} else {
			valueBytes, err := json.Marshal(valueIfc)
			if err != nil {
				Exit(Fmt("Genesis had invalid value %v:%v", value_, err.Error()))
			}
			value = string(valueBytes)
		}
		kvz = append(kvz, KeyValue{key,value})
	}
	return kvz
}