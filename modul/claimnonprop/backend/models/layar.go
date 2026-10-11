package models

// Untuk apa berkas ini: DEFINISI LAYAR - Section `OutstandingClaim` (FlowAction OutstandingClaim, Assignment2) dan
// `InputAcceptation` (FlowAction InputAcceptation, Assignment1) beserta section yang di-include-nya (`Catastrope_Sec`),
// plus daftar pilihan kode `associated` dan catatan OQ. Detail akseptasi dan modal ada di `layar_akseptasi.go`.
// Label, urutan, dan kondisi VERBATIM dari XML `Claim Non Prop/Section`.
//
// Aturan baca:
//   - `Aksi` = nama activity tanpa akhiran `_Act` / `_act` / `_ACT` (lihat `services/aksi.go`); "" = postValue saja
//     (nilai disimpan, turunan dihitung ulang). Aksi berawalan `Pilih` / `Buka` / `Lihat` membuka pop-up di layar.
//   - Tombol yang activity / harness-nya TIDAK diekspor tampil sesuai section tetapi NONAKTIF dengan `Catatan` OQ.
//   - Sel ber-visible `NEVER` / `1=2` dan placeholder `.pyTemplate*` tidak dibangun (container ALWAYS berisi NEVER
//     tetap tampil, aturan R2 Claim Prop).
//   - Kondisi tombol header grid (`r1`) dievaluasi atas pyWorkPage (konteks section), sel baris atas barisnya.

// Catatan OQ untuk rule / layanan yang tidak tersedia (docs/OQ.md).
const (
	OQSandiEditXOL     = "OQ-CNP-04: kata sandi Edit XOL Allocation menunggu jalur konfigurasi sah (tanpa sandi di repo)"
	OQGridSpreading    = "OQ-CNP-11: grid Spreading tampil hanya-baca"
	OQDLA              = "OQ-CNP-19: activity GenerateDLACNP tidak ada di ekspor Pega"
	OQViewMaster       = "OQ-CNP-20: harness InputTreatyInOffer (View Master) tidak ada di ekspor Pega"
	OQShareSpreading   = "OQ-CNP-21: CountSpreadingCNP_Act tidak ada di ekspor Pega - share hanya-baca"
	OQDokumen          = "OQ-CNP-22: stream HTML dokumen tidak diekspor - nomor dan data tersimpan, berkas belum dibuat"
	OQKelasAkseptasi   = "OQ-CNP-33: activity dipanggil dengan kelas Data-Adjustment, tidak ada di ekspor untuk kelas itu"
	OQLayananBaca      = "OQ-CNP-34: layanan REST luar (M_LINK_SERVICE) belum disetujui dipanggil dari aplikasi"
	OQKomiteTanpaBaris = "OQ-CNP-36: kasus komite CWP tanpa baris akseptasi tidak dapat ditulis (T_GENERAL_KOMITE.ADJUSTMENT_ID NOT NULL)"
)

// Kunci daftar pilihan (`services` mengisinya).
const (
	SumberMataUang  = "mataUang"  // BrowseCurrency_RD
	SumberPolis     = "polis"     // pageList Polis.pxResults (GetDataPolisNonProp_SQL, InputOutStandingClmTNP_PreAct 3)
	SumberAdjuster  = "adjuster"  // BrowseAdjusterConsultant
	SumberProvinsi  = "provinsi"  // BrowseProvince_RD Nation INDONESIA
	SumberLossAlloc = "lossAlloc" // pageList ListLossAllocation.pxResults (InputOutStandingClmTNP_PreAct 4)
	SumberTreaty    = "treaty"    // pageList ListTreaty.pxResults (InputOutStandingClmTNP_PreAct 3.6)
	SumberRekening  = "rekening"  // pageList Result.pxResults (SetPayableTreatyNP_Act / GetDataBankAccount*)
	SumberKlien     = "klien"     // pageList ClientName.pxResults (Payable To = 3, BrowseClientName_RD)
	AwalanKode      = "kode:"     // prompt values `associated` - tidak terbaca di XML; kode DB apa adanya
)

// KodePilihan - kode bersumber `associated` (prompt values tidak diekspor, OQ-CNP-14): kode yang disebut kondisi section /
// activity XML dan yang teramati di DEV (JSON_KLAIM CLMNP-, baca-saja 09-10-2026). Ditampilkan apa adanya.
var KodePilihan = map[string][]string{
	"ReportType":         {"1", "2", "3", "4", "5"},
	"ReporterStatus":     {"1", "2", "3"},
	"FormType":           {"1", "2"},                          // S30 `.FormType==1`, S37 `.FormType==2`
	"TypeDeductible":     {"1", "2"},                          // S36 `.TypeDeductible = 2`
	"MinMax":             {"1", "2"},                          // CNPDeducMinMax / CNPMinMax / CNPMinMax2 (SetDataDeductible_act, SetTPLNote_Act)
	"TPLFormat":          {"1", "2"},                          // InputDtlInterest S8 / S13
	"TPLType":            {"1", "2"},                          // InputDtlInterest TPLType / TPLType2
	"Payable":            {"1", "2", "3"},                     // AdjustmentDetailNP `.Payable==3`
	"PaymentType":        {"1", "2", "3", "4", "5", "6", "7"}, // SetInterimXOL_Act 2 = interim, 7 = cancellation
	"AdjustmentType":     {"1", "2", "3", "4"},
	"AcceptanceStatus":   {"1", "2"},
	"StsKatastrofe":      {"Catastrophe", "Non-Catastrophe"},
	"NonKatastrofeType":  {"Claim", "Big Claim"},
	"SubjectivityNote":   {"1"},
	PropKeputusanAnggota: {"1", "2", "0"},
}

