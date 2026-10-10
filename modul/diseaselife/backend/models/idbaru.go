package models

// ID baru - keputusan work owner 08-10-2026 D1.1: ID = TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL) (sequence aplikasi migrasi 080,
// mulai ID angka tertinggi + 1; DEV 197586).
//
// ⛔ Rumus Pega `'1' || LPAD(M_DISEASE_LIFE_SEQ.NEXTVAL, 5, '0')` (prosedur PEGA_DISEASE_LIFE) SENGAJA TIDAK dipakai:
// di DEV sequence itu menghasilkan 102051 yang sudah dipakai data impor (ID kembar `C718` / `TEST123`), dan 102052,
// 102053, 102054, ... juga sudah terpakai - setiap Add akan bertabrakan. M_DISEASE_LIFE_SEQ tidak disentuh.

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ErrIDTidakSah - nomor sequence tidak membentuk ID yang sah.
var ErrIDTidakSah = errors.New("models: ID Disease tidak dapat dibentuk")

var polaBulat = regexp.MustCompile(`^[0-9]+$`)

// BentukID - ID = teks nomor `TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL)` apa adanya; bukan angka atau lebih panjang dari
// VARCHAR2(BatasID) = galat yang jelas.
func BentukID(nomor string) (string, error) {
	nomor = strings.TrimSpace(nomor)
	if !polaBulat.MatchString(nomor) {
		return "", fmt.Errorf("%w: nomor sequence %q bukan angka", ErrIDTidakSah, nomor)
	}
	if len(nomor) > BatasID {
		return "", fmt.Errorf("%w: nomor sequence lebih dari %d angka", ErrIDTidakSah, BatasID)
	}
	return nomor, nil
}
