package repository

// Pembaca klaim ganda dan DOB peserta - `Save to RNM` langkah 11.
//
// Untuk apa berkas ini: `SaveOutStandingLife_Act` langkah 11.2-11.6 menolak
// klaim atas tertanggung yang sudah punya klaim di tabel warisan, dan 11.12
// menolak peserta tanpa tanggal lahir.
//
// `[terverifikasi]` `RDBList/CountPesertaAkseptasiLife_SQL.xml` (DEATH) dan
// `CountPesertaAkseptasiLifeHealth_SQL.xml` (bukan DEATH): keduanya menyaring
// `OS_AKSEPTASI_KLAIM_LIFE` dengan `CEDINGCO`, `NAME_OF_INSURED`, `DOB`,
// `CERTIFICATE_NO`, `PL_NUMBER`; yang kedua menambah `LAPSE_DATE` = DOL
// (tabel warisan menyimpan DOL di `LAPSE_DATE` - `UpdateDateClaimLife_SQL` b86).
//
// ⛔ NAMA DAN TANGGAL LAHIR TIDAK PERNAH MENINGGALKAN BASIS DATA. Aplikasi ini
// tidak menyalin keduanya ke tabel klaim (`kolomSalin`), jadi pencocokannya
// menggabung baris SUMBER peserta (`M_LIFE_PREMIUM_DETAIL`, lewat `SOURCE_ID`)
// di dalam SQL. Yang kembali ke Go hanya JAWABANNYA: status, ada/tidak, kosong/tidak.
//
// ⛔ BACA SAJA. Tabel warisan tidak ditulis di sini (brief GILIRAN-11 paket 1).
//
// ⚠️ `[terbuka - OQ-N2]` Baris warisan yang aplikasi ini sendiri tulis saat
// pendaftaran membiarkan `NAME_OF_INSURED`, `DOB`, dan `CEDINGCO` NULL
// (`pohonklaim.go`, 18 dari 55 kolom), sehingga klaim ganda antarklaim BARU
// tidak pernah cocok; yang tertangkap hanya baris era Pega.
//
// Dibaca sesudah: pesertapolis.go, pohonklaim.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// KunciPesertaSumber menunjuk SATU baris sumber peserta.
//
// ⛔ `PL_NUMBER` + `CERTIFICATE_NO` lebih dulu: itulah kolom ber-index di tabel
// 66,8 juta baris itu (penjaga `TestQueryTabelPesertaSelaluBerindexDanBerbatas`),
// dan itulah yang Pega pakai (`TempInputDetail.CARI4`/`CARI5`, langkah 11.1).
// `ID` (`SOURCE_ID` klaim) mengunci barisnya.
type KunciPesertaSumber struct {
	PLNumber, Sertifikat, SumberID string
}

// sqlDOBSumberKosong - langkah 11.12 `.DOB==""`, dijawab dari baris sumber.
func sqlDOBSumberKosong(sumber string) string {
	return fmt.Sprintf(`SELECT CASE WHEN DOB IS NULL THEN 1 ELSE 0 END FROM %s
	  WHERE PL_NUMBER = :1 AND CERTIFICATE_NO = :2 AND ID = :3
	  FETCH FIRST 1 ROWS ONLY`, sumber)
}

