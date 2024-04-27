package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Natanel-Ziv/stockTeleBot/internal/model"
	"github.com/Natanel-Ziv/stockTeleBot/internal/services/stock"
	"github.com/ggicci/httpin"
	"github.com/sirupsen/logrus"
)

type StockHandler struct {
	stockService stock.IStockService
	logger       *logrus.Logger
}

func NewStockHandler(logger *logrus.Logger) *StockHandler {
	return &StockHandler{
		stockService: stock.New(logger),
		logger:       logger,
	}
}

func (sh *StockHandler) GetMinMaxGraphForDuration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := r.Context().Value(httpin.Input).(*model.QuoteRequest)
		sh.logger.Info("got request with: ", input)

		resp, err := sh.stockService.GetMinMaxGraphForDuration(input.Symbol, input.Duration)
		if err != nil {
			sh.logger.WithError(err).Error("error getting min max graph: %+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		res := &model.QuoteResponse{
			Symbol: "aapl",
			Min:    resp.Min,
			Max:    resp.Max,
			Graph:  resp.Graph,
		}

		resData, err := json.Marshal(res)
		if err != nil {
			sh.logger.WithError(err).Error("error marshaling res")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resData)
	}
}

