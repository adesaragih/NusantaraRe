package models

// Untuk apa berkas ini: TANGGA PERSETUJUAN - port `Flow/InputRealizationTreatyIn`
// (22 shape, 32 connector; INVENTARIS-XML.md bab 4), `DecisionTable/isApproved`,
// (tiket 03, 04, 10; spec §5.2,
// §5.3; AC 1-10, 84, 92).
//
// ⭐ Keputusan work owner yang mengubah STRUKTUR flow (P13): tangga tiga anak
// tangga `ReasTreatyInAdmin -> ReasTreatyInSecHead -> ReasTreatyInDeptHead`.
// `ReasTreatyInGroupLeader` dan `ReasTreatyInDirector` tidak dibangun (AC 10);
// sub-graf klaim `Decision5 -> Assignment5 -> Decision7 -> Assignment1 ->
// Decision1` tidak punya connector masuk sama sekali (INVENTARIS bab 4.2) dan
// tidak dibangun (P5).
//
// ⭐ Assignment4 dan Assignment6 ("Acceptance by Head. Treaty") memakai
// workbasket yang SAMA (`ReasTreatyInSecHead`) dan flow action yang SAMA
// (`DeptHeadTreatyIn_UW`); yang membedakan keduanya hanya teks `NBStatus` berisi
// nama orang. Ketiga `When` pemilihnya - `IsSPVCreate`, `IsSPVTreaty1`,
// `IsTreaty1` - karena itu tidak mengubah ke mana berkas menunggu, dan tidak
// dibangun: berkas menunggu POSISI (AC 92).
//
// ⭐ Cabang sesudah Sec Head menyetujui - XML lawan keputusan work owner:
//
//	XML  Decision11/4 YES -> Decision13 `ToTREATYDEPTHEAD` (pyWorkPage.LetterNo ==
//	        "TREATYINDEPTHEAD", diisi `CekLimitTreatyAcc_Act` bila |TotalPremium|
//	        x kurs > 200.000.000) -> Assignment3 Dept Head
//	        [Else] -> Decision8 `NopolisEmpty` -> Assignment3 Dept Head
//	                            [Else] -> Utility1 simpan -> Utility2 -> End3
//	WO   P13 + spec §5.3 "Atasan terima -> naik", AC 8 "Sec Head menyetujui =>
//	     berkas naik ke Dept Head. Test yang menemukan tujuan lain gagal."
//
// ⛔ PERTENTANGAN WO LAWAN XML - DIIKUTI WO (aturan prompt implementasi: bila
// keputusan work owner bertentangan dengan XML, ikuti work owner, catat di
// laporan, RALAT tiket). Sec Head yang menyetujui SELALU menaikkan berkas ke
// Dept Head; cabang batas `CekLimitTreatyAcc_Act` dan tombol "Generate" Sec
// Head (`DetailDeptHeadTreatyIn_UW`, syarat `LetterNo==''`) karena itu TIDAK
// dibangun - di XML keduanya hanya menentukan apakah Sec Head boleh
// menyelesaikan sendiri, dan AC 8 menutup kemungkinan itu. Akibat sampingnya:
// kurs `TreatyIn.CurrencyList` (JSON master, P29) tidak dibutuhkan. Tercatat di
// tiket 03 (RALAT), INVENTARIS bab 6, dan HASIL-IMPLEMENTASI.

import (
	"errors"
	"strings"
)

// ErrPosisiTakDikenal - berkas berada di posisi di luar tiga anak tangga.
var ErrPosisiTakDikenal = errors.New("models: posisi bukan salah satu dari tiga anak tangga")

// Posisi (workbasket) tangga - VERBATIM nama workbasket di XML.
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

// Nilai `pyWorkPage.Position` yang ditulis connector - VERBATIM ("4" admin,
// "5" atasan). Nilai "6" (Director) tidak dibangun.
const (
	PositionAdmin  = "4"
	PositionAtasan = "5"
)

// Nama assignment - VERBATIM `pyMOName` shape; menjadi STATUS_WORK selama
// berkas menunggu (konvensi `T_WORK_POLIS`, modul premiumlistlife).
const (
	AssignmentAdmin    = "Input Realitation"
	AssignmentSecHead  = "Acceptance by Head. Treaty"
	AssignmentDeptHead = "Acceptance by Dept. Head"
)

// Status tutup. `[keputusan work owner]` P24 (AC 5): admin menolak => berkas
// DISELESAIKAN SEBAGAI DITOLAK. Shape End3 di XML tidak memasang status apa pun;
// kedua nilai ini mengikuti status penutup Pega yang dipakai tabel bersama
// `T_WORK_POLIS` (modul premiumlistlife, `StatusPolisDitolak/Selesai`).
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

// PositionPosisi memetakan posisi ke `pyWorkPage.Position` yang ditulis connector.
func PositionPosisi(p string) string {
	if p == PosisiAdmin {
		return PositionAdmin
	}
	return PositionAtasan
}

// ---------------------------------------------------------------- isApproved

// NilaiDitolak adalah satu-satunya baris `DecisionTable/isApproved`:
// `pyWorkPage.PolicyTreatyIn.IsApproved` (kolom `text`, `=`) "0" -> "No".
const NilaiDitolak = "0"

// Disetujui = `DecisionTable/isApproved`: "0" -> No; SELAIN ITU, TERMASUK
// KOSONG -> bawaan "YES" (AC 1, 2, 84).
//
// ⛔ Pembandingan TEKS (AC 3): " 0" dan "0.0" bukan "0". ⛔ `When/isApproved`
// (menguji `= 1`, kosong = tidak disetujui) TIDAK dipakai - nol kotak Decision
// menyambung ke sana (AC 4, 63).
func Disetujui(isApproved string) bool { return isApproved != NilaiDitolak }

