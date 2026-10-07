package models

// Untuk apa berkas ini: KASUS - bentuk satu work object Claim Prop di T_WORK_CLAIM, tahap alurnya, dan konstanta
// yang dibaca lintas lapisan.
//
// `[terverifikasi]` `Flow/Flow_TreatyIn.xml` (kelas `ASM-FW-GCNMFW-Work-ClaimTreaty`):
//
//	Start1 -> Assignment2 "Outstanding Claim"  (ToCurrentOperator - worklist pembuat; FlowAction OutstandingClaim)
//	       -> Assignment1 "Input Acceptation"  (ToWorkbasket; FlowAction InputAcceptation)
//	       -> Decision3 IsBackStage (.pyNote = "Back": kembali ke Assignment2 | Else: End1 Resolved-Completed)
//
// Penutupan yang terbaca di korpus berada di luar shape End: tombol Close Claim -> `CloseClaimProp` langkah 10
// (`ASMForceCaseClose`, WorkStatus Resolved-Completed).

import (
	"errors"
	"strings"
	"time"
)

// KelasKasus - `pxObjClass` kasus Claim Prop; kunci CLASS penghitung nomor (`SaveOutstanding_Act` langkah 16:
// `ParamSeq.CARI1 = pyWorkPage.pxObjClass`).
const KelasKasus = "ASM-FW-GCNMFW-Work-ClaimTreaty"

// KelasKasusKunci - awalan pzInsKey Pega kasus lama (`ASM-FW-GCNMFW-WORK CLMP-n`, katalog DEV OS_AKSEPTASI_KLAIM).
const KelasKasusKunci = "ASM-FW-GCNMFW-WORK"

// Awalan pengenal T_WORK_CLAIM lini PROP (STRUKTUR T3, keputusan work owner 18-09-2026). Awalan kasus komite TKMT-
// tidak dipakai: penyerahan ke komite = OQ-CP-16.
const (
	AwalanKlaim = "CLMP-"
	// LiniProp - T_WORK_CLAIM.LINI baris klaim.
	LiniProp = "PROP"
)

// Tahap kasus - T_WORK_CLAIM.TAHAP. Diisi NAMA FLOWACTION assignment (`Flow_TreatyIn` connector Transition4 /
// Transition2), bukan pyMOName shape-nya.
//
// ⚠️ SEBABNYA: pyMOName Assignment2 berbunyi "Outstanding Claim" - sama persis dengan tahap Claim Life
// (`claimlife/backend/models/tahap.go`), dan kotak masuk Claim Life menyaring `NVL(w.TAHAP, …) = :tahap` TANPA
// saringan LINI. Nilai FlowAction ("OutstandingClaim") berasal dari XML yang sama dan tidak bertabrakan, sehingga baris
// PROP tidak pernah muncul di kotak masuk Claim Life. Saringan LINI di modul itu = OQ work owner (tidak disunting).
const (
	TahapOutstanding = "OutstandingClaim"
	TahapAcceptation = "InputAcceptation"
)

// Label assignment - pyMOName shape Flow_TreatyIn, tampil di daftar kerja.
var LabelTahap = map[string]string{
	TahapOutstanding: "Outstanding Claim",
	TahapAcceptation: "Input Acceptation",
}

// WorkbasketAcceptation - pemegang Assignment1. `[keputusan work owner 07-10-2026]` `ReasKlaimTeknik` (ada di
// M_WORKBASKET), menggantikan `TreatyinPNCTeknik` XML yang tidak ada di DEV.
const WorkbasketAcceptation = "ReasKlaimTeknik"

// StatusSelesai - satu-satunya pyWorkStatus yang ditulis korpus (End1 dan CloseClaimProp langkah 10).
const StatusSelesai = "Resolved-Completed"

// Sumber baris T_GENERAL_CLAIM.SUMBER.
const (
	SumberGo   = "GO"
	SumberPega = "PEGA"
)

// Kasus adalah satu baris T_WORK_CLAIM klaim PROP.
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

// Tertutup - kasus sudah Resolved-Completed. Dibandingkan PERSIS.
func (k Kasus) Tertutup() bool { return k.StatusWork == StatusSelesai }

// KunciInstans - padanan `pyWorkPage.pzInsKey` kasus sistem baru di tabel proyeksi warisan (OS_AKSEPTASI_KLAIM.CASEID,
// JSON_KLAIM.IDPEGA, DOCUMENT_CLAIM.IDPEGA, DIRECTTOKASIR_LOG.IDPEGA, MONITORING_KLAIM_LOG.IDPEGA): ID T_WORK_CLAIM
// APA ADANYA. Diterapkan sama dengan keputusan work owner 06-10-2026 untuk NB Treaty In ("INTINYA KEY DARI T_WORK_POLIS
// JANGAN DI UBAH") - OQ penegasan untuk klaim di laporan.
func KunciInstans(id string) string { return id }

// KunciPegaLama - pzInsKey Pega kasus lama berpengenal `id` ("ASM-FW-GCNMFW-WORK CLMP-4894"). Pembaca tabel proyeksi
// mencari kedua bentuk supaya klaim hasil pemuat data lama tetap menemukan barisnya.
func KunciPegaLama(id string) string { return KelasKasusKunci + " " + id }

// IDDariKunciPega membaca pyID dari pzInsKey Pega ("ASM-FW-GCNMFW-WORK CLMP-4894" -> "CLMP-4894").
func IDDariKunciPega(kunci string) string {
	kunci = strings.TrimSpace(kunci)
	if i := strings.LastIndex(kunci, " "); i >= 0 {
		return kunci[i+1:]
	}
	return kunci
}

// Galat bersama models.
var (
	// ErrBarisTidakAda - indeks baris di luar daftar.
	ErrBarisTidakAda = errors.New("models: baris tidak ada")
	// ErrBarisBeku - baris data lama (IsOldData "Yes") tidak boleh diubah atau dihapus.
	ErrBarisBeku = errors.New("models: baris data lama tidak boleh diubah")
)
