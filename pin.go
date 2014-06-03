package main

import (
	"encoding/json"
	"log"
	"os"
)

type Pin struct {
	Blz    string `json:"blz"`
	UserId string `json:"uid"`
	Pin    string `json:"pin"`
}

func LoadPins(filename string) []Pin {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("%v", err)
		return nil
	}

	var pins []Pin
	err = json.NewDecoder(f).Decode(&pins)
	if err != nil {
		log.Fatal("%v", err)
		return nil
	}

	return pins
}
