package services_test

// Aturan Accounts TANPA Oracle (gudang tiruan). Sumber setiap aturan: dokumentasi paket services.

import (
	"context"
	"errors"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/accounts/backend/models"
	"nusantarare/modul/accounts/backend/services"
	"nusantarare/modul/accounts/backend/tiruan"
)

var (
	ctx   = context.Background()
	admin = inti.Pelaku{AkunID: "UJI-ADMIN"}
)

const (
	org1 = models.AwalanOrg + "ORG-101"
	org2 = models.AwalanOrg + "ORG-102"
)

func layanan(g *tiruan.Gudang) *services.Layanan { return services.BaruLayanan(g, tiruan.Transaksi) }

// Akun baru: ID dari SEQ_T_M_ACCOUNT, nama organisasi dan NOTE Group Business disalin dari tabelnya, Owner = akun
// pelaku, Create Date dari "SYSDATE".
func TestTambahAkunBaru(t *testing.T) {
	g := tiruan.Contoh()
	a, err := layanan(g).Tambah(ctx, admin, models.Isian{InsuredID: " " + org1 + " ", GroupBusinessID: "10002",
		Description: "  UJI catatan  "})
	if err != nil {
		t.Fatal(err)
	}
	mau := models.Account{ID: models.AwalanID + "ACC-4916562", IDView: "ACC-4916562", GroupBusinessID: "10002",
		GroupBusiness: "UJI GROUP B", InsuredID: org1, OrgID: "ORG-101", InsuredName: "UJI PT SATU",
		Description: "UJI catatan", CreateOp: "UJI-ADMIN", CreateDate: "2026-10-04 09:00"}
	if a != mau {
		t.Errorf("akun baru:\n dapat %+v\n mau   %+v", a, mau)
	}
}

// Nomor sequence yang sudah terpakai (mis. akun dibuat sebelum migrasi 842) dilewati.
func TestNomorTerpakaiDilewati(t *testing.T) {
	g := tiruan.Contoh()
	g.Nomor = 4916560
	a, err := layanan(g).Tambah(ctx, admin, models.Isian{InsuredID: org2, GroupBusinessID: "10001"})
	if err != nil || a.IDView != "ACC-4916562" {
		t.Errorf("melewati 4916560 dan 4916561: %+v %v", a, err)
	}
}

// Wajib isi VERBATIM tangkapan layar ("Value cannot be blank"), keduanya sekaligus; pilihan di luar tabelnya ditolak.
func TestIsianDitolak(t *testing.T) {
	l := layanan(tiruan.Contoh())
	_, err := l.Tambah(ctx, admin, models.Isian{InsuredID: " ", GroupBusinessID: ""})
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), "Insured Name: Value cannot be blank") ||
		!strings.Contains(err.Error(), "Group Business: Value cannot be blank") {
		t.Errorf("wajib isi: %v", err)
	}
	for _, k := range []struct {
		isi   models.Isian
		pesan string
	}{
		{models.Isian{InsuredID: "UJI-ORANG-1", GroupBusinessID: "10001"}, "Insured Name: the organization is not in"},
		{models.Isian{InsuredID: org1, GroupBusinessID: "99999"}, "Group Business: not in the business group list"},
		{models.Isian{InsuredID: org2, GroupBusinessID: "10001", Description: strings.Repeat("x", 4001)}, "longer than 4000"},
	} {
		if _, err := l.Tambah(ctx, admin, k.isi); !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(err.Error(), k.pesan) {
			t.Errorf("%+v: %v", k.isi, err)
		}
	}
}

// Pasangan Insured + Group Business yang sudah dimiliki akun yang ada ditolak (409).
func TestPasanganGandaDitolak(t *testing.T) {
	g := tiruan.Contoh()
	l := layanan(g)
	_, err := l.Tambah(ctx, admin, models.Isian{InsuredID: org1, GroupBusinessID: "10001"})
	if !errors.Is(err, services.ErrGanda) || !strings.Contains(err.Error(), "UJI PT SATU already has an account in UJI GROUP A (ACC-4916560)") {
		t.Errorf("tambah ganda: %v", err)
	}
	if len(g.Disisip) != 0 {
		t.Errorf("tidak ada baris baru: %d", len(g.Disisip))
	}
}

func TestTanpaPelakuDitolak(t *testing.T) {
	l := layanan(tiruan.Contoh())
	if _, err := l.Tambah(ctx, inti.Pelaku{}, models.Isian{InsuredID: org1, GroupBusinessID: "10002"}); !errors.Is(err, services.ErrTanpaPelaku) {
		t.Errorf("tambah: %v", err)
	}
	if _, err := l.Tambah(ctx, inti.Pelaku{AkunID: " "}, models.Isian{InsuredID: org1, GroupBusinessID: "10002"}); !errors.Is(err, services.ErrTanpaPelaku) {
		t.Errorf("tambah berpelaku spasi: %v", err)
	}
}

// Daftar: terbaru dulu, cari tanpa beda huruf di nama, Group Business, ACC-n, ORG-n; halaman dan ukuran dibatasi.
func TestDaftarCariDanHalaman(t *testing.T) {
	l := layanan(tiruan.Contoh())
	h, err := l.Daftar(ctx, "", 0, 0)
	if err != nil || h.Total != 2 || h.Halaman != 1 || h.Ukuran != services.UkuranBawaan || h.Daftar[0].IDView != "ACC-4916561" {
		t.Fatalf("daftar: %+v %v", h, err)
	}
	for kata, mau := range map[string]int{"uji pt dua": 1, "group a": 1, "acc-4916560": 1, "org-102": 1, "TIDAK ADA": 0} {
		if h, _ := l.Daftar(ctx, kata, 1, 20); h.Total != mau {
			t.Errorf("cari %q: %d, mau %d", kata, h.Total, mau)
		}
	}
	if h, _ := l.Daftar(ctx, "", 2, 1); len(h.Daftar) != 1 || h.Daftar[0].IDView != "ACC-4916560" {
		t.Errorf("halaman 2 ukuran 1: %+v", h.Daftar)
	}
	if h, _ := l.Daftar(ctx, "", 1, 5000); h.Ukuran != services.UkuranBawaan {
		t.Errorf("ukuran di atas batas: %d", h.Ukuran)
	}
}

// Pilihan Insured Name dibatasi; `lebih` menyatakan ada yang tidak ditampilkan.
func TestCariOrganisasiTerbatas(t *testing.T) {
	g := tiruan.Contoh()
	for i := 0; i < services.BatasOrganisasi+5; i++ {
		id := models.AwalanOrg + "ORG-9" + strings.Repeat("0", 3) + string(rune('A'+i%26)) + strings.Repeat("1", i/26)
		g.Org[id] = models.Organisasi{ID: id, IDView: "ORG-X", Nama: "UJI BANYAK " + id}
	}
	h, err := layanan(g).CariOrganisasi(ctx, "banyak")
	if err != nil || len(h.Daftar) != services.BatasOrganisasi || !h.Lebih {
		t.Errorf("terbatas: %d %v %v", len(h.Daftar), h.Lebih, err)
	}
	h, _ = layanan(g).CariOrganisasi(ctx, "pt dua")
	if len(h.Daftar) != 1 || h.Lebih || h.Daftar[0].IDView != "ORG-102" {
		t.Errorf("satu cocok: %+v", h)
	}
}
