package models

// Untuk apa berkas ini: PORT ACTIVITY ITEM OBJEK DAN ESTIMASI (FlowAction InputEstimasi; grid item
// `PropertyItemListGridEstimation*` / `ShowItemPA`, panel item `Estimasi` / `EstimasiPA` / `EstimasiMarine`). Kelas
// activity `ASM-FW-GCNMFW-Data-Object` (o = baris objek) / `Data-ObjectItem` (i = baris item) / `Data-Estimasi`
// (e = baris estimasi). Indeks berbasis satu, sama dengan Pega.
//
// ⚠️ Langkah Java "Hapus jika ada spreadingClaim yg sama" / "remove treaty yg sama" / "remove value yg sama" tidak
// diekspor. Isinya ditiru dari langkah agregasi SESUDAHNYA yang menjumlah per kunci (TreatyType; atau CurrencyID +
// TreatyType di GetCoverageAneka_Act 9): baris berkunci sama dibuang kecuali yang pertama (PARITAS `[inferensi]`).

import (
	"fmt"
	"strconv"

	"github.com/cockroachdb/apd/v3"
)

// Pesan validasi item / estimasi - VERBATIM.
const (
	PesanItemKembar          = "Can not add Object Item with the same name"             // ProtectionObjectItem_Act 1
	PesanEstimasiTanpaSpread = "Can not add estimation without spreading"               // ValidateInputEstimate_act 1
	PesanEstimasiKurangDOL   = "Estimation Date should not be less than Date of Loss"   // ValidateInputEstimate_act 1
	PesanEstimasiLebihHari   = "Estimation Date should not be more than todays date"    // ValidateInputEstimate_act 1
	PesanEstimasiLebihTSI    = "Total estimation value is more than TSI Nusare Limit"   // CheckEstimateValue 6
	PesanEstimasiNegatif     = "Estimation value cannot be filled with negative value"  // CheckEstimateValue 6
	PesanEstimasiNol         = "Estimation value cannot be filled with zero"            // CheckEstimateValue 6
	PesanEstimasiLebihLoL    = "Total estimation value is more than Limit of Liability" // CheckEstimateValue 6
)

// CheckTotalSpreadingPct_Act (pesan total persen spreading, prompt §5 butir 10) tidak dibangun: pemicunya `change` medan
// `.SharePercentage` grid Spreading Policy / Claim yang RO=ALWAYS di section Estimasi / EstimasiPA / EstimasiMarine -
// tidak pernah terpicu (PARITAS).

func barisDi(h *Halaman, daftar string, n int) (Baris, error) {
	d := h.AmbilDaftar(daftar)
	if n < 1 || n > len(d) {
		return nil, fmt.Errorf("%w: %s(%d)", ErrBarisTidakAda, daftar, n)
	}
	return d[n-1], nil
}

// Objek / Item - baris objek ke-o dan item ke-i.
func Objek(h *Halaman, o int) (Baris, error)   { return barisDi(h, DaftarObjek, o) }
func Item(h *Halaman, o, i int) (Baris, error) { return barisDi(h, DaftarItem(o), i) }

// ---------------------------------------------------------------- indeks turunan

// TurunkanIndeksItem menulis medan item yang XML turunkan dari objek induknya setiap panel dibuka / item ditambah
// (`SetIndex_Act` pra-proses masterDetail, `SetIndexObject_Act` tombol Add, `SetEstimation_DT` 6-9 pra-proses
// InputEstimasi, `CheckEstimateValue` 5). Medan ini tidak disimpan (fungsi murni objek induk):
//
//	semua lini     .ObjectIndex := o
//	selain MBU     .IndexObject := objek.ObjectID   (SetIndexObject_Act 7-8.1)
//	MBU            .IndexObject = .ObjectIndex = .IndexPropertyItem := o; .ObjectID := .Brand; LicensePlate /
//	               ChassisNumber / EngineNumber objek (SetIndexObject_Act 5.1, SetEstimation_DT 6)
//	MarineCargo    .IndexObject := o; .ObjectID := o (SetEstimation_DT 7, CheckEstimateValue 5.2.1)
//	PA / Travel    .ObjectID / .ObjectName objek (SetIndex_Act 8-9)
//	Aneka          .ObjectItemID := objek.ObjectIDCard (SetIndex_Act 6.2)
func TurunkanIndeksItem(h *Halaman) {
	mbu, marine, pa, travel, aneka := IsMBU(h), IsMarineCargo(h), IsPA(h), IsTravel(h), IsAneka(h)
	for o, ob := range h.AmbilDaftar(DaftarObjek) {
		on := strconv.Itoa(o + 1)
		for _, it := range h.AmbilDaftar(DaftarItem(o + 1)) {
			it["ObjectIndex"] = on
			switch {
			case mbu:
				it["IndexObject"], it["IndexPropertyItem"] = on, on
				it["ObjectID"] = ob["Brand"]
				it["LicensePlate"], it["ChassisNumber"], it["EngineNumber"] = ob["LicensePlate"], ob["ChassisNumber"],
					ob["EngineNumber"]
			case marine:
				it["IndexObject"], it["ObjectID"] = on, on
			default:
				it["IndexObject"] = ob["ObjectID"]
			}
			if pa || travel {
				it["ObjectID"], it["ObjectName"] = ob["ObjectID"], ob["ObjectName"]
			}
			if aneka {
				it["ObjectItemID"] = ob["ObjectIDCard"]
			}
		}
	}
}

