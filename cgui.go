package aqbanking

/*
#include <stdio.h>
#include <stdint.h>
#include <gwenhywfar/cgui.h>
#include <gwenhywfar/gui.h>

// forward declaration for exported go callback fnc
int aqbankingGetPasswordFn(const char *, char *, int, int);

// see gwenhywfar/src/gui/gui.h:698 for details about this callback
int goAqbankingGetPasswordFn_cgo(
    GWEN_GUI *gui,
    uint32_t flags,
    const char *token,
    const char *title,
    const char *text,
    char *buffer,
    int minLen,
    int maxLen,
    uint32_t guiid) {
  return aqbankingGetPasswordFn(token, buffer, minLen, maxLen);
}
*/
import "C"
