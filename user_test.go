package aqbanking

import (
	"fmt"
	"testing"
)

func TestRemoveUser(t *testing.T) {
	aq, freeAQ := newAQBankingTestInstance(t)
	if aq == nil {
		return
	}
	defer freeAQ()

	user := User{
		-1,
		"1",
		"2",
		"3",
		"4",
		"5",
		300,
		nil,
	}

	if err := aq.AddPinTanUser(&user); err != nil {
		t.Fatalf("unable to create aqbanking instance: %v", err)
	}

	if users, _ := aq.Users(); len(users.Users) != 1 {
		t.Fatalf("unable to create user.")
	}

	user.Remove(aq)

	if users, _ := aq.Users(); len(users.Users) != 0 {
		t.Fatalf("unable to remove user.")
	}
}

func TestAddUserAndListUsers(t *testing.T) {
	aq, freeAQ := newAQBankingTestInstance(t)
	if aq == nil {
		return
	}
	defer freeAQ()

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

	testAttrs := map[string][]string{
		"BankCode":    []string{user.BankCode, loadedUser.BankCode},
		"UserId":      []string{user.UserID, loadedUser.UserID},
		"CustomerId":  []string{user.CustomerID, loadedUser.CustomerID},
		"Name":        []string{user.Name, loadedUser.Name},
		"ServerUri":   []string{user.ServerURI, loadedUser.ServerURI},
		"HbciVersion": []string{fmt.Sprintf("%d", user.HbciVersion), fmt.Sprintf("%d", loadedUser.HbciVersion)},
	}

	for key, values := range testAttrs {
		if values[0] != values[1] {
			t.Fatalf("wrong value for attribute %v. expected '%v' got '%v'", key, values[0], values[1])
		}
	}
}
