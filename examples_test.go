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
