package models

// Untuk apa berkas ini: DEFINISI LAYAR INPUT ESTIMASI - Section `InputEstimasiAdmin` (FlowAction InputEstimasi) dengan
// `InputEstimasiDetail` (grid objek per lini, masterDetail), panel objek (`PropertyItemListGridEstimation` /
// `...MBU` / `...Marine` / `...Travel` / `ShowItemPA`: grid item, masterDetail), panel item (`Estimasi` / `EstimasiPA` /
// `EstimasiMarine`), dan pop-up `Pla_Dtl` (Print PLA). Label, urutan, kondisi VERBATIM dari XML.
//
// Panel masterDetail: grid ber-`Rincian` membuka panel baris `KunciPanel(prefiks, daftar, n)`:
//
//	PanelObjekEst  "est"      baris ClaimData.ObjectList(o)                      -> LayarObjekEstimasi(o)
//	PanelItemEst   "estitem"  baris ClaimData.ObjectList(o).ObjectItemList(i)    -> LayarItemEstimasi(o, i)
//
// Tidak dibangun (PARITAS): grid / tombol bergerbang `pyWorkPage.ClaimData.ExGratia = 1` (Spreading Policy / Claim yang
// dapat disunting, Save Spreading) - ExGratia tingkat klaim selalu 0 (InsertObjects_dt 11); baris Total
// `InputParam.CARI41 / CARI42` (halaman requestor tanpa penulis); grid estimasi Travel di section Estimasi (Travel memakai
// EstimasiPA).
//
// ⚠️ PERBAIKAN (PARITAS `[penyimpangan sadar]`): grid item MBU ber-EditAction `PropertyItemListGridEstimationMBU_FA` yang
// merender section grid item itu SENDIRI pada baris item - estimasi MBU tidak dapat diisi dari layar Pega. Di sini baris
// item MBU membuka section `Estimasi` (seperti Fire).

// Prefiks panel masterDetail layar Input Estimasi.
const (
	PanelObjekEst = "est"
	PanelItemEst  = "estitem"
)

// kondisi estimasi
var (
	bukanTreatyIn = atau(sama(JalurIsTreatyIn, "0"), sama(JalurIsTreatyIn, ""))
	bPFC          = bSama("PrintFaceClaim", "1")
	bBukanPFC     = func(_ *Halaman, b Baris) bool { return b["PrintFaceClaim"] != "1" }
	hAdj          = func(h *Halaman, _ Baris) bool { return h.Ambil(JalurIsAdjustment) == "1" }
)

// LayarEstimasi - Section InputEstimasiAdmin.
func LayarEstimasi() []Unsur {
	return []Unsur{
		sebaris("",
			tampilIsi(ro(medan(CD+"NoClaim", "No Claim", KTeks))),
			tombolOQ("ViewStatusPaymentPremi", "View Status Payment Premi", OQBayarPremi),
			naJika(tombol("SendToPIC", "Send to PIC Claim", "SetDisable"),
				atau(beda(JalurIsPicTransfer, "1"), errorLebih1)),
		),
		letak(LetakTab,
			bagian("Estimation", tampil(gridObjekEstimasi(), bukanTreatyIn)),
			bagian("Policy Detail & Claims History", tabRiwayat()...),
			bagian("Progress Claim", tabProgres()...),
			tampil(bagian("View Registration", detailRegister()...), beda(JalurIsTreatyIn, "1")),
		),
		tampil(tombol("Back", "Back", "BackToRegister"), beda(JalurIsCFS, "1")),
	}
}

