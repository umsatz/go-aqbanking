package aqbanking

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
#cgo linux CFLAGS: -I/usr/include/gwenhywfar4
#cgo linux CFLAGS: -I/usr/include/aqbanking5
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
		Purpose:           (*gwStringList)(C.AB_Transaction_GetPurpose(t)).toString(),
		Text:              C.GoString(C.AB_Transaction_GetTransactionText(t)),
		Status:            C.GoString(C.AB_Transaction_Status_toString(C.AB_Transaction_GetStatus(t))),
		CustomerReference: C.GoString(C.AB_Transaction_GetCustomerReference(t)),
		EndToEndReference: C.GoString(C.AB_Transaction_GetEndToEndReference(t)),
		MandateID:         C.GoString(C.AB_Transaction_GetMandateId(t)),

		Date:       (*gwTime)(C.AB_Transaction_GetDate(t)).toTime(),
		ValutaDate: (*gwTime)(C.AB_Transaction_GetValutaDate(t)).toTime(),

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
		RemoteName:          (*gwStringList)(C.AB_Transaction_GetRemoteName(t)).toString(),
	}

	if fees := C.AB_Transaction_GetFees(t); fees != nil {
		transaction.Fee = float32(C.AB_Value_GetValueAsDouble(fees))
		transaction.FeeCurrency = C.GoString(C.AB_Value_GetCurrency(fees))
	}

	return &transaction
}

// Transactions implements AB_JobGetTransactions_new from aqbanking, listing
// all transactions from a given aqbanking instance
func (ab *AQBanking) Transactions(acc *Account, from *time.Time, to *time.Time) ([]Transaction, error) {
	abJob := C.AB_JobGetTransactions_new(acc.ptr)
	defer C.AB_Job_free(abJob)

	if abJob == nil {
		return nil, errors.New("Unable to load transactions")
	}

	if err := C.AB_Job_CheckAvailability(abJob); err != 0 {
		return nil, fmt.Errorf("Transactions is not supported by backend: %d", err)
	}

	if from != nil {
		C.AB_JobGetTransactions_SetFromTime(abJob, (*C.GWEN_TIME)(newGwenTime(*from)))
	}
	if to != nil {
		C.AB_JobGetTransactions_SetToTime(abJob, (*C.GWEN_TIME)(newGwenTime(*to)))
	}

	abJobList := C.AB_Job_List2_new()
	defer C.AB_Job_List2_free(abJobList)
	C.AB_Job_List2_PushBack(abJobList, abJob)

	abContext := C.AB_ImExporterContext_new()
	defer C.AB_ImExporterContext_free(abContext)

	if err := C.AB_Banking_ExecuteJobs(ab.ptr, abJobList, abContext); err != 0 {
		return nil, fmt.Errorf("Unable to execute Transactions: %d", err)
	}

	status := C.AB_Job_GetStatus(abJob)
	if status == C.AB_Job_StatusError {
		return nil, errors.New(C.GoString(C.AB_Job_GetResultText(abJob)))
	}

	abInfo := C.AB_ImExporterContext_GetFirstAccountInfo(abContext)
	var transactions []Transaction

	if abInfo == nil {
		return nil, fmt.Errorf("Unable to get first account info")
	}

	for abInfo != nil {
		abTransaction := C.AB_ImExporterAccountInfo_GetFirstTransaction(abInfo)

		for abTransaction != nil {
			transaction := newTransaction(abTransaction)
			if transaction != nil {
				transactions = append(transactions, *transaction)
			}

			abTransaction = C.AB_ImExporterAccountInfo_GetNextTransaction(abInfo)
		}
		abInfo = C.AB_ImExporterContext_GetNextAccountInfo(abContext)
	}

	return transactions, nil
}

// AllTransactions implements AB_JobGetTransactions_new without filter
func (ab *AQBanking) AllTransactions(acc *Account) ([]Transaction, error) {
	return ab.Transactions(acc, nil, nil)
}
