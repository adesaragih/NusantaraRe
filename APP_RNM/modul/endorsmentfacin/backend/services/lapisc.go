package services

// Lapis C - salinan daftar baris dari polis lama, penanda baris warisan, dan
// penggabungan spreading. Tiket E10.
//
// Untuk apa berkas ini: `SetValueToEDMWork` langkah 14.7-14.13 memanggil tujuh
// varian `Activity/SetOLDValueToEDMWork_<LOB>.xml`, masing-masing bergerbang
// predikat lininya.
//
// Dibaca sesudah: beforeimage.go.
//
// ⛔ `[terverifikasi]` Temuan pembacaan korpus 01-10-2026 yang TIDAK tertulis
// di spec/tiket E10 (dicatat di `docs/LAPORAN-IMPLEMENTASI-BEFORE-IMAGE.md`
// §2, menunggu work owner). Bukti: pohon langkah ketujuh
// `Activity/SetOLDValueToEDMWork_<LOB>.xml`, diurai dengan
// `docs/alat/langkah.py`:
//
//  1. Tiap varian lebih dulu MENYALIN daftar barisnya dari OldData (langkah
//     1.1) - inilah satu-satunya tempat `LocationList`/`CargoList`/
//     `VehicleList`/`PersonList` kerja terisi; 54 salinan lapis A tidak
//     menyentuhnya. Lapis C bukan sekadar penanda.
//  2. Enam varian (semua kecuali LIFE) MENGGABUNGKAN baris spreading tiap
//     coverage per `TreatyType`. Langkahnya berlabel "EDM ADJ SPREADING" dan
//     bergerbang `Type=="4"`, tetapi `pyStepsPreCondition=false` - jadi
//     TETAP JALAN untuk semua jenis endorsement (P-11 tertutup 19-09-2026).
//  3. Penanda tingkat atas `FlagOldData` dipasang MC, MBU, PA, DAN LIFE -
//     bukan hanya LIFE seperti premis kasus uji `K046_Life_DuaPenandaOldData`.
//     FIRE, Aneka, dan Golf memakai `IsOldData` di tingkat lokasi.
//  4. Kedalaman FIRE berakhir di `SpreadingList`, bukan `LayerList` (nol
//     kemunculan `LayerList` di berkas FIRE). LIFE punya
//     `PersonList → CoverageList → SpreadingList`, bukan "PersonList saja".
//
// `[terverifikasi]` Cacah lama di bahan spec (FIRE 12, MC 8, Aneka 14, …) adalah cacah BARIS
// yang memuat token `IsOldData` - termasuk nama properti dan cache editor -
// bukan cacah penugasan. Penugasan bernilai `"old"`: FIRE 5, MC 4, Aneka 6,
// MBU 5, LIFE 3, PA 4, Golf 5 (dua cara: urai pohon langkah dan
// `grep -c '<PropertiesValue>"old"</PropertiesValue>'`, keduanya sepakat).

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/endorsmentfacin/backend/models"
)

// TypeAdjSpreading - `QuotationData.Type == "4"`, dilabeli "EDM ADJ SPREADING"
// di deskripsi langkah. Arti kode lain: OQ-020.
const TypeAdjSpreading models.JenisPenyesuaian = "4"

var (
	// ErrSpreadingTakTerjumlahkan - penggabungan spreading menyentuh nilai
	// kosong. Pega menjumlahkan properti desimal kosong dengan cara yang
	// belum terverifikasi; nilai tidak ditebak nol.
	ErrSpreadingTakTerjumlahkan = errors.New("endorsement: baris spreading kosong tidak dapat digabungkan")
	// ErrSkalaRasioBerbeda - dua SharePercentage berbeda jumlah desimal.
	ErrSkalaRasioBerbeda = errors.New("endorsement: dua rasio berbeda skala tidak dijumlahkan")
)

// tandaiBarisWarisan - langkah 14.7-14.13, URUTAN KORPUS: FIRE, MC, Aneka,
// MBU, LIFE, PA, Golf. Tiap varian bergerbang di pemanggil DAN di langkah 1-
// nya sendiri dengan predikat yang sama.
func tandaiBarisWarisan(k *models.KasusEndorsement, p Predikat) error {
	langkah := []struct {
		aktif bool
		jalan func(*models.KasusEndorsement) error
	}{
		{p.IsFire, varianFire},
		{p.IsMarineCargo, varianMarineCargo},
		{p.IsAneka, varianAneka},
		{p.IsMBU, varianMBU},
		{p.IsLife, varianLife},
		{p.IsPA, varianPA},
		{p.IsGolfInsurance, varianGolf},
	}
	for _, l := range langkah {
		if !l.aktif {
			continue
		}
		if err := l.jalan(k); err != nil {
			return err
		}
	}
	return nil
}

