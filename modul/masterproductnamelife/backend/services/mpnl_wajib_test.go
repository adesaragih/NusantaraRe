package services_test

// Wajib-isi (paket 5, tiket 05) - `SaveProductName_Act` langkah 2–5, pesan VERBATIM
// (`local.errMsg1..4` b431/b452/b473/b494), SEMUA dilaporkan sekaligus.

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/services"
)

func TestWajibIsiEmpatPesanVerbatimSekaligus(t *testing.T) {
	l, g := layananMaster()
	m := produkMasuk()
	m.Umum.ProductName, m.Umum.Ceding, m.Umum.CedingID = "", "", ""
	m.Umum.SOBName, m.Umum.SOBID = "  ", ""
	m.Inward.PolicyHolderName, m.Inward.PolicyHolder = "", ""
	_, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Fatalf("wajib-isi 422: %v", err)
	}
	if got := services.Pesan(err); got != "Product Name Empty; Ceding Empty; Policy Holder Empty; SOB Empty" {
		t.Errorf("urutan langkah 2-5, VERBATIM, semua sekaligus: %q", got)
	}
	if g.Komit != 0 || len(g.Produk) != 0 {
		t.Error("nol tulisan")
	}
}

func TestWajibIsiBerlakuJugaSaatUbah(t *testing.T) {
	l, g := layananMaster()
	g.IsiJSON("100007", `{"ID":"100007","PRODUCTNAME":"LAMA","CEDING":"UJI CEDING","CEDINGID":"L0UJI"}`, "")
	m := produkMasuk()
	m.ID = "100007"
	m.Umum.ProductName = ""
	_, err := l.SimpanProduk(context.Background(), pelakuUji, m, false)
	if services.Pesan(err) != "Product Name Empty" {
		t.Errorf("satu aturan untuk baru dan ubah: %v", err)
	}
}

func TestWajibIsiDigabungDenganPenolakanLain(t *testing.T) {
	l, _ := layananMaster()
	m := produkMasuk()
	m.Umum.ProductName = ""
	m.Umum.RIComm = "x"
	pesan := services.Pesan(func() error { _, err := l.SimpanProduk(context.Background(), pelakuUji, m, true); return err }())
	if !strings.HasPrefix(pesan, "Product Name Empty; ") || !strings.Contains(pesan, "Deduction (%)") {
		t.Errorf("semua penolakan sekaligus, wajib-isi lebih dulu: %q", pesan)
	}
}
