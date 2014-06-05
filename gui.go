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

type Gui struct {
	ptr    *C.struct_GWEN_GUI
	dbPins *C.GWEN_DB_NODE
}

type Pin interface {
	BankCode() string
	UserId() string
	Pin() string
}

func newGui(interactive bool) *Gui {
	var gui *C.struct_GWEN_GUI = C.GWEN_Gui_CGui_new()

	if !interactive {
		C.GWEN_Gui_SetFlags(gui, C.GWEN_GUI_FLAGS_ACCEPTVALIDCERTS|C.GWEN_GUI_FLAGS_NONINTERACTIVE)
	} else {
		C.GWEN_Gui_SetFlags(gui, C.GWEN_GUI_FLAGS_ACCEPTVALIDCERTS)
	}
	C.GWEN_Gui_SetCharSet(gui, C.CString("UTF-8"))
	C.GWEN_Gui_SetGui(gui)

	// C.GWEN_Logger_SetLevel(C.CString(C.AQBANKING_LOGDOMAIN), C.GWEN_LoggerLevel_Error)
	// C.GWEN_Logger_SetLevel(C.CString(C.GWEN_LOGDOMAIN), C.GWEN_LoggerLevel_Error)

	var dbPins *C.GWEN_DB_NODE = C.GWEN_DB_Group_new(C.CString("pins"))
	C.GWEN_Gui_CGui_SetPasswordDb(gui, dbPins, 1)

	return &Gui{
		gui,
		dbPins,
	}
}

func NewNonInteractiveGui() *Gui {
	return newGui(false)
}

func (g *Gui) Attach(aq *AQBanking) {
	C.AB_Gui_Extend(g.ptr, aq.ptr)
}

func (g *Gui) RegisterPin(pin Pin) {
	str := fmt.Sprintf("PIN_%v_%v=%v\n", pin.BankCode(), pin.UserId(), pin.Pin())
	pinLen := len(str)

	C.GWEN_DB_ReadFromString(g.dbPins, C.CString(str), C.int(pinLen), C.GWEN_PATH_FLAGS_CREATE_GROUP|C.GWEN_DB_FLAGS_DEFAULT)
}

func (g *Gui) Free() {
	C.GWEN_Gui_free(g.ptr)
}
