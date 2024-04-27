package main

import (
	"cmp"
	"context"
	"net"
	"net/http"
	"time"

	"github.com/Natanel-Ziv/stockTeleBot/config"
	"github.com/Natanel-Ziv/stockTeleBot/internal/delivery/http/server"
	"github.com/Natanel-Ziv/stockTeleBot/internal/infrastructure"
)

func main() {
    appCtx := infrastructure.AppContext
    config := config.New()

	port := config.GetString("APP_PORT")
   
    logger := infrastructure.NewLogger(cmp.Or(config.GetString("LOG_LEVEL"), "info"))
    srv := server.New(logger)
    httpServer := &http.Server{
        Addr: net.JoinHostPort("0.0.0.0", port),
        Handler: srv,
    }

    go func() {
        logger.Infof("listening on %s", httpServer.Addr)
        if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Panicf("error listening and serving: %s", err.Error())
        }
    }()

    <-appCtx.Done()
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
    defer cancel()
    if err := httpServer.Shutdown(shutdownCtx); err != nil {
        logger.Errorf("error shutting down http server: %s", err.Error())
    }
}
