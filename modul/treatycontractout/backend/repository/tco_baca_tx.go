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
// Dibaca sesudah: tco_identitas.go.

import (
	"context"
	"database/sql"

	"nusantarare/inti/backend/db"
)

// kuerierTCO - bagian `*sql.DB` / `*sql.Tx` yang dipakai pembaca modul.
type kuerierTCO interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type kunciBacaTxTCO struct{}

// DenganBacaTxTCO menandai ctx: pembaca modul memakai tx. tx nil (uji tanpa
// Oracle) mengembalikan ctx apa adanya.
func DenganBacaTxTCO(ctx context.Context, tx *db.Tx) context.Context {
	if tx == nil {
		return ctx
	}
	return context.WithValue(ctx, kunciBacaTxTCO{}, tx)
}

// bacaTCO memilih jalur baca: transaksi bila ctx menandainya, pool bila tidak.
//
// Refactor bentuk B (30-09-2026): fungsi, bukan lagi metode `*DB` - `DB`
// kini tinggal di paket bersama. Yang dikembalikan pembungkus `*Tx`/`*DB`,
// yang meneruskan kueri apa adanya ke `*sql.Tx`/`*sql.DB` di dalamnya.
func bacaTCO(ctx context.Context, d *db.DB) kuerierTCO {
	if tx, _ := ctx.Value(kunciBacaTxTCO{}).(*db.Tx); tx.Terisi() {
		return tx
	}
	return d
}
