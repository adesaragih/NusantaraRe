package services_test

// Regresi audit menyeluruh 02-10-2026 (permintaan work owner: "fokus perbaiki master productname").

import (
	"context"
	"strings"
	"testing"

	"nusantarare/modul/masterproductnamelife/backend/models"
)

// Keadaan DEV 02-10-2026: 45 dari 103 produk ber-UNDERWRITING LIMIT melampaui 4000 byte (terpanjang 84 baris; view
// PRODUCT_LIFE membacanya NULL sejak era Pega). Produk seperti itu TETAP dapat disimpan - sejak tabel flat
// (02-10-2026) setiap baris satu baris `M_PRODUCTNAME_LIFE_UWLIMIT`, tanpa batas 4000 byte larik.
func TestUWLimitPanjangSepertiDEVTetapDapatDisimpan(t *testing.T) {
	l, g := layananMaster()
	m := produkMasuk()
	for range 84 {
		m.UnderwritingLimit = append(m.UnderwritingLimit, models.BarisUWLimit{MinInsured: "100000000", MaxInsured: "250000000",
			MinAge: "18", MaxAge: "65", Medical: "UJI MEDIS", Description: strings.Repeat("D", 100)})
	}
	p, err := l.SimpanProduk(context.Background(), pelakuUji, m, true)
	if err != nil {
		t.Fatalf("produk ber-UW limit 84 baris tersimpan seperti di Pega: %v", err)
	}
	if len(g.Produk[p.ID].UnderwritingLimit) != 84 || len(p.UnderwritingLimit) != 84 {
		t.Errorf("84 baris UW limit tersimpan utuh: %d", len(p.UnderwritingLimit))
	}
}

// Keadaan DEV 02-10-2026: `M_PRODUCT_LIFE_SEQ` tertinggal dari data. Produk baru TIDAK lagi gagal 500 dan
// tidak menimpa produk yang ada: ID terpakai dilewati ke nomor bebas berikut.
func TestProdukBaruMelewatiIDTerpakaiSaatSequenceTertinggal(t *testing.T) {
	l, g := layananMaster()
	lama := g.IsiJSON("100044", `{"ID":"100044","PRODUCTNAME":"UJI LAMA"}`, "")
	g.IDWarisan["100045"] = true // ID baris JSON warisan (tabel inward) - juga tidak diterbitkan ulang
	p, err := l.SimpanProduk(context.Background(), pelakuUji, produkMasuk(), true)
	if err != nil {
		t.Fatalf("produk baru tersimpan walau sequence tertinggal: %v", err)
	}
	if p.ID != "100046" {
		t.Errorf("ID terpakai di salah satu tabel dilewati: dapat %q, mau 100046", p.ID)
	}
	if g.Produk["100044"].Umum.ProductName != lama.Umum.ProductName || len(g.Produk) != 2 {
		t.Error("produk lama ber-ID nomor sequence tidak tersentuh")
	}
}
