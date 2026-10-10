package models

// Untuk apa berkas ini: LAYAR KOMITE - flow action `ViewTransferDtl`, Section `ShowTransfer` (korpus `Komite Claim
// FacIn`) beserta section rinciannya `SpreadingDetail` (Object Detail expandPane), `DetailAdjustmentFac` (grid adjustment
// KMT expandPane), `ShowRetro_Sec` + `ShowSecurityReinsurer` (local action ShowRetro "View Retro", pra-proses
// `PreShowRetro_Act` / `PreSecurityReas_Act`). Urutan bagian, label, dan syarat tampil diambil dari section itu
// (VERBATIM). Nilai dibaca dari kasus klaim induk lewat kontrak (`pyWorkCover`) dan pra-proses `SetValueKomite`.
//
// Tata letak (keputusan work owner 09-10 Komite Prop): ubin ringkasan, kartu berjudul, pasangan label-nilai; tangga dan
// kartu keputusan DI BAWAH semua rincian. Tidak dibangun (PARITAS): sel `VIS 1=2` (LS2 CLAIM No., LS51 / LS52 / LS22 /
// LS23, Spreading Claim LS40 / LS44), ikon grid `pzPegaDefaultGridIcons`.

import (
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/kontrak"
)

// Jenis tampilan nilai (format di layar: angka 4 desimal `pxNumber DecimalPlaces=4`, tanggal dd-mm-yyyy).
const (
	JenisTeks        = "teks"
	JenisTeksPanjang = "teksPanjang"
	JenisAngka       = "angka"
	JenisTanggal     = "tanggal"
	JenisTanggalJam  = "tanggalJam"
	JenisCentang     = "centang"
	// JenisTautan - pxLink (nilai = kunci modal; kosong = tautan tidak tampil).
	JenisTautan = "tautan"
)

// Medan - satu sel hanya-baca berlabel.
type Medan struct {
	Label string `json:"label"`
	Nilai string `json:"nilai"`
	Jenis string `json:"jenis"`
}

// KolomGrid - satu kolom grid (judul VERBATIM).
type KolomGrid struct {
	Label    string `json:"label"`
	Properti string `json:"properti"`
	Jenis    string `json:"jenis"`
}

// Grid - satu grid hanya-baca. Rincian[i] = bagian expandPane baris ke-i (masterDetail); Kaki = baris kaki grid.
type Grid struct {
	Judul   string              `json:"judul"`
	Kolom   []KolomGrid         `json:"kolom"`
	Baris   []map[string]string `json:"baris"`
	Rincian [][]Bagian          `json:"rincian,omitempty"`
	Kaki    []Medan             `json:"kaki,omitempty"`
}

// Bagian - satu blok layar (urutan = urutan Section).
type Bagian struct {
	Kunci string  `json:"kunci"`
	Judul string  `json:"judul,omitempty"`
	Medan []Medan `json:"medan,omitempty"`
	Grid  []Grid  `json:"grid,omitempty"`
}

// Tombol - satu tombol layar.
type Tombol struct {
	Label  string `json:"label"`
	Aksi   string `json:"aksi"`
	Aktif  bool   `json:"aktif"`
	Alasan string `json:"alasan,omitempty"`
}

// IsianLayar - nilai awal isian keputusan dan aturan tampilnya.
type IsianLayar struct {
	Nilai Keputusan `json:"nilai"`
	// TampilUsul - dua Propose `VIS .TransferType = 2`; Terbuka - `NA .KomiteCount != '1'` salah.
	TampilUsul    bool              `json:"tampilUsul"`
	Terbuka       bool              `json:"terbuka"`
	PilihanTerima []Pilihan         `json:"pilihanTerima"`
	Label         map[string]string `json:"label"`
}

// Layar - satu kasus komite siap ditampilkan.
type Layar struct {
	Kasus  Kasus      `json:"kasus"`
	Judul  []string   `json:"judul"`
	Ubin   []Medan    `json:"ubin"`
	Bagian []Bagian   `json:"bagian"`
	Isian  IsianLayar `json:"isian"`
	Tombol []Tombol   `json:"tombol"`
	// Modal - isi pop-up (local action ShowRetro) per kunci tautan.
	Modal map[string][]Bagian `json:"modal,omitempty"`
	// BolehKerja - pelaku memegang tingkat berjalan dan kasus terbuka.
	BolehKerja bool     `json:"bolehKerja"`
	Pesan      []string `json:"pesan,omitempty"`
}

