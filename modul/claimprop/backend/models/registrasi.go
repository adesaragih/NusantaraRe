package models

// Untuk apa berkas ini: PORT ACTIVITY REGISTRASI KLAIM - pra-proses layar, validasi tanggal, polis, master treaty,
// lokasi, pelapor, adjuster/consultant, katastrofe, cause of loss (tiket 01-04; US 1-9; AC 96-122).
//
// Setiap fungsi menyebut berkas Activity aslinya; nomor langkah di komentar = nomor langkah Pega. Langkah berlabel
// `//` TIDAK ditulis (aturan induk spec). Langkah berprakondisi nonaktif ditulis TANPA syaratnya (aturan induk).

import (
	"fmt"
	"strings"
)

// ---------------------------------------------------------------- pesan

// Pesan validasi - teks VERBATIM korpus (salah ketik ikut).
const (
	PesanPolisTakAda        = "policy number not found, please re-choose policy number"
	PesanBelumAdaRNM        = "Belum ada Nomor RNM untuk master ini"
	PesanSudahAdaRNM        = "Sudah ada Nomor RNM untuk master ini"
	PesanMasterRevisi       = "ID MASTER Sedang dalam tahap Revisi, Tidak bisa Claim"
	PesanIsiEstimasi        = "Please fill Estimation list"
	PesanIsiSpreading       = "please Fill SpreadingList"
	PesanIsiTypeEstimasi    = "Please fill Type pada estimation list"
	PesanIsiShareSpreading  = "Please fill Share pada spreading"
	PesanIsiInterest        = "Please fill intrest fisrt"
	PesanIsiCauseOfLoss     = "Please fill Cause of loss"
	PesanIsiMaster          = "Please fill ID master treaty"
	PesanIsiLokasi          = "Please Input Loss Location"
	PesanIsiPolis           = "Please Insert Policy No"
	PesanDOLHariIni         = "Date of Loss should not be more than todays date"
	PesanDOLTerima          = "Date Of Loss should not be more than Received Date"
	PesanDOLLapor           = "Date Of Loss should not be more than Report Date"
	PesanDOLMulaiPolis      = "Date of Loss should not be more than Start Date Policy  "
	PesanDOLAkhirPolis      = "Date of Loss should not be more than End Date Policy "
	PesanDOLSerupa          = "There is a similar Date of Loss to this policy number"
	PesanIsiMulaiPolis      = "Please Fill Policy Start Date First"
	PesanTerimaHariIni      = "Date Received should not be more than todays date"
	PesanTerimaDOL          = "Date Received  should not be less than Date of Loss"
	PesanTerimaLapor        = "Date Received should not be less than Report Date"
	PesanLaporHariIni       = "ReportDate should not be more than todays date"
	PesanLaporDOL           = "ReportDate should not be less than Date of Loss"
	PesanLaporTerima        = "ReportDate should not be more than Received Date"
	PesanPeriodePolis       = "PeriodPolicy not match with period treaty"
	PesanFormatPolis        = "Format Policy No wrong, please fill the right one"
	TeksPilihMaster         = "Choose Master ID"
	TeksRegistrasi          = "Input Data Registration Claim"
	StatusMasterTuntas      = "Resolve Complete"
	AwalanPolisTreaty       = "RNM-Q"
	NilaiNonKatastrofeKlaim = "Claim"
	NilaiKatastrofe         = "Catastrophe"
)

// ---------------------------------------------------------------- pra-proses

// CheeckNoRNM meniru `Activity/CheeckNoRNM_Act.xml` - pra-proses FlowAction OutstandingClaim.
func CheeckNoRNM(k *Konteks, h *Halaman) error {
	h.BersihkanPesan() // 1 Page-Clear-Messages
	// 2 Property-Remove ClaimData.Attachment - halaman kerja lampiran tidak disimpan (STRUKTUR J1).
	// 4 RDB CekPolicyNumber_SQL bila PolicyNo terisi; 5 pesan bila hasilnya NOL baris - dievaluasi juga saat
	// langkah 4 dilewati (PolicyNo kosong), sehingga pesan muncul pula untuk polis kosong (AC 111, ditiru).
	ada := false
	if no := h.Ambil(CD + "PolicyData.PolicyNo"); no != "" {
		var err error
		if ada, err = k.Acuan.AdaPolisRealisasi(k.Ctxt(), no); err != nil {
			return err
		}
	}
	if !ada {
		h.TambahPesan("", PesanPolisTakAda)
	}
	pesanRealisasi(h) // 6-7
	// 8-10 TreatyGroupID dari TREATYBUSINESS bila kosong (DataF.CARI11 - tahun - tidak disetel activity ini).
	if h.Ambil(CD+"TreatyGroupID") == "" {
		id, err := k.Acuan.TreatyGroupBisnis(k.Ctxt(), h.Ambil(CD+"QuotationData.BusinessCode"), "")
		if err != nil {
			return err
		}
		h.Setel(CD+"TreatyGroupID", id)
	}
	CheckAnyAcceptationProp(h) // 11
	ProteksiData(h)            // 12
	// 13 SethistoryKlaimTreaty - urutan dan label riwayat diterapkan saat tampil (`TampilRiwayat`).
	return SetMOClaimTreaty(k, h) // 14
}

