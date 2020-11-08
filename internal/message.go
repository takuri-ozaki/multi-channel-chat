package internal

import "strconv"

type Message struct {
	Message     string `json:"message"`
	UserName    string `json:"user_name"`
	System      bool   `json:"system"`
	Transferred bool   `json:"transferred"`
}

func (m Message) ToTransferredJson() []byte {
	return []byte(`{"message":"` + m.Message + `","user_name":"` + m.UserName + `","system":` + strconv.FormatBool(m.System) + `,"transferred":true}`)
}
