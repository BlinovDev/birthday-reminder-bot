package birthdays_helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	// "telegram-webhook-bot/log_helper"
)

const filePath = "tasks.json" // File where tasks will be stored

// define struct of stored item
type Birthday struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	MessageId int    `json:"message_id"`
}

// Helper function to read tasks from file
func readTasks() ([]Birthday, error) {
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

// Helper function to write tasks to the file
func writeTasks(tasks []Birthday) error {
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

func AddTask(text string, message_id int) error {
	// Read the existing tasks
	tasks, err := readTasks()
	if err != nil {
		return err
	}

	// Generate a new task with a unique ID
	newID := 1
	if len(tasks) > 0 {
		newID = tasks[len(tasks)-1].ID + 1 // Increment ID from last task
	}

	newTask := Birthday{
		ID:        newID,
		Text:      text,
		MessageId: message_id,
	}

	// Append the new task
	tasks = append(tasks, newTask)

	// Write the updated tasks back to the file
	err = writeTasks(tasks)
	if err != nil {
		return err
	}

	// log_helper.PrintLog("MESSAGE", text)

	return nil
}

func GetTasks() ([]Birthday, error) {
	// Read tasks from the file
	tasks, err := readTasks()
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func Delete(message_id int) error {
	// Read the existing tasks
	tasks, err := readTasks()
	if err != nil {
		return err
	}

	// Find and remove the task by ID
	index := -1
	for i, task := range tasks {
		if task.MessageId == message_id {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("task with ID %d not found", message_id)
	}

	// Remove the task from the slice
	tasks = append(tasks[:index], tasks[index+1:]...)

	// Write the updated tasks back to the file
	err = writeTasks(tasks)
	if err != nil {
		return err
	}

	return nil
}
