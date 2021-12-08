package filters

import (
	"strings"
)

type FilterStrategy interface {
	Filter(done <-chan interface{}, messages <-chan string) <-chan string
}

type StringsFilter struct{}

func (sF *StringsFilter) Filter(done <-chan interface{}, messages <-chan string) <-chan string {
	out := make(chan string)
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
				if strings.Contains(message, "[") &&
					!strings.Contains(message, "added") &&
					!strings.Contains(message, "removed") &&
					!strings.Contains(message, "joined") &&
					!strings.Contains(message, "left") &&
					!strings.Contains(message, "changed") &&
					!strings.Contains(message, "created group") &&
					!strings.Contains(message, "disappearing messages") &&
					!strings.Contains(message, "You're now an admin") {
					out <- message
				}
			}
		}
	}()
	return out
}

type RegexFilter struct{}

func (rf *RegexFilter) Filter(done <-chan interface{}, messages <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case _, ok := <-messages:
				if !ok {
					return
				}
				// TODO: Pending implementation
				// learning about regex to get the right information
			}
		}
	}()
	return out
}
