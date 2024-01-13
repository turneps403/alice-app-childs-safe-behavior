package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

const (
	VERSION = 0.02
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
}

type Event struct {
	Version string   `json:"version"`
	Session struct{} `json:"session"`
	Request struct {
		OriginalUtterance string `json:"original_utterance"`
	} `json:"request"`
}

type Response struct {
	Version string   `json:"version"`
	Session struct{} `json:"session"`
	Result  struct {
		Text       string `json:"text"`
		Tts        string `json:"tts"`
		EndSession bool   `json:"end_session"`
	} `json:"response"`
}

func Handler(ctx context.Context, event []byte) (*Response, error) {
	log.Printf("VERSION := %v", VERSION)

	var input Event
	err := json.Unmarshal(event, &input)
	if err != nil {
		return nil, fmt.Errorf("an error has occurred when parsing event: %v", err)
	}

	text := "Привет малыш, как дела?"
	voiceText := `<speaker audio="dialogs-upload/e74da21f-e89c-4989-ab7a-f2138e7a7e33/f436ceb1-cee7-448f-a51b-a9c7fe0b4c0c.opus">`

	if input.Request.OriginalUtterance != "" {
		text = input.Request.OriginalUtterance
	}

	log.Printf("Input text: %v", input.Request.OriginalUtterance)

	return &Response{
		Version: input.Version,
		Session: input.Session,
		Result: struct {
			Text       string `json:"text"`
			Tts        string `json:"tts"`
			EndSession bool   `json:"end_session"`
		}{
			Text:       text,
			Tts:        voiceText,
			EndSession: true,
		},
	}, nil
}
