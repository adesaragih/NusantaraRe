package models

// Untuk apa berkas ini: DEFINISI LAYAR CHOOSE SURVEYOR / INPUT ADJUSTMENT - Section `ClaimSurvey` (FlowAction
// InputSurveyor) dengan tab Claim Details (`InputEstimasi`: Registration / Estimation / Adjustment & Acceptation),
// grid objek `ShowObjectAdj` (masterDetail), panel objek `ItemListEstimation*` (grid item), panel item `Adjusment_SC`
// (grid coverage Marine + grid Adjustment masterDetail), panel adjustment `InputAdjustment`, pop-up `PrintDLA_dtl`
// (dla.go) dan `ViewCedantPanel`. Label, urutan, kondisi VERBATIM dari XML.
//
//	PanelObjekAdj  "adj"      baris ClaimData.ObjectList(o)                                   -> LayarObjekAdj(o)
//	PanelItemAdj   "adjitem"  baris ClaimData.ObjectList(o).ObjectItemList(i)                 -> LayarItemAdj(o, i)
//	PanelAdj       "adjdtl"   baris ClaimData.ObjectList(o).ObjectItemList(i).Adjustment(a)   -> LayarAdjustment(o, i, a)
//	ModalCedant    "cedant"   kunci sama dengan PanelAdj                                      -> LayarCedant(o, i, a)
//
// Tidak dibangun (PARITAS): footer `CountSpread.CARI24-27` (halaman requestor tanpa penulis); grid `.KomiteList`
// ItemListEstimation (VIS 1=2); tombol "Send Claim to committee" grid item PA / Travel (harness ClaimCommittee tidak
// diekspor, OQ-CFI-25); "View Retro" spreading adjustment (local action ShowRetro kelas SpreadingRisk - section
// pemilihnya ShowRetro_Sec membaca halaman NB, OQ-CFI-25).

// Prefiks panel / modal layar Input Adjustment.
const (
	PanelObjekAdj = "adj"
	PanelItemAdj  = "adjitem"
	PanelAdj      = "adjdtl"
	ModalCedant   = "cedant"
)

// Sumber pilihan khas adjustment.
const (
	SumberMataUangAdj = "mataUangAdj" // `.CurencyAdjustment` baris adjustment (Choose Currency)
	SumberRekening    = "rekening"    // `Result.pxResults` SetPayable_Act (Name of Bank)
)

// OQ layar adjustment.
const (
	OQLayananBayar = "OQ-CFI-14: View Status Payment Claim / View Payment Attachment memanggil layanan REST luar " +
		"(M_LINK_SERVICE) yang belum disetujui"
	OQKomitePA = "OQ-CFI-25: harness ClaimCommittee dan local action ShowRetro tidak diekspor"
)

func init() {
	KodePilihan["PaymentType"] = []string{BayarFinal, BayarInterim, BayarSalvage, BayarFee, BayarAdjust, BayarExpense, "7"}
	KodePilihan["IndividualRiskType"] = []string{"1", "2", "3"}
	KodePilihan["Payable"] = []string{"1", "2", "3"}
	KodePilihan["AcceptanceStatus"] = []string{"1", "2"}
	// AcceptanceStatus: label kolom "Status" (work owner 11-10-2026 "STATUS NYA BENERIN JANGAN ANGKA") = label History
	// Adjustment layar Komite Claim Fac In (`komiteclaimfacin` LabelStatusAdjustment) - nilai yang sama, satu kosakata.
	LabelKode["AcceptanceStatus"] = map[string]string{"0": "Transfer to committee", "1": "Approve", "2": "Reject"}
	KodePilihan["ExGratia"] = []string{"0", "1"}
	// Kelas sama dengan Claim Prop (ASM-FW-GCNMFW-Data-Adjustment.IndividualRiskType, -Data-ClaimData.Payable) yang
	// labelnya diberikan work owner (08-10-2026) - `[inferensi]` satu property rule = satu prompt values. PaymentType:
	// prompt values property rule ASM-FW-GCNMFW-Data-Adjustment.PaymentType (screenshot work owner 11-10-2026). ExGratia:
	// "NO" dari layout lama (screenshot 11-10-2026, nilai bawaan 0); "YES" `[inferensi]` (OQ-CFI-37). Status komite:
	// komite.go.
	LabelKode["PaymentType"] = map[string]string{BayarFinal: "Final", BayarInterim: "Claim", BayarSalvage: "Salvage",
		BayarFee: "Adjuster Fee", BayarAdjust: "Adjustment", BayarExpense: "Consultant Fee", "7": "Cancellation"}
	LabelKode["ExGratia"] = map[string]string{"0": "NO", "1": "YES"}
	LabelKode["IndividualRiskType"] = map[string]string{"0": "Select..", "1": "% From claims", "2": "% FromTSI",
		"3": "Other"}
	LabelKode["Payable"] = map[string]string{"1": "Ceding Co Name", "2": "Broker Name", "3": "Others"}
}

