package models

// Untuk apa berkas ini: PENUTUPAN DAN PENOLAKAN KLAIM - local action `PreventRejectClaim` (tombol "Close Claim" layar
// Input Adjustment: `CloseClaim` langsung atau `SendCloseClaimToKomite` TT4 "close without payment") dan
// `RejectSurveyClaim` (tombol "Reject Claim" layar Input Register: ClaimComiteeReject -> SureRejectClaim ->
// `SendRejectClaimToKomite2` TT3).
//
// TT3 / TT4 melahirkan kasus komite TANPA adjustment (keputusan work owner 10-10-2026 KCF-03: `T_GENERAL_KOMITE.
// ADJUSTMENT_ID` boleh kosong, kolom `TRANSFER_TYPE` - migrasi komiteclaimfacin 641 / 642), satu tingkat
// `ReasClaimDeptHead` (pengganti akun + email orang tertulis mati di 7.2 / 7.3, prompt tahap 2 §5 butir 7). Keputusannya
// milik modul Komite Claim Fac In.

import (
	"strconv"
	"strings"
)

// Teks dan pesan VERBATIM.
const (
	PesanDLASebelumTutup  = "Please Print DLA Before Close This Claim"              // CloseClaim 4
	PesanTutupDiKomite    = "Can not close claim, there is adjustment in committee" // ValidationAdjustmentKomite 1
	PesanTanpaPembayaran  = "there is no payment for this claim, please upload a supporting file with the category Close Claim"
	PesanEstimasiBelumNol = "Total estimation is not null"
	PesanBelumCetakAksep  = "Can not close claim because there are adjustments that have not been printed acceptances"
	PesanKasirBelumMasuk  = "Cannot close claim, there is a direct to cashier that has not been successful."
	TeksTutupKlaim        = "Finish Adjustment (Close Claim)" // CloseClaim 9
	KategoriTutupKlaim    = "CloseClaim"                      // ValidationAdjustmentKomite 3.1.1
	StsOSTutup            = "4"                               // CloseClaim 8.1 CARI10
	// SendRejectClaimToKomite2 4 (VERBATIM).
	PesanSudahAkseptasi = "Cannot be rejected because there is already an acceptance"
	// SureRejectClaim_section LS14 / PreventRejectClaim LS28 (`Komite.CARI1 != ''`).
	TeksSuksesKomite = "Success Create Request to Committee"
	// Kronologi 7.6 (VERBATIM): "Request Reject claim " / "Request close claim without payment " + KMT.
	AwalanTolakKomite = "Request Reject claim "
	AwalanTutupKomite = "Request close claim without payment "
	// Tangga TT3 / TT4 satu tingkat (7.2 `IDKomite := "Claim Dept. Head"`, KCF-03 workbasket).
	JabatanTutupKomite    = "Claim Dept. Head"
	WorkbasketTutupKomite = "ReasClaimDeptHead"
	// PesanTutupKomiteGanda - permintaan TT3 / TT4 kedua selagi yang pertama menunggu komite (`[penyimpangan sadar]`,
	// pola satu adjustment per KMT tahap 1: Pega melahirkan kasus komite kedua). Bukan VERBATIM.
	PesanTutupKomiteGanda = "A reject / close request for this claim is already waiting for the committee"
)

// JalurKomiteBaru - `Komite.CARI1` (7.5 `@substring(pxCoveredInsKeys(<last>),19)`): kasus komite TT3 / TT4 yang baru
// lahir; halaman requestor, dibawa ke layar sesudah aksi (ModeLayar), tidak pernah diterima dari kiriman layar.
const JalurKomiteBaru = "Komite.CARI1"

// PesanHapusEstimasi - SendRejectClaimToKomite2 5.1.1.1 (VERBATIM; "Location" = indeks ITEM objek, `local.IdxObj` diisi
// subscript ObjectItemList di 5.1).
func PesanHapusEstimasi(e, i int) string {
	return "Please Delete Estimationlist " + strconv.Itoa(e) + " in Location " + strconv.Itoa(i)
}

// PesanAdjustmentDiKomite - SendCloseClaimToKomite 5.1.1.1 (VERBATIM).
func PesanAdjustmentDiKomite(a int) string {
	return "Can not close claim, there is adjustment " + strconv.Itoa(a) + " in comitee!"
}

