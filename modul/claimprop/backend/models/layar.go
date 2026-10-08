package models

// Untuk apa berkas ini: DEFINISI LAYAR - Section `OutstandingClaim` (FlowAction OutstandingClaim, Assignment2) dan
// `InputAcceptation` (FlowAction InputAcceptation, Assignment1) beserta section yang di-include-nya
// (`Catastrope_Sec`, `OutstandingClaim_Intrs`, `OutstandingClaim_Est`, `OutstandingClaim_Sprd`, `InputAcceptation_Est`,
// `InputAcceptation_Adjs`). Label, urutan, dan kondisi VERBATIM dari XML (pindai `cp/layar.md` bab 4.1-4.2).
//
// Aturan baca:
//   - `Aksi` = nama activity tanpa akhiran `_Act` / `_act` (lihat `services/aksi.go`); "" = postValue saja (nilai
//     disimpan, turunan dihitung ulang).
//   - Tombol yang activity-nya TIDAK diekspor (`BackToRegister_act`, `UpdateEstimasi_Act`, `AddSpreading_Act`,
//     `DeleteSpreading_Act`) atau harness-nya tidak diekspor (`DetailPolisRealization`, `InputTreatyInOffer`) tampil
//     sesuai section tetapi NONAKTIF dengan `Catatan` OQ.
//   - Sel ber-visible `NEVER` / `1=2` dan placeholder `.pyTemplate*` tidak dibangun.
//   - Medan tanggal Pega disimpan sebagai tanggal; kendali tanggal-waktu untuk DateTimeFormat DateTime-Short.

import "strings"

// ---------------------------------------------------------------- kondisi

func sama(j, v string) Kondisi   { return func(h *Halaman) bool { return h.Ambil(j) == v } }
func beda(j, v string) Kondisi   { return func(h *Halaman) bool { return h.Ambil(j) != v } }
func terisi(j string) Kondisi    { return func(h *Halaman) bool { return h.Ambil(j) != "" } }
func atau(ks ...Kondisi) Kondisi { return func(h *Halaman) bool { return anyK(h, ks) } }

func anyK(h *Halaman, ks []Kondisi) bool {
	for _, k := range ks {
		if k(h) {
			return true
		}
	}
	return false
}

func dan(ks ...Kondisi) Kondisi {
	return func(h *Halaman) bool {
		for _, k := range ks {
			if !k(h) {
				return false
			}
		}
		return true
	}
}

var (
	isOutstanding  = sama("IsOutstanding", "1")
	isAcceptation  = sama("IsAcceptation", "1")
	isAnyAccept    = sama("IsAnyAcceptation", "1")
	isEstimation   = sama("isEstimation", "1")
	nonCatastrophe = sama(CD+"StsKatastrofe", "Non-Catastrophe")
	editCatastrope = sama(CD+"EditCatastrope", "true")
	selalu         = func(*Halaman) bool { return true }
	pyNoteKosong   = sama("pyNote", "")
)

func bSama(p, v string) KondisiBaris { return func(_ *Halaman, b Baris) bool { return b[p] == v } }

// bTerisi - medan baris terisi.
func bTerisi(p string) KondisiBaris { return func(_ *Halaman, b Baris) bool { return b[p] != "" } }

// bAtau - salah satu kondisi baris benar.
func bAtau(ks ...KondisiBaris) KondisiBaris {
	return func(h *Halaman, b Baris) bool {
		for _, k := range ks {
			if k(h, b) {
				return true
			}
		}
		return false
	}
}

var (
	bIsAdjVal   = bSama("IsAdjVal", "Yes")
	bNote       = bSama("Note", "Yes")
	bOldData    = bSama("IsOldData", "Yes")
	bPrintFace  = bSama("PrintFaceClaim", "1")
	bSelalu     = func(*Halaman, Baris) bool { return true }
	hSelalu     = func(*Halaman, Baris) bool { return true }
	bOutstandng = func(h *Halaman, _ Baris) bool { return h.Ambil("IsOutstanding") == "1" }
)

// ---------------------------------------------------------------- pembentuk

func medan(jalur, label, kendali string) Unsur {
	return Unsur{Jenis: JenisMedan, Jalur: jalur, Label: label, Kendali: kendali}
}

func ro(u Unsur) Unsur                { u.HanyaBaca = selalu; return u }
func roJika(u Unsur, k Kondisi) Unsur { u.HanyaBaca = k; return u }
func naJika(u Unsur, k Kondisi) Unsur { u.Nonaktif = k; return u }
func tampil(u Unsur, k Kondisi) Unsur { u.Tampil = k; return u }
func wajibU(u Unsur) Unsur            { u.Wajib = selalu; return u }
func wajibJ(u Unsur, k Kondisi) Unsur { u.Wajib = k; return u }
func aksi(u Unsur, a string) Unsur    { u.Aksi = a; return u }
func sumber(u Unsur, s string) Unsur  { u.Sumber = s; return u }
func tampilIsi(u Unsur) Unsur         { u.Tampil = terisi(u.Jalur); return u }
func label(t string) Unsur            { return Unsur{Jenis: JenisLabel, Label: t} }
func catatan(u Unsur, c string) Unsur { u.Catatan = c; return u }

// tampilan - jalur teks yang ditampilkan medan ber-sumber (nilai Jalur tetap yang disimpan).
func tampilan(u Unsur, j string) Unsur { u.Tampilan = j; return u }

func bagian(judul string, anak ...Unsur) Unsur {
	return Unsur{Jenis: JenisBagian, Label: judul, Anak: anak}
}

// Letak layout Pega (`pyLayoutOtherFormat`, layout group) - dirender layar; "" = Stacked with labels left.
const (
	LetakDua     = "dua"     // Inline grid double: setiap anak satu sel, dua sel per baris
	LetakSebaris = "sebaris" // Inline / Inline labels left: anak sebaris; label bagian = label baris
	LetakTab     = "tab"     // layout group Tab: anak = bagian berjudul (satu tab per bagian)
	LetakJudul   = "judul"   // kepala layar: Inline grid triple dengan label di sel tengah
)

// Ikon tombol (kelas ikon / gambar tombol di XML); label tetap dikirim sebagai keterangan.
const (
	IkonTambah = "tambah" // pi-plus, webwb/pyWorkActionsAddWork.png
	IkonHapus  = "hapus"  // pi-trash
	IkonUbah   = "ubah"   // pi-pencil, webwb/pyEditIcon.png
	IkonSimpan = "simpan" // pi-check
)