// LayarSurveyor - Section ClaimSurvey (FlowAction InputSurveyor, label "Input Adjustment").
func LayarSurveyor() []Unsur {
	return []Unsur{
		sebaris("",
			tampilIsi(ro(medan(CD+"NoClaim", "Claim No", KTeks))),
			tombolOQ("ViewStatusPaymentPremi", "View Status Payment Premi", OQBayarPremi),
		),
		tampilIsi(ro(medan(CD+"Remark", "Remark from Inputor", KArea))),
		letak(LetakTab,
			bagian("Claim Details", letak(LetakTab,
				bagian("Registration", detailRegister()...),
				bagian("Estimation", tampil(gridObjekEstimasi(), bukanTreatyIn)),
				bagian("Adjustment & Acceptation", gridObjekAdj()),
			)),
			bagian("Policy Details and Claim History", tabRiwayat()...),
			tampil(bagian("Progress Claim", tabProgres()...), TabProgresTampil),
		),
		sebaris("",
			// Back disembunyikan begitu ada adjustment diakseptasi (work owner 11-10-2026 "TOMBOL BACK KALO UDAH ADA
			// AKSEPTASI, HIDE AJA"): kembali ke Input Estimasi sesudah akseptasi tidak lagi ditawarkan; server ikut
			// menolak aksinya (tombol tidak terbuka).
			tampil(tombol("Back", "Back", "BackToEstimasi"), func(h *Halaman) bool { return !AdaAkseptasi(h) }),
			tombol("CloseClaim", "Close Claim", "PreventRejectClaim"),
		),
	}
}

// gridObjekAdj - Section ShowObjectAdj: grid ObjectList per lini (masterDetail ke panel objek) + "Print DLA".
func gridObjekAdj() Unsur {
	dla := kTampil(kTombol("PrintDLA", "Print DLA", "BukaDLA",
		func(_ *Halaman, b Baris) bool { return b["DLAStatus"] != "0" && b["DLAStatus"] != "" }),
		bSama("IsFacretro", "1"))
	g := func(k Kondisi, kolom ...Unsur) Unsur {
		return tampil(Unsur{Jenis: JenisGrid, Jalur: DaftarObjek, Bernomor: true, Rincian: PanelObjekAdj,
			Kolom: append(kolom, dla)}, k)
	}
	return bagian("",
		g(atau(IsGolfInsurance, IsAneka, IsFire), kRO(kol("ObjectName", "Object Name", KTampil)),
			kRO(kol("ObjectLocation", "Location", KTampil))),
		g(IsMarineCargo, kRO(kol("ObjectName", "Trading Name", KTampil)), kRO(kol("ObjectLocation", "Packing Name", KTampil)),
			kRO(kol("ObjectJob", "Goods Name", KTampil)), kRO(kol("ObjectSurveyLocation", "Conveyance Name", KTampil))),
		g(IsPA, kRO(kol("ObjectName", "Object Name", KTampil)), kRO(kol("ObjectJob", "Job", KTampil)),
			kRO(kol("ObjectDateOfBirth", "Date of Birth", KTanggal)), kRO(kol("Gender", "Gender", KTampil))),
		g(IsTravel, kRO(kol("ObjectName", "Object Name", KTampil)),
			kSumber(kRO(kol("ObjectParticipantStatus", "Status", KPilih)), kode("ObjectParticipantStatus")),
			kRO(kol("ObjectIDCard", "IDCard", KTampil)), kRO(kol("ObjectDateOfBirth", "Date of Birth", KWaktu))),
		g(IsMBU, kRO(kol("BrandName", "Brand Name", KTampil)), kRO(kol("ModelName", "Model", KTampil)),
			kRO(kol("TypeName", "Type", KTampil))),
	)
}

