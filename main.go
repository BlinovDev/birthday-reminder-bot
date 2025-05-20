package main

import (
	"fmt"
	"log"

	"birthday-reminder-bot/commands_handler"
	"birthday-reminder-bot/reminder_sender"

	tgconfighelper "github.com/BlinovDev/go-tg-bot-config-helper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, new_update, cron_time_pattern := prepare_bot()
	fmt.Println("Bot prepared!")

	reminder_sender.Setup_cron_reminder(bot, cron_time_pattern)
	fmt.Println("Cron set up!")

	updates := bot.GetUpdatesChan(new_update)
	fmt.Println("Read updates...")

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Text {
			case "/start":
				commands_handler.HandleStart(&bot, update)
			case "/add_new_birthday":
				commands_handler.HandleNewBirthday(&bot, update)
			case "/show_saved_birthdays":
				commands_handler.HandleViewBirthdays(&bot, update)
			case "/update_birthday":
				commands_handler.HandleUpdateBirthday(&bot, update)
			case "/delete_birthday":
				commands_handler.HandleDeleteBirthday(&bot, update)
			default:
				commands_handler.HandleAnswerMessage(&bot, update)
			}
		}
	}
}

func prepare_bot() (tgbotapi.BotAPI, tgbotapi.UpdateConfig, string) {
	// Read configs from specified file
	config, err := tgconfighelper.LoadConfig()
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

	cron_time_pattern := config.Bot.FirstCron

	return *bot, new_update, cron_time_pattern
}
