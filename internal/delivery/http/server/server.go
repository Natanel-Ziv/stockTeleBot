package server

import (
	"net/http"

	"github.com/Natanel-Ziv/stockTeleBot/internal/delivery/http/handlers"
	"github.com/Natanel-Ziv/stockTeleBot/internal/delivery/http/middlewares"
	"github.com/Natanel-Ziv/stockTeleBot/internal/delivery/http/routes"
	"github.com/sirupsen/logrus"
)

type Server struct {
    port   int
    logger *logrus.Logger
}

func New(logger *logrus.Logger) http.Handler {
    mux := http.NewServeMux()
    routes.RegisterRoute(mux, logger, handlers.NewStockHandler(logger))
    
    var handler http.Handler = mux
    handler = middlewares.WithLogging(handler, logger)
    return handler
}