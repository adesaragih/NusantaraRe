// Package services memuat aturan endorsement Fac In.
//
// Untuk apa berkas ini: SEAM 4 (K-047) - menyiapkan "keadaan sebelum" saat
// kasus endorsement lahir. Tiket E06 (tracer) dan E07 (54 salinan lapis A).
//
// Dibaca sesudah: models/penawaran.go. Dibaca sebelum: lapisc.go,
// nilailama.go, porsiperiode.go.
//
// Asal (Pega): `Endorsment Fac In/Activity/SetValueToEDMWork.xml` langkah 14
// (13 sub-langkah) dan 15. Langkah 1-13 milik alur masuk (E03) dan langkah 16
// milik modul selisih (E13) - keduanya BUKAN isi berkas ini.
//
// ⛔ "Before-image" adalah nama konsep, bukan properti. Tiga lapisnya tetap
// tiga (spec §1):
//   - lapis A: `OfferFacIn.OldData` - dokumen polis lama utuh (berkas ini);
//   - lapis B: properti `*Old` per baris (nilailama.go);
//   - lapis C: penanda baris warisan (lapisc.go).
//
// ⛔ Predikat lini bisnis (`IsFire`, `IsMBU`, …) dan siklus (`IsEDM`) DITERIMA
// SEBAGAI MASUKAN, tidak dievaluasi di sini. Registry predikatnya milik
// tiket E01/NB-08 yang belum ada. Pilihan pemegang modul pada sesi
// implementasi 01-10-2026 (bukan keputusan work owner): irisan before-image
// dikerjakan tanpa registry lokal, supaya mesin NB tidak terduplikasi.
package services

import (
	"nusantarare/modul/endorsmentfacin/backend/models"
)

// Predikat - hasil rule `When` yang dipakai before-image, dievaluasi pemanggil.
//
// Tiap predikat berdiri SENDIRI, persis seperti gerbang Pega yang menguji
// satu per satu. Tidak diasumsikan saling eksklusif: bila dua lini bernilai
// benar, kedua cabang berjalan berurutan, seperti di Pega.
type Predikat struct {
	IsFire          bool // When/IsFire
	IsGolfInsurance bool // When/IsGolfInsurance
	IsAneka         bool // When/IsAneka
	IsPA            bool // When/IsPA
	IsMarineCargo   bool // When/IsMarineCargo
	IsMBU           bool // When/IsMBU
	IsTravel        bool // When/IsTravel
	IsLife          bool // When/IsLife
	// IsEDM - `pyWorkPage.OfferFacIn.QuotationData.StatusBusiness = 3`
	// (When/IsEDM). ⛔ Bukan kebalikan `IsNotEDM`, yang membaca
	// `pyWorkPage.Quotation.StatusBusiness` - halaman lain (E02).
	IsEDM bool
}

// MasukanBeforeImage - keadaan tepat sebelum langkah 14.
type MasukanBeforeImage struct {
	// Kasus - `newWorkPage` sesudah langkah 1-13 (alur masuk, E03).
	Kasus models.KasusEndorsement
	// PolisLama - hasil `RDBList/GetEDMOldData_SQL` yang diadopsi ke halaman
	// (14.1-14.2): versi polis terakhir, `PRODKE = COUNT(NOPOLIS)-1`.
	// nil = query tanpa baris; OldData tetap ada sebagai halaman kosong.
	PolisLama *models.OfferFacIn
	// TeamGroupMarketing - `ListMkt.pxResults(1).TeamGroup` dari master
	// marketing (14.4). Kosong = tidak ada nilai.
	TeamGroupMarketing string
	Predikat           Predikat
}

// HasilBeforeImage - objek kerja endorsement sesudah langkah 14-15, beserta
// dua tulisan ke halaman DI LUAR objek kerja itu.
type HasilBeforeImage struct {
	Kasus models.KasusEndorsement
	// PortalQuotationBusinessType - `pyWorkPage.Quotation.BusinessType`:
	// di SetValueToEDMWork, pyWorkPage adalah kasus portal, bukan kasus
	// endorsement. Penulisannya milik pemanggil (E03).
	PortalQuotationBusinessType string
	// InputDataCreditCARI6 - `InputDataCredit.CARI6` = nama marketing polis
	// lama. Halamannya bukan milik modul ini; penulisannya milik pemanggil.
	InputDataCreditCARI6 string
	// GalatPorsiPeriode - langkah 15 tidak dapat dihitung tanpa menebak
	// (`ErrSatuanSelisihWaktuBelumTerverifikasi`,
	// `ErrTanggalPorsiPeriodeKosong`). Bila tidak nil, `ProrateStartEDM` dan
	// `ProrateEDMEnd` KOSONG; sisa agregat tetap sah.
	GalatPorsiPeriode error
}

