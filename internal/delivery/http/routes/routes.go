package routes

import (
	"net/http"

	"github.com/Natanel-Ziv/stockTeleBot/internal/delivery/http/handlers"
	"github.com/Natanel-Ziv/stockTeleBot/internal/model"
	"github.com/ggicci/httpin"
	"github.com/justinas/alice"
	"github.com/sirupsen/logrus"
)

func RegisterRoute(mux *http.ServeMux, logger *logrus.Logger, stockHandler *handlers.StockHandler) {
    mux.Handle("/", http.NotFoundHandler())

	mux.Handle("/get-min-max-graph", alice.New(httpin.NewInput(model.QuoteRequest{})).ThenFunc(stockHandler.GetMinMaxGraphForDuration()))
}