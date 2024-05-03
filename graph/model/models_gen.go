// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Account struct {
	ID           string       `json:"id"`
	Raiting      string       `json:"raiting"`
	Borrows      []*Borrow    `json:"borrows"`
	Liquidations []*Liquidate `json:"liquidations"`
}

type AccountFilter struct {
	ID *string `json:"id,omitempty"`
}

type Borrow struct {
	AmountUsd string   `json:"amountUSD"`
	Account   *Account `json:"account"`
}

type Liquidate struct {
	AmountUsd  string   `json:"amountUSD"`
	Liquidatee *Account `json:"liquidatee"`
}

type Query struct {
}