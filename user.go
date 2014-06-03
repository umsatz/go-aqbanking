package main

import (
	"errors"
	"fmt"
)

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -laqhbci
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking5
#include <aqbanking/banking.h>
#include <aqbanking/banking_be.h>
#include <aqhbci/user.h>
#include <gwenhywfar/text.h>
#include <gwenhywfar/url.h>
#include <gwenhywfar/ct.h>
#include <gwenhywfar/ctplugin.h>
*/
import "C"

type User struct {
	Id         int
	UserId     string
	CustomerId string
	Name       string
	Country    string
	ServerUrl  string

	Ptr *C.AB_USER
}

type UserCollection struct {
	Users []User
	Ptr   *C.AB_USER_LIST2
}

func (ul *UserCollection) Free() {
	ul.Users = make([]User, 0)
	C.AB_User_List2_free(ul.Ptr)
}

func (ab *AQBanking) AddPinUser(userId, bankCode, name, serverUrl string) (User, error) {
	var aqUser *C.AB_USER
	var _ *C.AB_PROVIDER = C.AB_Banking_GetProvider(ab.Ptr, C.CString("aqhbci"))
	hbciVersion := 300

	if bankCode == "" {
		return User{}, errors.New("no bankCode given.")
	}
	if userId == "" {
		return User{}, errors.New("no userId given")
	}

	var supportHBCIVersions map[int]struct{} = map[int]struct{}{
		201: struct{}{},
		210: struct{}{},
		220: struct{}{},
		300: struct{}{},
	}
	if _, ok := supportHBCIVersions[hbciVersion]; ok != true {
		return User{}, errors.New(fmt.Sprintf("hbci version %d is not supported.", hbciVersion))
	}

	// TODO check for dups
	// user = AB_Banking_FindUser(ab, AH_PROVIDER_NAME,
	// 	"de",
	// 	lbankId, luserId, lcustomerId)
	// if user {
	// 	DBG_ERROR(0, "User %s already exists", luserId)
	// 	return 3
	// }

	aqUser = C.AB_Banking_CreateUser(ab.Ptr, C.CString(C.AH_PROVIDER_NAME))
	if aqUser == nil {
		return User{}, errors.New("unable to create user.")
	}

	var url *C.GWEN_URL = C.GWEN_Url_fromString(C.CString(serverUrl))
	if url == nil {
		return User{}, errors.New("invalid server url.")
	}
	C.GWEN_Url_SetProtocol(url, C.CString("https"))
	if C.GWEN_Url_GetPort(url) == 0 {
		C.GWEN_Url_SetPort(url, C.int(443))
	}
	defer C.GWEN_Url_free(url)

	C.AB_User_SetUserName(aqUser, C.CString(name))
	C.AB_User_SetCountry(aqUser, C.CString("de"))
	C.AB_User_SetBankCode(aqUser, C.CString(bankCode))
	C.AB_User_SetUserId(aqUser, C.CString(userId))
	C.AB_User_SetCustomerId(aqUser, C.CString(userId))

	// C.AH_User_SetTokenType(aqUser, C.CString("pintan"))
	C.AH_User_SetTokenContextId(aqUser, C.uint32_t(1))
	C.AH_User_SetCryptMode(aqUser, C.AH_CryptMode_Pintan)
	C.AH_User_SetHbciVersion(aqUser, C.int(hbciVersion))
	C.AH_User_SetServerUrl(aqUser, url)

	C.AB_Banking_AddUser(ab.Ptr, aqUser)

	return User{}, nil
}

func newUser(ptr *C.AB_USER) User {
	user := User{}
	user.Id = int(C.AB_User_GetUniqueId(ptr))

	user.UserId = C.GoString(C.AB_User_GetUserId(ptr))
	user.CustomerId = C.GoString(C.AB_User_GetCustomerId(ptr))
	user.Name = C.GoString(C.AB_User_GetUserName(ptr))
	user.Country = C.GoString(C.AB_User_GetCountry(ptr))

	user.Ptr = ptr
	return user
}

// implements AB_Banking_GetUsers
func (ab *AQBanking) Users() (*UserCollection, error) {
	var abUserList *C.AB_USER_LIST2 = C.AB_Banking_GetUsers(ab.Ptr)
	if abUserList == nil {
		// no users available
		return &UserCollection{}, nil
	}

	collection := &UserCollection{}
	collection.Users = make([]User, C.AB_Account_List2_GetSize(abUserList))

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
