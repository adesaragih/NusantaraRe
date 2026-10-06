package repository

// Untuk apa berkas ini: KASUS - baris `T_WORK_POLIS` (tabel kasus lintas-lini,
// premiumlistlife 050/059) dan baris generasi `T_GENERAL_POLIS_TREATY` yang berbagi
// kunci utama dengannya (spec-penyimpanan ID-7). Padanan work object
// `ASM-FW-GISFW-Work-NB` Pega: `createWork` tombol Create portal, perpindahan
// assignment `Flow/InputRealizationTreatyIn`, dan daftar portal
// `Section/SFAPortal_OpportunitiesList`.
//
// ⛔ Setiap ubah tahap MENYERTAKAN tahap lama di WHERE (pola premiumlistlife
// `sqlPindahTahapPolis`): dua permintaan serentak tidak dapat sama-sama
// memindahkan kasus yang sama.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
)

// IDKasusBerikut menerbitkan pengenal `NB-<n>` dari SEQ_WORK_POLIS.
func (g *Gudang) IDKasusBerikut(ctx context.Context, tx *db.Tx) (string, error) {
	urut, err := g.db.NomorBerikut(ctx, tx, sequenceKerjaPolis)
	if err != nil {
		return "", err
	}
	return models.RakitIDKasus(urut), nil
}

func sqlSisipKerja(t string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, LINI, POSITION, STATUS_WORK, FLAG_ONGOING_POLICY, COVER_KEY, CREATE_OP, CREATE_OP_NAME, TGL_CREATE, TGL_UPDATE)
		VALUES (:1, :2, :3, :4, :5, NULL, :6, :7, SYSDATE, SYSDATE)`, t)
}

func sqlSisipGenerasi(t string) string {
	// IDPEGA dibuang 06-10-2026 (keputusan work owner: sama dengan ID)
	return fmt.Sprintf(`INSERT INTO %s (ID, PRODKE, TGL_INPUT, USERNAME, POSITION_NOTE)
		VALUES (:1, 0, SYSDATE, :2, :3)`, t)
}

// SisipKasus melahirkan kasus di posisi admin - connector `Start1 ->
// Assignment2`: `.Position = 4; .FlagOnGoingPolicy = 1; .PositionNote =
// "ReasTreatyInAdmin"`. Generasi NB: PRODKE 0, OLD_POLIS_ID kosong (ID-8, ID-9).
//
// `pembuat` = identitas akses login (CREATE_OP); `namaPembuat` = nama tampilan
// dari M_LOGIN_GO (`pxCreateOpName`), dipakai `NBStatus` DT pasca (P40).
func (g *Gudang) SisipKasus(ctx context.Context, tx *db.Tx, id, pembuat, namaPembuat string) error {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return err
	}
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, "menyisipkan kasus", sqlSisipKerja(kerja),
		id, models.LiniKasus, models.PositionAdmin, models.AssignmentAdmin, "1",
		db.KosongJadiNil(pembuat), db.KosongJadiNil(namaPembuat))
	if err != nil {
		return err
	}
	if err := db.PastikanSatuBaris(hasil, "penyisipan kasus"); err != nil {
		return err
	}
	hasil, err = jalankan(ctx, tx, "menyisipkan generasi polis", sqlSisipGenerasi(gen),
		id, db.KosongJadiNil(pembuat), models.PosisiAdmin)
	if err != nil {
		return err
	}
	return db.PastikanSatuBaris(hasil, "penyisipan generasi polis")
}

func sqlKeadaan(kerja, gen string) string {
	return fmt.Sprintf(`SELECT w.ID, w.POSITION, w.STATUS_WORK, g.POSITION_NOTE, g.NOPOLIS,
	        CASE WHEN %s THEN 0 ELSE 1 END, w.CREATE_OP,
	        TO_CHAR(w.TGL_CREATE, '%s')
	   FROM %s w JOIN %s g ON g.ID = w.ID
	  WHERE w.ID = :1`, syaratTerbuka(gen), fmtTanggal, kerja, gen)
}

