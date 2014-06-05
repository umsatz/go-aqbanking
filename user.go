package main

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

type User struct {
	Id          int
	UserId      string // Benutzerkennung
	CustomerId  string // Kundennummer
	BankCode    string // BLZ
	Name        string
	ServerUri   string
	HbciVersion int

	ptr *C.AB_USER
}

func (u *User) Free() {
	C.AB_User_free(u.ptr)
}

type UserCollection struct {
	Users []User
	ptr   *C.AB_USER_LIST2
}

func (ul *UserCollection) Free() {
	for i, _ := range ul.Users {
		ul.Users[i].Free()
	}
	ul.Users = make([]User, 0)
	C.AB_User_List2_free(ul.ptr)
}

var supportHBCIVersions map[int]struct{} = map[int]struct{}{
	201: struct{}{},
	210: struct{}{},
	220: struct{}{},
	300: struct{}{},
}

// implements the simplified, pintan only workflow from
// src/plugins/backends/aqhbci/tools/aqhbci-tool/adduser.c
func (ab *AQBanking) AddPinTanUser(user *User) error {
	if user.BankCode == "" {
		return errors.New("no bankCode given.")
	}
	if user.UserId == "" {
		return errors.New("no userid given")
	}
	if user.ServerUri == "" {
		return errors.New("no server_url given")
	}

	if _, ok := supportHBCIVersions[user.HbciVersion]; ok != true {
		return errors.New(fmt.Sprintf("hbci version %d is not supported.", user.HbciVersion))
	}

	var aqUser *C.AB_USER

	var aqhbciProviderName *C.char = C.CString("aqhbci")
	defer C.free(unsafe.Pointer(aqhbciProviderName))

	var countryDe *C.char = C.CString("de")
	defer C.free(unsafe.Pointer(countryDe))

	var aqPinTan *C.char = C.CString("pintan")
	defer C.free(unsafe.Pointer(aqPinTan))

	var _ *C.AB_PROVIDER = C.AB_Banking_GetProvider(ab.ptr, aqhbciProviderName)

	var aqBankCode *C.char = C.CString(user.BankCode)
	defer C.free(unsafe.Pointer(aqBankCode))

	var aqUserId *C.char = C.CString(user.UserId)
	defer C.free(unsafe.Pointer(aqUserId))

	var aqName *C.char = C.CString(user.Name)
	defer C.free(unsafe.Pointer(aqName))

	aqUser = C.AB_Banking_FindUser(
		ab.ptr,
		C.CString(C.AH_PROVIDER_NAME),
		countryDe,
		aqBankCode,
		aqUserId,
		aqUserId,
	)
	if aqUser != nil {
		return errors.New(fmt.Sprintf("user %s already exists.", user.UserId))
	}

	aqUser = C.AB_Banking_CreateUser(ab.ptr, C.CString(C.AH_PROVIDER_NAME))
	if aqUser == nil {
		return errors.New("unable to create user.")
	}

	var url *C.GWEN_URL = C.GWEN_Url_fromString(C.CString(user.ServerUri))
	if url == nil {
		return errors.New("invalid server url.")
	}
	C.GWEN_Url_SetProtocol(url, C.CString("https"))
	if C.GWEN_Url_GetPort(url) == 0 {
		C.GWEN_Url_SetPort(url, C.int(443))
	}
	defer C.GWEN_Url_free(url)

	C.AB_User_SetUserName(aqUser, aqName)
	C.AB_User_SetCountry(aqUser, countryDe)
	C.AB_User_SetBankCode(aqUser, aqBankCode)
	C.AB_User_SetUserId(aqUser, aqUserId)
	C.AB_User_SetCustomerId(aqUser, aqUserId)

	C.AH_User_SetTokenType(aqUser, aqPinTan)
	C.AH_User_SetTokenContextId(aqUser, C.uint32_t(1)) // context
	C.AH_User_SetCryptMode(aqUser, C.AH_CryptMode_Pintan)
	C.AH_User_SetHbciVersion(aqUser, C.int(user.HbciVersion))
	C.AH_User_SetServerUrl(aqUser, url)

	C.AB_Banking_AddUser(ab.ptr, aqUser)
	user.ptr = aqUser

	return nil
}

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
		return errors.New(fmt.Sprintf("unable to delete user: %d\n", err))
	}
	return nil
}

func (u *User) FetchAccounts(aq *AQBanking) error {
	var ctx *C.AB_IMEXPORTER_CONTEXT = C.AB_ImExporterContext_new()

	var pro *C.AB_PROVIDER = C.AB_Banking_GetProvider(aq.ptr, C.CString("aqhbci"))
	if err := C.AH_Provider_GetAccounts(pro, u.ptr, ctx, 1, 0, 1); err != 0 {
		return errors.New(fmt.Sprintf("Error getting accounts (%d)", err))
	}

	C.AB_ImExporterContext_free(ctx)
	return nil
}

func newUser(ptr *C.AB_USER) User {
	user := User{}
	user.Id = int(C.AB_User_GetUniqueId(ptr))

	user.UserId = C.GoString(C.AB_User_GetUserId(ptr))
	user.CustomerId = C.GoString(C.AB_User_GetCustomerId(ptr))
	user.Name = C.GoString(C.AB_User_GetUserName(ptr))
	user.BankCode = C.GoString(C.AB_User_GetBankCode(ptr))

	var url *C.GWEN_URL = C.AH_User_GetServerUrl(ptr)
	if url != nil {
		var tbuf *C.GWEN_BUFFER = C.GWEN_Buffer_new(
			nil,
			C.uint32_t(256),
			C.uint32_t(0),
			C.int(1),
		)
		C.GWEN_Url_toString(url, tbuf)
		user.ServerUri = C.GoString(C.GWEN_Buffer_GetStart(tbuf))
		C.GWEN_Buffer_free(tbuf)
	}

	user.HbciVersion = int(C.AH_User_GetHbciVersion(ptr))

	user.ptr = ptr
	return user
}

// implements AB_Banking_GetUsers
func (ab *AQBanking) Users() (*UserCollection, error) {
	var abUserList *C.AB_USER_LIST2 = C.AB_Banking_GetUsers(ab.ptr)
	if abUserList == nil {
		// no users available
		return &UserCollection{}, nil
	}

	collection := &UserCollection{}
	collection.Users = make([]User, C.AB_Account_List2_GetSize(abUserList))
	collection.ptr = abUserList

	var abIterator *C.AB_USER_LIST2_ITERATOR = C.AB_User_List2_First(abUserList)
	if abIterator == nil {
		return nil, errors.New("Unable to get user iterator.")
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
