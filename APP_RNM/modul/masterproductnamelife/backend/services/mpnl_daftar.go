package services

// Daftar bersarang produk (paket 6–7) - gerbang per baris.
//
//	PLAN LIST b31557           `ProteksiPlanListLife` (onChange `.Plan` b33163, setiap `Delete` baris):
//	                           2.1.1 b565 `·` `.Plan==""`  → "Plan tidak boleh kosong" (b442)
//	                           2.1.2 b715 `·` plan sama di baris lain → "Plan tidak boleh sama" (b421)
//	                           2.1.3 b876 `·` `.RIRATE==""` → "RI/RATE tidak boleh kosong" (b463)
//	UNDERWRITING LIMIT b42075  angka desimal, Min ≤ Max (R17); `Medical` teks bebas (R16)
//
// ⛔ Pesan Pega dipasang di medan barisnya; di sini diawali judul grid dan
// nomor baris (mulai 1) supaya layar dapat menunjuknya.

import (
	"context"
	"fmt"
	"strings"
	"time"

	"nusantarare/modul/masterproductnamelife/backend/models"
)

// Pesan VERBATIM `ProteksiPlanListLife` 2 b347.
const (
	PesanPlanSama         = "Plan tidak boleh sama"      // b421
	PesanPlanKosong       = "Plan tidak boleh kosong"    // b442
	PesanRIRatePlanKosong = "RI/RATE tidak boleh kosong" // b463
)

// Judul grid VERBATIM.
const (
	judulPlan    = "PLAN LIST"          // b31557
	judulUWLimit = "UNDERWRITING LIMIT" // b42075
)

func (pk *periksa) baris(judul string, i int, format string, a ...any) {
	pk.tolak("%s row %d: %s", judul, i+1, fmt.Sprintf(format, a...))
}

// periksaPlan - gerbang murni `ProteksiPlanListLife`. Loop luar 2 b347 `·` memegang
// `local.plan` satu baris; loop dalam 2.1 b524 `·` menelusuri SELURUH daftar dan 2.1.2 b715
// `·` (PRE b836 `local.plan==.Plan && local.idx!=.pxListSubscript`) memasang "sama" di baris
// dalam. Akibatnya SETIAP baris yang plannya muncul di baris lain ditandai - kedua baris
// pasangan ganda, termasuk dua baris yang sama-sama kosong (audit 02-10-2026; dulu hanya
// kemunculan kedua dan seterusnya, baris kosong dilewati).
func periksaPlan(pk *periksa, daftar []models.BarisPlan) {
	jumlah := map[string]int{}
	for _, b := range daftar {
		jumlah[strings.TrimSpace(b.Plan)]++
	}
	for i, b := range daftar {
		plan := strings.TrimSpace(b.Plan)
		if plan == "" {
			pk.baris(judulPlan, i, "%s", PesanPlanKosong)
		}
		if jumlah[plan] > 1 {
			pk.baris(judulPlan, i, "%s", PesanPlanSama)
		}
		if strings.TrimSpace(b.RIRate) == "" {
			pk.baris(judulPlan, i, "%s", PesanRIRatePlanKosong)
		}
	}
}

// gerbangPlan - kapan `ProteksiPlanListLife` dijalankan saat simpan. Di Pega ia
// berjalan pada onChange `.Plan` b33163 dan setiap `Delete` baris `PLAN LIST` b35214,
// `FINANCIAL UNDERWRITING` b40024, `UNDERWRITING LIMIT` b45682 - bukan pada
// `SaveProductName_Act`. Padanannya di sini: `PLAN LIST` berubah dari yang tersimpan
// (baris ditambah, diubah, dihapus), atau baris salah satu grid lain berkurang.
// Produk lama yang daftar plannya tidak disentuh tetap dapat disimpan.
func gerbangPlan(m *models.Produk, tersimpan models.Produk) bool {
	if len(m.FinancialUnderwriting) < len(tersimpan.FinancialUnderwriting) ||
		len(m.UnderwritingLimit) < len(tersimpan.UnderwritingLimit) {
		return true
	}
	if len(m.PlanList) != len(tersimpan.PlanList) {
		return true
	}
	for i, b := range m.PlanList {
		s := tersimpan.PlanList[i]
		if b.Plan != s.Plan || b.PlanID != s.PlanID || b.Name != s.Name || b.Benefit != s.Benefit ||
			b.RIRate != s.RIRate || b.RIRateID != s.RIRateID {
			return true
		}
	}
	return false
}

// Judul grid VERBATIM untuk pesan `asli` dan pesan kolom flat.
const (
	judulLien     = "LIEN CLAUSE (Potongan Manfaat Klaim)" // b12201
	judulDokumen  = "DOCUMENT CLAIM"                       // b14601
	judulKomentar = "Comment"                              // popup `SaveProductName_Confirm` b1025
	judulOutward  = "On Retention"                         // checkbox b47312 → `GetReinsTypeOR_Life`
)

// pesanAsliAsing - `asli` baris (kunci JSON lama yang tidak dikelola layar) yang
// tidak berasal dari baris tersimpan produk ini.
const pesanAsliAsing = "carries stored keys (asli) that do not belong to any stored row of this product"

