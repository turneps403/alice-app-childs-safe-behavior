package main

import (
	"encoding/json"
	"fmt"
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
	Question Phrase   `json:"question"`
	Answers  []string `json:"answers"`
	Reaction Reaction `json:"reaction"`
}

type Place struct {
	Tokens  []string `json:"tokens"`
	Start   Phrase   `json:"start"`
	Riddles []Riddle `json:"riddles"`
}

type Dialog struct {
	Start  Phrase  `json:"start"`
	Places []Place `json:"place"`
}

var dialog *Dialog

func DialogInstance() *Dialog {
	if dialog == nil {
		dialog = &Dialog{}
	}
	return dialog
}

// func (d *Dialog) Start() Phrase {
// 	return d.Start
// }

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

	fmt.Printf("file was readed %v", d)
}

func main() {
	println(DialogInstance())
}
