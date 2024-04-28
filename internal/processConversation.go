package internal

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func ProcessConversation(conversation *Conversation) {
	log.Printf("Processing conversation %s", conversation.FacebookID)

	directory, err := filepath.Abs(filepath.Join(GetConfig().Entrypoint, "your_facebook_activity/messages/inbox", conversation.FacebookID))
	log.Printf("Processing directory: %s", directory)

	if err != nil {
		log.Fatalf("Couldn't process directory name of %s/%s", GetConfig().Entrypoint, conversation.FacebookID)
	}

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Printf("Couldn't find a directory %s", directory)
	}

	amountOfJsonFiles := 0

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			amountOfJsonFiles++
		}
	}

	log.Printf("Found %d json files", amountOfJsonFiles)

	if amountOfJsonFiles == 0 {
		return
	}

	for messageFileNumber := amountOfJsonFiles; messageFileNumber > 0; messageFileNumber-- {
		messageFile := filepath.Join(directory, "message_"+strconv.Itoa(messageFileNumber)+".json")

		log.Printf("Processing file %s (%d/%d)", messageFile, amountOfJsonFiles-messageFileNumber+1, amountOfJsonFiles)

		content, err := os.ReadFile(messageFile)

		if err != nil {
			log.Fatalf("Couldn't read file %s", messageFile)
		}

		ProcessMessageFile(content, conversation)
	}
}
