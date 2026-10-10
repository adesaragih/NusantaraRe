package services

// Selesai Policy (tiket 69a-2) - Flow\InputInwardFacultativeRISlip.xml Decision4 Else / Decision21 Else -> Utility1 "SAVE
// JSON_POLICY" (Activity\SaveJsonPolicyFacIn_Act.xml) jalur NB FIRE. Keputusan work owner K69a2-1..9 dan bukti di
// docs/issues/69-jalur-policy.md. Transisi ke Utility1 tidak men-SET PositionNote / NBStatus: POSITION dan NB_STATUS tetap,
// STATUS_WORK "Resolved-Completed" (K69-7) - case baca-saja lewat kasusSelesai. Email Policy dan service Arasapas
// (Utility7 / Utility2) = tiket 69b.
//
// Yang disusun di sini (repository menulis dalam transaksi Submit, repository/selesaipolicy.go):
//   - CEDING_FACINPRODUCTION - InsertCedingProduction langkah 2 per CedingCedantList;
//   - M_RISK_LOSS_PROFILE - SaveJsonPolicy langkah 19-20 (variabel Local ditimpa per iterasi - ditiru);
//   - FACINPRODUCTION - SaveFacinProdFireNB_Act langkah 4-13 (cabang FIRE);
//   - DATA_JSON JSON_POLIS minimal `{"IDNewBisnis": caseId}` (K69a2-1 / ADR-0028; work owner 09-10-2026 "sistem baru tidak pakai
//     data_json lgi" - isian kunci SQL korpus versi 08-10-2026 DIBATALKAN; saran akumulasi per polis membaca FACINPRODUCTION).

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

const (
	// statusSelesaiPolicy - K69-7.
	statusSelesaiPolicy = "Resolved-Completed"
	// statusBusinessNB - Quotation.StatusBusiness case NB ("1", seperti kasusTangga).
	statusBusinessNB = "1"
	// layerNoNull - FireNB 13.3.2.1 `Datain1.CARI12 = "NULL"` (teks, K69a2-4).
	layerNoNull = "NULL"
	// basisLayer - CoverageBasis "5" (13.3.2.2 per LayerList) - layering tidak dipakai di sistem baru (butir 30).
	basisLayer = "5"
)

// ErrBusinessOldIDTidakAda - BusinessOldId case tidak dapat ditentukan (BUSINESS 0 / > 1 baris / OLDID kosong): nomor polis
// tidak dibentuk (Pega memakai pxResults(1) / kosong - `[keputusan agent]` gagal terang). 409.
var ErrBusinessOldIDTidakAda = errors.New("BusinessOldId case tidak dapat ditentukan dari tabel BUSINESS (Class of Business / " +
	"Group Business); nomor polis tidak dapat dibentuk.")

// treatyQS - SaveJsonPolicy 19.1.2.2.1 (TreatyType QS); treatyOR - 19.1.2.2.2 (10001 / 10007); FACOUT 19.1.2.2.3 (10015).
var (
	treatyQS = map[string]bool{"10185": true, "10184": true, "10186": true, "10003": true, "10004": true, "10011": true,
		"10012": true, "10014": true, "10019": true, "10021": true, "10022": true, "10024": true, "10025": true, "10026": true,
		"10028": true}
	treatyOR       = map[string]bool{"10001": true, treatyORS: true}
	nolKomaDelapan = apd.New(8, -1)
	nolKomaDua     = apd.New(2, -1)
)