func letak(l string, anak ...Unsur) Unsur { return Unsur{Jenis: JenisBagian, Letak: l, Anak: anak} }

// sebaris - layout Inline; `lbl` = label baris (boleh kosong).
func sebaris(lbl string, anak ...Unsur) Unsur {
	u := letak(LetakSebaris, anak...)
	u.Label = lbl
	return u
}

// dua - layout Inline grid double; bagian tanpa judul di dalamnya = satu kolom Stacked with labels left.
func dua(sel ...Unsur) Unsur { return letak(LetakDua, sel...) }

func ikon(u Unsur, i string) Unsur { u.Ikon = i; return u }

// ModeLayar - penanda MODE layar Pega yang hidup di halaman kerja tetapi tidak punya kolom di tabel datar: ikon Edit /
// Save Catastrophe (`EditCatastrope`, Catastrope_Sec) dan ikon Edit RNM Share (`IsEditRNMShare`). Server mengirimnya
// di `Layar.Mode`, layar mengembalikannya di setiap aksi, server memasangnya SEBELUM tata (medan / aksi terbuka)
// dihitung. Aman: penanda hanya membuka mode edit yang memang dapat dinyalakan pemegang lewat tombolnya (temuan work
// owner 08-10-2026 "Catastrophe tidak berfungsi" - penanda hilang saat halaman dibaca ulang).
var ModeLayar = []string{CD + "EditCatastrope", "IsEditRNMShare"}

// PasangMode memasang penanda mode kiriman layar; hanya kunci ModeLayar dan nilai "true" / "false".
func PasangMode(h *Halaman, mode map[string]string) {
	for _, k := range ModeLayar {
		if v, ada := mode[k]; ada && (v == "true" || v == "false") {
			h.Setel(k, v)
		}
	}
}

// AmbilMode - penanda mode halaman untuk `Layar.Mode`.
func AmbilMode(h *Halaman) map[string]string {
	out := map[string]string{}
	for _, k := range ModeLayar {
		if v := h.Ambil(k); v != "" {
			out[k] = v
		}
	}
	return out
}

// Tombol "+" Consultant / Adjuster (Pega: harness MstAdjusterConsultant, tambah master): layar membuka popup tambah
// lalu menyimpan lewat rute pinjaman modul Adjuster Consultant (`POST /api/adjuster-consultant`, keputusan work owner
// 08-10-2026) dan mengisi ID baru ke medannya. Tidak ada aksi server Claim Prop; nonaktif selama ID hanya-baca.
const (
	AksiTambahKonsultan = "TambahKonsultan"
	AksiTambahAdjuster  = "TambahAdjuster"
)

// AksiLihatPolis - tombol View (harness DetailPolisRealization tidak diekspor): layar membuka berkas polis NB / EDM
// Treaty In di tab baru lewat `GET /berkas-polis` (keputusan work owner 08-10-2026); tidak ada aksi server.
const AksiLihatPolis = "LihatPolis"

func tombol(id, lbl, aksiNama string) Unsur {
	return Unsur{Jenis: JenisTombol, ID: id, Label: lbl, Aksi: aksiNama}
}

// tombolOQ - tombol yang tampil sesuai section tetapi nonaktif (rule tidak diekspor).
func tombolOQ(id, lbl, oq string) Unsur {
	return Unsur{Jenis: JenisTombol, ID: id, Label: lbl, Nonaktif: selalu, Catatan: oq}
}

func kol(prop, judul, kendali string) Unsur {
	return Unsur{Jenis: JenisMedan, Jalur: prop, Label: judul, Kendali: kendali}
}

func kRO(u Unsur) Unsur                  { u.HanyaBacaB = hSelalu; return u }
func kROJ(u Unsur, k KondisiBaris) Unsur { u.HanyaBacaB = k; return u }
func kNA(u Unsur, k KondisiBaris) Unsur  { u.NonaktifB = k; return u }
func kAksi(u Unsur, a string) Unsur      { u.Aksi = a; return u }
func kSumber(u Unsur, s string) Unsur    { u.Sumber = s; return u }
func kTombol(id, lbl, a string, na KondisiBaris) Unsur {
	return Unsur{Jenis: JenisTombol, ID: id, Label: lbl, Aksi: a, NonaktifB: na}
}

func ikonK(u Unsur, i string) Unsur { u.Ikon = i; return u }

// Catatan OQ untuk rule yang tidak diekspor.
const (
	OQTidakDiekspor   = "OQ-CP-01: activity tidak ada di ekspor Pega - perilakunya tidak ditebak"
	OQHarnessHilang   = "OQ-CP-02: harness tidak ada di ekspor Pega"
	OQLayananLuar     = "OQ-CP-03: layanan REST luar (M_LINK_SERVICE) belum disetujui dipanggil dari aplikasi"
	OQMasterLain      = "OQ-CP-04: pemeliharaan master milik menu lain - Claim Prop hanya memakai pemilihnya"
	OQDokumenPDF      = "OQ-CP-05: stream HTML dokumen tidak diekspor dan aplikasi belum punya mesin PDF"
	OQTutupTanpaBayar = "OQ-CP-06: kasus komite tanpa baris adjustment tidak dapat ditulis (T_GENERAL_KOMITE.ADJUSTMENT_ID NOT NULL)"
)

// Kunci daftar pilihan (`services` mengisinya).
const (
	SumberMataUang   = "mataUang"    // BrowseCurrency_RD
	SumberLimits     = "limits"      // pageList TreatyInMaster.Limits (.TreatyType)
	SumberJenisReas  = "jenisReas"   // BrowseReinsuranceType_RD
	SumberJenisReas4 = "jenisReas4"  // BrowseReinsuranceType_RD Type = 4
	SumberSpreading  = "spreading"   // pageList Spreading.pxResults (InputAcceptation_Est) = spreading polis klaim
	SumberAdjuster   = "adjuster"    // BrowseAdjusterConsultant
	SumberProvinsi   = "provinsi"    // BrowseProvince_RD Nation INDONESIA
	SumberShareRNM   = "shareRNM"    // pageList TreatyShare.pxResults (GetRNMShareTreaty)
	SumberMUAdj      = "mataUangAdj" // pageList .CurencyAdjustment
	SumberAllocation = "allocation"  // pageList .LossAllocation
	SumberRekening   = "rekening"    // pageList Result.pxResults (GetDataBankAccount*_sql)
	AwalanKode       = "kode:"       // prompt values `associated` - tidak terbaca di XML; kode DB apa adanya
)

