// Package models memuat bentuk data modul Accounts - satu baris `POOLDATA.T_M_ACCOUNT` dan pilihannya.
package models

import "strings"

// Awalan kunci Pega SFAGIS (katalog DEV 04-10-2026): `T_M_ACCOUNT.ID` = `ASM-SFAGIS-WORK-ACCOUNT ACC-n`,
// `T_M_ACCOUNT.INSUREDID` = `CLIENT.ID` = `ASM-SFAGIS-WORK-ORG ORG-n`.
const (
	AwalanID     = "ASM-SFAGIS-WORK-ACCOUNT "
	AwalanIDView = "ACC-"
	AwalanOrg    = "ASM-SFAGIS-WORK-ORG "
)

// Account adalah satu baris `T_M_ACCOUNT`. Semua teks; kosong = "".
type Account struct {
	// ID - `ASM-SFAGIS-WORK-ACCOUNT ACC-n`; tidak pernah berubah.
	ID string `json:"id"`
	// IDView - `ACC-n` (bagian akhir ID); kosong bila ID tidak berpola itu.
	IDView          string `json:"idView"`
	GroupBusinessID string `json:"groupBusinessId"`
	GroupBusiness   string `json:"groupBusiness"`
	// InsuredID - `CLIENT.ID` organisasi (`ASM-SFAGIS-WORK-ORG ORG-n`).
	InsuredID string `json:"insuredId"`
	// OrgID - `ORG-n` (bagian akhir InsuredID), yang tampil sebagai "Org ID".
	OrgID       string `json:"orgId"`
	InsuredName string `json:"insuredName"`
	Description string `json:"description"`
	// CreateOp - akun pelaku yang membuat baris ("Owner"); kosong untuk baris lama.
	CreateOp string `json:"createOp"`
	// CreateDate - `YYYY-MM-DD HH24:MI` jam Oracle; kosong untuk baris lama.
	CreateDate string `json:"createDate"`
}

// Isian adalah isian form tambah dan ubah. ID, Owner, dan Create Date tidak pernah datang dari klien.
type Isian struct {
	InsuredID       string `json:"insuredId"`
	GroupBusinessID string `json:"groupBusinessId"`
	Description     string `json:"description"`
}

// Organisasi adalah satu pilihan "Insured Name" - `CLIENT` ber-`FLAG` `Org`.
type Organisasi struct {
	ID     string `json:"id"`
	IDView string `json:"idView"`
	Nama   string `json:"nama"`
}

// GroupBusiness adalah satu pilihan "Group Business" - `POOLDATA.BUSINESSGROUP` (`ID`, `NOTE`).
type GroupBusiness struct {
	ID   string `json:"id"`
	Note string `json:"note"`
}

// Saringan adalah saringan daftar: kata cari dan halaman (mulai 1).
type Saringan struct {
	Cari    string
	Halaman int
	Ukuran  int
}

// BagianAkhir - `ACC-n` / `ORG-n` dari kunci penuh; kosong bila awalan tidak cocok.
func BagianAkhir(id, awalan string) string {
	if !strings.HasPrefix(id, awalan) {
		return ""
	}
	return strings.TrimPrefix(id, awalan)
}