// tombolObjekEstimasi - tombol baris grid objek InputEstimasiDetail (semua lini).
func tombolObjekEstimasi() []Unsur {
	cfs := kTombol("CLaimFaceSheet", "Download Claim Face Sheet", "CLaimFaceSheet",
		func(h *Halaman, b Baris) bool { return b["CFS"] != "1" || b["PrintFaceClaim"] == "1" || errorLebih1(h) })
	pla := kTampil(kTombol("PrintPLA", "Print PLA", "BukaPLA",
		func(h *Halaman, b Baris) bool {
			return b["IsFacretro"] != "1" || b["PlaStatus"] == "1" || errorLebih1(h)
		}),
		func(h *Halaman, _ Baris) bool { return h.Ambil(CD+"ExGratia") == "0" })
	return []Unsur{cfs, pla, kTombol("OutstandingSummary", "Outstanding Summary", "BukaOutstanding", nil)}
}

// gridObjekEstimasi - Section InputEstimasiDetail: grid ObjectList per lini (masterDetail ke panel objek).
func gridObjekEstimasi() Unsur {
	g := func(k Kondisi, kolom ...Unsur) Unsur {
		return tampil(Unsur{Jenis: JenisGrid, Jalur: DaftarObjek, Bernomor: true, Rincian: PanelObjekEst,
			Kolom: append(kolom, tombolObjekEstimasi()...)}, k)
	}
	return bagian("",
		g(IsMarineCargo, kRO(kol("ObjectName", "Trading Name", KTampil)), kRO(kol("ObjectLocation", "Packing Name", KTampil)),
			kRO(kol("ObjectJob", "Goods Name", KTampil)), kRO(kol("ObjectSurveyLocation", "Conveyance Name", KTampil))),
		g(atau(IsAneka, IsFire, IsGolfInsurance), kRO(kol("ObjectName", "Object Type", KTampil)),
			kRO(kol("ObjectLocation", "Object Location", KTampil))),
		g(IsMBU, kRO(kol("BrandName", "Brand Name", KTampil)), kRO(kol("ModelName", "Model", KTampil)),
			kRO(kol("TypeName", "Type", KTampil)), kRO(kol("LicensePlate", "License Plate", KTampil))),
		g(IsPA, kRO(kol("ObjectName", "Object Name", KTampil)), kRO(kol("ObjectJob", "Job", KTampil)),
			kRO(kol("ObjectDateOfBirth", "Birth Date", KTanggal)), kRO(kol("Gender", "Gender", KTampil))),
		g(IsTravel, kRO(kol("ObjectName", "Object Name", KTampil)), kRO(kol("ObjectParticipantStatus", "Status", KTampil)),
			kRO(kol("ObjectIDCard", "ID Card", KTampil)), kRO(kol("ObjectDateOfBirth", "Birth Date", KTanggal))),
	)
}

