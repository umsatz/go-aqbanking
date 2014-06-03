package main

import (
	"fmt"
	"log"
)

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking5
#include <gwenhywfar/cgui.h>
#include <aqbanking/banking.h>
*/
import "C"

func listAccounts(ab *AQBanking) {
	accountCollection, err := ab.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}
	defer accountCollection.Free()

	fmt.Println("%%\nAccounts")
	for _, account := range accountCollection.Accounts {
		fmt.Printf(`
## %v
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
	userCollection, err := ab.Users()
	if err != nil {
		log.Fatal("unable to list users: %v", err)
	}
	defer userCollection.Free()

	fmt.Println("%%\nUsers")
	for _, user := range userCollection.Users {
		fmt.Printf(`
## %v
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

func listTransactionsFor(ab *AQBanking, account *Account) {
	transactions, err := ab.Transactions(*account)
	if err != nil {
		log.Fatalf("unable to get transactions!: %v", err)
	}

	for _, t := range transactions {
		fmt.Printf(`
## %v
'%v'
Status: %v
MandateReference: %v
CustomerReference: %v
LocalBankCode: %v
LocalAccountNumber: %v
LocalIBAN: %v
LocalBIC: %v
LocalName: %v
RemoteBankCode: %v
RemoteAccountNumber: %v
RemoteIBAN: %v
RemoteBIC: %v
RemoteName: %v
Date: %v
ValutaDate: %v
Total: %2.2f %v
Fee: %2.2f %v
`, t.Purpose,
			t.Text,
			t.Status,
			t.MandateReference,
			t.CustomerReference,
			t.LocalBankCode,
			t.LocalAccountNumber,
			t.LocalIBAN,
			t.LocalBIC,
			t.LocalName,
			t.RemoteBankCode,
			t.RemoteAccountNumber,
			t.RemoteIBAN,
			t.RemoteBIC,
			t.RemoteName,
			t.Date,
			t.ValutaDate,
			t.Total, t.TotalCurrency,
			t.Fee, t.FeeCurrency,
		)
	}
}

func listTransactions(ab *AQBanking) {
	accountList, err := ab.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}
	defer accountList.Free()

	for _, account := range accountList.Accounts {
		listTransactionsFor(ab, &account)
	}
}

func main() {
	gui := NewNonInteractiveGui()
	defer gui.Free()

	aq, err := NewAQBanking("local")
	if err != nil {
		log.Fatal("unable to init aqbanking: %v", err)
	}
	defer aq.Free()

	fmt.Printf("using aqbanking %d.%d.%d\n",
		aq.Version.Major,
		aq.Version.Minor,
		aq.Version.Patchlevel,
	)

	gui.RegisterPins(aq, LoadPins("pins.json"))

	listUsers(aq)
	listAccounts(aq)
	listTransactions(aq)
}
