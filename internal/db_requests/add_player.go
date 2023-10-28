package dbrequests

import "strings"

type AddRequest struct {
	Player string `json:"player"`
}

func (r *AddRequest) SanitizedPlayer() string {
	return strings.Replace(r.Player, "'", "''", -1)
}
