package repository

// Layar Outstanding dirombak mengikuti Pega - 10-10-2026 (permintaan work owner "kenapa tidak ada detail peserta nya
// seperti yang di pega?").
//
// Untuk apa berkas ini: dua hal yang hanya Oracle dapat jawab untuk layar itu -
//   - nama dan tanggal lahir peserta untuk grid Participant Details (InputOSClaimLife b17203, b17380), dibaca saat
//     baca dari baris sumbernya, TIDAK disimpan di tabel klaim;
//   - penulisan tombol `Save` b19413: Accept peserta (IS_ACCEPT, migrasi 026) dan medan tab Register Claim di
//     T_GENERAL_CLAIM (STS_KATASTROFE, IS_KPR, KATASTROFE_NOTE).
//
// Dibaca sesudah: kolompeserta.go, gandawarisan.go.

import (
	"context"
	"database/sql"
	"fmt"

	"nusantarare/inti/backend/db"
)

// Tertanggung adalah nama dan tanggal lahir (YYYY-MM-DD) satu peserta, dari baris sumbernya.
type Tertanggung struct {
	Nama         string
	TanggalLahir string
	// SEX dan PLAN baris sumber - rincian peserta layar Outstanding (11-10-2026).
	JenisKelamin string
	Plan         string
}

// sqlTertanggungPeserta - satu kueri untuk seluruh peserta aktif sebuah klaim; subkueri berbatas (ROWNUM = 1) dan
// ber-index per peserta (penjaga tabel 66,8 juta baris). Kunci baris sumber sama dengan
// sqlIsiTertanggungCermin: PL_NUMBER + CERTIFICATE_NO + ID (= SOURCE_ID). Peserta tanpa baris sumber tidak muncul -
// layar menampilkan tanda pisah. ⛔ Hanya MEMBACA; kolomSalin tetap tanpa nama (data pribadi tidak disalin).
func sqlTertanggungPeserta(pes, sumber string) string {
	return fmt.Sprintf(`SELECT p.ID,
	       (SELECT m.NAME_OF_INSURED FROM %s m
	         WHERE m.PL_NUMBER = p.PL_NUMBER AND m.CERTIFICATE_NO = p.CERTIFICATE_NO AND m.ID = p.SOURCE_ID AND ROWNUM = 1),
	       (SELECT TO_CHAR(m2.DOB, 'YYYY-MM-DD') FROM %s m2
	         WHERE m2.PL_NUMBER = p.PL_NUMBER AND m2.CERTIFICATE_NO = p.CERTIFICATE_NO AND m2.ID = p.SOURCE_ID AND ROWNUM = 1),
	       (SELECT m3.SEX FROM %s m3
	         WHERE m3.PL_NUMBER = p.PL_NUMBER AND m3.CERTIFICATE_NO = p.CERTIFICATE_NO AND m3.ID = p.SOURCE_ID AND ROWNUM = 1),
	       (SELECT m4.PLAN FROM %s m4
	         WHERE m4.PL_NUMBER = p.PL_NUMBER AND m4.CERTIFICATE_NO = p.CERTIFICATE_NO AND m4.ID = p.SOURCE_ID AND ROWNUM = 1)
	   FROM %s p
	  WHERE p.CLAIM_ID = :1 AND p.STS_HAPUS IS NULL`, sumber, sumber, sumber, sumber, pes)
}

// TertanggungPeserta membaca nama dan tanggal lahir peserta aktif sebuah klaim, per pengenal peserta.
func (r *KlaimLife) TertanggungPeserta(ctx context.Context, klaimID string) (map[string]Tertanggung, error) {
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return nil, err
	}
	sumber, err := r.db.Qualify(namaTabelPeserta)
	if err != nil {
		return nil, err
	}
	q := sqlTertanggungPeserta(pes, sumber)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, q, klaimID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca tertanggung peserta: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := map[string]Tertanggung{}
	for rows.Next() {
		var id string
		var nama, lahir, sex, plan sql.NullString
		if err := rows.Scan(&id, &nama, &lahir, &sex, &plan); err != nil {
			return nil, fmt.Errorf("repository: membaca satu tertanggung peserta: %w", err)
		}
		out[id] = Tertanggung{Nama: nama.String, TanggalLahir: lahir.String, JenisKelamin: sex.String, Plan: plan.String}
	}
	return out, rows.Err()
}

// sqlSetelAcceptPeserta - checkbox Accept satu peserta (IS_ACCEPT, migrasi 026). Dikunci ID + CLAIM_ID + aktif.
func sqlSetelAcceptPeserta(pes string) string {
	return fmt.Sprintf(`UPDATE %s SET IS_ACCEPT = :1 WHERE ID = :2 AND CLAIM_ID = :3 AND STS_HAPUS IS NULL`, pes)
}

