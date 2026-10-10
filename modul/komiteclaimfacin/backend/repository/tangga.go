package repository

// Untuk apa berkas ini: TANGGA KOMITE - satu-satunya berkas Komite Claim Fac In yang menyebut tabel tangga (penjaga batas
// Claim Life `komite_statik_test.go`, izin work owner 10-10-2026, pola Komite Claim Prop / Non Prop). Isinya:
//
//	daftar kerja      KomiteRouter: baris tangga tingkat berjalan (KOMITE_URUT = KOMITE_COUNT) yang menunggu, akun
//	                  pelaku atau workbasket aktifnya (Pega S6.1 = baris menunggu pertama; pada alur berurutan sama)
//	baca tangga       grid "List of Committee" + bahan wewenang
//	perluasan         ApprovalKomite_Act S7.1 (anggota DEGREE > 1 ditambahkan saat tingkat 1 Submit, KCF-02)
//	tulis keputusan   KomitePost_Adjustment S7.2.1.8 (`KomiteList(KomiteCount)` = `ComiteeClaim` adjustment: SATU baris
//	                  yang sama di sistem baru), S7.2.1.9.1.1 (tolak: sisa baris menunggu), KomitePost_Reject /
//	                  KomitePost_CloseClaim S6 / S5
//
// Kolom keputusan disebut lewat nama kolom basis data; nama properti Pega tidak dipakai di kode.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/komiteclaimfacin/backend/models"
)

// tabelTangga - tangga anggota kasus komite (tabel bersama Komite Claim Life, Prop, Non Prop, dan Fac In).
const tabelTangga = "T_KOMITE_KOMITELIST"

// sqlTangga - anggota tangga satu kasus, urut jenjang (`KomiteList` urut tambah = KOMITE_URUT).
func sqlTangga(list string) string {
	return fmt.Sprintf(`SELECT ID, KOMITE_URUT, KOMITE_OPERATORID, KOMITE_JABATAN, KOMITE_EMAIL, KOMITE_APPROVAL,
		       KOMITE_COMMENT, %s
		  FROM %s WHERE DATA_KOMITE_ID = :1 ORDER BY KOMITE_URUT, ID`, fmt.Sprintf(db.FmtTanggalOracle, "DATE_APPROVE"), list)
}