// TambahItem = tombol Add grid item (`addRow` + `SetIndexObject_Act`). Mengembalikan indeks item baru.
func TambahItem(h *Halaman, o int) (int, error) {
	if _, err := Objek(h, o); err != nil {
		return 0, err
	}
	n := h.TambahBaris(DaftarItem(o), Baris{})
	TurunkanIndeksItem(h)
	return n, nil
}

// HapusItem = tombol Delete grid item (`deleteRow`, lalu `DeleteTest` / `DeleteObjectItemMBU`: IsError dikosongkan bila
// tidak ada item bernama sama; tanggal estimasi item lain diperiksa ulang `ProtectionDate_Act`).
func HapusItem(k *Konteks, h *Halaman, o, i int) error {
	if _, err := Item(h, o, i); err != nil {
		return err
	}
	h.HapusBaris(DaftarItem(o), i)
	kembar := false
	nama := map[string]bool{}
	for _, it := range h.AmbilDaftar(DaftarItem(o)) { // DeleteTest 2.2
		n := it["ObjectItemName"]
		if IsMBU(h) || IsPA(h) || IsTravel(h) {
			n = it["CoverageOLDID"]
		}
		if n != "" && nama[n] {
			kembar = true
		}
		nama[n] = true
	}
	if !kembar {
		h.Setel(JalurIsError, "")
	}
	for i2, it := range h.AmbilDaftar(DaftarItem(o)) { // DeleteTest 2.4 / DeleteObjectItemMBU 3.3
		if it["ObjectItemName"] == "" && it["CoverageOLDID"] == "" {
			continue
		}
		for e := range h.AmbilDaftar(DaftarDiItem(o, i2+1, AnakEstimasi)) {
			ProtectionDate(k, h, o, i2+1, e+1)
		}
	}
	return nil
}

// ---------------------------------------------------------------- pilih item / coverage (Fire)

// OpsiPropertyItem = D_FilteredPropertyItemList (`FilterPropertyItem_Act`, Param.ObjectID = .IndexObject): PropertyItemList
// objek ber-ObjectID sama; IndexPropertyItem kosong = urutan baris.
func OpsiPropertyItem(h *Halaman, indexObject string) []Baris {
	for o, ob := range h.AmbilDaftar(DaftarObjek) {
		if ob["ObjectID"] != indexObject {
			continue
		}
		var out []Baris
		for n, p := range h.AmbilDaftar(JalurAnak(DaftarObjek, o+1, "PropertyItemList")) {
			b := p.Salin()
			if b["IndexPropertyItem"] == "" { // 1.1.1
				b["IndexPropertyItem"] = strconv.Itoa(n + 1)
			}
			out = append(out, b)
		}
		return out
	}
	return nil
}

// OpsiCoverageFire = D_FilteredCoverageList (`FilterCoverage_Act`, ObjectID = .IndexObject, IndexPropertyItem =
// .ObjectItemID): CoverageList PropertyItem ber-PropertyItemNo = parameter.
func OpsiCoverageFire(h *Halaman, indexObject, noItem string) []Baris {
	for o, ob := range h.AmbilDaftar(DaftarObjek) {
		if ob["ObjectID"] != indexObject {
			continue
		}
		pl := JalurAnak(DaftarObjek, o+1, "PropertyItemList")
		for n, p := range h.AmbilDaftar(pl) {
			if p["PropertyItemNo"] == noItem {
				return SalinDaftar(h.AmbilDaftar(JalurAnak(pl, n+1, "CoverageList")))
			}
		}
	}
	return nil
}

