package main

import (
	"errors"
	"fmt"
)

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking5
#include <aqbanking/banking.h>
*/
import "C"

type Account struct {
	Name          string
	AccountNumber string
	BankCode      string
	BIC           string
	IBAN          string
	Owner         string
	Currency      string
	Country       string
	Bank          Bank

	ptr *C.AB_ACCOUNT
}

type Bank struct {
	Name string
}

func (a *Account) Free() {
	C.AB_Account_free(a.ptr)
}

type AccountCollection struct {
	Accounts []Account
	ptr      *C.AB_ACCOUNT_LIST2
}

func (al *AccountCollection) Free() {
	for i, _ := range al.Accounts {
		al.Accounts[i].Free()
	}

	al.Accounts = make([]Account, 0)
	C.AB_Account_List2_free(al.ptr)
}

func (a *Account) FirstUser() User {
	return newUser(C.AB_Account_GetFirstUser(a.ptr))
}

func newAccount(a *C.AB_ACCOUNT) Account {
	account := Account{}

	account.Name = C.GoString(C.AB_Account_GetAccountName(a))
	account.Owner = C.GoString(C.AB_Account_GetOwnerName(a))
	account.Currency = C.GoString(C.AB_Account_GetCurrency(a))
	account.Country = C.GoString(C.AB_Account_GetCountry(a))

	account.BankCode = C.GoString(C.AB_Account_GetBankCode(a))
	account.AccountNumber = C.GoString(C.AB_Account_GetAccountNumber(a))
	account.IBAN = C.GoString(C.AB_Account_GetIBAN(a))
	account.BIC = C.GoString(C.AB_Account_GetBIC(a))

	account.Bank = Bank{}
	account.Bank.Name = C.GoString(C.AB_Account_GetBankName(a))
	account.ptr = a

	return account
}

func (a *Account) Remove(aq *AQBanking) error {
	if err := C.AB_Banking_DeleteAccount(aq.ptr, a.ptr); err != 0 {
		return errors.New(fmt.Sprintf("unable to delete account: %d\n", err))
	}
	return nil
}

func (ab *AQBanking) AccountsFor(u *User) (*AccountCollection, error) {
	allAccountCollection, err := ab.Accounts()
	if err != nil {
		return nil, err
	}
	defer allAccountCollection.Free()

	var list *AccountCollection = &AccountCollection{}
	list.Accounts = make([]Account, 0)

	for _, account := range allAccountCollection.Accounts {
		accUser := account.FirstUser()
		if accUser.Id == u.Id {
			list.Accounts = append(list.Accounts, account)
		}
	}

	return list, nil
}

// implements AB_Banking_GetAccounts
func (ab *AQBanking) Accounts() (*AccountCollection, error) {
	var abAccountList *C.AB_ACCOUNT_LIST2 = C.AB_Banking_GetAccounts(ab.ptr)
	if abAccountList == nil {
		// no accounts available
		return &AccountCollection{}, nil
	}

	var list *AccountCollection = &AccountCollection{}
	list.Accounts = make([]Account, C.AB_Account_List2_GetSize(abAccountList))
	list.ptr = abAccountList

	var abIterator *C.AB_ACCOUNT_LIST2_ITERATOR = C.AB_Account_List2_First(abAccountList)
	if abIterator == nil {
		return nil, errors.New("Unable to get account iterator.")
	}

	var abAccount *C.AB_ACCOUNT
	abAccount = C.AB_Account_List2Iterator_Data(abIterator)

	for i := 0; abAccount != nil; i++ {
		list.Accounts[i] = newAccount(abAccount)

		abAccount = C.AB_Account_List2Iterator_Next(abIterator)
	}

	C.AB_Account_List2Iterator_free(abIterator)
	C.AB_Account_free(abAccount)

	return list, nil
}
