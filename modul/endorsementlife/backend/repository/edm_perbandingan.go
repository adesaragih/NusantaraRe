package repository

// Perbandingan summary sebelum/sesudah endorsement (permintaan work owner
// 08-10-2026, bukan korpus) - baca-saja. Peserta kasus dijumlah per
// (CURRENCY, kelas EDM_STATUS) dengan ekspresi rekap yang SAMA dengan
// `sqlSisipRekap`, sehingga jumlah seluruh kelas = rekap kasus.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
)

// ekspresiRekap - pasangan (kolom rekap, ekspresi SUM) untuk
// `models.KolomRekapKasus[1:]` per Type: PREMIUM/BALANCE dari
// `rumusPremiBalance`, sisanya dari `kolomJumlahRekap` - sumber yang sama
// dengan `sqlSisipRekap`.
func ekspresiRekap(tipe string) [][2]string {
	premi, balance := rumusPremiBalance(tipe)
	jumlah := map[string]string{}
	for _, k := range kolomJumlahRekap {
		jumlah[k[0]] = "SUM(" + nv(k[1]) + ")"
	}
	hasil := make([][2]string, 0, len(models.KolomRekapKasus)-1)
	for _, k := range models.KolomRekapKasus[1:] {
		switch k {
		case "PREMIUM":
			hasil = append(hasil, [2]string{k, premi})
		case "BALANCE":
			hasil = append(hasil, [2]string{k, balance})
		default:
			e, ada := jumlah[k]
			if !ada {
				e = "0"
			}
			hasil = append(hasil, [2]string{k, e})
		}
	}
	return hasil
}

// kelompokStatus - kelas `EDM_STATUS` (literal tetap dari konstanta models,
// bukan masukan): Delete/Batal → D, New → N, selain itu (Old, kosong) → O.
func kelompokStatus() string {
	return fmt.Sprintf(`CASE WHEN d.EDM_STATUS IN ('%s', '%s') THEN '%s' WHEN d.EDM_STATUS = '%s' THEN '%s' ELSE '%s' END`,
		models.StatusDelete, models.StatusBatal, models.KelompokHapus, models.StatusNew, models.KelompokBaru, models.KelompokLama)
}

// sqlPerbandingan - jumlah peserta kasus per (CURRENCY, kelas).
func sqlPerbandingan(peserta, tipe string) string {
	kel := kelompokStatus()
	kolom := []string{"d.CURRENCY", kel, "COUNT(d.ID)"}
	for _, e := range ekspresiRekap(tipe) {
		kolom = append(kolom, "TO_CHAR("+e[1]+", 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')")
	}
	return fmt.Sprintf(`SELECT %s FROM %s d WHERE d.PREMIUM_LIST_ID = :1 GROUP BY d.CURRENCY, %s ORDER BY d.CURRENCY`,
		strings.Join(kolom, ", "), peserta, kel)
}

// PerbandinganKasus membaca jumlah peserta kasus per (mata uang, kelas).
func (g *Gudang) PerbandinganKasus(ctx context.Context, tx *db.Tx, kasusID, tipe string) ([]models.KelompokRekap, error) {
	n, err := g.nama(tabelPeserta)
	if err != nil {
		return nil, err
	}
	q := sqlPerbandingan(n[0], tipe)
	if err := db.PeriksaSQL(q); err != nil {
		return nil, err
	}
	rows, err := g.pakai(tx).QueryContext(ctx, q, kasusID)
	if err != nil {
		return nil, fmt.Errorf("repository: membaca perbandingan kasus %q: %w", kasusID, err)
	}
	defer func() { _ = rows.Close() }()
	ekspresi := ekspresiRekap(tipe)
	var hasil []models.KelompokRekap
	for rows.Next() {
		var uang, kel sql.NullString
		var cacah int
		nilai := make([]sql.NullString, len(ekspresi))
		tujuan := []any{&uang, &kel, &cacah}
		for i := range nilai {
			tujuan = append(tujuan, &nilai[i])
		}
		if err := rows.Scan(tujuan...); err != nil {
			return nil, fmt.Errorf("repository: membaca perbandingan kasus %q: %w", kasusID, err)
		}
		k := models.KelompokRekap{Currency: uang.String, Kelompok: kel.String, Cacah: cacah, Nilai: map[string]string{}}
		for i, e := range ekspresi {
			k.Nilai[e[0]] = nilai[i].String
		}
		hasil = append(hasil, k)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: membaca perbandingan kasus %q: %w", kasusID, err)
	}
	return hasil, nil
}
