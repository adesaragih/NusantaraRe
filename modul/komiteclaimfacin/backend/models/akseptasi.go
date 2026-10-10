package models

// Untuk apa berkas ini: BARIS OS_AKSEPTASI_KLAIM yang ditulis komite dan penanda tingkat akhir (korpus `Komite Claim
// FacIn`). Murni; layanan yang menulis.
//
//	TT2 tingkat akhir setuju  SaveAcceptation_KMT S1-S5 (IsFacretro / DLAStatus)          -> `TandaAkseptasi`
//	                          SaveAccept_ACT S1-S9 + KomitePost_Adjustment S8              -> `SusunOSAkseptasi`
//	TT3 setuju                KomitePost_Reject S12 -> SaveReject_ACT_KMT (per estimasi)   -> `SusunOSTolak`
//	TT4 setuju                KomitePost_CloseClaim S12 (InputParamOs)                     -> `SusunOSTutup`
//
// Procedure `PEGA_JSON_OS_AKSEP_KLAIM` (ALL_SOURCE DEV 10-10-2026) ditulis ulang tanpa procedure, sama dengan Claim Fac In
// tahap 1: INSERT CASEID, NOCLAIM, DATA_JSON, TANGGAL (hari ini), NOPOLIS, STS_REJECT, STS_KONVERSI (CARI16), STS_DLA
// (CARI17), MASTERID (CARI18 - tidak diisi pemanggil mana pun), CLAIMOLD (CARI20 - idem); TGL_PROD diisi trigger.
//
// PERBAIKAN prompt §5 butir 5: KomitePost_Adjustment S8 menulis OS SEKALI per KMT dengan InputData adjustment TERAKHIR
// yang disusun SaveAccept_ACT; satu adjustment per KMT (OQ-CFI-28, KOMITE_ID UNIQUE) - hasilnya sama: satu baris untuk
// adjustment kasus komite itu. Medan `ClaimData.osAkseptasi` / `stsReject` klaim induk (SaveAccept_ACT S3-S9,
// SaveReject_ACT_KMT S2-S8) tidak ditulis: tanpa kolom dan tanpa pembaca (PARITAS).

import (
	"bytes"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"nusantarare/inti/backend/kontrak"
)

// Kelas halaman JSON (DATA_JSON OS `CLM-` DEV 10-10-2026: halaman dan baris AcceptationList osAkseptasi, baris
// CurrencyList Currency).
const (
	KelasOSAkseptasi = "ASM-FW-GCNMFW-Data-osAkseptasi"
	KelasMataUang    = "ASM-FW-GCNMFW-Data-Currency"
)

// Status OS_AKSEPTASI_KLAIM.STS_REJECT.
const (
	StsOSAkseptasi = "1" // SaveAccept_ACT `TempOSAkseptasi.Type := 1`
	StsOSFinal     = "4" // SaveAccept_ACT S8 CARI10 PaymentType 1; KomitePost_CloseClaim S12.1
	StsOSTolak     = "2" // SaveReject_ACT_KMT `Type := 2`
)

// STS_DLA (`InputData.CARI17`).
const (
	StsDLAFac   = "7"
	StsDLARetro = "8"
)

// BarisOS - satu baris OS_AKSEPTASI_KLAIM (kolom yang ditulis procedure).
type BarisOS struct {
	CaseID, NoClaim, NoPolis       string
	StsReject, StsKonversi, StsDLA string
	DataJSON                       string
}

// HalamanJSON - satu halaman untuk `@ASM.GetPageJSONString()`: nilai teks + daftar halaman (PageList).
type HalamanJSON struct {
	Nilai  map[string]string
	Daftar []DaftarJSON
}

// DaftarJSON - satu PageList bernama.
type DaftarJSON struct {
	Nama string
	Isi  []HalamanJSON
}