// BacaTangga membaca tangga kasus `id` (tx boleh nil).
func (g *Gudang) BacaTangga(ctx context.Context, tx *db.Tx, id string) ([]models.Anggota, error) {
	list, err := g.db.Qualify(tabelTangga)
	if err != nil {
		return nil, err
	}
	q := sqlTangga(list)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	var rows *sql.Rows
	if tx != nil {
		rows, err = tx.QueryContext(ctx, q, id)
	} else {
		rows, err = g.db.QueryContext(ctx, q, id)
	}
	if err != nil {
		return nil, fmt.Errorf("repository: membaca tangga komite: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []models.Anggota
	for rows.Next() {
		var n [7]sql.NullString
		var urut sql.NullInt64
		if err := rows.Scan(&n[0], &urut, &n[1], &n[2], &n[3], &n[4], &n[5], &n[6]); err != nil {
			return nil, fmt.Errorf("repository: memindai tangga komite: %w", err)
		}
		out = append(out, models.Anggota{ID: n[0].String, Urut: int(urut.Int64), OperatorID: n[1].String,
			Jabatan: n[2].String, Email: n[3].String, Keputusan: strings.TrimSpace(n[4].String), Komentar: n[5].String,
			Tanggal: n[6].String})
	}
	return out, rows.Err()
}

// sqlTambahAnggota - satu anggota tangga baru (perluasan S7.1; approval menunggu, tanpa komentar).
func sqlTambahAnggota(list string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, DATA_KOMITE_ID, KOMITE_URUT, KOMITE_OPERATORID, KOMITE_JABATAN, KOMITE_EMAIL,
		KOMITE_APPROVAL) VALUES (:1, :2, :3, :4, :5, :6, :7)`, list)
}

// TambahAnggota menyisipkan anggota perluasan tangga kasus `id` (ID dari SEQ_KOMITE_KOMITELIST, pola kelahiran Claim
// Fac In). Mengembalikan anggota beserta ID-nya.
func (g *Gudang) TambahAnggota(ctx context.Context, tx *db.Tx, id string, baru []models.Anggota) ([]models.Anggota, error) {
	if err := wajibTx(tx); err != nil {
		return nil, err
	}
	list, err := g.db.Qualify(tabelTangga)
	if err != nil {
		return nil, err
	}
	q := sqlTambahAnggota(list)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	out := make([]models.Anggota, 0, len(baru))
	for _, a := range baru {
		rid, err := g.db.NomorBerikut(ctx, tx, "SEQ_KOMITE_KOMITELIST")
		if err != nil {
			return nil, err
		}
		h, err := tx.ExecContext(ctx, q, rid, id, a.Urut, db.KosongJadiNil(a.OperatorID), db.KosongJadiNil(a.Jabatan),
			db.KosongJadiNil(a.Email), models.KeputusanMenunggu)
		if err != nil {
			return nil, fmt.Errorf("repository: menambah anggota tangga komite: %w", err)
		}
		if err := db.PastikanSatuBaris(h, "anggota tangga komite"); err != nil {
			return nil, err
		}
		a.ID = rid
		out = append(out, a)
	}
	return out, nil
}

// sqlTulisAnggota - keputusan satu baris tangga; bersyarat masih menunggu (dua klik tidak sama-sama menang). Baris yang
// diputus pelaku (`komentar`) menyimpan akun pemutusnya di KOMITE_OPERATORID (KomiteID workbasket -> akun yang memutus,
// pola Komite Claim Prop 09-10-2026); pemutus kosong = tidak ditimpa.
func sqlTulisAnggota(list string, komentar bool) string {
	if komentar {
		return fmt.Sprintf(`UPDATE %s SET KOMITE_APPROVAL = :1, KOMITE_COMMENT = :2, DATE_APPROVE = :3,
			       KOMITE_OPERATORID = NVL(:4, KOMITE_OPERATORID)
			 WHERE ID = :5 AND DATA_KOMITE_ID = :6 AND KOMITE_APPROVAL = :7`, list)
	}
	return fmt.Sprintf(`UPDATE %s SET KOMITE_APPROVAL = :1, DATE_APPROVE = :2
		 WHERE ID = :3 AND DATA_KOMITE_ID = :4 AND KOMITE_APPROVAL = :5`, list)
}

// TulisAnggota menulis satu ubahan tangga kasus `id` (baris harus masih menunggu).
func (g *Gudang) TulisAnggota(ctx context.Context, tx *db.Tx, id string, u models.UbahAnggota) error {
	if err := wajibTx(tx); err != nil {
		return err
	}
	list, err := g.db.Qualify(tabelTangga)
	if err != nil {
		return err
	}
	q := sqlTulisAnggota(list, u.IsiKomentar)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	args := []any{u.Keputusan, u.Tanggal, u.ID, id, models.KeputusanMenunggu}
	if u.IsiKomentar {
		args = []any{u.Keputusan, db.KosongJadiNil(u.Komentar), u.Tanggal, db.KosongJadiNil(u.Pemutus), u.ID, id,
			models.KeputusanMenunggu}
	}
	h, err := tx.ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("repository: menulis keputusan tangga komite: %w", err)
	}
	return satuAtauBersamaan(h, "baris tangga "+u.ID)
}

// sqlDariKerja - FROM daftar kerja: kasus, work, tangga, klaim induk dan baris adjustment TT2 (LEFT JOIN: kasus TT3 /
// TT4 tanpa adjustment dan kasus yang induknya hilang tetap terlihat).
func sqlDariKerja(gen, work, list, klaim, adj string) string {
	return fmt.Sprintf(` FROM %s g
		  JOIN %s w ON w.ID = g.ID
		  JOIN %s l ON l.DATA_KOMITE_ID = g.ID
		  LEFT JOIN %s c ON c.ID = w.COVER_KEY
		  LEFT JOIN %s a ON a.KOMITE_ID = g.ID`, gen, work, list, klaim, adj)
}

// sqlSaringKerja - baris tangga tingkat berjalan (KOMITE_URUT = KOMITE_COUNT) yang masih menunggu, ber-KomiteID akun
// pelaku ATAU salah satu dari `nPeran` workbasket aktifnya; kasus terbuka; LINI FACIN ketat dan awalan KMT-. Bind urut
// kemunculan: akun, peran..., menunggu, LINI, awalan.
func sqlSaringKerja(nPeran int) string {
	pemegang := "l.KOMITE_OPERATORID = :1"
	if nPeran > 0 {
		ph := make([]string, nPeran)
		for i := range ph {
			ph[i] = fmt.Sprintf(":%d", i+2)
		}
		pemegang = "(l.KOMITE_OPERATORID = :1 OR l.KOMITE_OPERATORID IN (" + strings.Join(ph, ", ") + "))"
	}
	n := nPeran + 2
	return fmt.Sprintf(`
		 WHERE %s
		   AND l.KOMITE_APPROVAL = :%d
		   AND l.KOMITE_URUT = g.KOMITE_COUNT
		   AND w.STATUS_WORK IS NULL
		   AND w.LINI = :%d
		   AND w.ID LIKE :%d`, pemegang, n, n+1, n+2)
}

// sqlDaftarKerja - daftar kerja satu penyetuju (`nPeran` workbasket aktif).
func sqlDaftarKerja(gen, work, list, klaim, adj string, nPeran int) string {
	return `SELECT g.ID, w.COVER_KEY, c.CLAIM_NO, g.TRANSFER_TYPE, l.KOMITE_URUT, g.KOMITE_COUNT, g.KOMITE_LOOP,
		       l.KOMITE_JABATAN, a.ACCEPTANCE_STATUS, w.STATUS_WORK, ` + fmt.Sprintf(db.FmtTanggalOracle, "w.TGL_UPDATE") +
		sqlDariKerja(gen, work, list, klaim, adj) + sqlSaringKerja(nPeran) + `
		 ORDER BY w.TGL_UPDATE DESC, g.ID DESC`
}

func (g *Gudang) tabelKerja() (gen, work, list, klaim, adj string, err error) {
	for _, p := range []struct {
		nama string
		ke   *string
	}{{"T_GENERAL_KOMITE", &gen}, {"T_WORK_CLAIM", &work}, {tabelTangga, &list}, {"T_GENERAL_CLAIM", &klaim},
		{"T_CLAIM_ADJUSTMENT", &adj}} {
		if *p.ke, err = g.db.Qualify(p.nama); err != nil {
			return
		}
	}
	return
}

// DaftarKerja membaca daftar kerja penyetuju `akun` pemegang workbasket aktif `peran` (sudah `models.PeranKerja`).
func (g *Gudang) DaftarKerja(ctx context.Context, akun string, peran []string) ([]models.BarisKerja, error) {
	gen, work, list, klaim, adj, err := g.tabelKerja()
	if err != nil {
		return nil, err
	}
	q := sqlDaftarKerja(gen, work, list, klaim, adj, len(peran))
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	args := []any{akun}
	for _, p := range peran {
		args = append(args, p)
	}
	args = append(args, models.KeputusanMenunggu, models.LiniFacIn, awalanLike())
	rows, err := g.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca daftar kerja komite: %w", err)
	}
	defer func() { _ = rows.Close() }()
	out := []models.BarisKerja{}
	for rows.Next() {
		var n [8]sql.NullString
		var urut, count, loop sql.NullInt64
		if err := rows.Scan(&n[0], &n[1], &n[2], &n[3], &urut, &count, &loop, &n[4], &n[5], &n[6], &n[7]); err != nil {
			return nil, fmt.Errorf("repository: memindai daftar kerja komite: %w", err)
		}
		out = append(out, models.BarisKerja{KasusID: n[0].String, KlaimID: n[1].String, NoKlaim: n[2].String,
			Jenis: models.JudulTransfer[strings.TrimSpace(n[3].String)], Tingkat: int(urut.Int64), Count: int(count.Int64),
			Loop: int(loop.Int64), Jabatan: n[4].String, StatusBaris: strings.TrimSpace(n[5].String),
			StatusWork: n[6].String, TglUpdate: waktuDB(n[7].String)})
	}
	return out, rows.Err()
}