// sqlStatusWarisanTerakhir - `CountPesertaAkseptasiLife_SQL`, baris pertama
// `ORDER BY ACCEPTATION_DATE DESC` (langkah 11.4 membaca `pxResults(1)`).
//
// Kelima kunci SQL aslinya VERBATIM; `DOB = TO_DATE(FormatDateTime(.DOB,
// 'dd/MM/YYYY'))` di Pega membandingkan HARI, jadi sisi sumber di-TRUNC.
//
// ⚠️ Urutan NULL mengikuti Oracle apa adanya (DESC: NULL lebih dulu), sama
// dengan SQL aslinya - tidak "diperbaiki".
func sqlStatusWarisanTerakhir(lama, sumber string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(o.STS_REJECT)
	   FROM %s o, %s m
	  WHERE m.PL_NUMBER = :1 AND m.CERTIFICATE_NO = :2 AND m.ID = :3
	    AND o.CEDINGCO = :4
	    AND o.NAME_OF_INSURED = m.NAME_OF_INSURED
	    AND o.DOB = TRUNC(m.DOB)
	    AND o.CERTIFICATE_NO = m.CERTIFICATE_NO
	    AND o.PL_NUMBER = m.PL_NUMBER
	  ORDER BY o.ACCEPTATION_DATE DESC
	  FETCH FIRST 1 ROWS ONLY`, lama, sumber)
}

// sqlAdaWarisanSamaDOL - `CountPesertaAkseptasiLifeHealth_SQL`, sebagai uji ADA.
//
// ⚠️ CACAT RULE WARISAN, ditiru MAKSUDnya: `LAPSE_DATE = TO_DATE({CARI6},
// 'DD/MM/YYYY')` padahal `CARI6 = .DATE_OF_LOSS` mentah, tanpa
// `@FormatDateTime` (SaveOutStandingLife_Act b3189). Maksudnya jelas: DOL yang
// sama. Dilaporkan OQ-N2.
func sqlAdaWarisanSamaDOL(lama, sumber string) string {
	return fmt.Sprintf(`SELECT 1
	   FROM %s o, %s m
	  WHERE m.PL_NUMBER = :1 AND m.CERTIFICATE_NO = :2 AND m.ID = :3
	    AND o.CEDINGCO = :4
	    AND o.NAME_OF_INSURED = m.NAME_OF_INSURED
	    AND o.DOB = TRUNC(m.DOB)
	    AND o.CERTIFICATE_NO = m.CERTIFICATE_NO
	    AND o.PL_NUMBER = m.PL_NUMBER
	    AND o.LAPSE_DATE = TO_DATE(:5, 'YYYY-MM-DD')
	  FETCH FIRST 1 ROWS ONLY`, lama, sumber)
}

func (r *KlaimLife) duaTabelGanda() (lama, sumber string, err error) {
	if lama, err = r.db.Qualify(namaTabelLama); err != nil {
		return "", "", err
	}
	if sumber, err = r.db.Qualify(namaTabelPeserta); err != nil {
		return "", "", err
	}
	return lama, sumber, nil
}

// DOBSumberKosong menjawab apakah tanggal lahir baris sumber peserta NULL.
func (r *KlaimLife) DOBSumberKosong(ctx context.Context, k KunciPesertaSumber) (bool, error) {
	sumber, err := r.db.Qualify(namaTabelPeserta)
	if err != nil {
		return false, err
	}
	q := sqlDOBSumberKosong(sumber)
	if err := PeriksaSQL(q); err != nil {
		return false, err
	}
	var kosong int
	if err := r.db.sql.QueryRowContext(ctx, q, k.PLNumber, k.Sertifikat, k.SumberID).Scan(&kosong); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("repository: baris sumber peserta %q tidak ada", k.SumberID)
		}
		return false, fmt.Errorf("repository: membaca DOB sumber: %w", err)
	}
	return kosong == 1, nil
}

// StatusWarisanTerakhir mengembalikan STS_REJECT baris warisan terbaru
// bertertanggung sama; kosong bila tidak ada.
func (r *KlaimLife) StatusWarisanTerakhir(ctx context.Context, k KunciPesertaSumber, cedingCo string) (string, error) {
	lama, sumber, err := r.duaTabelGanda()
	if err != nil {
		return "", err
	}
	q := sqlStatusWarisanTerakhir(lama, sumber)
	if err := PeriksaSQL(q); err != nil {
		return "", err
	}
	var status sql.NullString
	err = r.db.sql.QueryRowContext(ctx, q, k.PLNumber, k.Sertifikat, k.SumberID,
		kosongJadiNil(cedingCo)).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("repository: membaca klaim ganda warisan: %w", err)
	}
	return strings.TrimSpace(status.String), nil
}

// AdaWarisanSamaDOL menjawab apakah tertanggung yang sama sudah punya baris
// warisan ber-DOL sama. `dol` berbentuk `YYYY-MM-DD`.
func (r *KlaimLife) AdaWarisanSamaDOL(ctx context.Context, k KunciPesertaSumber, cedingCo, dol string) (bool, error) {
	lama, sumber, err := r.duaTabelGanda()
	if err != nil {
		return false, err
	}
	q := sqlAdaWarisanSamaDOL(lama, sumber)
	if err := PeriksaSQL(q); err != nil {
		return false, err
	}
	var satu int
	err = r.db.sql.QueryRowContext(ctx, q, k.PLNumber, k.Sertifikat, k.SumberID,
		kosongJadiNil(cedingCo), dol).Scan(&satu)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("repository: membaca klaim ganda health warisan: %w", err)
	}
	return true, nil
}
