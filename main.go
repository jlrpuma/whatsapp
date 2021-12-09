package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/whatsapp/counters"
	"github.com/whatsapp/filters"
	"github.com/whatsapp/transformers"
)

const messagesDir = "data/chat.txt"

var print = fmt.Println

func main() {
	printCounting()

}

func printCounting() error {
	file, err := os.Open(messagesDir)
	if err != nil {
		log.Println("File cannot be opened", err)
		return errors.New("File cannot be opened")
	}
	defer file.Close()

	done := make(chan interface{})
	defer close(done)

	filter := filters.StringsFilter{}
	transformer := transformers.DefaultTransformer{}
	counter := counters.ParticipantCounter{}

	messages := getMessagesChan(done, file)
	filteredMessages := filter.Filter(done, messages)
	messagesOnStructs := transformer.Transform(done, filteredMessages)

	counter.Count(done, messagesOnStructs)

	return nil
}

/*
 * Line by line read the file and send the line to a channel
 */
func getMessagesChan(done <-chan interface{}, file io.Reader) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			select {
			case <-done:
				return
			case out <- scanner.Text():
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()
	return out
}
