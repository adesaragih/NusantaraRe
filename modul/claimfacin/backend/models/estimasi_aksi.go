package models

// Untuk apa berkas ini: PORT TOMBOL LAYAR INPUT ESTIMASI di luar baris estimasi - Print PLA (`SetIndexObj_Act`,
// `CreateRemarksNotePLA_Act`, `GeneratePLA`, `GeneratePLATreaty_Act`), Send to PIC Claim (`SetDisable_ACT` +
// `SetDisable_DT`), Back (`BackFromRegister` + `BackToRegister_act`; `BackToEstimasi` di layar Adjustment), pra-proses
// (`InputEstimationPre` + `SetEstimation_DT`, `SetTypePDFAdjustment` + `SetStartDateAdjustment`, `SetMOClaim_Act`), dan
// PROGRESSCLAIM / SUBPROGRESSCLAIM (`InsertProgressClaim`).

import (
	"strconv"
	"strings"
	"time"
)

// Teks dan pesan VERBATIM.
const (
	TeksPrintPLA      = "Print PLA"                                         // GeneratePLA 1
	TeksFinishEst     = "Finish Estimation"                                 // SetDisable_ACT 2
	PesanPLADulu      = "Please print the PLA before sending it to the pic" // SetDisable_ACT 6
	AwalanRemarksPLA  = "- Refer to PLA : "                                 // CreateRemarksNotePLA_Act 2
	AwalanBackFrom    = "Back From "                                        // BackFromRegister 3
	CatatanBack       = "Back"                                              // BackFromRegister 1, BackToEstimasi 1
	JalurCatatan      = "pyNote"
	JalurAktifButton  = "AktifButton"
	DaftarPLA         = "PLAList" // anak baris objek (riwayat PLA: NoPLA, PLA, TipePLA, RetroList)
	JalurKirimSurel   = "ProtectPrint.CARI50"
	DaftarPenerima    = "Email.RecipientList"
	JalurPesanSurel   = "Email.Message"
	JalurSubjekSurel  = "Email.Subject"
	JalurCekLimit     = "CekLimit.CARI1" // halaman requestor; tak pernah 1 (OQ-CFI-24)
	TreatyQSHR        = "QS HR"
	TreatySurplus     = "SURPLUS"
	JenisTreatyQS     = "10003"
	JenisTreatySurpls = "10192"
)

// ---------------------------------------------------------------- Print PLA

// BukaPLA = tombol "Print PLA" baris objek: `SetIndexObj_Act` (RemarksPLA dibuang) lalu pra-proses harness
// `CreateRemarksNotePLA_Act` (RemarksPLA = "- Refer to PLA : " + NoPLA baris PLAList ber-TipePLA 10015 terakhir).
//
// PERBAIKAN prompt §5 butir 2 (PARITAS `[penyimpangan sadar]`): SetIndexObj_Act langkah 1 menyetel `.PlaStatus := 1`
// SEBELUM pop-up dibuka, sehingga pop-up yang ditutup tanpa Submit mematikan tombol Print PLA dan meloloskan Send to PIC
// tanpa PLA. Di sini PlaStatus baru menjadi 1 sesudah GeneratePLA berhasil (`GeneratePLA` langkah 9).
func BukaPLA(h *Halaman, o int) error {
	ob, err := Objek(h, o)
	if err != nil {
		return err
	}
	delete(ob, "RemarksPLA") // SetIndexObj_Act 2
	nomor := ""
	for _, p := range h.AmbilDaftar(JalurAnak(DaftarObjek, o, DaftarPLA)) { // CreateRemarksNotePLA_Act 1
		if p["TipePLA"] == TreatyFacRetro {
			nomor = p["NoPLA"]
		}
	}
	if nomor != "" { // 2
		ob["RemarksPLA"] = AwalanRemarksPLA + nomor
	}
	if ob["ShareRetro"] == "" { // Pla_Dtl: .ShareRetro DEF=1
		ob["ShareRetro"] = "1"
	}
	return nil
}

// spreadTreaty - GeneratePLA 5.1.1: `@contains(.TreatyName,"QS HR") || @contains(.TreatyName,"SURPLUS") ||
// @contains(.TreatyType,"10192") || @equals(.TreatyType,"10003")`.
func spreadTreaty(s Baris) bool {
	return strings.Contains(s["TreatyName"], TreatyQSHR) || strings.Contains(s["TreatyName"], TreatySurplus) ||
		strings.Contains(s["TreatyType"], JenisTreatySurpls) || s["TreatyType"] == JenisTreatyQS
}

