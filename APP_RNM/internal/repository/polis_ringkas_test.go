package repository

// Ringkasan polis untuk Claim Life - butir pl4. TANPA Oracle.

import (
	"strings"
	"testing"
)

// TestRingkasMembacaVersiBerjalan - `PROD_KE` TERBESAR, bukan sembarang baris.
//
// ⛔ Migrasi 051: satu baris = satu VERSI polis, dan seluruh versi hidup
// berdampingan. Tanpa urutan yang dinyatakan, Oracle bebas memberi versi mana
// pun yang paling murah dibacanya - dan yang paling murah bukan yang paling
// baru. Klaim yang dinilai dengan data polis versi lama dinilai salah, tanpa
// satu pun galat.
func TestRingkasMembacaVersiBerjalan(t *testing.T) {
	q := sqlPolisRingkas("SKEMAUJI.T_PREMIUM_LIST")
	if !strings.Contains(q, "ORDER BY NVL(p.PROD_KE, 0) DESC") {
		t.Errorf("query tidak mengurutkan versi menurun:\n%s", q)
	}
	if !strings.Contains(q, "FETCH FIRST 1 ROWS ONLY") {
		t.Errorf("query tidak berbatas satu baris:\n%s", q)
	}
	// ⚠️ `NVL` supaya baris new business yang PROD_KE-nya masih kosong tetap
	// ikut terurut, bukan terlempar ke ujung oleh NULL.
	if !strings.Contains(q, "NVL(p.PROD_KE, 0)") {
		t.Errorf("PROD_KE kosong tidak dinetralkan:\n%s", q)
	}
	if err := PeriksaSQL(q); err != nil {
		t.Errorf("%v\n%s", err, q)
	}
}

// TestRingkasHanyaKolomYangDipakai - bukan `SELECT *`.
func TestRingkasHanyaKolomYangDipakai(t *testing.T) {
	q := sqlPolisRingkas("SKEMAUJI.T_PREMIUM_LIST")
	if strings.Contains(q, "SELECT *") {
		t.Errorf("query membawa seluruh kolom:\n%s", q)
	}
	// Kesepuluh medan yang PUNYA kolom, plus versinya.
	for _, kolom := range []string{
		"p.NO_POLIS", "p.TYPE", "p.MARKETING_NAME", "p.CEDING_CO_NAME",
		"p.POLICY_HOLDER_NAME", "p.BUSINESS_NAME", "p.DATE_RECEIVED",
		"p.STATUSS", "p.STATUS_UPDATE", "p.PRODUCT_NAME_ID", "p.PRODUCT_NAME",
	} {
		if !strings.Contains(q, kolom) {
			t.Errorf("kolom %q tidak dibaca:\n%s", kolom, q)
		}
	}
}

// TestKolomRingkasAdaDiMigrasi051 - nol kolom hantu.
func TestKolomRingkasAdaDiMigrasi051(t *testing.T) {
	isi, err := berkasMigrasi.ReadFile("migrations/051_t_premium_list.sql")
	if err != nil {
		t.Fatalf("membaca migrasi 051: %v", err)
	}
	ada := map[string]bool{}
	for _, pernyataan := range strings.Split(string(isi), "\n/") {
		nama, kolom := KolomCreateTable(pernyataan)
		if !strings.HasSuffix(strings.ToUpper(nama), "T_PREMIUM_LIST") {
			continue
		}
		for _, k := range kolom {
			ada[strings.ToUpper(k)] = true
		}
	}
	if len(ada) == 0 {
		t.Fatal("nol kolom terbaca dari migrasi 051; pembacanya yang rusak")
	}
	for _, k := range []string{
		"NO_POLIS", "TYPE", "MARKETING_NAME", "CEDING_CO_NAME",
		"POLICY_HOLDER_NAME", "BUSINESS_NAME", "DATE_RECEIVED", "STATUSS",
		"STATUS_UPDATE", "PRODUCT_NAME_ID", "PRODUCT_NAME", "PROD_KE",
	} {
		if !ada[k] {
			t.Errorf("kolom %q tidak ada di migrasi 051", k)
		}
	}
	// ⛔ KETIGA MEDAN YANG DINYATAKAN TANPA SUMBER memang TIDAK ada. Bila
	// salah satunya kelak ditambahkan, penjaga ini berbunyi - dan catatan di
	// layar harus ikut berubah, bukan diam-diam menjadi bohong.
	for _, k := range []string{
		"TANGGAL_RESPON", "TANGGAL_KONFIRMASI", "TANGGAL_REALISASI",
	} {
		if ada[k] {
			t.Errorf("kolom %q ternyata ADA di migrasi 051 - catatan "+
				"`medanTanpaSumber` di services dan layar harus diperbarui", k)
		}
	}
}

// TestTanggalRingkasBerpolaDinyatakan - bukan dari NLS sesi.
func TestTanggalRingkasBerpolaDinyatakan(t *testing.T) {
	q := sqlPolisRingkas("SKEMAUJI.T_PREMIUM_LIST")
	if !strings.Contains(q, "TO_CHAR(p.DATE_RECEIVED, 'YYYY-MM-DD')") {
		t.Errorf("DATE_RECEIVED tidak dibungkus TO_CHAR berpola:\n%s", q)
	}
}
