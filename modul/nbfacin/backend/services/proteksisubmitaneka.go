package services

// Proteksi Submit case keluarga Aneka (tiket 73 P4). Port `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\Activity\` (dibaca
// 09-10-2026; rincian per langkah di docs/issues/73 "P4"):
//   - proteksiCoverageAneka = ProtectCoverage_Act (`ASM-FW-GISFW-WORK!PROTECTCOVERAGE_ACT`, dipanggil Protection_Act langkah 24)
//     langkah 5 ANEKA (IsAneka, L1844) dan langkah 6 GOLF (isGolfInsurance, L6626). Langkah 5.1.2-5.1.9 (Province / City /
//     ZipCode) ber-label `//` - tidak diport. Langkah 7 ProtectFIREMBUPA_Act bergerbang IsFire di dalamnya.
//   - pesanSpreadingAneka = CheckSpreadingProtectAnekaGolf_ACT (`ASM-FW-GISFW-WORK!CHECKSPREADINGPROTECTANEKAGOLF_ACT`, dari
//     CheckSpreadingProtect_ACT langkah 6): cabang ANEKA (IsAneka || IsBondingAndCustomBonds, skala 4) dan GOLF (skala 2,
//     K73-38 HALF_UP). Cabang Group / user-ID hardcode (L3937 / L4139) tidak diport (pola K48-6).
//
// Keputusan work owner 09-10-2026: K73-28 (proteksi FIRE hanya IsFire), K73-29 (CheckSpreadingProtect di Binding semua COB),
// K73-30 (IsBedaCurr Aneka `[dugaan]`), K73-31 (kasus ganda tanpa case sendiri), K73-32 (nomor objek Golf diisi).
//
// Penyimpangan `[keputusan agent]` (dilaporkan ke 9d): mata uang coverage kosong = mata uang objek (Count
// CountPremiAndTSIRNMFireAnekaGolf_ACT L2969 mengisinya dari objek untuk delapan lini; layar coverage Aneka hanya menampilkan
// Currency untuk lini daftar, K73 P2). Model Aneka menyimpan CURRENCY_CODE, bukan `.Currency.ID` / `.OldID`: "kosong" = kode
// kosong.

import (
	"context"
	"fmt"
	"strings"

	"nusantarare/modul/nbfacin/backend/models"
)

// dataAnekaProteksi - pohon objek Aneka + isian ProtectCoverage untuk Protection_Act.
type dataAnekaProteksi struct {
	lokasi []models.LokasiAneka
	isi    isianProteksiAneka
}

// dataProteksiAneka - nil bila case bukan keluarga Aneka (atau lini tak terbaca / penyimpan objek Aneka tidak terpasang):
// pemanggil menjalankan proteksi FIRE apa adanya. deductibleBinding = EmailTypeBinding != 4 && FlagOnGoingPolicy == 2.
func (s *Service) dataProteksiAneka(ctx context.Context, id string, kasus models.Kasus, deductibleBinding bool) (*dataAnekaProteksi, error) {
	if s.objekAneka == nil {
		return nil, nil
	}
	kb, ada := s.kasusBisnisDari(ctx, id)
	if !ada {
		return nil, nil
	}
	l := liniSpreadingDari(id, kb)
	if !l.aneka && !l.golf {
		return nil, nil
	}
	lokasi, err := s.objekAneka.BacaObjekAneka(ctx, id)
	if err != nil {
		return nil, galatKasus(err)
	}
	k, err := s.konteksSubmit(ctx, id)
	if err != nil {
		return nil, err
	}
	// ProtectCoverage_Act langkah 2 (L588-880): ParamInProdCase.CARI10-15; IsEDM / IsRenewal tidak ada di NB sistem baru.
	cari := &models.CariKasusGanda{InsuredName: strings.ReplaceAll(strings.ReplaceAll(k.InsuredName, ".", ""), " ", ""),
		SobName: k.SobName, BusinessCode: kb.businessCode, NoOfferSlip: kasus.General.ReffNumber,
		Mulai: tanggalAkumulasi(kasus.General.StartDateTime), Akhir: tanggalAkumulasi(kasus.General.EndDateTime), IDPegaSendiri: idPega(id)}
	return &dataAnekaProteksi{lokasi: lokasi, isi: isianProteksiAneka{lini: l, deductible: deductibleBinding && !s.grup.kasusBonding(k),
		cariGanda: cari}}, nil
}

