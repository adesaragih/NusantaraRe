package models

// Untuk apa berkas ini: PORT ACTIVITY REGISTRASI KLAIM Non Prop - pra/pasca-proses flow action, pilih master, polis,
// validasi tanggal, pelapor, adjuster/consultant, lokasi, katastrofe, cause of loss.
//
// Setiap fungsi menyebut berkas Activity aslinya di korpus `Claim Non Prop`; nomor langkah di komentar = nomor langkah
// Pega. Langkah berlabel `//` TIDAK ditulis. Kode arah kosong (`T:` tanpa angka) = lanjut (aturan baca Claim Prop 5b).

import (
	"fmt"
	"strings"
)

// Pesan validasi - teks VERBATIM korpus (salah ketik ikut).
const (
	PesanPolisKosong      = "Policy no can't be empty!"                          // InputOutStandingCTNP_PostAct 1
	PesanBisnisKosong     = "BusinessName tidak boleh Kosong"                    // CheckDateDOL_Act 7
	PesanIsiMulaiPolis    = "Please Fill Policy Start Date First"                // CheckDateDOL_Act 6
	PesanDOLHariIni       = "Date of Loss should not be more than todays date"   // CheckDateDOL_Act 6
	PesanDOLTerima        = "Date Of Loss should not be more than Received Date" // CheckDateDOL_Act 6
	PesanDOLLapor         = "Date Of Loss should not be more than Report Date"   // CheckDateDOL_Act 6
	PesanDOLMulaiPolis    = "Date of Loss should not be more than Start Date Policy  "
	PesanDOLAkhirPolis    = "Date of Loss should not be more than End Date Policy "
	PesanDOLTreaty        = "Date Of Loss Must not be more than Treaty Period"    // CheckDateDOL_Act 6 errMsg8
	PesanKlaimKembar      = "Claims with Master ID, DOL and Bisnis already exist" // CheckDateDOL_Act 6 errMsg7
	PesanTerimaHariIni    = "Date Received should not be more than todays date"   // CheckDateReceived_Act 2
	PesanTerimaDOL        = "Date Received  should not be less than Date of Loss"
	PesanTerimaLapor      = "Date Received should not be less than Report Date"
	PesanLaporHariIni     = "ReportDate should not be more than todays date" // CheckReportDate_Act 2
	PesanLaporDOL         = "ReportDate should not be less than Date of Loss"
	PesanLaporTerima      = "ReportDate should not be more than Received Date"
	PesanPeriodeRisk      = "Policy Period must not be more than Treaty Period" // SetEndDate_Act 4
	TeksRegistrasi        = "Input Data Registration Claim"                     // CheckDateDOL_Act 4
	StatusInputAcceptance = "INPUT ACCEPTATION CLAIM"                           // InputOutStandingCTNP_PostAct 2
	NilaiNonKatastrofe    = "Claim"
	NilaiKatastrofe       = "Catastrophe"
	ModeLoss              = "loss" // TreatyInMaster.AccountingModeNonProp
	ModeRisk              = "risk"
)

// ---------------------------------------------------------------- pra / pasca proses

// PilihanLossAllocation = InputOutStandingClmTNP_PreAct langkah 2 + 4: daftar dropdown Treaty Name Loss Allocation
// (`ListLossAllocation.pxResults.CARI1`) hanya berisi "OR" (langkah 5 hardcode CLMNP-232 dibuang, OQ-CNP-03; langkah 6
// ber-remark).
var PilihanLossAllocation = []string{"OR"}

// SaringTreatySpreading = InputOutStandingClmTNP_PreAct langkah 3.6: nama treaty GetTreatyName_SQL yang memuat "TRT"
// atau "ORS" (langkah 3.6.2 hardcode CLMNP-232 dibuang) - sumber dropdown Treaty Type Spreading List.
func SaringTreatySpreading(nama []string) []string {
	var out []string
	for _, n := range nama {
		if strings.Contains(n, "TRT") || strings.Contains(n, "ORS") {
			out = append(out, n)
		}
	}
	return out
}

