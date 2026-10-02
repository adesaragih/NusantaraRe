package services_test

// Sisi inward + simpan atomik (paket 4, tiket 03) - sejak 02-10-2026 satu baris induk flat + tujuh anak.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/services"
	"nusantarare/modul/masterproductnamelife/backend/tiruan"
)

func TestSimpanBaruMenulisInwardBerIDProduk(t *testing.T) {
	l, g := layananMaster()
	m := produkMasuk()
	m.Inward.CedingLimit, m.Inward.MaxExpiredClaim = "150000000,5", "180"
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil {
		t.Fatal(err)
	}
	if p.Inward.ID != "100044" || p.Inward.ProductID != "100044" {
		t.Errorf("R14: ID inward = PRODUCTID = ID produk: %+v", p.Inward)
	}
	s := g.Produk["100044"].Inward
	if s.ProductID != "100044" || s.CedingLimit != "150000000.5" || s.Begin != "2026-03-01" || s.PolicyHolder != "UJI-ORG-1" ||
		s.MaxExpiredClaim != "180" {
		t.Errorf("sisi inward di baris induk flat: %+v", s)
	}
	if g.Komit != 1 {
		t.Errorf("induk dan anak satu transaksi: %d komit", g.Komit)
	}
}

// Produk lama yang baris inward JSON-nya ber-ID lain (100009 → PRODUCTID 100007) dipindah alat pindah ke baris
// produknya sendiri: ID inward = ID produk (tiket 01 bab 02-10-2026). Medan inward tanpa form (tanpa kolom flat)
// tetap tidak pernah dari klien.
func TestUbahProdukHasilPindahInwardBerIDProduk(t *testing.T) {
	l, g := layananMaster()
	g.IsiJSON("100007", `{"ID":"100007","PRODUCTNAME":"LAMA","CREATEOP":"UJI-A"}`,
		`{"ID":"100009","PRODUCTID":"100007","EXPIRYAGE":"75"}`)
	m := produkMasuk()
	m.ID = "100007"
	m.Inward.Ceding = "KLIEN" // medan tanpa form - tidak dari klien
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Produk) != 1 || p.Inward.ID != "100007" || p.Inward.ProductID != "100007" || p.Inward.Ceding != "" ||
		g.Produk["100007"].Inward.Ceding != "" {
		t.Errorf("satu baris produk, inward ber-ID produk, medan tanpa form bukan dari klien: %+v", p.Inward)
	}
}

func TestUbahProdukTanpaInwardMenyisipInwardBerIDProduk(t *testing.T) {
	l, g := layananMaster()
	g.IsiJSON("100007", `{"ID":"100007","PRODUCTNAME":"LAMA"}`, "")
	m := produkMasuk()
	m.ID = "100007"
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, m, false); err != nil {
		t.Fatal(err)
	}
	if s := g.Produk["100007"].Inward; s.ProductID != "100007" || s.PolicyHolder != "UJI-ORG-1" {
		t.Errorf("sisi inward produk tanpa inward lama ditulis ber-ID produk: %+v", s)
	}
}

func TestSimpanAtomikIndukDanAnak(t *testing.T) {
	l, g := layananMaster()
	g.GagalTulisAnak = tiruan.ErrTiruan
	_, err := l.SimpanProduk(context.Background(), pelakuUji, produkMasuk(), true)
	if !errors.Is(err, tiruan.ErrTiruan) {
		t.Fatalf("galat penulis anak diteruskan: %v", err)
	}
	// `NEXTVAL` tidak ikut rollback (Oracle): nomor 44 terpakai, celahnya tetap - seperti DEV.
	if len(g.Produk) != 0 || g.Komit != 0 || g.Seq != 45 {
		t.Errorf("P4: gagal menulis anak membatalkan baris induk juga: produk %d komit %d seq %d",
			len(g.Produk), g.Komit, g.Seq)
	}
}

func TestInwardAngkaTanggalDanRentang(t *testing.T) {
	l, _ := layananMaster()
	m := produkMasuk()
	m.Inward.MinAge, m.Inward.MaxAge = "70", "22"
	m.Inward.MinSumInsured, m.Inward.MaxSumInsured = "2000000000", "1175000000"
	m.Inward.Begin = "31/02/2026"
	m.Inward.Brokerage = "x"
	_, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	pesan := services.Pesan(err)
	for _, w := range []string{"Minimum Age (Years) must not be greater than Maximum Age (Years)",
		"Min Sum Insured must not be greater than Max Sum Insured", `Begin Date "31/02/2026" is not a date`,
		`Brokerage Fee (%) "x" is not a number`} {
		if !strings.Contains(pesan, w) {
			t.Errorf("penolakan tanpa %q: %s", w, pesan)
		}
	}
	m = produkMasuk()
	m.Inward.Begin, m.Inward.Mature = "01/03/2026", "2027-02-28"
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil || p.Inward.Begin != "2026-03-01" || p.Inward.Mature != "2027-02-28" {
		t.Errorf("tanggal dd/MM/yyyy dan YYYY-MM-DD diterima: %+v %v", p.Inward, err)
	}
}

func TestPilihanInwardDiverifikasi(t *testing.T) {
	l, _ := layananMaster()
	m := produkMasuk()
	m.Inward.PolicyHolder = "UJI-TIDAK-ADA"
	m.Inward.CurrencyID, m.Inward.Currency = "1", "DIKETIK"
	_, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if !errors.Is(err, services.ErrMasukanTidakSah) || !strings.Contains(services.Pesan(err), `Policy Holder "UJI-TIDAK-ADA"`) {
		t.Errorf("pemegang polis di luar master: %v", err)
	}
	m = produkMasuk()
	m.Inward.Currency = "DIKETIK"
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil || p.Inward.Currency != "IDR" || p.Umum.PolicyHolderName != "UJI PEMEGANG" {
		t.Errorf("nama mata uang dari master, salinan pemegang polis ke umum: %+v %v", p.Inward, err)
	}
}
