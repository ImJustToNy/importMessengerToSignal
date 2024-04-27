package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"unicode/utf8"
)

type Message struct {
	SenderName            string `json:"sender_name"`
	TimestampMs           int64  `json:"timestamp_ms"`
	Content               string `json:"content,omitempty"`
	IsGeoblockedForViewer bool   `json:"is_geoblocked_for_viewer"`
	Reactions             []struct {
		Reaction string `json:"reaction"`
		Actor    string `json:"actor"`
	} `json:"reactions,omitempty"`
	Photos []struct {
		URI               string `json:"uri"`
		CreationTimestamp int    `json:"creation_timestamp"`
	} `json:"photos,omitempty"`
	Gifs []struct {
		URI string `json:"uri"`
	} `json:"gifs,omitempty"`
	Share struct {
		Link string `json:"link"`
	} `json:"share,omitempty"`
	Sticker struct {
		URI        string        `json:"uri"`
		AiStickers []interface{} `json:"ai_stickers"`
	} `json:"sticker,omitempty"`
	Videos []struct {
		URI               string `json:"uri"`
		CreationTimestamp int    `json:"creation_timestamp"`
	} `json:"videos,omitempty"`
	AudioFiles []struct {
		URI               string `json:"uri"`
		CreationTimestamp int    `json:"creation_timestamp"`
	} `json:"audio_files,omitempty"`
	IsUnsent     bool `json:"is_unsent,omitempty"`
	CallDuration int  `json:"call_duration,omitempty"`
}

type MessageFile struct {
	Participants []struct {
		Name string `json:"name"`
	} `json:"participants"`
	Messages           []Message     `json:"messages"`
	Title              string        `json:"title"`
	IsStillParticipant bool          `json:"is_still_participant"`
	ThreadPath         string        `json:"thread_path"`
	MagicWords         []interface{} `json:"magic_words"`
	Image              struct {
		URI               string `json:"uri"`
		CreationTimestamp int    `json:"creation_timestamp"`
	} `json:"image"`
	JoinableMode struct {
		Mode int    `json:"mode"`
		Link string `json:"link"`
	} `json:"joinable_mode"`
}

func ProcessMessageFile(content []byte, conversation *Conversation) {
	var data MessageFile

	err := json.Unmarshal(content, &data)

	if err != nil {
		log.Fatalf("Couldn't parse json file: %v", err)
	}

	for _, message := range data.Messages {
		message.SenderName = FixEncoding(message.SenderName)
		message.Content = FixEncoding(message.Content)

		//log.Printf("Processing message from %s: %s", message.SenderName, message.Content)

		person := getPersonByFacebookID(message.SenderName)

		if person == nil {
			//log.Printf("Couldn't find person with facebook id %s", message.SenderName)
			continue
		}

		SendMessage(*conversation, *person, message)
	}
}

// FixEncoding https://stackoverflow.com/a/77099052/3928847
func FixEncoding(s string) string {
	// Create a slice to hold the individual runes
	var runeSlice []rune
	// Convert the string to a slice of runes
	for _, r := range s {
		runeSlice = append(runeSlice, r)
	}

	// Create a byte slice from the rune slice
	byteSlice := make([]byte, len(runeSlice))
	for i, r := range runeSlice {
		byteSlice[i] = byte(r)
	}

	// Convert the byte slice to a UTF-8 string
	utf8String := string(byteSlice)

	// Validate that the string is valid UTF-8
	if !utf8.ValidString(utf8String) {
		// Handle invalid UTF-8
		fmt.Println("Invalid UTF-8 string")
		return ""
	}

	return utf8String
}

func SendMessage(conversation Conversation, person Person, message Message) {
	attachments := make([]string, 0)

	for _, photo := range message.Photos {
		attachments = append(attachments, photo.URI)
	}

	for _, gif := range message.Gifs {
		attachments = append(attachments, gif.URI)
	}

	for _, video := range message.Videos {
		attachments = append(attachments, video.URI)
	}

	for _, audioFile := range message.AudioFiles {
		attachments = append(attachments, audioFile.URI)
	}

	if len(attachments) > 0 {
		//log.Printf("Sending message with %d attachments", len(attachments))
		SendMessageToSignal(person.SignalNumber, conversation.SignalID, message.Content, attachments)
	}
}
