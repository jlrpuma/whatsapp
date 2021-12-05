package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

const messagesDir = "data/chat.txt"

var print = fmt.Println

type Message struct {
	SentOn  string
	Sender  string
	Message string
}

type UserMessages struct {
	Sender  string
	Counter int
}

func main() {
	printCounting()
}

// print the message counters by participant
func printCounting() error {
	file, err := os.Open(messagesDir)
	if err != nil {
		log.Println("File cannot be opened", err)
		return errors.New("File cannot be opened")
	}
	defer file.Close()

	/* done channel
	* 	Its purpose is based on the idea that you can terminate the go routines
	*	efectively passing a value throught this channel and using a select to
	*	force the go routine returns
	 */
	done := make(chan interface{})
	defer close(done)

	messages := getMessagesChan(done, file)
	filteredMessages := filterMessages(done, messages)
	messagesOnStructs := transformMessages(done, filteredMessages)

	/* handling go routines manually
	anyway := time.After(100 * time.Millisecond)
	go func() {
		<-anyway
		close(done)
	}()
	*/

	/* if you want to output the resutls on a file
	outputFile, err := os.Create("/tmp/output1")
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()
	*/

	usersMessagesCounterMap := make(map[string]int)
	for messageStruct := range messagesOnStructs {
		if usersMessagesCounterMap[messageStruct.Sender] == 0 {
			usersMessagesCounterMap[messageStruct.Sender] = 1
		} else if usersMessagesCounterMap[messageStruct.Sender] >= 1 {
			usersMessagesCounterMap[messageStruct.Sender] += 1
		}
	}

	var userMessagesSlice []UserMessages
	for k, v := range usersMessagesCounterMap {
		userMessagesSlice = append(userMessagesSlice, UserMessages{k, v})
	}

	sort.Slice(userMessagesSlice, func(i, j int) bool {
		return userMessagesSlice[i].Counter > userMessagesSlice[j].Counter
	})

	for _, userMessage := range userMessagesSlice {
		fmt.Printf("%-40s%d\n", userMessage.Sender, userMessage.Counter)
	}

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

/*
* Filter some messages that does not add any kind of value to the counter by participant
 */
func filterMessages(done <-chan interface{}, stringMessages <-chan string) <-chan string {
	out := make(chan string)
	go func(sM <-chan string) {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case stringMessage, ok := <-sM:
				if !ok {
					return
				}
				if strings.Contains(stringMessage, "[") &&
					!strings.Contains(stringMessage, "added") &&
					!strings.Contains(stringMessage, "removed") &&
					!strings.Contains(stringMessage, "joined") &&
					!strings.Contains(stringMessage, "left") &&
					!strings.Contains(stringMessage, "changed") &&
					!strings.Contains(stringMessage, "created group") &&
					!strings.Contains(stringMessage, "disappearing messages") &&
					!strings.Contains(stringMessage, "You're now an admin") {
					out <- stringMessage
				}
			}
		}
	}(stringMessages)
	return out
}

/*
* Transform messages from string to an apropiate structure
* TODO: Splits and Replace can be replaced by a fancy regex usage
 */
func transformMessages(done <-chan interface{}, stringMessages <-chan string) <-chan Message {
	out := make(chan Message)

	go func(sM <-chan string) {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case stringMessage, ok := <-sM:
				if !ok {
					return
				}
				sentOn := strings.Replace(strings.Split(stringMessage, "]")[0], "[", "", -1)
				sender := strings.Split(strings.Split(stringMessage, "]")[1], ":")[0]
				message := ""
				out <- Message{SentOn: sentOn, Sender: sender, Message: message}
			}
		}
	}(stringMessages)

	return out
}