// pesanRealisasi - `.Message` "Belum/Sudah ada Nomor RNM untuk master ini" menurut IsRealisation (CheeckNoRNM_Act
// langkah 6-7, CekPolisAvailable_Act langkah 4-5).
func pesanRealisasi(h *Halaman) {
	if h.Ambil(CD+"IDMaster") == "" {
		return
	}
	switch h.Ambil("IsRealisation") {
	case "0":
		h.Setel("Message", PesanBelumAdaRNM)
	case "1":
		h.Setel("Message", PesanSudahAdaRNM)
	}
}

// GetPICAdjutment meniru `Activity/GetPICAdjutment_Act.xml` - pra-proses FlowAction InputAcceptation.
func GetPICAdjutment(k *Konteks, h *Halaman) error {
	CheckAnyAcceptationProp(h)    // 3
	return SetMOClaimTreaty(k, h) // 5 (4 SethistoryKlaimTreaty = TampilRiwayat)
}

// CheckAnyAcceptationProp meniru `Activity/CheckAnyAcceptationProp.xml`: IsAnyAcceptation = 1 bila ada baris
// adjustment ber-AcceptanceStatus 1.
func CheckAnyAcceptationProp(h *Halaman) {
	n := 0
	for _, b := range h.AmbilDaftar(DaftarAdjustment) {
		if b["AcceptanceStatus"] == "1" {
			n++
		}
	}
	if n > 0 {
		h.Setel("IsAnyAcceptation", "1")
	} else {
		h.Setel("IsAnyAcceptation", "0")
	}
}

// SetMOClaimTreaty meniru `Activity/SetMOClaimTreaty_Act.xml`.
//
// Langkah 1 Exit-Activity: baris 1 `PolicyNo==""` T=5 (Skip Whens -> langkah JALAN = keluar), baris 2
// `MarketingData.BranchDetailID==""` T=3 F=2. Dibaca berantai: keluar bila polis kosong ATAU BranchDetailID TERISI.
func SetMOClaimTreaty(k *Konteks, h *Halaman) error {
	if h.Ambil(CD+"PolicyData.PolicyNo") == "" || h.Ambil(CD+"MarketingData.BranchDetailID") != "" {
		return nil
	}
	mo, ada, err := k.Acuan.MarketingPolis(k.Ctxt(), h.Ambil(CD+"PolicyData.PolicyNo")) // 3 (prakondisi nonaktif)
	if err != nil || !ada {
		return err
	}
	for i, p := range []string{"ID", "ClientID", "ClientName", "TeamGroup", "BranchDetailID", "BranchDetailName"} {
		h.Setel(CD+"MarketingData."+p, mo[i]) // 4
	}
	return nil
}

// ProteksiData meniru `Activity/ProteksiData_act.xml` - pesan kelengkapan klaim (post-processing FlowAction
// OutstandingClaim, CheeckNoRNM_Act langkah 12, SaveOutstanding_Act langkah 1). Langkah 18-19 ber-remark.
func ProteksiData(h *Halaman) {
	est := h.AmbilDaftar(DaftarEstimasi)
	spr := h.AmbilDaftar(DaftarSpreading)
	// 3 baris `@contains(StatusAkseptasi,"Resolve Complete")` T=3 F=(kosong): pesan bila TIDAK memuat.
	if !strings.Contains(h.Ambil(TM+"StatusAkseptasi"), StatusMasterTuntas) {
		h.TambahPesan("", PesanMasterRevisi)
	}
	if len(est) < 1 { // 4
		h.TambahPesan("", PesanIsiEstimasi)
	}
	if len(spr) < 1 { // 5
		h.TambahPesan("", PesanIsiSpreading)
	}
	for _, b := range est { // 6-7
		if b["Type"] == "" {
			h.TambahPesan("", PesanIsiTypeEstimasi)
			break
		}
	}
	for _, b := range spr { // 9-10
		if b["SharePercentage"] == "" {
			h.TambahPesan("", PesanIsiShareSpreading)
			break
		}
	}
	for _, b := range h.AmbilDaftar(DaftarInterest) { // 12-13
		if b["ObjectName"] == "" || b["CurrencyID"] == "" || b["TSIPerObject"] == "" {
			h.TambahPesan("", PesanIsiInterest)
			break
		}
	}
	if h.Ambil(CD+"CauseOfLoss") == "" { // 14
		h.TambahPesan("", PesanIsiCauseOfLoss)
	}
	if h.Ambil(CD+"IDMaster") == "" { // 15
		h.TambahPesan("", PesanIsiMaster)
	}
	if h.Ambil(CD+"Location") == "" { // 16
		h.TambahPesan("", PesanIsiLokasi)
	}
	if h.Ambil(CD+"PolicyData.PolicyNo") == "" { // 17
		h.TambahPesan("", PesanIsiPolis)
	}
}

