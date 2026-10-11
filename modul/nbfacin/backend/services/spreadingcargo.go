package services

// Tab Spreading MARINE CARGO (tiket 74 P3). Hierarki Cargo -> Coverage -> SpreadingList (tanpa lokasi); daun = FIRE. Bukti
// `[terverifikasi]` (`NB\` = `D:\migrasi\RNM\NB FacIn\`, riset sesi c3 10-10-2026):
//   - NR coverage = CountPremiAndTSIRNMCargo_ACT langkah 7 (L2425 / L2477): TSINusantaraRe = PS×TSI/100, PremiNusantaraRe =
//     PS×Premium/100 − PS×Discount/100 (skala 20);
//   - baris = DT `NB\DataTransform\SetTSIPremiSpreaded_FacIn.xml` langkah 12 (WHEN IsMarineCargo L6352) mode percent (L6564):
//     TSISpreaded = @Math.divide(share×TSINusantaraRe,100,20), PremiumSpreaded = @Math.divide(share×PremiNusantaraRe,100,20);
//     TSIGrossSpreaded / ClaimEstimation tidak di-SET (sel "Limit Of Liability" = isian baris);
//   - TANPA spreading otomatis: GetKapasitasTreaty tanpa cabang Marine; baris bawaan 10001 (SetSpreadingCoverage_DT langkah 8)
//     hanya lewat tombol tersembunyi - tidak diport;
//   - Copy To All = CopyToAllSpreadingCargoMBU_ACT langkah 2 (L488): template disalin ke setiap coverage (Page-Copy L2314), lalu
//     TSISpreaded = @Math.divide(share×(TSI×PS/100),100,20), PremiumSpreaded = @Math.divide(share×PremiNusantaraRe,100,20)
//     (L3099 / L3145); cabang EDM adj (2.2.3) tidak;
//   - ringkasan = SumTSIPremiSpreadedRNM_MARINECARGO_Act langkah 1 (L548-4447) lalu SumTSIPremiSpreadedRNM_Act langkah 17-24
//     bersama FIRE (ringkasan): per mata uang × jenis treaty ΣPremiumSpreaded, ΣTSISpreaded (cargo top risk dan bukan),
//     ClaimSpreaded = ΣTSISpreaded cargo top risk, ClaimEstimation = Σ (LoL coverage × PS / 100 HANYA cargo 1 - L1771) × share / 100;
//     share berbeda -> kosong (L3709); label CARI15 "Total Limit of Liability(" + PctLoL + "%)" (L4318, PctLoL coverage pertama
//     ber-LoL > 0: @divide(LoL×100, TSI>0?TSI:LoL, 4) L1791);
//   - cekSpreadingFactIn tanpa cabang Marine (generik: jenis treaty).
// K74-7 `[pertanyaan terbuka]`: IsThereAnyObjectLocation_Act langkah 6 men-SET Param.Type="percent" hanya bila LocationList >= 1 ||
// PersonList >= 1 (L1500) - untuk case cargo DT langkah 12 saat Save tidak menghitung ulang baris. Di sini Save menghitung ulang
// mode percent (sama dengan rantai on-change % Share / baris) `[keputusan agent]`, menunggu konfirmasi.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

// CargoSpreading - satu baris grid CargoList tab Spreading (`NB\Section\CoverageSpreadingList.xml` L15246: No., TradingNote,
// GoodNote, PackingNote, ConveyanceNote) dengan coverage-nya.
type CargoSpreading struct {
	No                                                 int
	TradingNote, GoodNote, PackingNote, ConveyanceNote string
	Coverages                                          []models.CoverageCargo
}

// TampilanSpreadingCargo - jawaban seluruh endpoint tab Spreading Marine Cargo. LabelLoL = ProtectSpreading.CARI15 (kosong bila
// tidak ada LoL).
type TampilanSpreadingCargo struct {
	PercentShare      *apd.Decimal
	Treaty            []models.TreatySpreading
	Template          []models.TemplateSpreading
	Cargo             []CargoSpreading
	RingkasanTreaty   []models.TotalSpreading
	RingkasanMataUang []models.RingkasanMataUangSpreading
	Pesan             []string
	LabelLoL          string
}

