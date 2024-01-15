package main

import (
	"encoding/json"
	"log"
	"os"
)

const DIALOG_FILE = "dialog.json"

type Phrase struct {
	Text string `json:"text"`
	TTS  string `json:"tts"`
}

type Reaction struct {
	Right       Phrase `json:"right"`
	Wrong       Phrase `json:"wrong"`
	Explanation Phrase `json:"explanation"`
}

type Riddle struct {
	Title    string   `json:"title"`
	Question Phrase   `json:"question"`
	Answers  []string `json:"answers"`
	Reaction Reaction `json:"reaction"`
}

type Place struct {
	Tokens   []string `json:"tokens"`
	Riddles  []Riddle `json:"riddles"`
	Prologue Phrase   `json:"prologue"`
	Epilogue Phrase   `json:"epilogue"`
}

type Dialog struct {
	Start           Phrase  `json:"start"`
	Places          []Place `json:"place"`
	PlacesAvaliable Phrase  `json:"places_avaliable"`
	Fail            Phrase  `json:"fail"`
}

var dialog *Dialog

func DialogInstance() *Dialog {
	if dialog == nil {
		dialog = &Dialog{}
	}
	return dialog
}

func init() {
	fileBytes, err := os.ReadFile(DIALOG_FILE)
	if err != nil {
		log.Fatal(err)
	}

	d := DialogInstance()
	err = json.Unmarshal(fileBytes, d)
	if err != nil {
		log.Fatal(err)
	}
}