// AwalanMaster = InputOutStandingClmTNP_PreAct langkah 3.1 `@substring(IDMaster,0,7)` (NOOFFER kedua GetDataPolisNonProp_SQL).
func AwalanMaster(idMaster string) string {
	r := []rune(idMaster)
	if len(r) > 7 {
		r = r[:7]
	}
	return string(r)
}

// PraOutstanding = `InputOutStandingClmTNP_PreAct` (pra-proses FlowAction OutstandingClaim): daftar polis dan treaty
// dibaca services saat dropdown dibuka (langkah 3); langkah 7 SethistoryKlaimTreaty = `TampilRiwayat`; langkah 8
// SetMOClaimTreaty_Act.
func PraOutstanding(k *Konteks, h *Halaman) error {
	return SetMOClaimTreaty(k, h)
}

// PascaOutstanding = `InputOutStandingCTNP_PostAct` (post-activity Submit Outstanding Claim).
func PascaOutstanding(h *Halaman) {
	h.Setel("CNPStatusCase", StatusInputAcceptance) // 2
	if h.Ambil(CD+"PolicyData.PolicyNo") == "" {    // 3
		h.TambahPesan("", PesanPolisKosong)
	}
}

// TanggalAkseptasiLama - batas akseptasi lama (InputAkseptasi_PreAct 3.1 `.AcceptedDate > @toDate("20201214")`).
const TanggalAkseptasiLama = "2020-12-14"

// PraAkseptasi = `InputAkseptasi_PreAct` (pra-proses FlowAction InputAcceptation): langkah 3 (CommentLOD, IsAcceptation),
// langkah 4 (InputOutStandingClmTNP_PreAct), langkah 5 (mesin XoL bila SpreadingRisk kosong), langkah 7 (MO).
func PraAkseptasi(k *Konteks, h *Halaman, m MasterTreaty) error {
	TandaiAkseptasiLama(h)
	for _, b := range h.AmbilDaftar(DaftarAdjustment) { // 3.2
		if b["AcceptanceStatus"] == "1" {
			h.Setel("IsAcceptation", "1")
		}
	}
	if err := PraOutstanding(k, h); err != nil {
		return err
	}
	if len(h.AmbilDaftar(DaftarXOL)) == 0 { // 5
		if err := HitungXOL(k, h, m, HitungXOLMod, 0, ""); err != nil {
			return err
		}
	}
	return nil
}

// TandaiAkseptasiLama = InputAkseptasi_PreAct langkah 3.1: akseptasi berstatus 1/2, bertanggal <= 14-12-2020, tanpa
// ListClaimAcceptation -> `CommentLOD = 1` (detail disembunyikan). Medan turunan, dihitung setiap halaman dimuat.
func TandaiAkseptasiLama(h *Halaman) {
	for i, b := range h.AmbilDaftar(DaftarAdjustment) {
		lama := (b["AcceptanceStatus"] == "1" || b["AcceptanceStatus"] == "2") &&
			!SesudahTanggal(b["AcceptedDate"], TanggalAkseptasiLama) &&
			len(h.AmbilDaftar(JalurAdj(i+1, AnakClaimAccept))) == 0
		if lama {
			b["CommentLOD"] = "1"
		} else {
			delete(b, "CommentLOD")
		}
	}
}

// SetMOClaimTreaty = `SetMOClaimTreaty_Act` langkah 1: baris 1 `PolicyNo==""` T=5 (Skip Whens -> keluar), baris 2
// `BranchDetailID==""` T=3 F=2 -> keluar bila polis kosong ATAU BranchDetailID TERISI.
func SetMOClaimTreaty(k *Konteks, h *Halaman) error {
	if h.Ambil(CD+"PolicyData.PolicyNo") == "" || h.Ambil(CD+"MarketingData.BranchDetailID") != "" {
		return nil
	}
	mo, ada, err := k.Acuan.MarketingPolis(k.Ctxt(), h.Ambil(CD+"PolicyData.PolicyNo")) // 3
	if err != nil || !ada {
		return err
	}
	for i, p := range []string{"ID", "ClientID", "ClientName", "TeamGroup", "BranchDetailID", "BranchDetailName"} {
		h.Setel(CD+"MarketingData."+p, mo[i]) // 4
	}
	return nil
}

// ---------------------------------------------------------------- master treaty