// JSONHalamanPega = `@ASM.GetPageJSONString()`. Fungsinya tidak diekspor; bentuknya dibaca dari DATA_JSON
// OS_AKSEPTASI_KLAIM `CLM-` DEV 10-10-2026 (DISALIN dari Claim Fac In `models/osakseptasi.go`, bukan impor):
//   - "{" LF, pasangan pertama, setiap pasangan berikutnya diawali LF ",", ditutup LF "}"; halaman puncak diakhiri LF;
//   - nilai teks dulu, urut kunci tanpa membedakan huruf besar; properti kosong tidak ditulis;
//   - PageList sesudah nilai teks: `"Nama":[ ` LF anggota LF `] ` (anggota dipisah LF ",").
func JSONHalamanPega(p HalamanJSON) string { return tulisHalamanJSON(p) + "\n" }

func tulisHalamanJSON(p HalamanJSON) string {
	kunci := make([]string, 0, len(p.Nilai))
	for k, v := range p.Nilai {
		if v != "" {
			kunci = append(kunci, k)
		}
	}
	sort.Slice(kunci, func(i, j int) bool {
		a, b := strings.ToLower(kunci[i]), strings.ToLower(kunci[j])
		if a != b {
			return a < b
		}
		return kunci[i] < kunci[j]
	})
	var bagian []string
	for _, k := range kunci {
		bagian = append(bagian, literalJSON(k)+":"+literalJSON(p.Nilai[k]))
	}
	for _, d := range p.Daftar {
		if len(d.Isi) == 0 {
			continue
		}
		var isi []string
		for _, a := range d.Isi {
			isi = append(isi, tulisHalamanJSON(a))
		}
		bagian = append(bagian, literalJSON(d.Nama)+":[ \n"+strings.Join(isi, "\n,")+"\n] ")
	}
	return "{\n" + strings.Join(bagian, "\n,") + "\n}"
}

// literalJSON - literal teks JSON tanpa escape HTML.
func literalJSON(s string) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s)
	return strings.TrimSuffix(buf.String(), "\n")
}

// tanggalStempel = `@getCurrentDateStamp()` (yyyyMMdd, hari Jakarta; DEV 8 digit).
func tanggalStempel(t time.Time) string { return t.In(Jakarta).Format("20060102") }

// noPolis - `TempOpenPage.OfferFacIn.PolicyData.PolicyNo` (CARI3).
func noPolis(kl kontrak.KlaimFacIn) string { return kl.Nilai["OfferFacIn.PolicyData.PolicyNo"] }

// TandaAkseptasi = SaveAcceptation_KMT S1-S5 (dipanggil KomitePost_Adjustment S7.2.1.16 bila tingkat akhir setuju dan
// `.IsPrintAccept == ""`): ada spreading item polis ber-TreatyType 10015 -> IsFacretro objek, adjustment, item = 1,
// selainnya 0 (S3 / S4); DLAStatus objek = 0 (S5).
func TandaAkseptasi(kl kontrak.KlaimFacIn) (objek, item, adj map[string]string) {
	retro := "0"
	for _, s := range kl.Daftar[DaftarDiItem(kl.Objek, kl.Item, AnakSpreadPL)] { // S2 SpreadingList item
		if s["TreatyType"] == TreatyRetro { // S2.1
			retro = "1"
		}
	}
	return map[string]string{"IsFacretro": retro, "DLAStatus": "0"}, map[string]string{"IsFacretro": retro},
		map[string]string{"IsFacRetro": retro}
}

// nilaiBayar - nilai AcceptationList menurut PaymentType (SaveAccept_ACT S3.5-S3.7 dan cabang lini lain).
func nilaiBayar(adj map[string]string) (value, deductible string) {
	switch adj["PaymentType"] {
	case "1", "2", "5":
		return adj["AdjustmentValue"], adj["IndividualRiskValue"]
	case "3":
		return adj["SalvageValue"], adj["IndividualRiskValue"]
	case "4", "6":
		return adj["AdjusterFeeValue"], adj["VATValue"]
	}
	return "", ""
}

