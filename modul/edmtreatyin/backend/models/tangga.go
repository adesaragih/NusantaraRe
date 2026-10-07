package models

// Untuk apa berkas ini: TANGGA PERSETUJUAN - port `Flow/InputAddendumTreatyIn` (kelas
// `ASM-FW-GISFW-Work-EndorsementTreaty`; 16 shape, 24 connector, nol shape yatim - INVENTARIS-XML.md bab 4) dan
// `DecisionTable/isApproved`. Pola dan sebagian kode: salinan `modul/nbtreatyin/backend/models/tangga.go`.
//
// Alur XML:
//
//	Start1 {Position "4", PositionNote ReasTreatyInAdmin, FlagOnGoingPolicy "1"} -> Assignment2 "Input Realitation"
//	Assignment2 -> Decision3 isApproved: No -> End3 (tanpa simpan produksi); YES -> Decision5 IsSPVCreate ->
//	  Decision10 IsSPVTreaty1 / Decision6 IsTreaty1 -> Assignment4 | Assignment7 "Acceptance by Sec Treaty"
//	  (KEDUANYA workbasket ReasTreatyInSecHead, flow action sama - beda hanya teks NBStatus nama orang)
//	Assignment4|7 -> Decision2|11 isApproved: No -> Assignment2; YES -> Decision9 ToTREATYDEPTHEAD
//	  When (LetterNo == "TREATYINDEPTHEAD") -> Assignment1 "Acceptance by Dept. Head"; Else -> Utility1
//	Assignment1 -> Decision1 isApproved: No -> Assignment2 {Position "5"}; YES -> Utility1
//	Utility1 SaveJsonPolisTreatyInEDM_Act -> Utility2 serviceInsertArasapas_act -> End3
//
// ⛔ PENYIMPANGAN SADAR - tangga TIGA jenjang (prompt eksekusi WO 06-10-2026 bab 5 "Tangga 3 jenjang Admin ->
// Sec Head -> Dept Head"; P13; ketetapan NB AC 8): Sec Head yang menyetujui SELALU menaikkan berkas ke Dept Head.
// Cabang Decision9 Else (Sec Head menyelesaikan sendiri bila `CekLimitTreatyAcc_Act` langkah 2.4 mengosongkan
// LetterNo, |TotalPremium| x kurs <= 200.000.000) TIDAK dibangun, begitu pula `CekLimitTreatyAcc_Act` dan
// LetterNo. Tercatat di HASIL-IMPLEMENTASI dan daftar pertanyaan.
//
// ⛔ Ketiga When pemilih Sec Head (`IsSPVCreate` 2 ID operator; `IsSPVTreaty1` / `IsTreaty1` kolom telepon
// operator) tidak dibangun: kedua assignment menunggu di workbasket yang SAMA, berkas menunggu POSISI (pola NB).
//
// ⛔ Status tutup: XML tidak menulis `pyStatusWork` di mana pun (penutup `ASMForceCaseClose`
// `serviceInsertArasapas_act` langkah 18 berlabel `//`). Dipakai status penutup T_WORK_POLIS NB / premiumlistlife
// (`Resolved-Rejected` / `Resolved-Completed`) - ketetapan NB P24.

import (
	"errors"
	"strings"
)

// ErrPosisiTakDikenal - berkas berada di posisi di luar tiga anak tangga.
var ErrPosisiTakDikenal = errors.New("models: posisi bukan salah satu dari tiga anak tangga")

// Posisi (workbasket) tangga - VERBATIM nama workbasket di XML (`ToWorkBasket` keempat assignment).
const (
	PosisiAdmin    = "ReasTreatyInAdmin"
	PosisiSecHead  = "ReasTreatyInSecHead"
	PosisiDeptHead = "ReasTreatyInDeptHead"
)

// PosisiTangga adalah ketiga posisi, urut jenjang.
var PosisiTangga = []string{PosisiAdmin, PosisiSecHead, PosisiDeptHead}

// AdalahPosisiTangga - posisi itu salah satu dari tiga.
func AdalahPosisiTangga(p string) bool {
	for _, x := range PosisiTangga {
		if x == p {
			return true
		}
	}
	return false
}

