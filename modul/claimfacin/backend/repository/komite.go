package repository

// Untuk apa berkas ini: BATAS KOMITE di repository Claim Fac In (pola Claim Prop) - penulis kasus komite KMT- TT2
// (`CreateKMTNo_Act` 8-10, `pxAddChildWork` ASM-FW-GCNMFW-Work-Komite flow Komite_Flow) dan pembacanya untuk grid
// "Committee Accept Status".
//
// ⛔ Satu-satunya berkas repository Claim Fac In yang menyebut tabel tangga komite. Penjaga batas Claim Life
// `komite_statik_test.go` mengecualikan berkas ini (izin work owner 09-10-2026). Keputusan anggota tetap DITULIS
// konteks Komite (tahap 2); berkas ini hanya menulis tangga awal (approval menunggu) dan membacanya.

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimfacin/backend/models"
)

// tabelTanggaKomite - tangga anggota kasus komite (tabel bersama Komite Claim Life).
const tabelTanggaKomite = "T_KOMITE_KOMITELIST"

// AnggotaTangga - satu anggota tangga komite yang ditulis.
type AnggotaTangga struct {
	Urut                       int
	OperatorID, Jabatan, Email string
}

// sqlSisipKasusKomite / sqlSisipKepalaKomite / sqlSisipAnggotaKomite - kelahiran kasus komite TKMT-.
func sqlSisipKasusKomite(work string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, COVER_KEY, LINI, TAHAP, POSITION, CREATE_OP, CREATE_OP_NAME, TGL_CREATE,
		TGL_UPDATE) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9)`, work)
}

// posisiAwal - T_WORK_CLAIM.POSITION kasus komite yang lahir = OPERATOR_ID anggota tangga tingkat pertama (prompt §6
// butir 1), seperti POSITION klaim = pemegang tahapnya. Modul Komite (tahap 2) memajukannya per tingkat.
func posisiAwal(anggota []AnggotaTangga) string {
	posisi, urut := "", 0
	for _, a := range anggota {
		if posisi == "" || a.Urut < urut {
			posisi, urut = a.OperatorID, a.Urut
		}
	}
	return posisi
}

func sqlSisipKepalaKomite(gen string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, ADJUSTMENT_ID, KOMITE_LOOP, KOMITE_COUNT) VALUES (:1, :2, :3, :4)`, gen)
}

// sqlSisipKepalaKomiteTutup - kepala kasus komite TT3 / TT4 (SendRejectClaimToKomite2 / SendCloseClaimToKomite 7.3):
// tanpa adjustment (ADJUSTMENT_ID NULL, migrasi komiteclaimfacin 641) dan ber-TRANSFER_TYPE (642). TT2 tetap
// `sqlSisipKepalaKomite` (DEFAULT '2').
func sqlSisipKepalaKomiteTutup(gen string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, ADJUSTMENT_ID, KOMITE_LOOP, KOMITE_COUNT, TRANSFER_TYPE)
		VALUES (:1, NULL, :2, :3, :4)`, gen)
}

// sqlKomiteTutupTerbuka - kasus komite TT3 / TT4 klaim `:1` yang masih menunggu (LINI FACIN, belum selesai).
func sqlKomiteTutupTerbuka(work, gen string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s w JOIN %s g ON g.ID = w.ID
		 WHERE w.COVER_KEY = :1 AND w.LINI = :2 AND w.TAHAP = :3 AND w.STATUS_WORK IS NULL
		   AND g.TRANSFER_TYPE IN (:4, :5)`, work, gen)
}

