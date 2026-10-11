package models

import "github.com/cockroachdb/apd/v3"

// Tab Payment kasus FIRE (tiket 51). Uang / persen = *apd.Decimal (ADR-0034), nil = kosong.

// CicilanPayment - satu `.Policy.Payment.ListInstallment(n)` (T_LISTINSTALLMENT). DueDate = teks Pega `YYYYMMDD`.
type CicilanPayment struct {
	InstallmentNo                                                                    int
	DueDate                                                                          string
	Persen                                                                           *apd.Decimal // InstallmentPercentage, skala sampai 20
	Premium, RICommision, Deduction2, BrokerageFee, PPh, PPN, Discount, PaymentTotal *apd.Decimal
}

// MataUangPayment - satu `.OfferFacIn.CurrencyList(n)` tab Payment: nilai Payment mata uang (tab Premium Deduction),
// SumTotalPayment / NetPremium, dan cicilannya.
type MataUangPayment struct {
	Baris                       BarisDeduction
	SumTotalPayment, NetPremium *apd.Decimal
	Cicilan                     []CicilanPayment
}

// KasusPayment - data tab Payment tersimpan: Installment ("" belum diisi) dan cicilan per kode mata uang (No, DueDate,
// persen; nilai uang dihitung ulang tiap baca).
type KasusPayment struct {
	Installment string
	Cicilan     map[string][]CicilanPayment
}

// SimpanPayment - isi Save tab Payment.
type SimpanPayment struct {
	Installment int
	MataUang    []MataUangPayment
}