// SpreadTreatyPLA - GeneratePLA 5.1.3: spreadTreaty + TreatyType 10201 / 10208.
func SpreadTreatyPLA(s Baris) bool {
	return spreadTreaty(s) || s["TreatyType"] == "10201" || s["TreatyType"] == "10208"
}

// RencanaPLA - hasil murni GeneratePLA untuk services.
type RencanaPLA struct {
	// Treaty - item (1..n) yang memanggil GeneratePLATreaty_Act (5.1.3: spreading treaty dan `CekLimit.CARI1 = 1`).
	Treaty []int
	// FacOQ - item ber-spreading 10015 (5.1.2 `ObjectItem.GeneratePLA` - rule TIDAK diekspor: nomor PLA fac tidak
	// diterbitkan, OQ-CFI-21).
	FacOQ []int
}

// GeneratePLA = tombol Submit pop-up Pla_Dtl (`GeneratePLA`, Data-Object, Param.idxobj = o): kronologi "Print PLA"
// (1-2), PLA per spreading item (5), IsTreatyOut (6). PlaStatus = 1 (9) disetel services sesudah seluruh PLA terbit.
//
// `CekLimit.CARI1` (5.1.3 / 5.1.5) adalah halaman requestor yang hanya ditulis CheckLimitSpreadingTreaty_Act 17.4 / 17.5
// - langkah di dalam perulangan `.Adjustment` halaman primer, padahal tombol pemanggilnya (Print DLA) berada di baris
// OBJEK yang tidak punya `.Adjustment` (nol rujukan `ObjectList(n).Adjustment` di korpus). Nilainya tidak pernah 1,
// sehingga PLA treaty (dan DLA treaty, ChooseDla_Act 3.1.4) tidak pernah terbit - ditiru apa adanya, OQ-CFI-24.
func GeneratePLA(k *Konteks, h *Halaman, o int) (RencanaPLA, error) {
	var r RencanaPLA
	if _, err := Objek(h, o); err != nil {
		return r, err
	}
	k.Kronologi(h, TeksPrintPLA) // 1-2
	treatyOut := ""
	cek := h.Ambil(JalurCekLimit) == "1"
	for i := range h.AmbilDaftar(DaftarItem(o)) { // 5
		fac, trt := false, false
		for _, s := range h.AmbilDaftar(DaftarDiItem(o, i+1, AnakSpreadPolis)) { // 5.1
			if spreadTreaty(s) { // 5.1.1
				treatyOut = "1"
			}
			if s["TreatyType"] == TreatyFacRetro { // 5.1.2
				fac = true
			}
			if SpreadTreatyPLA(s) && cek { // 5.1.3
				trt = true
			}
		}
		if fac {
			r.FacOQ = append(r.FacOQ, i+1)
		}
		if trt {
			r.Treaty = append(r.Treaty, i+1)
		}
	}
	h.Setel(JalurIsTreatyOut, treatyOut) // 6
	return r, nil
}

// NomorPLATreatyAda - GeneratePLA 4.2: `CekPLATreaty.START_DATE` = `.PLA` baris PLAList objek ber-TipePLA 10192 / 10003
// terakhir ("" = nomor baru diterbitkan GeneratePLATreaty_Act 5-8).
func NomorPLATreatyAda(h *Halaman, o int) string {
	no := ""
	for _, p := range h.AmbilDaftar(JalurAnak(DaftarObjek, o, DaftarPLA)) {
		if p["TipePLA"] == JenisTreatySurpls || p["TipePLA"] == JenisTreatyQS {
			no = p["PLA"]
		}
	}
	return no
}