// profilRisikoLoss - SaveJsonPolicy langkah 19-20 apa adanya: TSINusare / PremiNusare = coverage TERAKHIR; PremiumQs dikosongkan
// per coverage; HasilQs = PremiumQs × 0.8 baris QS terakhir; PremiumOR = PremiumQs (coverage itu, kosong = 0 `[dugaan]`) × 0.2 +
// PremiumSpreaded baris OR terakhir; PremiumFacout = baris FACOUT terakhir. Tak pernah diisi -> nil (NULL).
func profilRisikoLoss(sp models.KasusSpreading, cob string) models.ProfilRisikoLoss {
	p := models.ProfilRisikoLoss{COB: cob}
	for _, l := range sp.Lokasi {
		for _, it := range l.Items {
			for _, c := range it.Coverages {
				p.TSINusantaraRe, p.PremiNusantaraRe = c.TSINusantaraRe, c.PremiNusantaraRe
				var qs *apd.Decimal
				for _, b := range c.Spreading {
					switch {
					case treatyQS[b.TreatyType]:
						qs = b.PremiumSpreaded
						p.PremiumQS = kali(nolBila(qs), nolKomaDelapan)
					case treatyOR[b.TreatyType]:
						p.PremiumOR = tambah(kali(nolBila(qs), nolKomaDua), b.PremiumSpreaded)
					case b.TreatyType == treatyFACOUT:
						p.PremiumFacout = b.PremiumSpreaded
					}
				}
			}
		}
	}
	p.TSINusantaraRe, p.PremiNusantaraRe = simpan8(p.TSINusantaraRe), simpan8(p.PremiNusantaraRe)
	p.PremiumOR, p.PremiumQS, p.PremiumFacout = simpan8(p.PremiumOR), simpan8(p.PremiumQS), simpan8(p.PremiumFacout)
	return p
}

// barisCedingProduksi - InsertCedingProduction langkah 2 per CedingCedantList: 2.1 kode / nama tanpa ";"; PCT_SHARE berurutan
// (langkah sesudahnya menimpa): 2.2 GROSS ("0") ShareCeding / PercentShare (skala 20) × 100; 2.3 NET ("1") / kosong ShareCeding;
// 2.4 kosong + ShareCeding kosong 100 / jumlah baris (skala 4); 2.5 kosong + satu baris 100. Dibulatkan 8 desimal
// (angkaMasuk; kolom NUMBER(23,20)).
func barisCedingProduksi(k models.KasusCedant, percentShare *apd.Decimal) []models.BarisCedingProduksi {
	var hasil []models.BarisCedingProduksi
	n := len(k.Cedant)
	for _, c := range k.Cedant {
		var share *apd.Decimal
		switch k.ShareCedantType {
		case tipeCedantGross:
			if percentShare != nil && !percentShare.IsZero() && c.ShareCeding != nil {
				share = kali(bagiSkala(c.ShareCeding, percentShare, 20), seratus)
			}
		case tipeCedantShareRNM, "":
			share = c.ShareCeding
		}
		if k.ShareCedantType == "" && c.ShareCeding == nil {
			share = bagiSkala(seratus, apd.New(int64(n), 0), 4)
		}
		if k.ShareCedantType == "" && n == 1 {
			share = seratus
		}
		hasil = append(hasil, models.BarisCedingProduksi{CedingID: strings.ReplaceAll(c.CedingCo, ";", ""),
			CedingName: strings.ReplaceAll(c.CedingCoName, ";", ""), Share: simpan8(share)})
	}
	return hasil
}

// kodeOkupasiProduksi - FireNB 13.2: baris 1 = OccupationId; berikutnya `.OccupationId + ";" + Datain.CARI29` (yang baru DI
// DEPAN; pemisah ";" - berbeda dengan FACINOFFER " || ").
func kodeOkupasiProduksi(ok []models.OkupasiObjek) string {
	var hasil string
	for i, o := range ok {
		if i == 0 {
			hasil = o.OccupationID
			continue
		}
		hasil = o.OccupationID + ";" + hasil
	}
	return hasil
}

// dasarProduksi - nilai FireNB langkah 4-12 yang sama untuk semua baris.
type dasarProduksi struct {
	id, kodeBisnis, groupType string
	k                         models.KonteksSubmit
	kasus                     models.Kasus
	pctBrokerage              *apd.Decimal
}

