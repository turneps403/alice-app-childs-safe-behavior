package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	VERSION = 0.02
)

type Session struct {
	New bool `json:"new"`
}

type RequestEvent struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	State   struct {
		SessionState SessionState `json:"session"`
	} `json:"state"`
	Request struct {
		OriginalUtterance string `json:"original_utterance"`
	} `json:"request"`
}

type Result struct {
	Text       string `json:"text"`
	TTS        string `json:"tts"`
	EndSession bool   `json:"end_session"`
}

type SessionState struct {
	PlaceID  int `json:"place_id"`
	RiddleID int `json:"riddle_id"`
}

type Response struct {
	Version      string       `json:"version"`
	Session      Session      `json:"session"`
	Result       Result       `json:"response"`
	SessionState SessionState `json:"session_state"`
}

func Handler(ctx context.Context, event []byte) (*Response, error) {
	fmt.Printf("VERSION := %v", VERSION)

	var input RequestEvent
	err := json.Unmarshal(event, &input)
	if err != nil {
		return nil, fmt.Errorf("an error has occurred when parsing event: %v", err)
	}

	res := Result{}
	dialogue := DialogInstance()
	sessionData := input.State.SessionState

	userAnswer := input.Request.OriginalUtterance
	if amongTokens(userAnswer, dialogue.Interrupt.Tokens) {
		// interrupt game
		res.EndSession = true
		applyTextTTS(
			&res,
			dialogue.Interrupt.Speech,
		)
	} else if amongTokens(userAnswer, dialogue.Help.Tokens) {
		// help
		applyTextTTS(
			&res,
			dialogue.Help.Speech,
		)
	} else if amongTokens(userAnswer, dialogue.Abilities.Tokens) {
		// abilities
		applyTextTTS(
			&res,
			dialogue.Abilities.Speech,
		)
	} else if sessionData.PlaceID != 0 && amongTokens(userAnswer, dialogue.Repeat) {
		// repeat
		placeExtID, riddleID := sessionData.PlaceID, sessionData.RiddleID
		answer := []Phrase{}
		answer = append(answer, dialogue.Places[placeExtID-1].Riddles[riddleID].Question)
		applyTextTTS(
			&res,
			answer...,
		)
	} else if input.Session.New {
		// the very start
		applyTextTTS(
			&res,
			dialogue.Start,
			dialogue.PlacesAvaliable,
		)
	} else if sessionData.PlaceID == 0 {
		// choose a place
		placeExtID, riddleID := 0, 0
		// strings.ToLower
		for i, place := range dialogue.Places {
			if amongTokens(userAnswer, place.Tokens) {
				placeExtID = i + 1
				break
			}
		}

		if placeExtID == 0 {
			// userAnswer wasnt match to any place
			applyTextTTS(
				&res,
				dialogue.Fail,
				dialogue.PlacesAvaliable,
			)
		} else {
			sessionData.PlaceID = placeExtID
			sessionData.RiddleID = riddleID
			applyTextTTS(
				&res,
				dialogue.Places[placeExtID-1].Prologue,
				dialogue.Places[placeExtID-1].Riddles[riddleID].Question,
			)
		}
	} else {
		placeExtID, riddleID := sessionData.PlaceID, sessionData.RiddleID
		answer := []Phrase{}
		if amongTokens(userAnswer, dialogue.Places[placeExtID-1].Riddles[riddleID].Answers) {
			answer = append(answer, dialogue.Places[placeExtID-1].Riddles[riddleID].Reaction.Right)
		} else {
			answer = append(answer, dialogue.Places[placeExtID-1].Riddles[riddleID].Reaction.Wrong)
		}
		answer = append(answer, dialogue.Places[placeExtID-1].Riddles[riddleID].Reaction.Explanation)

		if riddleID+1 == len(dialogue.Places[placeExtID-1].Riddles) {
			res.EndSession = true
			answer = append(answer, dialogue.Places[placeExtID-1].Epilogue)
		} else {
			sessionData.RiddleID = riddleID + 1
			answer = append(answer, dialogue.Places[placeExtID-1].Riddles[riddleID+1].Question)
		}

		applyTextTTS(
			&res,
			answer...,
		)
	}

	return &Response{
		Version:      input.Version,
		Session:      input.Session,
		Result:       res,
		SessionState: sessionData,
	}, nil
}

func applyTextTTS(res *Result, phrases ...Phrase) {
	txt, tts := []string{}, []string{}
	for _, p := range phrases {
		txt = append(txt, p.Text)
		if p.TTS == "" {
			tts = append(tts, p.Text)
		} else {
			tts = append(tts, p.TTS)
		}
	}
	res.Text = strings.Join(txt, "\n---\n")
	res.TTS = strings.Join(tts, " sil<[1000]> ")
}

func amongTokens(v string, tokens []string) bool {
	// using regexp instead of Levenshtein distance
	vb := []byte(v)
	for _, t := range tokens {
		re := regexp.MustCompile(`(?i)` + `(?:\A|\s|[[:punct:]])` + regexp.QuoteMeta(t) + `(?:\s|\z|[[:punct:]])`)
		if re.Match(vb) {
			return true
		}
	}
	return false
}
