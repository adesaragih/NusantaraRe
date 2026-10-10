package repository

// Following Policy Number (tiket 70, kontrak sesi 9d; K70-1..3, A70-1..4 di docs/issues/70-following-policy.md):
//   - cari polis = RDBList\BrowseDataViewPolis1.xml atas POOLDATA.FACINPRODUCTION (DISTINCT, NOENDORS IS NULL) dengan klausa
//     GetListViewFollowing_act langkah 4.1-4.9 - BIND, tanpa {ASIS:} (A70-1); cadangan json_following tidak diport (K70-1);
//   - Copy Nopolis / Copy All = T_GENERAL_POLIS.FOLLOWING (+ PERCENT_SHARE / PAY_RI_COMMISION untuk Copy All). FollowingNB /
//     IDFollowingNB tidak punya kolom di sistem baru (belum dimigrasi - pertanyaan ke work owner).

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbfacin/backend/models"
)

// Jenis klausa cari (langkah GetListViewFollowing_act).
const (
	CariFollowingNoPolis   = "nopolis"  // 4.1 opsi 2
	CariFollowingIDPega    = "idpega"   // 4.2 opsi 1 + "NB-" / "RNW-"
	CariFollowingNoEndors  = "noendors" // 4.4 filter memuat "RNM-E"
	CariFollowingEDM       = "edm"      // 4.5 opsi 1 + "EDM-"
	CariFollowingReff      = "reff"     // 4.6 opsi 4
	CariFollowingInsured   = "insured"  // 4.7 opsi 5
	CariFollowingSOB       = "sob"      // 4.8 opsi 6
	CariFollowingQQ        = "qq"       // 4.9 opsi 7
	ukuranHalamanFollowing = 20
)

// PenyimpanFollowing - data Following Policy.
type PenyimpanFollowing interface {
	// CariFollowing - satu halaman (offset, batas) baris polis + jumlah seluruhnya. Jenis tak dikenal -> galat.
	CariFollowing(ctx context.Context, jenis, nilai string, offset, batas int) ([]models.BarisFollowing, int, error)
	// AdaPolisFollowing - FACINPRODUCTION punya baris IDPEGA case itu ber-NOPOLIS itu (NOENDORS kosong).
	AdaPolisFollowing(ctx context.Context, caseID, noPolis string) (bool, error)
	// TulisFollowing - FOLLOWING; nilai != nil -> juga PERCENT_SHARE / PAY_RI_COMMISION (Copy All). Baris General belum ada ->
	// disisipkan. Di transaksi pemanggil.
	TulisFollowing(ctx context.Context, tx *db.Tx, caseID, noPolis string, nilai *models.NilaiFollowing) error
}

// FollowingOracle - PenyimpanFollowing atas Oracle.
type FollowingOracle struct{ db *db.DB }

// NewFollowingOracle merakit penyimpan Following Policy.
func NewFollowingOracle(d *db.DB) *FollowingOracle { return &FollowingOracle{db: d} }

// klausaFollowing - WHERE langkah 4.x dengan bind :1 (nilai sudah berbentuk akhir: awalan IDPEGA / pola LIKE disusun services).
func klausaFollowing(prod, jenis string) (string, bool) {
	k, ada := map[string]string{
		CariFollowingNoPolis:  "PC0.NOPOLIS = :1",
		CariFollowingIDPega:   "PC0.IDPEGA = :1",
		CariFollowingNoEndors: "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM " + prod + " WHERE NOENDORS = :1)",
		CariFollowingEDM:      "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM " + prod + " WHERE IDPEGA = :1)",
		CariFollowingReff:     "PC0.NOPOLIS IN (SELECT DISTINCT NOPOLIS FROM " + prod + " WHERE RISLIPCEDING LIKE :1)",
		CariFollowingInsured:  "PC0.INSUREDNAME LIKE :1",
		CariFollowingSOB:      "PC0.SOB LIKE :1",
		CariFollowingQQ:       "PC0.QQNAME LIKE :1",
		CariPolisQQ:           "UPPER(PC0.QQNAME) LIKE UPPER(:1)", // tiket 71 opsi 7 (A71-2)
	}[jenis]
	return k + " AND PC0.NOENDORS IS NULL", ada
}

// sqlDasarFollowing - BrowseDataViewPolis1 (kolom + subkueri RiSlipCeding apa adanya; DISTINCT atas tanggal mentah).
func sqlDasarFollowing(prod, klausa string) string {
	return "SELECT DISTINCT SUBSTR(PC0.IDPEGA, 19, 10) AS CASEID, PC0.NOPOLIS, PC0.BEGINDATE, PC0.ENDDATE, PC0.SOB, PC0.INSUREDNAME," +
		" PC0.QQNAME, (SELECT MAX(B.RISLIPCEDING) KEEP (DENSE_RANK LAST ORDER BY B.SPREAD_DATE) FROM " + prod +
		" B WHERE B.NOPOLIS = PC0.NOPOLIS) AS RISLIP FROM " + prod + " PC0 WHERE " + klausa
}

func sqlJumlahFollowing(dasar string) string { return "SELECT COUNT(*) FROM (" + dasar + ")" }