// KodePilihan - kode yang teramati di data DEV untuk properti bersumber `associated` (prompt values tidak ada di
// korpus, OQ-CP-07). Ditampilkan apa adanya (keputusan work owner "jangan di singkat ikuti apa yang di DB").
var KodePilihan = map[string][]string{
	"ReportType":         {"1", "2", "3", "4", "5"},
	"ReporterStatus":     {"1", "2", "3"},
	"FormType":           {"1", "2"},
	"TypeDeductible":     {"1", "2"},
	"Payable":            {"1", "2", "3"},
	"StsKatastrofe":      {"Catastrophe", "Non-Catastrophe"},
	"NonKatastrofeType":  {"Claim", "Big Claim"},
	"EstimationType":     {"1", "2", "3", "4"},
	"AdjustmentType":     {"1", "2", "3", "4"},
	"IndividualRiskType": {"1", "2", "3"},
	"AcceptanceStatus":   {"1", "2"},
}

// LabelKode - label tampilan kode `associated` yang diberikan work owner (prompt values tidak ada di korpus; 08-10-2026).
// Nilai tersimpan tetap kodenya.
var LabelKode = map[string]map[string]string{
	"ReportType":     {"1": "Direct", "2": "Via Email", "3": "Via Fax", "4": "via Postal Mail/Courier", "5": "Via Telephone"},
	"ReporterStatus": {"1": "Ceding Co Name", "2": "SOB Name", "3": "Others"},
	// ASM-FW-GCNMFW-Data-Estimasi.Type - kolom Type grid Estimation List (screenshot work owner 08-10-2026)
	"EstimationType": {"1": "Claim", "2": "Adjuster Fee", "3": "Salvage", "4": "Consultant Fee"},
}

func kode(p string) string { return AwalanKode + p }

// ---------------------------------------------------------------- blok bersama

// blokTreaty - Layout S4 "Claim Treaty" (OutstandingClaim / InputAcceptation): baris tombol Inline, lalu Inline grid
// double - kolom kiri 8 medan, kolom kanan 5 medan.
func blokTreaty(tombolAtas ...Unsur) []Unsur {
	return []Unsur{
		sebaris("", tombolAtas...),
		dua(
			bagian("",
				ro(medan(CD+"IDMaster", "Treaty ID", KTeks)),
				ro(medan(CD+"TreatyName", "Treaty Name", KTeks)),
				ro(medan(TM+"ProportionType", "R/I Type", KTeks)),
				ro(medan(OQ+"BusinessName", "Class of Business", KTeks)),
				ro(medan(TM+"Ceding", "Ceding Name", KTeks)),
				ro(medan(TM+"LeadingReinsSource", "SOB Name", KTeks)),
				ro(medan(TM+"Bordeaux", "Bordereaux", KTampil)),
				ro(medan(TM+"BordereauxNote", "Bordereaux Note", KTeks)),
			),
			bagian("",
				ro(medan(CD+"YearofAccount", "Treaty Year", KTeks)),
				ro(medan(CD+"StartDateTreaty", "Treaty Start Date", KTanggal)),
				ro(medan(CD+"EndDateTreaty", "Treaty End Date", KTanggal)),
				ro(medan(TM+"AccountingMode", "Accounting Mode", KTampil)),
				ro(medan(TM+"TeritorialScope", "Teritorial Scope", KArea)),
			),
		),
	}
}

// blokKatastrofe - Section `Catastrope_Sec`.
func blokKatastrofe() Unsur {
	return bagian("",
		aksi(sumber(roJika(medan(CD+"StsKatastrofe", "Catastrophe", KRadio), beda(CD+"EditCatastrope", "true")),
			kode("StsKatastrofe")), "SetDefNonCatastrope"),
		aksi(sumber(tampil(roJika(medan(CD+"NonKatastrofeType", "", KRadio), beda(CD+"EditCatastrope", "true")),
			nonCatastrophe), kode("NonKatastrofeType")), "SetDefNonCatastrope"),
		tampil(bagian("",
			ro(medan(CD+"KatastrofeNote", "Catastrophe Note", KTeks)),
			tampil(tombol("CatastrofeList", "", "BukaKatastrofe"), editCatastrope),
		), atau(sama(CD+"StsKatastrofe", "Catastrophe"), sama(CD+"NonKatastrofeType", "Big Claim"))),
		tampil(ikon(tombol("EditCatastrope", "", "SetEditCatastrope:Edit"), IkonUbah), beda(CD+"EditCatastrope", "true")),
		tampil(ikon(tombol("SaveCatastrope", "", "SetEditCatastrope:Save"), IkonSimpan), editCatastrope),
	)
}

// blokPelapor - medan registrasi klaim (Section OutstandingClaim / InputAcceptation; `relasiStatus` = properti yang
// membuka Specify / Reporter Address: ReporterStatus di Outstanding, InsuredRelationship di InputAcceptation).
func blokPelapor(relasiStatus string) []Unsur {
	bukanTiga := beda(relasiStatus, "3")
	return []Unsur{
		aksi(wajibU(roJika(medan(CD+"DateOfLoss", "Date of Loss", KTanggal), isAnyAccept)), "CheckDateDOL"),
		aksi(wajibU(roJika(medan(CD+"ReportDate", "Report Date", KTanggal), isAnyAccept)), "CheckReportDate"),
		aksi(wajibJ(roJika(medan(CD+"DateReceived", "Received Date", KTanggal), isAnyAccept), pyNoteKosong),
			"CheckDateReceived"),
		wajibJ(roJika(medan(CD+"ReporterName", "Reporter Name", KTeks), isOutstanding), pyNoteKosong),
		wajibU(roJika(medan(CD+"ReporterTelp", "Reporter Phone Number", KTelepon), isOutstanding)),
		sumber(roJika(medan(CD+"ReportType", "Report Type", KPilih), isOutstanding), kode("ReportType")),
		aksi(sumber(naJika(roJika(medan(CD+"ReporterStatus", "Reporter Status", KPilih), isEstimation), isOutstanding),
			kode("ReporterStatus")), "GetReportStatus"),
		roJika(medan(CD+"InsuredRelationshipOthers", "Specify...", KTeks), bukanTiga),
		wajibJ(roJika(medan(CD+"ReportAddress", "Reporter Address", KArea), atau(isOutstanding, bukanTiga)),
			pyNoteKosong),
		blokKatastrofe(),
	}
}

