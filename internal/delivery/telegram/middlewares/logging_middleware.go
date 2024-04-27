package middlewares

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/sirupsen/logrus"
)

type loggingMiddeleware struct {
    logger *logrus.Logger
}

func New(logger *logrus.Logger) *loggingMiddeleware {
    return &loggingMiddeleware{
        logger: logger,
    }
}

func (l *loggingMiddeleware) LogMessage(next bot.HandlerFunc) bot.HandlerFunc {
    return func(ctx context.Context, b *bot.Bot, update *models.Update) {
        if update.Message != nil {
            l.logger.WithFields(logrus.Fields{
                "fromID": update.Message.From.ID,
                "fromUser": update.Message.From.FirstName,
            }).Info(update.Message.Text)
        }
        next(ctx, b, update)
    }
}