package services

// Tab Spreading keluarga Aneka (tiket 73 P3): Aneka / Bonding / Engineering (lokasi -> okupasi -> objek -> coverage) dan Golf
// (lokasi -> objek -> coverage). Port rantai Pega `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\` (dibaca 09-10-2026; rincian
// per langkah di docs/issues/73 "Backend P3"):
//   - hitungNRAneka         = CountPremiAndTSIRNMFireAnekaGolf_ACT 7.2.2.1.1.3.4 (L3853 / L3906) / 8 (L6472): TSI NR dari .TSI;
//   - terapkanPersenAneka   = DT SetTSIPremiSpreaded_FacIn blok 8 Golf (L3421) / 9 Engineering (L4094) / 10 Aneka (L4798) / 13
//                             Bonding (L7186) mode "percent": TSIGrossSpreaded TIDAK disentuh (K73-18);
//   - salinTemplateAneka    = CopyToAllSpreadingAnekaGolf_ACT (L1098 `local.TSI = .TSI*PercentShare/100`, L3404 / L3450);
//   - spreadingOtomatisGolf = GetKapasitasTreaty cabang isGolfInsurance (DDL L11453): ORS 10007 100 % per coverage;
//   - spreadingOtomatisEngineering = GetKapasitasTreaty cabang IsEngineering (DDL L21510), kerangka cabang IsFire dengan lima
//                             beda verbatim (lihat fungsinya);
//   - totalLokasiAneka      = SumTSIPremiSpreadedRNM_ANEKA_Act cabang 2.4 (tanpa top risk, L9275) / 2.5 (ada top risk, L15995);
//   - totalLokasiGolf       = SumTSIPremiSpreadedRNM_GOLF_Act 2.3;
//   - total case / ringkasan / batas treaty / jenis treaty = fungsi FIRE (totalMataUang, ringkasan, periksaJenisTreaty):
//     SumTSIPremiSpreadedRNM_Act langkah 17-22 tidak bergerbang lini.
//
// Penyimpangan sadar `[keputusan agent]` (dilaporkan ke work owner lewat 9d; semuanya hanya menggigit case multi-mata-uang):
//   - kunci mata uang total = mata uang coverage untuk lini daftar (bila kosong -> mata uang objek), selain itu mata uang objek;
//     Pega membawa CurrObj lintas coverage / objek (ANEKA L4597, GOLF tanpa reset) dan APPEND total Golf tanpa penjaga mata
//     uang (L5797);
//   - (K73-21 keputusan work owner 09-10-2026: penyimpangan multi-mata-uang di atas DISETUJUI);
//   - langkah 24 dispatcher (IsMBD) di-port apa adanya (K73-22): hapusPesanMBD;
//   - EDM / IsEdmAdjSpreading, syariah, FlagDelete: di luar cakupan (pola K48-6).

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// Daftar lini When Pega (verbatim, urutan sumber).
var (
	// liniMataUangCoverage24 - SumTSIPremiSpreadedRNM_ANEKA_Act 2.4 "SET CURR OBJ" (L4117-4347 / L4744-4974): lini ber-mata uang
	// total dari coverage.
	liniMataUangCoverage24 = []string{"IsMBD", "IsHE", "isBurglary", "IsObjectWithQuantityYear", "IsLiability", "isAllRisk",
		"IsMarineHull", "IsMarineHullOffshore", "IsAviationHull", "IsBuilderRisk", "IsPortRisk"}
	// liniMataUangCoverage25 - cabang 2.5 (L10410-10686 / L10936-11359): + IsElectronicEquipment, IsCMI.
	liniMataUangCoverage25 = append(append([]string{}, liniMataUangCoverage24...), "IsElectronicEquipment", "IsCMI")
	// liniTopRisk24 / liniTopRisk25 - TSI top risk (ClaimSpreaded) 2.4 (L5940-6009) / 2.5 (L12454-12592).
	liniTopRisk24 = []string{"IsMBD", "IsObjectWithQuantityYear", "IsLiability", "IsHE"}
	liniTopRisk25 = []string{"IsMBD", "IsObjectWithQuantityYear", "IsLiability", "IsElectronicEquipment", "IsCMI", "IsMarineHull", "IsHE"}
	// liniMataUangCount - CountPremiAndTSIRNMFireAnekaGolf_ACT 7.2.2.1.1.3.3 (L3477-3504): kunci mata uang = coverage; juga
	// 7.2.2.1.1.2 (L2556): TSI objek TIDAK dijumlah untuk lini ini.
	liniMataUangCount = []string{"IsMBD", "isBurglary", "IsAviationHull", "IsHE", "IsMarineHullOffshore", "IsMarineHull", "isAllRisk",
		"IsObjectWithQuantityYear", "IsBuilderRisk", "IsPortRisk"}
	// liniTSICoverageCount - 7.2.2.1.1.3.7 (L4324-4612): TSI dijumlah per coverage (BuilderRisk / PortRisk tidak - verbatim).
	liniTSICoverageCount = []string{"IsMBD", "IsObjectWithQuantityYear", "isAllRisk", "IsMarineHull", "IsMarineHullOffshore", "IsHE",
		"IsAviationHull", "isBurglary"}
)

// liniSpreadingAneka - When lini case yang dibaca rantai spreading Aneka (dievaluasi sekali per permintaan).
type liniSpreadingAneka struct {
	golf, engineering, aneka, mbd   bool // aneka = IsAneka || IsBondingAndCustomBonds (dispatcher langkah 10, L2090 / L2113)
	mataUangCov24, mataUangCov25    bool
	topRisk24, topRisk25            bool
	mataUangCount, tsiCoverageCount bool
}

func liniSpreadingDari(id string, kb kasusBisnis) liniSpreadingAneka {
	return liniSpreadingAneka{golf: salahSatu(id, kb, "isGolfInsurance"), engineering: salahSatu(id, kb, "IsEngineering"),
		aneka: salahSatu(id, kb, "IsAneka", "IsBondingAndCustomBonds"), mbd: salahSatu(id, kb, "IsMBD"),
		mataUangCov24: salahSatu(id, kb, liniMataUangCoverage24...), mataUangCov25: salahSatu(id, kb, liniMataUangCoverage25...),
		topRisk24: salahSatu(id, kb, liniTopRisk24...), topRisk25: salahSatu(id, kb, liniTopRisk25...),
		mataUangCount: salahSatu(id, kb, liniMataUangCount...), tsiCoverageCount: salahSatu(id, kb, liniTSICoverageCount...)}
}

// LokasiSpreadingAneka - satu lokasi tab Spreading Aneka: okupasi -> objek (Aneka / Bonding / Engineering), objek Golf, total
// per lokasi (`.Property.TotalTSIPremiSpreadRNM`; Currency = mata uang, Pega menyimpannya di TreatyName).
type LokasiSpreadingAneka struct {
	ObjectNo, Location string
	Okupasi            []models.OkupasiAneka
	ObjekGolf          []models.ObjekAneka
	Total              []models.TotalSpreading
}

