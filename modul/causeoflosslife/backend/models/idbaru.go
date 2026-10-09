package models

// ID baru - rumus prosedur warisan PEGA_M_CAUSEOFLOSS_LIFE (fakta WO 08-10-2026, WAJIB ditiru persis; K3):
//
//	id := '1' || LPAD(M_CAUSEOFLOSS_LIFE_SEQ.NEXTVAL, 5, '0')
//
// Huruf '1' TETAP (bukan id_site). ID DEV yang ada: 100001-100004 (sequence last_number 5; ID berikutnya 100005).
//
// ⛔ Satu penyimpangan dari LPAD Oracle, disengaja (pola benefitlife): nomor sequence yang lebih panjang dari 5 angka
// dipotong Oracle (`123456` -> `12345`, ID kembar diam-diam); di sini DITOLAK dengan galat yang jelas.

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// AwalanID - huruf tetap di depan ID (`'1' ||`, PEGA_M_CAUSEOFLOSS_LIFE).
	AwalanID = "1"
	// LebarNomorID - LPAD(..., 5, '0').
	LebarNomorID = 5
)

// ErrIDTidakSah - nomor sequence tidak membentuk ID yang sah.
var ErrIDTidakSah = errors.New("models: ID Cause of Loss tidak dapat dibentuk")

var polaBulat = regexp.MustCompile(`^[0-9]+$`)

// BentukID - '1' || LPAD(nomor, 5, '0'); nomor = teks `TO_CHAR(M_CAUSEOFLOSS_LIFE_SEQ.NEXTVAL)`.
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
