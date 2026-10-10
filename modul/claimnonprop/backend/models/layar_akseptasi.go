package models

// Untuk apa berkas ini: DEFINISI LAYAR DETAIL AKSEPTASI DAN MODAL - Section `AdjustmentDetailNP_Section`
// (`AdjustmentDetailNP` + `Subjectivity`; FlowAction AdjustmentDetailNP = masterDetail Acceptation List),
// `KomiteCLMNP` (harness KomiteCNP), `CloseClaimMD`, `CloseClaimNP`, `PreviewPLA` (FlowAction GeneratePLACNP),
// `InputDtlInterest` (masterDetail Insured Interests), `ReinstatementPremiumDetails` (FlowAction ShowDetailXOL,
// masterDetail XOL Allocation akseptasi).
//
// ⚠️ Kelas activity grid akseptasi (OQ-CNP-33, PARITAS): sel AdjustmentDetailNP memanggil `AddListClaimNP_Act`,
// `SetCurrency_Act`, `CountClaimTNP_Act`, `AddLossAlocation_Act`, `CountLossAllocation_act` dengan kelas
// ASM-FW-GCNMFW-Data-Adjustment, padahal activity itu hanya ada di kelas Work-ClaimTreatyNonProp - di Pega rule tidak
// ditemukan. Di sini: tombol Add nonaktif (OQ), sel postValue saja. Satu-satunya sel berkelas Work (`Rate of Exchange` ->
// CountClaimTNP_Act) menjalankan hitungan tingkat KLAIM (activity menulis jalur absolut pyWorkPage), persis XML.

// adjK membentuk kondisi atas baris akseptasi ke-n.
func adjK(n int, f func(b Baris) bool) Kondisi {
	return func(h *Halaman) bool {
		d := h.AmbilDaftar(DaftarAdjustment)
		if n < 1 || n > len(d) {
			return false
		}
		return f(d[n-1])
	}
}

func adjSama(n int, p, v string) Kondisi { return adjK(n, func(b Baris) bool { return b[p] == v }) }

