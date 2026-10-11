package models

// Untuk apa berkas ini: KASUS - bentuk satu work object Claim Fac In di T_WORK_CLAIM, tahap alurnya, dan konstanta yang
// dibaca lintas lapisan.
//
// `[terverifikasi]` `Claim Fac In/Flow/Register_Flow.xml` (ruleset ADESAMUEL@ 01-01-01 - checkout privat developer,
// prompt §1; kelas `ASM-FW-GCNMFW-Work-PNC`):
//
//	Start1 -> Decision4 "B2B" (When IsSPK -> Assignment3 | Else -> Assignment1)
//	Assignment1 "Input Register"   (worklist operator saat ini; FlowAction InputRegister)      -> Assignment7
//	Assignment7 "Input Estimasi"   (worklist, Custom;           FlowAction InputEstimasi)      -> Decision5
//	Decision5 IsBack (When IsBackStage -> Assignment1 | Else -> Assignment3)
//	Assignment3 "Choose Surveyor"  (workbasket, Custom;         FlowAction InputSurveyor)      -> Decision8
//	Decision8 IsBack (When IsBackStage -> Assignment7 | Else -> END52 Resolved-Completed)

import (
	"errors"
	"strings"
	"time"
)

// KelasKasus - `pxObjClass` kasus Claim Fac In; kunci CLASS penghitung nomor (`GetSequenceNumber_SQL`,
// `ParamSeq.CARI1 = pyWorkPage.pxObjClass`).
const KelasKasus = "ASM-FW-GCNMFW-Work-PNC"

// KelasKasusKunci - awalan pzInsKey Pega kasus lama (`ASM-FW-GCNMFW-WORK CLM-n`, DEV OS_AKSEPTASI_KLAIM 09-10-2026).
const KelasKasusKunci = "ASM-FW-GCNMFW-WORK"

// Awalan pengenal T_WORK_CLAIM (keputusan work owner 09-10-2026 OQ-CFI-02: awalan sama dengan Pega, nomor dari
// SEQ_WORK_CLAIM bersama semua lini). ⛔ Claim Life juga memakai `CLM-` - lini TIDAK PERNAH disaring lewat awalan, selalu
// lewat `T_WORK_CLAIM.LINI` (prompt §6 butir 1).
const (
	AwalanKlaim  = "CLM-"
	AwalanKomite = "KMT-"
	// LiniFacIn - T_WORK_CLAIM.LINI kedua baris (klaim dan komite) = STS_KLAIM roster EMAILKOMITE (`SetListKomite_act`
	// langkah 2), bawaan a prompt §3.
	LiniFacIn = "FACIN"
)

// Tahap kasus - T_WORK_CLAIM.TAHAP = NAMA FLOWACTION assignment (prompt §6 butir 1). ⛔ Bukan label assignment: label
// "Input Register" sama persis dengan nilai TAHAP Claim Life, dan inbox Claim Life tidak menyaring LINI.
const (
	TahapRegister = "InputRegister"
	TahapEstimasi = "InputEstimasi"
	TahapSurveyor = "InputSurveyor"
	// TahapKomite - kasus komite `KMT-` (flow `Komite_Flow`, korpus Komite Claim FacIn).
	TahapKomite = "Komite_Flow"
)

// LabelTahap - pyMOName shape Register_Flow, tampil di daftar kerja.
var LabelTahap = map[string]string{
	TahapRegister: "Input Register",
	TahapEstimasi: "Input Estimasi",
	TahapSurveyor: "Choose Surveyor",
}

// WorkbasketSurveyor - pemegang Assignment3 Choose Surveyor. Bawaan b prompt §3: `ReasKlaimTeknik` (ada di M_WORKBASKET
// DEV); `ReasPNCTeknik` XML tidak ada di DEV - keputusan yang sama dengan Claim Prop.
const WorkbasketSurveyor = "ReasKlaimTeknik"

// StatusSelesai - pyWorkStatus penutupan (END52, `CloseClaim` ASMForceCaseClose).
const StatusSelesai = "Resolved-Completed"

// Sumber baris T_GENERAL_CLAIM.SUMBER.
const SumberGo = "GO"

// Kasus adalah satu baris T_WORK_CLAIM klaim FACIN.
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
// JSON_KLAIM.IDPEGA, PROGRESSCLAIM.IDPEGA, DOCUMENT_CLAIM.IDPEGA): ID T_WORK_CLAIM APA ADANYA (`CLM-000001`, prompt §6
// butir 6).
func KunciInstans(id string) string { return id }

// KunciPegaLama - pzInsKey Pega kasus lama berpengenal `id` ("ASM-FW-GCNMFW-WORK CLM-n"). Pembaca tabel proyeksi mencari
// kedua bentuk kunci (prompt §6 butir 6).
func KunciPegaLama(id string) string { return KelasKasusKunci + " " + id }

// Galat bersama models.
var (
	// ErrBarisTidakAda - indeks baris di luar daftar.
	ErrBarisTidakAda = errors.New("models: baris tidak ada")
)
