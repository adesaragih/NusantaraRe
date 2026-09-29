package models

// Business pada kombinasi (tahun, grup, jenis) - tiket 07 Treaty Contract Out.
//
// Pohon yang ditiru, dibaca 29-09-2026 (nomor baris mentah):
//
//	`Section/ViewDetailTreatyBusinessGrid.xml`  `Active` b6499 (radio, wajib
//	                                           `pyRequired` b6542), `Business Name`
//	                                           b6241 (pemilih `BrowseFilterBusiness_RD`)
//	`Activity/SaveTreatyBusinessDetail_Act.xml` tanpa satu pun gerbang
//
// Dibaca sesudah: tco_reinsurer.go.

import (
	"errors"
	"fmt"
	"strings"
)

// Nilai `ISACTIVE`.
//
// `[terverifikasi]` hilir menyaring `isactive='1'` (`Claim Prop/RDBList/
// GetTreatyGroupID.xml`) dan grid `BrowseTreatyBusiness_RD` b670-b675
// `.IsActive = 1`. Nilai nonaktif `0` [keputusan work owner 29-09-2026] - "0 berarti nonaktif";
// daftar pilihan radio milik properti tidak diekspor (OQ-TCO-13, ditutup).
const (
	BusinessAktif    = "1"
	BusinessNonaktif = "0"
)

var (
	// ErrBusinessKodeKosong - `Business Name` wajib dipilih dari master.
	ErrBusinessKodeKosong = errors.New("models: BizCode wajib - pilih Business Name dari master BUSINESS")
	// ErrBusinessAktifTakSah - `Active` wajib (`pyRequired` b6542) dan hanya 1/0.
	ErrBusinessAktifTakSah = errors.New("models: IsActive wajib bernilai 1 (aktif) atau 0 (nonaktif)")
)

// PeriksaBusinessTCO menjalankan gerbang simpan baris bisnis.
func PeriksaBusinessTCO(b BusinessTreaty) error {
	if strings.TrimSpace(b.BizCode) == "" {
		return ErrBusinessKodeKosong
	}
	switch strings.TrimSpace(b.IsActive) {
	case BusinessAktif, BusinessNonaktif:
		return nil
	default:
		return fmt.Errorf("%w: %q", ErrBusinessAktifTakSah, b.IsActive)
	}
}

// BusinessAktifTCO - apakah baris dianggap aktif oleh hilir (`isactive='1'`).
func BusinessAktifTCO(b BusinessTreaty) bool {
	return strings.TrimSpace(b.IsActive) == BusinessAktif
}