// TampilanSpreadingAneka - jawaban seluruh endpoint tab Spreading Aneka. TampilTotalLokasi = K73-19: grid total per lokasi
// tampil hanya `IsFire || IsMBD` di tahap uw (`PropertyItemListCoverageSpreading_IsUW.xml` L17067; non-uw IsFire saja L20592).
type TampilanSpreadingAneka struct {
	PercentShare      *apd.Decimal
	Treaty            []models.TreatySpreading
	Template          []models.TemplateSpreading
	Lokasi            []LokasiSpreadingAneka
	RingkasanTreaty   []models.TotalSpreading
	RingkasanMataUang []models.RingkasanMataUangSpreading
	Pesan             []string
	TampilTotalLokasi bool
}

// ---- jelajah -----------------------------------------------------------------------------------------------------

// setiapCoverageAneka - fn untuk setiap coverage objek okupasi (okupasi=true) dan / atau objek Golf (golf=true), urut lokasi.
func setiapCoverageAneka(lokasi []models.LokasiAneka, okupasi, golf bool, fn func(i int, o *models.ObjekAneka, m int, c *models.CoverageAneka)) {
	for i := range lokasi {
		if okupasi {
			for j := range lokasi[i].Okupasi {
				for n := range lokasi[i].Okupasi[j].Objek {
					o := &lokasi[i].Okupasi[j].Objek[n]
					for m := range o.Coverages {
						fn(i, o, m, &o.Coverages[m])
					}
				}
			}
		}
		if golf {
			for n := range lokasi[i].ObjekGolf {
				o := &lokasi[i].ObjekGolf[n]
				for m := range o.Coverages {
					fn(i, o, m, &o.Coverages[m])
				}
			}
		}
	}
}

// salinDalamAneka - salinan yang tidak berbagi larik coverage / baris spreading dengan pemanggil.
func salinDalamAneka(lokasi []models.LokasiAneka) []models.LokasiAneka {
	hasil := salinAneka(lokasi)
	setiapCoverageAneka(hasil, true, true, func(_ int, o *models.ObjekAneka, m int, c *models.CoverageAneka) {
		c.Spreading = salinBaris(c.Spreading)
	})
	return hasil
}

// adaObjekAneka - sekurangnya satu objek okupasi / Golf.
func adaObjekAneka(lokasi []models.LokasiAneka) bool {
	for _, l := range lokasi {
		if len(l.ObjekGolf) > 0 {
			return true
		}
		for _, o := range l.Okupasi {
			if len(o.Objek) > 0 {
				return true
			}
		}
	}
	return false
}

// ---- rumus --------------------------------------------------------------------------------------------------------

// hitungNRAneka - TSINusantaraRe = PS × TSI / 100; PremiNusantaraRe = PS × Premium / 100 − PS × Discount / 100 (skala 20), setiap
// coverage objek okupasi dan Golf.
func hitungNRAneka(lokasi []models.LokasiAneka, ps *apd.Decimal) {
	setiapCoverageAneka(lokasi, true, true, func(_ int, _ *models.ObjekAneka, _ int, c *models.CoverageAneka) {
		c.TSINusantaraRe = persenDari(ps, c.TSI)
		c.PremiNusantaraRe = kurang(persenDari(ps, c.Premium), persenDari(ps, c.Discount))
	})
}

// terapkanPersenAneka - DT mode "percent": TSISpreaded = share × TSINusantaraRe / 100, PremiumSpreaded = share × PremiNusantaraRe
// / 100. Golf -> objek Golf (blok 8); Engineering / Aneka / Bonding -> objek okupasi (blok 9 / 10 / 13). TSIGrossSpreaded dan
// ClaimEstimation tidak disentuh (grep DT = 0).
func terapkanPersenAneka(lokasi []models.LokasiAneka, l liniSpreadingAneka) {
	setiapCoverageAneka(lokasi, l.aneka || l.engineering, l.golf, func(_ int, _ *models.ObjekAneka, _ int, c *models.CoverageAneka) {
		for s := range c.Spreading {
			b := &c.Spreading[s]
			b.TSISpreaded = persenDari(b.SharePercentage, c.TSINusantaraRe)
			b.PremiumSpreaded = persenDari(b.SharePercentage, c.PremiNusantaraRe)
		}
	})
}

// salinTemplateAneka - CopyToAllSpreadingAnekaGolf_ACT: setiap coverage (langkah 2 objek okupasi bila IsAneka / Bonding, langkah 3
// objek Golf bila isGolfInsurance) mendapat template; TSISpreaded = share × (TSI × PS / 100) / 100, PremiumSpreaded = share ×
// PremiNusantaraRe / 100; TSIGrossSpreaded / ClaimEstimation kosong.
func salinTemplateAneka(lokasi []models.LokasiAneka, ps *apd.Decimal, templat []models.TemplateSpreading, l liniSpreadingAneka) {
	setiapCoverageAneka(lokasi, l.aneka, l.golf, func(_ int, _ *models.ObjekAneka, _ int, c *models.CoverageAneka) {
		tsi := persenDari(c.TSI, ps)
		c.Spreading = make([]models.BarisSpreading, 0, len(templat))
		for _, t := range templat {
			c.Spreading = append(c.Spreading, models.BarisSpreading{TreatyType: t.TreatyType, TreatyName: t.TreatyName,
				SharePercentage: t.SharePercentage, TSISpreaded: persenDari(t.SharePercentage, tsi),
				PremiumSpreaded: persenDari(t.SharePercentage, c.PremiNusantaraRe)})
		}
	})
}

// spreadingOtomatisGolf - GetKapasitasTreaty cabang isGolfInsurance (DDL L10367-11168): per coverage objek Golf satu baris
// TreatyName "ORS", TreatyType "10007", share 100; TSISpreaded = share × (TSILiability × PS / 100) / 100 (TSILiability Golf
// kosong -> 0 sampai Save, K73-20), PremiumSpreaded = share × PremiNusantaraRe / 100, TSIGrossSpreaded = share × (TSI × PS / 100)
// / 100. Tanpa kapasitas, kurs, pesan.
func spreadingOtomatisGolf(lokasi []models.LokasiAneka, ps *apd.Decimal) {
	setiapCoverageAneka(lokasi, false, true, func(_ int, _ *models.ObjekAneka, _ int, c *models.CoverageAneka) {
		b := models.BarisSpreading{TreatyName: namaORS, TreatyType: treatyORS}
		isiBaris(&b, apd.New(100, 0), persenDari(c.TSILiability, ps), c.PremiNusantaraRe, persenDari(c.TSI, ps))
		c.Spreading = []models.BarisSpreading{b}
	})
}

