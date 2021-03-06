package types

import (
	"fmt"
	"strings"
	"bytes"
)

type Coin struct {
	Denom string 	`json:"denom"`
	Amount int64	`json:"amount"`
}

func (coin Coin) String() string {
	return fmt.Sprintf("(%v %v)", coin.Denom, coin.Amount)
}

type Coins []Coin

func (coins Coins) String() string {

	var buffer bytes.Buffer
	for _,pcoin := range coins {
		buffer.WriteString(pcoin.String())
	}
	return buffer.String()
	//return fmt.Sprintf("Coins(%v)", len(coins))
}

func (coinsA Coins) IsEqual(coinsB Coins) bool {
	if len(coinsA) != len(coinsB) {
		return false
	}
	for  i:=0; i <len(coinsA); i++ {
		if coinsA[i] != coinsB[i] {
			return false
		}
	}
	return true
}

func (coins Coins) Negative() Coins {
	res := make([]Coin, 0, len(coins))
	for _,coin := range coins {
		res = append(res, Coin{
			Denom: coin.Denom,
			Amount: -coin.Amount,
		})
	}
	return res
}

func (coinsA Coins) Plus(coinsB Coins) Coins {
	sum := []Coin{}
	indexA, indexB := 0,0
	lenA, lenB := len(coinsA), len(coinsB)
	for {
		if indexA == lenA {
			if indexB == lenB {
				return sum
			} else {
				return append(sum, coinsB[indexB:]...)
			}
		} else if indexB == lenB {
			return append(sum, coinsA[indexA:]...)
		}
		coinA, coinB := coinsA[indexA], coinsB[indexB]
		switch strings.Compare(coinA.Denom, coinB.Denom) {
		case -1:
			sum = append(sum, coinA)
			indexA += 1
		case 0:
			if coinA.Amount + coinB.Amount == 0 {
				// ignore 0 sum coin type
			} else {
				sum = append (sum, Coin {
					Denom: coinA.Denom,
					Amount: coinA.Amount + coinB.Amount,
				})
			}
			indexA += 1
			indexB += 1
		case 1:
			sum = append(sum, coinB)
			indexB +=1
		}
	}
	return sum
}

func (coinsA Coins) Minus (coinsB Coins) Coins {
	return coinsA.Plus(coinsB.Negative())
}

func (coins Coins) IsPositive() bool {
	if len(coins) == 0 {
		return false
	}
	for _, coinAmount := range coins {
		if coinAmount.Amount <= 0 {
			return false
		}
	}
	return true
}

func (coinsA Coins) IsGTE (coinsB Coins) bool {
	diff := coinsA.Minus (coinsB)
	if len(diff) == 0 {
		return true
	}

	return diff.IsPositive()
}
