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

// TestPungutEfekTanpaRowLimitingDenganForUpdate - bug DEV 10-10-2026 (ORA-02014): `FETCH FIRST` tidak boleh bersama
// `FOR UPDATE`; kandidat lewat subkueri ROWNUM berurut JADWAL_BERIKUT, kunci di tabel dasar, bind :1-:3 sekali masing-masing
// dan berurutan (bind posisi), saringan MODUL tetap di dalam subkueri.
func TestPungutEfekTanpaRowLimitingDenganForUpdate(t *testing.T) {
	q := sqlPungutEfek("S.L")
	satu := strings.Join(strings.Fields(q), " ")
	if strings.Contains(strings.ToUpper(q), "FETCH FIRST") {
		t.Errorf("FETCH FIRST bersama FOR UPDATE:\n%s", q)
	}
	mau := "SELECT ID, LINI, MODUL, JENIS_EFEK, RUJUKAN, MUATAN, PERCOBAAN FROM S.L WHERE ID = (SELECT ID FROM (SELECT ID FROM S.L" +
		" WHERE STATUS = :1 AND JADWAL_BERIKUT <= :2 AND MODUL = :3 ORDER BY JADWAL_BERIKUT ASC) WHERE ROWNUM = 1)" +
		" FOR UPDATE SKIP LOCKED"
	if satu != mau {
		t.Errorf("teks SQL:\n%s\nmau:\n%s", satu, mau)
	}
	i1, i2, i3 := strings.Index(q, ":1"), strings.Index(q, ":2"), strings.Index(q, ":3")
	if !(i1 < i2 && i2 < i3) || strings.Count(q, ":1") != 1 || strings.Count(q, ":2") != 1 || strings.Count(q, ":3") != 1 ||
		strings.Contains(q, ":4") {
		t.Errorf("bind posisi:\n%s", q)
	}
	if !strings.HasSuffix(strings.TrimSpace(q), "FOR UPDATE SKIP LOCKED") || !strings.Contains(q, "WHERE ROWNUM = 1)") {
		t.Errorf("kunci bukan di tabel dasar:\n%s", q)
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
