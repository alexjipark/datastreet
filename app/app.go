package datastreet

import (
	"github.com/tendermint/go-merkle"
	"github.com/tendermint/tmsp/types"
	. "github.com/tendermint/go-common"
	"strings"
)

type DataStreetApp struct {
	state merkle.Tree
}


func NewDataStreetApp() *DataStreetApp {
	state := merkle.NewIAVLTree(0, nil, )
	return &DataStreetApp{state:state}
}

func (app *DataStreetApp) Info() string {
	return Fmt("DataStreet size:%v", app.state.Size())
}


func (app *DataStreetApp) SetOption(key string, value string) (log string) {
	return ""
}

// tx is either "key=value" or just arbitrary bytes
func (app *DataStreetApp) AppendTx(tx []byte) types.Result {
	parts := strings.Split(string(tx),"=")
	if len(parts) == 2 {
		app.state.Set([]byte(parts[0]), []byte(parts[1]))
	} else {
		app.state.Set (tx, tx)
	}
	return types.OK
}

func (app *DataStreetApp) CheckTx(tx []byte) types.Result {
	return types.OK
}

func (app *DataStreetApp) Commit() types.Result {
	hash := app.state.Hash()
	return types.NewResultOK(hash,"")
}

func (app *DataStreetApp) Query(query []byte) types.Result {
	index, value, exists := app.state.Get(query)
	resStr := Fmt("Index=%v value=%v exists=%v", index, string(value), exists)
	return types.NewResultOK([]byte(resStr), "")
}
