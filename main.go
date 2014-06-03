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
	list, err := ab.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}
	defer list.Free()

	fmt.Println("%%\nAccounts")
	for _, account := range list.Accounts {
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

func listTransactions(ab *AQBanking) {
	accountList, err := ab.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}
	defer accountList.Free()
	// account := accountList.Accounts[len(accountList.Accounts)-1]
	account := accountList.Accounts[0]

	transactions, err := ab.Transactions(account)
	if err != nil {
		log.Fatalf("unable to get transactions!: %v", err)
	}

	for _, transaction := range transactions {
		fmt.Printf(`## %v
Currency: %v
Total: %2.2f
`, transaction.Purpose,
			transaction.Currency,
			transaction.Total)
	}
}

func main() {
	var gui *C.struct_GWEN_GUI = C.GWEN_Gui_CGui_new()
	C.GWEN_Gui_SetGui(gui)

	fmt.Printf("%d", AccountTypeBank)
	// C.GWEN_Gui_AddFlags(gui, C.GWEN_GUI_FLAGS_NONINTERACTIVE)

	// fmt.Println("%d", gui.flags)

	// var fnc *C.GWEN_GUI_PRINT_FN
	// fmt.Println("%#v", C.GWEN_Gui_SetPrintFn)
	// C.GWEN_Gui_SetPrintFn(gui, &C.ASDPrint)

	// GWEN_Gui_SetCheckCertFn
	// GWEN_Gui_SetReadDialogPrefsFn
	// GWEN_Gui_SetWriteDialogPrefsFn
	// GWEN_Gui_SetRunDialogFn
	// GWEN_Gui_SetGetPasswordFn
	// GWEN_Gui_SetSetPasswordStatusFn
	// GWEN_Gui_SetPrintFn

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
	// listUsers(ab)
	// listTransactions(ab)
}
