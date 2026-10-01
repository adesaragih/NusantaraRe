package repository

// `Add CSV Data` b10405 → `Activity/SaveCSVEDMLife.xml` (5 langkah, nol `//`).
//
// ⛔ Unggah ulang MENGGANTI, bukan menumpuk (langkah 2.1 b454): baris `New`
// kasus dibuang lebih dulu, di transaksi yang sama dengan penyisipan.
// Baris `Old`/`Delete`/`Batal` tidak tersentuh.
//
// Dibaca sesudah: edm_simpan.go.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
)

// sqlAcuanCSV - `PremiumListDetail(1)` (4.1 b2466/b2495): peserta pertama di
// urutan grid (`sqlPeserta`), di luar baris `New` yang dibuang langkah 2.1.
func sqlAcuanCSV(peserta string) string {
	return fmt.Sprintf(`SELECT d.PLAN, d.POLICY_HOLDER FROM %s d
	  WHERE d.PREMIUM_LIST_ID = :1 AND (d.EDM_STATUS IS NULL OR d.EDM_STATUS <> :2)
	  ORDER BY d.CERTIFICATE_NO, d.NAME_OF_INSURED, d.ID FETCH FIRST 1 ROWS ONLY`, peserta)
}

// AcuanCSV membaca baris acuan; false bila kasus tanpa peserta lama.
func (g *Gudang) AcuanCSV(ctx context.Context, tx *db.Tx, kasusID string) (models.AcuanCSV, bool, error) {
	n, err := g.nama(tabelPeserta)
	if err != nil {
		return models.AcuanCSV{}, false, err
	}
	q := sqlAcuanCSV(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return models.AcuanCSV{}, false, err
	}
	var plan, holder sql.NullString
	err = g.pakai(tx).QueryRowContext(ctx, q, kasusID, models.StatusNew).Scan(&plan, &holder)
	if errors.Is(err, sql.ErrNoRows) {
		return models.AcuanCSV{}, false, nil
	}
	if err != nil {
		return models.AcuanCSV{}, false, fmt.Errorf("repository: membaca acuan CSV kasus %q: %w", kasusID, err)
	}
	return models.AcuanCSV{Plan: plan.String, PolicyHolder: holder.String}, true, nil
}

// sqlHapusPesertaBaru - langkah 2.1 b454, prakondisi b547 `.EDMStatus=="New"`.
func sqlHapusPesertaBaru(peserta string) string {
	return fmt.Sprintf(`DELETE FROM %s d WHERE d.PREMIUM_LIST_ID = :1 AND d.EDM_STATUS = :2`, peserta)
}

// HapusPesertaBaru membuang baris `New` kasus; mengembalikan cacahnya.
func (g *Gudang) HapusPesertaBaru(ctx context.Context, tx *db.Tx, kasusID string) (int, error) {
	n, err := g.nama(tabelPeserta)
	if err != nil {
		return 0, err
	}
	q := sqlHapusPesertaBaru(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	h, err := tx.ExecContext(ctx, q, kasusID, models.StatusNew)
	if err != nil {
		return 0, fmt.Errorf("repository: membuang peserta New kasus %q: %w", kasusID, err)
	}
	c, err := h.RowsAffected()
	return int(c), err
}

// sqlSisipPesertaCSV - satu peserta `New` (4.3 b2715). `ID` = hash kasus +
// nomor baris CSV: unggahan yang sama menghasilkan pengenal yang sama, dan
// pengenal lama sudah dibuang 2.1 di transaksi ini. Penampung :1 kasus,
// :2 nomor, :3 kasus, :4 PL_NUMBER, :5 'New', :6… `models.KolomCSVTersimpan`.
func sqlSisipPesertaCSV(peserta string) string {
	kolom := []string{"ID", "PREMIUM_LIST_ID", "PL_NUMBER", "EDM_STATUS"}
	nilai := []string{"RAWTOHEX(STANDARD_HASH(:1 || '/N/' || :2, 'MD5'))", ":3", ":4", ":5"}
	for i, c := range models.KolomCSVTersimpan() {
		kolom = append(kolom, c.Kolom)
		p := ":" + strconv.Itoa(6+i)
		if c.Jenis == models.CSVTanggal {
			p = "TO_DATE(" + p + ", 'YYYY-MM-DD')"
		}
		nilai = append(nilai, p)
	}
	return fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, peserta, strings.Join(kolom, ", "), strings.Join(nilai, ", "))
}

// SisipPesertaCSV menyisipkan satu kelompok baris CSV ber-`New`.
//
// ⛔ Uang dikirim sebagai TEKS DESIMAL bertitik (pola penyisipan unggahan
// PremiumList) - tidak pernah `float64`. Kosong = NULL.
func (g *Gudang) SisipPesertaCSV(ctx context.Context, tx *db.Tx, kasusID, plNumber string, baris []models.BarisCSV) (int, error) {
	n, err := g.nama(tabelPeserta)
	if err != nil {
		return 0, err
	}
	q := sqlSisipPesertaCSV(n[0])
	if err := db.PeriksaSQL(q); err != nil {
		return 0, err
	}
	stmt, err := tx.PrepareContext(ctx, q)
	if err != nil {
		return 0, fmt.Errorf("repository: menyiapkan penyisipan CSV: %w", err)
	}
	defer func() { _ = stmt.Close() }()
	kolom := models.KolomCSVTersimpan()
	for _, b := range baris {
		args := []any{kasusID, strconv.Itoa(b.Nomor), kasusID, plNumber, models.StatusNew}
		for _, c := range kolom {
			args = append(args, db.KosongJadiNil(b.Nilai[c.Kolom]))
		}
		if _, err := stmt.ExecContext(ctx, args...); err != nil {
			return 0, fmt.Errorf("repository: menyisipkan baris CSV %d: %w", b.Nomor, err)
		}
	}
	return len(baris), nil
}
