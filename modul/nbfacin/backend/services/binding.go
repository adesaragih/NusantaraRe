package services

// Accept tahap Underwriting ke jalur akhir (tiket 63). Flow\InputInwardFacultativeOffer.xml: Decision23 Else (Transition35,
// tanpa SET) -> Utility16 "SAVE JSON_OFFER" (Activity\SaveJsonOfferFacIn_Act.xml: JSON_OFFER + langkah 20 Call
// SaveOfferProduction_Act -> FACINOFFER) -> Decision16 "It is Group?":
//   - WHEN IsGroup -> SubProcess1 POLICY (Transition60: PositionNote ReasFacInAdmin, flag 1) = tiket 64 -> DITAHAN pesanPolicyGroup
//     (keputusan work owner D5; tidak ada yang ditulis);
//   - ELSE -> Transition38 (`PositionNote = "ReasFacInMarketing"`, `.IsCedingConfirm = "Binding"`, `.FlagOnGoingPolicy = 2`,
//     NBStatus = prefix + " IS IN " + MarketingName + "'S INBOX") -> Utility6 SendEmailPolicy -> Assignment9 "MARKETING
//     (BINDING)" (ToWorkBasket ReasFacInMarketing). IsCedingConfirm tidak punya kolom di sistem baru: status Binding = FLAG 2
//     (kotak masuk). Case Binding baca-saja sampai tiket 64 (tahap marketing mensyaratkan flag 0).
// Penyimpangan sadar (keputusan work owner D1): JSON_OFFER.JSONDATA = dokumen minimal berisi kunci yang dibaca SQL korpus
// (`IDNewBisnis`, GetNoOffer_Sql / UpdateTgl_BindJsonOffer_SQL), bukan halaman OfferFacIn utuh (@ASM.GetPageJSONString) -
// sejalan ADR-0028; syaratnya tidak ada sistem hilir yang membaca isi JSONDATA (pertanyaan terbuka, issues/63).

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/nbfacin/backend/models"
	"nusantarare/modul/nbfacin/backend/repository"
)

const (
	// pesanPolicyGroup - Decision16 IsGroup -> SubProcess POLICY (tiket 64) belum tersedia (keputusan work owner D5).
	pesanPolicyGroup = "Accept meneruskan case ke jalur Policy (case group) yang belum tersedia."
	// flagBinding - Transition38 `.FlagOnGoingPolicy = 2`.
	flagBinding = "2"
	// bentukTanggalFacInOffer - padanan "dd/MM/yyyy hh:mm:ss" Asia/Jakarta SaveOfferProduction_Act, 24-jam (D3 b).
	bentukTanggalFacInOffer = "02/01/2006 15:04:05"
)

// ErrBindingTidakSejalan - objek dan spreading case tidak sejalan (jumlah / OLDID berbeda): FACINOFFER tidak dapat disusun.
var ErrBindingTidakSejalan = errors.New("services: objek dan spreading case tidak sejalan; FACINOFFER tidak dapat disusun")

// tujuanBinding - Transition38 (bukan group) atau ditahan pesanPolicyGroup.
func (s *Service) tujuanBinding(id string, k models.KonteksSubmit, pelaku inti.Pelaku) tujuanUW {
	if s.grup.kasusGroup(k, pelaku) {
		log.Printf("nbfacin: Accept UW kasus %s: Decision16 IsGroup - jalur Policy (tiket 64) ditahan", id)
		return tujuanUW{pesan: pesanPolicyGroup}
	}
	return tujuanUW{posisi: repository.WorkbasketMarketing, status: awalanNBStatus() + " IS IN " + k.MarketingName + "'S INBOX",
		binding: true}
}

// tanggalFacInOffer - teks Pega DateTime -> jam dinding WIB "DD/MM/YYYY HH24:MI:SS"; kosong / tak terurai -> "" (NULL).
func tanggalFacInOffer(teks string) string {
	t, err := time.Parse(BentukWaktuPega, strings.TrimSpace(teks))
	if err != nil {
		return ""
	}
	return t.In(WIB).Format(bentukTanggalFacInOffer)
}

// kodeOkupasiFacInOffer - SaveOfferProduction_Act 12.2: baris 1 = OccupationId; berikutnya `.OccupationId + " || " +
// Datain.CARI29` (yang baru DI DEPAN).
func kodeOkupasiFacInOffer(ok []models.OkupasiObjek) string {
	var hasil string
	for i, o := range ok {
		if i == 0 {
			hasil = o.OccupationID
			continue
		}
		hasil = o.OccupationID + " || " + hasil
	}
	return hasil
}

