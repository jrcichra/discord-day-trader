package stock

import (
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/quote"
)

//Stock -
type Stock struct {
	ticker    string
	lastQuote *finance.Quote
}

//New - make a new stock
func (s *Stock) New(ticker string) {
	s.ticker = ticker
}

//Get a quote from a ticker
func (s *Stock) Get(ticker string) *finance.Quote {
	q, err := quote.Get(ticker)
	if err != nil {
		panic(err)
	}
	s.lastQuote = q
	return q
}
