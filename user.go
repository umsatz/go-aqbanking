package aqbanking

import (
	"errors"
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -laqhbci
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking5
#cgo linux CFLAGS: -I/usr/include/gwenhywfar4
#cgo linux CFLAGS: -I/usr/include/aqbanking5
#include <stdlib.h>
#include <aqbanking/banking.h>
#include <aqbanking/banking_be.h>
#include <aqhbci/provider.h>
#include <aqhbci/user.h>
#include <gwenhywfar/text.h>
#include <gwenhywfar/url.h>
#include <gwenhywfar/ct.h>
#include <gwenhywfar/ctplugin.h>
*/
import "C"

// User represents an aqbanking user account
type User struct {
	ID          int
	UserID      string // Benutzerkennung
	CustomerID  string // Kundennummer
	BankCode    string // BLZ
	Name        string
	ServerURI   string
	HbciVersion int

	ptr *C.AB_USER
}

// Free frees the underlying aqbanking pointer for this user
func (u *User) Free() {
	C.AB_User_free(u.ptr)
}

// UserCollection wraps a list of aqbanking users to handle the collection pointer
type UserCollection struct {
	Users []User
	ptr   *C.AB_USER_LIST2
}

// Free frees all user accounts as well as the underlying collection pointer
func (ul *UserCollection) Free() {
	for i := range ul.Users {
		ul.Users[i].Free()
	}
	ul.Users = make([]User, 0)
	C.AB_User_List2_free(ul.ptr)
}

var supportedHBCIVersions = map[int]struct{}{
	201: struct{}{},
	210: struct{}{},
	220: struct{}{},
	300: struct{}{},
}

// AddPinTanUser implements the simplified, pintan only workflow from
// src/plugins/backends/aqhbci/tools/aqhbci-tool/adduser.c
func (ab *AQBanking) AddPinTanUser(user *User) error {
	if user.BankCode == "" {
		return errors.New("no bankCode given")
	}
	if user.UserID == "" {
		return errors.New("no userid given")
	}
	if user.ServerURI == "" {
		return errors.New("no server_url given")
	}

	if _, ok := supportedHBCIVersions[user.HbciVersion]; ok != true {
		return fmt.Errorf("hbci version %d is not supported", user.HbciVersion)
	}

	var aqUser *C.AB_USER

	aqhbciProviderName := C.CString("aqhbci")
	defer C.free(unsafe.Pointer(aqhbciProviderName))

	countryDe := C.CString("de")
	defer C.free(unsafe.Pointer(countryDe))

	aqPinTan := C.CString("pintan")
	defer C.free(unsafe.Pointer(aqPinTan))

	var _ *C.AB_PROVIDER = C.AB_Banking_GetProvider(ab.ptr, aqhbciProviderName)

	aqBankCode := C.CString(user.BankCode)
	defer C.free(unsafe.Pointer(aqBankCode))

	aqUserID := C.CString(user.UserID)
	defer C.free(unsafe.Pointer(aqUserID))

	aqName := C.CString(user.Name)
	defer C.free(unsafe.Pointer(aqName))

	aqUser = C.AB_Banking_FindUser(
		ab.ptr,
		C.CString(C.AH_PROVIDER_NAME),
		countryDe,
		aqBankCode,
		aqUserID,
		aqUserID,
	)
	if aqUser != nil {
		return fmt.Errorf("user %s already exists", user.UserID)
	}

	aqUser = C.AB_Banking_CreateUser(ab.ptr, C.CString(C.AH_PROVIDER_NAME))
	if aqUser == nil {
		return errors.New("unable to create user")
	}

	url := C.GWEN_Url_fromString(C.CString(user.ServerURI))
	if url == nil {
		return errors.New("invalid server url")
	}
	C.GWEN_Url_SetProtocol(url, C.CString("https"))
	if C.GWEN_Url_GetPort(url) == 0 {
		C.GWEN_Url_SetPort(url, C.int(443))
	}
	defer C.GWEN_Url_free(url)

	C.AB_User_SetUserName(aqUser, aqName)
	C.AB_User_SetCountry(aqUser, countryDe)
	C.AB_User_SetBankCode(aqUser, aqBankCode)
	C.AB_User_SetUserId(aqUser, aqUserID)
	C.AB_User_SetCustomerId(aqUser, aqUserID)

	C.AH_User_SetTokenType(aqUser, aqPinTan)
	C.AH_User_SetTokenContextId(aqUser, C.uint32_t(1)) // context
	C.AH_User_SetCryptMode(aqUser, C.AH_CryptMode_Pintan)
	C.AH_User_SetHbciVersion(aqUser, C.int(user.HbciVersion))
	C.AH_User_SetServerUrl(aqUser, url)

	C.AB_Banking_AddUser(ab.ptr, aqUser)
	user.ptr = aqUser

	return nil
}

// Remove removes a user from the given aqbanking database
func (u *User) Remove(aq *AQBanking) error {
	accountCollection, err := aq.AccountsFor(u)
	if err != nil {
		return err
	}

	for _, account := range accountCollection.Accounts {
		if err := account.Remove(aq); err != nil {
			return err
		}
	}

	if err := C.AB_Banking_DeleteUser(aq.ptr, u.ptr); err != 0 {
		return fmt.Errorf("unable to delete user: %d", err)
	}
	return nil
}

// FetchAccounts returns all accounts registered for a given aqbanking instance
func (u *User) FetchAccounts(aq *AQBanking) error {
	ctx := C.AB_ImExporterContext_new()

	pro := C.AB_Banking_GetProvider(aq.ptr, C.CString("aqhbci"))
	if err := C.AH_Provider_GetAccounts(pro, u.ptr, ctx, 1, 0, 1); err != 0 {
		return fmt.Errorf("Error getting accounts (%d)", err)
	}

	C.AB_ImExporterContext_free(ctx)
	return nil
}

func newUser(ptr *C.AB_USER) User {
	user := User{}
	user.ID = int(C.AB_User_GetUniqueId(ptr))

	user.UserID = C.GoString(C.AB_User_GetUserId(ptr))
	user.CustomerID = C.GoString(C.AB_User_GetCustomerId(ptr))
	user.Name = C.GoString(C.AB_User_GetUserName(ptr))
	user.BankCode = C.GoString(C.AB_User_GetBankCode(ptr))

	url := C.AH_User_GetServerUrl(ptr)
	if url != nil {
		tbuf := C.GWEN_Buffer_new(
			nil,
			C.uint32_t(256),
			C.uint32_t(0),
			C.int(1),
		)
		C.GWEN_Url_toString(url, tbuf)
		user.ServerURI = C.GoString(C.GWEN_Buffer_GetStart(tbuf))
		C.GWEN_Buffer_free(tbuf)
	}

	user.HbciVersion = int(C.AH_User_GetHbciVersion(ptr))

	user.ptr = ptr
	return user
}

// Users implements AB_Banking_GetUsers, returning all users registered
// with aqbanking
func (ab *AQBanking) Users() (*UserCollection, error) {
	abUserList := C.AB_Banking_GetUsers(ab.ptr)
	if abUserList == nil {
		// no users available
		return &UserCollection{}, nil
	}

	collection := &UserCollection{}
	collection.Users = make([]User, C.AB_User_List2_GetSize(abUserList))
	collection.ptr = abUserList

	abIterator := C.AB_User_List2_First(abUserList)
	if abIterator == nil {
		return nil, errors.New("Unable to get user iterator")
	}

	var abUser *C.AB_USER
	abUser = C.AB_User_List2Iterator_Data(abIterator)

	for i := 0; abUser != nil; i++ {
		collection.Users[i] = newUser(abUser)
		abUser = C.AB_User_List2Iterator_Next(abIterator)
	}

	C.AB_User_List2Iterator_free(abIterator)
	C.AB_User_free(abUser)

	return collection, nil
}