// ParamPilihMaster - parameter tombol Choose popup ChooseMasterTNonProp (ID, BusinessName, BusinessCode, TreatyGroup).
type ParamPilihMaster struct {
	ID, BusinessName, BusinessCode, TreatyGroup string
	// TreatyName - `.TREATYCONTRACTNAME` baris popup. `[inferensi]` (PARITAS): langkah 2 menulis `Param.TreatyName` yang
	// tidak dikirim tombol Choose (kosong di XML), padahal kasus DEV CLMNP- berisi nama kontrak master.
	TreatyName string
}

// PilihMaster = `SetValueClaimTNP_Act` (tombol Choose ChooseMasterTNonProp).
//
// ⚠️ Kelainan XML (PARITAS, OQ): tombol Choose mengirim `BusinessCode` (`.CLASSOFBUSINESSID`), tetapi activity tidak
// mendeklarasikan/membacanya, padahal GetTreatyName_SQL (daftar treaty spreading) dan GetLimitTONPPLA membaca
// `OfferFacIn.QuotationData.BusinessCode`. Di sini parameter itu disimpan ke BusinessCode. TreatyName / ProportionType /
// Ceding / dst. yang dideklarasikan activity tidak dikirim tombol - langkah 2 menulis kosong, lalu langkah 3-4
// (`adoptJSONObject`) mengisinya dari master.
func PilihMaster(k *Konteks, h *Halaman, p ParamPilihMaster) error {
	h.HapusAwalan("TreatyInMaster") // 1
	h.Hapus(CD + "InsuredInterest")
	h.Setel(CD+"IDMaster", p.ID) // 2
	h.Setel(CD+"TreatyName", p.TreatyName)
	h.Setel(OQ+"BusinessName", p.BusinessName)
	h.Setel(OQ+"BusinessCode", p.BusinessCode)
	h.Setel(TM+"TreatyGroup", p.TreatyGroup)
	m, _, err := k.Acuan.MasterTreaty(k.Ctxt(), p.ID) // 3-4
	if err != nil {
		return err
	}
	TerapkanMaster(h, m, true)
	h.Setel(TM+"TreatyGroup", p.TreatyGroup)
	h.Setel(CD+"StartDateTreaty", m.Commencement) // 5
	h.Setel(CD+"EndDateTreaty", m.Termination)
	h.Setel(CD+"YearofAccount", m.TreatyYear)
	ada, err := k.Acuan.AdaPolisMaster(k.Ctxt(), p.ID) // 6
	if err != nil {
		return err
	}
	if ada { // 7
		h.Setel("IsRealisation", "1")
	} else {
		h.Setel("IsRealisation", "0")
	}
	return nil
}

// TerapkanMaster menulis halaman `TreatyInMaster` dari master. `gantiShare=false` mempertahankan RNMShare tersimpan
// (salinan saat master dipilih - Pega menyimpan halaman TreatyInMaster di BLOB kasus) ketika berkas dibuka ulang.
func TerapkanMaster(h *Halaman, m MasterTreaty, gantiShare bool) {
	share, grup := h.Ambil(TM+"RNMShare"), h.Ambil(TM+"TreatyGroup")
	h.HapusAwalan("TreatyInMaster")
	for p, v := range map[string]string{
		"ID": m.ID, "ProportionType": m.ProportionType, "Ceding": m.Ceding, "CedingID": m.CedingID,
		"LeadingReinsSource": m.LeadingReinsSource, "LeadingReinsSourceID": m.LeadingReinsSourceID,
		"Bordeaux": m.Bordeaux, "BordereauxNote": m.BordereauxNote, "AccountingMode": m.AccountingMode,
		"AccountingModeNonProp": m.AccountingModeNonProp, "TeritorialScope": m.TeritorialScope,
		"Commencement": m.Commencement, "Termination": m.Termination, "TreatyYear": m.TreatyYear,
		"RNMShare": m.RNMShare, "EDMState": m.EDMState, "StatusAkseptasi": m.StatusAkseptasi,
	} {
		if v != "" {
			h.Setel(TM+p, v)
		}
	}
	if !gantiShare && share != "" {
		h.Setel(TM+"RNMShare", share)
	}
	if grup != "" {
		h.Setel(TM+"TreatyGroup", grup)
	}
}

