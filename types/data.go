package types

import "fmt"

type DataOwnership struct {
	DataHash 	string		`json:"datahash"`
	Amount		int64		`json:"amount"`
}

func (data DataOwnership) String() string {
	return fmt.Sprintf("(%v %v)", data.DataHash, data.Amount)
}
