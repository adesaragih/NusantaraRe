package services_test

// Regresi audit menyeluruh 02-10-2026 (permintaan work owner: "fokus perbaiki master productname").

import (
	"context"
	"testing"
)

// Keadaan DEV 02-10-2026: `M_PRODUCT_LIFE_SEQ` tertinggal dari data. Produk baru TIDAK lagi gagal 500 dan
// tidak menimpa produk yang ada: ID terpakai dilewati ke nomor bebas berikut.
func TestProdukBaruMelewatiIDTerpakaiSaatSequenceTertinggal(t *testing.T) {
	l, g := layananMaster()
	asli := `{"ID":"100044","UJI":"milik produk lama"}`
	g.Umum["100044"] = asli
	g.Inward["100045"] = `{"ID":"100045","PRODUCTID":"100045"}`
	p, err := l.SimpanProduk(context.Background(), pelakuUji, produkMasuk(), true)
	if err != nil {
		t.Fatalf("produk baru tersimpan walau sequence tertinggal: %v", err)
	}
	if p.ID != "100046" {
		t.Errorf("ID terpakai di salah satu tabel dilewati: dapat %q, mau 100046", p.ID)
	}
	if g.Umum["100044"] != asli {
		t.Error("produk lama ber-ID nomor sequence tidak tersentuh")
	}
}
