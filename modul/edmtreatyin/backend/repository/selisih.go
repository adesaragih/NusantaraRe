package repository

// Untuk apa berkas ini: TABEL PROYEKSI SELISIH (migrasi 360-363; `models/katalog_selisih.go`) - tulis dan baca
// `PolicyTreatyIn.TreatyDifference` dan `TreatyXOLDifferenceList` satu generasi endorsemen.
//
// ⛔ Tiga aturan proyeksi (spec-penyimpanan ID-26, AC 24-27): (1) hanya ditulis lapisan aplikasi, di transaksi yang
// SAMA dengan generasinya (pemanggil services.tulis); (2) bangun ulang hanya menyentuh baris SUMBER = 'GO' - baris
// 'PEGA' (hasil pemuat dokumen lama) tidak pernah dihapus `SimpanSelisih`; (3) sumber kebenaran tetap dua baris
// generasi. Rumus selisih nol di SQL (ID-23, AC 23): repository hanya menulis hasil `models`.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
)

// KunciSelisih - kunci saring yang disalin ke T_POLIS_DIFFERENCE (ID-27, AC 28).
type KunciSelisih struct {
	NoPolis string
	ProdKe  int
	EDMNo   string
	IDPega  string
	// Sumber - models.SumberGo (aplikasi) atau models.SumberPega (pemuat).
	Sumber string
}

// anak tabel selisih: kolom induk DIFFERENCE_ID.
var tabelAnakSelisih = []models.Tabel{models.TabelSelisihSpreading, models.TabelSelisihAngsuran, models.TabelSelisihLapisan}

// SimpanSelisih menulis ulang proyeksi selisih satu generasi (hapus lalu sisip) - HANYA bila baris induknya
// SUMBER = 'GO' atau belum ada; baris 'PEGA' tidak disentuh SIAPA PUN (juga pemuat: generasi yang sudah dimuat
// dilewati, bukan ditimpa) dan galatnya `ErrSelisihBeku` - SEBELUM satu baris anak pun dihapus.
func (g *Gudang) SimpanSelisih(ctx context.Context, tx *db.Tx, polisID string, h *models.Halaman, k KunciSelisih) error {
	t, err := g.nama(models.TabelSelisih.Nama)
	if err != nil {
		return err
	}
	lama, sumber, err := g.indukSelisih(ctx, tx, polisID)
	if err != nil {
		return err
	}
	if lama != "" {
		if sumber == models.SumberPega {
			return ErrSelisihBeku
		}
		if err := g.hapusSelisih(ctx, tx, lama); err != nil {
			return err
		}
	}
	id, err := idBaru()
	if err != nil {
		return err
	}
	kolom := []string{"ID", "POLIS_ID", "NOPOLIS", "PRODKE", "EDM_NO", "IDPEGA", "SUMBER"}
	nilai := []string{":1", ":2", ":3", ":4", ":5", ":6", ":7"}
	args := []any{id, polisID, db.KosongJadiNil(k.NoPolis), k.ProdKe, db.KosongJadiNil(k.EDMNo), db.KosongJadiNil(k.IDPega), k.Sumber}
	n := len(args) + 1
	for _, kol := range models.TabelSelisih.Kolom {
		eks, jml := ekspresiTulis(kol, n)
		v, err := nilaiTulis(kol, h.Ambil(kol.Properti))
		if err != nil {
			return err
		}
		kolom, nilai = append(kolom, kol.Kolom), append(nilai, eks)
		args = append(args, v...)
		n += jml
	}
	q := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, t, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
	if _, err := jalankan(ctx, tx, "menyimpan selisih", q, args...); err != nil {
		return err
	}
	baris := map[string][]models.Baris{
		models.TabelSelisihSpreading.Nama: h.AmbilDaftar(models.DaftarSelisihSpreading),
		models.TabelSelisihAngsuran.Nama:  h.AmbilDaftar(models.DaftarSelisihAngsuran),
		models.TabelSelisihLapisan.Nama:   models.DatarSelisihLapisan(h),
	}
	for _, a := range tabelAnakSelisih {
		if err := g.sisipAnakSelisih(ctx, tx, a, id, baris[a.Nama]); err != nil {
			return err
		}
	}
	return nil
}

// indukSelisih - ID dan SUMBER baris T_POLIS_DIFFERENCE generasi itu ("" bila belum ada).
func (g *Gudang) indukSelisih(ctx context.Context, tx *db.Tx, polisID string) (string, string, error) {
	t, err := g.nama(models.TabelSelisih.Nama)
	if err != nil {
		return "", "", err
	}
	q := fmt.Sprintf(`SELECT ID, SUMBER FROM %s WHERE POLIS_ID = :1`, t)
	if err := db.PeriksaSQL(q); err != nil {
		return "", "", err
	}
	var id, sumber sql.NullString
	err = g.pembaca(tx).QueryRowContext(ctx, q, polisID).Scan(&id, &sumber)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", nil
	}
	if err != nil {
		return "", "", fmt.Errorf("repository: membaca induk selisih: %w", err)
	}
	return teks(id), teks(sumber), nil
}

