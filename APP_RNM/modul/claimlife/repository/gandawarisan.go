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

	"nusantarare/inti/db"
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
//
// ⛔ OQ-N2 (GILIRAN-17): dua saringan TAMBAHAN, sebab cermin kini mengisi
// nama/DOB/CEDINGCO. (1) `o.STS_REJECT IS NOT NULL` - di Pega baris cermin
// baru LAHIR saat Save Outstanding berstatus '0' (`InsertJsonKlaimLife_sql`
// b176); baris aplikasi ini yang lahir saat pendaftaran (status dan
// ACCEPTATION_DATE NULL) tidak ada di dunia Pega, dan NULL-lebih-dulu akan
// menaruhnya di DEPAN baris era Pega yang sah. (2) klaim SENDIRI (`CASEID`)
// dikecualikan - di Pega pemeriksaan ini berjalan (langkah 11.x) SEBELUM
// barisnya sendiri ditulis (22.1.3.1). `CASEID` warisan boleh NULL, maka
// `IS NULL OR <>`.
func sqlStatusWarisanTerakhir(lama, sumber string) string {
	return fmt.Sprintf(`SELECT TO_CHAR(o.STS_REJECT)
	   FROM %s o, %s m
	  WHERE m.PL_NUMBER = :1 AND m.CERTIFICATE_NO = :2 AND m.ID = :3
	    AND o.CEDINGCO = :4
	    AND o.NAME_OF_INSURED = m.NAME_OF_INSURED
	    AND o.DOB = TRUNC(m.DOB)
	    AND o.CERTIFICATE_NO = m.CERTIFICATE_NO
	    AND o.PL_NUMBER = m.PL_NUMBER
	    AND o.STS_REJECT IS NOT NULL
	    AND (o.CASEID IS NULL OR o.CASEID <> :5)
	  ORDER BY o.ACCEPTATION_DATE DESC
	  FETCH FIRST 1 ROWS ONLY`, lama, sumber)
}

