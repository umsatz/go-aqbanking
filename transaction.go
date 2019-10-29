package aqbanking

import (
	"fmt"
	"time"
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

// Transaction represents an aqbanking transaction
type Transaction struct {
	Type              string
	Purpose           string
	Text              string
	Status            string
	Date              time.Time
	ValutaDate        time.Time
	CustomerReference string
	EndToEndReference string
	Total             float32
	TotalCurrency     string
	Fee               float32
	FeeCurrency       string

	MandateID     string
	BandReference string

	LocalBankCode      string
	LocalAccountNumber string
	LocalIBAN          string
	LocalBIC           string
	LocalName          string

	RemoteBankCode      string
	RemoteAccountNumber string
	RemoteIBAN          string
	RemoteBIC           string
	RemoteName          string
}

func newTransaction(t *C.AB_TRANSACTION) *Transaction {
	v := C.AB_Transaction_GetValue(t)

	if v == nil {
		return nil
	}

	transaction := Transaction{
		Type:              C.GoString(C.AB_Transaction_GetPurpose(t)),
		Purpose:           C.GoString(C.AB_Transaction_GetPurpose(t)),
		Text:              C.GoString(C.AB_Transaction_GetTransactionText(t)),
		Status:            C.GoString(C.AB_Transaction_Status_toString(C.AB_Transaction_GetStatus(t))),
		CustomerReference: C.GoString(C.AB_Transaction_GetCustomerReference(t)),
		EndToEndReference: C.GoString(C.AB_Transaction_GetEndToEndReference(t)),
		MandateID:         C.GoString(C.AB_Transaction_GetMandateId(t)),

		Date:       gwenDateToTime(C.AB_Transaction_GetDate(t)),
		ValutaDate: gwenDateToTime(C.AB_Transaction_GetValutaDate(t)),

		Total:         float32(C.AB_Value_GetValueAsDouble(v)),
		TotalCurrency: C.GoString(C.AB_Value_GetCurrency(v)),

		LocalIBAN:          C.GoString(C.AB_Transaction_GetLocalIban(t)),
		LocalBIC:           C.GoString(C.AB_Transaction_GetLocalBic(t)),
		LocalBankCode:      C.GoString(C.AB_Transaction_GetLocalBankCode(t)),
		LocalAccountNumber: C.GoString(C.AB_Transaction_GetLocalAccountNumber(t)),
		LocalName:          C.GoString(C.AB_Transaction_GetLocalName(t)),

		RemoteIBAN:          C.GoString(C.AB_Transaction_GetRemoteIban(t)),
		RemoteBIC:           C.GoString(C.AB_Transaction_GetRemoteBic(t)),
		RemoteBankCode:      C.GoString(C.AB_Transaction_GetRemoteBankCode(t)),
		RemoteAccountNumber: C.GoString(C.AB_Transaction_GetRemoteAccountNumber(t)),
		RemoteName:          C.GoString(C.AB_Transaction_GetRemoteName(t)),
	}

	if fees := C.AB_Transaction_GetFees(t); fees != nil {
		transaction.Fee = float32(C.AB_Value_GetValueAsDouble(fees))
		transaction.FeeCurrency = C.GoString(C.AB_Value_GetCurrency(fees))
	}

	return &transaction
}

// Transactions implements AB_TransactionGetTransactions_new from aqbanking, listing
// all transactions from a given aqbanking instance
func (ab *AQBanking) Transactions(acc *Account, from *time.Time, to *time.Time) ([]Transaction, error) {

	// create a list to which banking commands are added
	cmdList := C.AB_Transaction_List2_new()
	defer C.AB_Transaction_List2_free(cmdList)

	// create an online banking command
	t := C.AB_Transaction_new()
	C.AB_Transaction_SetCommand(t, C.AB_Transaction_CommandGetTransactions)
	C.AB_Transaction_SetUniqueAccountId(t, C.uint(acc.ID))

	if from != nil {
		C.AB_Transaction_SetFirstDate(t, (*C.GWEN_DATE)(newGwenDate(*from)))
	}
	if to != nil {
		C.AB_Transaction_SetLastDate(t, (*C.GWEN_DATE)(newGwenDate(*to)))
	}

	// add command to the list
	C.AB_Transaction_List2_PushBack(cmdList, t)

	ctx := C.AB_ImExporterContext_new()
	defer C.AB_ImExporterContext_free(ctx)

	if err := C.AB_Banking_SendCommands(ab.ptr, cmdList, ctx); err < 0 {
		return nil, newError("unable to send command", err)
	}

	ai := C.AB_ImExporterContext_GetFirstAccountInfo(ctx)

	if ai == nil {
		return nil, fmt.Errorf("unable to get first account info")
	}

	var transactions []Transaction
	for ai != nil {
		t = C.AB_ImExporterAccountInfo_GetFirstTransaction(ai, 0, 0)

		for t != nil {
			if transaction := newTransaction(t); transaction != nil {
				transactions = append(transactions, *transaction)
			}

			t = C.AB_Transaction_List_Next(t)
		}
		ai = C.AB_ImExporterAccountInfo_List_Next(ai)
	}

	return transactions, nil
}

// AllTransactions implements AB_TransactionGetTransactions_new without filter
func (ab *AQBanking) AllTransactions(acc *Account) ([]Transaction, error) {
	return ab.Transactions(acc, nil, nil)
}
