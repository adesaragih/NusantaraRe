package models

// Untuk apa berkas ini: PRA-PROSES flow action `ViewTransferDtl` - `SetValueKomite` (S1-S13) dan bagian bacaan
// `ApprovalKomite_Act` (S1-S3): adjustment kasus komite ini (`DataTempAdj`), riwayat adjustment seluruh klaim
// (`TempDataAdj`), total per mata uang + "Total in IDR" (`TempTotalAdj`), total adjustment IDR penentu tangga
// (`Local.TotalAdj`) dan penanda Fac Retro. Murni; perluasan tangga ada di `tangga.go`.
//
// TT3 / TT4: `SetValueKomite` S1 `.TransferType != 2 [T=6]` keluar - nol pra-proses (Section tanpa baris grid).

import (
	"fmt"
	"strconv"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
)

// Jalur pohon halaman klaim induk (ejaan halaman Claim Fac In).
const (
	DaftarObjek  = "ClaimData.ObjectList"
	AnakItem     = "ObjectItemList"
	AnakAdj      = "Adjustment"
	AnakSpread   = "SpreadingAdjustment"
	AnakQS       = "SpreadingQuotaShare"
	AnakSpreadPL = "SpreadingList"
	AnakEstimasi = "EstimationList"
	AnakCoverage = "CoverageList"
	// DaftarRetroPolis - `pyWorkCover.OfferFacIn.FacRetroList` (PreShowRetro_Act S2.2, PreSecurityReas_Act S2.2).
	DaftarRetroPolis = "OfferFacIn.FacRetroList"
)

// TreatyRetro - `.TreatyType == "10015"` (SaveAcceptation_KMT S2.1, DetailAdjustmentFac "View Retro").
const TreatyRetro = "10015"

func anak(induk string, n int, nama string) string { return fmt.Sprintf("%s(%d).%s", induk, n, nama) }

// DaftarItem / DaftarAdj / DaftarDiItem / DaftarDiAdj - jalur daftar pohon objek -> item -> adjustment.
func DaftarItem(o int) string                     { return anak(DaftarObjek, o, AnakItem) }
func DaftarAdj(o, i int) string                   { return anak(DaftarItem(o), i, AnakAdj) }
func DaftarDiItem(o, i int, nama string) string   { return anak(DaftarItem(o), i, nama) }
func DaftarDiAdj(o, i, a int, nama string) string { return anak(DaftarAdj(o, i), a, nama) }
func baris(kl kontrak.KlaimFacIn, daftar string, n int) map[string]string {
	rows := kl.Daftar[daftar]
	if n < 1 || n > len(rows) {
		return map[string]string{}
	}
	return rows[n-1]
}

// Objek / Item / Adjustment - baris pohon kasus komite ini (`IndexObject` / `IndexObjectItem` / `IndexAdjustment`).
func Objek(kl kontrak.KlaimFacIn) map[string]string { return baris(kl, DaftarObjek, kl.Objek) }
func Item(kl kontrak.KlaimFacIn) map[string]string  { return baris(kl, DaftarItem(kl.Objek), kl.Item) }
func Adjustment(kl kontrak.KlaimFacIn) map[string]string {
	return baris(kl, DaftarAdj(kl.Objek, kl.Item), kl.Adjustment)
}

// salinBaris - salinan baris (grid layar mengubahnya).
func salinBaris(b map[string]string) map[string]string {
	m := make(map[string]string, len(b))
	for k, v := range b {
		m[k] = v
	}
	return m
}

// PraProses - hasil `SetValueKomite` + bacaan `ApprovalKomite_Act` S1-S3.
type PraProses struct {
	// AdjKomite - `DataTempAdj.Adjustment` (S7.1.1 / S8.3): adjustment ber-KomiteNo kasus ini yang belum diakseptasi,
	// `pyNote` = "ObjectItem ke i, Adjustment ke a".
	AdjKomite []map[string]string
	// Posisi - (item, adjustment) setiap baris AdjKomite (rincian DetailAdjustmentFac).
	Posisi [][2]int
	// Riwayat - `TempDataAdj.Adjustment` (S6): seluruh adjustment seluruh objek / item klaim, diubah S11.1.1.1-S11.1.1.3.
	Riwayat []map[string]string
	// Total - `TempTotalAdj.pxResults` (S9-S11, S13): Currency, AdjustmentGross, AdjustmentValue; baris akhir
	// "Total in IDR".
	Total []map[string]string
	// TotalAdj - `ApprovalKomite_Act` S2.1 `Local.TotalAdj` (Σ ValueAdjustment AdjKomite, IDR).
	TotalAdj *apd.Decimal
	// Retro - S2.2 `.IsFacRetro == 1` (S3 keluar: tangga tidak diperluas).
	Retro bool
	// Inisial - S8.2 / S12 `Komite.Initial` = `.pxCreateOpName` adjustment yang belum diakseptasi.
	Inisial string
	// Okupasi - S4 `TempObject.OccupationName` (item kasus komite).
	Okupasi string
}

