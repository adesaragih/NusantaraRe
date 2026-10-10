package models

// ID baru - rumus prosedur warisan `PEGA_M_RICOMM_LIFE` / `PEGA_M_RICOMM_LIFE_SUMMARY` (dicek work owner 06-10-2026):
//
//	SELECT ID INTO id_site FROM M_SITE_DATABASE WHERE CURRENT_SITE='1';            -- b13
//	id := id_site || lpad(to_Char(M_RICOMM_LIFE_SEQ.nextval),6,'0');              -- b21
//
// Keputusan work owner 06-10-2026 butir 4: site dibaca saat berjalan (galat bila tidak tepat satu baris), sequence
// warisan `M_RICOMM_LIFE_SEQ` (rincian) / `M_RICOMM_LIFE_SUMMARY_SEQ` (ringkasan); nol sequence baru. M_SITE_DATABASE
// DEV: (1,'1') dan (2,'0') -> ID 7 karakter, mis. 1000043.
//
// ⛔ Satu penyimpangan dari LPAD Oracle, disengaja: nomor sequence lebih dari 6 angka dipotong LPAD (`1234567` ->
// `123456`, ID kembar diam-diam); di sini DITOLAK. ID lebih dari BatasID (VARCHAR2(10)) juga ditolak.

import (
	"errors"
	"fmt"
	"strings"
)

// LebarNomorID - LPAD(..., 6, '0').
const LebarNomorID = 6

// ErrIDTidakSah - site atau nomor sequence tidak membentuk ID yang sah.
var ErrIDTidakSah = errors.New("models: ID R/I Comm Life tidak dapat dibentuk")

// BentukID - site || LPAD(nomor, 6, '0'). Keduanya teks angka (site = `TO_CHAR(M_SITE_DATABASE.ID)`).
func BentukID(site, nomor string) (string, error) {
	site, nomor = strings.TrimSpace(site), strings.TrimSpace(nomor)
	if !polaBulat.MatchString(site) {
		return "", fmt.Errorf("%w: site %q bukan angka", ErrIDTidakSah, site)
	}
	if !polaBulat.MatchString(nomor) {
		return "", fmt.Errorf("%w: nomor sequence %q bukan angka", ErrIDTidakSah, nomor)
	}
	if len(nomor) > LebarNomorID {
		return "", fmt.Errorf("%w: nomor sequence %s lebih dari %d angka (LPAD akan memotongnya)", ErrIDTidakSah, nomor, LebarNomorID)
	}
	id := site + strings.Repeat("0", LebarNomorID-len(nomor)) + nomor
	if len(id) > BatasID {
		return "", fmt.Errorf("%w: ID %s lebih dari %d karakter (VARCHAR2(%d))", ErrIDTidakSah, id, BatasID, BatasID)
	}
	return id, nil
}