// blokAdjuster - Consultant | Adjuster (Section OutstandingClaim / InputAcceptation): Inline grid double, ikon
// `pyWorkActionsAddWork.png` di samping ID.
//
// [keputusan work owner 08-10-2026] "yang dropdown hanya dari namanya aja; untuk ID dihapus dari tampilan, tapi tetap
// simpan ID": medan ID (yang disimpan, aksi SetConsultant / SetAdjsuter tetap) berlabel nama dan menampilkan jalur
// nama; baris nama hanya-baca XML hanya muncul saat dropdown tersembunyi (IsAnyAcceptation = 1), supaya nama tidak
// tampil dua kali. Label XML "Consultant ID" / "Adjuster / Professional ID" tidak dipakai.
func blokAdjuster(namaTampil Kondisi, namaAdjTampil Kondisi) []Unsur {
	bukanAcc := beda("IsAnyAcceptation", "1")
	sudahAcc := sama("IsAnyAcceptation", "1")
	return []Unsur{dua(
		bagian("",
			tampilan(aksi(sumber(tampil(wajibU(roJika(medan(CD+"ConsultantID", "Consultant Name", KOtomatis), isOutstanding)),
				bukanAcc), SumberAdjuster), "SetConsultant"), CD+"ConsultantName"),
			ikon(tampil(naJika(tombol("AdjusterConsultantBaru1", "Add", AksiTambahKonsultan), isOutstanding), bukanAcc), IkonTambah),
			tampil(ro(medan(CD+"ConsultantName", "Consultant Name", KTeks)), dan(namaTampil, sudahAcc)),
		),
		bagian("",
			tampilan(aksi(sumber(tampil(wajibU(roJika(medan(CD+"AppointedADJID", "Adjuster / Professional Name", KOtomatis),
				isOutstanding)), bukanAcc), SumberAdjuster), "SetAdjsuter"), CD+"AppointedADJ"),
			ikon(tampil(naJika(tombol("AdjusterConsultantBaru2", "Add", AksiTambahAdjuster), isOutstanding), bukanAcc), IkonTambah),
			tampil(ro(medan(CD+"AppointedADJ", "Adjuster / Professional Name", KTeks)), dan(namaAdjTampil, sudahAcc)),
		),
	)}
}

// blokLokasi - Report Description dan Location of Loss (Occupation, Province, Zip Code ditulis per layar - kondisinya
// berbeda).
func blokLokasi() []Unsur {
	return []Unsur{
		aksi(roJika(medan(CD+"ReportDescription", "Report Description", KArea), isOutstanding), "MakeLowercase"),
		aksi(wajibU(roJika(medan(CD+"Location", "Location of Loss", KArea), isOutstanding)), "MakeLowercase"),
	}
}

// gridRiwayat - Layout "Claim History" (grid SuggestList, paging 5; urutan `TampilRiwayat`), tersimpan di
// T_VIEW_SUGGEST.
func gridRiwayat() Unsur {
	return Unsur{Jenis: JenisBagian, Label: "Claim History", Anak: []Unsur{{
		Jenis: JenisGrid, Jalur: DaftarRiwayatTampil, Bernomor: true, PerHalaman: 5,
		Kolom: []Unsur{kRO(kol("IsCedingConfirm", "Name", KTampil)), kRO(kol("DateSuggest", "Date", KWaktu)),
			kRO(kol("CommentSuggest", "Noted", KTampil))},
	}}}
}

// DaftarRiwayatTampil - salinan SuggestList untuk grid (label tingkat + urut menurun, `TampilRiwayat`).
const DaftarRiwayatTampil = "ClaimHistory"

// sectionInterest - Section `OutstandingClaim_Intrs` (dipakai juga InputAcceptation_Intrs).
func sectionInterest() Unsur {
	return bagian("",
		bagian("Insured Interests 100 %", Unsur{
			Jenis: JenisGrid, Jalur: DaftarInterest, Bernomor: true,
			Tambah: ptr(ikon(tombol("AddInterest", "Add", "AddInterest"), IkonTambah)),
			Kolom: []Unsur{
				kROJ(kol("ObjectName", "Insured Interest", KTeks), bIsAdjVal),
				kAksi(kSumber(kROJ(kol("CurrencyID", "Currency", KPilih), bIsAdjVal), SumberMataUang), "SetCurencyInterest"),
				kRO(kol("KursObjectItem", "Value In IDR", KAngka)),
				kAksi(kROJ(kol("TSIPerObject", "Value", KAngka), bIsAdjVal), "CountTotalInsterest"),
				ikonK(kTombol("DeleteInterest", "Delete", "DeleteInterest", bIsAdjVal), IkonHapus),
			},
		}),
		dua(
			bagian("Total in Original Currency", Unsur{
				Jenis: JenisGrid, Jalur: DaftarTotalTSI, Bernomor: true,
				Kolom: []Unsur{kRO(kol("Currency", "Total", KTampil)), kRO(kol("Value", "Value", KAngka))},
			}),
			bagian("Total In IDR", sebaris("", label("IDR"), ro(medan(CD+"TotalSumInsuredIDR", "", KAngka)))),
		),
		medan(CD+"InsuredInterest", "Description", KArea),
	)
}

func ptr(u Unsur) *Unsur { return &u }

// blokDeductible - Layout "Deductible Info" (OutstandingClaim_Est / InputAcceptation_Est; identik).
func blokDeductible() Unsur {
	ro1 := isOutstanding
	return bagian("Deductible Info",
		aksi(naJika(medan(CD+"DeductibleType", "Deductible", KCentang), ro1), "SetFormat"),
		sumber(tampil(roJika(medan(CD+"FormType", "Format", KPilih), ro1), sama(CD+"DeductibleType", "true")), kode("FormType")),
		tampil(bagian("",
			sumber(roJika(medan(CD+"CurrencyDeductible", "Currency", KPilih), ro1), SumberMataUang),
			roJika(medan(CD+"NetDeductibleValue", "Amount", KAngka), ro1),
		), sama(CD+"FormType", "1")),
		tampil(bagian("",
			aksi(roJika(medan(CD+"Amount", "%", KAngka), ro1), "CountDeductible"),
			label("of"),
			aksi(sumber(roJika(medan(CD+"TypeDeductible", "", KPilih), ro1), kode("TypeDeductible")), "CountDeductible"),
			tampil(aksi(roJika(medan(CD+"TSIDeductible", "TSI Amount", KAngka), ro1), "CountDeductible"),
				sama(CD+"TypeDeductible", "2")),
			label("Minimum"),
			aksi(sumber(roJika(medan(CD+"CurrencyDeductible", "Currency", KPilih), ro1), SumberMataUang), "CountDeductible"),
			aksi(roJika(medan(CD+"DeductibleValue", "Amount", KAngka), ro1), "CountDeductible"),
			ro(medan(CD+"NetDeductibleValue", "Deductible Value", KAngka)),
		), sama(CD+"FormType", "2")),
	)
}

