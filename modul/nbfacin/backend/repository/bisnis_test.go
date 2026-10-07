package repository

import (
	"strings"
	"testing"
)

// TestSQLKelasBisnis - teks SQL tiket 28: PERSIS dua kolom DDL BUSINESS, tabel yang sudah
// dikualifikasi, saringan group business lewat parameter terikat, NOTE kosong dibuang,
// urut NOTE lalu ID (A74, A75), tanpa paging.
func TestSQLKelasBisnis(t *testing.T) {
	mau := "SELECT ID, NOTE FROM UJI.BUSINESS WHERE BUSINESSGROUPID = :1 AND NOTE IS NOT NULL ORDER BY NOTE, ID"
	if got := sqlKelasBisnis("UJI.BUSINESS"); got != mau {
		t.Errorf("SQL\n%q\nmau\n%q", got, mau)
	}
	s := sqlKelasBisnis("X")
	for _, terlarang := range []string{"SELECT *", "UPPER", "LIKE", "OFFSET", "FETCH", "DESC"} {
		if strings.Contains(s, terlarang) {
			t.Errorf("SQL memuat %q: %q", terlarang, s)
		}
	}
	if TabelBisnis != "BUSINESS" {
		t.Error("nama tabel salah")
	}
}
