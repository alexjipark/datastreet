package main

import (
	. "github.com/tendermint/go-common"
	"fmt"
	"github.com/alexjipark/datastreet/test"
	"github.com/alexjipark/datastreet/types"
	"github.com/tendermint/go-wire"

	"github.com/alexjipark/datastreet/app"
	eyescli "github.com/tendermint/merkleeyes/client"
	"encoding/json"
	"reflect"
)

func main() {
	testNetwork()
}

func testNetwork(){

	eyesCli := eyescli.NewLocalClient()
	chainID := "test_chain_id"
	datastApp := datastreet.NewDataStreet(eyesCli)

	datastApp.SetOption("base/chainID", chainID)
	fmt.Println(datastApp.Info())

	kvz := loadGenesis("/Users/Park-jihun/Desktop/1_BlockChain/5_Tendermint/workspace/testnet_basecoin/bctest/app/genesis.json")
	for _, kv := range kvz {
		log := datastApp.SetOption(kv.Key, kv.Value)
		fmt.Println(Fmt("Set %v=%v. Log: %v", kv.Key, kv.Value, log))
	}

	// Get the root account
	root := test.PrivateAccountFromSecret("test")
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

	queryResultTest := datastApp.Query(queryBytes)
	fmt.Println("\n", queryResultTest)

	var queryAccTest *types.Account
	err := wire.ReadBinaryBytes( queryResultTest.Data, &queryAccTest)
	if err != nil {
		fmt.Println("Err in ReadBinaryBytes.. ", err.Error())
	}
	fmt.Printf("Test's Account : %X\n", queryAccTest.PubKey)
	fmt.Printf("Balance 1 : %v\n", queryAccTest.Balance)
	fmt.Printf("Sequece 1 : %v", queryAccTest.Sequence)

	var queryAcc *types.Account
	err = wire.ReadBinaryBytes (queryResultTest.Data, &queryAcc)
	if err != nil {
		fmt.Println("err in ReadBinaryBytes..", err.Error())
	}
	fmt.Printf("Account 1 : %X\n", queryAcc.PubKey)
	fmt.Printf("Balance 1 : %v\n", queryAcc.Balance)
	fmt.Printf("Sequece 1 : %v", queryAcc.Sequence)

	// ======= Send coins to each account
	tx := &types.SendTx{
		Inputs: []types.TxInput {
			types.TxInput {
				Address: root.Account.PubKey.Address(),
				PubKey: root.Account.PubKey,
				Coins:	types.Coins {{"USD", 10}},
				Sequence: queryAccTest.Sequence + 1,
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

	res := datastApp.AppendTx(txBytes)
	fmt.Println(res)
	if res.IsErr() {
		Exit(Fmt("Failed :%v", res.Error()))
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