// varianFire - `SetOLDValueToEDMWork_FIRE`.
func varianFire(k *models.KasusEndorsement) error {
	o := &k.OfferFacIn
	// 1.1
	o.LocationList = models.KlonLokasi(o.OldData.LocationList)
	// ⛔ HANYA varian ini menyetel penanda prorata tingkat agregat, dan ia
	// tidak termasuk 54 salinan lapis A. Alasannya: pertanyaan terbuka.
	o.IsProRate = models.NilaiProrate
	// 1.2 Lokasi → PropertyItem → Coverage → Spreading.
	for i := range o.LocationList {
		lok := &o.LocationList[i]
		lok.IsOldData = models.NilaiOld
		for j := range lok.Property.PropertyItemList {
			item := &lok.Property.PropertyItemList[j]
			item.IsOldData = models.NilaiOld
			for c := range item.CoverageList {
				cov := &item.CoverageList[c]
				cov.IsOldData = models.NilaiOld
				// Hanya varian FIRE yang mengosongkan `.EDM`.
				cov.EDM = ""
				if err := gabungSpreading(cov, buangBagianNolFire); err != nil {
					return fmt.Errorf("FIRE lokasi %d item %d coverage %d: %w", i+1, j+1, c+1, err)
				}
				tandaiSpreadingBilaBukanAdj(cov, o.QuotationData.Type)
			}
		}
	}
	return nil
}

// varianMarineCargo - `SetOLDValueToEDMWork_MC`.
func varianMarineCargo(k *models.KasusEndorsement) error {
	o := &k.OfferFacIn
	// 1.1 - selain CargoList, varian ini menyalin PolicyData UTUH, mengosongkan
	// nomor endorsement, dan membawa tiga daftar dari kargo PERTAMA ke akar
	// objek kerja.
	o.CargoList = models.KlonKargo(o.OldData.CargoList)
	o.PolicyData = o.OldData.PolicyData
	o.PolicyData.EndorsementNo = ""
	// [dugaan] Tanpa kargo pertama, ketiga daftar akar menjadi kosong; perilaku
	// Pega saat merujuk `CargoList(1)` yang tidak ada belum terverifikasi.
	k.TradingList, k.GoodsList, k.ConveyanceList = nil, nil, nil
	if len(o.OldData.CargoList) > 0 {
		pertama := o.OldData.CargoList[0]
		k.TradingList = models.KlonJSON(pertama.TradingList)
		// `newWorkPage.GoodsList = …CargoList(1).GoodList` - nama properti
		// sumber dan tujuan memang berbeda di korpus.
		k.GoodsList = models.KlonJSON(pertama.GoodList)
		k.ConveyanceList = models.KlonJSON(pertama.ConveyanceList)
	}
	// 1.2 Kargo → Coverage → Spreading.
	for i := range o.CargoList {
		kg := &o.CargoList[i]
		kg.FlagOldData = models.NilaiOld
		for c := range kg.CoverageList {
			cov := &kg.CoverageList[c]
			cov.IsOldData = models.NilaiOld
			if err := gabungSpreading(cov, buangBagianTakPositif); err != nil {
				return fmt.Errorf("MC kargo %d coverage %d: %w", i+1, c+1, err)
			}
			tandaiSpreadingBilaBukanAdj(cov, o.QuotationData.Type)
		}
	}
	return nil
}

// varianAneka - `SetOLDValueToEDMWork_Aneka`.
func varianAneka(k *models.KasusEndorsement) error {
	o := &k.OfferFacIn
	o.LocationList = models.KlonLokasi(o.OldData.LocationList)
	// Lokasi → Occupation → Aneka → Coverage → Spreading.
	for i := range o.LocationList {
		lok := &o.LocationList[i]
		lok.IsOldData = models.NilaiOld
		for j := range lok.Property.RiskLocation.OccupationList {
			okp := &lok.Property.RiskLocation.OccupationList[j]
			okp.IsOldData = models.NilaiOld
			for a := range okp.AnekaList {
				an := &okp.AnekaList[a]
				an.IsOldData = models.NilaiOld
				for c := range an.CoverageList {
					cov := &an.CoverageList[c]
					cov.IsOldData = models.NilaiOld
					if err := gabungSpreading(cov, buangBagianTakPositif); err != nil {
						return fmt.Errorf("Aneka lokasi %d okupasi %d aneka %d coverage %d: %w", i+1, j+1, a+1, c+1, err)
					}
					tandaiSpreadingBilaBukanAdj(cov, o.QuotationData.Type)
				}
			}
		}
	}
	return nil
}

