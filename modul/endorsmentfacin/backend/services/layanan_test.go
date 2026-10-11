package services

// Uji Layanan.Periksa di atas Gudang tiruan: validasi isian dan query yang dijalankan / dilompati.

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/repository"
	"nusantarare/modul/endorsmentfacin/backend/tiruan"
)

func gudangBersih() *tiruan.Gudang {
	return &tiruan.Gudang{JenisBisnisNilai: "Fire",
		Periode: &repository.PeriodeTeks{Mulai: "20250101000000", Akhir: "20260101000000"}}
}

func TestUraiMasukanPeriksaMenolak(t *testing.T) {
	for _, m := range []models.MasukanPeriksa{
		{PolicyNo: ""},
		{PolicyNo: "   "},
		{PolicyNo: polisUji, EdmType: "3"}, // K-029: 3 usang
		{PolicyNo: polisUji, EdmType: "x"},
		{PolicyNo: polisUji, EdmType: "1.0"},
		{PolicyNo: polisUji, EndorsementDate: "2025-01-01"},
		{PolicyNo: polisUji, EndorsementDate: "32-01-2025"},
		{PolicyNo: polisUji, EndorsementDate: "1-1-2025"},
	} {
		if _, err := UraiMasukanPeriksa(m); !errors.Is(err, ErrMasukanTidakSah) {
			t.Errorf("%+v: galat = %v", m, err)
		}
	}
}

func TestUraiMasukanPeriksaMenerima(t *testing.T) {
	for _, m := range []models.MasukanPeriksa{
		{PolicyNo: polisUji},
		{PolicyNo: polisUji, EdmType: "1", EndorsementDate: "29-02-2024"},
		{PolicyNo: polisUji, EdmType: "2", EndorsementDate: " 01-01-2025 "},
		{PolicyNo: polisUji, EdmType: "4", EdmTypeNew: "4"},
	} {
		d, err := UraiMasukanPeriksa(m)
		if err != nil {
			t.Fatalf("%+v: %v", m, err)
		}
		if (m.EndorsementDate == "") != (d.TanggalEndorsement == nil) {
			t.Errorf("%+v: tanggal = %v", m, d.TanggalEndorsement)
		}
	}
	// Nomor polis dipakai apa adanya (teks bebas).
	d, _ := UraiMasukanPeriksa(models.MasukanPeriksa{PolicyNo: " " + polisUji})
	if d.PolicyNo != " "+polisUji {
		t.Fatalf("polis diubah: %q", d.PolicyNo)
	}
	// Tanggal = tanggal kalender Jakarta.
	d, _ = UraiMasukanPeriksa(models.MasukanPeriksa{PolicyNo: polisUji, EndorsementDate: "01-01-2025"})
	if y, b, h := d.TanggalEndorsement.In(jakarta).Date(); y != 2025 || b != 1 || h != 1 {
		t.Fatalf("tanggal = %v", d.TanggalEndorsement)
	}
}

func TestPeriksaQueryYangDijalankan(t *testing.T) {
	for _, k := range []struct {
		nama  string
		m     models.MasukanPeriksa
		ada   []string
		absen []string
	}{
		{"biasa", models.MasukanPeriksa{PolicyNo: polisUji, EdmType: "4", EdmTypeNew: "1"},
			[]string{"StatusEDM", "Klaim", "AdaRenewal", "JenisBisnis"},
			[]string{"AdaPembayaran", "AdaFacout", "Pengecualian", "PeriodeProduksi"}},
		{"batal + retro + tanggal", models.MasukanPeriksa{PolicyNo: polisUji, EdmType: "1",
			EndorsementInternalRetro: "1", EndorsementDate: "01-06-2025"},
			[]string{"StatusEDM", "Klaim", "AdaPembayaran", "AdaRenewal", "AdaFacout", "JenisBisnis", "Pengecualian", "PeriodeProduksi"},
			nil},
		{"RI slip", models.MasukanPeriksa{PolicyNo: polisUji, EdmType: "4", EdmTypeNew: "4", EndorsementInternalRetro: "1"},
			[]string{"JenisBisnis"},
			[]string{"StatusEDM", "Klaim", "AdaPembayaran", "AdaRenewal", "AdaFacout"}},
		{"EdmType kosong", models.MasukanPeriksa{PolicyNo: polisUji},
			[]string{"StatusEDM", "Klaim", "AdaRenewal", "JenisBisnis"},
			[]string{"AdaPembayaran", "AdaFacout"}},
	} {
		g := gudangBersih()
		h, err := BaruLayanan(g).Periksa(context.Background(), k.m)
		if err != nil {
			t.Fatalf("%s: %v", k.nama, err)
		}
		for _, n := range k.ada {
			if !g.Dipanggil(n) {
				t.Errorf("%s: %s tidak dipanggil (%v)", k.nama, n, g.Panggilan)
			}
		}
		for _, n := range k.absen {
			if g.Dipanggil(n) {
				t.Errorf("%s: %s dipanggil (%v)", k.nama, n, g.Panggilan)
			}
		}
		if len(h.Gerbang) != 6 {
			t.Errorf("%s: %d gerbang", k.nama, len(h.Gerbang))
		}
	}
}

