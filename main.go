package main

import (
	// "fmt"
	"log"

	// "birthday-reminder-bot/birthdays_helper"
	"birthday-reminder-bot/config_helper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, _, new_update := prepare_bot()

	updates := bot.GetUpdatesChan(new_update)

	for update := range updates {
		if update.Message != nil {
			incoming_message := update.Message

			incoming_text := incoming_message.Text
			incoming_chat := incoming_message.Chat.ID

			sendMessage(&bot, incoming_chat, incoming_text)
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return
	}
}

// func deleteMessage(bot *tgbotapi.BotAPI, chatID int64, messageID int) {
// 	delMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
// 	if _, err := bot.Request(delMsg); err != nil {
// 		error_msg := fmt.Sprintf("Failed to delete message (ChatID: %d, MessageID: %d): %v", chatID, messageID, err)
// 		fmt.Println("ERROR", error_msg)
// 	}
// }

func prepare_bot() (tgbotapi.BotAPI, config_helper.Config, tgbotapi.UpdateConfig) {
	// Read configs from specified file
	config, err := config_helper.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialise bot instance
	bot, err := tgbotapi.NewBotAPI(config.Bot.Token)
	if err != nil {
		log.Panic(err)
	}

	new_update := tgbotapi.NewUpdate(0)
	new_update.Timeout = 60

	return *bot, *config, new_update
}