// setiapCoverageCargo - fn untuk setiap coverage cargo (urut cargo, coverage).
func setiapCoverageCargo(cargo []models.ObjekCargo, fn func(i int, c *models.CoverageCargo)) {
	for i := range cargo {
		for m := range cargo[i].Coverages {
			fn(i, &cargo[i].Coverages[m])
		}
	}
}

// hitungNRCargo - CountPremiAndTSIRNMCargo_ACT langkah 7 setiap coverage.
func hitungNRCargo(cargo []models.ObjekCargo, ps *apd.Decimal) {
	setiapCoverageCargo(cargo, func(_ int, c *models.CoverageCargo) { rnmCargo(c, nolBila(ps)) })
}

// terapkanPersenCargo - DT SetTSIPremiSpreaded_FacIn langkah 12 mode percent.
func terapkanPersenCargo(cargo []models.ObjekCargo) {
	setiapCoverageCargo(cargo, func(_ int, c *models.CoverageCargo) {
		for s := range c.Spreading {
			b := &c.Spreading[s]
			b.TSISpreaded = persenDari(b.SharePercentage, c.TSINusantaraRe)
			b.PremiumSpreaded = persenDari(b.SharePercentage, c.PremiNusantaraRe)
		}
	})
}

// salinTemplateCargo - CopyToAllSpreadingCargoMBU_ACT langkah 2: setiap coverage mendapat template (menggantikan baris lama, pola
// Copy To All Aneka); TSISpreaded = share × (TSI × PS / 100) / 100, PremiumSpreaded = share × PremiNusantaraRe / 100.
func salinTemplateCargo(cargo []models.ObjekCargo, ps *apd.Decimal, templat []models.TemplateSpreading) {
	setiapCoverageCargo(cargo, func(_ int, c *models.CoverageCargo) {
		tsi := persenDari(c.TSI, nolBila(ps))
		c.Spreading = make([]models.BarisSpreading, 0, len(templat))
		for _, t := range templat {
			c.Spreading = append(c.Spreading, models.BarisSpreading{TreatyType: t.TreatyType, TreatyName: t.TreatyName,
				SharePercentage: t.SharePercentage, TSISpreaded: persenDari(t.SharePercentage, tsi),
				PremiumSpreaded: persenDari(t.SharePercentage, c.PremiNusantaraRe)})
		}
	})
}

