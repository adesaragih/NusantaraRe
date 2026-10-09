package models

// Untuk apa berkas ini: PENUTUPAN DAN PENOLAKAN KLAIM - local action `PreventRejectClaim` (tombol "Close Claim" layar
// Input Adjustment: `CloseClaim` langsung atau `SendCloseClaimToKomite` TT4 "close without payment") dan
// `RejectSurveyClaim` (tombol "Reject Claim" layar Input Register: ClaimComiteeReject -> SureRejectClaim ->
// `SendRejectClaimToKomite2` TT3).
//
// ⛔ TT3 / TT4 melahirkan kasus komite TANPA adjustment, sedangkan `T_GENERAL_KOMITE.ADJUSTMENT_ID` NOT NULL dan tabel
// itu tidak punya kolom TransferType (katalog DEV 10-10-2026) - aturannya hanya dapat dibangun dengan MODIFY tabel
// bersama: OQ-CFI-27 (prompt §10), tombol "Yes" kedua jalur tampil sesuai section tetapi nonaktif dengan title OQ.
// Pemeriksaan sebelum kelahiran (pesan VERBATIM) tetap dibangun.

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
	OQKomiteTanpaAdj      = "OQ-CFI-27: kasus komite tanpa adjustment (TT3 Reject / TT4 Close Without Payment) " +
		"menuntut T_GENERAL_KOMITE.ADJUSTMENT_ID nullable dan kolom TransferType - MODIFY tabel bersama, menunggu keputusan"
)

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

// LayarTutup - Section PreventRejectClaim (local action PreventRejectClaim). `Komite.CARI1` (kasus komite TT4 lahir)
// tidak pernah terisi di sini (OQ-CFI-27), jadi bagian "Success Create Request to Committee" tidak dibangun.
func LayarTutup() []Unsur {
	alokasi := sama(JalurAlokasiSalvage, "true")
	return []Unsur{
		aksi(medan(JalurAlokasiSalvage, "Close Without Payment", KCentang), "SetelAlokasiSalvage"),
		medan(JalurTKKronologi, "Chronology", KArea),
		medan(JalurTKExtent, "Extent Of Loss", KArea),
		medan(JalurTKLiability, "Policy Liability", KArea),
		wajibU(medan(JalurTKRemarks, "Remarks", KArea)),
		tampil(bagian("", label("Are you sure want close this claim without payment?"),
			sebaris("", tombolOQ("SendCloseClaimToKomite", "Yes", OQKomiteTanpaAdj))), alokasi),
		tampil(bagian("", label("Are you sure want close this claim?"),
			sebaris("", tombol("CloseClaim", "Yes", "CloseClaim"))), func(h *Halaman) bool { return !alokasi(h) }),
	}
}

// LayarTolak - Section ClaimComiteeReject (local action RejectSurveyClaim, pra-proses SetRejectClaim_pre) dengan
// konfirmasi SureRejectClaim_section. "Yes" = SendRejectClaimToKomite2 (OQ-CFI-27).
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
		label("Are you sure want to reject this claim ?"),
		sebaris("",
			naJika(tombolOQ("SendRejectClaimToKomite2", "Yes", OQKomiteTanpaAdj), sama("IsReject", "1")),
		),
	}
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

// SendRejectClaimToKomite2 2-6 (Remark, proteksi akseptasi / estimasi belum CFS) dan SendCloseClaimToKomite 2-4 tidak
// dibangun bersama tombol "Yes"-nya (OQ-CFI-27, prompt §10): penyerahan TT3 / TT4 menuntut MODIFY tabel bersama.

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
