package services

// Selesai Policy case keluarga Aneka (tiket 73 P5). Port `[terverifikasi]` `D:\migrasi\RNM\NB FacIn\Activity\` (dibaca
// 10-10-2026; rincian docs/issues/73 "Backend P5"):
//   - barisFacInProductionAneka = SaveTreatyProduction_Act (`ASM-FW-GISFW-WORK!SAVETREATYPRODUCTION_ACT`) langkah 19 Aneka (WHEN
//     IsAneka L7539; Bonding ikut - langkah 20 Bonding berlabel `//`) dan 21 Golf (isGolfInsurance L13423): satu baris per lokasi ×
//     okupasi × objek × coverage × baris spreading (Golf: lokasi × objek × coverage × baris); kolom yang berbeda dari FIRE:
//     OBJECT_ITEM = `.ObjectName` objek (L4939 / L11297), CURR_ID = Currency objek, IsMBD = Currency coverage (L5292),
//     OCCUPATIONCODE = NAMA okupasi baris (L4614, K73-42; Golf "" L10940), NOSERIAL = `.SerialNo` (L5005; Golf "NULL"),
//     OBJECT_NAME_ID teks "NULL" (S16 L3490), MARKETINGOFFICERCODE = MOID case didahulukan (S11 L2492, K73-46), Golf LOL 0.0
//     (L12228); TOP_RISK / IDX_OBJECT_ITEM / TSI100_* tidak di-SET -> NULL (K73-43); DEDUCTION2_* NULL (K73-45); alamat Golf =
//     ASMCity / ASMDistrict / ASMRW seperti Aneka (K73-44; Pega Golf CityName / DistrictName / RWName L10800-10840);
//     tulis-backup Golf (InsertFacinProductionBackup_Sql L12730) tidak ditiru (K73-47). ACCUMULATIONCODE = "" (coverage Aneka
//     tidak menyimpan AccumulationCode; form Aneka tidak mengisinya).
//   - profilRisikoLossAneka = SaveJsonPolicyFacIn_Act langkah 19-21 untuk case Aneka: langkah 19 hanya menjelajah pohon FIRE
//     (L4965 / L5012), jadi Local tak pernah di-SET; nilai bawaan Local (pyLocalParameters): TSINusare String (L437 / L439) -> NULL,
//     PremiNusare / PremiumOR / HasilQs / PremiumFacout double (L443-469) -> 0 (K73-48 ikut Pega; double bawaan 0 `[dugaan:
//     platform]`).

import (
	"context"
	"strings"

	"nusantarare/modul/nbfacin/backend/models"
)

// barisFacInProductionAneka - lihat kepala berkas; dasar = kolom tingkat case (dasarBarisProduksi); objek = baris T_LOCATIONLIST
// yang sama dengan lokasi Aneka (alamat).
func barisFacInProductionAneka(d dasarProduksi, objek []models.ObjekFire, lokasi []models.LokasiAneka, l liniFacInOffer) ([]models.BarisFacInProduction, error) {
	if len(objek) != len(lokasi) {
		return nil, ErrBindingTidakSejalan
	}
	ri := nolBila(d.k.RICommision)
	dasar := dasarBarisProduksi(d)
	dasar.MOIDKasus = strings.TrimSpace(d.k.MOID)
	hasil := []models.BarisFacInProduction{}
	baris := func(lok models.BarisFacInProduction, o models.ObjekAneka, c models.CoverageAneka, golf bool) {
		for _, sb := range c.Spreading {
			b := lok
			b.ObjectItem, b.CurrID, b.ObjectNameID = o.ObjectName, o.Currency, "NULL"
			if l.mbd && !golf {
				b.CurrID = mataUangCoverage(o, c)
			}
			if !golf {
				b.NoSerial = strings.TrimSpace(o.SerialNo)
			}
			b.Coverage, b.CoverageNote, b.Rate, b.LayerNo = c.Coverage, c.CoverageNote, simpan8(c.Rate), layerNoNull
			b.PercentShareSpreaded, b.JnReas = simpan8(sb.SharePercentage), sb.TreatyType
			b.TSI, b.Premi = simpan8(sb.TSISpreaded), simpan8(sb.PremiumSpreaded)
			b.LoL = simpan8(nolBila(sb.ClaimEstimation))
			if golf {
				b.LoL = simpan8(nol()) // L12228 literal 0.0
			}
			b.RIComm = simpan8(persenDari(sb.PremiumSpreaded, ri))
			b.BrokerageFee = simpan8(persenDari(sb.PremiumSpreaded, d.pctBrokerage))
			hasil = append(hasil, b)
		}
	}
	for i, x := range lokasi {
		o := objek[i]
		lok := dasar
		lok.Province, lok.Kabupaten, lok.Kecamatan, lok.Kelurahan = o.Province, o.City, o.District, o.Territory
		lok.RoadName, lok.BuildingNo, lok.ZipCode, lok.ObjectNo = o.RoadType+o.RoadName, "NO. "+o.BuildingNo, o.ZipCode, o.ObjectNo
		lok.ObjectName, lok.RiskLocation = o.ObjectName, o.RiskLocation
		if l.aneka {
			for _, ok := range x.Okupasi {
				lokOk := lok
				lokOk.OccupationCode = ok.OccupationName
				for _, ob := range ok.Objek {
					for _, c := range ob.Coverages {
						baris(lokOk, ob, c, false)
					}
				}
			}
		}
		if l.golf {
			for _, ob := range x.ObjekGolf {
				for _, c := range ob.Coverages {
					baris(lok, ob, c, true)
				}
			}
		}
	}
	return hasil, nil
}

// profilRisikoLossAneka - lihat kepala berkas.
func profilRisikoLossAneka(cob string) models.ProfilRisikoLoss {
	return models.ProfilRisikoLoss{COB: cob, PremiNusantaraRe: nol(), PremiumOR: nol(), PremiumQS: nol(), PremiumFacout: nol()}
}

// liniProduksiAneka - cabang SaveTreatyProduction_Act case keluarga Aneka; ok=false = cabang FIRE (lini tak terbaca / bukan
// Aneka / penyimpan objek Aneka tidak terpasang). Bonding lewat cabang 19 (IsAneka).
func (s *Service) liniProduksiAneka(ctx context.Context, id string) (liniFacInOffer, []models.LokasiAneka, bool, error) {
	if s.objekAneka == nil {
		return liniFacInOffer{}, nil, false, nil
	}
	kb, ada := s.kasusBisnisDari(ctx, id)
	if !ada {
		return liniFacInOffer{}, nil, false, nil
	}
	l := liniFacInOffer{aneka: salahSatu(id, kb, "IsAneka"), golf: salahSatu(id, kb, "isGolfInsurance"), mbd: salahSatu(id, kb, "IsMBD")}
	if !l.aneka && !l.golf {
		return l, nil, false, nil
	}
	lokasi, err := s.objekAneka.BacaObjekAneka(ctx, id)
	if err != nil {
		return l, nil, false, galatKasus(err)
	}
	return l, lokasi, true, nil
}
