package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
		fmt.Printf(`
## %v
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

type Pin struct {
	Kto string `json:"kto"`
	Blz string `json:"blz"`
	Pin string `json:"pin"`
}

func loadPins() []Pin {
	f, err := os.Open("pins.json")
	if err != nil {
		log.Fatal("%v", err)
		return nil
	}

	var pins []Pin
	err = json.NewDecoder(f).Decode(&pins)
	if err != nil {
		log.Fatal("%v", err)
		return nil
	}

	return pins
}

func registerPins(aq *AQBanking, gui *C.struct_GWEN_GUI) {
	accountCollection, _ := aq.Accounts()
	pins := loadPins()
	var dbPins *C.GWEN_DB_NODE = C.GWEN_DB_Group_new(C.CString("pins"))

	for _, account := range accountCollection.Accounts {
		for _, pin := range pins {
			if pin.Blz == account.BankCode && pin.Kto == account.AccountNumber {
				user := account.FirstUser()
				str := fmt.Sprintf("PIN_%v_%v=%v\n", pin.Blz, user.CustomerId, pin.Pin)
				pinLen := len(str)

				C.GWEN_DB_ReadFromString(dbPins, C.CString(str), C.int(pinLen), C.GWEN_PATH_FLAGS_CREATE_GROUP|C.GWEN_DB_FLAGS_DEFAULT)
				break
			}
		}
	}

	C.GWEN_Gui_CGui_SetPasswordDb(gui, dbPins, 1)
}

func main() {
	var gui *C.struct_GWEN_GUI = C.GWEN_Gui_CGui_new()
	defer C.GWEN_Gui_free(gui)
	// var gui *C.struct_GWEN_GUI = C.GWEN_Gui_new()
	C.GWEN_Gui_SetFlags(gui, C.GWEN_GUI_FLAGS_ACCEPTVALIDCERTS|C.GWEN_GUI_FLAGS_NONINTERACTIVE)
	C.GWEN_Gui_SetGui(gui)

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

	registerPins(aq, gui)

	listAccounts(aq)
	listUsers(aq)
	listTransactions(aq)
}