// gridClaimAmount - grid ListClaimAmount (OutstandingClaim_Est: kolom Net Deductible; InputAcceptation_Est: tanpa).
func gridClaimAmount(denganNet bool) Unsur {
	kolom := []Unsur{
		kAksi(kSumber(kROJ(kol("CurrencyID", "Currency", KPilih), bNote), SumberMataUang), "SetCurencyList"),
		kAksi(kROJ(kol("ClaimAmount", "Claim Amount 100%", KAngka), bNote), "CountListClaimAmountIDR"),
	}
	if denganNet {
		kolom = append(kolom, kRO(kol("NetDeductibleValue", "Net Deductible", KAngka)))
	}
	kolom = append(kolom,
		kRO(kol("Value", "Claim Amount Ceding", KAngka)),
		kRO(kol("USD", "Claim Amount in IDR", KAngka)),
		ikonK(kTombol("DeleteListClaim", "Delete", "DeleteListClaim", bNote), IkonHapus),
	)
	return Unsur{Jenis: JenisGrid, Jalur: DaftarClaimAmount, Bernomor: true,
		Tambah: ptr(ikon(tombol("AddListClaimAmount", "Add", "AddListClaimAmount"), IkonTambah)), Kolom: kolom,
		Kaki: []Unsur{label("Total Claim Amount"), ro(medan(CD+"TotalListClaimAmount", "", KAngka)),
			ro(medan(CD+"TotalListClaimAmountIDR", "", KAngka))}}
}

// gridLossAllocation - Layout "Loss Allocation" (OutstandingClaim_Est / InputAcceptation_Est).
func gridLossAllocation(acc bool) Unsur {
	cur := kSumber(kROJ(kol("CurrencyID", "Curr", KPilih), bOldData), SumberMataUang)
	if acc {
		cur = kAksi(cur, "SetCurrency")
	}
	return bagian("Loss Allocation", Unsur{Jenis: JenisGrid, Jalur: DaftarLossAlloc, Bernomor: true,
		Tambah: ptr(ikon(tombol("AddLossAllocation", "Add", "AddLossAllocation"), IkonTambah)),
		Kolom: []Unsur{
			cur,
			kAksi(kSumber(kROJ(kol("TreatyName", "Treaty Type", KPilih), bOldData), SumberLimits), "SetNameTreaty"),
			kAksi(kROJ(kol("SharePercentage", "Share(%)", KAngka), bOldData), "CountPersen"),
			kRO(kol("ClaimSpreaded", "Result Claim", KAngka)),
			kRO(kol("ClaimEstimation", "Result Claim in IDR", KAngka)),
			ikonK(kTombol("RemoveLossAlloction", "Delete", "RemoveLossAlloction", bOldData), IkonHapus),
		}})
}

// gridEstimasi - Layout "Estimation List".
func gridEstimasi(aksiGross string) Unsur {
	return bagian("Estimation List", Unsur{Jenis: JenisGrid, Jalur: DaftarEstimasi, Bernomor: true,
		Tambah: ptr(ikon(tombol("AddEstimation", "Add", "AddEstimation"), IkonTambah)),
		Kolom: []Unsur{
			kRO(kol("TypeLoss", "", KTampil)),
			kAksi(kROJ(kol("EstimationDate", "Estimation Date", KTanggal), bPrintFace), "CheckEstimateDate"),
			kSumber(kROJ(kol("Type", "Type", KPilih), bPrintFace), kode("EstimationType")),
			kAksi(kSumber(kROJ(kol("CurrencyID", "Currency", KPilih), bPrintFace), SumberMataUang), "CurencyEstimation"),
			kRO(kol("KursValue", "Value In IDR", KAngka)),
			kAksi(kROJ(kol("GrossEstimationPct", "Gross Estimate Treaty (100%)", KAngka), bPrintFace), aksiGross),
			kRO(kol("EstimationValue", "Estimation RNM", KAngka)),
			kRO(kol("ConvertValue", "Estimation RNM in IDR", KAngka)),
			ikonK(kTombol("DeleteEstimation", "Delete", "DeleteEstimation", bPrintFace), IkonHapus),
		}})
}

// blokTotalEstimasi - "Total Original Currency Estimation" + "Total Estimation In IDR".
func blokTotalEstimasi() []Unsur {
	return []Unsur{
		bagian("Total Original Currency Estimation", Unsur{Jenis: JenisGrid, Jalur: DaftarTotalEst, Bernomor: true,
			Kolom: []Unsur{kRO(kol("Currency", "Currency", KTampil)), kRO(kol("IDR", "Gross Estimate Treaty (100%)", KAngka)),
				kRO(kol("Value", "Estimation RNM", KAngka))}}),
		bagian("Total Estimation In IDR",
			ro(medan(CD+"TotalGrossEstimateIDR", "Total Gross Estimate(100%) in IDR", KAngka)),
			ro(medan(CD+"TotalEstimasiIDR", "Total Estimation in IDR", KAngka)),
		),
	}
}

// gridSpreading - "Spreading List" pertama (SpreadingClaim), Section OutstandingClaim_Sprd / InputAcceptation_Est.
//
// [keputusan work owner 08-10-2026] Add / Delete AKTIF (membatalkan "non aktifkan" 07-10-2026); AddSpreading_Act /
// DeleteSpreading_Act tidak diekspor - perilakunya `AddSpreading` / `DeleteSpreading`. Add tanpa syarat, Delete
// nonaktif bila `.IsOldData='Yes'` (pyDisabledWhen XML). Ikon grid bawaan Pega tidak dibangun. Kolom Treaty Type di
// KEDUA layar = dropdown TreatyType bersumber spreading polis (XML Outstanding: pxTextInput .TreatyName - diubah karena
// baris dari Add harus memilih treaty), TERKUNCI bila baris sudah ber-TreatyType (dari polis, atau sudah dipilih
// sesudah Add) atau data lama (work owner "ini disable aja").
func gridSpreading() Unsur {
	tt := kAksi(kSumber(kROJ(kol("TreatyType", "Treaty Type", KPilih), bAtau(bOldData, bTerisi("TreatyType"))),
		SumberSpreading), "SetTreatyNameSpreading")
	return bagian("Spreading List", Unsur{Jenis: JenisGrid, Jalur: DaftarSpreading, Bernomor: true,
		Tambah: ptr(tombol("AddSpreading", "Add", "AddSpreading")),
		Kolom: []Unsur{
			tt,
			kAksi(kROJ(kol("SharePercentage", "Share(%)", KAngka), bOldData), "CountSpreading"),
			kRO(kol("Currency", "Currency", KTampil)),
			kAksi(kROJ(kol("ClaimSpreaded", "Claim Spreaded", KAngka), bOldData), "CountSpreading"),
			kTombol("DeleteSpreading", "Delete", "DeleteSpreading", bOldData),
		},
		Kaki: []Unsur{ro(medan(CD+"TotalEstimasi", "", KAngka))}})
}

