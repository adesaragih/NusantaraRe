package repository

// JSON_POLIS endorsement - satu baris per kasus yang diresmikan (Confirm), keputusan work owner 08-10-2026.
//
// Untuk apa berkas ini: rekam versi polis yang dibaca pihak hilir (dan jalur versi warisan modul lain) ditulis
// dengan SQL LANGSUNG di transaksi Confirm, SESUDAH kepala diresmikan dan NO_ENDORS terbit. Prosedur Pega
// `InsertJsonPolisEDM` tidak dipanggil.
//
// `[terverifikasi]` kolom = ALL_TAB_COLUMNS POOLDATA (diberikan user): NOT NULL hanya NOPOLIS dan IDPEGA;
// PRODKE VARCHAR2(5). Tabel WARISAN - nol DDL.
//
//	IDPEGA   ← ID kasus (`EDML-<n>`)              NOPOLIS ← PL Number
//	PRODKE   ← PRODKE baru sebagai teks            NOENDORS ← NO_ENDORS (`<PL>/<PRODKE dua digit>`)
//	TGL_PROD ← `models.TanggalProduksi` (TANGGAL_CLOSING)
//	USERNAME ← `CREATE_OP_NAME` kepala kasus - SAMA dengan CREATEOPNAME `LIFEINPRODUCTION` endorsement
//	TGL_INPUT ← SYSDATE
//
// ⛔ DATA_JSON, ERR_NOTE, TGL_KONVERSI, STS_KONVERSI, TGL_KONVERSI_RETRO, STS_KONVERSI_RETRO, OLDNOPOLIS: NULL.
// ⛔ Sekali per kasus: DELETE WHERE IDPEGA lalu INSERT. Nol COMMIT (ADR-U-0029).

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"nusantarare/inti/backend/db"
)

// KolomJSONPolisEDM - kolom yang DIISI, urut penampung.
var KolomJSONPolisEDM = []string{"IDPEGA", "NOPOLIS", "PRODKE", "NOENDORS", "TGL_PROD", "USERNAME", "TGL_INPUT"}

// JSONPolisTulis - bahan satu baris JSON_POLIS endorsement.
type JSONPolisTulis struct {
	KasusID    string
	NomorPolis string
	ProdKe     int
	Nomor      string
	TglProd    time.Time
}

func sqlHapusJSONPolisEDM(tujuan string) string {
	return fmt.Sprintf(`DELETE FROM %s WHERE IDPEGA = :1`, tujuan)
}

// sqlSisipJSONPolisEDM - :1 IDPEGA, :2 NOPOLIS, :3 PRODKE (teks), :4 NOENDORS, :5 TGL_PROD, :6 kasus.
func sqlSisipJSONPolisEDM(tujuan, polis string) string {
	return fmt.Sprintf(`INSERT INTO %s (%s) SELECT :1, :2, :3, :4, :5, p.CREATE_OP_NAME, SYSDATE FROM %s p
	  WHERE p.ID = :6`, tujuan, strings.Join(KolomJSONPolisEDM, ", "), polis)
}

// TulisJSONPolis mengganti baris JSON_POLIS kasus ini di tx; mengembalikan cacah baris (1).
func (g *Gudang) TulisJSONPolis(ctx context.Context, tx *db.Tx, r JSONPolisTulis) (int, error) {
	if strings.TrimSpace(r.NomorPolis) == "" || strings.TrimSpace(r.Nomor) == "" || r.ProdKe < 1 {
		return 0, fmt.Errorf("repository: menolak menulis JSON_POLIS kasus %q tanpa nomor/PRODKE", r.KasusID)
	}
	n, err := g.nama(tabelPolisWarisan, tabelPolis)
	if err != nil {
		return 0, err
	}
	qh, qs := sqlHapusJSONPolisEDM(n[0]), sqlSisipJSONPolisEDM(n[0], n[1])
	for _, q := range []string{qh, qs} {
		if err := db.PeriksaSQL(q); err != nil {
			return 0, err
		}
	}
	if _, err := tx.ExecContext(ctx, qh, r.KasusID); err != nil {
		return 0, fmt.Errorf("repository: menghapus JSON_POLIS kasus %q: %w", r.KasusID, err)
	}
	h, err := tx.ExecContext(ctx, qs, r.KasusID, r.NomorPolis, strconv.Itoa(r.ProdKe), r.Nomor, r.TglProd, r.KasusID)
	if err != nil {
		return 0, fmt.Errorf("repository: menulis JSON_POLIS kasus %q: %w", r.KasusID, err)
	}
	if c, err := h.RowsAffected(); err == nil && c != 1 {
		return 0, fmt.Errorf("repository: JSON_POLIS kasus %q menulis %d baris, mau 1", r.KasusID, c)
	}
	return 1, nil
}

// tabelTanggalClosing - hari tutup buku (warisan POOLDATA, kolom TANGGAL VARCHAR2).
const tabelTanggalClosing = "TANGGAL_CLOSING"

func sqlHariClosing(t string) string {
	return fmt.Sprintf(`SELECT TANGGAL FROM %s WHERE ROWNUM = 1`, t)
}

// ErrTanggalClosingKosong - TANGGAL_CLOSING kosong / tak terurai: periode produksi tidak dapat ditentukan.
var ErrTanggalClosingKosong = errors.New("repository: TANGGAL_CLOSING kosong atau bukan bilangan")

// HariClosing membaca hari tutup buku untuk periode produksi (TGL_PROD JSON_POLIS).
func (g *Gudang) HariClosing(ctx context.Context, tx *db.Tx) (int, error) {
	n, err := g.nama(tabelTanggalClosing)
	if err != nil {
		return 0, err
	}
	q := sqlHariClosing(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	var teks sql.NullString
	if err := g.pakai(tx).QueryRowContext(ctx, q).Scan(&teks); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrTanggalClosingKosong
		}
		return 0, fmt.Errorf("repository: membaca TANGGAL_CLOSING: %w", err)
	}
	hari, err := strconv.Atoi(strings.TrimSpace(teks.String))
	if err != nil || hari < 1 {
		return 0, fmt.Errorf("%w: %q", ErrTanggalClosingKosong, teks.String)
	}
	return hari, nil
}