// PrepareBeforeImage - SEAM 4. Lapis A (14.1-14.5), lapis C (14.7-14.13),
// porsi periode (15), lalu lapis B untuk pembukaan layar pertama.
//
// Masukan tidak diubah: hasilnya salinan.
//
// ⚠️ Lapis B (`IsiNilaiLama`) di Pega berjalan di pre-activity layar, SETIAP
// kali layar endorsement dibuka - bukan di activity ini. Ia dipanggil di
// sini supaya agregat yang keluar sudah seperti yang dilihat pengguna saat
// pertama membuka; pembukaan berikutnya memanggil `IsiNilaiLama` sendiri.
//
// 14.6 `Obj-Save` tidak dilakukan: penyimpanan milik repository (E17).
func PrepareBeforeImage(m MasukanBeforeImage) (HasilBeforeImage, error) {
	k := m.Kasus.Klon()
	o := &k.OfferFacIn

	// 14.1-14.2 - muat dokumen polis lama ke OldData. 14.1 ber-
	// `pyStepsPreCondition=false`: tetap jalan tanpa syarat (P-11 tertutup,
	// 19-09-2026).
	lama := models.OfferFacIn{}
	if m.PolisLama != nil {
		lama = m.PolisLama.Klon()
	}
	o.OldData = &lama

	portalBusinessType, cari6 := salinLapisA(&k)

	// 14.5 - TeamGroup dari master marketing, hanya bila terisi
	// (`@PropertyHasValue(ListMkt.pxResults(1).TeamGroup)`).
	if m.TeamGroupMarketing != "" {
		o.QuotationData.TeamGroup = m.TeamGroupMarketing
		k.Quotation.TeamGroup = m.TeamGroupMarketing
	}

	// 14.7-14.13 - lapis C.
	if err := tandaiBarisWarisan(&k, m.Predikat); err != nil {
		return HasilBeforeImage{}, err
	}

	// 15 - porsi periode. Galatnya TIDAK menggagalkan lapis A/B/C: nilainya
	// sementara (K-048), dan kedua rasio dibiarkan kosong - tidak ditebak.
	galatPorsi := hitungPorsiPeriode(o)

	if err := IsiNilaiLama(&k, m.Predikat); err != nil {
		return HasilBeforeImage{}, err
	}

	return HasilBeforeImage{
		Kasus:                       k,
		PortalQuotationBusinessType: portalBusinessType,
		InputDataCreditCARI6:        cari6,
		GalatPorsiPeriode:           galatPorsi,
	}, nil
}