// LayarDetailAkseptasi - Section `AdjustmentDetailNP_Section` untuk baris akseptasi ke-n.
func LayarDetailAkseptasi(n int) []Unsur {
	j := func(p string) string { return JalurAdj(n, p) }
	komite := adjSama(n, "IsKomite", "1")
	tombolOQKelas := func(id string) *Unsur { return ptr(ikon(tombolOQ(id, "Add", OQKelasAkseptasi), IkonTambah)) }
	pv := func(p, judul string) Unsur { return catatanK(kNA(kol(p, judul, KAngka), bTerkunci), OQKelasAkseptasi) }
	klaim := Unsur{Jenis: JenisGrid, Jalur: j(AnakClaimAccept), Bernomor: true, Tambah: tombolOQKelas("AddListClaimNPAkseptasi"),
		Kolom: []Unsur{
			catatanK(kSumber(kNA(kROJ(kol("CurrencyID", "Currency", KPilih), bNote), bTerkunci), SumberMataUang), OQKelasAkseptasi),
			kAksi(kNA(kol("AltValue", "Rate of Exchange", KAngka), bTerkunci), "CountClaimTNP"),
			kROJ(pv("Value", "Claim Amount"), bNote),
			kROJ(pv("TPL", "TPL"), bNote),
			pv("AdjusterFee", "Adjuster Fee"),
			pv("Salvage", "Salvage"),
			pv("CNPOthersFee", "Fee"),
			kRO(kol("PctProrateClaim", "Proportion (%)", KAngka)),
			kRO(kol("USD", "Claim Amount in IDR", KAngka)),
			kRO(kol("ClaimAmountCedant", "Claim Amount Cedant", KAngka)),
			ikonK(kTombol("HapusKlaimAkseptasi", "Delete", "HapusKlaimAkseptasi", bTerkunci), IkonHapus),
		},
		Kaki: []Unsur{label("Total Claim Amount"), ro(medan(CD+"TotalListClaimAmount", "", KAngka)),
			ro(medan(CD+"TotalListClaimAmountIDR", "", KAngka))}}
	loss := Unsur{Jenis: JenisGrid, Jalur: j(AnakLossAlloc), Bernomor: true, Tambah: tombolOQKelas("AddLossAlocationAkseptasi"),
		Kolom: []Unsur{
			kSumber(kRO(kol("CurrencyID", "Currency", KPilih)), SumberMataUang),
			kSumber(kNA(kol("TreatyName", "Treaty Name", KPilih), bTerkunci), SumberLossAlloc),
			pv("ClaimPercentage", "Share (%)"),
			pv("ClaimAmountAdjust", "Claim Amount"),
			pv("AdjusterFee", "Adjuster Fee"),
			pv("Salvage", "Salvage"),
			pv("CNPOthersFee", "Fee"),
			catatanK(kNA(kol("CNPFlagXOL", "To XOL", KCentang), bTerkunci), OQKelasAkseptasi),
			ikonK(kTombol("HapusLossAkseptasi", "Delete", "HapusLossAkseptasi", bTerkunci), IkonHapus),
		}}
	xol := func(jalur string, ubah bool) Unsur {
		total := kRO(kol("TotalClaim", "Claim Amount", KAngka))
		if ubah {
			total = kAksi(kNA(kol("TotalClaim", "Claim Amount", KAngka), bTerkunci), "AdjClaimCNP")
		}
		return Unsur{Jenis: JenisGrid, Jalur: jalur, Bernomor: true, Kolom: []Unsur{
			kSumber(kRO(kol("CurrencyID", "Currency", KPilih)), SumberMataUang),
			kRO(kol("TreatyName", "Treaty Name", KTampil)),
			total,
			kRO(kol("ClaimPercentage", "RNM Share (%)", KAngka)),
			kRO(kol("ClaimSpreaded", "Claim Amount RNM", KAngka)),
			kRO(kol("AdjusterFee", "Adjuster Fee", KAngka)),
			kRO(kol("Salvage", "Salvage", KAngka)),
			kRO(kol("CNPOthersFee", "Fee", KAngka)),
			kRO(kol("TotalClaimRNM", "Total Claim RNM", KAngka)),
			kRO(kol("CNPReinstatementRNM", "Reinstatement Premium", KAngka)),
		}}
	}
	spread := func(jalur string) Unsur {
		return Unsur{Jenis: JenisGrid, Jalur: jalur, Bernomor: true, Kolom: []Unsur{
			kRO(kol("Currency", "Currency", KTeks)),
			kRO(kol("TreatyName", "Treaty Type", KTeks)),
			kRO(kol("SharePercentage", "Share(%)", KAngka)),
			kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka)),
			kRO(kol("AdjusterFee", "Adjuster Fee", KAngka)),
			kRO(kol("Salvage", "Salvage", KAngka)),
			kRO(kol("CNPOthersFee", "Fee", KAngka)),
			kRO(kol("TotalClaim", "Total Claim RNM", KAngka)),
			kRO(kol("PremiumSpreaded", "Reinstatement Premium", KAngka)),
			kRO(kol("NetClaim", "RNM Net Claim", KAngka)),
		}}
	}
	lain := beda(CD+"Payable", "3")
	tiga := sama(CD+"Payable", "3")
	bank := func(sfx, aksiNama string) Unsur {
		mati := func(u Unsur) Unsur { return naJika(roJika(u, komite), selalu) }
		return bagian("",
			wajibU(aksi(sumber(roJika(medan(j("NameOfBank"+sfx), "Name of Bank", KOtomatis), komite), SumberRekening),
				aksiNama)),
			mati(medan(j("Currency"+sfx), "Currency", KTeks)),
			tampil(mati(medan(j("SwiftCode"+sfx), "Swift Code", KTeks)), adjK(n, func(b Baris) bool { return b["SwiftCode"+sfx] != "" })),
			mati(medan(j("BranchOfBank"+sfx), "Branch of Bank", KTeks)),
			wajibU(mati(medan(j("NoAccount"+sfx), "Account No", KTeks))),
			tampil(ro(medan(j("KomiteNo"), "", KTeks)), komite),
		)
	}
	detail := tampil(bagian("",
		dua(
			bagian("",
				wajibU(aksi(sumber(naJika(medan(j("PaymentType"), "Payment Type", KPilih), komite), kode("PaymentType")),
					"SetInterimXOL")),
				tampil(ro(medan(j("CNPIndexInterim"), "", KTampil)), adjSama(n, "PaymentType", "2")),
			),
			bagian("", roJika(medan(j("DLANoCeding"), "No DLA :", KTeks), komite)),
		),
		klaim,
		bagian("Loss Allocation", loss),
		bagian("XOL Allocation", xol(j(AnakXOL), true)),
		bagian("Previously Calculated", xol(j(AnakXOLDibayar), false)),
		bagian("",
			bagian("Spreading In", spread(j(AnakSpreadIn))),
			bagian("Spreading Out", spread(j(AnakSpreadOut))),
		),
		dua(
			bagian("",
				wajibU(aksi(sumber(roJika(medan(CD+"Payable", "Payable To", KPilih), komite), kode("Payable")),
					"SetPayableTreatyNP:Acc")),
				tampil(wajibU(aksi(roJika(medan(CD+"PayableTo", "Specify", KTeks), komite), "SetPayableTreatyNP")), lain),
				tampil(wajibU(aksi(sumber(roJika(medan(CD+"PayableTo", "Specify", KPilih), komite), SumberKlien),
					"SetPayableTreatyNP")), tiga),
				tampil(wajibU(roJika(medan(RC+"Name", "Name", KTeks), komite)), tiga),
			),
			bank("", "PilihRekening"),
			bagian(""),
			tampil(bank("2", "PilihRekening2"), adjSama(n, "FlagCurrency", "1")),
		),
		sebaris("",
			naJika(roJika(medan(j("DirectToKasir"), "Transfer Direct to Kasir", KCentang), komite),
				atau(komite, adjSama(n, "FlagErrorKasir", "1"))),
			tampil(ro(medan(j("StatusKasir"), "", KTeks)),
				adjK(n, func(b Baris) bool { return b["DirectToKasir"] == "true" && b["AcceptanceStatus"] == "1" })),
		),
		gridKomite(j(AnakKomite)),
		sebaris("",
			tombol("SaveAkseptasi", "Save", "Simpan"),
			tampil(naJika(tombol("SendToCommitte", "Send to Committe", "BukaKomite"),
				atau(komite, sama("IsSaveToOs", "0"))), adjK(n, func(b Baris) bool { return b["TotalKomite"] != "" })),
			naJika(tombol("Acceptation", "Acceptation", "HitServiceToKasir"),
				adjK(n, func(b Baris) bool { return b["AcceptanceStatus"] != "1" || b["IsPrintAccept"] == "1" })),
			tombolOQ("GenerateDLA", "Generate DLA", OQDLA),
			tombol("GenerateClaimAnalysis", "Generate Claim Analysis", "GenerateCACNP"),
			tombol("SavePreviouslyPaid", "Save Previously Paid", "SaveCNPLayerList"),
		),
	), adjK(n, func(b Baris) bool { return b["CommentLOD"] != "1" }))
	return []Unsur{
		detail,
		bagian("", sumber(ro(medan(j("SubjectivityNote"), "Subjectivity Note", KPilih)), kode("SubjectivityNote"))),
	}
}

