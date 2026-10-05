package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/treatydescription/backend/models"
	"nusantarare/modul/treatydescription/backend/services"
	"nusantarare/modul/treatydescription/backend/tiruan"
)

var ctx = context.Background()

var (
	penuh = services.Aktor{AkunID: "UJI-ADMIN", Penuh: true}
	lihat = services.Aktor{AkunID: "UJI-LIHAT"}
)

func layanan(g *tiruan.Gudang) *services.Layanan { return services.BaruLayanan(g, tiruan.Transaksi) }

// Add: ID '1' + TREATY_DESCRIPTION_SEQ 4 digit seperti PEGA_TREATYDESC, nomor yang sudah terpakai dilompati; nama huruf
// besar tanpa spasi tepi; jenis dan status tersimpan apa adanya ("0"/"1").
func TestAddIDDariSequence(t *testing.T) {
	g := tiruan.Contoh()
	d, err := layanan(g).Simpan(ctx, penuh, models.Isian{DescName: "  uji profit commission ", IsXOL: "1", StatusAktif: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if mau := (models.Desc{ID: "10018", DescName: "UJI PROFIT COMMISSION", IsXOL: "1", StatusAktif: "1"}); g.Baris["10018"] != mau {
		t.Errorf("tersimpan %+v, mau %+v", g.Baris["10018"], mau)
	}
	if d.ID != "10018" || !d.Aktif || g.Seq != 19 {
		t.Errorf("hasil %+v, sequence berikut %d", d, g.Seq)
	}
	// Status Inactive tersimpan "0".
	if d, err := layanan(g).Simpan(ctx, penuh, models.Isian{DescName: "UJI LAIN", IsXOL: "0", StatusAktif: "0"}); err != nil ||
		d.ID != "10019" || d.Aktif || g.Baris["10019"].StatusAktif != "0" {
		t.Errorf("inactive %+v %v", d, err)
	}
}

// Edit: nama, jenis, dan status ditulis; ID tetap; baris lama berstatus NULL menjadi "1" bila disimpan Active.
func TestEdit(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	d, err := l.Simpan(ctx, penuh, models.Isian{ID: "10001", DescName: "uji treaty limit baru", IsXOL: "0", StatusAktif: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if r := g.Baris["10001"]; r.DescName != "UJI TREATY LIMIT BARU" || r.StatusAktif != "1" || r.IsXOL != "0" {
		t.Errorf("tersimpan %+v", r)
	}
	if d.ID != "10001" || !d.Aktif {
		t.Errorf("hasil %+v", d)
	}
	// Nama sendiri boleh tetap (tidak dianggap kembar dengan dirinya), status nonaktif.
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10004", DescName: "UJI CASH LOSS LIMIT", IsXOL: "1", StatusAktif: "0"}); err != nil {
		t.Fatal(err)
	}
	if r := g.Baris["10004"]; r.IsXOL != "1" || r.StatusAktif != "0" {
		t.Errorf("jenis/status %+v", r)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "19999", DescName: "X", IsXOL: "0", StatusAktif: "1"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("edit ID tak ada: %v", err)
	}
}

func TestValidasi(t *testing.T) {
	for _, k := range []struct {
		nama  string
		isi   models.Isian
		pesan string
	}{
		{"nama kosong", models.Isian{DescName: " ", IsXOL: "0", StatusAktif: "1"}, "Description Name is required"},
		{"nama kembar tanpa beda huruf", models.Isian{DescName: "uji exclusion treaty", IsXOL: "0", StatusAktif: "1"},
			"Description Name is already used by 10013"},
		{"nama kembar dengan baris nonaktif", models.Isian{DescName: "UJI XOL LAYER", IsXOL: "1", StatusAktif: "1"},
			"already used by 10015"},
		{"nama terlalu panjang", models.Isian{DescName: strings.Repeat("A", 101), IsXOL: "0", StatusAktif: "1"},
			"longer than 100 characters"},
		{"jenis asing", models.Isian{DescName: "UJI", IsXOL: "2", StatusAktif: "1"}, "Type must be Non XOL or XOL"},
		{"status asing", models.Isian{DescName: "UJI", IsXOL: "0", StatusAktif: ""}, "Status must be Active or Inactive"},
	} {
		g := tiruan.Contoh()
		_, err := layanan(g).Simpan(ctx, penuh, k.isi)
		if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), k.pesan) {
			t.Errorf("%s: %v, mau %q", k.nama, err, k.pesan)
		}
		if len(g.Baris) != 5 || g.Seq != 17 {
			t.Errorf("%s: ditolak tetapi menulis atau mengambil nomor", k.nama)
		}
	}
	// Nama 100 byte pas diterima.
	if _, err := layanan(tiruan.Contoh()).Simpan(ctx, penuh, models.Isian{DescName: strings.Repeat("B", 100), IsXOL: "0", StatusAktif: "1"}); err != nil {
		t.Errorf("100 byte: %v", err)
	}
}

// View only dan tanpa akun tidak boleh menulis.
func TestHakTulis(t *testing.T) {
	g := tiruan.Contoh()
	isi := models.Isian{DescName: "UJI BARU", IsXOL: "0", StatusAktif: "1"}
	if _, err := layanan(g).Simpan(ctx, lihat, isi); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("View only: %v", err)
	}
	if _, err := layanan(g).Simpan(ctx, services.Aktor{Penuh: true}, isi); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("tanpa akun: %v", err)
	}
	if len(g.Baris) != 5 {
		t.Error("ditolak tetapi menulis")
	}
}

// Daftar: cari ID atau nama; jenis (Non XOL/XOL); status NULL dihitung Active; nilai saringan asing = tanpa saringan.
func TestDaftar(t *testing.T) {
	l := layanan(tiruan.Contoh())
	id := func(d []models.Desc) string {
		var s []string
		for _, r := range d {
			s = append(s, r.ID)
		}
		return strings.Join(s, ",")
	}
	for _, k := range []struct{ kata, xol, status, mau string }{
		{"", "", "", "10001,10004,10013,10015,10017"},
		{"limit", "", "", "10001,10004"},
		{"1001", "", "", "10013,10015,10017"},
		{"", "1", "", "10015"},
		{"", "", "1", "10001,10004,10013,10017"},
		{"", "", "0", "10015"},
		{"", "x", "y", "10001,10004,10013,10015,10017"},
	} {
		d, err := l.Daftar(ctx, k.kata, k.xol, k.status)
		if err != nil || id(d) != k.mau {
			t.Errorf("cari %q xol %q status %q: %s %v, mau %s", k.kata, k.xol, k.status, id(d), err, k.mau)
		}
	}
	if d, _ := l.Buka(ctx, "10001"); !d.Aktif || d.StatusAktif != "" {
		t.Errorf("baris berstatus NULL: %+v", d)
	}
	if _, err := l.Buka(ctx, "19999"); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("buka tak ada: %v", err)
	}
}
