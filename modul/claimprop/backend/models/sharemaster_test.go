package models_test

// RNM Share master per treaty group klaim (keputusan work owner 10-10-2026: RNM_SHARE view detail treaty group klaim,
// pengganti `TreatyInMaster.RNMShareP` JSON).

import (
	"testing"

	"nusantarare/modul/claimprop/backend/models"
)

func masterShare(detail ...models.DetailLimit) models.MasterTreaty {
	return models.MasterTreaty{ID: "UJI-M", Limits: []models.LimitMaster{{TreatyType: "QUOTA SHARE", Detail: detail}}}
}

func TestShareMasterPerTreatyGroupKlaim(t *testing.T) {
	dua := masterShare(models.DetailLimit{TreatyGroupID: "UJI-TG1", RNMShare: "2.5"}, models.DetailLimit{TreatyGroupID: "UJI-TG2", RNMShare: "40"})
	satu := masterShare(models.DetailLimit{TreatyGroupID: "UJI-TG1", RNMShare: "4"}, models.DetailLimit{TreatyGroupID: "UJI-TG2", RNMShare: "4"})
	for _, tt := range []struct {
		nama  string
		m     models.MasterTreaty
		grup  string
		harap string
	}{
		{"detail treaty group klaim", dua, "UJI-TG2", "40"},
		{"tanpa yang cocok, nilai berbeda -> kosong (pilih dari dropdown)", dua, "UJI-TG9", ""},
		{"tanpa yang cocok, satu nilai", satu, "UJI-TG9", "4"},
		{"treaty group kosong, satu nilai", satu, "", "4"},
		{"tanpa detail", models.MasterTreaty{}, "UJI-TG1", ""},
	} {
		if got := models.ShareMaster(tt.m, tt.grup); got != tt.harap {
			t.Errorf("%s: %q, harap %q", tt.nama, got, tt.harap)
		}
	}
}

func TestTerapkanMasterRNMShareDariTreatyGroupKlaim(t *testing.T) {
	h := models.HalamanBaru()
	h.Setel(models.CD+"TreatyGroupID", "UJI-TG2")
	m := masterShare(models.DetailLimit{TreatyGroupID: "UJI-TG1", RNMShare: "2.5"}, models.DetailLimit{TreatyGroupID: "UJI-TG2", RNMShare: "40"})
	models.TerapkanMaster(h, m, true)
	if got := h.Ambil(models.TM + "RNMShareP"); got != "40" {
		t.Fatalf("TreatyInMaster.RNMShareP = %q, harap 40 (detail UJI-TG2)", got)
	}
	// berkas dibuka ulang (gantiShare=false): pilihan tersimpan dipertahankan
	h.Setel(models.TM+"RNMShareP", "2.5")
	models.TerapkanMaster(h, m, false)
	if got := h.Ambil(models.TM + "RNMShareP"); got != "2.5" {
		t.Fatalf("RNMShareP tersimpan = %q, harap 2.5 dipertahankan", got)
	}
}