// ---------------------------------------------------------------- tanggal

// CheckDateDOL meniru `Activity/CheckDateDOL_Act.xml` (change Date of Loss).
//
// ⚠️ PENYIMPANGAN SADAR AC 102 / §15a (keputusan work owner): duplikat Date of Loss MEMBLOKIR. Di Pega langkah 25
// menyetel `.IsError = 1` yang tidak pernah lolos ambang `>1`; di sini penanda blokir terpisah `PropBlokirDOL`
// (§15 - IsError dipecah). Pengecualian klaim lama yang sudah ditolak tetap berlaku (langkah 23.1, AC 103).
// ⚠️ PERBAIKAN yang menyertainya (dicatat sebagai OQ penegasan): langkah 18 membaca `pyWorkPage.OfferFacIn.PolicyData.
// PolicyNo` - properti yang tidak pernah DITULIS di 329 berkas (spec §15b), sehingga di Pega riwayat tidak pernah
// ditemukan dan blokir AC 102 tidak mungkin terjadi. Di sini nomor polis klaim (`ClaimData.PolicyData.PolicyNo`).
// `Local.CLMNo != .BRANCH_NAME` (pyInsName lawan nomor klaim - alias berbohong §16) menjadi "bukan kasus ini".
// Langkah 19 (keluar bila polis = satu literal) dan langkah 20 (prakondisi nonaktif) tidak dimigrasikan (AC 104).
func CheckDateDOL(k *Konteks, h *Halaman) error {
	h.BersihkanPesan()           // 1
	k.Riwayat(h, TeksRegistrasi) // 2-3
	dol := h.Ambil(CD + "DateOfLoss")
	if h.Ambil(CD+"PolicyData.StartDateTime") == "" { // 5 (T kosong = lanjut, F=3)
		h.TambahPesan(CD+"PolicyData.StartDateTime", PesanIsiMulaiPolis)
	}
	if SesudahTanggal(dol, k.Hari()) { // 6-7
		h.TambahPesan(CD+"DateOfLoss", PesanDOLHariIni)
	}
	if SesudahTanggal(dol, h.Ambil(CD+"DateReceived")) { // 8-9
		h.TambahPesan(CD+"DateOfLoss", PesanDOLTerima)
	}
	if SesudahTanggal(dol, h.Ambil(CD+"ReportDate")) { // 10-11
		h.TambahPesan(CD+"DateOfLoss", PesanDOLLapor)
	}
	// 12-14 (12 prakondisi nonaktif -> selalu): StartDate polis > DOL -> pesan errMsg3.
	if SesudahTanggal(h.Ambil(CD+"PolicyData.StartDateTime"), dol) {
		h.TambahPesan(CD+"DateOfLoss", PesanDOLMulaiPolis)
	}
	// 16-17 (16 prakondisi nonaktif -> selalu): DOL > EndDate polis -> pesan errMsg4.
	if SesudahTanggal(dol, h.Ambil(CD+"PolicyData.EndDateTime")) {
		h.TambahPesan(CD+"DateOfLoss", PesanDOLAkhirPolis)
	}
	return PeriksaDOLSerupa(k, h) // 18-26
}

// PeriksaDOLSerupa - CheckDateDOL_Act langkah 18-26: klaim lain pada polis yang sama dengan Date of Loss sama (bukan
// klaim ini, bukan klaim yang ditolak) -> pesan + `PropBlokirDOL`. Dijalankan juga oleh Save to issue RNM supaya
// blokir AC 102 berlaku walau DOL tidak diubah ulang.
func PeriksaDOLSerupa(k *Konteks, h *Halaman) error {
	dol := h.Ambil(CD + "DateOfLoss")
	h.Hapus(PropBlokirDOL)
	nopol := h.Ambil(CD + "PolicyData.PolicyNo")
	if nopol == "" || dol == "" {
		return nil
	}
	riwayat, err := k.Acuan.RiwayatKlaimPolis(k.Ctxt(), nopol)
	if err != nil {
		return err
	}
	for _, r := range riwayat {
		if r.IDKasus == h.Ambil("pyID") || !SamaTanggal(r.DateOfLoss, dol) || r.Ditolak {
			continue
		}
		h.TambahPesan(CD+"DateOfLoss", PesanDOLSerupa) // 24
		h.Setel(PropBlokirDOL, "1")                    // 25 (§15: penanda blokir sendiri)
		break
	}
	return nil
}

