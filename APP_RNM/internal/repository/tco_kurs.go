package repository

// Master kurs - tiket 11 Treaty Contract Out.
//
// Untuk apa berkas ini: membaca `TREATYEXCHANGEYEARLY` (`[data DBA]` nama
// sebenarnya; kueri Pega menyebut `treatyexchange`, `RDBList/GetMasterKursList.xml`
// b85). DIBACA SAJA - nol tulisan ke master kurs (AC tiket 11).
//
// ⚠️ Tabel (lanjutan 6, "pakai tabel yang RDB hidup pakai"): `treatyexchange`
// hanya disebut `GetMasterKursList`; 12 RDB korpus di enam modul lain
// (Endorsment Fac In, NB FacIn, NB Treaty In, RNW Fac In, Treaty In, Treaty In
// Adjustment) membaca `treatyexchangeyearly`, dan DBA menyebut `TREATYEXCHANGE` tidak ada
// (`dba-procedures.md`). `[belum diverifikasi executor di katalog DEV]`.
//
// ⛔ Seluruh kolomnya VARCHAR2, termasuk kedua tanggal. Sejak lanjutan 6
// tanggalnya diurai ORACLE dengan ekspresi RDB yang sama, seperti Pega - Go
// menerima `DATE` (AC 53), bukan teks.
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
	"time"

	"nusantarare/internal/models"
	"nusantarare/inti/db"
)

// MasterKursTahunanTCO - nama tabel master kurs `[data DBA]`.
const MasterKursTahunanTCO = "TREATYEXCHANGEYEARLY"

// MasterMataUangTCO - master mata uang (`GetCurrencyID`), dibaca saja lewat
// `MataUang.Pengenal`.
const MasterMataUangTCO = "CURRENCY"

// MasterKursTCO membaca master kurs.
type MasterKursTCO struct{ db *db.DB }

// NewMasterKursTCO menyusun pembacanya.
func NewMasterKursTCO(db *db.DB) *MasterKursTCO { return &MasterKursTCO{db: db} }

// sqlBerlakuKursTCO - kueri `GetMasterKursList` (b85-b86), ekspresi tanggal
// VERBATIM, dengan dua beda yang disengaja (lanjutan 6):
//   - `DEFAULT NULL ON CONVERSION ERROR` (Oracle 12.2+; DEV 12.2.0.1): tanggal
//     yang Oracle tolak menjadi NULL dan barisnya disaring, bukan menggagalkan
//     seluruh kueri seperti di Pega;
//   - `BETWEEN` dihitung sebagai kolom `BERLAKU`, bukan saringan `WHERE`,
//     supaya baris yang ditolak dapat dicacah.
//
// ⛔ Bind urut kemunculan (go-ora mengikat menurut posisi): :1 tanggal, :2
// QUARTER, :3 IDCURRENCY.
func sqlBerlakuKursTCO(tabel string) string {
	return fmt.Sprintf(`SELECT TOIDR, STARTDATE, ENDDATE, IDCURRENCY, CURRENCY, QUARTER, MULAI, AKHIR,
	 CASE WHEN TO_DATE(:1, 'YYYYMMDD') BETWEEN MULAI AND AKHIR THEN 1 ELSE 0 END BERLAKU
	 FROM (SELECT TOIDR, STARTDATE, ENDDATE, IDCURRENCY, CURRENCY, QUARTER,
	  TRUNC(TO_TIMESTAMP_TZ(STARTDATE DEFAULT NULL ON CONVERSION ERROR, '%s')) MULAI,
	  TRUNC(TO_TIMESTAMP_TZ(ENDDATE DEFAULT NULL ON CONVERSION ERROR, '%s')) AKHIR
	  FROM %s WHERE QUARTER = :2 AND IDCURRENCY = :3)`, models.FormatTanggalKursTCO, models.FormatTanggalKursTCO, tabel)
}

// BacaBerlaku membaca baris kurs satu mata uang dan satu `QUARTER` yang
// berlaku pada tanggal itu, beserta baris yang tanggalnya ditolak Oracle.
func (m *MasterKursTCO) BacaBerlaku(ctx context.Context, idCurrency, quarter string, tanggal time.Time) (models.HasilMasterKursTCO, error) {
	var h models.HasilMasterKursTCO
	tabel, err := m.db.Qualify(MasterKursTahunanTCO)
	if err != nil {
		return h, err
	}
	q := sqlBerlakuKursTCO(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return h, err
	}
	// `to_date({InputData.CARI1},'YYYYMMDD')` - CARI1 = `Param.StartDate`.
	rows, err := bacaTCO(ctx, m.db).QueryContext(ctx, q, tanggal.Format("20060102"), quarter, idCurrency)
	if err != nil {
		return h, fmt.Errorf("repository: membaca master %s: %w", MasterKursTahunanTCO, err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var n [6]sql.NullString
		var mulai, akhir sql.NullTime
		var berlaku int64
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &n[4], &n[5], &mulai, &akhir, &berlaku); err != nil {
			return models.HasilMasterKursTCO{}, err
		}
		tambahBarisKursTCO(&h, n, mulai, akhir, berlaku)
	}
	return h, rows.Err()
}

// tambahBarisKursTCO memilah satu baris hasil: tanggal NULL = Oracle menolak
// teksnya (atau kosong) - dilaporkan, tidak dipakai; `BERLAKU` 1 - dibawa.
// `TOIDR` dibawa sebagai teks - diurai hanya bila barisnya terpilih.
//
// ⛔ Hari diambil dari medan DATE apa adanya: go-ora memberi jam dinding
// Oracle di lokasi zona server, dan `.UTC()` dapat menggeser harinya.
func tambahBarisKursTCO(h *models.HasilMasterKursTCO, n [6]sql.NullString, mulai, akhir sql.NullTime, berlaku int64) {
	switch {
	case !mulai.Valid:
		h.Ditolak = append(h.Ditolak, models.BarisKursDitolakTCO{Kolom: "STARTDATE", Teks: n[1].String})
	case !akhir.Valid:
		h.Ditolak = append(h.Ditolak, models.BarisKursDitolakTCO{Kolom: "ENDDATE", Teks: n[2].String})
	case berlaku == 1:
		h.Berlaku = append(h.Berlaku, models.KursTCO{TeksToIDR: n[0].String, Mulai: hariKursTCO(mulai.Time),
			Akhir: hariKursTCO(akhir.Time), IDCurrency: n[3].String, Currency: n[4].String, Quarter: n[5].String})
	}
}

func hariKursTCO(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
