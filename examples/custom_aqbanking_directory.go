package main

import (
	"fmt"
	"log"

	aq "github.com/umsatz/go-aqbanking"
)

func main() {
	gui := aq.NewNonInteractiveGui()
	defer gui.Free()

	aq, err := aq.NewAQBanking("custom", "./tmp")
	if err != nil {
		log.Fatal("unable to init aqbanking: %v", err)
	}
	gui.Attach(aq)
	defer aq.Free()

	fmt.Printf("using aqbanking %d.%d.%d\n",
		aq.Version.Major,
		aq.Version.Minor,
		aq.Version.Patchlevel,
	)
}
