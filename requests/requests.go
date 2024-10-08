package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/drunkleen/blum-bot/types"
)

var (
	getHeaders = map[string]string{
		"accept":          "application/json, text/plain, */*",
		"accept-language": "en-US,en;q=0.9",
		"origin":          "https://telegram.blum.codes",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
	}
)

func GetNewToken(queryID string) (string, error) {
	url := "https://user-domain.blum.codes/api/v1/auth/provider/PROVIDER_TELEGRAM_MINI_APP"

	headers := map[string]string{
		"accept":          "application/json, text/plain, */*",
		"accept-language": "en-US,en;q=0.9",
		"content-type":    "application/json",
		"origin":          "https://telegram.blum.codes",
		"priority":        "u=1, i",
		"referer":         "https://telegram.blum.codes/",
	}

	// Define the data to be sent in the POST request
	data := types.RequestBody{Query: queryID}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Try to get the token up to 3 Times
	for attempt := 1; attempt <= 3; attempt++ {
		fmt.Printf("Getting token (attempt %d)...\n", attempt)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %w", err)
		}

		for key, val := range headers {
			req.Header.Set(key, val)
		}

		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to send request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Token successfully created.")

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", fmt.Errorf("failed to read response body: %w", err)
			}

			var response types.ResponseBody
			if err := json.Unmarshal(body, &response); err != nil {
				return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
			}

			return response.Token.Refresh, nil
		} else {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("Failed to get token, attempt %d: %s\n", attempt, body)
		}

	}

	// If all attempts fail
	fmt.Println("Failed to get token after 3 attempts.")
	return "", nil
}

func GetUserInfo(token string, queryID string) (map[string]any, error) {
	url := "https://user-domain.blum.codes/api/v1/user/me"
	headers := getHeaders
	headers["Authorization"] = "Bearer " + token

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		var result map[string]any
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		return result, nil
	}

	var tokenResponse types.TokenResponseBody
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if tokenResponse.Message == "Token is invalid" {
		return nil, fmt.Errorf("token is invalid")
	} else {
		fmt.Println("Failed to get user information.")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

}

func GetUserBalance(token string) (types.UserBalance, error) {
	url := "https://game-domain.blum.codes/api/v1/user/balance"
	headers := getHeaders
	headers["Authorization"] = "Bearer " + token

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.UserBalance{}, fmt.Errorf("failed to send request: %w", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return types.UserBalance{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return types.UserBalance{}, fmt.Errorf("failed to read response body: %w", err)
		}
		var result types.UserBalance
		if err := json.Unmarshal(body, &result); err != nil {
			return types.UserBalance{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		return result, nil

	}

	return types.UserBalance{}, fmt.Errorf("failed to get balance")
}

func CheckDailyReward(token string) (map[string]any, error) {
	url := "https://game-domain.blum.codes/api/v1/daily-reward?offset=-420"
	headers := getHeaders
	headers["Authorization"] = "Bearer " + token
	headers["content-length"] = "0"
	headers["sec-ch-ua"] = `"Microsoft Edge";v="125", "Chromium";v="125", "Not.A/Brand";v="24", "Microsoft Edge WebView2";v="125"`
	headers["sec-ch-ua-mobile"] = "?0"
	headers["sec-ch-ua-platform"] = `Windows"`
	headers["sec-fetch-dest"] = "empty"
	headers["sec-fetch-mode"] = "cors"
	headers["sec-fetch-site"] = "same-site"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			fmt.Println("Failed to claim daily: Timeout")
			return nil, err
		}
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == 400 {
		var result map[string]any
		if err := json.Unmarshal(body, &result); err != nil {
			if string(body) == "OK" {
				return map[string]any{"message": "OK"}, nil
			}
			return nil, nil
		}
		return result, nil
	}

	var result map[string]any
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Json Error: %s\n", string(body))
		return nil, nil
	}

	return result, nil
}

