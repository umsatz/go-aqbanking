package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	aqb "github.com/umsatz/go-aqbanking"
)

func loadPins(filename string) []aqb.Pin {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var _pins []pin
	if err = json.NewDecoder(f).Decode(&_pins); err != nil {
		log.Fatal(err)
		return nil
	}

	var pins = make([]aqb.Pin, len(_pins))
	for i, p := range _pins {
		pins[i] = aqb.Pin(&p)
	}

	return pins
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

func main() {
	aq, err := aqb.DefaultAQBanking()
	// alternativly, customize the aqbanking path:
	// aq, err := aq.NewAQBanking("custom", "./tmp")
	if err != nil {
		log.Fatalf("unable to init aqbanking: %v", err)
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

	userCollection, err := aq.Users()
	if err != nil {
		log.Fatalf("unable to list users: %v", err)
	}
	defer userCollection.Free()
	fmt.Printf("found %d users.\n", len(userCollection.Users))

	for _, user := range userCollection.Users {
		fmt.Printf("User %d\n", user.ID)
	}

	accountCollection, err := aq.Accounts()
	if err != nil {
		log.Fatalf("unable to list accounts: %v", err)
	}

	fmt.Printf("found %d accounts.\n", len(accountCollection.Accounts))

	for _, account := range accountCollection.Accounts {
		fmt.Printf("Account %v\n", account.Name)

		transactions, _ := aq.Transactions(&account, nil, nil)
		fmt.Printf("%d transactions\n", len(transactions))
	}
}