// Nilai `pyWorkPage.Position` - VERBATIM connector / DT: "4" (Start1, Transition8 tolak Sec Head), "5"
// (Transition2 admin -> atasan, `InputPolicyTreatyInAddendum_preAddDT` langkah 5, Transition4 tolak Dept Head).
const (
	PositionAdmin  = "4"
	PositionAtasan = "5"
)

// Nama assignment - VERBATIM `pyMOName` shape; menjadi STATUS_WORK selama berkas menunggu.
const (
	AssignmentAdmin    = "Input Realitation"
	AssignmentSecHead  = "Acceptance by Sec Treaty"
	AssignmentDeptHead = "Acceptance by Dept. Head"
)

// Status tutup (ketetapan NB P24; T_WORK_POLIS premiumlistlife).
const (
	StatusDitolak = "Resolved-Rejected"
	StatusSelesai = "Resolved-Completed"
)

// AssignmentPosisi memetakan posisi ke nama assignment-nya.
func AssignmentPosisi(p string) string {
	switch p {
	case PosisiAdmin:
		return AssignmentAdmin
	case PosisiSecHead:
		return AssignmentSecHead
	case PosisiDeptHead:
		return AssignmentDeptHead
	}
	return ""
}

// ---------------------------------------------------------------- isApproved

// NilaiDitolak adalah satu-satunya baris `DecisionTable/isApproved`: `pyWorkPage.PolicyTreatyIn.IsApproved`
// (kolom teks, `=`) "0" -> "No".
const NilaiDitolak = "0"

// Disetujui = `DecisionTable/isApproved`: "0" -> No; SELAIN ITU, TERMASUK KOSONG -> bawaan "YES". ⛔
// Pembandingan TEKS: " 0" dan "0.0" bukan "0".
func Disetujui(isApproved string) bool { return isApproved != NilaiDitolak }

// ---------------------------------------------------------------- transisi

// Transisi adalah akibat satu putusan pada satu posisi.
type Transisi struct {
	// PosisiBaru - workbasket berikutnya; kosong bila berkas ditutup.
	PosisiBaru string
	// PositionBaru - `pyWorkPage.Position` yang ditulis connector (kosong = tidak berubah).
	PositionBaru string
	// StatusTutup - StatusDitolak / StatusSelesai bila berkas ditutup.
	StatusTutup string
	// Simpan - berkas melewati Utility1 (produksi) dan Utility2 (konversi Arasapas) sebelum End3.
	Simpan bool
	// KembaliKePembuat - penolakan atasan: berkas kembali ke Admin.
	KembaliKePembuat bool
	// NBStatusKePosisi - connector menulis "EDMT IS IN <...>'S INBOX" untuk posisi tujuan.
	NBStatusKePosisi bool
}

// Ditutup - berkas selesai (disetujui atau ditolak).
func (t Transisi) Ditutup() bool { return t.StatusTutup != "" }

// Langkah menghitung akibat putusan pada posisi `posisi`.
func Langkah(posisi, isApproved string) (Transisi, error) {
	setuju := Disetujui(isApproved)
	switch posisi {
	case PosisiAdmin:
		if !setuju {
			// Decision3 No -> End3 (Transition11) - tanpa Utility1.
			return Transisi{StatusTutup: StatusDitolak}, nil
		}
		// Decision3 YES -> Assignment4|7 (Transition2 menulis Position "5"; Transition5/17/18/19 PositionNote).
		return Transisi{PosisiBaru: PosisiSecHead, PositionBaru: PositionAtasan, NBStatusKePosisi: true}, nil
	case PosisiSecHead:
		if !setuju {
			// Decision2 No -> Assignment2 (Transition8: Position "4", PositionNote ReasTreatyInAdmin). Transition21
			// (Decision11 No, dari Assignment7) tidak menulis apa pun - kedua assignment satu posisi di sistem baru,
			// dipakai Transition8.
			return Transisi{PosisiBaru: PosisiAdmin, PositionBaru: PositionAdmin, KembaliKePembuat: true}, nil
		}
		// YES -> Dept Head SELALU (penyimpangan sadar tiga jenjang - kepala berkas).
		return Transisi{PosisiBaru: PosisiDeptHead, NBStatusKePosisi: true}, nil
	case PosisiDeptHead:
		if !setuju {
			// Decision1 No -> Assignment2 (Transition4: Position "5", PositionNote ReasTreatyInAdmin).
			return Transisi{PosisiBaru: PosisiAdmin, PositionBaru: PositionAtasan, KembaliKePembuat: true}, nil
		}
		// Decision1 YES (Transition12) -> Utility1 -> Utility2 -> End3.
		return Transisi{StatusTutup: StatusSelesai, Simpan: true}, nil
	}
	return Transisi{}, ErrPosisiTakDikenal
}

