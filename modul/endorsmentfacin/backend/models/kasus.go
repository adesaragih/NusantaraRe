package models

// Kasus Endorsment Fac In baru - irisan I3 (kontrak sesi utama 11-10-2026; keputusan work owner W-EDM-08…12).

import (
	"fmt"
	"time"
)

// Keadaan awal baris T_WORK_POLIS kasus EDM (keputusan work owner 11-10-2026).
const (
	// LiniEDMFacIn - `T_WORK_POLIS.LINI` kasus Endorsment Fac In. BUKAN "FAC": NB FacIn menyaring lininya hanya lewat
	// kolom ini (portal, kotak masuk, buka kasus), jadi kasus EDM ber-"FAC" akan tampil di layar NB.
	LiniEDMFacIn = "FACEDM"
	// PosisiAwalEDM - T_WORK_POLIS.POSITION = TAHAP (keputusan work owner W-EDM-18, mengganti W-EDM-11; pola NB tiket 77):
	// kasus EDM lahir di fase Policy (konektor Start2 `.IsCedingConfirm = Policy`).
	PosisiAwalEDM = "Policy"
	// CatatanPosisiAwalEDM - T_WORK_POLIS.POSITION_NOTE = WORKBASKET (migrasi nbfacin 254): `Assignment7` "MARKETING"
	// (flow InputAddendumFacultativeIn, konektor Start2 `PositionNote`). Sebelum 254 berjalan, nilai ini ditulis ke
	// POSITION (bentuk W-EDM-11) dan dikonversi migrasi 263.
	CatatanPosisiAwalEDM = "ReasFacInMarketing"
	// StatusAwalEDM - STATUS_WORK kasus EDM yang baru lahir.
	StatusAwalEDM = "Pending-Policy"
	// FlagAwalEDM - `.FlagOnGoingPolicy = "1"` konektor Start2.
	FlagAwalEDM = "1"
	// PolaStatusSelesai - STATUS_WORK kasus yang sudah selesai (pola LIKE): gerbang 2 mencari yang BELUM.
	PolaStatusSelesai = "Resolved-%"
)

// StatusEDMSelesaiProduksi - STATUS_WORK kasus EDM sistem baru yang SUDAH BERPRODUKSI; sumber gerbang 1 polis sistem
// baru (A51; padanan Pega: EdmType versi polis terakhir di JSON_POLIS, yang lahir saat EDM berproduksi).
// `[pertanyaan terbuka]` ditetapkan irisan produksi EDM. KOSONG = belum ada status itu: kueri tidak pernah cocok
// (Oracle membaca ” sebagai NULL), jadi belum ada kasus sistem baru yang terhitung berproduksi.
const StatusEDMSelesaiProduksi = ""

// AwalanKasusEDM - `SetValueToEDMWork` langkah 7 `param.IDPrefix = "EDM-"`.
const AwalanKasusEDM = "EDM-"

// RakitIDKasusEDM - "EDM-" + nomor sequence berpadding 5 digit (keputusan work owner 11-10-2026: "EDM-00001").
// ⚠️ Kasus Pega berbentuk `EDM-<n>` tanpa padding dan sudah melewati 10000; mulai EDM-10000 bentuknya sama (A42).
func RakitIDKasusEDM(nomor int64) string { return fmt.Sprintf("%s%05d", AwalanKasusEDM, nomor) }

// MasukanKasus - badan `POST /api/endorsmentfacin/kasus`: isian /periksa + medan form lainnya.
//
// ⛔ Tanpa AdminCharge / Survey (A50): `.AdminCharge` (L9361) dan `.Survey` (L9533) Section WorkPrimaryDetails.xml
// berada di kontainer `pyContainerVisibleWhen = never` (L9175) - Pega tidak pernah menampilkannya. Kolom
// ADMIN_CHARGE / SURVEY T_EDM_FACIN_KASUS (262, sudah jalan di DEV) tetap ada dan ditulis NULL.
type MasukanKasus struct {
	MasukanPeriksa
	Note                  string `json:"note"`
	EndorsementSource     string `json:"endorsementSource"`
	EndorsementSourceNote string `json:"endorsementSourceNote"`
	EndorsementStatus     string `json:"endorsementStatus"`
	DetailType            string `json:"detailType"`
}

// HasilKasus - jawaban /kasus. CaseID terisi = kasus lahir (201); kosong = ditolak / tidak lengkap (200, Periksa).
type HasilKasus struct {
	CaseID string `json:"caseId"`
	// Pesan - pesan OpenCase yang tetap dipasang walau kasus lahir (langkah 4); selalu larik ([] bila kosong).
	Pesan []PesanPeriksa `json:"pesan"`
	// Periksa - hasil pemeriksaan bila kasus TIDAK lahir.
	Periksa *HasilPeriksa `json:"periksa,omitempty"`
}

// KepalaKasus - satu baris T_EDM_FACIN_KASUS (migrasi 262).
type KepalaKasus struct {
	ID, OldPolicyNo, EdmType, EdmTypeNew, DetailType string
	// EndorsementDate - tanggal kalender WIB (jam diabaikan).
	EndorsementDate                                             time.Time
	EndorsementReason, EndorsementSource, EndorsementSourceNote string
	EndorsementStatus, AdminCharge, Survey, InternalRetro       string
	NBStatus, IsCedingConfirm                                   string
}

// PilihanEndorsementStatus - Prompt List properti `.EndorsementStatus` (nilai standar → teks layar). Sumber: tangkapan
// layar rule Property Pega kiriman work owner 11-10-2026 (properti tidak ada di korpus). Wajib diisi (Section
// WorkPrimaryDetails L8906 `pyRequiredWhen CARI2!='Life'`; Life di luar EDM Fac In, W-EDM-05).
var PilihanEndorsementStatus = map[string]string{
	"0": "Internal",
	"1": "External",
	"2": "Batal Internal",
	"3": "Batal External",
}

// PilihanInternalRetro - Prompt List properti `.EndorsementInternalRetro`; sumber sama (work owner 11-10-2026).
var PilihanInternalRetro = map[string]string{
	"1": "Yes",
	"0": "No",
}

// DetailTypeMenurutEdmTypeNew - pilihan Detail EDM Type per EdmTypeNew, `Endorsment Fac In\Activity\SetEdmType.xml`
// `[terverifikasi]` (CARI1; label CARI2): langkah 2 → `TypePenambahan` (Section EdmType1), 3 → `TypePengurangan`
// (EdmType2), 4 → `TypePerubahan` (EdmType3).
// ⛔ K-046 apa adanya: langkah 3 menambahkan "Adjustment PPN/PPH" (12) ke `TypePenambahan`, BUKAN ke
// `TypePengurangan` - pilihan pengurangan tidak memuat 12.
var DetailTypeMenurutEdmTypeNew = map[string][]JenisPenyesuaian{
	"1": {"1", "2", "4", "6", "7", "8", "11", "12"},
	"2": {"2", "4", "6", "7", "8", "11"},
	"3": {"9", "0", "11"},
}
