// Package aqbanking is a go wrapper around aqbanking version 5.0.24 to 5.5.1
//
// Details about aqbanking are available at the official homepage: http://www.aqbanking.de/
package aqbanking

import (
	"fmt"
	"unsafe"
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

// Version wraps AQBanking version informations
type Version struct {
	Major      int
	Minor      int
	Patchlevel int
}

// AQBanking represents a single aqbanking database path, located at a given path
type AQBanking struct {
	Name    string
	Version Version
	gui     *gui
	ptr     *C.AB_BANKING
}

// NewAQBanking creates a new AQBanking instance, given valid database path and name
func NewAQBanking(name string, dbPath string) (*AQBanking, error) {
	inst := &AQBanking{
		Name:    name,
		Version: getVersion(),
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	var cPath *C.char
	if dbPath != "" {
		cPath := C.CString(dbPath)
		defer C.free(unsafe.Pointer(cPath))
	}

	inst.ptr = C.AB_Banking_new(cName, cPath, 0)
	inst.gui = newNonInteractiveGui()
	inst.gui.attach(inst)

	if err := C.AB_Banking_Init(inst.ptr); err != 0 {
		return nil, newError("unable to initialize aqbanking", err)
	}

	// version string is not allowed to be longer than 5 characters
	inst.SetRuntimeConfig("fintsApplicationVersionString", fmt.Sprintf("%d.%d", inst.Version.Major, inst.Version.Minor))
	inst.SetRuntimeConfig("fintsRegistrationKey", "32F8A67FE34B57AB8D7E4FE70")

	return inst, nil
}

// DefaultAQBanking returns an aqbanking instance initialized with aqbankings default
// database path (most likely $HOME)
func DefaultAQBanking() (*AQBanking, error) {
	return NewAQBanking("local", "")
}

// SetRuntimeConfig sets a runtime configuration value
func (ab *AQBanking) SetRuntimeConfig(key, value string) {
	cKey := C.CString(key)
	cValue := C.CString(value)

	C.AB_Banking_RuntimeConfig_SetCharValue(ab.ptr, cKey, cValue)

	C.free(unsafe.Pointer(cKey))
	C.free(unsafe.Pointer(cValue))
}

// Free frees all underlying aqbanking pointers
func (ab *AQBanking) Free() error {
	ab.gui.free()

	C.AB_Banking_Fini(ab.ptr)
	C.AB_Banking_free(ab.ptr)

	return nil
}

func getVersion() Version {
	var major, minor, patchlevel, build C.int
	C.AB_Banking_GetVersion(&major, &minor, &patchlevel, &build)
	return Version{
		int(major),
		int(minor),
		int(patchlevel),
	}
}

func (version Version) String() string {
	return fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patchlevel)
}
