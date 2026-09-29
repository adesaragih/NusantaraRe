package repository

// Pencarian diagnosa atas `DISEASE_LIFE` — kelompok Medis.
//
// ⛔ 97.586 baris. Setiap query di berkas ini BERBATAS, dan batasnya diambil
// dari rule Pega (`models.BatasBarisPenyakit` = `pyMaxRecords` 500), bukan
// dikarang. Query tanpa batas atas tabel sebesar ini adalah cara paling mudah
// menahan basis data tanpa sengaja.
//
// ⛔ TABEL WARISAN, DIBACA SAJA. Tidak ada satu pun tulisan ke sana di
// seluruh berkas ini, dan tidak boleh ada: daftar penyakit bukan milik modul
// klaim.
//
// Dibaca sesudah: pesertapolis.go (pola query berbatasnya sama).

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/db"
	"nusantarare/modul/claimlife/models"
)

// Ketiga nama kolom `DISEASE_LIFE`.
//
// ✅ `[data DBA — katalog DEV 27-09-2026]` OQ-K.1 DITUTUP. Katalog
// `POOLDATA.DISEASE_LIFE`: `ID` VARCHAR2(100), `ICD_CODE` VARCHAR2(100),
// `DISEASE` VARCHAR2(1000); 97.586 baris; panjang isi maksimum `ICD_CODE` 7
// dan `DISEASE` 290.
//
// ⛔ RALAT ATAS TEBAKAN SAYA SENDIRI. Ronde pertama menulis
// `kolomNomorPenyakit = "NUMBER_"` dengan alasan *"NUMBER kata cadangan
// Oracle, jadi kolomnya pasti bernama lain, dan nama itu tidak ada di
// korpus"*. Separuh pertamanya benar; separuh keduanya KELIRU - namanya ADA
// di korpus, dan saya membacanya tanpa melihatnya:
//
//	`ReportDefinition/BrowseDiseaseLife_RD.xml`
//	  b598 `<pyFieldName>.Number</pyFieldName>`
//	  b599 `<pyFieldLabel>ID</pyFieldLabel>`      <- namanya, di baris berikutnya
//	  b676-677 pasangan yang sama, kedua kalinya
//
// Saya membaca `pyFieldLabel: ID` di keluaran grep saya sendiri dan
// memperlakukannya sebagai LABEL LAYAR, bukan sebagai nama kolom. Report
// definition memang memakai `pyFieldLabel` untuk keduanya ketika propertinya
// bernama lain dari kolomnya. Pelajarannya sempit dan tajam: sebelum
// menyatakan sesuatu "tidak ada di korpus", periksa apa yang sudah terbaca -
// bukan hanya apa yang sudah dicari.
//
// ⚠️ Lebar aslinya VARCHAR2(100)/(1000), jauh di atas isi terpanjangnya.
// Itu fakta tabel warisan yang dibaca, bukan yang kita buat; tidak ada yang
// perlu disesuaikan di sini.
const (
	kolomNomorPenyakit = "ID"
	kolomNamaPenyakit  = "DISEASE"
	kolomICDPenyakit   = "ICD_CODE"
)

// sqlCariPenyakit merakit query berbatas untuk `DISEASE_LIFE`.
//
// Meniru `ReportDefinition/BrowseDiseaseLife_RD.xml`:
//
//	b535/b754 `pyLogic` `A AND B`
//	  A  `.ICD_Code` `Contains` `Param.ICD_Code`
//	  B  `.Disease`  `Contains` `Param.Disease`
//	b598 urut `.Number` ASC
//	b659 `pyMaxRecords` 500
//
// ⛔ `A AND B`, bukan OR. Dua kata kunci mempersempit; menggantinya dengan OR
// akan mengembalikan seluruh penyakit yang kodenya cocok DITAMBAH seluruh
// yang namanya cocok - hasil yang lebih banyak dan lebih tidak berguna.
//
// ⛔ Kriteria yang KOSONG tidak menghasilkan klausa apa pun, bukan klausa
// `LIKE '%%'`. Keduanya bermakna sama di SQL, tetapi yang pertama membiarkan
// Oracle memakai indeks pada kolom yang satunya.
//
// ⚠️ Huruf besar dikerjakan di Go (`models.NormalkanKriteriaPenyakit`), bukan
// lewat `UPPER(:bind)`, supaya nilai yang DIKIRIM dan nilai yang
// DIBANDINGKAN adalah satu hal yang sama dan terlihat di log bind. Kolomnya
// tetap dibungkus `UPPER()` sebab isinya tidak dijamin berhuruf besar.
func sqlCariPenyakit(tabel string, k models.KriteriaPenyakit, batas int) (string, []any) {
	var syarat []string
	var arg []any

	if k.KodeICD != "" {
		arg = append(arg, k.KodeICD)
		syarat = append(syarat,
			"UPPER("+kolomICDPenyakit+") LIKE '%'||:"+strconv.Itoa(len(arg))+"||'%'")
	}
	if k.Nama != "" {
		arg = append(arg, k.Nama)
		syarat = append(syarat,
			"UPPER("+kolomNamaPenyakit+") LIKE '%'||:"+strconv.Itoa(len(arg))+"||'%'")
	}

	q := `SELECT ` + kolomNomorPenyakit + `, ` + kolomNamaPenyakit + `, ` +
		kolomICDPenyakit + ` FROM ` + tabel
	if len(syarat) > 0 {
		q += ` WHERE ` + strings.Join(syarat, " AND ")
	}
	// b598 `pySortOrder` 1 atas `.Number`, menaik.
	q += ` ORDER BY ` + kolomNomorPenyakit +
		` FETCH FIRST ` + strconv.Itoa(batas) + ` ROWS ONLY`
	return q, arg
}

// Penyakit membaca daftar diagnosa dari tabel warisan.
type Penyakit struct{ db *db.DB }

// NewPenyakit menyusun pembacanya.
func NewPenyakit(db *db.DB) *Penyakit { return &Penyakit{db: db} }

// Cari mengembalikan diagnosa yang cocok, SELALU berbatas.
func (r *Penyakit) Cari(ctx context.Context, k models.KriteriaPenyakit, batas int) (
	[]models.Penyakit, error) {

	tabel, err := r.db.Qualify("DISEASE_LIFE")
	if err != nil {
		return nil, err
	}
	q, arg := sqlCariPenyakit(tabel, k, models.BatasPenyakit(batas))
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := r.db.QueryContext(ctx, q, arg...)
	if err != nil {
		return nil, fmt.Errorf("repository: mencari diagnosa: %w", err)
	}
	defer baris.Close()

	// Daftar KOSONG, bukan nil: `encoding/json` menulis nil sebagai `null`.
	keluar := []models.Penyakit{}
	for baris.Next() {
		var nomor, nama, icd sql.NullString
		if err := baris.Scan(&nomor, &nama, &icd); err != nil {
			return nil, fmt.Errorf("repository: memindai diagnosa: %w", err)
		}
		keluar = append(keluar, models.Penyakit{
			Nomor: nomor.String, Nama: nama.String, KodeICD: icd.String,
		})
	}
	if err := baris.Err(); err != nil {
		return nil, fmt.Errorf("repository: mencari diagnosa: %w", err)
	}
	return keluar, nil
}