// LayarObjekAdj - panel baris objek o: ItemListEstimation (Fire / Aneka / Golf / Marine), ItemListEstimationMBU,
// ItemListEstimationPA, ItemListEstimationTravel (EditAction grid objek ShowObjectAdj per lini).
func LayarObjekAdj(o int) []Unsur {
	daftar := DaftarItem(o)
	angka := func(p, j string) Unsur { return kRO(kol(p, j, KAngka)) }
	g := func(k Kondisi, kolom ...Unsur) Unsur {
		return tampil(Unsur{Jenis: JenisGrid, Jalur: daftar, Rincian: PanelItemAdj, Kolom: kolom}, k)
	}
	nilai := func(idr string) []Unsur {
		return []Unsur{kRO(kol("Currency", "Currency", KTampil)), angka("KursObjectItem", "Currency Value(IDR)"),
			angka("TSIPerObject", "TSI Object"), angka("TSINusare", "TSI RNM"), angka("ValueTSINusareIDR", idr)}
	}
	komitePA := func(id, lbl string) Unsur { return kTombol(id, lbl, "", func(*Halaman, Baris) bool { return true }) }
	return []Unsur{
		tampil(bagian("",
			g(IsGolfInsurance, append([]Unsur{kRO(kol(PropNomorBaris, "Object No", KTampil)),
				kRO(kol("ObjectItemName", "Object Name", KTampil)), kRO(kol("CoverageNote", "Coverage Name", KTampil))},
				nilai("TSI RNM in IDR")...)...),
			g(atau(IsAneka, IsFire), append([]Unsur{kRO(kol(PropNomorBaris, "Object No", KTampil)),
				kRO(kol("ObjectItemName", "Object Name", KTampil)), kRO(kol("OccupationName", "Occupation", KTampil)),
				kRO(kol("CoverageNote", "Coverage Name", KTampil))}, nilai("TSI RNM in IDR")...)...),
			g(IsMarineCargo, kRO(kol(PropKlikAdjMarine, "Click Here to View Coverage and Input Adjustment", KTampil))),
		), beda(JalurIsTreatyIn, "1")),
		g(IsMBU, append([]Unsur{kRO(kol(PropNomorBaris, "Coverage No", KTampil)),
			kRO(kol("CoverageNote", "Coverage Name", KTampil))}, nilai("TSI RNM (IDR)")...)...),
		g(IsPA, append(append([]Unsur{kRO(kol(PropNomorBaris, "Object No", KTampil)),
			kRO(kol("CoverageNote", "Coverage Name", KTampil))}, nilai("TSI RNM in IDR")...),
			catatanK(komitePA("SendClaimCommittee", "Send Claim to committee"), OQKomitePA))...),
		g(IsTravel, append(append([]Unsur{kRO(kol(PropNomorBaris, "Coverage No", KTampil)),
			kRO(kol("CoverageID", "Coverage No", KTampil)), kRO(kol("CoverageNote", "Coverage Name", KTampil))},
			nilai("TSI RNM (IDR)")...), catatanK(komitePA("SentClaimCommittee", "Sent Claim to committee"), OQKomitePA))...),
	}
}

// catatanK - catatan (OQ) pada kolom tombol grid.
func catatanK(u Unsur, c string) Unsur { u.Catatan = c; return u }

// JudulNomorBaris - judul kolom nomor baris grid yang di XML tanpa judul (Adjustment, Spreading Adjustment, BreakDown
// QS, Committee Accept Status).
const JudulNomorBaris = "No"

// PropKlikAdjMarine - label baris grid item Marine Cargo layar Adjustment (ItemListEstimation LS17).
const PropKlikAdjMarine = "pyLabelAdjMarine"

// LayarItemAdj - Section Adjusment_SC (panel baris item i objek o).
func LayarItemAdj(o, i int) []Unsur {
	itemKomite := func(h *Halaman) bool {
		v := AmbilJalur(h, JalurItem(o, i)+".IsKomite")
		return v == "" || v == "0"
	}
	tambah := tampil(ikon(tombol("TambahAdjustment", "", "CountTotalEstimasi"), IkonTambah), itemKomite)
	coverage := Unsur{Jenis: JenisGrid, Jalur: DaftarDiItem(o, i, "CoverageList"), Kolom: []Unsur{
		kRO(kol("CoverageNote", "Coverage Name", KTampil)), kRO(kol("Currency.Name", "Currency", KTampil)),
		kRO(kol("Currency.KURS", "Currency Value (IDR)", KAngka)), kRO(kol("TSILiability", "TSI Object", KAngka)),
		kRO(kol("TSINusantaraRe", "TSI RNM", KAngka)), kRO(kol("TSIinIDR", "TSI RNM in IDR", KAngka)),
	}}
	diterima := func(_ *Halaman, b Baris) bool { return b["IsApproved"] == "1" }
	// Kolom nomor baris (`.pxListSubscript`, judul XML kosong) berjudul "No": tanpa judul angkanya terbaca sebagai nilai
	// kolom Status di sebelahnya (work owner 11-10-2026 "KODE INI PERBAIKI"). Sama untuk grid nomor panel adjustment.
	adj := Unsur{Jenis: JenisGrid, Jalur: DaftarAdj(o, i), Rincian: PanelAdj, Tambah: &tambah, Kolom: []Unsur{
		kRO(kol(PropNomorBaris, JudulNomorBaris, KTampil)),
		kTampil(kSumber(kRO(kol("AcceptanceStatus", "Status", KPilih)), kode("AcceptanceStatus")), diterima),
		kTampil(kRO(kol("AcceptedDate", "Accept Date", KTanggal)), diterima),
		kTampil(kRO(kol("AcceptedNo", "Accept No", KTampil)), diterima),
		kRO(kol("DLA_No", "DLA No", KTampil)),
		kRO(kol("pxCreateOpName", "PIC Name", KTampil)),
		kTampil(ikonK(kTombol("HapusAdjustment", "", "DisableSendComite", diterima), IkonHapus),
			func(_ *Halaman, b Baris) bool { return b["IsKomite"] != "1" }),
	}}
	return []Unsur{
		sebaris("",
			tombolOQ("PaymentClaim", "View Status Payment Claim", OQLayananBayar),
			tombolOQ("PaymentAttachment", "View Payment Attachment", OQLayananBayar),
		),
		tampil(tampil(coverage, IsMarineCargo), beda(JalurIsTreatyIn, "1")),
		adj,
	}
}

