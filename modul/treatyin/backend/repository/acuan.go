// Package repository memegang seluruh sentuhan basis data modul Treaty In.
//
// Arah ketergantungan: handlers -> services -> repository. Paket ini tidak
// pernah mengimpor handlers maupun services, dan tidak pernah mengimpor modul
// lain (kontrak lintas modul lewat `inti/backend/kontrak`, bab 7 panduan).
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
)

// tabelAcuan memetakan himpunan ke nama tabelnya.
//
// ⛔ Peta TERTUTUP, dan itu yang membuat nama tabel tidak pernah datang dari
// teks permintaan. Himpunan di luar peta ini ditolak sebelum satu pun kueri
// disusun - bukan diteruskan ke Oracle untuk ditolak di sana.
//
// Nama tabelnya adalah nama spec apa adanya (`KEPUTUSAN-PENYELARASAN-REPO.md`
// butir 3), bukan berawalan `T_`.
var tabelAcuan = map[models.Himpunan]string{
	models.HimpunanJenisPotongan:   "JENIS_POTONGAN",
	models.HimpunanKelasBisnis:     "KELAS_BISNIS",
	models.HimpunanKelompokTreaty:  "KELOMPOK_TREATY",
	models.HimpunanBahaya:          "BAHAYA",
	models.HimpunanJenisReasuransi: "JENIS_REASURANSI",
}

// kolomKunci memetakan himpunan ke nama kolom kunci utamanya. Kelimanya
// berbeda (`ID_MATA_UANG`, `ID_BAHAYA`, ...) - `KAMUS-KOLOM.md` §10.22.
var kolomKunci = map[models.Himpunan]string{
	models.HimpunanJenisPotongan:   "ID_JENIS_POTONGAN",
	models.HimpunanKelasBisnis:     "ID_KELAS_BISNIS",
	models.HimpunanKelompokTreaty:  "ID_KELOMPOK_TREATY",
	models.HimpunanBahaya:          "ID_BAHAYA",
	models.HimpunanJenisReasuransi: "ID_JENIS_REASURANSI",
}

// ErrHimpunanTidakDikenal - himpunan di luar keenam tabel acuan.
var ErrHimpunanTidakDikenal = fmt.Errorf("repository: himpunan acuan tidak dikenal")

// Gudang adalah seluruh sentuhan basis data modul ini.
type Gudang struct{ db *db.DB }

// Baru membuat Gudang di atas koneksi.
func Baru(d *db.DB) *Gudang { return &Gudang{db: d} }

// DaftarAcuan membaca seluruh baris satu tabel acuan, berurut menurut KODE.
//
// Urutannya KODE, bukan pengenal: pengenal datang dari sequence dan karena itu
// menyusun baris menurut urutan MASUKNYA, yang bukan urutan yang berarti bagi
// pembaca. KODE adalah kunci alaminya (INV-68).
func (g *Gudang) DaftarAcuan(ctx context.Context, h models.Himpunan) ([]models.Acuan, error) {
	tabel, ada := tabelAcuan[h]
	if !ada {
		return nil, fmt.Errorf("%w: %q", ErrHimpunanTidakDikenal, h)
	}
	nama, err := g.db.Qualify(tabel)
	if err != nil {
		return nil, err
	}
	induk := "NULL"
	if h.Bersusun() {
		induk = "ID_INDUK"
	}
	q := fmt.Sprintf("SELECT %s, KODE, NAMA, AKTIF, %s FROM %s ORDER BY KODE ASC", kolomKunci[h], induk, nama)
	baris, err := g.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", tabel, err)
	}
	defer func() { _ = baris.Close() }()

	// Daftar kosong adalah JAWABAN, bukan kegagalan: keenam tabel berdiri
	// kosong sampai tiket 44 memindahkan isinya dari sistem lama.
	keluar := []models.Acuan{}
	for baris.Next() {
		var a models.Acuan
		var idInduk sql.NullInt64
		if err := baris.Scan(&a.ID, &a.Kode, &a.Nama, &a.Aktif, &idInduk); err != nil {
			return nil, fmt.Errorf("repository: membaca baris %s: %w", tabel, err)
		}
		if idInduk.Valid {
			n := idInduk.Int64
			a.IDInduk = &n
		}
		keluar = append(keluar, a)
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca %s: %w", tabel, err)
	}
	return keluar, nil
}
