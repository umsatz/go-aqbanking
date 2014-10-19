package main

import (
	"encoding/json"
	"log"
	"os"
)

type Pin interface {
	BankCode() string
	UserID() string
	Pin() string
}

type pin struct {
	Blz string `json:"blz"`
	UID string `json:"uid"`
	PIN string `json:"pin"`
}

func (p *pin) BankCode() string {
	return p.Blz
}

func (p *pin) UserID() string {
	return p.UID
}
func (p *pin) Pin() string {
	return p.PIN
}

func LoadPins(filename string) []pin {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("%v", err)
		return nil
	}

	var pins []pin
	err = json.NewDecoder(f).Decode(&pins)
	if err != nil {
		log.Fatal("%v", err)
		return nil
	}

	return pins
}
