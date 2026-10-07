package models

// Untuk apa berkas ini: DEFINISI LAYAR BARIS ADJUSTMENT DAN MODAL - Section `AdjustmentDetail_Section`
// (`DtlDataCommitte`, `AdjustmentDetail`, `Subjectivity`; FlowAction edit baris AdjustmentDetail), `ComiteeClaimTreaty`
// (harness CommitteeTreaty), `PreventRejectClaimProp`, `GeneratePLA`, `generateDLATreaty` (pindai `cp/layar.md` bab
// 4.3-4.4).

// adjK membentuk kondisi atas baris adjustment ke-n.
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

// LayarAdjustment - Section `AdjustmentDetail_Section` untuk baris adjustment ke-n.
func LayarAdjustment(n int) []Unsur {
	j := func(p string) string { return JalurAdj(n, p) }
	komite := adjSama(n, "IsKomite", "1")
	subj := adjSama(n, "IsSubjectivity", "true")
	kunci := atau(komite, subj)
	bukanKomite := adjK(n, func(b Baris) bool { return b["IsKomite"] != "1" })
	tipe := func(v ...string) Kondisi {
		return adjK(n, func(b Baris) bool {
			for _, x := range v {
				if b["Type"] == x {
					return true
				}
			}
			return false
		})
	}
	diterima := adjSama(n, "AcceptanceStatus", "1")
	dtl := tampil(bagian("",
		naJika(ro(medan(j("DataCommitteeTreaty.CircumCauseOfLoss"), "Circumstanses", KArea)), tipe("3")),
		tampil(ro(medan(j("DataCommitteeTreaty.Salvage"), "Salvage", KArea)), tipe("3")),
		tampil(ro(medan(j("DataCommitteeTreaty.AdjusterFee"), "Adjuster Fee", KArea)), tipe("4")),
		ro(medan(j("DataCommitteeTreaty.Remarks"), "Remarks", KArea)),
	), kunci)
	tombolBaris := tampil(bagian("",
		tombol("SaveAdjustment", "Save", "Simpan"),
		tampil(naJika(tombol("SendToCommitte", "Send to Committe", "BukaKomite"),
			atau(komite, func(h *Halaman) bool { return angkaTeks(h.Ambil("IsError")) > 1 })),
			adjK(n, func(b Baris) bool { return b["TotalKomite"] != "" })),
		tampil(naJika(tombol("Acceptation", "Acceptation", "HitServiceToKasir"),
			adjK(n, func(b Baris) bool { return b["DirectToKasir"] == "false" || b["StatusKasir"] != "" })), diterima),
		tampil(naJika(tombol("GenerateDLA", "Generate DLA", "BukaDLA"),
			adjK(n, func(b Baris) bool { return b["IsFacRetro"] != "1" || b["DLA_No"] != "" })), diterima),
		tampil(catatan(naJika(tombol("PrintClaimAnalysis", "Print Claim Analysis", ""), selalu), OQDokumenPDF), diterima),
	), selalu)
	detail := []Unsur{
		aksi(sumber(wajibU(roJika(medan(j("Type"), "Type", KPilih), kunci)), kode("AdjustmentType")), "SetPayableTreaty"),
		aksi(sumber(wajibU(roJika(medan(j("CurrencyID"), "Currency", KPilih), kunci)), SumberMUAdj), "SetNameCurrency"),
		roJika(medan(j("FormType"), "Payment Type", KTeks), kunci),
		aksi(roJika(medan(j("DLANoCeding"), "DLA No. Ceding", KTeks), kunci), "SetDLACedingSOB"),
		aksi(roJika(medan(j("DLANoSOB"), "DLA No. SOB", KTeks), kunci), "SetDLACedingSOB"),
		ro(medan(j("KursIDR"), "Value In IDR", KAngka)),
		ro(medan(j("PersenRNM"), "RNM Share(%)", KAngka)),
		aksi(sumber(roJika(medan(j("TreatyName"), "Allocation", KOtomatis), kunci), SumberAllocation), "PilihAllocation"),
		ro(medan(j("ShareLossAllocation"), "Share(%)", KAngka)),
		bagian("",
			roJika(medan(j("DirectToKasir"), "DLA No Ceding", KCentang), kunci),
			tampil(ro(medan(j("StatusKasir"), "", KTeks)),
				adjK(n, func(b Baris) bool { return b["DirectToKasir"] == "true" && b["AcceptanceStatus"] == "1" })),
		),
		label("Dedutible Type"), label("Deductible (%)"), label("Currency"), label("Treaty Gross (100%)"),
		label("Currency"), label("RNM Share Claim"),
		ro(medan(j("Currency"), "", KTeks)),
		aksi(roJika(medan(j("GrossAdjustment"), "", KAngka), kunci), "CountGrossAdjTreaty"),
		ro(medan(j("Currency"), "", KTeks)),
		ro(medan(j("GrossValue"), "", KAngka)),
		aksi(sumber(wajibU(roJika(medan(j("IndividualRiskType"), "Text Input", KPilih), kunci)), kode("IndividualRiskType")),
			"CountGrossAdjTreaty"),
		aksi(roJika(medan(j("IndividualRiskPercentage"), "", KAngka), atau(kunci, adjSama(n, "IndividualRiskType", "3"))),
			"CountGrossAdjTreaty"),
		ro(medan(j("Currency"), "", KTeks)),
		aksi(roJika(medan(j("IndividualRiskValue"), "", KAngka), kunci), "CountGrossAdjTreaty"),
		ro(medan(j("Currency"), "", KTeks)),
		ro(medan(j("IndividualRiskRNM"), "", KAngka)),
		label("Total"),
		ro(medan(j("Currency"), "", KTeks)),
		ro(medan(j("ProposeAdjustmentValue"), "", KAngka)),
		ro(medan(j("Currency"), "", KTeks)),
		ro(medan(j("AdjustmentValue"), "", KAngka)),
		{Jenis: JenisGrid, Jalur: j(AnakLossAllocation), Bernomor: true, Kolom: []Unsur{
			kRO(kol("Currency", "Curr ID", KTampil)), kRO(kol("TreatyName", "Treaty Type", KTampil)),
			kRO(kol("SharePercentage", "Share (%)", KAngka)), kRO(kol("ClaimSpreaded", "Result Gross", KAngka)),
			kRO(kol("ClaimEstimation", "Result RNM", KAngka))}},
		tampil(bagian("",
			aksi(sumber(wajibU(roJika(medan(CD+"Payable", "Payable To", KPilih), kunci)), kode("Payable")), "SetPayableTreaty"),
			aksi(wajibU(roJika(medan(CD+"PayableTo", "Specify", KTeks), kunci)), "SetPayableTreaty"),
		), bukanKomite),
		tampil(bagian("",
			aksi(sumber(wajibU(roJika(medan(j("NameOfBank"), "Name of Bank", KOtomatis), kunci)), SumberRekening), "PilihRekening"),
			tampil(naJika(roJika(medan(j("SwiftCode"), "Swift Code", KTeks), komite), selalu),
				adjK(n, func(b Baris) bool { return b["SwiftCode"] != "" })),
			naJika(roJika(medan(j("BranchOfBank"), "Branch of Bank", KTeks), komite), selalu),
			wajibU(naJika(roJika(medan(j("NoAccount"), "Account No", KTeks), komite), selalu)),
		), bukanKomite),
		tampil(bagian("",
			sumber(ro(medan(j("Payable"), "Payable To", KPilih)), kode("Payable")),
			ro(medan(j("PayableTo"), "Specify", KTeks)),
			ro(medan(j("NameOfBank"), "Name of Bank", KTeks)),
			tampil(ro(medan(j("SwiftCode"), "Swift Code", KTeks)), adjK(n, func(b Baris) bool { return b["SwiftCode"] != "" })),
			ro(medan(j("BranchOfBank"), "Branch of Bank", KTeks)),
			ro(medan(j("NoAccount"), "Account No", KTeks)),
			tampil(tombol("ViewKomiteNo", "View Komite No", "SetKomiteNo"), adjSama(n, "KomiteNo", "")),
			ro(medan(j("KomiteNo"), "Komite Number", KTeks)),
		), komite),
		bagian("",
			bagian("Spreading In", Unsur{Jenis: JenisGrid, Jalur: j(AnakSpreadAdj), Bernomor: true, Kolom: []Unsur{
				kRO(kol("Currency", "Currency", KTeks)), kRO(kol("TreatyName", "Treaty Type", KTeks)),
				kRO(kol("SharePercentage", "Share(%)", KAngka)), kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka))}}),
			bagian("Spreading Out", Unsur{Jenis: JenisGrid, Jalur: j(AnakQuotaShare), Bernomor: true, Kolom: []Unsur{
				kRO(kol("Currency", "Currency", KTeks)), kRO(kol("TreatyName", "Treaty Type", KTeks)),
				kRO(kol("SharePercentage", "Share(%)", KAngka)), kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka))}}),
			bagian("Committe Accept Status", catatan(Unsur{Jenis: JenisGrid, Jalur: j("ComiteeClaim"), Bernomor: true,
				Kolom: []Unsur{kRO(kol("IDKomite", "Committee Name", KTampil))}}, OQBatasKomite)),
			tombolBaris,
		),
	}
	out := []Unsur{dtl}
	out = append(out, detail...)
	out = append(out, tampil(bagian("",
		sumber(ro(medan(j("SubjectivityNote"), "Subjectivity Note", KPilih)), kode("SubjectivityNote"))), subj))
	return out
}