// belumKomite - Komite.CARI1 kosong (kasus komite belum lahir di permintaan ini).
func belumKomite(h *Halaman) bool { return h.Ambil(JalurKomiteBaru) == "" }

// Jalur pop-up penutupan / penolakan (halaman requestor `TempCommiteClaim`, tidak disimpan).
const (
	JalurAlokasiSalvage = "TempCommiteClaim.AllocationShareSalvage"
	JalurTKRemarks      = "TempCommiteClaim.Remarks"
	JalurTKKronologi    = "TempCommiteClaim.CircumtansesCouseOfLoss"
	JalurTKExtent       = "TempCommiteClaim.ExtentOfLoss"
	JalurTKLiability    = "TempCommiteClaim.LegalLiability"
	JalurTKTipe         = "TempCommiteClaim.TypeComentAnalysis"
	JalurTKInisial      = "TempCommiteClaim.Initial"
)

// IsianPopUpTutup - isian teks pop-up ClaimComiteeReject LS4 / PreventRejectClaim LS23 (halaman requestor
// `TempCommiteClaim`, tidak disimpan). Di Pega nilainya tetap di clipboard sesudah "Yes" (ClaimComiteeReject LS4 tampil
// selalu), jadi layar balasan aksi membawanya (BawaSementara) - bukan penanda mode: tidak dikirim balik sebagai `mode`.
var IsianPopUpTutup = []string{JalurTKKronologi, JalurTKExtent, JalurTKLiability, JalurTKRemarks}

// LayarTutup - Section PreventRejectClaim (local action PreventRejectClaim): LS23 isian + konfirmasi selagi
// `Komite.CARI1 = ”`; LS28 "Success Create Request to Committee" sesudah kasus komite TT4 lahir.
func LayarTutup() []Unsur {
	alokasi := sama(JalurAlokasiSalvage, "true")
	sudah := func(h *Halaman) bool { return !belumKomite(h) }
	return []Unsur{
		aksi(medan(JalurAlokasiSalvage, "Close Without Payment", KCentang), "SetelAlokasiSalvage"), // LS21 (selalu)
		tampil(medan(JalurTKKronologi, "Chronology", KArea), belumKomite),                          // LS23
		tampil(medan(JalurTKExtent, "Extent Of Loss", KArea), belumKomite),
		tampil(medan(JalurTKLiability, "Policy Liability", KArea), belumKomite),
		tampil(wajibU(medan(JalurTKRemarks, "Remarks", KArea)), belumKomite),
		tampil(tampil(bagian("", label("Are you sure want close this claim without payment?"),
			sebaris("", tombol("SendCloseClaimToKomite", "Yes", "SendCloseClaimToKomite"))), alokasi), belumKomite),
		tampil(tampil(bagian("", label("Are you sure want close this claim?"),
			sebaris("", tombol("CloseClaim", "Yes", "CloseClaim"))), func(h *Halaman) bool { return !alokasi(h) }),
			belumKomite),
		tampil(label(TeksSuksesKomite), sudah), // LS28
	}
}

// LayarTolak - Section ClaimComiteeReject (local action RejectSurveyClaim, pra-proses SetRejectClaim_pre) dengan
// konfirmasi SureRejectClaim_section (LS14: tanya + "Yes" selagi `Komite.CARI1 = ”`, sesudahnya "Success Create Request
// to Committee").
func LayarTolak() []Unsur {
	return []Unsur{
		sebaris("",
			ro(medan("TempCommiteClaim.DateOfComitee", "Date Reject", KTanggal)),
			ro(medan(JalurTKInisial, "PIC Name", KTeks)),
			sumber(ro(medan(JalurTKTipe, "Tipe Analisis", KPilih)), kode("TypeComentAnalysis")),
		),
		medan(JalurTKKronologi, "Chronology", KArea),
		medan(JalurTKExtent, "Extent Of Loss", KArea),
		medan(JalurTKLiability, "Policy Liability", KArea),
		wajibU(medan(JalurTKRemarks, "Remarks", KArea)),
		tampil(label("Are you sure want to reject this claim ?"), belumKomite),
		tampil(sebaris("",
			naJika(tombol("SendRejectClaimToKomite2", "Yes", "SendRejectClaimToKomite2"), sama("IsReject", "1")),
		), belumKomite),
		tampil(label(TeksSuksesKomite), func(h *Halaman) bool { return !belumKomite(h) }),
	}
}

