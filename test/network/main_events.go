package main

import (
	"github.com/tendermint/go-rpc/client"
	. "github.com/tendermint/go-common"

	tmsp "github.com/tendermint/tmsp/types"

	tdtypes "github.com/tendermint/tendermint/types"

	"fmt"
	test 	"github.com/alexjipark/datastreet/test"
)

type ResultData struct {
	Result tmsp.Result `json:"result"`
}

func main() {
	ws := rpcclient.NewWSClient("35.160.145.128:46657", "/websocket")
	//chainID := "chain-AMUKE0"

	_,err := ws.Start()
	if err != nil {
		Exit(err.Error())
	}

	defer func() {
		ws.Close()
	}()

	event_id := tdtypes.EventStringNewBlock()
	if err = ws.Subscribe(event_id) ; err != nil {
		fmt.Println("Subscribe Error.. ", err.Error())
	}

	defer func() {
		ws.Unsubscribe(event_id)
	}()

	// waitForEvent(wsc, eid, true, func() {}, func(eid string, eventData interface{}) error {
	for i := 0 ; i < 3; i++ {
		test.WaitForEvent(ws, event_id, true, func(){}, func(eid string, eventData interface{}) error {

			block := eventData.(tdtypes.EventDataNewBlock).Block
			fmt.Println("\nReceived Block Height.. : ", block.Header.Height)

			return nil
		})
	}


}