// LabelTotalIDR - S13 baris total (`TempTotalAdj.pxResults(<APPEND>).Currency := Total in IDR`).
const LabelTotalIDR = "Total in IDR"

// SetValueKomite = SetValueKomite S1-S13 + ApprovalKomite_Act S1-S3 atas klaim induk `kl` kasus `k`. `kurs` = kurs IDR
// per CurrencyID bagi baris tanpa `.CurrencyDol` (Claim Fac In tidak menyimpan medan itu; sumbernya sama dengan
// `SetNilaiResikoSendiri` - CurrencyStandard).
func SetValueKomite(k Kasus, kl kontrak.KlaimFacIn, kurs map[string]string) (PraProses, error) {
	p := PraProses{TotalAdj: apd.New(0, 0)}
	if k.TransferType != TransferAdjustment { // S1
		return p, nil
	}
	o := kl.Objek
	p.Okupasi = Item(kl)["OccupationName"]   // S4
	for on := range kl.Daftar[DaftarObjek] { // S6
		for in := range kl.Daftar[DaftarItem(on+1)] {
			for _, b := range kl.Daftar[DaftarAdj(on+1, in+1)] {
				p.Riwayat = append(p.Riwayat, salinBaris(b))
			}
		}
	}
	for in := range kl.Daftar[DaftarItem(o)] { // S7 (objek IndexObject)
		for an, b := range kl.Daftar[DaftarAdj(o, in+1)] {
			if b["ID"] != k.AdjustmentID { // S7.1.1 `.KomiteNo == pyWorkPage.pyID` (satu adjustment per KMT, OQ-CFI-28)
				continue
			}
			if st := b["AcceptanceStatus"]; st != "" && st != KeputusanMenunggu { // S8.3
				continue
			}
			c := salinBaris(b)
			c["pyNote"] = "ObjectItem ke " + strconv.Itoa(in+1) + ", Adjustment ke " + strconv.Itoa(an+1)
			p.AdjKomite = append(p.AdjKomite, c)
			p.Posisi = append(p.Posisi, [2]int{in + 1, an + 1})
			p.Inisial = b["pxCreateOpName"] // S8.2
		}
	}
	var h hitung
	var mataUang []string // S9-S10 (unik, urut kemunculan)
	lihat := map[string]bool{}
	for _, b := range p.Riwayat {
		if !lihat[b["Currency"]] {
			lihat[b["Currency"]] = true
			mataUang = append(mataUang, b["Currency"])
		}
	}
	grossIDR, rnmIDR := apd.New(0, 0), apd.New(0, 0)
	for _, cur := range mataUang { // S11
		gross, rnm := apd.New(0, 0), apd.New(0, 0)
		for _, b := range p.Riwayat { // S11.1
			if b["AcceptanceStatus"] == KeputusanTolak || b["Currency"] != cur { // S11.1.1
				continue
			}
			if b["Currency"] == "" { // S11.1.1.1
				b["Currency"] = b["CurrencyEstimasi"]
			}
			switch b["PaymentType"] {
			case "4", "6": // S11.1.1.2 adjuster fee -> adjustment
				b["AdjustmentValue"] = b["AdjusterFeeValue"]
				b["GrossAdjustment"] = TeksAngka(h.tambah(h.dari("GrossAdjustment", b["GrossAdjustment"]),
					h.dari("VATValue", b["VATValue"])))
			case "3": // S11.1.1.3 salvage -> adjustment
				b["AdjustmentValue"] = b["SalvageValue"]
			}
			ks := b["CurrencyDol"]
			if ks == "" {
				ks = kurs[b["CurrencyID"]]
			}
			kd := h.dari("CurrencyDol", ks)
			g, r := h.dari("GrossAdjustment", b["GrossAdjustment"]), h.dari("AdjustmentValue", b["AdjustmentValue"])
			gross, rnm = h.tambah(gross, g), h.tambah(rnm, r) // S11.1.1.4
			grossIDR, rnmIDR = h.tambah(grossIDR, h.kali(g, kd)), h.tambah(rnmIDR, h.kali(r, kd))
		}
		p.Total = append(p.Total, map[string]string{"Currency": cur, "AdjustmentGross": TeksAngka(gross),
			"AdjustmentValue": TeksAngka(rnm)}) // S11.2
	}
	p.Total = append(p.Total, map[string]string{"Currency": LabelTotalIDR, "AdjustmentGross": TeksAngka(grossIDR),
		"AdjustmentValue": TeksAngka(rnmIDR)}) // S13
	for _, b := range p.AdjKomite { // ApprovalKomite_Act S2
		p.TotalAdj = h.tambah(p.TotalAdj, h.dari("ValueAdjustment", b["ValueAdjustment"])) // S2.1
		if b["IsFacRetro"] == "1" {                                                        // S2.2
			p.Retro = true
		}
	}
	return p, h.err
}
