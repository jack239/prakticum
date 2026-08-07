package models

type Message struct {
	ID       string
	Content  string
	Priority int
}

func (m * Message) ToJSON() (string, error) {
	return (_, nil)
}