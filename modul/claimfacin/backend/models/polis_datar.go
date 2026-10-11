package models

// Untuk apa berkas ini: POLIS DARI TABEL DATAR - perintah work owner 10-10-2026 ("pertama cek json_polis, jika json nya
// kosong ambil dari flat nya"). Polis NB Fac In sistem baru menulis JSON_POLIS hanya berisi penunjuk kasus
// (DEV 10-10-2026: `{"<kunci>":"NB-<n>"}`, 27 karakter) - isi polisnya di tabel datar nbfacin (`T_GENERAL_POLIS` =
// `T_WORK_POLIS.ID` lalu pohon PARENT_ID / SRC_PATH: QuotationData, CedingCoList, LocationList > Property >
// PropertyItemList / RiskLocation > OccupationList > AnekaList > CoverageList > SpreadingList, ...). Berkas ini menyusun
// ulang pohon itu menjadi bentuk DATA_JSON Pega (peta tabel / kolom / jalur: `polis_datar_peta.go`, disalin dari loader
// nbfacin) supaya `TerapkanPolis` dan seluruh port activity membaca halaman `OfferFacIn` yang sama dengan polis Pega.
//
// Medan yang TIDAK punya kolom datar diisi dari sumber yang sama dengan layar / rule NB Fac In (`LengkapiPolisDatar`):
//   - Quotation (`TambahanPolisDatar`, dibaca repository): BusinessCode / BusinessName / BusinessOldId / GroupPanel =
//     master BUSINESS; BusinessType = decision table `BusinessType_DeT` (SetBusinessType_Act 2-3); Insured dan Business
//     Status = opportunity ("New Business" = kode 1, Section Periode); Marketing = MARKETINGOFFICER menurut MOID;
//   - indeks baris = port DT NB `SetIndexCoverage_DT` (PropertyItemList / CoverageList Fire) dan `SetIndexObject_DT`
//     (AnekaList / CoverageList-nya, RiskLocation.OccupationList);
//   - CurrencyID PropertyItemList dan `Currency.ID` halaman Currency = CURRENCY.ID menurut kode (`idMataUang`).

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// KolomPolisDatar - satu kolom tabel datar dan medan dokumen polisnya (bertitik = halaman bersarang).
type KolomPolisDatar struct {
	Nama, Medan string
}

// TabelPolisDatar - satu tabel datar polis (rancangan; kolom yang benar-benar ada = katalog basis data).
type TabelPolisDatar struct {
	Daftar bool // ber-SEQ_NO: PageList; salah = Page tunggal (satu baris per induk)
	Src    bool // ber-SRC_PATH / PARENT_TABLE: induk dan jalur menurut baris
	Kolom  []KolomPolisDatar
}

// SisiPolisDatar - tabel anak, tabel induknya, dan jalur dokumen (lembar Jalur Sumber).
type SisiPolisDatar struct {
	Tabel, Induk string
	Jalur        []string
}

// BarisPolisDatar - satu baris tabel datar: ID, PARENT_ID, SRC_PATH, nilai per medan (kosong tidak dimuat).
type BarisPolisDatar struct {
	ID, Induk, Src string
	Nilai          map[string]string
}

// TabelAkarPolisDatar - akar dokumen polis (`OfferFacIn`), satu baris per kasus.
const TabelAkarPolisDatar = "T_GENERAL_POLIS"

// TabelPolisDatarDari - spesifikasi rancangan tabel datar `t`.
func TabelPolisDatarDari(t string) (TabelPolisDatar, bool) {
	s, ada := tabelPolisDatar[t]
	return s, ada
}

// NamaTabelPolisDatar - akar dan tabel yang dipakai klaim, urut nama (bahan saringan katalog).
func NamaTabelPolisDatar() []string {
	out := []string{TabelAkarPolisDatar}
	for t := range tabelPolisDatar {
		if tabelDipakaiKlaim[t] {
			out = append(out, t)
		}
	}
	sort.Strings(out)
	return out
}

