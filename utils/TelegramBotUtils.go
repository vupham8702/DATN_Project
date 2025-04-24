package utils

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

var bt *bot.Bot
var ctx context.Context

func CreateBot(contx context.Context, botToken string) {
	b, err := bot.New(botToken, bot.WithDefaultHandler(handle))
	if err != nil {
		panic(err)
	}
	bt = b
	ctx = contx

	go bt.Start(ctx)
}

func SendMessage(chatID int64, text string) {
	//bt.SendMessage(context.Background(), chatID, text)
	bt.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
}
func SendMessageTopic(chatID int64, threadId int, text string) {
	//bt.SendMessage(context.Background(), chatID, text)
	bt.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:          chatID,
		MessageThreadID: threadId,
		Text:            text,
	})
}

func handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Check if the update contains a message
	if update.Message != nil {
		chat := update.Message.Chat
		fmt.Printf("Chat ID: %d, Type: %s, Title: %s\n", chat.ID, chat.Type, chat.Title)

		// Example: Check if it's a supergroup
		if chat.Type == "supergroup" {
			fmt.Printf("Supergroup ID: %d\n", chat.ID)
		}
	}
}
