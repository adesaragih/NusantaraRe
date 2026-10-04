package outbox

import (
	"strings"
	"testing"
)

// TestPungutEfekHanyaModulnyaSendiri - temuan /code-review giliran 10.
//
// ⛔ Outbox `T_LOG_SERVICE_RNM` dipakai dua modul sejak tiket 06 PremiumList
// Life. Worker yang memungut baris modul lain menandainya gagal permanen
// karena pelaksananya tidak mengenal jenisnya.
func TestPungutEfekHanyaModulnyaSendiri(t *testing.T) {
	q := sqlPungutEfek("S.L")
	if !strings.Contains(q, "AND MODUL = :3") {
		t.Errorf("pemungutan efek tidak disaring modul:\n%s", q)
	}
	if !strings.Contains(q, "FOR UPDATE SKIP LOCKED") {
		t.Errorf("SKIP LOCKED hilang:\n%s", q)
	}
}

// TestEfekSudahSelesaiMengecualikanDirinya - tiket 07 Komite.
func TestEfekSudahSelesaiMengecualikanDirinya(t *testing.T) {
	q := sqlEfekSudahSelesai("S.L")
	for _, s := range []string{"MODUL = :1", "JENIS_EFEK = :2", "RUJUKAN = :3", "STATUS = :4", "ID <> :5"} {
		if !strings.Contains(q, s) {
			t.Errorf("anti-dobel tanpa %q:\n%s", s, q)
		}
	}
}
