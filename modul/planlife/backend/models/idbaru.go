package models

// ID baru - rumus prosedur warisan PEGA_M_PRODUCT_TYPE_LIFE (fakta WO 08-10-2026, ditiru persis; K3):
//
//	id := '1' || LPAD(M_PRODUCT_TYPE_LIFE_SEQ.NEXTVAL, 5, '0')
//
// Huruf '1' TETAP; lebar 6 = `PRODUCT_TYPE_LIFE.ID` VARCHAR2(6). ⛔ Penyimpangan disengaja (pola benefitlife): nomor
// sequence yang lebih panjang dari 5 angka DITOLAK (LPAD Oracle memotongnya diam-diam -> ID kembar).

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// AwalanID - huruf tetap di depan ID (`'1' ||`).
	AwalanID = "1"
	// LebarNomorID - LPAD(..., 5, '0').
	LebarNomorID = 5
)

// ErrIDTidakSah - nomor sequence tidak membentuk ID yang sah.
var ErrIDTidakSah = errors.New("models: ID Plan tidak dapat dibentuk")

var polaBulat = regexp.MustCompile(`^[0-9]+$`)

// BentukID - '1' || LPAD(nomor, 5, '0'); nomor = teks `TO_CHAR(M_PRODUCT_TYPE_LIFE_SEQ.NEXTVAL)`.
func BentukID(nomor string) (string, error) {
	nomor = strings.TrimSpace(nomor)
	if !polaBulat.MatchString(nomor) {
		return "", fmt.Errorf("%w: nomor sequence %q bukan angka", ErrIDTidakSah, nomor)
	}
	if len(nomor) > LebarNomorID {
		return "", fmt.Errorf("%w: nomor sequence %s lebih dari %d angka (LPAD akan memotongnya)", ErrIDTidakSah, nomor, LebarNomorID)
	}
	return AwalanID + strings.Repeat("0", LebarNomorID-len(nomor)) + nomor, nil
}
