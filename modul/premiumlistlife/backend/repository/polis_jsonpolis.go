package repository

// JSON_POLIS - satu baris per polis yang di-Confirm (Resolved-Completed),
// keputusan work owner 06-10-2026.
//
// Untuk apa berkas ini: meniru `InsertJsonPolisLife_Act` langkah 10
// (`InsertJsonPolis` → `POOLDATA.INSERTJSONPOLISLIFE`) dengan SQL LANGSUNG di
// transaksi Confirm milik pemanggil, SESUDAH PL Number terbit dan LIFEINPRODUCTION
// (polis_produksi.go). Prosedurnya TIDAK dipanggil.
//
// `[terverifikasi]` nama kolom = ALL_TAB_COLUMNS POOLDATA (diberikan user
// 06-10-2026): NOT NULL hanya NOPOLIS dan IDPEGA. Tabel WARISAN - nol DDL.
//
//	IDPEGA ← p_IDPEGA pzInsKey = pengenal kasus (SAMA dengan LIFEINPRODUCTION)
//	NOPOLIS ← p_NOPOLIS PremiumListSummary.PL_NUMBER
//	TGL_PROD ← P_TGL_PROD CARI36 - `models.TanggalProduksi` (aturan periode PL Number)
//	USERNAME ← P_USERNAME CARI13 = CREATE_OP_NAME (SAMA dengan CREATEOPNAME LIFEINPRODUCTION)
//	TGL_INPUT ← SYSDATE
//	PRODKE ← '0' (VARCHAR2(5), teks) - keputusan work owner 08-10-2026; dulu NULL
//
// ⚠️ Dampak PRODKE '0' ke Endorsement Life (dibaca, kode endorsement TIDAK diubah):
// `sqlVersiWarisan` memakai NVL(PRODKE, 1) untuk VERSI → versi warisan 0 < versi
// aplikasi 1, sehingga baris aplikasi tetap menang; PRODKE mentah untuk NOMOR
// endorsement memperlakukan kosong = 0 → '0' memberi nomor yang sama dengan kosong.
//
// ⛔ DATA_JSON (P_JSONDATA) TIDAK diisi (NULL). NOENDORS, ERR_NOTE, TGL_KONVERSI,
// STS_KONVERSI, TGL_KONVERSI_RETRO, STS_KONVERSI_RETRO, OLDNOPOLIS: NULL - diisi
// layanan hilir.
// ⛔ Sekali per kasus: DELETE lalu INSERT - Confirm ulang tidak menggandakan.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
)

// TabelJSONPolis - tabel warisan yang ditulis.
const TabelJSONPolis = "JSON_POLIS"

// KolomJSONPolis - kolom yang DIISI, urut penampung.
var KolomJSONPolis = []string{"IDPEGA", "NOPOLIS", "TGL_PROD", "USERNAME", "TGL_INPUT", "PRODKE"}

func sqlTglProdLama(tujuan string) string {
	return fmt.Sprintf(`SELECT MAX(TGL_PROD) FROM %s WHERE IDPEGA = :1`, tujuan)
}

func sqlHapusJSONPolis(tujuan string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE IDPEGA = :1`, tujuan)
}

// sqlSisipJSONPolis - :1 IDPEGA, :2 NOPOLIS, :3 TGL_PROD, :4 kasus (baris utama).
// PRODKE literal teks '0' (08-10-2026) - tetap, bukan penampung.
func sqlSisipJSONPolis(tujuan, polis string) string {
	return fmt.Sprintf(`INSERT INTO %s (%s) SELECT :1, :2, :3, p.CREATE_OP_NAME, SYSDATE, '0' FROM %s p
	  WHERE p.ID = :4 AND p.ID = p.ID_PEGA`, tujuan, strings.Join(KolomJSONPolis, ", "), polis)
}

// JSONPolis menulis JSON_POLIS.
type JSONPolis struct{ db *db.DB }

// NewJSONPolis menyusunnya.
func NewJSONPolis(db *db.DB) *JSONPolis { return &JSONPolis{db: db} }

// TglProdLama - TGL_PROD baris JSON_POLIS kasus ini yang sudah ada (nil bila tidak ada),
// dipakai Confirm ulang atas polis yang nomornya sudah terbit lebih dahulu.
func (r *JSONPolis) TglProdLama(ctx context.Context, tx *db.Tx, kasusID string) (*time.Time, error) {
	tujuan, err := r.db.Qualify(TabelJSONPolis)
	if err != nil {
		return nil, err
	}
	q := sqlTglProdLama(tujuan)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	var t sql.NullTime
	if err := tx.QueryRowContext(ctx, q, kasusID).Scan(&t); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("repository: membaca TGL_PROD %s kasus %q: %w", TabelJSONPolis, kasusID, err)
	}
	if !t.Valid {
		return nil, nil
	}
	return &t.Time, nil
}

// Ganti menghapus baris JSON_POLIS kasus ini lalu menyisipkan satu baris baru di tx.
func (r *JSONPolis) Ganti(ctx context.Context, tx *db.Tx, kasusID, plNumber string, tglProd time.Time) error {
	if strings.TrimSpace(plNumber) == "" {
		return errors.New("repository: menolak menulis JSON_POLIS tanpa PL Number")
	}
	tujuan, err := r.db.Qualify(TabelJSONPolis)
	if err != nil {
		return err
	}
	polis, err := r.db.Qualify("T_PREMIUM_LIST")
	if err != nil {
		return err
	}
	qh, qs := sqlHapusJSONPolis(tujuan), sqlSisipJSONPolis(tujuan, polis)
	for _, q := range []string{qh, qs} {
		if err := db.PeriksaSQL(q); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, qh, kasusID); err != nil {
		return fmt.Errorf("repository: menghapus %s kasus %q: %w", TabelJSONPolis, kasusID, err)
	}
	h, err := tx.ExecContext(ctx, qs, kasusID, plNumber, tglProd, kasusID)
	if err != nil {
		return fmt.Errorf("repository: menulis %s kasus %q: %w", TabelJSONPolis, kasusID, err)
	}
	if n, err := h.RowsAffected(); err == nil && n != 1 {
		return fmt.Errorf("repository: %s kasus %q menulis %d baris, mau 1", TabelJSONPolis, kasusID, n)
	}
	return nil
}
