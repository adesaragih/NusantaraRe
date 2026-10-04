package services_test

// Copy Old (permintaan work owner 03-10-2026): popup daftar produk lama yang belum ada di tabel flat, `Process Copy`
// menyalin yang dicentang - per produk, satu transaksi per produk.

import (
	"context"
	"errors"
	"fmt"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

func TestDaftarProdukLamaHanyaYangBelumDiFlat(t *testing.T) {
	l, g := layananUji()
	g.IsiLama("100044", `{"ID":"100044","PRODUCTNAME":"UJI LAMA SATU","CEDING":"UJI CEDING"}`, "")
	g.IsiLama("100045", `{"ID":"100045","PRODUCTNAME":"UJI LAMA DUA"}`, "")
	g.IsiJSON("100045", `{"ID":"100045","PRODUCTNAME":"UJI LAMA DUA"}`, "") // sudah disalin
	d, err := l.DaftarProdukLama(context.Background(), pelakuUji)
	if err != nil {
		t.Fatal(err)
	}
	if len(d) != 1 || d[0].ID != "100044" || d[0].ProductName != "UJI LAMA SATU" || d[0].Ceding != "UJI CEDING" || !d[0].BolehDisalin {
		t.Errorf("hanya produk lama yang belum ada di tabel flat: %+v", d)
	}
	if _, err := l.DaftarProdukLama(context.Background(), inti.Pelaku{}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
}

func TestSalinProdukLamaPerProduk(t *testing.T) {
	l, g := layananUji()
	g.IsiLama("100044", `{"ID":"100044","PRODUCTNAME":"UJI LAMA SATU"}`, `{"ID":"100044","PRODUCTID":"100044","INSURED":"UJI INSURED"}`)
	g.IsiLama("100045", `{"ID":"100045","PRODUCTNAME":"UJI DITOLAK"}`, "")
	g.TolakLama("100045", "TYPE terisi tetapi tidak punya kolom flat")
	g.IsiLama("100046", `{"ID":"100046","PRODUCTNAME":"UJI GAGAL DB"}`, "")
	g.GagalSalinLama["100046"] = errors.New("UJI ORA-03113")
	g.IsiLama("100047", `{"ID":"100047","PRODUCTNAME":"UJI LAMA DUA"}`, "")

	j, err := l.SalinProdukLama(context.Background(), pelakuUji, []string{"100044", "100045", "100046", "100047", "100044", " ", "100099"})
	if err != nil {
		t.Fatal(err)
	}
	mau := []struct{ id, status string }{
		{"100044", models.SalinDisalin}, {"100045", models.SalinDitolak}, {"100046", models.SalinGagal},
		{"100047", models.SalinDisalin}, {"100099", models.SalinSudahAda},
	}
	if len(j.Hasil) != len(mau) || j.Disalin != 2 {
		t.Fatalf("hasil per ID unik, urutan permintaan: %+v", j)
	}
	for i, m := range mau {
		if j.Hasil[i].ID != m.id || j.Hasil[i].Status != m.status {
			t.Errorf("hasil %d: %+v, mau %s %s", i, j.Hasil[i], m.id, m.status)
		}
	}
	if len(j.Hasil[1].Pesan) != 1 || j.Hasil[1].Pesan[0] != "TYPE terisi tetapi tidak punya kolom flat" {
		t.Errorf("alasan penolakan disampaikan: %+v", j.Hasil[1])
	}
	// Tertulis di tabel flat (tiruan) - utuh, termasuk sisi inward; yang ditolak/gagal tidak.
	if p, ada := g.Produk["100044"]; !ada || p.Inward.Insured != "UJI INSURED" || p.Umum.ProductName != "UJI LAMA SATU" {
		t.Errorf("produk lama tersalin utuh: %+v", p)
	}
	for _, id := range []string{"100045", "100046"} {
		if _, ada := g.Produk[id]; ada {
			t.Errorf("%s tidak boleh tertulis", id)
		}
	}
	// Sesudah disalin, produk itu tidak lagi ada di daftar; salin ulang = sudah ada.
	j2, err := l.SalinProdukLama(context.Background(), pelakuUji, []string{"100044"})
	if err != nil || len(j2.Hasil) != 1 || j2.Hasil[0].Status != models.SalinSudahAda || j2.Disalin != 0 {
		t.Errorf("salin ulang: %+v %v", j2, err)
	}
}

func TestSalinProdukLamaValidasi(t *testing.T) {
	l, _ := layananUji()
	var v services.GalatValidasi
	if _, err := l.SalinProdukLama(context.Background(), pelakuUji, []string{" ", ""}); !errors.As(err, &v) {
		t.Errorf("tanpa ID terpilih: %v", err)
	}
	banyak := make([]string, services.MaksSalinLama+1)
	for i := range banyak {
		banyak[i] = fmt.Sprintf("U%05d", i)
	}
	if _, err := l.SalinProdukLama(context.Background(), pelakuUji, banyak); !errors.As(err, &v) {
		t.Errorf("lebih dari MaksSalinLama: %v", err)
	}
	if _, err := l.SalinProdukLama(context.Background(), inti.Pelaku{}, []string{"100044"}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
}

// Penerbitan ID produk baru melewati ID produk lama yang belum disalin - produk lama itu tetap dapat disalin dengan ID-nya.
func TestProdukBaruMelewatiIDProdukLama(t *testing.T) {
	l, g := layananMaster()
	g.IsiLama("100044", `{"ID":"100044","PRODUCTNAME":"UJI LAMA"}`, "")
	p, err := l.SimpanProduk(context.Background(), pelakuUji, produkMasuk(), true)
	if err != nil || p.ID != "100045" {
		t.Fatalf("ID lama dilewati: %q %v", p.ID, err)
	}
	if j, err := l.SalinProdukLama(context.Background(), pelakuUji, []string{"100044"}); err != nil || j.Disalin != 1 {
		t.Errorf("produk lama tetap dapat disalin dengan ID-nya: %+v %v", j, err)
	}
}
