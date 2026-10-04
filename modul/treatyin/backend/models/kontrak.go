package models

// Bentuk kontrak dan versinya - tiket 14.
//
// ⛔ VERSI_KONTRAK punya 49 kolom; yang dimodelkan di sini hanya yang WAJIB
// ISI, ditambah pengenal dan nomor urutnya. Ke-37 kolom nullable lainnya
// adalah milik tiket yang membuatnya - mata uang (20), retensi (22), layer
// (31), dan seterusnya - dan menariknya ke sini sekarang berarti membuat
// masukan untuk ruas yang belum punya aturan. Tiket 14 menyebut lingkupnya
// sendiri: "membuat sebuah kontrak baru, mengisi KEPALA versinya".

import "github.com/cockroachdb/apd/v3"

// Dua nilai sah SIFAT_PROPORSI (INV-29).
//
// ⛔ Himpunan TERTUTUP, dan ia tidak disimpan sebagai tabel acuan - berbeda
// dengan keenam himpunan ADR-0038 yang DAPAT bertambah. Basis data tidak
// menegakkannya (ADR-0056); di sinilah ia ditegakkan.
const (
	SifatProporsional    = "PROPORSIONAL"
	SifatNonProporsional = "NON_PROPORSIONAL"
)

// Kontrak adalah LAPISAN BEKU yang seluruh versinya bagi tanpa menyalinnya
// (ADR-0040): cedant, asal bisnis, sifat proporsi, dan periode.
type Kontrak struct {
	ID                   int64  `json:"id"`
	NomorKontrakWarisan  string `json:"nomorKontrakWarisan,omitempty"`
	IDKontrakDisalinDari *int64 `json:"idKontrakDisalinDari,omitempty"`
	IDCedant             int64  `json:"idCedant"`
	IDAsalBisnis         int64  `json:"idAsalBisnis"`
	SifatProporsi        string `json:"sifatProporsi"`
	// Batas INKLUSIF keduanya (ADR-0022).
	TanggalMulai    string `json:"tanggalMulai"`
	TanggalBerakhir string `json:"tanggalBerakhir"`
}

// VersiKontrak adalah kepala satu versi - kolom wajib isinya saja.
//
// NomorUrutVersi penunjuk sebab migrasi 441 membuatnya boleh kosong: kosong
// berarti baris warisan yang belum dinomori ulang (tiket 10 Adjustment),
// bukan nol.
type VersiKontrak struct {
	ID                 int64        `json:"id"`
	IDKontrak          int64        `json:"idKontrak"`
	NomorUrutVersi     *int64       `json:"nomorUrutVersi"`
	KeadaanSiklusHidup string       `json:"keadaanSiklusHidup"`
	NamaKontrak        string       `json:"namaKontrak"`
	KodeMataUangKontak int64        `json:"kodeMataUangKontrak"`
	PersenBagianNure   *apd.Decimal `json:"-"`
	BagianNureSeragam  string       `json:"bagianNureSeragam"`
	MemakaiBordereaux  string       `json:"memakaiBordereaux"`
	CaraPembukuan      string       `json:"caraPembukuan"`
	MemakaiProrata     string       `json:"memakaiProrata"`
	RetroBerganda      string       `json:"retroBerganda"`
}

// KontrakDenganVersi adalah apa yang dibaca kembali dengan pengenalnya -
// lapisan beku SEKALI, versinya banyak.
type KontrakDenganVersi struct {
	Kontrak Kontrak        `json:"kontrak"`
	Versi   []VersiKontrak `json:"versi"`
}
