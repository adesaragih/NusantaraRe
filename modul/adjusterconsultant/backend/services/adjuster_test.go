package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/adjusterconsultant/backend/models"
	"nusantarare/modul/adjusterconsultant/backend/services"
	"nusantarare/modul/adjusterconsultant/backend/tiruan"
)

var ctx = context.Background()

var (
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	lihat = services.Aktor{AkunID: "UJI-LIHAT"}
)

func layanan(g *tiruan.Gudang) *services.Layanan { return services.BaruLayanan(g, tiruan.Transaksi) }

// Add seperti SaveAdjusterConsultant_Act: NAME dan ADDRESS huruf besar, TELPNO apa adanya, USERNAME akun login, ID =
// situs + sequence 4 digit dan nomor yang sudah terpakai dilompati.
func TestAddSepertiPega(t *testing.T) {
	g := tiruan.Contoh()
	a, err := layanan(g).Simpan(ctx, penuh, models.Isian{Name: "  uji baru ", Address: " jl. uji 1 ", TelpNo: " +62 21 123 "})
	if err != nil {
		t.Fatal(err)
	}
	mau := models.Adjuster{ID: "10004", Name: "UJI BARU", Address: "JL. UJI 1", TelpNo: "+62 21 123", Username: "UJI-ADMIN",
		EditDate: g.Jam, Active: true}
	if a != mau {
		t.Errorf("baru %+v", a)
	}
}

func TestEditDanValidasi(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	a, err := l.Simpan(ctx, penuh, models.Isian{ID: "10001", Name: "uji adjuster satu", Address: "baru"})
	if err != nil || a.Name != "UJI ADJUSTER SATU" || a.Address != "BARU" || a.Username != "UJI-ADMIN" {
		t.Errorf("edit nama sendiri %+v %v", a, err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "TIDAK-ADA", Name: "X"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("edit ID tak ada: %v", err)
	}
	_, err = l.Simpan(ctx, penuh, models.Isian{Name: " ", TelpNo: strings.Repeat("9", models.BatasTelp+1)})
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), "Name is required") ||
		!strings.Contains(err.Error(), "Telp No is longer than 50") {
		t.Errorf("validasi: %v", err)
	}
}

// Nama ganda ditolak (keputusan work owner 05-10-2026), juga terhadap baris nonaktif - pesannya menyarankan Activate.
func TestNamaGandaDitolak(t *testing.T) {
	l := layanan(tiruan.Contoh())
	if _, err := l.Simpan(ctx, penuh, models.Isian{Name: "uji adjuster dua"}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "already used by 10002") {
		t.Errorf("ganda aktif: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{Name: "Uji Konsultan Lama"}); err == nil ||
		!strings.Contains(err.Error(), "10003, which is inactive - activate it") {
		t.Errorf("ganda nonaktif: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10001", Name: "UJI ADJUSTER DUA"}); err == nil {
		t.Error("edit ke nama baris lain harus ditolak")
	}
}

// Pengganti hapus: Activate / Deactivate; View only tidak boleh menulis apa pun.
func TestAktifNonaktifDanViewOnly(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	a, err := l.SetelAktif(ctx, penuh, "10001", false)
	if err != nil || a.Active || a.Username != "UJI-ADMIN" {
		t.Errorf("nonaktif %+v %v", a, err)
	}
	if a, err := l.SetelAktif(ctx, penuh, "10003", true); err != nil || !a.Active {
		t.Errorf("aktif %+v %v", a, err)
	}
	if _, err := l.SetelAktif(ctx, penuh, "TIDAK-ADA", true); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("tak ada: %v", err)
	}
	if _, err := l.SetelAktif(ctx, lihat, "10002", false); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("View only nonaktif: %v", err)
	}
	if _, err := l.Simpan(ctx, lihat, models.Isian{Name: "UJI X"}); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("View only simpan: %v", err)
	}
	if !g.Baris["10002"].Active {
		t.Error("View only tidak boleh mengubah data")
	}
}

func TestDaftarSaring(t *testing.T) {
	l := layanan(tiruan.Contoh())
	for status, mau := range map[string]int{"": 3, "active": 2, "inactive": 1} {
		d, err := l.Daftar(ctx, "", status)
		if err != nil || len(d) != mau {
			t.Errorf("status %q: %d %v", status, len(d), err)
		}
	}
	if d, _ := l.Daftar(ctx, "konsultan", ""); len(d) != 1 || d[0].ID != "10003" {
		t.Errorf("cari %+v", d)
	}
	if _, err := l.Daftar(ctx, "", "hapus"); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("status asing: %v", err)
	}
}
