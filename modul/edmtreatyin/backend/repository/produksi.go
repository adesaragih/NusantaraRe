package repository

// Untuk apa berkas ini: TABEL PRODUKSI LAMA sesudah Dept Head menyetujui endorsemen (Utility1
// `SaveJsonPolisTreatyInEDM_Act`, pemetaan di models/produksi.go; `[keputusan work owner 06-10-2026]` "JSON-nya
// tidak disimpan, tapi tetap insert kolom lainnya"). Pengganti (korpus `EDM Treaty In\RDBList\`):
//
//	TreatyInSearchProdKe        select count(*) as CARI1 from pooldata.json_polis where Nopolis = {InputData.CARI2}
//	SavePolisTreatyInEDM_SQL    -> POOLDATA.PEGA_JSON_POLIS_TREATYIN: INSERT json_polis (DATA_JSON TIDAK ditulis)
//	GetPolicyNoByCaseId         SELECT nopolis FROM json_polis WHERE idpega = {TempSearch.CARI1}
//	SaveAchievementSQL          -> POOLDATA.InsertUpdateAchievment: INSERT ACHIEVEMENT (..., TGL_PROD = sysdate)
//	GetDataTreatyInProd_SQL     select IDPEGA from treatyinproduction where IDPEGA = {TempSearch.CARI1}
//	GetReinstypeIDbyName_SQL    Select id as CARI1 from REINSURANCETYPE where note = {TempSearch.CARI10}
//	InsertTreatyInProdEDMT_SQL  INSERT INTO TREATYINPRODUCTION (58 kolom)
//
// Isi kedua prosedur dibaca tim NB dari ALL_SOURCE DEV (06-10-2026): masing-masing SATU INSERT; COUNT di
// InsertUpdateAchievment tidak dipakai (selalu menyisip). ⛔ TABEL WARISAN - tidak dibuat, tidak diubah
// strukturnya. ⛔ Nol prosedur, nol COMMIT: ditulis di transaksi services; skema eksplisit lewat `g.nama`
// (`Qualify`). Trigger warisan tetap jalan apa adanya.
//
// Bentuk: BACA dulu (`prBaca`, urutan RDB-List Pega), lalu SUSUN pernyataan tulis secara murni (`prSusunTulisan`,
// diuji tanpa Oracle), lalu JALANKAN berurutan.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
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

// prSQLCacahNoPolis = `TreatyInSearchProdKe` (SaveJsonPolisTreatyInEDM_Act 11): jumlah baris json_polis bernomor
// polis sama, dibaca SEBELUM sisip langkah 13 -> PRODKE baris baru.
func prSQLCacahNoPolis(t string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE NOPOLIS = :1`, t)
}

// prSQLIDJenisReas = `GetReinstypeIDbyName_SQL` (InsetTreatyInProdAddendum_Act 10.3.2) lalu `pxResults(1)`
// (10.3.3). Pega tanpa ORDER BY: bila NOTE kembar, baris pertama tak tentu - di sini ID terkecil supaya tetap.
func prSQLIDJenisReas(t string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(ID) FROM %s WHERE NOTE = :1 ORDER BY ID FETCH FIRST 1 ROWS ONLY`, t)
}

func sqlSisipJSONPolis(t string) string {
	return sqlSisipKolom(t, models.KolomJSONPolis, kolomSQL{"TGL_INPUT", "SYSDATE"})
}

func sqlSisipCapaian(t string) string {
	return sqlSisipKolom(t, models.KolomCapaian, kolomSQL{"TGL_PROD", "SYSDATE"})
}

func sqlSisipProduksi(t string) string { return sqlSisipKolom(t, models.KolomProduksi) }

// prTabelProduksi - nama berskema ketiga tabel tulis.
type prTabelProduksi struct{ jsonPolis, capaian, produksi string }

// prBacaanProduksi - hasil RDB-List baca yang menentukan tulisan Utility1.
type prBacaanProduksi struct {
	// adaJSON, adaProduksi - baris ber-IDPEGA yang sudah ada (kunci utama json_polis; GetDataTreatyInProd_SQL).
	adaJSON, adaProduksi int
	// prodKe - TreatyInSearchProdKe; dibaca hanya bila json_polis akan disisip.
	prodKe int
	// noPolisJSON - GetPolicyNoByCaseId atas baris yang SUDAH ada (json_polis tidak disisip ulang).
	noPolisJSON string
	// idJenisReas - GetReinstypeIDbyName_SQL; dibaca hanya bila produksi ditulis dan models meminta.
	idJenisReas string
}

// prTulisan - satu pernyataan tulis.
type prTulisan struct {
	apa, q string
	arg    []any
}

func prSalinBaris(b models.Baris) models.Baris {
	out := models.Baris{}
	for k, v := range b {
		out[k] = v
	}
	return out
}