// spreadingOtomatisEngineering - GetKapasitasTreaty cabang IsEngineering (DDL L12531-21510) atas objek okupasi; kerangka sama
// dengan spreadingOtomatis (FIRE). Beda verbatim `[terverifikasi]` yang dipertahankan:
//  1. template SpreadingList TIDAK dibuang per lokasi (FIRE 6.2.1-6.2.3 tidak ada) - satu template untuk seluruh case;
//  4. QS dipotong (.19.3): pembagi share = TSISpreaded sebelum dipotong (CekLimit.CARID1, L19128);
//  5. baris SPL tanpa top risk (.19.4) hanya bergerbang "QS dipotong" (L19795), tanpa TSIRNM > maxtsi.
//
// Perbaikan keputusan work owner 09-10-2026 (lewat 9d; docs/issues/73):
//   - K73-23 (A): fase A mengisi TSI RNM top risk persis cabang IsFire - tsitoprisk += TSILiability (× kurs bila bukan IDR, L2441)
//     dan TSIRNMTopRisk1 += TSILiability (× kurs) × PS / 100 (L2494), coverage pertama lokasi top risk (L2595); verbatim Engineering
//     L13793 tanpa kurs dan TSIRNMTopRisk1 yang tak pernah diisi (L16245 membaca 0) DIGANTI; L16245 sendiri tetap;
//   - K73-25 (P1): pembilang share QS dipotong = Local.maxtsiQS (L16086 / L16384, kurs yang sama), bukan MaxLimitQSIDR mentah
//     (L19216 / L19237) - hasil IDR tidak berubah;
//   - K73-26 (c): pemicu pemotongan .19.3 (`.TSISpreaded > Local.maxtsiQS`, L19342) membandingkan KEDUA nilai sesudah dibulatkan
//     8 desimal (sisa skala 20 pecahan tak berujung tidak lagi memicu pemotongan);
//   - K73-27 (ii, MENYIMPANG dari Pega): pemicu .19.3 hanya untuk baris berjenis treaty QS (TreatyType = IDTreatyQS). WHEN L19342
//     Pega tanpa saringan jenis, sehingga baris SPL dibandingkan dengan limit QS lalu dipotong dan "beda 5" menambah SPL lagi.
//     Pemicu top risk .20.3 tetap verbatim.
//
// Kurs, KAPASITAS_TREATY tanpa baris, dan pesan kapasitas sama dengan FIRE (K48-3 / K48-4).
func spreadingOtomatisEngineering(ctx context.Context, src sumberKapasitas, k models.KasusSpreading, lokasi []models.LokasiAneka) ([]string, error) {
	ps := nolBila(k.PercentShare)
	tanggal := tanggalMulai(k)
	cache := map[string]*apd.Decimal{}
	var gagal []string
	setiapCoverageAneka(lokasi, true, false, func(_ int, o *models.ObjekAneka, m int, _ *models.CoverageAneka) {
		if m == 0 {
			gagal = append(gagal, o.Currency)
		}
	})
	for _, u := range gagal {
		_, ok, err := kursKeIDR(ctx, src, cache, u)
		if err != nil {
			return nil, err
		}
		if !ok {
			return []string{"Kurs " + u + " ke IDR tidak ada di TREATYEXCHANGEYEARLY - spreading otomatis tidak dibuat"}, nil
		}
	}
	// Fase A (9.1): TSI RNM 100 % (IDR) dan TSI top risk dari coverage pertama tiap objek; spreading lama dibuang.
	tsiRNM100, tsiTopRisk, tsiRNMTopRisk1 := nol(), nol(), nol()
	var errA error
	setiapCoverageAneka(lokasi, true, false, func(i int, o *models.ObjekAneka, m int, c *models.CoverageAneka) {
		if m == 0 && errA == nil {
			kurs, _, err := kursKeIDR(ctx, src, cache, o.Currency)
			if err != nil {
				errA = err
				return
			}
			dasar := nolBila(c.TSILiability)
			if o.Currency != namaIDR {
				dasar = kali(dasar, kurs)
			}
			tsiRNM100 = tambah(tsiRNM100, persenDari(dasar, ps))
			if lokasi[i].IsTopRisk { // K73-23: L2441 / L2494
				tsiTopRisk = tambah(tsiTopRisk, dasar)
				tsiRNMTopRisk1 = tambah(tsiRNMTopRisk1, persenDari(dasar, ps))
			}
		}
		c.Spreading = []models.BarisSpreading{}
	})
	if errA != nil {
		return nil, errA
	}
	// Fase B (9.2).
	var templat []models.BarisSpreading
	tsiRNMTopRisk := nol()
	for i := range lokasi {
		top := lokasi[i].IsTopRisk
		for j := range lokasi[i].Okupasi {
			for n := range lokasi[i].Okupasi[j].Objek {
				o := &lokasi[i].Okupasi[j].Objek[n]
				kurs, _, err := kursKeIDR(ctx, src, cache, o.Currency)
				if err != nil {
					return nil, err
				}
				for m := range o.Coverages {
					c := &o.Coverages[m]
					tsi := nolBila(c.TSILiability)
					if top {
						tsi = tsiTopRisk
					}
					tsiRNM, tsiGross, premi := persenDari(c.TSILiability, ps), persenDari(c.TSI, ps), nolBila(c.PremiNusantaraRe)
					kap, ada, err := src.KapasitasTreaty(ctx, simpan8(tsi), tanggal)
					if err != nil {
						return nil, err
					}
					if !ada {
						return []string{"Kapasitas treaty (KAPASITAS_TREATY) untuk TSI " + teksAngka(tsi) + " pada " + tanggal +
							" tidak ada - spreading otomatis berhenti"}, nil
					}
					maxQS, maxSPL := nolBila(kap.MaxLimitQSIDR), nolBila(kap.MaxLimitSPLIDR)
					if top { // .7 L16245: TSIRNMTopRisk = TSIRNMTopRisk1 (IDR) / kurs
						tsiRNMTopRisk = tsiRNMTopRisk1
						if o.Currency != namaIDR {
							tsiRNMTopRisk = bagiSkala(tsiRNMTopRisk1, kurs, 20)
						}
					}
					if o.Currency != namaIDR {
						maxQS, maxSPL = bagiSkala(maxQS, kurs, 20), bagiSkala(maxSPL, kurs, 20)
						tsiRNM100 = bagiSkala(tsiRNM100, kurs, 20) // .8 L16450 - dibagi ulang tiap coverage, seperti FIRE
					}
					if len(templat) == 0 {
						templat = []models.BarisSpreading{{TreatyName: kap.TreatyNameQS, TreatyType: kap.IDTreatyQS, SharePercentage: apd.New(100, 0)}}
					}
					c.Spreading = salinBaris(templat)
					maxTSI := tambah(maxQS, maxSPL)
					if top && lebihBesar(tsiRNMTopRisk, maxTSI) {
						return []string{pesanKapasitasTop}, nil
					}
					pecah := func(dasar *apd.Decimal) {
						pctQS := kali(bagiSkala(maxQS, dasar, 20), seratus)
						pctSPL := bagiSkala(kurang(seratus, pctQS), apd.New(1, 0), 20)
						if lebihBesar(dasar, maxTSI) {
							pctSPL = kali(bagiSkala(maxSPL, dasar, 20), seratus)
						}
						templat = []models.BarisSpreading{
							{TreatyName: kap.TreatyNameQS, TreatyType: kap.IDTreatyQS, SharePercentage: pctQS, TSISpreaded: maxQS},
							{TreatyName: kap.TreatyNameSPL, TreatyType: kap.IDTreatySPL, SharePercentage: pctSPL},
						}
						c.Spreading = salinBaris(templat)
					}
					if top && lebihBesar(tsiRNMTopRisk, maxQS) {
						pecah(tsiRNMTopRisk)
					}
					if !top && lebihBesar(tsiRNM100, maxQS) {
						pecah(tsiRNM100)
					}
					dipotong := false
					var pctQS *apd.Decimal
					banyak := len(c.Spreading)
					for s := 0; s < banyak; s++ {
						b := &c.Spreading[s]
						isiBaris(b, b.SharePercentage, tsiRNM, premi, tsiGross)
						potong := lebihBesar(b.TSISpreaded, maxQS) // .20.3 L20772 (top risk): verbatim
						if !top {
							potong = b.TreatyType == kap.IDTreatyQS && lebihBesar(simpan8(b.TSISpreaded), simpan8(maxQS)) // .19.3 L19342, K73-26 / K73-27
						}
						if potong {
							dipotong = true
							pembagi := b.TSISpreaded // beda 4 (CekLimit.CARID1)
							if top {
								pembagi = tsiRNMTopRisk // .20.3 L20646: maxtsiQS / TSIRNMTopRisk
							}
							pctQS = kali(bagiSkala(maxQS, pembagi, 20), seratus) // K73-25: maxtsiQS
							isiBaris(b, pctQS, tsiRNM, premi, tsiGross)
							b.TSISpreaded = maxQS
						}
						tambahSPL := dipotong // beda 5
						if top {
							tambahSPL = dipotong && !lebihBesar(tsiRNMTopRisk, maxQS)
						}
						if tambahSPL {
							pembagi := tsiRNM
							if top {
								pembagi = tsiRNMTopRisk
							}
							c.Spreading = append(c.Spreading, barisSPL(kap, pctQS, tsiRNM, pembagi, maxSPL, premi, tsiGross, top, tsiRNMTopRisk))
						}
						if !top && lebihBesar(tsiRNM, maxTSI) {
							return []string{"TSI in Location " + strconv.Itoa(i+1) + " exceeds treaty capacity"}, nil
						}
					}
				}
			}
		}
	}
	return nil, nil
}

