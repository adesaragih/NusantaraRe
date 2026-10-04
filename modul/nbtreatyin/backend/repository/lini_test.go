package repository

// T_GENERAL_POLIS / T_WORK_POLIS BERSAMA lini lain (keputusan work owner
// 04-10-2026) - TANPA Oracle. Setiap GERBANG baca/tulis kasus menyaring
// `T_WORK_POLIS.LINI = models.LiniKasus` ('NONLIFE'), sehingga baris FacIn
// (`LINI = 'FAC'`) tidak tampil di daftar portal dan tidak terjangkau dari
// pengenal kasus. Gerbangnya: `DaftarKasus` (portal), `Keadaan` (setiap
// pembacaan dan tindakan layanan mulai di sini), `KunciKasus` (setiap
// transaksi tulis mengunci di sini). Perilaku lawan Oracle:
// `lini_db_test.go` (tag db).

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

// argLini - nomor penampung yang membandingkan LINI, atau "" bila tidak ada.
func argLini(q string) string {
	m := regexp.MustCompile(`(?i)\bLINI\s*=\s*:(\d+)`).FindStringSubmatch(q)
	if m == nil {
		return ""
	}
	return m[1]
}

func TestGerbangKasusMenyaringLini(t *testing.T) {
	if models.LiniKasus != "NONLIFE" {
		t.Fatalf("LiniKasus = %q, harap NONLIFE", models.LiniKasus)
	}
	// Portal: penampung :1 = LiniKasus.
	q, args := sqlDaftarKasus("S.W", "S.G", "S.Q", models.SaringanKasus{})
	if n := argLini(q); n != "1" || args[0] != models.LiniKasus {
		t.Errorf("daftar portal: LINI di penampung %q, argumen pertama %v:\n%s", n, args[0], q)
	}
	// Keadaan: :1 = ID, :2 = LINI.
	if n := argLini(sqlKeadaan("S.W", "S.G")); n != "2" {
		t.Errorf("Keadaan tidak menyaring LINI di :2 (%q):\n%s", n, sqlKeadaan("S.W", "S.G"))
	}
	// KunciKasus: :1 = ID, :2 = LINI, tetap FOR UPDATE.
	k := sqlKunciKasus("S.W")
	if n := argLini(k); n != "2" || !strings.Contains(k, "FOR UPDATE") {
		t.Errorf("KunciKasus tidak menyaring LINI di :2 (%q) atau tanpa FOR UPDATE:\n%s", n, k)
	}
	for _, s := range []string{sqlKeadaan("S.W", "S.G"), k} {
		if c := len(regexp.MustCompile(`:(\d+)`).FindAllString(s, -1)); c != 2 {
			t.Errorf("harap tepat dua penampung (ID, LINI), dapat %d:\n%s", c, s)
		}
	}
}
