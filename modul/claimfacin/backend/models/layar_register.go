package models

// Untuk apa berkas ini: DEFINISI LAYAR INPUT REGISTER - Section `InputRegister` (FlowAction InputRegister) beserta
// `InputRegisterDetail`, `ViewHistoryClaim` (+ `InputInwardFacultativeDtl`), `ProgresClaim_Sec`, `Catastrope_Sec`, dan
// pop-up `ViewPolis` (harness ChoosePolis), `ProtectDOL`. Label, urutan, dan kondisi VERBATIM dari XML
// `Claim Fac In/Section`.
//
// Aturan baca (pola Claim Prop / Claim Non Prop):
//   - `Aksi` = nama activity tanpa akhiran `_Act` / `_act` / `_ACT` (services/aksi.go); "" = postValue saja.
//   - Tombol yang activity / harness-nya TIDAK diekspor atau memanggil layanan luar yang belum disetujui tampil sesuai
//     section tetapi NONAKTIF dengan `Catatan` OQ.
//   - Sel ber-visible `NEVER` / `1=2` / `1==2` tidak dibangun.

// Catatan OQ (docs/OQ.md).
const (
	OQBayarPremi    = "OQ-CFI-14: View Status Payment Premi memanggil layanan REST luar (M_LINK_SERVICE) yang belum disetujui"
	OQProgresManual = "OQ-CFI-15: FlowAction InputSubProgressClaim tanpa post-activity dan sumber Progres1 / Progres2 di ekspor"
	OQDokumenPDF    = "OQ-CFI-20: aliran HTML dokumen (CFS / PLA / DLA) tidak diekspor - data dan nomor tersimpan, berkas belum dibuat"
	OQPLAFac        = "OQ-CFI-21: ObjectItem.GeneratePLA (PLA fac retro) tidak diekspor - nomor PLA fac belum dapat diterbitkan"
)

// Kunci daftar pilihan (`services` mengisinya).
const (
	SumberMataUang  = "mataUang"  // BrowseCurrency_RD (ID, nama)
	SumberAdjuster  = "adjuster"  // BrowseAdjusterConsultant (ID, NAME)
	SumberNegara    = "negara"    // BrowseCountry_RD
	SumberProvinsi  = "provinsi"  // BrowseProvince_RD (RW.PROVINCENAME menurut negara)
	SumberKota      = "kota"      // BrowseCity_RD (RW.CITYNAME menurut provinsi)
	SumberDistrik   = "distrik"   // BrowseDistrictInputC_RD (DISTRICT menurut CityID)
	SumberRW        = "rw"        // BrowseRWInput_RD (RW menurut DistrictID)
	SumberJenisReas = "jenisReas" // BrowseReinsuranceType_RD
	SumberOkupasi   = "okupasi"   // D_OccupationList (halaman polis)
	SumberItemProp  = "itemProp"  // D_FilteredPropertyItemList
	SumberCovFire   = "covFire"   // D_FilteredCoverageList
	SumberAneka     = "aneka"     // D_AnekaList
	SumberCovAneka  = "covAneka"  // D_FilteredCoverageAnekaList
	SumberCovObjek  = "covObjek"  // D_CoverageMBUClaimList / D_CoveragePAClaimList / D_CoverageTravelClaimList
	AwalanKode      = "kode:"     // prompt values `associated` - tidak terbaca di XML; kode DB apa adanya
)

// KodePilihan - kode bersumber `associated` (prompt values tidak diekspor, OQ-CFI-18): kode yang disebut kondisi section /
// activity XML. Ditampilkan apa adanya kecuali ada label di `LabelKode`.
var KodePilihan = map[string][]string{
	"ReportType":          {"1", "2", "3", "4", "5"},
	"InsuredRelationship": {"1", "2", "3"}, // GetReportStatus_Act: 1 ceding, 2 SOB, 3 lain (Specify...)
	"FormType":            {"1", "2"},      // Estimasi LS5 `.FormType==1`, LS11 `.FormType==2`
	"TypeDeductible":      {"1", "2"},      // Estimasi LS15 `.TypeDeductible = 2` / `= 1`
	"StsKatastrofe":       {"Catastrophe", "Non-Catastrophe"},
	"NonKatastrofeType":   {"Claim", "Big Claim"},
	"SearchType":          {CariNoPolis, CariCeding, CariInsured, CariQQ}, // SearchPolis_act 1-4
	"ShareRetro":          {"1", "2"},                                     // Pla_Dtl `.ShareRetro` DEF=1
}

