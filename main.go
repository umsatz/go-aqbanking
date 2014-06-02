package main

import (
	"errors"
	"fmt"
	"log"
)

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo CFLAGS: -I/usr/local/include/aqbanking5
#include <aqbanking5/aqbanking/banking.h>
*/
import "C"

type AQBankingVersion struct {
	Major      int
	Minor      int
	Patchlevel int
}

type AQBanking struct {
	Name    string
	Version AQBankingVersion

	Ptr *C.AB_BANKING
	gui *C.GWEN_GUI
}

func NewAQBanking(name string) (*AQBanking, error) {
	inst := &AQBanking{}
	inst.Name = name

	inst.Ptr = C.AB_Banking_new(C.CString(inst.Name), nil, 0)
	if err := C.AB_Banking_Init(inst.Ptr); err != 0 {
		return nil, errors.New(fmt.Sprintf("unable to initialized aqbanking: %d", err))
	}
	if err := C.AB_Banking_OnlineInit(inst.Ptr); err != 0 {
		return nil, errors.New(fmt.Sprintf("unable to initialized aqbanking: %d", err))
	}

	inst.gui = C.GWEN_Gui_new()
	C.GWEN_Gui_SetGui(inst.gui)

	inst.loadVersion()

	return inst, nil
}

func (ab *AQBanking) loadVersion() {
	var major, minor, patchlevel, build C.int
	C.AB_Banking_GetVersion(&major, &minor, &patchlevel, &build)
	ab.Version = AQBankingVersion{int(major), int(minor), int(patchlevel)}
}

func (ab *AQBanking) Free() error {
	if err := C.AB_Banking_OnlineFini(ab.Ptr); err != 0 {
		return errors.New(fmt.Sprintf("unable to free aqbanking online: %d\n", err))
	}

	C.AB_Banking_Fini(ab.Ptr)
	C.AB_Banking_free(ab.Ptr)
	return nil
}

func main() {
	//
	// SETUP
	//
	acc, err := NewAQBanking("golib")
	if err != nil {
		fmt.Printf("unable to init aqbanking: %v", err)
	}
	defer acc.Free()

	fmt.Printf("using aqbanking %d.%d.%d\n", acc.Version.Major, acc.Version.Minor, acc.Version.Patchlevel)

	// list known accounts
	var account_list *C.AB_ACCOUNT_LIST2
	account_list = C.AB_Banking_GetAccounts(acc.Ptr)
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

	fmt.Printf("Hello, World!\n")
}
