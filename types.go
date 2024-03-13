package main

import (
	"github.com/gofrs/uuid/v5"
	"github.com/shopspring/decimal"
)

type Movie struct {
	Id          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	ReleaseYear int             `json:"release_year"`
	Rating      decimal.Decimal `json:"rating"`
	Genres      []string        `json:"genres"`
	Director    string          `json:"director"`
}
