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
#include <aqbanking/jobgettransactions.h>
#include <aqbanking/banking.h>
#include <aqbanking/banking_ob.h>
*/
import "C"

type Transaction struct {
	Purpose           string
	Category          string
	Type              string // AB_Transaction_Type
	SubType           string // AB_TRANSACTION_SUBTYPE
	Text              string
	MandateReference  string
	CustomerReference string
	Currency          string
	Total             float32
	LocalBankCode     string
	RemoteBankCode    string
	TransactionPeriod string
}

// implements AB_JobGetTransactions_new
func (ab *AQBanking) Transactions(acc Account) ([]Transaction, error) {
	fmt.Println("before get transactions")
	var abJob *C.AB_JOB = C.AB_JobGetTransactions_new(acc.Ptr)
	fmt.Println("after get transactions")
	if abJob == nil {
		return nil, errors.New("Unable to load transactions.")
	}

	if err := C.AB_Job_CheckAvailability(abJob); err != 0 {
		return nil, errors.New(fmt.Sprintf("Transactions is not supported by backend: %d", err))
	}

	// TODO set arguments?
	// AB_JobGetTransactions_SetFromTime
	// AB_JobGetTransactions_SetToTime

	var abJobList *C.AB_JOB_LIST2 = C.AB_Job_List2_new()
	C.AB_Job_List2_PushBack(abJobList, abJob)
	var abContext *C.AB_IMEXPORTER_CONTEXT = C.AB_ImExporterContext_new()

	if err := C.AB_Banking_ExecuteJobs(ab.Ptr, abJobList, abContext); err != 0 {
		return nil, errors.New(fmt.Sprintf("Unable to execute Transactions: %d", err))
	}

	var abInfo *C.AB_IMEXPORTER_ACCOUNTINFO = C.AB_ImExporterContext_GetFirstAccountInfo(abContext)
	var transactions []Transaction = make([]Transaction, 0)

	for abInfo != nil {
		var abTransaction *C.AB_TRANSACTION = C.AB_ImExporterAccountInfo_GetFirstTransaction(abInfo)

		for abTransaction != nil {
			var abValue *C.AB_VALUE
			abValue = C.AB_Transaction_GetValue(abTransaction)
			if abValue != nil {
				var transaction Transaction = Transaction{}

				transaction.Purpose = (*GwStringList)(C.AB_Transaction_GetPurpose(abTransaction)).toString()
				transaction.Category = (*GwStringList)(C.AB_Transaction_GetCategory(abTransaction)).toString()
				transaction.Text = C.GoString(C.AB_Transaction_GetTransactionText(abTransaction))
				transaction.MandateReference = C.GoString(C.AB_Transaction_GetMandateReference(abTransaction))
				transaction.CustomerReference = C.GoString(C.AB_Transaction_GetMandateReference(abTransaction))

				transaction.Currency = C.GoString(C.AB_Value_GetCurrency(abValue))
				transaction.Total = float32(C.AB_Value_GetValueAsDouble(abValue))

				transaction.Type = C.GoString(C.AB_Transaction_Type_toString(C.AB_Transaction_GetType(abTransaction)))
				transaction.SubType = C.GoString(C.AB_Transaction_SubType_toString(C.AB_Transaction_GetSubType(abTransaction)))

				transaction.TransactionPeriod = C.GoString(C.AB_Transaction_Period_toString(C.AB_Transaction_GetPeriod(abTransaction)))

				transactions = append(transactions, transaction)
			}

			abTransaction = C.AB_ImExporterAccountInfo_GetNextTransaction(abInfo)
		}
		abInfo = C.AB_ImExporterContext_GetNextAccountInfo(abContext)
	}

	C.AB_Job_free(abJob)

	return transactions, nil
}
