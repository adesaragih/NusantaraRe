// Package models memuat bentuk data modul R/I Comm Life (`ricommlife`) - section Pega `InboxSummaryRIComm` (kelas
// `ASM-FW-GISFW-Int-RICOMM_LIFE_SUMMARY`, judul "R/I COMM SUMMARY" b382): ringkasan komisi reasuransi life
// (kolom `M_RICOMM_LIFE_SUMMARY`) dan rinciannya (section `InboxRIComm`, kelas `ASM-FW-GISFW-Int-RI_COMM_LIFE`, judul
// "R/I COMM DETAIL" b349) di kolom `M_RICOMM_LIFE` - satu tabel per jenis data (migrasi inti 931-934, RALAT R1).
// Perintah work owner 06-10-2026: "membuat modul baru di Master Treaty dengan nama R/I Comm Life ... konsepnya hampir
// sama dengan menu R/I Rate, membuat CRUD, dan detail bisa di save dan edit".
package models

// Ringkasan - satu baris `M_RICOMM_LIFE_SUMMARY` (grid `BrowseRICommSummary`).
type Ringkasan struct {
	// ID - `.ID` b1091 (pxTextInput, disabled b1109-b1110); site || LPAD(seq, 6, '0').
	ID string `json:"id"`
	// UsedBy - `.USEDBY` b1274 (wajib b1289), kolom grid "R/I COMM NAME" b7220.
	UsedBy string `json:"usedby"`
	// OperatorID - `.OPERATORID` b8348, kolom "MODIFY OPERATOR" b7377.
	OperatorID string `json:"operatorId"`
	// ModifiedDate - `.MODIFIEDDATE` b8526 (format `Date-Short-Custom-YYYY`), kolom "MODIFY DATE" b7521; teks Pega
	// `yyyyMMdd'T'HHmmss.SSS 'GMT'` apa adanya.
	ModifiedDate string `json:"modifiedDate"`
	// Diubah - tampilan WIB `DD-MM-YYYY` dari ModifiedDate; bentuk lain apa adanya.
	Diubah string `json:"diubah"`
	// JumlahKomisi - baris `M_RICOMM_LIFE` ber-IDUSEDBY = ID (hanya diisi Buka, untuk dialog Delete).
	JumlahKomisi *int `json:"jumlahKomisi,omitempty"`
}

// Isian - isian form ringkasan (Save b1799 -> `AddToListSummary_Act` b1823). ID dari jalur (Edit) atau dibentuk server.
type Isian struct {
	UsedBy string `json:"usedby"`
}

// Komisi - satu baris `M_RICOMM_LIFE` (grid `BrowseRICommLife_RD`). Angka sebagai teks KANONIK (titik desimal,
// tanpa nol depan / nol ekor - bentuk yang dikembalikan Oracle `TO_CHAR(..., 'TM9')`), kosong = NULL.
type Komisi struct {
	ID string `json:"id"`
	// IDUsedBy - `TempIDUsedBy.ID` b1509: ID ringkasan pemilik.
	IDUsedBy string `json:"idUsedBy"`
	// UsedBy - `TempIDUsedBy.USEDBY` b1718: salinan nama ringkasan.
	UsedBy   string `json:"usedby"`
	Contract string `json:"contract"`
	Year     string `json:"year"`
	Comm     string `json:"comm"`
}

// Saringan - filter dan urutan grid ringkasan (50 baris per halaman, `pyPageSize` b10081).
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

// UkuranHalaman - `pyPageSize` 50: ringkasan b10081 dan R/I COMM DETAIL (`InboxRIComm` b9716).
const UkuranHalaman = 50

// Kolom urut grid ringkasan (nama kunci kueri -> kolom tabel).
const (
	UrutID       = "id"
	UrutUsedBy   = "usedby"
	UrutOperator = "operatorid"
	UrutTanggal  = "modifieddate"
	// UrutBawaan - ID menaik: sort kolom pertama ASC, urutan 1 (`InboxSummaryRIComm` b9857/b9863).
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

// HasilHapus - Delete (`DeleteSummaryDetail` b9704, DeleteID=.ID b9719): ringkasan beserta rinciannya.
type HasilHapus struct {
	ID             string `json:"id"`
	KomisiTerhapus int    `json:"komisiTerhapus"`
}

// BatasNama - panjang R/I COMM NAME (byte) = lebar `M_RICOMM_LIFE.USEDBY` dan `M_RICOMM_LIFE_SUMMARY.USEDBY`
// VARCHAR2(200) (933 / 931, sama dengan 924). XML tidak memberi batas
// (`.USEDBY` b1274 tanpa pyMax) - preseden riratelife (200), [penyimpangan sadar] STRUKTUR-TABEL-RICOMMLIFE.md.
const BatasNama = 200

// BatasID - `M_RICOMM_LIFE.ID` / `.IDUSEDBY` dan `M_RICOMM_LIFE_SUMMARY.ID` VARCHAR2(10).
const BatasID = 10