// TerapkanPLATreaty = GeneratePLATreaty_Act 9-13 + 19 (item i objek o) dan GeneratePLA 5.1.5: revisi PLA treaty item
// (`PLARev_Treaty` + 1, kosong -> 1), nomor revisi `nomor/rev`, baris PLAList objek, IsPicTransfer 1.
//
// Tidak dibawa (PARITAS): 14-15 menyetel TotalEstimasi / TotalGrossEstimasi item = 0 (lokal yang baru dinolkan; indeks
// `Param.idxobjlist` = ObjectID objek) - penghapusan nilai tanpa maksud yang terbaca (`[penyimpangan sadar]`); 16-18.26
// HTML -> PDF per reasuradur (aliran HTML tidak diekspor, OQ-CFI-20).
func TerapkanPLATreaty(k *Konteks, h *Halaman, o, i int, nomor string, s Baris) error {
	it, err := Item(h, o, i)
	if err != nil {
		return err
	}
	rev := 1 // 11-12
	if it["PLARev_Treaty"] != "" {
		var kk Kalkulator
		n := kk.B(it, "PLARev_Treaty")
		if err := kk.Galat(); err != nil {
			return err
		}
		v, _ := n.Int64()
		rev = int(v) + 1
	}
	it["PLARev_Treaty"] = strconv.Itoa(rev)
	revisi := nomor + "/" + strconv.Itoa(rev) // 13
	var kk Kalkulator
	nilai := kk.B(it, "TotalEstimationValueinIDR")
	retro := Baris{}
	if fr := h.AmbilDaftar(DaftarFacRetro); len(fr) > 0 {
		retro = fr[0]
	}
	total := kk.Persen(nilai, kk.B(retro, "PctShareAllObj"), kk.B(s, "SharePercentage"))
	if err := kk.Galat(); err != nil {
		return err
	}
	h.TambahBaris(JalurAnak(DaftarObjek, o, DaftarPLA), Baris{"NoPLA": revisi, "PLA": nomor, // GeneratePLA 5.1.5
		"TipePLA": s["TreatyType"], "pxCreateOperator": k.Pelaku, "pxCreateDateTime": k.Waktu()})
	h.SetelDaftar(JalurAnak(JalurAnak(DaftarObjek, o, DaftarPLA), len(h.AmbilDaftar(JalurAnak(DaftarObjek, o,
		DaftarPLA))), "RetroList"), []Baris{{"AdditionalInfo": retro["AdditionalInfo"], "Attention": retro["Attention"],
		"PctShareAllObj": retro["PctShareAllObj"], "ReinsurerID": retro["ReinsurerID"],
		"ReinsurerName": retro["ReinsurerName"], "TotalClaim": Teks(total)}})
	h.Setel(JalurIsPicTransfer, "1") // 19
	return nil
}

// SelesaiPLA = GeneratePLA langkah 9: PlaStatus objek 1 (sesudah PLA terbit - perbaikan prompt §5 butir 2).
func SelesaiPLA(h *Halaman, o int) error {
	ob, err := Objek(h, o)
	if err != nil {
		return err
	}
	ob["PlaStatus"] = "1"
	return nil
}

// ---------------------------------------------------------------- Send to PIC Claim

// SetDisable = tombol "Send to PIC Claim" (`SetDisable_ACT` + `SetDisable_DT`): pyNote kosong (1), kronologi "Finish
// Estimation" (2-3), EndDateEstimation, AktifButton 0, seluruh objek PrintFaceClaim 1 / CFS 0 (DT), IsAdjustment 1 (5),
// objek retro tanpa PLA -> pesan (7.1), objek IsPrintAccept 1 / IsKomite 1 / RemarksDLA "-" (7.2).
func SetDisable(k *Konteks, h *Halaman) {
	h.Setel(JalurCatatan, "")     // 1
	k.Kronologi(h, TeksFinishEst) // 2-3
	h.Setel("EndDateEstimation", k.Waktu())
	h.Setel(JalurAktifButton, "0")
	for _, ob := range h.AmbilDaftar(DaftarObjek) { // DT 3
		ob["PrintFaceClaim"], ob["CFS"] = "1", "0"
	}
	h.Setel(JalurIsAdjustment, "1")                 // 5
	for o, ob := range h.AmbilDaftar(DaftarObjek) { // 7
		if ob["IsFacretro"] == "1" && ob["PlaStatus"] != "1" { // 7.1
			h.TambahPesan(JalurAnak(DaftarObjek, o+1, "PlaStatus"), PesanPLADulu)
		}
		ob["IsPrintAccept"], ob["IsKomite"], ob["RemarksDLA"] = "1", "1", "-" // 7.2
	}
}

// BolehKirimPIC - syarat finishAssignment tombol Send to PIC Claim: `ObjectList(1).IsFacretro != "1" OR
// ObjectList(1).PlaStatus = "1"`.
func BolehKirimPIC(h *Halaman) bool {
	d := h.AmbilDaftar(DaftarObjek)
	if len(d) == 0 {
		return true
	}
	return d[0]["IsFacretro"] != "1" || d[0]["PlaStatus"] == "1"
}

// ---------------------------------------------------------------- Back