// SusunOSAkseptasi = SaveAccept_ACT S1-S9 (adjustment kasus komite, sesudah S7.2.1.5 menulis AcceptedNo) dan
// KomitePost_Adjustment S8 SaveOSClaim_SQL. `adj` = baris adjustment sesudah keputusan, `retro` = IsFacretro objek
// sesudah SaveAcceptation_KMT (S8 CARI17).
func SusunOSAkseptasi(kl kontrak.KlaimFacIn, adj map[string]string, retro, komiteID string, saat time.Time) BarisOS {
	cd := func(p string) string { return kl.Nilai["ClaimData."+p] }
	ob, it := Objek(kl), Item(kl)
	cur, curID := adj["Currency"], adj["CurrencyID"]
	if curID == "" { // S2: mata uang dari nama mata uang adjustment dan ID item
		cur, curID = adj["CurrencyName"], it["CurrencyID"]
	}
	tgl := tanggalStempel(saat)
	value, ded := nilaiBayar(adj)
	totalDed := adj["IndividualRiskValue"] // S3.3 @if(PT 4 / 6, VATValue, IndividualRiskValue)
	if p := adj["PaymentType"]; p == "4" || p == "6" {
		totalDed = adj["VATValue"]
	}
	akseptasi := map[string]string{"pxObjClass": KelasOSAkseptasi, "ObjectID": ob["ObjectID"], "EstimationDate": tgl,
		"NetForCollection": adj["NetForCollection"]}
	switch adj["PaymentType"] { // S3.5-S3.7 (PT 7 tanpa cabang)
	case "1", "2", "3", "4", "5", "6":
		akseptasi["Value"], akseptasi["GrossValue"], akseptasi["Currency"] = value, adj["GrossAdjustment"], cur
		akseptasi["CurrencyID"], akseptasi["DeductibleValue"] = curID, ded
	}
	cov := baris(kl, DaftarDiItem(kl.Objek, kl.Item, AnakCoverage), 1)
	switch LiniOS(kl.Nilai) {
	case LiniFireAneka: // S3.4
		akseptasi["ObjectName"], akseptasi["CoverageID"] = ob["ObjectName"], it["CoverageID"]
		akseptasi["CoverageName"], akseptasi["ObjectItemName"] = it["CoverageNote"], it["ObjectItemName"]
	case LiniMarine: // S4.4
		akseptasi["ConveyanceName"], akseptasi["CoverageID"] = ob["ObjectSurveyLocation"], cov["Coverage"]
		akseptasi["CoverageName"], akseptasi["Goodsname"] = cov["CoverageNote"], ob["ObjectJob"]
		akseptasi["PackingName"], akseptasi["TradingName"] = ob["ObjectLocation"], ob["ObjectName"]
		akseptasi["ObjectItemName"] = it["ObjectItemName"]
	case LiniMBU: // S5.4
		akseptasi["BrandID"], akseptasi["BrandName"], akseptasi["CoverageID"] = ob["Brand"], ob["BrandName"], it["CoverageID"]
		akseptasi["CoverageName"], akseptasi["Currency"], akseptasi["CurrencyID"] = it["CoverageNote"], cur, curID
		akseptasi["LicensePlate"], akseptasi["Model"], akseptasi["ModelName"] = ob["LicensePlate"], ob["Model"], ob["ModelName"]
		akseptasi["TypeName"] = ob["TypeName"]
	case LiniTravel: // S6.4
		akseptasi["ObjectName"], akseptasi["CoverageID"], akseptasi["CoverageName"] = ob["ObjectName"], it["CoverageID"],
			it["CoverageNote"]
		akseptasi["Currency"], akseptasi["CurrencyID"] = cur, curID
		akseptasi["ObjectDateOfBirth"], akseptasi["ObjectIDCard"] = ob["ObjectDateOfBirth"], ob["ObjectIDCard"]
		akseptasi["ObjectParticipantStatus"] = ob["ObjectParticipantStatus"]
	case LiniPA: // S7.4
		akseptasi["ObjectName"], akseptasi["CoverageID"], akseptasi["CoverageName"] = ob["ObjectName"], it["CoverageID"],
			it["CoverageNote"]
		akseptasi["Currency"], akseptasi["CurrencyID"] = cur, curID
		akseptasi["ObjectDateOfBirth"], akseptasi["ObjectIDCard"] = ob["ObjectDateOfBirth"], ob["ObjectIDCard"]
		akseptasi["ObjectParticipantStatus"], akseptasi["Gender"] = ob["ObjectParticipantStatus"], ob["Gender"]
		akseptasi["ObjectJob"] = ob["ObjectJob"]
	default: // lini di luar lima cabang: tanpa baris AcceptationList (cabang tidak berjalan)
		akseptasi = nil
	}
	top := map[string]string{"pxObjClass": KelasOSAkseptasi, "CauseOfLoss": cd("CauseOfLoss"),
		"CauseOfLossID": cd("CauseOfLossID"), "Type": StsOSAkseptasi, "NoDla": adj["DLA_No"],
		"PersenRNM": adj["PersenRNM"], "NoClaim": cd("NoClaim"), "AcceptedNo": adj["AcceptedNo"], "EstimationDate": tgl,
		// S8
		"pxCreateOperator": adj["pxCreateOperator"], "KomiteNo": komiteID, "PaymentType": adj["PaymentType"]}
	hal := HalamanJSON{Nilai: top}
	if akseptasi != nil { // S3.3 / S3.4: CurrencyList(1) + AcceptationList(1)
		mu := HalamanJSON{Nilai: map[string]string{"pxObjClass": KelasMataUang, "CurrencyID": curID, "Currency": cur,
			"TotalGrossValue": adj["GrossAdjustment"], "TotalDeductibleValue": totalDed},
			Daftar: []DaftarJSON{{Nama: "AcceptationList", Isi: []HalamanJSON{{Nilai: akseptasi}}}}}
		hal.Daftar = []DaftarJSON{{Nama: "CurrencyList", Isi: []HalamanJSON{mu}}}
	}
	sts := StsOSAkseptasi
	if adj["PaymentType"] == "1" { // S8 CARI10
		sts = StsOSFinal
	}
	dla := StsDLAFac
	if retro == "1" { // S8 CARI17
		dla = StsDLARetro
	}
	return BarisOS{CaseID: KunciInstans(kl.Nilai["pyID"]), NoClaim: cd("NoClaim"), NoPolis: noPolis(kl), StsReject: sts,
		StsDLA: dla, DataJSON: JSONHalamanPega(hal)}
}

