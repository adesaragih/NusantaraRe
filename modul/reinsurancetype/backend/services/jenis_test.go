package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/reinsurancetype/backend/models"
	"nusantarare/modul/reinsurancetype/backend/services"
	"nusantarare/modul/reinsurancetype/backend/tiruan"
)

var ctx = context.Background()

var (
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	lihat = services.Aktor{AkunID: "UJI-LIHAT"}
	jam   = time.Date(2026, 10, 5, 3, 4, 5, 600e6, time.UTC)
)

func layanan(g *tiruan.Gudang) *services.Layanan {
	return services.BaruLayanan(g, tiruan.Transaksi).DenganJam(func() time.Time { return jam })
}

// Add: huruf besar, Code kosong = 00, ID `1` + sequence 4 digit (nomor terpakai dilompati), TGLUPDATE Pega GMT, USERID.
func TestAdd(t *testing.T) {
	g := tiruan.Contoh()
	j, err := layanan(g).Simpan(ctx, penuh, models.Isian{Name: " uji xl 6th layer ", Type: "4", SoaName: " uji xl ", Flag: "Active",
		GroupType: "qs"})
	if err != nil {
		t.Fatal(err)
	}
	mau := models.Jenis{ID: "10263", Name: "UJI XL 6TH LAYER", Type: "4", SoaName: "UJI XL", Code: "00", Flag: "active",
		GroupType: "QS", UserID: "UJI-ADMIN", TglUpdate: "20261005T030405.600 GMT"}
	if g.Baris["10263"] != mau {
		t.Errorf("tersimpan\n%+v\nmau\n%+v", g.Baris["10263"], mau)
	}
	if j.Diubah != "05-10-2026 10:04" {
		t.Errorf("tampilan %+v", j)
	}
}

// Nilai warisan Flag `1` / kosong dan Type kosong boleh tetap selama tidak diubah; nilai baru harus pilihan sah.
func TestNilaiWarisanTetap(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10196", Name: "uji qs", Type: "2", Code: "00", Flag: "1"}); err != nil {
		t.Errorf("Flag 1 tidak diubah: %v", err)
	}
	if g.Baris["10196"].Flag != "1" {
		t.Error("Flag 1 berubah")
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10033", Name: "UJI POOL SURPLUS", SoaName: "UJI POOL", Code: "00"}); err != nil {
		t.Errorf("Type dan Flag kosong tidak diubah: %v", err)
	}
	_, err := l.Simpan(ctx, penuh, models.Isian{ID: "10007", Name: "UJI ORS", Type: "1", Code: "00", Flag: "1"})
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), "Flag must be active or inactive") {
		t.Errorf("ganti ke Flag 1: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{Name: "UJI BARU", Code: "00", Flag: ""}); err == nil ||
		!strings.Contains(err.Error(), "Type must be 1") || !strings.Contains(err.Error(), "Flag must be active") {
		t.Errorf("baru tanpa Type / Flag: %v", err)
	}
}

func TestValidasiDanKembar(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	_, err := l.Simpan(ctx, penuh, models.Isian{Name: " ", Type: "5", Code: "A1", Flag: "aktif", NoUrut: "x", GroupType: "XX",
		SoaName: strings.Repeat("A", models.BatasTeks+1)})
	for _, mau := range []string{"Name is required", "Type must be 1", "Flag must be active", "Group Type must be OR",
		"Code must be a number", "No Urut must be a number", "SOA Name is longer than 100"} {
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), mau) {
			t.Errorf("tanpa %q: %v", mau, err)
		}
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{Name: "Uji Ors", Type: "1", Flag: "active"}); err == nil ||
		!strings.Contains(err.Error(), "Name UJI ORS is already used by 10007") {
		t.Errorf("nama kembar: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10262", Name: "UJI QS", Type: "2", Flag: "active"}); err == nil {
		t.Error("edit ke nama baris lain harus ditolak")
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "19999", Name: "UJI X", Type: "1", Flag: "active"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("edit ID tak ada: %v", err)
	}
	if _, err := l.Simpan(ctx, lihat, models.Isian{Name: "UJI Y", Type: "1", Flag: "active"}); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("View only: %v", err)
	}
	if len(g.Baris) != 5 {
		t.Error("ditolak tetapi tertulis")
	}
}

// Nonaktif = Flag inactive lewat Edit (pengganti hapus).
func TestNonaktifLewatEdit(t *testing.T) {
	g := tiruan.Contoh()
	if _, err := layanan(g).Simpan(ctx, penuh, models.Isian{ID: "10007", Name: "UJI ORS", Type: "1", Code: "00", Flag: "inactive",
		GroupType: "OR"}); err != nil {
		t.Fatal(err)
	}
	if g.Baris["10007"].Flag != "inactive" {
		t.Errorf("flag %q", g.Baris["10007"].Flag)
	}
}

func TestDaftarSaring(t *testing.T) {
	l := layanan(tiruan.Contoh())
	d, err := l.Daftar(ctx, "", "", "")
	if err != nil || len(d) != 5 || d[0].ID != "10033" {
		t.Errorf("urutan %+v %v", d, err)
	}
	if d, _ := l.Daftar(ctx, "", "2", "active"); len(d) != 1 || d[0].ID != "10262" {
		t.Errorf("saring %+v", d)
	}
	if d, _ := l.Daftar(ctx, "quota", "", ""); len(d) != 1 || d[0].ID != "10262" {
		t.Errorf("cari SOA %+v", d)
	}
	if j, err := l.Buka(ctx, "10030"); err != nil || j.Diubah != "29-10-2019 23:40" {
		t.Errorf("buka %+v %v", j, err)
	}
}