// BackFromRegister = tombol Back layar Input Estimasi (pre-DT `BackFromRegister`, lalu `BackToRegister_act`):
// pyNote "Back" (1), kronologi "Back From " + label langkah (2-4), isEstimation kosong (12), ObjectList dibuang bila ada
// item tanpa nama (13), ReceiverClaim(1) dibuang (BackToRegister_act 1).
//
// PERBAIKAN prompt §5 butir 3 (PARITAS `[penyimpangan sadar]`): langkah 5-11 mengisi DOL / tanggal terima / tanggal lapor
// dengan waktu sekarang, nama pelapor "Test", telepon "0812", sebab kerugian "Test", mata uang "IDR" bila kosong - nilai
// uji yang masuk data produksi. Tidak dibawa; medan wajib registrasi diperiksa lagi saat Submit Input Register.
func BackFromRegister(k *Konteks, h *Halaman) {
	h.Setel(JalurCatatan, CatatanBack) // 1
	if k.Langkah != "" {               // 2
		h.Setel("Data.CARI11", k.Langkah)
	}
	k.Kronologi(h, AwalanBackFrom+h.Ambil("Data.CARI11")) // 3-4
	h.Setel(JalurIsEstimation, "")                        // 12
	kosong := false
	for o := range h.AmbilDaftar(DaftarObjek) { // 13
		for _, it := range h.AmbilDaftar(DaftarItem(o + 1)) {
			if it["ObjectItemName"] == "" {
				kosong = true
			}
		}
	}
	if kosong { // 13.1.1.1 REMOVE pyWorkPage.ClaimData.ObjectList
		h.HapusAwalanDaftar(DaftarObjek)
	}
	h.SetelDaftar(CD+"ReceiverClaim", nil) // BackToRegister_act 1
}

// BackToEstimasi = tombol Back layar Adjustment (pre-DT `BackToEstimasi`): pyNote "Back", IsAdjustment kosong.
func BackToEstimasi(h *Halaman) {
	h.Setel(JalurCatatan, CatatanBack)
	h.Setel(JalurIsAdjustment, "")
}

// ---------------------------------------------------------------- pra-proses

// PraEstimasi = pra-proses FlowAction InputEstimasi: `InputEstimationPre` (Attachment dibuang, CheckAnyAcceptation;
// SetTypePDFAdjustment / SetchronologyKlaimFacIn / SetMOClaim_Act / InsertProgressClaim lewat services) lalu DT
// `SetEstimation_DT` (isEstimation 1, IsRegister 1, IsAdjustment kosong, pyNote kosong, indeks item turunan).
func PraEstimasi(h *Halaman) {
	h.HapusAwalanDaftar(CD + "Attachment")
	CheckAnyAcceptation(h)
	h.Setel(JalurIsEstimation, "1")
	h.Setel(JalurIsRegister, "1")
	h.Setel(JalurIsAdjustment, "")
	h.Setel(JalurCatatan, "")
	TurunkanIndeksItem(h)
}

// PraAdjustment = pra-proses FlowAction InputSurveyor: `SetTypePDFAdjustment` (Attachment dibuang, CheckAnyAcceptation;
// langkah 1-3 membuang subpohon halaman polis / salinan GISFW yang tidak disimpan) lalu DT `SetStartDateAdjustment`
// (StartDateAdjustment, IsAdjustment / isEstimation / IsRegister 1, item IsKomite 0 bila AktifButton 0 / kosong).
//
// ⚠️ Langkah 1 DT `SetStartDateAdjustment` menimpa StartDateAdjustment SETIAP layar dibuka (Pega: setiap assignment
// dibuka); di sini hanya bila kosong - pra-proses sistem ini berjalan setiap muat (PARITAS `[penyimpangan sadar]`).
func PraAdjustment(k *Konteks, h *Halaman) {
	h.HapusAwalanDaftar(CD + "Attachment")
	CheckAnyAcceptation(h)
	if h.Ambil("StartDateAdjustment") == "" {
		h.Setel("StartDateAdjustment", k.Waktu())
	}
	h.Setel(JalurIsAdjustment, "1")
	h.Setel(JalurIsEstimation, "1")
	h.Setel(JalurIsRegister, "1")
	aktif := h.Ambil(JalurAktifButton)
	for o := range h.AmbilDaftar(DaftarObjek) {
		for _, it := range h.AmbilDaftar(DaftarItem(o + 1)) {
			if aktif == "0" || aktif == "" {
				it["IsKomite"] = "0"
			}
		}
	}
	TurunkanIndeksItem(h)
}

