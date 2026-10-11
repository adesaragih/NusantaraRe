package services_test

// Lampiran "Reas" kasus (keputusan work owner 08-10-2026): Upload / Delete "semua bisa asal belum resolve"; lampiran
// Pega lama berkunci pzInsKey `ASM-FW-GISFW-WORK <pyID>` ikut dibaca. Fixture UJI-.

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
	"nusantarare/modul/nbtreatyin/backend/tiruan"
)

func TestKasusLampiranBelumResolve(t *testing.T) {
	g := tiruan.Baru()
	g.Kasus["NB-1"] = models.Kasus{ID: "NB-1", StatusWork: models.AssignmentAdmin}
	g.Kasus["NB-2"] = models.Kasus{ID: "NB-2", StatusWork: models.StatusSelesai}
	g.Kasus["NB-3"] = models.Kasus{ID: "NB-3", StatusWork: models.StatusDitolak}
	l := services.Baru(g, nil)
	ctx := context.Background()
	p := inti.Pelaku{AkunID: "UJI-SIAPA-SAJA"}
	k, err := l.KasusLampiran(ctx, p, "NB-1")
	if err != nil || !k.BolehUbah || k.Tulis != "NB-1" || len(k.Baca) != 2 || k.Baca[1] != "ASM-FW-GISFW-WORK NB-1" {
		t.Errorf("kasus berjalan: %+v %v", k, err)
	}
	for _, id := range []string{"NB-2", "NB-3"} {
		if k, err := l.KasusLampiran(ctx, p, id); err != nil || k.BolehUbah {
			t.Errorf("%s Resolve: %+v %v", id, k, err)
		}
	}
	if _, err := l.KasusLampiran(ctx, p, "NB-404"); !errors.Is(err, services.ErrKasusTidakAda) {
		t.Errorf("kasus tidak ada: %v", err)
	}
	if _, err := l.KasusLampiran(ctx, inti.Pelaku{}, "NB-1"); err == nil {
		t.Error("tanpa identitas diterima")
	}
}
