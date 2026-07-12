package messages

type Role int

const (
	AGENT Role = iota
	USER
	SYSTEM
	DEVELOPER
)