// ---------------------------------------------------------------- modal

// AksiSegarKomite - change textarea pop-up komite ("Post value" tanpa activity) = aksi Simpan tanpa langkah.
const AksiSegarKomite = "Simpan"

// LayarKomite - Section `KomiteCLMNP` (harness KomiteCNP) untuk akseptasi ke-n. Gerbang tombol Send =
// `CekError.CARI1..10` (ProteksiSendKomiteCNP_Act + ProteksiNilaiClaim, dihitung ulang server saat aksi `BukaKomite`
// dan saat dikirim) diringkas `CekError.Lolos`, ditambah `Payable` dan `Occupation` kosong.
func LayarKomite(n int) []Unsur {
	j := func(p string) string { return JalurAdj(n, p) }
	return []Unsur{
		bagian("",
			naJika(medan(JalurTanggalKomite, "Date", KTanggal), selalu),
			ro(medan(JalurPICKomite, "Initial", KTeks)),
		),
		label("Adjustment"),
		bagian("",
			aksi(naJika(medan(j("DataCommitteeTreaty.CircumCauseOfLoss"), "Circumstanses", KArea), adjSama(n, "Type", "3")),
				AksiSegarKomite),
			aksi(medan(j("DataCommitteeTreaty.Remarks"), "Remarks", KArea), AksiSegarKomite),
		),
		sebaris("",
			naJika(tombol("SendClaimToCommittee", "Send Claim to Committee", "CreateChildKomiteCNP"),
				atau(beda(JalurLolosKomite, "1"), sama(CD+"Payable", ""), sama(CD+"Occupation", ""))),
			tombol("CancelKomite", "Cancel", "TutupModal"),
		),
	}
}

