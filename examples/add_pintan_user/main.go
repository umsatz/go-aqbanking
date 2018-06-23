package main

import (
	"fmt"
	"log"

	aqb "github.com/umsatz/go-aqbanking"
	examples "github.com/umsatz/go-aqbanking/examples"
)

func main() {
	aq, err := aqb.DefaultAQBanking()
	if err != nil {
		log.Fatalf("unable to init aqbanking: %v", err)
	}
	defer aq.Free()

	fmt.Printf("using aqbanking %d.%d.%d\n",
		aq.Version.Major,
		aq.Version.Minor,
		aq.Version.Patchlevel,
	)

	for _, pin := range examples.LoadPins("pins.json") {
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
		log.Fatalf("unable to list accounts: %v", err)
	}
	fmt.Printf("found %d accounts.\n", len(accountCollection.Accounts))
}