// LabelKode - label tampilan kode untuk properti yang KELASNYA SAMA dengan properti Claim Prop yang labelnya diberikan
// work owner (screenshot 08-10-2026, ekspor `Komite Claim Prop/KomiteAproval.xml`): ASM-FW-GCNMFW-Data-ClaimData
// (ReportType, ReporterStatus, Payable), ASM-FW-GCNMFW-Data-Adjustment.Type, ASM-FW-GCNMFW-Data-Comitee.KomiteAproval.
// Satu property rule Pega = satu daftar prompt values (PARITAS `[inferensi]`). Properti khas Non Prop tanpa label.
var LabelKode = map[string]map[string]string{
	"ReportType":         {"1": "Direct", "2": "Via Email", "3": "Via Fax", "4": "via Postal Mail/Courier", "5": "Via Telephone"},
	"ReporterStatus":     {"1": "Ceding Co Name", "2": "SOB Name", "3": "Others"},
	"Payable":            {"1": "Ceding Co Name", "2": "Broker Name", "3": "Others"},
	"AdjustmentType":     {"1": "Claim", "2": "Adjuster Fee", "3": "Salvage", "4": "Consultant Fee"},
	PropKeputusanAnggota: {"1": "Approved", "2": "Reject", "0": "Waiting"},
}

func kode(p string) string { return AwalanKode + p }

func ptr(u Unsur) *Unsur { return &u }

// ---------------------------------------------------------------- blok bersama

// blokTreaty - Layout "Claim Treaty" (Outstanding S5-S10 / InputAcceptation S6-S10): baris tombol Inline, lalu dua
// kolom master. `proporsi` = InputAcceptation menampilkan R/I Type; `bukanTO` = Outstanding menyembunyikan Ceding / SOB
// bila `IDMasterTONP` terisi.
func blokTreaty(tombolAtas []Unsur, akseptasi bool) Unsur {
	tanpaTO := sama(CD+"IDMasterTONP", "")
	labelID := "Master Treaty ID"
	if akseptasi {
		labelID = "Treaty ID"
	}
	kiri := []Unsur{
		ro(medan(CD+"IDMaster", labelID, KTeks)),
		ro(medan(CD+"TreatyName", "Treaty Name", KTeks)),
	}
	if akseptasi {
		kiri = append(kiri, ro(medan(TM+"ProportionType", "R/I Type", KTeks)))
	}
	kiri = append(kiri, ro(medan(OQ+"BusinessName", "Class of Business", KTeks)))
	ceding, sob := ro(medan(TM+"Ceding", "Ceding Name", KTeks)), ro(medan(TM+"LeadingReinsSource", "SOB Name", KTeks))
	if !akseptasi {
		ceding, sob = tampil(ceding, tanpaTO), tampil(sob, tanpaTO)
	}
	kiri = append(kiri, ceding, sob,
		ro(medan(TM+"Bordeaux", "Bordereaux", KTampil)),
		ro(medan(TM+"BordereauxNote", "Bordereaux Note", KTeks)))
	return bagian("Claim Treaty",
		sebaris("", tombolAtas...),
		dua(
			bagian("", kiri...),
			bagian("",
				ro(medan(CD+"YearofAccount", "Treaty Year", KTeks)),
				ro(medan(CD+"StartDateTreaty", "Treaty Start Date", KTanggal)),
				ro(medan(CD+"EndDateTreaty", "Treaty End Date", KTanggal)),
				ro(medan(TM+"AccountingMode", "Accounting Mode", KTampil)),
				ro(medan(TM+"TeritorialScope", "Teritorial Scope", KArea)),
			),
		),
	)
}

// blokKatastrofe - Section `Catastrope_Sec` (kelas ASM-FW-GCNMFW-Work, dipakai bersama Claim Prop).
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

// gridRiwayat - Layout "Claim History" (grid SuggestList, pyGridPaginator; urutan `TampilRiwayat`), tersimpan di
// T_VIEW_SUGGEST. InputAcceptation menambah kolom Initial.
func gridRiwayat(inisial bool) Unsur {
	kolom := []Unsur{kRO(kol("IsCedingConfirm", "Name", KTampil))}
	if inisial {
		kolom = append(kolom, kRO(kol("Initial", "", KTampil)))
	}
	kolom = append(kolom, kRO(kol("DateSuggest", "Date", KWaktu)), kRO(kol("CommentSuggest", "Noted", KTampil)))
	return bagian("Claim History", Unsur{Jenis: JenisGrid, Jalur: DaftarRiwayatTampil, Bernomor: true, PerHalaman: 5,
		Kolom: kolom})
}

// DaftarRiwayatTampil - salinan SuggestList untuk grid (urut menurun, `TampilRiwayat`).
const DaftarRiwayatTampil = "ClaimHistory"