// angkaTeks - bilangan bulat teks (0 bila bukan angka).
func angkaTeks(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return n
		}
		n = n*10 + int(c-'0')
	}
	return n
}

// LayarKomite - Section `ComiteeClaimTreaty` (harness CommitteeTreaty) untuk baris ke-n. `lolos` = `Protect.CARI1 = 1
// && Protect.CARI2 = 1` (AttachmentProtect_ACT + CekPremiLunas_Act).
func LayarKomite(n int, lolos bool) []Unsur {
	j := func(p string) string { return JalurAdj(n, p) }
	tipe := func(v ...string) Kondisi {
		return adjK(n, func(b Baris) bool {
			for _, x := range v {
				if b["Type"] == x {
					return true
				}
			}
			return false
		})
	}
	isi := func(p string) Kondisi { return adjK(n, func(b Baris) bool { return b[p] != "" }) }
	kirimTampil := atau(
		dan(tipe("1"), terisi(CD+"Occupation"), isi("DataCommitteeTreaty.Remarks"), isi("DataCommitteeTreaty.CircumCauseOfLoss")),
		dan(tipe("2", "4"), isi("DataCommitteeTreaty.AdjusterFee"), isi("DataCommitteeTreaty.Remarks")),
		dan(tipe("3"), isi("DataCommitteeTreaty.Salvage"), isi("DataCommitteeTreaty.Remarks")),
	)
	return []Unsur{
		bagian("",
			naJika(medan("TempCommiteClaim.DateOfComitee", "Date", KTanggal), selalu),
			ro(medan("TreatyExchangeYearly.UserName", "Initial", KTeks)),
		),
		bagian("",
			tampil(label("Adjustment"), tipe("1")),
			tampil(label("Salvage"), tipe("3")),
			tampil(label("Adjuster / Consultant Fee"), tipe("2", "4")),
		),
		bagian("",
			tampil(wajibU(medan(j("DataCommitteeTreaty.CircumCauseOfLoss"), "Circumstanses", KArea)), tipe("1")),
			tampil(wajibJ(medan(j("DataCommitteeTreaty.Salvage"), "Salvage", KArea), tipe("3")), tipe("3")),
			tampil(wajibJ(medan(j("DataCommitteeTreaty.AdjusterFee"), "Adjuster / Consultant Fee", KArea), tipe("2", "4")),
				tipe("2", "4")),
			wajibU(medan(j("DataCommitteeTreaty.Remarks"), "Remarks", KArea)),
		),
		bagian("",
			tampil(bagian("",
				tampil(tombolOQ("SendClaimToCommittee", "Send Claim to Committee", OQBatasKomite), kirimTampil),
			), func(*Halaman) bool { return lolos }),
			tombol("CancelKomite", "Cancel", "TutupModal"),
		),
	}
}

