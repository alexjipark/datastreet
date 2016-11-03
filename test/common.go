package test

import (
	"github.com/alexjipark/datastreet/types"
	"github.com/tendermint/go-crypto"
	"time"
	"github.com/tendermint/go-rpc/client"
	rpcctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/go-wire"
	"fmt"
)

// Create a PrivAccount from secret
// The amount is not set
func PrivateAccountFromSecret(secret string) types.PrivateAccount {
	privateKey := crypto.GenPrivKeyEd25519FromSecret([]byte(secret))
	privateAccount := types.PrivateAccount{
		PrivKey : privateKey,
		Account : types.Account{
			PubKey: privateKey.PubKey(),
			Sequence: 0,
		},
	}
	return privateAccount
}

var Fmt = fmt.Sprintf

// wait for an event; do things that might trigger events, and check them when they are received
// the check function takes an event id and the byte slice read off the ws
func WaitForEvent(wsc *rpcclient.WSClient, eventid string, dieOnTimeout bool, f func(), check func(string, interface{}) error) {
	// go routine to wait for webscoket msg
	goodCh := make(chan interface{})
	errCh := make(chan error)

	// Read message
	go func() {
		var err error
		LOOP:
		for {
			select {
			case r := <-wsc.ResultsCh:
				result := new(rpcctypes.TMResult)
				wire.ReadJSONPtr(result, r, &err)
				if err != nil {
					errCh <- err
					break LOOP
				}
				event, ok := (*result).(*rpcctypes.ResultEvent)
				if ok && event.Name == eventid {
					goodCh <- event.Data
					break LOOP
				}
			case err := <-wsc.ErrorsCh:
				errCh <- err
				break LOOP
			case <-wsc.Quit:
				break LOOP
			}
		}
	}()

	// do stuff (transactions)
	f()

	// wait for an event or timeout
	timeout := time.NewTimer(10 * time.Second)
	select {
	case <-timeout.C:
		if dieOnTimeout {
			wsc.Stop()
			panic(Fmt("%s event was not received in time", eventid))
		}
	// else that's great, we didn't hear the event
	// and we shouldn't have
	case eventData := <-goodCh:
		if dieOnTimeout {
			// message was received and expected
			// run the check
			if err := check(eventid, eventData); err != nil {
				panic(err) // Show the stack trace.
			}
		} else {
			wsc.Stop()
			panic(Fmt("%s event was not expected", eventid))
		}
	case err := <-errCh:
		panic(err) // Show the stack trace.

	}
}