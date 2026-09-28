package repository

// Penerjemah kode mata uang menjadi pengenalnya - celah sensus 28-09-2026.
//
// Untuk apa berkas ini: `CURRENCYID` pada baris adjustment selama ini hanya
// DIBAWA (disalin dari baris sebelumnya, `services/adjustment.go`) dan tidak
// pernah DITERBITKAN. Akibatnya baris pertama sebuah klaim - yang tidak punya
// baris sebelumnya untuk disalin - lahir dengan `CURRENCYID` kosong, dan
// `models.HitungTotalPeserta` menolak menjumlahkannya.
//
// Ditemukan sensus paritas: `SetCurrencyID_Act` punya pemanggil di XML
// (`Section/AdjustmentDetail_Section.xml` b1206 `pyPreDataTransform`) tetapi
// nol padanan di repositori ini.
//
// Pohon rule-nya:
//
//	Activity/SetCurrencyID_Act.xml  kelas `Data-AdjustmentLife` b61
//	  b253-254  `InputData.CARI1` = `.CURRENCY`
//	  b361      `RDB-List` -> b419 `GetCurrencyID`
//	  b564-565  `.CURRENCYID` = `CurrencyList.pxResults(1).CARI1`
//
//	RDBList/GetCurrencyID.xml  b85
//	  SELECT ID AS CARI1 FROM POOLDATA.CURRENCY WHERE CURRENCY = {InputData.CARI1}
//
// ⛔ TABEL WARISAN, DIBACA SAJA. Nol tulisan ke `CURRENCY` di berkas ini, dan
// tidak boleh ada: daftar mata uang bukan milik modul klaim.
//
// Dibaca sesudah: penyakit.go (pola pembacaan tabel warisan yang sama).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ErrMataUangTidakDikenal - kode mata uang tidak ada di tabel warisan.
//
// ⛔ Galat, bukan teks kosong. Pengenal mata uang yang diam-diam kosong akan
// lolos sampai ke penjumlahan total, dan di sana ia menjadi "mata uang
// beragam" - kalimat yang benar tentang hal yang salah, jauh dari sebabnya.
var ErrMataUangTidakDikenal = errors.New("repository: kode mata uang tidak dikenal")

// MataUang membaca daftar mata uang warisan.
type MataUang struct{ db *DB }

// NewMataUang menyusun pembacanya.
func NewMataUang(db *DB) *MataUang { return &MataUang{db: db} }

// sqlPengenalMataUang merakit query-nya.
//
// ⚠️ `FETCH FIRST 1 ROWS ONLY` DITAMBAHKAN, dan rule aslinya tidak punya:
// Pega membaca `pxResults(1)` saja, yaitu baris pertama apa pun cacahnya.
// Menuliskannya di SQL membuat niat itu terbaca, dan menahan pembacaan tabel
// yang ternyata berisi kode kembar.
func sqlPengenalMataUang(tabel string) string {
	return fmt.Sprintf(
		`SELECT ID FROM %s WHERE CURRENCY = :1 FETCH FIRST 1 ROWS ONLY`, tabel)
}

// Pengenal menerjemahkan kode mata uang menjadi pengenalnya.
//
// Kode KOSONG mengembalikan pengenal kosong tanpa galat: peserta yang mata
// uangnya memang belum diisi bukan peserta yang mata uangnya salah
// (ADR-U-0027). Kode yang TERISI tetapi tidak dikenal adalah galat.
func (r *MataUang) Pengenal(ctx context.Context, kode string) (string, error) {
	kode = strings.TrimSpace(kode)
	if kode == "" {
		return "", nil
	}
	tabel, err := r.db.Qualify("CURRENCY")
	if err != nil {
		return "", err
	}
	q := sqlPengenalMataUang(tabel)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	var id sql.NullString
	err = r.db.sql.QueryRowContext(ctx, q, kode).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("%w: %q", ErrMataUangTidakDikenal, kode)
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca pengenal mata uang: %w", err)
	}
	return id.String, nil
}