// PilihItemProperti = autocomplete Object Name grid item Fire (isi `.IndexPropertyItem`), lalu `SetObjectItem_Act`
// (ObjectItemID, Currency, TSIPerObject, CurrencyID, TotalGrossPremi, TotalPremiumNusantaraRe dari PropertyItem),
// setValue CoverageNote / IndexCoverage = "", `SetIndex_Act`, `SetConvertValueKurs_ObjectItem` (kurs + proteksi nama).
func PilihItemProperti(k *Konteks, h *Halaman, o, i int, indexPropItem string) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	TurunkanIndeksItem(h)
	var pilih Baris
	for _, p := range OpsiPropertyItem(h, it["IndexObject"]) {
		if p["IndexPropertyItem"] == indexPropItem {
			pilih = p
		}
	}
	if pilih == nil {
		return fmt.Errorf("%w: property item %q", ErrBarisTidakAda, indexPropItem)
	}
	it["ObjectItemName"] = pilih["ItemType"]
	it["IndexPropertyItem"] = indexPropItem
	no := pilih["PropertyItemNo"] // SetObjectItem_Act 1.1.1.1.1
	if no == "" {
		no = indexPropItem
	}
	it["ObjectItemID"] = no // 1.1.1.1.2
	it["Currency"], it["TSIPerObject"], it["CurrencyID"] = pilih["Currency"], pilih["TSIObjectItem"], pilih["CurrencyID"]
	it["TotalGrossPremi"], it["TotalPremiumNusantaraRe"] = pilih["TotalGrossPremi"], pilih["TotalPremiumNusantaraRe"]
	it["CoverageNote"], it["IndexCoverage"] = "", "" // setValue
	if err := SetKursItem(k, h, o, i); err != nil {
		return err
	}
	return ProteksiItemKembar(h, o, i, it["ObjectItemName"])
}

// SetKursItem = `SetConvertValueKurs_ObjectItem` / `GetKursObjectItem_Act` 1-3: KursObjectItem = kurs standar CurrencyID.
func SetKursItem(k *Konteks, h *Halaman, o, i int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	v, err := k.Acuan.KursStandar(k.Ctxt(), it["CurrencyID"])
	if err != nil {
		return err
	}
	it["KursObjectItem"] = v
	return nil
}

// ProteksiItemKembar = `ProtectionObjectItem_Act` (Param.ObjectName): item lain pada objek yang sama bernama sama ->
// IsError = jumlahnya; >1 -> pesan + item terakhir yang sama dibuang (4.1); tidak ada -> IsError 0, pesan dibersihkan.
//
// ⚠️ Langkah 2.1 melewati baris yang SAMA dengan item yang sedang diubah, maka hitungan = item LAIN bernama sama; langkah
// 3.1 (`pre=false`) memasang pesan bila hitungan > 0.
func ProteksiItemKembar(h *Halaman, o, i int, nama string) error {
	n, akhir := 0, 0
	for j, it := range h.AmbilDaftar(DaftarItem(o)) {
		if j+1 == i { // 2.1
			continue
		}
		if nama != "" && it["ObjectItemName"] == nama { // 2.1.1
			n++
			akhir = j + 1
		}
	}
	if n == 0 { // 2.1.3, 5
		h.Setel(JalurIsError, "0")
		return nil
	}
	h.Setel(JalurIsError, strconv.Itoa(n))                              // 2.1.2
	h.TambahPesan(JalurAnak(DaftarItem(o), i, "Test"), PesanItemKembar) // 3.1
	if n > 1 {                                                          // 4.1
		h.HapusBaris(DaftarItem(o), akhir)
	}
	return nil
}

// PilihCoverageFire = autocomplete Coverage Name grid item Fire (D_FilteredCoverageList nilai .OLDID; isi CoverageNote,
// IndexCoverage, CoverageID (.Coverage), TSINusare (.TSINusantaraRe), PremiNusare, LimitofLiability), lalu
// `CopySpreading_Act` dan `ConvertTSINusare_Act`.
func PilihCoverageFire(h *Halaman, o, i int, oldID string) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	TurunkanIndeksItem(h)
	var c Baris
	for _, b := range OpsiCoverageFire(h, it["IndexObject"], it["ObjectItemID"]) {
		if b["OLDID"] == oldID {
			c = b
		}
	}
	if c == nil {
		return fmt.Errorf("%w: coverage %q", ErrBarisTidakAda, oldID)
	}
	it["CoverageOLDID"] = oldID
	it["CoverageNote"], it["IndexCoverage"], it["CoverageID"] = c["CoverageNote"], c["IndexCoverage"], c["Coverage"]
	it["TSINusare"], it["PremiNusare"], it["LimitofLiability"] = c["TSINusantaraRe"], c["PremiNusantaraRe"],
		c["LimitofLiability"]
	CopySpreading(h, o, i)
	return ConvertTSINusare(h, o, i)
}

// ConvertTSINusare = `ConvertTSINusare_Act`: ValueTSINusareIDR = KursObjectItem * TSINusare.
func ConvertTSINusare(h *Halaman, o, i int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	var k Kalkulator
	v := k.Kali(k.B(it, "KursObjectItem"), k.B(it, "TSINusare"))
	if err := k.Galat(); err != nil {
		return err
	}
	it["ValueTSINusareIDR"] = Teks(v)
	return nil
}