// teksAngka - nilai 8 desimal sebagai teks pesan (sama dengan pesan kapasitas FIRE).
func teksAngka(d *apd.Decimal) string {
	return utils.FormatDecimal(simpan8(d))
}

// kunciMataUang - mata uang total satu coverage: lini daftar -> mata uang coverage (kosong -> objek, `[keputusan agent]`, Count
// L2969 mengisinya dari objek untuk delapan lini); selain itu mata uang objek.
func kunciMataUang(o *models.ObjekAneka, c *models.CoverageAneka, perCoverage bool) string {
	if perCoverage && strings.TrimSpace(c.Currency) != "" {
		return c.Currency
	}
	return o.Currency
}

// totalLokasiAneka - SumTSIPremiSpreadedRNM_ANEKA_Act (objek okupasi). Tanpa top risk di case mana pun -> cabang 2.4, selain itu
// 2.5. Per baris spreading:
//   - premi selalu dijumlah; TSI dijumlah hanya coverage pertama (L6374 / L12957), tetapi baris BARU memakai TSISpreaded baris
//     berjalan (L7519 / L14171, verbatim);
//   - LoL = LimitOfLiability × PS / 100 coverage pertama ber-LoL > 0 (L5116); LoL spread = LoL × share / 100: 2.4 untuk semua
//     baris dan DITULIS ke ClaimEstimation baris (L5701), 2.5 hanya lokasi top risk dan tidak ditulis (L12266);
//   - TSI top risk (ClaimSpreaded) hanya 2.5, coverage pertama, lokasi top risk, lini liniTopRisk25 (L12408-12592);
//   - share beda -> kosong (L6719 / L13371).
//
// Mengembalikan total per lokasi dan "ada top risk".
func totalLokasiAneka(lokasi []models.LokasiAneka, ps *apd.Decimal, l liniSpreadingAneka) ([][]models.TotalSpreading, bool) {
	adaTop := false
	for _, x := range lokasi {
		adaTop = adaTop || x.IsTopRisk
	}
	perCoverage, liniTop := l.mataUangCov24, false
	if adaTop {
		perCoverage, liniTop = l.mataUangCov25, l.topRisk25
	}
	total := make([][]models.TotalSpreading, len(lokasi))
	for i := range total {
		total[i] = []models.TotalSpreading{}
	}
	if !l.aneka {
		return total, adaTop
	}
	setiapCoverageAneka(lokasi, true, false, func(i int, o *models.ObjekAneka, m int, c *models.CoverageAneka) {
		pertama, top := m == 0, lokasi[i].IsTopRisk
		mataUang := kunciMataUang(o, c, perCoverage)
		tsiLoL := nol()
		if pertama && lebihBesar(c.LimitOfLiability, nol()) {
			tsiLoL = persenDari(c.LimitOfLiability, ps)
		}
		for s := range c.Spreading {
			b := &c.Spreading[s]
			share := nolBila(b.SharePercentage)
			tsiObj, lolSpread, tsiTop := nol(), nol(), nol()
			if !adaTop {
				tsiObj = nolBila(b.TSISpreaded)
				lolSpread = persenDari(tsiLoL, share)
				b.ClaimEstimation = lolSpread
			} else {
				if pertama {
					tsiObj = nolBila(b.TSISpreaded)
				}
				if top {
					lolSpread = persenDari(tsiLoL, share)
				}
				if pertama && top && liniTop {
					tsiTop = nolBila(b.TSISpreaded)
				}
			}
			t := &total[i]
			if n := cariTotal(*t, b.TreatyType, mataUang); n >= 0 {
				x := &(*t)[n]
				x.PremiumSpreaded = tambah(x.PremiumSpreaded, b.PremiumSpreaded)
				if pertama {
					x.TSISpreaded = tambah(x.TSISpreaded, tsiObj)
					x.ClaimEstimation = tambah(x.ClaimEstimation, lolSpread)
				}
				if adaTop && top && liniTop {
					x.ClaimSpreaded = tambah(x.ClaimSpreaded, tsiTop)
				}
				if x.SharePercentage == nil || share.Cmp(x.SharePercentage) != 0 {
					x.SharePercentage = nil
				}
				continue
			}
			baru := models.TotalSpreading{Currency: mataUang, TreatyType: b.TreatyType, SharePercentage: b.SharePercentage,
				TSISpreaded: nolBila(b.TSISpreaded), PremiumSpreaded: nolBila(b.PremiumSpreaded), ClaimEstimation: lolSpread}
			if adaTop && top && liniTop {
				baru.ClaimSpreaded = tsiTop
			}
			*t = append(*t, baru)
		}
	})
	return total, adaTop
}

