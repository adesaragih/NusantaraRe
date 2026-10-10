package services_test

// Laporan pemakai 8 Oktober 2026 (tangkapan tab Limits Non-Prop): kolom
// "Reinstatement Premium Amount IDR (MDP x % Add Premium)" kosong.
//
// Ekspor `Treaty In`: `AdditionalAmount1/2` HANYA ditulis `SetReinstatementPct`
// (medan Reinstatement berubah) dan DT `ReCalculateReinstatement` (% Additional
// Premium berubah). `DetailCalculation(mdp)` TIDAK menyentuh grid itu — grid
// yang dibangun SEBELUM MDP ada tetap kosong, persis Pega.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

func layerLayar(mdp []services.NilaiMataUang) services.LayerNP {
	return services.LayerNP{
		Limit: "16000000000", Currency: "IDR", AdjRate: "4.183", MDPPct: "90", ReinstatementValue: "3",
		PremiumEarnedList: []services.NilaiMataUang{{Currency: "IDR", Value: "2300650000"}},
		MDPList:           mdp,
	}
}

func TestReinstatementPremiumAmountDariMDP(t *testing.T) {
	mdp := []services.NilaiMataUang{{Currency: "IDR", Value: "2070585000"}}

	// Reinstatement diisi SESUDAH MDP ada → ketiga baris berisi MDP × 100%.
	h := services.HitungLimitNP(services.MasukanLimitNP{Aksi: "reinstatement", Layers: []services.LayerNP{layerLayar(mdp)}})
	for _, b := range h.Layers[0].ReinstatementList {
		if b.AdditionalAmount1 != "2070585000" || b.AdditionalAmount2 != "0" {
			t.Errorf("baris %s: %q/%q", b.ReinstatementValue, b.AdditionalAmount1, b.AdditionalAmount2)
		}
	}

	// Reinstatement diisi SEBELUM MDP ada → kosong (SetReinstatementPct
	// langkah 3.2/3.3 dilewati).
	awal := services.HitungLimitNP(services.MasukanLimitNP{Aksi: "reinstatement", Layers: []services.LayerNP{layerLayar(nil)}})
	l := awal.Layers[0]
	if got := l.ReinstatementList[0].AdditionalAmount1; got != "" {
		t.Errorf("tanpa MDP: %q", got)
	}
	// ⭐ Simpangan WAKTU (keputusan pemakai 8 Oktober 2026): MDP dihitung
	// ulang (`DetailCalculation(mdp)`) → ReCalculateReinstatement — rumus
	// yang SAMA — menyegarkan setiap baris: 100% × 2.070.585.000.
	sesudahMDP := services.HitungLimitNP(services.MasukanLimitNP{Aksi: "mdp", Layers: []services.LayerNP{l}})
	for _, b := range sesudahMDP.Layers[0].ReinstatementList {
		if b.AdditionalAmount1 != "2070585000" || b.AdditionalAmount2 != "0" {
			t.Errorf("sesudah MDP, baris %s: %q/%q", b.ReinstatementValue, b.AdditionalAmount1, b.AdditionalAmount2)
		}
	}
	l.MDPList = mdp

	// % Additional Premium diubah → ReCalculateReinstatement mengisinya.
	ubah := services.HitungLimitNP(services.MasukanLimitNP{Aksi: "reinst-tambahan", Layers: []services.LayerNP{l}, Baris: 1})
	if got := ubah.Layers[0].ReinstatementList[1].AdditionalAmount1; got != "2070585000" {
		t.Errorf("ReCalculateReinstatement baris 2 = %q", got)
	}
}
