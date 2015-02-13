package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	aqb "github.com/umsatz/go-aqbanking"
)

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

func loadPins(filename string) []aqb.Pin {
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

	var pins = make([]aqb.Pin, len(_pins))
	for i, p := range _pins {
		pins[i] = aqb.Pin(&p)
	}

	return pins
}

func main() {
	aq, err := aqb.DefaultAQBanking()
	if err != nil {
		log.Fatal("unable to init aqbanking: %v", err)
	}
	defer aq.Free()

	fmt.Printf("using aqbanking %d.%d.%d\n",
		aq.Version.Major,
		aq.Version.Minor,
		aq.Version.Patchlevel,
	)

	for _, pin := range loadPins("pins.json") {
		aq.RegisterPin(pin)
	}

	var user aqb.User
	if err := aq.AddPinTanUser(&user); err != nil {
		fmt.Printf("unable to add user. %v\n", err)
	} else {
		user.FetchAccounts(aq)
	}

	accountCollection, err := aq.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}
	fmt.Printf("found %d accounts.\n", len(accountCollection.Accounts))
}
