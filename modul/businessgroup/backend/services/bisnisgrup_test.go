package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/businessgroup/backend/models"
	"nusantarare/modul/businessgroup/backend/services"
	"nusantarare/modul/businessgroup/backend/tiruan"
)

var ctx = context.Background()

var (
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	lihat = services.Aktor{AkunID: "UJI-LIHAT"}
)

func layanan(g *tiruan.Gudang) *services.Layanan { return services.BaruLayanan(g, tiruan.Transaksi) }

// Add: TREATYNAME disalin dari Treaty Group; huruf besar; Alias kosong = Name; ID situs + sequence, nomor terpakai
// dilompati.
func TestAdd(t *testing.T) {
	g := tiruan.Contoh()
	b, err := layanan(g).Simpan(ctx, penuh, models.Isian{TopID: "10009", Name: " uji hull "})
	if err != nil {
		t.Fatal(err)
	}
	mau := models.BisnisGrup{ID: "10015", Name: "UJI HULL", Alias: "UJI HULL", TopID: "10009", TreatyName: "UJI MARINE CARGO"}
	if b != mau || g.Baris["10015"] != mau {
		t.Errorf("baru %+v, mau %+v", b, mau)
	}
	if b, _ := layanan(g).Simpan(ctx, penuh, models.Isian{TopID: "10007", Name: "uji flood", Alias: " uji banjir "}); b.Alias != "UJI BANJIR" {
		t.Errorf("alias %+v", b)
	}
}

// Edit: Treaty Group tidak diganti = TREATYNAME dibiarkan (juga TOPID yatim); diganti = disalin ulang.
func TestEditSalinanTreaty(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10014", TopID: "10007", Name: "uji earthquake", Alias: "uji gempa"}); err != nil {
		t.Fatal(err)
	}
	if b := g.Baris["10014"]; b.TreatyName != "UJI PROPERTY LAMA" || b.Alias != "UJI GEMPA" {
		t.Errorf("tidak diganti %+v", b)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10028", TopID: "10019", Name: "uji yatim"}); err != nil {
		t.Errorf("TOPID yatim tidak diganti: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10014", TopID: "10009", Name: "uji earthquake"}); err != nil {
		t.Fatal(err)
	}
	if b := g.Baris["10014"]; b.TopID != "10009" || b.TreatyName != "UJI MARINE CARGO" {
		t.Errorf("diganti %+v", b)
	}
}

// SYARIAH: tidak tampil, tidak dapat dibuka atau diubah, nama baru berakhiran SYARIAH ditolak.
func TestTanpaSyariah(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	d, _ := l.Daftar(ctx, "", "")
	for _, b := range d {
		if strings.HasSuffix(b.Name, "SYARIAH") {
			t.Errorf("SYARIAH tampil: %+v", b)
		}
	}
	if len(d) != 3 {
		t.Errorf("daftar %d, mau 3", len(d))
	}
	if _, err := l.Buka(ctx, "10020"); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("buka SYARIAH: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10020", TopID: "10007", Name: "UJI FIRE BARU"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("ubah SYARIAH: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{TopID: "10007", Name: "uji motor syariah "}); err == nil ||
		!strings.Contains(err.Error(), "ending with SYARIAH") {
		t.Errorf("nama SYARIAH: %v", err)
	}
	if g.Baris["10020"].Name != "UJI FIRE SYARIAH" {
		t.Error("baris SYARIAH berubah")
	}
}

func TestValidasi(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, models.Isian{Alias: strings.Repeat("A", models.BatasTeks+1)}); err == nil ||
		!strings.Contains(err.Error(), "Treaty Group is required") || !strings.Contains(err.Error(), "Name is required") ||
		!strings.Contains(err.Error(), "Alias Name is longer than 4000") {
		t.Errorf("wajib dan panjang: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{TopID: "19999", Name: "UJI X"}); err == nil || !strings.Contains(err.Error(), "Treaty Group 19999 is not in") {
		t.Errorf("treaty asing: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{TopID: "10007", Name: "uji fire syariah x"}); err != nil {
		t.Errorf("SYARIAH di tengah nama boleh: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{TopID: "10007", Name: "Uji Fire"}); err == nil || !strings.Contains(err.Error(), "already used by 10013") {
		t.Errorf("nama kembar: %v", err)
	}
	if _, err := l.Simpan(ctx, lihat, models.Isian{TopID: "10007", Name: "UJI Y"}); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("View only: %v", err)
	}
	if len(g.Baris) != 5 {
		t.Errorf("baris %d, mau 5 (satu Add sah)", len(g.Baris))
	}
}

func TestDaftarSaring(t *testing.T) {
	l := layanan(tiruan.Contoh())
	if d, _ := l.Daftar(ctx, "", "10019"); len(d) != 1 || d[0].ID != "10028" {
		t.Errorf("saring treaty %+v", d)
	}
	if d, _ := l.Daftar(ctx, "gempa", ""); len(d) != 0 {
		t.Errorf("cari %+v", d)
	}
	if d, _ := l.Daftar(ctx, "eq", ""); len(d) != 1 || d[0].ID != "10014" {
		t.Errorf("cari alias %+v", d)
	}
	if p, _ := l.Pilihan(ctx); len(p) != 2 {
		t.Errorf("pilihan %+v", p)
	}
}
