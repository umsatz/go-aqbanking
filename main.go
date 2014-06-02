package main

import (
	"fmt"
	"log"
	"unsafe"
)

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo CFLAGS: -I/usr/local/include/aqbanking5
#include <aqbanking5/aqbanking/banking.h>
*/
import "C"

// type AQBanking struct {
// 	Name string
// 	Ptr  *C.AB_BANKING
// }

// func NewAQBanking(name string) *AQBanking {
// 	inst := &AQBanking{name, nil}
// }

func main() {
	var ab *C.AB_BANKING
	appName := C.CString("golib")
	defer C.free(unsafe.Pointer(appName))

	//
	// SETUP
	//
	ab = C.AB_Banking_new(appName, nil, 0)
	if err := C.AB_Banking_Init(ab); err != 0 {
		log.Fatal("unable to initialize aqbanking: %d\n", err)
	}
	if err := C.AB_Banking_OnlineInit(ab); err != 0 {
		log.Fatal("unable to initialize aqbanking online: %d\n", err)
	}

	var gui *C.GWEN_GUI
	gui = C.GWEN_Gui_new()
	C.GWEN_Gui_SetGui(gui)

	// list version, debug stuff
	var major, minor, patchlevel, build C.int
	C.AB_Banking_GetVersion(&major, &minor, &patchlevel, &build)
	fmt.Printf("using aqbanking %d.%d.%d\n", major, minor, patchlevel)

	// list known accounts
	var account_list *C.AB_ACCOUNT_LIST2
	account_list = C.AB_Banking_GetAccounts(ab)
	if account_list == nil {
		fmt.Println("Unable to load accounts.")
	}

	var iterator *C.AB_ACCOUNT_LIST2_ITERATOR
	iterator = C.AB_Account_List2_First(account_list)
	if iterator == nil {
		log.Fatal("Unable to get account iterator.")
	}

	var a *C.AB_ACCOUNT
	a = C.AB_Account_List2Iterator_Data(iterator)

	for a != nil {
		var account_number, bank_code string
		account_number = C.GoString(C.AB_Account_GetAccountNumber(a))
		bank_code = C.GoString(C.AB_Account_GetBankCode(a))
		fmt.Printf("kto: %v, blz: %v\n", account_number, bank_code)
		a = C.AB_Account_List2Iterator_Next(iterator)
	}

	C.AB_Account_List2Iterator_free(iterator)
	C.AB_Account_free(a)
	C.AB_Account_List2_FreeAll(account_list)

	//
	// TEARDOWN
	//
	// fmt.Printf("%d\n", ab.onlineInitCount)
	if err := C.AB_Banking_OnlineFini(ab); err != 0 {
		log.Fatal("unable to free aqbanking online: %d\n", err)
	}
	C.AB_Banking_Fini(ab)

	C.AB_Banking_free(ab)

	fmt.Printf("Hello, World!\n")
}