// spreadPolis - SpreadingList coverage polis Fire (CopySpreading_Act 5-6): LocationList(.IndexObject).Property.
// PropertyItemList(.IndexPropertyItem).CoverageList(.IndexCoverage).SpreadingList, atau LayerList(1).SpreadingList bila
// CoverageBasis 5.
func spreadPolis(h *Halaman, it Baris) []Baris {
	lo, err1 := strconv.Atoi(it["IndexObject"])
	pi, err2 := strconv.Atoi(it["IndexPropertyItem"])
	ci, err3 := strconv.Atoi(it["IndexCoverage"])
	if err1 != nil || err2 != nil || err3 != nil {
		return nil
	}
	cov := JalurBaris(JalurAnak(JalurAnak(DaftarLokasiPolis, lo, "Property.PropertyItemList"), pi, "CoverageList"), ci)
	b, ok := barisJalur(h, cov)
	if !ok {
		return nil
	}
	if b["CoverageBasis"] == "5" { // 6
		return SalinDaftar(h.AmbilDaftar(JalurAnak(cov+".LayerList", 1, "SpreadingList")))
	}
	return SalinDaftar(h.AmbilDaftar(cov + ".SpreadingList"))
}

// barisJalur - baris di jalur "daftar(n)".
func barisJalur(h *Halaman, j string) (Baris, bool) {
	d, n, _, ok := pecahJalurBaris(j + ".x")
	if !ok {
		return nil, false
	}
	rows := h.AmbilDaftar(d)
	if n < 1 || n > len(rows) {
		return nil, false
	}
	return rows[n-1], true
}

// agregasiSpread - langkah Java "hapus yang sama" + jumlah SharePercentage per kunci + buang yang 0 (CopySpreading_Act 8-9,
// GetCurencyCoverage_Act 6-7, GetSpreadingMarine_Act 2.2-2.3).
func agregasiSpread(rows []Baris, kunci func(Baris) string) ([]Baris, error) {
	var k Kalkulator
	var out []Baris
	lihat := map[string]bool{}
	for _, b := range rows {
		key := kunci(b)
		if lihat[key] {
			continue
		}
		lihat[key] = true
		total := apd.New(0, 0)
		for _, b1 := range rows {
			if kunci(b1) == key {
				total = k.Tambah(total, k.B(b1, "SharePercentage"))
			}
		}
		nb := b.Salin()
		nb["SharePercentage"] = Teks(total)
		if Nol(total) {
			continue
		}
		out = append(out, nb)
	}
	return out, k.Galat()
}

func kunciTreaty(b Baris) string { return b["TreatyType"] }

// totalSpreadList = langkah "get total spreading list" (CopySpreading_Act 13-14, GetCurencyCoverage_Act 9-10):
// ClaimSpreaded = `Local.Total` (ValueTSINusareIDR item, CopySpreading_Act 10; activity lain tidak mengisi Local.Total
// -> kosong), total Share / TSISpreaded / PremiumSpreaded ke item.
func totalSpreadList(it Baris, rows []Baris, total string) error {
	var k Kalkulator
	share, tsi, premi := apd.New(0, 0), apd.New(0, 0), apd.New(0, 0)
	for _, b := range rows {
		b["ClaimSpreaded"] = total
		share = k.Tambah(share, k.B(b, "SharePercentage"))
		tsi = k.Tambah(tsi, k.B(b, "TSISpreaded"))
		premi = k.Tambah(premi, k.B(b, "PremiumSpreaded"))
	}
	if err := k.Galat(); err != nil {
		return err
	}
	it["TotalSharePercentage"], it["TotalTSISpreaded"], it["TotalPremiumSpreaded"] = Teks(share), Teks(tsi), Teks(premi)
	return nil
}

// CopySpreading = `CopySpreading_Act` (Fire): spreading polis coverage -> SpreadingList dan SpreadingClaim item.
// ExGratia klaim selalu 0 di modul ini (InsertObjects_dt 11), cabang ORS langkah 7 tidak terjangkau; cabang TreatyIn
// langkah 12 tidak terjangkau.
func CopySpreading(h *Halaman, o, i int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	if it["IndexPropertyItem"] == "" { // 1
		it["IndexPropertyItem"] = it["ObjectItemID"]
	}
	rows := spreadPolis(h, it) // 5-6
	agg, err := agregasiSpread(rows, kunciTreaty)
	if err != nil {
		return err
	}
	if err := totalSpreadList(it, agg, it["ValueTSINusareIDR"]); err != nil { // 10, 13-14
		return err
	}
	h.SetelDaftar(DaftarDiItem(o, i, AnakSpreadPolis), agg)              // 11
	h.SetelDaftar(DaftarDiItem(o, i, AnakSpreadKlaim), SalinDaftar(agg)) // 11
	return ConvertTSINusare(h, o, i)                                     // 15
}

// ---------------------------------------------------------------- Aneka / Golf