// periksaAsli - `asli` hanya boleh DIKEMBALIKAN klien, tidak dikarang: setiap
// `asli` berisi harus sama dengan `asli` salah satu baris tersimpan daftar yang
// sama (produk yang diubah, atau produk asal `Copy`). Tanpa ini klien dapat
// menyisipkan kunci sembarang ke `JSONDATA` - termasuk kunci yang dibaca hilir.
func periksaAsli(pk *periksa, m *models.Produk, tersimpan models.Produk) {
	cek := func(judul string, kiriman, simpan []string) {
		sah := map[string]bool{}
		for _, a := range simpan {
			sah[a] = true
		}
		for i, a := range kiriman {
			if a != "" && !sah[a] {
				pk.baris(judul, i, "%s", pesanAsliAsing)
			}
		}
	}
	cek(judulLien, asliDari(m.LienClause, func(b models.BarisLien) string { return string(b.Asli) }),
		asliDari(tersimpan.LienClause, func(b models.BarisLien) string { return string(b.Asli) }))
	cek(judulDokumen, asliDari(m.DocumentClaim, func(b models.BarisDokumen) string { return string(b.Asli) }),
		asliDari(tersimpan.DocumentClaim, func(b models.BarisDokumen) string { return string(b.Asli) }))
	cek(judulPlan, asliDari(m.PlanList, func(b models.BarisPlan) string { return string(b.Asli) }),
		asliDari(tersimpan.PlanList, func(b models.BarisPlan) string { return string(b.Asli) }))
	cek(judulFinUW, asliDari(m.FinancialUnderwriting, func(b models.BarisFinUW) string { return string(b.Asli) }),
		asliDari(tersimpan.FinancialUnderwriting, func(b models.BarisFinUW) string { return string(b.Asli) }))
	cek(judulUWLimit, asliDari(m.UnderwritingLimit, func(b models.BarisUWLimit) string { return string(b.Asli) }),
		asliDari(tersimpan.UnderwritingLimit, func(b models.BarisUWLimit) string { return string(b.Asli) }))
}

func asliDari[T any](daftar []T, ambil func(T) string) []string {
	hasil := make([]string, len(daftar))
	for i, b := range daftar {
		hasil[i] = ambil(b)
	}
	return hasil
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
// diterima apa adanya. R/I Rate BARU (K1 keputusan work owner 01-10-2026, OQ-MPNL-03) wajib ada di view
// `RATE_LIFE_SUMMARY`; namanya = `.USEDBY` master, seperti `SetRIRate` b249 (`.RIRATE ← usedby`).
func (l *Layanan) periksaPilihanPlan(ctx context.Context, pk *periksa, daftar []models.BarisPlan, lama []models.BarisPlan) error {
	rateLama, planLama := map[[2]string]bool{}, map[[2]string]bool{}
	// Satu pembacaan view per RIRATEID per simpan (code review #15) - baris plan sering ber-R/I Rate sama.
	type hasilRate struct {
		v   models.NilaiMaster
		ada bool
	}
	rateDibaca := map[string]hasilRate{}
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
		b.RIRateID = strings.TrimSpace(b.RIRateID)
		if rateLama[[2]string{b.RIRateID, b.RIRate}] {
			continue
		}
		if b.RIRateID != "" {
			r, sudah := rateDibaca[b.RIRateID]
			if !sudah {
				v, ada, err := l.gudang.AmbilMaster(ctx, models.MasterRIRate, b.RIRateID)
				if err != nil {
					return err
				}
				r = hasilRate{v, ada}
				rateDibaca[b.RIRateID] = r
			}
			v, ada := r.v, r.ada
			if !ada {
				pk.baris(judulPlan, i, "R/I Rate %q is not in the master list", b.RIRateID)
			} else {
				b.RIRate = v.Nama
			}
		} else if strings.TrimSpace(b.RIRate) != "" {
			pk.baris(judulPlan, i, "R/I Rate %q must be chosen from the master list", b.RIRate)
		}
	}
	return nil
}

// Judul grid VERBATIM paket 7.
const judulFinUW = "FINANCIAL UNDERWRITING" // b37148

// periksaFinUW - angka dan rentang baris `FINANCIAL UNDERWRITING` (`Min Insured`
// b37670, `Max Insured` b37818 = `pxNumber`; `Employee`, `Non-Employee` teks).
func periksaFinUW(pk *periksa, daftar []models.BarisFinUW) {
	for i := range daftar {
		b := &daftar[i]
		var sub periksa
		minSI, maxSI := sub.desimal("Min Insured", &b.MinInsured), sub.desimal("Max Insured", &b.MaxInsured)
		sub.tidakLebihBesar("Min Insured", minSI, "Max Insured", maxSI)
		for _, s := range sub.pesan {
			pk.baris(judulFinUW, i, "%s", s)
		}
	}
}

// bentukWaktuPega - `@CurrentDateTime()` Pega (`YYYYMMDDTHHMMSS.SSS GMT`).
const bentukWaktuPega = "20060102T150405.000 GMT"

// barisKomentar - `AddCommentList_Act` 1 b235 `·`: `Date = @CurrentDateTime()`,
// `OperatorName = OperatorID.pxInsName` (akun pelaku), `IsApproved = param.status`
// (tidak dikirim `SaveProductName_Act` 7 b1515 → kosong), `Suggest = param.comment`
// (`ProductName.Comment`).
func barisKomentar(saat time.Time, akun, komentar string) models.BarisKomentar {
	return models.BarisKomentar{Date: saat.UTC().Format(bentukWaktuPega), OperatorName: akun, Suggest: komentar}
}
