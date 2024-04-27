package main

import (
	"cmp"

	"github.com/Natanel-Ziv/stockTeleBot/config"
	"github.com/Natanel-Ziv/stockTeleBot/internal/delivery/telegram"
	"github.com/Natanel-Ziv/stockTeleBot/internal/infrastructure"
	"github.com/Natanel-Ziv/stockTeleBot/internal/services/stock"
)

func main() {
	appCtx := infrastructure.AppContext
	config := config.New()

    logger := infrastructure.NewLogger(cmp.Or(config.GetString("LOG_LEVEL"), "info"))

    telegramBot, err := telegram.New(config.GetString("TG_TOKEN"), logger)
    if err != nil {
        logger.Panic(err)
    }

    stockService := stock.New(logger)
    telegramBot.RegisterHandlers(stockService)

    go telegramBot.Start(appCtx)
    logger.Info("Telegram bot running")
    <-appCtx.Done()
}