// PropBlokirDOL - penanda "kesalahan yang memblokir" duplikat Date of Loss (pengganti `.IsError = 1`, §15/AC 105).
// Tidak disimpan; dihitung ulang setiap DOL berubah dan diperiksa lagi saat Save to issue RNM.
const PropBlokirDOL = "BlokirDOL"

// CheckDateReceived meniru `Activity/CheckDateReceived_Act.xml`. Langkah 10 dan 13 ber-remark.
func CheckDateReceived(k *Konteks, h *Halaman) {
	h.BersihkanPesan()
	terima := h.Ambil(CD + "DateReceived")
	if SesudahTanggal(terima, k.Hari()) { // 3-4
		h.TambahPesan(CD+"DateReceived", PesanTerimaHariIni)
	}
	if SesudahTanggal(h.Ambil(CD+"DateOfLoss"), terima) { // 5-6
		h.TambahPesan(CD+"DateReceived", PesanTerimaDOL)
	}
	if SesudahTanggal(h.Ambil(CD+"ReportDate"), terima) { // 7-8
		h.TambahPesan(CD+"DateReceived", PesanTerimaLapor)
	}
}

// CheckReportDate meniru `Activity/CheckReportDate_Act.xml`. Langkah 11 dan 14 ber-remark.
// ⚠️ Langkah 3 MENIMPA `DateReceived` dengan hari ini setiap Report Date berubah - ditiru apa adanya.
func CheckReportDate(k *Konteks, h *Halaman) {
	h.BersihkanPesan()
	lapor := h.Ambil(CD + "ReportDate")
	lewatHariIni := SesudahTanggal(lapor, k.Hari())
	h.Setel(CD+"DateReceived", k.Hari()) // 3
	if lewatHariIni {                    // 4
		h.TambahPesan(CD+"ReportDate", PesanLaporHariIni)
	}
	if SesudahTanggal(h.Ambil(CD+"DateOfLoss"), lapor) { // 5-6
		h.TambahPesan(CD+"ReportDate", PesanLaporDOL)
	}
	if SesudahTanggal(lapor, h.Ambil(CD+"DateReceived")) { // 7-8
		h.TambahPesan(CD+"ReportDate", PesanLaporTerima)
	}
}

// SetEndDate meniru `Activity/SetEndDate_Act.xml`: EndDateTime = StartDateTime + 1 tahun, lalu pesan bila awal polis
// di luar periode treaty. Rantai langkah 4-7 (CompareDate/Comparewithend/EqualDate) setara: pesan bila
// mulai < StartDateTreaty ATAU mulai > EndDateTreaty. Langkah 8 (CheckPeriodPolicy_Act) ber-remark.
func SetEndDate(h *Halaman) {
	h.BersihkanPesan()
	mulai := h.Ambil(CD + "PolicyData.StartDateTime")
	h.Setel(CD+"PolicyData.EndDateTime", TambahTahun(mulai, 1)) // 2-3
	awalT, akhirT := h.Ambil(CD+"StartDateTreaty"), h.Ambil(CD+"EndDateTreaty")
	if mulai == "" || awalT == "" {
		return
	}
	if SesudahTanggal(awalT, mulai) || (akhirT != "" && SesudahTanggal(mulai, akhirT)) {
		h.TambahPesan(CD+"PolicyData.StartDateTime", PesanPeriodePolis) // 7
	}
}

// CheckPeriodPolicy - `Activity/CheckPeriodPolicy_Act.xml`: langkah 1 (satu-satunya pengisi Local.CompareDate /
// EqualDate) ber-remark, sehingga langkah 2-4 tidak pernah memenuhi syarat - activity tanpa efek (spec butir terbuka 7).
// Ditulis kosong supaya aksi layar Policy End Date tetap terpetakan.
func CheckPeriodPolicy(*Halaman) {}

// ---------------------------------------------------------------- lokasi dan pelapor

