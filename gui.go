package aqbanking

/*
#cgo LDFLAGS: -laqbanking
#cgo LDFLAGS: -lgwenhywfar
#cgo darwin CFLAGS: -I/usr/local/include/gwenhywfar4
#cgo darwin CFLAGS: -I/usr/local/include/aqbanking5
#include <gwenhywfar/cgui.h>
#include <aqbanking/abgui.h>
#include <gwenhywfar/gwenhywfar.h>

// forward declaration to allow cgui.go to use our go callback fnc
int goAqbankingGetPasswordFn_cgo(
		GWEN_GUI *gui,
		uint32_t flags,
		const char *token,
		const char *title,
		const char *text,
		char *buffer,
		int minLen,
		int maxLen,
		uint32_t guiid
	);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type gui struct {
	ptr *C.struct_GWEN_GUI
}

//export aqbankingGetPasswordFn
func aqbankingGetPasswordFn(token *C.char, buffer unsafe.Pointer, minLen, maxLen int) int {
	pin, ok := knownAqbankingPins[C.GoString(token)]
	if ok {
		C.memcpy(buffer, unsafe.Pointer(C.CString(pin)), 6)
	}
	return 0
}

var knownAqbankingPins = map[string]string{}

func newGui(interactive bool) *gui {
	abGui := C.GWEN_Gui_CGui_new()

	if !interactive {
		C.GWEN_Gui_SetFlags(abGui, C.GWEN_GUI_FLAGS_ACCEPTVALIDCERTS|C.GWEN_GUI_FLAGS_NONINTERACTIVE)
	} else {
		C.GWEN_Gui_SetFlags(abGui, C.GWEN_GUI_FLAGS_ACCEPTVALIDCERTS)
	}
	C.GWEN_Gui_SetCharSet(abGui, C.CString("UTF-8"))
	C.GWEN_Gui_SetGui(abGui)

	C.GWEN_Gui_SetGetPasswordFn(abGui, (C.GWEN_GUI_GETPASSWORD_FN)(unsafe.Pointer(C.goAqbankingGetPasswordFn_cgo)))

	return &gui{
		abGui,
	}
}

func newNonInteractiveGui() *gui {
	return newGui(false)
}

func (g *gui) attach(aq *AQBanking) {
	C.AB_Gui_Extend(g.ptr, aq.ptr)
}

// RegisterPin registers a given Pin code with the aqbanking gui.
// pins can be added at any point in time prior to making a request.
// Pins are stored in memory only. When the process exits, all pins are forgotten.
func (ab *AQBanking) RegisterPin(pin Pin) {
	key := fmt.Sprintf("PIN_%v_%v", pin.BankCode(), pin.UserID())
	knownAqbankingPins[key] = pin.Pin()
}

func (g *gui) free() {
	C.GWEN_Gui_free(g.ptr)
}
