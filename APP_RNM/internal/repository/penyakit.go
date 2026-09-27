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

	"nusantarare/internal/models"
)

// Ketiga nama kolom `DISEASE_LIFE`.
//
// ⛔ `[terbuka - OQ-K]` DAN INI HARUS DIBACA SEBELUM DIPAKAI DI ORACLE MANA
// PUN. Ekspor yang kami terima memuat NAMA PROPERTI Pega (`.Number`,
// `.Disease`, `.ICD_Code`) tetapi **tidak** memuat pemetaan kelas-ke-tabelnya -
// tidak ada `Rule-Obj-Class` untuk `Int-DISEASE_LIFE` di seluruh korpus.
//
//   - `DISEASE` dan `ICD_CODE` `[dugaan kuat]`: keduanya nama kolom yang sama
//     persis dengan `T_CLAIMLF_PREMIUMLIST_DETAIL` (migrasi 003 b42-43), dan
//     konvensi korpus ini memetakan properti ke kolom senama berhuruf besar.
//   - `NUMBER_` `[terbuka]`: `NUMBER` adalah kata cadangan Oracle, sehingga
//     kolomnya PASTI bernama lain - dan nama itu tidak ada di korpus.
//
// Ketiganya dikumpulkan DI SINI, satu tempat, supaya koreksi dari DBA adalah
// satu suntingan dan bukan perburuan. Uji `TestKolomPenyakitBelumDipastikan`
// menagih OQ-K supaya pertanyaannya tidak hilang bersama giliran ini.
const (
	kolomNomorPenyakit = "NUMBER_"
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
type Penyakit struct{ db *DB }

// NewPenyakit menyusun pembacanya.
func NewPenyakit(db *DB) *Penyakit { return &Penyakit{db: db} }

// Cari mengembalikan diagnosa yang cocok, SELALU berbatas.
func (r *Penyakit) Cari(ctx context.Context, k models.KriteriaPenyakit, batas int) (
	[]models.Penyakit, error) {

	tabel, err := r.db.Qualify("DISEASE_LIFE")
	if err != nil {
		return nil, err
	}
	q, arg := sqlCariPenyakit(tabel, k, models.BatasPenyakit(batas))
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	baris, err := r.db.sql.QueryContext(ctx, q, arg...)
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
