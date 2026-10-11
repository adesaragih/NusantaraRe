package repository

// Find Insured membawa DOB (keputusan work owner 10-10-2026 "policy no ini ganti saja dengan dob") - teks ISO lewat
// TO_CHAR di daftar pilih luar (kebal NLS), kolom mentah DOB di subkueri peringkat; EDMSTATUS tetap di ekor.

import (
	"strings"
	"testing"
)

func TestCariPesertaMembawaDOBSebagaiTeksISO(t *testing.T) {
	q, _ := sqlCariPeserta("S.M", "UJI-PL", "", "", 10)
	s := strings.Join(strings.Fields(q), " ")
	if !strings.HasPrefix(s, "SELECT PL_NUMBER, POLICY_NO, CERTIFICATE_NO, NAME_OF_INSURED, CURRENCY, IDPEGA, TO_CHAR(DOB, 'YYYY-MM-DD'), EDMSTATUS FROM (") {
		t.Errorf("daftar pilih luar:\n%s", s)
	}
	if !strings.Contains(s, "FROM (SELECT PL_NUMBER, POLICY_NO, CERTIFICATE_NO, NAME_OF_INSURED, CURRENCY, IDPEGA, DOB, EDMSTATUS, ") {
		t.Errorf("subkueri tanpa kolom mentah DOB:\n%s", s)
	}
}
