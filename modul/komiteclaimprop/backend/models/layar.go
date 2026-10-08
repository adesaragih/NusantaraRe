package models

// Untuk apa berkas ini: LAYAR KOMITE - flow action `ViewTransferDtl`, Section `ShowTransfer` wajah TT 2 (ADJUSTMENT).
// Urutan bagian, label, dan syarat tampil diambil dari `Section/ShowTransfer.xml` (VERBATIM, termasuk ejaan
// "Dedutible Type" dan "Committe"). Nilai dibaca dari kasus klaim induk lewat kontrak (`pyWorkCover`, dan
// `pyWorkPage.Adjustment` / `pyWorkPage.Komite` - salinan baris adjustment yang dibuat `AddKomiteTreatyChild_ACT`
// langkah 13-17; baris itu beku selama diserahkan ke komite, jadi baris aslinya yang dibaca).
//
// ⚠️ Tidak dibangun (PARITAS):
//   - judul "CLOSE" (TT 4, OQ-CP-06) dan "REJECT" (TT 3, tanpa penulis di Claim Prop);
//   - sel ber-`pyVisible` NEVER: No Claim di kepala "Claim Treaty", "sample text", Total Original Currency Gross
//     Estimate(100%), Total Original Currency Estimation (blok Total Estimation In IDR), grid "Spreading Claim" (1=2);
//   - sel yang di XML dapat disunting tetapi tidak pernah disimpan `KomitePost*` (Payment Type / Komite No History
//     Adjustment, Share / Claim Spreaded Spreading Out, Currency Total Original Currency Estimation) tampil hanya-baca.

import (
	"strings"

	"nusantarare/inti/backend/kontrak"
)

