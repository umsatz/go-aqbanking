package aqbanking

// Pin is a interface to support pluggable pin loading.
// The examples read the pin from a pins.json, which
// is extremely insecure and should never be used in production
type Pin interface {
	BankCode() string
	UserID() string
	Pin() string
}
