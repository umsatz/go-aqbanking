package aqbanking

// Pin is a interface to support pluggable pin loading.
// The examples read the pin from a pins.json, which
// is extremely insecure and should never be used in production
type Pin interface {
	BankCode() string
	UserID() string
	Pin() string
}

type pin struct {
	Blz string `json:"blz"`
	UID string `json:"uid"`
	PIN string `json:"pin"`
}

func (p *pin) BankCode() string {
	return p.Blz
}

func (p *pin) UserID() string {
	return p.UID
}

func (p *pin) Pin() string {
	return p.PIN
}
