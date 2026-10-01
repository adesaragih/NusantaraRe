package services

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPenyimpananLokalKirimIdempotenHapusTanpaGalat(t *testing.T) {
	akar := t.TempDir()
	p := PenyimpananLokal(akar)
	ctx := context.Background()
	if err := p.SimpanAntrean(ctx, "ABC123", strings.NewReader("isi")); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := p.Kirim(ctx, "ABC123"); err != nil {
			t.Fatalf("kirim ke-%d: %v", i+1, err)
		}
	}
	r, err := p.Buka(ctx, "ABC123")
	if err != nil {
		t.Fatal(err)
	}
	isi, _ := io.ReadAll(r)
	_ = r.Close()
	if string(isi) != "isi" {
		t.Errorf("isi: %q", isi)
	}
	if _, err := os.Stat(filepath.Join(akar, folderStub, "simpan", "ABC123")); err != nil {
		t.Errorf("berkas di folder stub: %v", err)
	}
	if err := p.Hapus(ctx, "ABC123"); err != nil {
		t.Errorf("hapus: %v", err)
	}
	if err := p.Hapus(ctx, "ABC123"); err != nil {
		t.Errorf("hapus yang sudah tidak ada bukan galat: %v", err)
	}
	if err := p.Kirim(ctx, "ABC123"); !errors.Is(err, ErrBerkasSumberHilang) {
		t.Errorf("sumber hilang: %v", err)
	}
	if err := p.SimpanAntrean(ctx, "../keluar", strings.NewReader("x")); !errors.Is(err, errImageIDTidakSah) {
		t.Errorf("kunci objek berjalur ditolak: %v", err)
	}
	if err := PenyimpananLokal(" ").Kirim(ctx, "X"); !errors.Is(err, ErrPenyimpananBelumDisetel) {
		t.Errorf("UNGGAHAN_DIR kosong gagal terang: %v", err)
	}
}

func TestNamaObjekPolaPega(t *testing.T) {
	// `InsertGoogleStorage_Act` 8 b1339: Asia/Jakarta, `yyyyMMdd-hhmmss-S` (jam 12-an, milidetik tanpa nol depan).
	folder, file := namaObjek(time.Date(2026, 1, 5, 17, 4, 5, 7_000_000, time.UTC), "x.pdf")
	if folder != "Contract/Doc/2026/01/" || file != "20260106-120405-7 - x.pdf" {
		t.Errorf("%q %q", folder, file)
	}
}
