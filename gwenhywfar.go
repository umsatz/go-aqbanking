package aqbanking

/*
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar5
#cgo linux CFLAGS: -I/usr/include/gwenhywfar5
#include <gwenhywfar/gwendate.h>
*/
import "C"
import "time"

type gwDate C.GWEN_DATE

func (date *gwDate) toTime() time.Time {
	return time.Unix(int64(C.GWEN_Date_toLocalTime((*C.GWEN_DATE)(date))), 0)
}

func (date *gwDate) String() string {
	return C.GoString(C.GWEN_Date_GetString((*C.GWEN_DATE)(date)))
}

func newGwenDate(date time.Time) *gwDate {
	return (*gwDate)(C.GWEN_Date_fromGregorian(
		C.int(date.Year()),
		C.int(date.Month()),
		C.int(date.Day()),
	))
}

func gwenDateToTime(in *C.GWEN_DATE) time.Time {
	return (*gwDate)(in).toTime()
}
