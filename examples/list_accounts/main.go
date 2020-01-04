package main

import (
	"fmt"
	"log"

	aqb "github.com/umsatz/go-aqbanking"
	"github.com/umsatz/go-aqbanking/examples"
)

func main() {
	aq, err := aqb.DefaultAQBanking()
	// alternativly, customize the aqbanking path:
	// aq, err := aq.NewAQBanking("custom", "./tmp")
	if err != nil {
		log.Fatalf("unable to init aqbanking: %v", err)
	}
	defer aq.Free()

	fmt.Println("using aqbanking", aq.Version)

	for _, pin := range examples.LoadPins("pins.json") {
		aq.RegisterPin(pin)
	}

	accountCollection, err := aq.Accounts()
	if err != nil {
		log.Fatalf("unable to list accounts: %v", err)
	}

	fmt.Printf("found %d accounts.\n", len(accountCollection))

	for _, account := range accountCollection {
		fmt.Printf("Account %v\n", account.Name)

		transactions, _ := aq.Transactions(&account, nil, nil)
		fmt.Printf("%d transactions\n", len(transactions))
	}
}
