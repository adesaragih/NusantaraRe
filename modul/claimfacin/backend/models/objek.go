package models

// Untuk apa berkas ini: OBJEK KLAIM - calon objek dari polis (`InsertObjects_dt` langkah 8, pra-proses FlowAction
// InputRegister), objek terpilih (`InsertObjectItemList_DT`, pre-DT tombol Submit Input Register), dan pelengkapan medan
// salinan polis saat kasus dimuat.
//
// ⛔ Medan objek yang XML hanya SALIN dari polis (ObjectName, ObjectLocation, Brand.., PropertyItemList, OccupationList,
// AnekaList, CoverageList - 29 medan jenis objek STRUKTUR butir 11) TIDAK disimpan (prompt §6 butir 2): objek menyimpan
// `KunciPolis` (letak calonnya di halaman polis), dan setiap muat `LengkapiObjek` menyalin ulang dari polis. Yang
// disimpan hanya medan yang DITULIS klaim (CFS, PrintFaceClaim, PlaStatus, DLAStatus, IsFacretro, ...).
//
// ⚠️ `[penyimpangan sadar]` InsertObjects_dt 8.6.25: `FOR EACH LocationList -> ObjectList(<LAST>).OccupationList :=
// .Property.RiskLocation.OccupationList` menulis HANYA calon terakhir (setiap iterasi menimpa <LAST>), sehingga calon
// lain tanpa OccupationList. Calon Aneka lahir satu per lokasi (8.6.x.1.2 APPEND sekali per lokasi), maka maksudnya
// pasti: OccupationList lokasi calon itu sendiri. Diperbaiki (prompt §3 OQ-CFI-03), dicatat di PARITAS.

import (
	"strconv"
	"strings"
)

// Jalur objek.
const (
	DaftarObjek = CD + "ObjectList"
	// DaftarCalon - `TempClaimData.ClaimData.ObjectList` (grid pilih objek Input Register; halaman requestor Pega, tidak
	// disimpan - pilihan dibawa `Layar.Mode`).
	DaftarCalon = "TempClaimData.ClaimData.ObjectList"
	// JalurUjiCalon - `TempClaimData.ClaimData.Test` (CheckListEstimasi_Act: "1" bila ada objek terpilih).
	JalurUjiCalon = "TempClaimData.ClaimData.Test"
	AnakItem      = "ObjectItemList"
	// PropKunciPolis - letak calon objek di halaman polis ("LocationList(3)", "VehicleList(1)", ...).
	PropKunciPolis = "KunciPolis"
	PropDipilih    = "Selected"
)

// JalurObjek / JalurItem - jalur baris objek ke-o dan item ke-i.
func JalurObjek(o int) string   { return JalurBaris(DaftarObjek, o) }
func DaftarItem(o int) string   { return JalurAnak(DaftarObjek, o, AnakItem) }
func JalurItem(o, i int) string { return JalurBaris(DaftarItem(o), i) }
func DaftarDiItem(o, i int, anak string) string {
	return JalurAnak(DaftarItem(o), i, anak)
}

// Calon - satu calon objek polis.
type Calon struct {
	// Medan - medan objek yang diisi InsertObjects_dt (ObjectID, ObjectName, ...).
	Medan Baris
	// Kunci - letak di halaman polis (relatif `OfferFacIn.`).
	Kunci string
	// Daftar - anak daftar objek -> jalur daftar sumber di halaman polis (PropertyItemList, OccupationList, AnekaList,
	// CoverageList).
	Daftar map[string]string
	// CoverageItem - jalur daftar CoverageList calon untuk item pertama (`ObjectItemList(<APPEND>).CoverageList`
	// MBU 8.2.1.11 / MarineCargo 8.7.1.9 / Travel InsertObjectItemList_DT 2.2.4.7.1.1.8).
	CoverageItem string
}

func sumberPolis(rel string) string { return AwalanPolis + "." + rel }

