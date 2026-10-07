package models

// Untuk apa berkas ini: COPY OLD - perintah work owner 07-10-2026 ("BUATKAN TOMBOL COPY OLD SAMA SEPERTI MASTER
// PRODUCTNAME LIFE, KHUSUS BUAT SUPERUSER"): tombol di samping Create membuka popup berisi dokumen endorsemen Treaty In
// lama di JSON_POLIS yang BELUM ada di tabel flat; yang dicentang disalin lewat `Process Copy` dengan aturan pemuat
// dokumen lama (`services/pemuat.go` `muat`), satu transaksi per dokumen. Pola: Copy Old
// `modul/masterproductnamelife` (03-10-2026) + gerbang superadmin `modul/bordereaux` (04-10-2026). JSON_POLIS hanya
// DIBACA.

import "errors"

// DokumenLama - satu baris popup Copy Old. `Alasan` (kenapa tidak dapat disalin) menyebut sebab dan jalur medan,
// tidak pernah nilai dokumen.
type DokumenLama struct {
	// ID - nomor kasus `EDMT-<n>` dari IDPEGA (kosong bila IDPEGA tak terbaca).
	ID           string   `json:"id"`
	NoPolis      string   `json:"noPolis"`
	EDMNo        string   `json:"edmNo"`
	ProdKe       int      `json:"prodKe"`
	EDMType      string   `json:"edmType"`
	SOBName      string   `json:"sobName"`
	CedingCoName string   `json:"cedingCoName"`
	TglProd      string   `json:"tglProd"`
	BolehDisalin bool     `json:"bolehDisalin"`
	Alasan       []string `json:"alasan"`
}

// Status satu dokumen sesudah `Process Copy` (nilai sama dengan Copy Old Product Name Life).
const (
	// SalinDisalin - generasi, proyeksi selisih 'PEGA', penanda, dan SuggestList tertulis di tabel flat.
	SalinDisalin = "disalin"
	// SalinSudahAda - generasi sudah ada di tabel flat (disalin sebelumnya, oleh pemakai lain, atau alat pemuat).
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

// JawabanSalinLama - jawaban `Process Copy`: hasil per ID, urutan generasi (NOPOLIS lalu PRODKE).
type JawabanSalinLama struct {
	Hasil   []HasilSalinLama `json:"hasil"`
	Disalin int              `json:"disalin"`
}

// alasanSalin - teks layar (bahasa Inggris, pola Copy Old Product Name Life) galat pemecah / pemuat. Teks galat
// aslinya memuat nomor polis dan nilai dokumen - tidak pernah dikirim ke layar.
var alasanSalin = []struct {
	err  error
	teks string
}{
	{ErrGenerasiSebelumnyaTidakAda, "the previous generation is not in the new tables yet - copy it first (generation 1 needs the NB policy)"},
	{ErrPercabangan, "the previous generation already has an endorsement in the new tables"},
	{ErrKeutuhan, "spreading rows of the previous generation are missing"},
	{ErrOldDataTakSesuai, "Old Data EDM number does not match the previous generation"},
	{ErrIDKasusBentrok, "this EDM number is already used by another policy in the new tables"},
	{ErrDokumenGanda, "more than one old document for this EDM number"},
	{ErrDokumenRusak, "the old JSON cannot be read"},
	{ErrIDPega, "the Pega case ID is not an EDMT case"},
	{ErrProdKe, "generation number (PRODKE) is empty or not a number"},
	{ErrProdKeBeda, "generation number (PRODKE) does not match the document"},
	{ErrNoPolisKosong, "policy number is empty"},
	{ErrNoPolisBeda, "policy number does not match the document"},
	{ErrEDMNoKosong, "EDM number is empty"},
	{ErrEDMNoBeda, "EDM number does not match the document"},
	{ErrNilaiKolom, "a value does not fit its column"},
}

// AlasanSalinLama - teks layar satu galat pemecah / pemuat; "" = bukan galat dokumen (galat basis data).
func AlasanSalinLama(err error) string {
	for _, a := range alasanSalin {
		if errors.Is(err, a.err) {
			return a.teks
		}
	}
	return ""
}
