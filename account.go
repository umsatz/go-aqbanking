package aqbanking

import (
	"errors"
)

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar5
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking6
#cgo linux CFLAGS: -I/usr/include/gwenhywfar5
#cgo linux CFLAGS: -I/usr/include/aqbanking6
#include <aqbanking/banking.h>
*/
import "C"

// Account represents an aqbanking account
// Right now Paypal and CreditCards are not supported, even though
// aqbanking supports them.
type Account struct {
	ID               int // ID is the unique account ID
	Name             string
	AccountNumber    string
	SubAccountNumber string
	BankCode         string
	IBAN             string
	BIC              string
	Owner            string
	Currency         string
	Country          string
	BackendName      string

	ptr *C.AB_ACCOUNT_SPEC
}

// Free frees the underlying aqbanking account pointer
func (a *Account) Free() {
	C.AB_AccountSpec_free(a.ptr)
}

// AccountCollection wraps working with multiple accounts, e.g. when searching by banking code.
// Necessary to support proper freeing of the underlying aqbanking collection pointer
type AccountCollection []Account

func newAccount(a *C.AB_ACCOUNT_SPEC) Account {
	return Account{
		ptr:              a,
		ID:               int(C.AB_AccountSpec_GetUniqueId(a)),
		Name:             C.GoString(C.AB_AccountSpec_GetAccountName(a)),
		Owner:            C.GoString(C.AB_AccountSpec_GetOwnerName(a)),
		IBAN:             C.GoString(C.AB_AccountSpec_GetIban(a)),
		BIC:              C.GoString(C.AB_AccountSpec_GetBic(a)),
		Currency:         C.GoString(C.AB_AccountSpec_GetCurrency(a)),
		Country:          C.GoString(C.AB_AccountSpec_GetCountry(a)),
		AccountNumber:    C.GoString(C.AB_AccountSpec_GetAccountNumber(a)),
		SubAccountNumber: C.GoString(C.AB_AccountSpec_GetSubAccountNumber(a)),
		BankCode:         C.GoString(C.AB_AccountSpec_GetBankCode(a)),
		BackendName:      C.GoString(C.AB_AccountSpec_GetBackendName(a)),
	}
}

// Remove an Account from aqbanking files
func (a *Account) Remove(aq *AQBanking) error {
	return errors.New("not implemented")
}

// Accounts returns all accounts registered with the given AQBanking instance
func (ab *AQBanking) Accounts() (AccountCollection, error) {
	asl := C.AB_AccountSpec_List_new()
	defer C.AB_AccountSpec_List_free(asl)

	rv := C.AB_Banking_GetAccountSpecList(ab.ptr, &asl)
	if rv < 0 {
		// no accounts available
		return nil, nil
	}

	list := make(AccountCollection, 0, C.AB_AccountSpec_List_GetCount(asl))

	as := C.AB_AccountSpec_List_First(asl)
	for as != nil {
		list = append(list, newAccount(as))
		as = C.AB_AccountSpec_List_Next(as)
	}

	return list, nil
}