// ---------------------------------------------------------------- polis

// CheckNoPolicy = `Activity/CheckNoPolicy.xml` (pilih polis: autocomplete Policy No / tautan ViewListPolicyCNP).
// Langkah 2-9 ber-remark.
func CheckNoPolicy(h *Halaman, nopolis, cob string) {
	h.Setel(CD+"PolicyData.PolicyNo", nopolis)
	h.Setel(CD+"PolicyData.TreatyGroup", cob)
}

// ---------------------------------------------------------------- tanggal

// CheckDateDOL = `CheckDateDOL_Act` (change Date of Loss). Langkah 26-28, 32, 38-40 ber-remark; langkah 33-37
// (CariHistoryClaim + RejectedClaim_RD) hanya mengisi `Local.Flagerr` yang dibaca langkah ber-remark - tanpa efek,
// tidak ditulis.
func CheckDateDOL(k *Konteks, h *Halaman) error {
	h.BersihkanPesan()                    // 1
	k.Riwayat(h, TeksRegistrasi)          // 4-5
	if h.Ambil(OQ+"BusinessName") == "" { // 8-9
		h.Setel("ProtectCOB.CARI1", "1")
		h.TambahPesan(OQ+"BusinessName", PesanBisnisKosong)
	}
	dol := h.Ambil(CD + "DateOfLoss")
	if h.Ambil(CD+"PolicyData.StartDateTime") == "" { // 10 (T kosong = lanjut)
		h.TambahPesan(CD+"PolicyData.StartDateTime", PesanIsiMulaiPolis)
	}
	if SesudahTanggal(dol, k.Hari()) { // 11-12
		h.TambahPesan(CD+"DateOfLoss", PesanDOLHariIni)
	}
	if SesudahTanggal(dol, h.Ambil(CD+"DateReceived")) { // 13-14
		h.TambahPesan(CD+"DateOfLoss", PesanDOLTerima)
	}
	if SesudahTanggal(dol, h.Ambil(CD+"ReportDate")) { // 15-16
		h.TambahPesan(CD+"DateOfLoss", PesanDOLLapor)
	}
	if SesudahTanggal(h.Ambil(CD+"PolicyData.StartDateTime"), dol) { // 17-19
		h.TambahPesan(CD+"DateOfLoss", PesanDOLMulaiPolis)
	}
	if SesudahTanggal(dol, h.Ambil(CD+"PolicyData.EndDateTime")) { // 20-22
		h.TambahPesan(CD+"DateOfLoss", PesanDOLAkhirPolis)
	}
	if ProteksiDOL(h) { // 23-25
		h.TambahPesan(CD+"DateOfLoss", PesanDOLTreaty)
	}
	return PeriksaKlaimKembar(k, h) // 29-31
}

// ProteksiDOL = CheckDateDOL_Act langkah 23-24 (`ProtectDOL.CARI1`): mode akuntansi "loss" dan DOL melewati akhir treaty.
// Halaman ProtectDOL di Pega hidup di requestor; di sini dihitung ulang dari halaman tersimpan (`HitungTurunan`).
func ProteksiDOL(h *Halaman) bool {
	return h.Ambil(TM+"AccountingModeNonProp") == ModeLoss &&
		SesudahTanggal(h.Ambil(CD+"DateOfLoss"), h.Ambil(CD+"EndDateTreaty"))
}

// ProteksiMulaiPolis = SetEndDate_Act langkah 8-9 (`ProtectStartDate.CARI1`): mode "risk" dan awal polis melewati akhir
// treaty.
func ProteksiMulaiPolis(h *Halaman) bool {
	return h.Ambil(TM+"AccountingModeNonProp") == ModeRisk &&
		SesudahTanggal(h.Ambil(CD+"PolicyData.StartDateTime"), h.Ambil(CD+"EndDateTreaty"))
}

