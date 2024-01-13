package main

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	VERSION = 0.02
)

type Session struct {
	New bool `json:"new"`
}

type Event struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Request struct {
		OriginalUtterance string `json:"original_utterance"`
	} `json:"request"`
}

type Result struct {
	Text       string `json:"text"`
	TTS        string `json:"tts"`
	EndSession bool   `json:"end_session"`
}

type Response struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Result  Result  `json:"response"`
}

func Handler(ctx context.Context, event []byte) (*Response, error) {
	fmt.Printf("VERSION := %v", VERSION)

	var input Event
	err := json.Unmarshal(event, &input)
	if err != nil {
		return nil, fmt.Errorf("an error has occurred when parsing event: %v", err)
	}

	res := Result{}

	dialogue := DialogInstance()

	if input.Session.New {
		res.Text = dialogue.Start.Text
		res.TTS = dialogue.Start.TTS
	}

	// text := "Привет малыш, как дела?"
	// voiceText := `<speaker audio="dialogs-upload/e74da21f-e89c-4989-ab7a-f2138e7a7e33/f436ceb1-cee7-448f-a51b-a9c7fe0b4c0c.opus">`

	// if input.Request.OriginalUtterance != "" {
	// 	text = input.Request.OriginalUtterance
	// }

	fmt.Printf("Input text: %v", input.Request.OriginalUtterance)

	return &Response{
		Version: input.Version,
		Session: input.Session,
		Result:  res,
	}, nil
}