// ---------------------------------------------------------------- NBStatus

// TeksNBStatusBaru = `Activity/CreateEDMT` langkah 14: "EDM IS IN " + OperatorID.pxInsName + "'S INBOX" -
// pxInsName operator = identitas akses login pembuat (data, bukan nama tertanam).
func TeksNBStatusBaru(akunPembuat string) string {
	return "EDM IS IN " + akunPembuat + "'S INBOX"
}

// TeksNBStatusDitolakOleh = `DataTransform/InboxPolicyTreatyInAddendum_postDT` langkah 2 (PositionNote ==
// ReasTreatyInAdmin && IsApproved == "0"): "EDM WAS DECLINED BY  " + @toUpperCase(OperatorID.pyUserName) - DUA
// spasi sesudah BY, verbatim.
func TeksNBStatusDitolakOleh(namaTampilan string) string {
	return "EDM WAS DECLINED BY  " + hurufBesar(namaTampilan)
}

// TeksNBStatusKotakMasuk = connector Transition5/17/18/19 (ke Sec Head) dan Transition23 (ke Dept Head):
// "EDMT IS IN " + <nama> + "'S INBOX". ⛔ Nama orang tertanam di connector diganti DATA (ketetapan NB P40):
// `NamaKotakMasuk` pemegang workbasket tujuan.
func TeksNBStatusKotakMasuk(nama string) string {
	return "EDMT IS IN " + hurufBesar(nama) + "'S INBOX"
}

// TeksNBStatusKembali - penolakan atasan (Transition8/21/4). XML tidak menulis NBStatus (teks "EDMT IS IN
// <atasan>'S INBOX" tertinggal walau berkas sudah di Admin). ⛔ PENYIMPANGAN SADAR - ketetapan NB
// `[keputusan work owner 06-10-2026]` "NBStatus tidak pernah kosong / petunjuk berkas ada di siapa": nama
// pembuat berkas, dengan awalan EDM (`CreateEDMT` langkah 14).
func TeksNBStatusKembali(namaPembuat string) string {
	return "EDM IS IN " + hurufBesar(namaPembuat) + "'S INBOX"
}

// PemegangKotakMasuk - pemegang AKTIF di luar divisi IT workbasket tujuan (`M_LOGIN_GO_WORKBASKET` x
// `M_LOGIN_GO`): jumlahnya, nama akun bila hanya satu, dan nama workbasket (`M_WORKBASKET.NAME`). Pola NB.
type PemegangKotakMasuk struct {
	Jumlah         int
	NamaAkun       string
	NamaWorkbasket string
}

// NamaKotakMasuk - nama di NBStatus (pola NB `[keputusan work owner 06-10-2026]`): tepat satu pemegang aktif =
// nama akunnya; nol atau lebih dari satu = nama workbasket; nama kosong jatuh ke berikutnya, terakhir ID
// workbasket.
func NamaKotakMasuk(workbasket string, p PemegangKotakMasuk) string {
	if p.Jumlah == 1 && strings.TrimSpace(p.NamaAkun) != "" {
		return p.NamaAkun
	}
	if strings.TrimSpace(p.NamaWorkbasket) != "" {
		return p.NamaWorkbasket
	}
	return workbasket
}

// hurufBesar = `@toUpperCase` (Java `toUpperCase`, Unicode).
func hurufBesar(s string) string { return strings.ToUpper(s) }
