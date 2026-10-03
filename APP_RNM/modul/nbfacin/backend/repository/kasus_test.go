package repository

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// kolomMigrasi - kolom CREATE TABLE satu tabel di satu berkas migrasi modul ini.
func kolomMigrasi(t *testing.T, berkas, tabel string) map[string]bool {
	t.Helper()
	blok := regexp.MustCompile(`(?s)CREATE TABLE \{skema\}\.` + tabel + ` \((.*?)\n\)`).FindStringSubmatch(bacaMigrasi(t, berkas))
	if blok == nil {
		t.Fatalf("CREATE TABLE %s tidak terbaca di %s", tabel, berkas)
	}
	kolom := map[string]bool{}
	for _, b := range strings.Split(blok[1], "\n") {
		if f := strings.Fields(b); len(f) > 0 && f[0] != "CONSTRAINT" {
			kolom[f[0]] = true
		}
	}
	return kolom
}

// TestSQLKasusMemakaiKolomMigrasi - setiap kolom g./q. yang dibaca dan setiap kolom yang
// ditulis SimpanGeneral ADA di migrasi 182/183 (dua arah: kolom migrasi yang tidak
// dipakai juga dilaporkan, kecuali kolom sistem IDPEGA/COB_GROUP yang diisi loader).
func TestSQLKasusMemakaiKolomMigrasi(t *testing.T) {
	gen := kolomMigrasi(t, "182_t_general_polis.sql", "T_GENERAL_POLIS")
	quo := kolomMigrasi(t, "183_t_quotationdata.sql", "T_QUOTATIONDATA")
	baca := sqlBacaKasus("W", "O", "A", "G", "Q")
	tulisG := sqlUbahGeneral("G") + " " + sqlSisipGeneral("G")
	tulisQ := sqlUbahQuotation("Q") + " " + sqlSisipQuotation("Q")
	pakai := map[string]map[string]bool{"G": {}, "Q": {}}
	for _, m := range regexp.MustCompile(`\b([gq])\.([A-Z_]+)`).FindAllStringSubmatch(baca, -1) {
		pakai[strings.ToUpper(m[1])][m[2]] = true
	}
	for _, m := range regexp.MustCompile(`\b([A-Z][A-Z_]+)\b`).FindAllStringSubmatch(tulisG, -1) {
		if gen[m[1]] {
			pakai["G"][m[1]] = true
		}
	}
	for _, m := range regexp.MustCompile(`\b([A-Z][A-Z_]+)\b`).FindAllStringSubmatch(tulisQ, -1) {
		if quo[m[1]] {
			pakai["Q"][m[1]] = true
		}
	}
	for tb, kolom := range map[string]map[string]bool{"G": gen, "Q": quo} {
		for k := range pakai[tb] {
			if !kolom[k] {
				t.Errorf("%s.%s dipakai SQL, tidak ada di migrasi", tb, k)
			}
		}
		for k := range kolom {
			if !pakai[tb][k] && k != "IDPEGA" && k != "COB_GROUP" {
				t.Errorf("%s.%s di migrasi tidak dipakai SQL", tb, k)
			}
		}
	}
	// SQL tulis tidak menyentuh medan tampil-saja.
	for _, k := range []string{"SOB_NAME", "CEDING_CO_NAME", "GROUP_NAME", "FOLLOWING"} {
		if strings.Contains(tulisG+tulisQ, k) {
			t.Errorf("%s tampil-saja tetapi ditulis", k)
		}
	}
}

// TestSQLKasus - teks SQL tiket 31: case dibatasi LINI (bind :2), nama tertanggung lewat
// subkueri MAX (T_M_ACCOUNT tanpa PK), sentuh case = kunci baris + 404, bind bernomor.
func TestSQLKasus(t *testing.T) {
	baca := sqlBacaKasus("UJI.W", "UJI.O", "UJI.A", "UJI.G", "UJI.Q")
	for _, harus := range []string{"FROM UJI.W w", "LEFT JOIN UJI.O o ON o.ID = w.ID", "LEFT JOIN UJI.G g ON g.ID = w.ID",
		"LEFT JOIN UJI.Q q ON q.PARENT_ID = w.ID", "(SELECT MAX(a.INSUREDNAME) FROM UJI.A a WHERE a.ID = o.ACCOUNT_ID)",
		"WHERE w.ID = :1 AND w.LINI = :2"} {
		if !strings.Contains(baca, harus) {
			t.Errorf("baca tanpa %q", harus)
		}
	}
	pilih := baca[len("SELECT "):strings.Index(baca, "\nFROM")]
	if n := len(strings.Split(pilih, ",")); n != len(kolomBacaKasus) || n != 32 {
		t.Errorf("SELECT %d kolom, daftar %d, mau 32 (3 case + 14 opportunity + nama + 5 general + 9 quotation)", n, len(kolomBacaKasus))
	}
	for got, mau := range map[string]string{
		sqlSentuhCase("UJI.W"):   "UPDATE UJI.W SET TGL_UPDATE = SYSDATE WHERE ID = :1 AND LINI = :2",
		sqlUbahGeneral("UJI.G"):  "UPDATE UJI.G SET START_DATE_TIME = :1, OFFERING_DATE = :2, END_DATE_TIME = :3 WHERE ID = :4",
		sqlSisipGeneral("UJI.G"): "INSERT INTO UJI.G (ID, START_DATE_TIME, OFFERING_DATE, END_DATE_TIME) VALUES (:1, :2, :3, :4)",
		sqlUbahQuotation("UJI.Q"): "UPDATE UJI.Q SET NO_OFFER_SLIP = :1, QQ_NAME = :2, POLICY_TYPE = :3, MOID = :4, EDM_DAY = :5," +
			" TYPE_FACULTATIVE = :6 WHERE PARENT_ID = :7",
		sqlSisipQuotation("UJI.Q"): "INSERT INTO UJI.Q (ID, PARENT_ID, NO_OFFER_SLIP, QQ_NAME, POLICY_TYPE, MOID, EDM_DAY, TYPE_FACULTATIVE)" +
			" VALUES (:1, :2, :3, :4, :5, :6, :7, :8)",
	} {
		if got != mau {
			t.Errorf("SQL\n%q\nmau\n%q", got, mau)
		}
	}
	if !strings.Contains(bacaMigrasi(t, "183_t_quotationdata.sql"), "CREATE SEQUENCE {skema}."+SequenceQuotationData+" START WITH 1") {
		t.Error("183 tanpa " + SequenceQuotationData)
	}
}

// TestBacaKasusMemakaiKolomBernama - setiap kunci v("...")/teks["..."] di BacaKasus ada
// di kolomBacaKasus (salah ketik = panic saat Oracle), dan setiap kolom teks terbaca.
func TestBacaKasusMemakaiKolomBernama(t *testing.T) {
	b, err := os.ReadFile("kasus.go")
	if err != nil {
		t.Fatal(err)
	}
	ada := map[string]bool{}
	for _, k := range kolomBacaKasus {
		ada[k] = true
	}
	dibaca := map[string]bool{}
	for _, m := range regexp.MustCompile(`(?:v\(|teks\[)"([^"]+)"`).FindAllStringSubmatch(string(b), -1) {
		dibaca[m[1]] = true
		if !ada[m[1]] {
			t.Errorf("kunci %q dibaca tetapi tidak ada di kolomBacaKasus", m[1])
		}
	}
	for _, k := range kolomBacaKasus {
		if k != kolomTanggalKasus && !dibaca[k] {
			t.Errorf("kolom %s dipilih tetapi tidak dibaca", k)
		}
	}
}
