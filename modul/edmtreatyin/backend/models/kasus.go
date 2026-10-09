package models

// Untuk apa berkas ini: BENTUK DATA KASUS yang dipertukarkan repository, services, dan handlers - keadaan kerja
// satu kasus endorsemen (`T_WORK_POLIS` + generasi `T_GENERAL_POLIS_TREATY`), baris daftar portal, dan pengenal
// kasus. Asal: salinan `modul/nbtreatyin/backend/models/kasus.go` (06-10-2026), disesuaikan EDM.
//
// Pengenal kasus `EDMT-<n>`: `[terverifikasi]` `Activity/EDMChooseBusiness_Act` langkah 11-12
// (`pyWorkPage.pyWorkIDPrefix=="EDMT-"`) - kelas `ASM-FW-GISFW-Work-EndorsementTreaty`; rule pendefinisi
// prefiksnya (case type) tidak ada di korpus. Angkanya dari sequence bersama `SEQ_WORK_POLIS` (premiumlistlife
// 057), seperti `NB-<n>` NB Treaty In - pengenal T_WORK_POLIS unik lintas awalan.

import (
	"fmt"
	"strings"
)

// AwalanKasus adalah `pyWorkIDPrefix` kasus endorsemen treaty.
const AwalanKasus = "EDMT-"

// RakitIDKasus menyusun pengenal kasus dari angka sequence.
func RakitIDKasus(urut string) string { return AwalanKasus + strings.TrimSpace(urut) }

// LiniKasus - `T_WORK_POLIS.LINI` kasus ini: lini polis treaty inward yang di-endorse (NB Treaty In menulis
// `NONLIFE`, sejalan `KODE_PRODUKSI.TYPE = 'NONLIFE'`). Daftar NB menyaring `PRODKE = 0`, daftar EDM
// `PRODKE >= 1` - kedua modul tidak saling melihat kasusnya.
const LiniKasus = "NONLIFE"

// Kasus adalah keadaan kerja satu kasus.
type Kasus struct {
	ID string `json:"id"`
	// Position - `pyWorkPage.Position` ("4" admin, "5" atasan).
	Position string `json:"position"`
	// StatusWork - nama assignment selama berjalan, `Resolved-*` saat tertutup.
	StatusWork string `json:"statusWork"`
	// PositionNote - workbasket tempat kasus menunggu (POSISI).
	PositionNote string `json:"positionNote"`
	// NoPolis - NOPOLIS generasi ini: KOSONG selama endorsemen belum selesai, diisi nomor polis induk saat
	// Utility1 (`SaveJsonPolisTreatyInEDM_Act`) - lihat repository/generasi.go.
	NoPolis string `json:"noPolis"`
	// ProdKe - PRODKE generasi ini (>= 1).
	ProdKe int `json:"prodKe"`
	// OldPolisID - OLD_POLIS_ID: ID generasi tepat sebelumnya (kosong sesudah admin menolak - generasi dilepas
	// dari rantai polis, repository.LepasGenerasi).
	OldPolisID string `json:"oldPolisId"`
	// GenerasiTertutup - generasi sudah punya penerus (OLD_POLIS_ID baris lain menunjuknya; ID-10).
	GenerasiTertutup bool   `json:"generasiTertutup"`
	CreateOp         string `json:"createOp"`
	TglCreate        string `json:"tglCreate"`
}

// Tertutup - kasus sudah diselesaikan (disetujui atau ditolak).
func (k Kasus) Tertutup() bool {
	return k.StatusWork == StatusDitolak || k.StatusWork == StatusSelesai
}

// RingkasanKasus adalah satu baris daftar portal - kolom grid `Section/SFAPortal_Endorsement_Treaty`
// (RD `InboxEDM_RD2`), urut VERBATIM: EDM Number (`A.pyID`), Offer No (`PolicyTreatyIn.NoOffer`), Policy
// Number (`PolicyTreatyIn.PolicyNo`), EDM Number (`PolicyTreatyIn.EDMNo`), SOB (`SOBName`), Ceding
// (`CedingCoName`), EDM Type (`EDMType`), Proportional Type (`Quotation.ProportionalType`), Marketing
// (`Quotation.MarketingName`), Status (`NBStatus`). Ditambah PositionNote (syarat tautan buka) dan kolom
// kotak masuk Beranda (StartDate, TglUpdate).
type RingkasanKasus struct {
	ID               string `json:"id"`
	NoOffer          string `json:"noOffer"`
	NoPolis          string `json:"noPolis"`
	EDMNo            string `json:"edmNo"`
	SOBName          string `json:"sobName"`
	CedingCoName     string `json:"cedingCoName"`
	EDMType          string `json:"edmType"`
	ProportionalType string `json:"proportionalType"`
	MarketingName    string `json:"marketingName"`
	NBStatus         string `json:"nbStatus"`
	StatusWork       string `json:"statusWork"`
	PositionNote     string `json:"positionNote"`
	TglCreate        string `json:"tglCreate"`
	StartDate        string `json:"startDate"`
	TglUpdate        string `json:"tglUpdate"`
	// TglProd - T_GENERAL_POLIS_TREATY.TGL_PROD (`PolicyTreatyIn.ProductionDate`, Utility1): kolom "Production Date"
	// tab Resolved portal (keputusan work owner 07-10-2026, bukan kolom RD InboxEDM_RD2); kosong bila belum selesai.
	TglProd string `json:"tglProd"`
}

