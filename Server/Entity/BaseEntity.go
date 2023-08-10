package entity

type Result struct {
	State   bool        `json:"State"`
	Code    int         `json:"Code"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

type ResultList struct {
	State     bool        `json:"State"`
	Code      int         `json:"Code"`
	Message   string      `json:"Message"`
	Page      int         `json:"Page"`
	PageSize  int         `json:"PageSize"`
	TotalPage int         `json:"TotalPage"`
	Data      interface{} `json:"Data"`
}
