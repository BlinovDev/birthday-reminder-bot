package main

import (
	// "fmt"
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
			message := update.Message
			chat := message.Chat.ID
			text := message.Text

			if text == "/start" {
				// Send a minimal message with the inline keyboard
				start_msg := "Welcome! Nice to see you in birthday-reminder-bot.\n\nHere you can store your relatives and friends birthdays and get reminders about their birthdays. Use keyboard to call main bot functions.\n\nEnjoy!"

				msg := tgbotapi.NewMessage(chat, start_msg)
				msg.ReplyMarkup = getPresetMessageKeyboard()
				bot.Send(msg)
			} else {
				sendMessage(&bot, chat, text)
			}
		}
	}
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	// TODO: process commands like /add_bd, /all_bd, /next_bd
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

func getPresetMessageKeyboard() tgbotapi.ReplyKeyboardMarkup {
	// Define the reply keyboard with preset messages
	buttons := [][]tgbotapi.KeyboardButton{
		{
			tgbotapi.NewKeyboardButton("Add birthday"),
			tgbotapi.NewKeyboardButton("Saved Birthdays"),
		},
		// { },
	}

	// Create and return the keyboard markup
	return tgbotapi.NewReplyKeyboard(buttons...)
}
