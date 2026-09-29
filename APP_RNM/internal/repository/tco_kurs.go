package repository

// Master kurs - tiket 11 Treaty Contract Out.
//
// Untuk apa berkas ini: membaca `TREATYEXCHANGEYEARLY` (`[data DBA]` nama
// sebenarnya; kueri Pega menyebut `treatyexchange`, `RDBList/GetMasterKursList.xml`
// b85). DIBACA SAJA - nol tulisan ke master kurs (AC tiket 11).
//
// ⛔ Seluruh kolomnya VARCHAR2, termasuk kedua tanggal. Pega mengurainya di
// SETIAP kueri (`TO_TIMESTAMP_TZ(STARTDATE, ...)`); di sini baris kandidat
// dibaca sekali dan tanggal/nilainya diurai di batas baca menjadi tanggal dan
// desimal (AC 53, ADR-0003) - teks yang tidak terurai adalah galat terang.
//
// ⛔ Saringan `QUARTER` dan `IDCURRENCY` di-bind - tidak ada identitas mata
// uang di teks SQL (AC 47); pengenalnya dibaca dari master `CURRENCY` lewat
// `MataUang.Pengenal` (kode bersama, tidak diubah).
//
// Dibaca sesudah: models/tco_kurs.go, matauangid.go.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/internal/models"
)

// MasterKursTahunanTCO - nama tabel master kurs `[data DBA]`.
const MasterKursTahunanTCO = "TREATYEXCHANGEYEARLY"

// MasterMataUangTCO - master mata uang (`GetCurrencyID`), dibaca saja lewat
// `MataUang.Pengenal`.
const MasterMataUangTCO = "CURRENCY"

// MasterKursTCO membaca master kurs.
type MasterKursTCO struct{ db *DB }

// NewMasterKursTCO menyusun pembacanya.
func NewMasterKursTCO(db *DB) *MasterKursTCO { return &MasterKursTCO{db: db} }

// sqlDaftarKursTCO - kandidat kurs satu mata uang dan satu `QUARTER`; rentang
// tanggal dipilih sesudah diurai (`models.PilihKursBerlakuTCO`).
func sqlDaftarKursTCO(tabel string) string {
	return fmt.Sprintf(`SELECT TOIDR, STARTDATE, ENDDATE, IDCURRENCY, CURRENCY, QUARTER FROM %s
	 WHERE QUARTER = :1 AND IDCURRENCY = :2`, tabel)
}

// Daftar membaca dan mengurai baris kurs kandidat.
func (m *MasterKursTCO) Daftar(ctx context.Context, idCurrency, quarter string) ([]models.KursTCO, error) {
	tabel, err := m.db.Qualify(MasterKursTahunanTCO)
	if err != nil {
		return nil, err
	}
	q := sqlDaftarKursTCO(tabel)
	if err := PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := m.db.sql.QueryContext(ctx, q, quarter, idCurrency)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca master %s: %w", MasterKursTahunanTCO, err)
	}
	defer func() { _ = rows.Close() }()
	var hasil []models.KursTCO
	for rows.Next() {
		var n [6]sql.NullString
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5]); err != nil {
			return nil, err
		}
		k, err := uraiBarisKursTCO(n)
		if err != nil {
			return nil, err
		}
		hasil = append(hasil, k)
	}
	return hasil, rows.Err()
}

// uraiBarisKursTCO mengurai satu baris master; galatnya menyebut baris itu.
func uraiBarisKursTCO(n [6]sql.NullString) (models.KursTCO, error) {
	k := models.KursTCO{IDCurrency: n[3].String, Currency: n[4].String, Quarter: n[5].String}
	var err error
	if k.ToIDR, err = models.UraiNilaiKursTCO(n[0].String); err != nil {
		return models.KursTCO{}, fmt.Errorf("%w (baris STARTDATE %q)", err, n[1].String)
	}
	if k.Mulai, err = models.UraiTanggalKursTCO("STARTDATE", n[1].String); err != nil {
		return models.KursTCO{}, err
	}
	if k.Akhir, err = models.UraiTanggalKursTCO("ENDDATE", n[2].String); err != nil {
		return models.KursTCO{}, err
	}
	return k, nil
}
