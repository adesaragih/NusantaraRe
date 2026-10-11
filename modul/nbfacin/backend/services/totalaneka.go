package services

// Total per mata uang tab Coverage Aneka (tiket 73 P2) - port `NB FacIn\Activity\SumTotalTSIPremiGross_Act.xml`
// (ASM-FW-GISFW-DATA-COVERAGE!SUMTOTALTSIPREMIGROSS_ACT) `[terverifikasi]` (riset P2 bagian 1, dicek ulang sesi c3): hasil per
// baris LocationList x nama mata uang (okupasi dalam satu lokasi dilebur), urutan kemunculan pertama; premi = `Coverage.Premium`
// mentah (tanpa diskon); TSI / Premium dijumlah tanpa pembulatan; Rate = `@if(.TSI=0,0,@divide(.Premium,.TSI,20)*100)` (persen,
// L3058 / L4583; FIRE x1000).
//   - Langkah 2 GOLF (`IsGolfInsurance` L3284): `.Property.RiskLocation.AnekaList` - TSI = `Aneka.TSI` per mata uang Aneka, premi
//     coverage masuk ke mata uang ANEKA.
//   - Langkah 3 ANEKA (`IsAneka` || `IsBondingAndCustomBonds`): `.OccupationList.AnekaList` - bila TIDAK satu pun dari 10 When
//     khusus (3a): seperti Golf; bila salah satu (3b): TSI dan premi per COVERAGE, mata uang coverage.
// Ke-10 When membaca BusinessType / BusinessCode CASE, jadi pilihan 3a / 3b sama untuk semua objek case. Bila gerbang 2 dan 3
// sama-sama terbuka, langkah 3 (yang lebih akhir, menghapus list lebih dulu) menang. Tanpa gerbang terbuka: list kosong.
// Tidak disimpan (baca-saja, pola TotalPerMataUang FIRE); nilai simpan8.

import (
	"github.com/cockroachdb/apd/v3"

	"nusantarare/modul/nbfacin/backend/models"
)

// khususTotalAneka - 10 When langkah 3.3.2 (L4768-L4975), nama registry.
var khususTotalAneka = []string{"IsObjectWithQuantityYear", "isAllrisk", "IsMarineHull", "IsMarineHullOffshore", "IsMBD", "IsHE",
	"isAviationHull", "isBurglary", "IsPortRisk", "IsBuilderRisk"}

// akumulatorTotal - baris per mata uang, urutan kemunculan pertama.
type akumulatorTotal struct {
	baris  []models.TotalCoverage
	posisi map[string]int
}

func (a *akumulatorTotal) ke(mataUang string) *models.TotalCoverage {
	i, ada := a.posisi[mataUang]
	if !ada {
		i = len(a.baris)
		a.posisi[mataUang] = i
		a.baris = append(a.baris, models.TotalCoverage{Currency: mataUang, TSI: apd.New(0, 0), Premium: apd.New(0, 0)})
	}
	return &a.baris[i]
}

func (a *akumulatorTotal) hasil() []models.TotalCoverage {
	for i := range a.baris {
		t := &a.baris[i]
		if t.TSI.IsZero() {
			t.Rate = apd.New(0, 0)
		} else if r, err := bagi(t.Premium, t.TSI, skalaBagi); err == nil {
			t.Rate = simpan8(kali(r, seratus))
		}
		t.TSI, t.Premium = simpan8(t.TSI), simpan8(t.Premium)
	}
	return a.baris
}

// tambahObjek - satu objek: per objek (TSI objek, premi ke mata uang objek) atau per coverage (`perCoverage`, 3b).
func (a *akumulatorTotal) tambahObjek(o models.ObjekAneka, perCoverage bool) {
	if !perCoverage {
		t := a.ke(o.Currency)
		konteksCoverage.Add(t.TSI, t.TSI, nolBila(o.TSI))
		for _, c := range o.Coverages {
			konteksCoverage.Add(t.Premium, t.Premium, nolBila(c.Premium))
		}
		return
	}
	for _, c := range o.Coverages {
		t := a.ke(c.Currency)
		konteksCoverage.Add(t.TSI, t.TSI, nolBila(c.TSI))
		konteksCoverage.Add(t.Premium, t.Premium, nolBila(c.Premium))
	}
}

// totalAnekaLokasi - total per mata uang satu lokasi menurut gerbang case.
func totalAnekaLokasi(id string, kb kasusBisnis, l models.LokasiAneka) []models.TotalCoverage {
	a := &akumulatorTotal{posisi: map[string]int{}}
	switch {
	case salahSatu(id, kb, "IsAneka", "IsBondingAndCustomBonds"):
		khusus := salahSatu(id, kb, khususTotalAneka...)
		for _, ok := range l.Okupasi {
			for _, o := range ok.Objek {
				a.tambahObjek(o, khusus)
			}
		}
	case salahSatu(id, kb, "IsGolfInsurance"):
		for _, o := range l.ObjekGolf {
			a.tambahObjek(o, false)
		}
	}
	return a.hasil()
}
