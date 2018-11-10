package aqbanking

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDefaultAQBankingInstance(t *testing.T) {
	aq, err := DefaultAQBanking()
	if err != nil {
		t.Fatalf("unable to create aqbanking instance")
	}
	aq.Free()
}

func TestNewAQBankingInstance(t *testing.T) {
	tmp, err := ioutil.TempDir("", "aqbanking")
	if err != nil {
		t.Errorf("unable to create temporary dir: %v", err)
		return
	}
	defer os.RemoveAll(tmp)

	aq, err := NewAQBanking("example", tmp)

	if err != nil {
		t.Fatalf("unable to create custom aqbanking instance")
	}
	aq.Free()
}