// totalCargo - SumTSIPremiSpreadedRNM_MARINECARGO_Act langkah 1: TotalSpreading per mata uang (urut kemunculan - Pega:
// OfferFacIn.CurrencyList `[dugaan]` sama) × jenis treaty, dan label CARI15. Cargo ber-FlagDelete "1" tidak dijumlah
// `[keputusan agent]` (pola CountPremiAndTSIRNMCargo_ACT 7.2.2.1.4).
func totalCargo(cargo []models.ObjekCargo, ps *apd.Decimal) ([]models.TotalSpreading, string) {
	var urut []string
	lihat := map[string]bool{}
	setiapCoverageCargo(cargo, func(_ int, c *models.CoverageCargo) {
		if !lihat[c.Currency] {
			lihat[c.Currency] = true
			urut = append(urut, c.Currency)
		}
	})
	hasil := []models.TotalSpreading{}
	var pctLoL *apd.Decimal
	for _, mu := range urut {
		awal := len(hasil)
		for i, x := range cargo {
			if strings.TrimSpace(x.FlagDelete) == "1" {
				continue
			}
			for _, c := range x.Coverages {
				if c.Currency != mu {
					continue
				}
				tsiLoL := nol()
				if i == 0 {
					tsiLoL = persenDari(c.LimitOfLiability, nolBila(ps))
				}
				if pctLoL == nil && lebihBesar(c.LimitOfLiability, nol()) {
					dasar := c.LimitOfLiability
					if lebihBesar(c.TSI, nol()) {
						dasar = c.TSI
					}
					if p, err := bagi(kali(c.LimitOfLiability, seratus), dasar, 4); err == nil {
						pctLoL = p
					}
				}
				for _, b := range c.Spreading {
					share := nolBila(b.SharePercentage)
					lolSpread := persenDari(tsiLoL, share)
					top := nol()
					if x.IsTopRisk {
						top = nolBila(b.TSISpreaded)
					}
					n := -1
					for k := awal; k < len(hasil); k++ {
						if hasil[k].TreatyType == b.TreatyType {
							n = k
							break
						}
					}
					if n < 0 {
						hasil = append(hasil, models.TotalSpreading{Currency: mu, TreatyType: b.TreatyType, SharePercentage: share,
							TSISpreaded: nolBila(b.TSISpreaded), PremiumSpreaded: nolBila(b.PremiumSpreaded), ClaimEstimation: lolSpread,
							ClaimSpreaded: top})
						continue
					}
					t := &hasil[n]
					t.PremiumSpreaded = tambah(t.PremiumSpreaded, b.PremiumSpreaded)
					t.ClaimEstimation = tambah(t.ClaimEstimation, lolSpread)
					t.TSISpreaded = tambah(t.TSISpreaded, b.TSISpreaded)
					t.ClaimSpreaded = tambah(t.ClaimSpreaded, top)
					if t.SharePercentage == nil || share.Cmp(t.SharePercentage) != 0 {
						t.SharePercentage = nil // L3709 `.SharePercentage=""`
					}
				}
			}
		}
	}
	label := ""
	if pctLoL != nil && pctLoL.Sign() > 0 {
		label = "Total Limit of Liability(" + utils.FormatDecimal(pctLoL) + "%)"
	}
	return hasil, label
}

// bulatkanSimpanCargo - nilai yang ditulis: 8 desimal (NUMBER(38,8)).
func bulatkanSimpanCargo(cargo []models.ObjekCargo) {
	setiapCoverageCargo(cargo, func(_ int, c *models.CoverageCargo) {
		bulatkanCoverageCargo(c)
		for s := range c.Spreading {
			b := &c.Spreading[s]
			b.SharePercentage, b.TSISpreaded, b.TSIGrossSpreaded = simpan8(b.SharePercentage), simpan8(b.TSISpreaded), simpan8(b.TSIGrossSpreaded)
			b.PremiumSpreaded, b.ClaimEstimation = simpan8(b.PremiumSpreaded), simpan8(b.ClaimEstimation)
		}
	})
}

// kasusSpreadingCargo - pohon cargo sebagai KasusSpreading untuk ringkasanRNM (Payment / Premium Deduction / Cedant, tiket 74 P4):
// CountPremiAndTSIRNMCargo_ACT langkah 7: per mata uang coverage, premi = Σ PremiNusantaraRe, TSI = Σ TSILiability (FlagDelete !=
// 1). Satu item per coverage; TSILiability coverage pertama item = TSI item. Penyimpangan sadar `[keputusan agent]` (sama dengan
// Aneka): Local.totaltsi Pega TIDAK direset per mata uang (7.2.1) - di sini per mata uang.
func kasusSpreadingCargo(k models.KasusSpreading, cargo []models.ObjekCargo) models.KasusSpreading {
	hasil := models.KasusSpreading{PercentShare: k.PercentShare, StartDateTime: k.StartDateTime}
	for i, x := range cargo {
		if strings.TrimSpace(x.FlagDelete) == "1" {
			continue
		}
		ls := models.LokasiSpreading{ObjectNo: fmt.Sprint(i + 1), IsTopRisk: x.IsTopRisk}
		for _, c := range x.Coverages {
			ls.Items = append(ls.Items, models.ItemSpreading{Currency: c.Currency, TSI: c.TSILiability,
				Coverages: []models.CoverageSpreading{{Premium: c.Premium, Discount: c.Discount, TSI: c.TSI, TSILiability: c.TSILiability}}})
		}
		hasil.Lokasi = append(hasil.Lokasi, ls)
	}
	return hasil
}