// Label tombol dan judul (VERBATIM ShowTransfer).
const (
	JudulKomite = "CLAIM COMMITTEE -"
	TombolLihat = "View more details"
	TombolBatal = "Cancel"
	TombolKirim = "Submit"
)

// "View more details" = klaim induk hanya-baca (`bukaKasus.hanyaLihat`, `GET /api/claim-fac-in/kasus/{id}?lihat=1`;
// pola Komite Claim Prop - harness klaim `pyWorkCover` tidak diekspor korpus ini).

func kol(label, prop, jenis string) KolomGrid {
	return KolomGrid{Label: label, Properti: prop, Jenis: jenis}
}

func salinDaftar(rows []map[string]string) []map[string]string {
	out := make([]map[string]string, 0, len(rows))
	for _, b := range rows {
		out = append(out, salinBaris(b))
	}
	return out
}

func nomori(rows []map[string]string) []map[string]string {
	for i, b := range rows {
		b["No"] = strconv.Itoa(i + 1) // GridNumbering (nomor baris 1..n)
	}
	return rows
}

// SusunLayar menyusun layar `ShowTransfer` kasus `k` atas klaim induk `kl` dan pra-proses `pr`. `perluasan` = anggota
// tangga KCF-02 yang akan disimpan tingkat 1 (ditampilkan, belum ditulis); `akun` / `peran` = pelaku.
func SusunLayar(k Kasus, kl kontrak.KlaimFacIn, pr PraProses, perluasan []Anggota, namaJenisReas map[string]string,
	akun string, peran []string) Layar {
	v := func(j string) string { return kl.Nilai[j] }
	ly := Layar{Kasus: k, Judul: []string{JudulKomite}, BolehKerja: k.Pemegang(akun, peran), Modal: map[string][]Bagian{}}
	if j := JudulTransfer[k.TransferType]; j != "" { // LS1
		ly.Judul = append(ly.Judul, j)
	}
	// ubin ringkasan `[tidak ada di korpus]` (pola Komite Non Prop): label sejajar kolom tabel Committee inbox Claim Fac In
	ly.Ubin = []Medan{
		{Label: "Committee No", Nilai: k.ID, Jenis: JenisTeks},
		{Label: "Claim ID", Nilai: k.KlaimID, Jenis: JenisTeks},
		{Label: "Claim No", Nilai: v("ClaimData.NoClaim"), Jenis: JenisTeks},
		{Label: "Level", Nilai: strconv.Itoa(k.Count) + " / " + strconv.Itoa(k.Loop+len(perluasan)), Jenis: JenisTeks},
	}
	tt2 := k.TransferType == TransferAdjustment
	if tt2 {
		// LS3-LS9: grid item objek kasus komite (TempObjectList = ObjectList(IndexObject)); TT3 / TT4 tanpa pra-proses.
		polis := []Medan{}
		treaty := v("IsTreatyIn") == "1"
		if !treaty { // LS6 `VIS pyWorkCover.IsTreatyIn == 0`
			polis = append(polis, Medan{Label: "Name of Insured", Nilai: v(OQ + "InsuredName"), Jenis: JenisTeks})
		}
		polis = append(polis,
			Medan{Label: "QQ Name", Nilai: v(OQ + "QQName"), Jenis: JenisTeks},
			Medan{Label: "Line of Business", Nilai: v(OQ + "BusinessName"), Jenis: JenisTeks},
			Medan{Label: "Policy No.", Nilai: v("OfferFacIn.PolicyData.PolicyNo"), Jenis: JenisTeks},
			Medan{Label: "Ceding Co Name", Nilai: v(OQ + "CedingCoName"), Jenis: JenisTeks},
			Medan{Label: "SOB Name", Nilai: v(OQ + "SobName"), Jenis: JenisTeks})
		if !treaty { // LS8 `pyWorkCover.Policy.StartDateTime` "-" EndDateTime (Policy = OfferFacIn.PolicyData, [inferensi])
			polis = append(polis, Medan{Label: "Period", Nilai: periode(v("OfferFacIn.PolicyData.StartDateTime"),
				v("OfferFacIn.PolicyData.EndDateTime")), Jenis: JenisTeks})
		}
		ly.Bagian = append(ly.Bagian, Bagian{Kunci: "polis", Judul: "Policy Detail", Medan: polis})
		ly.Bagian = append(ly.Bagian, bagianObjek(kl, namaJenisReas)) // LS9 VIS TransferType = 2
		ly.Bagian = append(ly.Bagian, bagianKlaim(k, kl, pr))
		ly.Bagian = append(ly.Bagian, bagianAdjKomite(kl, pr, namaJenisReas, ly.Modal)) // LS24
		ly.Bagian = append(ly.Bagian,
			Bagian{Kunci: "riwayat", Judul: "History Adjustment", Grid: []Grid{{Kolom: []KolomGrid{ // LS28
				kol("No", "No", JenisTeks), kol("Payment Type", "PaymentTypeLabel", JenisTeks),
				kol("Currency", "Currency", JenisTeks), kol("Adjustment Gross (100%)", "GrossAdjustment", JenisAngka),
				kol("Adjustment RNM", "AdjustmentValue", JenisAngka), kol("Status", "StatusLabel", JenisTeks),
				kol("Komite No", "KomiteID", JenisTeks), kol("Accepted Date", "AcceptedDate", JenisTanggal),
				kol("Accepted No", "AcceptedNo", JenisTeks)}, Baris: barisRiwayat(pr.Riwayat)}}},
			Bagian{Kunci: "total", Judul: "Total Adjustment", Grid: []Grid{{Kolom: []KolomGrid{ // LS33
				kol("Currency", "Currency", JenisTeks), kol("Total Gross (100%)", "AdjustmentGross", JenisAngka),
				kol("Total Adjustment RNM", "AdjustmentValue", JenisAngka)}, Baris: salinDaftar(pr.Total)}}})
	}
	ly.Bagian = append(ly.Bagian, bagianTeksKomite(k, kl)) // LS38-LS39
	if tt2 {                                               // LS41 "List of Committee" `.TransferType == 2`
		ly.Bagian = append(ly.Bagian, Bagian{Kunci: "tangga", Judul: "List of Committee", Grid: []Grid{{Kolom: []KolomGrid{
			kol("No", "No", JenisTeks), kol("Committee", "jabatan", JenisTeks), kol("Status", "keputusan", JenisTeks),
			kol("Date Approve", "tanggal", JenisTanggalJam), kol("Comment", "komentar", JenisTeks)},
			Baris: barisTangga(append(append([]Anggota{}, k.Tangga...), perluasan...))}}})
	}
	ly.Isian = IsianLayar{Nilai: nilaiAwal(k), TampilUsul: tt2, Terbuka: IsianTerbuka(k), PilihanTerima: LabelTerima,
		Label: map[string]string{"acceptStatus": LabelAcceptStatus, "usulTutup": LabelProposeClose,
			"usulCadang": LabelProposeReserved, "comment": LabelNote}}
	kirim := Tombol{Label: TombolKirim, Aksi: "putuskan", Aktif: ly.BolehKerja}
	if tt2 && Adjustment(kl)["AcceptedNo"] != "" { // LS48 `NA ... || pyWorkPage.Adjustment.AcceptedNo != ''`
		kirim.Aktif, kirim.Alasan = false, "Adjustment already accepted"
	}
	ly.Tombol = []Tombol{{Label: TombolLihat, Aksi: "lihat", Aktif: true}, {Label: TombolBatal, Aksi: "batal", Aktif: true},
		kirim}
	if len(ly.Modal) == 0 {
		ly.Modal = nil
	}
	return ly
}

