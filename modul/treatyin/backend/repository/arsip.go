package repository

// Jalur arsip muatan keluar - tiket 42 dan 74.
//
// ⛔ DUA pernyataan SQL di seluruh berkas ini, dan hanya yang pertama menyentuh
// kolom `MUATAN` - untuk MENULISnya. Pernyataan kedua menghitung dan
// mengelompokkan; ia tidak pernah memilih `MUATAN`. `INV-61`: arsip tidak
// punya jalur baca.
//
// Dijaga `TestArsipTidakPunyaJalurBaca` di `migrasi_invarian_test.go`, yang
// menyapu SELURUH berkas Go modul ini dan mencetak angkanya.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"nusantarare/modul/treatyin/backend/models"
)

// SimpanArsipMuatanKeluar menyisipkan satu pengiriman - tiket 42.
func (g *Gudang) SimpanArsipMuatanKeluar(ctx context.Context, a models.ArsipMuatanKeluar) error {
	tabel, err := g.db.Qualify("ARSIP_MUATAN_KELUAR")
	if err != nil {
		return err
	}
	tx, err := g.db.Mulai(ctx)
	if err != nil {
		return fmt.Errorf("repository: membuka transaksi: %w", err)
	}
	selesai := false
	defer func() {
		if !selesai {
			_ = tx.Rollback()
		}
	}()

	id, err := nomor(ctx, g, tx, "SEQ_TRIN_ARSIP_MUATAN_KELUAR")
	if err != nil {
		return err
	}
	berhasil := "0"
	if a.Berhasil {
		berhasil = "1"
	}
	q := fmt.Sprintf("INSERT INTO %s (ID_ARSIP_MUATAN_KELUAR,ID_KONTRAK,TUJUAN,DIKIRIM_PADA,"+
		"BERHASIL,MUATAN) VALUES (:1,:2,:3,TO_DATE(:4,'YYYY-MM-DD\"T\"HH24:MI:SS'),:5,:6)", tabel)
	// Cap waktu dipotong ke detik: kolomnya `DATE`, dan Oracle `DATE` memang
	// beresolusi detik. Pemotongan dinyatakan di sini, bukan dibiarkan driver
	// menebaknya.
	t, err := time.Parse(time.RFC3339, a.DikirimPada)
	if err != nil {
		return fmt.Errorf("repository: waktu pengiriman %q bukan RFC 3339: %w", a.DikirimPada, err)
	}
	if _, err := tx.ExecContext(ctx, q, id, a.IDKontrak, a.Tujuan,
		t.UTC().Format("2006-01-02T15:04:05"), berhasil, a.Muatan); err != nil {
		return fmt.Errorf("repository: menyisipkan ARSIP_MUATAN_KELUAR kontrak %d: %w", a.IDKontrak, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("repository: menutup transaksi: %w", err)
	}
	selesai = true
	return nil
}

// BuktiArsipKontrak menghitung arsip sebuah kontrak - tiket 42.
//
// ⛔ Daftar kolom yang dipilih: `COUNT(*)`, `MAX(DIKIRIM_PADA)`, dan
// `LISTAGG(DISTINCT TUJUAN)`. Nol di antaranya `MUATAN`, dan tidak ada
// `SELECT *` di sini yang dapat menariknya ikut.
func (g *Gudang) BuktiArsipKontrak(ctx context.Context, idKontrak int64) (models.BuktiArsip, error) {
	tabel, err := g.db.Qualify("ARSIP_MUATAN_KELUAR")
	if err != nil {
		return models.BuktiArsip{}, err
	}
	q := fmt.Sprintf(`SELECT COUNT(*),
		TO_CHAR(MAX(DIKIRIM_PADA),'YYYY-MM-DD"T"HH24:MI:SS'),
		LISTAGG(DISTINCT TUJUAN, ',') WITHIN GROUP (ORDER BY TUJUAN)
		FROM %s WHERE ID_KONTRAK = :1`, tabel)

	var cacah int64
	var terakhir, tujuan sql.NullString
	if err := g.db.QueryRowContext(ctx, q, idKontrak).Scan(&cacah, &terakhir, &tujuan); err != nil {
		return models.BuktiArsip{}, fmt.Errorf("repository: menghitung arsip kontrak %d: %w", idKontrak, err)
	}
	b := models.BuktiArsip{Cacah: cacah, Tujuan: []string{}}
	if terakhir.Valid {
		// Dikembalikan sebagai RFC 3339 UTC; kolomnya `DATE`, jadi zona
		// waktunya tidak tersimpan dan tidak dikarang di sini.
		b.TerakhirDikirim = terakhir.String + "Z"
	}
	if tujuan.Valid && tujuan.String != "" {
		b.Tujuan = strings.Split(tujuan.String, ",")
	}
	return b, nil
}
