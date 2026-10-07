package models

// Untuk apa berkas ini: BENTUK DATA KASUS yang dipertukarkan repository,
// services, dan handlers - keadaan kerja satu kasus (`T_WORK_POLIS` +
// `T_GENERAL_POLIS_TREATY`), baris daftar portal, dan pengenal kasus.
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
	// NamaPembuat - T_WORK_POLIS.CREATE_OP_NAME, kolom "User Create" portal (keputusan work owner 06-10-2026).
	NamaPembuat string `json:"createOpName"`
	// JenisProporsi - T_POLIS_QUOTATION.PROPORTIONAL_TYPE, kolom "Type" portal (keputusan work owner 06-10-2026).
	JenisProporsi string `json:"proportionalType"`
	// Kolom daftar kotak masuk Beranda (keputusan work owner 06-10-2026): Ceding Company, Inception Date
	// (T_GENERAL_POLIS_TREATY.CEDING_CO_NAME / START_DATE), Time Since Last Update (T_WORK_POLIS.TGL_UPDATE).
	CedingCoName string `json:"cedingCoName"`
	StartDate    string `json:"startDate"`
	TglUpdate    string `json:"tglUpdate"`
	// TglProduksi - T_GENERAL_POLIS_TREATY.TGL_PROD, kolom "Production Date" tab Resolved portal (perintah work owner 07-10-2026: "KALO DAH RESOLVE TAMBAHIN KOLOM NOPOLISNYA" dan "SEKALIAN KELUARIN TANGGAL PRODUKSINYA AJA DD-MM-YYYY").
	// Nomor polisnya = NoPolis (g.NOPOLIS); keduanya terisi saat realisasi selesai, kosong untuk berkas ditolak.
	TglProduksi string `json:"productionDate"`
}

// AntreanKotakMasuk - satu baris kotak masuk Beranda (keputusan work owner 06-10-2026: "beranda menunjukkan berapa
// banyak case yang masuk di akun dia, mengikuti workbasket"): satu workbasket tangga yang DIPEGANG akun dan jumlah
// berkas yang MENUNGGU dia di sana.
type AntreanKotakMasuk struct {
	Workbasket string `json:"workbasket"`
	// Nama - M_WORKBASKET.NAME; tanpa master = ID workbasket.
	Nama   string `json:"nama"`
	Jumlah int    `json:"jumlah"`
}

// BatasDaftarPortal - baris terbanyak daftar portal = `pyMaxRecords` RD
// `GetListOpportunity` (500).
const BatasDaftarPortal = 500

// SaringanKasus - saringan daftar portal.
type SaringanKasus struct {
	// Cari - teks pencarian `.FilterTermForOpportunity` -> `Param.Search` RD
	// `GetListOpportunity`: filter G `.TextNoQuotation Contains Param.Search`
	// (tanpa beda huruf besar/kecil) = pengenal kasus (`CocokCariPortal`).
	// ⛔ Filter C `.Name Contains Param.Search` tidak dibangun - `.Name`
	// milik kelas CRM, ditulis nol rule korpus, tak berkolom di diagram.
	Cari string
	// Posisi - workbasket; kosong = semua posisi.
	Posisi string
	// Antrean - bila tidak kosong, hanya kasus yang menunggu di salah satu
	// workbasket ini. Diisi gerbang portal `services.DaftarKasus` bagi pelaku
	// di luar wadah grid `ReasTreatyInAdmin` (P8); nil = tanpa batas antrean.
	Antrean []string
	// Pembuat - filter A RD `GetListOpportunity` (`A.pxCreateOperator = Param.UserIdentifier`; RALAT keputusan
	// work owner 06-10-2026): hanya berkas buatan akun ini (CREATE_OP), ATAU berkas di `Antrean` bila diisi.
	// Kosong = tanpa saringan pembuat.
	Pembuat string
	// PembuatPosisi - bila diisi, saringan pembuat hanya berlaku bagi berkas di posisi ini (kotak masuk Beranda:
	// buatan akun yang MASIH di Admin = menunggu dia).
	PembuatPosisi string
	// Selesai - switch portal (keputusan work owner 06-10-2026: "switch untuk lihat yang lagi proses atau
	// resolve, default ke proses"): true = berkas Resolved (Completed / Rejected), false = yang masih proses.
	Selesai bool
}

// CocokCariPortal = filter G RD `GetListOpportunity`: `.TextNoQuotation
// Contains Param.Search`, `pyCaseInsensitive=true`. Kosong = tanpa saringan.
// Repository menulis padanannya di SQL (`UPPER(w.ID) LIKE`).
func CocokCariPortal(id, cari string) bool {
	cari = strings.TrimSpace(cari)
	return cari == "" || strings.Contains(strings.ToUpper(id), strings.ToUpper(cari))
}

// JalurAnak - kunci daftar bersarang di halaman: `<induk>(<n>).<anak>`,
// n berbasis 1 seperti `pxListSubscript` Pega.
func JalurAnak(induk string, n int, anak string) string {
	return fmt.Sprintf("%s(%d).%s", induk, n, anak)
}

// KunciInstans - IDPEGA berkas sistem baru (padanan `pyWorkPage.pzInsKey`) di HISTORYAKSEPTASIPEGA.ID_PEGA,
// HISTORYAKSEPTASIPRODUCTION, json_polis, ACHIEVEMENT, TREATYINPRODUCTION dan muatan konversi: ID T_WORK_POLIS APA
// ADANYA (`NB-22445`). `[keputusan work owner 06-10-2026]` "INTINYA KEY DARI T_WORK_POLIS JANGAN DI UBAH" - RALAT
// bentuk lama `ASM-FW-GISFW-WORK-NB NB-x` (kelas + spasi + pyID); baris DEV berawalan itu diubah ke ID polos.
// Dokumen Pega lama tetap berkunci pzInsKey Pega (`ASM-FW-GISFW-WORK NB-x`, `IDKasusDariIDPega`).
func KunciInstans(id string) string { return id }

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