// dasarBarisProduksi - kolom FACINPRODUCTION tingkat case (langkah 1-17 / FireNB 4-12, sama untuk semua cabang COB). Langkah
// 11-12: RICommision kosong -> "0" (PERCENT_RI_COMM / PCT_RI_COMM_SELISIH); ricomm / brokerage memakai nilai mentah (kosong
// dihitung 0 `[dugaan]`); % brokerage kosong -> NULL.
func dasarBarisProduksi(d dasarProduksi) models.BarisFacInProduction {
	g := d.kasus.General
	ids := make([]string, 0, len(g.CedingList))
	for _, c := range g.CedingList {
		ids = append(ids, c.ID)
	}
	ri := nolBila(d.k.RICommision)
	return models.BarisFacInProduction{IDPega: idPega(d.id), BusinessCode: d.kodeBisnis, CedingCo: g.CedingCoName,
		SOB: d.k.SobName, BeginDate: tanggalFacInOffer(g.StartDateTime), EndDate: tanggalFacInOffer(g.EndDateTime),
		PercentRNM: simpan8(d.k.PercentShare), MarketingCode: d.k.MarketingCode, InsuredName: d.k.InsuredName,
		PercentRIComm: simpan8(ri), PctRICommSelisih: simpan8(ri), RISlipCeding: g.ReffNumber, GroupType: d.groupType,
		QQName: g.QQName, PolicyType: kodePolicyType(g.PolicyType), CedingID: strings.Join(ids, ";"), SOBID: g.SourceOfBusinessID,
		InsuredID: d.kasus.Opportunity.InsuredID, StatusBusiness: statusBusinessNB, GroupName: g.GroupName,
		PctBrokerageFee: simpan8(d.pctBrokerage)}
}

// barisFacInProduction - FireNB langkah 13 (IsFire, NB): satu baris per lokasi × item × coverage × SpreadingList (13.3.2.3,
// CoverageBasis != "5"); coverage basis "5" mengulang LayerList (13.3.2.2) yang tidak ada di sistem baru (butir 30) -> tanpa
// baris, seperti LayerList kosong. Objek dan spreading dipasangkan seperti FACINOFFER (ErrBindingTidakSejalan).
func barisFacInProduction(d dasarProduksi, objek []models.ObjekFire, sp models.KasusSpreading) ([]models.BarisFacInProduction, error) {
	if len(objek) != len(sp.Lokasi) {
		return nil, ErrBindingTidakSejalan
	}
	ri := nolBila(d.k.RICommision)
	dasar := dasarBarisProduksi(d)
	hasil := []models.BarisFacInProduction{}
	for i, o := range objek {
		lok := dasar
		lok.Province, lok.Kabupaten, lok.Kecamatan, lok.Kelurahan = o.Province, o.City, o.District, o.Territory
		lok.RoadName, lok.BuildingNo, lok.ZipCode, lok.ObjectNo = o.RoadType+o.RoadName, "NO. "+o.BuildingNo, o.ZipCode, o.ObjectNo
		lok.ObjectName, lok.RiskLocation, lok.TopRisk = o.ObjectName, o.RiskLocation, teksBoolPega(o.IsTopRisk)
		lok.OccupationCode = kodeOkupasiProduksi(o.Occupations)
		if len(o.Items) != len(sp.Lokasi[i].Items) {
			return nil, ErrBindingTidakSejalan
		}
		for j, it := range o.Items {
			spIt := sp.Lokasi[i].Items[j]
			if len(it.Coverages) != len(spIt.Coverages) {
				return nil, ErrBindingTidakSejalan
			}
			for m, c := range it.Coverages {
				spC := spIt.Coverages[m]
				if spC.OldID != c.OldID {
					return nil, ErrBindingTidakSejalan
				}
				if strings.TrimSpace(c.CoverageBasis) == basisLayer {
					log.Printf("nbfacin: selesai Policy kasus %s: coverage %d.%d.%d basis 5 tanpa LayerList - tanpa baris FACINPRODUCTION",
						d.id, i+1, j+1, m+1)
					continue
				}
				for _, sb := range spC.Spreading {
					b := lok
					b.ObjectItem, b.CurrID, b.ObjectNameID, b.IdxObjectItem = it.ItemType, it.Currency, it.ItemTypeID, strconv.Itoa(j+1)
					b.Coverage, b.CoverageNote, b.AccumulationCode, b.Rate = c.Coverage, c.CoverageNote, c.AccumulationCode, simpan8(c.Rate)
					b.LayerNo = layerNoNull
					b.PercentShareSpreaded, b.JnReas = simpan8(sb.SharePercentage), sb.TreatyType
					b.TSI, b.Premi, b.TSI100 = simpan8(sb.TSISpreaded), simpan8(sb.PremiumSpreaded), simpan8(sb.TSIGrossSpreaded)
					b.LoL = simpan8(nolBila(sb.ClaimEstimation))
					b.RIComm = simpan8(persenDari(sb.PremiumSpreaded, ri))
					b.BrokerageFee = simpan8(persenDari(sb.PremiumSpreaded, d.pctBrokerage))
					hasil = append(hasil, b)
				}
			}
		}
	}
	return hasil, nil
}