// GetAdders meniru `Activity/GetAdders_Act.xml` (change Zip Code): wilayah dari RW + CITY, lalu Location dirangkai
// `TempLocation.CARI1 + ", " + RW + ", " + District + ", " + City + ", " + Province`.
// `TempLocation.CARI1` diisi `MakeLowercase_Act` (halaman sementara): di sini nilai Location sebelum langkah 1.
func GetAdders(k *Konteks, h *Halaman) error {
	temp := h.Ambil(PropTempLokasi)
	if temp == "" {
		temp = h.Ambil(CD + "Location")
	}
	h.Hapus(CD + "Location") // 1
	baris, err := k.Acuan.Wilayah(k.Ctxt(), h.Ambil(CD+"PostalCode"))
	if err != nil {
		return err
	}
	for _, w := range baris { // 4 - baris terakhir menang
		h.Setel(CD+"RW", w.RW)
		h.Setel(CD+"District", w.District)
		h.Setel(CD+"City", w.City)
		h.Setel(CD+"Province", w.Province)
		h.Setel(CD+"RWID", w.RWID)
		h.Setel(CD+"DistrictID", w.DistrictID)
		h.Setel(CD+"CityID", w.CityID)
		h.Setel(CD+"ProvinceID", w.ProvinceID)
	}
	h.Setel(CD+"Location", temp+", "+h.Ambil(CD+"RW")+", "+h.Ambil(CD+"District")+", "+h.Ambil(CD+"City")+", "+h.Ambil(CD+"Province")) // 5
	return nil
}

// PropTempLokasi - `TempLocation.CARI1` (halaman sementara, tidak disimpan).
const PropTempLokasi = "TempLocation.CARI1"

// MakeLowercase meniru `Activity/MakeLowercase_Act.xml`.
func MakeLowercase(h *Halaman) {
	h.Setel(CD+"Location", strings.ToLower(h.Ambil(CD+"Location")))
	h.Setel(CD+"ReportDescription", strings.ToLower(h.Ambil(CD+"ReportDescription")))
	h.Setel(PropTempLokasi, h.Ambil(CD+"Location"))
}

// GetReportStatus meniru `Activity/GetReportStatus_Act.xml` (change Reporter Status).
func GetReportStatus(k *Konteks, h *Halaman) error {
	var cari1 string
	switch h.Ambil(CD + "ReporterStatus") {
	case "1": // 1 ceding
		h.Setel(CD+"InsuredRelationshipOthers", h.Ambil(TM+"Ceding"))
		cari1 = h.Ambil(TM + "CedingID")
	case "2": // 2 SOB
		h.Setel(CD+"InsuredRelationshipOthers", h.Ambil(TM+"LeadingReinsSource"))
		cari1 = h.Ambil(TM + "LeadingReinsSource")
	}
	h.Hapus(CD + "ReportAddress") // 3
	// 4-5 GetLeaderReport mengisi InputData.CARI2 yang hanya dibaca langkah 7 - tanpa efek di Claim Prop, tidak ditulis.
	// 6 GetAddressCeding memakai InputData.CARI1 (agen) - pada status 1 berisi CedingID, status 2 berisi NAMA SOB.
	// ⚠️ Ditiru apa adanya: bacaan dengan nama SOB sebagai ID agen tidak menemukan baris.
	alamat, ada, err := k.Acuan.AlamatKlien(k.Ctxt(), cari1)
	if err != nil {
		return err
	}
	// 7 GetAddressTreatyIn hanya bila pyWorkPage.IsTreatyIn=="1" - properti itu tidak pernah ditulis korpus Claim Prop.
	if ada {
		h.Setel(CD+"ReportAddress", strings.Join(alamat[:], ", ")) // 8-9
	} else {
		h.Setel(CD+"ReportAddress", ", , , ")
	}
	return nil
}

// SetAdjuster meniru `Activity/SetAdjsuter_act.xml` / SetConsultant meniru `SetConsultant_Act.xml`: nama dari RD
// BrowseAdjusterConsultant; ID atau nama kosong -> keduanya dikosongkan.
func SetAdjuster(k *Konteks, h *Halaman) error {
	return setNamaAdjuster(k, h, CD+"AppointedADJID", CD+"AppointedADJ")
}

// SetConsultant - lihat SetAdjuster.
func SetConsultant(k *Konteks, h *Halaman) error {
	return setNamaAdjuster(k, h, CD+"ConsultantID", CD+"ConsultantName")
}

func setNamaAdjuster(k *Konteks, h *Halaman, jalurID, jalurNama string) error {
	nama, err := k.Acuan.NamaAdjuster(k.Ctxt(), h.Ambil(jalurID))
	if err != nil {
		return err
	}
	h.Setel(jalurNama, nama)
	if h.Ambil(jalurID) == "" || nama == "" {
		h.Setel(jalurID, "")
		h.Setel(jalurNama, "")
	}
	return nil
}

// ---------------------------------------------------------------- katastrofe dan cause of loss

// SetDefNonCatastrope meniru `Activity/SetDefNonCatastrope_Act.xml`.
// Langkah 2: baris 1 `StsKatastrofe=="Catastrophe"` T=3 F=2, baris 2 `NonKatastrofeType==""` T=2 F=3 -> AND.
func SetDefNonCatastrope(h *Halaman) {
	if h.Ambil(CD+"StsKatastrofe") == NilaiKatastrofe {
		h.Setel(CD+"NonKatastrofeType", "")
	}
	if h.Ambil(CD+"StsKatastrofe") != NilaiKatastrofe && h.Ambil(CD+"NonKatastrofeType") == "" {
		h.Setel(CD+"NonKatastrofeType", NilaiNonKatastrofeKlaim)
	}
	if h.Ambil(CD+"NonKatastrofeType") == NilaiNonKatastrofeKlaim {
		h.Setel(CD+"KatastrofeNote", "")
		h.Setel(CD+"KatastrofeID", "")
	}
}