func catatanK(u Unsur, c string) Unsur { u.Catatan = c; return u }

// gridBreakQS - "Spreading List" kedua (SpreadingBreakQS), hanya-baca.
func gridBreakQS(acc bool) Unsur {
	tt := kRO(kol("TreatyName", "Treaty Type", KTampil))
	if acc {
		tt = kSumber(kRO(kol("TreatyType", "Treaty Type", KPilih)), SumberJenisReas)
	}
	return bagian("Spreading List", Unsur{Jenis: JenisGrid, Jalur: DaftarBreakQS, Bernomor: true,
		Kolom: []Unsur{tt, kRO(kol("SharePercentage", "Share(%)", KAngka)), kRO(kol("Currency", "Currency", KTampil)),
			kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka))},
		Kaki: []Unsur{ro(medan(CD+"TotalEstimasi", "", KAngka))}})
}

// ---------------------------------------------------------------- Outstanding Claim

// LayarOutstanding - Section `OutstandingClaim` (FlowAction OutstandingClaim).
func LayarOutstanding() []Unsur {
	nomorKosong := dan(sama(CD+"NoClaim", ""), sama(CD+"ClaimNo", ""))
	// Inline grid double Claim Information: kolom kiri (polis), kolom kanan (pelapor), sel ketiga Cause of Loss.
	kiri := bagian("",
		tampil(bagian("Information", catatan(ro(medan("Message", "", KTeks)), OQHarnessHilang)), terisi("Message")),
		sebaris("",
			tampil(naJika(tombol("ChoosePolicy", "Choose Policy No", "PilihPolis"), isOutstanding), beda("IsOutstanding", "1")),
			tampil(tombol("ViewPolicy", "View", AksiLihatPolis), terisi(CD+"PolicyData.PolicyNo")),
			tampil(tombolOQ("PaymentPremi", "View Status Payment Premi", OQLayananLuar), terisi(CD+"PolicyData.PolicyNo")),
		),
		ro(medan(CD+"PolicyData.PolicyNo", "Policy No", KTeks)),
		sebaris("Quarter/Year",
			label("Q"), ro(medan(CD+"Quater", "", KTeks)), label("/"), ro(medan(CD+"YearofQuartal", "", KTeks)),
			label("U/Y"), ro(medan(CD+"TreatyYear", "", KTeks))),
		roJika(medan(CD+"PolicyNo", "Policy No Ceding", KTeks), isOutstanding),
		roJika(medan(CD+"InsuredName", "Insured Name", KTeks), isAnyAccept),
		roJika(medan(CD+"PlaNoCeding", "Pla No  Ceding", KTeks), isOutstanding),
		roJika(medan(CD+"PlaNoSOB", "Pla No SOB", KTeks), isOutstanding),
		naJika(medan(CD+"PeriodPolicyTBA", "Policy Period TBA ?", KCentang), isAnyAccept),
		aksi(wajibU(roJika(medan(CD+"PolicyData.StartDateTime", "Policy Start Ceding", KTanggal), isAnyAccept)), "SetEndDate"),
		aksi(wajibU(roJika(medan(CD+"PolicyData.EndDateTime", "Policy End Ceding", KTanggal), isAnyAccept)), "CheckPeriodPolicy"),
	)
	out := []Unsur{
		letak(LetakJudul,
			label("Outstanding Claim"),
			tampil(label("Claim No      .........."), nomorKosong),
			tampil(bagian("",
				tampilIsi(ro(medan(CD+"NoClaim", "Claim No", KTeks))),
				tampilIsi(ro(medan(CD+"ClaimNo", "Claim No", KTeks))),
			), atau(terisi(CD+"NoClaim"), terisi(CD+"ClaimNo"))),
		),
		bagian("Claim Treaty", blokTreaty(naJika(tombol("ChooseMaster", "Choose Master", "PilihMaster"), isOutstanding))...),
		label("Claim Information"),
		tampil(tombol("SummaryOutstanding", "Summary Outstanding Claim", "RingkasanOS"), isOutstanding),
		dua(kiri, bagian("", blokPelapor(CD+"ReporterStatus")...), dua(
			bagian("", ro(medan(CD+"CauseOfLoss", "Cause of Loss", KTampil))),
			bagian("", naJika(tombol("ChooseCauseOfLoss", "Choose Cause of Loss", "PilihSebab"), isAnyAccept)),
		)),
	}
	out = append(out, blokAdjuster(terisi(CD+"ConsultantName"), terisi(CD+"AppointedADJ"))...)
	out = append(out, blokLokasi()...)
	out = append(out,
		roJika(aksi(medan(CD+"Occupation", "Occupation", KArea), "MakeLowercase"), isOutstanding),
		dua(
			bagian("", sumber(wajibU(medan(CD+"Province", "Province", KOtomatis)), SumberProvinsi)),
			bagian("", aksi(medan(CD+"PostalCode", "Zip Code", KTeks), "GetAdders")),
		),
		sebaris("",
			tampil(ro(medan(TM+"RNMShareP", "RNM Share", KAngka)), beda("IsEditRNMShare", "true")),
			tampil(aksi(sumber(medan("TreatyShareTemp.CARI1", "RNM Share", KPilih), SumberShareRNM), "DisableEditRNMShare"),
				sama("IsEditRNMShare", "true")),
			label("%"),
			tampil(ikon(tombol("EditRNMShare", "", "GetRNMShareTreaty"), IkonUbah), estimasiPertamaBelumTerkirim),
		),
		letak(LetakTab,
			bagian("Interest", sectionInterest()),
			bagian("Estimation", sectionEstimasiOutstanding()),
			bagian("Spreading", bagian("Spreading Claim", gridSpreading(), gridBreakQS(false))),
		),
		tampil(tombol("Save", "Save", "SetOutstanding"), beda("IsAcceptation", "1")),
		naJika(tombol("SaveToIssueRNM", "Save to issue RNM", "SaveOutstanding"), sama("IsCFS", "")),
		naJika(tombol("PrintPLA", "PRINT PLA", "BukaPLA"), beda(CD+"IsPLA", "1")),
		tampil(naJika(tombol("SendToAcceptation", "Send to Acceptation", "CheckNopolicy"), isAcceptation), isOutstanding),
		tampil(tombol("Submit", "Submit", "SubmitOutstanding"), isAcceptation),
		gridRiwayat(),
	)
	return out
}

