package services

// Penjaga jalur tulis - perbaikan /code-review 01-10-2026: kunci induk (K2 tanpa FK), jejak audit yang
// menyebut pelaku (tiket 08 AC 29, tiket 09), lebar kolom `[data DBA]`, dan kalimat dampak yang terbaca.

import (
	"context"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/mastercontractretrolife/backend/models"
)

// Lebar kolom VARCHAR2 dalam BYTE (`ddl-tables-from-dba.md` `[data DBA]`).
const (
	// lebarKode - `ID`, `TREATYYEAR`, `UNDERWRITINGYEAR`, `RIRATEID`, dan `USERID` tahun/kontrak/business.
	lebarKode = 100
	// lebarTeks - `RIRATE`, nama, dan `USERID` reinsurer/security.
	lebarTeks = 1000
)

// muat - teks harus muat di kolomnya; lebih = 422 menyebut medannya (bukan ORA-12899 → 500, K8).
func muat(label, v string, lebar int) error {
	if len(v) > lebar {
		return fmt.Errorf("%w: %s is longer than %d bytes (%d)", ErrMasukanTidakSah, label, lebar, len(v))
	}
	return nil
}

// pelakuMuat - akun pelaku ditulis ke `USERID`; akun yang tidak muat ditolak sebelum menulis.
func pelakuMuat(p inti.Pelaku, lebar int) error {
	return muat("user account", p.AkunID, lebar)
}

// oleh - akhiran jejak audit: siapa pelakunya (kapan = stempel log server).
func oleh(p inti.Pelaku) string {
	return " by account " + p.AkunID
}

// teksDampak - cacah baris ikut terhapus dalam kalimat, untuk pesan layar dan log.
func teksDampak(d models.Dampak) string {
	return fmt.Sprintf("%d reinsurer, %d security, %d business rows", d.Reinsurer, d.Security, d.Business)
}

// galatTidakAda - galat 404 untuk satu jenis baris; nil = jenis tidak dikenal.
func galatTidakAda(jenis JenisHapus) error {
	switch jenis {
	case HapusKontrak:
		return ErrKontrakTidakAda
	case HapusReinsurer:
		return ErrReinsurerTidakAda
	case HapusSecurity:
		return ErrSecurityTidakAda
	case HapusBusiness:
		return ErrBusinessTidakAda
	}
	return nil
}

// kunci mengunci baris `jenis` ber-ID id sampai transaksi selesai; tidak ada = 404 jenis itu.
func (l *Layanan) kunci(ctx context.Context, tx *db.Tx, jenis JenisHapus, id string) error {
	jadi := galatTidakAda(jenis)
	if jadi == nil {
		return fmt.Errorf("%w: %q", ErrJenisHapusTidakAda, jenis)
	}
	return tidakAda(l.gudang.KunciBaris(ctx, tx, string(jenis), id), jadi, id)
}
