package models

// Untuk apa berkas ini: LAYAR KOMITE - flow action `ViewTransferDtl`, Section `ShowTransfer` (korpus `Komite Claim Non
// Prop`). Urutan bagian, label, dan syarat tampil diambil dari Section itu (VERBATIM, termasuk ejaan "Committe").
// Nilai dibaca dari kasus klaim induk lewat kontrak (`pyWorkCover`, dan `pyWorkPage.Adjustment` / `pyWorkPage.Komite` -
// salinan baris akseptasi yang dibuat `CreateChildKomiteCNP_Act`; baris itu beku selama diserahkan ke komite, jadi baris
// aslinya yang dibaca). Pola disalin dari Komite Claim Prop (bukan impor).
//
// Tidak dibangun (PARITAS):
//   - judul "CLOSE" / "REJECT" dan seluruh cabang `IsCloseFile` / `IsReject` kasus komite (jalur CWP, OQ-CNP-36) - syarat
//     tampil `IsCloseFile!=1 && IsReject!=1` selalu benar;
//   - sel ber-`pyVisible` NEVER ("sample text", ExtentOfLoss, LegalLiability, tombol "spacer");
//   - tombol Cancel / Submit Section ada di kontainer NEVER (`pyShowFAButtons` false) - dibangun seperti Komite Claim
//     Prop (tombol layar), tanpa syarat nonaktif tambahan.

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
	// JenisCentang - pxCheckbox hanya-baca ("true" / "1" = dicentang).
	JenisCentang = "centang"
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
	Terbuka       bool      `json:"terbuka"`
	PilihanTerima []Pilihan `json:"pilihanTerima"`
	// PilihanSubjectivityNote - dropdown "Subjectivity Note" (SubjectivityNote.xml).
	PilihanSubjectivityNote []Pilihan         `json:"pilihanSubjectivityNote"`
	Label                   map[string]string `json:"label"`
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
	TombolLihat      = "View Claim"
	TombolBatal      = "Cancel"
	TombolKirim      = "Submit"
	JudulClaimTreaty = "Claim Analysis"
)

// LabelTerima - pilihan `.AcceptStatus` (`pyListSource associated`, prompt values tidak diekspor). Teks dari
// `SetDataAcceptationTreaty_Act` S2-S3 (`ParamDataTreaty.CARI7`): 1 "Approve", 2 "Reject".
var LabelTerima = []Pilihan{{Nilai: KeputusanSetuju, Label: "Approve"}, {Nilai: KeputusanTolak, Label: "Reject"}}

// "View Claim" = harness `ViewDtlClaimKmt` atas `pyWorkCover` (klaim induk; harness tidak diekspor). Padanannya berkas
// Claim Non Prop klaim induk, dibuka hanya-baca di jendela di atas layar komite (`PropsRute.onLihatBerkas`).

// labelTipeAdjustment - `.Adjustment.Type` (dropdown `associated`): label Claim Non Prop `labelKode["AdjustmentType"]`.
var labelTipeAdjustment = map[string]string{"1": "Claim", "2": "Adjuster Fee", "3": "Salvage", "4": "Consultant Fee"}

