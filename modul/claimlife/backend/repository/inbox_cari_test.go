package repository

import (
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// Kotak cari inbox Claim Life (keputusan work owner 12-10-2026).

func TestPolaCariInbox(t *testing.T) {
	for _, u := range []struct{ masuk, mau string }{
		{"", ""},
		{"   ", ""},
		{" clmlf-12 ", "%CLMLF-12%"},
		{"uji nama", "%UJI NAMA%"},
		// Karakter khusus LIKE diloloskan (ESCAPE '\'): tidak menjadi wildcard.
		{"50%", `%50\%%`},
		{"a_b", `%A\_B%`},
		{`x\y`, `%X\\Y%`},
	} {
		if got := PolaCariInbox(u.masuk); got != u.mau {
			t.Errorf("PolaCariInbox(%q) = %q, mau %q", u.masuk, got, u.mau)
		}
	}
}

func TestSQLInboxCari(t *testing.T) {
	cari := strings.Join(strings.Fields(sqlInboxCari("S.P", "S.M")), " ")
	for _, mau := range []string{
		"UPPER(w.ID) LIKE :cariKasus ESCAPE '\\'",
		"OR UPPER(h.CLAIM_NO) LIKE :cariKlaim ESCAPE '\\'",
		"FROM S.P pc WHERE pc.CLAIM_ID = w.ID AND pc.STS_HAPUS IS NULL AND UPPER(TRIM(pc.PL_NUMBER)) LIKE :cariPL ESCAPE '\\'",
		"FROM S.P pn WHERE pn.CLAIM_ID = w.ID AND pn.STS_HAPUS IS NULL",
		// Nama dari baris sumber lewat kuncinya - per peserta klaim, bukan memindai tabel sumber.
		"FROM S.M m WHERE m.PL_NUMBER = pn.PL_NUMBER AND m.CERTIFICATE_NO = pn.CERTIFICATE_NO AND m.ID = pn.SOURCE_ID",
		"UPPER(m.NAME_OF_INSURED) LIKE :cariNama ESCAPE '\\'",
	} {
		if !strings.Contains(cari, mau) {
			t.Errorf("saringan cari tanpa %q:\n%s", mau, cari)
		}
	}
	if strings.Count(cari, "NO_UNNEST") != 3 {
		t.Errorf("mau 3 NO_UNNEST (penelusuran dari kasus antrian), dapat %d", strings.Count(cari, "NO_UNNEST"))
	}
	// Empat bind cari; nama tidak pernah di SELECT.
	if strings.Count(cari, ":cari") != 4 {
		t.Errorf("mau 4 bind cari, dapat %d", strings.Count(cari, ":cari"))
	}
	for _, q := range []string{sqlInbox("S.W", "S.H", "S.P", "S.J", sqlInboxCari("S.P", "S.M"), true),
		sqlCacahInbox("S.W", "S.H", sqlInboxCari("S.P", "S.M"), true)} {
		satu := strings.Join(strings.Fields(q), " ")
		if err := db.PeriksaSQL(q); err != nil {
			t.Errorf("PeriksaSQL: %v", err)
		}
		// Saringan cari sesudah saringan tahap/akun, sebelum ORDER BY (urutan bind: tahap, akun, cari, offset).
		iAkun, iCari := strings.Index(satu, ":akun"), strings.Index(satu, ":cariKasus")
		if iAkun < 0 || iCari < iAkun {
			t.Errorf("bind cari harus sesudah :akun:\n%s", satu)
		}
		if i := strings.LastIndex(satu, "ORDER BY"); i >= 0 && i < iCari {
			t.Errorf("saringan cari sesudah ORDER BY:\n%s", satu)
		}
		sel := satu[:strings.Index(satu, " FROM S.W w")]
		if strings.Contains(sel, "NAME_OF_INSURED") {
			t.Errorf("nama tertanggung tidak boleh dikeluarkan:\n%s", sel)
		}
	}
	// Tanpa kata kunci: tanpa saringan cari.
	if strings.Contains(sqlInbox("S.W", "S.H", "S.P", "S.J", "", false), ":cari") {
		t.Error("tanpa kata kunci tetap memasang saringan cari")
	}
}
