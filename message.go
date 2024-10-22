package gollum

import (
	"time"
)

type Message struct {
	Role    string
	Content string
}

type Response struct {
	Id       string
	Created  time.Time
	Model    string
	Messages []Message
}
