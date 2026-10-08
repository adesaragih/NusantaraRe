// Package models memegang bentuk data modul Treaty In (`treatyin`).
package models

// Himpunan menyebut salah satu dari enam tabel acuan tiket 15.
//
// ⛔ Keenamnya DISEBUT satu per satu, bukan diturunkan dari teks permintaan.
// ADR-0038 menyimpan ISI himpunan sebagai data; nama TABELNYA tetap tertutup,
// dan membiarkannya terbuka berarti nama tabel datang dari luar.
type Himpunan string

// LIMA himpunan acuan - `KAMUS-KOLOM.md` §10.22, tiket 15.
//
// ⛔ `mata-uang` DICABUT 4 Oktober 2026 (migrasi 434). Daftar mata uang dan
// kursnya dibaca dari `TREATYEXCHANGEYEARLY`, tabel warisan yang sudah hidup -
// keputusan pemilik proses. Tabel acuan `MATA_UANG` yang pernah berdiri di
// sini nol baris sepanjang hidupnya.
const (
	HimpunanJenisPotongan   Himpunan = "jenis-potongan"
	HimpunanKelasBisnis     Himpunan = "kelas-bisnis"
	HimpunanKelompokTreaty  Himpunan = "kelompok-treaty"
	HimpunanBahaya          Himpunan = "bahaya"
	HimpunanJenisReasuransi Himpunan = "jenis-reasuransi"
)

// Acuan adalah satu baris tabel acuan: kode, nama, dan penanda aktifnya.
//
// IDInduk hanya terisi pada `JENIS_REASURANSI`, yang bersusun (§10.6). Pada
// keempat himpunan lain ia SELALU kosong - dan kosong di sana berarti "tabel ini
// memang tidak bersusun", bukan "induknya belum diisi".
type Acuan struct {
	ID      int64  `json:"id"`
	Kode    string `json:"kode"`
	Nama    string `json:"nama"`
	Aktif   string `json:"aktif"`
	IDInduk *int64 `json:"idInduk,omitempty"`
}

// Bersusun menyatakan apakah himpunan ini membawa kolom `ID_INDUK`.
func (h Himpunan) Bersusun() bool { return h == HimpunanJenisReasuransi }