// varianMBU - `SetOLDValueToEDMWork_MBU`.
func varianMBU(k *models.KasusEndorsement) error {
	o := &k.OfferFacIn
	o.VehicleList = models.KlonKendaraan(o.OldData.VehicleList)
	// Kendaraan → Coverage (+ AdditionalCoverage) → Spreading.
	for i := range o.VehicleList {
		v := &o.VehicleList[i]
		v.FlagOldData = models.NilaiOld
		for c := range v.CoverageList {
			cov := &v.CoverageList[c]
			cov.IsOldData = models.NilaiOld
			for a := range cov.AdditionalCoverage {
				cov.AdditionalCoverage[a].IsOldData = models.NilaiOld
			}
			if err := gabungSpreading(cov, buangBagianTakPositif); err != nil {
				return fmt.Errorf("MBU kendaraan %d coverage %d: %w", i+1, c+1, err)
			}
			tandaiSpreadingBilaBukanAdj(cov, o.QuotationData.Type)
		}
	}
	return nil
}

// varianLife - `SetOLDValueToEDMWork_LIFE`. Satu-satunya varian TANPA
// penggabungan spreading dan tanpa gerbang `Type`.
func varianLife(k *models.KasusEndorsement) error {
	o := &k.OfferFacIn
	o.PersonList = models.KlonOrang(o.OldData.PersonList)
	// Orang → CoverageList → Spreading.
	for i := range o.PersonList {
		p := &o.PersonList[i]
		p.FlagOldData = models.NilaiOld
		for c := range p.CoverageList {
			cov := &p.CoverageList[c]
			cov.IsOldData = models.NilaiOld
			for s := range cov.SpreadingList {
				cov.SpreadingList[s].IsOldData = models.NilaiOld
			}
		}
	}
	return nil
}

// varianPA - `SetOLDValueToEDMWork_PA`.
func varianPA(k *models.KasusEndorsement) error {
	o := &k.OfferFacIn
	o.PersonList = models.KlonOrang(o.OldData.PersonList)
	// Orang → ASMCoverage → Spreading.
	for i := range o.PersonList {
		p := &o.PersonList[i]
		p.FlagOldData = models.NilaiOld
		for c := range p.ASMCoverage {
			cov := &p.ASMCoverage[c]
			cov.IsOldData = models.NilaiOld
			if err := gabungSpreading(cov, buangBagianTakPositif); err != nil {
				return fmt.Errorf("PA orang %d coverage %d: %w", i+1, c+1, err)
			}
			tandaiSpreadingBilaBukanAdj(cov, o.QuotationData.Type)
		}
	}
	return nil
}

// varianGolf - `SetOLDValueToEDMWork_GOLF`.
func varianGolf(k *models.KasusEndorsement) error {
	o := &k.OfferFacIn
	o.LocationList = models.KlonLokasi(o.OldData.LocationList)
	// Lokasi → RiskLocation.Aneka → Coverage → Spreading.
	for i := range o.LocationList {
		lok := &o.LocationList[i]
		lok.IsOldData = models.NilaiOld
		for a := range lok.Property.RiskLocation.AnekaList {
			an := &lok.Property.RiskLocation.AnekaList[a]
			an.IsOldData = models.NilaiOld
			for c := range an.CoverageList {
				cov := &an.CoverageList[c]
				cov.IsOldData = models.NilaiOld
				if err := gabungSpreading(cov, buangBagianTakPositif); err != nil {
					return fmt.Errorf("Golf lokasi %d aneka %d coverage %d: %w", i+1, a+1, c+1, err)
				}
				tandaiSpreadingBilaBukanAdj(cov, o.QuotationData.Type)
			}
		}
	}
	return nil
}

// aturanBuang - kapan baris hasil penggabungan dibuang.
type aturanBuang func(share *apd.Decimal) bool

// buangBagianNolFire - FIRE: `Page-Remove` bergerbang `.SharePercentage<>0`
// WhenTrue=lewati → baris dibuang bila bagiannya TEPAT nol. Bagian negatif
// TETAP.
func buangBagianNolFire(share *apd.Decimal) bool { return share.Sign() == 0 }