// barisPolis - baris ke-n daftar polis beserta jalur daftarnya.
type barisPolis struct {
	jalur string
	n     int
	b     Baris
}

func daftarPolis(h *Halaman, rel string) []barisPolis {
	j := sumberPolis(rel)
	var out []barisPolis
	for i, b := range h.AmbilDaftar(j) {
		out = append(out, barisPolis{jalur: j, n: i + 1, b: b})
	}
	return out
}

func (p barisPolis) anak(prop string) string { return JalurAnak(p.jalur, p.n, prop) }

func (p barisPolis) kunci() string {
	return strings.TrimPrefix(JalurBaris(p.jalur, p.n), AwalanPolis+".")
}

// CalonObjek = InsertObjects_dt langkah 8 (`pyWorkPage.IsTreatyIn != 1`; polis dari Choose Polis selalu RNM-F, lihat
// `IsTreatyInDari`). Cabang 8.1-8.7 diuji BERURUTAN dan TIDAK eksklusif, persis DT (Aneka 8.6.24 `BusinessType ==
// "Aneka"` dapat menambah calon kedua bagi kasus yang juga lolos sub-cabang Aneka bersandi bisnis).
func CalonObjek(h *Halaman) []Calon {
	var out []Calon
	if IsFire(h) { // 8.1
		for _, p := range daftarPolis(h, "LocationList") {
			out = append(out, Calon{Kunci: p.kunci(), Medan: Baris{
				"ObjectName": p.b["Property.ObjectName"], "ObjectID": p.b["Property.ObjectNo"], "ObjectStatus": "",
				"ObjectLocation": p.b["Property.RiskLocation.ASMAddress"],
			}, Daftar: map[string]string{
				"PropertyItemList": p.anak("Property.PropertyItemList"),
				"OccupationList":   p.anak("Property.OccupationList"),
			}})
		}
	}
	if IsMBU(h) { // 8.2
		for _, p := range daftarPolis(h, "VehicleList") {
			m := Baris{}
			for _, k := range []string{"Brand", "BrandName", "Model", "ModelName", "Type", "TypeName", "LicensePlate",
				"ChassisNumber", "EngineNumber"} {
				m[k] = p.b[k]
			}
			out = append(out, Calon{Kunci: p.kunci(), Medan: m,
				Daftar: map[string]string{"CoverageList": p.anak("CoverageList")}, CoverageItem: p.anak("CoverageList")})
		}
	}
	if IsPA(h) { // 8.3
		for _, p := range daftarPolis(h, "PersonList") {
			out = append(out, Calon{Kunci: p.kunci(), Medan: Baris{
				"ObjectName": p.b["pyFullName"], "ObjectID": p.b["ASMIDCard"], "ObjectJob": p.b["ASMJobName"],
				"ObjectWeight": p.b["ASMWeight"], "ObjectHeight": p.b["ASMHeight"], "ObjectLeftHanded": p.b["ASMLeftHanded"],
				"ObjectDateOfBirth": p.b["ASMDateOfBirth"], "Gender": p.b["ASMGender"],
			}, Daftar: map[string]string{"CoverageList": p.anak("ASMCoverage")}})
		}
	}
	if IsTravel(h) { // 8.4
		for _, p := range daftarPolis(h, "PersonList") {
			out = append(out, Calon{Kunci: p.kunci(), Medan: Baris{
				"ObjectName": p.b["pyFullName"], "ObjectID": strconv.Itoa(p.n),
				"ObjectParticipantStatus": p.b["ASMParticipantStatus"], "ObjectIDCard": p.b["ASMIDCard"],
				"ObjectDateOfBirth": p.b["ASMDateOfBirth"],
			}, Daftar: map[string]string{"CoverageList": p.anak("ASMCoverage")}, CoverageItem: p.anak("ASMCoverage")})
		}
	}
	if IsGolfInsurance(h) { // 8.5
		for _, p := range daftarPolis(h, "LocationList") {
			aneka := p.anak("Property.RiskLocation.AnekaList")
			c := Calon{Kunci: p.kunci(), Medan: Baris{
				"ObjectID": p.b["Property.ObjectNo"], "ObjectName": "OTHERS", "ObjectIDCard": strconv.Itoa(p.n),
				"ObjectLocation": p.b["Property.RiskLocation.ASMAddress"],
			}, Daftar: map[string]string{"AnekaList": aneka}}
			if n := len(h.AmbilDaftar(aneka)); n > 0 { // 8.5.1.7.1 - AnekaList terakhir
				c.Daftar["CoverageList"] = JalurAnak(aneka, n, "CoverageList")
			}
			out = append(out, c)
		}
	}
	if IsAneka(h) { // 8.6
		out = append(out, calonAneka(h)...)
	}
	if IsMarineCargo(h) { // 8.7
		for _, p := range daftarPolis(h, "CargoList") {
			out = append(out, Calon{Kunci: p.kunci(), Medan: Baris{
				"ObjectID": p.b["TradingID"], "ObjectName": p.b["TradingNote"], "ObjectIDCard": p.b["GoodID"],
				"ObjectJob": p.b["GoodNote"], "ObjectWeight": p.b["PackingID"], "ObjectSurveyor": p.b["PackingNote"],
				"ObjectSurveyID": p.b["ConveyanceID"], "ObjectSurveyLocation": p.b["ConveyanceNote"],
			}, CoverageItem: p.anak("CoverageList")})
		}
	}
	return out
}

