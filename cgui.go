package main

/*
#include <stdio.h>
#include <stdint.h>
#include <gwenhywfar/cgui.h>
#include <gwenhywfar/gui.h>

int callMeOnGo_cgo(GWEN_GUI *gui,
    uint32_t flags,
    const char *token,
    const char *title,
    const char *text,
    char *buffer,
    int minLen,
    int maxLen,
    uint32_t guiid) {
  return callMeOnGo(token, buffer, minLen, maxLen);
}
*/
import "C"