// gridInterest - Layout "Insured Interests 100 %" (Outstanding: masterDetail -> InputDtlInterest, Add
// AddInterestListCNP_Act, Delete deleteRow; InputAcceptation: readOnly, tanpa Add / Delete / TPL).
func gridInterest(akseptasi bool) Unsur {
	kolom := []Unsur{
		kRO(kol("ObjectName", "Insured Interest", KTeks)),
		kSumber(kRO(kol("CurrencyID", "Currency", KPilih)), SumberMataUang),
		kRO(kol("KursObjectItem", "Value In IDR", KAngka)),
		kRO(kol("TSIPerObject", "Value", KAngka)),
	}
	g := Unsur{Jenis: JenisGrid, Jalur: DaftarInterest, Bernomor: true}
	if !akseptasi {
		kolom = append(kolom, kRO(kol("TPLAmount", "TPL", KAngka)),
			ikonK(kTombol("HapusInterest", "Delete", "HapusInterest", nil), IkonHapus))
		g.Tambah = ptr(ikon(tombol("AddInterestListCNP", "Add", "AddInterestListCNP"), IkonTambah))
	}
	g.Kolom = kolom
	return bagian("Insured Interests 100 %", g)
}

// sectionInterest - tab Interests (Outstanding S20 / InputAcceptation S12).
func sectionInterest(akseptasi bool) Unsur {
	desk := medan(CD+"InsuredInterest", "Description", KArea)
	anak := []Unsur{
		gridInterest(akseptasi),
		dua(
			bagian("Total in Original Currency", Unsur{Jenis: JenisGrid, Jalur: DaftarTotalTSI, Bernomor: true,
				Kolom: []Unsur{kRO(kol("Currency", "Total", KTampil)), kRO(kol("Value", "Value", KAngka))}}),
			bagian("Total In IDR", sebaris("", label("IDR"), ro(medan(CD+"TotalSumInsuredIDR", "", KAngka)))),
		),
	}
	if akseptasi { // Outstanding: Description tampil[1=3] - tidak dibangun
		anak = append(anak, ro(desk))
	}
	return bagian("", anak...)
}

// blokDeductible - Layout "Additional Info" (Outstanding S25-S43 / InputAcceptation). `kunci` = RO, `mati` = NA.
// Outstanding: RO `IsOutstanding==1`, NA `IsCFS==1`, aksi SetFormat / SetDataDeductible; InputAcceptation: RO / NA
// `IsAcceptation = 1`, Format dan Type Deductible postValue saja.
func blokDeductible(kunci, mati Kondisi, outstanding bool) Unsur {
	f := func(u Unsur) Unsur { return naJika(roJika(u, kunci), mati) }
	// Outstanding: Format -> SetDataDeductible_act; Type Deductible (FormType 1) dan Min/Max (FormType 2) ->
	// SetDataDeductible_act Posisi "TypeDeductible"; selebihnya postValue.
	formAksi, tipe1, minmax2 := "", "", ""
	if outstanding {
		formAksi, tipe1, minmax2 = "SetDataDeductible", "SetDataDeductible:TypeDeductible", "SetDataDeductible:TypeDeductible"
	}
	deduct := aksi(naJika(roJika(medan(CD+"DeductibleType", "Deductible", KCentang), kunci), mati), "SetFormat")
	if !outstanding {
		deduct = aksi(naJika(medan(CD+"DeductibleType", "Deductible", KCentang), kunci), "SetFormat")
	}
	mataUang := func() Unsur { return sumber(f(medan(CD+"CurrencyDeductible", "Currency", KPilih)), SumberMataUang) }
	tsi := func() Unsur {
		return tampil(f(medan(CD+"TSIDeductible", "TSI Amount", KAngka)), sama(CD+"TypeDeductible", "2"))
	}
	var bentuk1, bentuk2 []Unsur
	bentuk1 = append(bentuk1, mataUang(), f(medan(CD+"DeductibleValue", "Amount", KAngka)))
	if outstanding {
		bentuk1 = append(bentuk1, sumber(f(medan(CD+"CNPDeducMinMax", "", KPilih)), kode("MinMax")))
	} else {
		bentuk1 = append(bentuk1, label("Minimum"))
	}
	bentuk1 = append(bentuk1, f(medan(CD+"Amount", "%", KAngka)), label("of"),
		aksi(sumber(f(medan(CD+"TypeDeductible", "", KPilih)), kode("TypeDeductible")), tipe1), tsi())
	bentuk2 = append(bentuk2, f(medan(CD+"Amount", "%", KAngka)), label("of"),
		sumber(f(medan(CD+"TypeDeductible", "", KPilih)), kode("TypeDeductible")), tsi())
	if outstanding {
		bentuk2 = append(bentuk2, aksi(sumber(f(medan(CD+"CNPDeducMinMax", "", KPilih)), kode("MinMax")), minmax2))
	} else {
		bentuk2 = append(bentuk2, label("Minimum"))
	}
	bentuk2 = append(bentuk2, mataUang(), f(medan(CD+"DeductibleValue", "Amount", KAngka)))
	return bagian("Additional Info",
		deduct,
		aksi(sumber(tampil(f(medan(CD+"FormType", "Format", KPilih)), sama(CD+"DeductibleType", "true")), kode("FormType")),
			formAksi),
		tampil(sebaris("", bentuk1...), sama(CD+"FormType", "1")),
		tampil(sebaris("", bentuk2...), sama(CD+"FormType", "2")),
	)
}