// SusunOSTolak = KomitePost_Reject S12 (AcceptStatus 1): setiap estimasi ber-PrintFaceClaim 1 seluruh objek / item ->
// SaveReject_ACT_KMT S9-S20 (STS 2). Halaman TempOSAkseptasi BERTAHAN antar estimasi (tidak dibuang di activity) -
// properti cabang lini yang sama ditimpa setiap kali.
func SusunOSTolak(kl kontrak.KlaimFacIn, saat time.Time) []BarisOS {
	cd := func(p string) string { return kl.Nilai["ClaimData."+p] }
	persen := kl.Nilai["OfferFacIn.PercentShare"]
	tgl := tanggalStempel(saat)
	temp := map[string]string{"pxObjClass": KelasOSAkseptasi}
	var out []BarisOS
	for on, ob := range kl.Daftar[DaftarObjek] {
		for in, it := range kl.Daftar[DaftarItem(on+1)] {
			cov := baris(kl, DaftarDiItem(on+1, in+1, AnakCoverage), 1)
			item1 := baris(kl, DaftarItem(on+1), 1)
			for _, e := range kl.Daftar[DaftarDiItem(on+1, in+1, AnakEstimasi)] {
				if e["PrintFaceClaim"] != "1" { // S12.2.2.1
					continue
				}
				umum := map[string]string{"CauseOfLoss": cd("CauseOfLoss"), "CauseOfLossID": cd("CauseOfLossID"),
					"PersenRNM": persen, "ObjectID": ob["ObjectID"], "Type": StsOSTolak, "NoClaim": cd("NoClaim"),
					"EstimationDate": tgl, "Value": e["EstimationValue"], "GrossValue": e["GrossEstimationPct"],
					"Currency": e["Currency"], "CurrencyID": e["CurrencyID"]}
				for k, v := range umum {
					temp[k] = v
				}
				v := kl.Nilai
				if IsFire(v) { // S9
					temp["CoverageID"], temp["CoverageName"] = it["CoverageID"], it["CoverageNote"]
					temp["ObjectName"], temp["ObjectItemName"] = ob["ObjectName"], it["ObjectItemName"]
				}
				if IsAneka(v) { // S10
					temp["CoverageName"], temp["ObjectName"] = cov["CoverageNote"], ob["ObjectName"]
					temp["ObjectItemName"] = item1["ObjectItemName"]
				}
				if IsGolfInsurance(v) { // S11
					temp["CoverageName"], temp["ObjectName"] = cov["CoverageNote"], ob["ObjectName"]
					temp["ObjectItemName"] = it["ObjectItemName"]
				}
				if IsMarineCargo(v) { // S12
					temp["TradingName"], temp["PackingName"] = ob["ObjectName"], ob["ObjectLocation"]
					temp["Goodsname"], temp["ConveyanceName"] = ob["ObjectJob"], ob["ObjectSurveyLocation"]
					temp["CoverageName"] = cov["CoverageNote"]
				}
				if IsMBU(v) { // S13
					temp["BrandID"], temp["BrandName"], temp["Model"] = ob["Brand"], ob["BrandName"], ob["Model"]
					temp["ModelName"], temp["TypeName"], temp["LicensePlate"] = ob["ModelName"], ob["TypeName"], ob["LicensePlate"]
					temp["TypeID"], temp["CoverageID"], temp["CoverageName"] = ob["Type"], it["CoverageID"], it["CoverageNote"]
					temp["GrossValue"] = e["GrossEstimationPctMBU"]
				}
				if IsTravel(v) { // S14
					temp["ObjectName"], temp["ObjectIDCard"] = ob["ObjectName"], ob["ObjectIDCard"]
					temp["ObjectParticipantStatus"] = ob["ObjectParticipantStatus"]
					temp["CoverageID"], temp["CoverageName"] = it["CoverageID"], it["CoverageNote"]
				}
				if IsPA(v) { // S15
					temp["ObjectName"], temp["ObjectJob"], temp["Gender"] = ob["ObjectName"], ob["ObjectJob"], ob["Gender"]
					temp["ObjectIDCard"], temp["ObjectParticipantStatus"] = ob["ObjectIDCard"], ob["ObjectParticipantStatus"]
					temp["CoverageID"], temp["CoverageName"] = it["CoverageID"], it["CoverageNote"]
				}
				konversi := "" // S16-S17 CARI16: Value kosong -> "1"
				if strings.TrimSpace(temp["Value"]) == "" {
					konversi = "1"
				}
				dla := StsDLAFac // S16 CARI17 @if(objek.PlaStatus == "1", 8, 7)
				if ob["PlaStatus"] == "1" {
					dla = StsDLARetro
				}
				salin := make(map[string]string, len(temp))
				for k, x := range temp {
					salin[k] = x
				}
				out = append(out, BarisOS{CaseID: KunciInstans(kl.Nilai["pyID"]), NoClaim: cd("NoClaim"), NoPolis: noPolis(kl),
					StsReject: StsOSTolak, StsKonversi: konversi, StsDLA: dla, DataJSON: JSONHalamanPega(HalamanJSON{Nilai: salin})})
			}
		}
	}
	return out
}

