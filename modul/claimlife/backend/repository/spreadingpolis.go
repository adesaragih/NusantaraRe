package repository

// Spreading peserta klaim DIAMBIL dari PremiumList Life — keputusan work owner 11-10-2026: "untuk spreading klaim datanya
// ambil saja dari premium list, tidak perlu perhitungan lagi. di premium list life sudah ada spreading pesertanya, tinggal
// ambil".
//
// Tabel milik PremiumList Life (migrasi 052-054), DIBACA SAJA; SQL disalin bentuknya dari
// premiumlistlife/backend/repository/polis_spreading.go (sqlRincianSpreadingPeserta), paket modul itu TIDAK diimpor:
//
//	T_PREMIUM_LIST_DETAIL           peserta versi polis; kunci PREMIUM_LIST_ID (IDX_PLD_PL)
//	T_PREMIUM_LIST_SPREADING        lapisan per peserta; DETAIL_ID (IDX_PLS_DETAIL)
//	T_PREMIUM_LIST_SPREADING_RETRO  reinsurer per lapisan; SPREADING_ID (IDX_PLSR_SPR)
//
// Kunci peserta: header IDPEGA klaim (= T_PREMIUM_LIST.ID versi polis = PREMIUM_LIST_ID) + CERTIFICATE_NO (TRIM), tanpa
// baris EDM berstatus 'Delete'. ⛔ TIDAK lewat ID: M_LIFE_PREMIUM_DETAIL.ID (SOURCE_ID peserta klaim) adalah sequence
// sendiri (premiumlistlife polis_warisan.go), bukan T_PREMIUM_LIST_DETAIL.ID.
//
// Versi EDM: Endorsement Life MENYALIN spreading peserta ke versi barunya (endorsementlife edm_buat.go sqlSalinSpreading,
// AC 59), jadi IDPEGA EDML- punya spreading sendiri. Versi warisan Pega tidak punya spreading di tabel aplikasi — kosong.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/uang"
	"nusantarare/modul/claimlife/backend/models"
)

const (
	tabelDetailPolis    = "T_PREMIUM_LIST_DETAIL"
	tabelSpreadingPolis = "T_PREMIUM_LIST_SPREADING"
	tabelSpreadingRetro = "T_PREMIUM_LIST_SPREADING_RETRO"
	// edmStatusHapus - peserta yang dicabut di versi EDM (endorsementlife: EDM_STATUS 'Delete').
	edmStatusHapus = "Delete"
)

// sqlDetailPolisPeserta - peserta versi polis :1 bersertifikat :2. Bila beberapa baris cocok, yang ber-ID terkecil.
func sqlDetailPolisPeserta(detail string) string {
	return fmt.Sprintf(`SELECT d.ID, d.CURRENCY FROM %s d
	 WHERE d.PREMIUM_LIST_ID = :1 AND TRIM(d.CERTIFICATE_NO) = TRIM(:2)
	   AND (d.EDM_STATUS IS NULL OR d.EDM_STATUS <> :3)
	 ORDER BY d.ID FETCH FIRST 1 ROWS ONLY`, detail)
}

// sqlSpreadingPolisPeserta - lapisan -> reinsurer satu peserta polis, urut ID (urut lapisan lalu reinsurer, seperti
// PremiumList). LEFT JOIN: lapisan tanpa reinsurer tetap tampil.
func sqlSpreadingPolisPeserta(spreading, retro string) string {
	return fmt.Sprintf(`SELECT s.ID, s.TREATY_TYPE_NAME, s.TREATY_YEAR_LIFE, `+db.FmtDesimal+`,
	       r.REINSURER_NAME, `+db.FmtDesimal+`, `+db.FmtDesimal+`, `+db.FmtDesimal+`, `+db.FmtDesimal+`, `+db.FmtDesimal+`,
	       `+db.FmtDesimal+`, `+db.FmtDesimal+`
	  FROM %s s LEFT JOIN %s r ON r.SPREADING_ID = s.ID
	 WHERE s.DETAIL_ID = :1
	 ORDER BY s.ID, r.ID`,
		"s.RETROCADED_SHARE", "r.PERCENT_SHARE", "r.AMOUNT", "r.RATE", "r.COMMISION", "r.OVR_COMM",
		"r.PREMIUM_SPREADED_GROSS", "r.PREMIUM_SPREADED_NET", spreading, retro)
}

