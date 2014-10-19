package main

import (
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

// AQBankingVersion wraps AQBanking version informations
type AQBankingVersion struct {
	Major      int
	Minor      int
	Patchlevel int
}

// AQBanking represents a single aqbanking database path, located at a given path
type AQBanking struct {
	Name    string
	Version AQBankingVersion
	gui     *gui
	ptr     *C.AB_BANKING
}

// NewAQBanking creates a new AQBanking instance, given valid database path and name
func NewAQBanking(name string, dbPath string) (*AQBanking, error) {
	inst := &AQBanking{}
	inst.Name = name

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	if dbPath == "" {
		inst.ptr = C.AB_Banking_new(cName, nil, 0)
	} else {
		cPath := C.CString(dbPath)
		defer C.free(unsafe.Pointer(cPath))

		inst.ptr = C.AB_Banking_new(cName, cPath, 0)
	}

	if err := C.AB_Banking_Init(inst.ptr); err != 0 {
		return nil, fmt.Errorf("unable to initialized aqbanking: %d", err)
	}
	if err := C.AB_Banking_OnlineInit(inst.ptr); err != 0 {
		return nil, fmt.Errorf("unable to initialized aqbanking: %d", err)
	}

	inst.loadVersion()

	inst.gui = newNonInteractiveGui()
	inst.gui.attach(inst)

	return inst, nil
}

// DefaultAQBanking returns an aqbanking instance initialized with aqbankings default
// database path (most likely $HOME)
func DefaultAQBanking() (*AQBanking, error) {
	return NewAQBanking("local", "")
}

func (ab *AQBanking) loadVersion() {
	var major, minor, patchlevel, build C.int
	C.AB_Banking_GetVersion(&major, &minor, &patchlevel, &build)
	ab.Version = AQBankingVersion{int(major), int(minor), int(patchlevel)}
}

// Free frees all underlying aqbanking pointers
func (ab *AQBanking) Free() error {
	if err := C.AB_Banking_OnlineFini(ab.ptr); err != 0 {
		return fmt.Errorf("unable to free aqbanking online: %d\n", err)
	}

	ab.gui.free()

	C.AB_Banking_Fini(ab.ptr)
	C.AB_Banking_free(ab.ptr)

	return nil
}
