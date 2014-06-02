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

func (ab *AQBanking) Accounts() ([]Account, error) {
	var abAccountList *C.AB_ACCOUNT_LIST2 = C.AB_Banking_GetAccounts(ab.Ptr)
	if abAccountList == nil {
		return nil, errors.New("Unable to load accounts.")
	}

	var accounts []Account = make([]Account, C.AB_Account_List2_GetSize(abAccountList))

	var abIterator *C.AB_ACCOUNT_LIST2_ITERATOR = C.AB_Account_List2_First(abAccountList)
	if abIterator == nil {
		return nil, errors.New("Unable to get account iterator.")
	}

	var abAccount *C.AB_ACCOUNT
	abAccount = C.AB_Account_List2Iterator_Data(abIterator)

	for i := 0; abAccount != nil; i++ {
		account := Account{}
		account.AccountNumber = C.GoString(C.AB_Account_GetAccountNumber(abAccount))
		account.BankCode = C.GoString(C.AB_Account_GetBankCode(abAccount))
		accounts[i] = account
		abAccount = C.AB_Account_List2Iterator_Next(abIterator)
	}

	C.AB_Account_List2Iterator_free(abIterator)
	C.AB_Account_free(abAccount)
	C.AB_Account_List2_FreeAll(abAccountList)

	return accounts, nil
}

type Account struct {
	AccountNumber string
	BankCode      string
}

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
		fmt.Printf("kto: %v, blz: %v\n", account.AccountNumber, account.BankCode)
	}
}
