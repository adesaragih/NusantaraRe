package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/treatygroupojk/backend/models"
	"nusantarare/modul/treatygroupojk/backend/services"
	"nusantarare/modul/treatygroupojk/backend/tiruan"
)

var ctx = context.Background()

var (
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	lihat = services.Aktor{AkunID: "UJI-LIHAT"}
)

func layanan(g *tiruan.Gudang) *services.Layanan { return services.BaruLayanan(g, tiruan.Transaksi) }

// Add: ID = nomor tertinggi + 1 dua digit, Order No = Order No tertinggi + 1 (keduanya tidak diketik; perintah work owner
// 05-10-2026: "ORDERNO hide aja, isi sesuai max dari order no"), sesudah seluruh baris dikunci; nama huruf besar.
func TestAddIDDanOrderBerikutnya(t *testing.T) {
	g := tiruan.Contoh()
	o, err := layanan(g).Simpan(ctx, penuh, models.Isian{Name: " uji cargo ", NameIDN: " uji pengangkutan "})
	if err != nil {
		t.Fatal(err)
	}
	if mau := (models.Ojk{ID: "10", Name: "UJI CARGO", NameIDN: "UJI PENGANGKUTAN", OrderNo: "12"}); o != mau {
		t.Errorf("baru %+v, mau %+v", o, mau)
	}
	if g.Dikunci != 1 {
		t.Errorf("Add mengunci %d kali, mau 1", g.Dikunci)
	}
	g.Baris["99"] = models.Ojk{ID: "99", Name: "UJI Y", NameIDN: "UJI Y", OrderNo: "lama"}
	if o, err := layanan(g).Simpan(ctx, penuh, models.Isian{Name: "UJI Z", NameIDN: "UJI Z"}); err != nil || o.ID != "100" || o.OrderNo != "13" {
		t.Errorf("sesudah 99 (Order No bukan angka diabaikan): %+v %v", o, err)
	}
}

// Edit: ID dan Order No tetap, tanpa mengunci; Name boleh tetap milik sendiri.
func TestEdit(t *testing.T) {
	g := tiruan.Contoh()
	o, err := layanan(g).Simpan(ctx, penuh, models.Isian{ID: "09", Name: "uji engineering", NameIDN: "uji rekayasa baru"})
	if err != nil || o.NameIDN != "UJI REKAYASA BARU" || o.ID != "09" || o.OrderNo != "2" {
		t.Errorf("edit %+v %v", o, err)
	}
	if g.Dikunci != 0 {
		t.Error("Edit tidak boleh mengunci seluruh tabel")
	}
	if _, err := layanan(g).Simpan(ctx, penuh, models.Isian{ID: "77", Name: "UJI", NameIDN: "UJI"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("edit ID tak ada: %v", err)
	}
}

// Validasi: wajib, panjang kolom, Name tidak kembar.
func TestValidasiDanKembar(t *testing.T) {
	l := layanan(tiruan.Contoh())
	_, err := l.Simpan(ctx, penuh, models.Isian{Name: " ", NameIDN: strings.Repeat("A", models.BatasNama+1)})
	for _, mau := range []string{"Name is required", "Name (IDN) is longer than 50"} {
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), mau) {
			t.Errorf("validasi tanpa %q: %v", mau, err)
		}
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{Name: "Uji Property", NameIDN: "UJI"}); err == nil ||
		!strings.Contains(err.Error(), "Name is already used by 01") {
		t.Errorf("nama kembar: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "01", Name: "UJI ENGINEERING", NameIDN: "UJI"}); err == nil {
		t.Error("edit ke nama baris lain harus ditolak")
	}
}

// View only dan tanpa login tidak boleh menulis.
func TestViewOnly(t *testing.T) {
	g := tiruan.Contoh()
	for _, a := range []services.Aktor{lihat, {Penuh: true}} {
		if _, err := layanan(g).Simpan(ctx, a, models.Isian{Name: "UJI", NameIDN: "UJI"}); !errors.Is(err, services.ErrDilarang) {
			t.Errorf("%+v: %v", a, err)
		}
	}
	if len(g.Baris) != 3 {
		t.Error("ditolak tetapi tertulis")
	}
}

// Daftar: urut Order No sebagai angka (`2` sebelum `11`), cari ID/Name/Name IDN.
func TestDaftar(t *testing.T) {
	l := layanan(tiruan.Contoh())
	d, err := l.Daftar(ctx, "")
	if err != nil || len(d) != 3 || d[0].ID != "01" || d[1].ID != "09" || d[2].ID != "05" {
		t.Errorf("urutan %+v %v", d, err)
	}
	if d, _ := l.Daftar(ctx, "rekayasa"); len(d) != 1 || d[0].ID != "09" {
		t.Errorf("cari %+v", d)
	}
}
