package models

// Untuk apa berkas ini: PORT ACTIVITY REGISTRASI KLAIM (FlowAction InputRegister, Section InputRegister /
// InputRegisterDetail): pra-proses, pemeriksaan tanggal, klaim kembar, Reporter Status, kode pos, adjuster / consultant,
// katastrofe, cause of loss, pop-up Choose Polis, dan Submit. Setiap fungsi menyebut berkas Activity aslinya di korpus
// `Claim Fac In`; nomor langkah di komentar = nomor langkah XML (`//` = ter-remark, tidak dibangun).

import (
	"context"
	"strings"
)

// Pesan validasi - teks VERBATIM korpus (salah ketik dan spasi ikut).
const (
	PesanDOLLebihHariIni    = "Date of Loss should not be more than todays date"         // CheckDate_Act 2
	PesanDOLLebihTerima     = "Date Of Loss should not be more than Received Date"       // CheckDate_Act 2
	PesanDOLLebihLapor      = "Date Of Loss should not be more than Report Date"         // CheckDate_Act 2
	PesanDOLLebihMulaiPolis = "Date of Loss should not be more than Start Date Policy  " // CheckDate_Act 2
	PesanDOLLebihAkhirPolis = "Date of Loss should not be more than End Date Policy "    // CheckDate_Act 2
	PesanLaporLebihHariIni  = "ReportDate should not be more than todays date"           // CheckDateReport_Act 1
	PesanLaporKurangDOL     = "ReportDate should not be less than Date of Loss"          // CheckDateReport_Act 1
	PesanLaporLebihTerima   = "ReportDate should not be more than Received Date"         // CheckDateReport_Act 1
	PesanTerimaLebihHariIni = "Date Received should not be more than todays date"        // CheckDateReceived_Act 1
	PesanTerimaKurangDOL    = "Date Received  should not be less than Date of Loss"      // CheckDateReceived_Act 1
	PesanTerimaKurangLapor  = "Date Received should not be less than Report Date"        // CheckDateReceived_Act 1
	PesanSebabKosong        = "Cause of Loss can't be blank"                             // ProteksiDataRegister_Act 2
	PesanTeleponKosong      = "Reporter Phone Number can't be blank"                     // ProteksiDataRegister_Act 2
	PesanPelaporKosong      = "Reporter Name can't be blank"                             // ProteksiDataRegister_Act 2
	PesanPilihObjek         = "Please, select at least one object!"                      // CheckListEstimasi_Act 1
)

// CheckPeriodePolicy (Submit pop-up ViewPolis) tidak dibangun: langkah 1 `Exit-Activity` pre=false keluar tanpa syarat,
// sehingga cek periode / endorsement langkah 2-6 tidak pernah jalan. Maksudnya tidak pasti -> OQ-CFI-34 (PARITAS).

// Jalur registrasi.
const (
	DaftarRiwayatPolis = "HistoryClaim" // Primary.HistoryClaim (CallActivityInputRegister 9.1, grid "Claim History")
	JalurIsError       = "IsError"
	JalurIsEstimation  = "isEstimation"
	JalurIsAdjustment  = "IsAdjustment"
	JalurIsRegister    = "IsRegister"
	JalurIsAnyAccept   = "IsAnyAcceptation"
	JalurIsTreatyIn    = "IsTreatyIn"
	// JalurLokasiSementara - `TempLocaion.CARI1` (SetTempLocation_Act; halaman requestor - dibawa `Layar.Mode`).
	JalurLokasiSementara = "TempLocaion.CARI1"
)

