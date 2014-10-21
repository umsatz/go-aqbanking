package main

import (
	"errors"
	"fmt"
	"time"
)

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking5
#include <aqbanking/jobgettransactions.h>
#include <aqbanking/banking.h>
#include <aqbanking/job.h>
#include <aqbanking/banking_ob.h>
*/
import "C"

// Transaction represents an aqbanking transaction
type Transaction struct {
	Purpose           string
	Text              string
	Status            string
	Date              time.Time
	ValutaDate        time.Time
	MandateReference  string
	CustomerReference string
	Total             float32
	TotalCurrency     string
	Fee               float32
	FeeCurrency       string

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

func newTransaction(t *C.AB_TRANSACTION) (Transaction, bool) {
	var v *C.AB_VALUE
	v = C.AB_Transaction_GetValue(t)

	if v == nil {
		return Transaction{}, false
	}

	transaction := Transaction{}

	transaction.Purpose = (*gwStringList)(C.AB_Transaction_GetPurpose(t)).toString()
	transaction.Text = C.GoString(C.AB_Transaction_GetTransactionText(t))
	transaction.Status = C.GoString(C.AB_Transaction_Status_toString(C.AB_Transaction_GetStatus(t)))
	transaction.MandateReference = C.GoString(C.AB_Transaction_GetMandateReference(t))
	transaction.CustomerReference = C.GoString(C.AB_Transaction_GetCustomerReference(t))
	transaction.Date = (*gwTime)(C.AB_Transaction_GetDate(t)).toTime()
	transaction.ValutaDate = (*gwTime)(C.AB_Transaction_GetValutaDate(t)).toTime()

	transaction.Total = float32(C.AB_Value_GetValueAsDouble(v))
	transaction.TotalCurrency = C.GoString(C.AB_Value_GetCurrency(v))

	f := C.AB_Transaction_GetFees(t)
	if f != nil {
		transaction.Fee = float32(C.AB_Value_GetValueAsDouble(f))
		transaction.FeeCurrency = C.GoString(C.AB_Value_GetCurrency(f))
	}

	transaction.LocalIBAN = C.GoString(C.AB_Transaction_GetLocalIban(t))
	transaction.LocalBIC = C.GoString(C.AB_Transaction_GetLocalBic(t))
	transaction.LocalBankCode = C.GoString(C.AB_Transaction_GetLocalBankCode(t))
	transaction.LocalAccountNumber = C.GoString(C.AB_Transaction_GetLocalAccountNumber(t))
	transaction.LocalName = C.GoString(C.AB_Transaction_GetLocalName(t))

	transaction.RemoteIBAN = C.GoString(C.AB_Transaction_GetRemoteIban(t))
	transaction.RemoteBIC = C.GoString(C.AB_Transaction_GetRemoteBic(t))
	transaction.RemoteBankCode = C.GoString(C.AB_Transaction_GetRemoteBankCode(t))
	transaction.RemoteAccountNumber = C.GoString(C.AB_Transaction_GetRemoteAccountNumber(t))
	transaction.RemoteName = (*gwStringList)(C.AB_Transaction_GetRemoteName(t)).toString()

	return transaction, true
}

// Transactions implements AB_JobGetTransactions_new from aqbanking, listing
// all transactions from a given aqbanking instance
func (ab *AQBanking) Transactions(acc *Account, from *time.Time, to *time.Time) ([]Transaction, error) {
	abJob := C.AB_JobGetTransactions_new(acc.ptr)

	if abJob == nil {
		return nil, errors.New("Unable to load transactions.")
	}

	if err := C.AB_Job_CheckAvailability(abJob); err != 0 {
		return nil, fmt.Errorf("Transactions is not supported by backend: %d", err)
	}

	if from != nil && to != nil {
		C.AB_JobGetTransactions_SetFromTime(abJob, (*C.GWEN_TIME)(newGwenTime(*from)))
		C.AB_JobGetTransactions_SetToTime(abJob, (*C.GWEN_TIME)(newGwenTime(*to)))
	}

	abJobList := C.AB_Job_List2_new()
	C.AB_Job_List2_PushBack(abJobList, abJob)
	abContext := C.AB_ImExporterContext_new()

	if err := C.AB_Banking_ExecuteJobs(ab.ptr, abJobList, abContext); err != 0 {
		return nil, fmt.Errorf("Unable to execute Transactions: %d", err)
	}

	if C.AB_Job_GetStatus(abJob) == C.AB_Job_StatusError {
		return nil, errors.New(C.GoString(C.AB_Job_GetResultText(abJob)))
	}

	abInfo := C.AB_ImExporterContext_GetFirstAccountInfo(abContext)
	var transactions = make([]Transaction, 0)

	for abInfo != nil {
		abTransaction := C.AB_ImExporterAccountInfo_GetFirstTransaction(abInfo)

		for abTransaction != nil {
			transaction, ok := newTransaction(abTransaction)

			if ok {
				transactions = append(transactions, transaction)
			}

			abTransaction = C.AB_ImExporterAccountInfo_GetNextTransaction(abInfo)
		}
		abInfo = C.AB_ImExporterContext_GetNextAccountInfo(abContext)
	}

	C.AB_Job_free(abJob)

	return transactions, nil
}

// AllTransactions implements AB_JobGetTransactions_new without filter
func (ab *AQBanking) AllTransactions(acc *Account) ([]Transaction, error) {
	return ab.Transactions(acc, nil, nil)
}