// totalLokasiGolf - SumTSIPremiSpreadedRNM_GOLF_Act 2.3 (objek Golf): TSI dan premi SETIAP coverage dijumlah (L3584), tanpa LoL /
// top risk; share beda -> kosong (L3740).
func totalLokasiGolf(lokasi []models.LokasiAneka, l liniSpreadingAneka) [][]models.TotalSpreading {
	total := make([][]models.TotalSpreading, len(lokasi))
	for i := range total {
		total[i] = []models.TotalSpreading{}
	}
	if !l.golf {
		return total
	}
	setiapCoverageAneka(lokasi, false, true, func(i int, o *models.ObjekAneka, _ int, c *models.CoverageAneka) {
		for _, b := range c.Spreading {
			t := &total[i]
			if n := cariTotal(*t, b.TreatyType, o.Currency); n >= 0 {
				x := &(*t)[n]
				x.TSISpreaded = tambah(x.TSISpreaded, b.TSISpreaded)
				x.PremiumSpreaded = tambah(x.PremiumSpreaded, b.PremiumSpreaded)
				if x.SharePercentage == nil || nolBila(b.SharePercentage).Cmp(x.SharePercentage) != 0 {
					x.SharePercentage = nil
				}
				continue
			}
			*t = append(*t, models.TotalSpreading{Currency: o.Currency, TreatyType: b.TreatyType, SharePercentage: b.SharePercentage,
				TSISpreaded: nolBila(b.TSISpreaded), PremiumSpreaded: nolBila(b.PremiumSpreaded)})
		}
	})
	return total
}

// totalAneka - total per lokasi (Aneka + Golf) sebagai KasusSpreading berlokasi (IsTopRisk, Total) agar totalMataUang FIRE dapat
// dipakai; "ada top risk" untuk totalMataUang.
func totalAneka(lokasi []models.LokasiAneka, k models.KasusSpreading, l liniSpreadingAneka) (models.KasusSpreading, bool) {
	aneka, adaTop := totalLokasiAneka(lokasi, k.PercentShare, l)
	golf := totalLokasiGolf(lokasi, l)
	hasil := models.KasusSpreading{PercentShare: k.PercentShare, StartDateTime: k.StartDateTime,
		Lokasi: make([]models.LokasiSpreading, len(lokasi))}
	for i := range lokasi {
		hasil.Lokasi[i] = models.LokasiSpreading{ObjectNo: lokasi[i].ObjectNo, Location: lokasi[i].RiskLocation,
			IsTopRisk: lokasi[i].IsTopRisk, Total: append(aneka[i], golf[i]...)}
	}
	return hasil, adaTop
}

// hapusPesanMBD - SumTSIPremiSpreadedRNM_Act langkah 24 "recheck tsi psreaddnya lagi" (WHEN IsMBD L10448 / L10396; K73-22 work
// owner 09-10-2026: port apa adanya) `[terverifikasi]`: untuk setiap baris GetTreatyName (24.1 TreatyType = .CARI1, limit = .CARI3
// L9002 / L9049; ORS 10007 = 150000000000.00 L9163, atau 181500000000.00 bila Begin date >= "20260630T170000.000 GMT" L9307),
// setiap total per lokasi berjenis sama dengan limit <= TSISpreaded (mata uang asli, tanpa kurs, L9849 / L9895 / L10069) ->
// Page-Clear-Messages pyWorkPage (L10014): SELURUH pesan yang sudah terkumpul dihapus. Di luar port: 24.4 kasus group / pyID
// hardcode (pola K48-6), ProtectSpreading.CARI1 = 0 (L10186, penanda yang tidak dimodelkan). Limit yang tidak terurai sebagai
// angka tidak dibandingkan `[dugaan]` (Pega membandingkan properti desimal kosong - perilakunya belum terverifikasi).
func hapusPesanMBD(k models.KasusSpreading, lokasi []models.LokasiSpreading, batas []models.TreatySpreading) bool {
	mulai, adaMulai := waktuMulai(k)
	for _, b := range batas {
		limit, err := utils.ParseDecimal(strings.TrimSpace(b.Limit))
		if b.ID == treatyORS {
			limit, err = batasORS150, nil
			if adaMulai && !mulai.Before(mulaiORS1815) {
				limit = batasORS1815
			}
		}
		if err != nil {
			continue
		}
		for _, l := range lokasi {
			for _, t := range l.Total {
				if t.TreatyType == b.ID && limit.Cmp(nolBila(t.TSISpreaded)) <= 0 {
					return true
				}
			}
		}
	}
	return false
}

// bulatkanSimpanAneka - nilai yang ditulis: 8 desimal (NUMBER(38,8)).
func bulatkanSimpanAneka(lokasi []models.LokasiAneka) {
	setiapCoverageAneka(lokasi, true, true, func(_ int, _ *models.ObjekAneka, _ int, c *models.CoverageAneka) {
		c.TSINusantaraRe, c.PremiNusantaraRe = simpan8(c.TSINusantaraRe), simpan8(c.PremiNusantaraRe)
		for s := range c.Spreading {
			b := &c.Spreading[s]
			b.SharePercentage, b.TSISpreaded, b.TSIGrossSpreaded = simpan8(b.SharePercentage), simpan8(b.TSISpreaded), simpan8(b.TSIGrossSpreaded)
			b.PremiumSpreaded, b.ClaimEstimation = simpan8(b.PremiumSpreaded), simpan8(b.ClaimEstimation)
		}
	})
}

// hitungSaveSpreadingAneka - hitungan Save (IsThereAnyObjectLocation_Act): % Share RNM, NR, DT percent, total (ClaimEstimation
// baris cabang 2.4), 8 desimal. Satu jalur untuk PUT, hitung-baris, dan Save tab Object.
func hitungSaveSpreadingAneka(lokasi []models.LokasiAneka, ps *apd.Decimal, l liniSpreadingAneka) {
	hitungNRAneka(lokasi, ps)
	terapkanPersenAneka(lokasi, l)
	totalLokasiAneka(lokasi, ps, l)
	bulatkanSimpanAneka(lokasi)
}

