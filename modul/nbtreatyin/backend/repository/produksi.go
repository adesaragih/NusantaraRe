package repository

// Untuk apa berkas ini: TABEL PRODUKSI LAMA sesudah realisasi selesai (Utility1 `SaveJsonPolisTreatyIn_Act`,
// pemetaan di models/produksi.go; `[keputusan work owner 06-10-2026]` "JSON-nya tidak disimpan, tapi tetap insert
// kolom lainnya"). Pengganti:
//
//	SavePolisTreatyIn_SQL      -> POOLDATA.PEGA_JSON_POLIS_TREATYIN: INSERT json_polis (DATA_JSON TIDAK ditulis)
//	GetPolicyNoByCaseId        SELECT nopolis FROM json_polis WHERE idpega = {TempSearch.CARI1}
//	SaveAchievementSQL         -> POOLDATA.InsertUpdateAchievment: INSERT ACHIEVEMENT (..., TGL_PROD = sysdate)
//	GetDataTreatyInProd_SQL    select IDPEGA from treatyinproduction where IDPEGA = {TempSearch.CARI1}
//	InsertTreatyInProd_SQL     INSERT INTO TREATYINPRODUCTION (58 kolom)
//
// Isi kedua prosedur dibaca dari ALL_SOURCE DEV (06-10-2026): masing-masing SATU INSERT; COUNT di
// InsertUpdateAchievment tidak dipakai (selalu menyisip). ⛔ TABEL WARISAN - tidak dibuat, tidak diubah
// strukturnya (MODUL.md). ⛔ Nol prosedur, nol COMMIT: ditulis di transaksi submit (spec-penyimpanan AC 48; AC
// 29, 83); skema eksplisit lewat `Qualify` (AC 47). Trigger warisan tetap jalan apa adanya
// (`UPDATE_JSON_POLIS_LOG` -> JSONPOLISLOG, `TRG_TREATYINPRODUCTION_INSERT` -> TREATYINPRODUCTION_BACKUP).

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// kolomSQL - kolom yang nilainya ekspresi SQL tetap (bukan penampung).
type kolomSQL struct{ kolom, ekspresi string }

// sqlSisipKolom merakit INSERT satu baris: kolom katalog berpenampung (`ekspresiTulis`), lalu kolom tetap.
func sqlSisipKolom(t string, ks []models.Kolom, tetap ...kolomSQL) string {
	var kol, nilai []string
	n := 1
	for _, k := range ks {
		e, pakai := ekspresiTulis(k, n)
		kol, nilai = append(kol, k.Kolom), append(nilai, e)
		n += pakai
	}
	for _, x := range tetap {
		kol, nilai = append(kol, x.kolom), append(nilai, x.ekspresi)
	}
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", t, strings.Join(kol, ", "), strings.Join(nilai, ", "))
}

// argumenKolom - argumen bind satu baris menurut urutan katalog.
func argumenKolom(ks []models.Kolom, b models.Baris) ([]any, error) {
	var out []any
	for _, k := range ks {
		v, err := nilaiTulis(k, b[k.Properti])
		if err != nil {
			return nil, err
		}
		out = append(out, v...)
	}
	return out, nil
}