// mataUangCoverage - Currency coverage, kosong -> Currency objek (lihat kepala berkas).
func mataUangCoverage(o models.ObjekAneka, c models.CoverageAneka) string {
	if u := strings.TrimSpace(c.Currency); u != "" {
		return u
	}
	return strings.TrimSpace(o.Currency)
}

// deductibleBedaMataUangAneka - K73-30 `[dugaan]` (SaveDeductibleAneka tidak ada di ekspor): ada deductible ber-Currency terisi
// (bukan 'UNKNOWN') yang beda dari Currency coverage-nya; padanan deductibleBedaMataUang FIRE.
func deductibleBedaMataUangAneka(lokasi []models.LokasiAneka) bool {
	beda := false
	setiapCoverageAneka(lokasi, true, true, func(_ int, o *models.ObjekAneka, _ int, c *models.CoverageAneka) {
		for _, d := range c.Deductibles {
			if cur := strings.TrimSpace(d.Currency); cur != "" && cur != mataUangKosongAneka && cur != mataUangCoverage(*o, *c) {
				beda = true
			}
		}
	})
	return beda
}

// mataUangKosongAneka - DEFAULT kolom CURRENCY deductible (A169).
const mataUangKosongAneka = "UNKNOWN"

// deductibleKosong - ProtectCoverage_Act 5.1.10.1.1.9.2 (L5149): `.TypeDeductible == "" && .Condition == "" && .MinMax == "" &&
// .Amount == "" && .TimeExcess == ""`.
func deductibleKosong(d models.Deductible) bool {
	return strings.TrimSpace(d.TypeDeductible) == "" && strings.TrimSpace(d.Condition) == "" && strings.TrimSpace(d.MinMax) == "" &&
		d.Amount == nil && d.TimeExcess == nil
}

// isianProteksiAneka - data ProtectCoverage cabang Aneka di luar pohon objek.
type isianProteksiAneka struct {
	lini liniSpreadingAneka
	// deductible - 5.1.10.1.1.7 / .9: EmailTypeBinding != 4 && FlagOnGoingPolicy == 2 (= deductibleBinding FIRE) && !IsTBonding.
	deductible bool
	// cariGanda - nil = kasus ganda tidak diperiksa (data tidak lengkap / penyimpan tidak ada).
	cariGanda *models.CariKasusGanda
}

// proteksiCoverageAneka - ProtectCoverage_Act langkah 5 (ANEKA) dan 6 (GOLF), pesan VERBATIM (termasuk spasi depan):
//   - per coverage (5.1.10.1.1): Coverage / CoverageNote kosong (L3779), Currency kosong (L4229), deductible belum ditambah
//     (L4529, flag CARI4), baris deductible kosong (L5030, flag CARI2 - DUA spasi sebelum "DeductibleList");
//   - per objek (5.1.10.1.2): Currency objek kosong (L5482; `.IdxLocation` = nomor lokasi `[dugaan]`, SetIndexObject_DT);
//   - per lokasi (5.1.12-14): GetDataDoubleCase -> "Double Input Case <IDPEGA>" (L6147) lalu LOMPAT ke label AA (L6372):
//     lokasi berikutnya dan blok GOLF dilewati;
//   - GOLF per coverage L6920 / per objek L7281.
//
// Nomor = .pxListSubscript (mulai 1): lokasi, objek di dalam okupasinya (IdxAneka), coverage, deductible.
func (s *Service) proteksiCoverageAneka(ctx context.Context, lokasi []models.LokasiAneka, isi isianProteksiAneka, p *pengumpulPesan) error {
	if isi.lini.aneka {
		for i, l := range lokasi {
			for _, ok := range l.Okupasi {
				for k, o := range ok.Objek {
					for m, c := range o.Coverages {
						if strings.TrimSpace(c.Coverage) == "" || strings.TrimSpace(c.CoverageNote) == "" {
							p.tambahFlag(fmt.Sprintf(" Coverage on location %d can't be empty!", i+1), flagMarine2)
						}
						if mataUangCoverage(o, c) == "" {
							p.tambahFlag(fmt.Sprintf(" Currency on location %d object item %d Coverage %d can't be empty!", i+1, k+1, m+1), flagMarine2)
						}
						if !isi.deductible {
							continue
						}
						if len(c.Deductibles) == 0 {
							p.tambahFlag(fmt.Sprintf(" Deductible on location %d object item %d Coverage %d should be added!", i+1, k+1, m+1), flagDeductible4)
						}
						for n, d := range c.Deductibles {
							if deductibleKosong(d) {
								p.tambahFlag(fmt.Sprintf(" Deductible on location %d Coverage %d  DeductibleList %d can't be empty!", i+1, m+1, n+1), flagDeductible2)
							}
						}
					}
					if strings.TrimSpace(o.Currency) == "" {
						p.tambahFlag(fmt.Sprintf(" Currency on location %d can't be empty!", i+1), flagMarine2)
					}
				}
			}
			if isi.cariGanda != nil && s.submit != nil {
				id, err := s.submit.KasusGanda(ctx, *isi.cariGanda)
				if err != nil {
					return err
				}
				if len(id) > 0 {
					p.tambahFlag("Double Input Case "+id[0], flagMarine2)
					return nil // TRANS L6372 T=1 -> AA
				}
			}
		}
	}
	if isi.lini.golf {
		for i, l := range lokasi {
			for _, o := range l.ObjekGolf {
				for _, c := range o.Coverages {
					if strings.TrimSpace(c.Coverage) == "" || strings.TrimSpace(c.CoverageNote) == "" {
						p.tambahFlag(fmt.Sprintf(" Coverage on location %d can't be empty!", i+1), flagMarine2)
					}
				}
				if strings.TrimSpace(o.Currency) == "" {
					p.tambahFlag(fmt.Sprintf(" Currency on location %d can't be empty!", i+1), flagMarine2)
				}
			}
		}
	}
	return nil
}

