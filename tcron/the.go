package tcron

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Schedule string
	Command  string
}

func Runner() {
    // this needs to read from a file 
	tasks := []Task{
		{Schedule: "*/5 * * * *", Command: "echo 'Task 1 executed'"},
		{Schedule: "0 0 * * *", Command: "echo 'Task 2 executed'"},
	}

    // this rruns every task as a goroutine
	for _, task := range tasks {
		go runTask(task)
	}

	// Keep the program running
	select {}
}

func runTask(task Task) {
	schedule, err := parseSchedule(task.Schedule)
	if err != nil {
		log.Printf("Failed to parse schedule for command '%s': %s", task.Command, err)
		return
	}

	for {
		now := time.Now().UTC()
		nextRun := schedule.Next(now)

		// Calculate the duration until the next run
		duration := nextRun.Sub(now)

		// Sleep until the next run
		time.Sleep(duration)

		// Execute the command
		cmd := exec.Command("sh", "-c", task.Command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to execute command '%s': %s", task.Command, err)
		} else {
			log.Printf("Executed command '%s'", task.Command)
		}
	}
}

// Parse the cron schedule string and return a cron.Schedule instance
func parseSchedule(schedule string) (cron.Schedule, error) {
	// You can use a third-party cron library or implement your own parser here
	// This example uses the github.com/robfig/cron package

	parsedSchedule, err := cron.ParseStandard(schedule)
	if err != nil {
		return nil, err
	}

	return parsedSchedule, nil
}

// this should be an internal detail
type Schedule struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
}

// adding a scheduler interface

func parseSchedule(scheduleString string) (Schedule, error) {
	parts := strings.Fields(scheduleString)

	if len(parts) != 5 {
		return Schedule{}, fmt.Errorf("invalid schedule format")
	}

	return Schedule{
		Minute:     parts[0],
		Hour:       parts[1],
		DayOfMonth: parts[2],
		Month:      parts[3],
		DayOfWeek:  parts[4],
	}, nil
}
