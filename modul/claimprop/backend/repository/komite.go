package repository

// Untuk apa berkas ini: BATAS KOMITE di repository Claim Prop - penulis tangga kasus komite (`AddKomiteTreatyChild_ACT`
// langkah 22-25, `pxAddChildWork` ASM-FW-GCNMFW-Work-KomiteTreaty) dan pembacanya untuk grid "Committe Accept Status".
//
// ⛔ Satu-satunya berkas repository Claim Prop yang menyebut tabel tangga komite. Penjaga batas Claim Life
// `komite_statik_test.go` mengecualikan berkas ini (keputusan work owner 07-10-2026, opsi B). Keputusan anggota tetap
// DITULIS konteks Komite; berkas ini hanya menulis tangga awal (approval menunggu) dan membacanya.

import (
	"context"
	"fmt"
	"time"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimprop/backend/models"
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

// posisiAwal - T_WORK_CLAIM.POSITION kasus komite yang lahir = KomiteID anggota tangga tingkat pertama (workbasket
// tingkat 1 sejak tangga PROP ke workbasket, migrasi 537; keputusan work owner 09-10-2026), seperti POSITION klaim =
// workbasket pemegang tahapnya. Modul Komite memajukannya per tingkat.
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

// sqlSisipKepalaKomiteTutup - kepala kasus komite Close Without Payment (`SendCloseClaimToKomite` 5.3
// `childPageKomite.TransferType := 4`): tanpa adjustment (ADJUSTMENT_ID NULL) dan ber-TRANSFER_TYPE 4, teks pop-up
// Chronology (`childPageKomite.Komite.CircumtansesCouseOfLoss`, 5.2) - kolom bersama dari migrasi komiteclaimfacin
// 641 / 642 / 643 (izin work owner 10-10-2026). Extent Of Loss / Legal Liability tidak ada di pop-up Claim Prop.
func sqlSisipKepalaKomiteTutup(gen string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, ADJUSTMENT_ID, KOMITE_LOOP, KOMITE_COUNT, TRANSFER_TYPE,
		KOMITE_CIRCUM_CAUSE_OF_LOSS) VALUES (:1, NULL, :2, :3, :4, :5)`, gen)
}

// panjangTeksKomite - VARCHAR2(4000) kolom teks kepala kasus komite (migrasi komiteclaimfacin 643).
const panjangTeksKomite = 4000

// sqlKomiteTutupTerbuka - kasus komite Close Without Payment klaim `:1` yang masih menunggu (LINI PROP, belum selesai).
// Dikenali dari ADJUSTMENT_ID kosong (Claim Prop tanpa penulis TT 3), sehingga kueri tidak bergantung kolom
// TRANSFER_TYPE.
func sqlKomiteTutupTerbuka(work, gen string) string {
	return fmt.Sprintf(`SELECT COUNT(*) FROM %s w JOIN %s g ON g.ID = w.ID
		 WHERE w.COVER_KEY = :1 AND w.LINI = :2 AND w.TAHAP = :3 AND w.STATUS_WORK IS NULL AND g.ADJUSTMENT_ID IS NULL`,
		work, gen)
}

// sqlTutupKomiteTerbuka - `ASMForceCaseClose` / `pxForceCaseClose` `CloseAllSubCases true` (CloseClaimProp 10,
// KomitePost_Close S17): kasus komite klaim `:3` yang masih terbuka ikut Resolved-Completed tanpa pemegang, kecuali
// kasus `:6` (kasus komite yang sedang memutus; kosong = tanpa pengecualian - NVL karena `ID <> NULL` tidak pernah
// benar).
func sqlTutupKomiteTerbuka(work string) string {
	return fmt.Sprintf(`UPDATE %s SET STATUS_WORK = :1, TGL_UPDATE = :2, POSITION = NULL
		 WHERE COVER_KEY = :3 AND LINI = :4 AND TAHAP = :5 AND STATUS_WORK IS NULL AND ID <> NVL(:6, '-')`, work)
}

func sqlSisipAnggotaKomite(list string) string {
	return fmt.Sprintf(`INSERT INTO %s (ID, DATA_KOMITE_ID, KOMITE_URUT, KOMITE_OPERATORID, KOMITE_JABATAN, KOMITE_EMAIL,
		KOMITE_APPROVAL) VALUES (:1, :2, :3, :4, :5, :6, :7)`, list)
}

// sqlTanggaKomite - `ComiteeClaim` baris adjustment: anggota tangga SEMUA putaran komite baris adjustment yang sama
// (kasus `:1` dan kasus terdahulunya - klaim, adjustment, dan LINI sama), urut putaran lalu jenjang.
// AddKomiteTreatyChild_ACT S17 menghapus ComiteeClaim hanya pada penyerahan pertama; penyerahan ulang subjectivity
// MENAMBAHKAN anggota baru (S22.2) di belakang keputusan putaran terdahulu.
func sqlTanggaKomite(list, gen, work string) string {
	return fmt.Sprintf(`SELECT l.ID, l.KOMITE_OPERATORID, l.KOMITE_JABATAN, l.KOMITE_APPROVAL, l.KOMITE_COMMENT, %s
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
// (KOMITE_ID kosong) atau penyerahan ulang baris subjectivity dari kasus komite yang dibaca (`:3`).
func sqlSetelKomiteAdjustment(adj string) string {
	return fmt.Sprintf(`UPDATE %s SET KOMITE_ID = :1 WHERE ID = :2 AND (KOMITE_ID IS NULL OR KOMITE_ID = :3)`, adj)
}

// BuatKasusKomite = AddKomiteTreatyChild_ACT langkah 22-25: T_WORK_CLAIM TKMT- (COVER_KEY = klaim, LINI PROP, TAHAP
// KomiteTreaty_Flow, POSITION = `posisiAwal`), T_GENERAL_KOMITE (KOMITE_LOOP = cacah tangga, KOMITE_COUNT 1), tangga satu baris per anggota
// (approval menunggu). Mengembalikan ID kasus komite.
func (g *Gudang) BuatKasusKomite(ctx context.Context, tx *db.Tx, klaimID, adjID, pembuat, namaPembuat string,
	anggota []AnggotaTangga, saat time.Time) (string, error) {
	return g.lahirkanKomite(ctx, tx, klaimID, pembuat, namaPembuat, anggota, saat, func(gen, id string) (string, []any) {
		return sqlSisipKepalaKomite(gen), []any{id, adjID, len(anggota), 1}
	})
}

// BuatKasusKomiteTutup = SendCloseClaimToKomite langkah 5.1-5.4 (Close Without Payment): kasus komite TKMT- tanpa
// adjustment, TRANSFER_TYPE 4, satu tingkat (`KomiteLoop 1`, `KomiteCount 1`), teks Chronology pop-up. Mengembalikan
// ID kasus komite.
func (g *Gudang) BuatKasusKomiteTutup(ctx context.Context, tx *db.Tx, klaimID, kronologi, pembuat, namaPembuat string,
	anggota []AnggotaTangga, saat time.Time) (string, error) {
	kron, err := nilaiTulis(models.Kolom{Properti: models.PropKronologiTutup, Golongan: models.GolTeks,
		Panjang: panjangTeksKomite}, kronologi)
	if err != nil {
		return "", err
	}
	return g.lahirkanKomite(ctx, tx, klaimID, pembuat, namaPembuat, anggota, saat, func(gen, id string) (string, []any) {
		return sqlSisipKepalaKomiteTutup(gen), append([]any{id, len(anggota), 1, models.TransferTutup}, kron...)
	})
}

// lahirkanKomite - T_WORK_CLAIM TKMT- + kepala T_GENERAL_KOMITE (`kepala` merakit pernyataan dan nilainya) + tangga.
func (g *Gudang) lahirkanKomite(ctx context.Context, tx *db.Tx, klaimID, pembuat, namaPembuat string,
	anggota []AnggotaTangga, saat time.Time, kepala func(gen, id string) (string, []any)) (string, error) {
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
	hasil, err := tx.ExecContext(ctx, q, id, klaimID, models.LiniProp, models.TahapKomiteTreaty,
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
	q, args := kepala(gen, id)
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

// KomiteTutupTerbuka - cacah kasus komite Close Without Payment klaim `klaimID` yang masih menunggu keputusan.
func (g *Gudang) KomiteTutupTerbuka(ctx context.Context, tx *db.Tx, klaimID string) (int, error) {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return 0, err
	}
	gen, err := g.db.Qualify("T_GENERAL_KOMITE")
	if err != nil {
		return 0, err
	}
	q := sqlKomiteTutupTerbuka(work, gen)
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var n int
	if err := tx.QueryRowContext(ctx, q, klaimID, models.LiniProp, models.TahapKomiteTreaty).Scan(&n); err != nil {
		return 0, fmt.Errorf("repository: menghitung kasus komite close klaim: %w", err)
	}
	return n, nil
}

// TutupKomiteTerbuka = `CloseAllSubCases true`: kasus komite klaim `klaimID` yang masih terbuka ditutup Resolved-Completed,
// kecuali `kecuali` (kosong = semuanya). Nol baris bukan galat.
func (g *Gudang) TutupKomiteTerbuka(ctx context.Context, tx *db.Tx, klaimID, kecuali string, saat time.Time) error {
	work, err := g.db.Qualify("T_WORK_CLAIM")
	if err != nil {
		return err
	}
	q := sqlTutupKomiteTerbuka(work)
	if err := db.PeriksaSQL(q); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, q, models.StatusSelesai, saat, klaimID, models.LiniProp, models.TahapKomiteTreaty,
		teksAtauNil(kecuali)); err != nil {
		return fmt.Errorf("repository: menutup kasus komite terbuka klaim: %w", err)
	}
	return nil
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

// TanggaKomite - anggota tangga kasus komite (beserta putaran terdahulu baris adjustment yang sama) untuk grid
// "Committe Accept Status".
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
	rows, err := a.banyak(ctx, sqlTanggaKomite(t, gen, work), 6, komiteID)
	if err != nil {
		return nil, err
	}
	out := []models.AnggotaKomite{}
	for _, r := range rows {
		out = append(out, models.AnggotaKomite{ID: r[0], OperatorID: r[1], Jabatan: r[2], Approval: r[3], Comment: r[4],
			TanggalSetuju: r[5]})
	}
	return out, nil
}
