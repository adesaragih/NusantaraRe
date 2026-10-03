package services

// OQ-N13 (GILIRAN-18) - penyetel status cermin di Save to RNM. TANPA Oracle.

import (
	"os"
	"strings"
	"testing"
)

// TestSimpanRNMMenyetelCerminDiTransaksiYangSama - `InsertJsonKlaimLife_sql`
// b176 menulis `'0'` di langkah 22.1.3.1, baris demi baris adjustment yang
// BARU disimpan (`.PrintFaceClaim==""`).
//
// ⛔ Di transaksi yang SAMA dengan penanda baris (22.1.3.2): cermin berstatus
// tanpa barisnya - atau sebaliknya - membuat pemeriksa klaim ganda klaim lain
// melihat keadaan yang tidak pernah ter-commit.
func TestSimpanRNMMenyetelCerminDiTransaksiYangSama(t *testing.T) {
	isi, err := os.ReadFile("simpanrnm.go")
	if err != nil {
		t.Fatal(err)
	}
	badan := tubuhFungsi(string(isi), "Simpan")
	urut := []string{
		"x.svc.DalamTransaksi(ctx",
		"baca.TandaiBarisOutstanding(ctx, tx, b.ID)",
		"baca.SetelCerminOutstanding(ctx, tx, b.ID, caseID)",
		"baca.CerminkanHeader(ctx, tx,",
	}
	lalu := -1
	for _, jejak := range urut {
		i := strings.Index(badan, jejak)
		if i < 0 {
			t.Fatalf("Simpan (Save to RNM) tidak memuat %s", jejak)
		}
		if i < lalu {
			t.Errorf("%s berada sebelum langkah sebelumnya; urutan: %v", jejak, urut)
		}
		lalu = i
	}
}

// TestLangkahCerminSaveRNM - baris baru ditandai DAN dicerminkan; baris yang
// sudah Outstanding hanya dicerminkan (cermin NULL milik klaim yang disimpan
// sebelum N13 pulih saat tombolnya ditekan ulang, OQ-N1); baris berkeputusan
// tidak disentuh.
func TestLangkahCerminSaveRNM(t *testing.T) {
	for kode, mau := range map[string][2]bool{
		"": {true, true}, " ": {true, true}, "0": {false, true}, "1": {false, false}, "2": {false, false},
	} {
		tandai, cermin := langkahCerminSaveRNM(kode)
		if tandai != mau[0] || cermin != mau[1] {
			t.Errorf("kode %q: tandai %v cermin %v, mau %v", kode, tandai, cermin, mau)
		}
	}
}
