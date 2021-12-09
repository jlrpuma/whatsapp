package transformers

import (
	"strings"

	msgs "github.com/whatsapp/messages"
)

type Transformer interface {
	Transform(done <-chan interface{}, messages <-chan string)
}

type DefaultTransformer struct{}

func (dt *DefaultTransformer) Transform(done <-chan interface{}, messages <-chan string) <-chan msgs.Message {
	out := make(chan msgs.Message)

	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case msg, ok := <-messages:
				if !ok {
					return
				}
				sentOn := strings.Replace(strings.Split(msg, "]")[0], "[", "", -1)
				sender := strings.Split(strings.Split(msg, "]")[1], ":")[0]
				content := ""
				out <- msgs.Message{SentOn: sentOn, Sender: sender, Message: content}
			}
		}
	}()

	return out
}