// gridClaimAmount - grid ListClaimAmount (Outstanding S44 / InputAcceptation; identik). Kondisi Add header
// `.CNPFlagOuts==1` atas pyWorkPage hanya di Outstanding.
func gridClaimAmount(addTerkunci bool) Unsur {
	ubah := func(p, j string) Unsur {
		return kAksi(kNA(kol(p, j, KAngka), bTerkunci), "CountClaimTNP")
	}
	add := ikon(tombol("AddListClaimNP", "Add", "AddListClaimNP:Claim"), IkonTambah)
	if addTerkunci {
		add = naJika(add, sama("CNPFlagOuts", "1"))
	}
	return Unsur{Jenis: JenisGrid, Jalur: DaftarClaimAmount, Bernomor: true, Tambah: ptr(add),
		Kolom: []Unsur{
			kAksi(kSumber(kNA(kROJ(kol("CurrencyID", "Currency", KPilih), bNote), bTerkunci), SumberMataUang), "SetCurrency:Claim"),
			ubah("AltValue", "Rate of Exchange"),
			kROJ(ubah("Value", "Claim Amount"), bNote),
			kROJ(ubah("TPL", "TPL"), bNote),
			ubah("AdjusterFee", "Adjuster Fee"),
			ubah("Salvage", "Salvage"),
			ubah("CNPOthersFee", "Fee"),
			kRO(kol("PctProrateClaim", "Proportion (%)", KAngka)),
			kRO(kol("USD", "Claim Amount in IDR", KAngka)),
			kRO(kol("ClaimAmountCedant", "Claim Amount Cedant", KAngka)),
			ikonK(kTombol("HapusKlaim", "Delete", "HapusKlaim", bTerkunci), IkonHapus),
		},
		Kaki: []Unsur{label("Total Claim Amount"), ro(medan(CD+"TotalListClaimAmount", "", KAngka)),
			ro(medan(CD+"TotalListClaimAmountIDR", "", KAngka))}}
}

// gridLossAllocation - Layout "Loss Allocation" (CNPSpreadLoss). Outstanding: Adjuster Fee / Salvage / Fee postValue
// saja; InputAcceptation identik.
func gridLossAllocation(addTerkunci bool) Unsur {
	add := ikon(tombol("AddLossAlocation", "Add", "AddLossAlocation"), IkonTambah)
	if addTerkunci {
		add = naJika(add, sama("CNPFlagOuts", "1"))
	}
	return bagian("Loss Allocation", Unsur{Jenis: JenisGrid, Jalur: DaftarLossAlloc, Bernomor: true, Tambah: ptr(add),
		Kolom: []Unsur{
			kSumber(kRO(kol("CurrencyID", "Currency", KPilih)), SumberMataUang),
			kSumber(kNA(kol("TreatyName", "Treaty Name", KPilih), bTerkunci), SumberLossAlloc),
			kAksi(kNA(kol("ClaimPercentage", "Share (%)", KAngka), bTerkunci), "CountLossAllocation:pct"),
			kAksi(kNA(kol("ClaimAmountAdjust", "Claim Amount", KAngka), bTerkunci), "CountLossAllocation:amount"),
			kNA(kol("AdjusterFee", "Adjuster Fee", KAngka), bTerkunci),
			kNA(kol("Salvage", "Salvage", KAngka), bTerkunci),
			kNA(kol("CNPOthersFee", "Fee", KAngka), bTerkunci),
			kAksi(kNA(kol("CNPFlagXOL", "To XOL", KCentang), bTerkunci), "CountLossAllocation:CountXOL"),
			ikonK(kTombol("HapusLossAlloc", "Delete", "HapusLossAlloc", bTerkunci), IkonHapus),
		}})
}

// gridXOL - Layout "XOL Allocation" (SpreadingRisk). Outstanding: Claim Amount NA `.CNPFlagOuts==1`; InputAcceptation:
// tanpa kondisi.
func gridXOL(akseptasi bool) Unsur {
	total := kAksi(kol("TotalClaim", "Claim Amount", KAngka), "AdjClaimAmount")
	if !akseptasi {
		total = kNA(total, bTerkunci)
	}
	return bagian("XOL Allocation", Unsur{Jenis: JenisGrid, Jalur: DaftarXOL, Bernomor: true, Kolom: []Unsur{
		kSumber(kRO(kol("CurrencyID", "Currency", KPilih)), SumberMataUang),
		kRO(kol("TreatyName", "Treaty Name", KTampil)),
		total,
		kRO(kol("ClaimPercentage", "RNM Share (%)", KAngka)),
		kRO(kol("ClaimSpreaded", "Claim Amount RNM", KAngka)),
		kRO(kol("AdjusterFee", "Adjuster Fee", KAngka)),
		kRO(kol("Salvage", "Salvage", KAngka)),
		kRO(kol("CNPOthersFee", "Fee", KAngka)),
	}})
}