// labelKodeKlaim - dropdown `associated` panel klaim induk (label Claim Non Prop `labelKode`). Kode lain tampil apa
// adanya.
var labelKodeKlaim = map[string]map[string]string{
	"ClaimData.ReportType":     {"1": "Direct", "2": "Via Email", "3": "Via Fax", "4": "via Postal Mail/Courier", "5": "Via Telephone"},
	"ClaimData.ReporterStatus": {"1": "Ceding Co Name", "2": "SOB Name", "3": "Others"},
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

// SusunLayar menyusun layar `ShowTransfer` (korpus Komite Claim Non Prop) kasus `k` atas klaim induk `kl`; `akun` =
// pelaku, `peran` = workbasket aktifnya. Urutan bagian = urutan Section (S2-S46).
func SusunLayar(k Kasus, kl kontrak.KlaimTreaty, akun string, peran []string) Layar {
	p := pembaca{kl: kl, adj: AdjustmentKlaim(kl)}
	ly := Layar{Kasus: k, Judul: []string{JudulKomite}, BolehKerja: k.Pemegang(akun, peran)}
	n := kl.Adjustment

	// S5-S6 kolom kiri.
	kiri := []Medan{
		p.medan("Class of Business", "OfferFacIn.QuotationData.BusinessName", JenisTeks),
		p.medan("Ceding", "TreatyInMaster.Ceding", JenisTeks),
		p.medan("Source of Business", "TreatyInMaster.LeadingReinsSource", JenisTeks),
		p.medan("Policy No", "ClaimData.PolicyData.PolicyNo", JenisTeks),
		p.medan("Policy No Ceding", "ClaimData.PolicyNo", JenisTeks),
		p.medan("Insured Name", "ClaimData.InsuredName", JenisTeks),
		// S7 Inline: StartDateTime "-" EndDateTime, sel tanpa label.
		{Label: "", Nilai: periode(p.v("ClaimData.PolicyData.StartDateTime"), p.v("ClaimData.PolicyData.EndDateTime")),
			Jenis: JenisTeks},
		p.medan("Cause of Loss", "ClaimData.CauseOfLoss", JenisTeks),
	}
	kiri = tampilBila(kiri, strings.TrimSpace(p.v("ClaimData.Occupation")) != "",
		p.medan("Occupation", "ClaimData.Occupation", JenisTeksPanjang))
	kiri = append(kiri, Medan{Label: "Report Type", Nilai: labelKode(labelKodeKlaim["ClaimData.ReportType"],
		p.v("ClaimData.ReportType")), Jenis: JenisTeks})
	kiri = tampilBila(kiri, p.v("ClaimData.AppointedADJ") != "",
		p.medan("Adjuster / Professional", "ClaimData.AppointedADJ", JenisTeks))
	kiri = tampilBila(kiri, p.v("ClaimData.ConsultantName") != "",
		p.medan("Consultant", "ClaimData.ConsultantName", JenisTeks))
	kiri = append(kiri, p.medan("Location of Loss", "ClaimData.Location", JenisTeksPanjang),
		Medan{Label: "Type", Nilai: labelKode(labelTipeAdjustment, p.adj["Type"]), Jenis: JenisTeks})

	// S8 kolom kanan.
	kanan := []Medan{
		{Label: "PIC Name", Nilai: k.PembuatNama, Jenis: JenisTeks},
		{Label: "Date", Nilai: FormatWaktu(k.TglCreate), Jenis: JenisTanggal},
	}
	kanan = tampilBila(kanan, strings.TrimSpace(p.v("ClaimData.NoClaim")) != "",
		Medan{Label: "", Nilai: p.v("ClaimData.NoClaim") + " / " + k.KlaimID, Jenis: JenisTeks})
	persen := p.adj["PersenRNM"]
	if strings.TrimSpace(persen) == "" { // `.Adjustment.PersenRNM` tanpa penulis di Claim Non Prop: RNM Share master
		persen = p.v("TreatyInMaster.RNMShare")
	}
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
		Medan{Label: "RNM Share (%)", Nilai: persen, Jenis: JenisAngka},
	)
	ly.Bagian = append(ly.Bagian,
		Bagian{Kunci: "klaim", Judul: JudulClaimTreaty, Medan: kiri},
		Bagian{Kunci: "klaimKanan", Medan: kanan},
		// S10 (IsCloseFile kasus komite != 1 - selalu, CWP tidak dibangun).
		Bagian{Kunci: "kejadian", Medan: []Medan{p.a("Circumstances", "DataCommitteeTreaty.CircumCauseOfLoss",
			JenisTeksPanjang)}},
		// S11 Claim Acceptation, lalu kaki Total Claim Amount (tanpa penulis di korpus - kosong, OQ-CNP-37).
		Bagian{Kunci: "klaimAkseptasi", Grid: []Grid{{Judul: "", Kolom: []KolomGrid{
			kol("Currency", "Currency", JenisTeks), kol("Rate of Exchange", "AltValue", JenisAngka),
			kol("Claim Amount", "Value", JenisAngka), kol("TPL", "TPL", JenisAngka),
			kol("Adjuster Fee", "AdjusterFee", JenisAngka), kol("Salvage", "Salvage", JenisAngka),
			kol("Fee", "CNPOthersFee", JenisAngka), kol("Proportion (%)", "PctProrateClaim", JenisAngka),
			kol("Claim Amount in IDR", "USD", JenisAngka), kol("Claim Amount Cedant", "ClaimAmountCedant", JenisAngka)},
			Baris: barisNama(p.daftar(jalurAdj(n, AnakClaimAcc)))}}},
		Bagian{Kunci: "totalKlaim", Medan: []Medan{
			{Label: "Total Claim Amount", Nilai: p.v("ClaimData.TotalListClaimAmount"), Jenis: JenisAngka},
			{Label: "", Nilai: p.v("ClaimData.TotalListClaimAmountIDR"), Jenis: JenisAngka}}},
		Bagian{Kunci: "alokasi", Grid: []Grid{
			{Judul: "Loss Allocation", Kolom: []KolomGrid{kol("Currency", "Currency", JenisTeks),
				kol("Treaty Name", "TreatyName", JenisTeks), kol("Share (%)", "ClaimPercentage", JenisAngka),
				kol("Claim Amount", "ClaimAmountAdjust", JenisAngka), kol("Adjuster Fee", "AdjusterFee", JenisAngka),
				kol("Salvage", "Salvage", JenisAngka), kol("Fee", "CNPOthersFee", JenisAngka),
				kol("To XOL", "CNPFlagXOL", JenisCentang)}, Baris: barisNama(p.daftar(jalurAdj(n, AnakLossAlloc)))},
			{Judul: "XOL Allocation", Kolom: kolomXOL(), Baris: barisNama(p.daftar(jalurAdj(n, AnakXOL)))},
		}},
	)
	if dibayar := p.daftar(jalurAdj(n, AnakXOLDibayar)); len(dibayar) > 0 { // S23 IsPrevious == 1
		ly.Bagian = append(ly.Bagian, Bagian{Kunci: "dibayar", Grid: []Grid{{Judul: "Previously Calculated",
			Kolom: kolomXOL(), Baris: barisNama(dibayar)}}})
	}
	ly.Bagian = append(ly.Bagian, Bagian{Kunci: "spreading", Grid: []Grid{ // S27
		{Judul: "Spreading In", Kolom: kolomSpreading(), Baris: p.daftar(jalurAdj(n, AnakSpreadIn))},
		{Judul: "Spreading Out", Kolom: kolomSpreading(), Baris: p.daftar(jalurAdj(n, AnakSpreadOut))},
	}})
	if len(p.kl.Daftar[jalurAdj(n, AnakXOLDibayar)]) > 0 { // S36 tampil[IsPrevious == 1 ...] - VERBATIM
		bayar := []Medan{{Label: "Payable To", Nilai: labelKode(LabelPayable, p.v("ClaimData.Payable")), Jenis: JenisTeks},
			p.medan("Specify", "ClaimData.PayableTo", JenisTeks)}
		ly.Bagian = append(ly.Bagian, Bagian{Kunci: "bayar", Medan: bayar},
			Bagian{Kunci: "bank", Medan: bankAkseptasi(p, "")})
		if p.adj["FlagCurrency"] == "1" { // S40
			ly.Bagian = append(ly.Bagian, Bagian{Kunci: "bank2", Medan: bankAkseptasi(p, "2")})
		}
	}
	ly.Bagian = append(ly.Bagian,
		Bagian{Kunci: "tangga", Grid: []Grid{{Judul: "Committe Accept Status", Kolom: []KolomGrid{ // S41
			kol("Committe Name", "jabatan", JenisTeks), kol("", "inisial", JenisTeks), kol("Status", "keputusan", JenisTeks),
			kol("Date", "tanggal", JenisTanggalJam), kol("Comment", "komentar", JenisTeks)},
			Baris: barisTangga(k.Tangga)}}},
		Bagian{Kunci: "teksKomite", Medan: []Medan{p.a("Remarks", "DataCommitteeTreaty.Remarks", JenisTeksPanjang)}}, // S45
	)

	ly.Isian = IsianLayar{Nilai: nilaiAwal(k), Terbuka: IsianTerbuka(k), PilihanTerima: LabelTerima,
		PilihanSubjectivityNote: PilihanSubjectivityNote,
		Label: map[string]string{"acceptStatus": LabelAcceptStatus, "isSubjectivity": LabelSubjectivity,
			"subjectivityNote": LabelSubjectivityNote, "usulTutup": LabelProposeClose, "usulCadang": LabelProposeReserved,
			"comment": LabelNote}}
	ly.Tombol = []Tombol{
		{Label: TombolLihat, Aksi: "lihat", Aktif: true},
		{Label: TombolBatal, Aksi: "batal", Aktif: true},
		{Label: TombolKirim, Aksi: "putuskan", Aktif: ly.BolehKerja},
	}
	return ly
}