// periode - LS8 "<mulai> - <akhir>"; kedua medan `FormatType=date` (Claim Fac In menyimpan "2006-01-02", NormalkanTanggalPolis)
// tampil dd-mm-yyyy; nilai lain apa adanya.
func periode(mulai, akhir string) string {
	if strings.TrimSpace(mulai) == "" && strings.TrimSpace(akhir) == "" {
		return ""
	}
	return tanggalTampil(mulai) + " - " + tanggalTampil(akhir)
}

func tanggalTampil(s string) string {
	s = strings.TrimSpace(s)
	if len(s) < 10 {
		return s
	}
	if d, err := time.Parse("2006-01-02", s[:10]); err == nil {
		return d.Format("02-01-2006")
	}
	return s
}

// namaTreaty - kolom Treaty Type dropdown `BrowseReinsuranceType_RD` (`REINSURANCETYPE.NOTE`): nama tersimpan, nama
// master, atau kode.
func namaTreaty(b map[string]string, master map[string]string) string {
	if n := strings.TrimSpace(b["TreatyName"]); n != "" {
		return n
	}
	if n := master[b["TreatyType"]]; n != "" {
		return n
	}
	return b["TreatyType"]
}

// bagianObjek - LS9 "Object Detail" (item objek) + rincian `SpreadingDetail` (Spreading Policy, VIS !IsTravel).
func bagianObjek(kl kontrak.KlaimFacIn, master map[string]string) Bagian {
	items := salinDaftar(kl.Daftar[DaftarItem(kl.Objek)])
	g := Grid{Kolom: []KolomGrid{kol("No", "No", JenisTeks), kol("Name", "ObjectItemName", JenisTeks),
		kol("Coverage Note", "CoverageNote", JenisTeks), kol("Curr", "Currency", JenisTeks),
		kol("Curr Value in IDR", "KursObjectItem", JenisAngka), kol("TSI Object", "TSIPerObject", JenisAngka),
		kol("TSI RNM", "TSINusare", JenisAngka), kol("TSI RNM in IDR", "ValueTSINusareIDR", JenisAngka)},
		Baris: nomori(items)}
	for i, it := range items {
		var rinci []Bagian
		if !IsTravel(kl.Nilai) { // LS35 `VIS? !IsTravel`
			rows := nomori(salinDaftar(kl.Daftar[DaftarDiItem(kl.Objek, i+1, AnakSpreadPL)]))
			for _, s := range rows {
				s["TreatyLabel"] = namaTreaty(s, master)
			}
			rinci = append(rinci, Bagian{Kunci: "spreadingPolis", Grid: []Grid{{Judul: "Spreading Policy",
				Kolom: []KolomGrid{kol("Spreading No", "No", JenisTeks), kol("Treaty Type", "TreatyLabel", JenisTeks),
					kol("Share %", "SharePercentage", JenisAngka), kol("TSI Spreaded", "TSISpreaded", JenisAngka)},
				Baris: rows, Kaki: []Medan{{Label: "Total", Nilai: it["TotalSharePercentage"], Jenis: JenisAngka}}}}},
				Bagian{Kunci: "spreadingTotal", Medan: []Medan{ // LS39
					{Label: "Total Share Percentage", Nilai: it["TotalSharePercentage"], Jenis: JenisAngka},
					{Label: "Total TSI Spreaded", Nilai: it["TotalTSISpreaded"], Jenis: JenisAngka}}})
		}
		g.Rincian = append(g.Rincian, rinci)
	}
	return Bagian{Kunci: "objek", Judul: "Object Detail", Grid: []Grid{g}}
}

