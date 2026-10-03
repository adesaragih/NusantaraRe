// Package models memuat bentuk data modul Marketing Officer.
//
// Satu baris `POOLDATA.MARKETINGOFFICER` (14 kolom warisan Pega, tanpa PK) dan pilihan form: akun login
// (`M_LOGIN_GO`, milik inti, dibaca saja) dan cabang (`BRANCH`, dibaca saja).
package models

// NilaiLeader - `CLIENTID2` baris yang adalah leader (`SaveMarketingOfficer_Act` "Jika leader").
const NilaiLeader = "LEADER"

// Nilai `MOSTATUS` (caption form Pega "Active"; DEV 03-10-2026: 29 baris `1`, 45 baris `2`).
const (
	StatusAktif    = "1"
	StatusNonaktif = "2"
)

// MarketingOfficer adalah satu baris `MARKETINGOFFICER`.
type MarketingOfficer struct {
	ID string `json:"id"`
	// ClientID - Marketing Code; disalin ke data transaksi (quotation, GENERALOFFER, JSON_FOLLOWING, ...), jadi
	// TIDAK PERNAH diubah sesudah baris dibuat.
	ClientID   string `json:"clientId"`
	ClientName string `json:"clientName"`
	// ClientID2 - `LEADER` bila baris ini leader, selain itu ID baris leadernya.
	ClientID2        string `json:"clientId2"`
	MOLeader         string `json:"moLeader"`
	MOStatus         string `json:"moStatus"`
	BranchParent     string `json:"branchParent"`
	BranchDetailID   string `json:"branchDetailId"`
	BranchDetailName string `json:"branchDetailName"`
	TeamGroup        string `json:"teamGroup"`
	// BranchStatus - tidak ada di form Pega; dibaca saja, tidak pernah ditulis modul ini.
	BranchStatus string `json:"branchStatus"`
	// AksesLogin - `M_LOGIN_GO.LOGIN_ID` akun MO (di Pega: Operator ID, dipakai `SendEmailPolicy`).
	AksesLogin string `json:"aksesLogin"`
	UserUpdate string `json:"userUpdate"`
	// Tanggal - `TANGGAL` `YYYY-MM-DD HH:MI` jam Oracle; kosong = tidak pernah diisi.
	Tanggal string `json:"tanggal"`
}

// Leader menjawab apakah baris ini leader.
func (m MarketingOfficer) Leader() bool { return m.ClientID2 == NilaiLeader }

// Aktif menjawab apakah baris ini aktif.
func (m MarketingOfficer) Aktif() bool { return m.MOStatus == StatusAktif }

// Akun adalah satu akun `M_LOGIN_GO`.
type Akun struct {
	LoginID   string `json:"loginId"`
	Nama      string `json:"nama"`
	ContactID string `json:"contactId"`
	Email     string `json:"email"`
	Aktif     bool   `json:"aktif"`
}

// Cabang adalah satu baris `BRANCH` (kelas Pega `ASM-FW-GISFW-Int-BRANCHDETAIL`).
type Cabang struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
	// Induk - `BRANCHPARENTID`: nilai `BRANCHPARENT` di baris MO.
	Induk       string `json:"induk"`
	KanwilGroup string `json:"kanwilGroup"`
	Aktif       bool   `json:"-"`
}

// Isian adalah isian form tambah dan ubah.
type Isian struct {
	// AksesLogin - LOGIN_ID akun yang dipilih; kosong = tanpa akun.
	AksesLogin string `json:"aksesLogin"`
	// Leader - centang "Set as a leader".
	Leader bool `json:"leader"`
	// LeaderID - ID baris leader (dropdown "Leader"); diabaikan bila Leader.
	LeaderID       string `json:"leaderId"`
	BranchParent   string `json:"branchParent"`
	BranchDetailID string `json:"branchDetailId"`
	// Aktif - radio "Active".
	Aktif bool `json:"aktif"`
}

// AksiGo - `MARKETINGOFFICER_LOG.ACTION` baris log perubahan lewat aplikasi Go: AKSES_LOGIN-nya sah (migrasi 760).
// Kosong = baris dari trigger saja (Pega atau lama) - AKSES_LOGIN tidak diketahui.
const AksiGo = "UPDATE-GO"

// BarisLog adalah satu baris `MARKETINGOFFICER_LOG`: keadaan LAMA satu MO sebelum satu UPDATE (trigger warisan).
type BarisLog struct {
	MarketingOfficer
	Aksi string `json:"aksi"`
	// LogTime - `LOG_TIME` `YYYY-MM-DD HH:MI:SS`; kosong = baris lama sebelum migrasi 760 (urutan perkiraan).
	LogTime string `json:"logTime"`
}

// AksesDiketahui menjawab apakah AKSES_LOGIN baris log ini sah.
func (b BarisLog) AksesDiketahui() bool { return b.Aksi == AksiGo }
