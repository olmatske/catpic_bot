package main

import (
	"strings"
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
)

var userState = make(map[int64]string)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	token:= os.Getenv("TOKEN")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID:= update.Message.Chat.ID
	text:=   update.Message.Text

	switch userState[chatID] {
	case "":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text: "Hi! Would you like to recieve some cat pictures?",
		})
		userState[chatID] = "await.answer"

	case "await.answer":
		answer := text
		if !strings.EqualFold(answer, "yes") && !strings.EqualFold(answer, "no") {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text: "Please answer with 'Yes' or 'No'.",
			})
			return
		} else if strings.EqualFold(answer, "yes") {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text: "Yippie!!",
			})
		} else if strings.EqualFold(answer, "no") {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text: "Oh, okay :(",
			})
		}



		userState[chatID] = ""
	}

}
