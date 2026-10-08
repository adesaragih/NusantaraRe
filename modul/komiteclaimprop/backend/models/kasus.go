// Package models - domain Komite Claim Prop (kasus `ASM-FW-GCNMFW-Work-KomiteTreaty`, `Flow/KomiteTreaty_Flow.xml`).
//
// Untuk apa berkas ini: bentuk kasus komite dan tangganya, serta nilai tetap yang dibaca XML. Nama properti keputusan
// anggota memakai kosakata kami (Keputusan / Komentar / Tanggal), bukan ejaan Pega - penjaga batas Claim Life
// `komite_statik_test.go` hanya mengizinkan penulis tangga (`repository/tangga.go`) menyebut nama tabelnya.
package models

import (
	"time"

	"nusantarare/inti/backend/penomor"
	"nusantarare/inti/backend/utils"
)

// Nilai tetap kasus komite.
const (
	// AwalanKomite - awalan ID kasus komite Prop (STRUKTUR §T_WORK_CLAIM, keputusan work owner 18-09-2026).
	AwalanKomite = "TKMT-"
	// AwalanKlaim - awalan kasus klaim induknya (COVER_KEY).
	AwalanKlaim = "CLMP-"
	// LiniProp - `T_WORK_CLAIM.LINI` kasus komite Prop (keputusan work owner 07-10-2026: saringan KETAT, tanpa
	// `OR LINI IS NULL`).
	LiniProp = "PROP"
	// TahapKomite - `T_WORK_CLAIM.TAHAP` kasus komite (nama flow, Claim Prop `models.TahapKomiteTreaty`).
	TahapKomite = "KomiteTreaty_Flow"
	// StatusSelesai - `KomiteTreaty_Flow` End: Resolved-Completed.
	StatusSelesai = "Resolved-Completed"
	// KelasKlaim - `pyWorkCover.pxObjClass` (KomitePostAdjustment S16.6 `ParamSeq.CARI1`).
	KelasKlaim = "ASM-FW-GCNMFW-Work-ClaimTreaty"
	// KelasKomite - kelas kasus komite (kunci instans HISTORYAKSEPTASIPEGA.ID_KOMITE).
	KelasKomite = "ASM-FW-GCNMFW-Work-KomiteTreaty"
	// TransferAdjustment - `.TransferType` jalur penyesuaian (TT 2). TT 3 (reject claim) tanpa penulis di Claim Prop,
	// TT 4 (close) ditunda OQ-CP-06.
	TransferAdjustment = "2"
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
	ID           string    `json:"id"`
	KlaimID      string    `json:"klaimId"`
	AdjustmentID string    `json:"adjustmentId"`
	Loop         int       `json:"komiteLoop"`
	Count        int       `json:"komiteCount"`
	AcceptStatus string    `json:"acceptStatus"`
	UsulTutup    string    `json:"usulTutup"`
	UsulCadang   string    `json:"usulCadang"`
	Tahap        string    `json:"tahap"`
	StatusWork   string    `json:"statusWork"`
	PembuatID    string    `json:"pembuatId"`
	PembuatNama  string    `json:"pembuatNama"`
	TglCreate    time.Time `json:"tglCreate"`
	TglUpdate    time.Time `json:"tglUpdate"`
	Tangga       []Anggota `json:"tangga"`
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

// Giliran = `KomiteRouter` S6 (S1-S5, S7 ter-remark): `AssignTo` = operator baris tangga PERTAMA ber-keputusan 0
// (S6.1, transisi kode 6 = keluar sesudah yang pertama). Nol workbasket; tanpa baris menunggu = tanpa pemegang.
func (k Kasus) Giliran() (Anggota, bool) {
	for _, a := range k.Tangga {
		if a.Keputusan == KeputusanMenunggu {
			return a, true
		}
	}
	return Anggota{}, false
}

// Pemegang menjawab apakah `akun` memegang assignment kasus ini (KomiteRouter S6.1, keputusan 30 ADR-0014).
func (k Kasus) Pemegang(akun string) bool {
	if k.Tertutup() || akun == "" {
		return false
	}
	a, ada := k.Giliran()
	return ada && a.OperatorID == akun
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

// TingkatAkhir - `pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop` (dievaluasi sebelum S40 menaikkan KomiteCount).
func (k Kasus) TingkatAkhir() bool { return k.Count == k.Loop }

// MasihBerjalan = `When/IsKomiteLoop` (`.AcceptStatus = "1" AND .KomiteCount <= .KomiteLoop`): Decision `KomiteLoop`
// sesudah KomitePost - benar = kembali ke assignment KomiteRouter, salah = Resolved-Completed.
func MasihBerjalan(acceptStatus string, count, loop int) bool {
	return acceptStatus == KeputusanSetuju && count <= loop
}