// buangBagianTakPositif - MC, Aneka, MBU, PA, Golf: gerbang `.SharePercentage>0`
// WhenTrue=lewati → baris dibuang bila bagiannya nol ATAU negatif. Berbeda
// dari FIRE; diport apa adanya.
func buangBagianTakPositif(share *apd.Decimal) bool { return share.Sign() <= 0 }

// gabungSpreading - langkah "EDM ADJ SPREADING" (1.2.x.2 di tiap varian).
//
//	TempSpread = halaman baru
//	untuk tiap baris s di .SpreadingList:
//	    s.IsOldData = "old"
//	    untuk tiap t di TempSpread: bila t.TreatyType == s.TreatyType →
//	        t.TSISpreaded += s, t.PremiumSpreaded += s, t.SharePercentage += s
//	    bila tidak ada yang cocok → tambahkan baris baru berisi LIMA properti
//	untuk tiap t di TempSpread: Page-Remove menurut aturanBuang
//	Property-Remove .SpreadingList; .SpreadingList = TempSpread
//
// ⚠️ `[dugaan]` `Page-Remove TempSpread.pxResults(Local.Idx)` dijalankan DI
// DALAM loop atas daftar yang sama. Apakah Pega lalu melewatkan baris
// sesudah yang dibuang belum terverifikasi; di sini SEMUA baris yang memenuhi
// aturan dibuang.
func gabungSpreading(cov *models.Coverage, buang aturanBuang) error {
	var gabung []models.Spreading
	for i := range cov.SpreadingList {
		s := &cov.SpreadingList[i]
		s.IsOldData = models.NilaiOld
		cocok := false
		for j := range gabung {
			t := &gabung[j]
			if t.TreatyType != s.TreatyType {
				continue
			}
			var err error
			if t.TSISpreaded, err = jumlahUang(t.TSISpreaded, s.TSISpreaded); err != nil {
				return fmt.Errorf("TSISpreaded treaty %q: %w", s.TreatyType, err)
			}
			if t.PremiumSpreaded, err = jumlahUang(t.PremiumSpreaded, s.PremiumSpreaded); err != nil {
				return fmt.Errorf("PremiumSpreaded treaty %q: %w", s.TreatyType, err)
			}
			if t.SharePercentage, err = jumlahRasio(t.SharePercentage, s.SharePercentage); err != nil {
				return fmt.Errorf("SharePercentage treaty %q: %w", s.TreatyType, err)
			}
			cocok = true
		}
		if !cocok {
			gabung = append(gabung, models.Spreading{
				IsOldData:       s.IsOldData,
				PremiumSpreaded: s.PremiumSpreaded,
				SharePercentage: s.SharePercentage,
				TreatyType:      s.TreatyType,
				TSISpreaded:     s.TSISpreaded,
			})
		}
	}
	sisa := []models.Spreading{}
	for _, t := range gabung {
		if t.SharePercentage.Kosong() {
			return fmt.Errorf("SharePercentage treaty %q: %w", t.TreatyType, ErrSpreadingTakTerjumlahkan)
		}
		if !buang(t.SharePercentage.Value) {
			sisa = append(sisa, t)
		}
	}
	cov.SpreadingList = sisa
	return nil
}

// tandaiSpreadingBilaBukanAdj - langkah sesudah penggabungan: bergerbang
// `Type=="4"` dengan WhenTrue=lewati (prakondisi AKTIF), jadi penanda baris
// spreading dipasang hanya bila jenis endorsement BUKAN "4".
func tandaiSpreadingBilaBukanAdj(cov *models.Coverage, jenis models.JenisPenyesuaian) {
	if jenis == TypeAdjSpreading {
		return
	}
	for s := range cov.SpreadingList {
		cov.SpreadingList[s].IsOldData = models.NilaiOld
	}
}

func jumlahUang(a, b uang.Money) (uang.Money, error) {
	if a.Kosong() || b.Kosong() {
		return uang.Money{}, ErrSpreadingTakTerjumlahkan
	}
	return a.Add(b)
}

func jumlahRasio(a, b uang.Ratio) (uang.Ratio, error) {
	if a.Kosong() || b.Kosong() {
		return uang.Ratio{}, ErrSpreadingTakTerjumlahkan
	}
	if a.Scale != b.Scale {
		return uang.Ratio{}, fmt.Errorf("%w: %d dan %d", ErrSkalaRasioBerbeda, a.Scale, b.Scale)
	}
	hasil := new(apd.Decimal)
	if _, err := utils.DecimalContext().Add(hasil, a.Value, b.Value); err != nil {
		return uang.Ratio{}, err
	}
	return uang.Ratio{Value: hasil, Scale: a.Scale}, nil
}