// gridSummaryXOL - "Summary XOL Allocation" (Outstanding S59, Claim Amount = .IDR) / "Total XOL Allocation"
// (InputAcceptation, Claim Amount = .CNPTotalClaim).
func gridSummaryXOL(akseptasi bool) Unsur {
	judul, klaim := "Summary XOL Allocation", "IDR"
	if akseptasi {
		judul, klaim = "Total XOL Allocation", "CNPTotalClaim"
	}
	return bagian(judul, Unsur{Jenis: JenisGrid, Jalur: DaftarSummaryXOL, Bernomor: true, Kolom: []Unsur{
		kRO(kol("Currency", "Currency", KTampil)),
		kRO(kol(klaim, "Claim Amount", KAngka)),
		kRO(kol("Value", "Claim Amount RNM", KAngka)),
		kRO(kol("AdjusterFee", "Adjuster Fee", KAngka)),
		kRO(kol("Salvage", "Salvage", KAngka)),
		kRO(kol("CNPOthersFee", "Fee", KAngka)),
	}})
}

// gridSpreadingKlaim - "Spreading List" (SpreadingClaim) dan kedua (SpreadingBreakQS). Bawaan OQ-CNP-11: hanya-baca;
// Add (AddListClaimNP_Act Table=Spreading) dan Delete (deleteRow) Outstanding tampil nonaktif; share = OQ-CNP-21.
func gridSpreadingKlaim(denganTombol bool) []Unsur {
	fee := []Unsur{kRO(kol("AdjusterFee", "Adjuster Fee", KAngka)), kRO(kol("Salvage", "Salvage", KAngka)),
		kRO(kol("CNPOthersFee", "Fee", KAngka))}
	atas := Unsur{Jenis: JenisGrid, Jalur: DaftarSpreading, Bernomor: true,
		Kolom: append([]Unsur{
			kSumber(kRO(kol("TreatyName", "Treaty Type", KPilih)), SumberTreaty),
			catatanK(kRO(kol("SharePercentage", "Share(%)", KAngka)), OQShareSpreading),
			kRO(kol("Currency", "Currency", KTampil)),
			kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka)),
		}, fee...),
		Kaki: []Unsur{ro(medan(CD+"TotalEstimasi", "", KAngka))}}
	if denganTombol {
		atas.Tambah = ptr(tombolOQ("AddListClaimNPSpreading", "Add", OQGridSpreading))
		atas.Kolom = append(atas.Kolom, Unsur{Jenis: JenisTombol, ID: "HapusSpreading", Label: "Delete",
			NonaktifB: hSelalu, Catatan: OQGridSpreading})
	}
	bawah := Unsur{Jenis: JenisGrid, Jalur: DaftarBreakQS, Bernomor: true,
		Kolom: append([]Unsur{
			kRO(kol("TreatyName", "Treaty Type", KTeks)),
			kRO(kol("SharePercentage", "Share(%)", KAngka)),
			kRO(kol("Currency", "Currency", KTampil)),
			kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka)),
		}, fee...),
		Kaki: []Unsur{ro(medan(CD+"TotalEstimasi", "", KAngka))}}
	return []Unsur{bagian("Spreading List", atas), bagian("Spreading List", bawah)}
}

func catatanK(u Unsur, c string) Unsur { u.Catatan = c; return u }

// ---------------------------------------------------------------- Outstanding Claim

// proteksiIssue - `ProtectDOL.CARI1 = 1 || ProtectStartDate.CARI1 = 1 || ProtectEndDate.CARI1 = 1` (tombol Save to issue
// RNM). ProtectEndDate tidak punya penulis aktif (CheckPeriodPolicy_Act ber-remark) - selalu kosong.
func proteksiIssue(h *Halaman) bool { return ProteksiDOL(h) || ProteksiMulaiPolis(h) }