// PraRegister = pra-proses FlowAction InputRegister: `CallActivityInputRegister` (riwayat klaim polis, IsAnyAcceptation)
// lalu DT `InsertObjects_dt` (cabang FacIn; calon objek dari polis). `riwayat` = BrowseHistoryClaim polis kasus.
//
// InsertProgressClaim (CallActivityInputRegister 12) menulis PROGRESSCLAIM - efek basis data, dijalankan services di
// transaksi aksi (`TulisProgress`), bukan di pra-proses baca.
func PraRegister(k *Konteks, h *Halaman, riwayat []RiwayatKlaimPolis) {
	h.Setel("pyNote", "") // CallActivityInputRegister 5
	SetelRiwayatPolis(h, riwayat)
	CheckAnyAcceptation(h) // 11
	// InsertObjects_dt
	if objek := h.AmbilDaftar(DaftarObjek); len(objek) > 0 && objek[0][PropKunciPolis] == "" &&
		objek[0]["ObjectID"] == "" && objek[0]["ObjectName"] == "" { // 1
		h.HapusAwalanDaftar(DaftarObjek)
	}
	h.Hapus(JalurIsEstimation)              // 2
	h.Setel(JalurIsAdjustment, "")          // 3
	SusunCalon(h)                           // 4, 8
	h.Setel(CD+"ClaimNo", h.Ambil("pyID"))  // 9
	if h.Ambil("StartDateRegister") == "" { // 10
		h.Setel("StartDateRegister", k.Waktu())
	}
	h.Setel(CD+"ExGratia", "0")   // 11
	h.Setel(JalurIsRegister, "1") // 12
}

// SetelRiwayatPolis = CallActivityInputRegister 8-9.1: grid "Claim History" (`.HistoryClaim`: ClaimNo, DateOfLoss).
// Langkah 8 bergerbang nomor polis tertulis mati - DIBUANG (bawaan c prompt §3): berlaku untuk setiap polis.
func SetelRiwayatPolis(h *Halaman, riwayat []RiwayatKlaimPolis) {
	var rows []Baris
	for _, r := range riwayat {
		rows = append(rows, Baris{"ClaimNo": r.ClaimNo, "DateOfLoss": r.DateOfLoss})
	}
	h.SetelDaftar(DaftarRiwayatPolis, rows)
}

// CheckAnyAcceptation = `CheckAnyAcceptation`: IsAnyAcceptation = 1 bila ada adjustment ber-AcceptanceStatus 1.
func CheckAnyAcceptation(h *Halaman) {
	n := 0
	for o := range h.AmbilDaftar(DaftarObjek) {
		for i := range h.AmbilDaftar(DaftarItem(o + 1)) {
			for _, a := range h.AmbilDaftar(DaftarDiItem(o+1, i+1, AnakAdj)) {
				if a["AcceptanceStatus"] == "1" {
					n++
				}
			}
		}
	}
	h.Setel(JalurIsAnyAccept, nolSatu(n > 0))
}