// sqlAdaWarisanSamaDOL - `CountPesertaAkseptasiLifeHealth_SQL`, sebagai uji ADA.
//
// ⚠️ CACAT RULE WARISAN, ditiru MAKSUDnya: `LAPSE_DATE = TO_DATE({CARI6},
// 'DD/MM/YYYY')` padahal `CARI6 = .DATE_OF_LOSS` mentah, tanpa
// `@FormatDateTime` (SaveOutStandingLife_Act b3189). Maksudnya jelas: DOL yang
// sama. Dilaporkan OQ-N2.
//
// ⛔ OQ-N2 (GILIRAN-17): klaim SENDIRI dikecualikan, alasan yang sama dengan
// `sqlStatusWarisanTerakhir`. Baris cermin milik aplikasi tidak menulis
// LAPSE_DATE, jadi saringan status tidak diperlukan di sini.
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
	    AND (o.CASEID IS NULL OR o.CASEID <> :6)
	  FETCH FIRST 1 ROWS ONLY`, lama, sumber)
}

// sqlIsiTertanggungCermin - OQ-N2 DITUTUP 29-09-2026 (GILIRAN-17): cermin
// `OS_AKSEPTASI_KLAIM_LIFE` mengisi NAME_OF_INSURED, DOB, CEDINGCO seperti Pega
// (`SaveOutStandingLife_Act` 22.1.1 b8747: CARI6 `.NAME_OF_INSURED` b8906,
// CARI8 `@FormatDateTime(.DOB,"dd/MM/YYYY")` b8946, CARI27
// `pyWorkPage.PolicyDataLife.CedingCo` b9226 -> `InsertJsonKlaimLife_sql`
// b93/b95/b114).
//
// ⛔ DI DALAM SQL: nama dan tanggal lahir disalin Oracle dari baris sumber
// `M_LIFE_PREMIUM_DETAIL` - tidak pernah dibaca, diikat, dicatat, atau diuji
// dari Go. `TRUNC(m.DOB)` = hari, padanan `@FormatDateTime(..,"dd/MM/YYYY")`.
// CEDINGCO pernyataan tersendiri (`sqlIsiCedingCermin`, pohonklaim.go).
//
// `EXISTS` menjadikan baris sumber yang tidak ada GALAT (nol baris), bukan
// NULL diam-diam. Tiap subkueri berbatas (penjaga tabel 66,8 juta baris).
func sqlIsiTertanggungCermin(lama, sumber string) string {
	return fmt.Sprintf(`UPDATE %s o
	   SET (NAME_OF_INSURED, DOB) =
	       (SELECT m.NAME_OF_INSURED, TRUNC(m.DOB) FROM %s m
	         WHERE m.PL_NUMBER = :1 AND m.CERTIFICATE_NO = :2 AND m.ID = :3 AND ROWNUM = 1)
	 WHERE o.ID = :4 AND o.CASEID = :5
	   AND EXISTS (SELECT 1 FROM %s m2
	                WHERE m2.PL_NUMBER = :6 AND m2.CERTIFICATE_NO = :7 AND m2.ID = :8 AND ROWNUM = 1)`,
		lama, sumber, sumber)
}

// IsiTertanggungCermin mengisi NAME_OF_INSURED, DOB, CEDINGCO satu baris
// cermin (`adjID`) dari baris sumbernya - lihat `sqlIsiTertanggungCermin`.
// Di dalam transaksi pemanggilnya; baris sumber yang tidak ada = galat.
//
// ⛔ Dikunci `ID` DAN `CASEID` (temuan /code-review): pengenal baris aplikasi
// adalah angka sequence, dan baris era Pega berpengenal sama milik klaim lain
// tidak boleh tersentuh.
func (r *KlaimLife) IsiTertanggungCermin(ctx context.Context, tx *db.Tx, adjID, caseID string,
	k KunciPesertaSumber, nomorPolis string) error {

	lama, sumber, err := r.duaTabelGanda()
	if err != nil {
		return err
	}
	q := sqlIsiTertanggungCermin(lama, sumber)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, k.PLNumber, k.Sertifikat, k.SumberID,
		adjID, caseID, k.PLNumber, k.Sertifikat, k.SumberID)
	if err != nil {
		return fmt.Errorf("repository: mengisi tertanggung cermin %s: %w", adjID, err)
	}
	if err := db.PastikanSatuBaris(hasil, "pengisian tertanggung cermin"); err != nil {
		return fmt.Errorf("%w: baris sumber peserta %q atau baris cermin %q tidak ada", err, k.SumberID, adjID)
	}
	return r.isiCedingCermin(ctx, tx, lama, adjID, caseID, nomorPolis)
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
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var kosong int
	if err := r.db.QueryRowContext(ctx, q, k.PLNumber, k.Sertifikat, k.SumberID).Scan(&kosong); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("repository: baris sumber peserta %q tidak ada", k.SumberID)
		}
		return false, fmt.Errorf("repository: membaca DOB sumber: %w", err)
	}
	return kosong == 1, nil
}

// StatusWarisanTerakhir mengembalikan STS_REJECT baris warisan terbaru
// bertertanggung sama; kosong bila tidak ada.
//
// `caseID` - klaim yang sedang diperiksa; barisnya sendiri dikecualikan (OQ-N2).
func (r *KlaimLife) StatusWarisanTerakhir(ctx context.Context, k KunciPesertaSumber, cedingCo, caseID string) (string, error) {
	lama, sumber, err := r.duaTabelGanda()
	if err != nil {
		return "", err
	}
	q := sqlStatusWarisanTerakhir(lama, sumber)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	var status sql.NullString
	err = r.db.QueryRowContext(ctx, q, k.PLNumber, k.Sertifikat, k.SumberID,
		db.KosongJadiNil(cedingCo), caseID).Scan(&status)
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
func (r *KlaimLife) AdaWarisanSamaDOL(ctx context.Context, k KunciPesertaSumber, cedingCo, dol, caseID string) (bool, error) {
	lama, sumber, err := r.duaTabelGanda()
	if err != nil {
		return false, err
	}
	q := sqlAdaWarisanSamaDOL(lama, sumber)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var satu int
	err = r.db.QueryRowContext(ctx, q, k.PLNumber, k.Sertifikat, k.SumberID,
		db.KosongJadiNil(cedingCo), dol, caseID).Scan(&satu)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("repository: membaca klaim ganda health warisan: %w", err)
	}
	return true, nil
}