// SetEditCatastrope meniru `Activity/SetEditCatastrope.xml` (ikon Edit / Save Catastrope_Sec). Mengembalikan true
// bila langkah 3 Obj-Save jalan (Action Save).
func SetEditCatastrope(h *Halaman, aksi string) bool {
	switch aksi {
	case "Edit":
		h.Setel(CD+"EditCatastrope", "true")
	case "Save":
		h.Setel(CD+"EditCatastrope", "false")
		return true
	}
	return false
}

// InputCatastrope meniru `Activity/InputCatastrope.xml`: ParamInput.CARI1 1 (Add) / 0 (Cancel) - mode popup
// CatastrofeList_Sec.
func InputCatastrope(h *Halaman, aksi string) {
	switch aksi {
	case "Add":
		h.Setel("ParamInput.CARI1", "1")
	case "Cancel":
		h.Setel("ParamInput.CARI1", "0")
	}
}

// SetCatastrope meniru `Activity/SetCatastrope_act.xml` (tombol Choose grid katastrofe).
func SetCatastrope(h *Halaman, id, catatan string) {
	h.Setel(CD+"KatastrofeNote", catatan)
	h.Setel(CD+"KatastrofeID", id)
}

// KatastrofeBaru - baris CATASTROPHE dari `Activity/SaveCatasrtope_Act.xml` langkah 1 (Obj-Save langkah 3; SQL langkah
// 2 ber-remark).
type KatastrofeBaru struct {
	ID, KlaimType, UserInput, Note, StsKatastrofe, NonKatastrofeType string
}

// SusunKatastrofe meniru langkah 1: ID "CTS-" + waktu sekarang tanpa [GMT.] (`@DateTime.CurrentDateTime()` berbentuk
// "yyyyMMddTHHmmss.SSS GMT" -> "yyyyMMddTHHmmssSSS"), NOTE huruf besar, KLAIMTYPE "NON-LIFE".
func SusunKatastrofe(k *Konteks, h *Halaman, catatan string) KatastrofeBaru {
	t := k.Sekarang.UTC()
	waktu := fmt.Sprintf("%sT%s%03d", t.Format("20060102"), t.Format("150405"), t.Nanosecond()/1e6)
	return KatastrofeBaru{
		ID:                "CTS-" + waktu,
		KlaimType:         "NON-LIFE",
		UserInput:         k.Pelaku,
		Note:              strings.ToUpper(catatan),
		StsKatastrofe:     h.Ambil(CD + "StsKatastrofe"),
		NonKatastrofeType: h.Ambil(CD + "NonKatastrofeType"),
	}
}

// GetNameCauseofLoss meniru `Activity/GetNameCauseofLoss_Act.xml` (Choose popup Cause of Loss): CauseOfLoss = `.Info`
// baris terpilih, lalu Obj-Save.
// ⚠️ `.Info` bukan kolom view V_D_CAUSE_OF_LOSS_BUSINESS (RD BrowseCouseOfLoss_Business). `[data DEV 07-10-2026]`
// 5 dari 5 dokumen klaim CLMP berisi CauseOfLoss = DESCRIPTION view itu dan CauseOfLossID terisi: maka DESCRIPTION
// dan D_COL_ID baris terpilih.
func GetNameCauseofLoss(h *Halaman, deskripsi, id string) {
	h.Setel(CD+"CauseOfLoss", deskripsi)
	h.Setel(CD+"CauseOfLossID", id)
}

// ---------------------------------------------------------------- polis

// CheckNoPolicy meniru `Activity/CheckNoPolicy.xml` (double-click popup Data Polis). Langkah 3-4 ber-remark.
// Langkah 2: `@contains(PolicyNo,"RNM-Q")` T=3 F=2 - pesan muncul saat nomor TIDAK memuat "RNM-Q" (AC 108).
func CheckNoPolicy(k *Konteks, h *Halaman, nopolis, quarter, treatyYear, prodke string) error {
	h.Setel(CD+"PolicyData.PolicyNo", nopolis) // 1
	if !strings.Contains(nopolis, AwalanPolisTreaty) {
		h.TambahPesan(CD+"PolicyData.PolicyNo", PesanFormatPolis)
	}
	h.Setel(CD+"Quater", quarter) // 5
	h.Setel(CD+"TreatyYear", treatyYear)
	if prodke == "" { // 6
		prodke = "0"
	}
	yq, err := k.Acuan.YearOfQuartal(k.Ctxt(), nopolis, prodke) // 7
	if err != nil {
		return err
	}
	h.Setel(CD+"YearofQuartal", yq) // 8; 9 Obj-Save (services)
	return nil
}

