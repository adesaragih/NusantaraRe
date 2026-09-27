package repository

// Pembaca roster Komite - TANPA Oracle.
//
// Pemilik: A2 (butir af).

import (
	"strings"
	"testing"
)

// TestSQLRosterMeniruFilterReportDefinition.
//
// `[terverifikasi]` `Claim Life/ReportDefinition/FilterEmailKomiteWithLimit.xml`
// pecahan baris 662 `pyFilterLogic = A AND C AND B`, dengan
// A (670-680) `.LIMIT_BOTTOM <= Param.LIMIT_BOTTOM`,
// C (683-697) `.STS_KLAIM = Param.STS_KLAIM`,
// B (701-715) `.STS_AKTIF = "1"`,
// dan pengurutan `.DEGREE` `ASC` (baris 824-831, `pySortOrder 1`).
func TestSQLRosterMeniruFilterReportDefinition(t *testing.T) {
	q := sqlRosterKomite("UJI.EMAILKOMITE")
	if err := PeriksaSQL(q); err != nil {
		t.Fatalf("pernyataan roster tidak lolos PeriksaSQL: %v", err)
	}
	for _, bagian := range []string{
		"LIMIT_BOTTOM <=", "STS_KLAIM =", "STS_AKTIF =", "ORDER BY", "DEGREE",
	} {
		if !strings.Contains(q, bagian) {
			t.Errorf("pernyataan roster tidak memuat %q: %s", bagian, q)
		}
	}
	// ⛔ Ketiga filter XML, tidak lebih dan tidak kurang: menambah filter
	// berarti mengarang aturan, menguranginya berarti memanggil orang yang
	// tidak berhak memutuskan.
	if n := strings.Count(strings.ToUpper(q), " AND "); n != 2 {
		t.Errorf("filter disambung %d kali AND, mau 2 (tiga filter): %s", n, q)
	}
	// Urutan menaik - tangga komite berjenjang dari yang terendah.
	if !strings.Contains(strings.ToUpper(q), "ORDER BY DEGREE ASC") {
		t.Errorf("urutan bukan DEGREE ASC: %s", q)
	}
	if !strings.Contains(q, "UJI.EMAILKOMITE") {
		t.Errorf("nama tabel tidak berskema: %s", q)
	}
	if strings.Contains(strings.ToUpper(q), "COMMIT") {
		t.Error("pernyataan roster memuat COMMIT")
	}
}

// TestRosterHanyaMengambilKolomYangDipakai - pagar data orang.
//
// ⚠️ `EMAILKOMITE` memuat 18 kolom, sebagian **data orang** (`NAME`, `EMAIL`,
// `JABATAN`, `OPERATOR_ID`). Yang diambil hanya yang benar-benar dipakai
// membangun tangga; `NAME` dan `JABATAN` TIDAK ikut, sebab tidak ada yang
// membutuhkannya dan kolom yang dibaca cenderung ikut tercatat di log.
func TestRosterHanyaMengambilKolomYangDipakai(t *testing.T) {
	q := sqlRosterKomite("UJI.EMAILKOMITE")
	// ⛔ RALAT: ronde pertama menuntut `ID` diambil dan melarang `OPERATOR_ID`
	// serta `JABATAN` - dugaan saya, bukan bacaan. XML menang:
	// `[terverifikasi]` `CreateKMTLife_Act.xml` 866/972 `.KomiteID =
	// .OPERATOR_ID` dan 952/1041 `.IDKomite = .JABATAN`. Nama kolomnya
	// menipu ke DUA arah sekaligus.
	for _, wajib := range []string{"OPERATOR_ID", "JABATAN", "EMAIL", "DEGREE", "LIMIT_BOTTOM"} {
		if !strings.Contains(q, wajib) {
			t.Errorf("kolom %q yang dipakai tangga tidak diambil: %s", wajib, q)
		}
	}
	// `NAME` memang tidak dipakai satu pun langkah XML, dan ia data orang.
	if strings.Contains(q, "NAME") {
		t.Errorf("kolom NAME ikut diambil padahal nol langkah XML memakainya; "+
			"kolom yang dibaca cenderung ikut tercatat: %s", q)
	}
}
