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
#include <aqbanking/abgui.h>
int ASDPrint(GWEN_GUI *gui,
			const char *docTitle,
			const char *docType,
			const char *descr,
			const char *text,
			uint32_t guiid){

  return 0;
}
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
		fmt.Printf(`## %v
Owner: %v
Type: %d
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
			account.Type,
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

func listTransactions(ab *AQBanking) {
	accountList, err := ab.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}
	defer accountList.Free()
	account := accountList.Accounts[len(accountList.Accounts)-1]
	// account := accountList.Accounts[0]

	transactions, err := ab.Transactions(account)
	if err != nil {
		log.Fatalf("unable to get transactions!: %v", err)
	}

	for _, transaction := range transactions {
		fmt.Printf(`## %v
'%v'
MandateReference: %v
CustomerReference: %v
Date: %v
ValutaDate: %v
Category: %v
Period: %v
Type: %v
SubType: %v
Currency: %v
Total: %2.2f
`, transaction.Purpose,
			transaction.Text,
			transaction.MandateReference,
			transaction.CustomerReference,
			transaction.Date,
			transaction.ValutaDate,
			transaction.Category,
			transaction.TransactionPeriod,
			transaction.Type,
			transaction.SubType,
			transaction.Currency,
			transaction.Total)
	}
}

func main() {
	var gui *C.struct_GWEN_GUI = C.GWEN_Gui_CGui_new()
	// var gui *C.struct_GWEN_GUI = C.GWEN_Gui_new()
	C.GWEN_Gui_SetGui(gui)

	// C.GWEN_Gui_SetFlags(gui, C.GWEN_GUI_FLAGS_NONINTERACTIVE)

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

	// listAccounts(aq)
	listUsers(aq)

	// userCollection, err := aq.Users()
	// if err != nil {
	// 	log.Fatal("unable to list users: %v", err)
	// }
	// defer userCollection.Free()

	// user := userCollection.Users[0]
	// var pw *C.char = C.CString("123456")
	// C.GWEN_DB_SetCharValue(user.Ptr, C.GWEN_DB_FLAGS_OVERWRITE_VARS, C.CString("password"), pw)

	listTransactions(aq)
}