// CekPolisAvailable meniru `Activity/CekPolisAvailable_Act.xml` (deferred load blok Information).
// Langkah 1: `NoClaim=="" && IDMaster!=""` T=2 F=6 - selain itu keluar.
func CekPolisAvailable(k *Konteks, h *Halaman) error {
	if !(h.Ambil(CD+"NoClaim") == "" && h.Ambil(CD+"IDMaster") != "") {
		return nil
	}
	ada, err := k.Acuan.AdaPolisMaster(k.Ctxt(), h.Ambil(CD+"IDMaster"))
	if err != nil {
		return err
	}
	if ada {
		h.Setel("IsRealisation", "1")
	} else {
		h.Setel("IsRealisation", "0")
	}
	pesanRealisasi(h)
	return nil
}

// AwalanMaster meniru `Activity/SetMasterID.xml` langkah 1: `@substring(.ClaimData.IDMaster,0,7)` - saringan NOOFFER
// popup Data Polis (RDB SetPolicyTreatyProp).
func AwalanMaster(idMaster string) string {
	r := []rune(idMaster)
	if len(r) > 7 {
		r = r[:7]
	}
	return string(r)
}

// ---------------------------------------------------------------- master treaty

// ParamPilihMaster - parameter tombol Choose grid "Data Master TreatyIn" (Section MasterTreatyInList).
type ParamPilihMaster struct {
	COB, COBID, TreatyContractName, ProportionType, IDMaster, TreatyGroupID, TreatyGroupName string
}

// Template Insured Interest menurut TreatyGroupID (SetValueToClaim_Act langkah 9-12; AC 73 - salah ketik
// "Constrution" dan template Marine Cargo = proyek ikut dibawa).
const (
	templatProperty = "Constrution Class : Occupation:  Risk Category: Risk Location: Coverage :"
	templatProyek   = "Project Name: Risk Location: Construction Period: Project Type:"
	templatMotor    = "Brand: Type/Year: Serial No/Engine No: Model:"
)

var kodeTreatyGroupAneka = map[string]bool{"10012": true, "10011": true, "10055": true, "10059": true, "10002": true,
	"10003": true, "10004": true, "10006": true, "10008": true, "10010": true, "10013": true, "10026": true,
	"10051": true, "10057": true, "10058": true, "10054": true}

// SetValueToClaim meniru `Activity/SetValueToClaim_Act.xml` (Choose master treaty). Obj-Save langkah 18 di services.
func SetValueToClaim(k *Konteks, h *Halaman, p ParamPilihMaster) error {
	h.HapusAwalan("TreatyInMaster") // 1
	h.Hapus(CD + "InsuredInterest")
	h.Setel(CD+"TreatyName", p.TreatyContractName) // 2
	h.Setel(OQ+"BusinessCode", p.COBID)
	h.Setel(CD+"QuotationData.BusinessCode", p.COBID)
	h.Setel(CD+"TreatyGroupID", p.TreatyGroupID)
	h.Setel(CD+"TreatyGroupName", p.TreatyGroupName)
	h.Setel(OQ+"BusinessName", p.COB)
	h.Setel(CD+"QuotationData.BusinessName", p.COB)
	h.Setel(CD+"ProportionType", p.ProportionType)
	h.Setel(CD+"IDMaster", p.IDMaster)
	m, _, err := k.Acuan.MasterTreaty(k.Ctxt(), p.IDMaster) // 3-4
	if err != nil {
		return err
	}
	TerapkanMaster(h, m, true)
	h.Setel(CD+"StartDateTreaty", m.Commencement) // 5
	h.Setel(CD+"EndDateTreaty", m.Termination)
	h.Setel(CD+"YearofAccount", m.TreatyYear)
	tg, err := k.Acuan.TreatyGroupBisnis(k.Ctxt(), p.COBID, m.TreatyYear) // 6-7
	if err != nil {
		return err
	}
	if h.Ambil(CD+"TreatyGroupID") == "" { // 8
		h.Setel(CD+"TreatyGroupID", tg)
	}
	h.Setel(CD+"TreatyGroupName", p.TreatyGroupName)
	h.Setel(CD+"QuotationData.CedingCo", m.CedingID)
	h.Setel(CD+"QuotationData.CedingCoName", m.Ceding)
	h.Setel(CD+"QuotationData.SobName", m.LeadingReinsSource)
	h.Setel(CD+"QuotationData.SobLeader0", m.LeadingReinsSourceID)
	switch { // 9-12 menurut InputData.CARI27 (= TreatyGroupID hasil langkah 7, AC 71)
	case tg == "10007":
		h.Setel(CD+"InsuredInterest", templatProperty)
	case kodeTreatyGroupAneka[tg]:
		h.Setel(CD+"InsuredInterest", templatProyek)
	case tg == "10016":
		h.Setel(CD+"InsuredInterest", templatMotor)
	case tg == "10009":
		h.Setel(CD+"InsuredInterest", templatProyek)
	}
	k.Riwayat(h, TeksPilihMaster)                            // 13
	ada, err := k.Acuan.AdaPolisMaster(k.Ctxt(), p.IDMaster) // 14-17
	if err != nil {
		return err
	}
	h.Setel(CD+"ClaimID", h.Ambil("pyID"))
	if ada {
		h.Setel("IsRealisation", "1")
	} else {
		h.Setel("IsRealisation", "0")
	}
	return nil
}

