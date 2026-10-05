//go:build db

package repository_test

// Uji seam repository lawan Oracle SUNGGUHAN untuk pembaca master XOL (K8, F1):
// SQL `BrowseTreatyIn` / `BrowseTreatyInJoinEDM` membaca `M_TREATY_IN` lalu
// `M_TREATY_IN_EDM` berkunci nomor kontrak, HANYA medan daftar tertutup
// `models.MedanMasterXOL` yang keluar, dan nol baris = `ErrMasterXOLTidakAda`.
// ⛔ Belum pernah dijalankan (K11 kosong - skema uji belum ada).
//
// Kedua tabel master adalah tabel WARISAN `POOLDATA` (bukan milik modul ini,
// tidak ada di migrasi). Bila skema uji tidak memuatnya, keduanya dibuat
// sebagai TIRUAN lewat `buatTiruan` (keputusan WO U1: berpagar POOLDATA,
// dibuang sesudah uji); tabel yang disediakan DBA dipakai apa adanya dan
// hanya baris UJI- yang dihapus. Fixture berawalan UJI-.

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

const (
	// tiruan uji, bukan tabel aplikasi (bila DBA tidak menyediakannya di skema
	// uji): master treaty in warisan `POOLDATA.M_TREATY_IN`.
	tabelMasterTreatyUji = "M_TREATY_IN"
	// tiruan uji, bukan tabel aplikasi (bila DBA tidak menyediakannya di skema
	// uji): master treaty in EDM warisan `POOLDATA.M_TREATY_IN_EDM`.
	tabelMasterTreatyEDMUji = "M_TREATY_IN_EDM"
)

// dokumenMasterUji memuat medan daftar (sebagian) beserta medan DI LUAR daftar
// tertutup di setiap tingkat: skalar akar, keluaran `TreatyXOLList` (F1),
// medan polis `FlagPPH`/`TypeTax`, anggota baris lain, daftar bersarang lain.
const dokumenMasterUji = `{"RNMShare": 10, "RIOGR": "UJI-LUAR", "FlagPPH": "1", "TypeTax": "Inclusive",
 "TreatyXOLList": [{"GrossPremi": "1", "ValueList": [{"Layer": "1"}]}],
 "Share": [{"Layer": "1", "UjiLuar": "UJI-LUAR", "GrossPremiumList": [{"Currency": "IDR", "Value": 5, "UjiLuar": "x"}],
   "UjiAnakLuar": [{"A": "UJI-LUAR"}]}],
 "Installment": [{"Currency": "IDR", "InstallmentList": [{"PaymentDate": "20261101", "UjiLuar": "x"}]}],
 "UjiDaftarLuar": [{"A": "UJI-LUAR"}]}`

var subskripUji = regexp.MustCompile(`\(\d+\)`)

func TestMasterXOLDariJSONMembacaKeduaTabelMaster(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	kolom := []string{"ID VARCHAR2(100)", "JSONDATA CLOB"}
	dibuat := map[string]bool{
		tabelMasterTreatyUji:    buatTiruan(t, ctx, sqlDB, skema, tabelMasterTreatyUji, kolom),
		tabelMasterTreatyEDMUji: buatTiruan(t, ctx, sqlDB, skema, tabelMasterTreatyEDMUji, kolom),
	}
	for nama, isi := range map[string]string{tabelMasterTreatyUji: dokumenMasterUji, tabelMasterTreatyEDMUji: `{"EDMState": "2"}`} {
		if !dibuat[nama] { // tabel DBA: hanya baris UJI- yang dibuang
			nama := nama
			t.Cleanup(func() {
				_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s.%s WHERE ID = 'UJI-T1'`, skema, nama))
			})
		}
		if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`INSERT INTO %s.%s (ID, JSONDATA) VALUES ('UJI-T1', :1)`, skema, nama), isi); err != nil {
			t.Fatalf("baris UJI- %s: %v", nama, err)
		}
	}
	g := repository.Baru(d)
	m, err := g.MasterXOLDariJSON(ctx, "UJI-T1")
	if err != nil {
		t.Fatal(err)
	}
	if m.Nilai["RNMShare"] != "10" || m.Nilai["EDMState"] != "2" {
		t.Fatalf("kedua tabel dibaca: %v", m.Nilai)
	}
	if r := m.Daftar["Installment(1).InstallmentList"]; len(r) != 1 || r[0]["PaymentDate"] != "2026-11-01" {
		t.Fatalf("InstallmentList %v", r)
	}
	// F1: tidak satu pun medan di luar daftar tertutup keluar dari pembaca.
	boleh := map[string]bool{}
	for _, j := range models.MedanMasterXOL() {
		boleh[j] = true
		for i := strings.LastIndex(j, "."); i > 0; i = strings.LastIndex(j[:i], ".") {
			boleh[j[:i]] = true
		}
	}
	for k := range m.Nilai {
		if !boleh[k] {
			t.Errorf("skalar %s di luar daftar tertutup (F1)", k)
		}
	}
	for k, bs := range m.Daftar {
		jalur := subskripUji.ReplaceAllString(k, "")
		if !boleh[jalur] {
			t.Errorf("daftar %s di luar daftar tertutup (F1)", k)
		}
		for _, b := range bs {
			for a := range b {
				if !boleh[jalur+"."+a] {
					t.Errorf("medan %s.%s di luar daftar tertutup (F1)", k, a)
				}
			}
		}
	}
	if _, err := g.MasterXOLDariJSON(ctx, "UJI-TIDAK-ADA"); !errors.Is(err, repository.ErrMasterXOLTidakAda) {
		t.Fatalf("kontrak tanpa master: %v", err)
	}
}
