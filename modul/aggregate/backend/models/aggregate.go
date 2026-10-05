// Package models memuat bentuk data modul Aggregate.
//
// Satu baris `POOLDATA.AGGREGATE` = satu zona penilaian (Assessment Zone) x satu coverage x satu treaty type untuk
// satu ceding pada satu tanggal posisi (As At). Baris lahir dari unggahan CSV (Pega `UploadCSVAggregate_Act`) dan
// disimpan `SaveAggregate_Act`. Baris dipertukarkan sebagai peta kolom Oracle -> teks: angka dalam bentuk desimal
// bertitik tanpa pemisah ribuan, tanggal `DD-MM-YYYY`.
package models

// Jenis isi satu kolom AGGREGATE.
type Jenis int

// Jenis kolom.
const (
	Teks Jenis = iota
	Angka
	Tanggal
)

// Kolom - satu kolom AGGREGATE di grid pratinjau dan rincian.
type Kolom struct {
	// Nama - nama kolom Oracle, juga kunci baris.
	Nama  string
	Jenis Jenis
	// Lebar - VARCHAR2 dalam BYTE (kolom teks).
	Lebar int
	// CSV - nomor kolom berkas CSV, 1-38 (`UploadCSVAggregate_Act` langkah 6: `.COLUMN1`-`.COLUMN38`); 0 = bukan dari CSV.
	CSV int
	// Ubah - dapat diubah di grid pratinjau (`ShowAggregateList`, sel readOnly=false).
	Ubah bool
}

// KolomGrid - kolom grid pratinjau `ShowAggregateList` (TempCSV.pxResults), urutan layar Pega.
var KolomGrid = []Kolom{
	{Nama: "ASSESMENT_ZONE", Jenis: Teks, Lebar: 100, CSV: 1},
	{Nama: "M_TREATY_ID", Jenis: Teks, Lebar: 15},
	{Nama: "TREATY_TYPE", Jenis: Teks, Lebar: 100, CSV: 2, Ubah: true},
	{Nama: "COVERAGE", Jenis: Teks, Lebar: 100, CSV: 3, Ubah: true},
	{Nama: "AS_AT", Jenis: Tanggal, CSV: 4},
	{Nama: "UW_YEAR", Jenis: Teks, Lebar: 10, CSV: 5, Ubah: true},
	{Nama: "TREATYYEAR", Jenis: Teks, Lebar: 10, Ubah: true},
	{Nama: "CEDING_CODE", Jenis: Teks, Lebar: 20, CSV: 6, Ubah: true},
	{Nama: "CEDING_NAME", Jenis: Teks, Lebar: 100, CSV: 7, Ubah: true},
	{Nama: "CURRENCY", Jenis: Teks, Lebar: 10, CSV: 8, Ubah: true},
	{Nama: "TO_USD", Jenis: Angka, CSV: 9},
	{Nama: "NOR_BUILDINGS", Jenis: Angka, CSV: 10},
	{Nama: "BUILDINGS", Jenis: Angka, CSV: 11},
	{Nama: "NOR_STOCKS", Jenis: Angka, CSV: 12},
	{Nama: "STOCKS", Jenis: Angka, CSV: 13},
	{Nama: "NOR_MACHINERY", Jenis: Angka, CSV: 14},
	{Nama: "MACHINERY", Jenis: Angka, CSV: 15},
	{Nama: "NOR_OTHER_CONTENTS", Jenis: Angka, CSV: 16},
	{Nama: "OTHER_CONTENTS", Jenis: Angka, CSV: 17},
	{Nama: "NOR_CONSEQUENTIAL_LOSS", Jenis: Angka, CSV: 18},
	{Nama: "CONSEQUENTIAL_LOSS", Jenis: Angka, CSV: 19},
	{Nama: "NOR_RESIDENTIAL", Jenis: Angka, CSV: 20},
	{Nama: "RESIDENTIAL", Jenis: Angka, CSV: 21},
	{Nama: "NOR_COMMERCIAL", Jenis: Angka, CSV: 22},
	{Nama: "COMMERCIAL", Jenis: Angka, CSV: 23},
	{Nama: "NOR_INDUSTRIAL", Jenis: Angka, CSV: 24},
	{Nama: "INDUSTRIAL", Jenis: Angka, CSV: 25},
	{Nama: "NOR_AGRICULTURE", Jenis: Angka, CSV: 26},
	{Nama: "AGRICULTURE", Jenis: Angka, CSV: 27},
	{Nama: "NOR_MISCELLANEOUS", Jenis: Angka, CSV: 28},
	{Nama: "MISCELLANEOUS", Jenis: Angka, CSV: 29},
	{Nama: "NOR_UTILITIES", Jenis: Angka, CSV: 30},
	{Nama: "UTILITIES", Jenis: Angka, CSV: 31},
	{Nama: "TOTAL_NO_OF_RISK", Jenis: Angka, CSV: 32},
	{Nama: "TOTAL_IN_AMOUNT", Jenis: Angka, CSV: 33},
	{Nama: "TOTAL_IN_AMOUNT_IN_USD", Jenis: Angka, CSV: 34},
	{Nama: "RNM_SHARE", Jenis: Angka, CSV: 35},
	{Nama: "RNM_VALUE", Jenis: Angka, CSV: 36},
	{Nama: "RNM_VALUE_IN_USD", Jenis: Angka, CSV: 37},
	{Nama: "REMARK", Jenis: Teks, Lebar: 4000, CSV: 38, Ubah: true},
}

