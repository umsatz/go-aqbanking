package main

import (
	"fmt"
	"log"

	aqb "github.com/umsatz/go-aqbanking"
)

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

	for _, pin := range aqb.LoadPins("pins.json") {
		aq.RegisterPin(&pin)
	}

	user := aqb.User{Id: 0, UserId: "userid", CustomerId: "userid", BankCode: "bankcode", Name: "name of bank", ServerUri: "https://your hbci server url", HbciVersion: 300}

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