// adaCoverageCargo - sekurangnya satu coverage cargo.
func adaCoverageCargo(cargo []models.ObjekCargo) bool {
	for _, c := range cargo {
		if len(c.Coverages) > 0 {
			return true
		}
	}
	return false
}

// dataSpreadingCargo - % Share RNM + Begin date (pohon FIRE) dan pohon cargo tersimpan (salinan dalam).
type dataSpreadingCargo struct {
	kasus models.KasusSpreading
	cargo []models.ObjekCargo
}

func (s *Service) bacaDataSpreadingCargo(ctx context.Context, id string) (dataSpreadingCargo, error) {
	if s.spreading == nil || s.objekCargo == nil {
		return dataSpreadingCargo{}, ErrSpreadingTanpaDatabase
	}
	if !idKasusSah(id) {
		return dataSpreadingCargo{}, ErrKasusTidakAda
	}
	k, err := s.bacaSpreading(ctx, nil, id)
	if err != nil {
		return dataSpreadingCargo{}, err
	}
	cargo, err := s.objekCargo.BacaObjekCargo(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return dataSpreadingCargo{}, ErrKasusTidakAda
	}
	if err != nil {
		return dataSpreadingCargo{}, err
	}
	if !s.adalahMarineCargo(ctx, id) {
		return dataSpreadingCargo{}, ErrBukanMarineCargo
	}
	return dataSpreadingCargo{kasus: models.KasusSpreading{PercentShare: k.PercentShare, StartDateTime: k.StartDateTime},
		cargo: salinCargo(cargo)}, nil
}

// tampilanSpreadingCargo - total, ringkasan (SumTSIPremiSpreadedRNM_Act 17-24 bersama FIRE), pesan; periksaJenis =
// cekSpreadingFactIn langkah 8.
func (s *Service) tampilanSpreadingCargo(ctx context.Context, id string, d dataSpreadingCargo, pesan []string, periksaJenis bool) (TampilanSpreadingCargo, error) {
	cargo := salinCargo(d.cargo)
	dropdown, batas, err := s.treatyKasus(ctx, id, d.kasus)
	if err != nil {
		return TampilanSpreadingCargo{}, err
	}
	perMataUang, label := totalCargo(cargo, d.kasus.PercentShare)
	treaty, semua, pesanBatas, err := ringkasan(ctx, s.spreading, d.kasus, perMataUang, dropdown, batas)
	if err != nil {
		return TampilanSpreadingCargo{}, err
	}
	pesan = append(append([]string{}, pesan...), pesanBatas...)
	if periksaJenis {
		pesan = append(pesan, periksaJenisTreaty(perMataUang, batas)...)
	}
	bulatkanSimpanCargo(cargo)
	hasil := TampilanSpreadingCargo{PercentShare: d.kasus.PercentShare, Treaty: dropdown, Template: []models.TemplateSpreading{},
		Cargo: make([]CargoSpreading, 0, len(cargo)), RingkasanTreaty: bulatkanTotal(treaty), RingkasanMataUang: semua,
		Pesan: append([]string{}, pesan...), LabelLoL: label}
	for i, x := range cargo {
		hasil.Cargo = append(hasil.Cargo, CargoSpreading{No: i + 1, TradingNote: x.TradingNote, GoodNote: x.GoodNote,
			PackingNote: x.PackingNote, ConveyanceNote: x.ConveyanceNote, Coverages: x.Coverages})
	}
	for i := range semua {
		r := &semua[i]
		r.TSITopRisk, r.TSI, r.LoL, r.Premium = simpan8(r.TSITopRisk), simpan8(r.TSI), simpan8(r.LoL), simpan8(r.Premium)
	}
	return hasil, nil
}