// LayarAdjustment - Section InputAdjustment (panel baris adjustment a item i objek o).
func LayarAdjustment(o, i, a int) []Unsur {
	p := JalurAdj(o, i, a) + "."
	j := func(prop string) string { return p + prop }
	v := func(h *Halaman, prop string) string { return AmbilJalur(h, j(prop)) }
	komite := func(h *Halaman) bool { return v(h, "IsKomite") == "1" }
	bukanKomite := func(h *Halaman) bool { return !komite(h) }
	ada := func(prop string) Kondisi { return func(h *Halaman) bool { return v(h, prop) != "" } }
	pt := func(h *Halaman) string { return v(h, "PaymentType") }
	dasar := func(h *Halaman) bool { t := pt(h); return t != BayarFee && t != BayarSalvage && t != BayarExpense }
	fee := func(h *Halaman) bool { t := pt(h); return t == BayarFee || t == BayarExpense }
	salvage := func(h *Halaman) bool { return pt(h) == BayarSalvage }
	diproses := func(h *Halaman) bool { return v(h, "AcceptanceStatus") != "" }
	analis := func(h *Halaman) bool { return v(h, "IsAnalisatorTransfer") == "1" && h.Ambil(JalurBisnisPA) == "PA" }
	exGratia := func(h *Halaman) bool { return v(h, "ExGratia") == "1" }
	bExGratia := func(h *Halaman, _ Baris) bool { return v(h, "ExGratia") != "1" }
	qs := func(h *Halaman) bool {
		rows := h.AmbilDaftar(DaftarDiAdj(o, i, a, AnakAdjQS))
		return len(rows) > 0 && rows[0]["TreatyName"] != ""
	}
	payable := func(h *Halaman) bool { return h.Ambil(CD+"Payable") == "1" }
	cedingAda := func(h *Halaman) bool { return AmbilJalur(h, JalurBaris(DaftarCedingCo, 1)+".CedingCoName") != "" }
	spreadKosong := func(h *Halaman) bool {
		rows := h.AmbilDaftar(DaftarDiAdj(o, i, a, AnakAdjSpread))
		return len(rows) == 0 || rows[0]["TreatyType"] == ""
	}
	angkaRO := func(prop, lbl string) Unsur { return ro(medan(j(prop), lbl, KAngka)) }

	// LS3 + LS4: masing-masing tiga medan sebaris -> satu Inline grid triple berlabel atas (work owner 11-10-2026
	// "RAPIHKAN!"; dahulu dua baris Inline label kiri yang berdesakan). Label "Lega lLiability" XML diketik ulang
	// "Legal Liability" seperti pop-up Comittee yang mengisi medan yang sama.
	nilaiKomite := func(prop, lbl string) Unsur {
		return naJika(medan(j("DataCommitteFacin."+prop), lbl, KArea), selalu)
	}
	kronologi := tampil(bagian("Detail Chronology to Comitee",
		tiga(
			nilaiKomite("CircumCauseOfLoss", "Chronology"),
			nilaiKomite("ExtentOfLoss", "Extent Of Loss"),
			nilaiKomite("LegalLiability", "Legal Liability"),
			nilaiKomite("Salvage", "Salvage"),
			nilaiKomite("AdjusterFee", "Adjuster Fee"),
			nilaiKomite("Remarks", "Remarks"),
		),
	), komite)

	cedant := Unsur{Jenis: JenisGrid, Jalur: DaftarDiAdj(o, i, a, AnakCedingCedant), Kolom: []Unsur{
		kRO(kol("CedingCoName", "Ceding Co Name", KTeks)), kRO(kol("ShareCeding", "Share (%)", KAngka)),
	}}
	cedingPolis := tampil(tampil(Unsur{Jenis: JenisGrid, Jalur: DaftarCedingCo, Kolom: []Unsur{
		kRO(kol("CedingCoName", "Ceding Co Name", KTampil)),
		kTombol("PilihPayableTo", "Choose", "SetPayableTo", nil),
	}}, dan(bukanKomite, payable)), cedingAda)

	// Rupa layout lama Pega (screenshot work owner 11-10-2026 "PERBAIKI TAMPILAN ADJUSTMENT! INI SAYA LAMPIRKAN LAYOUT
	// LAMA"): LS12 dan LS15 = Inline grid triple (tiga sel per baris, label di atas isian); setiap sel = kontainer XML-nya
	// (LS16-LS38) - medan sejajar berlabel di atas (`kolom`), Payable / bank Stacked with labels left (`kiri`). Kontainer
	// tersembunyi tidak memakan sel. Medan, label, kondisi, dan aksi tetap InputAdjustment.
	detail := []Unsur{
		tombol("ViewCedantPanel", "View Inc Fac Cedant Panel", "BukaCedant"),
		cedant,
		tiga( // LS12
			bagian("", // LS13
				aksi(roJika(medan(j("DLANoCeding"), "DLA No Ceding", KTeks), komite), "SetDLACedingSOB"),
				aksi(roJika(medan(j("DLANoSOB"), "DLA No SOB", KTeks), komite), "SetDLACedingSOB"),
			),
			naJika(roJika(medan(j("NetForCollection"), "Net For Collection", KCentang), komite), komite),
			bagian("", // LS14 (GetStatusKasir_Act)
				naJika(roJika(medan(j("DirectToKasir"), "Transfer Direct to Kasir", KCentang), komite),
					atau(komite, func(h *Halaman) bool { return v(h, "IsFacRetro") == "1" })),
				tampil(ro(medan(j("StatusKasir"), "", KTeks)),
					func(h *Halaman) bool { return v(h, "DirectToKasir") == "true" && v(h, "AcceptanceStatus") == "1" }),
			),
		),
		tiga( // LS15 "Inline grid triple"
			kolom( // LS16
				tampil(aksi(sumber(wajibU(roJika(medan(j("PaymentType"), "Payment Type", KPilih), atau(diproses, analis))),
					kode("PaymentType")), "SetAdjTypePayment"), IsPA),
				tampil(aksi(sumber(wajibU(roJika(medan(j("PaymentType"), "Payment Type", KPilih), komite)),
					kode("PaymentType")), "SetAdjTypePayment"), func(h *Halaman) bool { return !IsPA(h) }),
				roJika(medan(j("FormType"), "Payment Type Detail", KTeks), komite),
			),
			kolom( // LS17
				ro(medan(j("CurrencyEstimasi"), "Currency Estimation", KTeks)),
				naJika(angkaRO("EstimationValue", "Estimation RNM"), diproses),
				naJika(angkaRO("TotalEstimasiValue", "Total Estimation RNM in IDR"), diproses),
			),
			kolom(angkaRO("PersenRNM", "RNM Share(%)")), // LS18
			kolom( // LS19
				tampil(aksi(sumber(wajibU(roJika(medan(j("UploadLOD"), "Choose Currency", KPilih), komite)),
					SumberMataUangAdj), "CheckCurrency"), bukanKomite),
				tampil(ro(medan(j("Currency"), "Currency Adjustment", KTeks)), ada("Currency")),
			),
			kolom( // LS20
				aksi(wajibU(roJika(medan(j("GrossAdjustment"), "Gross Adjustment (100%)", KAngka), komite)),
					"SetGrossAdjustment"),
				angkaRO("GrossValue", "Gross  Adjustment RNM"),
			),
			tampil(kolom( // LS21
				aksi(sumber(wajibU(roJika(medan(j("IndividualRiskType"), "Deductible Type", KPilih), komite)),
					kode("IndividualRiskType")), "SetNilaiResikoSendiri"),
				naJika(aksi(roJika(medan(j("IndividualRiskPercentage"), "Deductible (%)", KAngka),
					func(h *Halaman) bool { return v(h, "IndividualRiskType") == "3" || diproses(h) || komite(h) }),
					"SetNilaiResikoSendiri"), atau(diproses, analis)),
			), dasar),
			tampil(kolom( // LS22
				aksi(roJika(medan(j("IndividualRiskValue"), "Deductible (100%)", KAngka),
					func(h *Halaman) bool {
						t := v(h, "IndividualRiskType")
						return t == "1" || t == "2" || diproses(h) || komite(h)
					}), "SetNilaiResikoSendiri"),
				angkaRO("IndividualRiskRNM", "Deductible RNM"),
			), dasar),
			tampil(kolom( // LS24
				aksi(wajibU(roJika(medan(j("VAT"), "VAT (%)", KAngka), atau(diproses, komite))), "SetValueAdjusterFee"),
			), fee),
			tampil(kolom(angkaRO("VATValue", "VAT Value(100%)"), angkaRO("IndividualRiskRNM", "VAT Value RNM")), fee), // LS25
			tampil(kolom(angkaRO("SalvageValue", "Salvage RNM"), angkaRO("ValueAdjustment", "Salvage RNM in IDR")), // LS26
				salvage),
			tampil(kolom( // LS27
				angkaRO("ProposeAdjustmentValue", "Nett Adjustment (100%)"),
				angkaRO("AdjustmentValue", "Adjustment RNM"),
				angkaRO("ValueAdjustment", "Adjustment RNM in IDR"),
			), dasar),
			tampil(kolom( // LS28
				angkaRO("ProfessionalFee", "Nett Professional Fee(100%)"),
				angkaRO("AdjusterFeeValue", "Adjuster Fee RNM"),
				angkaRO("ValueAdjustment", "Adjuster Fee RNM in IDR"),
			), fee),
			tampil(kiri( // LS29
				sumber(ro(medan(j("Payable"), "Payable To", KPilih)), kode("Payable")),
				ro(medan(j("PayableTo"), "Specify", KTeks)),
			), komite),
			tampil(kiri( // LS30 + LS31 (ceding polis)
				aksi(sumber(wajibU(medan(CD+"Payable", "Payable To", KPilih)), kode("Payable")), "SetPayable"),
				aksi(medan(CD+"PayableTo", "Specify", KTeks), "SetPayable"),
				cedingPolis,
			), bukanKomite),
			kiri( // LS37
				aksi(sumber(roJika(medan(j("NameOfBank"), "Name of Bank", KOtomatis), komite), SumberRekening),
					"PilihRekening"),
				tampil(ro(medan(j("SwiftCode"), "Swift Code", KTeks)), ada("SwiftCode")),
				ro(medan(j("BranchOfBank"), "Branch of Bank", KTeks)),
				ro(medan(j("NoAccount"), "Account No", KTeks)),
			),
			kolom( // LS38 - label KomiteNo XML "Text Input" (bawaan); layout lama "Komite Number"
				tampil(naJika(tombolOQ("ViewKMTNo", "View KMT NO", OQKomiteTahap2), ada("KomiteNo")),
					func(h *Halaman) bool { return komite(h) && !diproses(h) }),
				tampil(ro(medan(j("KomiteNo"), "Komite Number", KTeks)), dan(komite, ada("KomiteNo"))),
			),
		),
	}

	spread := tampil(Unsur{Jenis: JenisGrid, Jalur: DaftarDiAdj(o, i, a, AnakAdjSpread), Label: "Spreading Adjustment",
		Tambah: ptr(tampil(ikon(tombol("TambahSpreadAdj", "Add", "TambahSpreadAdj"), IkonTambah), exGratia)),
		Kolom: []Unsur{
			kRO(kol(PropNomorBaris, JudulNomorBaris, KTampil)),
			kROJ(kSumber(kol("TreatyType", "Treaty Type", KPilih), SumberJenisReas), bExGratia),
			kROJ(kAksi(kol("SharePercentage", "Share %", KAngka), "CountSpreadingAdjustment"), bExGratia),
			kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka)),
			kTampil(catatanK(kTombol("ViewRetro", "View Retro", "", func(*Halaman, Baris) bool { return true }),
				OQKomitePA), bSama("TreatyType", TreatyFacRetro)),
			kTampil(ikonK(kTombol("HapusSpreadAdj", "Delete", "HapusSpreadAdj", nil), IkonHapus),
				func(h *Halaman, _ Baris) bool { return exGratia(h) }),
		}, Kaki: []Unsur{ // LS46 (footer body3 XML: Total Share / Total Claim Spread)
			angkaRO("TotalSharePersen", "Total Share Percentage(%)"),
			angkaRO("TotalSpreadAdjustment", "Total Claim Spread"),
		}}, func(h *Halaman) bool { return !IsTravel(h) })
	quota := tampil(Unsur{Jenis: JenisGrid, Jalur: DaftarDiAdj(o, i, a, AnakAdjQS),
		Label: "BreakDown Spreading Quota Share (QS)", Kolom: []Unsur{
			kRO(kol(PropNomorBaris, JudulNomorBaris, KTampil)),
			kSumber(kRO(kol("TreatyType", "Treaty Type", KPilih)), SumberJenisReas),
			kRO(kol("SharePercentage", "Share %", KAngka)), kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka)),
		}, Kaki: []Unsur{ // LS51
			angkaRO("TotalSharePersen", "Total Share Percentage(%)"),
			angkaRO("TotalSpreadBreakQs", "Total Claim Spread"),
		}}, qs)
	komiteGrid := gridKomite(DaftarDiAdj(o, i, a, AnakKomiteAdj))

	// LS1 layout group TAB (work owner 11-10-2026 "INI BENTUK TAB!"): tab "Detail Chronology to Comitee" (LS2, VIS
	// IsKomite = 1) dan "Adjusment Detail" (LS5); LS6 - spreading (LS39), Committee Accept Status (LS52), tombol (LS57) -
	// berada DI DALAM tab Adjusment Detail.
	adjDetail := append(detail,
		bagian("",
			kolom(aksi(sumber(roJika(medan(j("ExGratia"), "Ex Gratia", KRadio), komite), kode("ExGratia")), // LS41
				"CekExGratia")),
			spread,
			quota,
		),
		komiteGrid,
		sebaris("",
			tombol("SimpanAdjustment", "Save", "SimpanAdjustment"),
			naJika(tombol("SendToCommitte", "Send to Committe", "SendPICProtect"), atau(komite, spreadKosong)),
			tampil(naJika(tombol("Acceptation", "Acceptation", "Acceptation"),
				func(h *Halaman) bool { return !BolehAkseptasi(baris(h, j)) }),
				func(h *Halaman) bool { return v(h, "AcceptanceStatus") == "1" }),
		),
	)
	return []Unsur{
		letak(LetakTab, kronologi, bagian("Adjusment Detail", adjDetail...)),
	}
}