// periode - Section S7 Inline: "<mulai> - <akhir>".
func periode(mulai, akhir string) string {
	if strings.TrimSpace(mulai) == "" && strings.TrimSpace(akhir) == "" {
		return ""
	}
	return TanggalSurat(mulai) + " - " + TanggalSurat(akhir)
}

// bankAkseptasi - S38 / S40: Name of Bank, Currency, Swift Code (NOTBLANK), Branch of Bank, Account No (akhiran "" / "2").
func bankAkseptasi(p pembaca, sfx string) []Medan {
	bank := []Medan{p.a("Name of Bank", "NameOfBank"+sfx, JenisTeks), p.a("Currency", "Currency"+sfx, JenisTeks)}
	bank = tampilBila(bank, p.adj["SwiftCode"+sfx] != "", p.a("Swift Code", "SwiftCode"+sfx, JenisTeks))
	return append(bank, p.a("Branch of Bank", "BranchOfBank"+sfx, JenisTeks), p.a("Account No", "NoAccount"+sfx, JenisTeks))
}

// nilaiAwal - isian yang tersimpan di `pyWorkPage` dari tingkat sebelumnya (Pega tidak mengosongkannya antar tingkat):
// `.AcceptStatus` = keputusan terakhir, `.Comment` = komentar tingkat sebelumnya, Subjectivity / catatannya / dua
// Propose = nilai tersimpan di header.
func nilaiAwal(k Kasus) Keputusan {
	n := Keputusan{AcceptStatus: k.AcceptStatus, UsulTutup: k.UsulTutup == UsulYa, UsulCadang: k.UsulCadang == UsulYa,
		IsSubjectivity: k.Subjectivity == UsulYa, SubjectivityNote: k.SubjectivityNote}
	if i := k.barisBerjalan(); i > 0 {
		n.Comment = k.Tangga[i-1].Komentar
	} else if i == 0 {
		n.Comment = k.KomentarAwal // CreateChildKomiteCNP_Act: putaran terdahulu; selainnya kosong
	}
	return n
}

