package dbrequests

type ChangePlayerRequest struct {
	Player string `json:"player"`
	Val    int    `json:"value"`
}