// LabelKode - label tampilan kode. ReportType: kelas properti sama dengan Claim Prop (ASM-FW-GCNMFW-Data-ClaimData) yang
// labelnya diberikan work owner (screenshot 08-10-2026) - `[inferensi]` satu property rule = satu prompt values.
// SearchType: kolom yang dicari SearchPolis_act langkah 1-4 dengan judul kolom grid ViewPolis.
var LabelKode = map[string]map[string]string{
	"ReportType": {"1": "Direct", "2": "Via Email", "3": "Via Fax", "4": "via Postal Mail/Courier", "5": "Via Telephone"},
	"SearchType": {CariNoPolis: "Policy Number", CariCeding: "Ceding Co Name", CariInsured: "Name of Insured",
		CariQQ: "QQ Name"},
}

func kode(p string) string { return AwalanKode + p }

func ptr(u Unsur) *Unsur { return &u }

// kondisi registrasi
var (
	bukanEstimasi   = beda(JalurIsEstimation, "1")
	kunciAkseptasi  = adaAkseptasi
	adaPolis        = terisi(JalurNoPolis)
	tanpaPolis      = sama(JalurNoPolis, "")
	negaraIndonesia = sama(CD+"Country", "INDONESIA")
)

// errorDOL / errorNol - IsError Register (CheckDoubleClaim_Act): 1 = klaim kembar, 0 atau kosong = lolos.
// `[inferensi]` Pega membandingkan properti Integer kosong dengan 0 sebagai sama; tanpa itu kasus yang Claim Estimate-nya
// tidak pernah diubah tidak punya tombol Submit.
func errorDOL(h *Halaman) bool { return h.Ambil(JalurIsError) == "1" }
func errorNol(h *Halaman) bool { return h.Ambil(JalurIsError) == "0" || h.Ambil(JalurIsError) == "" }

// errorLebih1 - `pyWorkPage.IsError > 1`.
func errorLebih1(h *Halaman) bool {
	ok, _ := lebihDari(h.Ambil(JalurIsError), "1")
	return ok
}

// LayarRegister - Section InputRegister.
func LayarRegister() []Unsur {
	return []Unsur{
		sebaris("",
			tombolOQ("ViewStatusPaymentPremi", "View Status Payment Premi", OQBayarPremi),
			ikon(tombol("RejectClaim", "Reject Claim", "BukaRejectClaim"), IkonHapus),
		),
		letak(LetakTab,
			bagian("Register", detailRegister()...),
			bagian("Policy Detail & Claims History", tabRiwayat()...),
			bagian("Progress Claim", tabProgres()...),
		),
		sebaris("",
			tampil(tombol("Simpan", "Save", "Simpan"), tanpaPolis),
			tampil(tombol("SubmitDOL", "Submit", "BukaProtectDOL"), dan(errorDOL, adaPolis)),
			tampil(tombol("Submit", "Submit", "Submit"), dan(errorNol, adaPolis)),
		),
	}
}