// OpsiAneka = D_AnekaList (`GetAnekaList`, IdxLoc = .IndexObject): AnekaList lokasi ber-ObjectNo sama (Golf:
// Property.RiskLocation.AnekaList; Aneka: AnekaList setiap occupation), ber-IdxAneka = urutan baris.
func OpsiAneka(h *Halaman, idxLoc string) []Baris {
	var out []Baris
	for _, p := range daftarPolis(h, "LocationList") {
		if p.b["Property.ObjectNo"] != idxLoc {
			continue
		}
		if IsGolfInsurance(h) { // 2
			for n, a := range h.AmbilDaftar(p.anak("Property.RiskLocation.AnekaList")) {
				b := a.Salin()
				b["IdxAneka"] = strconv.Itoa(n + 1)
				out = append(out, b)
			}
			continue
		}
		okup := p.anak("Property.RiskLocation.OccupationList") // 1
		for oi := range h.AmbilDaftar(okup) {
			for n, a := range h.AmbilDaftar(JalurAnak(okup, oi+1, "AnekaList")) {
				b := a.Salin() // SetIndexObject_DT 1.1.2.2.2-4
				b["IdxLocation"], b["IdxOccupation"], b["IdxAneka"] = strconv.Itoa(p.n), strconv.Itoa(oi+1),
					strconv.Itoa(n+1)
				out = append(out, b)
			}
		}
	}
	return out
}

// PilihAneka = autocomplete Object Id grid item Aneka / Golf (D_AnekaList nilai .IdxAneka), lalu setValue CoverageNote /
// IndexCoverage = "", `SetIndex_Act`, `GettsiAneka_Act` (TSIPerObject, Currency, IndexOccupation, IndexAneka,
// ObjectItemName dari aneka terpilih), `GetKursObjectItem_Act` (kurs + proteksi nama).
func PilihAneka(k *Konteks, h *Halaman, o, i int, idxAneka string) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	it["IndexAneka"] = idxAneka
	it["CoverageNote"], it["IndexCoverage"] = "", ""
	TurunkanIndeksItem(h)
	gettsiAneka(h, it)
	if err := SetKursItem(k, h, o, i); err != nil {
		return err
	}
	return ProteksiItemKembar(h, o, i, it["ObjectItemName"])
}

// gettsiAneka = `GettsiAneka_Act` (Param ObjectNo = .IndexObject, OcupationID = .ObjectItemID, IndexAneka).
//
// Selain Golf (1): lokasi ObjectNo, occupation OccupationId, setiap aneka: IdxOccupation terisi DAN IndexAneka ==
// IdxAneka (1.1.1.1.2.1; bukan MBD) / IdxOccupation kosong (1.1.1.1.2.2) / MBD dengan IndexAneka == IdxAneka (2.3).
// Golf (2): aneka lokasi ber-urutan = IndexAneka. Nilai TERAKHIR yang cocok menang (loop menimpa Local).
func gettsiAneka(h *Halaman, it Baris) {
	var hasil Baris
	idxOk, idxAneka := "", ""
	for _, p := range daftarPolis(h, "LocationList") {
		if p.b["Property.ObjectNo"] != it["IndexObject"] {
			continue
		}
		if IsGolfInsurance(h) {
			for n, a := range h.AmbilDaftar(p.anak("Property.RiskLocation.AnekaList")) {
				if strconv.Itoa(n+1) == it["IndexAneka"] { // 2.1.2.1
					hasil, idxOk, idxAneka = a, strconv.Itoa(n+1), strconv.Itoa(n+1)
				}
			}
			continue
		}
		okup := p.anak("Property.RiskLocation.OccupationList")
		for oi, ob := range h.AmbilDaftar(okup) {
			if ob["OccupationId"] != it["ObjectItemID"] {
				continue
			}
			for n, a := range h.AmbilDaftar(JalurAnak(okup, oi+1, "AnekaList")) {
				idx := strconv.Itoa(n + 1) // IdxAneka (SetIndexObject_DT)
				switch {
				case IsMBD(h) && it["IndexAneka"] == idx: // 1.1.1.1.2.3
					hasil, idxOk, idxAneka = a, idx, idx
				case !IsMBD(h) && a["IdxOccupation"] != "" && it["IndexAneka"] == idx: // 1.1.1.1.2.1
					hasil, idxOk, idxAneka = a, a["IdxOccupation"], idx
				case !IsMBD(h) && a["IdxOccupation"] == "": // 1.1.1.1.2.2
					hasil, idxOk, idxAneka = a, idx, idx
				}
			}
		}
	}
	if hasil == nil {
		return
	}
	it["TSIPerObject"], it["Currency"], it["CurrencyID"] = hasil["TSI"], hasil["Currency.Name"], hasil["Currency.ID"]
	it["IndexOccupation"], it["IndexAneka"], it["ObjectItemName"] = idxOk, idxAneka, hasil["ObjectName"]
}