// Kolom jejak - diisi saat disimpan, tidak ada di grid pratinjau.
const (
	KolomID           = "ID"
	KolomTanggalInput = "TANGGAL_INPUT"
	KolomUserInput    = "USER_INPUT"
	// LebarUserInput - AGGREGATE.USER_INPUT VARCHAR2(100).
	LebarUserInput = 100
)

// Nilai tetap Pega.
const (
	// ZonaTotal - Assessment Zone baris jumlah per mata uang (`UploadCSVAggregate_Act` langkah 9); tidak disimpan.
	ZonaTotal = "Total :"
	// AwalanID - `ID` = AwalanID + nomor (`pzGenerateUniqueID(tools, "AGG")`).
	AwalanID = "AGG-"
	// MataUangIDR - kurs IDR = 1 / kurs IDR baris USD (`UploadCSVAggregate_Act` langkah 7.13-7.16).
	MataUangIDR = "IDR"
	MataUangUSD = "USD"
)

// TreatyTypeSah - Treaty Type yang diterima `SaveAggregate_Act` langkah 3.1.3.
var TreatyTypeSah = []string{"OR", "QS", "SURPLUS"}

// Baris - satu baris AGGREGATE: kolom Oracle -> teks.
type Baris map[string]string

// MasterTreaty - satu baris view `TREATYINDETAILJOINEDM` pilihan Master ID (`ChooseMasterID`).
type MasterTreaty struct {
	// ID - kunci baris view (RowKey `ID`).
	ID                 string `json:"id"`
	TreatyID           string `json:"treatyId"`
	CedingID           string `json:"cedingId"`
	Ceding             string `json:"ceding"`
	RnmShare           string `json:"rnmShare"`
	TreatyYear         string `json:"treatyYear"`
	ProportionType     string `json:"proportionType"`
	TreatyContractName string `json:"treatyContractName"`
	TreatyGroup        string `json:"treatyGroup"`
	SOB                string `json:"sob"`
}

// Kunci - satu baris daftar `GridDasbordAgg`: Tanggal Input (hari), Ceding Code, Ceding Name, Treaty Type, As At,
// UW Year. Rincian dan Delete bekerja atas seluruh baris AGGREGATE berkunci sama.
type Kunci struct {
	TanggalInput string `json:"tanggalInput"`
	CedingCode   string `json:"cedingCode"`
	CedingName   string `json:"cedingName"`
	TreatyType   string `json:"treatyType"`
	AsAt         string `json:"asAt"`
	UwYear       string `json:"uwYear"`
}

// Kelompok - satu baris daftar.
type Kelompok struct {
	Kunci
	// JumlahBaris - baris AGGREGATE berkunci ini.
	JumlahBaris int `json:"jumlahBaris"`
	// InputTerakhir - TANGGAL_INPUT terakhir, `DD-MM-YYYY HH:MM`.
	InputTerakhir string `json:"inputTerakhir"`
}

// IrisanRingkasan - satu daun chart: jumlah RNM Value (USD) satu Ceding, Treaty Type, dan Coverage pada satu As At.
// Layar menjumlahkannya lagi per Ceding lalu per Treaty Type (chart bertingkat, work owner 04-10-2026).
type IrisanRingkasan struct {
	CedingCode    string `json:"cedingCode"`
	CedingName    string `json:"cedingName"`
	TreatyType    string `json:"treatyType"`
	Coverage      string `json:"coverage"`
	RnmValueInUSD string `json:"rnmValueInUsd"`
}