// subAneka - sub-cabang 8.6.1-8.6.24 (urutan XML). `aneka` = sub-cabang menyalin AnekaList occupation (8.6.1 / 8.6.3).
var subAneka = []struct {
	uji   func(*Halaman) bool
	aneka bool
	pohon bool // 8.6.2 GrowingTrees
}{
	{IsLiability, true, false}, {IsGrowingTrees, false, true}, {IsMBD, true, false}, {IsEar, false, false},
	{IsElectronicEquipment, false, false}, {IsMarineHull, false, false}, {IsAviationHull, false, false},
	{IsCustomBond, false, false}, {IsExclusion, false, false}, {IsMaintenance, false, false}, {IsCAR, false, false},
	{IsGlass, false, false}, {IsAllRisk, false, false}, {IsFidelity, false, false},
	{IsBillboardNeonSyariah, false, false}, {IsBurglary, false, false}, {IsCIT, false, false}, {IsCIS, false, false},
	{IsBoiler, false, false}, {IsHE, false, false}, {IsLandRig, false, false},
	{IsContractorsPlantMachinery, false, false}, {IsBonding, false, false},
	{func(h *Halaman) bool { return jenisSama(h, "Aneka") }, false, false},
}

// calonAneka - InsertObjects_dt 8.6: satu calon per lokasi; medan occupation/aneka TERAKHIR yang tinggal (setiap
// iterasi menimpa <LAST>); 8.6.25 OccupationList lokasi (lihat penyimpangan sadar di kepala berkas).
func calonAneka(h *Halaman) []Calon {
	var out []Calon
	for _, s := range subAneka {
		if !s.uji(h) {
			continue
		}
		for _, p := range daftarPolis(h, "LocationList") {
			lokasi := p.b["Property.RiskLocation.ASMAddress"]
			c := Calon{Kunci: p.kunci(), Medan: Baris{"ObjectID": p.b["Property.ObjectNo"]},
				Daftar: map[string]string{"OccupationList": p.anak("Property.RiskLocation.OccupationList")}}
			okup := p.anak("Property.RiskLocation.OccupationList")
			for oi, ob := range h.AmbilDaftar(okup) {
				anekaJ := JalurAnak(okup, oi+1, "AnekaList")
				if s.pohon { // 8.6.2.1.3.1
					for ai, ab := range h.AmbilDaftar(anekaJ) {
						c.Medan["ObjectName"] = ab["ObjectName"]
						c.Daftar["CoverageList"] = JalurAnak(anekaJ, ai+1, "CoverageList")
						c.Medan["ObjectVehicleChasis"] = ab["VehicleHE.ChassisNumber"]
						c.Medan["ObjectVehicleType"] = ab["VehicleHE.TypeName"]
					}
					c.Medan["ObjectLocation"] = lokasi // 8.6.2.1.3.2
					c.Medan["ObjectIDCard"] = ob["OccupationId"]
					continue
				}
				c.Medan["ObjectName"] = ob["OccupationName"]
				c.Medan["ObjectIDCard"] = ob["OccupationId"]
				c.Medan["ObjectLocation"] = lokasi
				if s.aneka {
					c.Daftar["AnekaList"] = anekaJ
				}
				if n := len(h.AmbilDaftar(anekaJ)); n > 0 {
					c.Daftar["CoverageList"] = JalurAnak(anekaJ, n, "CoverageList")
				}
			}
			out = append(out, c)
		}
	}
	return out
}

