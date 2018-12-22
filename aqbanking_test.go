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
	aq, freeAQ := newAQBankingTestInstance(t)
	if aq == nil {
		return
	}
	freeAQ()
}

// Creates a new AQBanking with a temporary dir.
// Caller must call the returned free function for freeing resources.
func newAQBankingTestInstance(t *testing.T) (aq *AQBanking, free func()) {
	tmp, err := ioutil.TempDir("", "aqbanking")
	if err != nil {
		t.Fatalf("unable to create temporary dir: %v", err)
		return
	}

	aq, err = NewAQBanking("temporary", tmp)
	if err != nil {
		os.RemoveAll(tmp)
		t.Fatalf("unable to create aqbanking instance: %v", err)
		return
	}

	free = func() {
		os.RemoveAll(tmp)
		err = aq.Free()
		if err != nil {
			t.Fatalf("unable to free aqbanking instance: %v", err)
		}
	}

	return
}