// TampilanSpreadingCargo - GET /api/nbfacin/kasus/{caseId}/spreading-cargo.
func (s *Service) TampilanSpreadingCargo(ctx context.Context, id string) (TampilanSpreadingCargo, error) {
	d, err := s.bacaDataSpreadingCargo(ctx, id)
	if err != nil {
		return TampilanSpreadingCargo{}, err
	}
	return s.tampilanSpreadingCargo(ctx, id, d, nil, false)
}

// HitungShareSpreadingCargo - POST …/spreading-cargo/hitung-share: NR + baris (DT percent); TIDAK menyimpan; tanpa spreading
// otomatis.
func (s *Service) HitungShareSpreadingCargo(ctx context.Context, id string, ps *apd.Decimal) (TampilanSpreadingCargo, error) {
	if err := periksaShareKasus(ps); err != nil {
		return TampilanSpreadingCargo{}, err
	}
	d, err := s.bacaDataSpreadingCargo(ctx, id)
	if err != nil {
		return TampilanSpreadingCargo{}, err
	}
	d.kasus.PercentShare = ps
	hitungNRCargo(d.cargo, ps)
	terapkanPersenCargo(d.cargo)
	return s.tampilanSpreadingCargo(ctx, id, d, nil, false)
}

// SalinSpreadingCargo - POST …/spreading-cargo/salin: Copy To All + total; TIDAK menyimpan. Tahap uw + workbasket posisi dan
// pemeriksaan template sama dengan FIRE / Aneka.
func (s *Service) SalinSpreadingCargo(ctx context.Context, pelaku inti.Pelaku, id string, ps *apd.Decimal, templat []models.TemplateSpreading) (TampilanSpreadingCargo, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanSpreadingCargo{}, err
	}
	if err := s.wajibTahapUW(ctx, pelaku, id); err != nil {
		return TampilanSpreadingCargo{}, err
	}
	if err := periksaShareKasus(ps); err != nil {
		return TampilanSpreadingCargo{}, err
	}
	if len(templat) == 0 {
		return TampilanSpreadingCargo{}, fmt.Errorf("%w: template wajib berisi paling sedikit satu baris", ErrMasukanSpreading)
	}
	jenis, nama, share := make([]string, len(templat)), make([]string, len(templat)), make([]*apd.Decimal, len(templat))
	for i, t := range templat {
		jenis[i], nama[i], share[i] = t.TreatyType, t.TreatyName, t.SharePercentage
	}
	if m := periksaBarisShare("template", jenis, nama, share); len(m) > 0 {
		return TampilanSpreadingCargo{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(m, "; "))
	}
	d, err := s.bacaDataSpreadingCargo(ctx, id)
	if err != nil {
		return TampilanSpreadingCargo{}, err
	}
	m, err := s.proteksiTemplate(ctx, templat)
	if err != nil {
		return TampilanSpreadingCargo{}, err
	}
	if len(m) > 0 {
		return TampilanSpreadingCargo{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(m, "; "))
	}
	d.kasus.PercentShare = ps
	hitungNRCargo(d.cargo, ps)
	salinTemplateCargo(d.cargo, ps, templat)
	t, err := s.tampilanSpreadingCargo(ctx, id, d, nil, true)
	if err != nil {
		return t, err
	}
	t.Template = templat
	return t, nil
}

// samaBentukCargo - jumlah cargo / coverage dan kode coverage sama menurut posisi.
func samaBentukCargo(simpan, badan []models.ObjekCargo) bool {
	if len(simpan) != len(badan) {
		return false
	}
	for i := range simpan {
		if len(simpan[i].Coverages) != len(badan[i].Coverages) {
			return false
		}
		for m := range simpan[i].Coverages {
			if simpan[i].Coverages[m].Coverage != badan[i].Coverages[m].Coverage {
				return false
			}
		}
	}
	return true
}

