package services_test

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/masterproductnamelife/backend/models"
	"nusantarare/modul/masterproductnamelife/backend/services"
)

func TestCariMasterHurufBesarDanRIRateMenunggu(t *testing.T) {
	l, g := layananUji()
	g.Master[models.MasterRIRisk] = []models.NilaiMaster{{ID: "1", Nama: "UJI RISK"}}
	d, err := l.CariMaster(context.Background(), pelakuUji, models.MasterRIRisk, "  risk ", 0)
	if err != nil || len(d) != 1 || g.CariTerakhir != "RISK" {
		t.Errorf("SearchPolicyHolder_act b236 huruf besar: %v %v %q", d, err, g.CariTerakhir)
	}
	if _, err := l.CariMaster(context.Background(), pelakuUji, models.MasterRIRate, "", 0); !errors.Is(err, services.ErrRIRateMenungguPersetujuan) {
		t.Errorf("R/I Rate: %v", err)
	}
	if _, err := l.CariMaster(context.Background(), inti.Pelaku{}, models.MasterCeding, "", 0); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
}