// SetMOClaim = `SetMOClaim_Act`: MarketingData klaim dari QuotationData polis (1); bila ID terisi dan cabang kosong,
// dilengkapi baris MARKETINGOFFICER menurut ID (3-5).
func SetMOClaim(k *Konteks, h *Halaman) error {
	m := CD + "MarketingData."
	h.Setel(m+"ID", h.Ambil(OQ+"MOID"))
	h.Setel(m+"ClientID", h.Ambil(OQ+"MarketingCode"))
	h.Setel(m+"ClientName", h.Ambil(OQ+"MarketingName"))
	h.Setel(m+"TeamGroup", h.Ambil(OQ+"TeamGroup"))
	h.Setel(m+"BranchDetailID", h.Ambil(OQ+"BranchCode"))
	h.Setel(m+"BranchDetailName", h.Ambil(OQ+"BranchName"))
	if h.Ambil(m+"ID") == "" || h.Ambil(m+"BranchDetailID") != "" { // 2
		return nil
	}
	mo, _, err := k.Acuan.MarketingOfficer(k.Ctxt(), h.Ambil(m+"ID")) // 4 (tanpa baris -> pxResults(1) kosong)
	if err != nil {
		return err
	}
	for i, p := range []string{"ID", "ClientID", "ClientName", "TeamGroup", "BranchDetailID", "BranchDetailName"} {
		h.Setel(m+p, mo[i]) // 5
	}
	return nil
}

// ---------------------------------------------------------------- InsertProgressClaim

// Progres - satu panggilan PEGA_PROGRESSCLAIM (CaseID, ProdKe, DaftarObjek, StatusClaim).
type Progres struct {
	IDPega, Posisi, Status string
	Tanggal                time.Time
}

// SubProgres - satu panggilan PEGA_SUBPROGRESSCLAIM.
type SubProgres struct {
	IDProgres, IDPega, Jenis, Posisi1, Posisi2, Pengguna, Status string
	Tanggal, TindakLanjut                                        time.Time
}

// ParamProgres - parameter InsertProgressClaim cabang ADJUSTMENT (Param.CASEID / Progres1 / Progres2 / Comment).
type ParamProgres struct {
	CaseID, Progres1, Progres2, Comment string
}

// Posisi progres klaim.
const (
	PosisiRegister   = "REGISTER"
	PosisiEstimation = "ESTIMATION"
	PosisiAdjustment = "ADJUSTMENT"
	StatusOnProgress = "On Progress"
)

// RencanaProgres = `InsertProgressClaim`: cabang menurut IsRegister / isEstimation / IsAdjustment kasus. `ada` = false
// bila tidak satu cabang pun cocok.
//
// PERBAIKAN prompt §5 butir 9 (PARITAS `[penyimpangan sadar]`): cabang REGISTER langkah 2.1 menulis `ProdKe := REGISTER`
// tanpa kutip - rujukan properti kosong, sehingga PROGRESSCLAIM.POSITION register kosong; di sini "REGISTER".
func RencanaProgres(k *Konteks, h *Halaman, kasusID string, p ParamProgres) (Progres, SubProgres, bool) {
	reg, est, adj := h.Ambil(JalurIsRegister), h.Ambil(JalurIsEstimation), h.Ambil(JalurIsAdjustment)
	posisi := ""
	switch {
	case reg == "1" && est == "" && adj == "":
		posisi = PosisiRegister
	case reg == "1" && est == "1" && adj == "":
		posisi = PosisiEstimation
	case reg == "1" && est == "1" && adj == "1":
		posisi = PosisiAdjustment
	default:
		return Progres{}, SubProgres{}, false
	}
	pr := Progres{IDPega: kasusID, Posisi: posisi, Status: StatusOnProgress, Tanggal: k.Sekarang}
	sub := SubProgres{IDProgres: kasusID, IDPega: kasusID, Jenis: posisi, Posisi1: posisi, Pengguna: k.Pelaku,
		Status: "Auto Create " + map[string]string{PosisiRegister: "Register", PosisiEstimation: "Estimation",
			PosisiAdjustment: "Adjustment"}[posisi], Tanggal: k.Sekarang, TindakLanjut: k.Sekarang.AddDate(0, 0, 1)}
	if posisi == PosisiAdjustment { // 4.2
		if p.CaseID != "" {
			sub.IDPega = p.CaseID
		}
		if p.Progres1 != "" {
			sub.Posisi1 = p.Progres1
		}
		sub.Posisi2 = p.Progres2
		if p.Comment != "" {
			sub.Status = p.Comment
		}
	}
	return pr, sub, true
}
