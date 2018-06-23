package examples

import (
	"encoding/json"
	"log"
	"os"

	aqb "github.com/umsatz/go-aqbanking"
)

// Pin stores bank credentials
type Pin struct {
	Blz string `json:"blz"`
	UID string `json:"uid"`
	PIN string `json:"pin"`
}

// LoadPins loads pins from a JSON file
func LoadPins(filename string) []aqb.Pin {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var _pins []Pin
	if err = json.NewDecoder(f).Decode(&_pins); err != nil {
		log.Fatal(err)
		return nil
	}

	var pins = make([]aqb.Pin, len(_pins))
	for i, pin := range _pins {
		pins[i] = aqb.Pin(&pin)
	}

	return pins
}

// BankCode returns the BankCode
func (p *Pin) BankCode() string {
	return p.Blz
}

// UserID returns the UserID
func (p *Pin) UserID() string {
	return p.UID
}

// Pin returns the Pin
func (p *Pin) Pin() string {
	return p.PIN
}
