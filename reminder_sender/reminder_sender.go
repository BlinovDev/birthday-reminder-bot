package reminder_sender

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"birthday-reminder-bot/birthdays_helper"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot_instanse tgbotapi.BotAPI

func Setup_cron_reminder(bot tgbotapi.BotAPI, cron_time_pattern string) {
	c := cron.New()

	// Schedule the task to run daily at midnight
	c.AddFunc(cron_time_pattern, notifyUpcomingBirthdays) // Runs every day at set time

	// Start the scheduler
	c.Start()

	bot_instanse = bot
}

func notifyUpcomingBirthdays() {
	birthdays, err := birthdays_helper.GetBirthdays()
	if err != nil {
		fmt.Println("Error reading birthdays:", err)
		return
	}

	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)

	for _, birthday := range birthdays {
		if isBirthdayTomorrow(birthday.Birthday, tomorrow) {
			sendBirthdayNotification(birthday)
		}
	}
}

func isBirthdayTomorrow(birthday time.Time, tomorrow time.Time) bool {
	return birthday.Month() == tomorrow.Month() && birthday.Day() == tomorrow.Day()
}

func sendBirthdayNotification(birthday birthdays_helper.Birthday) {
	msg := tgbotapi.NewMessage(int64(birthday.ChatID), "Reminder: Tomorrow is "+birthday.Name+"'s birthday!")
	_, err := bot_instanse.Send(msg)
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}
