package message

import "encoding/json"

type Message struct {
	Key     []byte `json:"key,omitempty"`
	Value   any    `json:"value"`
	Headers map[string][]byte
}

func NewMessage(value map[string]interface{}) *Message {
	return &Message{Value: value}
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
