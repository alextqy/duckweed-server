package entity

type Result struct {
	State   bool        `json:"State"`
	Code    int         `json:"Code"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}
