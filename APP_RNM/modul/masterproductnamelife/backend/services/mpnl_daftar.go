package services

// Daftar bersarang produk (paket 6–7) - gerbang per baris.
//
//	PLAN LIST b31560           `ProteksiPlanListLife` (onChange `.Plan` b33163, setiap `Delete` baris):
//	                           2.1.1 b563 `·` `.Plan==""`  → "Plan tidak boleh kosong" (b442)
//	                           2.1.2 b713 `·` plan sama di baris lain → "Plan tidak boleh sama" (b421)
//	                           2.1.3 b874 `·` `.RIRATE==""` → "RI/RATE tidak boleh kosong" (b463)
//	UNDERWRITING LIMIT b42078  angka desimal, Min ≤ Max (R17); `Medical` teks bebas (R16)
//
// ⛔ Pesan Pega dipasang di medan barisnya; di sini diawali judul grid dan
// nomor baris (mulai 1) supaya layar dapat menunjuknya.

import (
	"context"
	"fmt"
	"strings"

	"nusantarare/modul/masterproductnamelife/backend/models"
)

// Pesan VERBATIM `ProteksiPlanListLife` 2 b345.
const (
	PesanPlanSama         = "Plan tidak boleh sama"      // b421
	PesanPlanKosong       = "Plan tidak boleh kosong"    // b442
	PesanRIRatePlanKosong = "RI/RATE tidak boleh kosong" // b463
)

// Judul grid VERBATIM.
const (
	judulPlan    = "PLAN LIST"          // b31560
	judulUWLimit = "UNDERWRITING LIMIT" // b42078
)

func (pk *periksa) baris(judul string, i int, format string, a ...any) {
	pk.tolak("%s row %d: %s", judul, i+1, fmt.Sprintf(format, a...))
}

// periksaPlan - gerbang murni `ProteksiPlanListLife`. Duplikat hanya untuk nilai
// berisi (dua baris kosong sudah ditolak "kosong").
func periksaPlan(pk *periksa, daftar []models.BarisPlan) {
	pertama := map[string]int{}
	for i, b := range daftar {
		plan := strings.TrimSpace(b.Plan)
		if plan == "" {
			pk.baris(judulPlan, i, "%s", PesanPlanKosong)
		} else if _, ada := pertama[plan]; ada {
			pk.baris(judulPlan, i, "%s", PesanPlanSama)
		} else {
			pertama[plan] = i
		}
		if strings.TrimSpace(b.RIRate) == "" {
			pk.baris(judulPlan, i, "%s", PesanRIRatePlanKosong)
		}
	}
}

// periksaUWLimit - angka dan rentang baris `UNDERWRITING LIMIT`.
func periksaUWLimit(pk *periksa, daftar []models.BarisUWLimit) {
	for i := range daftar {
		b := &daftar[i]
		var sub periksa
		minSI, maxSI := sub.desimal("Min Insured", &b.MinInsured), sub.desimal("Max Insured", &b.MaxInsured)
		minUsia, maxUsia := sub.desimal("Min Age", &b.MinAge), sub.desimal("Max Age", &b.MaxAge)
		sub.tidakLebihBesar("Min Insured", minSI, "Max Insured", maxSI)
		sub.tidakLebihBesar("Min Age", minUsia, "Max Age", maxUsia)
		for _, s := range sub.pesan {
			pk.baris(judulUWLimit, i, "%s", s)
		}
	}
}

// periksaPilihanPlan - tiket 06: R/I Rate baris plan dari master R/I Rate (bukan
// R/I Risk); plan dari master `PRODUCT_TYPE_LIFE`. Pasangan yang sudah tersimpan
// diterima apa adanya; R/I Rate BARU tidak dapat diverifikasi sampai OQ-MPNL-03.
func (l *Layanan) periksaPilihanPlan(ctx context.Context, pk *periksa, daftar []models.BarisPlan, lama []models.BarisPlan) error {
	rateLama, planLama := map[[2]string]bool{}, map[[2]string]bool{}
	for _, b := range lama {
		rateLama[[2]string{b.RIRateID, b.RIRate}] = true
		planLama[[2]string{b.PlanID, b.Plan}] = true
	}
	for i := range daftar {
		b := &daftar[i]
		b.PlanID = strings.TrimSpace(b.PlanID)
		if !planLama[[2]string{b.PlanID, b.Plan}] && b.PlanID != "" {
			v, ada, err := l.gudang.AmbilPlan(ctx, b.PlanID)
			if err != nil {
				return err
			}
			if !ada {
				pk.baris(judulPlan, i, "Plan Name %q is not in the master list", b.PlanID)
			} else {
				b.Plan, b.Name, b.Benefit = v.CoverName, v.Business, v.Benefit
			}
		} else if b.PlanID == "" && strings.TrimSpace(b.Plan) != "" && !planLama[[2]string{b.PlanID, b.Plan}] {
			pk.baris(judulPlan, i, "Plan Name %q must be chosen from the master list", b.Plan)
		}
		if strings.TrimSpace(b.RIRate) != "" && !rateLama[[2]string{b.RIRateID, b.RIRate}] {
			pk.baris(judulPlan, i, "R/I Rate %q cannot be chosen until reading the R/I Rate master is approved (OQ-MPNL-03)",
				b.RIRate)
		}
	}
	return nil
}