// PeriksaKlaimKembar = CheckDateDOL_Act langkah 29-31 (dan SaveDataToOSAksep_Act langkah 6): klaim lain (bukan kasus ini,
// tidak di CLAIMREJECTED) pada polis yang sama dengan Master ID, DOL, dan bisnis yang sama -> pesan.
func PeriksaKlaimKembar(k *Konteks, h *Halaman) error {
	nopol := h.Ambil(CD + "PolicyData.PolicyNo")
	if nopol == "" {
		return nil
	}
	riwayat, err := k.Acuan.RiwayatKlaimPolis(k.Ctxt(), nopol)
	if err != nil {
		return err
	}
	id := h.Ambil("pyID")
	for _, r := range riwayat {
		if r.Ditolak || r.IDKasus == id || r.KunciKasus == KunciPegaLama(id) || r.KunciKasus == id {
			continue
		}
		if r.IDMaster == h.Ambil(CD+"IDMaster") && SamaTanggal(r.DateOfLoss, h.Ambil(CD+"DateOfLoss")) &&
			r.Bisnis == h.Ambil(OQ+"BusinessName") {
			h.TambahPesan(CD+"DateOfLoss", PesanKlaimKembar)
			break
		}
	}
	return nil
}

// CheckDateReceived = `CheckDateReceived_Act`. Langkah 10 dan 13 ber-remark.
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

// CheckReportDate = `CheckReportDate_Act`. Langkah 11 dan 14 ber-remark.
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

// SetEndDate = `SetEndDate_Act`: EndDateTime = StartDateTime + 1 tahun (5-6), lalu mode "risk" dengan awal polis melewati
// akhir treaty -> pesan (8-12). Langkah 13-14 ber-remark.
func SetEndDate(h *Halaman) {
	h.BersihkanPesan() // 1
	h.Setel(CD+"PolicyData.EndDateTime", TambahTahun(h.Ambil(CD+"PolicyData.StartDateTime"), 1))
	if ProteksiMulaiPolis(h) {
		h.TambahPesan(CD+"PolicyData.StartDateTime", PesanPeriodeRisk)
	}
}

// CheckPeriodPolicy = `CheckPeriodPolicy_Act`: langkah 1 (satu-satunya pengisi Local.CompareDate / EqualDate) dan
// langkah 5, 8, 9 ber-remark, sehingga langkah 6, 7, 10 tidak pernah memenuhi syarat - activity tanpa efek. Ditulis
// kosong supaya aksi layar Policy End tetap terpetakan.
func CheckPeriodPolicy(*Halaman) {}

// ---------------------------------------------------------------- lokasi dan pelapor

// PropTempLokasi - `TempLocation.CARI1` (halaman sementara, tidak disimpan).
const PropTempLokasi = "TempLocation.CARI1"

// MakeLowercase = `MakeLowercase_Act`.
func MakeLowercase(h *Halaman) {
	h.Setel(CD+"Location", strings.ToLower(h.Ambil(CD+"Location")))
	h.Setel(CD+"ReportDescription", strings.ToLower(h.Ambil(CD+"ReportDescription")))
	h.Setel(PropTempLokasi, h.Ambil(CD+"Location"))
}

// GetAdders = `GetAdders_Act` (change Zip Code): wilayah dari RW + CITY, lalu Location dirangkai
// `TempLocation.CARI1 + ", " + RW + ", " + District + ", " + City + ", " + Province`.
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

// GetReportStatus = `GetReportStatus_Act` (change Reporter Status). Status 1: InputData.CARI1 = CedingID; status 2:
// InputData.CARI1 = NAMA SOB dan CARI2 = SOB ID. Langkah 8 GetAddressCeding membaca alamat klien milik agen CARI1
// (M_CLIENT JSON, RDB XML) - pada status 2 nama SOB dipakai sebagai ID agen sehingga tidak menemukan baris (ditiru apa
// adanya). Langkah 9 GetAddressTreatyIn hanya bila `pyWorkPage.IsTreatyIn=="1"` - properti tanpa penulis.
func GetReportStatus(k *Konteks, h *Halaman) error {
	h.Setel(CD+"ReportAddress", "") // 1
	cari1 := ""
	switch h.Ambil(CD + "ReporterStatus") {
	case "1": // 3
		h.Setel(CD+"InsuredRelationshipOthers", h.Ambil(TM+"Ceding"))
		cari1 = h.Ambil(TM + "CedingID")
	case "2": // 4
		h.Setel(CD+"InsuredRelationshipOthers", h.Ambil(TM+"LeadingReinsSource"))
		cari1 = h.Ambil(TM + "LeadingReinsSource")
	}
	alamat, err := k.Acuan.AlamatAgen(k.Ctxt(), cari1) // 8, 10
	if err != nil {
		return err
	}
	h.Setel(CD+"ReportAddress", strings.Join(alamat[:], ", ")) // 11
	return nil
}

