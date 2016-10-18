package main

import (
	"flag"
	"github.com/tendermint/tmsp/server"
	"github.com/tendermint/go-common"
	"github.com/alexjipark/datastreet/app"
)

func main() {
	addrPtr := flag.String("addr", "tcp://0.0.0.0:46658", "Listen Address")
	tmspPtr := flag.String("tmsp", "socket", "socket | grpc")
	flag.Parse()

	// Start the Listener
	_, err := server.NewServer(*addrPtr, *tmspPtr, datastreet.NewDataStreetApp())
	if err != nil {
		//os.Exit(err.Error())
	}

	// Wait Forever
	common.TrapSignal(func() {

	})
}