// LayarObjekEstimasi - panel baris objek o (section menurut EditAction grid objek per lini).
func LayarObjekEstimasi(o int) []Unsur {
	daftar := DaftarItem(o)
	tambah := func(vis Kondisi) *Unsur {
		u := tombol("TambahItem", "Add", "TambahItem")
		if vis != nil {
			u = tampil(u, vis)
		}
		return ptr(ikon(u, IkonTambah))
	}
	bukanAdj := beda(JalurIsAdjustment, "1")
	hapusFire := kTampil(ikonK(kTombol("HapusItem", "", "HapusItem", bPFC), IkonHapus),
		func(h *Halaman, _ Baris) bool { return h.Ambil(JalurIsAdjustment) != "1" })
	hapusObjek := kTampil(ikonK(kTombol("HapusItem", "Delete", "HapusItem", nil), IkonHapus), bBukanPFC)
	angka := func(p, j string) Unsur { return kRO(kol(p, j, KAngka)) }
	nilaiTSI := []Unsur{angka("KursObjectItem", "Currency Value (IDR)"), angka("TSIPerObject", "TSI Object"),
		angka("TSINusare", "TSI RNM"), angka("ValueTSINusareIDR", "TSI RNM in IDR")}
	okupasi := kNA(kROJ(kAksi(kSumber(kol("OccupationName", "Occupation", KOtomatis), SumberOkupasi), "PilihOkupasi"),
		bPFC), hAdj)
	aneka := kROJ(kAksi(kSumber(kol("IndexAneka", "Object Id", KOtomatis), SumberAneka), "PilihAneka"), bPFC)
	covAneka := kROJ(kAksi(kSumber(kol("CoverageNote", "Coverage Name", KOtomatis), SumberCovAneka),
		"PilihCoverageAneka"), bPFC)
	g := func(k Kondisi, tb *Unsur, kolom ...Unsur) Unsur {
		return tampil(Unsur{Jenis: JenisGrid, Jalur: daftar, Rincian: PanelItemEst, Tambah: tb, Kolom: kolom}, k)
	}
	covObjek := kROJ(kAksi(kSumber(kol("CoverageOLDID", "Coverage ID", KOtomatis), SumberCovObjek), "PilihCoverageObjek"),
		bPFC)
	return []Unsur{tampil(bagian("",
		g(IsAneka, tambah(bukanAdj), append([]Unsur{aneka, kRO(kol("ObjectItemName", "Object Name", KTampil)), okupasi,
			covAneka, kRO(kol("Currency", "Currency", KTampil))}, append(nilaiTSI, hapusFire)...)...),
		g(IsGolfInsurance, tambah(bukanAdj), append([]Unsur{aneka, kRO(kol("ObjectItemName", "Object Name", KTampil)),
			covAneka, kRO(kol("Currency", "Currency", KTampil))}, append(nilaiTSI, hapusFire)...)...),
		g(IsFire, tambah(bukanAdj),
			kRO(kol(PropNomorBaris, "Object No", KTampil)),
			kNA(kROJ(kAksi(kSumber(kol("ObjectItemName", "Object Name", KOtomatis), SumberItemProp), "SetObjectItem"),
				bPFC), hAdj),
			okupasi,
			kAksi(kSumber(kol("CoverageOLDID", "Coverage Name", KOtomatis), SumberCovFire), "PilihCoverageFire"),
			kRO(kol("Currency", "Currency", KTampil)),
			kol("KursObjectItem", "Currency Value(IDR)", KAngka),
			angka("TSIPerObject", "TSI Object"), angka("TSINusare", "TSI RNM"),
			kTampil(kol("ValueTSINusareIDR", "TSI RNM in IDR", KAngka), bTerisi("ValueTSINusareIDR")),
			angka("LimitofLiability", "Limit Of Liability"),
			hapusFire),
		g(IsMBU, tambah(nil), append([]Unsur{kRO(kol(PropNomorBaris, "Coverage No", KTampil)), covObjek,
			kRO(kol("CoverageNote", "Coverage Name", KTampil)), kRO(kol("Currency", "Currency", KTampil))},
			append(nilaiTSI, hapusObjek)...)...),
		g(atau(IsTravel, IsPA), tambah(nil), append([]Unsur{kRO(kol(PropNomorBaris, "Coverage No", KTampil)),
			covObjek, kRO(kol("CoverageNote", "Coverage Name", KTampil)), kRO(kol("Currency", "Currency", KTampil))},
			append(nilaiTSI, hapusObjek)...)...),
		g(IsMarineCargo, nil, kRO(kol(PropKlikMarine, "Click Here to View Coverage and Input Estimation", KTampil))),
	), bukanTreatyIn)}
}

// Properti tampilan baris grid item (tidak disimpan).
const (
	// PropNomorBaris - nomor baris Pega (`.pxListSubscript`, layar menampilkan nomor baris). Dinamai lain karena literal
	// nama Pega itu dilarang di kode oleh penjaga batas komite Claim Life (indeks posisi bukan kunci rujukan, AC 62);
	// medan ini hanya tampilan, dihitung ulang setiap muat.
	PropNomorBaris = "NomorBaris"
	// PropKlikMarine - label tetap baris grid item Marine Cargo.
	PropKlikMarine = "KlikMarine"
)