// LayarTutupKlaim - Section `CloseClaimMD` (local action Close Claim).
func LayarTutupKlaim() []Unsur {
	return []Unsur{
		label("Are you sure want to close this claim?"),
		wajibU(medan("Message", "Remarks", KArea)),
		sebaris("",
			tombol("CloseNo", "No", "TutupModal"),
			tombol("CloseYes", "Yes", "CloseClaimTNonProp"),
		),
	}
}

// LayarCWP - Section `CloseClaimNP` (local action Close Without Payment, pra-proses CloseClaimNP_preAct). "Yes"
// (CreateChildKomiteCloseNP_Act) nonaktif: kasus komite tanpa baris akseptasi tidak dapat ditulis (OQ-CNP-36).
func LayarCWP() []Unsur {
	return []Unsur{
		naJika(medan(JalurTanggalKomite, "Date", KTanggal), selalu),
		ro(medan(JalurPICKomite, "Initial", KTeks)),
		label("Are you sure want to close this claim?"),
		wajibU(aksi(medan("TempCommiteClaim.Remarks", "Remarks", KArea), AksiSegarKomite)),
		sebaris("",
			tombol("CwpNo", "No", "TutupModal"),
			tombolOQ("CwpYes", "Yes", OQKomiteTanpaBaris),
		),
	}
}

// LayarPLA - Section `PreviewPLA` (FlowAction GeneratePLACNP, activity GeneratePlaCNP_Act membaca `PreviewPLA.CARI14`).
// ⚠️ `pyShowFAButtons=false` - tombol Submit / Cancel bawaan modal Pega (`[inferensi]`, PARITAS).
func LayarPLA() []Unsur {
	return []Unsur{
		medan("PreviewPLA.CARI14", "Remarks", KArea),
		sebaris("",
			tombol("SubmitPLA", "Submit", "GeneratePlaCNP"),
			tombol("CancelPLA", "Cancel", "TutupModal"),
		),
	}
}

