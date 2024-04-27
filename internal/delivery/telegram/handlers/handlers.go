package handlers

import (
	"context"
	"os"
	"strings"

	"github.com/Natanel-Ziv/stockTeleBot/internal/services/stock"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/wcharczuk/go-chart/v2"
)

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Welcome to Stock TeleBot :)\nSay /help for full command list",
	})
}

func HelpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
    b.SendMessage(ctx, &bot.SendMessageParams{
        ChatID: update.Message.Chat.ID,
        Text: "/symbol <symbol> - Get data about last 4 weeks",
    })
}

type stockHandler struct {
	stockService stock.IStockService
}

func NewStockHandler(stockService stock.IStockService) *stockHandler {
	return &stockHandler{stockService: stockService}
}

func (sh *stockHandler) SymbolHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	dataAfterCommand := strings.Split(update.Message.Text, "/symbol ")
	
	if len(dataAfterCommand) < 1{
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: "Symbol is missing! Usage:\n/symbol <symbol>",
		})
		return
	}

	if len(dataAfterCommand) > 2 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: "Too many args! Usage:\n/symbol <symbol>",
		})
		return
	}

	symbol := dataAfterCommand[len(dataAfterCommand)-1]
	stockData, err := sh.stockService.GetMinMaxGraphForDuration(symbol, "4w")
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: "Failed to retrieve data :(",
		})
		return
	}

	f, err := os.CreateTemp(".", "*")
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: "Something went wrong :(",
		})
		return
	}
	defer os.Remove(f.Name())

	err = stockData.Graph.Render(chart.PNG, f)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: "Something went wrong :(",
		})
		return
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text: "Something went wrong :(",
		})
		return
	}
	
	params := &bot.SendPhotoParams{
		ChatID:  update.Message.Chat.ID,
		Photo:   &models.InputFileUpload{Filename: "resp.png", Data: f},
	}
	defer f.Close()

    b.SendPhoto(ctx, params)
}