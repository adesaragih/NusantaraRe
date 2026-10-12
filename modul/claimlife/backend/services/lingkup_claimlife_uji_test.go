package services

import "testing"

// Lingkup SQL penjaga Claim Life (11-10-2026): menyebut tabel Claim Life lewat FOREIGN KEY BUKAN sasaran; ALTER /
// CREATE / UPDATE atas tabel Claim Life dari modul mana pun tetap diperiksa.
func TestSqlLingkupClaimLife(t *testing.T) {
	lain := "modul/claimuji/backend/migrations/900_uji.sql"
	fk := "CREATE TABLE {skema}.T_UJI_ITEM (\n  TOTAL_UJI NUMBER,\n" +
		"  CONSTRAINT FK_UJI FOREIGN KEY (CLAIM_ID) REFERENCES {skema}.T_GENERAL_CLAIM (ID)\n)\n/\n"
	if got := sqlLingkupClaimLife(lain, fk); len(got) != 0 {
		t.Errorf("FK ke T_GENERAL_CLAIM ikut diperiksa: %q", got)
	}
	ubah := fk + "ALTER TABLE {skema}.T_CLAIMLF_PREMIUMLIST_DETAIL ADD (TOTAL_UJI NUMBER)\n/\n" +
		"UPDATE {skema}.T_GENERAL_CLAIM SET X = 1; COMMENT ON COLUMN {skema}.T_UJI_ITEM.TOTAL_UJI IS 'x';\n"
	if got := sqlLingkupClaimLife(lain, ubah); len(got) != 2 {
		t.Errorf("mau 2 pernyataan bersasaran Claim Life (ALTER + UPDATE), dapat %d: %q", len(got), got)
	}
	komentar := "-- ALTER TABLE {skema}.T_CLAIMLF_X ADD (TOTAL_UJI NUMBER)\nSELECT 1 FROM DUAL\n/\n"
	if got := sqlLingkupClaimLife(lain, komentar); len(got) != 0 {
		t.Errorf("komentar ikut diperiksa: %q", got)
	}
	if got := sqlLingkupClaimLife("modul/claimlife/backend/migrations/001.sql", fk); len(got) != 1 || got[0] != fk {
		t.Errorf("berkas Claim Life wajib diperiksa UTUH: %q", got)
	}
}
