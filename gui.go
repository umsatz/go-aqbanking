package main

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking5
#include <gwenhywfar/cgui.h>
*/
import "C"
import "fmt"

type Gui C.struct_GWEN_GUI

func newGui(interactive bool) *Gui {
	var gui *C.struct_GWEN_GUI = C.GWEN_Gui_CGui_new()
	if !interactive {
		C.GWEN_Gui_SetFlags(gui, C.GWEN_GUI_FLAGS_ACCEPTVALIDCERTS|C.GWEN_GUI_FLAGS_NONINTERACTIVE)
	}
	C.GWEN_Gui_SetGui(gui)
	return (*Gui)(gui)
}

func NewNonInteractiveGui() *Gui {
	return newGui(false)
}

func (g *Gui) RegisterPins(aq *AQBanking, pins []Pin) {
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