// PeriksaTolakKomite = SendRejectClaimToKomite2 4-6: sudah ada akseptasi; estimasi belum dicetak Face Claim
// (`.PrintFaceClaim == ""`) harus dihapus. `Local.Error` DITIMPA setiap temuan - hanya pesan TERAKHIR yang tampil (6).
func PeriksaTolakKomite(h *Halaman) string {
	galat := ""
	if h.Ambil(JalurIsAnyAccept) == "1" { // 4
		galat = PesanSudahAkseptasi
	}
	for o := range h.AmbilDaftar(DaftarObjek) { // 5
		for i := range h.AmbilDaftar(DaftarItem(o + 1)) { // 5.1 local.IdxObj = subscript item
			for e, est := range h.AmbilDaftar(DaftarDiItem(o+1, i+1, AnakEstimasi)) {
				if est["PrintFaceClaim"] == "" { // 5.1.1.1
					galat = PesanHapusEstimasi(e+1, i+1)
				}
			}
		}
	}
	return galat
}

// PeriksaTutupKomite = SendCloseClaimToKomite 5-6: adjustment ber-AcceptanceStatus kosong menahan penutupan tanpa
// pembayaran; pesan TERAKHIR yang tampil.
func PeriksaTutupKomite(h *Halaman) string {
	galat := ""
	for o := range h.AmbilDaftar(DaftarObjek) {
		for i := range h.AmbilDaftar(DaftarItem(o + 1)) {
			for a, b := range h.AmbilDaftar(DaftarAdj(o+1, i+1)) {
				if b["AcceptanceStatus"] == "" { // 5.1.1.1
					galat = PesanAdjustmentDiKomite(a + 1)
				}
			}
		}
	}
	return galat
}

// SalinCatatanTutup = SendRejectClaimToKomite2 / SendCloseClaimToKomite 2: Remark dan Remark_Close := Remarks pop-up.
// `[penyimpangan sadar]` (PARITAS): di TT3 Pega langkah 3 Obj-Open-By-Handle membuka ulang pyWorkPage sehingga isian 2
// hilang sebelum 7.8 Obj-Save; di sini disimpan - Remark klaim satu-satunya tempat teks komite TT3 / TT4 tersimpan
// (`Komite.Remarks` kasus komite = `CLAIMREJECTED.REMARK`, KomitePost_Reject S14.1).
func SalinCatatanTutup(h *Halaman) {
	r := h.Ambil(JalurTKRemarks)
	h.Setel(CD+"Remark", r)
	h.Setel(CD+"Remark_Close", r)
}

// TandaiKirimTutup = 7.5-7.7 sesudah kasus komite `kmt` lahir: Komite.CARI1, kronologi "Request ... " + KMT.
func TandaiKirimTutup(k *Konteks, h *Halaman, transfer, kmt string) {
	h.Setel(JalurKomiteBaru, kmt)
	awal := AwalanTolakKomite
	if transfer == TransferTutup {
		awal = AwalanTutupKomite
	}
	k.Kronologi(h, awal+kmt)
}

func init() {
	KodePilihan["TypeComentAnalysis"] = []string{"1", "5"} // SetRejectClaim_Cancel 2 / SetRejectClaim_pre 2
}

// SetRejectClaimPre = pra-proses RejectSurveyClaim (`SetRejectClaim_pre` 2): Tipe Analisis 5, PIC = pembuat kasus,
// tanggal sekarang.
func SetRejectClaimPre(k *Konteks, h *Halaman, pembuat string) {
	h.Setel(JalurTKTipe, "5")
	h.Setel(JalurTKInisial, pembuat)
	h.Setel("TempCommiteClaim.DateOfComitee", k.Hari())
}

