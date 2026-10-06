// ⛔⛔ BERKAS INI TIDAK LAGI DIPAKAI JALUR BACA, sejak 6 Oktober 2026.
//
// Keputusan pemilik proses melarang keras menarik nilai dari `JSONDATA`.
// Keempat tab kini dibaca dari tabel pendaratan lewat `pendaratan_layer.go`;
// NOL kode produksi memanggil fungsi di berkas ini, hanya ujinya sendiri.
//
// ⚠️ Ia DITINGGALKAN satu ronde dengan sengaja, bukan karena terlupa:
// ketiga belas tabel pendaratan masih NOL BARIS, sehingga jalur barunya
// belum pernah terbukti di layar. Membuang jalur lama sebelum penggantinya
// terbukti berarti membuang satu-satunya pembanding ketika hasilnya
// berselisih.
//
// ⭐ Begitu tabelnya terisi dan layarnya terbukti, berkas ini beserta
// ujinya DIBUANG — dua sumber untuk satu tab berarti salah satunya akan
// basi tanpa suara, dan itu doktrin modul ini sendiri.

package repository

// POHON tab Limits cabang proporsional — `Limits[] → Detail[] → daftar`,
// apa adanya dari `M_TREATY_IN.JSONDATA`.
//
// ⛔ BACA SAJA. Berdampingan dengan `LayerDariDokumen` (bentuk PIPIH yang
// melayani Share, Event Limits, RNM Share), bukan menggantikannya.
//
// # Kenapa pohon, dan dari mana bentuknya
//
// `Section/TreatyInTabsProportional.xml` tab Limits (@524276) adalah grid
// `Kind of Treaty` (`.TreatyType`) ber-`pyEditingMode=expandPane` dan
// `pyEditAction=LimitProportional` (@604564). Baris yang dibuka menampilkan
// `Section/LimitProportional.xml` (kelas `TreatyInLimits`): `Treaty Type` +
// grid `.Detail` — yang JUGA expandPane, `pyEditAction=DetailLimits`
// (@87354). `Section/DetailLimits.xml` (kelas `TreatyInLimitsDetail`) memuat
// medan, empat grid, dan sebelas tab.
//
// Bentuk pipih tidak dapat menyatakan tiga tingkat itu — dan bentuk pipih
// itulah sebab layar sebelumnya menaruh medan LAYER non-prop (Cover, MDP,
// ROL …) di bawah `Kind of Treaty`, tempat Section prop tidak punya satu pun.
//
// # Daftar-izin per tingkat — kunci yang Section ikat, BUKAN seluruh dokumen
//
// Terukur 6 Oktober 2026 atas dokumen proporsional: 1.360 elemen `Limits[]`,
// 2.862 `Detail[]`; `TreatyTypeID` terisi 19, `TreatyGroupID` 2.861,
// `QSPct` 1.469, `Surplus` 1.390, `IOOLimitList` 2.870 elemen.

import (
	"bytes"
	"encoding/json"
)

// kunciLimit - medan tingkat `Kind of Treaty` (grid + LimitProportional).
var kunciLimit = []string{"TreatyType", "TreatyTypeID"}

// kunciDetail - medan skalar `DetailLimits.xml`.
var kunciDetail = []string{
	"TreatyGroup", "TreatyGroupID", "TreatyType", "QSPct", "Surplus", "RetentionPct", "CessionPct",
	"CurrencyRSMD", "RSMDLimit", "CurrencyEarthquake", "Earthquake", "CurrencyFloodJab", "FloodJab",
	"CurrencyFloodNat", "FloodNation", "RIOGR", "RIONR", "PremiumReservePct",
	"ProfitCommision", "ProfitME", "ProfitYDCF", "LowerBand", "UpperBand", "ReisuredParticipant", "Periode",
}

// larikDetail - grid `DetailLimits.xml`, beserta kunci tiap barisnya.
var larikDetail = map[string][]string{
	"COBList":            {"ClassOfBusiness"},
	"IOOLimitList":       {"Currency", "Value", "Layer", "Note"},
	"RetentionList":      {"Currency", "Value", "Layer", "Note"},
	"CessionList":        {"Currency", "Value"},
	"DeductionList":      {"Comment", "CurrencyID", "Currency", "Deduction", "DeductionPct"},
	"DeductionTotalList": {"Currency", "Value"},
	"ReserveList":        {"Currency", "Value"},
	"PLAList":            {"Currency", "Value"},
	"CashLossList":       {"Currency", "Value"},
	"ClaimCoopList":      {"Currency", "Value"},
	"EPIList":            {"Currency", "Value"},
	"AchievementLists": {
		"Quarter", "QUARTERYEAR", "Currency", "PREMIUM", "RICOMM", "BROKERAGE", "NETPREMIUM",
		"PaidClaim", "CASHCALL", "OutstandingClaim", "IncuredClaim", "Total", "LossRatio",
	},
}

// PohonLimitsDariDokumen membaca `Limits` satu dokumen menjadi pohon.
//
// Bentuk keluaran: tiap simpul `map[string]any` — nilai skalar TEKS apa
// adanya (angka dengan digit aslinya), larik `[]map[string]any`. Kunci yang
// TIDAK ADA di dokumen tidak dimasukkan; yang ada bernilai `null` masuk
// sebagai teks kosong.
func PohonLimitsDariDokumen(dok []byte) []map[string]any {
	var akar struct {
		Limits json.RawMessage `json:"Limits"`
	}
	if json.Unmarshal(dok, &akar) != nil {
		return []map[string]any{}
	}
	out := []map[string]any{}
	for _, l := range elemen(akar.Limits) {
		s := petikSkalar(l, kunciLimit)
		detail := []map[string]any{}
		for _, d := range elemen(l["Detail"]) {
			sd := petikSkalar(d, kunciDetail)
			for nama, kunci := range larikDetail {
				r, ada := d[nama]
				if !ada {
					continue
				}
				baris := []map[string]any{}
				for _, e := range elemen(r) {
					baris = append(baris, petikSkalar(e, kunci))
				}
				sd[nama] = baris
			}
			detail = append(detail, sd)
		}
		s["Detail"] = detail
		out = append(out, s)
	}
	return out
}

// elemen mengurai larik objek; yang bukan larik objek → nol elemen.
func elemen(r json.RawMessage) []map[string]json.RawMessage {
	var mentah []json.RawMessage
	if json.Unmarshal(r, &mentah) != nil {
		return nil
	}
	out := make([]map[string]json.RawMessage, 0, len(mentah))
	for _, m := range mentah {
		d := json.NewDecoder(bytes.NewReader(m))
		d.UseNumber()
		var o map[string]json.RawMessage
		if d.Decode(&o) == nil {
			out = append(out, o)
		}
	}
	return out
}

// petikSkalar membawa kunci daftar-izin yang ADA dan bernilai skalar.
func petikSkalar(o map[string]json.RawMessage, kunci []string) map[string]any {
	s := map[string]any{}
	for _, k := range kunci {
		r, ada := o[k]
		if !ada {
			continue
		}
		t := bytes.TrimSpace(r)
		switch {
		case len(t) == 0 || t[0] == '{' || t[0] == '[':
			continue
		case string(t) == "null":
			s[k] = ""
		case t[0] == '"':
			var v string
			if json.Unmarshal(t, &v) == nil {
				s[k] = v
			}
		default:
			s[k] = string(t)
		}
	}
	return s
}
