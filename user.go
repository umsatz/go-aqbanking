package main

import "errors"

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo CFLAGS: -I/usr/local/include/aqbanking5
#include <aqbanking5/aqbanking/banking.h>
*/
import "C"

type User struct {
	UserId     string
	CustomerId string
	Name       string
}

func (ab *AQBanking) Users() ([]User, error) {
	var abUserList *C.AB_USER_LIST2 = C.AB_Banking_GetUsers(ab.Ptr)
	if abUserList == nil {
		return nil, errors.New("Unable to load users.")
	}

	var users []User = make([]User, C.AB_Account_List2_GetSize(abUserList))

	var abIterator *C.AB_USER_LIST2_ITERATOR = C.AB_User_List2_First(abUserList)
	if abIterator == nil {
		return nil, errors.New("Unable to get user iterator.")
	}

	var abUser *C.AB_USER
	abUser = C.AB_User_List2Iterator_Data(abIterator)

	for i := 0; abUser != nil; i++ {
		user := User{}

		user.UserId = C.GoString(C.AB_User_GetUserId(abUser))
		user.CustomerId = C.GoString(C.AB_User_GetCustomerId(abUser))
		user.Name = C.GoString(C.AB_User_GetUserName(abUser))

		users[i] = user
		abUser = C.AB_User_List2Iterator_Next(abIterator)
	}

	C.AB_User_List2Iterator_free(abIterator)
	C.AB_User_free(abUser)
	C.AB_User_List2_freeAll(abUserList)

	return users, nil
}
