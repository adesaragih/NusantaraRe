package models

// ID baru - DUA rumus prosedur warisan (fakta WO 08-10-2026, teks prosedur Pega - WAJIB ditiru persis):
//
//	PEGA_M_RIRISK_LIFE_SUMMARY: id := id_site || LPAD(M_RIRISK_LIFE_SUMMARY_SEQ.NEXTVAL, 6, '0')
//	                            id_site = SELECT ID FROM M_SITE_DATABASE WHERE CURRENT_SITE='1'
//	PEGA_M_RIRISK_LIFE:         id := '1' || LPAD(M_RIRISK_LIFE_SEQ.NEXTVAL, 5, '0')
//
// Rincian: huruf '1' TETAP (bukan id_site) dan lebar ID 6 (`RIRISK_LIFE.ID` VARCHAR2(6)). Ringkasan: site saat
// berjalan (galat bila tidak tepat satu baris), lebar ID VARCHAR2(10). Sequence warisan; nol sequence baru.
//
// ⛔ Satu penyimpangan dari LPAD Oracle, disengaja: nomor sequence yang lebih panjang dari lebar LPAD dipotong Oracle
// (`123456` -> `12345` pada LPAD 5, ID kembar diam-diam); di sini DITOLAK. ID yang melewati lebar kolomnya juga ditolak.

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// LebarNomorID - LPAD(..., 6, '0') ringkasan.
	LebarNomorID = 6
	// LebarNomorIDRincian - LPAD(..., 5, '0') rincian.
	LebarNomorIDRincian = 5
	// AwalanIDRincian - huruf tetap di depan ID rincian (`'1' ||`, PEGA_M_RIRISK_LIFE).
	AwalanIDRincian = "1"
)

// ErrIDTidakSah - site atau nomor sequence tidak membentuk ID yang sah.
var ErrIDTidakSah = errors.New("models: ID R/I Risk tidak dapat dibentuk")

func lpad(nomor string, lebar int) (string, error) {
	nomor = strings.TrimSpace(nomor)
	if !polaBulat.MatchString(nomor) {
		return "", fmt.Errorf("%w: nomor sequence %q bukan angka", ErrIDTidakSah, nomor)
	}
	if len(nomor) > lebar {
		return "", fmt.Errorf("%w: nomor sequence %s lebih dari %d angka (LPAD akan memotongnya)", ErrIDTidakSah, nomor, lebar)
	}
	return strings.Repeat("0", lebar-len(nomor)) + nomor, nil
}

// BentukID - ID RINGKASAN: site || LPAD(nomor, 6, '0'). Keduanya teks angka (site = `TO_CHAR(M_SITE_DATABASE.ID)`).
func BentukID(site, nomor string) (string, error) {
	site = strings.TrimSpace(site)
	if !polaBulat.MatchString(site) {
		return "", fmt.Errorf("%w: site %q bukan angka", ErrIDTidakSah, site)
	}
	n, err := lpad(nomor, LebarNomorID)
	if err != nil {
		return "", err
	}
	id := site + n
	if len(id) > BatasID {
		return "", fmt.Errorf("%w: ID %s lebih dari %d karakter (VARCHAR2(%d))", ErrIDTidakSah, id, BatasID, BatasID)
	}
	return id, nil
}

// BentukIDRincian - ID RINCIAN: '1' || LPAD(nomor, 5, '0') - lebar 6 = `RIRISK_LIFE.ID` VARCHAR2(6).
func BentukIDRincian(nomor string) (string, error) {
	n, err := lpad(nomor, LebarNomorIDRincian)
	if err != nil {
		return "", err
	}
	id := AwalanIDRincian + n
	if len(id) > BatasIDRincian {
		return "", fmt.Errorf("%w: ID %s lebih dari %d karakter (VARCHAR2(%d))", ErrIDTidakSah, id, BatasIDRincian, BatasIDRincian)
	}
	return id, nil
}
