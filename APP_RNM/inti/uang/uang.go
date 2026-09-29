// Package uang memuat tipe dasar uang dan rasio yang dipakai bersama.
//
// Refactor bentuk B (30-09-2026): dulu `models/money.go`. Fase 0 - scaffold:
// hanya tipe dasar uang dan rasio. Nol aturan dagang.
package uang

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/utils"
)

// ⛔ float64 tidak boleh muncul di jalur uang mana pun - termasuk JSON keluar
// dan React (ADR-U-0003, ADR-U-0016). Karena itu tidak satu pun tipe di
// berkas ini memakai float, dan MarshalJSON menulis desimal sebagai TEKS.

var (
	// ErrMataUangBerbeda muncul bila dua Money dengan mata uang berbeda
	// dijumlahkan. Ia tidak pernah didiamkan menjadi konversi diam-diam:
	// kurs adalah keputusan tersendiri, bukan efek samping penjumlahan.
	ErrMataUangBerbeda = errors.New("mata uang berbeda tidak dapat dijumlahkan")
	// ErrUangKosong muncul bila operasi aritmetika menyentuh nilai kosong.
	ErrUangKosong = errors.New("nilai uang kosong")
	// ErrAngkaJSON muncul bila JSON masuk membawa uang sebagai angka, bukan
	// teks. Angka JSON dibaca sebagai float64 oleh pustaka mana pun, dan itu
	// persis yang dilarang.
	ErrAngkaJSON = errors.New("nilai uang harus berupa teks desimal, bukan angka JSON")
)

// Money adalah nilai uang beserta mata uangnya.
//
// Money dan Ratio adalah DUA TIPE BERBEDA yang tidak dapat dijumlahkan satu
// sama lain (ADR-F-0004). Go tidak punya operator untuk struct, sehingga
// satu-satunya jalan aritmetika adalah metode di bawah - dan metode itu
// memaksa kesamaan mata uang.
type Money struct {
	Amount   *apd.Decimal
	Currency string
}

// NewMoney membaca nilai uang dari teks desimal (ADR-U-0022: konversi sekali
// saat masuk).
func NewMoney(amount, currency string) (Money, error) {
	d, err := utils.ParseDecimal(amount)
	if err != nil {
		return Money{}, err
	}
	return Money{Amount: d, Currency: currency}, nil
}

// Kosong menyatakan nilai belum terisi. Kolom kosong dan kolom bernilai nol
// adalah dua hal berbeda (ADR-U-0027).
func (m Money) Kosong() bool { return m.Amount == nil }

// Add menjumlahkan dua nilai uang bermata uang sama.
func (m Money) Add(other Money) (Money, error) {
	if m.Kosong() || other.Kosong() {
		return Money{}, ErrUangKosong
	}
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("%w: %q + %q", ErrMataUangBerbeda, m.Currency, other.Currency)
	}
	hasil := new(apd.Decimal)
	if _, err := utils.DecimalContext().Add(hasil, m.Amount, other.Amount); err != nil {
		return Money{}, err
	}
	return Money{Amount: hasil, Currency: m.Currency}, nil
}

// Sub mengurangkan nilai uang bermata uang sama.
func (m Money) Sub(other Money) (Money, error) {
	if m.Kosong() || other.Kosong() {
		return Money{}, ErrUangKosong
	}
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("%w: %q - %q", ErrMataUangBerbeda, m.Currency, other.Currency)
	}
	hasil := new(apd.Decimal)
	if _, err := utils.DecimalContext().Sub(hasil, m.Amount, other.Amount); err != nil {
		return Money{}, err
	}
	return Money{Amount: hasil, Currency: m.Currency}, nil
}

// String menulis nilai uang sebagai teks.
func (m Money) String() string {
	if m.Kosong() {
		return ""
	}
	return utils.FormatDecimal(m.Amount) + " " + m.Currency
}

type uangJSON struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// MarshalJSON menulis jumlah sebagai TEKS desimal, tidak pernah sebagai angka
// JSON. Angka JSON akan dibaca sebagai float64 di sisi React.
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(uangJSON{Amount: utils.FormatDecimal(m.Amount), Currency: m.Currency})
}

// bacaDesimalJSON membaca satu medan desimal dari JSON.
//
// Medan yang ABSEN dan medan bernilai null sama-sama berarti kosong, bukan
// galat: kolom kosong dan kolom bernilai nol adalah dua hal berbeda
// (ADR-U-0027). Yang ditolak hanya angka JSON, sebab angka JSON dibaca
// sebagai float64 oleh pustaka mana pun.
func bacaDesimalJSON(mentah json.RawMessage) (*apd.Decimal, error) {
	teksMentah := string(mentah)
	if len(mentah) == 0 || teksMentah == "null" {
		return nil, nil
	}
	if mentah[0] != '"' {
		return nil, fmt.Errorf("%w: %s", ErrAngkaJSON, teksMentah)
	}
	var teks string
	if err := json.Unmarshal(mentah, &teks); err != nil {
		return nil, err
	}
	if teks == "" {
		return nil, nil
	}
	return utils.ParseDecimal(teks)
}

// UnmarshalJSON menolak angka JSON secara tegas.
func (m *Money) UnmarshalJSON(b []byte) error {
	var mentah struct {
		Amount   json.RawMessage `json:"amount"`
		Currency string          `json:"currency"`
	}
	if err := json.Unmarshal(b, &mentah); err != nil {
		return err
	}
	d, err := bacaDesimalJSON(mentah.Amount)
	if err != nil {
		return err
	}
	m.Currency = mentah.Currency
	m.Amount = d
	return nil
}

// Ratio adalah pangsa atau persentase - BUKAN uang, dan tidak dapat
// dijumlahkan dengan Money (ADR-F-0004).
//
// Scale datang dari satu resolver, bukan ditebak per pemakaian. Resolver itu
// belum ada di Fase 0: ia lahir bersama tiket yang memerlukannya.
type Ratio struct {
	Value *apd.Decimal
	Scale int32
}

// NewRatio membaca rasio dari teks desimal.
func NewRatio(value string, scale int32) (Ratio, error) {
	d, err := utils.ParseDecimal(value)
	if err != nil {
		return Ratio{}, err
	}
	return Ratio{Value: d, Scale: scale}, nil
}

// Kosong menyatakan rasio belum terisi.
func (r Ratio) Kosong() bool { return r.Value == nil }

// String menulis rasio sebagai teks.
func (r Ratio) String() string { return utils.FormatDecimal(r.Value) }

// MarshalJSON menulis rasio sebagai teks desimal, sealasan dengan Money.
func (r Ratio) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Value string `json:"value"`
		Scale int32  `json:"scale"`
	}{Value: utils.FormatDecimal(r.Value), Scale: r.Scale})
}

// UnmarshalJSON menolak angka JSON, sama tegasnya dengan Money. Rasio ikut
// jalur uang: pangsa yang dibaca sebagai float64 merusak nilai yang
// dikalikannya.
func (r *Ratio) UnmarshalJSON(b []byte) error {
	var mentah struct {
		Value json.RawMessage `json:"value"`
		Scale int32           `json:"scale"`
	}
	if err := json.Unmarshal(b, &mentah); err != nil {
		return err
	}
	d, err := bacaDesimalJSON(mentah.Value)
	if err != nil {
		return err
	}
	r.Value = d
	r.Scale = mentah.Scale
	return nil
}
