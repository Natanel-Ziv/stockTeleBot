package model

type QuoteRequest struct {
    Symbol string `json:"symbol" in:"query=symbol;required"`
    Duration string `json:"duration" in:"query=duration;default=4w"`
}

type QuoteResponse struct {
    Symbol string `json:"symbol"`
    Min float64 `json:"min"`
    Max float64 `json:"max"`
    Graph string `json:"graph"`
}