func CheckTasks(token string) {
	url := "https://game-domain.blum.codes/api/v1/tasks"

	// Assuming getHeaders is a function, invoke it.
	headers := getHeaders
	headers["Authorization"] = "Bearer " + token
	headers["content-length"] = "0"
	headers["priority"] = "u=1, i"
	headers["sec-ch-ua"] = `"Microsoft Edge";v="125", "Chromium";v="125", "Not.A/Brand";v="24", "Microsoft Edge WebView2";v="125"`
	headers["sec-ch-ua-mobile"] = "?0"
	headers["sec-ch-ua-platform"] = `"Windows"`
	headers["sec-fetch-dest"] = "empty"
	headers["sec-fetch-mode"] = "cors"
	headers["sec-fetch-site"] = "same-site"

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return
	}

	// Set headers for the request
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	// Send the request
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v\n", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		var tasks []map[string]any
		if err := json.Unmarshal(body, &tasks); err != nil {
			log.Printf("Failed to unmarshal JSON: %v\n", err)
			return
		}
		for _, task := range tasks {
			taskTitle, ok := task["title"].(string)
			if !ok {
				continue
			}
			taskID, ok := task["id"].(string)
			if !ok {
				continue
			}
			tasksList, ok := task["tasks"].([]interface{})
			if !ok {
				continue
			}
			for _, t := range tasksList {
				taskMap, ok := t.(map[string]interface{})
				if !ok {
					continue
				}
				taskStatus, ok := taskMap["status"].(string)
				if !ok {
					continue
				}
				if taskStatus == "NOT_STARTED" {
					fmt.Printf("Starting Task: %s\n", taskTitle)
					startTask(token, taskID, taskTitle)
					claimTask(token, taskID, taskTitle)
				} else {
					fmt.Printf("[Task %s | Reward: %v] %s\n", taskStatus, task["reward"], taskTitle)
				}
			}
		}
	}
}

func startTask(token, taskID, taskTitle string) {
	url := "https://game-domain.blum.codes/api/v1/tasks/" + taskID + "/start"
	headers := getHeaders
	headers["Authorization"] = "Bearer " + token
	headers["content-length"] = "0"
	headers["priority"] = "u=1, i"
	headers["sec-ch-ua"] = `"Microsoft Edge";v="125", "Chromium";v="125", "Not.A/Brand";v="24", "Microsoft Edge WebView2";v="125"`
	headers["sec-ch-ua-mobile"] = "?0"
	headers["sec-ch-ua-platform"] = `"Windows"`
	headers["sec-fetch-dest"] = "empty"
	headers["sec-fetch-mode"] = "cors"
	headers["sec-fetch-site"] = "same-site"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("failed to send request: %v\n", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("[Task Started] %s\n", taskTitle)
	} else {
		fmt.Printf("[Staring Task Failed] %s\n", taskTitle)
	}

}

func claimTask(token, taskID, taskTitle string) {
	url := "https://game-domain.blum.codes/api/v1/tasks/" + taskID + "/claim"
	headers := getHeaders
	headers["Authorization"] = "Bearer " + token
	headers["content-length"] = "0"
	headers["priority"] = "u=1, i"
	headers["sec-ch-ua"] = `"Microsoft Edge";v="125", "Chromium";v="125", "Not.A/Brand";v="24", "Microsoft Edge WebView2";v="125"`
	headers["sec-ch-ua-mobile"] = "?0"
	headers["sec-ch-ua-platform"] = `"Windows"`
	headers["sec-fetch-dest"] = "empty"
	headers["sec-fetch-mode"] = "cors"
	headers["sec-fetch-site"] = "same-site"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Printf("failed to send request: %v\n", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to send request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("[Task Claimed] %s\n", taskTitle)
	} else {
		fmt.Printf("[Claim Task Failed] %s\n", taskTitle)
	}

}

func PlayGame(token string) (map[string]any, error) {
	url := "https://game-domain.blum.codes/api/v1/game/play"
	headers := getHeaders
	headers["Authorization"] = "Bearer " + token
	headers["content-length"] = "0"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to play the game due to connection problems: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil

}