// SpreadingPolisPeserta membaca spreading PremiumList satu peserta. IDPEGA atau sertifikat kosong, atau peserta tidak ada
// di versi itu = nil (layar: "No spreading data for this insured in the policy").
func (r *KlaimLife) SpreadingPolisPeserta(ctx context.Context, idPega, sertifikat string) ([]models.Spreading, error) {
	if idPega == "" || sertifikat == "" {
		return nil, nil
	}
	detail, err := r.db.Qualify(tabelDetailPolis)
	if err != nil {
		return nil, err
	}
	spreading, err := r.db.Qualify(tabelSpreadingPolis)
	if err != nil {
		return nil, err
	}
	retro, err := r.db.Qualify(tabelSpreadingRetro)
	if err != nil {
		return nil, err
	}
	q1 := sqlDetailPolisPeserta(detail)
	if err := db.PeriksaSQL(q1); err != nil {
		return nil, err
	}
	var detailID string
	var mataUang sql.NullString
	err = r.db.QueryRowContext(ctx, q1, idPega, sertifikat, edmStatusHapus).Scan(&detailID, &mataUang)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository: membaca peserta polis untuk spreading: %w", err)
	}
	q2 := sqlSpreadingPolisPeserta(spreading, retro)
	if err := db.PeriksaSQL(q2); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q2, detailID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca spreading peserta polis: %w", err)
	}
	defer func() { _ = rows.Close() }()
	kurs := mataUang.String
	var out []models.Spreading
	terakhir := ""
	for rows.Next() {
		var id string
		var k [11]sql.NullString
		if err := rows.Scan(&id, &k[0], &k[1], &k[2], &k[3], &k[4], &k[5], &k[6], &k[7], &k[8], &k[9], &k[10]); err != nil {
			return nil, fmt.Errorf("repository: membaca satu baris spreading polis: %w", err)
		}
		if id != terakhir || len(out) == 0 {
			terakhir = id
			bagian, err := uraiUang(id, "RETROCADED_SHARE", k[2], kurs)
			if err != nil {
				return nil, err
			}
			out = append(out, models.Spreading{ID: id, TreatyTypeName: k[0].String, TreatyYearLife: k[1].String,
				RetrocadedShare: bagian, Currency: kurs})
		}
		if !k[3].Valid {
			continue // lapisan tanpa reinsurer
		}
		l := &out[len(out)-1]
		rt := models.SpreadingRetro{ReinsurerName: k[3].String}
		for _, x := range []struct {
			nama string
			v    sql.NullString
			ke   *uang.Ratio
		}{{"PERCENT_SHARE", k[4], &rt.PercentShare}, {"RATE", k[6], &rt.Rate}, {"COMMISION", k[7], &rt.Commision},
			{"OVR_COMM", k[8], &rt.OvrComm}} {
			v, err := uraiRasio(id, x.nama, x.v)
			if err != nil {
				return nil, err
			}
			*x.ke = v
		}
		for _, x := range []struct {
			nama string
			v    sql.NullString
			ke   *uang.Money
		}{{"AMOUNT", k[5], &rt.Amount}, {"PREMIUM_SPREADED_GROSS", k[9], &rt.PremiumSpreadedGross},
			{"PREMIUM_SPREADED_NET", k[10], &rt.PremiumSpreadedNet}} {
			v, err := uraiUang(id, x.nama, x.v, kurs)
			if err != nil {
				return nil, err
			}
			*x.ke = v
		}
		l.Retro = append(l.Retro, rt)
	}
	return out, rows.Err()
}