// tabelDipakaiKlaim - tabel datar yang halamannya DIBACA port activity Claim Fac In (CopyNB / InsertObjects_dt /
// estimasi / DLA / spreading) dan tampilan polis NB hanya-baca (sub-tab Object Address / Surrounding Risk / Occupation /
// FEA / Loss Record, grid Deductible - perintah work owner 10-10-2026 "IKUTI MODUL NB FACIN"; tab per lini 11-10-2026):
// sisi lain (mis. salinan FacRetro T_FR_*) tidak dibaca. Pengukuran DEV 10-10-2026: setiap sisi menyapu tabel
// bersegmen besar (ratusan MB, 2-4 baris) - sisi yang tidak dipakai = waktu terbuang.
var tabelDipakaiKlaim = map[string]bool{
	"T_QUOTATIONDATA": true, "T_CEDINGCOLIST": true, "T_CEDINGCEDANTLIST": true, "T_FACRETROLIST": true,
	"T_LOCATIONLIST": true, "T_PROPERTY": true, "T_RISKLOCATION": true, "T_OCCUPATIONLIST": true, "T_ANEKALIST": true,
	"T_VEHICLEHE": true, "T_PROPERTYITEMLIST": true, "T_COVERAGELIST": true, "T_SPREADINGLIST": true,
	"T_CARGOLIST": true, "T_VEHICLELIST": true, "T_PERSONLIST": true, "T_CURRENCY": true,
	// tampilan polis NB (`tabelTampilanNB`)
	"T_BUILDINGCONSTRUCTION": true, "T_SURROUNDINGRISK": true, "T_LISTCAUSEOFLOSS": true, "T_COINSDATA": true,
	"T_TABLEOFLIMIT": true, "T_FEALIST": true, "T_DEDUCTIBLELIST": true,
	// form objek Aneka tampilan polis NB (work owner 11-10-2026 "tampilan view polis yang belum tergenerate")
	"T_MARINEHULL": true, "T_AVIATIONHULL": true, "T_MEMORANDUMACCEPT": true,
	// tab Object / Coverage Marine Cargo tampilan polis NB
	"T_SHIP": true, "T_LC": true, "T_SURVEYAGENT": true, "T_ADDITIONALSHIP": true, "T_ADDITIONALCOVERAGE": true,
	// tab Premium Deduction / Payment (nilai tersimpan, keputusan work owner 11-10-2026)
	"T_CURRENCYLIST": true, "T_LISTINSTALLMENT": true,
	"T_CLAUSELIST": true, // tab Clauses FIRE
	// tab Scoring Risk FIRE (nilai tersimpan)
	"T_SCORINGRISK": true, "T_SCORINGRESULT": true, "T_DATASCORINGRISKLIST": true, "T_SCORING_FACTOR": true,
	"T_SCORING_OPTION": true,
}

// tabelTampilanNB - tabel datar yang HANYA dibaca tampilan polis NB hanya-baca (tidak dibaca activity klaim): dimuat
// HANYA bila katalog menemukan indeks induknya (`PARENT_ID` / `PARENT_TABLE, PARENT_ID`). Pengukuran DEV 10-10-2026:
// T_LISTCAUSEOFLOSS ~1 juta baris tanpa indeks PARENT_ID = ~10 detik per kueri (sub-tab Loss Record), T_FEALIST tanpa
// indeks; begitu indeksnya dibuat (DBA), tabel itu ikut dimuat tanpa perubahan kode.
var tabelTampilanNB = map[string]bool{
	"T_BUILDINGCONSTRUCTION": true, "T_SURROUNDINGRISK": true, "T_LISTCAUSEOFLOSS": true, "T_COINSDATA": true,
	"T_TABLEOFLIMIT": true, "T_FEALIST": true, "T_DEDUCTIBLELIST": true,
	"T_MARINEHULL": true, "T_AVIATIONHULL": true, "T_MEMORANDUMACCEPT": true,
	"T_SHIP": true, "T_LC": true, "T_SURVEYAGENT": true, "T_ADDITIONALSHIP": true, "T_ADDITIONALCOVERAGE": true,
	"T_CURRENCYLIST": true, "T_LISTINSTALLMENT": true, "T_CLAUSELIST": true,
	"T_SCORINGRISK": true, "T_SCORINGRESULT": true, "T_DATASCORINGRISKLIST": true, "T_SCORING_FACTOR": true,
	"T_SCORING_OPTION": true,
}

// TabelTampilanNB - `tabelTampilanNB`: tabel ini dimuat hanya bila berindeks induk.
func TabelTampilanNB(t string) bool { return tabelTampilanNB[t] }