// LayarOutstanding - Section `OutstandingClaim` (FlowAction OutstandingClaim).
func LayarOutstanding() []Unsur {
	ro1 := isOutstanding
	kiri := bagian("",
		tampil(bagian("Information", ro(medan("Message", "", KTeks))), terisi("Message")),
		wajibU(aksi(sumber(roJika(medan(CD+"PolicyData.PolicyNo", "Policy No", KOtomatis),
			dan(ro1, terisi(CD+"PolicyData.PolicyNo"))), SumberPolis), "CheckNoPolicy")),
		ro(medan(CD+"PolicyData.TreatyGroup", "Class of Business", KTeks)),
		sebaris("",
			tombol("ViewListPolicy", "View List Policy", "PilihPolis"),
			tampil(tombol("ViewPolicy", "View", AksiLihatPolis), terisi(CD+"PolicyData.PolicyNo")),
			tampil(tombolOQ("ViewPaymentStatus", "View Payment Status", OQLayananBaca), terisi(CD+"PolicyData.PolicyNo")),
		),
		roJika(medan(CD+"PolicyNo", "Policy No Ceding", KTeks), dan(terisi(CD+"PolicyNo"), ro1)),
		roJika(medan(CD+"InsuredName", "Insured Name", KTeks), ro1),
		roJika(medan(CD+"CNPReinsuranceSlip", "Reinsurance Slip", KTeks), ro1),
		roJika(medan(CD+"CNPClmNoCedant", "Claim No Ceding", KTeks), ro1),
		roJika(medan(CD+"PlaNoCeding", "Pla No  Ceding", KTeks), ro1),
		roJika(medan(CD+"DLANoCeding", "Dla No  Ceding", KTeks), ro1),
		roJika(medan(CD+"PlaNoSOB", "Pla No SOB", KTeks), ro1),
		wajibU(aksi(roJika(medan(CD+"PolicyData.StartDateTime", "Policy Start Ceding", KTanggal), ro1), "SetEndDate")),
		wajibU(roJika(medan(CD+"PolicyData.EndDateTime", "Policy End Ceding", KTanggal), ro1)),
	)
	kanan := bagian("",
		wajibU(aksi(roJika(medan(CD+"DateOfLoss", "Date of Loss", KTanggal), ro1), "CheckDateDOL")),
		wajibU(aksi(roJika(medan(CD+"ReportDate", "Report Date", KTanggal), ro1), "CheckReportDate")),
		wajibJ(aksi(roJika(medan(CD+"DateReceived", "Received Date", KTanggal), ro1), "CheckDateReceived"), pyNoteKosong),
		wajibJ(roJika(medan(CD+"ReporterName", "Reporter Name", KTeks), ro1), pyNoteKosong),
		wajibU(roJika(medan(CD+"ReporterTelp", "Reporter Phone Number", KTelepon), ro1)),
		blokKatastrofe(),
	)
	info := bagian("",
		dua(
			kiri,
			kanan,
			bagian("",
				ro(medan(CD+"CauseOfLoss", "Cause of Loss", KTampil)),
				tombol("ChooseCauseOfLoss", "Choose Cause of Loss", "PilihSebab"),
			),
			bagian("",
				aksi(sumber(naJika(roJika(medan(CD+"ReporterStatus", "Reporter Status", KPilih), isEstimation), ro1),
					kode("ReporterStatus")), "GetReportStatus"),
				roJika(medan(CD+"InsuredRelationshipOthers", "Specify...", KTeks), beda(CD+"ReporterStatus", "3")),
			),
			bagian("", sumber(roJika(medan(CD+"ReportType", "Report Type", KPilih), ro1), kode("ReportType"))),
			bagian("", wajibJ(roJika(medan(CD+"ReportAddress", "Reporter Address", KArea),
				atau(ro1, beda(CD+"ReporterStatus", "3"))), pyNoteKosong)),
			bagian("",
				aksi(sumber(roJika(medan(CD+"AppointedADJID", "Adjuster / Professional ID", KOtomatis), ro1),
					SumberAdjuster), "SetAdjsuter"),
				tampilIsi(ro(medan(CD+"AppointedADJ", "Adjuster / Professional Name", KTeks))),
			),
			bagian("",
				aksi(sumber(roJika(medan(CD+"ConsultantID", "Consultant ID", KOtomatis), ro1), SumberAdjuster),
					"SetConsultant"),
				tampilIsi(ro(medan(CD+"ConsultantName", "Consultant Name", KTeks))),
			),
		),
		bagian("",
			wajibU(medan(CD+"CNPCircumtances", "Circumtances", KArea)),
			medan(CD+"Occupation", "Occupation", KArea),
			wajibU(roJika(medan(CD+"Location", "Location of Loss", KArea), ro1)),
			dua(
				bagian("", wajibU(sumber(medan(CD+"Province", "Province", KOtomatis), SumberProvinsi))),
				bagian("", aksi(medan(CD+"PostalCode", "Zip Code", KTeks), "GetAdders")),
			),
		),
	)
	estimasi := bagian("",
		roJika(medan(CD+"ShareCeding", "Share Ceding(%)", KAngka), ro1),
		blokDeductible(ro1, isCFS, true),
		gridClaimAmount(true),
		sebaris("", ro(medan(TM+"RNMShare", "RNM Share", KAngka)), label("%")),
		gridLossAllocation(true),
		tombolOQ("EditXOLAllocation", "Edit XOL Allocation", OQSandiEditXOL),
		gridXOL(false),
		gridSummaryXOL(false),
	)
	return []Unsur{
		letak(LetakJudul,
			label("Outstanding Claim"),
			tampil(label("Claim No      .........."), sama(CD+"NoClaim", "")),
			tampilIsi(ro(medan(CD+"NoClaim", "Claim No", KTeks))),
		),
		blokTreaty([]Unsur{
			tombol("ChooseMasterIn", "Choose Master In", "PilihMaster:IN"),
			tombolOQ("ViewMaster", "View Master", OQViewMaster),
		}, false),
		// Rupa ikut Claim Prop (perintah work owner 09-10-2026 "ikuti tampilan klaim prop"): Claim Information kartu
		// sendiri di atas tab, tombol layar satu baris aksi tanpa kartu - di XML Claim Information tab pertama layout group
		// (`[penyimpangan sadar]` tata letak, PARITAS §9; isi, kondisi, dan aksi tidak berubah).
		bagian("Claim Information", info),
		sebaris("", ro(medan(TM+"RNMShare", "RNM Share", KAngka)), label("%")),
		letak(LetakTab,
			bagian("Interests", sectionInterest(false)),
			bagian("Estimation", estimasi),
			bagian("Spreading", bagian("Spreading Claim", gridSpreadingKlaim(true)...)),
		),
		tombol("Save", "Save", "SaveDataToJClaim"),
		naJika(tombol("SaveToIssueRNM", "Save to issue RNM", "SaveDataToOSAksep"), atau(ro1, proteksiIssue)),
		naJika(tombol("PrintCFS", "Print CFS", "GenerateCFS"), atau(beda("IsOutstanding", "1"), isCFS)),
		naJika(tombol("PrintPLA", "Print PLA", "BukaPLA"), atau(sama(CD+"IsPLA", "1"), sama(CD+"IsPLA", "1A"))),
		naJika(tombol("Submit", "Submit", "Submit"), atau(beda("IsCFS", "1"), beda("IsOutstanding", "1"))),
		gridRiwayat(false),
	}
}