// sqlHalamanFollowing - Pega tanpa ORDER BY; urutan stabil NOPOLIS, CASEID, lalu kolom lain `[keputusan agent]`.
func sqlHalamanFollowing(dasar string) string {
	return "SELECT CASEID, NOPOLIS, TO_CHAR(BEGINDATE, 'DD-MM-YYYY'), TO_CHAR(ENDDATE, 'DD-MM-YYYY'), SOB, INSUREDNAME, QQNAME, RISLIP" +
		" FROM (" + dasar + ") ORDER BY NOPOLIS, CASEID, BEGINDATE, ENDDATE, SOB, INSUREDNAME, QQNAME OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY"
}

func sqlAdaPolisFollowing(prod string) string {
	return "SELECT COUNT(*) FROM " + prod + " WHERE IDPEGA = :1 AND NOPOLIS = :2 AND NOENDORS IS NULL"
}

func sqlFollowing(general string) string {
	return "UPDATE " + general + " SET FOLLOWING = :1 WHERE ID = :2"
}

func sqlSisipFollowing(general string) string {
	return "INSERT INTO " + general + " (ID, FOLLOWING) VALUES (:1, :2)"
}

func sqlSalinFollowing(general string) string {
	return "UPDATE " + general + " SET FOLLOWING = :1, PERCENT_SHARE = " + angkaMasuk(":2") + ", PAY_RI_COMMISION = " +
		angkaMasuk(":3") + " WHERE ID = :4"
}

func sqlSisipSalinFollowing(general string) string {
	return "INSERT INTO " + general + " (ID, FOLLOWING, PERCENT_SHARE, PAY_RI_COMMISION) VALUES (:1, :2, " + angkaMasuk(":3") + ", " +
		angkaMasuk(":4") + ")"
}

// CariFollowing - lihat PenyimpanFollowing.
func (r *FollowingOracle) CariFollowing(ctx context.Context, jenis, nilai string, offset, batas int) ([]models.BarisFollowing, int, error) {
	prod, err := r.db.Qualify(TabelFacInProduction)
	if err != nil {
		return nil, 0, err
	}
	klausa, ada := klausaFollowing(prod, jenis)
	if !ada {
		return nil, 0, fmt.Errorf("repository: jenis cari Following %q tidak dikenal", jenis)
	}
	dasar := sqlDasarFollowing(prod, klausa)
	q := sqlJumlahFollowing(dasar)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, 0, err
	}
	var total int
	if err := r.db.QueryRowContext(ctx, q, nilai).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("repository: jumlah polis Following: %w", err)
	}
	q = sqlHalamanFollowing(dasar)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, 0, err
	}
	rows, err := r.db.QueryContext(ctx, q, nilai, offset, batas)
	if err != nil {
		return nil, 0, fmt.Errorf("repository: cari polis Following: %w", err)
	}
	defer rows.Close()
	hasil := []models.BarisFollowing{}
	for rows.Next() {
		var v [8]sql.NullString
		if err := rows.Scan(&v[0], &v[1], &v[2], &v[3], &v[4], &v[5], &v[6], &v[7]); err != nil {
			return nil, 0, fmt.Errorf("repository: cari polis Following: %w", err)
		}
		hasil = append(hasil, models.BarisFollowing{CaseID: v[0].String, PolicyNo: v[1].String, StartPeriod: v[2].String,
			EndPeriod: v[3].String, SourceOfBusiness: v[4].String, InsuredName: v[5].String, QQName: v[6].String,
			ReffNumber: v[7].String})
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("repository: cari polis Following: %w", err)
	}
	return hasil, total, nil
}

// AdaPolisFollowing - lihat PenyimpanFollowing.
func (r *FollowingOracle) AdaPolisFollowing(ctx context.Context, caseID, noPolis string) (bool, error) {
	prod, err := r.db.Qualify(TabelFacInProduction)
	if err != nil {
		return false, err
	}
	q := sqlAdaPolisFollowing(prod)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := r.db.QueryRowContext(ctx, q, AwalanIDPega+caseID, noPolis).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: polis Following: %w", err)
	}
	return n > 0, nil
}

// TulisFollowing - lihat PenyimpanFollowing.
func (r *FollowingOracle) TulisFollowing(ctx context.Context, tx *db.Tx, caseID, noPolis string, nilai *models.NilaiFollowing) error {
	general, err := r.db.Qualify(TabelGeneralPolis)
	if err != nil {
		return err
	}
	d := func(x *apd.Decimal) any { return ikatDesimal(x) }
	ubah, sisip := sqlFollowing(general), sqlSisipFollowing(general)
	argUbah, argSisip := []any{noPolis, caseID}, []any{caseID, noPolis}
	if nilai != nil {
		ubah, sisip = sqlSalinFollowing(general), sqlSisipSalinFollowing(general)
		argUbah = []any{noPolis, d(nilai.PercentShare), d(nilai.RICommision), caseID}
		argSisip = []any{caseID, noPolis, d(nilai.PercentShare), d(nilai.RICommision)}
	}
	hasil, err := jalankan(ctx, tx, ubah, "FOLLOWING", argUbah...)
	if err != nil {
		return err
	}
	if n, err := hasil.RowsAffected(); err != nil {
		return fmt.Errorf("repository: FOLLOWING: %w", err)
	} else if n == 0 {
		_, err = jalankan(ctx, tx, sisip, "FOLLOWING", argSisip...)
		return err
	}
	return nil
}
