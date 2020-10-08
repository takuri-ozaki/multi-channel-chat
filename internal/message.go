package internal

type Message struct {
	Message  string `json:"message"`
	UserName string `json:"user_name"`
	System   bool   `json:"system"`
}
