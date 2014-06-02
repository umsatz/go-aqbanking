package main

import (
	"fmt"
	"log"
)

func listAccounts(ab *AQBanking) {
	accounts, err := ab.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}

	fmt.Println("%%\nAccounts")
	for _, account := range accounts {
		fmt.Printf(`## %v
Owner: %v
Currency: %v
Country: %v
AccountNumber: %v
BankCode: %v
Bank: %v
IBAN: %v
BIC: %v

`,
			account.Name,
			account.Owner,
			account.Currency,
			account.Country,
			account.AccountNumber,
			account.BankCode,
			account.Bank.Name,
			account.IBAN,
			account.BIC,
		)
	}
}

func listUsers(ab *AQBanking) {
	users, err := ab.Users()
	if err != nil {
		log.Fatal("unable to list users: %v", err)
	}

	fmt.Println("%%\nUsers")
	for _, user := range users {
		fmt.Printf(`## %v
Name: %v
UserId: %v
CustomerId: %v

`,
			user.Id,
			user.Name,
			user.UserId,
			user.CustomerId,
		)
	}
}

func main() {
	ab, err := NewAQBanking("local")
	if err != nil {
		log.Fatal("unable to init aqbanking: %v", err)
	}
	defer ab.Free()

	fmt.Printf("using aqbanking %d.%d.%d\n",
		ab.Version.Major,
		ab.Version.Minor,
		ab.Version.Patchlevel,
	)

	listAccounts(ab)
	listUsers(ab)
}