// Keadaan membaca keadaan kerja satu kasus. Di dalam transaksi bila `tx`
// terisi, supaya keputusan dan tulisannya melihat baris yang sama.
func (g *Gudang) Keadaan(ctx context.Context, tx *db.Tx, id string) (models.Kasus, error) {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return models.Kasus{}, err
	}
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return models.Kasus{}, err
	}
	q := sqlKeadaan(kerja, gen)
	if err := db.PeriksaSQL(q); err != nil {
		return models.Kasus{}, err
	}
	var k models.Kasus
	var pos, status, catatan, nopol, op, lahir sql.NullString
	var tutup int
	err = g.pembaca(tx).QueryRowContext(ctx, q, id).Scan(&k.ID, &pos, &status, &catatan, &nopol, &tutup, &op, &lahir)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Kasus{}, ErrKasusTidakAda
	}
	if err != nil {
		return models.Kasus{}, fmt.Errorf("repository: membaca keadaan kasus: %w", err)
	}
	k.Position, k.StatusWork, k.PositionNote = teks(pos), teks(status), teks(catatan)
	k.NoPolis, k.CreateOp, k.TglCreate = teks(nopol), teks(op), teks(lahir)
	k.GenerasiTertutup = tutup == 1
	return k, nil
}

// KunciKasus mengunci baris kasus (`FOR UPDATE`) di transaksi pemanggil dan
// memastikan tahapnya masih `statusHarap` - tahap yang DIBACA sebelum
// transaksi. Setiap tindakan tulis memanggilnya lebih dulu: tanpa itu simpan
// draf yang berbalapan dengan submit dapat menulis halaman kasus yang sudah
// berpindah posisi, atau menerbitkan nomor polis bagi kasus yang baru
// ditolak.
func (g *Gudang) KunciKasus(ctx context.Context, tx *db.Tx, id, statusHarap string) error {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`SELECT STATUS_WORK FROM %s WHERE ID = :1 FOR UPDATE`, kerja)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	var status sql.NullString
	err = tx.QueryRowContext(ctx, q, id).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrKasusTidakAda
	}
	if err != nil {
		return fmt.Errorf("repository: mengunci kasus: %w", err)
	}
	if status.String != statusHarap {
		return ErrTahapBerubah
	}
	return nil
}

func sqlPindahKerja(t string) string {
	return fmt.Sprintf(`UPDATE %s SET POSITION = :1, STATUS_WORK = :2, TGL_UPDATE = SYSDATE
		  WHERE ID = :3 AND STATUS_WORK = :4`, t)
}

func sqlPindahGenerasi(t string) string {
	return fmt.Sprintf(`UPDATE %s g SET POSITION_NOTE = :1 WHERE g.ID = :2 AND %s`, t, syaratTerbuka(t))
}

// PindahPosisi memindahkan kasus ke posisi lain - tugas properti connector
// (`pyWorkPage.Position`, `PositionNote`) dan assignment tujuannya.
func (g *Gudang) PindahPosisi(ctx context.Context, tx *db.Tx, id, statusLama, posisiBaru string) error {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return err
	}
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, "memindahkan kasus", sqlPindahKerja(kerja),
		models.PositionPosisi(posisiBaru), models.AssignmentPosisi(posisiBaru), id, statusLama)
	if err != nil {
		return err
	}
	if n, _ := hasil.RowsAffected(); n != 1 {
		return ErrTahapBerubah
	}
	hasil, err = jalankan(ctx, tx, "memindahkan posisi generasi", sqlPindahGenerasi(gen), posisiBaru, id)
	if err != nil {
		return err
	}
	if n, _ := hasil.RowsAffected(); n != 1 {
		return ErrGenerasiTertutup
	}
	return nil
}

func sqlTutupKerja(t string) string {
	return fmt.Sprintf(`UPDATE %s SET STATUS_WORK = :1, POSITION = NULL, TGL_UPDATE = SYSDATE
		  WHERE ID = :2 AND STATUS_WORK = :3`, t)
}

// TutupKasus menyelesaikan kasus (End3) dengan status akhirnya. POSITION
// dikosongkan: kasus tertutup tidak menunggu di antrean mana pun (pola
// premiumlistlife `TutupKasus`).
func (g *Gudang) TutupKasus(ctx context.Context, tx *db.Tx, id, statusLama, statusAkhir string) error {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return err
	}
	hasil, err := jalankan(ctx, tx, "menutup kasus", sqlTutupKerja(kerja), statusAkhir, id, statusLama)
	if err != nil {
		return err
	}
	if n, _ := hasil.RowsAffected(); n != 1 {
		return ErrTahapBerubah
	}
	return nil
}

