package services_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/services"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

// gudangKasusCSV - kasus Perubahan Data `EDMLF-1` (QR) dengan dua peserta Old.
func gudangKasusCSV(edmType string) *tiruan.Gudang {
	g := tiruan.Baru()
	g.Polis["EDMLF-1"] = &tiruan.Polis{ID: "EDMLF-1", OldPolicyNo: "UJI-PL-1", EdmType: edmType, ProdKe: 2,
		Kepala: map[string]string{"TYPE": models.TypeQR}}
	for _, id := range []string{"UJI-D1", "UJI-D2"} {
		g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: id, PolisID: "EDMLF-1", PLNumber: "UJI-PL-1", EdmStatus: models.StatusOld,
			Nilai: map[string]string{"CERTIFICATE_NO": id, "PLAN": "UJI-PLAN", "POLICY_HOLDER": "UJI-PH", "CURRENCY": "IDR", "GROSS_PREMIUM": "100"}})
	}
	return g
}

const judulCSV = "CERTIFICATE_NO,PLAN,POLICY_HOLDER,CURRENCY,GROSS_PREMIUM,DOB\n"

func csvUji(baris ...string) *strings.Reader {
	return strings.NewReader(judulCSV + strings.Join(baris, "\n") + "\n")
}

func pesertaBaru(g *tiruan.Gudang) []*tiruan.Peserta {
	var hasil []*tiruan.Peserta
	for _, p := range g.Peserta {
		if p.EdmStatus == models.StatusNew {
			hasil = append(hasil, p)
		}
	}
	return hasil
}

func TestPeriksaCSVTanpaMenulis(t *testing.T) {
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	l, j := layananJejak(g)
	h, err := l.PeriksaCSV(context.Background(), pelakuUji, "EDMLF-1", csvUji(
		"UJI-N1,UJI-PLAN,UJI-PH,IDR,\"10,5\",02/01/1990",
		"UJI-N2,UJI-LAIN,UJI-PH,IDR,1.234.567,02/01/1990",
		"",
		"UJI-N3,UJI-PLAN,UJI-PH,IDR,7,1990-01-02",
	))
	if err != nil {
		t.Fatal(err)
	}
	if h.Total != 3 || h.Ditolak != 2 || len(pesertaBaru(g)) != 0 || len(j.catatan) != 0 {
		t.Fatalf("tinjauan %+v; baris New %d", h, len(pesertaBaru(g)))
	}
	kolom := map[string]int{}
	for _, p := range h.Pesan {
		kolom[p.Kolom] = p.Baris
		if p.Kolom == "PLAN" && p.Pesan != "Plan di CSV tidak sesuai, mohon di cek kembali" {
			t.Errorf("pesan PLAN tidak VERBATIM: %q", p.Pesan)
		}
	}
	if kolom["PLAN"] != 2 || kolom["GROSS_PREMIUM"] != 2 || kolom["DOB"] != 3 {
		t.Errorf("baris/kolom penolakan %v", h.Pesan)
	}
}

func TestTambahCSVMenyisipBaruLaluMengunci(t *testing.T) {
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	l, j := layananJejak(g)
	ctx := context.Background()
	h, err := l.TambahCSV(ctx, pelakuUji, "EDMLF-1", csvUji(
		"UJI-N1,UJI-PLAN,UJI-PH,IDR,\"10,5\",02/01/1990",
		"UJI-N2,UJI-PLAN,UJI-PH,USD,,",
	))
	if err != nil {
		t.Fatal(err)
	}
	baru := pesertaBaru(g)
	if h.Disimpan != 2 || len(baru) != 2 || baru[0].Nilai["GROSS_PREMIUM"] != "10.5" || baru[0].Nilai["DOB"] != "1990-01-02" ||
		baru[1].Nilai["GROSS_PREMIUM"] != "0" || baru[1].Nilai["DOB"] != "" || baru[0].PLNumber != "UJI-PL-1" {
		t.Fatalf("%+v; baris New %+v %+v", h, baru[0], baru[1])
	}
	if len(j.catatan) != 1 || j.catatan[0].Komentar != "Add CSV Data" {
		t.Errorf("jejak %+v", j.catatan)
	}
	k, _ := l.BacaKasus(ctx, pelakuUji, "EDMLF-1")
	if !k.CSVTerkunci {
		t.Error("Upload CSV tidak terkunci sesudah Add CSV Data (.EditInput1 = 1)")
	}
	// Unggahan kedua ditolak: tidak menumpuk.
	if _, err := l.TambahCSV(ctx, pelakuUji, "EDMLF-1", csvUji("UJI-N9,UJI-PLAN,UJI-PH,IDR,1,")); !errors.Is(err, services.ErrCSVTerkunci) {
		t.Fatalf("unggahan kedua: %v", err)
	}
	if len(pesertaBaru(g)) != 2 {
		t.Errorf("baris New %d sesudah unggahan kedua, mau tetap 2", len(pesertaBaru(g)))
	}
}

