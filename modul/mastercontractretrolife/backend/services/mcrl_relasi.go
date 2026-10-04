package services

// Relasi lima tabel _LIFE di basis data (migrasi 100, keputusan work owner
// 01-10-2026): PK, FK ke induk langsung ON DELETE CASCADE, dan FK penjaga
// salinan DEFERRABLE INITIALLY DEFERRED.
//
// Untuk apa berkas ini: penolakan Oracle atas relasi itu - termasuk yang baru
// muncul saat COMMIT karena constraint-nya ditunda - sampai ke layar sebagai
// kalimat (K8), bukan 500. Aturan layanan (salinan dari induk, kaskade di Go,
// kunci induk) tetap lapis pertama; basis data adalah lapis kedua yang menolak
// apa pun yang lolos, termasuk tulisan dari luar aplikasi.

import (
	"context"
	"errors"
	"strings"

	"nusantarare/inti/backend/db"
)

// ErrRelasiDitolak - basis data menolak perubahan karena melanggar relasi
// tahun, kontrak, reinsurer, security, dan business.
var ErrRelasiDitolak = errors.New("services: the change breaks the relation between treaty year, contract, " +
	"reinsurer, security reinsurer, and business")

// kodeRelasi - galat Oracle yang berarti relasi atau kunci dilanggar.
//
//	ORA-02291  induk tidak ada, atau kolom salinan tidak sama dengan induknya
//	           (FK; FK penjaga salinan yang ditunda melaporkannya saat COMMIT
//	           bersama ORA-02091)
//	ORA-02292  induk masih punya anak (FK tanpa kaskade)
//	ORA-00001  kunci ganda (PK atau kunci unik)
var kodeRelasi = []struct{ kode, pesan string }{
	{"ORA-02291", "The parent row no longer exists, or a copied column no longer matches its parent. " +
		"Reload the screen and try again."},
	{"ORA-02292", "This row still has child rows. Delete it from the screen so its child rows are deleted with it."},
	{"ORA-00001", "A row with the same key already exists. Reload the screen and try again."},
}

// galatRelasi membawa kalimat layar dan galat aslinya (untuk log server).
type galatRelasi struct {
	pesan string
	asal  error
}

func (g galatRelasi) Error() string      { return ErrRelasiDitolak.Error() + ": " + g.asal.Error() }
func (g galatRelasi) PesanLayar() string { return g.pesan }
func (g galatRelasi) Unwrap() []error    { return []error{ErrRelasiDitolak, g.asal} }

// TerjemahRelasi mengubah penolakan relasi Oracle menjadi ErrRelasiDitolak
// berkalimat; galat lain dikembalikan apa adanya.
func TerjemahRelasi(err error) error {
	if err == nil || errors.Is(err, ErrRelasiDitolak) {
		return err
	}
	teks := err.Error()
	for _, k := range kodeRelasi {
		if strings.Contains(teks, k.kode) {
			return galatRelasi{pesan: k.pesan, asal: err}
		}
	}
	return err
}

// transaksiBerelasi membungkus Transaksi: galat dari pernyataan mana pun DAN
// dari COMMIT diterjemahkan di satu tempat.
func transaksiBerelasi(tx Transaksi) Transaksi {
	if tx == nil {
		return nil
	}
	return func(ctx context.Context, fn func(tx *db.Tx) error) error {
		return TerjemahRelasi(tx(ctx, fn))
	}
}
