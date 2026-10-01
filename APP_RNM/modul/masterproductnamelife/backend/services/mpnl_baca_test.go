package services_test

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/masterproductnamelife/backend/services"
	"nusantarare/modul/masterproductnamelife/backend/tiruan"
)

var pelakuUji = inti.Pelaku{AkunID: "UJI-PELAKU"}

func layananUji() (*services.Layanan, *tiruan.Gudang) {
	g := tiruan.Baru()
	return services.BaruLayanan(g, g.Transaksi, nil), g
}

func TestBacaMenuntutIdentitas(t *testing.T) {
	l, _ := layananUji()
	if _, err := l.DaftarProduk(context.Background(), inti.Pelaku{}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("daftar tanpa identitas: %v", err)
	}
	if _, err := l.AmbilProduk(context.Background(), inti.Pelaku{}, "UJI-1"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("ambil tanpa identitas: %v", err)
	}
}

func TestAmbilProdukTidakAdaJadiGalatEntitas(t *testing.T) {
	l, _ := layananUji()
	_, err := l.AmbilProduk(context.Background(), pelakuUji, "UJI-TIDAK-ADA")
	if !errors.Is(err, services.ErrProdukTidakAda) {
		t.Errorf("mau ErrProdukTidakAda, dapat %v", err)
	}
	if services.Pesan(err) != "product not found: UJI-TIDAK-ADA" {
		t.Errorf("kalimat layar tanpa awalan lapisan: %q", services.Pesan(err))
	}
}

func TestGagalBacaDiteruskanTidakDitelan(t *testing.T) {
	l, g := layananUji()
	g.GagalBaca = tiruan.ErrTiruan
	if _, err := l.DaftarProduk(context.Background(), pelakuUji); !errors.Is(err, tiruan.ErrTiruan) {
		t.Errorf("galat gudang harus diteruskan: %v", err)
	}
}
