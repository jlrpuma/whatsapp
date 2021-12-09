package counters

import (
	"fmt"
	"sort"

	m "github.com/whatsapp/models"
)

type Counter interface {
	Count(done <-chan interface{}, messages <-chan m.Message)
}

type ParticipantCounter struct{}

func (pC *ParticipantCounter) Count(done <-chan interface{}, messages <-chan m.Message) {
	usersMessagesCounterMap := make(map[string]int)

	for messageStruct := range messages {
		if usersMessagesCounterMap[messageStruct.Sender] == 0 {
			usersMessagesCounterMap[messageStruct.Sender] = 1
		} else if usersMessagesCounterMap[messageStruct.Sender] >= 1 {
			usersMessagesCounterMap[messageStruct.Sender] += 1
		}
	}

	var userMessagesSlice []m.UserMessages
	for k, v := range usersMessagesCounterMap {
		userMessagesSlice = append(userMessagesSlice, m.UserMessages{k, v})
	}

	sort.Slice(userMessagesSlice, func(i, j int) bool {
		return userMessagesSlice[i].Counter > userMessagesSlice[j].Counter
	})

	for _, userMessage := range userMessagesSlice {
		fmt.Printf("%-40s%d\n", userMessage.Sender, userMessage.Counter)
	}
}
