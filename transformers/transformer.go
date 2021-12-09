package transformers

import (
	"strings"

	m "github.com/whatsapp/models"
)

type Transformer interface {
	Transform(done <-chan interface{}, messages <-chan string)
}

type DefaultTransformer struct{}

func (dt *DefaultTransformer) Transform(done <-chan interface{}, messages <-chan string) <-chan m.Message {
	out := make(chan m.Message)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case message, ok := <-messages:
				if !ok {
					return
				}
				sentOn := strings.Replace(strings.Split(message, "]")[0], "[", "", -1)
				sender := strings.Split(strings.Split(message, "]")[1], ":")[0]
				msg := ""
				out <- m.Message{SentOn: sentOn, Sender: sender, Message: msg}
			}
		}
	}()

	return out
}
