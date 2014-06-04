package main

import (
	"fmt"
	"log"

	aq "github.com/umsatz/go-aqbanking"
)

func main() {
	gui := aq.NewNonInteractiveGui()
	defer gui.Free()

	aq, err := aq.DefaultAQBanking()
	// alternativly, customize the aqbanking path:
	// aq, err := aq.NewAQBanking("custom", "./tmp")
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

	userCollection, err := aq.Users()
	if err != nil {
		log.Fatal("unable to list users: %v", err)
	}
	defer userCollection.Free()
	fmt.Printf("found %d users.\n", len(userCollection.Users))

	accountCollection, err := aq.Accounts()
	if err != nil {
		log.Fatal("unable to list accounts: %v", err)
	}
	defer accountCollection.Free()

	fmt.Printf("found %d accounts.\n", len(accountCollection.Accounts))
}