// estimasiPertamaBelumTerkirim = `pyWorkPage.ClaimData.EstimationList(1).PrintFaceClaim = ”`.
func estimasiPertamaBelumTerkirim(h *Halaman) bool {
	d := h.AmbilDaftar(DaftarEstimasi)
	return len(d) == 0 || d[0]["PrintFaceClaim"] == ""
}

// sectionEstimasiOutstanding - Section `OutstandingClaim_Est`.
func sectionEstimasiOutstanding() Unsur {
	anak := []Unsur{
		roJika(medan(CD+"ShareCeding", "Share Ceding(%)", KAngka), isOutstanding),
		blokDeductible(),
		gridClaimAmount(true),
		gridLossAllocation(false),
		sebaris("", ro(medan(TM+"RNMShareP", "RNM Share", KAngka)), label("%")),
		gridEstimasi("CountEstimation"),
	}
	anak = append(anak, blokTotalEstimasi()...)
	return bagian("", anak...)
}

// ---------------------------------------------------------------- Input Acceptation

// LayarAkseptasi - Section `InputAcceptation` (FlowAction InputAcceptation).
func LayarAkseptasi() []Unsur {
	kiri := bagian("",
		sebaris("",
			tampil(tombol("ViewPolicy", "View", AksiLihatPolis), terisi(CD+"PolicyData.PolicyNo")),
			tampil(tombolOQ("PaymentPremi", "View Status Payment Premi", OQLayananLuar), terisi(CD+"PolicyData.PolicyNo")),
		),
		ro(medan(CD+"PolicyData.PolicyNo", "Policy No", KTeks)),
		roJika(medan(CD+"PolicyNo", "Policy No Ceding", KTeks), isOutstanding),
		roJika(medan(CD+"InsuredName", "Insured Name", KTeks), isAnyAccept),
		roJika(medan(CD+"PlaNoCeding", "Pla No Ceding", KTeks), isAnyAccept),
		roJika(medan(CD+"PlaNoSOB", "Pla No SOB", KTeks), isAnyAccept),
		naJika(medan(CD+"PeriodPolicyTBA", "Policy Period TBA ?", KCentang), isAnyAccept),
		aksi(roJika(medan(CD+"PolicyData.StartDateTime", "Policy Start", KTanggal), isAnyAccept), "SetEndDate"),
		aksi(roJika(medan(CD+"PolicyData.EndDateTime", "Policy End", KTanggal), isAnyAccept), "CheckPeriodPolicy"),
		dua(
			bagian("", ro(medan(CD+"CauseOfLoss", "Cause of Loss", KTampil))),
			bagian("", naJika(tombol("ChooseCauseOfLoss", "Choose Cause of Loss", "PilihSebab"), isAnyAccept)),
		),
	)
	out := []Unsur{
		letak(LetakJudul,
			label("Acceptation Claim"),
			tampilIsi(ro(medan(CD+"NoClaim", "Claim No", KTeks))),
		),
		bagian("Claim Treaty", blokTreaty(
			tombol("SummaryOutstanding", "Summary Outstanding Claim", "RingkasanOS"),
			tombolOQ("PaymentPremiTreaty", "View Status Payment Premi", OQLayananLuar),
			tombol("CloseClaim", "Close Claim", "BukaTutupKlaim"),
		)...),
		label("Claim Information"),
		dua(kiri, bagian("", blokPelapor(CD+"InsuredRelationship")...)),
	}
	out = append(out, blokAdjuster(atau(terisi(CD+"ConsultantName"), isAnyAccept), atau(terisi(CD+"AppointedADJ"), isAnyAccept))...)
	out = append(out, blokLokasi()...)
	out = append(out,
		wajibU(roJika(aksi(medan(CD+"Occupation", "Occupation", KArea), "MakeLowercase"), isOutstanding)),
		dua(
			bagian("", roJika(sumber(wajibU(medan(CD+"Province", "Province", KOtomatis)), SumberProvinsi), isOutstanding)),
			bagian("", roJika(aksi(medan(CD+"PostalCode", "Zip Code", KTeks), "GetAdders"), isOutstanding)),
		),
		letak(LetakTab,
			bagian("Interests", sectionInterest()),
			bagian("Estimation", sectionEstimasiAkseptasi(), tombol("SaveEstimation", "Save", "Simpan")),
			bagian("Acceptation", sectionAdjs(), tombol("SaveAcceptation", "Save", "Simpan")),
		),
		gridRiwayat(),
	)
	return out
}

// sectionEstimasiAkseptasi - Section `InputAcceptation_Est`.
func sectionEstimasiAkseptasi() Unsur {
	anak := []Unsur{
		roJika(medan(CD+"ShareCeding", "Share Ceding(%)", KAngka), isAnyAccept),
		blokDeductible(),
		gridClaimAmount(false),
		sebaris("", ro(medan(TM+"RNMShareP", "RNM Share", KAngka)), label("%")),
		gridLossAllocation(true),
		gridEstimasi("CountEstimation"),
	}
	anak = append(anak, blokTotalEstimasi()...)
	anak = append(anak,
		bagian("Spreading Claim", gridSpreading(), gridBreakQS(true)),
		naJika(tombol("SaveToIssueRNM", "Save to issue RNM", "SaveOutstanding"), beda("ReCFS", "1")),
		naJika(tombol("PrintPLA", "PRINT PLA", "BukaPLA"), beda(CD+"IsPLA", "1")),
	)
	return bagian("", anak...)
}

