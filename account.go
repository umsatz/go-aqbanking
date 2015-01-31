package aqbanking

import (
	"errors"
	"fmt"
)

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo CFLAGS: -I/usr/local/include/aqbanking5
#include <aqbanking/banking.h>
*/
import "C"

// Account represents an aqbanking account
// Right now Paypal and CreditCards are not supported, even though
// aqbanking supports them.
type Account struct {
	Name          string
	AccountNumber string
	BIC           string
	IBAN          string
	Owner         string
	Currency      string
	Country       string
	Bank          Bank

	ptr *C.AB_ACCOUNT
}

// Bank represents a credit institute
type Bank struct {
	Name     string
	BankCode string
}

// Free frees the underlying aqbanking account pointer
func (a *Account) Free() {
	C.AB_Account_free(a.ptr)
}

// AccountCollection wraps working with multiple accounts, e.g. when searching by banking code.
// Necessary to support proper freeing of the underlying aqbanking collection pointer
type AccountCollection struct {
	Accounts []Account
}

// FirstUser returns the first user associated with a given account
func (a *Account) FirstUser() User {
	return newUser(C.AB_Account_GetFirstUser(a.ptr))
}

func newAccount(a *C.AB_ACCOUNT) Account {
	account := Account{
		ptr:           a,
		Name:          C.GoString(C.AB_Account_GetAccountName(a)),
		Owner:         C.GoString(C.AB_Account_GetOwnerName(a)),
		Currency:      C.GoString(C.AB_Account_GetCurrency(a)),
		Country:       C.GoString(C.AB_Account_GetCountry(a)),
		AccountNumber: C.GoString(C.AB_Account_GetAccountNumber(a)),
		IBAN:          C.GoString(C.AB_Account_GetIBAN(a)),
		BIC:           C.GoString(C.AB_Account_GetBIC(a)),
		Bank: Bank{
			Name:     C.GoString(C.AB_Account_GetBankName(a)),
			BankCode: C.GoString(C.AB_Account_GetBankCode(a)),
		},
	}

	return account
}

// Remove an Account from aqbanking files
func (a *Account) Remove(aq *AQBanking) error {
	if err := C.AB_Banking_DeleteAccount(aq.ptr, a.ptr); err != 0 {
		return fmt.Errorf("unable to delete account: %d\n", err)
	}
	return nil
}

// AccountsFor returns all accounts associated with a given user
func (ab *AQBanking) AccountsFor(u *User) (*AccountCollection, error) {
	allAccountCollection, err := ab.Accounts()
	if err != nil {
		return nil, err
	}

	list := &AccountCollection{}
	list.Accounts = make([]Account, 0)

	for _, account := range allAccountCollection.Accounts {
		accUser := account.FirstUser()
		if accUser.ID == u.ID {
			list.Accounts = append(list.Accounts, account)
		}
	}

	return list, nil
}

// Accounts returns all accounts registered with the given AQBanking instance
func (ab *AQBanking) Accounts() (*AccountCollection, error) {
	abAccountList := C.AB_Banking_GetAccounts(ab.ptr)
	if abAccountList == nil {
		// no accounts available
		return &AccountCollection{}, nil
	}

	list := &AccountCollection{}
	list.Accounts = make([]Account, C.AB_Account_List2_GetSize(abAccountList))

	abIterator := C.AB_Account_List2_First(abAccountList)
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
	C.AB_Account_List2_free(abAccountList)

	return list, nil
}
