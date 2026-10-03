package models

// Untuk apa berkas ini: BENTUK DATA KASUS yang dipertukarkan repository,
// services, dan handlers - keadaan kerja satu kasus (`T_WORK_POLIS` +
// `T_GENERAL_POLIS`), baris daftar portal, dan pengenal kasus.
//
// Pengenal kasus `NB-<n>`: `[terverifikasi]` filter RD `GetListOpportunity`
// (`.TextNoQuotation Contains "NB-"` bersama `A.Quotation.BusinessFac = T`) -
// kasus realisasi treaty di portal berawalan `NB-`. Angkanya dari sequence
// bersama `SEQ_WORK_POLIS` (premiumlistlife 057), tanpa nol di depan.

import (
	"fmt"
	"strings"
)

// AwalanKasus adalah `pyWorkIDPrefix` kasus realisasi treaty.
const AwalanKasus = "NB-"

// RakitIDKasus menyusun pengenal kasus dari angka sequence.
func RakitIDKasus(urut string) string { return AwalanKasus + strings.TrimSpace(urut) }

// LiniKasus - `T_WORK_POLIS.LINI` kasus ini. `inti` baru mengenal `LIFE`;
// `NONLIFE` sejalan `KODE_PRODUKSI.TYPE = 'NONLIFE'` yang dibaca
// `GeneratePolicyNoTreaty_Act` (RDB `GetKodeProdNonLife_SQL`). Permintaan
// konstanta bersama dicatat di docs/PERMINTAAN-TIM-INTI.md.
const LiniKasus = "NONLIFE"

// BisnisTreaty - nilai `Quotation.BusinessFac` kasus treaty (filter E
// `GetListOpportunity`: `A.Quotation.BusinessFac = T`). Rule pengisinya ada di
// kelas CRM di luar korpus; sistem baru mengisinya saat kasus lahir.
const BisnisTreaty = "T"

// Kasus adalah keadaan kerja satu kasus.
type Kasus struct {
	ID string `json:"id"`
	// Position - `pyWorkPage.Position` ("4" admin, "5" atasan).
	Position string `json:"position"`
	// StatusWork - nama assignment selama berjalan, `Resolved-*` saat tertutup.
	StatusWork string `json:"statusWork"`
	// PositionNote - workbasket tempat kasus menunggu (POSISI, AC 92).
	PositionNote string `json:"positionNote"`
	NoPolis      string `json:"noPolis"`
	// GenerasiTertutup - generasi sudah punya penerus (OLD_POLIS_ID baris lain
	// menunjuknya; ID-10).
	GenerasiTertutup bool   `json:"generasiTertutup"`
	CreateOp         string `json:"createOp"`
	TglCreate        string `json:"tglCreate"`
}

// Tertutup - kasus sudah diselesaikan (disetujui atau ditolak).
func (k Kasus) Tertutup() bool {
	return k.StatusWork == StatusDitolak || k.StatusWork == StatusSelesai
}

// RingkasanKasus adalah satu baris daftar portal
// (`Section/SFAPortal_OpportunitiesList`: nomor, bisnis, tertanggung,
// marketing, NBStatus).
type RingkasanKasus struct {
	ID            string `json:"id"`
	BusinessName  string `json:"businessName"`
	InsuredName   string `json:"insuredName"`
	MarketingName string `json:"marketingName"`
	NBStatus      string `json:"nbStatus"`
	StatusWork    string `json:"statusWork"`
	PositionNote  string `json:"positionNote"`
	NoPolis       string `json:"noPolis"`
	TglCreate     string `json:"tglCreate"`
}

// SaringanKasus - saringan daftar portal.
type SaringanKasus struct {
	// Cari - teks pencarian (`.Name Contains Param.Search`,
	// `.TextNoQuotation Contains Param.Search` di `GetListOpportunity`).
	Cari string
	// Posisi - workbasket; kosong = semua posisi.
	Posisi string
	// Antrean - bila tidak kosong, hanya kasus yang menunggu di salah satu
	// workbasket ini. Diisi gerbang portal `services.DaftarKasus` bagi pelaku
	// di luar wadah grid `ReasTreatyInAdmin` (P8); nil = tanpa batas antrean.
	Antrean []string
}

// JalurAnak - kunci daftar bersarang di halaman: `<induk>(<n>).<anak>`,
// n berbasis 1 seperti `pxListSubscript` Pega.
func JalurAnak(induk string, n int, anak string) string {
	return fmt.Sprintf("%s(%d).%s", induk, n, anak)
}

// KunciInstans = `pyWorkPage.pzInsKey`: kelas work huruf besar + spasi + pyID
// (bentuk kunci instans Pega). Dipakai `HISTORYAKSEPTASIPEGA.ID_PEGA`
// (`InsertHistory.CARI1 = pyWorkPage.pzInsKey`).
func KunciInstans(id string) string { return strings.ToUpper(KelasDeret) + " " + id }

// Pilihan adalah satu opsi daftar pilihan layar (RD mata uang, MO, jenis
// reasuransi).
type Pilihan struct {
	Nilai string `json:"nilai"`
	Label string `json:"label"`
}

// Riwayat adalah satu baris HISTORYAKSEPTASIPEGA
// (`Activity/InsertHistoryAkseptasiPega`).
type Riwayat struct {
	// IDPega - `pyWorkPage.pzInsKey` (`KunciInstans`).
	IDPega string `json:"idPega"`
	// Status - ACCEPT / REJECT / "" (`StatusRiwayat`).
	Status string `json:"status"`
	// Username - `OperatorID.pyUserName` = NAMA TAMPILAN (AC 40).
	Username string `json:"username"`
	// Workbasket - `pyWorkPage.PositionNote` saat submit.
	Workbasket string `json:"workbasket"`
	// OperatorID - IDENTITAS AKSES LOGIN, selalu terisi (P4, AC 39, 41).
	OperatorID  string `json:"operatorId"`
	TglTransfer string `json:"tglTransfer"`
}
