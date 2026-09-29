package repository

// Pembaca lewat transaksi - Treaty Contract Out.
//
// Untuk apa berkas ini: selama `DenganBacaTxTCO`, setiap pembaca modul
// (`bacaTCO`) membaca lewat transaksi penulisnya, bukan pool - dipakai
// penulis yang membaca ulang baris SESUDAH menguncinya (kaskade hapus, ganti
// jenis kontrak, ganti tahun/grup, hitung ulang anak klausul).
//
// Sejak OQ-TCO-19 (keputusan work owner 29-09-2026) simpan utuh lintas enam
// tabel dibuang; identitas sementaranya ikut dibuang. Yang tersisa hanya
// pembaca ini. Transaksi tetap dibuka dan di-commit sekali oleh
// `Service.DalamTransaksi` pemanggil (nol COMMIT di teks SQL).
//
// Dibaca sesudah: tco_identitas.go, tco_jejak.go.

import (
	"context"
	"database/sql"
)

// kuerierTCO - bagian `*sql.DB` / `*sql.Tx` yang dipakai pembaca modul.
type kuerierTCO interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type kunciBacaTxTCO struct{}

// DenganBacaTxTCO menandai ctx: pembaca modul memakai tx. tx nil (uji tanpa
// Oracle) mengembalikan ctx apa adanya.
func DenganBacaTxTCO(ctx context.Context, tx *Tx) context.Context {
	if tx == nil {
		return ctx
	}
	return context.WithValue(ctx, kunciBacaTxTCO{}, tx)
}

// bacaTCO memilih jalur baca: transaksi bila ctx menandainya, pool bila tidak.
func (d *DB) bacaTCO(ctx context.Context) kuerierTCO {
	if tx, _ := ctx.Value(kunciBacaTxTCO{}).(*Tx); tx != nil && tx.tx != nil {
		return tx.tx
	}
	return d.sql
}
