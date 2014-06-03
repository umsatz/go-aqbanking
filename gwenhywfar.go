package main

/*
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#include <gwenhywfar/stringlist.h>
*/
import "C"

type GwStringList C.GWEN_STRINGLIST

func (list *GwStringList) toString() string {
	str := ""
	size := int(C.GWEN_StringList_Count((*C.GWEN_STRINGLIST)(list)))
	for i := 0; i < size; i++ {
		part := C.GoString(C.GWEN_StringList_StringAt((*C.GWEN_STRINGLIST)(list), C.int(i)))
		str += part
	}
	return str
}
