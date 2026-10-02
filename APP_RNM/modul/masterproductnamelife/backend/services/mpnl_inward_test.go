package services_test

// Sisi inward + simpan atomik dua tabel (paket 4, tiket 03).

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
	isi := g.Inward["100044"]
	for _, w := range []string{`"PRODUCTID":"100044"`, `"CEDINGLIMIT":"150000000.5"`, `"BEGIN":"01/03/2026"`, `"POLICYHODER":"UJI-ORG-1"`} {
		if !strings.Contains(isi, w) {
			t.Errorf("JSON inward tanpa %s: %s", w, isi)
		}
	}
	if g.Komit != 1 {
		t.Errorf("kedua tabel satu transaksi: %d komit", g.Komit)
	}
}

func TestUbahMemperbaruiBarisInwardLamaBerIDSendiri(t *testing.T) {
	l, g := layananMaster()
	g.Umum["100007"] = `{"ID":"100007","PRODUCTNAME":"LAMA","CREATEOP":"UJI-A"}`
	g.Inward["100009"] = `{"ID":"100009","PRODUCTID":"100007","CEDING":"UJI CEDING LAMA","EXPIRYAGE":"75"}`
	m := produkMasuk()
	m.ID = "100007"
	m.Inward.Ceding = "KLIEN" // medan tanpa form - tidak dari klien
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Inward) != 1 || p.Inward.ID != "100009" || p.Inward.ProductID != "100007" || p.Inward.Ceding != "UJI CEDING LAMA" {
		t.Errorf("baris inward lama diperbarui, tidak digandakan: %+v %v", p.Inward, g.Inward)
	}
}

func TestUbahProdukTanpaInwardMenyisipInwardBerIDProduk(t *testing.T) {
	l, g := layananMaster()
	g.Umum["100007"] = `{"ID":"100007","PRODUCTNAME":"LAMA"}`
	m := produkMasuk()
	m.ID = "100007"
	if _, err := l.SimpanProduk(context.Background(), pelakuUji, m, false); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(g.Inward["100007"], `"PRODUCTID":"100007"`) {
		t.Errorf("inward baru ber-ID produk: %v", g.Inward)
	}
}

func TestSimpanAtomikDuaTabel(t *testing.T) {
	l, g := layananMaster()
	g.GagalTulisInward = tiruan.ErrTiruan
	_, err := l.SimpanProduk(context.Background(), pelakuUji, produkMasuk(), true)
	if !errors.Is(err, tiruan.ErrTiruan) {
		t.Fatalf("galat inward diteruskan: %v", err)
	}
	// `NEXTVAL` tidak ikut rollback (Oracle): nomor 44 terpakai, celahnya tetap - seperti DEV.
	if len(g.Umum) != 0 || len(g.Inward) != 0 || g.Komit != 0 || g.Seq != 45 {
		t.Errorf("P4: gagal sisi inward membatalkan sisi umum juga: umum %d inward %d komit %d seq %d",
			len(g.Umum), len(g.Inward), g.Komit, g.Seq)
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
