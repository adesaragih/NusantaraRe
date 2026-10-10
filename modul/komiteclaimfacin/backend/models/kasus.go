// Package models - domain Komite Claim Fac In (kasus `ASM-FW-GCNMFW-Work-Komite`, `Flow/Komite_Flow.xml` korpus
// `Komite Claim FacIn`: assignment "KomiteRouter" -> flow action `ViewTransferDtl` (pra-proses `SetValueKomite`,
// Section `ShowTransfer`, pasca-proses `KomitePostAct`) -> decision `IsKomiteLoop`). Pola disalin dari Komite Claim
// Prop / Non Prop (bukan impor). Murni: nol SQL, nol HTTP.
//
// Untuk apa berkas ini: bentuk kasus komite dan tangganya, wewenang pemegang tingkat berjalan, dan nilai tetap yang
// dibaca XML. Nama properti keputusan anggota memakai kosakata kami (Keputusan / Komentar / Tanggal), bukan ejaan Pega -
// penjaga batas Claim Life `komite_statik_test.go` hanya mengizinkan penulis tangga (`repository/tangga.go`) menyebut
// nama tabelnya.
package models

import (
	"slices"
	"time"

	"nusantarare/inti/backend/penomor"
	"nusantarare/inti/backend/utils"
)

// Nilai tetap kasus komite.
const (
	// AwalanKomite - awalan ID kasus komite Fac In (`KMT-`, awalan Pega, OQ-CFI-02; dilahirkan Claim Fac In). Lini
	// disaring lewat `T_WORK_CLAIM.LINI`, bukan awalan (fixture Claim Life juga memakai `KMT-`).
	AwalanKomite = "KMT-"
	// LiniFacIn - `T_WORK_CLAIM.LINI` kasus komite Fac In (saringan KETAT, pola Komite Claim Prop).
	LiniFacIn = "FACIN"
	// TahapKomite - `T_WORK_CLAIM.TAHAP` kasus komite (nama flow, Claim Fac In `models.TahapKomite`).
	TahapKomite = "Komite_Flow"
	// StatusSelesai - END52 `Komite_Flow` (decision IsKomiteLoop "NoLoop").
	StatusSelesai = "Resolved-Completed"
	// KelasKlaim - `pyWorkCover.pxObjClass` (KomitePost_Adjustment S7.2.1.2.5 `ParamSeq.CARI1`).
	KelasKlaim = "ASM-FW-GCNMFW-Work-PNC"
)

// Jenis penyerahan `pyWorkPage.TransferType` (kolom `T_GENERAL_KOMITE.TRANSFER_TYPE`, migrasi 642; KomitePostAct S2 / S4
// / S5 - TT1 survey ter-remark S3).
const (
	TransferAdjustment = "2" // CreateKMTNo_Act 6.9
	TransferReject     = "3" // SendRejectClaimToKomite2 7.3
	TransferClose      = "4" // SendCloseClaimToKomite 7.3
)

// JudulTransfer - label Section ShowTransfer LS1 menurut TransferType (VERBATIM).
var JudulTransfer = map[string]string{TransferAdjustment: "ADJUSTMENT", TransferReject: "REJECT", TransferClose: "CLOSE"}

// Keputusan anggota tangga (`KomiteList(n).KomiteAproval` / kolom `KOMITE_APPROVAL`) dan `.AcceptStatus`.
const (
	KeputusanMenunggu = "0"
	KeputusanSetuju   = "1"
	KeputusanTolak    = "2"
)

// Penanda usul di header kasus komite (`KOMITE_USUL_TUTUP` / `KOMITE_USUL_CADANG`, migrasi komiteclaimprop 680).
const (
	UsulYa    = "1"
	UsulTidak = "0"
)

// Workbasket roster komite FACIN (migrasi 640, keputusan work owner 10-10-2026 KCF-01).
const (
	WorkbasketSPVA     = "ReasClaimSPVA"
	WorkbasketSPVB     = "ReasClaimSPVB"
	WorkbasketDeptHead = "ReasClaimDeptHead"
)

// Kasus - satu kasus komite beserta tangganya.
type Kasus struct {
	ID      string `json:"id"`
	KlaimID string `json:"klaimId"`
	// AdjustmentID - baris T_CLAIM_ADJUSTMENT yang diputus (TT2); kosong untuk TT3 / TT4 (migrasi 641).
	AdjustmentID string `json:"adjustmentId"`
	TransferType string `json:"transferType"`
	Loop         int    `json:"komiteLoop"`
	Count        int    `json:"komiteCount"`
	AcceptStatus string `json:"acceptStatus"`
	UsulTutup    string `json:"usulTutup"`
	UsulCadang   string `json:"usulCadang"`
	Tahap        string `json:"tahap"`
	StatusWork   string `json:"statusWork"`
	// PembuatID / PembuatNama - `pyWorkPage.pxCreateOperator` (penerima email "(Approval)" / "(Reject)").
	PembuatID   string    `json:"pembuatId"`
	PembuatNama string    `json:"pembuatNama"`
	TglCreate   time.Time `json:"tglCreate"`
	TglUpdate   time.Time `json:"tglUpdate"`
	Tangga      []Anggota `json:"tangga"`
	// Kronologi / Extent / Liability - teks pop-up TT3 / TT4 (`Komite.CircumtansesCouseOfLoss` / `ExtentOfLoss` /
	// `LegalLiability`, SendRejectClaimToKomite2 / SendCloseClaimToKomite 7.2; migrasi 643, OQ-KCFI-03). Tampil lewat
	// blok teks komite, bukan JSON kasus.
	Kronologi string `json:"-"`
	Extent    string `json:"-"`
	Liability string `json:"-"`
}