// ---------------------------------------------------------------- Input Acceptation

// LayarAkseptasi - Section `InputAcceptation` (FlowAction InputAcceptation). Tanpa Submit (OQ-CNP-10); Back (S28 NEVER)
// tidak dibangun (OQ-CNP-09).
func LayarAkseptasi() []Unsur {
	acc := isAcceptation
	ro1 := isOutstanding
	info := bagian("",
		dua(
			ro(medan(CD+"PolicyData.PolicyNo", "Policy No", KTeks)),
			wajibU(aksi(roJika(medan(CD+"DateOfLoss", "Date of Loss", KTanggal), acc), "CheckDateDOL")),
			tampil(tombolOQ("ViewPaymentStatus", "View Payment Status", OQLayananBaca), terisi(CD+"PolicyData.PolicyNo")),
			wajibU(roJika(medan(CD+"ReportDate", "Report Date", KTanggal), ro1)),
			roJika(medan(CD+"PolicyNo", "Policy No Ceding", KTeks), acc),
			wajibJ(roJika(medan(CD+"DateReceived", "Received Date", KTanggal), ro1), pyNoteKosong),
			roJika(medan(CD+"InsuredName", "Insured Name", KTeks), acc),
			wajibJ(roJika(medan(CD+"ReporterName", "Reporter Name", KTeks), ro1), pyNoteKosong),
			roJika(medan(CD+"PlaNoCeding", "Pla No Ceding", KTeks), acc),
			roJika(medan(CD+"CNPClmNoCedant", "Claim No Ceding", KTeks), ro1),
			roJika(medan(CD+"PlaNoSOB", "Pla No SOB", KTeks), acc),
			wajibU(roJika(medan(CD+"ReporterTelp", "Reporter Phone Number", KTelepon), ro1)),
			aksi(roJika(medan(CD+"PolicyData.StartDateTime", "Policy Start", KTanggal), acc), "SetEndDate"),
			aksi(roJika(medan(CD+"PolicyData.EndDateTime", "Policy End", KTanggal), acc), "CheckPeriodPolicy"),
			bagian("",
				ro(medan(CD+"CauseOfLoss", "Cause of Loss", KTampil)),
				naJika(tombol("ChooseCauseOfLoss", "Choose Cause of Loss", "PilihSebab"), acc),
			),
			sumber(naJika(roJika(medan(CD+"ReporterStatus", "Reporter Status", KPilih), isEstimation), ro1),
				kode("ReporterStatus")),
			sumber(roJika(medan(CD+"ReportType", "Report Type", KPilih), ro1), kode("ReportType")),
			blokKatastrofe(),
			roJika(medan(CD+"InsuredRelationshipOthers", "Specify...", KTeks), beda(CD+"InsuredRelationship", "3")),
			aksi(roJika(medan(CD+"Location", "Location of Loss", KArea), acc), "MakeLowercase"),
			wajibJ(roJika(medan(CD+"ReportAddress", "Reporter Address", KArea),
				atau(ro1, beda(CD+"InsuredRelationship", "3"))), pyNoteKosong),
			tampil(aksi(sumber(roJika(medan(CD+"AppointedADJID", "Adjuster / Professional ID", KOtomatis), ro1),
				SumberAdjuster), "SetAdjsuter"), sama(CD+"AppointedADJID", "")),
			roJika(medan(CD+"ReportDescription", "Report Description", KArea), ro1),
			tampil(aksi(sumber(roJika(medan(CD+"ConsultantID", "Consultant ID", KOtomatis), ro1), SumberAdjuster),
				"SetConsultant"), sama(CD+"ConsultantID", "")),
			tampilIsi(ro(medan(CD+"AppointedADJ", "Adjuster / Professional Name", KTeks))),
			tampilIsi(ro(medan(CD+"ConsultantName", "Consultant Name", KTeks))),
		),
		bagian("",
			medan(CD+"CNPCircumtances", "Circumtances", KArea),
			wajibU(medan(CD+"Occupation", "Occupation", KArea)),
			medan(CD+"CNPSupportDoc", "Supporting Document", KArea),
		),
	)
	estimasi := bagian("",
		aksi(roJika(medan(CD+"ShareCeding", "Share Ceding(%)", KAngka), acc), "IntIsSaveToOs"),
		blokDeductible(acc, nil, false),
		sebaris("",
			aksi(medan("FlagActualPremium", "Waiting For Actual Premium", KCentang), "SetActualPremium"),
			tampil(tombol("ViewActualPremium", "View", "LihatSelisihAktual"), sama("FlagActualPremium", "true")),
		),
		tombol("ViewOldAllocation", "View Old Allocation", "LihatAlokasiLama"),
		gridClaimAmount(false),
		sebaris("", ro(medan(TM+"RNMShare", "RNM Share", KAngka)), label("%")),
		gridLossAllocation(false),
		gridXOL(true),
		gridSummaryXOL(true),
		bagian("Spreading Claim Out", gridSpreadingKlaim(false)...),
		tampil(bagian("",
			tampilIsi(ro(medan("KomiteNo", "KomiteNo", KTeks))),
			gridKomite(DaftarKomiteKlaim),
		), atau(sama("IsCloseFile", "1"), sama("IsReject", "1"))),
		naJika(tombol("SaveToOS", "Save To OS", "SaveToOS"), sama("IsSaveToOs", "1")),
	)
	return []Unsur{
		letak(LetakJudul,
			label("Acceptation Claim"),
			tampilIsi(ro(medan(CD+"NoClaim", "Claim No", KTeks))),
		),
		blokTreaty([]Unsur{
			tombolOQ("ViewMaster", "View Master", OQViewMaster),
			tombol("CloseClaim", "Close Claim", "BukaTutupKlaim"),
			naJika(tombol("CloseWithoutPayment", "Close Without Payment (CWP)", "BukaCWP"), acc),
			naJika(tampil(tombol("ChooseMasterIn", "Choose Master", "PilihMaster:IN"), beda("IsAcceptation", "1")), acc),
			naJika(tampil(tombol("ChooseMasterInEDM", "Choose Master", "PilihMaster:INEDM"), acc), beda("IsAcceptation", "1")),
			tampil(tombol("ViewPaymentAttachment", "View Payment Attachment", "LihatLampiranBayar"), acc),
			tombol("ClaimHistoryMasterID", "Claim History Master ID", "LihatRiwayatMaster"),
		}, true),
		bagian("Claim Information", info), // rupa ikut Claim Prop (lihat LayarOutstanding)
		letak(LetakTab,
			bagian("Interests", sectionInterest(true)),
			bagian("Estimation", estimasi),
			bagian("Acceptation", sectionAcceptation()),
		),
		gridRiwayat(true),
	}
}

