package models

// Untuk apa berkas ini: COPY OLD - perintah work owner 07-10-2026 ("nb ttreatyin tobol copy untuk data lama mana?" ->
// "langsung anda kerjakan!", dikerjakan sesi EDM TREATY IN dengan izin WO): tombol di samping Create membuka popup
// berisi dokumen polis NB Treaty In lama di JSON_POLIS (generasi 0) yang BELUM ada di tabel flat; yang dicentang
// disalin lewat `Process Copy` dengan aturan pemuat dokumen lama (`services/pemuat.go` `muat`, tiket 22), satu
// transaksi per dokumen. Pola sama dengan Copy Old EDM Treaty In (`modul/edmtreatyin`), Product Name Life, dan gerbang
// superadmin Bordereaux. JSON_POLIS hanya DIBACA.

import "errors"

// DokumenLama - satu baris popup Copy Old. `Alasan` (kenapa tidak dapat disalin) menyebut sebab dan jalur medan,
// tidak pernah nilai dokumen.
type DokumenLama struct {
	// ID - nomor kasus `NB-<n>` dari IDPEGA (kosong bila IDPEGA tak terbaca).
	ID           string   `json:"id"`
	NoOffer      string   `json:"noOffer"`
	NoPolis      string   `json:"noPolis"`
	InsuredName  string   `json:"insuredName"`
	BusinessName string   `json:"businessName"`
	SOBName      string   `json:"sobName"`
	CedingCoName string   `json:"cedingCoName"`
	TglProd      string   `json:"tglProd"`
	BolehDisalin bool     `json:"bolehDisalin"`
	Alasan       []string `json:"alasan"`
}

// Status satu dokumen sesudah `Process Copy` (nilai sama dengan Copy Old EDM Treaty In / Product Name Life).
const (
	// SalinDisalin - kasus, generasi, tabel T_POLIS_*, dan SuggestList tertulis di tabel flat.
	SalinDisalin = "disalin"
	// SalinSudahAda - kasus sudah ada di tabel flat (disalin sebelumnya, oleh pemakai lain, atau alat pemuat).
	SalinSudahAda = "sudahAda"
	// SalinDitolak - penjaga pemuat menolak dokumen ini (`Pesan` = alasannya); tidak ada yang ditulis.
	SalinDitolak = "ditolak"
	// SalinGagal - galat basis data; transaksi dokumen ini dibatalkan, dokumen lain tidak terpengaruh.
	SalinGagal = "gagal"
)

// HasilSalinLama - hasil satu ID yang diminta.
type HasilSalinLama struct {
	ID     string   `json:"id"`
	Status string   `json:"status"`
	Pesan  []string `json:"pesan"`
}

// JawabanSalinLama - jawaban `Process Copy`: hasil per ID, urutan daftar popup.
type JawabanSalinLama struct {
	Hasil   []HasilSalinLama `json:"hasil"`
	Disalin int              `json:"disalin"`
}

// alasanSalin - teks layar (bahasa Inggris, pola Copy Old EDM) galat pemecah. Teks galat aslinya dapat memuat nilai
// dokumen - tidak pernah dikirim ke layar.
var alasanSalin = []struct {
	err  error
	teks string
}{
	{ErrDokumenGanda, "more than one old document for this NB number"},
	{ErrDokumenRusak, "the old JSON cannot be read"},
	{ErrIDPega, "the Pega case ID is not an NB case"},
	{ErrProdKe, "generation number (PRODKE) is empty or not a number"},
	{ErrNoPolisKosong, "policy number is empty"},
	{ErrNoPolisBeda, "policy number does not match the document"},
	{ErrNilaiKolom, "a value does not fit its column"},
	{ErrTanggalAmbigu, "a date is ambiguous (day/month order cannot be determined)"},
	{ErrFormatTanggal, "a date has an unknown format"},
}

// AlasanSalinLama - teks layar satu galat pemecah; "" = bukan galat dokumen.
func AlasanSalinLama(err error) string {
	for _, a := range alasanSalin {
		if errors.Is(err, a.err) {
			return a.teks
		}
	}
	return ""
}
