//go:build db

package repository_test

// Uji seam repository lawan Oracle SUNGGUHAN untuk pembaca master XOL (K8):
// SQL `BrowseTreatyIn` / `BrowseTreatyInJoinEDM` membaca `M_TREATY_IN` lalu
// `M_TREATY_IN_EDM` berkunci nomor kontrak, menyaring medan K8, dan nol baris =
// `ErrMasterXOLTidakAda`. ⛔ Belum pernah dijalankan (K11 kosong - skema uji
// belum ada).
//
// Kedua tabel master adalah tabel WARISAN `POOLDATA` (bukan milik modul ini,
// tidak ada di migrasi). Di skema uji keduanya dibuat sebagai TIRUAN sementara
// - pola `endorsementlife/edm_db_test.go` (`JSON_POLIS`) - dan dibuang lagi.
// Fixture berawalan UJI-.

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/repository"
)

func TestMasterXOLDariJSONMembacaKeduaTabelMaster(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	for _, tb := range []string{"M_TREATY_IN", "M_TREATY_IN_EDM"} {
		q := `CREATE TABLE ` + skema + `.` + tb + ` (ID VARCHAR2(100), JSONDATA CLOB)`
		if _, err := sqlDB.ExecContext(ctx, q); err != nil && !strings.Contains(err.Error(), "ORA-00955") {
			t.Fatalf("tiruan %s: %v", tb, err)
		}
		tb := tb
		t.Cleanup(func() { _, _ = sqlDB.ExecContext(ctx, `DROP TABLE `+skema+`.`+tb+` PURGE`) })
	}
	if _, err := sqlDB.ExecContext(ctx, `INSERT INTO `+skema+`.M_TREATY_IN (ID, JSONDATA) VALUES ('UJI-T1', :1)`,
		`{"RNMShare": 10, "RIOGR": "UJI-TIDAK", "Installment": [{"Currency": "IDR", "InstallmentList": [{"PaymentDate": "20261101"}]}]}`); err != nil {
		t.Fatal(err)
	}
	if _, err := sqlDB.ExecContext(ctx, `INSERT INTO `+skema+`.M_TREATY_IN_EDM (ID, JSONDATA) VALUES ('UJI-T1', :1)`,
		`{"EDMState": "2"}`); err != nil {
		t.Fatal(err)
	}
	g := repository.Baru(d)
	m, err := g.MasterXOLDariJSON(ctx, "UJI-T1")
	if err != nil {
		t.Fatal(err)
	}
	if m.Nilai["RNMShare"] != "10" || m.Nilai["EDMState"] != "2" {
		t.Fatalf("kedua tabel dibaca: %v", m.Nilai)
	}
	if _, ada := m.Nilai["RIOGR"]; ada {
		t.Fatal("medan di luar K8 terbaca")
	}
	if r := m.Daftar["Installment(1).InstallmentList"]; len(r) != 1 || r[0]["PaymentDate"] != "2026-11-01" {
		t.Fatalf("InstallmentList %v", r)
	}
	if _, err := g.MasterXOLDariJSON(ctx, "UJI-TIDAK-ADA"); !errors.Is(err, repository.ErrMasterXOLTidakAda) {
		t.Fatalf("kontrak tanpa master: %v", err)
	}
}
