package state

import (

	bctypes "github.com/tendermint/basecoin/types"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-events"
	tmsp "github.com/tendermint/tmsp/types"
	"github.com/alexjipark/datastreet/types"

	//"github.com/tendermint/basecoin/state"
	//"log"
	"fmt"
	"bytes"
)

// If the tx is invalid, a TMSP error will be returned.
func ExecTx(state *State, pgz *bctypes.Plugins, tx bctypes.Tx, isCheckTx bool, evc events.Fireable) tmsp.Result {

	// TODO: do something with fees
	fees := types.Coins{}
	chainID := state.GetChainID()

	// Exec tx
	switch tx := tx.(type) {
	case *types.SendTx:
		// Validate inputs and outputs, basic
		res := validateInputsBasic(tx.Inputs)		// To Be Implemented..
		if res.IsErr() {
			return res.PrependLog("in validateInputsBasic()")
		}
		res = validateOutputsBasic(tx.Outputs)		// To Be Implemented
		if res.IsErr() {
			return res.PrependLog("in validateOutputsBasic()")
		}

		// Get inputs
		accounts, res := getInputs(state, tx.Inputs)
		if res.IsErr() {
			return res.PrependLog("in getInputs()")
		}

		// alexjipark - if Tx.Input = Tx.Output.. it would be related to "Data Ownership"
		// Get or make outputs.
		accounts, res = getOrMakeOutputs(state, accounts, tx.Outputs)
		if res.IsErr() {
			return res.PrependLog("in getOrMakeOutputs()")
		}

		// Validate inputs and outputs, advanced
		signBytes := tx.SignBytes(chainID)
		inTotal, res := validateInputsAdvanced(accounts, signBytes, tx.Inputs)
		if res.IsErr() {
			return res.PrependLog("in validateInputsAdvanced()")
		}
		outTotal := sumOutputs(tx.Outputs)

		// Have to be Resolved.. alexjipark
		if tx.Fee != 0 {
			if !inTotal.IsEqual(outTotal.Plus(types.Coins{{"", tx.Fee}})) {
				return tmsp.ErrBaseInvalidOutput.AppendLog("Input total != output total + fees")
			}
			fees = fees.Plus(types.Coins{{"", tx.Fee}})
		}

		// TODO: Fee validation for SendTx

		// Good! Adjust accounts
		adjustByInputs(state, accounts, tx.Inputs)
		adjustByOutputs(state, accounts, tx.Outputs, isCheckTx)

		/*
			// Fire events
			if !isCheckTx {
				if evc != nil {
					for _, i := range tx.Inputs {
						evc.FireEvent(types.EventStringAccInput(i.Address), types.EventDataTx{tx, nil, ""})
					}
					for _, o := range tx.Outputs {
						evc.FireEvent(types.EventStringAccOutput(o.Address), types.EventDataTx{tx, nil, ""})
					}
				}
			}
		*/

		return tmsp.OK
/*
	case *types.AppTx:
		// Validate input, basic
		res := tx.Input.ValidateBasic()
		if res.IsErr() {
			return res
		}

		// Get input account
		inAcc := state.GetAccount(tx.Input.Address)
		if inAcc == nil {
			return tmsp.ErrBaseUnknownAddress
		}
		if tx.Input.PubKey != nil {
			inAcc.PubKey = tx.Input.PubKey
		}

		// Validate input, advanced
		signBytes := tx.SignBytes(chainID)
		res = validateInputAdvanced(inAcc, signBytes, tx.Input)
		if res.IsErr() {
			//log.Info(Fmt("validateInputAdvanced failed on %X: %v", tx.Input.Address, res))
			return res.PrependLog("in validateInputAdvanced()")
		}
		if !tx.Input.Coins.IsGTE(types.Coins{{"", tx.Fee}}) {
			//log.Info(Fmt("Sender did not send enough to cover the fee %X", tx.Input.Address))
			return tmsp.ErrBaseInsufficientFunds
		}

		// Validate call address
		plugin := pgz.GetByByte(tx.Type)
		if plugin == nil {
			return tmsp.ErrBaseUnknownAddress.AppendLog(
				Fmt("Unrecognized type byte %v", tx.Type))
		}

		// Good!
		coins := tx.Input.Coins.Minus(types.Coins{{"", tx.Fee}})
		inAcc.Sequence += 1
		inAcc.Balance = inAcc.Balance.Minus(tx.Input.Coins)

		// If this is a CheckTx, stop now.
		if isCheckTx {
			state.SetAccount(tx.Input.Address, inAcc)
			return tmsp.OK
		}

		// Create inAcc checkpoint
		inAccDeducted := inAcc.Copy()

		// Run the tx.
		// XXX cache := types.NewStateCache(state)
		cache := state.CacheWrap()
		cache.SetAccount(tx.Input.Address, inAcc)
		ctx := types.NewCallContext(tx.Input.Address, coins)
		res = plugin.RunTx(cache, ctx, tx.Data)
		if res.IsOK() {
			cache.CacheSync()
			//log.Info("Successful execution")
		} else {
			//log.Info("AppTx failed", "error", res)
			// Just return the coins and return.
			inAccDeducted.Balance = inAccDeducted.Balance.Plus(coins)
			// But take the gas
			// TODO
			state.SetAccount(tx.Input.Address, inAccDeducted)
		}
		return res
*/
	default:
		return tmsp.ErrBaseEncodingError.SetLog("Unknown tx type")
	}
}

