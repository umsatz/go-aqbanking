package aqbanking_test

import (
	"fmt"
	"log"

	"github.com/umsatz/go-aqbanking"
)

// Returns the default aqbanking instance. Files are normally stored at $HOME/.aqbanking
func ExampleAQBanking_initialization() {
	var aq, err = aqbanking.DefaultAQBanking()
	if err != nil {
		fmt.Printf("Failed to initializate aqbanking instance %q", err)
	}
	defer aq.Free()
}

func ExampleAQBanking_addUser() {
	var aq, err = aqbanking.DefaultAQBanking()
	if err != nil {
		fmt.Printf("Failed to initializate aqbanking instance %q", err)
	}
	defer aq.Free()

	fmt.Printf("using aqbanking %d.%d.%d\n",
		aq.Version.Major,
		aq.Version.Minor,
		aq.Version.Patchlevel,
	)

	user := aqbanking.User{
		ID:          0,
		UserID:      "userid",
		CustomerID:  "userid",
		BankCode:    "bankcode",
		Name:        "name of bank",
		ServerURI:   "https://your hbci server url",
		HbciVersion: 300,
	}
	if err = aq.AddPinTanUser(&user); err != nil {
		fmt.Printf("unable to add user. %v\n", err)
	}

	// Next a call to aq.RegisterPin() is required to allow PinTan authentication
}

func ExampleAQBanking_listAccounts() {
	var aq, err = aqbanking.DefaultAQBanking()
	if err != nil {
		fmt.Printf("Failed to initializate aqbanking instance %q", err)
	}
	defer aq.Free()

	accountCollection, err := aq.Accounts()
	if err != nil {
		log.Fatalf("unable to list accounts: %v", err)
	}

	for _, account := range accountCollection {
		fmt.Printf("%v", account)
	}
}

func ExampleAQBanking_listUsers() {
	var aq, err = aqbanking.DefaultAQBanking()
	if err != nil {
		fmt.Printf("Failed to initializate aqbanking instance %q", err)
	}
	defer aq.Free()

	userCollection, err := aq.Users()
	if err != nil {
		log.Fatalf("unable to list users: %v", err)
	}
	defer userCollection.Free()

	for _, user := range userCollection.Users {
		fmt.Printf("%v", user)
	}
}