// SusunOSTutup = KomitePost_CloseClaim S12.1-S12.3 (AcceptStatus 1): halaman InputParamOs (CauseOfLoss, CauseOfLossID,
// NoClaim; `IDMasterTreaty` = TreatyInMaster.ID tanpa penulis di Fac In - kosong), STS 4, NOPOLIS
// `ClaimData.PolicyData.PolicyNo` (= salinan `OfferFacIn.PolicyData` InsertJsonClaimNonMBU_act S1), STS_DLA / STS_KONVERSI
// tidak diisi (CARI16 / CARI17 tidak disetel). Bentuk sama dengan OS CloseClaim Claim Fac In tahap 1 (DEV STS 4).
func SusunOSTutup(kl kontrak.KlaimFacIn) BarisOS {
	cd := func(p string) string { return kl.Nilai["ClaimData."+p] }
	p := map[string]string{"pxObjClass": KelasOSAkseptasi, "CauseOfLoss": cd("CauseOfLoss"),
		"CauseOfLossID": cd("CauseOfLossID"), "NoClaim": cd("NoClaim")}
	return BarisOS{CaseID: KunciInstans(kl.Nilai["pyID"]), NoClaim: cd("NoClaim"), NoPolis: noPolis(kl),
		StsReject: StsOSFinal, DataJSON: JSONHalamanPega(HalamanJSON{Nilai: p})}
}
