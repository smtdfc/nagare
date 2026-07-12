package messages

type Text struct {
	Content string `json:"content"`
	Role    Role   `json:"role"`
}

func (t *Text) Kind() string {
	return "Text"
}