// bagianKlaim - LS14 "Claim Details" (badan grid DataTempAdj: tampil bila ada adjustment kasus komite) + LS19-LS21.
func bagianKlaim(k Kasus, kl kontrak.KlaimFacIn, pr PraProses) Bagian {
	v := func(j string) string { return kl.Nilai[j] }
	if len(pr.AdjKomite) == 0 {
		return Bagian{Kunci: "klaim", Judul: "Claim Details"}
	}
	m := []Medan{
		{Label: "Input Date", Nilai: FormatWaktu(k.TglCreate), Jenis: JenisTanggal},                           // S12 Komite.DateOfComitee
		{Label: "PIC Name", Nilai: pr.Inisial, Jenis: JenisTeks},                                              // S8.2 / S12 Komite.Initial
		{Label: "Claim No", Nilai: v("ClaimData.NoClaim") + " / " + v("ClaimData.ClaimNo"), Jenis: JenisTeks}, // LS50
		{Label: "Occupation", Nilai: pr.Okupasi, Jenis: JenisTeks},
		{Label: "Date of Loss", Nilai: v("ClaimData.DateOfLoss"), Jenis: JenisTanggal},
		{Label: "Cause of Loss", Nilai: v("ClaimData.CauseOfLoss"), Jenis: JenisTeks},
		{Label: "Location of Loss", Nilai: v("ClaimData.Location"), Jenis: JenisTeks},
		{Label: "Catastrophe", Nilai: v("ClaimData.StsKatastrofe"), Jenis: JenisTeks},
	}
	if v("ClaimData.StsKatastrofe") == "Non-Catastrophe" { // LS20
		m = append(m, Medan{Label: "Non Catastrophe Type", Nilai: v("ClaimData.NonKatastrofeType"), Jenis: JenisTeks})
	}
	if v("ClaimData.StsKatastrofe") == "Catastrophe" || v("ClaimData.NonKatastrofeType") == "Big Claim" { // LS21
		m = append(m, Medan{Label: "Catastrophe Note", Nilai: v("ClaimData.KatastrofeNote"), Jenis: JenisTeks})
	}
	return Bagian{Kunci: "klaim", Judul: "Claim Details", Medan: m}
}