// sectionAdjs - Section `InputAcceptation_Adjs`. Grid ClaimAmount / Loss Allocation / Estimation / Spreading di sini
// hanya-baca selalu (read-only selalu di XML); `IsEditEstimation` visible NEVER; tombol "Update Estimation" dan "Save"
// di Layout S33 visible NEVER - tidak dibangun.
func sectionAdjs() Unsur {
	semuaRO := func(us ...Unsur) []Unsur {
		for i := range us {
			if us[i].Jenis == JenisMedan {
				us[i].HanyaBacaB = hSelalu
			}
		}
		return us
	}
	return bagian("",
		label("Claim Estimation"),
		bagian("Count Claim Amount", Unsur{Jenis: JenisGrid, Jalur: DaftarClaimAmount, Bernomor: true,
			Kolom: semuaRO(kSumber(kol("CurrencyID", "Currency", KPilih), SumberMataUang), kol("Value", "Claim Amount", KAngka),
				kol("USD", "Claim Amount in IDR", KAngka))}),
		bagian("Loss Allocation", Unsur{Jenis: JenisGrid, Jalur: DaftarLossAlloc, Bernomor: true,
			Kolom: semuaRO(kSumber(kol("CurrencyID", "Curr", KPilih), SumberMataUang),
				kSumber(kol("TreatyType", "Treaty Type", KPilih), SumberJenisReas4), kol("SharePercentage", "Share(%)", KAngka),
				kol("ClaimSpreaded", "Result Claim", KAngka), kol("ClaimEstimation", "Result Claim In IDR", KAngka))}),
		bagian("Estimation List", Unsur{Jenis: JenisGrid, Jalur: DaftarEstimasi, Bernomor: true,
			Kolom: semuaRO(kol("TypeLoss", "", KTampil), kol("EstimationDate", "Estimation Date", KTanggal),
				kSumber(kol("Type", "Type", KPilih), kode("EstimationType")), kSumber(kol("CurrencyID", "Currency", KPilih), SumberMataUang),
				kol("KursValue", "Value In IDR", KAngka), kol("GrossEstimationPct", "Gross Estimate Treaty (100%)", KAngka),
				kol("EstimationValue", "Estimation RNM", KAngka), kol("ConvertValue", "Estimation RNM in IDR", KAngka))}),
		blokTotalEstimasi()[0], blokTotalEstimasi()[1],
		bagian("Spreading Claim",
			bagian("Spreading List", Unsur{Jenis: JenisGrid, Jalur: DaftarSpreading, Bernomor: true,
				Kolom: []Unsur{
					kROJ(kol("TreatyName", "Treaty Type", KTeks), spreadAdjRO),
					kROJ(kol("SharePercentage", "Share(%)", KAngka), spreadAdjRO),
					kRO(kol("Currency", "Currency", KTampil)),
					kROJ(kol("ClaimSpreaded", "Claim Spreaded", KAngka), bOutstandng),
				}, Kaki: []Unsur{ro(medan(CD+"TotalEstimasi", "", KAngka))}}),
			bagian("Spreading List", Unsur{Jenis: JenisGrid, Jalur: DaftarBreakQS, Bernomor: true,
				Kolom: []Unsur{kSumber(kRO(kol("TreatyType", "Treaty Type", KPilih)), SumberJenisReas),
					kRO(kol("SharePercentage", "Share(%)", KAngka)), kRO(kol("Currency", "Currency", KTampil)),
					kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka))},
				Kaki: []Unsur{ro(medan(CD+"TotalEstimasi", "", KAngka))}}),
		),
		label("Acceptance Information"),
		tombolOQ("PaymentClaim", "View Status Payment Claim", OQLayananLuar),
		tombolOQ("PaymentAttachment", "View Payment Attachment", OQLayananLuar),
		bagian("Acceptation List", Unsur{Jenis: JenisGrid, Jalur: DaftarAdjustment, Bernomor: true,
			Tambah: ptr(ikon(tampil(tombol("AddAdjustment", "Add", "AddAdjustment"),
				atau(sama("AktifButton", "0"), sama("AktifButton", ""))), IkonTambah)),
			Kolom: []Unsur{
				kSumber(kRO(kol("Type", "Type", KPilih)), kode("AdjustmentType")),
				kSumber(kRO(kol("AcceptanceStatus", "Status", KPilih)), kode("AcceptanceStatus")),
				kRO(kol("AcceptedNo", "Acceptation No", KTampil)),
				kRO(kol("AcceptedDate", "Acceptation Date", KTanggal)),
				kRO(kol("pxCreateOpName", "PIC Name", KTampil)),
				ikonK(kTombol("DeleteAjsutment", "Delete", "DeleteAjsutment", func(_ *Halaman, b Baris) bool {
					return b["IsKomite"] == "1" || b["IsSubjectivity"] == "true"
				}), IkonHapus),
			}}),
		bagian("Spreading Adjustment Total",
			bagian("Spreading In", Unsur{Jenis: JenisGrid, Jalur: DaftarSpreadAdj, Bernomor: true,
				Kolom: []Unsur{kRO(kol("TreatyName", "Treaty Type", KTeks)), kRO(kol("SharePercentage", "Share(%)", KAngka)),
					kRO(kol("Currency", "Currency", KTampil)), kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka))},
				Kaki: []Unsur{ro(medan(CD+"TotalEstimasi", "", KAngka))}}),
			bagian("Spreading Out", Unsur{Jenis: JenisGrid, Jalur: DaftarSpreadAdjQS, Bernomor: true,
				Kolom: []Unsur{kSumber(kRO(kol("TreatyType", "Treaty Type", KPilih)), SumberJenisReas),
					kRO(kol("SharePercentage", "Share(%)", KAngka)), kRO(kol("Currency", "Currency", KTampil)),
					kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka))},
				Kaki: []Unsur{ro(medan(CD+"TotalEstimasi", "", KAngka))}}),
		),
	)
}

// spreadAdjRO = `pyWorkPage.IsOutstanding = 1 && pyWorkPage.IsEditEstimation != 'true'`.
func spreadAdjRO(h *Halaman, _ Baris) bool {
	return h.Ambil("IsOutstanding") == "1" && h.Ambil("IsEditEstimation") != "true"
}

// ---------------------------------------------------------------- pilihan

// AdaKode - kunci sumber kode `associated`.
func AdaKode(s string) (string, bool) {
	if !strings.HasPrefix(s, AwalanKode) {
		return "", false
	}
	p := strings.TrimPrefix(s, AwalanKode)
	_, ada := KodePilihan[p]
	return p, ada
}
