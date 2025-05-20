package commands_handler

import (
	"birthday-reminder-bot/birthdays_helper"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var userState = make(map[int64]string)
var tempBirthdays = make(map[int64]birthdays_helper.Birthday)

// helper function to process name input
func handleNameInput(bot *tgbotapi.BotAPI, userID, chatID int64, text, nextState, prompt string) {
	birthday := tempBirthdays[userID]
	birthday.Name = text
	tempBirthdays[userID] = birthday
	msg := tgbotapi.NewMessage(chatID, prompt)
	bot.Send(msg)
	userState[userID] = nextState
}

// helper function to process birthday input
func handleBirthdayInput(bot *tgbotapi.BotAPI, userID, chatID int64, text, nextState, prompt string) {
	birthday := tempBirthdays[userID]
	bday, err := time.Parse("2006-01-02", text)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Invalid date format. Please enter the birthday in YYYY-MM-DD format:")
		bot.Send(msg)
		return
	}

	birthday.Birthday = bday
	tempBirthdays[userID] = birthday
	msg := tgbotapi.NewMessage(chatID, prompt)
	bot.Send(msg)
	userState[userID] = nextState
}

// helper function to process telegram name and finalise saving
func handleTgNameInput(bot *tgbotapi.BotAPI, userID, chatID int64, text, nextState, prompt, action string) {
	birthday := tempBirthdays[userID]
	if text == "skip" {
		birthday.TgName = ""
	} else {
		birthday.TgName = text
	}

	err := birthdays_helper.AddBirthday(birthday.Name, birthday.TgName, birthday.Birthday, birthday.ChatID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "An error occurred while saving the birthday.")
		bot.Send(msg)
	} else {
		successMsg := "Birthday added successfully!\nName: " + birthday.Name +
			"\nDate: " + birthday.Birthday.Format("2006-01-02") +
			"\nTelegram: " + birthday.TgName
		if action == "update" {
			successMsg = "Birthday updated successfully!\nName: " + birthday.Name +
				"\nDate: " + birthday.Birthday.Format("2006-01-02") +
				"\nTelegram: " + birthday.TgName
		}
		msg := tgbotapi.NewMessage(chatID, successMsg)
		bot.Send(msg)
	}

	delete(userState, userID)
	delete(tempBirthdays, userID)
	_ = nextState
	_ = prompt
}

func HandleStart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	message := update.Message
	chat := message.Chat.ID

	start_msg := "Welcome! Nice to see you in birthday-reminder-bot.\n\nHere you can store your relatives and friends birthdays and get reminders about their birthdays. Use keyboard to call main bot functions.\n\nEnjoy!"

	msg := tgbotapi.NewMessage(chat, start_msg)
	msg.ReplyMarkup = getPresetMessageKeyboard()
	bot.Send(msg)
}

func HandleNewBirthday(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	// Initialize a new birthday struct from birthdays_helper
	tempBirthday := birthdays_helper.Initialise()
	tempBirthday.ChatID = int(chatID) // Set the ChatID after initializing

	// Store the modified struct back into the map
	tempBirthdays[userID] = tempBirthday

	// Start the process and ask for the person's name
	msg := tgbotapi.NewMessage(chatID, "Please enter the person's name:")
	bot.Send(msg)
	userState[userID] = "waiting_for_name"
	return

}

func HandleUpdateBirthday(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	tempBirthdays[userID] = birthdays_helper.Birthday{ChatID: int(chatID)}

	msg := tgbotapi.NewMessage(chatID, "Whose birthday do you want to update?")
	bot.Send(msg)
	userState[userID] = "update_waiting_for_name"
}

func HandleDeleteBirthday(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID

	msg := tgbotapi.NewMessage(chatID, "Enter the name to delete:")
	bot.Send(msg)
	userState[userID] = "delete_waiting_for_name"
}

func HandleAnswerMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	userID := update.Message.From.ID
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	switch userState[userID] {
	case "waiting_for_name":
		handleNameInput(bot, userID, chatID, text, "waiting_for_birthday", "Please enter the person's birthday (YYYY-MM-DD):")
	case "update_waiting_for_name":
		handleNameInput(bot, userID, chatID, text, "update_waiting_for_birthday", "Please enter the person's birthday (YYYY-MM-DD):")
	case "waiting_for_birthday":
		handleBirthdayInput(bot, userID, chatID, text, "waiting_for_tg_name", "Optionally, enter the person's Telegram username (or type 'skip'):")
	case "update_waiting_for_birthday":
		handleBirthdayInput(bot, userID, chatID, text, "update_waiting_for_tg_name", "Optionally, enter the person's Telegram username (or type 'skip'):")
	case "waiting_for_tg_name":
		handleTgNameInput(bot, userID, chatID, text, "", "", "add")
	case "update_waiting_for_tg_name":
		handleTgNameInput(bot, userID, chatID, text, "", "", "update")
	}
}

func HandleViewBirthdays(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	birthdays, err := birthdays_helper.GetBirthdays()
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Error retrieving birthdays.")
		bot.Send(msg)
		return
	}

	if len(birthdays) == 0 {
		msg := tgbotapi.NewMessage(chatID, "No birthdays found.")
		bot.Send(msg)
	} else {
		response := "Here are the stored birthdays:\n"
		for _, b := range birthdays {
			response += "ID: " + strconv.Itoa(b.ID) + ", Name: " + b.Name + ", Birthday: " + b.Birthday.Format("2006-01-02")
			if b.TgName != "" {
				response += ", Telegram: " + b.TgName
			}
			response += "\n"
		}
		msg := tgbotapi.NewMessage(chatID, response)
		bot.Send(msg)
	}
}

func getPresetMessageKeyboard() tgbotapi.ReplyKeyboardMarkup {
	// Define the reply keyboard with preset messages
	buttons := [][]tgbotapi.KeyboardButton{
		{tgbotapi.NewKeyboardButton("/add_new_birthday")},
		{tgbotapi.NewKeyboardButton("/show_saved_birthdays")},
		{tgbotapi.NewKeyboardButton("/update_birthday")},
		{tgbotapi.NewKeyboardButton("/delete_birthday")},
	}

	// Create and return the keyboard markup
	return tgbotapi.NewReplyKeyboard(buttons...)
}