// LayarTutupKlaim - Section `PreventRejectClaimProp`. Layout S22 (Remark, visible never) dan S28-29 (pesan sukses
// komite) tidak dibangun; "Yes" tutup tanpa pembayaran nonaktif (OQ-CP-06).
func LayarTutupKlaim() []Unsur {
	cwp := "TempCommiteClaim.AllocationShareSalvage"
	return []Unsur{
		medan(cwp, "Close Without Payment", KCentang),
		medan("TempCommiteClaim.CircumtansesCouseOfLoss", "Chronology", KArea),
		wajibU(medan("TempCommiteClaim.Remarks", "Remarks", KArea)),
		tampil(bagian("",
			label("Are you sure want close this claim without payment?"),
			tombolOQ("CloseWithoutPayment", "Yes", OQTutupTanpaBayar),
			tombol("CloseNo", "No", "TutupModal"),
		), sama(cwp, "true")),
		tampil(bagian("",
			label("Are you sure want close this claim?"),
			tombol("CloseYes", "Yes", "CloseClaimProp"),
			tombol("CloseNo2", "No", "TutupModal"),
		), beda(cwp, "true")),
	}
}

// LayarPLA - Section `GeneratePLA` (FlowAction GeneratePLA, post-activity TryMakePLA_Act).
func LayarPLA() []Unsur {
	return []Unsur{
		medan(CD+"Remark", "Remarks", KArea),
		tombol("SubmitPLA", "Submit", "TryMakePLA"),
		tombol("CancelPLA", "Cancel", "TutupModal"),
	}
}

// LayarDLA - Section `generateDLATreaty` (FlowAction GenerateDLATreaty, post-activity PrintDLATreatyIn).
func LayarDLA(n int) []Unsur {
	return []Unsur{
		wajibU(medan(JalurAdj(n, "RemarksDLA"), "Remarks", KArea)),
		tombol("SubmitDLA", "Submit", "PrintDLATreatyIn"),
		tombol("CancelDLA", "Cancel", "TutupModal"),
	}
}

// PesanPrintPla - Section `PrintFile` (local action sesudah "Save to issue RNM" bila IsPLA = 1). Section `PrintFileDLA`
// memakai `PesanCetakDLA`.
const PesanPrintPla = "Please Print Pla"
