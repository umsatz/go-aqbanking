package main

import (
	"fmt"
	"log"

	aqb "github.com/umsatz/go-aqbanking"
)

func main() {
	aq, err := aqb.DefaultAQBanking()
	// alternativly, customize the aqbanking path:
	// aq, err := aq.NewAQBanking("custom", "./tmp")
	if err != nil {
		log.Fatal("unable to init aqbanking: %v", err)
	}
	defer aq.Free()

	fmt.Printf("using aqbanking %d.%d.%d\n",
		aq.Version.Major,
		aq.Version.Minor,
		aq.Version.Patchlevel,
	)

	for _, pin := range aqb.LoadPins("pins.json") {
		aq.RegisterPin(&pin)
	}

	userCollection, err := aq.Users()
	if err != nil {
		log.Fatal("unable to list users: %v", err)
	}
	defer userCollection.Free()
	fmt.Printf("found %d users.\n", len(userCollection.Users))

	for _, user := range userCollection.Users {
		fmt.Printf("User %d\n", user.Id)
	}

	accountCollection, err := aq.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}

	fmt.Printf("found %d accounts.\n", len(accountCollection.Accounts))

	for _, account := range accountCollection.Accounts {
		fmt.Printf("Account %v\n", account.Name)

		transactions, _ := aq.Transactions(&account, nil, nil)
		fmt.Printf("%d transactions\n", len(transactions))
	}
}
