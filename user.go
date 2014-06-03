package main

import "errors"

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking5
#include <aqbanking/banking.h>
*/
import "C"

type User struct {
	Id         int
	UserId     string
	CustomerId string
	Name       string
	Country    string

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

func NewUser(ptr *C.AB_USER) User {
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
		return nil, errors.New("Unable to load users.")
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
		collection.Users[i] = NewUser(abUser)
		abUser = C.AB_User_List2Iterator_Next(abIterator)
	}

	C.AB_User_List2Iterator_free(abIterator)
	C.AB_User_free(abUser)

	return collection, nil
}