// bacaBentukBadanCargo - % Share RNM 0..100, data tersimpan, coverage wajib ada ("Object can't be empty!", pola Aneka), bentuk
// badan = tersimpan (409).
func (s *Service) bacaBentukBadanCargo(ctx context.Context, id string, ps *apd.Decimal, cargo []models.ObjekCargo) (dataSpreadingCargo, error) {
	if err := periksaShareKasus(ps); err != nil {
		return dataSpreadingCargo{}, err
	}
	d, err := s.bacaDataSpreadingCargo(ctx, id)
	if err != nil {
		return d, err
	}
	if !adaCoverageCargo(d.cargo) {
		return d, fmt.Errorf("%w: %s", ErrMasukanSpreading, pesanObjekKosong)
	}
	if !samaBentukCargo(d.cargo, cargo) {
		return d, ErrSpreadingBerubah
	}
	return d, nil
}

// pasangBarisBadanCargo - treatyType / treatyName / sharePercentage baris badan menggantikan baris tersimpan (periksaBarisShare,
// 400).
func pasangBarisBadanCargo(simpan, badan []models.ObjekCargo) error {
	var masalah []string
	for i := range simpan {
		for m := range simpan[i].Coverages {
			baris := badan[i].Coverages[m].Spreading
			jenis, nama, share := make([]string, len(baris)), make([]string, len(baris)), make([]*apd.Decimal, len(baris))
			salin := make([]models.BarisSpreading, len(baris))
			for x, r := range baris {
				jenis[x], nama[x], share[x] = r.TreatyType, r.TreatyName, r.SharePercentage
				salin[x] = models.BarisSpreading{TreatyType: strings.TrimSpace(r.TreatyType), TreatyName: strings.TrimSpace(r.TreatyName),
					SharePercentage: r.SharePercentage, ClaimEstimation: r.ClaimEstimation, TSIGrossSpreaded: r.TSIGrossSpreaded}
			}
			masalah = append(masalah, periksaBarisShare(fmt.Sprintf("cargo[%d].coverages[%d].spreading", i, m), jenis, nama, share)...)
			simpan[i].Coverages[m].Spreading = salin
		}
	}
	if len(masalah) > 0 {
		return fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(masalah, "; "))
	}
	return nil
}

// hitungSaveSpreadingCargo - NR, DT percent (K74-7), 8 desimal.
func hitungSaveSpreadingCargo(cargo []models.ObjekCargo, ps *apd.Decimal) {
	hitungNRCargo(cargo, ps)
	terapkanPersenCargo(cargo)
	bulatkanSimpanCargo(cargo)
}

