package dbrequests

import "strings"

// Add new player
type AddRequest struct {
	Player string `json:"player"`
}

func (r *AddRequest) SanitizedPlayer() string {
	return strings.Replace(r.Player, "'", "''", -1)
}

// Change player record
type ChangePlayerRequest struct {
	Player int `json:"player"`
	Val    int `json:"record"`
}

// Get record list
type Record struct {
	Id     int    `json:"id"`
	Player string `json:"player"`
	Val    int    `json:"record"`
}

type RecordList struct {
	List []Record `json:"record_list"`
}
