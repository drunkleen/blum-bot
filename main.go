package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/drunkleen/blum-bot/requests"
	"github.com/drunkleen/blum-bot/utils"

	"github.com/fatih/color"
)

var (
	startTime       time.Time
	checkedTasksMap map[string]bool   = map[string]bool{}
	tokenMap        map[string]string = map[string]string{}
	checkTaskEnable bool
	claimReffEnable bool

	initStage bool = true

	printText string = ""

	bold = color.New(color.Bold).SprintFunc()

	red    = color.New(color.FgRed).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

func init() {
	startTime = time.Now()
}

func main() {
	utils.ClearScreen()
	utils.PrintLogo()
	checkTaskEnable, claimReffEnable = utils.ParseArgs()
	fmt.Printf("Task: %v\n", checkTaskEnable)
	fmt.Printf("Reff: %v\n", claimReffEnable)

	queryList, err := utils.ParseQueries()
	if err != nil {
		log.Fatalf("QueryList Error: %v\n", err)
	}

	mainLoop(queryList)
}

func mainLoop(queryList []string) {
	utils.ClearScreen()
	for { // start infinit loop

		for _, queryID := range queryList { // start query loop
			// get Token if not exists
			val, exists := tokenMap[queryID]
			if !exists || val == "" {

				newToken, err := requests.GetNewToken(queryID)
				if err != nil {
					log.Printf("Error: %v\n", err)
				}
				tokenMap[queryID] = newToken
			}
			token := tokenMap[queryID]

			// fetching user info
			userInfo, err := requests.GetUserInfo(token, queryID)
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
			username := userInfo["username"]
			printText += fmt.Sprintf("---\nUser: %s\n", bold(username))
			if initStage {
				fmt.Printf("---\nUser: %s\n", bold(username))
			}

			// fetching balance info
			balanceInfo, err := requests.GetUserBalance(token)
			if err != nil {
				log.Printf("Error: %v\n", err)
			}

			userBalanceStr := fmt.Sprintf("%v", balanceInfo["availableBalance"])
			userBalanceStr = strings.Replace(userBalanceStr, ",", ".", -1)

			availableBalance, err := strconv.ParseFloat(userBalanceStr, 64)
			if err != nil {
				log.Printf("Failed to parse available balance: %v", err)
			}
			userBalance := fmt.Sprintf("%.1f", availableBalance)
			printText += fmt.Sprintf("[Balance]: %v\n", bold(cyan(userBalance)))
			if initStage {
				fmt.Printf("[Balance]: %v\n", bold(cyan(userBalance)))
			}
			// Check Daily Rewards
			dailyRewardResponse, err := requests.CheckDailyReward(token)
			if err != nil {
				log.Printf("Error: %v\n", err)
			}

			switch dailyRewardResponse["message"] {
			case "same day":
				printText += fmt.Sprintf("[Daily Reward]: already claimed today\n")
			case "OK":
				printText += fmt.Sprintf("[Daily Reward] successfully claimed!\n")
			default:
				printText += fmt.Sprintf("[Daily Reward] Failed to check daily reward!\n")
			}

			if checkTaskEnable {
				checked, exists := checkedTasksMap[queryID]
				if !exists || !checked {
					fmt.Println("Checking tasks...")
					requests.CheckTasks(token)
					checkedTasksMap[queryID] = true
				}
			}

			// if claimReffEnable {
			// 	friendBalance, err := requests.CheckBalanceFriend(token)
			// 	if err != nil {
			// 		log.Printf("[Referrals Balance] Failed to get friend's balance: %v\n", err)
			// 	}
			// }

			playPasses := int(balanceInfo["playPasses"].(float64))
			printText += fmt.Sprintf("[%v]: %v tickets\n", bold(cyan("Play Passes")), playPasses)

		gameLoop:
			for { // game loop

				if playPasses <= 0 {
					break gameLoop
				}
				gameResponse, err := requests.PlayGame(token)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					fmt.Printf("Game played successfully: %v\n", gameResponse)
				}

				if gameId, ok := gameResponse["gameId"].(string); ok {
					requests.ClaimGame(token, gameId, queryID, 2000)
				}

			} // end game loop

		} // end queryList loop

		if initStage {
			initStage = !initStage
		}
		time.Sleep(5 * time.Second)
		utils.ClearScreen()
		utils.PrintLogo()
		fmt.Printf(
			"%v\n----------- Up Time: %v -----------",
			printText, yellow(utils.FormatUpTime(time.Since(startTime))),
		)
		printText = ""
	} // end infinit loop

}