// OpsiCoverageAneka = D_FilteredCoverageAnekaList (`GetCoverageAnekaList`): selain Golf - CoverageList aneka ber-
// ObjectName = .ObjectItemName pada occupation .ObjectItemID lokasi .IndexObject (1.1.1.1.1.3); Golf - CoverageList aneka
// ke-IndexAneka (2.1.1.3).
func OpsiCoverageAneka(h *Halaman, it Baris) []Baris {
	var out []Baris
	for _, p := range daftarPolis(h, "LocationList") {
		if p.b["Property.ObjectNo"] != it["IndexObject"] {
			continue
		}
		if IsGolfInsurance(h) {
			aj := p.anak("Property.RiskLocation.AnekaList")
			for n := range h.AmbilDaftar(aj) {
				if strconv.Itoa(n+1) == it["IndexAneka"] {
					out = append(out, SalinDaftar(h.AmbilDaftar(JalurAnak(aj, n+1, "CoverageList")))...)
				}
			}
			continue
		}
		okup := p.anak("Property.RiskLocation.OccupationList")
		for oi, ob := range h.AmbilDaftar(okup) {
			if ob["OccupationId"] != it["ObjectItemID"] {
				continue
			}
			aj := JalurAnak(okup, oi+1, "AnekaList")
			for n, a := range h.AmbilDaftar(aj) {
				if a["ObjectName"] == it["ObjectItemName"] {
					out = append(out, SalinDaftar(h.AmbilDaftar(JalurAnak(aj, n+1, "CoverageList")))...)
				}
			}
		}
	}
	return out
}

// PilihCoverageAneka = autocomplete Coverage Name grid item Aneka / Golf (nilai .CoverageNote; isi .Coverage ->
// CoverageID), lalu `GetCoverageAneka_Act`.
func PilihCoverageAneka(k *Konteks, h *Halaman, o, i int, note string) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	TurunkanIndeksItem(h)
	for _, c := range OpsiCoverageAneka(h, it) {
		if c["CoverageNote"] == note {
			it["CoverageID"] = c["Coverage"]
		}
	}
	it["CoverageNote"] = note
	return getCoverageAneka(h, o, i)
}

// getCoverageAneka = `GetCoverageAneka_Act`: premi / TSI RNM dan spreading coverage aneka (ObjectName + TSI sama, nama
// coverage sama), agregasi per CurrencyID + TreatyType (9), total (11-12), ValueTSINusareIDR (13), proteksi nama (15).
func getCoverageAneka(h *Halaman, o, i int) error {
	it, _ := Item(h, o, i)
	var cov Baris
	var spread []Baris
	for _, p := range daftarPolis(h, "LocationList") {
		if p.b["Property.ObjectNo"] != it["IndexObject"] {
			continue
		}
		var aneka []string
		if IsGolfInsurance(h) { // 4
			aj := p.anak("Property.RiskLocation.AnekaList")
			for n := range h.AmbilDaftar(aj) {
				if strconv.Itoa(n+1) == it["IndexOccupation"] { // Param.AnekaID = .IndexOccupation
					aneka = append(aneka, JalurAnak(aj, n+1, "CoverageList"))
				}
			}
		} else { // 2 / 3 (MBD)
			okup := p.anak("Property.RiskLocation.OccupationList")
			for oi, ob := range h.AmbilDaftar(okup) {
				if ob["OccupationId"] != it["ObjectItemID"] {
					continue
				}
				aj := JalurAnak(okup, oi+1, "AnekaList")
				for n, a := range h.AmbilDaftar(aj) {
					if a["ObjectName"] == it["ObjectItemName"] && samaAngka(a["TSI"], it["TSIPerObject"]) {
						aneka = append(aneka, JalurAnak(aj, n+1, "CoverageList"))
					}
				}
			}
		}
		for _, cj := range aneka {
			for n, c := range h.AmbilDaftar(cj) {
				if IsGolfInsurance(h) || c["CoverageNote"] == it["CoverageNote"] {
					cov = c
					spread = SalinDaftar(h.AmbilDaftar(JalurAnak(cj, n+1, "SpreadingList")))
				}
			}
		}
	}
	if cov != nil { // 5
		it["PremiNusantaraRe"], it["PremiNusare"], it["TSINusare"] = cov["PremiNusantaraRe"], cov["Premium"],
			cov["TSINusantaraRe"]
	}
	agg, err := agregasiSpread(spread, func(b Baris) string { return b["CurrencyID"] + "\x00" + b["TreatyType"] })
	if err != nil {
		return err
	}
	if err := totalSpreadList(it, agg, ""); err != nil { // 11-12 (Local.Total tidak diisi activity ini)
		return err
	}
	h.SetelDaftar(DaftarDiItem(o, i, AnakSpreadPolis), agg)              // 10
	h.SetelDaftar(DaftarDiItem(o, i, AnakSpreadKlaim), SalinDaftar(agg)) // 10
	var k Kalkulator
	v := k.Kali(k.B(it, "TSINusare"), k.B(it, "KursObjectItem")) // 13
	if err := k.Galat(); err != nil {
		return err
	}
	it["ValueTSINusareIDR"] = Teks(v)
	return ProteksiItemKembar(h, o, i, it["ObjectItemName"]) // 15
}

// ---------------------------------------------------------------- MBU / PA / Travel