func sqlCacahIDPega(t string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE IDPEGA = :1`, t)
}

func sqlNoPolisJSON(t string) string {
	return fmt.Sprintf(`SELECT NOPOLIS FROM %s WHERE IDPEGA = :1`, t)
}

func sqlSisipJSONPolis(t string) string {
	return sqlSisipKolom(t, models.KolomJSONPolis, kolomSQL{"TGL_INPUT", "SYSDATE"})
}

func sqlSisipCapaian(t string) string {
	return sqlSisipKolom(t, models.KolomCapaian, kolomSQL{"TGL_PROD", "SYSDATE"})
}

func sqlSisipProduksi(t string) string { return sqlSisipKolom(t, models.KolomProduksi) }

// SimpanPolisProduksi menulis baris Utility1 satu kasus, urutan Pega: json_polis (langkah 11), ACHIEVEMENT
// (langkah 13), TREATYINPRODUCTION (langkah 17).
func (g *Gudang) SimpanPolisProduksi(ctx context.Context, tx *db.Tx, s models.SimpananPolis) error {
	jp, err := g.nama(models.TabelJSONPolis)
	if err != nil {
		return err
	}
	capai, err := g.nama(models.TabelCapaian)
	if err != nil {
		return err
	}
	prod, err := g.nama(models.TabelProduksi)
	if err != nil {
		return err
	}

	// 11 - kunci utama IDPEGA: baris yang sudah ada membuat prosedur Pega gagal dan menelan galatnya; di sini
	// dilewati dengan sengaja, bukan digagalkan.
	ada, err := g.cacahIDPega(ctx, tx, jp, s.IDPega)
	if err != nil {
		return err
	}
	if ada == 0 {
		arg, err := argumenKolom(models.KolomJSONPolis, s.JSONPolis)
		if err != nil {
			return err
		}
		if _, err := jalankan(ctx, tx, "menulis json_polis", sqlSisipJSONPolis(jp), arg...); err != nil {
			return err
		}
	}

	// 13 - SetAchivementValue: NOPOLIS = GetPolicyNoByCaseId
	if len(s.Capaian) > 0 {
		q := sqlNoPolisJSON(jp)
		if err := db.PeriksaSQL(q); err != nil {
			return err
		}
		var nopol sql.NullString
		if err := tx.QueryRowContext(ctx, q, s.IDPega).Scan(&nopol); err != nil {
			return fmt.Errorf("repository: membaca nomor polis json_polis: %w", err)
		}
		qc := sqlSisipCapaian(capai)
		for _, b := range s.Capaian {
			salin := models.Baris{}
			for k, v := range b {
				salin[k] = v
			}
			salin["NOPOLIS"] = nopol.String
			arg, err := argumenKolom(models.KolomCapaian, salin)
			if err != nil {
				return err
			}
			if _, err := jalankan(ctx, tx, "menulis ACHIEVEMENT", qc, arg...); err != nil {
				return err
			}
		}
	}

	// 17 - InsetTreatyInProd_Act 6-8: IDPEGA yang sudah punya baris keluar tanpa menulis
	adaProd, err := g.cacahIDPega(ctx, tx, prod, s.IDPega)
	if err != nil {
		return err
	}
	if adaProd > 0 {
		return nil
	}
	qp := sqlSisipProduksi(prod)
	for _, b := range s.Produksi {
		arg, err := argumenKolom(models.KolomProduksi, b)
		if err != nil {
			return err
		}
		if _, err := jalankan(ctx, tx, "menulis TREATYINPRODUCTION", qp, arg...); err != nil {
			return err
		}
	}
	return nil
}

// CacahProduksi - jumlah baris json_polis / ACHIEVEMENT / TREATYINPRODUCTION satu IDPEGA (alat `simpanproduksi`:
// laporan uji-kering dan penjaga ganda ACHIEVEMENT). Di dalam transaksi bila `tx` terisi.
func (g *Gudang) CacahProduksi(ctx context.Context, tx *db.Tx, idPega string) (models.CacahProduksi, error) {
	var c models.CacahProduksi
	for _, x := range []struct {
		tabel string
		ke    *int
	}{{models.TabelJSONPolis, &c.JSONPolis}, {models.TabelCapaian, &c.Capaian}, {models.TabelProduksi, &c.Produksi}} {
		t, err := g.nama(x.tabel)
		if err != nil {
			return c, err
		}
		if *x.ke, err = g.cacahIDPega(ctx, tx, t, idPega); err != nil {
			return c, err
		}
	}
	return c, nil
}

func (g *Gudang) cacahIDPega(ctx context.Context, tx *db.Tx, t, idPega string) (int, error) {
	q := sqlCacahIDPega(t)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int
	if err := g.pembaca(tx).QueryRowContext(ctx, q, idPega).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: mencacah baris %s: %w", t, err)
	}
	return n, nil
}