// ValidasiTutup = ValidationAdjustmentKomite (CloseClaim 1). `lampiranTutup` = cacah dokumen klaim berkategori
// "CloseClaim" (`ClaimData.Attachment.pyCategory`, lampiran sistem baru). Pesan dipasang di halaman.
func ValidasiTutup(h *Halaman, lampiranTutup int) error {
	var kk Kalkulator
	adaAdj, menunggu, belumCetak, kasir := 0, 0, 0, 0
	sumEst := ""
	for o := range h.AmbilDaftar(DaftarObjek) { // 2
		for i, it := range h.AmbilDaftar(DaftarItem(o + 1)) { // 2.1
			adj := h.AmbilDaftar(DaftarAdj(o+1, i+1))
			adaAdj += len(adj) // 2.1.1
			if v := it["TotalEstimationValueinIDR"]; v != "" {
				if sumEst == "" {
					sumEst = "0"
				}
				sumEst = Teks(kk.Tambah(kk.Teks("CountSumEstimasi", sumEst), kk.Teks("TotalEstimationValueinIDR", v)))
			}
			for a, b := range adj { // 2.1.2
				if b["AcceptanceStatus"] == "0" { // 2.1.2.1
					menunggu++
				}
				if b["IsPrintAccept"] == "" && b["AcceptanceStatus"] == "1" { // 2.1.2.2
					belumCetak++
				}
				if b["AcceptanceStatus"] == "1" && b["DirectToKasir"] == "true" && b["StatusKasir"] != StatusKasirSukses {
					kasir++ // 2.1.2.4
				}
				menunggu += KomiteMenunggu(h.AmbilDaftar(DaftarDiAdj(o+1, i+1, a+1, AnakKomiteAdj))) // 2.1.2.5
			}
		}
	}
	if err := kk.Galat(); err != nil {
		return err
	}
	if menunggu > 0 { // 4
		h.TambahPesan("", PesanTutupDiKomite)
		h.Setel(JalurIsError, strconv.Itoa(menunggu))
	}
	if adaAdj == 0 && lampiranTutup == 0 { // 5
		h.TambahPesan(JalurIsError, PesanTanpaPembayaran)
	}
	if adaAdj == 0 && !(sumEst == "" || AngkaNol(sumEst)) { // 7
		h.TambahPesan(JalurIsError, PesanEstimasiBelumNol)
	}
	if belumCetak > 0 { // 10
		h.TambahPesan("", PesanBelumCetakAksep)
	}
	if kasir > 0 { // 12
		h.TambahPesan("", PesanKasirBelumMasuk)
	}
	return nil
}

// CekDLATutup = CloseClaim 4-5: Remark / Remark_Close = Remarks pop-up; setiap objek (sesudah ada adjustment di objek
// mana pun sebelumnya - `Local.CountAdj` tidak pernah dinolkan) retro ber-RemarksDLA kosong -> pesan.
func CekDLATutup(h *Halaman) {
	h.Setel(CD+"Remark", h.Ambil(JalurTKRemarks))
	h.Setel(CD+"Remark_Close", h.Ambil(JalurTKRemarks))
	n := 0
	for o, ob := range h.AmbilDaftar(DaftarObjek) {
		for i := range h.AmbilDaftar(DaftarItem(o + 1)) {
			n += len(h.AmbilDaftar(DaftarAdj(o+1, i+1)))
		}
		if n > 0 && ob["IsFacretro"] == "1" && strings.TrimSpace(ob["RemarksDLA"]) == "" {
			h.TambahPesan("", PesanDLASebelumTutup)
		}
	}
}

// BarisOSTutup = CloseClaim 8 (SaveOSClaim_SQL -> PEGA_JSON_OS_AKSEP_KLAIM, CARI10 "4"): DATA_JSON = halaman
// `InputParamOs` (CauseOfLoss, CauseOfLossID, NoClaim, IDMasterTreaty - TreatyInMaster tidak ada di Fac In) dalam format
// GetPageJSONString, kelas `ASM-FW-GCNMFW-Data-osAkseptasi` (2.700 baris STS 4 `CLM-` DEV, 10-10-2026). STS_DLA kosong
// (CARI16 / CARI17 tidak diisi).
func BarisOSTutup(h *Halaman, kasusID string) BarisOS {
	p := map[string]string{"pxObjClass": KelasOSAkseptasi, "CauseOfLoss": h.Ambil(CD + "CauseOfLoss"),
		"CauseOfLossID": h.Ambil(CD + "CauseOfLossID"), "NoClaim": h.Ambil(CD + "NoClaim")}
	return BarisOS{CaseID: kasusID, NoClaim: h.Ambil(CD + "NoClaim"), NoPolis: h.Ambil(JalurNoPolis),
		StsReject: StsOSTutup, DataJSON: JSONHalamanPega(HalamanJSON{Nilai: p})}
}

// SelesaiTutup = CloseClaim 9-10: kronologi "Finish Adjustment (Close Claim)".
func SelesaiTutup(k *Konteks, h *Halaman) { k.Kronologi(h, TeksTutupKlaim) }