// medanTerpilih - InsertObjectItemList_DT 2.2.4.x: medan calon yang disalin ke `ClaimData.ObjectList` per lini, dan
// isi item pertama (`ObjectItemList(<APPEND>)`): IndexObject := .ObjectID atau CoverageList calon.
type medanLini struct {
	medan      []string
	daftar     []string
	itemIndeks bool // ObjectItemList(<APPEND>).IndexObject := .ObjectID
	itemCover  bool // ObjectItemList(<APPEND>).CoverageList := .CoverageList
	lokasiSurv bool // MarineCargo 2.2.4.2.1.1.5 ObjectLocation := .ObjectSurveyor
}

// liniTerpilih - urutan cabang InsertObjectItemList_DT 2.2.4.1-2.2.4.7 (Fire, MarineCargo, Aneka, MBU, PA, Golf, Travel).
var liniTerpilih = []struct {
	uji func(*Halaman) bool
	m   medanLini
}{
	{IsFire, medanLini{medan: []string{"ObjectName", "ObjectID", "ObjectStatus", "ObjectLocation", "ObjectSurveyLocation"},
		daftar: []string{"PropertyItemList", "OccupationList"}, itemIndeks: true}},
	{IsMarineCargo, medanLini{medan: []string{"ObjectName", "ObjectID", "ObjectIDCard", "ObjectJob", "ObjectSurveyID",
		"ObjectSurveyLocation"}, itemCover: true, lokasiSurv: true}},
	{IsAneka, medanLini{medan: []string{"ObjectName", "ObjectID", "ObjectLocation", "ObjectIDCard"},
		daftar: []string{"CoverageList", "AnekaList", "OccupationList"}, itemIndeks: true}},
	{IsMBU, medanLini{medan: []string{"Brand", "BrandName", "Model", "ModelName", "Type", "TypeName", "LicensePlate",
		"ChassisNumber", "EngineNumber"}, daftar: []string{"CoverageList"}, itemCover: true}},
	{IsPA, medanLini{medan: []string{"ObjectName", "ObjectID", "ObjectJob", "ObjectWeight", "ObjectHeight", "ASMClassID",
		"ObjectLeftHanded", "Gender", "ObjectDateOfBirth", "ObjectIDCard", "ObjectParticipantStatus"},
		daftar: []string{"CoverageList"}, itemIndeks: true}},
	{IsGolfInsurance, medanLini{medan: []string{"ObjectName", "ObjectID", "ObjectLocation", "ObjectIDCard"},
		daftar: []string{"CoverageList", "AnekaList"}, itemIndeks: true}},
	{IsTravel, medanLini{medan: []string{"ObjectName", "ObjectID", "ObjectParticipantStatus", "ObjectIDCard",
		"ObjectDateOfBirth"}, daftar: []string{"CoverageList", "OccupationList"}, itemCover: true}},
}