// sqlDaftarKasus merakit daftar portal beserta argumen ikatnya. Penampung
// UNIK - klausa pembatas baris memecah pengikatan penampung berulang (penjaga
// `TestNolPenampungBerulangDiSQLBerpembatasBaris`).
//
// RD `GetListOpportunity` (logika `B AND D AND E AND F AND A AND (C OR G) AND F1`):
//
//	D, F  `A.pyStatusWork` != Resolved-Completed / Resolved-Rejected
//	C, G  `.Name` / `.TextNoQuotation` Contains `Param.Search` (pyCaseInsensitive)
//	      - G = pengenal kasus `NB-<n>` (filter B `.TextNoQuotation Contains "NB-"`);
//	      ⛔ C tidak dibangun: `.Name` milik kelas CRM
//	      `ASM-FW-SFAGISFW-Work-Opportunity`, ditulis NOL rule korpus dan tak
//	      berkolom di diagram grilling (butir terbuka, bab 0 butir 11)
//	pyMaxRecords 500 -> `models.BatasDaftarPortal`
//
// `s.Antrean` = workbasket gerbang portal (services.DaftarKasus); kosong =
// tanpa batas antrean.
func sqlDaftarKasus(kerja, gen, quot string, s models.SaringanKasus) (string, []any) {
	var b strings.Builder
	args := []any{models.LiniKasus, models.StatusDitolak, models.StatusSelesai}
	pen := func(v any) string {
		args = append(args, v)
		return fmt.Sprintf(":%d", len(args))
	}
	fmt.Fprintf(&b, `SELECT w.ID, q.BUSINESS_NAME, q.INSURED_NAME, q.MARKETING_NAME, g.NB_STATUS,
	        w.STATUS_WORK, g.POSITION_NOTE, g.NOPOLIS, TO_CHAR(w.TGL_CREATE, '%s')
	   FROM %s w
	   JOIN %s g ON g.ID = w.ID
	   LEFT JOIN %s q ON q.POLIS_ID = g.ID
	  WHERE g.PRODKE = 0 AND w.LINI = :1
	    AND (w.STATUS_WORK IS NULL OR w.STATUS_WORK NOT IN (:2, :3))`, fmtTanggal, kerja, gen, quot)
	if cari := strings.TrimSpace(s.Cari); cari != "" {
		fmt.Fprintf(&b, `
	    AND UPPER(w.ID) LIKE %s`, pen("%"+strings.ToUpper(cari)+"%"))
	}
	if s.Posisi != "" {
		fmt.Fprintf(&b, `
	    AND g.POSITION_NOTE = %s`, pen(s.Posisi))
	}
	if len(s.Antrean) > 0 {
		var ps []string
		for _, a := range s.Antrean {
			ps = append(ps, pen(a))
		}
		fmt.Fprintf(&b, `
	    AND g.POSITION_NOTE IN (%s)`, strings.Join(ps, ", "))
	}
	fmt.Fprintf(&b, `
	  ORDER BY w.TGL_CREATE DESC
	  FETCH FIRST %d ROWS ONLY`, models.BatasDaftarPortal)
	return b.String(), args
}

// DaftarKasus membaca kasus terbuka untuk portal.
//
// Saringan dari RD `GetListOpportunity` (lihat `sqlDaftarKasus`).
// ⚠️ Logika filter RD itu juga memuat "A" (`pxCreateOperator =
// Param.UserIdentifier`, daftar per pembuat) - bertentangan dengan antrean
// bersama (AC 11, 92); yang dipakai: SEMUA kasus terbuka, opsional per posisi,
// dibatasi `s.Antrean` bila gerbang portal mengisinya (services.DaftarKasus).
func (g *Gudang) DaftarKasus(ctx context.Context, s models.SaringanKasus) ([]models.RingkasanKasus, error) {
	kerja, err := g.nama(tabelKerja)
	if err != nil {
		return nil, err
	}
	gen, err := g.nama(models.TabelGeneralPolis.Nama)
	if err != nil {
		return nil, err
	}
	quot, err := g.nama(models.TabelQuotation.Nama)
	if err != nil {
		return nil, err
	}
	q, args := sqlDaftarKasus(kerja, gen, quot, s)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca daftar kasus: %w", err)
	}
	defer rows.Close()
	var out []models.RingkasanKasus
	for rows.Next() {
		var r models.RingkasanKasus
		var bis, ins, mkt, nb, st, pn, np, tg sql.NullString
		if err := rows.Scan(&r.ID, &bis, &ins, &mkt, &nb, &st, &pn, &np, &tg); err != nil {
			return nil, fmt.Errorf("repository: membaca baris daftar kasus: %w", err)
		}
		r.BusinessName, r.InsuredName, r.MarketingName = teks(bis), teks(ins), teks(mkt)
		r.NBStatus, r.StatusWork, r.PositionNote, r.NoPolis, r.TglCreate = teks(nb), teks(st), teks(pn), teks(np), teks(tg)
		out = append(out, r)
	}
	return out, rows.Err()
}