// prSusunTulisan menyusun pernyataan tulis urutan Pega: json_polis (SaveJson 13), ACHIEVEMENT (15),
// TREATYINPRODUCTION (19).
func prSusunTulisan(t prTabelProduksi, s models.SimpananPolis, b prBacaanProduksi) ([]prTulisan, error) {
	var out []prTulisan
	// 13 - kunci utama IDPEGA: baris yang sudah ada membuat prosedur Pega gagal dan menelan galatnya (14
	// ERRMSG); di sini dilewati dengan sengaja, bukan digagalkan.
	noPolis := b.noPolisJSON
	if b.adaJSON == 0 {
		jp := prSalinBaris(s.JSONPolis)
		jp["PRODKE"] = strconv.Itoa(b.prodKe) // 12 CARI5 = ProdKe.pxResults(1).CARI1
		arg, err := argumenKolom(models.KolomJSONPolis, jp)
		if err != nil {
			return nil, err
		}
		out = append(out, prTulisan{"menulis json_polis", sqlSisipJSONPolis(t.jsonPolis), arg})
		noPolis = jp["NOPOLIS"] // GetPolicyNoByCaseId sesudah sisip = baris ini
	}

	// 15 - SetEDMAchivementValue 2: NOPOLIS = GetPolicyNoByCaseId; tanpa penjaga (selalu menyisip)
	for _, c := range s.Capaian {
		salin := prSalinBaris(c)
		salin["NOPOLIS"] = noPolis
		arg, err := argumenKolom(models.KolomCapaian, salin)
		if err != nil {
			return nil, err
		}
		out = append(out, prTulisan{"menulis ACHIEVEMENT", sqlSisipCapaian(t.capaian), arg})
	}

	// 19 - InsetTreatyInProdAddendum_Act 6-8: IDPEGA yang sudah punya baris keluar tanpa menulis
	if b.adaProduksi > 0 {
		return out, nil
	}
	s.TerapkanJenisReasXOL(b.idJenisReas) // 10.3.3 (hanya bila models meminta pencarian)
	for _, p := range s.Produksi {
		arg, err := argumenKolom(models.KolomProduksi, p)
		if err != nil {
			return nil, err
		}
		out = append(out, prTulisan{"menulis TREATYINPRODUCTION", sqlSisipProduksi(t.produksi), arg})
	}
	return out, nil
}

// SimpanPolisProduksi menulis baris Utility1 satu kasus EDM di dalam `tx`, urutan Pega: json_polis, ACHIEVEMENT,
// TREATYINPRODUCTION.
func (g *Gudang) SimpanPolisProduksi(ctx context.Context, tx *db.Tx, s models.SimpananPolis) error {
	t, err := g.prNamaTabel()
	if err != nil {
		return err
	}
	b, err := g.prBaca(ctx, tx, t, s)
	if err != nil {
		return err
	}
	tulisan, err := prSusunTulisan(t, s, b)
	if err != nil {
		return err
	}
	for _, w := range tulisan {
		if _, err := jalankan(ctx, tx, w.apa, w.q, w.arg...); err != nil {
			return err
		}
	}
	return nil
}

func (g *Gudang) prNamaTabel() (prTabelProduksi, error) {
	var t prTabelProduksi
	var err error
	if t.jsonPolis, err = g.nama(models.TabelJSONPolis); err != nil {
		return t, err
	}
	if t.capaian, err = g.nama(models.TabelCapaian); err != nil {
		return t, err
	}
	t.produksi, err = g.nama(models.TabelProduksi)
	return t, err
}

// prBaca menjalankan RDB-List baca Utility1 di dalam `tx`.
func (g *Gudang) prBaca(ctx context.Context, tx *db.Tx, t prTabelProduksi, s models.SimpananPolis) (prBacaanProduksi, error) {
	var b prBacaanProduksi
	var err error
	if b.adaJSON, err = g.cacahIDPega(ctx, tx, t.jsonPolis, s.IDPega); err != nil {
		return b, err
	}
	if b.adaJSON == 0 { // 11 TreatyInSearchProdKe (InputData.CARI2 = PolicyTreatyIn.PolicyNo)
		if b.prodKe, err = g.prCacah(ctx, tx, prSQLCacahNoPolis(t.jsonPolis), s.JSONPolis["NOPOLIS"]); err != nil {
			return b, err
		}
	} else if len(s.Capaian) > 0 { // SetEDMAchivementValue 2 GetPolicyNoByCaseId
		if b.noPolisJSON, err = g.prTeksSatu(ctx, tx, sqlNoPolisJSON(t.jsonPolis), s.IDPega); err != nil {
			return b, err
		}
	}
	if b.adaProduksi, err = g.cacahIDPega(ctx, tx, t.produksi, s.IDPega); err != nil {
		return b, err
	}
	if b.adaProduksi == 0 && s.NotaJenisReasXOL != "" && len(s.Produksi) > 0 { // 10.3.1-10.3.2
		tj, err := g.nama(tabelJenisReas)
		if err != nil {
			return b, err
		}
		if b.idJenisReas, err = g.prTeksSatu(ctx, tx, prSQLIDJenisReas(tj), s.NotaJenisReasXOL); err != nil {
			return b, err
		}
	}
	return b, nil
}

// prCacah - satu COUNT(*) berparameter tunggal.
func (g *Gudang) prCacah(ctx context.Context, tx *db.Tx, q, arg string) (int, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int
	if err := g.pembaca(tx).QueryRowContext(ctx, q, arg).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: mencacah produksi lama: %w", err)
	}
	return n, nil
}

// prTeksSatu - satu nilai teks baris pertama; tanpa baris = "" (Pega `pxResults(1).X` atas daftar kosong).
func (g *Gudang) prTeksSatu(ctx context.Context, tx *db.Tx, q, arg string) (string, error) {
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var v sql.NullString
	err := g.pembaca(tx).QueryRowContext(ctx, q, arg).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca acuan produksi lama: %w", err)
	}
	return v.String, nil
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
