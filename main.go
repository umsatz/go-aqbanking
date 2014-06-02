package main

import (
	"fmt"
	"log"
)

func main() {
	acc, err := NewAQBanking("golib")
	if err != nil {
		log.Fatal("unable to init aqbanking: %v", err)
	}
	defer acc.Free()

	fmt.Printf("using aqbanking %d.%d.%d\n", acc.Version.Major, acc.Version.Minor, acc.Version.Patchlevel)

	accounts, err := acc.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}
	for _, account := range accounts {
		fmt.Printf(`Owner: %v
Name: %v
Currency: %v
Country: %v
AccountNumber: %v
BankCode: %v
Bank: %v
IBAN: %v

`,
			account.Owner,
			account.Name,
			account.Currency,
			account.Country,
			account.AccountNumber,
			account.BankCode,
			account.Bank.Name,
			account.IBAN,
		)
	}
}