func nolSatu(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// CheckDate = `CheckDate_Act` (dipanggil ProteksiDataRegister_Act 3 bila DOL terisi). Langkah 10, 14, 17, 19
// ber-remark; cabang TreatyIn tidak terjangkau (polis Choose Polis selalu RNM-F).
//
// ⚠️ Langkah 9 / 13 `Param.StartDate := Policy.StartDateTime+1` / `EndDateTime+1`: tanggal polis sesudah CopyNB_Act
// langkah 11 berbentuk teks "yyyyMMdd"; `+1` atas teks di ekspresi Pega = penggabungan teks ("201801011"), tanggal yang
// tak terurai sehingga `@CompareDates` langkah 11 selalu salah - pesan "Start Date Policy" tidak pernah terpasang.
// Maksudnya tidak pasti (+1 hari atau tanpa +1) -> OQ-CFI-12; dibangun mengikuti hasil XML: pemeriksaan awal polis
// tidak memasang pesan. Langkah 15 membandingkan DOL dengan EndDateTime APA ADANYA (bukan Param.EndDate) - dibangun.
func CheckDate(k *Konteks, h *Halaman) {
	dol := h.Ambil(CD + "DateOfLoss")
	if SesudahTanggal(dol, k.Hari()) { // 3-4
		h.TambahPesan(CD+"DateOfLoss", PesanDOLLebihHariIni)
	}
	if SesudahTanggal(dol, h.Ambil(CD+"DateReceived")) { // 5-6
		h.TambahPesan(CD+"DateOfLoss", PesanDOLLebihTerima)
	}
	if SesudahTanggal(dol, h.Ambil(CD+"ReportDate")) { // 7-8
		h.TambahPesan(CD+"DateOfLoss", PesanDOLLebihLapor)
	}
	if h.Ambil(JalurIsTreatyIn) != "0" { // 13, 15 bergerbang IsTreatyIn==0
		return
	}
	lewat := SesudahTanggal(dol, h.Ambil(JalurAkhirPolis)) // 15
	if IsAneka(h) {                                        // 16 MaintenancePeriod
		bulan := h.Ambil(JalurAnak(JalurAnak(JalurAnak(DaftarLokasiPolis, 1, "Property.RiskLocation.OccupationList"), 1,
			"AnekaList"), 1, "MaintenancePeriod"))
		lewat = SesudahTanggal(dol, TambahBulan(h.Ambil(JalurAkhirPolis), bulan))
	}
	if lewat { // 18
		h.TambahPesan(CD+"DateOfLoss", PesanDOLLebihAkhirPolis)
	}
}

// CheckDateReport = `CheckDateReport_Act` (ProteksiDataRegister_Act 4). Langkah 8, 11, 14 ber-remark.
func CheckDateReport(k *Konteks, h *Halaman) {
	lapor := h.Ambil(CD + "ReportDate")
	if SesudahTanggal(lapor, k.Hari()) { // 2-3
		h.TambahPesan(CD+"ReportDate", PesanLaporLebihHariIni)
	}
	if SesudahTanggal(h.Ambil(CD+"DateOfLoss"), lapor) { // 4-5
		h.TambahPesan(CD+"ReportDate", PesanLaporKurangDOL)
	}
	if SesudahTanggal(lapor, h.Ambil(CD+"DateReceived")) { // 6-7
		h.TambahPesan(CD+"ReportDate", PesanLaporLebihTerima)
	}
}

// CheckDateReceived = `CheckDateReceived_Act` (change Received Date; ProteksiDataRegister_Act 5). Langkah 8, 11, 14
// ber-remark.
func CheckDateReceived(k *Konteks, h *Halaman) {
	terima := h.Ambil(CD + "DateReceived")
	if SesudahTanggal(terima, k.Hari()) { // 2-3
		h.TambahPesan(CD+"DateReceived", PesanTerimaLebihHariIni)
	}
	if SesudahTanggal(h.Ambil(CD+"DateOfLoss"), terima) { // 4-5
		h.TambahPesan(CD+"DateReceived", PesanTerimaKurangDOL)
	}
	if SesudahTanggal(h.Ambil(CD+"ReportDate"), terima) { // 6-7
		h.TambahPesan(CD+"DateReceived", PesanTerimaKurangLapor)
	}
}

// ProteksiDataRegister = `ProteksiDataRegister_Act` (aktivitas FlowAction InputRegister saat Submit): langkah 1
// (CallActivityInputRegister) dijalankan services sebelumnya; 3-5 pemeriksaan tanggal bila terisi; 6-8 medan wajib.
// ⚠️ Pesan 6-8 terpasang pada medan CauseOfLoss (ARG Field langkah 6, 7, 8 - verbatim).
func ProteksiDataRegister(k *Konteks, h *Halaman) {
	if h.Ambil(CD+"DateOfLoss") != "" { // 3
		CheckDate(k, h)
	}
	if h.Ambil(CD+"ReportDate") != "" { // 4
		CheckDateReport(k, h)
	}
	if h.Ambil(CD+"DateReceived") != "" { // 5
		CheckDateReceived(k, h)
	}
	if h.Ambil(CD+"ReporterName") == "" { // 6
		h.TambahPesan(CD+"CauseOfLoss", PesanPelaporKosong)
	}
	if h.Ambil(CD+"ReporterTelp") == "" { // 7
		h.TambahPesan(CD+"CauseOfLoss", PesanTeleponKosong)
	}
	if h.Ambil(CD+"CauseOfLoss") == "" { // 8
		h.TambahPesan(CD+"CauseOfLoss", PesanSebabKosong)
	}
}

// CheckListEstimasi = `CheckListEstimasi_Act` (change Selected / Submit): Test = 1 bila ada objek terpilih; tanpa objek
// terpilih -> pesan (langkah 3). Langkah 2.2 membersihkan pesan halaman bila isinya pesan yang sama.
func CheckListEstimasi(h *Halaman) {
	h.Setel(JalurUjiCalon, "") // 1
	for _, b := range h.AmbilDaftar(DaftarCalon) {
		if b[PropDipilih] == "true" { // 2.1
			h.Setel(JalurUjiCalon, "1")
			delete(h.Pesan, CD+"Test") // 2.2
		}
	}
	if h.Ambil(JalurUjiCalon) != "1" { // 3
		h.TambahPesan(CD+"Test", PesanPilihObjek)
	}
}

// PascaRegister = Submit Input Register: pre-DT `InsertObjectItemList_DT` (objek terpilih; langkah 3-10 membuang
// ObjectList bila ada pesan atau medan wajib kosong), `CheckListEstimasi_Act`, lalu `ProteksiDataRegister_Act` dan post-DT
// `InputRegisterPostDT` (`ObjectCoverageList(..).IsKomiteApproveExGratia := ""` - daftar salinan polis, tidak disimpan:
// tanpa efek). Validate `ValidateDate` FlowAction tidak diekspor (OQ-CFI-13).
//
// `TempClaimData.ClaimData.Test` adalah halaman requestor yang ditulis CheckListEstimasi_Act setiap centang Selected
// berubah; di sini tidak disimpan, maka nilainya dihitung ulang dari centang yang sama sebelum pre-DT membacanya.
func PascaRegister(k *Konteks, h *Halaman) {
	h.Setel(JalurUjiCalon, "")
	for _, b := range h.AmbilDaftar(DaftarCalon) {
		if b[PropDipilih] == "true" {
			h.Setel(JalurUjiCalon, "1")
		}
	}
	TerapkanObjekTerpilih(k, h) // InsertObjectItemList_DT 2
	CheckListEstimasi(h)
	ProteksiDataRegister(k, h)
	if h.AdaPesan() { // InsertObjectItemList_DT 3
		h.HapusAwalanDaftar(DaftarObjek)
		return
	}
	for _, j := range []string{"DateOfLoss", "ReportDate", "DateReceived", "ReporterName", "ReporterTelp", "CauseOfLoss",
		"Currency"} { // 4-10
		if h.Ambil(CD+j) == "" {
			h.HapusAwalanDaftar(DaftarObjek)
			return
		}
	}
	if h.Ambil("EndDateRegister") == "" { // 11
		h.Setel("EndDateRegister", k.Waktu())
	}
	if h.Ambil("StartDateEstimation") == "" { // 12
		h.Setel("StartDateEstimation", k.Waktu())
	}
}

// CheckDoubleClaim = `CheckDoubleClaim_Act` (change Claim Estimate): klaim lain pada polis yang sama dengan DOL, Cause of
// Loss, dan Claim Estimate sama, dan tidak ada di CLAIMREJECTED -> IsError 1 (tombol Submit membuka ProtectDOL); selainnya
// IsError 0. Langkah 2 (polis tertulis mati -> keluar) DIBUANG (bawaan c prompt §3). Mengembalikan kunci kasus kembar
// (grid ProtectDOL `TempHistoryClaim.pxResults` .BRANCH_NAME = nomor klaim).
func CheckDoubleClaim(ctx context.Context, a Acuan, h *Halaman) ([]RiwayatKlaimPolis, error) {
	rs, err := a.RiwayatKlaimPolis(ctx, h.Ambil(JalurNoPolis))
	if err != nil {
		return nil, err
	}
	dol, kasus := h.Ambil(CD+"DateOfLoss"), h.Ambil("pyID")
	sebab, estimasi := h.Ambil(CD+"CauseOfLoss"), h.Ambil(CD+"ClaimEstimate")
	var kembar []RiwayatKlaimPolis
	for _, r := range rs { // 3
		if !(SamaTanggal(r.DateOfLoss, dol) && kasus != r.ClaimNo && sebab == r.CauseOfLoss &&
			samaAngka(estimasi, r.ClaimEstimate)) { // 3.1
			continue
		}
		ditolak, err := a.KlaimDitolak(ctx, r.KunciKasus) // 3.2-3.4.1
		if err != nil {
			return nil, err
		}
		if !ditolak {
			kembar = append(kembar, r)
		}
	}
	h.Setel(JalurIsError, nolSatu(len(kembar) > 0)) // 4-5
	return kembar, nil
}

// samaAngka - `Local.Estimasi == .TSI` (TSI = REPLACE(ClaimEstimate, ',', '.')): pembandingan nilai.
func samaAngka(a, b string) bool {
	da, ea := AngkaTeks("", a)
	db, eb := AngkaTeks("", strings.ReplaceAll(b, ",", "."))
	if ea != nil || eb != nil {
		return strings.TrimSpace(a) == strings.TrimSpace(b)
	}
	return da.Cmp(db) == 0
}

// GetReportStatus = `GetReportStatus_Act` (change Reporter Status `.ClaimData.InsuredRelationship`): 1 = ceding (satu
// ceding: QuotationData.CedingCo; CedingCoList satu baris: baris 1), 2 = SOB; alamat dari klien agen (GetLeaderReport bila
// kode agen kosong, lalu GetAddressCeding). Cabang TreatyIn (4, 6, 11) tidak terjangkau.
func GetReportStatus(k *Konteks, h *Halaman) error {
	var klien, agen string
	jumlah := len(h.AmbilDaftar(DaftarCedingCo)) // 1
	switch h.Ambil(CD + "InsuredRelationship") {
	case "1":
		if jumlah == 0 { // 2
			h.Setel(CD+"InsuredRelationshipOthers", h.Ambil(OQ+"CedingCoName"))
			klien = h.Ambil(OQ + "CedingCo")
		}
		if jumlah == 1 { // 3
			nama := h.AmbilDaftar(DaftarCedingCo)[0]["CedingCoName"]
			h.Setel(CD+"InsuredRelationshipOthers", nama)
			klien = nama
		}
	case "2": // 5
		h.Setel(CD+"InsuredRelationshipOthers", h.Ambil(OQ+"SobName"))
		klien, agen = h.Ambil(OQ+"SourceOfBusiness"), h.Ambil(OQ+"SobLeader0")
	}
	h.Hapus(CD + "ReportAddress") // 7
	return setelAlamatPelapor(k, h, klien, agen)
}

// GetCeding = `GetCeding_act` (tombol Choose grid CedingCoList, Param.Ceding = .CedingCoName).
func GetCeding(k *Konteks, h *Halaman, ceding string) error {
	h.Setel(CD+"InsuredRelationshipOthers", ceding) // 1
	return setelAlamatPelapor(k, h, ceding, "")
}

// setelAlamatPelapor = GetLeaderReport (kode agen dari nama klien bila kosong) + GetAddressCeding + alamat dirangkai
// "a, b, c, d" (GetReportStatus_Act 8-13 / GetCeding_act 2-6).
func setelAlamatPelapor(k *Konteks, h *Halaman, klien, agen string) error {
	if agen == "" {
		id, err := k.Acuan.AgenKlien(k.Ctxt(), klien)
		if err != nil {
			return err
		}
		agen = id
	}
	al, err := k.Acuan.AlamatAgen(k.Ctxt(), agen)
	if err != nil {
		return err
	}
	h.Setel(CD+"ReportAddress", al[0]+", "+al[1]+", "+al[2]+", "+al[3])
	return nil
}

// InputCurrencyValue = `InputCurrencyValueAct_Register` (change Currency): DollarCurrencyVal = kurs standar.
func InputCurrencyValue(k *Konteks, h *Halaman) error {
	v, err := k.Acuan.KursStandar(k.Ctxt(), h.Ambil(CD+"Currency"))
	if err != nil {
		return err
	}
	h.Setel(CD+"DollarCurrencyVal", v)
	return nil
}

// SetTempLocation = `SetTempLocation_Act` (change Location of Loss): langkah 1 (`pre=false` - gerbang tidak berlaku).
func SetTempLocation(h *Halaman) { h.Setel(JalurLokasiSementara, h.Ambil(CD+"Location")) }

// InputKodePos = `InputKodePos1_Act` (change Zip Code): wilayah dari kode pos (baris terakhir menimpa), Location
// dirangkai dari lokasi sementara + RW, District, City, Province.
func InputKodePos(k *Konteks, h *Halaman) error {
	h.Hapus(CD + "Location") // 1
	ws, err := k.Acuan.Wilayah(k.Ctxt(), h.Ambil(CD+"PostalCode"))
	if err != nil {
		return err
	}
	for _, w := range ws { // 4.1
		h.Setel(CD+"RW", w.RW)
		h.Setel(CD+"District", w.District)
		h.Setel(CD+"City", w.City)
		h.Setel(CD+"Province", w.Province)
		h.Setel(CD+"RWID", w.RWID)
		h.Setel(CD+"DistrictID", w.DistrictID)
		h.Setel(CD+"CityID", w.CityID)
		h.Setel(CD+"ProvinceID", w.ProvinceID)
	}
	h.Setel(CD+"Location", h.Ambil(JalurLokasiSementara)+", "+h.Ambil(CD+"RW")+", "+h.Ambil(CD+"District")+", "+
		h.Ambil(CD+"City")+", "+h.Ambil(CD+"Province")) // 5
	return nil
}

// SetAdjuster = `SetAdjsuter_act` / SetConsultant = `SetConsultant_Act`: nama dari RD BrowseAdjusterConsultant; ID atau
// nama kosong -> keduanya kosong.
func SetAdjuster(k *Konteks, h *Halaman) error {
	return setNamaAdjuster(k, h, CD+"AppointedADJID", CD+"AppointedADJ")
}

// SetConsultant - lihat SetAdjuster.
func SetConsultant(k *Konteks, h *Halaman) error {
	return setNamaAdjuster(k, h, CD+"ConsultantID", CD+"ConsultantName")
}

func setNamaAdjuster(k *Konteks, h *Halaman, jalurID, jalurNama string) error {
	id := h.Ambil(jalurID)
	nama := ""
	if id != "" {
		n, err := k.Acuan.NamaAdjuster(k.Ctxt(), id)
		if err != nil {
			return err
		}
		nama = n
	}
	h.Setel(jalurNama, nama)    // 3
	if id == "" || nama == "" { // 4
		h.Setel(jalurID, "")
		h.Setel(jalurNama, "")
	}
	return nil
}

// ---------------------------------------------------------------- katastrofe & cause of loss (kelas Work, bersama)

// SetDefNonCatastrope = `SetDefNonCatastrope_Act`.
func SetDefNonCatastrope(h *Halaman) {
	if h.Ambil(CD+"StsKatastrofe") == "Catastrophe" { // 1
		h.Setel(CD+"NonKatastrofeType", "")
	} else if h.Ambil(CD+"NonKatastrofeType") == "" { // 2
		h.Setel(CD+"NonKatastrofeType", "Claim")
	}
	if h.Ambil(CD+"NonKatastrofeType") == "Claim" { // 3
		h.Setel(CD+"KatastrofeNote", "")
		h.Setel(CD+"KatastrofeID", "")
	}
}

// SetEditCatastrope = `SetEditCatastrope` (ikon Edit / Save, Param.Action). Mode layar (tidak berkolom, pola Claim Prop).
func SetEditCatastrope(h *Halaman, aksi string) {
	switch aksi {
	case "Edit":
		h.Setel(CD+"EditCatastrope", "true")
	case "Save":
		h.Setel(CD+"EditCatastrope", "false")
	}
}

// InputCatastrope = `InputCatastrope` (param.Act Add / Cancel): form katastrofe baru (`ParamInput.CARI1`).
func InputCatastrope(h *Halaman, act string) {
	switch act {
	case "Add":
		h.Setel("ParamInput.CARI1", "1")
	case "Cancel":
		h.Setel("ParamInput.CARI1", "0")
	}
}

// SetCatastrope = `SetCatastrope_act` (Choose grid katastrofe).
func SetCatastrope(h *Halaman, id, note string) {
	h.Setel(CD+"KatastrofeNote", note)
	h.Setel(CD+"KatastrofeID", id)
}

// BarisKatastrofe - baris CATASTROPHE baru (`SaveCatasrtope_Act` langkah 1, Obj-Save langkah 3).
type BarisKatastrofe struct {
	ID, KlaimType, UserInput, Note, StsKatastrofe, NonKatastrofeType string
}

// SusunKatastrofe = SaveCatasrtope_Act langkah 1: ID "CTS-" + stempel waktu tanpa [GMT.] (`@pxReplaceAllViaRegex(
// CurrentDateTime(),"[GMT.]","")` = "yyyyMMddTHHmmss" + milidetik), NOTE huruf besar.
func SusunKatastrofe(k *Konteks, h *Halaman, note string) BarisKatastrofe {
	t := k.Sekarang.UTC()
	return BarisKatastrofe{ID: "CTS-" + t.Format("20060102T150405") + t.Format(".000")[1:], KlaimType: "NON-LIFE",
		UserInput: k.Pelaku, Note: strings.ToUpper(note), StsKatastrofe: h.Ambil(CD + "StsKatastrofe"),
		NonKatastrofeType: h.Ambil(CD + "NonKatastrofeType")}
}

// GetNameCauseofLoss = `GetNameCauseofLoss_Act` (Choose grid List Cause of Loss): `CauseOfLoss := .Info`. Baris grid RD
// BrowseCouseOfLoss_Business tidak punya `.Info` (kelas activity Int-CLAUSE) - dipakai DESCRIPTION baris terpilih dan
// ID-nya disimpan (pola Claim Prop / Non Prop, PARITAS).
func GetNameCauseofLoss(h *Halaman, deskripsi, id string) {
	h.Setel(CD+"CauseOfLoss", deskripsi)
	h.Setel(CD+"CauseOfLossID", id)
}

// ---------------------------------------------------------------- Choose Polis

// PilihPolis = `CopyNB_Act` (tautan Policy Number grid Choose Polis) + pemuatan dokumen: nomor polis dan prodke
// disimpan; IsTreatyIn menurut awalan (4-5); dokumen polis diterapkan (`TerapkanPolis`); InputParam.CARI4 = "true"
// (detail polis tampil, 12).
func PilihPolis(h *Halaman, nopolis, prodke string, dokumen []byte) error {
	if prodke == "" { // 2
		prodke = "0"
	}
	if v := IsTreatyInDari(nopolis); v != "" { // 4-5
		h.Setel(JalurIsTreatyIn, v)
	}
	if err := TerapkanPolis(h, dokumen); err != nil { // 6-7
		return err
	}
	NormalkanTanggalPolis(h)       // 11
	h.Setel(JalurNoPolis, nopolis) // 12
	h.Setel(JalurProdke, prodke)
	h.Setel("InputParam.CARI4", "true") // 12
	return nil
}

// SubmitPilihPolis = `SetInputParam_Act` (tombol Submit pop-up Choose Polis): InputParam.CARI30 (1), kronologi "Create
// New Claim" (2-3), `RemoveDataClaimData_DT` (4: menulis "" ke medan yang sudah kosong - tanpa efek).
func SubmitPilihPolis(k *Konteks, h *Halaman) {
	h.Setel("InputParam.CARI30", "true")
	k.Kronologi(h, "Create New Claim")
}

// TerapkanPolisTersimpan menulis dokumen polis lalu mengembalikan isian kasus atas medan polis yang dapat diketik
// (InputRegisterDetail: Insured Name / CedingCoName / Source Of Business hanya-baca bila polis terpilih, QQName hanya-baca
// bila sudah terisi) - nilai polis menang bila terisi.
func TerapkanPolisTersimpan(h *Halaman, dokumen []byte) error {
	simpan := map[string]string{}
	for _, j := range MedanPolisKetik {
		simpan[j] = h.Ambil(j)
	}
	if err := TerapkanPolis(h, dokumen); err != nil {
		return err
	}
	NormalkanTanggalPolis(h)
	for j, v := range simpan {
		if h.Ambil(j) == "" && v != "" {
			h.Setel(j, v)
		}
	}
	return nil
}

// MedanPolisKetik - medan halaman polis yang dapat diketik di InputRegisterDetail (berkolom di T_GENERAL_CLAIM).
var MedanPolisKetik = []string{OQ + "InsuredName", OQ + "QQName", OQ + "CedingCoName", OQ + "SobName"}