func claimGameRequest(token, gameID string, points int) (*http.Response, error) {
	url := "https://game-domain.blum.codes/api/v1/game/claim"
	headers := map[string]string{
		"accept":          "application/json, text/plain, */*",
		"accept-language": "en-US,en;q=0.9",
		"authorization":   "Bearer " + token,
		"content-type":    "application/json",
		"origin":          "https://telegram.blum.codes",
		"priority":        "u=1, i",
		"user-agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
	}

	payload := types.GameClaimRequest{
		GameID: gameID,
		Points: points,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to claim game rewards due to connection problems: %w", err)
	}

	return resp, nil

}

func ClaimGame(token, gameID, queryID string, points int) {
	for {
		resp, err := claimGameRequest(token, gameID, points)
		if err != nil {
			fmt.Printf("Failed to claim game, try again: %v\n", err)
			return
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Failed to read response: %v\n", err)
			return
		}

		bodyString := string(body)

		if bodyString == `{"message":"game session not finished"}` {
			fmt.Println("Playing drop game...")
			time.Sleep(1 * time.Second)
			continue
		} else if bodyString == `{"message":"game session not found"}` {
			fmt.Println("The game is over")
			break
		} else {
			var response map[string]interface{}
			if err := json.Unmarshal(body, &response); err == nil {
				if response["message"] == "Token is invalid" {
					fmt.Println("Invalid token, get new token...")
					token, _ = GetNewToken(queryID)
					continue
				}
			}
			fmt.Printf("Game finished: %s\n", bodyString)
			break
		}
	}
}

func CheckBalanceFriend(token string) (types.FriendsBalance, error) {
	url := "https://gateway.blum.codes/v1/friends/balance"
	headers := map[string]string{
		"Authorization":      "Bearer " + token,
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9",
		"origin":             "https://telegram.blum.codes",
		"priority":           "u=1, i",
		"sec-ch-ua":          `"Microsoft Edge";v="125", "Chromium";v="125", "Not.A/Brand";v="24", "Microsoft Edge WebView2";v="125"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"Windows"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return types.FriendsBalance{}, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return types.FriendsBalance{}, fmt.Errorf("failed to get friend's balance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.FriendsBalance{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var friendsBalance types.FriendsBalance
	if err := json.NewDecoder(resp.Body).Decode(&friendsBalance); err != nil {
		return types.FriendsBalance{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return friendsBalance, nil
}

func ClaimBalanceFriend(token string) (bool, error) {
	url := "https://gateway.blum.codes/v1/friends/claim"
	headers := map[string]string{
		"Authorization":      "Bearer " + token,
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9",
		"content-length":     "0",
		"origin":             "https://telegram.blum.codes",
		"priority":           "u=1, i",
		"sec-ch-ua":          `"Microsoft Edge";v="125", "Chromium";v="125", "Not.A/Brand";v="24", "Microsoft Edge WebView2";v="125"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"Windows"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to claim friend's balance: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	//var claimFriendBalance map[string]any
	//if err := json.NewDecoder(resp.Body).Decode(&claimFriendBalance); err != nil {
	//	return false, fmt.Errorf("failed to decode response: %w", err)
	//}

	return true, nil
}

func ClaimFarm(token string) (bool, error) {
	url := "https://game-domain.blum.codes/api/v1/farming/claim"
	headers := getHeaders
	headers["Authorization"] = "Bearer " + token

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to claim farm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return true, nil
}

func StartFarm(token string) (bool, error) {
	url := "https://game-domain.blum.codes/api/v1/farming/start"
	headers := getHeaders
	headers["Authorization"] = "Bearer " + token

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to start farm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return true, nil
}

func GetDailyRewards(token string) (bool, error) {
	url := "https://game-domain.blum.codes/api/v1/daily-reward"
	headers := map[string]string{
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9",
		"origin":             "https://telegram.blum.codes",
		"priority":           "u=1, i",
		"sec-ch-ua":          `"Microsoft Edge";v="125", "Chromium";v="125", "Not.A/Brand";v="24", "Microsoft Edge WebView2";v="125"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"Windows"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
	}
	headers["Authorization"] = "Bearer " + token

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to make the request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return false, fmt.Errorf("failed to create request: %w", err)
		}
		for key, value := range headers {
			req.Header.Set(key, value)
		}
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return false, fmt.Errorf("failed to claim reward: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return true, nil
		}

	}
	return false, nil

}
