package main

import (
	eyescli "github.com/tendermint/merkleeyes/client"
	"github.com/alexjipark/datastreet/app"
	"fmt"
	"github.com/alexjipark/datastreet/test"
	"github.com/alexjipark/datastreet/types"
	"github.com/tendermint/go-wire"
	. "github.com/tendermint/go-common"
	"encoding/json"
	"reflect"

)

func main() {
	testSendTx()
	testQuery()
}

func testQuery(){

}

func testSendTx() {

	/*
	testStr := "010000000000000000000000000000000001010114D9B727742AA29FA638DC63D70813C976014C4CE001010103555344000000000000000A010101F610525139E31087E36F8FD438F13BC16F477E99391298845A985C2D2E614B0C79CFD08661F4F3DA5EE3705FA0774C0F09EF1F47A9FFB7CD9146E4A43F6E69010167D3B5EAF0C0BF6B5A602D359DAECC86A7A74053490EC37AE08E71360587C870000101011412BB36B57DA6E4EC8229F4D99E14567F1E528B0F01010103555344000000000000000A"
	hexBytes, err := hex.DecodeString(testStr)
	fmt.Println(hexBytes)

	var curTx *types.Tx
	errtx := wire.ReadBinaryBytes(hexBytes, &curTx)
	if errtx != nil {
		fmt.Println(errtx.Error())
	}
	*/


	eyesCli := eyescli.NewLocalClient()
	chainID := "test_chain_id"
	datastApp := datastreet.NewDataStreet(eyesCli)

	datastApp.SetOption("base/chainID", chainID)
	fmt.Println(datastApp.Info())

	test1PrivateAcc := test.PrivateAccountFromSecret("test")
	test2PrivateAcc := test.PrivateAccountFromSecret("test1")

	kvz := loadGenesis("/Users/Park-jihun/Desktop/1_BlockChain/5_Tendermint/workspace/testnet_basecoin/bctest/app/genesis.json")
	for _, kv := range kvz {
		log := datastApp.SetOption(kv.Key, kv.Value)
		fmt.Println(Fmt("Set %v=%v. Log: %v", kv.Key, kv.Value, log))
	}
/*
	// Seed DataStreetApp with account
	test1Acc := test1PrivateAcc.Account
	test1Acc.Balance = types.Coins{{"USD",1000}}
	fmt.Println(datastApp.SetOption("base/account", string(wire.JSONBytes(test1Acc))))
*/
	// Construct a SendTx signature
	tx := &types.SendTx{
		Fee: 0,
		Gas: 0,
		Inputs : []types.TxInput {
			types.TxInput{
				Address: test1PrivateAcc.Account.PubKey.Address(),
				PubKey:  test1PrivateAcc.Account.PubKey,
				Coins: types.Coins{{"KRW",10}},
				Sequence: 1,
			},
		},
		Outputs: []types.TxOutput {
			types.TxOutput{
				Address: test2PrivateAcc.Account.PubKey.Address(),
				Coins: types.Coins{{"KRW",10}},
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
	fmt.Println("\nAppendTx :",res)
	if res.IsErr() {
		Exit(Fmt("Failed :%v", res.Error()))
	}

	res = datastApp.Commit()
	fmt.Println("\nCommit: ", res)

	//==== Query ====//
	addrTest1 := test1PrivateAcc.Account.PubKey.Address()
	queryTest1 := make([]byte, 1+wire.ByteSliceSize(addrTest1))
	buf := queryTest1
	buf[0] = 0x01
	buf = buf[1:]
	wire.PutByteSlice(buf, addrTest1)

	queryResultTest1 := datastApp.Query(queryTest1)
	fmt.Println("\n", queryResultTest1)

	var queryAcc1 *types.Account
	err := wire.ReadBinaryBytes (queryResultTest1.Data, &queryAcc1)
	if err != nil {
		fmt.Println("err in ReadBinaryBytes..", err.Error())
	}
	fmt.Printf("Account 1 : %X\n", queryAcc1.PubKey)
	fmt.Printf("Balance 1 : %v\n", queryAcc1.Balance)
	fmt.Printf("Sequece 1 : %v", queryAcc1.Sequence)


	addrTest2 := test2PrivateAcc.Account.PubKey.Address()
	queryTest2 := make([]byte, 1+wire.ByteSliceSize(addrTest2))
	buf2 := queryTest2
	buf2[0] = 0x01
	buf2 = buf2[1:]
	wire.PutByteSlice(buf2, addrTest2)

	queryResultTest2 := datastApp.Query(queryTest2)
	fmt.Println("\n", queryResultTest2)

	var queryAcc2 *types.Account
	err = wire.ReadBinaryBytes (queryResultTest2.Data, &queryAcc2)
	if err != nil {
		fmt.Println("err in ReadBinaryBytes..", err.Error())
	}
	fmt.Printf("Account 2 : %X\n", queryAcc2.PubKey)
	fmt.Printf("Balance 2 : %v\n", queryAcc2.Balance)
	fmt.Printf("Sequece 2 : %v", queryAcc2.Sequence)

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