// Jenis tampilan nilai (format di layar: angka 4 desimal `pyDecimalPlaces 4`, tanggal `Date-Short-Custom-YYYY`,
// tanggal-jam `DateTime-Short-YYYY-Custom`).
const (
	JenisTeks        = "teks"
	JenisTeksPanjang = "teksPanjang"
	JenisAngka       = "angka"
	JenisTanggal     = "tanggal"
	JenisTanggalJam  = "tanggalJam"
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

// Grid - satu grid hanya-baca.
type Grid struct {
	Judul string              `json:"judul"`
	Kolom []KolomGrid         `json:"kolom"`
	Baris []map[string]string `json:"baris"`
}

// Bagian - satu blok layar (urutan = urutan Section).
type Bagian struct {
	Kunci string  `json:"kunci"`
	Judul string  `json:"judul,omitempty"`
	Medan []Medan `json:"medan,omitempty"`
	Grid  []Grid  `json:"grid,omitempty"`
	// Sel - tabel bebas (blok deductible 6 x 4); sel tanpa nilai = Medan kosong.
	Sel [][]Medan `json:"sel,omitempty"`
}

// Pilihan - satu opsi dropdown.
type Pilihan struct {
	Nilai string `json:"nilai"`
	Label string `json:"label"`
}

// Tombol - satu tombol layar.
type Tombol struct {
	Label  string `json:"label"`
	Aksi   string `json:"aksi"`
	Aktif  bool   `json:"aktif"`
	Alasan string `json:"alasan,omitempty"`
}

// IsianLayar - nilai awal isian keputusan dan aturan tampilnya (dievaluasi ulang di klien saat isian berubah, sama
// seperti `pyIsClientWhen true` di XML).
type IsianLayar struct {
	Nilai Keputusan `json:"nilai"`
	// Terbuka - `pyDisabledWhen .KomiteCount!='1'` salah: Subjectivity, catatannya, dan dua Propose aktif.
	Terbuka       bool              `json:"terbuka"`
	PilihanTerima []Pilihan         `json:"pilihanTerima"`
	Label         map[string]string `json:"label"`
}

// Layar - satu kasus komite siap ditampilkan.
type Layar struct {
	Kasus  Kasus      `json:"kasus"`
	Judul  []string   `json:"judul"`
	Bagian []Bagian   `json:"bagian"`
	Isian  IsianLayar `json:"isian"`
	Tombol []Tombol   `json:"tombol"`
	// BolehKerja - pelaku memegang assignment (KomiteRouter S6.1) dan kasus terbuka.
	BolehKerja bool     `json:"bolehKerja"`
	Pesan      []string `json:"pesan,omitempty"`
}

// Label tombol dan judul (VERBATIM ShowTransfer).
const (
	JudulKomite      = "CLAIM COMMITTEE -"
	JudulAdjustment  = "ADJUSTMENT"
	TombolLihat      = "View more details"
	TombolBatal      = "Cancel"
	TombolKirim      = "Submit"
	JudulClaimTreaty = "Claim Analysis"
)

// LabelTerima - pilihan `.AcceptStatus` (`pyListSource associated`, prompt values tidak diekspor). Teks dari
// `SetDataAcceptationTreaty_Act` S2-S3 (`ParamDataTreaty.CARI7`): 1 "Approve", 2 "Reject".
var LabelTerima = []Pilihan{{Nilai: KeputusanSetuju, Label: "Approve"}, {Nilai: KeputusanTolak, Label: "Reject"}}

// AlasanLihatNonaktif - "View more details" membuka harness `ViewClaimFormKomite` yang tidak diekspor (prompt §6
// butir 11).
const AlasanLihatNonaktif = "Harness ViewClaimFormKomite tidak ada di ekspor Pega (OQ)"

// AlasanSudahBernomor - tombol Submit `pyDisabledWhen pyWorkPage.Adjustment.AcceptedNo != ”`.
const AlasanSudahBernomor = "Adjustment ini sudah memiliki Accepted No"

// labelTipeAdjustment - `.Adjustment.Type` (dropdown `associated`): teks `SetDataAcceptationTreaty_Act` S4-S7, sama
// dengan prompt values Claim Prop dari screenshot work owner 08-10-2026.
var labelTipeAdjustment = map[string]string{"1": "Claim", "2": "Adjuster Fee", "3": "Salvage", "4": "Consultant Fee"}

// labelKodeKlaim - dropdown `associated` panel klaim induk yang labelnya diberikan work owner (Claim Prop
// `LabelKode`, 08-10-2026). Kode lain tampil apa adanya.
var labelKodeKlaim = map[string]map[string]string{
	"ClaimData.ReportType":     {"1": "Direct", "2": "Via Email", "3": "Via Fax", "4": "via Postal Mail/Courier", "5": "Via Telephone"},
	"ClaimData.ReporterStatus": {"1": "Ceding Co Name", "2": "SOB Name", "3": "Others"},
	"IndividualRiskType":       {"0": "Select..", "1": "% From claims", "2": "% FromTSI", "3": "Other"},
}

func labelKode(peta map[string]string, v string) string {
	if l, ok := peta[v]; ok {
		return l
	}
	return v
}

// pembaca - nilai halaman klaim induk dan baris adjustment yang diputus.
type pembaca struct {
	kl  kontrak.KlaimTreaty
	adj map[string]string
}

func (p pembaca) v(j string) string { return p.kl.Nilai[j] }

func (p pembaca) medan(label, j, jenis string) Medan {
	return Medan{Label: label, Nilai: p.v(j), Jenis: jenis}
}

func (p pembaca) a(label, prop, jenis string) Medan {
	return Medan{Label: label, Nilai: p.adj[prop], Jenis: jenis}
}

func (p pembaca) daftar(j string) []map[string]string {
	rows := p.kl.Daftar[j]
	out := make([]map[string]string, 0, len(rows))
	for _, b := range rows {
		m := make(map[string]string, len(b))
		for k, v := range b {
			m[k] = v
		}
		out = append(out, m)
	}
	return out
}

func kol(label, prop, jenis string) KolomGrid {
	return KolomGrid{Label: label, Properti: prop, Jenis: jenis}
}

func tampilBila(m []Medan, ya bool, x ...Medan) []Medan {
	if ya {
		return append(m, x...)
	}
	return m
}

// SusunLayar menyusun layar `ShowTransfer` (TT 2) kasus `k` atas klaim induk `kl`; `total` = hasil
// `SetKomiteList_Act`, `akun` = pelaku.
func SusunLayar(k Kasus, kl kontrak.KlaimTreaty, total []TotalMataUang, akun string) Layar {
	p := pembaca{kl: kl, adj: AdjustmentKlaim(kl)}
	ly := Layar{Kasus: k, Judul: []string{JudulKomite, JudulAdjustment}, BolehKerja: k.Pemegang(akun)}

	// Panel "Claim Treaty" - kolom kiri (Stacked with labels left).
	kiri := []Medan{
		p.medan("Class of Business", "OfferFacIn.QuotationData.BusinessName", JenisTeks),
		p.medan("Ceding", "TreatyInMaster.Ceding", JenisTeks),
		p.medan("Source of Business", "TreatyInMaster.LeadingReinsSource", JenisTeks),
		p.medan("Policy No", "ClaimData.PolicyData.PolicyNo", JenisTeks),
		p.medan("Policy No Ceding", "ClaimData.PolicyNo", JenisTeks),
		p.medan("Insured Name", "ClaimData.InsuredName", JenisTeks),
		// Inline labels left: label = deskripsi properti (aturan properti tidak diekspor) -> nama properti.
		p.medan("StartDateTime", "ClaimData.PolicyData.StartDateTime", JenisTanggal),
		p.medan("EndDateTime", "ClaimData.PolicyData.EndDateTime", JenisTanggal),
		p.medan("Cause of Loss", "ClaimData.CauseOfLoss", JenisTeks),
		{Label: "Report Type", Nilai: labelKode(labelKodeKlaim["ClaimData.ReportType"], p.v("ClaimData.ReportType")),
			Jenis: JenisTeks},
	}
	kiri = tampilBila(kiri, p.v("ClaimData.AppointedADJ") != "",
		p.medan("Adjuster / Professional", "ClaimData.AppointedADJ", JenisTeks))
	kiri = tampilBila(kiri, p.v("ClaimData.ConsultantName") != "",
		p.medan("Consultant", "ClaimData.ConsultantName", JenisTeks))
	kiri = append(kiri, p.medan("Location of Loss", "ClaimData.Location", JenisTeks),
		// `.Adjustment.Type` berlabel "Policy No" (pyVisible ALWAYS) - VERBATIM.
		Medan{Label: "Policy No", Nilai: labelKode(labelTipeAdjustment, p.adj["Type"]), Jenis: JenisTeks},
		p.medan("Catastrophe", "ClaimData.StsKatastrofe", JenisTeks))
	kiri = tampilBila(kiri, p.v("ClaimData.StsKatastrofe") == "Non-Catastrophe",
		p.medan("NonKatastrofeType", "ClaimData.NonKatastrofeType", JenisTeks))
	kiri = tampilBila(kiri, p.v("ClaimData.StsKatastrofe") == "Catastrophe" ||
		p.v("ClaimData.NonKatastrofeType") == "Big Claim",
		p.medan("Catastrophe Note", "ClaimData.KatastrofeNote", JenisTeksPanjang))

	// Kolom kanan.
	kanan := []Medan{
		{Label: "PIC Name", Nilai: k.PembuatNama, Jenis: JenisTeks},
		{Label: "Date", Nilai: FormatWaktu(k.TglCreate), Jenis: JenisTanggal},
	}
	kanan = tampilBila(kanan, strings.TrimSpace(p.v("ClaimData.NoClaim")) != "",
		Medan{Label: "NoClaim", Nilai: p.v("ClaimData.NoClaim") + " / " + k.KlaimID, Jenis: JenisTeks})
	kanan = append(kanan,
		p.medan("PLA / DLA Number", "ClaimData.NoPla", JenisTeks),
		p.medan("Date of Loss", "ClaimData.DateOfLoss", JenisTanggal),
		p.medan("Report Date", "ClaimData.ReportDate", JenisTanggal),
		p.medan("Received Date", "ClaimData.DateReceived", JenisTanggal),
		p.medan("Reporter Name", "ClaimData.ReporterName", JenisTeks),
		p.medan("Reporter Phone Number", "ClaimData.ReporterTelp", JenisTeks),
		Medan{Label: "Reporter Status", Nilai: labelKode(labelKodeKlaim["ClaimData.ReporterStatus"],
			p.v("ClaimData.ReporterStatus")), Jenis: JenisTeks},
		p.medan("Specify...", "ClaimData.InsuredRelationshipOthers", JenisTeks),
		p.medan("Reporter Address", "ClaimData.ReportAddress", JenisTeksPanjang),
		p.medan("Report Description", "ClaimData.ReportDescription", JenisTeksPanjang),
		p.a("RNM Share (%)", "PersenRNM", JenisAngka),
	)
	ly.Bagian = append(ly.Bagian,
		Bagian{Kunci: "klaim", Judul: JudulClaimTreaty, Medan: kiri},
		Bagian{Kunci: "klaimKanan", Medan: kanan},
		Bagian{Kunci: "kerugian", Grid: []Grid{
			{Judul: "Insured Interests 100 %", Kolom: []KolomGrid{kol("Insured Interest", "ObjectName", JenisTeks),
				kol("Currency", "CurrencyID", JenisTeks), kol("Value In IDR", "KursObjectItem", JenisAngka),
				kol("Value", "TSIPerObject", JenisAngka)}, Baris: p.daftar("ClaimData.InterestList")},
			{Judul: "Count Claim Amount", Kolom: []KolomGrid{kol("Currency", "CurrencyID", JenisTeks),
				kol("Claim Amount", "Value", JenisAngka), kol("Claim Amount in IDR", "USD", JenisAngka)},
				Baris: p.daftar("ClaimData.ListClaimAmount")},
			{Judul: "Loss Allocation", Kolom: []KolomGrid{kol("Curr", "CurrencyID", JenisTeks),
				kol("Treaty Type", "TreatyType", JenisTeks), kol("Share(%)", "SharePercentage", JenisAngka),
				kol("Result Claim", "ClaimSpreaded", JenisAngka), kol("Result Claim In IDR", "ClaimEstimation", JenisAngka)},
				Baris: p.daftar("ClaimData.SpreadingRisk")},
		}},
		Bagian{Kunci: "estimasi", Grid: []Grid{
			{Judul: "Estimation List", Kolom: []KolomGrid{kol("", "TypeLoss", JenisTeks),
				kol("Estimation Date", "EstimationDate", JenisTanggal), kol("Type", "Type", JenisTeks),
				kol("Currency", "CurrencyID", JenisTeks), kol("Value In IDR", "KursValue", JenisAngka),
				kol("Gross Estimate Treaty (100%)", "GrossEstimationPct", JenisAngka),
				kol("Estimation RNM", "EstimationValue", JenisAngka),
				kol("Estimation RNM in IDR", "ConvertValue", JenisAngka)}, Baris: p.daftar("ClaimData.EstimationList")},
			// ContainerVisibleWhen .TransferType = 2.
			{Judul: "Total Original Currency Estimation", Kolom: []KolomGrid{kol("Currency", "Currency", JenisTeks),
				kol("Gross Estimate Treaty (100%)", "IDR", JenisAngka), kol("Estimation RNM", "Value", JenisAngka)},
				Baris: p.daftar("ClaimData.ListTotalEstimation")},
		}},
		Bagian{Kunci: "totalEstimasi", Judul: "Total Estimation In IDR", Medan: []Medan{
			p.medan("Total Gross Estimate(100%) in IDR", "ClaimData.TotalGrossEstimateIDR", JenisAngka),
			p.medan("Total Estimation in IDR", "ClaimData.TotalEstimasiIDR", JenisAngka),
		}},
		Bagian{Kunci: "riwayatAdjustment", Grid: []Grid{
			{Judul: "History Adjustment", Kolom: []KolomGrid{kol("Payment Type", "Type", JenisTeks),
				kol("Treaty Type", "TreatyName", JenisTeks), kol("Currency", "Currency", JenisTeks),
				kol("Adjustment Gross (100%)", "GrossAdjustment", JenisAngka),
				kol("Adjustment RNM", "AdjustmentValue", JenisAngka), kol("Status", "AcceptanceStatus", JenisTeks),
				kol("Komite No", "KomiteNo", JenisTeks), kol("Accepted Date", "AcceptedDate", JenisTanggal),
				kol("Accepted No", "AcceptedNo", JenisTeks)}, Baris: riwayatAdjustment(p.daftar(DaftarAdjustment))},
			{Judul: "Total Adjustment", Kolom: []KolomGrid{kol("Currency", "currency", JenisTeks),
				kol("Total Gross (100%)", "adjustmentGross", JenisAngka),
				kol("Total Adjustment RNM", "adjustmentValue", JenisAngka)}, Baris: barisTotal(total)},
		}},
		Bagian{Kunci: "deductible", Sel: selDeductible(p)},
		Bagian{Kunci: "spreading", Grid: []Grid{
			{Judul: "Spreading In", Kolom: kolomSpreading(), Baris: p.daftar(jalurAdj(kl.Adjustment, "SpreadingAdjustment"))},
			{Judul: "Spreading Out", Kolom: kolomSpreading(), Baris: p.daftar(jalurAdj(kl.Adjustment, "SpreadingQuotaShare"))},
		}},
	)
	bayar := []Medan{p.a("Payable To", "Payable", JenisTeks), p.a("Specify", "PayableTo", JenisTeks)}
	bank := []Medan{p.a("Name of Bank", "NameOfBank", JenisTeks)}
	bank = tampilBila(bank, p.adj["SwiftCode"] != "", p.a("Swift Code", "SwiftCode", JenisTeks))
	bank = append(bank, p.a("Branch of Bank", "BranchOfBank", JenisTeks), p.a("Account No", "NoAccount", JenisTeks))
	// Teks komite: empat pertama `pyVisible NOTBLANK`, Remarks selalu.
	var teks []Medan
	for _, x := range []Medan{
		p.a("Circumstances", "DataCommitteeTreaty.CircumCauseOfLoss", JenisTeksPanjang),
		p.medan("Occupation", "ClaimData.Occupation", JenisTeksPanjang),
		p.a("Salvage", "DataCommitteeTreaty.Salvage", JenisTeksPanjang),
		p.a("Adjuster / Consultant Fee", "DataCommitteeTreaty.AdjusterFee", JenisTeksPanjang),
	} {
		teks = tampilBila(teks, strings.TrimSpace(x.Nilai) != "", x)
	}
	teks = append(teks, p.a("Remarks", "DataCommitteeTreaty.Remarks", JenisTeksPanjang))
	ly.Bagian = append(ly.Bagian,
		Bagian{Kunci: "bayar", Medan: bayar},
		Bagian{Kunci: "bank", Medan: bank},
		Bagian{Kunci: "teksKomite", Medan: teks},
		Bagian{Kunci: "tangga", Grid: []Grid{{Judul: "Committe Accept Status", Kolom: []KolomGrid{
			kol("Committe Name", "jabatan", JenisTeks), kol("Status", "keputusan", JenisTeks),
			kol("Date Approve", "tanggal", JenisTanggalJam), kol("Comment", "komentar", JenisTeks)},
			Baris: barisTangga(k.Tangga)}}},
	)

	ly.Isian = IsianLayar{Nilai: nilaiAwal(k), Terbuka: IsianTerbuka(k), PilihanTerima: LabelTerima,
		Label: map[string]string{"acceptStatus": LabelAcceptStatus, "isSubjectivity": LabelSubjectivity,
			"subjectivityNote": LabelSubjectivityNote, "usulTutup": LabelProposeClose, "usulCadang": LabelProposeReserved,
			"comment": LabelNote}}

	bernomor := strings.TrimSpace(p.adj["AcceptedNo"]) != ""
	kirim := Tombol{Label: TombolKirim, Aksi: "putuskan", Aktif: ly.BolehKerja && !bernomor}
	if bernomor {
		kirim.Alasan = AlasanSudahBernomor
	}
	ly.Tombol = []Tombol{
		{Label: TombolLihat, Aksi: "lihat", Aktif: false, Alasan: AlasanLihatNonaktif},
		{Label: TombolBatal, Aksi: "batal", Aktif: true},
		kirim,
	}
	return ly
}

// nilaiAwal - isian yang tersimpan di `pyWorkPage` dari tingkat sebelumnya (Pega tidak mengosongkannya antar tingkat):
// `.AcceptStatus` = keputusan terakhir, `.Comment` = komentar tingkat sebelumnya (tingkat 1: `AddKomiteTreatyChild_ACT`
// S16, hanya kirim ulang subjectivity), Subjectivity / catatannya / dua Propose = nilai tersimpan di header.
func nilaiAwal(k Kasus) Keputusan {
	n := Keputusan{AcceptStatus: k.AcceptStatus, UsulTutup: k.UsulTutup == UsulYa, UsulCadang: k.UsulCadang == UsulYa,
		IsSubjectivity: k.Subjectivity == UsulYa, SubjectivityNote: k.SubjectivityNote}
	if i := k.barisBerjalan(); i > 0 {
		n.Comment = k.Tangga[i-1].Komentar
	} else if i == 0 {
		n.Comment = k.KomentarAwal // S16 kirim ulang subjectivity; selainnya kosong
	}
	return n
}

func kolomSpreading() []KolomGrid {
	return []KolomGrid{kol("Currency", "Currency", JenisTeks), kol("Treaty Type", "TreatyName", JenisTeks),
		kol("Share(%)", "SharePercentage", JenisAngka), kol("Claim Spreaded", "ClaimSpreaded", JenisAngka)}
}

// riwayatAdjustment - `TempClaimData.AdjustmentList` (SetKomiteList_Act S4 Page-Copy `ClaimData`): kolom "Komite No" =
// `.KomiteNo` (SetKomiteNo_Act: ID kasus komite) - di sistem baru ID kasus komite baris itu.
func riwayatAdjustment(rows []map[string]string) []map[string]string {
	for _, b := range rows {
		if strings.TrimSpace(b["KomiteNo"]) == "" {
			b["KomiteNo"] = b["KomiteID"]
		}
	}
	return rows
}

func barisTotal(t []TotalMataUang) []map[string]string {
	out := make([]map[string]string, 0, len(t))
	for _, x := range t {
		out = append(out, map[string]string{"currency": x.Currency, "adjustmentGross": x.AdjustmentGross,
			"adjustmentValue": x.AdjustmentValue})
	}
	return out
}

func barisTangga(t []Anggota) []map[string]string {
	out := make([]map[string]string, 0, len(t))
	for _, a := range t {
		out = append(out, map[string]string{"jabatan": a.Jabatan, "keputusan": a.Keputusan, "tanggal": a.Tanggal,
			"komentar": a.Komentar})
	}
	return out
}

// selDeductible - tabel 6 x 4 (`ContainerVisibleWhen .TransferType = 2`).
func selDeductible(p pembaca) [][]Medan {
	h := func(s string) Medan { return Medan{Label: s, Jenis: JenisTeks} }
	cur := p.a("", "Currency", JenisTeks)
	return [][]Medan{
		{h("Dedutible Type"), h("Deductible (%)"), h("Currency"), h("Treaty Gross (100%)"), h("Currency"),
			h("RNM Share Claim")},
		{{}, {}, cur, p.a("", "GrossAdjustment", JenisAngka), cur, p.a("", "GrossValue", JenisAngka)},
		{{Nilai: labelKode(labelKodeKlaim["IndividualRiskType"], p.adj["IndividualRiskType"]), Jenis: JenisTeks},
			p.a("", "IndividualRiskPercentage", JenisAngka), cur, p.a("", "IndividualRiskValue", JenisAngka), cur,
			p.a("", "IndividualRiskRNM", JenisAngka)},
		{h("Total"), {}, cur, p.a("", "ProposeAdjustmentValue", JenisAngka), cur, p.a("", "AdjustmentValue", JenisAngka)},
	}
}
