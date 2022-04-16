package models

type QuotesMsg struct {
	Timestamp string `json:"timestamp"`
	Date      string `json:"date"`
	USD       string `json:"usd"`
	EUR       string `json:"eur"`
}
