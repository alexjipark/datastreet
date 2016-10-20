package datastreet

import (
	"github.com/tendermint/go-merkle"
	tmsp "github.com/tendermint/tmsp/types"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/merkleeyes/client"
	"github.com/alexjipark/datastreet/types"
	"github.com/tendermint/go-wire"
	"fmt"
)

const (
	maxTxSize = 10240

	PluginTypeByteBase = 0x01
	PluginTypeByteEyes = 0x02
	PluginTypeByteGov  = 0x03
)

type DataStreetApp struct {
	state merkle.Tree
	eyesCli *eyes.Client
}


func NewDataStreetApp() *DataStreetApp {
	state := merkle.NewIAVLTree(0, nil, )
	return &DataStreetApp{state:state}
}

func NewDataStreet(eyesCli *eyes.Client) *DataStreetApp {
	state := merkle.NewIAVLTree(0, nil,)
	return &DataStreetApp {
		eyesCli: eyesCli,
		state: state,
	}

}

func (app *DataStreetApp) Info() string {
	return Fmt("DataStreet size:%v", app.state.Size())
}


func (app *DataStreetApp) SetOption(key string, value string) (log string) {
	return ""
}

// tx is either "key=value" or just arbitrary bytes
// basecoin - TMSP::AppendTx
func (app *DataStreetApp) AppendTx(txBytes []byte) tmsp.Result {

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

func (app *DataStreetApp) CheckTx(txBytes []byte) tmsp.Result {
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

	return tmsp.OK
}

func (app *DataStreetApp) Query(query []byte) tmsp.Result {

	if len(query) == 0 {
		return tmsp.ErrEncodingError.SetLog("Query cannot be zero length")
	}
	typeByte := query[0]
	query = query[1:]
	switch typeByte {
	case PluginTypeByteBase:
		return tmsp.OK.SetLog("This type of query not yet supported")
	case PluginTypeByteEyes:
		// Should be implemented soon..
		return tmsp.OK.SetLog("Ok but not yet implemented")
	}
	return tmsp.ErrBaseUnknownPlugin.SetLog(
		Fmt("Unknown plugin with type byte %X", typeByte))

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

}

// TMSP::EndBlock
func (app *DataStreetApp) EndBlock(height uint64) (diffs []*tmsp.Validator){
	// TBD .. should be implmented soon
	return
}