// bagianAdjKomite - LS24 grid DataTempAdj (`.pyNote`) + rincian `DetailAdjustmentFac`.
func bagianAdjKomite(kl kontrak.KlaimFacIn, pr PraProses, master map[string]string, modal map[string][]Bagian) Bagian {
	g := Grid{Kolom: []KolomGrid{kol("No", "No", JenisTeks), kol("Adjustment", "pyNote", JenisTeks)},
		Baris: nomori(salinDaftar(pr.AdjKomite))}
	for n, b := range pr.AdjKomite {
		pos := pr.Posisi[n]
		g.Rincian = append(g.Rincian, rincianAdjustment(kl, b, pos[0], pos[1], master, modal))
	}
	return Bagian{Kunci: "adjustment", Grid: []Grid{g}}
}

// rincianAdjustment = Section DetailAdjustmentFac.
func rincianAdjustment(kl kontrak.KlaimFacIn, b map[string]string, i, a int, master map[string]string,
	modal map[string][]Bagian) []Bagian {
	cur := b["Currency"]
	baris3 := func(label, cur1, v1, cur2, v2 string) map[string]string {
		return map[string]string{"Deskripsi": label, "Cur1": cur1, "Nilai1": v1, "Cur2": cur2, "Nilai2": v2}
	}
	kolom3 := func(l1, l2 string) []KolomGrid {
		return []KolomGrid{kol("DESCRIPTION", "Deskripsi", JenisTeks), kol("CURRENCY", "Cur1", JenisTeks),
			kol(l1, "Nilai1", JenisAngka), kol("CURRENCY", "Cur2", JenisTeks), kol(l2, "Nilai2", JenisAngka)}
	}
	out := []Bagian{{Kunci: "penanda", Medan: []Medan{ // LS1
		{Label: "Net For Collection", Nilai: b["NetForCollection"], Jenis: JenisCentang},
		{Label: "Direct To Kasir", Nilai: b["DirectToKasir"], Jenis: JenisCentang}}}}
	switch b["PaymentType"] {
	case "1", "2", "5": // LS2 "Claim Adjustment"
		out = append(out, Bagian{Kunci: "nilai", Judul: "Claim Adjustment", Grid: []Grid{{Kolom: []KolomGrid{
			kol("DESCRIPTION", "Deskripsi", JenisTeks), kol("CURRENCY", "Cur0", JenisTeks),
			kol("ESTIMATION RNM", "Nilai0", JenisAngka), kol("CURRENCY", "Cur1", JenisTeks),
			kol("ADJUSTMENT 100 %", "Nilai1", JenisAngka), kol("CURRENCY", "Cur2", JenisTeks),
			kol("ADJUSTMENT RNM", "Nilai2", JenisAngka)}, Baris: []map[string]string{
			{"Deskripsi": "Gross Adjustment(100%)", "Cur1": cur, "Nilai1": b["GrossAdjustment"], "Cur2": cur,
				"Nilai2": b["GrossValue"]},
			{"Deskripsi": "Deductible Value", "Cur1": cur, "Nilai1": b["IndividualRiskValue"], "Cur2": cur,
				"Nilai2": b["IndividualRiskRNM"]},
			{"Deskripsi": "Claim Value", "Cur0": b["CurrencyEstimasi"], "Nilai0": b["EstimationValue"], "Cur1": cur,
				"Nilai1": b["ProposeAdjustmentValue"], "Cur2": cur, "Nilai2": b["AdjustmentValue"]},
			{"Deskripsi": "Total Claim Value", "Cur1": cur, "Nilai1": b["ProposeAdjustmentValue"], "Cur2": cur,
				"Nilai2": b["AdjustmentValue"]}}}}})
	case "4", "6": // LS3 "Adjuster Fee"
		out = append(out, Bagian{Kunci: "nilai", Judul: "Adjuster Fee", Grid: []Grid{{Kolom: append([]KolomGrid{
			kol("DESCRIPTION", "Deskripsi", JenisTeks), kol("%", "Persen", JenisAngka)}, kolom3("AMOUNT 100%",
			"AMOUNT RNM")[1:]...), Baris: []map[string]string{
			baris3("Gross Adjustment(100%)", cur, b["GrossAdjustment"], cur, b["GrossValue"]),
			{"Deskripsi": "VAT", "Persen": b["VAT"], "Cur1": cur, "Nilai1": b["VATValue"], "Cur2": cur,
				"Nilai2": b["IndividualRiskRNM"]},
			baris3("Professional Fee RNM", cur, b["ProfessionalFee"], cur, b["AdjusterFeeValue"]),
			baris3("Total Adjuster Fee RNM", cur, b["ProfessionalFee"], cur, b["AdjusterFeeValue"])}}}})
	case "3": // LS4 "Salvage"
		out = append(out, Bagian{Kunci: "nilai", Judul: "Salvage", Grid: []Grid{{Kolom: append([]KolomGrid{
			kol("DESCRIPTION", "Deskripsi", JenisTeks), kol("%", "Persen", JenisAngka)}, kolom3("AMOUNT 100%",
			"AMOUNT RNM")[1:]...), Baris: []map[string]string{
			baris3("Gross Adjustment(100%)", cur, b["GrossAdjustment"], cur, b["GrossValue"]),
			baris3("Total Salvage  RNM", cur, b["GrossAdjustment"], cur, b["SalvageValue"])}}}})
	}
	spread := nomori(salinDaftar(kl.Daftar[DaftarDiAdj(kl.Objek, i, a, AnakSpread)]))
	var h hitung
	totPersen, totKlaim := apd.New(0, 0), apd.New(0, 0)
	for n, s := range spread {
		s["TreatyLabel"] = namaTreaty(s, master)
		totPersen = h.tambah(totPersen, h.dari("SharePercentage", s["SharePercentage"]))
		totKlaim = h.tambah(totKlaim, h.dari("ClaimSpreaded", s["ClaimSpreaded"]))
		if s["TreatyType"] == TreatyRetro { // "View Retro" `VIS .TreatyType = '10015'`
			kunci := "retro:" + strconv.Itoa(i) + ":" + strconv.Itoa(a) + ":" + strconv.Itoa(n+1)
			s["ViewRetro"] = kunci
			modal[kunci] = modalRetro(kl, s)
		}
	}
	if !IsTravel(kl.Nilai) { // LS7 `VIS? !IsTravel`
		out = append(out, Bagian{Kunci: "spreading", Grid: []Grid{{Judul: "Spreading Adjustment", Kolom: []KolomGrid{
			kol("No", "No", JenisTeks), kol("Treaty Type", "TreatyLabel", JenisTeks),
			kol("Share %", "SharePercentage", JenisAngka), kol("Claim Spreaded", "ClaimSpreaded", JenisAngka),
			kol("", "ViewRetro", JenisTautan)}, Baris: spread, Kaki: []Medan{ // body3 CountSpread.CARI27 / CARI26
			{Label: "Total Share(%)", Nilai: TeksAngka(totPersen), Jenis: JenisAngka},
			{Label: "Total Claim Spread", Nilai: TeksAngka(totKlaim), Jenis: JenisAngka}}}}},
			Bagian{Kunci: "spreadingTotal", Medan: []Medan{ // LS11
				{Label: "Total Share Percentage(%)", Nilai: b["TotalSharePersen"], Jenis: JenisAngka},
				{Label: "Total Claim Spread", Nilai: b["TotalSpreadAdjustment"], Jenis: JenisAngka}}})
	}
	qs := nomori(salinDaftar(kl.Daftar[DaftarDiAdj(kl.Objek, i, a, AnakQS)]))
	if len(qs) > 0 && qs[0]["TreatyName"] != "" { // LS12 / LS16 `.SpreadingQuotaShare(1).TreatyName != ''`
		qp, qk := apd.New(0, 0), apd.New(0, 0)
		for _, s := range qs {
			s["TreatyLabel"] = namaTreaty(s, master)
			qp = h.tambah(qp, h.dari("SharePercentage", s["SharePercentage"]))
			qk = h.tambah(qk, h.dari("ClaimSpreaded", s["ClaimSpreaded"]))
		}
		out = append(out, Bagian{Kunci: "qs", Grid: []Grid{{Judul: "BreakDown Spreading Quota Share (QS)",
			Kolom: []KolomGrid{kol("No", "No", JenisTeks), kol("Treaty Type", "TreatyLabel", JenisTeks),
				kol("Share %", "SharePercentage", JenisAngka), kol("Claim Spreaded", "ClaimSpreaded", JenisAngka)},
			Baris: qs, Kaki: []Medan{{Label: "Total Share (%)", Nilai: TeksAngka(qp), Jenis: JenisAngka},
				{Label: "Total Claim Spread", Nilai: TeksAngka(qk), Jenis: JenisAngka}}}}},
			Bagian{Kunci: "qsTotal", Medan: []Medan{
				{Label: "Total Share Percentage(%)", Nilai: b["TotalSharePersen"], Jenis: JenisAngka},
				{Label: "Total Claim Spread", Nilai: b["TotalSpreadBreakQs"], Jenis: JenisAngka}}})
	}
	bayar := []Medan{{Label: "Payable To", Nilai: labelKode(LabelPayable, b["Payable"]), Jenis: JenisTeks}, // LS18
		{Label: "Specify", Nilai: b["PayableTo"], Jenis: JenisTeks},
		{Label: "Name of Bank", Nilai: b["NameOfBank"], Jenis: JenisTeks}}
	if b["SwiftCode"] != "" { // LS19 `VIS .SwiftCode != ''`
		bayar = append(bayar, Medan{Label: "Swift Code", Nilai: b["SwiftCode"], Jenis: JenisTeks})
	}
	bayar = append(bayar, Medan{Label: "Branch of Bank", Nilai: b["BranchOfBank"], Jenis: JenisTeks},
		Medan{Label: "Account No", Nilai: b["NoAccount"], Jenis: JenisTeks})
	return append(out, Bagian{Kunci: "bayar", Medan: bayar})
}

