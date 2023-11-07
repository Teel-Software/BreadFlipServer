package dbrequests

import "strings"

type ChangePlayerRequest struct {
	Player string `json:"player"`
	Val    int    `json:"record"`
}

func (r *ChangePlayerRequest) SanitizedPlayer() string {
	return strings.Replace(r.Player, "'", "''", -1)
}
