package main

import (
	"fmt"
	"log"
	"strconv"
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
	claimRefEnable  bool

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
	checkTaskEnable, claimRefEnable = utils.ParseArgs()
	fmt.Printf("Task: %v\n", checkTaskEnable)
	fmt.Printf("Reff: %v\n", claimRefEnable)

	queryList, err := utils.ParseQueries()
	if err != nil {
		log.Fatalf(red("QueryList Error: %v\n"), err)
	}

	mainLoop(queryList)
}

func mainLoop(queryList []string) {
	utils.ClearScreen()
	for { // start infinite loop

		for _, queryID := range queryList { // start query loop
			// get Token if not exists
			val, exists := tokenMap[queryID]
			if !exists || val == "" {
				newToken, err := requests.GetNewToken(queryID)
				if err != nil {
					log.Printf(red("Error: %v\n"), err)
				}
				tokenMap[queryID] = newToken
			}
			token := tokenMap[queryID]

			// fetching user info
			userInfo, err := requests.GetUserInfo(token, queryID)
			if err != nil {
				log.Printf(red("Error: %v\n"), err)
			}
			username := userInfo["username"]
			printText += fmt.Sprintf("---\n"+"["+bold(cyan("User"))+"] "+"%s\n", bold(green(username)))
			if initStage {
				fmt.Printf("---\nUser: %s\n", bold(username))
			}

			// fetching balance info
			balanceInfo, err := requests.GetUserBalance(token)
			if err != nil {
				log.Printf(red("Error: %v\n"), err)
			}

			printText += fmt.Sprintf("["+bold(cyan("Balance"))+"] "+"%v\n", balanceInfo.AvailableBalance)
			if initStage {
				fmt.Printf("[Balance] %v\n", bold(cyan(balanceInfo.AvailableBalance)))
			}

			nextFarmingTime, _ := utils.TimeLeft(balanceInfo.Farming.EndTime)

			if balanceInfo.Farming.Balance != "57.6" {
				printText += fmt.Sprintf("["+bold(cyan("Farming"))+"] "+"next claim %v remaining", nextFarmingTime)
				printText += fmt.Sprintf(" | Earned: %v\n", balanceInfo.Farming.Balance)
			} else {
				ok, err := requests.ClaimFarm(token)
				if err != nil {
					log.Printf(red("Error: %v\n"), err)
				}
				if ok {
					printText += fmt.Sprintf("[" + bold(cyan("Farming")) + "] " + "claimed successfully!\n")
				} else {
					printText += fmt.Sprintf("[" + bold(cyan("Farming")) + "] " + "Failed to claim farm!\n")
				}

				ok, err = requests.StartFarm(token)
				if err != nil {
					log.Printf(red("Error: %v\n"), err)
				}
				if ok {
					printText += fmt.Sprintf("[" + bold(cyan("Farming")) + "] " + "started farming successfully!\n")
				} else {
					printText += fmt.Sprintf("[" + bold(cyan("Farming")) + "] " + "Failed to start farming!\n")
				}
			}

			// Check Daily Rewards
			dailyRewardResponse, err := requests.CheckDailyReward(token)
			if err != nil {
				log.Printf(red("Error: %v\n"), err)
			}

			switch dailyRewardResponse["message"] {
			case "same day":
				printText += fmt.Sprintf("[" + bold(cyan("Daily Reward")) + "] " + " already claimed today\n")
			case "OK":
				printText += fmt.Sprintf("[" + bold(cyan("Daily Reward")) + "] " + " successfully claimed!\n")
			default:
				printText += fmt.Sprintf("[" + bold(cyan("Daily Reward")) + "] " + " Failed to check daily reward!\n")
			}

			if checkTaskEnable {
				checked, exists := checkedTasksMap[queryID]
				if !exists || !checked {
					fmt.Println("Checking tasks...")
					requests.CheckTasks(token)
					checkedTasksMap[queryID] = true
				}
			}

			if claimRefEnable {
				friendsBalance, err := requests.CheckBalanceFriend(token)

				if friendsBalance.AmountForClaim != "0" || friendsBalance.CanClaim {

					if err != nil {
						log.Printf(red("[Referrals Balance] Failed to get friend's balance: %v\n"), err)
					}

					printText += fmt.Sprintf("[" + bold(cyan("Referrals")) + "] ")
					printText += fmt.Sprintf(yellow("amount: %v"), friendsBalance.AmountForClaim)
					printText += fmt.Sprintf(yellow(" | Claimable: %v"), friendsBalance.CanClaim)

					var claimTime int64
					if friendsBalance.CanClaimAt != "" {
						claimTime, err = strconv.ParseInt(friendsBalance.CanClaimAt, 10, 64)
						if err != nil {
							log.Printf(red("[Referrals] Failed to parse claim time: %v\n"), err)
						}
						remainingClaimTime, err := utils.TimeLeft(claimTime)
						if err != nil {
							log.Printf(red("[Referrals] Failed to calculate remaining claim time: %v\n"), err)
						}
						printText += fmt.Sprintf(yellow(" | %v remaining\n"), remainingClaimTime)
					} else {
						printText += fmt.Sprintf("\n")
					}

					if friendsBalance.CanClaim {
						ok, err := requests.ClaimBalanceFriend(token)
						if err != nil {
							log.Printf(red("[Referrals] Failed to claim friend's balance: %v\n"), err)
						}
						if ok {
							printText += fmt.Sprintf(bold(cyan("[Referrals]")) + " successfully claimed!\n")
						}
					}

				}
			}

			printText += fmt.Sprintf("[%v] %v tickets\n", bold(cyan("Play Passes")), balanceInfo.PlayPasses)

			for balanceInfo.PlayPasses > 0 { // game loop
				gameResponse, err := requests.PlayGame(token)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					fmt.Printf("Game played successfully: %v\n", gameResponse)
				}

				if gameId, ok := gameResponse["gameId"].(string); ok {
					requests.ClaimGame(token, gameId, queryID, 2000)
				}
				balanceInfo.PlayPasses--

			} // end game loop

		} // end queryList loop

		if initStage {
			initStage = !initStage
		} else {
			time.Sleep(5 * time.Second)
		}
		h, m, _ := time.Now().Clock()
		utils.ClearScreen()
		utils.PrintLogo()
		fmt.Printf(
			"-------- Time: %d:%d --------\n%v\n------------ Up Time: %v ------------",
			h, m, printText, yellow(utils.FormatUpTime(time.Since(startTime))),
		)
		printText = ""
	} // end infinite loop

}