func TestPeriksaKlepDariGudang(t *testing.T) {
	klep := []string{"AdaProteksi:1", "AdaProteksi:2", "AdaProteksi:3", "AdaProteksi:4"}
	// Bukan RI slip: keempat klep dibaca, urut Param.Type, walau tidak satu gerbang pun menolak (seperti Pega).
	g := gudangBersih()
	if _, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji}); err != nil {
		t.Fatal(err)
	}
	var dibaca []string
	for _, p := range g.Panggilan {
		if len(p) > len("AdaProteksi") && p[:len("AdaProteksi")] == "AdaProteksi" {
			dibaca = append(dibaca, p)
		}
	}
	if !samaIrisan(dibaca, klep) {
		t.Fatalf("klep dibaca %v", dibaca)
	}
	// RI slip: nol klep.
	g = gudangBersih()
	if _, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji, EdmType: "4", EdmTypeNew: "4"}); err != nil {
		t.Fatal(err)
	}
	for _, k := range klep {
		if g.Dipanggil(k) {
			t.Errorf("RI slip membaca %s", k)
		}
	}
	// Klep 4 terisi membuka gerbang 5.
	g = gudangBersih()
	g.Renewal, g.Proteksi = true, map[int]bool{4: true}
	h, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji})
	if err != nil {
		t.Fatal(err)
	}
	if h.Gerbang[4].Hasil != models.GerbangDiklep || h.Ditolak {
		t.Fatalf("gerbang 5 = %+v ditolak %v", h.Gerbang[4], h.Ditolak)
	}
	// Tabel tidak ada: tolak penuh + catatan.
	g.GalatProteksi = fmt.Errorf("%w: ORA-00942", repository.ErrBelumAda)
	h, err = BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji})
	if err != nil {
		t.Fatal(err)
	}
	if h.Gerbang[4].Hasil != T || h.Gerbang[4].Catatan != catatanKlepTakTerbaca {
		t.Fatalf("tabel tidak ada: gerbang 5 = %+v", h.Gerbang[4])
	}
}

func TestPeriksaArasapasTakTerbaca(t *testing.T) {
	g := gudangBersih()
	g.GalatPembayaran = fmt.Errorf("%w: ORA-00942", repository.ErrSkemaLainTakTerbaca)
	h, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji, EdmType: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if h.Gerbang[3].Hasil != BD || h.Lengkap {
		t.Fatalf("gerbang4 = %+v lengkap %v", h.Gerbang[3], h.Lengkap)
	}
	// Galat lain tetap galat.
	g.GalatPembayaran = errors.New("ORA-03113")
	if _, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji, EdmType: "1"}); err == nil {
		t.Fatal("galat pembayaran lain ditelan")
	}
}

func TestPeriksaPengecualianBelumAda(t *testing.T) {
	g := gudangBersih()
	g.GalatPengecualian = fmt.Errorf("%w: ORA-00942", repository.ErrBelumAda)
	h, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji, EndorsementDate: "01-06-2025"})
	if err != nil {
		t.Fatal(err)
	}
	if h.Tanggal == nil || h.Tanggal.Hasil != L || h.Tanggal.Catatan != catatanPengecualianBelumAda {
		t.Fatalf("tanggal = %+v", h.Tanggal)
	}
	g.GalatPengecualian = errors.New("ORA-03113")
	if _, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji, EndorsementDate: "01-06-2025"}); err == nil {
		t.Fatal("galat pengecualian lain ditelan")
	}
}

func TestPeriksaPengecualianDariTabel(t *testing.T) {
	g := gudangBersih()
	g.PengecualianPer = map[string][]string{JenisPengecualianTanggal: {polisUji}, "BATAL": {"UJI.99.9999"}}
	h, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji, EndorsementDate: "01-01-2030"})
	if err != nil {
		t.Fatal(err)
	}
	if h.Tanggal.Hasil != D {
		t.Fatalf("tanggal = %+v", h.Tanggal)
	}
}

func TestPeriksaGalatQueryDiteruskan(t *testing.T) {
	boom := errors.New("ORA-03113")
	for _, set := range []func(*tiruan.Gudang){
		func(g *tiruan.Gudang) { g.GalatStatusEDM = boom },
		func(g *tiruan.Gudang) { g.GalatKlaim = boom },
		func(g *tiruan.Gudang) { g.GalatRenewal = boom },
		func(g *tiruan.Gudang) { g.GalatFacout = boom },
		func(g *tiruan.Gudang) { g.GalatJenisBisnis = boom },
		func(g *tiruan.Gudang) { g.GalatPeriode = boom },
		func(g *tiruan.Gudang) { g.GalatProteksi = boom },
	} {
		g := gudangBersih()
		set(g)
		_, err := BaruLayanan(g).Periksa(context.Background(), models.MasukanPeriksa{PolicyNo: polisUji,
			EndorsementInternalRetro: "1", EndorsementDate: "01-06-2025"})
		if !errors.Is(err, boom) {
			t.Errorf("galat = %v (panggilan %v)", err, g.Panggilan)
		}
	}
}