// detailRegister - Section InputRegisterDetail (juga tab "View Registration" layar Input Estimasi, `ro` = seluruhnya
// hanya-baca bila layar itu dikunci).
func detailRegister() []Unsur {
	roPolis := terisi(JalurNoPolis)
	dateRO := kunciAkseptasi
	return []Unsur{
		dua(
			bagian("",
				sebaris("",
					ro(medan(JalurNoPolis, "Policy No", KTampil)),
					naJika(tombol("BukaPolis", "Choose Polis", "BukaPolis"), isEstimasi),
				),
				wajibU(roJika(medan(OQ+"InsuredName", "Insured Name", KTeks), roPolis)),
				// change SetQQName (`.Quotation.QQName := .OfferFacIn.QuotationData.QQName`): Quotation = halaman yang
				// sama di sini, nilai ketikan langsung tersimpan (QQ_NAME) - tanpa aksi.
				roJika(medan(OQ+"QQName", "QQ Name", KTeks), terisi(OQ+"QQName")),
				wajibU(roJika(medan(OQ+"CedingCoName", "Ceding Co Name", KTeks), roPolis)),
				wajibU(roJika(medan(OQ+"SobName", "Source Of Business", KTeks), roPolis)),
				aksi(wajibU(roJika(medan(CD+"DateOfLoss", "Date of Loss", KTanggal), dateRO)), "CheckDate"),
				aksi(wajibU(roJika(medan(CD+"ReportDate", "Report Date", KTanggal), dateRO)), "CheckDateReport"),
				aksi(wajibU(roJika(medan(CD+"DateReceived", "Received Date", KTanggal), dateRO)), "CheckDateReceived"),
				wajibU(roJika(medan(CD+"ReporterName", "Reporter Name", KTeks), dateRO)),
				wajibU(roJika(medan(CD+"ReporterTelp", "Reporter Phone Number", KTelepon), dateRO)),
				sebaris("",
					naJika(tombol("ChooseCauseOfLoss", "Choose Cause of Loss", "BukaSebab"), kunciAkseptasi),
					ro(medan(CD+"CauseOfLoss", "Cause of Loss", KTampil)),
				),
				blokKatastrofe(),
			),
			bagian("",
				roJika(medan(CD+"PlaNoCeding", "PLA No Ceding", KTeks), dateRO),
				sumber(roJika(medan(CD+"ReportType", "Report Type", KPilih), dateRO), kode("ReportType")),
				roJika(medan(CD+"PlaNoSOB", "PLA No SOB", KTeks), dateRO),
				sebaris("",
					aksi(sumber(roJika(medan(CD+"InsuredRelationship", "Reporter Status", KPilih), isEstimasi),
						kode("InsuredRelationship")), "GetReportStatus"),
					roJika(medan(CD+"InsuredRelationshipOthers", "Specify...", KTeks),
						beda(CD+"InsuredRelationship", "3")),
				),
				tampil(gridCedingCo(), dan(bukanEstimasi, sama(CD+"InsuredRelationship", "1"),
					func(h *Halaman) bool { return AmbilJalur(h, JalurBaris(DaftarCedingCo, 1)+".CedingCoName") != "" })),
			),
		),
		bagian("",
			roJika(medan(CD+"ReportDescription", "Description Report", KArea), dateRO),
			wajibU(roJika(medan(CD+"ReportAddress", "Reporter Address", KArea),
				atau(dateRO, beda(CD+"InsuredRelationship", "3")))),
		),
		tampil(bagian("", gridCalon()...), bukanEstimasi),
		dua(
			bagian("",
				tampil(sebaris("",
					aksi(sumber(withTampilan(medan(CD+"ConsultantID", "Consultant ID", KOtomatis), CD+"ConsultantName"),
						SumberAdjuster), "SetConsultant"),
					ikon(tombolOQ("TambahKonsultan", "", OQTambahAdjuster), IkonTambah),
				), beda(JalurIsAnyAccept, "1")),
				tampilIsi(ro(medan(CD+"ConsultantName", "Consultant Name", KTeks))),
			),
			bagian("",
				tampil(sebaris("",
					aksi(sumber(withTampilan(medan(CD+"AppointedADJID", "Adjuster / Professional ID", KOtomatis),
						CD+"AppointedADJ"), SumberAdjuster), "SetAdjsuter"),
					ikon(tombolOQ("TambahAdjuster", "", OQTambahAdjuster), IkonTambah),
				), beda(JalurIsAnyAccept, "1")),
				tampilIsi(ro(medan(CD+"AppointedADJ", "Adjuster / Professional Name", KTeks))),
			),
		),
		aksi(wajibU(naJika(medan(CD+"Location", "Location of Loss", KArea), kunciAkseptasi)), "SetTempLocation"),
		bagian("",
			aksi(sumber(roJika(medan(CD+"Country", "Country", KOtomatis), isEstimasi), SumberNegara), "PilihNegara"),
			tampil(aksi(sumber(roJika(medan(CD+"Province", "Province", KOtomatis), isEstimasi), SumberProvinsi), "PilihProvinsi"), negaraIndonesia),
			tampil(aksi(sumber(roJika(medan(CD+"City", "City", KOtomatis), isEstimasi), SumberKota), "PilihKota"), negaraIndonesia),
			tampil(aksi(sumber(roJika(medan(CD+"District", "District", KOtomatis), isEstimasi), SumberDistrik), "PilihDistrik"), negaraIndonesia),
			tampil(aksi(sumber(roJika(medan(CD+"RW", "Region", KOtomatis), isEstimasi), SumberRW), "PilihRW"), negaraIndonesia),
			tampil(aksi(roJika(medan(CD+"PostalCode", "Zip Code", KTeks), isEstimasi), "InputKodePos1"), negaraIndonesia),
		),
		dua(
			aksi(sumber(wajibU(roJika(medan(CD+"Currency", "Currency", KPilih), isEstimasi)), SumberMataUang),
				"InputCurrencyValueAct_Register"),
			bagian("",
				roJika(medan(CD+"GrossEstimate", "Gross Estimate (100%)", KAngka), dateRO),
				aksi(roJika(medan(CD+"ClaimEstimate", "Claim Estimate", KAngka), dateRO), "CheckDoubleClaim"),
			),
		),
		tampil(ro(medan(CD+"DollarCurrencyVal", "Value in IDR", KAngka)), sama(CD+"PostalCode", "1")),
	}
}

