package services

// Proteksi Submit case MARINE CARGO (tiket 74 P5, keputusan work owner 10-10-2026 K74-7…K74-13). Port `[terverifikasi]`
// `D:\migrasi\RNM\NB FacIn\Activity\` (riset sesi c3 10-10-2026):
//   - proteksiCoverageCargo = ProtectCoverage_Act langkah 4 MARINE (IsMarineCargo L1808; L1228-1673): per cargo × coverage
//     " Coverage "+IndexCargo+" can't be empty!" bila Coverage / CoverageNote kosong, dan "Currency on object "+IndexCargo+" coverage
//     "+n+" must be choose!" bila Currency.ID / OldID kosong (flag ProtectMarine.CARI2). Nomor cargo = urutan 1..n, coverage kosong
//     (CoverageList size 0) TIDAK diproteksi - tiru Pega (K74-12).
//   - proteksiKapalCargo = Protection_Act langkah 27 MARINE CARGO (L4841-5571): 27.1 ProtectShipData_Act (IsMOP != "MOP" - polis
//     individual, K74-3; L676-1074) " Ship data CargoList " + IdxCargo + " can't  be empty!" (DUA spasi, verbatim) bila Ship.ID /
//     NM_SHIP kosong (flag ProtectMarine.CARI1); 27.2-27.4 "Multi Location Must Have Toprisk" bila > 1 cargo dan tidak satu pun
//     IsTopRisk.
//   - pesanSpreadingCargo = CheckSpreadingProtectMCargoMBU_ACT langkah 2 (IsMarineCargo L3501, IsMOP != "MOP" L3524): Σshare TIDAK
//     dibulatkan vs 100 (L1150, pesan L1964); Σ PremiumSpreaded dan PremiNusantaraRe @divide(x,1,2) (L1342 / L1394, pesan L2135)
//     (K74-13). Cabang IsGroup / operator literal (A14) tidak diport. SpreadingList kosong -> Σ 0 -> pesan 100 %.
//   - IsBedaCurr Team Leader: addDeductibleCargo_ACT men-SET `pyWorkPage.IsBedaCurr = 0` (L248) - case cargo tidak pernah menahan
//     "Deductible Acceptance cannot be empty".
// Langkah IsFire (18 Scoring Risk, 21 ProtectCurrencyTSI, 24 ProtectFIREMBUPA) tidak berjalan untuk case cargo.

import (
	"context"
	"fmt"
	"strings"

	"nusantarare/modul/nbfacin/backend/models"
)

// flagMarine1 - flag pesan ProtectShipData_Act (flagMarine2 di proteksi.go).
const flagMarine1 flagProteksi = "ProtectMarine.CARI1"

// dataProteksiCargo - pohon cargo case MARINE CARGO; nil bila bukan cargo (atau penyimpan cargo tidak terpasang): pemanggil
// menjalankan proteksi lain apa adanya.
func (s *Service) dataProteksiCargo(ctx context.Context, id string) ([]models.ObjekCargo, bool, error) {
	if s.objekCargo == nil || !s.adalahMarineCargo(ctx, id) {
		return nil, false, nil
	}
	cargo, err := s.objekCargo.BacaObjekCargo(ctx, id)
	if err != nil {
		return nil, false, galatKasus(err)
	}
	return cargo, true, nil
}

// proteksiCoverageCargo - ProtectCoverage_Act langkah 4 MARINE (lihat kepala berkas). Currency.ID / OldID = master CURRENCY atas
// nama tersimpan; tidak ada di master = kosong.
func (s *Service) proteksiCoverageCargo(ctx context.Context, cargo []models.ObjekCargo, p *pengumpulPesan) error {
	ada := map[string]bool{}
	for i, c := range cargo {
		for n, v := range c.Coverages {
			if strings.TrimSpace(v.Coverage) == "" || strings.TrimSpace(v.CoverageNote) == "" {
				p.tambahFlag(fmt.Sprintf(" Coverage %d can't be empty!", i+1), flagMarine2)
			}
			mu := strings.TrimSpace(v.Currency)
			sah, dikenal := ada[mu]
			if mu != "" && !dikenal {
				m, ok, err := s.objekCargo.MataUangCargo(ctx, mu)
				if err != nil {
					return err
				}
				sah = ok && m.ID != "" && m.OldID != ""
				ada[mu] = sah
			}
			if mu == "" || !sah {
				p.tambahFlag(fmt.Sprintf("Currency on object %d coverage %d must be choose!", i+1, n+1), flagMarine2)
			}
		}
	}
	return nil
}

// proteksiKapalCargo - Protection_Act langkah 27 (lihat kepala berkas).
func proteksiKapalCargo(cargo []models.ObjekCargo, p *pengumpulPesan) {
	for i, c := range cargo {
		if strings.TrimSpace(c.Ship.ID) == "" || strings.TrimSpace(c.Ship.NmShip) == "" {
			p.tambahFlag(fmt.Sprintf(" Ship data CargoList %d can't  be empty!", i+1), flagMarine1)
		}
	}
	if len(cargo) > 1 {
		for _, c := range cargo {
			if c.IsTopRisk {
				return
			}
		}
		p.tambah("Multi Location Must Have Toprisk")
	}
}

// pesanSpreadingCargo - CheckSpreadingProtectMCargoMBU_ACT langkah 2 per cargo × coverage (lihat kepala berkas).
func pesanSpreadingCargo(cargo []models.ObjekCargo) []string {
	pesan := []string{}
	for i, c := range cargo {
		for n, v := range c.Coverages {
			total, premi := nol(), nol()
			for _, b := range v.Spreading {
				total, premi = tambah(total, b.SharePercentage), tambah(premi, b.PremiumSpreaded)
			}
			if total.Cmp(seratus) != 0 {
				pesan = append(pesan, fmt.Sprintf("%% Total spreading on cargo %d coverage %d must be 100%%!", i+1, n+1))
			}
			if bagiSkala(premi, satu, 2).Cmp(bagiSkala(nolBila(v.PremiNusantaraRe), satu, 2)) != 0 {
				pesan = append(pesan, fmt.Sprintf("Total premium spreaded on cargo %d coverage %d is not equals premium RNM!", i+1, n+1))
			}
		}
	}
	return pesan
}
