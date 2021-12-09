package counters

import (
	"sort"

	msgs "github.com/whatsapp/messages"
)

type Counter interface {
	Count(messages <-chan msgs.Message)
}

type ParticipantCounter struct{}

func (pC *ParticipantCounter) Count(messages <-chan msgs.Message) []msgs.UserMessages {
	usersMessagesCounting := make(map[string]int)

	for messageStruct := range messages {
		if usersMessagesCounting[messageStruct.Sender] == 0 {
			usersMessagesCounting[messageStruct.Sender] = 1
		}
		if usersMessagesCounting[messageStruct.Sender] >= 1 {
			usersMessagesCounting[messageStruct.Sender] += 1
		}
	}

	var userMessages []msgs.UserMessages
	for k, v := range usersMessagesCounting {
		userMessages = append(userMessages, msgs.UserMessages{k, v})
	}

	sort.Slice(userMessages, func(i, j int) bool {
		return userMessages[i].Counter > userMessages[j].Counter
	})

	return userMessages
}