// ---- adaptor hilir (Payment / Premium Deduction / Cedant) -------------------------------------------------------------

// kasusSpreadingAneka - pohon Aneka sebagai KasusSpreading untuk ringkasanRNM (premi / TSI per mata uang): sumber Pega
// CountPremiAndTSIRNMFireAnekaGolf_ACT 7.2.5 (L5211-5277) / 8.2.5 (L7512-7578). Satu item per kelompok:
//   - Aneka lini liniMataUangCount: satu item per coverage (mata uang coverage, kosong -> objek); TSI = TSI coverage untuk
//     liniTSICoverageCount, selain itu kosong (BuilderRisk / PortRisk tidak menyumbang TSI - verbatim L2556 / L4324);
//   - Aneka lini lain dan Golf: satu item per objek (mata uang objek, TSI objek).
//
// TSILiability coverage pertama item = TSI item (dasar TSIRNM ringkasanRNM = PS × TSI / 100, L4841 / L7140). Penyimpangan sadar
// `[keputusan agent]`: Local.totaltsi Pega TIDAK direset per mata uang (kumulatif) - di sini per mata uang.
func kasusSpreadingAneka(k models.KasusSpreading, lokasi []models.LokasiAneka, l liniSpreadingAneka) models.KasusSpreading {
	hasil := models.KasusSpreading{PercentShare: k.PercentShare, StartDateTime: k.StartDateTime}
	item := func(mataUang string, tsi *apd.Decimal, cov []models.CoverageAneka) models.ItemSpreading {
		it := models.ItemSpreading{Currency: mataUang, TSI: tsi, Coverages: make([]models.CoverageSpreading, 0, len(cov))}
		for m, c := range cov {
			cs := models.CoverageSpreading{OldID: c.OldID, Premium: c.Premium, Discount: c.Discount, TSI: c.TSI}
			if m == 0 {
				cs.TSILiability = tsi
			}
			it.Coverages = append(it.Coverages, cs)
		}
		return it
	}
	for _, x := range lokasi {
		ls := models.LokasiSpreading{ObjectNo: x.ObjectNo, Location: x.RiskLocation, IsTopRisk: x.IsTopRisk}
		for _, ok := range x.Okupasi {
			if !l.aneka { // langkah 7 WHEN IsAneka (L5462)
				break
			}
			for _, o := range ok.Objek {
				if !l.mataUangCount {
					ls.Items = append(ls.Items, item(o.Currency, o.TSI, o.Coverages))
					continue
				}
				for _, c := range o.Coverages {
					var tsi *apd.Decimal
					if l.tsiCoverageCount {
						tsi = c.TSI
					}
					o2 := o
					ls.Items = append(ls.Items, item(kunciMataUang(&o2, &c, true), tsi, []models.CoverageAneka{c}))
				}
			}
		}
		for _, o := range x.ObjekGolf {
			if !l.golf { // langkah 8 WHEN isGolfInsurance (L7763)
				break
			}
			ls.Items = append(ls.Items, item(o.Currency, o.TSI, o.Coverages))
		}
		hasil.Lokasi = append(hasil.Lokasi, ls)
	}
	return hasil
}

// adaItemFire - pohon FIRE berisi sekurangnya satu PropertyItem.
func adaItemFire(k models.KasusSpreading) bool {
	for _, l := range k.Lokasi {
		if len(l.Items) > 0 {
			return true
		}
	}
	return false
}

// bacaSpreadingRingkas - sumber premi / TSI per mata uang tab Payment / Premium Deduction / Cedant: pohon FIRE apa adanya
// (perilaku tiket 48-51 TIDAK berubah) bila memuat item; bila tidak, pohon Aneka (tiket 73 P3) bila ada objek. Pohon Aneka
// dibaca di luar tx (penyimpan objek Aneka tidak menerima tx).
func (s *Service) bacaSpreadingRingkas(ctx context.Context, tx *db.Tx, id string) (models.KasusSpreading, error) {
	k, err := s.spreading.BacaSpreading(ctx, tx, id)
	if err != nil || s.objekAneka == nil || adaItemFire(k) {
		return k, err
	}
	lokasi, err := s.objekAneka.BacaObjekAneka(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return k, nil
	}
	if err != nil {
		return k, err
	}
	if !adaObjekAneka(lokasi) {
		return k, nil
	}
	kb, _ := s.kasusBisnisDari(ctx, id)
	return kasusSpreadingAneka(k, lokasi, liniSpreadingDari(id, kb)), nil
}

// ---- alur ---------------------------------------------------------------------------------------------------------

// dataSpreadingAneka - % Share RNM + Begin date (pohon FIRE, tanpa item) dan pohon Aneka tersimpan (salinan dalam) + lini.
type dataSpreadingAneka struct {
	kasus  models.KasusSpreading
	lokasi []models.LokasiAneka
	lini   liniSpreadingAneka
	mbdUW  bool
}

func (s *Service) bacaDataSpreadingAneka(ctx context.Context, id string) (dataSpreadingAneka, error) {
	if s.spreading == nil || s.objekAneka == nil {
		return dataSpreadingAneka{}, ErrSpreadingTanpaDatabase
	}
	k, err := s.bacaSpreading(ctx, nil, id)
	if err != nil {
		return dataSpreadingAneka{}, err
	}
	lokasi, err := s.objekAneka.BacaObjekAneka(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return dataSpreadingAneka{}, ErrKasusTidakAda
	}
	if err != nil {
		return dataSpreadingAneka{}, err
	}
	kb, _ := s.kasusBisnisDari(ctx, id)
	l := liniSpreadingDari(id, kb)
	t, ada, err := s.tahapKasus(ctx, id)
	if err != nil {
		return dataSpreadingAneka{}, err
	}
	return dataSpreadingAneka{kasus: models.KasusSpreading{PercentShare: k.PercentShare, StartDateTime: k.StartDateTime},
		lokasi: salinDalamAneka(lokasi), lini: l, mbdUW: l.mbd && (!ada || tahapUW(t))}, nil
}

