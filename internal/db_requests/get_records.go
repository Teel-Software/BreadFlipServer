package dbrequests

type Record struct {
	Player string `json:"player"`
	Val    int    `json:"record"`
}

type RecordList struct {
	List []Record `json:"record_list"`
}