// OpsiCoverageObjek = D_CoverageMBUClaimList (`GetCoverageListMBUClaim`: VehicleList Brand + LicensePlate + Chassis +
// Engine sama) / D_CoveragePAClaimList / D_CoverageTravelClaimList (PersonList urutan = ID dan pyFullName = Name):
// CoverageList / ASMCoverage polis objek item.
func OpsiCoverageObjek(h *Halaman, it Baris) []Baris {
	switch {
	case IsMBU(h):
		for _, p := range daftarPolis(h, "VehicleList") {
			if p.b["Brand"] == it["ObjectID"] && p.b["LicensePlate"] == it["LicensePlate"] &&
				p.b["ChassisNumber"] == it["ChassisNumber"] && p.b["EngineNumber"] == it["EngineNumber"] {
				return SalinDaftar(h.AmbilDaftar(p.anak("CoverageList")))
			}
		}
	case IsPA(h), IsTravel(h):
		id := it["ObjectID"] // Travel: ID = .ObjectID; PA: ID = .IndexObject (= objek.ObjectID)
		if IsPA(h) {
			id = it["IndexObject"]
		}
		for _, p := range daftarPolis(h, "PersonList") {
			if strconv.Itoa(p.n) == id && p.b["pyFullName"] == it["ObjectName"] {
				return SalinDaftar(h.AmbilDaftar(p.anak("ASMCoverage")))
			}
		}
	}
	return nil
}

// PilihCoverageObjek = autocomplete Coverage ID grid item MBU / PA / Travel (nilai .Coverage), lalu `SetCoverageID_DT`
// (CoverageID, CoverageNote), `SetSpreading_Act` -> `ProtectCoverage_Act` (coverage kembar), `GetCurencyCoverage_Act`
// (spreading per TreatyType, TSI, mata uang; MBU: `GetKursObjectItemMBU_Act`).
func PilihCoverageObjek(k *Konteks, h *Halaman, o, i int, coverage string) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	TurunkanIndeksItem(h)
	var c Baris
	for _, b := range OpsiCoverageObjek(h, it) {
		if b["Coverage"] == coverage {
			c = b
		}
	}
	if c == nil {
		return fmt.Errorf("%w: coverage %q", ErrBarisTidakAda, coverage)
	}
	it["CoverageOLDID"], it["CoverageID"], it["CoverageNote"] = coverage, coverage, c["CoverageNote"] // SetCoverageID_DT
	if IsMBU(h) {
		it["IndexCoverage"] = c["IndexCoverage"]
	}
	if proteksiCoverageKembar(h, o, i, coverage) { // SetSpreading_Act -> ProtectCoverage_Act
		return nil
	}
	// GetCurencyCoverage_Act 1-3: coverage terpilih
	it["TSIPerObject"], it["TSINusare"] = c["TSILiability"], c["TSINusantaraRe"]
	it["Currency"], it["CurrencyID"] = c["Currency.Name"], c["Currency.ID"]
	agg, err := agregasiSpread(SalinDaftar(h.AmbilDaftar(coverageSpread(h, it, c))), kunciTreaty) // 5-7
	if err != nil {
		return err
	}
	if err := totalSpreadList(it, agg, ""); err != nil { // 9-10
		return err
	}
	h.SetelDaftar(DaftarDiItem(o, i, AnakSpreadPolis), agg)              // 8
	h.SetelDaftar(DaftarDiItem(o, i, AnakSpreadKlaim), SalinDaftar(agg)) // 8
	if err := SetKursItem(k, h, o, i); err != nil {                      // 13 GetKursObjectItemMBU_Act 1-3
		return err
	}
	var kk Kalkulator
	v := kk.Kali(kk.B(it, "TSINusare"), kk.B(it, "KursObjectItem")) // GetKursObjectItemMBU_Act 4-5
	if err := kk.Galat(); err != nil {
		return err
	}
	it["ValueTSINusareIDR"] = Teks(v)
	return nil
}

// coverageSpread - jalur SpreadingList coverage terpilih di halaman polis.
func coverageSpread(h *Halaman, it, c Baris) string {
	switch {
	case IsMBU(h):
		for _, p := range daftarPolis(h, "VehicleList") {
			if p.b["Brand"] != it["ObjectID"] || p.b["LicensePlate"] != it["LicensePlate"] {
				continue
			}
			for n, b := range h.AmbilDaftar(p.anak("CoverageList")) {
				if b["Coverage"] == c["Coverage"] && b["IndexCoverage"] == it["IndexCoverage"] {
					return JalurAnak(p.anak("CoverageList"), n+1, "SpreadingList")
				}
			}
		}
	default:
		for _, p := range daftarPolis(h, "PersonList") {
			for n, b := range h.AmbilDaftar(p.anak("ASMCoverage")) {
				if b["Coverage"] == c["Coverage"] && p.b["pyFullName"] == it["ObjectName"] {
					return JalurAnak(p.anak("ASMCoverage"), n+1, "SpreadingList")
				}
			}
		}
	}
	return ""
}

