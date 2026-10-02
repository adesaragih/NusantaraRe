//go:build db

package repository_test

// Satu pintu skema uji untuk uji db repository PremiumList Life.
//
// Untuk apa berkas ini: setiap uji db di paket ini membuka skema uji dengan
// urutan yang sama - buka, lewati HANYA bila Oracle belum dikonfigurasi, ping,
// pasang dari keadaan bersih, bongkar sesudahnya. Satu pemanggil
// `skemauji.Buka()` untuk semuanya: penjaga Claim Life
// `TestSetiapPemanggilBukaMemeriksaBolehDilewati` menghitung pemanggil itu di
// seluruh repo (brief seragam kolom 01-10-2026, migrasi 059).

import (
	"context"
	"database/sql"
	"testing"

	"nusantarare/uji/skemauji"
)

// pasangSkemaUji membuka dan memasang skema uji, atau melewati test.
func pasangSkemaUji(t *testing.T) (*sql.DB, string, context.Context) {
	t.Helper()
	sqlDB, skema, err := skemauji.Buka()
	if err != nil {
		if !skemauji.BolehDilewati(err) {
			t.Fatalf("skema uji menolak: %v", err)
		}
		t.Skipf("lewati: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	if err := skemauji.Pasang(ctx, sqlDB, skema); err != nil {
		t.Fatalf("memasang skema uji: %v", err)
	}
	// t.Cleanup berjalan terbalik: bongkar dulu, baru koneksinya ditutup.
	t.Cleanup(func() { _ = skemauji.Bongkar(ctx, sqlDB, skema) })
	return sqlDB, skema, ctx
}