// SisiPolisDatarSemua - sisi induk-anak tabel yang dipakai klaim, urut topologis (induk sebelum anak).
func SisiPolisDatarSemua() []SisiPolisDatar {
	urut := urutTabelDatar()
	var out []SisiPolisDatar
	for _, s := range sisiPolisDatar {
		if tabelDipakaiKlaim[s.Tabel] && (s.Induk == TabelAkarPolisDatar || tabelDipakaiKlaim[s.Induk]) {
			out = append(out, s)
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return urut[out[i].Tabel] < urut[out[j].Tabel] })
	return out
}

// urutTabelDatar - kedalaman setiap tabel dari akar (anak selalu lebih dalam dari setiap induknya).
func urutTabelDatar() map[string]int {
	d := map[string]int{TabelAkarPolisDatar: 0}
	for ubah := true; ubah; {
		ubah = false
		for _, s := range sisiPolisDatar {
			if di, ada := d[s.Induk]; ada && d[s.Tabel] <= di {
				d[s.Tabel], ubah = di+1, true
			}
		}
	}
	return d
}

// DokumenPolisBerisi - DATA_JSON memuat isi polis. Salah hanya untuk dokumen kosong atau objek JSON tanpa halaman polis
// (penunjuk kasus sistem baru); teks yang tidak terurai dianggap berisi supaya galat uraiannya tetap terlihat
// (`TerapkanPolis`).
func DokumenPolisBerisi(dok []byte) bool {
	if strings.TrimSpace(string(dok)) == "" {
		return false
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(dok, &m); err != nil {
		return true
	}
	for _, k := range []string{"PolicyData", "QuotationData", "LocationList"} {
		if _, ada := m[k]; ada {
			return true
		}
	}
	return false
}

// polaKasusDatar - pengenal kasus work Pega / sistem baru (`NB-184351`, `EDM-12`).
var polaKasusDatar = regexp.MustCompile(`^[A-Z]{2,5}-[0-9]+$`)

// KasusPenunjukPolis - ID kasus di dokumen penunjuk JSON_POLIS sistem baru (nilai teks pertama, urut kunci, yang berbentuk
// pengenal kasus); kosong bila tidak ada.
func KasusPenunjukPolis(dok []byte) string {
	var m map[string]any
	if json.Unmarshal(dok, &m) != nil {
		return ""
	}
	for _, k := range kunciUrut(m) {
		if s, ok := m[k].(string); ok && polaKasusDatar.MatchString(strings.TrimSpace(s)) {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

// simpulDatar - baris yang sudah ditempatkan di dokumen.
type simpulDatar struct {
	jalur string
	m     map[string]any
}

// SusunPolisDatar menyusun dokumen polis dari baris akar `akar` (T_GENERAL_POLIS) dan pembaca anak `baca` (baris tabel
// sisi `s` ber-PARENT_ID di `induk`, urut SEQ_NO; tabel yang tidak ada di basis data = tanpa baris).
func SusunPolisDatar(akar BarisPolisDatar, baca func(s SisiPolisDatar, induk []string) ([]BarisPolisDatar, error)) (
	map[string]any, error) {
	doc := map[string]any{}
	isiSimpul(doc, akar.Nilai)
	simpul := map[string]map[string]simpulDatar{TabelAkarPolisDatar: {akar.ID: {jalur: "", m: doc}}}
	for _, s := range SisiPolisDatarSemua() {
		induk := simpul[s.Induk]
		if len(induk) == 0 {
			continue
		}
		spec := tabelPolisDatar[s.Tabel]
		ids := make([]string, 0, len(induk))
		for id := range induk {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		rows, err := baca(s, ids)
		if err != nil {
			return nil, err
		}
		for _, b := range rows {
			p, ada := induk[b.Induk]
			if !ada {
				continue
			}
			jalur := s.Jalur[0]
			if spec.Src && b.Src != "" {
				jalur = b.Src
			}
			m := map[string]any{}
			isiSimpul(m, b.Nilai)
			if !tempatkan(p.m, relatifDatar(jalur, p.jalur), m, spec.Daftar) {
				continue
			}
			if simpul[s.Tabel] == nil {
				simpul[s.Tabel] = map[string]simpulDatar{}
			}
			simpul[s.Tabel][b.ID] = simpulDatar{jalur: jalur, m: m}
		}
	}
	return doc, nil
}

// relatifDatar - jalur baris relatif terhadap jalur induknya; jalur yang tidak berawalan induk memakai ruas terakhir.
func relatifDatar(jalur, induk string) string {
	if induk == "" {
		return jalur
	}
	if strings.HasPrefix(jalur, induk+"/") {
		return strings.TrimPrefix(jalur, induk+"/")
	}
	return jalur[strings.LastIndex(jalur, "/")+1:]
}

// tempatkan memasang simpul `m` di `induk` pada jalur relatif `rel` (ruas antara = halaman); daftar = larik. Ruas
// antara yang sudah berisi selain halaman (mis. larik) tidak ditimpa - barisnya dilewati.
func tempatkan(induk map[string]any, rel string, m map[string]any, daftar bool) bool {
	ruas := strings.Split(rel, "/")
	for _, r := range ruas[:len(ruas)-1] {
		v, ada := induk[r]
		if !ada {
			anak := map[string]any{}
			induk[r] = anak
			induk = anak
			continue
		}
		anak, ok := v.(map[string]any)
		if !ok {
			return false
		}
		induk = anak
	}
	k := ruas[len(ruas)-1]
	if daftar {
		v, ada := induk[k]
		xs, ok := v.([]any)
		if ada && !ok {
			return false
		}
		induk[k] = append(xs, m)
		return true
	}
	if _, ada := induk[k]; ada { // halaman tunggal: baris pertama yang dipakai
		return false
	}
	induk[k] = m
	return true
}

// isiSimpul menulis nilai per medan (medan bertitik = halaman bersarang).
func isiSimpul(m map[string]any, nilai map[string]string) {
	for _, medan := range kunciUrutTeks(nilai) {
		if v := nilai[medan]; v != "" {
			isiBila(m, medan, v)
		}
	}
}

func kunciUrutTeks(m map[string]string) []string {
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

// TambahanPolisDatar - medan QuotationData tanpa kolom datar (sumber: lihat kepala berkas).
type TambahanPolisDatar struct {
	BusinessCode, BusinessName, BusinessOldId, GroupPanel string // BUSINESS (ID, NOTE, OLDID, GROUPPANEL)
	InsuredID, InsuredName, StatusBusiness                string // T_NB_OPPORTUNITY (+ T_M_ACCOUNT.INSUREDNAME)
	MarketingCode, MarketingName, TeamGroup               string // MARKETINGOFFICER menurut MOID (CLIENTID, CLIENTNAME, TEAMGROUP)
	BranchCode, BranchName                                string // MARKETINGOFFICER BRANCHDETAILID / BRANCHDETAILNAME
}

// kodeStatusBisnis - label Business Status opportunity NB -> kode `QuotationData.StatusBusiness` Pega (Section
// `Periode`: `.QuotationData.StatusBusiness = 1` bernilai "New Business"; satu-satunya opsi opportunity NB).
var kodeStatusBisnis = map[string]string{"New Business": "1"}

// LengkapiPolisDatar mengisi medan yang tidak berkolom datar (hanya bila kosong): PolicyData.PolicyNo = nomor polis
// JSON_POLIS, medan Quotation `t`, BusinessType (`JenisBisnisNB`), indeks baris (SetIndexCoverage_DT /
// SetIndexObject_DT), CurrencyID / Currency.ID menurut `idMataUang` (kode -> CURRENCY.ID).
func LengkapiPolisDatar(doc map[string]any, nopolis string, t TambahanPolisDatar, idMataUang map[string]string) {
	isiBila(doc, "PolicyData.PolicyNo", nopolis)
	q, ok := doc["QuotationData"].(map[string]any)
	if !ok {
		q = map[string]any{}
		doc["QuotationData"] = q
	}
	status := t.StatusBusiness
	if k, ada := kodeStatusBisnis[status]; ada {
		status = k
	}
	for _, kv := range [][2]string{{"BusinessCode", t.BusinessCode}, {"BusinessName", t.BusinessName},
		{"BusinessOldId", t.BusinessOldId}, {"GroupPanel", t.GroupPanel}, {"InsuredID", t.InsuredID},
		{"InsuredName", t.InsuredName}, {"StatusBusiness", status}, {"MarketingCode", t.MarketingCode},
		{"MarketingName", t.MarketingName}, {"TeamGroup", t.TeamGroup}, {"BranchCode", t.BranchCode},
		{"BranchName", t.BranchName}} {
		isiBila(q, kv[0], kv[1])
	}
	if teksDok(q["BusinessType"]) == "" { // SetBusinessType_Act 2-3
		jenis := JenisBisnisNB(teksDok(q["GroupPanel"]), teksDok(q["BusinessOldId"]))
		q["BusinessType"] = jenis
		q["BusinessType2"] = map[bool]string{true: "LIFE", false: "NONLIFE"}[jenis == "Life"]
	}
	for l, lok := range petaDaftar(doc, "LocationList") {
		prop, _ := lok["Property"].(map[string]any)
		ls := strconv.Itoa(l + 1)
		for pi, it := range petaDaftar(prop, "PropertyItemList") { // SetIndexCoverage_DT 1.2
			ps := strconv.Itoa(pi + 1)
			isiBila(it, "IndexProperty", ls)
			isiBila(it, "IndexPropertyItem", ps)
			isiBila(it, "CurrencyID", idMataUang[teksDok(it["Currency"])])
			for ci, c := range petaDaftar(it, "CoverageList") { // 1.2.4
				isiBila(c, "IndexCoverage", strconv.Itoa(ci+1))
				isiBila(c, "IndexProperty", ls)
				isiBila(c, "IndexPropertyItem", ps)
			}
		}
		risk, _ := prop["RiskLocation"].(map[string]any)
		for oi, okp := range petaDaftar(risk, "OccupationList") { // SetIndexObject_DT 1.1.2
			os := strconv.Itoa(oi + 1)
			for ai, a := range petaDaftar(okp, "AnekaList") { // 1.1.2.2
				as := strconv.Itoa(ai + 1)
				isiBila(a, "IdxLocation", ls)
				isiBila(a, "IdxOccupation", os)
				isiBila(a, "IdxAneka", as)
				for _, c := range petaDaftar(a, "CoverageList") { // 1.1.2.2.5
					isiBila(c, "IndexLocation", ls)
					isiBila(c, "IndexOccupation", os)
					isiBila(c, "IndexAneka", as)
				}
			}
		}
	}
	jelajahHalaman(doc, func(m map[string]any) { // halaman Currency: ID menurut Name
		if cur, ok := m["Currency"].(map[string]any); ok {
			isiBila(cur, "ID", idMataUang[teksDok(cur["Name"])])
		}
	})
}

// KodeMataUangDatar - kode mata uang dokumen: PropertyItemList.Currency dan Name setiap halaman Currency (bahan
// `idMataUang` LengkapiPolisDatar).
func KodeMataUangDatar(doc map[string]any) []string {
	lihat := map[string]bool{}
	var out []string
	tambah := func(c string) {
		if c != "" && c != KodeMataUangTakDikenal && !lihat[c] {
			lihat[c] = true
			out = append(out, c)
		}
	}
	jelajahHalaman(doc, func(m map[string]any) {
		if cur, ok := m["Currency"].(map[string]any); ok {
			tambah(teksDok(cur["Name"]))
		}
		for _, it := range petaDaftar(m, "PropertyItemList") {
			tambah(teksDok(it["Currency"]))
		}
	})
	sort.Strings(out)
	return out
}

// KodeMataUangTakDikenal - DEFAULT kolom CURRENCY / CURRENCY_CODE tabel datar nbfacin (A169, sentinel K-069): dibaca
// sebagai kosong.
const KodeMataUangTakDikenal = "UNKNOWN"

// MOIDDatar - QuotationData.MOID dokumen (kunci MARKETINGOFFICER).
func MOIDDatar(doc map[string]any) string {
	q, _ := doc["QuotationData"].(map[string]any)
	return teksDok(q["MOID"])
}

// petaDaftar - baris halaman larik `m[nama]`.
func petaDaftar(m map[string]any, nama string) []map[string]any {
	xs, _ := m[nama].([]any)
	out := make([]map[string]any, 0, len(xs))
	for _, x := range xs {
		if r, ok := x.(map[string]any); ok {
			out = append(out, r)
		}
	}
	return out
}

func isiBila(m map[string]any, medan, v string) {
	if v == "" || m == nil {
		return
	}
	ruas := strings.Split(medan, ".")
	cur := m
	for _, r := range ruas[:len(ruas)-1] {
		anak, ok := cur[r].(map[string]any)
		if !ok {
			if _, ada := cur[r]; ada {
				return
			}
			anak = map[string]any{}
			cur[r] = anak
		}
		cur = anak
	}
	if teksDok(cur[ruas[len(ruas)-1]]) == "" {
		cur[ruas[len(ruas)-1]] = v
	}
}

func teksDok(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}

// jelajahHalaman memanggil `f` untuk setiap halaman (objek) dokumen.
func jelajahHalaman(v any, f func(m map[string]any)) {
	switch x := v.(type) {
	case map[string]any:
		f(x)
		for _, k := range kunciUrut(x) {
			jelajahHalaman(x[k], f)
		}
	case []any:
		for _, e := range x {
			jelajahHalaman(e, f)
		}
	}
}

// ---------------------------------------------------------------- BusinessType_DeT

// barisJenisBisnis - satu baris decision table `NB FacIn/DecisionTable/BusinessType_DeT.xml`: GroupPanel `=`, daftar
// BusinessOldId (kosong = semua; sel + pyOrConditions baris itu), hasil.
type barisJenisBisnis struct {
	panel string
	lama  []string
	jenis string
}

// tabelJenisBisnis - BusinessType_DeT VERBATIM (36 baris, urut; pyRowNum OR-condition berbasis 0).
var tabelJenisBisnis = []barisJenisBisnis{
	{"001", nil, "Medicare"},
	{"002", []string{"03"}, "PA"},
	{"003", strings.Fields("36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 C2 C3 C4 C5 C6 C7 C8 C9 D1 D2 D3 D4 D5 D6 D7 D8 D9 F1 F2 F3 F4 F5"), "Bonding"},
	{"003", strings.Fields("95 96 97 98 99 A1 A2 A3 A4 A5 A6 A7 A8 A9 B4 B5 B6 B8 B9 C1"), "BondingKBG"},
	{"003", []string{"20"}, "HE"},
	{"003", []string{"05", "55"}, "MarineHull"},
	{"003", []string{"21"}, "Glass"},
	{"003", strings.Fields("15 24 89 90 91 E7 E9"), "Liability"},
	{"003", []string{"16"}, "AllRisk"},
	{"003", []string{"06"}, "AviationHull"},
	{"003", []string{"14"}, "Burglary"},
	{"003", []string{"07"}, "Car"},
	{"003", []string{"08"}, "Ear"},
	{"003", []string{"09"}, "ElectronicEquipment"},
	{"003", []string{"19"}, "Fidelity"},
	{"003", []string{"27"}, "GolfInsurance"},
	{"003", []string{"17"}, "CIT"},
	{"003", []string{"18"}, "CIS"},
	{"003", []string{"11"}, "Boiler"},
	{"003", []string{"93"}, "Workmen"},
	{"003", []string{"10"}, "MBD"},
	{"003", []string{"12"}, "ContractorsPM"},
	{"003", strings.Fields("79 80 81 82 83 84 85 86 87 88"), "CustomBond"},
	{"003", []string{"28"}, "FireStyle1"},
	{"003", []string{"B3"}, "LandRig"},
	{"003", nil, "Aneka"},
	{"004", nil, "MarineCargo"},
	{"005", nil, "Travel"},
	{"006", []string{"78"}, "OilGas"},
	{"006", []string{"22", "26", "29", "57"}, "FireStyle1"},
	{"006", []string{"01", "34", "59"}, "FireStyle2"},
	{"006", nil, "Fire"},
	{"007", []string{"02"}, "MBUCar"},
	{"007", []string{"60"}, "MBUMotorCycle"},
	{"009", strings.Fields("L1 L2 L3 L4 L5 L6 L7 L8 L9 L10 L11 L12 L13 L14 L15 L16 L17 L18 L19 L20 L21"), "Life"},
	{"008", nil, "Medicare"},
}

// JenisBisnisNB = BusinessType_DeT (`Quotation.GroupPanel`, `Quotation.BusinessOldId`): baris pertama yang cocok,
// bawaan "UNKNOWN" (pyDefaultResult).
func JenisBisnisNB(groupPanel, businessOldID string) string {
	for _, b := range tabelJenisBisnis {
		if b.panel != groupPanel {
			continue
		}
		if len(b.lama) == 0 {
			return b.jenis
		}
		for _, l := range b.lama {
			if l == businessOldID {
				return b.jenis
			}
		}
	}
	return "UNKNOWN"
}