// baris - salinan medan baris adjustment dari jalurnya (untuk kondisi tombol).
func baris(h *Halaman, j func(string) string) Baris {
	b := Baris{}
	for _, p := range []string{"AcceptanceStatus", "IsPrintAccept", "DirectToKasir", "StatusKasir"} {
		b[p] = AmbilJalur(h, j(p))
	}
	return b
}

// OQKomiteTahap2 - "View KMT NO" (ViewKomite_act: `@substring(pxCoveredInsKeys(<last>), 19)`) membaca kasus komite
// tertaut - KomiteNo ditulis saat KMT lahir (models/komite.go); tombolnya tidak diperlukan.
const OQKomiteTahap2 = "KomiteNo diisi saat kasus KMT lahir (CreateKMTNo_Act); ViewKomite_act hanya membaca ulang nomor itu"

// JalurBisnisPA - `pyWorkPage.Policy.Quotation.BusinessType` (When isAnalistorTransfer B). Halaman `Policy` tidak
// pernah diisi rule Claim Fac In - selalu kosong, sehingga isAnalistorTransfer selalu salah (ditiru apa adanya).
const JalurBisnisPA = "Policy.Quotation.BusinessType"

// ---------------------------------------------------------------- Cedant Panel

// LayarCedant - Section ViewCedantPanel (harness ViewCedantpanels) untuk adjustment a: ringkasan polis, grid
// CedingCedantList polis (Choose / Choose All saat adjustment belum berkomite).
func LayarCedant(o, i, a int) []Unsur {
	komite := func(h *Halaman) bool { return AmbilJalur(h, JalurAdj(o, i, a)+".IsKomite") == "1" }
	daftar := AwalanPolis + ".CedingCedantList"
	pilih := tampil(Unsur{Jenis: JenisGrid, Jalur: daftar,
		Tambah: ptr(tombol("PilihSemuaCedant", "Choose All", "SetCedant:2")), Kolom: []Unsur{
			kRO(kol("CedingCoName", "Ceding", KTampil)), kRO(kol("ShareCeding", "% Share", KAngka)),
			kTombol("PilihCedant", "Choose", "SetCedant:1", nil),
		}}, func(h *Halaman) bool { return !komite(h) })
	lihat := tampil(Unsur{Jenis: JenisGrid, Jalur: daftar, Kolom: []Unsur{
		kRO(kol("CedingCoName", "Ceding", KTampil)), kRO(kol("ShareCeding", "% Share", KAngka)),
	}}, komite)
	return []Unsur{
		sebaris("",
			ro(medan(OQ+"SobName", "Source of Business", KTeks)),
			ro(medan(AwalanPolis+".ShareCedantType", "Share Cedant Type", KTeks)),
			ro(medan(AwalanPolis+".PercentShare", "% Share RNM", KAngka)),
			ro(medan(AwalanPolis+".TotalTSINusaReSpreading", "Total TSI RNM", KAngka)),
			ro(medan(AwalanPolis+".TotalPremiNusaRe", "Total PRemi RNM", KAngka)),
		),
		pilih, lihat,
	}
}

