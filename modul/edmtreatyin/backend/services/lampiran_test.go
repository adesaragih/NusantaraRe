package services_test

// Lampiran "Reas" kasus (keputusan work owner 08-10-2026): Upload / Delete "semua bisa asal belum resolve"; lampiran
// Pega lama berkunci pzInsKey `ASM-FW-GISFW-WORK <pyID>` ikut dibaca. Fixture UJI-.

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/services"
	"nusantarare/modul/edmtreatyin/backend/tiruan"
)

func TestKasusLampiranBelumResolve(t *testing.T) {
	g := tiruan.Baru()
	g.Kasus["EDMT-1"] = models.Kasus{ID: "EDMT-1", StatusWork: models.AssignmentAdmin}
	g.Kasus["EDMT-2"] = models.Kasus{ID: "EDMT-2", StatusWork: models.StatusSelesai}
	g.Kasus["EDMT-3"] = models.Kasus{ID: "EDMT-3", StatusWork: models.StatusDitolak}
	for _, id := range []string{"EDMT-1", "EDMT-2", "EDMT-3"} {
		g.Generasi[id] = &tiruan.Generasi{} // tiruan EDM: setiap kasus punya generasi polis
	}
	l := services.Baru(g, nil)
	ctx := context.Background()
	p := inti.Pelaku{AkunID: "UJI-SIAPA-SAJA"}
	k, err := l.KasusLampiran(ctx, p, "EDMT-1")
	if err != nil || !k.BolehUbah || k.Tulis != "EDMT-1" || len(k.Baca) != 2 || k.Baca[1] != "ASM-FW-GISFW-WORK EDMT-1" {
		t.Errorf("kasus berjalan: %+v %v", k, err)
	}
	for _, id := range []string{"EDMT-2", "EDMT-3"} {
		if k, err := l.KasusLampiran(ctx, p, id); err != nil || k.BolehUbah {
			t.Errorf("%s Resolve: %+v %v", id, k, err)
		}
	}
	if _, err := l.KasusLampiran(ctx, p, "EDMT-404"); !errors.Is(err, services.ErrKasusTidakAda) {
		t.Errorf("kasus tidak ada: %v", err)
	}
	if _, err := l.KasusLampiran(ctx, inti.Pelaku{}, "EDMT-1"); err == nil {
		t.Error("tanpa identitas diterima")
	}
}
