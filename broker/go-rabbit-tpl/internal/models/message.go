package models

import "encoding/json"

type Message struct {
	ID       int    `json:"id"`
	Content  string `json:"content"`
	Priority int    `json:"priority,omitempty"`
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func FromJSON(data []byte) (*Message, error) {
	var m Message
	return &m, json.Unmarshal(data, &m)
}
