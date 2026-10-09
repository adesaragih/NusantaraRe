// Package models - domain Komite Claim Non Prop (kasus `ASM-FW-GCNMFW-Work-KomiteTreatyNonProp`,
// `Flow/KomiteTreaty_Flow.xml` korpus `Komite Claim Non Prop`). Pola disalin dari Komite Claim Prop (bukan impor).
//
// Untuk apa berkas ini: bentuk kasus komite dan tangganya, serta nilai tetap yang dibaca XML. Nama properti keputusan
// anggota memakai kosakata kami (Keputusan / Komentar / Tanggal), bukan ejaan Pega - penjaga batas Claim Life
// `komite_statik_test.go` hanya mengizinkan penulis tangga (`repository/tangga.go`) menyebut nama tabelnya.
package models

import (
	"slices"
	"time"

	"nusantarare/inti/backend/penomor"
	"nusantarare/inti/backend/utils"
)

// Nilai tetap kasus komite.
const (
	// AwalanKomite - awalan ID kasus komite Non Prop (`KMTNP-`, awalan Pega; dilahirkan Claim Non Prop).
	AwalanKomite = "KMTNP-"
	// LiniNonProp - `T_WORK_CLAIM.LINI` kasus komite Non Prop (saringan KETAT, pola Komite Claim Prop).
	LiniNonProp = "NONPROP"
	// TahapKomite - `T_WORK_CLAIM.TAHAP` kasus komite (nama flow, Claim Non Prop `models.TahapKomite`).
	TahapKomite = "KomiteTreaty_Flow"
	// StatusSelesai - KomitePostAdjustment S28 `ASMForceCaseClose` WorkStatus Resolved-Completed.
	StatusSelesai = "Resolved-Completed"
	// KelasKlaim - `TempMainWork.pxObjClass` (KomitePostAdjustment S14.9 `ParamSeq.CARI1`).
	KelasKlaim = "ASM-FW-GCNMFW-Work-ClaimTreatyNonProp"
)

// Keputusan anggota tangga (`KomiteList(n)` / kolom `KOMITE_APPROVAL`) dan `.AcceptStatus`.
const (
	KeputusanMenunggu = "0"
	KeputusanSetuju   = "1"
	KeputusanTolak    = "2"
)

// Penanda usul di header kasus komite (`KOMITE_USUL_TUTUP` / `KOMITE_USUL_CADANG`, keputusan 19-09-2026).
const (
	UsulYa    = "1"
	UsulTidak = "0"
)

// Kasus - satu kasus komite beserta tangganya.
type Kasus struct {
	ID           string `json:"id"`
	KlaimID      string `json:"klaimId"`
	AdjustmentID string `json:"adjustmentId"`
	Loop         int    `json:"komiteLoop"`
	Count        int    `json:"komiteCount"`
	AcceptStatus string `json:"acceptStatus"`
	UsulTutup    string `json:"usulTutup"`
	UsulCadang   string `json:"usulCadang"`
	// Subjectivity / SubjectivityNote - isian tingkat 1 yang disimpan antar tingkat (`KOMITE_SUBJECTIVITY` '1'/'0',
	// `KOMITE_SUBJECTIVITY_NOTE`; kolom migrasi komiteclaimprop 682, pola Komite Claim Prop).
	Subjectivity     string `json:"subjectivity"`
	SubjectivityNote string `json:"subjectivityNote"`
	// KomentarAwal - `CreateChildKomiteCNP_Act` (`ChildWorkPage.Comment`): `.Comment` awal =
	// `ComiteeClaim(1).KomiteComment`, komentar anggota pertama putaran komite pertama baris akseptasi yang sama (kosong
	// di putaran pertama).
	KomentarAwal string    `json:"-"`
	Tahap        string    `json:"tahap"`
	StatusWork   string    `json:"statusWork"`
	PembuatID    string    `json:"pembuatId"`
	PembuatNama  string    `json:"pembuatNama"`
	TglCreate    time.Time `json:"tglCreate"`
	TglUpdate    time.Time `json:"tglUpdate"`
	Tangga       []Anggota `json:"tangga"`
}

// Kepala - tulisan kepala kasus komite sesudah satu Submit (S25 / S40, isian tingkat 1).
type Kepala struct {
	Count                                                               int
	AcceptStatus, UsulTutup, UsulCadang, Subjectivity, SubjectivityNote string
}

// Tertutup - kasus sudah Resolved-Completed.
func (k Kasus) Tertutup() bool { return k.StatusWork == StatusSelesai }

// Anggota - satu baris tangga (penyetuju).
type Anggota struct {
	ID         string `json:"id"`
	Urut       int    `json:"urut"`
	OperatorID string `json:"operatorId"`
	// Jabatan - `IDKomite` Pega; label "Committe Name" grid dan `IsCedingConfirm` riwayat.
	Jabatan   string `json:"jabatan"`
	Email     string `json:"-"`
	Keputusan string `json:"keputusan"`
	Komentar  string `json:"komentar"`
	// Tanggal - tanggal diputuskan ("2006-01-02 15:04:05", zona Jakarta); kosong = belum.
	Tanggal string `json:"tanggal"`
}

// Jakarta - zona waktu bisnis (`Asia/Jakarta`, XML `@CurrentDate(...,"Asia/Jakarta")`), diambil dari SATU sumber zona
// repo (`penomor.DiJakarta`).
var Jakarta = penomor.DiJakarta(time.Time{}).Location()

// FormatWaktu - teks waktu halaman (`utils.TanggalWaktu`, Jakarta).
func FormatWaktu(t time.Time) string { return t.In(Jakarta).Format(utils.TanggalWaktu) }

// Giliran = `KomiteRouter` S6 (S1-S5, S7 ter-remark): `AssignTo` = KomiteID baris tangga PERTAMA ber-keputusan 0
// (S6.1, transisi kode 6 = keluar sesudah yang pertama). Tanpa baris menunggu = tanpa pemegang.
func (k Kasus) Giliran() (Anggota, bool) {
	for _, a := range k.Tangga {
		if a.Keputusan == KeputusanMenunggu {
			return a, true
		}
	}
	return Anggota{}, false
}

// Pemegang menjawab apakah `akun` (pemegang workbasket aktif `peran`) memegang assignment kasus ini (KomiteRouter
// S6.1). KomiteID tingkat berjalan = akun itu, ATAU workbasket yang dipegangnya (roster komite NONPROP ke workbasket,
// migrasi claimnonprop 611, keputusan work owner 09-10-2026). Tanpa larangan rangkap: "1 akun memang
// tidak boleh memiliki 2 jabatan dalam komite" dijaga pengaturan akun (Kelola User), bukan di sini (WO 09-10-2026).
func (k Kasus) Pemegang(akun string, peran []string) bool {
	if k.Tertutup() || akun == "" {
		return false
	}
	a, ada := k.Giliran()
	return ada && (a.OperatorID == akun || slices.Contains(peran, a.OperatorID))
}

// barisBerjalan - anggota tangga tingkat `KomiteCount` (`KomiteList(local.count)`, KomitePostAdjustment S5-S6).
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