// Kepala - tulisan kepala kasus komite sesudah satu Submit (KomitePost_Adjustment S14 / S24, KomitePost_Reject S16;
// Loop = `ApprovalKomite_Act` S8 sesudah perluasan).
type Kepala struct {
	Count, Loop                         int
	AcceptStatus, UsulTutup, UsulCadang string
}

// Tertutup - kasus sudah Resolved-Completed.
func (k Kasus) Tertutup() bool { return k.StatusWork != "" }

// Anggota - satu baris tangga (penyetuju).
type Anggota struct {
	ID         string `json:"id"`
	Urut       int    `json:"urut"`
	OperatorID string `json:"operatorId"`
	// Jabatan - `IDKomite` Pega (JABATAN roster); kolom "Committee" grid List of Committee dan teks kronologi.
	Jabatan   string `json:"jabatan"`
	Email     string `json:"-"`
	Keputusan string `json:"keputusan"`
	Komentar  string `json:"komentar"`
	// Tanggal - tanggal diputuskan ("2006-01-02 15:04:05", zona Jakarta); kosong = belum.
	Tanggal string `json:"tanggal"`
}

// Jakarta - zona waktu bisnis (`Asia/Jakarta`), dari SATU sumber zona repo (`penomor.DiJakarta`).
var Jakarta = penomor.DiJakarta(time.Time{}).Location()

// FormatWaktu - teks waktu halaman (`utils.TanggalWaktu`, Jakarta).
func FormatWaktu(t time.Time) string { return t.In(Jakarta).Format(utils.TanggalWaktu) }

// barisBerjalan - indeks anggota tangga tingkat `KomiteCount` (`KomiteList(pyWorkPage.KomiteCount)`, KomitePost_Adjustment
// S3 / S7.2.1.8); -1 bila tidak ada.
func (k Kasus) barisBerjalan() int {
	for i, a := range k.Tangga {
		if a.Urut == k.Count {
			return i
		}
	}
	if k.Count >= 1 && k.Count <= len(k.Tangga) {
		return k.Count - 1
	}
	return -1
}

// Giliran - anggota tangga tingkat berjalan yang masih menunggu.
//
// `KomiteRouter` S6.1 menugaskan baris menunggu PERTAMA (transisi pasca-langkah `true` -> keluar); PERBAIKAN kelainan
// XML (prompt §5 butir 1, OQ-CFI-03): `SetProteksiSubmiteKomite` L2.1 meloloskan akun yang cocok dengan baris menunggu
// MANA SAJA, padahal `KomitePost_*` menulis keputusan ke `KomiteList(KomiteCount)`. Di sini pemutus = anggota tingkat
// berjalan (`KomiteCount`), selaras penulisnya (pada alur berurutan = baris menunggu pertama).
func (k Kasus) Giliran() (Anggota, bool) {
	i := k.barisBerjalan()
	if i < 0 || k.Tangga[i].Keputusan != KeputusanMenunggu {
		return Anggota{}, false
	}
	return k.Tangga[i], true
}

// Pemegang menjawab apakah `akun` (pemegang workbasket aktif `peran`) boleh memutus tingkat berjalan kasus ini: KomiteID
// tingkat berjalan = akun itu, ATAU workbasket yang dipegangnya (roster ke workbasket, migrasi 640). Anggota
// `ReasClaimSPVB` juga boleh memutus tingkat ber-KomiteID `ReasClaimSPVA` (cadangan SPV A, KCF-01). Tanpa larangan
// rangkap (pola Komite Claim Prop 09-10-2026); pintu belakang "IT Developer" dan pemetaan akun `SetProteksiSubmiteKomite`
// L1 / L2.1 tidak ditiru (prompt §5 butir 1).
func (k Kasus) Pemegang(akun string, peran []string) bool {
	if k.Tertutup() || akun == "" {
		return false
	}
	a, ada := k.Giliran()
	if !ada {
		return false
	}
	if a.OperatorID == akun || slices.Contains(peran, a.OperatorID) {
		return true
	}
	return a.OperatorID == WorkbasketSPVA && slices.Contains(peran, WorkbasketSPVB)
}

// PeranKerja - workbasket yang dicocokkan daftar kerja `peran`: anggota ReasClaimSPVB juga melihat tingkat ber-KomiteID
// ReasClaimSPVA (KCF-01).
func PeranKerja(peran []string) []string {
	out := slices.Clone(peran)
	if slices.Contains(out, WorkbasketSPVB) && !slices.Contains(out, WorkbasketSPVA) {
		out = append(out, WorkbasketSPVA)
	}
	return out
}

// AnggotaSPVB - pemutus anggota workbasket ReasClaimSPVB SAJA (`OperatorID.pyPosition == "SPV B"`, ApprovalKomite_Act
// L5; KCF-01: keanggotaan workbasket, bukan jabatan di profil). Anggota ReasClaimSPVA + ReasClaimSPVB sekaligus = SPV A,
// pita tidak berlaku (jawaban work owner 10-10-2026 OQ-KCFI-06; di produksi rangkap dijaga Kelola User).
func AnggotaSPVB(peran []string) bool {
	return slices.Contains(peran, WorkbasketSPVB) && !slices.Contains(peran, WorkbasketSPVA)
}