// PropPilihCedant - pilihan Cedant Panel adjustment yang DISIMPAN (CEDANT_CHOICE): "*" = Choose All, selainnya kode
// CedingCo. `.CedingCedantList` adjustment adalah salinan baris CedingCedantList polis (SetCedant_act 4 / 5.2) -
// polis dimuat ulang setiap aksi, sehingga hanya pilihannya yang disimpan dan daftarnya disusun ulang
// (`LengkapiCedant`). `[inferensi]` (PARITAS): tanpa tabel baru.
const PropPilihCedant = "CedantChoice"

// PilihSemuaCedant - nilai PropPilihCedant untuk Choose All.
const PilihSemuaCedant = "*"

// SetCedant = SetCedant_act (Param.ChooseType 2 = Choose All, 1 = Choose satu ber-Param.ID): CedingCedantList
// adjustment dari polis; Choose satu juga menyetel PersenRNM = ShareCeding baris itu (6).
func SetCedant(h *Halaman, o, i, a int, tipe, id string) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	switch tipe {
	case "2": // 4
		b[PropPilihCedant] = PilihSemuaCedant
	case "1": // 5-6
		b[PropPilihCedant] = id
		share := ""
		for _, c := range h.AmbilDaftar(AwalanPolis + ".CedingCedantList") {
			if c["CedingCo"] == id {
				share = c["ShareCeding"]
			}
		}
		b["PersenRNM"] = share
	default:
		return ErrBarisTidakAda
	}
	lengkapiCedantBaris(h, o, i, a, b)
	return nil
}

