// Package models memuat bentuk data modul R/I Rate Life (`riratelife`) - section Pega `InboxSummaryRIRate` (kelas
// `ASM-FW-GISFW-Int-RATE_LIFE_SUMMARY`, judul "R/I RATE SUMMARY" b382): ringkasan rate reasuransi life
// (kolom `M_RATE_LIFE_SUMMARY`, migrasi inti 927/928 - RALAT R6) dan rincian rate-nya (kolom `M_RATE_LIFE`, migrasi
// inti 929/930 - RALAT R7). Perintah work owner 05-10-2026: "Buat Menu baru Namanya R/I Rate Life pada Master Treaty, menu ini
// bisa CRUD untuk simpan data ke tabel RATE_LIFE_SUMMARY, panduannya xml yang saya berikan".
package models

// Ringkasan - satu baris `M_RATE_LIFE_SUMMARY` (grid `BrowseRateLifeSummary` b9568).
type Ringkasan struct {
	// ID - `.ID` b1089 (pxTextInput, disabled selalu b1104-b1105); VARCHAR2(10).
	ID string `json:"id"`
	// UsedBy - `.USEDBY` b1273, label "R/I RATE NAME" b1242, wajib b1285.
	UsedBy string `json:"usedby"`
	// OperatorID - `.OPERATORID` b10709, kolom "MODIFY OPERATOR" b9748.
	OperatorID string `json:"operatorId"`
	// ModifiedDate - `.MODIFIEDDATE` b10879, kolom "MODIFY DATE" b9894 (format `Date-Short-Custom-YYYY` b10929). Teks
	// apa adanya; baris yang ditulis modul ini berformat Pega `YYYYMMDDTHHMMSS.mmm GMT` (ASUMSI, MODUL.md).
	ModifiedDate string `json:"modifiedDate"`
	// Diubah - tampilan tanggal WIB `DD-MM-YYYY` dari ModifiedDate; bentuk lain apa adanya.
	Diubah string `json:"diubah"`
	// JumlahRate - baris `M_RATE_LIFE` ber-IDUSEDBY = ID (hanya diisi Buka, untuk dialog Delete).
	JumlahRate *int `json:"jumlahRate,omitempty"`
}

// Isian - isian form ringkasan (Save b1809 -> `AddToListSummary_Act` b1833). ID dari jalur (Edit) atau sequence.
type Isian struct {
	UsedBy string `json:"usedby"`
}

// Rate - satu baris `M_RATE_LIFE` (kolom sama dengan view `RATE_LIFE` lama, RALAT R7). Seluruhnya teks.
type Rate struct {
	ID       string `json:"id"`
	IDUsedBy string `json:"idUsedBy"`
	UsedBy   string `json:"usedby"`
	Gender   string `json:"gender"`
	Contract string `json:"contract"`
	Age      string `json:"age"`
	// Rate - teks berkoma desimal seperti mayoritas data DEV (90.436 dari 98.305 baris).
	Rate string `json:"rate"`
}

// Saringan - filter dan urutan grid ringkasan (filter + sort `BrowseRateLifeSummary`, 50 baris per halaman b12645).
type Saringan struct {
	ID     string
	UsedBy string
	// Urut - salah satu KolomUrut; kosong = UrutBawaan.
	Urut string
	// Turun - urutan menurun.
	Turun bool
	// Halaman - mulai 1.
	Halaman int
}

// UkuranHalaman - `pyPageSize` 50 b12645; dipakai juga grid Rate Detail.
const UkuranHalaman = 50

// Kolom urut grid ringkasan (nama kunci kueri -> kolom view).
const (
	UrutID       = "id"
	UrutUsedBy   = "usedby"
	UrutOperator = "operatorid"
	UrutTanggal  = "modifieddate"
	// UrutBawaan - ID menurun (ringkasan terbaru di atas). ASUMSI: urutan bawaan RD Pega tidak ada di XML.
	UrutBawaan = UrutID
)

// KolomUrut - kolom yang boleh diurutkan.
var KolomUrut = []string{UrutID, UrutUsedBy, UrutOperator, UrutTanggal}

// Halaman - satu halaman grid.
type Halaman[T any] struct {
	Daftar  []T `json:"daftar"`
	Total   int `json:"total"`
	Halaman int `json:"halaman"`
	Ukuran  int `json:"ukuran"`
}

// HasilHapus - Delete (`DeleteSummaryDetail` b12262, DeleteID=.ID b12279): ringkasan beserta rate-nya.
type HasilHapus struct {
	ID           string `json:"id"`
	RateTerhapus int    `json:"rateTerhapus"`
}

// BatasNama - panjang R/I RATE NAME (byte). ASUMSI: kolom JSON tanpa batas; 200 cukup untuk nama DEV.
const BatasNama = 200

// BatasID - `M_RATE_LIFE_SUMMARY.ID` dan `M_RATE_LIFE.ID` VARCHAR2(10).
const BatasID = 10