// MedanSalinanObjek - seluruh medan objek salinan polis (diisi ulang `LengkapiObjek`, tidak disimpan).
var MedanSalinanObjek = func() map[string]bool {
	out := map[string]bool{"ObjectSurveyor": true, "ObjectVehicleChasis": true, "ObjectVehicleType": true}
	for _, l := range liniTerpilih {
		for _, m := range l.m.medan {
			out[m] = true
		}
	}
	return out
}()

// DaftarSalinanObjek - anak daftar objek salinan polis.
var DaftarSalinanObjek = []string{"PropertyItemList", "OccupationList", "AnekaList", "CoverageList"}

// SusunCalon menulis grid pilih objek (`TempClaimData.ClaimData.ObjectList`) dari polis. Pilihan (`Selected`) yang dibawa
// layar dipertahankan menurut kunci polis.
func SusunCalon(h *Halaman) {
	lama := map[string]string{}
	for _, b := range h.AmbilDaftar(DaftarCalon) {
		lama[b[PropKunciPolis]] = b[PropDipilih]
	}
	var rows []Baris
	for _, c := range CalonObjek(h) {
		b := c.Medan.Salin()
		b[PropKunciPolis] = c.Kunci
		if v := lama[c.Kunci]; v != "" {
			b[PropDipilih] = v
		}
		rows = append(rows, b)
	}
	h.SetelDaftar(DaftarCalon, rows)
}

// calonMenurutKunci - calon ber-kunci polis `k`.
func calonMenurutKunci(cs []Calon, k string) (Calon, bool) {
	for _, c := range cs {
		if c.Kunci == k {
			return c, true
		}
	}
	return Calon{}, false
}

// SalinSubpohon menyalin daftar `dari` beserta seluruh daftar bersarang di bawah barisnya ke `ke` (Property-Set PageList
// := PageList).
func SalinSubpohon(h *Halaman, dari, ke string) {
	h.HapusAwalanDaftar(ke)
	h.SetelDaftar(ke, SalinDaftar(h.AmbilDaftar(dari)))
	awal := dari + "("
	for k, v := range h.Daftar {
		if strings.HasPrefix(k, awal) {
			h.Daftar[ke+k[len(dari):]] = SalinDaftar(v)
		}
	}
}

// HapusAwalanDaftar membuang daftar `jalur` beserta daftar bersarang di bawah barisnya.
func (h *Halaman) HapusAwalanDaftar(jalur string) {
	h.pastikan()
	delete(h.Daftar, jalur)
	for k := range h.Daftar {
		if strings.HasPrefix(k, jalur+"(") {
			delete(h.Daftar, k)
		}
	}
}

// LengkapiObjek menyalin ulang medan salinan polis setiap objek klaim dari calonnya (kunci polis). Objek yang calonnya
// tidak lagi ada di polis dibiarkan apa adanya (medan salinan kosong).
func LengkapiObjek(h *Halaman) {
	cs := CalonObjek(h)
	for o, b := range h.AmbilDaftar(DaftarObjek) {
		c, ada := calonMenurutKunci(cs, b[PropKunciPolis])
		if !ada {
			continue
		}
		m := medanUntuk(h)
		for _, k := range m.medan {
			b[k] = c.Medan[k]
		}
		if m.lokasiSurv {
			b["ObjectLocation"] = c.Medan["ObjectSurveyor"]
		}
		for _, k := range []string{"ObjectSurveyor", "ObjectVehicleChasis", "ObjectVehicleType"} {
			if v := c.Medan[k]; v != "" {
				b[k] = v
			}
		}
		for _, anak := range m.daftar {
			if src := c.Daftar[anak]; src != "" {
				SalinSubpohon(h, src, JalurAnak(DaftarObjek, o+1, anak))
			}
		}
		if m.itemCover && c.CoverageItem != "" && len(h.AmbilDaftar(DaftarItem(o+1))) > 0 {
			SalinSubpohon(h, c.CoverageItem, DaftarDiItem(o+1, 1, "CoverageList"))
		}
	}
}