// SetelAcceptPeserta menulis IS_ACCEPT satu peserta di transaksi pemanggil; tepat satu baris.
func (r *KlaimLife) SetelAcceptPeserta(ctx context.Context, tx *db.Tx, klaimID, pesertaID, nilai string) error {
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return err
	}
	q := sqlSetelAcceptPeserta(pes)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(nilai), pesertaID, klaimID)
	if err != nil {
		return fmt.Errorf("repository: menulis Accept peserta %s: %w", pesertaID, err)
	}
	return db.PastikanSatuBaris(hasil, "Accept peserta "+pesertaID)
}

// sqlSetelKatastrofe - tiga medan tab Register Claim, HANYA selama klaim belum bernomor: StsKatastrofe
// `pyDisabledWhen CLAIM_NO !=""` b3515 dan KatastrofeNote `pyReadOnlyCondition CLAIM_NO !=""` b5059. Syarat
// `CLAIM_NO IS NULL` di DALAM pernyataan: Save to RNM yang menyela tidak dapat diikuti penulisan yang sudah terlarang.
func sqlSetelKatastrofe(hdr string) string {
	return fmt.Sprintf(`UPDATE %s SET STS_KATASTROFE = :1, IS_KPR = :2, KATASTROFE_NOTE = :3
	 WHERE ID = :4 AND CLAIM_NO IS NULL`, hdr)
}

// sqlSetelKPR - Is KPR tanpa syarat bernomor (b3810 tanpa pyDisabledWhen / pyReadOnlyCondition).
func sqlSetelKPR(hdr string) string {
	return fmt.Sprintf(`UPDATE %s SET IS_KPR = :1 WHERE ID = :2`, hdr)
}

// SetelMedanRegisterClaim menulis medan tab Register Claim di transaksi pemanggil. `bernomor` = CLAIM_NO sudah
// terisi saat dibaca layanan: hanya IS_KPR yang ditulis. Mengembalikan false bila klaim ternyata sudah bernomor saat
// ditulis (nol baris pada varian belum-bernomor) - layanan menjawabnya 409.
func (r *KlaimLife) SetelMedanRegisterClaim(ctx context.Context, tx *db.Tx, klaimID string, bernomor bool,
	stsKatastrofe, isKPR, catatan string) (bool, error) {

	hdr, err := r.db.Qualify("T_GENERAL_CLAIM")
	if err != nil {
		return false, err
	}
	q, arg := sqlSetelKPR(hdr), []any{db.KosongJadiNil(isKPR), klaimID}
	if !bernomor {
		q = sqlSetelKatastrofe(hdr)
		arg = []any{db.KosongJadiNil(stsKatastrofe), db.KosongJadiNil(isKPR), db.KosongJadiNil(catatan), klaimID}
	}
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	hasil, err := tx.ExecContext(ctx, q, arg...)
	if err != nil {
		return false, fmt.Errorf("repository: menulis medan Register Claim %s: %w", klaimID, err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("repository: mencacah medan Register Claim: %w", err)
	}
	return n == 1, nil
}

// sqlSetelCatatanAdmin - ADMIN NOTES satu peserta (migrasi 027). Dibekukan bila peserta sudah diputus
// (`.STS_REJECT=='1' || =='2'` b2628) - syaratnya di DALAM pernyataan.
func sqlSetelCatatanAdmin(pes string) string {
	return fmt.Sprintf(`UPDATE %s SET ADMIN_NOTES = :1
	 WHERE ID = :2 AND CLAIM_ID = :3 AND STS_HAPUS IS NULL AND (STS_REJECT IS NULL OR STS_REJECT NOT IN ('1', '2'))`, pes)
}

// SetelCatatanAdmin menulis ADMIN_NOTES satu peserta di transaksi pemanggil; false bila peserta sudah diputus.
func (r *KlaimLife) SetelCatatanAdmin(ctx context.Context, tx *db.Tx, klaimID, pesertaID, catatan string) (bool, error) {
	pes, err := r.db.Qualify("T_CLAIMLF_PREMIUMLIST_DETAIL")
	if err != nil {
		return false, err
	}
	q := sqlSetelCatatanAdmin(pes)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	hasil, err := tx.ExecContext(ctx, q, db.KosongJadiNil(catatan), pesertaID, klaimID)
	if err != nil {
		return false, fmt.Errorf("repository: menulis Admin Notes peserta %s: %w", pesertaID, err)
	}
	n, err := hasil.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("repository: mencacah Admin Notes: %w", err)
	}
	return n == 1, nil
}
