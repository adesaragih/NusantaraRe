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

func TestSQLInboxSemuaKasusTerbaruDahulu(t *testing.T) {
	q := sqlInbox(uPolis)
	if n := penampungUnik(t, "sqlInbox", q); n != 3 {
		t.Fatalf("sqlInbox: %d penampung, mau 3 (pola, offset, ukuran)", n)
	}
	// Kasus tertutup ikut terdaftar (09-10-2026): tanpa saringan STATUSS.
	if strings.Contains(q, "STATUSS IS NULL") || strings.Contains(sqlCacahInbox(uPolis), "STATUSS") {
		t.Errorf("kotak masuk masih menyaring STATUSS: %s", q)
	}
	for _, wajib := range []string{"FROM " + uPolis + " p", "WHERE p.ID LIKE :1\n",
		"ORDER BY p.TGL_INPUT DESC, p.ID", "OFFSET :2 ROWS FETCH NEXT :3 ROWS ONLY"} {
		if !strings.Contains(q, wajib) {
			t.Errorf("sqlInbox tanpa %q", wajib)
		}
	}
	// 07-10-2026: Endorsement No dari NO_ENDORS (Policy No tetap OLD_POLICY_NO); empat kolom dihapus.
	if !strings.Contains(q, "SELECT p.ID, p.NO_ENDORS, p.OLD_POLICY_NO, p.TYPE, p.EDM_TYPE, p.CEDING_CO_NAME,") {
		t.Errorf("sqlInbox bukan urutan kolom kotak masuk 07-10-2026: %s", q)
	}
	for _, buang := range []string{"SOB_NAME", "POLICY_HOLDER_NAME", "MARKETING_NAME", "EDM_TYPE_BATAL"} {
		if strings.Contains(q, buang) {
			t.Errorf("sqlInbox masih membaca %s", buang)
		}
	}
	if !strings.HasSuffix(sqlCacahInbox(uPolis), "WHERE p.ID LIKE :1") {
		t.Error("cacah kotak masuk tidak memakai saringan yang sama")
	}
	if polaPengenalKasus != "EDML-%" {
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
	// Sub-jenis EDM Type (483) sebelum kolom kepala salinan; cacah kolom tetap sepakat dengan pemindai.
	if !strings.Contains(biasa, "p.CREATE_OP_NAME, p.EDM_TYPE_PERUBAHAN_DATA, p.EDM_TYPE_BATAL, ") {
		t.Error("sqlKasus tidak membaca sub-jenis EDM Type")
	}
	// Koma di dalam kurung (NVL, TO_CHAR) bukan pemisah kolom.
	n, dalam := 1, 0
	for _, c := range kolomKasus {
		switch {
		case c == '(':
			dalam++
		case c == ')':
			dalam--
		case c == ',' && dalam == 0:
			n++
		}
	}
	if n != jumlahKolomKasus {
		t.Errorf("kolomKasus %d kolom, jumlahKolomKasus %d", n, jumlahKolomKasus)
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
			// PRODKE 0-based (08-10-2026); warisan diurut sebagai BILANGAN (VARCHAR2: "9" > "10" bila teks).
			if !strings.HasPrefix(urut, "ORDER BY NVL(p.PROD_KE, 0) DESC") &&
				!strings.HasPrefix(urut, "ORDER BY NVL(TO_NUMBER(TRIM(j.PRODKE) DEFAULT NULL ON CONVERSION ERROR), 0) DESC") {
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

// TestVersiPRODKEMulaiNol - keputusan work owner 08-10-2026 (migrasi 485): SATU hitungan PRODKE, kosong = 0;
// penghitung terpisah urutan Pega (K3) dihapus - tiap jalur hanya ID, PRODKE, EDM Type.
func TestVersiPRODKEMulaiNol(t *testing.T) {
	pilihan := func(q string) string {
		return strings.Join(strings.Fields(q[len("SELECT "):strings.Index(q, " FROM ")]), " ")
	}
	for nama, k := range map[string][2]string{
		"sqlVersiEDM":     {sqlVersiEDM(uPolis, false), "p.ID, NVL(p.PROD_KE, 0), p.EDM_TYPE"},
		"sqlVersiNB":      {sqlVersiNB(uPolis, uPeserta, false), "p.ID, NVL(p.PROD_KE, 0), p.EDM_TYPE"},
		"sqlVersiWarisan": {sqlVersiWarisan(uJSON, false), "j.IDPEGA, " + prodKeWarisan + ", JSON_VALUE(j.DATA_JSON, '$.EdmType' NULL ON ERROR)"},
	} {
		if got := pilihan(k[0]); got != k[1] {
			t.Errorf("%s: SELECT %q, mau %q", nama, got, k[1])
		}
		if strings.Contains(k[0], ", 1)") || strings.Contains(k[0], "PL_NUMBER_EDM") {
			t.Errorf("%s: masih memakai skema lama (NVL …, 1 / nomor sebagai urutan)", nama)
		}
	}
	if !strings.Contains(sqlKasus(uPolis, false), "NVL(p.PROD_KE, 0)") {
		t.Error("kepala kasus membaca PROD_KE kosong bukan sebagai 0")
	}
}

// TestVersiWarisanTanpaBarisSistemBaru - JSON_POLIS kini juga memuat baris PremiumList (`NBLF-`, PRODKE '0') dan
// endorsement modul ini (`EDML-`); keduanya dibaca jalur aplikasi, jadi jalur warisan mengecualikannya.
func TestVersiWarisanTanpaBarisSistemBaru(t *testing.T) {
	for _, sebelum := range []bool{false, true} {
		w := sqlVersiWarisan(uJSON, sebelum)
		if !strings.Contains(w, "j.NOPOLIS = :1 AND j.IDPEGA NOT LIKE :2 AND j.IDPEGA NOT LIKE :3") {
			t.Errorf("warisan tidak mengecualikan IDPEGA sistem baru: %s", w)
		}
		if sebelum && !strings.Contains(w, prodKeWarisan+" < :4") {
			t.Errorf("batas sebelum warisan: %s", w)
		}
	}
	if polaIDPegaAplikasi[0] != "NBLF-%" || polaIDPegaAplikasi[1] != "EDML-%" {
		t.Errorf("pola IDPEGA aplikasi = %q", polaIDPegaAplikasi)
	}
}

// TestSQLJSONPolisEndorsement - Confirm menulis JSON_POLIS (08-10-2026): kolom bernama, DATA_JSON tidak diisi.
func TestSQLJSONPolisEndorsement(t *testing.T) {
	if q := sqlHapusJSONPolisEDM(uJSON); q != "DELETE FROM "+uJSON+" WHERE IDPEGA = :1" {
		t.Errorf("hapus = %q", q)
	}
	q := sqlSisipJSONPolisEDM(uJSON, uPolis)
	penampungUnik(t, "sqlSisipJSONPolisEDM", q)
	for _, k := range []string{
		"INSERT INTO " + uJSON + " (IDPEGA, NOPOLIS, PRODKE, NOENDORS, TGL_PROD, USERNAME, TGL_INPUT)",
		"SELECT :1, :2, :3, :4, :5, p.CREATE_OP_NAME, SYSDATE FROM " + uPolis + " p", "WHERE p.ID = :6",
	} {
		if !strings.Contains(q, k) {
			t.Errorf("sisip tanpa %q: %s", k, q)
		}
	}
	for _, larangan := range []string{"DATA_JSON", "KONVERSI", "OLDNOPOLIS", "ERR_NOTE", "COMMIT", "SELECT *"} {
		if strings.Contains(q, larangan) {
			t.Errorf("sisip memuat %q", larangan)
		}
	}
	if q := sqlHariClosing("S.TANGGAL_CLOSING"); q != "SELECT TANGGAL FROM S.TANGGAL_CLOSING WHERE ROWNUM = 1" {
		t.Errorf("hari closing = %q", q)
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
