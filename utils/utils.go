package utils

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	yellow = color.New(color.FgYellow).SprintFunc()
)

func ParseArgs() (bool, bool) {
	tasks := flag.String("task", "", "Check and Claim Task (Y/n)")
	reffs := flag.String("reff", "", "Do you want to claim referrals? (Y/n)")
	flag.Parse()

	var tasksEnable bool
	var reffsEnable bool

	if *tasks == "" {
		fmt.Print("Do you want to check and claim tasks? (Y/n): ")
		var taskInput string
		fmt.Scanln(&taskInput)
		taskInput = strings.TrimSpace(strings.ToLower(taskInput))

		switch taskInput {
		case "y":
			tasksEnable = true
		case "n":
			tasksEnable = false
		default:
			tasksEnable = true
		}
	}

	if *reffs == "" {
		fmt.Print("Do you want to claim Referrals? (Y/n): ")
		var reffInput string
		fmt.Scanln(&reffInput)
		reffInput = strings.TrimSpace(strings.ToLower(reffInput))

		switch reffInput {
		case "y":
			reffsEnable = true
		case "no":
			reffsEnable = false
		default:
			reffsEnable = true
		}
	}

	return tasksEnable, reffsEnable
}

func PrintLogo() {

	fmt.Printf(yellow(" ____  _    _   _ __  __ ____   ___ _____ \n| __ )| |  | | | |  \\/  | __ ) / _ \\_   _|\n|  _ \\| |  | | | | |\\/| |  _ \\| | | || |  \n| |_) | |__| |_| | |  | | |_) | |_| || |  \n|____/|_____\\___/|_|  |_|____/ \\___/ |_|  \n"))
}

func ClearScreen() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		return
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func loadListFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var parseList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			parseList = append(parseList, strings.TrimSpace(scanner.Text()))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(parseList) < 1 {
		return nil, errors.New(fmt.Sprintf("\"%v\" is empty!", fileName))
	}

	return parseList, nil

}

func ParseQueries() ([]string, error) {
	return loadListFile("./configs/query_list.conf")
}

func FormatUpTime(d time.Duration) string {
	totalSeconds := int(d.Seconds())

	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
}

func TimeLeft(futureTimestamp int64) (string, error) {

	seconds := futureTimestamp / 1000
	nanoseconds := (futureTimestamp % 1000) * 1e6

	t := time.Unix(seconds, nanoseconds)

	currentTime := time.Now()

	duration := t.Sub(currentTime)

	if duration < 0 {
		return "00:00:00", nil
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds = int64(duration.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds), nil
}
