package main

import (
	"fmt"
	"os"
	"testing"
)

func TestAddUserAndListUsers(t *testing.T) {
	defer os.RemoveAll("./tmp")

	aq, err := NewAQBanking("user tests", "./tmp")
	if err != nil {
		t.Fatalf("unable to create aqbanking instance: %v", err)
	}
	defer aq.Free()

	user := User{
		-1,
		"123456789",
		"123456789",
		"12030000",
		"A Bank",
		"https://hbci.example.com/hbciservlet",
		300,
		nil,
	}

	if err := aq.AddPinTanUser(&user); err != nil {
		t.Fatalf("unable to create aqbanking instance: %v", err)
	}

	users, err := aq.Users()
	if err != nil {
		t.Fatalf("unable to list users: %v", err)
	}

	if len(users.Users) != 1 {
		t.Fatalf("wrong number of users returned. expected 1, got %d", len(users.Users))
	}

	loadedUser := users.Users[0]

	var testAttrs map[string][]string = map[string][]string{
		"BankCode":    []string{user.BankCode, loadedUser.BankCode},
		"UserId":      []string{user.UserId, loadedUser.UserId},
		"CustomerId":  []string{user.CustomerId, loadedUser.CustomerId},
		"Name":        []string{user.Name, loadedUser.Name},
		"ServerUri":   []string{user.ServerUri, loadedUser.ServerUri},
		"HbciVersion": []string{fmt.Sprintf("%d", user.HbciVersion), fmt.Sprintf("%d", loadedUser.HbciVersion)},
	}

	for key, values := range testAttrs {
		if values[0] != values[1] {
			t.Fatalf("wrong value for attribute %v. expected '%v' got '%v'", key, values[0], values[1])
		}
	}
	// if loadedUser.BankCode != user.BankCode {

	// }
}
