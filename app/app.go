package datastreet

import (
	"github.com/tendermint/go-merkle"
	tmsp "github.com/tendermint/tmsp/types"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/merkleeyes/client"
	"github.com/alexjipark/datastreet/types"
	"github.com/tendermint/go-wire"
	"fmt"
	"strings"
	sm "github.com/alexjipark/datastreet/state"
)

const (
	maxTxSize = 10240

	PluginTypeByteBase = 0x01
	PluginTypeByteEyes = 0x02
	//PluginTypeByteGov  = 0x03

	PluginNameBase = "base"
	PluginNameEyes = "eyes"
	//PluginNameGov  = "gov"
)

type DataStreetApp struct {
	state 		merkle.Tree
	bcstate 	*sm.State
	cacheState	*sm.State
	eyesCli 	*eyes.Client
}


func NewDataStreetApp() *DataStreetApp {
	state := merkle.NewIAVLTree(0, nil, )
	return &DataStreetApp{state:state}

}

func NewDataStreet(eyesCli *eyes.Client) *DataStreetApp {
	state := merkle.NewIAVLTree(0, nil,)
	bcstate := sm.NewState(eyesCli)
	return &DataStreetApp {
		eyesCli: eyesCli,
		state: state,
		bcstate: bcstate,
		cacheState: nil,
	}

}

func (app *DataStreetApp) Info() string {
	return Fmt("DataStreet size:%v", app.state.Size())
}


func (app *DataStreetApp) SetOption(key string, value string) (log string) {

	PluginName, key := splitKey(key)
	if PluginName != PluginNameBase {
		//Set option on plugin
		// TBD.. To be developed soon!
	} else {
		//Set option on Basecoin
		switch key {
		case "chainID":
			app.bcstate.SetChainID(value)
			return "Success"
		case "account":
			var err error
			var acc *types.Account
			wire.ReadJSONPtr(&acc, []byte(value), &err)
			if err != nil {
				return "Error decoding acc message: " + err.Error()
			}

			//====== Check Account =====
			fmt.Printf("Public Key  : %X\n", acc.PubKey.Address())
			fmt.Printf("Public Byte : %X\n", acc.PubKey.Bytes())
			//==========================

			app.bcstate.SetAccount(acc.PubKey.Address(), acc)
			return "Success"
		}
	}
	return "Unrecoginzed option key " + key

}

// tx is either "key=value" or just arbitrary bytes
// basecoin - TMSP::AppendTx
func (app *DataStreetApp) AppendTx(txBytes []byte) (res tmsp.Result) {

	// basecoin
	if len(txBytes) > maxTxSize {
		return tmsp.ErrBaseEncodingError.AppendLog("Tx size exceeds maximum")
	}

	// basecoin - Decode tx
	var tx types.Tx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return tmsp.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
	}

	fmt.Println("About to validate and exec..")

	// Validate and Exec Tx
	// TBD .. should be implemented
	res = sm.ExecTx(app.bcstate, nil, tx,false, nil)	// plugin is not used in SendTx in sm
								// But it's utilized with AppTx (Smart Contract)
	if  res.IsErr() {
		return res.PrependLog("Error in AppendTx")
	}
/*
	parts := strings.Split(string(tx),"=")
	if len(parts) == 2 {
		app.state.Set([]byte(parts[0]), []byte(parts[1]))
	} else {
		app.state.Set (tx, tx)
	}
*/

	return tmsp.OK
}

func (app *DataStreetApp) CheckTx(txBytes []byte) (res tmsp.Result) {
	if len(txBytes) > maxTxSize {
		return tmsp.ErrBaseEncodingError.AppendLog("Tx size exceeds maximum")
	}

	// Decode tx
	var tx types.Tx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return tmsp.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
	}

	// Validate tx
	// TBD .. should be implemented
	res = sm.ExecTx(app.cacheState, nil, tx, true, nil)

	if res.IsErr() {
		return res.PrependLog("Error in CheckTx")
	}

	return tmsp.OK
}

func (app *DataStreetApp) Query(query []byte) tmsp.Result {

	if len(query) == 0 {
		return tmsp.ErrEncodingError.SetLog("Query cannot be zero length")
	}
	return app.eyesCli.QuerySync(query)

	/*
	typeByte := query[0]
	query = query[1:]
	switch typeByte {
	case PluginTypeByteBase:
		return tmsp.OK.SetLog("This type of query not yet supported")
	case PluginTypeByteEyes:
		// Should be implemented soon..

		return app.eyesCli.QuerySync(query)
		//return tmsp.OK.SetLog("Ok but not yet implemented")
	}
	return tmsp.ErrBaseUnknownPlugin.SetLog(
		Fmt("Unknown plugin with type byte %X", typeByte))
	*/
/*
	index, value, exists := app.state.Get(query)
	resStr := Fmt("Index=%v value=%v exists=%v", index, string(value), exists)
	return tmsp.NewResultOK([]byte(resStr), "")
*/
}

// TMSP::Commit
func (app *DataStreetApp) Commit() tmsp.Result {

	//Commit eyes
	res := app.eyesCli.CommitSync()
	if res.IsErr() {
		PanicSanity("Error Getting Hash: " + res.Error())
	}
	return res
/*
	hash := app.state.Hash()
	return tmsp.NewResultOK(hash,"")
*/
}

// TMSP::InitChain
func (app *DataStreetApp) InitChain(validators []*tmsp.Validator) {
	// TBD .. should be implemented soon

}

// TMSP::BeginBlock
func (app *DataStreetApp) BeginBlock(height uint64) {
	// TBD .. should be implemented soon
	app.cacheState = app.bcstate.CacheWrap()
}

// TMSP::EndBlock
func (app *DataStreetApp) EndBlock(height uint64) (diffs []*tmsp.Validator){
	// TBD .. should be implmented soon
	return
}


//--------------------
// splits the string at the first '/'
// if there are none, the second string is nil
func splitKey(key string) (prefix string, suffix string) {
	if strings.Contains(key, "/") {
		keyParts := strings.SplitN(key, "/", 2)
		return keyParts[0], keyParts[1]
	}
	return key, ""
}