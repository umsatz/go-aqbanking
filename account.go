package main

import "errors"

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking5
#include <aqbanking/banking.h>
*/
import "C"

type AccountType int

const (
	AccountTypeUnknown     AccountType = iota
	AccountTypeBank        AccountType = iota
	AccountTypeCreditCard  AccountType = iota
	AccountTypeChecking    AccountType = iota
	AccountTypeSavings     AccountType = iota
	AccountTypeInvestment  AccountType = iota
	AccountTypeCash        AccountType = iota
	AccountTypeMoneyMarket AccountType = iota
)

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
	Type          AccountType

	Ptr *C.AB_ACCOUNT
}

type Bank struct {
	Name string
}

type AccountCollection struct {
	Accounts []Account
	Ptr      *C.AB_ACCOUNT_LIST2
}

func (al *AccountCollection) Free() {
	al.Accounts = make([]Account, 0)
	C.AB_Account_List2_free(al.Ptr)
}

// implements AB_Banking_GetAccounts
func (ab *AQBanking) Accounts() (*AccountCollection, error) {
	var abAccountList *C.AB_ACCOUNT_LIST2 = C.AB_Banking_GetAccounts(ab.Ptr)
	if abAccountList == nil {
		return nil, errors.New("Unable to load accounts.")
	}

	var list *AccountCollection = &AccountCollection{}
	list.Accounts = make([]Account, C.AB_Account_List2_GetSize(abAccountList))
	list.Ptr = abAccountList

	var abIterator *C.AB_ACCOUNT_LIST2_ITERATOR = C.AB_Account_List2_First(abAccountList)
	if abIterator == nil {
		return nil, errors.New("Unable to get account iterator.")
	}

	var abAccount *C.AB_ACCOUNT
	abAccount = C.AB_Account_List2Iterator_Data(abIterator)

	for i := 0; abAccount != nil; i++ {
		account := Account{}

		account.Name = C.GoString(C.AB_Account_GetAccountName(abAccount))
		account.Owner = C.GoString(C.AB_Account_GetOwnerName(abAccount))
		account.Currency = C.GoString(C.AB_Account_GetCurrency(abAccount))
		account.Country = C.GoString(C.AB_Account_GetCountry(abAccount))

		account.BankCode = C.GoString(C.AB_Account_GetBankCode(abAccount))
		account.AccountNumber = C.GoString(C.AB_Account_GetAccountNumber(abAccount))
		account.IBAN = C.GoString(C.AB_Account_GetIBAN(abAccount))
		account.BIC = C.GoString(C.AB_Account_GetBIC(abAccount))
		account.Type = AccountType(C.AB_Account_GetAccountType(abAccount))

		account.Bank = Bank{}
		account.Bank.Name = C.GoString(C.AB_Account_GetBankName(abAccount))
		account.Ptr = abAccount

		list.Accounts[i] = account
		abAccount = C.AB_Account_List2Iterator_Next(abIterator)
	}

	C.AB_Account_List2Iterator_free(abIterator)
	C.AB_Account_free(abAccount)

	return list, nil
}
