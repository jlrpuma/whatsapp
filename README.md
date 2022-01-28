# Whatsapp Group Message Counter

This is a learning project to get familiar with some topics related to Golang.

make sure to have your exported file on ./data/

```bash
go run main.go
```

## Cleaning data exported from whatsapp

```bash
perl -pi -e 's/[^[:ascii:]]//g' chat_modified.txt
grep "\S" chat_modified.txt
sed -n '/^[[].*/p' chat_modified.txt > output.txt
```
make sure tu use the last file generated output.txt.

## Motivation

Searching through the WhatsApp options available for groups I didn't find a tool that allows me to know how many messages were sent from every sender on the chat, so... I decided to export the the file that contains all the WhatsApp group messages

Please feel free to let comments on things where an improve can take place (code/performance)