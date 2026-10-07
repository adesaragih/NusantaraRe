package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

func TestCariMasterHurufBesarDanRIRateDibaca(t *testing.T) {
	l, g := layananUji()
	g.Master[models.MasterRIRisk] = []models.NilaiMaster{{ID: "1", Nama: "UJI RISK"}}
	d, err := l.CariMaster(context.Background(), pelakuUji, models.MasterRIRisk, "  risk ", 0)
	if err != nil || len(d) != 1 || g.CariTerakhir != "RISK" {
		t.Errorf("SearchPolicyHolder_act b236 huruf besar: %v %v %q", d, err, g.CariTerakhir)
	}
	// K1 01-10-2026 (OQ-MPNL-03): pemilih R/I Rate membaca tabel `M_RATE_LIFE_SUMMARY`, kata cari dihurufbesarkan.
	g.Master[models.MasterRIRate] = []models.NilaiMaster{{ID: "R1", Nama: "UJI RATE"}}
	if d, err := l.CariMaster(context.Background(), pelakuUji, models.MasterRIRate, " rate", 0); err != nil || len(d) != 1 ||
		g.CariTerakhir != "RATE" {
		t.Errorf("R/I Rate: %v %v %q", d, err, g.CariTerakhir)
	}
	if _, err := l.CariMaster(context.Background(), inti.Pelaku{}, models.MasterCeding, "", 0); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
}

// View Rate (K1 01-10-2026): view `RATE_LIFE` satu RIRATEID, urut `ID DESC, RATE ASC`, RATE teks apa
// adanya, 500 baris + terpotong; tanpa RIRATEID ditolak (nol pembacaan tanpa kunci).
func TestDaftarRateDariViewRate(t *testing.T) {
	l, g := layananUji()
	g.Rate["R1"] = []models.BarisRate{{ID: "UJI-A", Rate: "0,5"}, {ID: "UJI-B", Rate: "1.25"}, {ID: "UJI-B", Rate: "0,75"}}
	r, err := l.DaftarRate(context.Background(), pelakuUji, "R1")
	if err != nil || r.Total != 3 || r.Terpotong || r.Daftar[0].Rate != "0,75" || r.Daftar[2].ID != "UJI-A" {
		t.Errorf("DaftarRate: %v %+v", err, r)
	}
	for i := 0; i < 501; i++ {
		g.Rate["R2"] = append(g.Rate["R2"], models.BarisRate{ID: fmt.Sprintf("UJI-%04d", i)})
	}
	if r, err := l.DaftarRate(context.Background(), pelakuUji, "R2"); err != nil || r.Total != 500 || !r.Terpotong {
		t.Errorf("501 baris: %v total %d terpotong %v", err, r.Total, r.Terpotong)
	}
	if r, err := l.DaftarRate(context.Background(), pelakuUji, "R-KOSONG"); err != nil || r.Daftar == nil || r.Total != 0 {
		t.Errorf("tanpa baris = daftar kosong: %v %+v", err, r)
	}
	if _, err := l.DaftarRate(context.Background(), pelakuUji, " "); !errors.Is(err, services.ErrMasukanTidakSah) {
		t.Errorf("tanpa RIRATEID: %v", err)
	}
	if _, err := l.DaftarRate(context.Background(), inti.Pelaku{}, "R1"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
}
