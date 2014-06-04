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

type Gui C.struct_GWEN_GUI

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

	return (*Gui)(gui)
}

func NewNonInteractiveGui() *Gui {
	return newGui(false)
}

func (g *Gui) Attach(aq *AQBanking) {
	C.AB_Gui_Extend((*C.struct_GWEN_GUI)(g), aq.Ptr)
}

func (g *Gui) RegisterPins(pins []Pin) {
	var dbPins *C.GWEN_DB_NODE = C.GWEN_DB_Group_new(C.CString("pins"))

	for _, pin := range pins {
		str := fmt.Sprintf("PIN_%v_%v=%v\n", pin.Blz, pin.UserId, pin.Pin)
		pinLen := len(str)

		C.GWEN_DB_ReadFromString(dbPins, C.CString(str), C.int(pinLen), C.GWEN_PATH_FLAGS_CREATE_GROUP|C.GWEN_DB_FLAGS_DEFAULT)
		break
	}

	C.GWEN_Gui_CGui_SetPasswordDb((*C.struct_GWEN_GUI)(g), dbPins, 1)
}

func (g *Gui) Free() {
	C.GWEN_Gui_free((*C.struct_GWEN_GUI)(g))
}