// modalRetro = local action ShowRetro (pra-proses PreShowRetro_Act S2.2-S3) atas satu baris spreading 10015: reasuradur
// FacRetroList polis, Share "x% of y%", Amount = ClaimSpreaded x PctShareAllObj / 100 (4 desimal), Currency; rincian
// ShowSecurityReinsurer (PreSecurityReas_Act S2.2-S3): security reinsurer baris FacRetroList ber-ReinsurerID sama,
// Premium = Amount x PctShare / 100 (2 desimal).
func modalRetro(kl kontrak.KlaimFacIn, s map[string]string) []Bagian {
	var h hitung
	seratus := apd.New(100, 0)
	retro := kl.Daftar[DaftarRetroPolis]
	g := Grid{Kolom: []KolomGrid{kol("Reinsurer Name", "ReinsurerName", JenisTeks), kol("Share", "Comment", JenisTeks),
		kol("Currency", "Attention", JenisTeks), kol("Amount", "TotalClaim", JenisAngka)}}
	for n, r := range retro {
		total := h.bagiBulat(h.kali(h.dari("ClaimSpreaded", s["ClaimSpreaded"]), h.dari("PctShareAllObj",
			r["PctShareAllObj"])), seratus, 4)
		share := TeksAngka(h.bagiBulat(h.dari("SharePercentage", s["SharePercentage"]), apd.New(1, 0), 2)) + "% of " +
			TeksAngka(h.bagiBulat(h.dari("PctShareAllObj", r["PctShareAllObj"]), apd.New(1, 0), 2)) + "%"
		g.Baris = append(g.Baris, map[string]string{"ReinsurerName": r["ReinsurerName"], "Comment": share,
			"Attention": s["Currency"], "TotalClaim": TeksAngka(total)})
		sec := Grid{Kolom: []KolomGrid{kol("Security Reinsurer", "ReinsurerName", JenisTeks),
			kol("Share", "PctShare", JenisAngka), kol("Premium", "Premium", JenisAngka)}}
		for _, x := range retro { // S2.2.1 Local.IDBroker == .ReinsurerID
			if x["ReinsurerID"] != r["ReinsurerID"] {
				continue
			}
			for _, sr := range kl.Daftar[anak(anak(DaftarRetroPolis, n+1, "PrintRISlip.FacOfferList"), 1,
				"SecurityReinsurer")] {
				prem := h.bagiBulat(h.kali(total, h.dari("PctShare", sr["PctShare"])), seratus, 2)
				sec.Baris = append(sec.Baris, map[string]string{"ReinsurerName": sr["ReinsurerName"],
					"PctShare": sr["PctShare"], "Premium": TeksAngka(prem)})
			}
			break
		}
		g.Rincian = append(g.Rincian, []Bagian{{Kunci: "security", Grid: []Grid{sec}}})
	}
	if h.err != nil {
		return []Bagian{{Kunci: "retro", Judul: "ShowRetro", Medan: []Medan{{Label: "Error", Nilai: h.err.Error(),
			Jenis: JenisTeks}}}}
	}
	return []Bagian{{Kunci: "retro", Judul: "ShowRetro", Grid: []Grid{g}}}
}

