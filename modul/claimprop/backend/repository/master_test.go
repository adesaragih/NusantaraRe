package repository

// Perakit SQL popup "Data Master TreatyIn": saringan XML PROPORTIONTYPE = 'Proportional', filter per kolom (AND,
// mengandung, tanpa beda huruf), urutan terbaru dulu, batas 500 (keputusan work owner 08-10-2026).

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/modul/claimprop/backend/models"
)

func TestSqlDaftarMasterTanpaFilter(t *testing.T) {
	q, args := sqlDaftarMaster("S.CLAIM_MASTER_TREATY", models.SaringanMaster{})
	for _, w := range []string{
		"FROM S.CLAIM_MASTER_TREATY",
		"WHERE PROPORTIONTYPE = :1",
		"ORDER BY TREATYYEAR DESC, TREATYID DESC",
		"FETCH FIRST 500 ROWS ONLY",
	} {
		if !strings.Contains(q, w) {
			t.Fatalf("SQL tanpa %q:\n%s", w, q)
		}
	}
	if strings.Contains(q, "LIKE") {
		t.Fatalf("tanpa filter tidak boleh ada LIKE:\n%s", q)
	}
	if !reflect.DeepEqual(args, []any{"Proportional"}) {
		t.Fatalf("argumen %v", args)
	}
}

func TestSqlDaftarMasterFilterPerKolom(t *testing.T) {
	q, args := sqlDaftarMaster("S.V", models.SaringanMaster{TreatyID: " 1002 ", ContractName: "quota", TreatyYear: "2025"})
	for _, w := range []string{
		"PROPORTIONTYPE = :1",
		"AND UPPER(TREATYID) LIKE :2",
		"AND UPPER(TREATYCONTRACTNAME) LIKE :3",
		"AND UPPER(TREATYYEAR) LIKE :4",
	} {
		if !strings.Contains(q, w) {
			t.Fatalf("SQL tanpa %q:\n%s", w, q)
		}
	}
	if strings.Contains(q, " OR ") {
		t.Fatalf("filter per kolom digabung AND, bukan OR:\n%s", q)
	}
	if !reflect.DeepEqual(args, []any{"Proportional", "%1002%", "%QUOTA%", "%2025%"}) {
		t.Fatalf("argumen %v", args)
	}
}

// Tombol View (keputusan work owner 08-10-2026): berkas NB / EDM Treaty In terbaru untuk satu nomor polis - generasi
// PRODKE terbesar yang benar-benar ada di T_WORK_POLIS.
func TestSqlBerkasPolis(t *testing.T) {
	q := sqlBerkasPolis("S.T_GENERAL_POLIS_TREATY", "S.T_WORK_POLIS")
	for _, w := range []string{
		"SELECT g.ID, g.PRODKE FROM S.T_GENERAL_POLIS_TREATY g JOIN S.T_WORK_POLIS w ON w.ID = g.ID",
		"WHERE g.NOPOLIS = :1",
		"ORDER BY g.PRODKE DESC",
		"FETCH FIRST 1 ROWS ONLY",
	} {
		if !strings.Contains(q, w) {
			t.Fatalf("SQL tanpa %q:\n%s", w, q)
		}
	}
}

// Reporter Address dari CLIENT_ADDRESS, bukan JSON M_CLIENT (work owner 08-10-2026 "ubah jangan dari json, ambil dari
// client address"): klien milik agen, baris Kantor (tipe 2) dulu, Email (tipe 7) tidak pernah, urutan stabil
// PXCREATEDATETIME dengan ROWID hanya pemutus seri terakhir.
func TestSqlAlamatKlienDariClientAddress(t *testing.T) {
	q := sqlAlamatKlien("S.CLIENT_ADDRESS", "S.AGENT")
	for _, w := range []string{
		"SELECT ASMADDRESS, RWNAME, DISTRICTNAME, CITYNAME FROM",
		"FROM S.CLIENT_ADDRESS a",
		"a.CLIENTID = (SELECT MAX(CLIENTID) FROM S.AGENT WHERE ID = :1)",
		"NVL(a.ASMADDRESSTYPE, '-') <> '7'",
		"ORDER BY CASE WHEN a.ASMADDRESSTYPE = '2' THEN 0 ELSE 1 END, a.PXCREATEDATETIME NULLS LAST, a.ROWID",
		"WHERE URUT = 1",
	} {
		if !strings.Contains(q, w) {
			t.Fatalf("SQL tanpa %q:\n%s", w, q)
		}
	}
	if strings.Contains(strings.ToUpper(q), "JSON") || strings.Contains(q, "M_CLIENT") {
		t.Fatalf("Reporter Address tidak boleh dari JSON M_CLIENT:\n%s", q)
	}
}
