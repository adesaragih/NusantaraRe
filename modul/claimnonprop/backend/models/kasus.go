package models

// Untuk apa berkas ini: KASUS - bentuk satu work object Claim Non Prop di T_WORK_CLAIM, tahap alurnya, dan konstanta
// yang dibaca lintas lapisan.
//
// `[terverifikasi]` `Claim Non Prop/Flow/Flow_TreatyIn.xml` (kelas `ASM-FW-GCNMFW-Work-ClaimTreatyNonProp`):
//
//	Start -> Assignment2 "Outstanding Claim"  (worklist pembuat; FlowAction OutstandingClaim)
//	      -> Assignment1 "Input Acceptation"  (workbasket TreatyinPNCTeknik; FlowAction InputAcceptation)
//	      -> Decision IsBackStage (.pyNote = "Back" -> Assignment2 | Else -> End Resolved-Completed)
//
// Penutupan yang terbaca di korpus: Close Claim (`CloseClaimTNonProp` langkah 9 `ASMForceCaseClose` Resolved-Completed)
// dan Close Without Payment (kasus komite satu tingkat, tahap 2). Jalur Back tidak dibangun (OQ-CNP-09); Input
// Acceptation tanpa Submit (OQ-CNP-10).

import (
	"errors"
	"strings"
	"time"
)

// KelasKasus - `pxObjClass` kasus Claim Non Prop; kunci CLASS penghitung nomor (`SaveDataToOSAksep_Act` langkah
// 15.6.6: `ParamSeq.CARI1 = pyWorkPage.pxObjClass`).
const KelasKasus = "ASM-FW-GCNMFW-Work-ClaimTreatyNonProp"

// KelasKasusKunci - awalan pzInsKey Pega kasus lama (`ASM-FW-GCNMFW-WORK CLMNP-n`, DEV OS_AKSEPTASI_KLAIM 09-10-2026).
const KelasKasusKunci = "ASM-FW-GCNMFW-WORK"

// Awalan pengenal T_WORK_CLAIM lini NONPROP (prompt §6 butir 1; `When/IsCLMNP.xml`, awalan kasus komite Pega KMTNP-).
const (
	AwalanKlaim  = "CLMNP-"
	AwalanKomite = "KMTNP-"
	// LiniNonProp - T_WORK_CLAIM.LINI kedua baris (klaim dan komite).
	LiniNonProp = "NONPROP"
)

// Tahap kasus - T_WORK_CLAIM.TAHAP. Diisi NAMA FLOWACTION assignment (`Flow_TreatyIn` Assignment2 / Assignment1), pola
// Claim Prop: nilai FlowAction tidak bertabrakan dengan tahap Claim Life ("Outstanding Claim"), dan kotak masuk Claim
// Prop / Komite Claim Prop menyaring LINI (bukti di docs/PARITAS.md bab Kotak masuk modul lain).
const (
	TahapOutstanding = "OutstandingClaim"
	TahapAcceptation = "InputAcceptation"
	// TahapKomite - kasus komite `KMTNP-` (flow `KomiteTreaty_Flow`, korpus Komite Claim Non Prop).
	TahapKomite = "KomiteTreaty_Flow"
)

// LabelTahap - pyMOName shape Flow_TreatyIn, tampil di daftar kerja.
var LabelTahap = map[string]string{
	TahapOutstanding: "Outstanding Claim",
	TahapAcceptation: "Input Acceptation",
}

// WorkbasketAcceptation - pemegang Assignment1. Bawaan OQ-CNP-08 `ReasKlaimTeknik` (ada di M_WORKBASKET DEV),
// menggantikan `TreatyinPNCTeknik` XML yang tidak ada di DEV - sama dengan Claim Prop.
const WorkbasketAcceptation = "ReasKlaimTeknik"

// StatusSelesai - pyWorkStatus penutupan (`CloseClaimTNonProp` langkah 9, shape End Flow_TreatyIn).
const StatusSelesai = "Resolved-Completed"

// Sumber baris T_GENERAL_CLAIM.SUMBER.
const SumberGo = "GO"

// Kasus adalah satu baris T_WORK_CLAIM klaim NONPROP.
type Kasus struct {
	ID          string    `json:"id"`
	Tahap       string    `json:"tahap"`
	Posisi      string    `json:"posisi"`
	StatusWork  string    `json:"statusWork"`
	PembuatID   string    `json:"pembuatId"`
	PembuatNama string    `json:"pembuatNama"`
	TglCreate   time.Time `json:"tglCreate"`
	TglUpdate   time.Time `json:"tglUpdate"`
	Sumber      string    `json:"sumber"`
}

// Tertutup - kasus sudah berstatus (Resolved-*). STATUS_WORK hanya diisi saat kasus ditutup.
func (k Kasus) Tertutup() bool { return strings.TrimSpace(k.StatusWork) != "" }

// KunciInstans - padanan `pyWorkPage.pzInsKey` kasus sistem baru di tabel proyeksi warisan (OS_AKSEPTASI_KLAIM.CASEID,
// JSON_KLAIM.IDPEGA): ID T_WORK_CLAIM APA ADANYA (`CLMNP-000001`) - keputusan Claim Prop yang berlaku (prompt §6 butir 4).
func KunciInstans(id string) string { return id }

// KunciPegaLama - pzInsKey Pega kasus lama berpengenal `id` ("ASM-FW-GCNMFW-WORK CLMNP-3998"). Pembaca tabel proyeksi
// mencari kedua bentuk kunci (prompt §6 butir 4).
func KunciPegaLama(id string) string { return KelasKasusKunci + " " + id }

// Galat bersama models.
var (
	// ErrBarisTidakAda - indeks baris di luar daftar.
	ErrBarisTidakAda = errors.New("models: baris tidak ada")
)

// IDDariKunciPega membaca pyID dari pzInsKey Pega ("ASM-FW-GCNMFW-WORK CLMNP-3998" -> "CLMNP-3998"); kunci sistem baru
// (ID apa adanya) dikembalikan utuh.
func IDDariKunciPega(kunci string) string {
	if i := strings.LastIndex(kunci, " "); i >= 0 {
		return kunci[i+1:]
	}
	return kunci
}