// tampilanAneka - total, ringkasan, pesan (pola tampilan FIRE); periksaJenis = cekSpreadingFactIn langkah 8 (Copy To All /
// hitung-baris).
func (s *Service) tampilanSpreadingAneka(ctx context.Context, id string, d dataSpreadingAneka, pesan []string, periksaJenis bool) (TampilanSpreadingAneka, error) {
	lokasi := salinDalamAneka(d.lokasi)
	dropdown, batas, err := s.treatyKasus(ctx, id, d.kasus)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	tot, adaTop := totalAneka(lokasi, d.kasus, d.lini)
	perMataUang := totalMataUang(tot, adaTop)
	treaty, semua, pesanBatas, err := ringkasan(ctx, s.spreading, d.kasus, perMataUang, dropdown, batas)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	pesan = append(append([]string{}, pesan...), pesanBatas...)
	if d.lini.mbd && hapusPesanMBD(d.kasus, tot.Lokasi, batas) {
		pesan = []string{}
	}
	if periksaJenis {
		pesan = append(pesan, periksaJenisTreaty(perMataUang, batas)...)
	}
	bulatkanSimpanAneka(lokasi)
	hasil := TampilanSpreadingAneka{PercentShare: d.kasus.PercentShare, Treaty: dropdown, Template: []models.TemplateSpreading{},
		Lokasi: make([]LokasiSpreadingAneka, len(lokasi)), RingkasanTreaty: bulatkanTotal(treaty), RingkasanMataUang: semua,
		Pesan: append([]string{}, pesan...), TampilTotalLokasi: d.mbdUW}
	for i, x := range lokasi {
		total := tot.Lokasi[i].Total
		for j := range total {
			total[j].TreatyName = namaTreaty(dropdown, total[j].TreatyType)
		}
		hasil.Lokasi[i] = LokasiSpreadingAneka{ObjectNo: x.ObjectNo, Location: x.RiskLocation, Okupasi: x.Okupasi, ObjekGolf: x.ObjekGolf,
			Total: bulatkanTotal(total)}
	}
	for i := range semua {
		r := &semua[i]
		r.TSITopRisk, r.TSI, r.LoL, r.Premium = simpan8(r.TSITopRisk), simpan8(r.TSI), simpan8(r.LoL), simpan8(r.Premium)
	}
	return hasil, nil
}

// TampilanSpreadingAneka - GET /api/nbfacin/kasus/{caseId}/spreading-aneka.
func (s *Service) TampilanSpreadingAneka(ctx context.Context, id string) (TampilanSpreadingAneka, error) {
	d, err := s.bacaDataSpreadingAneka(ctx, id)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	return s.tampilanSpreadingAneka(ctx, id, d, nil, false)
}

// HitungShareSpreadingAneka - POST …/spreading-aneka/hitung-share: NR + spreading; TIDAK menyimpan. Spreading otomatis hanya tahap
// uw (pola tiket 59) dan hanya Golf / Engineering (K73-17; GetKapasitasTreaty tanpa cabang Aneka / Bonding): baris lain
// dipertahankan dan dihitung ulang (DT percent).
func (s *Service) HitungShareSpreadingAneka(ctx context.Context, id string, ps *apd.Decimal) (TampilanSpreadingAneka, error) {
	if err := periksaShareKasus(ps); err != nil {
		return TampilanSpreadingAneka{}, err
	}
	t, ada, err := s.tahapKasus(ctx, id)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	d, err := s.bacaDataSpreadingAneka(ctx, id)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	d.kasus.PercentShare = ps
	hitungNRAneka(d.lokasi, ps)
	var pesan []string
	switch uw := !ada || tahapUW(t); {
	case uw && d.lini.golf:
		spreadingOtomatisGolf(d.lokasi, ps)
	case uw && d.lini.engineering:
		if pesan, err = spreadingOtomatisEngineering(ctx, s.spreading, d.kasus, d.lokasi); err != nil {
			return TampilanSpreadingAneka{}, err
		}
	default:
		terapkanPersenAneka(d.lokasi, d.lini)
	}
	return s.tampilanSpreadingAneka(ctx, id, d, pesan, false)
}

// SalinSpreadingAneka - POST …/spreading-aneka/salin: Copy To All Spreading + total; TIDAK menyimpan. Tahap uw + workbasket
// posisi (sama dengan FIRE, tiket 59); template diperiksa sama dengan FIRE (400).
func (s *Service) SalinSpreadingAneka(ctx context.Context, pelaku inti.Pelaku, id string, ps *apd.Decimal, templat []models.TemplateSpreading) (TampilanSpreadingAneka, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanSpreadingAneka{}, err
	}
	if err := s.wajibTahapUW(ctx, pelaku, id); err != nil {
		return TampilanSpreadingAneka{}, err
	}
	if err := periksaShareKasus(ps); err != nil {
		return TampilanSpreadingAneka{}, err
	}
	if len(templat) == 0 {
		return TampilanSpreadingAneka{}, fmt.Errorf("%w: template wajib berisi paling sedikit satu baris", ErrMasukanSpreading)
	}
	jenis, nama, share := make([]string, len(templat)), make([]string, len(templat)), make([]*apd.Decimal, len(templat))
	for i, t := range templat {
		jenis[i], nama[i], share[i] = t.TreatyType, t.TreatyName, t.SharePercentage
	}
	if m := periksaBarisShare("template", jenis, nama, share); len(m) > 0 {
		return TampilanSpreadingAneka{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(m, "; "))
	}
	d, err := s.bacaDataSpreadingAneka(ctx, id)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	m, err := s.proteksiTemplate(ctx, templat)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	if len(m) > 0 {
		return TampilanSpreadingAneka{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(m, "; "))
	}
	d.kasus.PercentShare = ps
	hitungNRAneka(d.lokasi, ps)
	salinTemplateAneka(d.lokasi, ps, templat, d.lini)
	t, err := s.tampilanSpreadingAneka(ctx, id, d, nil, true)
	if err != nil {
		return t, err
	}
	t.Template = templat
	return t, nil
}

// samaBentukAneka - jumlah lokasi / okupasi / objek / objek Golf / coverage dan oldId coverage sama menurut posisi.
func samaBentukAneka(simpan, badan []models.LokasiAneka) bool {
	samaObjek := func(a, b []models.ObjekAneka) bool {
		if len(a) != len(b) {
			return false
		}
		for n := range a {
			if len(a[n].Coverages) != len(b[n].Coverages) {
				return false
			}
			for m := range a[n].Coverages {
				if a[n].Coverages[m].OldID != b[n].Coverages[m].OldID {
					return false
				}
			}
		}
		return true
	}
	if len(simpan) != len(badan) {
		return false
	}
	for i := range simpan {
		if len(simpan[i].Okupasi) != len(badan[i].Okupasi) || !samaObjek(simpan[i].ObjekGolf, badan[i].ObjekGolf) {
			return false
		}
		for j := range simpan[i].Okupasi {
			if !samaObjek(simpan[i].Okupasi[j].Objek, badan[i].Okupasi[j].Objek) {
				return false
			}
		}
	}
	return true
}

