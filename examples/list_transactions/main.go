package main

import (
	"fmt"
	"log"
	"time"

	aqb "github.com/umsatz/go-aqbanking"
	examples "github.com/umsatz/go-aqbanking/examples"
)

func main() {
	aq, err := aqb.NewAQBanking("custom", "./tmp")
	if err != nil {
		log.Fatalf("unable to init aqbanking: %v", err)
	}
	defer aq.Free()

	fmt.Println("using aqbanking", aq.Version)

	for _, pin := range examples.LoadPins("pins.json") {
		aq.RegisterPin(pin)
	}

	listAccounts(aq)
	listTransactions(aq)
}

func listAccounts(ab *aqb.AQBanking) {
	accountCollection, err := ab.Accounts()
	if err != nil {
		log.Fatalf("unable to list accounts: %v", err)
	}

	fmt.Println("%%\nAccounts")
	for _, account := range accountCollection {
		fmt.Printf(`
## %v
Owner: %v
Currency: %v
Country: %v
AccountNumber: %v
BankCode: %v
IBAN: %v
BIC: %v
`,
			account.Name,
			account.Owner,
			account.Currency,
			account.Country,
			account.AccountNumber,
			account.BankCode,
			account.IBAN,
			account.BIC,
		)
	}
}

func listTransactionsFor(ab *aqb.AQBanking, account *aqb.Account) {
	fromDate := time.Date(2014, 05, 14, 0, 0, 0, 0, time.UTC)
	toDate := time.Date(2014, 05, 16, 0, 0, 0, 0, time.UTC)
	transactions, err := ab.Transactions(account, &fromDate, &toDate)
	// or
	// transactions, err := ab.AllTransactions(account)
	if err != nil {
		log.Fatalf("unable to get transactions!: %v", err)
	}

	for _, t := range transactions {
		fmt.Printf(`
## %v
'%v'
Purpose: %v
Status: %v
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
Value: %2.2f %v
Fee: %2.2f %v
`, t.Type,
			t.Text,
			t.Purpose,
			t.Status,
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
			t.Value.Amount, t.Value.Currency,
			t.Fee.Amount, t.Fee.Currency,
		)
	}
}

func listTransactions(ab *aqb.AQBanking) {
	accounts, err := ab.Accounts()
	if err != nil {
		log.Fatalf("unable to list accounts: %v", err)
	}

	for _, account := range accounts {
		listTransactionsFor(ab, &account)
	}
}
