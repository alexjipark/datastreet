package main

import (
	"flag"
	"github.com/tendermint/tmsp/server"
	. "github.com/tendermint/go-common"
	"github.com/alexjipark/datastreet/app"
	"github.com/tendermint/merkleeyes/client"
	"fmt"
	"encoding/json"
	"reflect"
	/*
	"github.com/alexjipark/datastreet/test"
	"github.com/tendermint/go-wire"
	*/
)

func main() {
	addrPtr := flag.String("address", "tcp://0.0.0.0:46658", "Listen Address")
	tmspPtr := flag.String("tmsp", "socket", "socket | grpc")

	// basecoin
	eyesPtr := flag.String("eyes", "local", "MerkleEyes Address, or 'local' for embedded")
	genesisFilePath := flag.String("genesis", "/Users/Park-jihun/Desktop/1_BlockChain/5_Tendermint/workspace/testnet_basecoin/bctest/app/genesis.json", "Genesis File, if any")
	flag.Parse()

	// basecoin, connect to MerkleEyes
	eyesCli, err := eyes.NewClient(*eyesPtr, *tmspPtr)
	if err != nil {
		Exit("connect to MerkleEyes: " + err.Error())
	}

	// basecoin, Create DataStreet App
	dataStreet := datastreet.NewDataStreet(eyesCli)

	// basecoin, if genesis file was specified, set key-value options
	if *genesisFilePath != "" {
		kvz := loadGenesis(*genesisFilePath)
		for _, kv := range kvz {
			log := dataStreet.SetOption(kv.Key, kv.Value)
			fmt.Println(Fmt("Set %v=%v. Log: %v", kv.Key, kv.Value, log))
		}
	}

	//===== Test =====
	/*
	root := test.PrivateAccountFromSecret("test")

	//====== Check Account
	fmt.Printf("Private Key : %X\n", root.PrivKey)
	fmt.Printf("Public Byte : %X\n", root.Account.PubKey.Bytes())
	fmt.Printf("Public Addr : %X\n", root.Account.PubKey.Address())

	addrBytes := root.Account.PubKey.Address()
	queryBytes := make([]byte, 1+wire.ByteSliceSize(addrBytes))
	buf := queryBytes
	buf[0] = 0x01
	buf = buf[1:]
	wire.PutByteSlice(buf, addrBytes)

	//queryResult := dataStreet.Query(queryBytes))
	*/
	//================


	// Start the Listener
	svr, err := server.NewServer(*addrPtr, *tmspPtr, dataStreet)
	if err != nil {
		Exit("create Listener: " + err.Error())
	}

	// Wait Forever
	TrapSignal(func() {
		svr.Stop()
	})
}

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