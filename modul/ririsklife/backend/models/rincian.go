// Package models memuat bentuk data modul R/I Risk (`ririsklife`) - section Pega `InboxSummaryRIRisk` (kelas
// `ASM-FW-GISFW-Int-RIRISK_LIFE_SUMMARY`, judul "R/I RISK SUMMARY" b375): ringkasan risk reasuransi life (tabel
// `RIRISK_LIFE_SUMMARY`) dan rinciannya (section `InboxRIRisk`, kelas `ASM-FW-GISFW-Int-RI_RISK_LIFE`, judul
// "R/I RISK DETAIL" b382) di tabel `RIRISK_LIFE` - satu tabel per jenis data (keputusan work owner 08-10-2026 K1/K2,
// migrasi inti 935-940: tabel Pega M_RIRISK_LIFE* berganti nama, JSONDATA dibuang).
package models

// Ringkasan - satu baris `RIRISK_LIFE_SUMMARY` (grid `BrowseRIRiskSummary` b7127).
type Ringkasan struct {
	// ID - `.ID` b1073 (pxTextInput, disabled); site || LPAD(seq, 6, '0').
	ID string `json:"id"`
	// UsedBy - `.USEDBY` b1256 (wajib b1274, label "R/I RISK NAME" b1225), kolom grid "R/I RISK NAME" b7147.
	UsedBy string `json:"usedby"`
	// OperatorID - `.OPERATORID` b8180, kolom "MODIFY OPERATOR" b7295.
	OperatorID string `json:"operatorId"`
	// ModifiedDate - `.MODIFIEDDATE` b8358, kolom "MODIFY DATE" b7396; teks Pega `yyyyMMdd'T'HHmmss.SSS 'GMT'` apa
	// adanya (data DEV `20191113T025753.044 GMT`).
	ModifiedDate string `json:"modifiedDate"`
	// Diubah - tampilan WIB `DD-MM-YYYY` dari ModifiedDate; bentuk lain apa adanya.
	Diubah string `json:"diubah"`
	// JumlahRincian - baris `RIRISK_LIFE` ber-IDUSEDBY = ID (hanya diisi Buka, untuk dialog Delete).
	JumlahRincian *int `json:"jumlahRincian,omitempty"`
}

// Isian - isian form ringkasan (Save b1796 -> `AddToListSummary_Act` b1819; Edit -> `EditListSummary_DT` b8616). ID dari
// jalur (Edit) atau dibentuk server. Save / Cancel DITAMPILKAN atas keputusan work owner 08-10-2026 (di XML wadahnya
// `1=2` b1511 - RALAT R1 MODUL.md).
type Isian struct {
	UsedBy string `json:"usedby"`
}

// Rincian - satu baris `RIRISK_LIFE` (grid `BrowseRIRiskLife_RD` b7612). CONTRACT, YEAR, MONTH teks angka (kolom
// VARCHAR2(10) warisan), RISK teks KANONIK (titik desimal, tanpa nol depan / ekor - `TO_CHAR(..., 'TM9')` ber-NLS
// eksplisit), kosong = NULL. AGE kolom view lama yang tidak pernah terisi dan tidak ada di XML - tidak dibaca/ditulis.
type Rincian struct {
	ID string `json:"id"`
	// IDUsedBy - `TempIDUsedBy.ID` b1524: ID ringkasan pemilik.
	IDUsedBy string `json:"idUsedBy"`
	// UsedBy - `TempIDUsedBy.USEDBY` b1727: salinan nama ringkasan.
	UsedBy   string `json:"usedby"`
	Contract string `json:"contract"`
	Year     string `json:"year"`
	Month    string `json:"month"`
	Risk     string `json:"risk"`
}

// Saringan - filter dan urutan grid ringkasan (50 baris per halaman, `pyPageSize` b10030).
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

// UkuranHalaman - `pyPageSize` 50 grid ringkasan (`InboxSummaryRIRisk` b10030).
const UkuranHalaman = 50

// UkuranHalamanRincian - grid R/I RISK DETAIL: `pyPageSize` Other b10243 = `pyPageSizeOther` 200 b10169 (`InboxRIRisk`)
// - BERBEDA dengan R/I Comm Life (50).
const UkuranHalamanRincian = 200

// Kolom urut grid ringkasan (nama kunci kueri -> kolom tabel).
const (
	UrutID       = "id"
	UrutUsedBy   = "usedby"
	UrutOperator = "operatorid"
	UrutTanggal  = "modifieddate"
	// UrutBawaan - ID menaik: sort kolom pertama ASC, urutan 1 (`InboxSummaryRIRisk` b9806/b9812).
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

// HasilHapus - Delete (`DeleteSummaryDetail` b9648, DeleteID=.ID b9664): ringkasan beserta rinciannya.
type HasilHapus struct {
	ID              string `json:"id"`
	RincianTerhapus int    `json:"rincianTerhapus"`
}

// BatasNama - panjang R/I RISK NAME (byte) = lebar `RIRISK_LIFE_SUMMARY.USEDBY` VARCHAR2(200) (migrasi inti 936, pola
// ricommlife 931; data DEV maks 81 byte). `RIRISK_LIFE.USEDBY` warisan VARCHAR2(1000) memuatnya. XML tanpa pyMax.
const BatasNama = 200

// BatasID - `RIRISK_LIFE_SUMMARY.ID` dan `RIRISK_LIFE.IDUSEDBY` (VARCHAR2(10) / VARCHAR2(100) warisan; ID ringkasan
// paling panjang 10).
const BatasID = 10

// BatasIDRincian - `RIRISK_LIFE.ID` VARCHAR2(6) warisan.
const BatasIDRincian = 6