// bagianTeksKomite - LS39 (`VIS NOTBLANK`, hanya-baca). TT2: `pyWorkPage.Komite.*` = `DataCommitteFacin.*` adjustment
// (CreateKMTNo_Act 7). TT3 / TT4: Legal / Chronology / Extent = teks pop-up tersimpan di kepala kasus komite (7.2,
// migrasi 643, jawaban work owner 10-10-2026 OQ-KCFI-03); `Komite.Remarks` = `ClaimData.Remark` klaim induk (2 / 7.2).
func bagianTeksKomite(k Kasus, kl kontrak.KlaimFacIn) Bagian {
	var m []Medan
	tambah := func(label, nilai string) {
		if strings.TrimSpace(nilai) != "" {
			m = append(m, Medan{Label: label, Nilai: nilai, Jenis: JenisTeksPanjang})
		}
	}
	if k.TransferType == TransferAdjustment {
		b := Adjustment(kl)
		d := func(p string) string { return b["DataCommitteFacin."+p] }
		tambah("Legal Liability / Policy Liability", d("LegalLiability"))
		tambah("Chronology", d("CircumCauseOfLoss"))
		tambah("Extent Of Loss", d("ExtentOfLoss"))
		if b["PaymentType"] == "4" { // `pyWorkPage.Type == '4' && .Komite.AdjusterFee != ''` ([inferensi] Type = PaymentType)
			tambah("Adjuster Fee", d("AdjusterFee"))
		}
		tambah("Salvage", d("Salvage"))
		tambah("Remarks", d("Remarks"))
	} else {
		tambah("Legal Liability / Policy Liability", k.Liability)
		tambah("Chronology", k.Kronologi)
		tambah("Extent Of Loss", k.Extent)
		tambah("Remarks", kl.Nilai["ClaimData.Remark"])
	}
	return Bagian{Kunci: "teksKomite", Medan: m}
}

