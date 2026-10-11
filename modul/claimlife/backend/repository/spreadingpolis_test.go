package repository

import (
	"context"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// 11-10-2026 - spreading peserta dibaca dari PremiumList Life: berkunci IDPEGA (PREMIUM_LIST_ID, IDX_PLD_PL) +
// sertifikat, lalu DETAIL_ID (IDX_PLS_DETAIL) / SPREADING_ID (IDX_PLSR_SPR). Baca saja.
func TestSQLSpreadingPolisPeserta(t *testing.T) {
	d := strings.Join(strings.Fields(sqlDetailPolisPeserta("S.D")), " ")
	for _, mau := range []string{"FROM S.D d", "WHERE d.PREMIUM_LIST_ID = :1 AND TRIM(d.CERTIFICATE_NO) = TRIM(:2)",
		"(d.EDM_STATUS IS NULL OR d.EDM_STATUS <> :3)", "FETCH FIRST 1 ROWS ONLY"} {
		if !strings.Contains(d, mau) {
			t.Errorf("SQL peserta polis tanpa %q:\n%s", mau, d)
		}
	}
	s := strings.Join(strings.Fields(sqlSpreadingPolisPeserta("S.S", "S.R")), " ")
	for _, mau := range []string{"FROM S.S s LEFT JOIN S.R r ON r.SPREADING_ID = s.ID", "WHERE s.DETAIL_ID = :1",
		"ORDER BY s.ID, r.ID", "r.REINSURER_NAME", "PREMIUM_SPREADED_GROSS", "PREMIUM_SPREADED_NET", "r.COMMISION", "r.OVR_COMM"} {
		if !strings.Contains(s, mau) {
			t.Errorf("SQL spreading polis tanpa %q:\n%s", mau, s)
		}
	}
	for _, q := range []string{d, s} {
		if err := db.PeriksaSQL(q); err != nil {
			t.Error(err)
		}
		u := strings.ToUpper(q)
		if strings.Contains(u, "INSERT") || strings.Contains(u, "UPDATE") || strings.Contains(u, "DELETE") {
			t.Errorf("pembaca spreading polis menulis:\n%s", q)
		}
	}
	if edmStatusHapus != "Delete" {
		t.Errorf("status EDM dicabut = %q, mau Delete (endorsementlife)", edmStatusHapus)
	}
}

// Tanpa IDPEGA atau sertifikat: nol lapisan TANPA menyentuh basis data.
func TestSpreadingPolisTanpaKunci(t *testing.T) {
	for _, k := range [][2]string{{"", "001"}, {"NBLF-UJI-1", ""}} {
		got, err := (&KlaimLife{}).SpreadingPolisPeserta(context.Background(), k[0], k[1])
		if err != nil || got != nil {
			t.Errorf("%v: %v, %v", k, got, err)
		}
	}
}
