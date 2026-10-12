package repository

// Migrasi 025 (keputusan work owner 10-10-2026 "iya tambahkan") - bentuk berkas dan kecocokan dengan pembaca/penulis.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrasi025KolomIdPegaDanIsiUlang(t *testing.T) {
	baca := func(nama string) string {
		b, err := os.ReadFile(filepath.Join("..", "migrations", nama))
		if err != nil {
			t.Fatal(err)
		}
		return strings.Join(strings.Fields(buangKomentarSumber(nama, string(b))), " ")
	}
	maju := baca("025_kolom_idpega_peserta.sql")
	for _, wajib := range []string{
		// kolom baru lewat ALTER biasa (penjaga STRUKTUR membaca jalur maju), lebar = M_LIFE_PREMIUM_DETAIL.IDPEGA
		"ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL ADD ( IDPEGA VARCHAR2(50) )",
		// isi ulang aman diulang, berkunci index PL_NUMBER + CERTIFICATE_NO + ID = SOURCE_ID, satu baris
		"UPDATE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL d SET d.IDPEGA = (SELECT m.IDPEGA FROM {skema}.M_LIFE_PREMIUM_DETAIL m",
		"WHERE m.PL_NUMBER = d.PL_NUMBER AND m.CERTIFICATE_NO = d.CERTIFICATE_NO AND m.ID = d.SOURCE_ID AND ROWNUM = 1)",
		"WHERE d.IDPEGA IS NULL AND d.SOURCE_ID IS NOT NULL",
	} {
		if !strings.Contains(maju, wajib) {
			t.Errorf("025 tanpa %q:\n%s", wajib, maju)
		}
	}
	if strings.Contains(strings.ToUpper(maju), "COMMIT") {
		t.Error("025 memuat COMMIT")
	}
	mundur := baca("025_kolom_idpega_peserta_down.sql")
	for _, wajib := range []string{
		"TABLE_NAME = 'T_CLAIMLF_PREMIUMLIST_DETAIL' AND COLUMN_NAME = 'IDPEGA'",
		"IF n > 0 THEN EXECUTE IMMEDIATE 'ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL DROP COLUMN IDPEGA'",
	} {
		if !strings.Contains(mundur, wajib) {
			t.Errorf("025 mundur tanpa %q:\n%s", wajib, mundur)
		}
	}
	// Penulis dan pembaca peserta klaim membawa kolom yang dibuat 025.
	ada := false
	for _, k := range kolomPeserta {
		ada = ada || k.Nama == "IDPEGA"
	}
	if !ada {
		t.Error("kolomPeserta tidak menulis/membaca IDPEGA")
	}
}
