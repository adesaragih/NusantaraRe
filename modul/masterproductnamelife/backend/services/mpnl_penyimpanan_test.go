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

	"nusantarare/modul/masterproductnamelife/backend/models"
)

func TestPenyimpananLokalKirimIdempotenHapusTanpaGalat(t *testing.T) {
	akar := t.TempDir()
	p := PenyimpananLokal(akar)
	ctx := context.Background()
	if err := p.SimpanAntrean(ctx, "ABC123", strings.NewReader("isi")); err != nil {
		t.Fatal(err)
	}
	o := models.ObjekPenyimpanan{ImageID: "ABC123", AppFolder: "Contract/Doc/2026/10/", FileName: "x - UJI.pdf", AppName: "UJI-APP"}
	for i := 0; i < 2; i++ {
		hasil, err := p.Kirim(ctx, o, "pdf", "application/pdf")
		if err != nil {
			t.Fatalf("kirim ke-%d: %v", i+1, err)
		}
		// Stub: objek tercatat apa adanya - URLPUBLIC kosong, alamat penyimpanan tidak ada.
		if hasil != o {
			t.Errorf("objek stub: %+v", hasil)
		}
	}
	r, baru, err := p.Buka(ctx, o)
	if err != nil || baru != nil {
		t.Fatal(err, baru)
	}
	isi, _ := io.ReadAll(r)
	_ = r.Close()
	if string(isi) != "isi" {
		t.Errorf("isi: %q", isi)
	}
	if _, err := os.Stat(filepath.Join(akar, folderStub, "simpan", "ABC123")); err != nil {
		t.Errorf("berkas di folder stub: %v", err)
	}
	if err := p.Hapus(ctx, "ABC123", &o); err != nil {
		t.Errorf("hapus: %v", err)
	}
	if err := p.Hapus(ctx, "ABC123", nil); err != nil {
		t.Errorf("hapus yang sudah tidak ada bukan galat: %v", err)
	}
	if _, err := p.Kirim(ctx, o, "pdf", "application/pdf"); !errors.Is(err, ErrBerkasSumberHilang) {
		t.Errorf("sumber hilang: %v", err)
	}
	if err := p.SimpanAntrean(ctx, "../keluar", strings.NewReader("x")); !errors.Is(err, errImageIDTidakSah) {
		t.Errorf("kunci objek berjalur ditolak: %v", err)
	}
	if _, err := PenyimpananLokal(" ").Kirim(ctx, models.ObjekPenyimpanan{ImageID: "X"}, "pdf", ""); !errors.Is(err, ErrPenyimpananBelumDisetel) {
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
