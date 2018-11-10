package aqbanking

/*
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo linux CFLAGS: -I/usr/include/gwenhywfar4
#include <gwenhywfar/stringlist.h>
#include <gwenhywfar/gwentime.h>
*/
import "C"
import "time"

type gwStringList C.GWEN_STRINGLIST

func (list *gwStringList) toSlice() []string {
	var s []string
	size := int(C.GWEN_StringList_Count((*C.GWEN_STRINGLIST)(list)))
	for i := 0; i < size; i++ {
		l := C.GoString(C.GWEN_StringList_StringAt((*C.GWEN_STRINGLIST)(list), C.int(i)))
		s = append(s, l)
	}
	return s
}

type gwTime C.GWEN_TIME

func (gt *gwTime) toTime() time.Time {
	seconds := int64(C.GWEN_Time_Seconds((*C.GWEN_TIME)(gt)))
	return time.Unix(seconds, 0)
}

func newGwenTime(date time.Time) *gwTime {
	utcDate := date.UTC()
	return (*gwTime)(C.GWEN_Time_fromSeconds(C.uint32_t(utcDate.Unix())))
}