// kolomXOL - S19 / S23 (XOL Allocation, Previously Calculated).
func kolomXOL() []KolomGrid {
	return []KolomGrid{kol("Currency", "Currency", JenisTeks), kol("Treaty Name", "TreatyName", JenisTeks),
		kol("Claim Amount", "TotalClaim", JenisAngka), kol("RNM Share (%)", "ClaimPercentage", JenisAngka),
		kol("Claim Amount RNM", "ClaimSpreaded", JenisAngka), kol("Adjuster Fee", "AdjusterFee", JenisAngka),
		kol("Salvage", "Salvage", JenisAngka), kol("Fee", "CNPOthersFee", JenisAngka),
		kol("Reinstatement Premium", "CNPReinstatementRNM", JenisAngka)}
}

// kolomSpreading - S28 / S32 (Spreading In / Out).
func kolomSpreading() []KolomGrid {
	return []KolomGrid{kol("Currency", "Currency", JenisTeks), kol("Treaty Type", "TreatyName", JenisTeks),
		kol("Share(%)", "SharePercentage", JenisAngka), kol("Claim Spreaded", "ClaimSpreaded", JenisAngka),
		kol("Adjuster Fee", "AdjusterFee", JenisAngka), kol("Salvage", "Salvage", JenisAngka),
		kol("Fee", "CNPOthersFee", JenisAngka), kol("Total Claim", "TotalClaim", JenisAngka),
		kol("Reinstatement Premium", "PremiumSpreaded", JenisAngka)}
}

// namaKode - kolom nama tampilan grid klaim dan kolom kodenya (nama kosong = kode yang tampil).
var namaKode = map[string]string{"Currency": "CurrencyID", "TreatyName": "TreatyType"}

// barisNama - nilai tampilan grid klaim: mata uang / treaty = NAMA-nya (kode hanya bila nama kosong). Baris `p.daftar`
// sudah salinan - diubah di tempat.
func barisNama(rows []map[string]string) []map[string]string {
	for _, b := range rows {
		for nama, kode := range namaKode {
			if strings.TrimSpace(b[nama]) == "" && b[kode] != "" {
				b[nama] = b[kode]
			}
		}
	}
	return rows
}

// barisTangga - grid "Committe Accept Status": jabatan roster (`KomitePost`), inisial (`Initial` - tertulis mati per
// nama orang di SetKomiteList_Act, dibuang OQ-CNP-03), keputusan, tanggal, komentar.
func barisTangga(t []Anggota) []map[string]string {
	out := make([]map[string]string, 0, len(t))
	for _, a := range t {
		out = append(out, map[string]string{"jabatan": a.Jabatan, "inisial": "",
			"keputusan": labelKode(LabelKeputusanAnggota, a.Keputusan), "tanggal": a.Tanggal, "komentar": a.Komentar})
	}
	return out
}