// dasarFacInOffer - kolom FACINOFFER tingkat case (semua cabang COB).
func dasarFacInOffer(id string, k models.KonteksSubmit, g models.General, kodeBisnis, groupType string) models.BarisFacInOffer {
	return models.BarisFacInOffer{IDPega: id, BusinessCode: kodeBisnis, CedingCo: g.CedingCoName, SOB: k.SobName,
		BeginDate: tanggalFacInOffer(g.StartDateTime), EndDate: tanggalFacInOffer(g.EndDateTime), PercentRNM: k.PercentShare,
		MarketingOfficerCode: k.MOID, InsuredName: k.InsuredName, PercentRIComm: k.RICommision, RISlipCeding: g.ReffNumber,
		GroupType: groupType, GroupName: g.GroupName}
}

// barisFacInOffer - cabang FIRE langkah 12: satu baris per SpreadingList coverage (coverage tanpa spreading tidak menulis
// apa pun). Objek (alamat, okupasi, coverage) dan spreading dipasangkan per indeks lokasi / item / coverage + OLDID.
func barisFacInOffer(id string, k models.KonteksSubmit, g models.General, kodeBisnis, groupType string, objek []models.ObjekFire,
	sp models.KasusSpreading) ([]models.BarisFacInOffer, error) {
	if len(objek) != len(sp.Lokasi) {
		return nil, fmt.Errorf("%w: %d lokasi objek, %d lokasi spreading", ErrBindingTidakSejalan, len(objek), len(sp.Lokasi))
	}
	dasar := dasarFacInOffer(id, k, g, kodeBisnis, groupType)
	hasil := []models.BarisFacInOffer{}
	for i, o := range objek {
		lok := dasar
		lok.Province, lok.Kabupaten, lok.Kecamatan, lok.Kelurahan = o.Province, o.City, o.District, o.Territory
		lok.RoadName, lok.BuildingNo, lok.ZipCode, lok.ObjectNo = o.RoadType+o.RoadName, "NO. "+o.BuildingNo, o.ZipCode, o.ObjectNo
		lok.ObjectName, lok.RiskLocation, lok.TopRisk = o.ObjectName, o.RiskLocation, teksBoolPega(o.IsTopRisk)
		lok.OccupationCode = kodeOkupasiFacInOffer(o.Occupations)
		if len(o.Items) != len(sp.Lokasi[i].Items) {
			return nil, fmt.Errorf("%w: lokasi %d", ErrBindingTidakSejalan, i+1)
		}
		for j, it := range o.Items {
			spIt := sp.Lokasi[i].Items[j]
			if len(it.Coverages) != len(spIt.Coverages) {
				return nil, fmt.Errorf("%w: lokasi %d item %d", ErrBindingTidakSejalan, i+1, j+1)
			}
			for m, c := range it.Coverages {
				spC := spIt.Coverages[m]
				if spC.OldID != c.OldID {
					return nil, fmt.Errorf("%w: lokasi %d item %d coverage %d", ErrBindingTidakSejalan, i+1, j+1, m+1)
				}
				for _, sb := range spC.Spreading {
					b := lok
					b.ObjectItem, b.CurrID = it.ItemType, it.Currency
					b.Coverage, b.CoverageNote, b.AccumulationCode = c.Coverage, c.CoverageNote, c.AccumulationCode
					if c.Rate != nil {
						b.Rate = utils.FormatDecimal(c.Rate)
					}
					b.PercentShareSpreaded, b.JnReas, b.TSI, b.Premi = sb.SharePercentage, sb.TreatyType, sb.TSISpreaded, sb.PremiumSpreaded
					hasil = append(hasil, b)
				}
			}
		}
	}
	return hasil, nil
}