// dokumenJSONPolis - DATA_JSON minimal `{"IDNewBisnis": caseId}` (pola JSON_OFFER tiket 63 D1). Work owner 09-10-2026 "sistem
// baru tidak pakai data_json lgi": isi polis TIDAK disalin ke DATA_JSON; pembaca sistem baru memakai tabel (FACINPRODUCTION).
func dokumenJSONPolis(id string) (string, error) {
	b, err := json.Marshal(map[string]string{"IDNewBisnis": id})
	return string(b), err
}

// tulisanPolicy - TulisanPolicy case (dibaca di luar transaksi; posisi dijaga TulisSubmit).
func (s *Service) tulisanPolicy(ctx context.Context, id string, k models.KonteksSubmit, pelaku inti.Pelaku) (*models.TulisanPolicy, error) {
	if s.kasus == nil || s.cedant == nil || s.deduction == nil {
		return nil, ErrSubmitTanpaDatabase
	}
	kasus, err := s.kasus.BacaKasus(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	oldID, err := s.oldIDBisnisKasus(ctx, id)
	if err != nil {
		return nil, err
	}
	if oldID == "" {
		return nil, ErrBusinessOldIDTidakAda
	}
	kode, err := s.kodeBisnisKasus(ctx, id)
	if err != nil {
		return nil, err
	}
	objek, err := s.BacaObjek(ctx, id)
	if err != nil {
		return nil, err
	}
	sp, err := s.bacaSpreading(ctx, nil, id)
	if err != nil {
		return nil, err
	}
	cedant, err := s.cedant.BacaCedant(ctx, id)
	if err != nil {
		return nil, galatKasus(err)
	}
	ded, err := s.deduction.BacaDeduction(ctx, id)
	if err != nil {
		return nil, galatKasus(err)
	}
	// GROUP_TYPE: FireNB 7-9 menurut QuotationData.IsGroup ("Group" 1 / "NonGroup" 2; "SemiGroup" tidak lahir di sistem baru).
	groupType := "2"
	if s.grup.pembuatGroup(k) {
		groupType = "1"
	}
	dp := dasarProduksi{id: id, kodeBisnis: kode, groupType: groupType, k: k, kasus: kasus, pctBrokerage: ded.PctBrokerageFee}
	// Tiket 73 P5: SaveTreatyProduction_Act cabang 19 / 21 (keluarga Aneka) + profil risiko Aneka; selain itu FireNB (tiket 69a-2).
	lini, lokasi, aneka, err := s.liniProduksiAneka(ctx, id)
	if err != nil {
		return nil, err
	}
	var prod []models.BarisFacInProduction
	risiko := profilRisikoLoss(sp, k.ClassOfBusiness)
	if aneka {
		prod, err = barisFacInProductionAneka(dp, objek, lokasi, lini)
		risiko = profilRisikoLossAneka(k.ClassOfBusiness)
	} else {
		prod, err = barisFacInProduction(dp, objek, sp)
	}
	if err != nil {
		return nil, err
	}
	doc, err := dokumenJSONPolis(id)
	if err != nil {
		return nil, err
	}
	p := &models.TulisanPolicy{BusinessOldID: oldID, Sekarang: s.sekarang(), Username: pelaku.AkunID, JSONData: doc,
		Ceding: barisCedingProduksi(cedant, k.PercentShare), RisikoLoss: risiko, Produksi: prod}
	// PolicyData.StartDateTime (T_GENERAL_POLIS.START_DATE_TIME); kosong / tak terurai -> nol (langkah 5.2 dilewati).
	if t, err := time.Parse(BentukWaktuPega, strings.TrimSpace(kasus.General.StartDateTime)); err == nil {
		p.MulaiPolis = t
	}
	return p, nil
}
