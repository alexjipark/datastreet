package types

import "fmt"

type Coin struct {
	Denom string 	`json:"denom"`
	Amount int64	`json:"amount"`
}

func (coin Coin) String() string {
	return fmt.Sprintf("(%v %v)", coin.Denom, coin.Amount)
}

type Coins []Coin