// SetAdjuster = `SetAdjsuter_act` / SetConsultant = `SetConsultant_Act`: nama dari RD BrowseAdjusterConsultant.
// Pola tampilan Claim Prop: dropdown memilih ID dan menampilkan nama (PARITAS); ID atau nama kosong -> keduanya kosong.
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

// SetDefNonCatastrope = `SetDefNonCatastrope_Act` (langkah 2: `StsKatastrofe=="Catastrophe"` T=3 F=2 lalu
// `NonKatastrofeType==""` T=2 F=3 -> AND).
func SetDefNonCatastrope(h *Halaman) {
	if h.Ambil(CD+"StsKatastrofe") == NilaiKatastrofe {
		h.Setel(CD+"NonKatastrofeType", "")
	}
	if h.Ambil(CD+"StsKatastrofe") != NilaiKatastrofe && h.Ambil(CD+"NonKatastrofeType") == "" {
		h.Setel(CD+"NonKatastrofeType", NilaiNonKatastrofe)
	}
	if h.Ambil(CD+"NonKatastrofeType") == NilaiNonKatastrofe {
		h.Setel(CD+"KatastrofeNote", "")
		h.Setel(CD+"KatastrofeID", "")
	}
}

// SetEditCatastrope = `SetEditCatastrope` (ikon Edit / Save Catastrope_Sec).
func SetEditCatastrope(h *Halaman, aksi string) {
	switch aksi {
	case "Edit":
		h.Setel(CD+"EditCatastrope", "true")
	case "Save":
		h.Setel(CD+"EditCatastrope", "false")
	}
}

// InputCatastrope = `InputCatastrope`: ParamInput.CARI1 1 (Add) / 0 (Cancel) - mode popup CatastrofeList_Sec.
func InputCatastrope(h *Halaman, aksi string) {
	switch aksi {
	case "Add":
		h.Setel("ParamInput.CARI1", "1")
	case "Cancel":
		h.Setel("ParamInput.CARI1", "0")
	}
}

// SetCatastrope = `SetCatastrope_act` (Choose grid katastrofe).
func SetCatastrope(h *Halaman, id, catatan string) {
	h.Setel(CD+"KatastrofeNote", catatan)
	h.Setel(CD+"KatastrofeID", id)
}

// KatastrofeBaru - baris CATASTROPHE `SaveCatasrtope_Act` langkah 1 + 3 (Obj-Save ParamCat; SQL langkah 2 ber-remark).
type KatastrofeBaru struct {
	ID, KlaimType, UserInput, Note, StsKatastrofe, NonKatastrofeType string
}

// SusunKatastrofe = SaveCatasrtope_Act langkah 1: ID "CTS-" + `@DateTime.CurrentDateTime()` tanpa [GMT.], NOTE huruf
// besar, KLAIMTYPE "NON-LIFE".
func SusunKatastrofe(k *Konteks, h *Halaman, catatan string) KatastrofeBaru {
	t := k.Sekarang.UTC()
	waktu := fmt.Sprintf("%sT%s%03d", t.Format("20060102"), t.Format("150405"), t.Nanosecond()/1e6)
	return KatastrofeBaru{ID: "CTS-" + waktu, KlaimType: "NON-LIFE", UserInput: k.Pelaku, Note: strings.ToUpper(catatan),
		StsKatastrofe: h.Ambil(CD + "StsKatastrofe"), NonKatastrofeType: h.Ambil(CD + "NonKatastrofeType")}
}

// GetNameCauseofLoss = `GetNameCauseofLoss_Act` (Choose popup Cause of Loss): DESCRIPTION dan D_COL_ID baris terpilih.
func GetNameCauseofLoss(h *Halaman, deskripsi, id string) {
	h.Setel(CD+"CauseOfLoss", deskripsi)
	h.Setel(CD+"CauseOfLossID", id)
}
