package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)
var (
	userState = make(map[int64]string)
	catAPIkey string
)

type catImage struct {
	URL string `json:"url"`
}

func main() {
	token:= os.Getenv("TOKEN")
	catAPIkey:= os.Getenv("CAT")
	if token == "" || catAPIkey == "" {
		log.Fatal("TOKEN or CAT env var missing")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		panic(err)
	}

	_, err = b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: []models.BotCommand{
			// {Command: "start", Description: "Start the bot"},
			{Command: "cat", Description: "Recieve a cat picture"},
		},
	})
	if err != nil {
		log.Fatalf("SetMyCommands error: %v", err)
	}

	b.Start(ctx)
}



func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	chatID:= update.Message.Chat.ID
	text:=   update.Message.Text

	b.RegisterHandler(bot.HandlerTypeMessageText, "cat", bot.MatchTypeCommand, catHandler)

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
		} 
		if strings.EqualFold(answer, "yes") {
			fileContent, err := os.ReadFile("cat.png")
			if err != nil {
				catpic, err := getCatURL()
				if err != nil {
					return
				}
				imgData, err := getByteArr(catpic)
				if err != nil {
					return
				}
				fileContent = imgData
			}
			params := &bot.SendPhotoParams {
				ChatID: chatID,
				Photo: &models.InputFileUpload{Filename: "cat.png", Data: bytes.NewReader(fileContent)},
				Caption: "That's a cat 😎",
			}
			b.SendPhoto(ctx, params)
			if _, err := b.SendPhoto(ctx, params); err != nil {
				log.Printf("SendPhoto error: %v", err)
			}
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text: "Oh, okay :(",
			})
		}

		userState[chatID] = ""
	}

}


func catHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	url, err := getCatURL()
	if err != nil {
		log.Printf("cat api error: %v", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "Could not fetch a cat right now.",
		})
		return
	}

	params := &bot.SendPhotoParams{
		ChatID: chatID,
		Photo: &models.InputFileString{Data: url},
		Caption: "Here's your cat 😎",
	}
	if _, err := b.SendPhoto(ctx, params); err != nil {
		log.Printf("SendPhoto error: %v", err)
	}
}



func getCatURL() (string, error) {
	req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/images/search", nil)

	if err != nil {
		return "", err
	}
	req.Header.Set("x-api-key", catAPIkey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var images []catImage
	if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
		return "", err
	}
	if len(images) == 0 {
		return "", fmt.Errorf("no images returned")
	}
	return images[0].URL, nil
}



func getByteArr(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// random cat facts api
// commands