// DaftarKomiteKlaim - `pyWorkPage.ClaimData.ClaimComitee` (grid "Committe Accept Status" InputAcceptation, tampil bila
// kasus ditutup / ditolak komite): tangga akseptasi terakhir yang diserahkan (medan turunan).
const DaftarKomiteKlaim = CD + "ClaimComitee"

// sectionAcceptation - tab Acceptation (InputAcceptation S14): Acceptation List (masterDetail -> AdjustmentDetailNP_Section)
// dan "Spreading Claim Out" (Σ Spreading In / Out seluruh akseptasi, hanya-baca).
func sectionAcceptation() Unsur {
	spread := func(jalur string, pilih bool) Unsur {
		tt := kRO(kol("TreatyName", "Treaty Type", KTeks))
		if pilih {
			tt = kSumber(kRO(kol("TreatyName", "Treaty Type", KPilih)), SumberTreaty)
		}
		return Unsur{Jenis: JenisGrid, Jalur: jalur, Bernomor: true, Kolom: []Unsur{
			tt,
			kRO(kol("SharePercentage", "Share(%)", KAngka)),
			kRO(kol("Currency", "Currency", KTampil)),
			kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka)),
			kRO(kol("AdjusterFee", "Adjuster Fee", KAngka)),
			kRO(kol("Salvage", "Salvage", KAngka)),
			kRO(kol("CNPOthersFee", "Fee", KAngka)),
			kRO(kol("TotalClaim", "Total Claim", KAngka)),
			kRO(kol("PremiumSpreaded", "Reinstatement Premium", KAngka)),
		}, Kaki: []Unsur{ro(medan(CD+"TotalEstimasi", "", KAngka))}}
	}
	return bagian("",
		bagian("Acceptation List", Unsur{Jenis: JenisGrid, Jalur: DaftarAdjustment, Bernomor: true,
			Tambah: ptr(ikon(naJika(tombol("AddAkseptasiCNP", "Add", "AddAkseptasiCNP"), sama("IsSaveToOs", "0")), IkonTambah)),
			Kolom: []Unsur{
				kSumber(kol("Type", "Type", KPilih), kode("AdjustmentType")),
				kRO(kol("AcceptedNo", "Acceptation No", KTampil)),
				kRO(kol("AcceptedDate", "Acceptation Date", KTanggal)),
				kSumber(kRO(kol("AcceptanceStatus", "Status", KPilih)), kode("AcceptanceStatus")),
				ikonK(kTombol("DeleteAkseptasi", "Delete", "DeleteAkseptasi", func(_ *Halaman, b Baris) bool {
					return b["AcceptanceStatus"] != ""
				}), IkonHapus),
			}}),
		bagian("Spreading Claim Out",
			bagian("Spreading List", spread(DaftarSpreadAdj, true)),
			bagian("Spreading List", spread(DaftarSpreadAdjQS, false)),
		),
	)
}
