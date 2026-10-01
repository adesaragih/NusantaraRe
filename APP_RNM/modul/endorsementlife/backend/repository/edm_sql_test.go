package repository

// Bentuk SQL modul Endorsement Life - tanpa Oracle.

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/modul/endorsementlife/backend/models"
)

const (
	uPolis   = "UJISKEMA.T_PREMIUM_LIST"
	uPeserta = "UJISKEMA.T_PREMIUM_LIST_DETAIL"
	uJSON    = "UJISKEMA.JSON_POLIS"
)

var polaPenampung = regexp.MustCompile(`:(\d+)`)

// penampungUnik - setiap penampung `:n` muncul sekali dan berurutan 1..N
// (ORA-01008 bila berulang bersama OFFSET/FETCH).
func penampungUnik(t *testing.T, nama, q string) int {
	t.Helper()
	m := polaPenampung.FindAllStringSubmatch(q, -1)
	for i, x := range m {
		if x[1] != strconv.Itoa(i+1) {
			t.Errorf("%s: penampung ke-%d adalah :%s (berulang atau tak berurutan)", nama, i+1, x[1])
		}
	}
	if err := db.PeriksaSQL(q); err != nil {
		t.Errorf("%s: %v", nama, err)
	}
	return len(m)
}

func nullString(s string) sql.NullString { return sql.NullString{String: s, Valid: s != ""} }

func TestSQLInboxMenyaringKasusTerbukaTerbaruDahulu(t *testing.T) {
	q := sqlInbox(uPolis)
	if n := penampungUnik(t, "sqlInbox", q); n != 3 {
		t.Fatalf("sqlInbox: %d penampung, mau 3 (pola, offset, ukuran)", n)
	}
	for _, wajib := range []string{"FROM " + uPolis + " p", "p.ID LIKE :1", "p.STATUSS IS NULL",
		"ORDER BY p.TGL_INPUT DESC, p.ID", "OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY"} {
		if !strings.Contains(q, wajib) {
			t.Errorf("sqlInbox tanpa %q", wajib)
		}
	}
	if !strings.Contains(sqlCacahInbox(uPolis), "p.ID LIKE :1 AND p.STATUSS IS NULL") {
		t.Error("cacah kotak masuk tidak memakai saringan yang sama")
	}
	if polaPengenalKasus != "EDMLF-%" {
		t.Errorf("pola pengenal %q", polaPengenalKasus)
	}
}

func TestSQLKasusMembacaKepalaSalinanDanMengunci(t *testing.T) {
	biasa, kunci := sqlKasus(uPolis, false), sqlKasus(uPolis, true)
	penampungUnik(t, "sqlKasus", biasa)
	if strings.Contains(biasa, "FOR UPDATE") || !strings.HasSuffix(kunci, " FOR UPDATE") {
		t.Error("FOR UPDATE hanya bila diminta")
	}
	if !strings.Contains(biasa, "WHERE p.ID = :1 AND p.ID LIKE :2") {
		t.Error("sqlKasus tidak membatasi ke kasus EDM")
	}
	for _, k := range models.KolomKepalaSalin {
		if !strings.Contains(biasa, "p."+k.Kolom) {
			t.Errorf("sqlKasus tidak membaca %s", k.Kolom)
		}
		if k.Tanggal && !strings.Contains(biasa, "TO_CHAR(p."+k.Kolom+", 'YYYY-MM-DD')") {
			t.Errorf("tanggal %s tidak dibaca berpola", k.Kolom)
		}
	}
}

func TestSQLPesertaBerhalamanDanBerurutanStabil(t *testing.T) {
	q := sqlPeserta(uPeserta, models.KolomPesertaGrid)
	if n := penampungUnik(t, "sqlPeserta", q); n != 3 {
		t.Fatalf("%d penampung", n)
	}
	for _, wajib := range []string{"d.PREMIUM_LIST_ID = :1", "ORDER BY d.CERTIFICATE_NO, d.NAME_OF_INSURED, d.ID",
		"TO_CHAR(d.DOB, 'YYYY-MM-DD')", "TO_CHAR(d.ENTRY_AGE, 'TM9', 'NLS_NUMERIC_CHARACTERS=''.,''')"} {
		if !strings.Contains(q, wajib) {
			t.Errorf("sqlPeserta tanpa %q", wajib)
		}
	}
}

// TestPRODKESatuUrutan - AC 27: setiap pembaca `PRODKE` memakai `PROD_KE DESC`;
// test ini GAGAL bila ada jalur yang mengurut `TGL_INPUT` lebih dulu
// (`GetProdKeOldData_SQL` b84, penyimpangan sadar 1).
func TestPRODKESatuUrutan(t *testing.T) {
	for _, sebelum := range []bool{false, true} {
		for nama, q := range map[string]string{
			"sqlVersiEDM":     sqlVersiEDM(uPolis, sebelum),
			"sqlVersiNB":      sqlVersiNB(uPolis, uPeserta, sebelum),
			"sqlVersiWarisan": sqlVersiWarisan(uJSON, sebelum),
		} {
			penampungUnik(t, nama, q)
			urut := q[strings.Index(q, "ORDER BY"):]
			if !strings.HasPrefix(urut, "ORDER BY NVL(p.PROD_KE, 1) DESC") && !strings.HasPrefix(urut, "ORDER BY NVL(j.PRODKE, 1) DESC") {
				t.Errorf("%s: urutan %q bukan PRODKE DESC", nama, urut)
			}
			if !strings.Contains(q, "FETCH FIRST 1 ROWS ONLY") {
				t.Errorf("%s: tanpa batas satu baris", nama)
			}
			if sebelum != strings.Contains(q, "< :") {
				t.Errorf("%s: batas versi sebelum = %v tidak tercermin", nama, sebelum)
			}
		}
	}
	if !strings.Contains(sqlVersiEDM(uPolis, false), "p.NO_POLIS = :1 AND p.STATUSS = :2") {
		t.Error("versi EDM resmi tidak dicari lewat NO_POLIS + Resolved-Completed")
	}
	if !strings.Contains(sqlVersiNB(uPolis, uPeserta, false), "p.EDM_TYPE IS NULL") ||
		!strings.Contains(sqlVersiNB(uPolis, uPeserta, false), "d.PL_NUMBER = :1") {
		t.Error("versi NB tidak dicari lewat PL_NUMBER pesertanya")
	}
	w := sqlVersiWarisan(uJSON, false)
	if !strings.Contains(w, "j.NOPOLIS = :1") || !strings.Contains(w, "JSON_VALUE(j.DATA_JSON, '$.EdmType' NULL ON ERROR)") {
		t.Error("versi warisan tidak membaca NOPOLIS + EdmType JSON_POLIS")
	}
}

func TestRapikanAngkaTM9(t *testing.T) {
	for masuk, mau := range map[string]string{".5": "0.5", "-.25": "-0.25", "1234.5": "1234.5", "-1000": "-1000", "": "", "0": "0"} {
		g, err := rapikanAngka("X", nullString(masuk))
		if err != nil || g != mau {
			t.Errorf("rapikanAngka(%q) = %q, %v; mau %q", masuk, g, err, mau)
		}
	}
	if _, err := rapikanAngka("X", nullString("abc")); err == nil {
		t.Error("angka rusak tidak bergalat")
	}
}
