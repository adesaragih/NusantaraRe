package services_test

import "testing"

// Salinan paket services_test (11-10-2026): lingkup Go dan SQL penjaga Claim Life.
func TestLingkupClaimLifeLuar(t *testing.T) {
	if !goLingkupClaimLife("modul/claimlife/backend/x.go", "package x") {
		t.Error("berkas Claim Life wajib selalu diperiksa")
	}
	if goLingkupClaimLife("modul/claimprop/backend/models/layar.go", `const a = "Komite`+`Aproval"`) {
		t.Error("berkas modul lain tanpa tabel Claim Life ikut diperiksa")
	}
	if !goLingkupClaimLife("modul/komiteuji/backend/r.go", `q := "UPDATE T_GENERAL_CLAIM SET A = 1"`) {
		t.Error("berkas modul lain yang menyebut tabel Claim Life wajib diperiksa")
	}
	lain := "modul/claimuji/backend/migrations/900_uji.sql"
	if got := sqlLingkupClaimLifeLuar(lain, "CREATE TABLE {skema}.T_UJI (APPROVAL VARCHAR2(1),\n"+
		"  CONSTRAINT FK FOREIGN KEY (C) REFERENCES {skema}.T_GENERAL_CLAIM (ID))\n/\n"); len(got) != 0 {
		t.Errorf("FK ke T_GENERAL_CLAIM ikut diperiksa: %q", got)
	}
	if got := sqlLingkupClaimLifeLuar(lain, "ALTER TABLE {skema}.T_WORK_CLAIM ADD (APPROVAL VARCHAR2(1))\n/\n"); len(got) != 1 {
		t.Errorf("ALTER atas T_WORK_CLAIM dari modul lain wajib diperiksa: %q", got)
	}
}