// LayarInterest - Section `InputDtlInterest` (masterDetail Insured Interests Outstanding) untuk baris interest ke-i.
// ⚠️ `view.CARI21` (RO / tampil) tidak punya penulis di korpus - dianggap kosong (OQ-CNP-35).
func LayarInterest(i int) []Unsur {
	j := func(p string) string { return JalurAnak(DaftarInterest, i, p) }
	ro1 := isOutstanding
	tpl := func(p, aksiNama string, k string) Unsur { return aksi(roJika(medan(j(p), "", k), ro1), aksiNama) }
	return []Unsur{
		dua(
			bagian("",
				wajibU(medan(j("ObjectName"), "Interest Insured", KTeks)),
				wajibU(medan(j("KursObjectItem"), "Value In IDR", KAngka)),
			),
			bagian("",
				wajibU(aksi(sumber(medan(j("CurrencyID"), "Currency", KPilih), SumberMataUang), "SetCurrency:Interest")),
				wajibU(aksi(medan(j("TSIPerObject"), "Value", KAngka), "CountTotalInterest")),
			),
		),
		bagian("",
			roJika(medan(j("IsTPL"), "Deductible", KCentang), ro1),
			sumber(roJika(medan(j("TPLFormat"), "Format", KPilih), ro1), kode("TPLFormat")),
		),
		tampil(sebaris("",
			tpl("TPLPct", "SetTPLNote", KAngka),
			label("% of"),
			sumber(tpl("TPLType", "SetTPLNote", KPilih), kode("TPLType")),
			sumber(tpl("CNPMinMax", "SetTPLNote", KPilih), kode("MinMax")),
			tampilIsi(ro(medan(j("Currency"), "", KTampil))),
			tpl("TPLAmount", "SetTPLNote", KAngka),
		), barisSama(DaftarInterest, i, "TPLFormat", "2")),
		tampil(sebaris("",
			tampilIsi(ro(medan(j("Currency"), "", KTampil))),
			tpl("TPLAmount2", "SetTPLNote", KAngka),
			sumber(tpl("CNPMinMax2", "SetTPLNote", KPilih), kode("MinMax")),
			tpl("TPLPct2", "SetTPLNote", KAngka),
			label("% of"),
			sumber(tpl("TPLType2", "SetTPLNote", KPilih), kode("TPLType")),
		), barisSama(DaftarInterest, i, "TPLFormat", "1")),
		sebaris("",
			tombol("SubmitInterest", "Submit", "Simpan"),
			tombol("CancelInterest", "Cancel", "TutupModal"),
		),
	}
}

// barisSama - kondisi atas baris ke-i daftar.
func barisSama(daftar string, i int, p, v string) Kondisi {
	return func(h *Halaman) bool {
		d := h.AmbilDaftar(daftar)
		return i >= 1 && i <= len(d) && d[i-1][p] == v
	}
}

// LayarReinstatement - Section `ReinstatementPremiumDetails` (FlowAction ShowDetailXOL, masterDetail XOL Allocation
// akseptasi) untuk layer ke-i akseptasi ke-n. `FlagProrate` tidak punya penulis di korpus (kosong = varian 0, perbandingan
// numerik Pega); varian 2 / 3 tidak dibangun (OQ-CNP-12), varian 1 dibangun sesuai section.
func LayarReinstatement(n, i int) []Unsur {
	jalur := JalurAdj(n, AnakXOL)
	j := func(p string) string { return JalurAnak(jalur, i, p) }
	rumus := func() Unsur {
		return sebaris("",
			ro(medan(j("TotalClaim"), "Claim Layer", KAngka)), label("/"),
			ro(medan(j("CNPLimit"), "Limit Layer", KAngka)), label("*"),
			ro(medan(j("CNPMDP"), "Premium Layer", KAngka)), label("*"),
			ro(medan(j("CNPPctReinstate"), "Reinstatement (%)", KAngka)), label("="),
			ro(medan(j("CNPReinstatement"), "Reinstatement Premium", KAngka)))
	}
	prorate := func(v string) Kondisi { return barisSama(jalur, i, "FlagProrate", v) }
	return []Unsur{
		tampil(bagian("Reinstatement Calculation", rumus()), atau(prorate("0"), prorate(""))),
		tampil(bagian("Reinstatement Calculation", rumus()), prorate("1")),
		bagian("Reinstatement For RNM", sebaris("",
			ro(medan(j("CNPReinstatement"), "Reinstatement Premium", KAngka)), label("*"),
			ro(medan(j("ClaimPercentage"), "Share RNM", KAngka)), label("="),
			ro(medan(j("CNPReinstatementRNM"), "Reinstatement Premium RNM", KAngka)))),
		tombol("TutupReinstatement", "Cancel", "TutupModal"),
	}
}
