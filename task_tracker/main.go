package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const tasksFile = "tasks.json"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide a task description")
			os.Exit(1)
		}
		addTask(os.Args[2])
	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Error: Please provide task ID and new description")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: Invalid task ID")
			os.Exit(1)
		}
		updateTask(id, os.Args[3])
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide task ID")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: Invalid task ID")
			os.Exit(1)
		}
		deleteTask(id)
	case "mark-in-progress":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide task ID")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: Invalid task ID")
			os.Exit(1)
		}
		markTask(id, "in-progress")
	case "mark-done":
		if len(os.Args) < 3 {
			fmt.Println("Error: Please provide task ID")
			os.Exit(1)
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: Invalid task ID")
			os.Exit(1)
		}
		markTask(id, "done")
	case "list":
		status := ""
		if len(os.Args) > 2 {
			status = os.Args[2]
		}
		listTasks(status)
	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Task Tracker CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  task-cli add <description>")
	fmt.Println("  task-cli update <id> <description>")
	fmt.Println("  task-cli delete <id>")
	fmt.Println("  task-cli mark-in-progress <id>")
	fmt.Println("  task-cli mark-done <id>")
	fmt.Println("  task-cli list [done|todo|in-progress]")
}


func loadTasks() ([]Task, error) {
	if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
		return []Task{}, nil
	}

	data, err := os.ReadFile(tasksFile)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	if len(data) == 0 {
		return []Task{}, nil
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTasks(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(tasksFile, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func getNextID(tasks []Task) int {
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}

func addTask(description string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	newTask := Task{
		ID:          getNextID(tasks),
		Description: description,
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tasks = append(tasks, newTask)

	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.ID)
}

func updateTask(id int, description string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	found := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Error: Task with ID %d not found\n", id)
		os.Exit(1)
	}

	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Task %d updated successfully\n", id)
}

func deleteTask(id int) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	found := false
	newTasks := []Task{}
	for _, task := range tasks {
		if task.ID == id {
			found = true
		} else {
			newTasks = append(newTasks, task)
		}
	}

	if !found {
		fmt.Printf("Error: Task with ID %d not found\n", id)
		os.Exit(1)
	}

	err = saveTasks(newTasks)
	if err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Task %d deleted successfully\n", id)
}

func markTask(id int, status string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	found := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Error: Task with ID %d not found\n", id)
		os.Exit(1)
	}

	err = saveTasks(tasks)
	if err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Task %d marked as %s\n", id, status)
}

func listTasks(status string) {
	tasks, err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		os.Exit(1)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	filteredTasks := []Task{}
	for _, task := range tasks {
		if status == "" || task.Status == status {
			filteredTasks = append(filteredTasks, task)
		}
	}

	if len(filteredTasks) == 0 {
		fmt.Printf("No tasks found with status '%s'\n", status)
		return
	}

	fmt.Println("\nTasks:")
	fmt.Println("------")
	for _, task := range filteredTasks {
		fmt.Printf("ID: %d\n", task.ID)
		fmt.Printf("Description: %s\n", task.Description)
		fmt.Printf("Status: %s\n", task.Status)
		fmt.Printf("Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("Updated: %s\n", task.UpdatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println("------")
	}
}
