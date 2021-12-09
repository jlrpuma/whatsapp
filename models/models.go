package models

type Message struct {
	SentOn  string
	Sender  string
	Message string
}

type UserMessages struct {
	Sender  string
	Counter int
}
