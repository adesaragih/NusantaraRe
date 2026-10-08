package models_test

import (
	"strings"
	"testing"

	"nusantarare/modul/claimprop/backend/models"
)

// isiTab - jalur medan / grid dan label di dalam tab berjudul `judul` layar InputAcceptation.
func isiTab(t *testing.T, judul string) map[string]bool {
	t.Helper()
	ts := models.Evaluasi(models.HalamanBaru(), models.LayarAkseptasi(), false)
	var tab *models.Tata
	var cari func(xs []models.Tata)
	cari = func(xs []models.Tata) {
		for i := range xs {
			if xs[i].Letak == models.LetakTab {
				tab = &xs[i]
				return
			}
			cari(xs[i].Anak)
		}
	}
	cari(ts)
	if tab == nil {
		t.Fatal("layar akseptasi tanpa tab")
	}
	isi := map[string]bool{}
	var kumpul func(xs []models.Tata)
	kumpul = func(xs []models.Tata) {
		for _, x := range xs {
			if x.Jalur != "" {
				isi[x.Jalur] = true
			}
			if x.Label != "" && x.Jenis == models.JenisLabel {
				isi["label:"+x.Label] = true
			}
			kumpul(x.Anak)
			kumpul(x.Kaki)
		}
	}
	var judulTab []string
	for _, a := range tab.Anak {
		judulTab = append(judulTab, a.Label)
		if a.Label == judul {
			kumpul(a.Anak)
		}
	}
	if len(isi) == 0 {
		t.Fatalf("tab %q kosong atau tidak ada (tab: %s)", judul, strings.Join(judulTab, "|"))
	}
	return isi
}

// Work owner 08-10-2026: "pada bagian akseptasi, semua data estimasi yang muncul di tab akseptasi di hapus" - tab
// Acceptation tidak lagi mengulang data estimasi (Claim Estimation InputAcceptation_Adjs); datanya tetap di tab
// Estimation.
func TestTabAcceptationTanpaDataEstimasi(t *testing.T) {
	acc := isiTab(t, "Acceptation")
	for _, j := range []string{"label:Claim Estimation", models.DaftarClaimAmount, models.DaftarLossAlloc,
		models.DaftarEstimasi, models.DaftarTotalEst, models.DaftarSpreading, models.DaftarBreakQS,
		models.CD + "TotalGrossEstimateIDR", models.CD + "TotalEstimasiIDR"} {
		if acc[j] {
			t.Errorf("tab Acceptation masih menampilkan %s", j)
		}
	}
	// data akseptasi tetap: Acceptation List dan Spreading Adjustment Total (kakinya TotalEstimasi, seperti XML)
	if !acc[models.DaftarAdjustment] || !acc["label:Acceptance Information"] || !acc[models.DaftarSpreadAdj] ||
		!acc[models.DaftarSpreadAdjQS] {
		t.Errorf("tab Acceptation kehilangan data akseptasi: %v", acc)
	}
	est := isiTab(t, "Estimation")
	for _, j := range []string{models.DaftarClaimAmount, models.DaftarLossAlloc, models.DaftarEstimasi,
		models.DaftarSpreading, models.DaftarBreakQS} {
		if !est[j] {
			t.Errorf("tab Estimation kehilangan %s", j)
		}
	}
}
