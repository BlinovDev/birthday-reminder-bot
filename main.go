package main

import (
	"log"

	"birthday-reminder-bot/commands_handler"
	"birthday-reminder-bot/config_helper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, new_update := prepare_bot()

	updates := bot.GetUpdatesChan(new_update)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				commands_handler.HandleStart(&bot, update)
			case "/add_new_birthday":
				commands_handler.HandleNewBirthday(&bot, update)
			case "/show_saved_birthdays":
				commands_handler.HandleViewBirthdays(&bot, update)
			default:
				commands_handler.HandleAnswerMessage(&bot, update)
			}
		}
	}
}

func prepare_bot() (tgbotapi.BotAPI, tgbotapi.UpdateConfig) {
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

	return *bot, new_update
}