// proteksiCoverageKembar = `ProtectCoverage_Act`: lebih dari satu item objek ber-CoverageOLDID sama -> IsError, item
// (terakhir yang sama) dibuang, pesan. Mengembalikan true bila item dibuang.
func proteksiCoverageKembar(h *Halaman, o, i int, coverage string) bool {
	n, akhir := 0, 0
	for j, it := range h.AmbilDaftar(DaftarItem(o)) {
		if it["CoverageOLDID"] == coverage {
			n++
			akhir = j + 1
		}
	}
	if n <= 1 {
		return false
	}
	h.Setel(JalurIsError, strconv.Itoa(n))
	h.HapusBaris(DaftarItem(o), akhir)
	h.TambahPesan(JalurAnak(DaftarObjek, o, "Brand"), PesanItemKembar)
	return akhir == i
}

// ---------------------------------------------------------------- Marine (EstimasiMarine)

// SpreadingMarine = `GetSpreadingMarine_Act` (pra-proses EstimasiMarine_FA): spreading CoverageList item (coverage
// TERAKHIR yang menimpa), agregasi per TreatyType, total, lalu `GetKursObjectItemMarine_Act` atas coverage itu (Currency,
// kurs, TSINusare, ValueTSINusareIDR). Cabang ExGratia 1 tidak terjangkau.
func SpreadingMarine(k *Konteks, h *Halaman, o, i int) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	cj := DaftarDiItem(o, i, "CoverageList")
	covs := h.AmbilDaftar(cj)
	if len(covs) == 0 {
		return nil
	}
	ci := len(covs) // 2.1.1 Local.indexCoverage = baris terakhir
	agg, err := agregasiSpread(SalinDaftar(h.AmbilDaftar(JalurAnak(cj, ci, "SpreadingList"))), kunciTreaty)
	if err != nil {
		return err
	}
	if err := totalSpreadList(it, agg, ""); err != nil {
		return err
	}
	h.SetelDaftar(DaftarDiItem(o, i, AnakSpreadPolis), agg)
	h.SetelDaftar(DaftarDiItem(o, i, AnakSpreadKlaim), SalinDaftar(agg))
	c := covs[ci-1] // GetKursObjectItemMarine_Act
	it["Currency"], it["CurrencyID"] = c["Currency.Name"], c["Currency.ID"]
	v, err := k.Acuan.KursStandar(k.Ctxt(), c["Currency.ID"])
	if err != nil {
		return err
	}
	it["KursObjectItem"] = v
	var kk Kalkulator
	idr := kk.Kali(kk.B(c, "TSINusantaraRe"), kk.Teks("KursObjectItem", v))
	if err := kk.Galat(); err != nil {
		return err
	}
	it["TSINusare"], it["ValueTSINusareIDR"] = c["TSINusantaraRe"], Teks(idr)
	return nil
}

// LengkapiMarine = pra-proses `EstimasiMarine_FA` (`GetSpreadingMarine_Act`) item Marine Cargo. Pega menjalankannya
// setiap pane item dibuka; pra-proses di sini berjalan setiap aksi, sehingga dijalankan hanya bila Spreading Policy item
// masih kosong - menjalankannya ulang menimpa Spreading Claim hasil estimasi (`[penyimpangan sadar]`, PARITAS).
func LengkapiMarine(k *Konteks, h *Halaman) error {
	if !IsMarineCargo(h) {
		return nil
	}
	for o := range h.AmbilDaftar(DaftarObjek) {
		for i := range h.AmbilDaftar(DaftarItem(o + 1)) {
			if len(h.AmbilDaftar(DaftarDiItem(o+1, i+1, AnakSpreadPolis))) > 0 {
				continue
			}
			if err := SpreadingMarine(k, h, o+1, i+1); err != nil {
				return err
			}
		}
	}
	return nil
}

// ---------------------------------------------------------------- okupasi

// OpsiOkupasi = D_OccupationList (`Occupation_Act`, Param.ObjectID = .IndexObject item): OccupationList objek ber-ObjectID
// sama (salinan polis, `LengkapiObjek`).
func OpsiOkupasi(h *Halaman, indexObject string) []Baris {
	var out []Baris
	for o, ob := range h.AmbilDaftar(DaftarObjek) {
		if ob["ObjectID"] == indexObject {
			out = append(out, h.AmbilDaftar(JalurAnak(DaftarObjek, o+1, "OccupationList"))...)
		}
	}
	return out
}

// PilihOkupasi = autocomplete Occupation grid item (nilai .OccupationName, isi .OccupationId -> .OccupationId): baris
// D_OccupationList ber-OccupationId terpilih.
func PilihOkupasi(h *Halaman, o, i int, id string) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	for _, b := range OpsiOkupasi(h, it["IndexObject"]) {
		if b["OccupationId"] == id {
			it["OccupationName"], it["OccupationId"] = b["OccupationName"], b["OccupationId"]
			return nil
		}
	}
	return fmt.Errorf("%w: okupasi %q", ErrBarisTidakAda, id)
}