// tulisanBinding - JSON_OFFER (dokumen minimal) + FACINOFFER + flag, dibaca di luar transaksi (posisi dijaga TulisSubmit).
func (s *Service) tulisanBinding(ctx context.Context, id string, k models.KonteksSubmit) (*models.TulisanBinding, error) {
	if s.kasus == nil {
		return nil, ErrSubmitTanpaDatabase
	}
	kasus, err := s.kasus.BacaKasus(ctx, id)
	if errors.Is(err, repository.ErrKasusTidakAda) {
		return nil, ErrKasusTidakAda
	}
	if err != nil {
		return nil, err
	}
	kode, err := s.kodeBisnisKasus(ctx, id)
	if err != nil {
		return nil, err
	}
	objek, err := s.BacaObjek(ctx, id)
	if err != nil {
		return nil, err
	}
	// GROUP_TYPE: SaveOfferProduction_Act langkah 7-8 menurut QuotationData.IsGroup ("Group" -> "1", "NonGroup" -> "2") = A.
	groupType := "2"
	if s.grup.pembuatGroup(k) {
		groupType = "1"
	}
	// Tiket 73 P4: cabang 13 / 14 / 15 (Aneka / Bonding / Golf) bila case keluarga Aneka; selain itu cabang FIRE (tiket 63).
	offer, aneka, err := s.barisFacInOfferKasusAneka(ctx, id, dasarFacInOffer(id, k, kasus.General, kode, groupType), objek)
	if err != nil {
		return nil, err
	}
	if !aneka {
		sp, err := s.bacaSpreading(ctx, nil, id)
		if err != nil {
			return nil, err
		}
		if offer, err = barisFacInOffer(id, k, kasus.General, kode, groupType, objek, sp); err != nil {
			return nil, err
		}
	}
	doc, err := json.Marshal(map[string]string{"IDNewBisnis": id})
	if err != nil {
		return nil, err
	}
	return &models.TulisanBinding{FlagBaru: flagBinding, JSONData: string(doc), BulanNoOffer: s.sekarang().In(WIB).Format("01.2006"),
		Offer: offer}, nil
}

// ---- email "Waiting for Bind" (Activity\SendEmailPolicy.xml) --------------------------------------------------------

// domainInternal - SendEmailPolicy langkah 5: Emailto bukan alamat internal / kosong -> diganti CC.
const domainInternal = "@nusantarare.com"

// surelBinding - langkah 2-8, 11, 13: To = email pembuat case (bukan internal / kosong -> CC), CC = email marketing + atasan
// (GetMKTandLeader_SQL), Subject langkah 8 verbatim. BCC langkah 8 (2 alamat individu ter-hardcode) TIDAK diport (D4).
// `[belum terverifikasi]` Badan HTML NotifEmailPolicy_HTML tidak ada di korpus (hanya dirujuk langkah 11): badan ringkas
// `[keputusan agent]`. ok false = langkah 12 Exit-Activity (Emailto kosong).
func surelBinding(id string, k models.KonteksSubmit, nbStatus, pembuat string, cc []string) (models.Surel, bool) {
	to := pembuat
	if to == "" || !strings.Contains(strings.ToLower(to), domainInternal) {
		to = strings.Join(cc, ",")
	}
	if to == "" {
		return models.Surel{}, false
	}
	subjek := "Waiting for Bind ID Pega : " + id + " - " + k.InsuredName
	badan := "<p>" + htmlAman(subjek) + "</p><p>" + htmlAman(nbStatus) + "</p>"
	return models.Surel{To: strings.Split(to, ","), Cc: cc, Subject: subjek, HTML: badan}, true
}

// htmlAman - escape teks untuk badan HTML.
func htmlAman(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

// kirimSurelBinding - langkah 14 `Call SendEmailWithAttachments` hanya bila IsPEGAPROD (IS_PEGA_PROD, ADR-U-0005) dan SMTP
// terpasang; dijalankan SESUDAH commit Submit. Galat hanya dicatat (Pega tidak membaca hasil kirim) - Submit tetap berhasil.
func (s *Service) kirimSurelBinding(ctx context.Context, id string, k models.KonteksSubmit, nbStatus string) {
	if !s.produksi || s.surel == nil || s.submit == nil {
		return
	}
	pembuat, cc, err := s.submit.PenerimaBinding(ctx, k.CreateOp, k.MOID)
	if err != nil {
		log.Printf("nbfacin: email Waiting for Bind kasus %s: penerima: %v", id, err)
		return
	}
	sur, ok := surelBinding(id, k, nbStatus, pembuat, cc)
	if !ok {
		log.Printf("nbfacin: email Waiting for Bind kasus %s: tanpa penerima - tidak dikirim", id)
		return
	}
	if err := s.surel.KirimSurel(ctx, sur); err != nil {
		log.Printf("nbfacin: email Waiting for Bind kasus %s GAGAL ke %d alamat: %v", id, len(sur.To)+len(sur.Cc), err)
		return
	}
	log.Printf("nbfacin: email Waiting for Bind kasus %s terkirim ke %d alamat", id, len(sur.To)+len(sur.Cc))
}
