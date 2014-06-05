package main

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking5
#include <gwenhywfar/cgui.h>
#include <aqbanking/abgui.h>
*/
import "C"
import "fmt"

type gui struct {
	ptr    *C.struct_GWEN_GUI
	dbPins *C.GWEN_DB_NODE
}

func newGui(interactive bool) *gui {
	var abGui *C.struct_GWEN_GUI = C.GWEN_Gui_CGui_new()

	if !interactive {
		C.GWEN_Gui_SetFlags(abGui, C.GWEN_GUI_FLAGS_ACCEPTVALIDCERTS|C.GWEN_GUI_FLAGS_NONINTERACTIVE)
	} else {
		C.GWEN_Gui_SetFlags(abGui, C.GWEN_GUI_FLAGS_ACCEPTVALIDCERTS)
	}
	C.GWEN_Gui_SetCharSet(abGui, C.CString("UTF-8"))
	C.GWEN_Gui_SetGui(abGui)

	// C.GWEN_Logger_SetLevel(C.CString(C.AQBANKING_LOGDOMAIN), C.GWEN_LoggerLevel_Error)
	// C.GWEN_Logger_SetLevel(C.CString(C.GWEN_LOGDOMAIN), C.GWEN_LoggerLevel_Error)

	var dbPins *C.GWEN_DB_NODE = C.GWEN_DB_Group_new(C.CString("pins"))
	C.GWEN_Gui_CGui_SetPasswordDb(abGui, dbPins, 1)

	return &gui{
		abGui,
		dbPins,
	}
}

func newNonInteractiveGui() *gui {
	return newGui(false)
}

func (g *gui) attach(aq *AQBanking) {
	C.AB_Gui_Extend(g.ptr, aq.ptr)
}

func (ab *AQBanking) RegisterPin(pin Pin) {
	str := fmt.Sprintf("PIN_%v_%v=%v\n", pin.BankCode(), pin.UserId(), pin.Pin())
	pinLen := len(str)

	C.GWEN_DB_ReadFromString(ab.gui.dbPins, C.CString(str), C.int(pinLen), C.GWEN_PATH_FLAGS_CREATE_GROUP|C.GWEN_DB_FLAGS_DEFAULT)
}

func (g *gui) free() {
	C.GWEN_Gui_free(g.ptr)
}
