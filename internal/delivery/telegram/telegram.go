package telegram

import (
	"context"

	"github.com/Natanel-Ziv/stockTeleBot/internal/delivery/telegram/handlers"
	"github.com/Natanel-Ziv/stockTeleBot/internal/delivery/telegram/middlewares"
	"github.com/Natanel-Ziv/stockTeleBot/internal/services/stock"
	"github.com/go-telegram/bot"
	"github.com/sirupsen/logrus"
)

type ITelegramBot interface {
    RegisterHandlers(stockService stock.IStockService)
	Start(ctx context.Context)
}

type telegramBot struct {
    bot *bot.Bot
	logger *logrus.Logger
}

func (t *telegramBot) Start(ctx context.Context) {
    t.bot.Start(ctx)
}

func (t *telegramBot) RegisterHandlers(stockService stock.IStockService) {
    stockHandler := handlers.NewStockHandler(stockService)

    t.bot.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handlers.HelpHandler)
    t.bot.RegisterHandler(bot.HandlerTypeMessageText, "/symbol", bot.MatchTypePrefix, stockHandler.SymbolHandler)
}

func New(token string, logger *logrus.Logger) (ITelegramBot, error) {
    loggingMiddleware := middlewares.New(logger)
    opts := []bot.Option{
        bot.WithMiddlewares(loggingMiddleware.LogMessage),
        bot.WithDefaultHandler(handlers.DefaultHandler),
    }

    b, err := bot.New(token, opts...)
    if err != nil {
        return nil, err
    }

	return &telegramBot{
        bot: b,
		logger: logger,
	}, nil
}
