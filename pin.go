package main

import (
	"encoding/json"
	"log"
	"os"
)

// Pin is a interface to support pluggable pin loading.
// The examples read the pin from a pins.json, which
// is extremely insecure and should never be used in production
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

// LoadPins deserializes every pin specified in the given file, even if they might not
// contain all required attributes.
func LoadPins(filename string) []Pin {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("%v", err)
		return nil
	}

	var _pins []pin
	if err = json.NewDecoder(f).Decode(&_pins); err != nil {
		log.Fatal("%v", err)
		return nil
	}

	var pins = make([]Pin, len(_pins))
	for i, pin := range _pins {
		pins[i] = Pin(&pin)
	}

	return pins
}