// AntreanKotakMasuk - satu baris kotak masuk Beranda: satu workbasket tangga yang DIPEGANG akun dan jumlah
// berkas yang MENUNGGU dia di sana (kontrak `MenuModul.antreanBeranda`, pola NB Treaty In).
type AntreanKotakMasuk struct {
	Workbasket string `json:"workbasket"`
	// Nama - M_WORKBASKET.NAME; tanpa master = ID workbasket.
	Nama   string `json:"nama"`
	Jumlah int    `json:"jumlah"`
}

// BatasDaftarPortal - baris terbanyak daftar portal = `pyMaxRecords` RD `InboxEDM_RD2` (500).
const BatasDaftarPortal = 500

// SaringanKasus - saringan daftar portal RD `InboxEDM_RD2` (logika `B AND C AND A AND D AND E`).
type SaringanKasus struct {
	// Cari - kotak saring portal. XML filter B hanya `.OfferFacIn.QuotationData.OldPolicyNo Contains
	// Param.FilterTermForEndorsement`; ⛔ perintah work owner 07-10-2026 ("pencarian ... buat bisa mencari nomor
	// nb/edm, insured name dll, intinya buat searchnya itu sangat berguna"): setiap kata (`KataCari`) cocok dengan
	// salah satu kolom portal - nomor kasus, Offer No, nomor polis, EDM No, insured, group business, SOB, ceding,
	// marketing, treaty group, class of business, nama pembuat (`repository.kolomCariPortal`; seragam dengan NB).
	Cari string
	// Pembuat - filter A `.pxCreateOperator = Param.UserNameID`. Kosong = tanpa saringan pembuat.
	Pembuat string
	// PembuatPosisi - bila diisi, saringan pembuat hanya berlaku bagi berkas di posisi ini (kotak masuk Beranda).
	PembuatPosisi string
	// Posisi - workbasket; kosong = semua posisi (kotak masuk Beranda).
	Posisi string
	// Antrean - bila tidak kosong, berkas yang menunggu di salah satu workbasket ini (kotak masuk Beranda).
	Antrean []string
	// Selesai - switch portal In Progress / Resolved (aturan portal NB Treaty In, keputusan work owner 07-10-2026
	// "YA"): false = filter C `.pyStatusWork != "Resolved-Completed"` (berkas yang masih proses), true = HANYA berkas
	// selesai (Completed / Rejected).
	Selesai bool
}

// JalurAnak - kunci daftar bersarang di halaman: `<induk>(<n>).<anak>`, n berbasis 1 seperti
// `pxListSubscript` Pega.
func JalurAnak(induk string, n int, anak string) string {
	return fmt.Sprintf("%s(%d).%s", induk, n, anak)
}

// KunciInstans - IDPEGA berkas sistem baru (padanan `pyWorkPage.pzInsKey`) di HISTORYAKSEPTASIPEGA.ID_PEGA,
// HISTORYAKSEPTASIPRODUCTION, json_polis, ACHIEVEMENT, TREATYINPRODUCTION, T_POLIS_DIFFERENCE.IDPEGA dan muatan
// konversi: ID T_WORK_POLIS APA ADANYA (`EDMT-990001`) - ketetapan NB Treaty In `[keputusan work owner
// 06-10-2026]` "INTINYA KEY DARI T_WORK_POLIS JANGAN DI UBAH".
func KunciInstans(id string) string { return id }

// Pilihan adalah satu opsi daftar pilihan layar (RD mata uang, MO, jenis reasuransi).
type Pilihan struct {
	Nilai string `json:"nilai"`
	Label string `json:"label"`
}

// Riwayat adalah satu baris HISTORYAKSEPTASIPEGA (`Activity/InsertHistoryAkseptasiPega`).
type Riwayat struct {
	// IDPega - `pyWorkPage.pzInsKey` (`KunciInstans`).
	IDPega string `json:"idPega"`
	// Status - ACCEPT / REJECT / "" (`StatusRiwayat`).
	Status string `json:"status"`
	// Username - `OperatorID.pyUserName` = NAMA TAMPILAN.
	Username string `json:"username"`
	// Workbasket - `pyWorkPage.PositionNote` saat submit.
	Workbasket string `json:"workbasket"`
	// OperatorID - IDENTITAS AKSES LOGIN, selalu terisi (P4).
	OperatorID  string `json:"operatorId"`
	TglTransfer string `json:"tglTransfer"`
}