// SimpanSpreadingCargo - PUT …/spreading-cargo (Save tab): penjaga tahap = PUT spreading FIRE / Aneka; bentuk 409; hitungan
// Save; di uw jenis treaty tidak sah 400 (K59-8) dan langkah 11 (K59-9 / K59-10) seperti FIRE / Aneka; tulis pohon cargo
// (spreading ikut) + % Share RNM + langkah 11 satu transaksi.
func (s *Service) SimpanSpreadingCargo(ctx context.Context, pelaku inti.Pelaku, id string, ps *apd.Decimal, cargo []models.ObjekCargo) (TampilanSpreadingCargo, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanSpreadingCargo{}, err
	}
	barisBadan, uw, err := s.penjagaPutSpreading(ctx, pelaku, id)
	if err != nil {
		return TampilanSpreadingCargo{}, err
	}
	d, err := s.bacaBentukBadanCargo(ctx, id, ps, cargo)
	if err != nil {
		return TampilanSpreadingCargo{}, err
	}
	if barisBadan {
		if err := pasangBarisBadanCargo(d.cargo, cargo); err != nil {
			return TampilanSpreadingCargo{}, err
		}
	}
	d.kasus.PercentShare = ps
	hitungSaveSpreadingCargo(d.cargo, ps)
	if s.transaksi == nil {
		return TampilanSpreadingCargo{}, ErrSpreadingTanpaDatabase
	}
	if err := s.cekKunciUnggah(ctx, id); err != nil { // tiket 74 P6: Save menulis ulang pohon cargo
		return TampilanSpreadingCargo{}, err
	}
	var l11 *dasarLangkah11
	var baris11 []models.BarisDeduction
	if uw {
		_, batas, err := s.treatyKasus(ctx, id, d.kasus)
		if err != nil {
			return TampilanSpreadingCargo{}, err
		}
		perMataUang, _ := totalCargo(salinCargo(d.cargo), ps)
		if pesan := periksaJenisTreaty(perMataUang, batas); len(pesan) > 0 {
			return TampilanSpreadingCargo{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, pesan[0])
		}
		if l11, err = s.bacaDasarLangkah11(ctx, id); err != nil {
			return TampilanSpreadingCargo{}, err
		}
		baris11 = barisLangkah11(kasusSpreadingCargo(d.kasus, d.cargo), l11)
		if pesan := pesanDeduction2(baris11); len(pesan) > 0 {
			return TampilanSpreadingCargo{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, pesan[0])
		}
		var negatif []string
		for _, b := range baris11 {
			if netPremium(b).Sign() < 0 {
				negatif = append(negatif, pesanNetNegatif(b.Currency))
			}
		}
		if len(negatif) > 0 {
			return TampilanSpreadingCargo{}, fmt.Errorf("%w: %s", ErrMasukanSpreading, strings.Join(negatif, "; "))
		}
	}
	err = s.transaksi(ctx, func(tx *db.Tx) error {
		if err := s.objekCargo.GantiObjekCargo(ctx, tx, id, d.cargo); err != nil {
			return err
		}
		if err := s.objekCargo.UbahPercentShare(ctx, tx, id, simpan8(ps)); err != nil {
			return err
		}
		return s.tulisLangkah11UW(ctx, tx, id, l11, baris11)
	})
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return TampilanSpreadingCargo{}, ErrKasusTidakAda
	}
	if err != nil {
		return TampilanSpreadingCargo{}, err
	}
	return s.TampilanSpreadingCargo(ctx, id)
}

// HitungBarisSpreadingCargo - POST …/spreading-cargo/hitung-baris: hitungan Save tahap uw tanpa menulis + cekSpreadingFactIn
// langkah 8-9. Tahap uw + workbasket posisi.
func (s *Service) HitungBarisSpreadingCargo(ctx context.Context, pelaku inti.Pelaku, id string, ps *apd.Decimal, cargo []models.ObjekCargo) (TampilanSpreadingCargo, error) {
	if err := inti.WajibIdentitas(pelaku); err != nil {
		return TampilanSpreadingCargo{}, err
	}
	if err := s.wajibTahapUW(ctx, pelaku, id); err != nil {
		return TampilanSpreadingCargo{}, err
	}
	d, err := s.bacaBentukBadanCargo(ctx, id, ps, cargo)
	if err != nil {
		return TampilanSpreadingCargo{}, err
	}
	if err := pasangBarisBadanCargo(d.cargo, cargo); err != nil {
		return TampilanSpreadingCargo{}, err
	}
	d.kasus.PercentShare = ps
	hitungNRCargo(d.cargo, ps)
	terapkanPersenCargo(d.cargo)
	return s.tampilanSpreadingCargo(ctx, id, d, nil, true)
}

// spreadingRingkasCargo - tiket 74 P4: case MARINE CARGO tanpa item FIRE / objek Aneka -> pohon cargo (kasusSpreadingCargo) sebagai
// sumber Payment / Premium Deduction / Cedant; selain itu k apa adanya (FIRE / Aneka tidak berubah). Dibaca di luar tx.
func (s *Service) spreadingRingkasCargo(ctx context.Context, k models.KasusSpreading, id string) (models.KasusSpreading, error) {
	if s.objekCargo == nil || !s.adalahMarineCargo(ctx, id) {
		return k, nil
	}
	cargo, err := s.objekCargo.BacaObjekCargo(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return k, nil
	}
	if err != nil {
		return k, err
	}
	if !adaCoverageCargo(cargo) {
		return k, nil
	}
	return kasusSpreadingCargo(k, cargo), nil
}