// salinLapisA - langkah 14.3: 54 penugasan, URUTAN KORPUS.
//
// [terverifikasi] 54 = 51 dari OldData + 3 bukan. Dari 51: 45 berpola
// `OfferFacIn.X = OldData.X`, 6 menulis ke halaman lain (ditandai CERMIN).
// Arsip lintas-siklus menyebut ke-54 seragam; itu keliru untuk 9 di antaranya
// (spec §3) - mengikutinya meninggalkan empat halaman kosong.
//
// Mengembalikan dua tulisan ke halaman di luar objek kerja:
// `pyWorkPage.Quotation.BusinessType` dan `InputDataCredit.CARI6`.
func salinLapisA(k *models.KasusEndorsement) (portalBusinessType, cari6 string) {
	o := &k.OfferFacIn
	lama := o.OldData
	q, qLama := &o.QuotationData, &lama.QuotationData

	q.BusinessOldId = qLama.BusinessOldId
	q.BusinessCode = qLama.BusinessCode
	q.BusinessFac = qLama.BusinessFac
	q.BusinessName = qLama.BusinessName
	q.BusinessType = qLama.BusinessType
	q.CedingCo = qLama.CedingCo
	q.CedingCoName = qLama.CedingCoName
	q.SourceOfBusiness = qLama.SourceOfBusiness
	q.SobName = qLama.SobName
	q.MarketingCode = qLama.MarketingCode
	q.MarketingName = qLama.MarketingName
	q.MOID = qLama.MOID
	q.TeamGroup = qLama.TeamGroup
	o.PolicyData.StartDateTime = lama.PolicyData.StartDateTime
	o.PolicyData.EndDateTime = lama.PolicyData.EndDateTime
	// CERMIN 1 (L3135) - newWorkPage.Quotation.BusinessType.
	k.Quotation.BusinessType = qLama.BusinessType
	q.InsuredName = qLama.InsuredName
	q.InsuredID = qLama.InsuredID
	o.PolicyData.OfferingDate = lama.PolicyData.OfferingDate
	q.NoOfferSlip = qLama.NoOfferSlip
	o.Currency = lama.Currency
	q.IsGroup = qLama.IsGroup
	q.QQName = qLama.QQName
	o.CurrencyList = models.KlonBarisMataUang(lama.CurrencyList)
	o.Parameters = models.KlonJSON(lama.Parameters)
	o.InwardScale = lama.InwardScale
	o.PercentShare = lama.PercentShare
	o.AdditionalCapital = lama.AdditionalCapital
	o.MaxPctTreatyCapacity = lama.MaxPctTreatyCapacity
	o.MaxTreatyCapacity = lama.MaxTreatyCapacity
	o.CurrentYear = lama.CurrentYear
	o.ProRatePercent = lama.ProRatePercent
	// CERMIN 2-4 (L3490, L3531, L3572) - tiga field pembayaran ditulis DUA
	// kali: ke newWorkPage.Policy.Payment dan ke OfferFacIn.PolicyData.Payment.
	k.Policy.Payment.Installment = lama.PolicyData.Payment.Installment
	o.PolicyData.Payment.Installment = lama.PolicyData.Payment.Installment
	k.Policy.Payment.RICommision = lama.PolicyData.Payment.RICommision
	o.PolicyData.Payment.RICommision = lama.PolicyData.Payment.RICommision
	k.Policy.Payment.PctBrokerageFee = lama.PolicyData.Payment.PctBrokerageFee
	o.PolicyData.Payment.PctBrokerageFee = lama.PolicyData.Payment.PctBrokerageFee
	o.PolicyData.ProdDateTime = lama.PolicyData.ProdDateTime
	// BUKAN dari OldData 1 (L3632) - pyWorkPage.Quotation.BusinessType.
	portalBusinessType = k.Quotation.BusinessType
	q.PolicyType = qLama.PolicyType
	o.CedingCedantList = models.KlonCedant(lama.CedingCedantList)
	o.ShareCedantType = lama.ShareCedantType
	o.IsSpecialAcceptance = lama.IsSpecialAcceptance
	// BUKAN dari OldData 2 (L3733) - newWorkPage.Quotation menjadi salinan
	// halaman OfferFacIn.QuotationData (yang barusan terisi).
	// [dugaan] Penugasan halaman ke halaman di Pega MENGGANTI isi tujuan
	// seluruhnya; bila Pega menggabungkan (properti tujuan yang tidak ada di
	// sumber tetap hidup), hasilnya berbeda. Belum terverifikasi.
	k.Quotation = o.QuotationData
	o.ProRateType = lama.ProRateType
	o.IsB2B = lama.IsB2B
	// CERMIN 5 (L3793) - akar objek kerja.
	k.IsB2B = lama.IsB2B
	// CERMIN 6 (L3813) - halaman parameter kredit.
	cari6 = qLama.MarketingName
	q.EDMDay = qLama.EDMDay
	// Salinan halaman Quotation (L3733) terjadi SEBELUM EDMDay disalin - urutan
	// korpus. newWorkPage.Quotation.EDMDay karena itu memegang nilai
	// OfferFacIn.QuotationData.EDMDay dari langkah 1-13, bukan dari OldData.
	o.BinderRNM = lama.BinderRNM
	o.PPnCheck = lama.PPnCheck
	o.PolicyMasterNumber = lama.PolicyMasterNumber
	// BUKAN dari OldData 3 (L3914).
	k.PPnCheck = o.PPnCheck
	return portalBusinessType, cari6
}