// BuangSelisihGenerasi membuang proyeksi selisih satu generasi yang DILEPAS dari rantai (admin menolak,
// `LepasGenerasi`) - tinjauan kode 06-10-2026: tanpa ini baris selisihnya tersisa berkunci saring (PRODKE, EDM_NO)
// sama dengan endorsemen berikutnya. Hanya baris 'GO'; baris 'PEGA' = `ErrSelisihBeku`.
func (g *Gudang) BuangSelisihGenerasi(ctx context.Context, tx *db.Tx, polisID string) error {
	lama, sumber, err := g.indukSelisih(ctx, tx, polisID)
	if err != nil || lama == "" {
		return err
	}
	if sumber != models.SumberGo {
		return ErrSelisihBeku
	}
	return g.hapusSelisih(ctx, tx, lama)
}

// hapusSelisih menghapus anak lalu induk satu baris T_POLIS_DIFFERENCE.
func (g *Gudang) hapusSelisih(ctx context.Context, tx *db.Tx, id string) error {
	for _, a := range tabelAnakSelisih {
		t, err := g.nama(a.Nama)
		if err != nil {
			return err
		}
		if _, err := jalankan(ctx, tx, "menghapus "+a.Nama, fmt.Sprintf(`DELETE FROM %s WHERE DIFFERENCE_ID = :1`, t), id); err != nil {
			return err
		}
	}
	t, err := g.nama(models.TabelSelisih.Nama)
	if err != nil {
		return err
	}
	_, err = jalankan(ctx, tx, "menghapus selisih", fmt.Sprintf(`DELETE FROM %s WHERE ID = :1 AND SUMBER = :2`, t), id, models.SumberGo)
	return err
}

func (g *Gudang) sisipAnakSelisih(ctx context.Context, tx *db.Tx, a models.Tabel, induk string, baris []models.Baris) error {
	if len(baris) == 0 {
		return nil
	}
	t, err := g.nama(a.Nama)
	if err != nil {
		return err
	}
	kolom := []string{"ID", "DIFFERENCE_ID", "NOURUT"}
	nilai := []string{":1", ":2", ":3"}
	n := 4
	for _, k := range a.Kolom {
		eks, jml := ekspresiTulis(k, n)
		kolom, nilai = append(kolom, k.Kolom), append(nilai, eks)
		n += jml
	}
	q := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, t, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
	for i, b := range baris {
		idBaris, err := idBaru()
		if err != nil {
			return err
		}
		args := []any{idBaris, induk, i + 1}
		for _, k := range a.Kolom {
			v, err := nilaiTulis(k, b[k.Properti])
			if err != nil {
				return fmt.Errorf("%s baris %d: %w", a.Nama, i+1, err)
			}
			args = append(args, v...)
		}
		if _, err := jalankan(ctx, tx, "menyimpan "+a.Nama, q, args...); err != nil {
			return err
		}
	}
	return nil
}

// BacaSelisih memuat proyeksi selisih satu generasi ke halaman `h` di bawah awalan `awalan` (jalur relatif
// `PolicyTreatyIn.`; "" = generasi ini, "OldData." = generasi sebelumnya untuk tab Old Data `PropOldData2`).
// Induk TreatyXOLDifferenceList dibangun ulang dari lapisannya (`BangunIndukSelisihXOL`).
func (g *Gudang) BacaSelisih(ctx context.Context, tx *db.Tx, polisID string, h *models.Halaman, awalan string) error {
	id, _, err := g.indukSelisih(ctx, tx, polisID)
	if err != nil || id == "" {
		return err
	}
	t, err := g.nama(models.TabelSelisih.Nama)
	if err != nil {
		return err
	}
	ks := models.TabelSelisih.Kolom
	q := fmt.Sprintf(`SELECT %s FROM %s WHERE ID = :1`, daftarBaca(ks, ""), t)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	nilai := make([]sql.NullString, len(ks))
	tujuan := make([]any, len(nilai))
	for i := range nilai {
		tujuan[i] = &nilai[i]
	}
	if err := g.pembaca(tx).QueryRowContext(ctx, q, id).Scan(tujuan...); err != nil {
		return fmt.Errorf("repository: membaca selisih: %w", err)
	}
	sementara := models.HalamanBaru()
	for i, k := range ks {
		if v := nilaiBaca(k, nilai[i]); v != "" {
			sementara.Setel(k.Properti, v)
		}
	}
	for _, a := range tabelAnakSelisih {
		c, err := g.nama(a.Nama)
		if err != nil {
			return err
		}
		qa := fmt.Sprintf(`SELECT ID, %s FROM %s WHERE DIFFERENCE_ID = :1 ORDER BY NOURUT`, daftarBaca(a.Kolom, ""), c)
		baris, _, err := g.bacaBaris(ctx, tx, qa, a.Kolom, id)
		if err != nil {
			return err
		}
		if len(baris) == 0 {
			continue
		}
		if a.Nama == models.TabelSelisihLapisan.Nama {
			if err := models.BangunIndukSelisihXOL(sementara, baris); err != nil {
				return err
			}
			continue
		}
		sementara.SetelDaftar(a.Daftar, baris)
	}
	models.PindahAwalan(sementara, h, models.HalamanPolis+".", models.HalamanPolis+"."+awalan)
	return nil
}
