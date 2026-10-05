// Package models memuat bentuk data modul Treaty Exchange Yearly - tabel warisan `POOLDATA.TREATYEXCHANGEYEARLY`
// (kelas Pega `ASM-FW-GISFW-Int-TREATYEXCHANGEYEARLY`, dibaca `BrowseTreatyExchangeYearly_RD`): kurs tahunan per tahun
// treaty, mata uang, dan quarter. Modul di luar korpus (perintah work owner 05-10-2026).
package models

// Kurs - satu baris `TREATYEXCHANGEYEARLY`, ditambah nama mata uang (view `CURRENCY`, dibaca saja).
type Kurs struct {
	// Kunci - ROWID baris: ID warisan tidak unik (10114-10116 masing-masing dua baris di DEV), jadi Edit memakai ini.
	Kunci        string `json:"kunci"`
	ID           string `json:"id"`
	TreatyYear   string `json:"treatyYear"`
	IDCurrency   string `json:"idCurrency"`
	Currency     string `json:"currency"`
	CurrencyName string `json:"currencyName"`
	// StartDate / EndDate - teks warisan format Pega `YYYYMMDDTHHMMSS.mmm GMT`; Mulai / Akhir - tampilan `DD-MM-YYYY`.
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Mulai     string `json:"mulai"`
	Akhir     string `json:"akhir"`
	ToIDR     string `json:"toIdr"`
	ToUSD     string `json:"toUsd"`
	Quarter   string `json:"quarter"`
	UserID    string `json:"userId"`
	// DateIU - waktu ubah, DateIn - waktu input (format Pega); Diubah - tampilan DateIU WIB `DD-MM-YYYY HH:MM`.
	DateIU string `json:"dateIu"`
	DateIn string `json:"dateIn"`
	Diubah string `json:"diubah"`
}

// MataUang - satu baris view `CURRENCY` (pilihan Currency; dibaca saja).
type MataUang struct {
	ID   string `json:"id"`
	Kode string `json:"kode"`
	Nama string `json:"nama"`
}

// Isian - isian form Add / Edit. Kunci kosong = baris baru. Tanggal dari layar `YYYY-MM-DD`.
type Isian struct {
	Kunci      string `json:"kunci"`
	TreatyYear string `json:"treatyYear"`
	IDCurrency string `json:"idCurrency"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
	ToIDR      string `json:"toIdr"`
	ToUSD      string `json:"toUsd"`
	Quarter    string `json:"quarter"`
}

// BatasTeks - kolom selain ID VARCHAR2(255) (katalog DEV 05-10-2026).
const BatasTeks = 255

// Quarter yang sah: `0` tahunan (139 dari 140 baris DEV; pembaca Pega menyaring `Quarter = "0"`), `1`-`4` triwulan.
var Quarter = []string{"0", "1", "2", "3", "4"}