// barisRiwayat - LS29: label Payment Type / Status.
func barisRiwayat(rows []map[string]string) []map[string]string {
	out := nomori(salinDaftar(rows))
	for _, b := range out {
		b["PaymentTypeLabel"] = labelKode(LabelJenisBayar, b["PaymentType"])
		b["StatusLabel"] = LabelStatusBaris(b["AcceptanceStatus"])
	}
	return out
}

// barisTangga - grid "List of Committee": jabatan roster (`IDKomite`), keputusan, tanggal, komentar.
func barisTangga(t []Anggota) []map[string]string {
	out := make([]map[string]string, 0, len(t))
	for i, a := range t {
		out = append(out, map[string]string{"No": strconv.Itoa(i + 1), "jabatan": a.Jabatan,
			"keputusan": labelKode(LabelKeputusanAnggota, a.Keputusan), "tanggal": a.Tanggal, "komentar": a.Komentar})
	}
	return out
}

// nilaiAwal - isian yang tersimpan di `pyWorkPage` dari tingkat sebelumnya (Pega tidak mengosongkannya antar tingkat;
// `KomitePostAct` S6 ber-remark): `.AcceptStatus` terakhir, `.Comment` = komentar tingkat sebelumnya, dua Propose = nilai
// tersimpan di header.
func nilaiAwal(k Kasus) Keputusan {
	n := Keputusan{AcceptStatus: k.AcceptStatus, UsulTutup: k.UsulTutup == UsulYa, UsulCadang: k.UsulCadang == UsulYa}
	if i := k.barisBerjalan(); i > 0 {
		n.Comment = k.Tangga[i-1].Komentar
	}
	return n
}
