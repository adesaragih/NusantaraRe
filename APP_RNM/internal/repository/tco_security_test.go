package repository

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// AC 18-20: kolom bernama, kunci surrogate, nol pemangkas spasi.
func TestSQLSecurityTCO(t *testing.T) {
	sisip := sqlSisipSecurityTCO("S.T")
	if !strings.Contains(sisip, "("+strings.Join(kolomSecurityTCO, ", ")+")") {
		t.Errorf("sisip tanpa daftar kolom lengkap: %s", sisip)
	}
	if strings.Count(sisip, "NULL") != 3 {
		t.Errorf("TOP_ID, TP_TREATY, USER_ID harus NULL eksplisit: %s", sisip)
	}
	perbarui, hapus := sqlPerbaruiSecurityTCO("S.T"), sqlHapusSecurityTCO("S.T")
	if !strings.Contains(perbarui, "WHERE ID = :4 AND REAS_ID = :5") || !strings.Contains(hapus, "WHERE ID = :1 AND REAS_ID = :2") {
		t.Errorf("kunci bukan ID surrogate + REAS_ID:\n%s\n%s", perbarui, hapus)
	}
	for nama, q := range map[string]string{"perbarui": perbarui, "hapus": hapus} {
		if strings.Contains(q[strings.Index(q, "WHERE"):], "REAS_SECURITY") {
			t.Errorf("%s memakai nama security sebagai kunci: %s", nama, q)
		}
	}
	semua := []string{sisip, perbarui, hapus, sqlCariDobelSecurityTCO("S.T"),
		sqlDaftarSecurityTCO("S.T", "S.A"), sqlAmbilSecurityTCO("S.T", "S.A")}
	for _, q := range semua {
		if strings.Contains(strings.ToUpper(q), "TRIM(") {
			t.Errorf("SQL security memakai TRIM: %s", q)
		}
		if err := PeriksaSQL(q); err != nil {
			t.Errorf("PeriksaSQL: %v", err)
		}
	}
	if !strings.Contains(sqlDaftarSecurityTCO("S.T", "S.A"), "WHERE s.REAS_ID = :1 AND s.THN_TREATY = :2") {
		t.Error("daftar tidak menyaring A AND D SelectSecurityReinsurer")
	}
}

// Daftar kolom sisip = kolom DDL tiket 01, satu per satu.
func TestKolomSecuritySamaDenganDDL(t *testing.T) {
	isi, err := os.ReadFile("migrations/303_t_mtreatysecurity.sql")
	if err != nil {
		t.Fatal(err)
	}
	pernyataan := strings.Split(string(isi), "\n/\n")[0]
	_, kolom := KolomCreateTable(pernyataan)
	if strings.Join(kolom, ",") != strings.Join(kolomSecurityTCO, ",") {
		t.Errorf("DDL %v, sisip %v", kolom, kolomSecurityTCO)
	}
}

// AC 19: tidak ada penulisan berposisi di sumber modul mana pun - INSERT yang
// langsung disusul VALUES (tanpa daftar kolom) menggagalkan uji ini.
func TestTCONolInsertPosisional(t *testing.T) {
	pola := regexp.MustCompile(`(?is)INSERT\s+INTO\s+(%s|\S+)\s+VALUES`)
	diperiksa := 0
	for nama, isi := range berkasSumberProduksi(t) {
		if !strings.HasSuffix(nama, ".go") || !strings.Contains(nama, "/tco_") {
			continue
		}
		kode := buangKomentarSumber(nama, isi)
		if strings.Contains(strings.ToUpper(kode), "INSERT INTO") {
			diperiksa++
		}
		if m := pola.FindString(kode); m != "" {
			t.Errorf("%s: INSERT tanpa daftar kolom: %q", nama, m)
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol berkas penulis terbaca; pembacanya yang rusak")
	}
	if pola.FindString("INSERT INTO %s VALUES (:1)") == "" {
		t.Fatal("pola penjaga tidak menggigit")
	}
}

func TestPindaiSecurityTCO(t *testing.T) {
	baris := barisPalsu{nilai: []any{"1000001", "2026", "", "", "1000007", "12.5", "", "UJI-AG1", "UJI SECURITY"}}
	s, err := pindaiSecurityTCO(baris)
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != "1000001" || s.ReasID != "1000007" || s.ReasSecurity != "UJI-AG1" || s.ClientName != "UJI SECURITY" ||
		s.PctShare.Text('f') != "12.5" {
		t.Errorf("pindai: %+v", s)
	}
}
