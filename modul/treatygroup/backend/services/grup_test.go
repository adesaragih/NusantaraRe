package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/treatygroup/backend/models"
	"nusantarare/modul/treatygroup/backend/services"
	"nusantarare/modul/treatygroup/backend/tiruan"
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

// Add: OJK disalin dari TREATYGROUPOJK; COAID dari grup lain se-OJK (opsi A); ID situs + sequence, nomor terpakai
// dilompati; TGLUPDATE format Pega GMT; USERID akun; nama huruf besar.
func TestAddSalinOjkDanCoa(t *testing.T) {
	g := tiruan.Contoh()
	d, err := layanan(g).Simpan(ctx, penuh, models.Isian{OjkID: "01", Name: " uji livestock ", SoaName: " uji ls "})
	if err != nil {
		t.Fatal(err)
	}
	r := g.Baris[d.ID]
	mau := models.Grup{ID: "10010", OjkID: "01", OjkName: "UJI PROPERTY", OjkNameIDN: "UJI HARTA BENDA", OrderNo: "1",
		Name: "UJI LIVESTOCK", SoaName: "UJI LS", TglUpdate: "20261005T030405.600 GMT", UserID: "UJI-ADMIN", CoaID: "10013"}
	if r != mau {
		t.Errorf("tersimpan\n%+v\nmau\n%+v", r, mau)
	}
	if d.CoaName != "UJI FIRE" || d.Diubah != "05-10-2026 10:04" {
		t.Errorf("tampilan %+v", d.Grup)
	}
	// OJK tanpa grup lain: COAID kosong.
	if d, err := layanan(g).Simpan(ctx, penuh, models.Isian{OjkID: "03", Name: "UJI CARGO"}); err != nil || g.Baris[d.ID].CoaID != "" {
		t.Errorf("OJK baru: %+v %v", d, err)
	}
}

// Edit: OJK tidak diganti = salinan OJK dan COAID lama dibiarkan (juga yang sudah beda dari TREATYGROUPOJK); OJK diganti
// = salinan OJK disalin ulang dan COAID ikut OJK baru (perintah work owner 05-10-2026: "ubah pas pilih OJK Business"),
// kosong bila OJK itu belum punya grup ber-COAID. OLDID tidak pernah disentuh.
func TestEditOjkDanCoa(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10009", OjkID: "09", Name: "uji engineering baru"}); err != nil {
		t.Fatal(err)
	}
	if r := g.Baris["10009"]; r.OjkName != "UJI ENGINEERING LAMA" || r.Name != "UJI ENGINEERING BARU" || r.CoaID != "10019" {
		t.Errorf("OJK tidak diganti %+v", r)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10009", OjkID: "01", Name: "uji engineering baru"}); err != nil {
		t.Fatal(err)
	}
	if r := g.Baris["10009"]; r.OjkID != "01" || r.OjkName != "UJI PROPERTY" || r.CoaID != "10013" {
		t.Errorf("OJK diganti ke OJK ber-COA %+v", r)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "10007", OjkID: "03", Name: "UJI PROPERTY"}); err != nil {
		t.Fatal(err)
	}
	if r := g.Baris["10007"]; r.OjkID != "03" || r.OjkName != "UJI CARGO" || r.OrderNo != "5" || r.CoaID != "" || r.OldID != "01" {
		t.Errorf("OJK diganti ke OJK tanpa COA %+v", r)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{ID: "19999", OjkID: "01", Name: "X"}); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("edit ID tak ada: %v", err)
	}
}

func TestValidasi(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	if _, err := l.Simpan(ctx, penuh, models.Isian{SoaName: strings.Repeat("A", models.BatasTeks+1)}); err == nil ||
		!strings.Contains(err.Error(), "OJK Business is required") || !strings.Contains(err.Error(), "Treaty Group Name is required") ||
		!strings.Contains(err.Error(), "SOA Name is longer than 1000") {
		t.Errorf("wajib dan panjang: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{OjkID: "77", Name: "UJI X"}); !errors.Is(err, services.ErrMasukanTidakSah) ||
		!strings.Contains(err.Error(), "OJK Business 77 is not in") {
		t.Errorf("OJK asing: %v", err)
	}
	if _, err := l.Simpan(ctx, penuh, models.Isian{OjkID: "01", Name: "uji agriculture"}); err == nil ||
		!strings.Contains(err.Error(), "already used by 10056") {
		t.Errorf("nama kembar: %v", err)
	}
	if _, err := l.Simpan(ctx, lihat, models.Isian{OjkID: "01", Name: "UJI Y"}); !errors.Is(err, services.ErrDilarang) {
		t.Errorf("View only: %v", err)
	}
	if len(g.Baris) != 3 {
		t.Error("ditolak tetapi tertulis")
	}
}

// View: grup bisnis anak tanpa yang berakhiran SYARIAH; daftar urut Order No dan tersaring OJK.
func TestBukaDanDaftar(t *testing.T) {
	l := layanan(tiruan.Contoh())
	d, err := l.Buka(ctx, "10007")
	if err != nil || len(d.Anak) != 2 || d.Anak[0].ID != "10014" || d.Anak[1].ID != "10013" || d.Diubah != "05-01-2022 20:56" {
		t.Errorf("buka %+v %v", d, err)
	}
	if _, err := l.Buka(ctx, "19999"); !errors.Is(err, services.ErrTidakAda) {
		t.Errorf("tak ada: %v", err)
	}
	daftar, _ := l.Daftar(ctx, "", "")
	if len(daftar) != 3 || daftar[0].ID != "10056" || daftar[2].ID != "10009" {
		t.Errorf("urutan %+v", daftar)
	}
	if daftar, _ := l.Daftar(ctx, "", "09"); len(daftar) != 1 || daftar[0].ID != "10009" {
		t.Errorf("saring OJK %+v", daftar)
	}
	// Pilihan OJK membawa COA OJK itu (ditampilkan form begitu OJK dipilih): yang paling sering, lalu terkecil.
	p, _ := l.Pilihan(ctx)
	coa := map[string]string{}
	for _, o := range p {
		coa[o.ID] = o.CoaID + "|" + o.CoaName
	}
	if len(p) != 3 || coa["01"] != "10013|UJI FIRE" || coa["09"] != "10019|UJI CAR" || coa["03"] != "|" {
		t.Errorf("pilihan %+v", p)
	}
}

func TestTampilWaktu(t *testing.T) {
	if s := services.TampilWaktu("20220105T235959.000 GMT"); s != "06-01-2022 06:59" {
		t.Errorf("WIB %q", s)
	}
	if s := services.TampilWaktu("lain"); s != "lain" {
		t.Errorf("bentuk lain %q", s)
	}
}