func lengkapiCedantBaris(h *Halaman, o, i, a int, b Baris) {
	pil := b[PropPilihCedant]
	var out []Baris
	for _, c := range h.AmbilDaftar(AwalanPolis + ".CedingCedantList") {
		if pil == PilihSemuaCedant || (pil != "" && c["CedingCo"] == pil) {
			out = append(out, Baris{"CedingCo": c["CedingCo"], "CedingCoName": c["CedingCoName"],
				"ShareCeding": c["ShareCeding"]})
		}
	}
	h.SetelDaftar(DaftarDiAdj(o, i, a, AnakCedingCedant), out)
}

// LengkapiCedant menyusun ulang `.CedingCedantList` setiap adjustment dari pilihannya (sesudah polis dimuat).
func LengkapiCedant(h *Halaman) {
	for o := range h.AmbilDaftar(DaftarObjek) {
		for i := range h.AmbilDaftar(DaftarItem(o + 1)) {
			for a, b := range h.AmbilDaftar(DaftarAdj(o+1, i+1)) {
				lengkapiCedantBaris(h, o+1, i+1, a+1, b)
			}
		}
	}
}

// ---------------------------------------------------------------- spreading adjustment Ex Gratia

// TambahSpreadAdj = tombol "Add" grid Spreading Adjustment (addRow, VIS ExGratia = 1).
func TambahSpreadAdj(h *Halaman, o, i, a int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	if b["ExGratia"] != "1" {
		return ErrBarisTidakAda
	}
	h.TambahBaris(DaftarDiAdj(o, i, a, AnakAdjSpread), Baris{})
	return nil
}

// HapusSpreadAdj = tombol "Delete" baris n (deleteRow) lalu CountSpreadingAdjustment_ACT.
func HapusSpreadAdj(h *Halaman, o, i, a, n int) error {
	b, err := Adj(h, o, i, a)
	if err != nil {
		return err
	}
	if b["ExGratia"] != "1" {
		return ErrBarisTidakAda
	}
	if _, err := barisDi(h, DaftarDiAdj(o, i, a, AnakAdjSpread), n); err != nil {
		return err
	}
	h.HapusBaris(DaftarDiAdj(o, i, a, AnakAdjSpread), n)
	return CountSpreadingAdjustment(h, o, i, a)
}

// LabelAdjMarine - label baris grid item Marine Cargo layar Adjustment (dihitung HitungTurunan).
func LabelAdjMarine(h *Halaman) {
	if !IsMarineCargo(h) {
		return
	}
	for o := range h.AmbilDaftar(DaftarObjek) {
		for _, it := range h.AmbilDaftar(DaftarItem(o + 1)) {
			it[PropKlikAdjMarine] = "Click Here to View Coverage and Input Adjustment"
		}
	}
}
