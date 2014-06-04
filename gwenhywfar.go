package main

/*
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#include <gwenhywfar/stringlist.h>
#include <gwenhywfar/gwentime.h>
*/
import "C"
import "time"

type gwStringList C.GWEN_STRINGLIST

func (list *gwStringList) toString() string {
	str := ""
	size := int(C.GWEN_StringList_Count((*C.GWEN_STRINGLIST)(list)))
	for i := 0; i < size; i++ {
		part := C.GoString(C.GWEN_StringList_StringAt((*C.GWEN_STRINGLIST)(list), C.int(i)))
		str += part
	}
	return str
}

type gwTime C.GWEN_TIME

func (gt *gwTime) toTime() time.Time {
	var seconds int64 = int64(C.GWEN_Time_Seconds((*C.GWEN_TIME)(gt)))
	return time.Unix(seconds, 0)
}
