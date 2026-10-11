package models

// Policy No di Create Addendum boleh diisi NOMOR ENDORSEMENT (keputusan work owner 10-10-2026, bukan korpus):
// akhiran "/<angka>" paling akhir dibuang, PL Number-nya yang dipakai. Versi yang ditarik TETAP versi PRODKE
// terbesar (`VersiBerjalan`) - akhiran yang diketik tidak memilih versi.

import (
	"fmt"
	"regexp"
	"strings"
)

// polaAkhiranEndorsement - "/" + satu atau lebih digit di posisi paling akhir.
var polaAkhiranEndorsement = regexp.MustCompile(`^(.+)/(\d+)$`)

// NomorPolisDariMasukan - masukan Policy No (di-trim) → nomor polis tanpa akhiran endorsement, akhirannya, dan
// apakah ada akhiran. Hanya "/<digit>" paling akhir yang dibuang; "/" lain dan akhiran bukan digit apa adanya.
func NomorPolisDariMasukan(s string) (nomorPolis, akhiran string, adaAkhiran bool) {
	t := strings.TrimSpace(s)
	m := polaAkhiranEndorsement.FindStringSubmatch(t)
	if m == nil {
		return t, "", false
	}
	return strings.TrimSpace(m[1]), m[2], true
}

// AkhiranVersi - akhiran nomor endorsement versi ber-PRODKE ini ("%02d", pola `NomorEndorsement`); kosong untuk
// versi new business (PRODKE 0).
func AkhiranVersi(prodKe int) string {
	if prodKe < 1 {
		return ""
	}
	return fmt.Sprintf("%02d", prodKe)
}