func sqlSisipAnggotaKomite(list string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, DATA_KOMITE_ID, KOMITE_URUT, KOMITE_OPERATORID, KOMITE_JABATAN, KOMITE_EMAIL,
		KOMITE_APPROVAL) VALUES (:1, :2, :3, :4, :5, :6, :7)`, list)
}

// sqlTanggaKomite - `ComiteeClaim` baris adjustment: anggota tangga semua kasus komite baris adjustment yang sama
// (kasus `:1` dan kasus terdahulunya - klaim, adjustment, dan LINI sama), urut kasus lalu jenjang (pola Claim Prop).
func sqlTanggaKomite(list, gen, work string) string {
	return fmt.Sprintf(`SELECT l.ID, l.KOMITE_OPERATORID, l.KOMITE_JABATAN, l.KOMITE_APPROVAL, l.KOMITE_COMMENT, l.KOMITE_EMAIL, %s
		  FROM %s l
		  JOIN %s g ON g.ID = l.DATA_KOMITE_ID
		  JOIN %s w ON w.ID = g.ID
		  JOIN %s g0 ON g0.ID = :1
		  JOIN %s w0 ON w0.ID = g0.ID
		 WHERE g.ADJUSTMENT_ID = g0.ADJUSTMENT_ID AND w.COVER_KEY = w0.COVER_KEY AND w.LINI = w0.LINI
		 ORDER BY w.TGL_CREATE, g.ID, l.KOMITE_URUT, l.ID`, fmt.Sprintf(db.FmtTanggalOracle, "l.DATE_APPROVE"), list, gen, work,
		gen, work)
}

// sqlSetelKomiteAdjustment - penautan baris adjustment ke kasus komitenya (KOMITE_ID UNIQUE): penyerahan pertama
// (KOMITE_ID kosong) atau penyerahan ulang dari kasus komite yang dibaca (`:3`).
func sqlSetelKomiteAdjustment(adj string) string {
	return fmt.Sprintf(`UPDATE %s SET KOMITE_ID = :1 WHERE ID = :2 AND (KOMITE_ID IS NULL OR KOMITE_ID = :3)`, adj)
}

// BuatKasusKomite = CreateKMTNo_Act 8-10 (TT2, `transfer` "2" / kosong) dan SendRejectClaimToKomite2 /
// SendCloseClaimToKomite 7.1-7.4 (TT3 / TT4, `adjID` kosong): T_WORK_CLAIM KMT- (COVER_KEY = klaim, LINI FACIN, TAHAP
// Komite_Flow, POSITION = `posisiAwal`), T_GENERAL_KOMITE (ADJUSTMENT_ID, KOMITE_LOOP = cacah tangga, KOMITE_COUNT 1,
// TRANSFER_TYPE TT3 / TT4), tangga satu baris per anggota (approval menunggu). Mengembalikan ID kasus komite.
func (g *Gudang) BuatKasusKomite(ctx context.Context, tx *db.Tx, klaimID, adjID, transfer, pembuat, namaPembuat string,
	anggota []AnggotaTangga, saat time.Time) (string, error) {
	tutup := transfer == models.TransferTolak || transfer == models.TransferTutup
	if tutup != (adjID == "") {
		return "", fmt.Errorf("repository: kasus komite TT %q dengan adjustment %q tidak sah", transfer, adjID)
	}
	if len(anggota) == 0 {
		return "", fmt.Errorf("repository: tangga komite kosong; kasus tanpa anggota tidak dapat diputuskan siapa pun")
	}
	id, err := g.IDKasusBerikut(ctx, tx, models.AwalanKomite)
	if err != nil {
		return "", err
	}
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return "", err
	}
	q := sqlSisipKasusKomite(work)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	hasil, err := tx.ExecContext(ctx, q, id, klaimID, models.LiniFacIn, models.TahapKomite,
		teksAtauNil(posisiAwal(anggota)), teksAtauNil(pembuat), teksAtauNil(namaPembuat), saat, saat)
	if err != nil {
		return "", fmt.Errorf("repository: melahirkan kasus komite: %w", err)
	}
	if err := db.PastikanSatuBaris(hasil, "kelahiran kasus komite"); err != nil {
		return "", err
	}
	gen, err := g.db.Qualify("T_GENERAL_KOMITE")
	if err != nil {
		return "", err
	}
	q, args := sqlSisipKepalaKomite(gen), []any{id, adjID, len(anggota), 1}
	if tutup {
		q, args = sqlSisipKepalaKomiteTutup(gen), []any{id, len(anggota), 1, transfer}
	}
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	if hasil, err = tx.ExecContext(ctx, q, args...); err != nil {
		return "", fmt.Errorf("repository: melahirkan baris komite: %w", err)
	}
	if err := db.PastikanSatuBaris(hasil, "kelahiran baris komite"); err != nil {
		return "", err
	}
	list, err := g.db.Qualify(tabelTanggaKomite)
	if err != nil {
		return "", err
	}
	q = sqlSisipAnggotaKomite(list)
	if err := db.PeriksaSQL(q); err != nil {
		return "", err
	}
	for _, a := range anggota {
		rid, err := g.db.NomorBerikut(ctx, tx, "SEQ_KOMITE_KOMITELIST")
		if err != nil {
			return "", err
		}
		hasil, err := tx.ExecContext(ctx, q, rid, id, a.Urut, teksAtauNil(a.OperatorID), teksAtauNil(a.Jabatan),
			teksAtauNil(a.Email), models.ApprovalKomiteMenunggu)
		if err != nil {
			return "", fmt.Errorf("repository: menulis anggota tangga komite: %w", err)
		}
		if err := db.PastikanSatuBaris(hasil, "anggota tangga komite"); err != nil {
			return "", err
		}
	}
	return id, nil
}

// SetelKomiteAdjustment menautkan baris adjustment ke kasus komitenya (T_CLAIM_ADJUSTMENT.KOMITE_ID, UNIQUE);
// `komiteLama` = kasus komite yang tertaut saat dibaca (kosong pada penyerahan pertama).
func (g *Gudang) SetelKomiteAdjustment(ctx context.Context, tx *db.Tx, adjID, komiteID, komiteLama string) error {
	tabel, err := g.db.Qualify(models.TabelAdjustment.Nama)
	if err != nil {
		return err
	}
	q := sqlSetelKomiteAdjustment(tabel)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, komiteID, adjID, teksAtauNil(komiteLama))
	if err != nil {
		return fmt.Errorf("repository: menautkan adjustment ke komite: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "penautan adjustment ke komite")
}

// TanggaKomite - anggota tangga kasus komite (beserta kasus terdahulu baris adjustment yang sama) untuk grid
// "Committee Accept Status".
func (a *Acuan) TanggaKomite(ctx context.Context, komiteID string) ([]models.AnggotaKomite, error) {
	t, err := a.q(tabelTanggaKomite)
	if err != nil {
		return nil, err
	}
	gen, err := a.q("T_GENERAL_KOMITE")
	if err != nil {
		return nil, err
	}
	work, err := a.q("T_WORK_CLAIM")
	if err != nil {
		return nil, err
	}
	rows, err := a.banyak(ctx, sqlTanggaKomite(t, gen, work), 7, komiteID)
	if err != nil {
		return nil, err
	}
	out := []models.AnggotaKomite{}
	for _, r := range rows {
		out = append(out, models.AnggotaKomite{ID: r[0], OperatorID: r[1], Jabatan: r[2], Approval: r[3], Comment: r[4],
			Email: r[5], TanggalSetuju: r[6]})
	}
	return out, nil
}

// sqlUbahAdjustmentKomite - UPDATE kolom milik komite (`Kolom.MilikKomite`, yang dilewati simpan halaman klaim) baris
// adjustment `ID` + `CLAIM_ID`; `kolom` urut nama properti (penampung urut kemunculan = urut argumen).
func sqlUbahAdjustmentKomite(tabel string, kolom []models.Kolom) string {
	var set []string
	n := 1
	for _, k := range kolom {
		e, jml := ekspresiTulis(k, n)
		set = append(set, k.Kolom+" = "+e)
		n += jml
	}
	return fmt.Sprintf(`UPDATE %s SET %s WHERE ID = :%d AND CLAIM_ID = :%d`, tabel, strings.Join(set, ", "), n, n+1)
}

// UbahAdjustmentKomite menulis kolom milik komite baris adjustment `adjID` klaim `klaimID` (`KomitePost_Adjustment`
// S7.2.1.5-S7.2.1.7, lewat kontrak `kontrak.KlaimFacInKomite`). Properti yang bukan kolom milik komite ditolak (kolom
// lain ditulis simpan halaman).
func (g *Gudang) UbahAdjustmentKomite(ctx context.Context, tx *db.Tx, klaimID, adjID string, nilai map[string]string) error {
	if len(nilai) == 0 {
		return nil
	}
	props := make([]string, 0, len(nilai))
	for p := range nilai {
		props = append(props, p)
	}
	sort.Strings(props)
	var kolom []models.Kolom
	var args []any
	for _, p := range props {
		k, ok := kolomKomiteAdjustment(p)
		if !ok {
			return fmt.Errorf("repository: %q bukan kolom milik komite baris adjustment", p)
		}
		v, err := nilaiTulis(k, nilai[p])
		if err != nil {
			return err
		}
		kolom = append(kolom, k)
		args = append(args, v...)
	}
	tabel, err := g.db.Qualify(models.TabelAdjustment.Nama)
	if err != nil {
		return err
	}
	q := sqlUbahAdjustmentKomite(tabel, kolom)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	hasil, err := tx.ExecContext(ctx, q, append(args, adjID, klaimID)...)
	if err != nil {
		return fmt.Errorf("repository: menulis keputusan komite baris adjustment: %w", err)
	}
	return db.PastikanSatuBaris(hasil, "keputusan komite baris adjustment")
}

// kolomKomiteAdjustment - kolom katalog milik komite berproperti `p`.
func kolomKomiteAdjustment(p string) (models.Kolom, bool) {
	for _, k := range models.TabelAdjustment.Kolom {
		if k.Properti == p && k.MilikKomite {
			return k, true
		}
	}
	return models.Kolom{}, false
}

// AdaKomiteTutupTerbuka - klaim `klaimID` punya kasus komite TT3 / TT4 yang masih menunggu (penjaga permintaan ganda,
// models.PesanTutupKomiteGanda).
func (g *Gudang) AdaKomiteTutupTerbuka(ctx context.Context, tx *db.Tx, klaimID string) (bool, error) {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return false, err
	}
	gen, err := g.db.Qualify("T_GENERAL_KOMITE")
	if err != nil {
		return false, err
	}
	q := sqlKomiteTutupTerbuka(work, gen)
	if err := db.PeriksaSQL(q); err != nil {
		return false, err
	}
	var n int
	if err := tx.QueryRowContext(ctx, q, klaimID, models.LiniFacIn, models.TahapKomite, models.TransferTolak,
		models.TransferTutup).Scan(&n); err != nil {
		return false, fmt.Errorf("repository: memeriksa kasus komite TT3 / TT4: %w", err)
	}
	return n > 0, nil
}
