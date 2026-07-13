package messages

type Reasoning struct {
	Content string `json:"content"`
}

func (t *Reasoning) Kind() string {
	return "Reasoning"
}