// OQTambahAdjuster - ikon tambah Consultant / Adjuster (harness MstAdjusterConsultant, kelas Data-Portal). Claim Prop
// menyimpan master barunya lewat rute pinjaman `POST /api/adjuster-consultant`, yang di `cmd/api/rakit.go`
// (`ruteDipinjam`) hanya dipinjamkan ke claimprop - tampil sesuai section, nonaktif sampai rute itu dipinjamkan juga.
const OQTambahAdjuster = "OQ-CFI-32: rute pinjaman POST /api/adjuster-consultant belum dipinjamkan ke Claim Fac In " +
	"(cmd/api/rakit.go ruteDipinjam) - tambahkan Adjuster / Consultant dari menu Adjuster Consultant"

func withTampilan(u Unsur, j string) Unsur { u.Tampilan = j; return u }

// gridCedingCo - LS132 (CedingCoList polis, tombol Choose = GetCeding_act Param.Ceding).
func gridCedingCo() Unsur {
	return Unsur{Jenis: JenisGrid, Jalur: DaftarCedingCo, Kolom: []Unsur{
		kRO(kol("CedingCoName", "Ceding Co Name", KTeks)),
		kTombol("GetCeding", "Choose", "GetCeding", nil),
	}}
}

// gridCalon - LS137: grid pilih objek per lini (`TempClaimData.ClaimData.ObjectList`, CheckListEstimasi_Act).
func gridCalon() []Unsur {
	pilih := kAksi(kol(PropDipilih, "", KCentang), "CheckListEstimasi")
	g := func(k Kondisi, kolom ...Unsur) Unsur {
		return tampil(Unsur{Jenis: JenisGrid, Jalur: DaftarCalon, Bernomor: true, PerHalaman: 10,
			Kolom: append([]Unsur{pilih}, kolom...)}, k)
	}
	return []Unsur{
		g(atau(IsGolfInsurance, IsAneka), kRO(kol("ObjectName", "Object", KTampil)), kRO(kol("ObjectLocation", "Location", KTampil))),
		g(IsFire, kRO(kol("ObjectName", "Object", KTampil)), kRO(kol("ObjectLocation", "Location", KTampil))),
		g(IsMarineCargo, kRO(kol("ObjectName", "Trading", KTampil)), kRO(kol("ObjectJob", "Goods", KTampil)),
			kRO(kol("ObjectSurveyor", "Packing", KTampil)), kRO(kol("ObjectSurveyLocation", "Conveyance", KTampil))),
		g(IsMBU, kRO(kol("BrandName", "Brand Name", KTampil)), kRO(kol("ModelName", "Model", KTampil)),
			kRO(kol("TypeName", "Type", KTampil)), kRO(kol("LicensePlate", "License Plate", KTampil)),
			kRO(kol("ChassisNumber", "Chassis Number", KTampil)), kRO(kol("EngineNumber", "Engine Number", KTampil))),
		g(IsTravel, kRO(kol("ObjectName", "Participant Name", KTampil)),
			kRO(kol("ObjectParticipantStatus", "Status", KTampil)), kRO(kol("ObjectIDCard", "ID/Pasport", KTampil)),
			kRO(kol("ObjectDateOfBirth", "Date of Birth", KTanggal))),
		// PA: `.ObjectDateOfBirth` tanpa RO di XML - salinan polis, tidak disimpan; ditampilkan hanya-baca.
		g(IsPA, kRO(kol("ObjectName", "Name", KTampil)), kRO(kol("ObjectJob", "Job", KTampil)),
			kRO(kol("ObjectDateOfBirth", "Date of Birth", KTanggal)), kRO(kol("Gender", "Gender", KTampil))),
	}
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

// tabRiwayat - Section ViewHistoryClaim: Policy Detail (InputInwardFacultativeDtl, IsTreatyIn 0 / kosong), Claim Status
// (kronologi, urut menurun), Claim History (klaim lain polis ini).
func tabRiwayat() []Unsur {
	return []Unsur{
		letak(LetakTab,
			bagian("Policy Detail", tampil(detailPolis(), atau(sama(JalurIsTreatyIn, "0"), sama(JalurIsTreatyIn, "")))),
			bagian("Claim Status", Unsur{Jenis: JenisGrid, Jalur: DaftarKronologiTampil, Bernomor: true, Kolom: []Unsur{
				kRO(kol("ASMDateTimeChronology", "Date Input", KTanggal)),
				kRO(kol("ASMUserID", "User", KTampil)),
				kRO(kol("ASMNoteType", "Status", KTampil)),
				kRO(kol("pyNote", "Note", KTampil)),
			}}),
			bagian("Claim History", Unsur{Jenis: JenisGrid, Jalur: DaftarRiwayatPolis, Kolom: []Unsur{
				kRO(kol("ClaimNo", "CLAIM NO", KTampil)),
				kRO(kol("DateOfLoss", "DATE OF LOSS", KWaktu)),
			}}),
		),
	}
}

// DaftarKronologiTampil - salinan Chronology untuk grid "Claim Status" (urut menurun, `TampilKronologi`).
const DaftarKronologiTampil = "ClaimStatus"

// detailPolis - Section InputInwardFacultativeDtl (layout "General").
func detailPolis() Unsur {
	return bagian("General",
		ro(medan(OQ+"NoOfferSlip", "Offer  Slip Number", KTeks)),
		ro(medan(JalurNoPolis, "Policy No", KTeks)),
		ro(medan(OQ+"BusinessName", "Class Of Business", KTampil)),
		ro(medan(OQ+"StatusBusiness", "Business Status", KTampil)),
		ro(medan(OQ+"SobName", "Source Of Business", KTampil)),
		ro(medan(OQ+"InsuredName", "Insured Name", KTeks)),
		ro(medan(OQ+"CedingCoName", "Ceding Co Name", KTampil)),
		ro(medan(OQ+"QQName", "QQ Name", KTampil)),
		ro(medan(JalurMulaiPolis, "Begin date", KTanggal)),
		ro(medan(JalurAkhirPolis, "End Date", KTanggal)),
		ro(medan(AwalanPolis+".PolicyData.OfferingDate", "Offering date", KTanggal)),
		tampilIsi(ro(medan(OQ+"Period", "Period", KTeks))),
		ro(medan(AwalanPolis+".IsProRate", "Prorate / Short period", KTampil)),
		ro(medan(OQ+"MarketingName", "Marketing Name", KTeks)),
		tampil(tombol("ViewRetroList", "View Retro List", "BukaRetroList"),
			func(h *Halaman) bool { return AmbilJalur(h, JalurBaris(DaftarFacRetro, 1)+".ReinsurerName") != "" }),
	)
}

// tabProgres - Section ProgresClaim_Sec (GetProgresClaim_ACT: PROGRESSCLAIM kasus + SUBPROGRESSCLAIM per posisi).
func tabProgres() []Unsur {
	return []Unsur{
		Unsur{Jenis: JenisGrid, Jalur: DaftarProgres, Rincian: PanelProgres, Kolom: []Unsur{
			kRO(kol("CaseID", "Case ID", KTampil)),
			kRO(kol("ProdKe", "Position", KTampil)),
			kRO(kol("AnalystTransferDate", "Date", KWaktu)),
			kRO(kol("StatusClaim", "Status", KTampil)),
		}},
	}
}

// Progres klaim (halaman baca, tidak disimpan di halaman kasus).
const (
	DaftarProgres = "ProgresClaim"
	AnakSubProg   = "ObjectList"
	PanelProgres  = "progres"
)

// LayarSubProgres - Section SubProgresClaim_Sec (baris ke-n grid Progress Claim).
func LayarSubProgres(n int) []Unsur {
	return []Unsur{
		tombolOQ("InputProgresClaim", "Input Progres Claim", OQProgresManual),
		Unsur{Jenis: JenisGrid, Jalur: JalurAnak(DaftarProgres, n, AnakSubProg), Kolom: []Unsur{
			kRO(kol("ObjekTanggal", "Input Date", KWaktu)),
			kRO(kol("Notes", "Case ID", KTampil)),
			kRO(kol("RemarksPLA", "Status Progress 1", KTampil)),
			kRO(kol("RemarksDLA", "Status Progress 2", KTampil)),
			kRO(kol("TypeName", "User Input", KTampil)),
			kRO(kol("TanggalFolloup", "Next Follow Up Date", KWaktu)),
			kRO(kol("Comment", "Note", KTampil)),
		}},
	}
}

// LayarPilihPolis - Section ViewPolis (harness ChoosePolis): pencarian (SearchType / SearchName), grid hasil diisi layar
// dari pilihan server (`SearchPolis_act`), tautan nomor polis = `CopyNB_Act`, detail polis bila sudah dipilih
// (`InputParam.CARI4 = 'true'`), Cancel / Submit (`SetInputParam_Act`).
func LayarPilihPolis() []Unsur {
	return []Unsur{
		sebaris("",
			sumber(medan(JalurJenisCari, "Search Type", KPilih), kode("SearchType")),
			medan(JalurTeksCari, "", KTeks),
			tombol("SearchPolis", "Search", "SearchPolis"),
		),
		tampil(detailPolis(), dan(sama(JalurDetailPolis, "true"), atau(sama(JalurIsTreatyIn, "0"),
			sama(JalurIsTreatyIn, "")))),
		sebaris("",
			tombol("BatalPolis", "Cancel", ""),
			tombol("SetInputParam", "Submit", "SetInputParam"),
		),
	}
}

// LayarProtectDOL - Section ProtectDOL (local action ProtectDOL): kasus kembar (TempHistoryClaim .BRANCH_NAME = nomor
// klaim) lalu Cancel / Submit (`InsertObjectItemList_DT` + `CheckListEstimasi_Act` + finishAssignment).
func LayarProtectDOL() []Unsur {
	return []Unsur{
		bagian("Sudah ada data yang sama dengan ID Pega berikut", Unsur{Jenis: JenisGrid, Jalur: DaftarKembar,
			Kolom: []Unsur{kRO(kol("BRANCH_NAME", "", KTampil))}}),
		bagian("Lanjutkan?", sebaris("",
			tombol("BatalDOL", "Cancel", ""),
			tombol("SubmitDOL2", "Submit", "Submit"),
		)),
	}
}

// DaftarKembar - `TempHistoryClaim.pxResults` (CheckDoubleClaim_Act): diisi ulang setiap layar dibuka.
const DaftarKembar = "TempHistoryClaim"
