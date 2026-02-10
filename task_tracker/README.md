# Task Tracker CLI

A simple command-line interface (CLI) application to track and manage your tasks. Built with Go.

**Project URL:** https://roadmap.sh/projects/task-tracker

## Features

- Add, update, and delete tasks
- Mark tasks as in-progress or done
- List all tasks or filter by status (todo, in-progress, done)
- Tasks stored in JSON format
- Automatic timestamp tracking (created and updated times)

## Requirements

- Go 1.21 or higher

## Installation

1. Clone or download this project
2. Navigate to the project directory:
   ```bash
   cd task_tracker
   ```

3. Build the application:
   ```bash
   go build -o task-cli
   ```

## Usage

### Add a new task
```bash
./task-cli add "Buy groceries"
# Output: Task added successfully (ID: 1)
```

### Update a task
```bash
./task-cli update 1 "Buy groceries and cook dinner"
# Output: Task 1 updated successfully
```

### Delete a task
```bash
./task-cli delete 1
# Output: Task 1 deleted successfully
```

### Mark a task as in-progress
```bash
./task-cli mark-in-progress 1
# Output: Task 1 marked as in-progress
```

### Mark a task as done
```bash
./task-cli mark-done 1
# Output: Task 1 marked as done
```

### List all tasks
```bash
./task-cli list
```

### List tasks by status
```bash
# List completed tasks
./task-cli list done

# List pending tasks
./task-cli list todo

# List in-progress tasks
./task-cli list in-progress
```

## Task Properties

Each task has the following properties:
- **id**: Unique identifier for the task
- **description**: Short description of the task
- **status**: Current status (todo, in-progress, done)
- **createdAt**: Timestamp when the task was created
- **updatedAt**: Timestamp when the task was last updated

## Data Storage

Tasks are stored in a `tasks.json` file in the current directory. The file is created automatically if it doesn't exist.

## Example Workflow

```bash
# Add some tasks
./task-cli add "Write project documentation"
./task-cli add "Review pull requests"
./task-cli add "Deploy to production"

# Start working on a task
./task-cli mark-in-progress 1

# Complete a task
./task-cli mark-done 1

# View all in-progress tasks
./task-cli list in-progress

# Update a task description
./task-cli update 2 "Review and merge pull requests"

# Delete a task
./task-cli delete 3
```

## Error Handling

The application handles common errors gracefully:
- Invalid task IDs
- Missing required arguments
- File read/write errors
- JSON parsing errors

## Project Structure

```
task_tracker/
├── main.go       # Main application code
├── go.mod        # Go module file
├── README.md     # This file
└── tasks.json    # Task data (created automatically)
```

## License

This project is open source and available for educational purposes.