//--------------------------------------------------------------------------------

// The accounts from the TxInputs must either already have
// crypto.PubKey.(type) != nil, (it must be known),
// or it must be specified in the TxInput.
func getInputs(state types.AccountGetter, ins []types.TxInput) (map[string]*types.Account, tmsp.Result) {
	accounts := map[string]*types.Account{}
	for _, in := range ins {
		// Account shouldn't be duplicated
		if _, ok := accounts[string(in.Address)]; ok {
			return nil, tmsp.ErrBaseDuplicateAddress
		}
		acc := state.GetAccount(in.Address)
		if acc == nil {
			return nil, tmsp.ErrBaseUnknownAddress
		}
		if in.PubKey != nil {
			acc.PubKey = in.PubKey
		}
		accounts[string(in.Address)] = acc
	}
	return accounts, tmsp.OK
}

func getOrMakeOutputs(state types.AccountGetter, accounts map[string]*types.Account, outs []types.TxOutput) (map[string]*types.Account, tmsp.Result) {
	if accounts == nil {
		accounts = make(map[string]*types.Account)
	}

	for _, out := range outs {
		// Account shouldn't be duplicated
		if _, ok := accounts[string(out.Address)]; ok {
			return nil, tmsp.ErrBaseDuplicateAddress
		}
		acc := state.GetAccount(out.Address)
		// output account may be nil (new)
		if acc == nil {
			acc = &types.Account{
				PubKey:   nil,
				Sequence: 0,
			}
		}
		accounts[string(out.Address)] = acc
	}
	return accounts, tmsp.OK
}

// Validate inputs basic structure
func validateInputsBasic(ins []types.TxInput) (res tmsp.Result) {
	for _, in := range ins {
		// Check TxInput basic
		if res := in.ValidateBasic(); res.IsErr() {
			return res
		}
	}
	return tmsp.OK
}

// Validate inputs and compute total amount of coins
func validateInputsAdvanced(accounts map[string]*types.Account, signBytes []byte, ins []types.TxInput) (total types.Coins, res tmsp.Result) {
	for _, in := range ins {
		acc := accounts[string(in.Address)]
		if acc == nil {
			PanicSanity("validateInputsAdvanced() expects account in accounts")
		}
		res = validateInputAdvanced(acc, signBytes, in)
		if res.IsErr() {
			return
		}
		// Good. Add amount to total
		total = total.Plus(in.Coins)
	}
	return total, tmsp.OK
}

func validateInputAdvanced(acc *types.Account, signBytes []byte, in types.TxInput) (res tmsp.Result) {
	// Check sequence/coins
	seq, balance := acc.Sequence, acc.Balance
	if seq+1 != in.Sequence {
		fmt.Printf("Err..Got %v, expected %v. (acc.seq=%v)", in.Sequence, seq+1, acc.Sequence)
		//return tmsp.ErrBaseInvalidSequence.AppendLog(Fmt("Got %v, expected %v. (acc.seq=%v)", in.Sequence, seq+1, acc.Sequence))
	}
	// Check amount
	if !balance.IsGTE(in.Coins) {
		return tmsp.ErrBaseInsufficientFunds
	}
	// Check signatures
	/* ================= Have to be Resolved.. ===========//
	if !acc.PubKey.VerifyBytes(signBytes, in.Signature) {
		return tmsp.ErrBaseInvalidSignature.AppendLog(Fmt("SignBytes: %X", signBytes))
	}
	*/
	return tmsp.OK
}

func validateOutputsBasic(outs []types.TxOutput) (res tmsp.Result) {
	for _, out := range outs {
		// Check TxOutput basic
		if res := out.ValidateBasic(); res.IsErr() {
			return res
		}
	}
	return tmsp.OK
}

// alexjipark - Temp code. should be Re-Implemented with data.go
func checkRequestForDataOwnership(ins []types.TxInput, outs []types.TxOutput) (res tmsp.Result) {
	if len(ins) != 1 {
		res = tmsp.ErrInternalError
	}

	if len(outs) != 1 {
		res = tmsp.ErrInternalError
	}

	if bytes.Compare( ins[0].Address, outs[0].Address) != 0 {
		res = tmsp.OK
	}

	return res
}


func sumOutputs(outs []types.TxOutput) (total types.Coins) {
	for _, out := range outs {
		total = total.Plus(out.Coins)
	}
	return total
}

func adjustByInputs(state types.AccountSetter, accounts map[string]*types.Account, ins []types.TxInput) {
	for _, in := range ins {
		acc := accounts[string(in.Address)]
		if acc == nil {
			PanicSanity("adjustByInputs() expects account in accounts")
		}
		if !acc.Balance.IsGTE(in.Coins) {
			PanicSanity("adjustByInputs() expects sufficient funds")
		}
		acc.Balance = acc.Balance.Minus(in.Coins)
		acc.Sequence += 1
		state.SetAccount(in.Address, acc)
	}
}

func adjustByOutputs(state types.AccountSetter, accounts map[string]*types.Account, outs []types.TxOutput, isCheckTx bool) {
	for _, out := range outs {
		acc := accounts[string(out.Address)]
		if acc == nil {
			PanicSanity("adjustByOutputs() expects account in accounts")
		}
		acc.Balance = acc.Balance.Plus(out.Coins)
		if !isCheckTx {
			state.SetAccount(out.Address, acc)
		}
	}
}