// pesanSpreadingAneka - CheckSpreadingProtectAnekaGolf_ACT per coverage, urutan Pega: baris share / treaty kosong (L1730 /
// Golf L5623), Σ share ≠ 100 (L2269 / L6162), Σ PremiumSpreaded ≠ PremiNusantaraRe (L2289 / L6188; keduanya @divide(x,1,4),
// Golf @divide(x,1,2) L5918 / L5970). Aneka: "object item" pesan baris kosong = nomor OKUPASI (idxobjitem, verbatim); Golf:
// idxobjitem tidak di-SET Pega -> diisi nomor objek (K73-32). SpreadingList kosong -> Σ 0 -> pesan 100 %.
func pesanSpreadingAneka(lokasi []models.LokasiAneka, l liniSpreadingAneka) []string {
	pesan := []string{}
	cek := func(c models.CoverageAneka, skala int32, kosong, seratusPersen, premiRNM string) {
		total, premi := nol(), nol()
		for _, b := range c.Spreading {
			total, premi = tambah(total, b.SharePercentage), tambah(premi, b.PremiumSpreaded)
			if b.SharePercentage == nil || strings.TrimSpace(b.TreatyType) == "" {
				pesan = append(pesan, kosong)
			}
		}
		if total.Cmp(seratus) != 0 {
			pesan = append(pesan, seratusPersen)
		}
		if bagiSkala(premi, satu, skala).Cmp(bagiSkala(nolBila(c.PremiNusantaraRe), satu, skala)) != 0 {
			pesan = append(pesan, premiRNM)
		}
	}
	if l.aneka {
		for i, x := range lokasi {
			for j, ok := range x.Okupasi {
				for k, o := range ok.Objek {
					for m, c := range o.Coverages {
						cek(c, 4, fmt.Sprintf("%% spreading on location %d object item %d coverage %d can't be empty!", i+1, j+1, m+1),
							fmt.Sprintf("%% Total spreading on location %d occupation %d object %d coverage %d must be 100%%!", i+1, j+1, k+1, m+1),
							fmt.Sprintf("Total premium spreaded on location %d occupation %d object %d coverage %d is not equals premium RNM!", i+1, j+1, k+1, m+1))
					}
				}
			}
		}
	}
	if l.golf {
		for i, x := range lokasi {
			for k, o := range x.ObjekGolf {
				for m, c := range o.Coverages {
					cek(c, 2, fmt.Sprintf("%% spreading on location %d object item %d coverage %d can't be empty!", i+1, k+1, m+1),
						fmt.Sprintf("%% Total spreading on location %d object %d coverage %d must be 100%%!", i+1, k+1, m+1),
						fmt.Sprintf("Total premium spreaded on location %d object %d coverage %d is not equals premium RNM!", i+1, k+1, m+1))
				}
			}
		}
	}
	return pesan
}
