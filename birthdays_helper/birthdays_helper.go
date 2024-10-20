package birthdays_helper

import (
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"os"
	"time"
	// "telegram-webhook-bot/log_helper"
)

const filePath = "birthdays.json" // File where tasks will be stored

// define struct of stored item
type Birthday struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	TgName   string    `json:"tg_name"`
	Birthday time.Time `json:"birthday"`
}

// Helper function to read birthdays from file
func readBirthdays() ([]Birthday, error) {
	var birthdays []Birthday

	// Check if file exists, if not, create an empty one
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := ioutil.WriteFile(filePath, []byte("[]"), 0644) // Create an empty JSON array
		if err != nil {
			return nil, err
		}
	}

	// Read the file
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into the tasks slice
	err = json.Unmarshal(data, &birthdays)
	if err != nil {
		return nil, err
	}

	return birthdays, nil
}

// Helper function to write birthdays to the file
func writeBirthdays(tasks []Birthday) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func AddBirthday(name string, tg_name string, birthday time.Time) error {
	// Read the existing birthdays
	birthdays, err := readBirthdays()
	if err != nil {
		return err
	}

	// Generate a new birthday with a unique ID
	newID := 1
	if len(birthdays) > 0 {
		newID = birthdays[len(birthdays)-1].ID + 1 // Increment ID from last birthday
	}

	newTask := Birthday{
		ID:       newID,
		Name:     name,
		TgName:   tg_name,
		Birthday: birthday,
	}

	// Append the new birthday
	birthdays = append(birthdays, newTask)

	// Write the updated birthday back to the file
	err = writeBirthdays(birthdays)
	if err != nil {
		return err
	}

	return nil
}

func GetBirthdays() ([]Birthday, error) {
	// Read tasks from the file
	birthdays, err := readBirthdays()
	if err != nil {
		return nil, err
	}

	return birthdays, nil
}

func Delete(name string) error {
	// Read the existing tasks
	birthdays, err := readBirthdays()
	if err != nil {
		return err
	}

	// Find and remove the task by ID
	index := -1
	for i, birthday := range birthdays {
		if birthday.Name == name {
			index = i
			break
		}
	}

	// if index == -1 {
	// 	return fmt.Errorf("task with ID %d not found", message_id)
	// }

	// Remove the task from the slice
	birthdays = append(birthdays[:index], birthdays[index+1:]...)

	// Write the updated tasks back to the file
	err = writeBirthdays(birthdays)
	if err != nil {
		return err
	}

	return nil
}
