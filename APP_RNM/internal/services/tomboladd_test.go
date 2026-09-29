package services_test

// Tombol `Add` grid adjustment - `ClaimLifeDetailGCNM.xml` b17937. TANPA Oracle.
//
// ⛔ Sejak GILIRAN-14 butir bp `Add` = jalur putaran SAJA; baris pertama lahir
// saat Submit Register (`barispendaftaran_test.go`).

import (
	"os"
	"strings"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/internal/services"
)

// TestPemegangClaimAnalisSamaDenganPeranPutaran - dua gerbang, satu peran.
//
// Layar menampilkan `Add` hanya di tahap Claim Analis
// (`pyWorkPage.pyPosition =='ReasLifeSPV'` b18160); layanan putaran menuntut
// `PeranSimpanOutstanding`. Bila kelak keduanya berbeda, tombolnya tampil
// kepada peran yang tidak dilayani.
func TestPemegangClaimAnalisSamaDenganPeranPutaran(t *testing.T) {
	peran, ada := models.PeranPemegangTahap(models.TahapClaimAnalis)
	if !ada || peran != services.PeranSimpanOutstanding {
		t.Errorf("pemegang Claim Analis %q (%v), peran putaran %q - satu tombol, dua peran",
			peran, ada, services.PeranSimpanOutstanding)
	}
}

// TestTombolAddAdjustmentVERBATIMDariKorpus - bentuk `Add` dan `Delete`.
//
// ⚠️ Melewati bila korpus tidak terjangkau - dan mengatakannya.
//
// ⚠️ b17991 (activity pewaris delapan kolom) SENGAJA tidak diperiksa di sini:
// namanya memuat pola indeks posisi yang dilarang penjaga AC 62
// (`komite_statik_test.go`) di luar komentar. Ia dikutip di komentar
// `BarisPertama`, bukan dikelabui dengan memecah teksnya.
func TestTombolAddAdjustmentVERBATIMDariKorpus(t *testing.T) {
	const letak = `D:\XML\RNM_BRD\Claim Life\Section\ClaimLifeDetailGCNM.xml`
	isi, err := os.ReadFile(letak)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v); tombol tidak terperiksa", err)
	}
	baris := strings.Split(string(isi), "\n")
	for nomor, mau := range map[int]string{
		17126: `<pyPageListProperty>.AdjustmentList</pyPageListProperty>`,
		17937: `<pyLabel>Add</pyLabel>`,
		17947: `<pyAction>addRow</pyAction>`,
		18160: `<pyCondition>pyWorkPage.pyPosition =='ReasLifeSPV'</pyCondition>`,
		19120: `<pyLabel>Delete</pyLabel>`,
		19130: `<pyAction>deleteRow</pyAction>`,
		19399: `<pyCondition>.PrintFaceClaim == '' &amp;&amp; pyWorkPage.pyPosition =='ReasLifeSPV'</pyCondition>`,
	} {
		if got := strings.TrimSpace(baris[nomor-1]); got != mau {
			t.Errorf("b%d = %q, mau %q", nomor, got, mau)
		}
	}
}
