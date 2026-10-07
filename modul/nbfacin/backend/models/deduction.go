package models

import "github.com/cockroachdb/apd/v3"

// Tab Premium Deduction / Brokerage Fee kasus FIRE (tiket 50). Uang / persen = *apd.Decimal (ADR-0034), nil = kosong.

// BarisDeduction - satu `.OfferFacIn.CurrencyList(n).Policy.Payment` (T_CURRENCYLIST).
type BarisDeduction struct {
	Currency                                               string
	Premium, Commision, Deduction2, BrokerageFee, PPh, PPN *apd.Decimal
}

// KasusDeduction - data tab Premium Deduction yang tersimpan: % Premium Deduction (`RICommision`), % Brokerage Fee,
// centang Exclude PPN (teks Pega apa adanya: "" belum pernah, "true", "false"; butir 68.1), dan potret CurrencyList saat Save.
type KasusDeduction struct {
	RICommision, PctBrokerageFee *apd.Decimal
	PPnCheck                     string
	MataUang                     []BarisDeduction
}

// TarifPajak - baris POOLDATA.TAX ber-TANGGAL terbaru (persen).
type TarifPajak struct {
	PPH, PPN *apd.Decimal
}
