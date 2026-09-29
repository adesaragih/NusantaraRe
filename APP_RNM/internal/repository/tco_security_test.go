package repository

import (
	"errors"
	"regexp"
	"strings"
	"testing"
)

// AC 19: kolom bernama. tco4 (RALAT AC 18/20 tiket 06): tabel warisan tanpa
// identitas - kunci (REAS_ID, TRIM(REAS_SECURITY)) PERSIS UpdateMTreatySecurity
// b89 dan DeleteSecurityReinsurer b85.
func TestSQLSecurityTCO(t *testing.T) {
	sisip := sqlSisipSecurityTCO("S.T")
	if !strings.Contains(sisip, "("+strings.Join(kolomSecurityTCO, ", ")+")") {
		t.Errorf("sisip tanpa daftar kolom lengkap: %s", sisip)
	}
	if strings.Count(sisip, "NULL") != 3 {
		t.Errorf("TOP_ID, TP_TREATY, USER_ID harus NULL eksplisit: %s", sisip)
	}
	perbarui, hapus := sqlPerbaruiSecurityTCO("S.T"), sqlHapusSecurityTCO("S.T")
	if !strings.Contains(perbarui, "SET THN_TREATY = :1, PCT_SHARE = :2, REAS_SECURITY = :3") ||
		!strings.Contains(perbarui, "WHERE REAS_ID = :4 AND TRIM(REAS_SECURITY) = TRIM(:5)") {
		t.Errorf("perbarui bukan UpdateMTreatySecurity b85-b89:\n%s", perbarui)
	}
	if !strings.Contains(hapus, "WHERE REAS_ID = :1 AND TRIM(REAS_SECURITY) = TRIM(:2)") {
		t.Errorf("hapus bukan DeleteSecurityReinsurer b85:\n%s", hapus)
	}
	semua := []string{sisip, perbarui, hapus, sqlCariDobelSecurityTCO("S.T"),
		sqlDaftarSecurityTCO("S.T", "S.A"), sqlAmbilSecurityTCO("S.T", "S.A")}
	for _, q := range semua {
		if strings.Contains(q, " ID,") || strings.Contains(q, "s.ID") {
			t.Errorf("SQL security menyebut ID yang tidak ada di MTREATYSECURITY: %s", q)
		}
		if err := PeriksaSQL(q); err != nil {
			t.Errorf("PeriksaSQL: %v", err)
		}
	}
	if !strings.Contains(sqlDaftarSecurityTCO("S.T", "S.A"), "WHERE s.REAS_ID = :1 AND s.THN_TREATY = :2") {
		t.Error("daftar tidak menyaring A AND D SelectSecurityReinsurer")
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
	// CHAR berekor spasi dipangkas; PCT_SHARE teks warisan berkoma diterima.
	baris := barisPalsu{nilai: []any{"2026", "", "  ", "1000007", "12,5", "", "UJI-AG1   ", "UJI SECURITY"}}
	s, err := pindaiSecurityTCO(baris)
	if err != nil {
		t.Fatal(err)
	}
	if s.ID != "UJI-AG1" || s.ReasID != "1000007" || s.ReasSecurity != "UJI-AG1" || s.ClientName != "UJI SECURITY" ||
		s.PctShare.Text('f') != "12.5" || s.TpTreaty != "" {
		t.Errorf("pindai: %+v", s)
	}
}

type hasilUjiSecurity int64

func (h hasilUjiSecurity) LastInsertId() (int64, error) { return 0, nil }
func (h hasilUjiSecurity) RowsAffected() (int64, error) { return int64(h), nil }

// Temuan /code-review lanjutan 3: UPDATE/DELETE berkunci nama mengenai SEMUA
// baris senama seperti Pega - duplikat warisan tetap dapat diubah dan dihapus.
func TestPalingSedikitSatuSecurityTCO(t *testing.T) {
	if err := palingSedikitSatuSecurityTCO(hasilUjiSecurity(0)); !errors.Is(err, ErrSecurityTidakAda) {
		t.Errorf("nol baris: %v", err)
	}
	for _, n := range []int64{1, 2} {
		if err := palingSedikitSatuSecurityTCO(hasilUjiSecurity(n)); err != nil {
			t.Errorf("%d baris: %v", n, err)
		}
	}
}
