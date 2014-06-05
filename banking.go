package main

import (
	"errors"
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking5
#include <aqbanking/banking.h>
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

	ptr *C.AB_BANKING
}

func NewAQBanking(name string, dbPath string) (*AQBanking, error) {
	inst := &AQBanking{}
	inst.Name = name

	var cName *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	if dbPath == "" {
		inst.ptr = C.AB_Banking_new(cName, nil, 0)
	} else {
		var cPath *C.char = C.CString(dbPath)
		defer C.free(unsafe.Pointer(cPath))

		inst.ptr = C.AB_Banking_new(cName, cPath, 0)
	}

	if err := C.AB_Banking_Init(inst.ptr); err != 0 {
		return nil, errors.New(fmt.Sprintf("unable to initialized aqbanking: %d", err))
	}
	if err := C.AB_Banking_OnlineInit(inst.ptr); err != 0 {
		return nil, errors.New(fmt.Sprintf("unable to initialized aqbanking: %d", err))
	}

	inst.loadVersion()

	return inst, nil
}

func DefaultAQBanking() (*AQBanking, error) {
	return NewAQBanking("local", "")
}

func (ab *AQBanking) loadVersion() {
	var major, minor, patchlevel, build C.int
	C.AB_Banking_GetVersion(&major, &minor, &patchlevel, &build)
	ab.Version = AQBankingVersion{int(major), int(minor), int(patchlevel)}
}

func (ab *AQBanking) Free() error {
	if err := C.AB_Banking_OnlineFini(ab.ptr); err != 0 {
		return errors.New(fmt.Sprintf("unable to free aqbanking online: %d\n", err))
	}

	C.AB_Banking_Fini(ab.ptr)
	C.AB_Banking_free(ab.ptr)
	return nil
}