// LayarItemEstimasi - panel baris item (o, i): Estimasi (Fire / Aneka / Golf / MBU), EstimasiPA (PA / Travel),
// EstimasiMarine (Marine Cargo).
func LayarItemEstimasi(o, i int) []Unsur {
	p := JalurItem(o, i) + "."
	m := func(prop string) string { return p + prop }
	itemSama := func(prop, v string) Kondisi { return func(h *Halaman) bool { return AmbilJalur(h, m(prop)) == v } }
	adj := isAdjustment
	est := DaftarDiItem(o, i, AnakEstimasi)
	tglEst := kNA(kAksi(kol("EstimationDate", "Estimation Date", KTanggal), "ProtectionDate"), bPFC)
	mu := func(na KondisiBaris) Unsur {
		return kNA(kROJ(kAksi(kSumber(kol("CurrencyID", "Currency", KPilih), SumberMataUang),
			"SetConvertValueKurs_Estimation"), bPFC), na)
	}
	gross := kNA(kROJ(kAksi(kol("GrossEstimationPct", "Estimation Gross (100%)", KAngka), "CheckEstimateValue"), bPFC), bPFC)
	hapusEst := kTampil(ikonK(kTombol("DeleteValueEstimation", "Delete", "DeleteValueEstimation", nil), IkonHapus),
		bBukanPFC)
	tambahEst := func(na Kondisi) *Unsur {
		u := ikon(tombol("ValidateInputEstimate", "Add", "ValidateInputEstimate"), IkonTambah)
		if na != nil {
			u = naJika(u, na)
		}
		return ptr(u)
	}
	totalEst := tampil(bagian("",
		ro(medan(m("TotalGrossEstimasi"), "Total Estimation Gross", KAngka)),
		ro(medan(m("TotalGrossEstimasiIDR"), "Total Estimation Gross in IDR", KAngka)),
		ro(medan(m("TotalClaimSpreaded"), "Total Estimation RNM", KAngka)),
		ro(medan(m("TotalEstimationValueinIDR"), "Total Estimation RNM in IDR", KAngka)),
	), func(h *Halaman) bool { return !IsTravel(h) })
	spreadPolis := func(premi bool) Unsur {
		kolom := []Unsur{kRO(kol(PropNomorBaris, "Spreading No", KTampil)),
			kRO(kSumber(kol("TreatyType", "Treaty Type", KPilih), SumberJenisReas)),
			kRO(kol("SharePercentage", "Share %", KTampil)), kRO(kol("TSISpreaded", "TSI Spreaded", KAngka))}
		if premi {
			kolom = append(kolom, kRO(kol("PremiumSpreaded", "Premium Spreaded", KAngka)))
		}
		return tampil(bagian("Spreading Policy",
			Unsur{Jenis: JenisGrid, Jalur: DaftarDiItem(o, i, AnakSpreadPolis), Kolom: kolom},
			ro(medan(m("TotalSharePercentage"), "Total Share Percentage", KTampil)),
			ro(medan(m("TotalTSISpreaded"), "Total TSI Spreaded", KAngka)),
		), dan(func(h *Halaman) bool { return !IsTravel(h) }, sama(CD+"ExGratia", "0")))
	}
	spreadKlaim := func(retro bool) Unsur {
		kolom := []Unsur{kRO(kol(PropNomorBaris, "Spreading No", KTampil)), kRO(kol("Currency", "Currency", KTampil)),
			kRO(kSumber(kol("TreatyType", "Treaty Type", KPilih), SumberJenisReas)),
			kRO(kol("SharePercentage", "Share %", KTampil)), kRO(kol("ClaimSpreaded", "Claim Spreaded", KAngka))}
		if retro {
			kolom = append(kolom, kTampil(kTombol("ShowRetro", "View Retro", "BukaRetro", nil),
				bSama("TreatyType", TreatyFacRetro)))
		}
		return tampil(bagian("Spreading Claim",
			Unsur{Jenis: JenisGrid, Jalur: DaftarDiItem(o, i, AnakSpreadKlaim), Kolom: kolom},
			ro(medan(m("TotalSharePercentage"), "Total Share Percentage", KTampil)),
			tampil(ro(medan(m("TotalClaimSpreaded"), "Total Claim Spreaded", KAngka)),
				func(h *Halaman) bool {
					return AmbilJalur(h, m("TotalClaimSpreaded")) == AmbilJalur(h, m("TotalEstimasi"))
				}),
		), dan(func(h *Halaman) bool { return !IsTravel(h) }, sama(CD+"ExGratia", "0")))
	}
	deductible := bagian("Deductible Info",
		sebaris("",
			naJika(medan(m("DeductibleType"), "", KCentang), adj),
			tampil(naJika(sumber(medan(m("FormType"), "Format", KPilih), kode("FormType")), adj),
				itemSama("DeductibleType", "true")),
		),
		tampil(sebaris("",
			naJika(sumber(medan(m("CurrencyDeductible"), "Currency", KPilih), SumberMataUang), adj),
			naJika(medan(m("NetDeductibleValue"), "Amount", KAngka), adj),
		), itemSama("FormType", "1")),
		tampil(sebaris("",
			aksi(naJika(medan(m("Amount"), "%", KAngka), adj), "CountTSI"),
			label("of"),
			aksi(naJika(sumber(medan(m("TypeDeductible"), "", KPilih), kode("TypeDeductible")), adj), "CountTSI"),
			tampil(ro(medan(m("TSIDeductible"), "TSI Amount", KAngka)), itemSama("TypeDeductible", "2")),
			tampil(naJika(medan(m("ClaimDeductible"), "Claim Amount", KAngka), adj), itemSama("TypeDeductible", "1")),
			label("Minimum"),
			naJika(roJika(sumber(medan(m("CurrencyDeductible"), "Currency", KPilih), SumberMataUang),
				itemSama("TypeDeductible", "2")), adj),
			aksi(naJika(medan(m("DeductibleValue"), "Amount", KAngka), adj), "CountTSI"),
			ro(medan(m("NetDeductibleValue"), "Deductible Value", KAngka)),
		), itemSama("FormType", "2")),
	)
	estimasi := tampil(Unsur{Jenis: JenisGrid, Jalur: est, Tambah: tambahEst(itemSama("CoverageNote", "")), Kolom: []Unsur{
		kRO(kol(PropNomorBaris, "Estimation No", KTampil)), tglEst, mu(bPFC),
		kRO(kol("KursValue", "Currency Value (IDR)", KAngka)), kRO(kol("PersenRNM", "Share RNM %", KAngka)), gross,
		kRO(kol("Deductible", "Deductible", KAngka)), kRO(kol("NetEstimationValue", "Net Estimation Gross (100%)", KAngka)),
		kRO(kol("EstimationValue", "Net Estimation RNM", KAngka)), kRO(kol("ConvertValue", "Net Estimation Value in IDR", KAngka)),
		hapusEst,
	}}, func(h *Halaman) bool { return !IsTravel(h) })
	estimasiPA := Unsur{Jenis: JenisGrid, Jalur: est, Tambah: tambahEst(nil), Kolom: []Unsur{
		kRO(kol(PropNomorBaris, "Estimation No", KTampil)),
		kROJ(kAksi(kol("EstimationDate", "Estimation Date", KTanggal), "ProtectionDate"), bPFC),
		mu(func(h *Halaman, _ Baris) bool { return h.Ambil(JalurIsAdjustment) == "1" }),
		kRO(kol("KursValue", "Currency Value (IDR)", KAngka)), kRO(kol("PersenRNM", "Share RNM %", KAngka)),
		kROJ(kAksi(kol("GrossEstimationPct", "Estimation Gross (100%)", KAngka), "CheckEstimateValue"), bPFC),
		kRO(kol("EstimationValue", "Estimation Value", KAngka)), kRO(kol("ConvertValue", "Estimation Value in IDR", KAngka)),
		hapusEst,
	}}
	estimasiMarine := tampil(Unsur{Jenis: JenisGrid, Jalur: est, Tambah: tambahEst(nil), Kolom: []Unsur{
		kNA(kAksi(kol("EstimasiMoreThanTSI", "Estimasi > TSI", KCentang), "CheckEstimateValue"), bPFC),
		kROJ(kAksi(kol("EstimationDate", "Estimation Date", KTanggal), "ProtectionDate"), bPFC),
		kROJ(kAksi(kSumber(kol("CurrencyID", "Currency", KPilih), SumberMataUang), "SetConvertValueKurs_Estimation"), bPFC),
		kRO(kol("KursValue", "Currency Value (IDR)", KAngka)), kRO(kol("PersenRNM", "Share RNM %", KAngka)),
		kROJ(kAksi(kol("GrossEstimationPct", "Gross Estimation (%)", KAngka), "CheckEstimateValue"), bPFC),
		kRO(kol("EstimationValue", "Estimation Value", KAngka)), kRO(kol("ConvertValue", "Estimation Value in IDR", KAngka)),
		hapusEst,
	}, Kaki: []Unsur{sebaris("Total", ro(medan(m("TotalClaimSpreaded"), "", KAngka)),
		ro(medan(m("TotalEstimationValueinIDR"), "", KAngka)))}}, func(h *Halaman) bool { return !IsTravel(h) })
	coverageMarine := Unsur{Jenis: JenisGrid, Jalur: DaftarDiItem(o, i, "CoverageList"), Kolom: []Unsur{
		kRO(kol("CoverageNote", "Coverage Name", KTampil)), kRO(kol("Currency.Name", "Currency", KTampil)),
		kRO(kol("Currency.KURS", "Currency Value (IDR)", KAngka)), kRO(kol("TSILiability", "TSI Object", KAngka)),
		kRO(kol("TSINusantaraRe", "TSI RNM", KAngka)), kRO(kol("TSIinIDR", "TSI RNM in IDR", KAngka)),
	}}
	return []Unsur{
		tampil(bagian("", deductible, estimasi, totalEst, spreadPolis(false), spreadKlaim(true)),
			atau(IsFire, IsAneka, IsGolfInsurance, IsMBU)),
		tampil(bagian("", estimasiPA, totalEst, spreadKlaim(false), spreadPolis(false)), atau(IsPA, IsTravel)),
		tampil(bagian("", coverageMarine, estimasiMarine, totalEst, spreadPolis(true), spreadKlaim(false)), IsMarineCargo),
	}
}

