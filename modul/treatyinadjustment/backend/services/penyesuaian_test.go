package services_test

// Uji seam layar Adjustment: identitas, pengenal kosong, penyesuaian yang
// tidak ada, galat yang diteruskan apa adanya, dan nilai yang TIDAK
// diterjemahkan di services.

import (
	"context"
	"errors"
	"fmt"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/repository"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

func (g *gudangTiruan) DaftarPenyesuaianWarisan(context.Context) ([]models.BarisPenyesuaian, error) {
	g.disentuh++
	return g.daftarPenyesuaian, g.galatPenyesuaian
}

func (g *gudangTiruan) BacaPenyesuaianPendaratan(_ context.Context, id string) (models.Penyesuaian, error) {
	g.disentuh++
	if g.galatPenyesuaian != nil {
		return models.Penyesuaian{}, g.galatPenyesuaian
	}
	p, ada := g.penyesuaian[id]
	if !ada {
		return models.Penyesuaian{}, fmt.Errorf("%w: %s", repository.ErrPenyesuaianTidakAda, id)
	}
	return p, nil
}

func TestPenyesuaianMenolakTanpaIdentitasSebelumGudang(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	if _, err := l.DaftarPenyesuaianWarisan(context.Background(), inti.Pelaku{}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("DaftarPenyesuaianWarisan: mau ErrTanpaIdentitas, dapat %v", err)
	}
	if _, err := l.PenyesuaianWarisan(context.Background(), inti.Pelaku{}, "1000080/R02"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("PenyesuaianWarisan: mau ErrTanpaIdentitas, dapat %v", err)
	}
	if g.disentuh != 0 {
		t.Errorf("gudang tersentuh %d kali walau identitas tidak ada", g.disentuh)
	}
}

func TestPenyesuaianPengenalKosongDitolakSebelumGudang(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)
	for _, id := range []string{"", "   ", "\t"} {
		if _, err := l.PenyesuaianWarisan(context.Background(), pelakuAda, id); !errors.Is(err, services.ErrIDTidakSah) {
			t.Errorf("id %q: mau ErrIDTidakSah, dapat %v", id, err)
		}
	}
	if g.disentuh != 0 {
		t.Errorf("gudang tersentuh %d kali untuk pengenal kosong", g.disentuh)
	}
}

// ⛔ Penyesuaian yang tidak ada berbunyi sebagai TIDAK ADA - bukan galat
// umum, bukan dokumen kosong yang terbuka diam-diam.
func TestPenyesuaianTidakAdaDikenaliSebagaiTidakAda(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{penyesuaian: map[string]models.Penyesuaian{}})
	_, err := l.PenyesuaianWarisan(context.Background(), pelakuAda, "9999999/R01")
	if !services.PenyesuaianTidakAda(err) {
		t.Fatalf("mau galat tidak-ada, dapat %v", err)
	}
}

// ⛔ Galat gudang DITERUSKAN apa adanya - tidak ditelan, tidak diganti.
func TestPenyesuaianGalatGudangDiteruskan(t *testing.T) {
	sebab := errors.New("ORA-03113: end-of-file on communication channel")
	l := services.LayananDengan(&gudangTiruan{galatPenyesuaian: sebab})
	if _, err := l.DaftarPenyesuaianWarisan(context.Background(), pelakuAda); !errors.Is(err, sebab) {
		t.Errorf("DaftarPenyesuaianWarisan: galat tidak diteruskan, dapat %v", err)
	}
	if _, err := l.PenyesuaianWarisan(context.Background(), pelakuAda, "1000080/R02"); !errors.Is(err, sebab) {
		t.Errorf("PenyesuaianWarisan: galat tidak diteruskan, dapat %v", err)
	}
	if services.PenyesuaianTidakAda(sebab) {
		t.Error("galat koneksi terbaca sebagai 'tidak ada'")
	}
}

// ⭐ Positif: kedua sisi tiba UTUH, dan services TIDAK menerjemahkan apa
// pun - tanggal `20180101` tetap `20180101`, angka tetap digit aslinya.
// Pemformatnya satu, di layar, untuk kedua panel.
func TestPenyesuaianPositifKeduaSisiApaAdanya(t *testing.T) {
	asli := models.Penyesuaian{
		ID: "1000080/R02", IDAsal: "1000080/R01",
		Baru: models.SisiPenyesuaian{
			Medan: map[string]string{"Commencement": "20180101", "TreatyContractName": "BARU"},
			Larik: map[string][]map[string]string{"CurrencyList": {{"Currency": "USD", "Conversion": "14250.5"}}},
		},
		Lama: models.SisiPenyesuaian{
			Medan: map[string]string{"Commencement": "20170101", "TreatyContractName": "LAMA"},
			Larik: map[string][]map[string]string{},
		},
	}
	l := services.LayananDengan(&gudangTiruan{penyesuaian: map[string]models.Penyesuaian{asli.ID: asli}})
	p, err := l.PenyesuaianWarisan(context.Background(), pelakuAda, asli.ID)
	if err != nil {
		t.Fatal(err)
	}
	if p.Baru.Medan["Commencement"] != "20180101" || p.Lama.Medan["Commencement"] != "20170101" {
		t.Errorf("tanggal diterjemahkan di services: baru %q lama %q", p.Baru.Medan["Commencement"], p.Lama.Medan["Commencement"])
	}
	if p.Baru.Larik["CurrencyList"][0]["Conversion"] != "14250.5" {
		t.Errorf("angka berubah: %q", p.Baru.Larik["CurrencyList"][0]["Conversion"])
	}
	if p.Lama.Medan["TreatyContractName"] != "LAMA" || p.Baru.Medan["TreatyContractName"] != "BARU" {
		t.Error("sisi Old dan New tertukar")
	}
}

func TestDaftarPenyesuaianKosongBukanGalat(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{daftarPenyesuaian: []models.BarisPenyesuaian{}})
	d, err := l.DaftarPenyesuaianWarisan(context.Background(), pelakuAda)
	if err != nil || d == nil || len(d) != 0 {
		t.Fatalf("mau daftar kosong tanpa galat, dapat %v %v", d, err)
	}
}