// ---------------------------------------------------------------- transisi

// Transisi adalah akibat satu putusan pada satu posisi.
type Transisi struct {
	// PosisiBaru - workbasket berikutnya; kosong bila berkas ditutup.
	PosisiBaru string
	// StatusTutup - StatusDitolak / StatusSelesai bila berkas ditutup.
	StatusTutup string
	// Simpan - berkas melewati Utility1 (simpan polis) dan Utility2 (layanan
	// Arasapas) sebelum End3.
	Simpan bool
	// KosongkanNBStatus - connector menulis `NBStatus = ""` (penolakan atasan).
	KosongkanNBStatus bool
	// NBStatusKePosisi - connector menulis "NB IS IN <...>'S INBOX" untuk posisi
	// tujuan; nama di dalamnya diambil dari DATA (P40), lihat `TeksNBStatus*`.
	NBStatusKePosisi bool
}

// Ditutup - berkas selesai (disetujui atau ditolak).
func (t Transisi) Ditutup() bool { return t.StatusTutup != "" }

// Langkah menghitung akibat putusan pada posisi `posisi`.
//
//	isApproved     - `PolicyTreatyIn.IsApproved` saat submit
//	adaNomorPolis  - `PolicyTreatyIn.PolicyNo` terisi (`NopolisEmpty` palsu)
func Langkah(posisi, isApproved string, adaNomorPolis bool) (Transisi, error) {
	setuju := Disetujui(isApproved)
	switch posisi {
	case PosisiAdmin:
		// Assignment2 -> Decision3 isApproved
		if !setuju {
			// Decision3 No -> End3. P24: diselesaikan sebagai ditolak (AC 5).
			return Transisi{StatusTutup: StatusDitolak}, nil
		}
		// Decision3 YES -> Decision6/12/10 -> Assignment4|6 Sec Head (AC 7)
		return Transisi{PosisiBaru: PosisiSecHead, NBStatusKePosisi: true}, nil
	case PosisiSecHead:
		// Assignment4|6 -> Decision11|4 isApproved
		if !setuju {
			// No -> Assignment2, NBStatus "" (AC 6)
			return Transisi{PosisiBaru: PosisiAdmin, KosongkanNBStatus: true}, nil
		}
		// YES -> Dept Head SELALU (AC 8, keputusan work owner; lihat kepala
		// berkas - Decision13/Decision8 XML tidak dipakai di posisi ini).
		return Transisi{PosisiBaru: PosisiDeptHead, NBStatusKePosisi: true}, nil
	case PosisiDeptHead:
		// Assignment3 -> Decision2 isApproved
		if !setuju {
			// No -> Assignment2, NBStatus "" (AC 6)
			return Transisi{PosisiBaru: PosisiAdmin, KosongkanNBStatus: true}, nil
		}
		// YES -> Decision8. ⛔ Tugas properti connector Decision2 YES
		// (Position "6", PositionNote ReasTreatyInDirector) milik posisi yang
		// dibuang P13 - tidak ditulis (AC 10).
		return sesudahDecision8(adaNomorPolis), nil
	}
	return Transisi{}, ErrPosisiTakDikenal
}

// sesudahDecision8 - `NopolisEmpty` -> Assignment3 Dept Head; selain itu simpan
// lalu selesai (AC 9: tidak ada jenjang keempat).
func sesudahDecision8(adaNomorPolis bool) Transisi {
	if !adaNomorPolis {
		// Connector Decision8 -> Assignment3 menulis PositionNote
		// "ReasTreatyInGroupLeader" - posisi yang dibuang P13; workbasket
		// assignment-nya `ReasTreatyInDeptHead`, dan itulah yang dipakai.
		return Transisi{PosisiBaru: PosisiDeptHead, NBStatusKePosisi: true}
	}
	return Transisi{StatusTutup: StatusSelesai, Simpan: true}
}

// ---------------------------------------------------------------- NBStatus

// TeksNBStatusDitolakOleh = `InboxPolicyTreatyIn_postDT` / `DeptHeadTreatyIn_UW_postDT`
// langkah "PositionNote == ReasTreatyInAdmin && IsApproved == 0":
// "NB WAS DECLINED BY  " + @toUpperCase(OperatorID.pyUserName) - DUA spasi
// sesudah BY, verbatim.
func TeksNBStatusDitolakOleh(namaTampilan string) string {
	return "NB WAS DECLINED BY  " + hurufBesar(namaTampilan)
}

// TeksNBStatusKotakMasuk = "NB IS IN " + <nama> + "'S INBOX".
//
// ⭐ P40 `[keputusan work owner]`: nama di teks pemberitahuan DIAMBIL DARI DATA.
// Di DT pasca, cabang tolak menulis `@toUpperCase(pyWorkPage.pxCreateOpName)`
// (data); connector flow menulis nama orang tertanam - di sistem baru diganti
// DATA KASUS: nama posisi tujuan, sebab berkas menunggu posisi, bukan orang
// (AC 92). Rincian di tiket 10 (RALAT).
func TeksNBStatusKotakMasuk(nama string) string {
	return "NB IS IN " + hurufBesar(nama) + "'S INBOX"
}

// hurufBesar = `@toUpperCase` (Java `toUpperCase`, Unicode).
func hurufBesar(s string) string { return strings.ToUpper(s) }