// LayarPLA - Section Pla_Dtl (harness Pla_Dtl_Harness) untuk objek o.
//
// Grid penerima `Email.RecipientList` dan subjek tidak dibangun: surel PLA berjalan lewat outbox produksi yang muatannya
// hanya pengenal (OQ-CFI-22). `Email.Message` hanya-baca bila objek `.IsEditable != 1` - properti tanpa penulis di
// korpus, sehingga Submit ber-"Send Email" selalu nonaktif (sama dengan Pega).
func LayarPLA(o int) []Unsur {
	p := JalurObjek(o) + "."
	pesanKosong := sama(JalurPesanSurel, "")
	kirim := sama(JalurKirimSurel, "true")
	return []Unsur{
		medan(JalurKirimSurel, "Send Email", KCentang),
		tampil(roJika(medan(JalurPesanSurel, "Message", KArea),
			func(h *Halaman) bool { return AmbilJalur(h, p+"IsEditable") != "1" }), kirim),
		sumber(medan(p+"ShareRetro", "Share Retro", KRadio), kode("ShareRetro")),
		wajibU(medan(p+"RemarksPLA", "Remarks", KArea)),
		naJika(tombol("GeneratePLA", "Submit", "GeneratePLA"), dan(kirim, pesanKosong)),
	}
}