// bacaBentukBadanAneka - % Share RNM 0..100, data tersimpan, objek wajib ada ("Object can't be empty!"), bentuk badan = tersimpan
// (409).
func (s *Service) bacaBentukBadanAneka(ctx context.Context, id string, ps *apd.Decimal, lokasi []models.LokasiAneka) (dataSpreadingAneka, error) {
	if err := periksaShareKasus(ps); err != nil {
		return dataSpreadingAneka{}, err
	}
	d, err := s.bacaDataSpreadingAneka(ctx, id)
	if err != nil {
		return d, err
	}
	if !adaObjekAneka(d.lokasi) {
		return d, fmt.Errorf("%w: %s", ErrMasukanSpreading, pesanObjekKosong)
	}
	if !samaBentukAneka(d.lokasi, lokasi) {
		return d, ErrSpreadingBerubah
	}
	return d, nil
}

// pasangBarisBadanAneka - treatyType / treatyName / sharePercentage baris badan menggantikan baris tersimpan; setiap daftar baris
// coverage diperiksa periksaBarisShare (400).
func pasangBarisBadanAneka(simpan, badan []models.LokasiAneka) error {
	var masalah []string
	pasang := func(jalur string, a, b []models.ObjekAneka) {
		for n := range a {
			for m := range a[n].Coverages {
				baris := b[n].Coverages[m].Spreading
				jenis, nama, share := make([]string, len(baris)), make([]string, len(baris)), make([]*apd.Decimal, len(baris))
				salin := make([]models.BarisSpreading, len(baris))
				for x, r := range baris {
					jenis[x], nama[x], share[x] = r.TreatyType, r.TreatyName, r.SharePercentage
					salin[x] = models.BarisSpreading{TreatyType: strings.TrimSpace(r.TreatyType), TreatyName: strings.TrimSpace(r.TreatyName),
						SharePercentage: r.SharePercentage}
				}
				masalah = append(masalah, periksaBarisShare(fmt.Sprintf("%s[%d].coverages[%d].spreading", jalur, n, m), jenis, nama, share)...)
				a[n].Coverages[m].Spreading = salin
			}
		}
	}
	for i := range simpan {
		for j := range simpan[i].Okupasi {
			pasang(fmt.Sprintf("lokasi[%d].okupasi[%d].objek", i, j), simpan[i].Okupasi[j].Objek, badan[i].Okupasi[j].Objek)
		}
		pasang(fmt.Sprintf("lokasi[%d].objekGolf", i), simpan[i].ObjekGolf, badan[i].ObjekGolf)
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(masalah, "; "))
	}
	return nil
}

// SimpanSpreadingAneka - PUT …/spreading-aneka (Save tab): penjaga tahap = PUT spreading FIRE (penjagaPutSpreading: baris badan
// hanya di uw; di marketing / Team Leader baris tersimpan dihitung ulang); bentuk 409; hitungSaveSpreadingAneka; di uw jenis
// treaty tidak sah 400 (K59-8) dan langkah 11 (Premium Deduction / Payment, K59-9 / K59-10) seperti FIRE; tulis % Share RNM +
// pohon objek Aneka (coverage beserta spreading) dalam satu transaksi; jawab baca ulang.
func (s *Service) SimpanSpreadingAneka(ctx context.Context, pelaku inti.Pelaku, id string, ps *apd.Decimal, lokasi []models.LokasiAneka) (TampilanSpreadingAneka, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanSpreadingAneka{}, err
	}
	barisBadan, uw, err := s.penjagaPutSpreading(ctx, pelaku, id)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	d, err := s.bacaBentukBadanAneka(ctx, id, ps, lokasi)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	if barisBadan {
		if err := pasangBarisBadanAneka(d.lokasi, lokasi); err != nil {
			return TampilanSpreadingAneka{}, err
		}
	}
	d.kasus.PercentShare = ps
	hitungSaveSpreadingAneka(d.lokasi, ps, d.lini)
	if s.transaksi == nil {
		return TampilanSpreadingAneka{}, ErrSpreadingTanpaDatabase
	}
	var l11 *dasarLangkah11
	var baris11 []models.BarisDeduction
	if uw {
		_, batas, err := s.treatyKasus(ctx, id, d.kasus)
		if err != nil {
			return TampilanSpreadingAneka{}, err
		}
		tot, adaTop := totalAneka(d.lokasi, d.kasus, d.lini)
		if pesan := periksaJenisTreaty(totalMataUang(tot, adaTop), batas); len(pesan) > 0 {
			return TampilanSpreadingAneka{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, pesan[0])
		}
		if l11, err = s.bacaDasarLangkah11(ctx, id); err != nil {
			return TampilanSpreadingAneka{}, err
		}
		baris11 = barisLangkah11(kasusSpreadingAneka(d.kasus, d.lokasi, d.lini), l11)
		if pesan := pesanDeduction2(baris11); len(pesan) > 0 {
			return TampilanSpreadingAneka{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, pesan[0])
		}
		var negatif []string
		for _, b := range baris11 {
			if netPremium(b).Sign() < 0 {
				negatif = append(negatif, pesanNetNegatif(b.Currency))
			}
		}
		if len(negatif) > 0 {
			return TampilanSpreadingAneka{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(negatif, "; "))
		}
	}
	proRate, err := s.objekAneka.ProRateKasus(ctx, id)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		if err := s.objekAneka.GantiObjekAneka(ctx, tx, id, d.lokasi, proRate); err != nil {
			return err
		}
		if err := s.objekAneka.UbahPercentShare(ctx, tx, id, simpan8(ps)); err != nil {
			return err
		}
		return s.tulisLangkah11UW(ctx, tx, id, l11, baris11)
	})
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return TampilanSpreadingAneka{}, ErrKasusTidakAda
	}
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	return s.TampilanSpreadingAneka(ctx, id)
}

// HitungBarisSpreadingAneka - POST …/spreading-aneka/hitung-baris: hitungan Save tahap uw tanpa menulis + cekSpreadingFactIn
// langkah 8-9 (pola HitungBarisSpreading FIRE). Tahap uw + workbasket posisi.
func (s *Service) HitungBarisSpreadingAneka(ctx context.Context, pelaku inti.Pelaku, id string, ps *apd.Decimal, lokasi []models.LokasiAneka) (TampilanSpreadingAneka, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanSpreadingAneka{}, err
	}
	if err := s.wajibTahapUW(ctx, pelaku, id); err != nil {
		return TampilanSpreadingAneka{}, err
	}
	d, err := s.bacaBentukBadanAneka(ctx, id, ps, lokasi)
	if err != nil {
		return TampilanSpreadingAneka{}, err
	}
	if err := pasangBarisBadanAneka(d.lokasi, lokasi); err != nil {
		return TampilanSpreadingAneka{}, err
	}
	d.kasus.PercentShare = ps
	hitungSaveSpreadingAneka(d.lokasi, ps, d.lini)
	return s.tampilanSpreadingAneka(ctx, id, d, nil, true)
}
