package services

// FACINOFFER case keluarga Aneka (tiket 73 P4 butir 5). Port `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\Activity\
// SaveOfferProduction_Act.xml` (dibaca 09-10-2026): cabang 13 Aneka (WHEN IsAneka L8787), 14 Bonding (IsBonding L11798), 15
// Golf (isGolfInsurance L14443); satu INSERT per baris SpreadingList, alamat lokasi = cabang FIRE (L5650-5856):
//   - OCCUPATIONCODE (CARI29): Aneka / Bonding = nama okupasi digabung bertahap - loop AnekaList bersarang DI DALAM loop okupasi
//     (13.2.3), jadi baris okupasi ke-n memuat `nama_n || … || nama_1` (13.2.1 `.pxListSubscript==1` L6233; 13.2.2 `>1` L6375);
//     Golf = "" (L12182);
//   - OBJECT_ITEM (CARI16) = `.ObjectName` objek (L6558); CURR_ID (CARI23) = `.Currency.Name` objek (L6604), Aneka IsMBD = Currency
//     coverage (L7240, WHEN L7317; kosong -> objek, CountPremiAndTSIRNMFireAnekaGolf_ACT L2969);
//   - NOSERIAL (Datain1.CARI5) = `.SerialNo` objek, "NULL" bila kosong (L6624 / L6729) - Aneka saja; Bonding / Golf "NULL";
//   - TOP_RISK (CARI48) tidak di-SET cabang 13-15 -> NULL `[dugaan: Datain baru]`.
//
// K73-35 (keputusan work owner 09-10-2026): case Bonding (IsAneka DAN IsBonding benar) SATU baris per spreading lewat cabang
// Bonding saja (Pega: cabang 13 dan 14 sama-sama jalan -> baris ganda `[dugaan]`).

import (
	"context"
	"fmt"
	"strings"

	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
)

// liniFacInOffer - cabang SaveOfferProduction_Act yang dijalankan.
type liniFacInOffer struct {
	aneka, bonding, golf, mbd bool
}

// barisFacInOfferAneka - baris FACINOFFER cabang 13 / 14 / 15; dasar = kolom tingkat case (sama dengan FIRE). objek = baris
// T_LOCATIONLIST yang sama dengan lokasi Aneka (alamat); jumlah wajib sama.
func barisFacInOfferAneka(dasar models.BarisFacInOffer, objek []models.ObjekFire, lokasi []models.LokasiAneka, l liniFacInOffer) ([]models.BarisFacInOffer, error) {
	if len(objek) != len(lokasi) {
		return nil, fmt.Errorf("%w: %d lokasi objek, %d lokasi Aneka", ErrBindingTidakSejalan, len(objek), len(lokasi))
	}
	hasil := []models.BarisFacInOffer{}
	baris := func(lok models.BarisFacInOffer, o models.ObjekAneka, c models.CoverageAneka, serial bool) {
		for _, sb := range c.Spreading {
			b := lok
			b.ObjectItem, b.CurrID = o.ObjectName, o.Currency
			if l.mbd && !l.bonding {
				b.CurrID = mataUangCoverage(o, c)
			}
			if serial {
				b.NoSerial = strings.TrimSpace(o.SerialNo) // "" = "NULL"
			}
			b.Coverage, b.CoverageNote = c.Coverage, c.CoverageNote
			if c.Rate != nil {
				b.Rate = utils.FormatDecimal(c.Rate)
			}
			b.PercentShareSpreaded, b.JnReas, b.TSI, b.Premi = sb.SharePercentage, sb.TreatyType, sb.TSISpreaded, sb.PremiumSpreaded
			hasil = append(hasil, b)
		}
	}
	for i, x := range lokasi {
		o := objek[i]
		lok := dasar
		lok.Province, lok.Kabupaten, lok.Kecamatan, lok.Kelurahan = o.Province, o.City, o.District, o.Territory
		lok.RoadName, lok.BuildingNo, lok.ZipCode, lok.ObjectNo = o.RoadType+o.RoadName, "NO. "+o.BuildingNo, o.ZipCode, o.ObjectNo
		lok.ObjectName, lok.RiskLocation, lok.TopRisk = o.ObjectName, o.RiskLocation, ""
		if l.aneka || l.bonding {
			okupasi := ""
			for j, ok := range x.Okupasi {
				if j == 0 {
					okupasi = ok.OccupationName
				} else {
					okupasi = ok.OccupationName + " || " + okupasi
				}
				lokOk := lok
				lokOk.OccupationCode = okupasi
				for _, ob := range ok.Objek {
					for _, c := range ob.Coverages {
						baris(lokOk, ob, c, !l.bonding)
					}
				}
			}
		}
		if l.golf {
			for _, ob := range x.ObjekGolf {
				for _, c := range ob.Coverages {
					baris(lok, ob, c, false)
				}
			}
		}
	}
	return hasil, nil
}

// barisFacInOfferKasusAneka - nil bila case bukan keluarga Aneka (pemanggil memakai cabang FIRE apa adanya).
func (s *Service) barisFacInOfferKasusAneka(ctx context.Context, id string, dasar models.BarisFacInOffer, objek []models.ObjekFire) ([]models.BarisFacInOffer, bool, error) {
	if s.objekAneka == nil {
		return nil, false, nil
	}
	kb, ada := s.kasusBisnisDari(ctx, id)
	if !ada {
		return nil, false, nil
	}
	l := liniFacInOffer{aneka: salahSatu(id, kb, "IsAneka"), bonding: salahSatu(id, kb, "IsBonding"), golf: salahSatu(id, kb, "isGolfInsurance"),
		mbd: salahSatu(id, kb, "IsMBD")}
	if !l.aneka && !l.bonding && !l.golf {
		return nil, false, nil
	}
	if l.bonding {
		l.aneka = false // K73-35
	}
	lokasi, err := s.objekAneka.BacaObjekAneka(ctx, id)
	if err != nil {
		return nil, false, galatKasus(err)
	}
	b, err := barisFacInOfferAneka(dasar, objek, lokasi, l)
	return b, true, err
}
