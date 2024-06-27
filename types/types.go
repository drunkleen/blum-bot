package types

type UserBalance struct {
	AvailableBalance string `json:"availableBalance"`
	PlayPasses       int    `json:"playPasses"`
	Timestamp        int64  `json:"timestamp"`
	Farming          struct {
		StartTime    int64  `json:"startTime"`
		EndTime      int64  `json:"endTime"`
		EarningsRate string `json:"earningsRate"`
		Balance      string `json:"balance"`
	} `json:"farming"`
}

//type Farming struct {
//	StartTime    int64   `json:"startTime"`
//	EndTime      int64   `json:"endTime"`
//	EarningsRate float64 `json:"earningsRate"`
//	Balance      float64 `json:"balance"`
//}
//
//// UserBalance represents the response from the '/v1/user/balance' endpoint.
//type UserBalance struct {
//	AvailableBalance string `json:"availableBalance"`
//	PlayPasses       uint64 `json:"playPasses"`
//	Timestamp        int64  `json:"timestamp"`
//	Farming
//}

// FriendsBalance represents the response from the '/v1/friends/balance' endpoint.
type FriendsBalance struct {
	LimitInvitation             string  `json:"limitInvitation"`
	UsedInvitation              string  `json:"usedInvitation"`
	AmountForClaim              string  `json:"amountForClaim"`
	ReferralToken               string  `json:"referralToken"`
	PercentFromFriends          int     `json:"percentFromFriends"`
	PercentFromFriendsOfFriends float64 `json:"percentFromFriendsOfFriends"`
	CanClaim                    bool    `json:"canClaim"`
	CanClaimAt                  string  `json:"canClaimAt"`
}

type RequestBody struct {
	Query string `json:"query"`
}

type ResponseBody struct {
	Token struct {
		Refresh string `json:"refresh"`
	} `json:"token"`
}

type TokenResponseBody struct {
	Message string `json:"message"`
	ResponseBody
}

type GameClaimRequest struct {
	GameID string `json:"gameId"`
	Points int    `json:"points"`
}