func TestTambahCSVSatuPenolakanMenolakSeluruhBerkas(t *testing.T) {
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	l, j := layananJejak(g)
	_, err := l.TambahCSV(context.Background(), pelakuUji, "EDMLF-1", csvUji(
		"UJI-N1,UJI-PLAN,UJI-PH,IDR,1,",
		"UJI-N2,UJI-PLAN,UJI-PH-LAIN,IDR,1,",
	))
	var gc services.GalatCSV
	if !errors.As(err, &gc) || gc.Ditolak != 1 || gc.Total != 2 || gc.Pesan[0].Kolom != "POLICY_HOLDER" || gc.Pesan[0].Baris != 2 {
		t.Fatalf("%v %+v", err, gc)
	}
	if len(pesertaBaru(g)) != 0 || len(j.catatan) != 0 {
		t.Error("baris sebelum pelanggaran ikut tersimpan (R24: Pega berhenti SESUDAH menambah)")
	}
}

// TestTambahCSVTanpaBatasBaris - AC 36: di atas 50.000 baris tetap diterima.
func TestTambahCSVTanpaBatasBaris(t *testing.T) {
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	l, _ := layananJejak(g)
	var b strings.Builder
	b.WriteString(judulCSV)
	const n = 50_001
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&b, "UJI-N%d,UJI-PLAN,UJI-PH,IDR,1,\n", i)
	}
	h, err := l.TambahCSV(context.Background(), pelakuUji, "EDMLF-1", strings.NewReader(b.String()))
	if err != nil || h.Disimpan != n || len(pesertaBaru(g)) != n {
		t.Fatalf("%d baris: %+v %v", n, h, err)
	}
}

func TestTambahCSVMenghitungUlangRekapSesudahSave(t *testing.T) {
	g := gudangKasusCSV(models.EdmTypePerubahanData)
	l, _ := layananJejak(g)
	ctx := context.Background()
	if _, err := l.Simpan(ctx, pelakuUji, "EDMLF-1", models.PilihanHapus{}); err != nil {
		t.Fatal(err)
	}
	h, err := l.TambahCSV(ctx, pelakuUji, "EDMLF-1", csvUji("UJI-N1,UJI-PLAN,UJI-PH,IDR,50,"))
	if err != nil {
		t.Fatal(err)
	}
	// R31: rekap Pega basi bila CSV ditambah sesudah Save; di sini ikut 100 + 100 + 50.
	if len(h.Rekap) != 1 || h.Rekap[0]["PREMIUM"] != "250" {
		t.Fatalf("rekap %v", h.Rekap)
	}
}

func TestCSVDitolak(t *testing.T) {
	ctx := context.Background()
	tertutup := gudangKasusCSV(models.EdmTypePerubahanData)
	tertutup.Polis["EDMLF-1"].Status = models.StatusKasusDitolak
	kosong := gudangKasusCSV(models.EdmTypePerubahanData)
	kosong.Peserta = nil
	gagal := gudangKasusCSV(models.EdmTypePerubahanData)
	gagal.GagalSisipKe = 1
	for _, k := range []struct {
		nama  string
		g     *tiruan.Gudang
		isi   string
		galat error
	}{
		{"Batal", gudangKasusCSV(models.EdmTypeBatal), judulCSV + "UJI-N1,UJI-PLAN,UJI-PH,IDR,1,\n", services.ErrCSVBukanPerubahanData},
		{"tertutup", tertutup, judulCSV + "UJI-N1,UJI-PLAN,UJI-PH,IDR,1,\n", services.ErrKasusTertutup},
		{"tanpa peserta lama", kosong, judulCSV + "UJI-N1,UJI-PLAN,UJI-PH,IDR,1,\n", services.ErrCSVTanpaAcuan},
		{"kosong", gudangKasusCSV("1"), judulCSV, services.ErrCSVKosong},
		{"tanpa judul", gudangKasusCSV("1"), "", services.ErrCSVKosong},
		{"tanpa POLICY_HOLDER", gudangKasusCSV("1"), "PLAN,CURRENCY\nUJI-PLAN,IDR\n", services.ErrMasukanTidakSah},
		{"kutip rusak", gudangKasusCSV("1"), judulCSV + "UJI-N1,\"UJI-PLAN,UJI-PH,IDR,1,\n", services.ErrCSVRusak},
	} {
		l, _ := layananJejak(k.g)
		if _, err := l.TambahCSV(ctx, pelakuUji, "EDMLF-1", strings.NewReader(k.isi)); !errors.Is(err, k.galat) {
			t.Errorf("%s: %v, mau %v", k.nama, err, k.galat)
		}
	}
	l, j := layananJejak(gagal)
	if _, err := l.TambahCSV(ctx, pelakuUji, "EDMLF-1", csvUji("UJI-N1,UJI-PLAN,UJI-PH,IDR,1,")); err == nil || len(j.catatan) != 0 {
		t.Errorf("galat sisip tidak diteruskan: %v", err)
	}
}