// TerapkanMaster menulis halaman `TreatyInMaster` dari master. `gantiShare=false` mempertahankan RNMShareP yang sudah
// tersimpan (pilihan layar DisableEditRNMShare) saat master dimuat ulang ketika berkas dibuka.
func TerapkanMaster(h *Halaman, m MasterTreaty, gantiShare bool) {
	share := h.Ambil(TM + "RNMShareP")
	h.HapusAwalan("TreatyInMaster")
	for p, v := range map[string]string{
		"ID": m.ID, "TreatyContractName": m.TreatyContractName, "ProportionType": m.ProportionType,
		"Ceding": m.Ceding, "CedingID": m.CedingID, "LeadingReinsSource": m.LeadingReinsSource,
		"LeadingReinsSourceID": m.LeadingReinsSourceID, "Bordeaux": m.Bordeaux, "BordereauxNote": m.BordereauxNote,
		"AccountingMode": m.AccountingMode, "TeritorialScope": m.TeritorialScope, "Commencement": m.Commencement,
		"Termination": m.Termination, "TreatyYear": m.TreatyYear, "RNMShareP": ShareMaster(m, h.Ambil(CD+"TreatyGroupID")),
		"StatusAkseptasi": m.StatusAkseptasi,
	} {
		if v != "" {
			h.Setel(TM+p, v)
		}
	}
	if !gantiShare && share != "" {
		h.Setel(TM+"RNMShareP", share)
	}
	var lim []Baris
	for _, l := range m.Limits {
		lim = append(lim, Baris{"TreatyType": l.TreatyType})
	}
	h.SetelDaftar(TM+"Limits", lim)
}

// ShareMaster - RNM Share master untuk klaim (keputusan work owner 10-10-2026: RNM_SHARE view detail treaty group
// klaim, menggantikan `TreatyInMaster.RNMShareP` JSON): RNMShare Detail pertama ber-TreatyGroupID klaim; tanpa yang
// cocok, nilai bila SELURUH detail berbagi satu nilai; selain itu kosong - dipilih dari dropdown RNM Share
// (`PilihanShareRNM`).
func ShareMaster(m MasterTreaty, grup string) string {
	satu, banyak := "", false
	for _, l := range m.Limits {
		for _, d := range l.Detail {
			if grup != "" && d.TreatyGroupID == grup && d.RNMShare != "" {
				return d.RNMShare
			}
			switch {
			case d.RNMShare == "":
			case satu == "":
				satu = d.RNMShare
			case satu != d.RNMShare:
				banyak = true
			}
		}
	}
	if banyak {
		return ""
	}
	return satu
}

// PilihanShareRNM meniru `Activity/GetRNMShareTreaty.xml`: seluruh `Limits(n).Detail(m).RNMShare` tanpa duplikat
// (Java langkah 3 membuang duplikat dari belakang - kemunculan PERTAMA dipertahankan), lalu IsEditRNMShare = true.
func PilihanShareRNM(h *Halaman, m MasterTreaty) []string {
	var out []string
	sudah := map[string]bool{}
	for _, l := range m.Limits {
		for _, d := range l.Detail {
			if !sudah[d.RNMShare] {
				sudah[d.RNMShare] = true
				out = append(out, d.RNMShare)
			}
		}
	}
	h.Setel("IsEditRNMShare", "true")
	return out
}

// DisableEditRNMShare meniru `DataTransform/DisableEditRNMShare.xml` (change dropdown RNM Share):
// IsEditRNMShare = "false", TreatyInMaster.RNMShareP = pilihan.
func DisableEditRNMShare(h *Halaman, pilihan string) {
	h.Setel("IsEditRNMShare", "false")
	if pilihan != "" {
		if d, err := AngkaTeks("RNMShare", pilihan); err == nil {
			h.Setel(TM+"RNMShareP", Teks(d))
		}
	}
}