// medanUntuk - medan salinan lini kasus ini (cabang PERTAMA yang lolos; InsertObjectItemList_DT cabang lain menulis medan
// yang sama untuk lini yang sama).
func medanUntuk(h *Halaman) medanLini {
	for _, l := range liniTerpilih {
		if l.uji(h) {
			return l.m
		}
	}
	return medanLini{}
}

// TerapkanObjekTerpilih = InsertObjectItemList_DT langkah 2 (Submit Input Register): ObjectList DIBUANG lalu dibangun
// ulang dari calon terpilih (`Param.CFS != 1` selalu benar - tidak ada pemanggil yang mengirim CFS, sehingga cabang
// selalu berjalan; Back hanya mungkin sebelum CFS, `InputEstimasiAdmin` Back VIS `.IsCFS != 1`). Kronologi "Finish
// Registration" bila ada objek terpilih (2.2.1). Mengembalikan true bila ObjectList dibangun.
func TerapkanObjekTerpilih(k *Konteks, h *Halaman) bool {
	h.HapusAwalanDaftar(DaftarObjek) // 2.1
	if h.Ambil(JalurUjiCalon) != "1" {
		return false
	}
	k.Kronologi(h, "Finish Registration") // 2.2.1-2.2.2
	m := medanUntuk(h)
	cs := CalonObjek(h)
	for _, b := range h.AmbilDaftar(DaftarCalon) {
		if b[PropDipilih] != "true" {
			continue
		}
		c, ada := calonMenurutKunci(cs, b[PropKunciPolis])
		if !ada {
			continue
		}
		baru := Baris{PropKunciPolis: c.Kunci}
		for _, f := range m.medan {
			baru[f] = c.Medan[f]
		}
		if m.lokasiSurv {
			baru["ObjectLocation"] = c.Medan["ObjectSurveyor"]
		}
		o := h.TambahBaris(DaftarObjek, baru)
		for _, anak := range m.daftar {
			if src := c.Daftar[anak]; src != "" {
				SalinSubpohon(h, src, JalurAnak(DaftarObjek, o, anak))
			}
		}
		item := Baris{}
		if m.itemIndeks {
			item["IndexObject"] = c.Medan["ObjectID"]
		}
		h.TambahBaris(DaftarItem(o), item)
		if m.itemCover && c.CoverageItem != "" {
			SalinSubpohon(h, c.CoverageItem, DaftarDiItem(o, 1, "CoverageList"))
		}
	}
	return true
}

// Anak daftar baris item objek dan adjustment.
const (
	AnakEstimasi     = "EstimationList"
	AnakSpreadPolis  = "SpreadingList"
	AnakSpreadKlaim  = "SpreadingClaim"
	AnakBreakQS      = "SpreadingAdjustment" // tingkat item (T_CLAIM_BREAK_QS)
	AnakAdj          = "Adjustment"
	AnakAdjSpread    = "SpreadingAdjustment" // tingkat adjustment (T_CLAIM_ADJ_SPREADING)
	AnakAdjQS        = "SpreadingQuotaShare"
	AnakFacRetro     = "FacRetroList"
	AnakCedingCedant = "CedingCedantList"
	AnakKomiteAdj    = "ComiteeClaim"
)

// DaftarAdj / JalurAdj - daftar adjustment item (o, i) dan baris ke-a.
func DaftarAdj(o, i int) string   { return DaftarDiItem(o, i, AnakAdj) }
func JalurAdj(o, i, a int) string { return JalurBaris(DaftarAdj(o, i), a) }
func DaftarDiAdj(o, i, a int, anak string) string {
	return JalurAnak(DaftarAdj(o, i), a, anak)
}
