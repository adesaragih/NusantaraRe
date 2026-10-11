package services_test

// Uji pilihan dropdown kepala — nilai TERSIMPAN berpasangan label tampil.
//
// Sebab uji ini ada: dropdown `Accounting Mode` berbunyi
// "Underwriting Year (tidak ada di daftar referensi)" karena form mengisinya
// dengan LABEL, sementara pilihannya NILAI.

import (
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func pasangan(o []models.Opsi) [][2]string {
	out := [][2]string{}
	for _, x := range o {
		out = append(out, [2]string{x.Nilai, x.Label})
	}
	return out
}

func samakan(t *testing.T, nama string, got, mau [][2]string) {
	t.Helper()
	if len(got) != len(mau) {
		t.Fatalf("%s: %v, mau %v", nama, got, mau)
	}
	for i := range mau {
		if got[i] != mau[i] {
			t.Errorf("%s[%d] = %v, mau %v", nama, i, got[i], mau[i])
		}
	}
}

func TestOpsiKepalaNilaiTersimpanBerpasanganLabel(t *testing.T) {
	o := services.OpsiKepalaKontrak()
	samakan(t, "caraPembukuan", pasangan(o.CaraPembukuan),
		[][2]string{{"underwriting", "Underwriting Year"}, {"accounting", "Accounting Year"}})
	// ⛔ `risk` dan `nonreporting` TIDAK ditebak labelnya — berlabel nilainya.
	samakan(t, "caraPembukuanNonProp", pasangan(o.CaraPembukuanNonProp),
		[][2]string{{"loss", "Loss Occuring"}, {"risk", "risk"}})
	samakan(t, "bordereaux", pasangan(o.Bordereaux),
		[][2]string{{"reporting", "Reporting"}, {"nonreporting", "nonreporting"}})
}

// ⛔ Label datang dari penerjemah yang SAMA dengan medan kontrak — nol
// penerjemah kedua.
func TestOpsiKepalaMemakaiPenerjemahYangSama(t *testing.T) {
	for _, x := range services.OpsiKepalaKontrak().CaraPembukuan {
		if x.Label != services.CaraPembukuanTampil(x.Nilai) {
			t.Errorf("%s: label %q, penerjemah %q", x.Nilai, x.Label, services.CaraPembukuanTampil(x.Nilai))
		}
	}
}

func TestOpsiKepalaMenolakTanpaIdentitas(t *testing.T) {
	l := services.LayananDengan(nil)
	if _, err := l.OpsiKepala(inti.Pelaku{}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("mau ErrTanpaIdentitas, dapat %v", err)
	}
	if o, err := l.OpsiKepala(inti.Pelaku{AkunID: "AKUN-UJI"}); err != nil || len(o.CaraPembukuan) != 2 {
		t.Errorf("positif: %v %v", o, err)
	}
}
