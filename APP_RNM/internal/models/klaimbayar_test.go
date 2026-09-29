package models

import (
	"errors"
	"os"
	"strings"
	"testing"

	"nusantarare/inti/uang"
	"nusantarare/inti/utils"
)

func uangUji(t *testing.T, s string) uang.Money {
	t.Helper()
	m, err := uang.NewMoney(s, "IDR")
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func rasioUji(t *testing.T, s string) uang.Ratio {
	t.Helper()
	r, err := uang.NewRatio(s, 5)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

// TestClaimPaidPersenDibagiSeratusLimaDesimal - langkah 3 b702/b723.
func TestClaimPaidPersenDibagiSeratusLimaDesimal(t *testing.T) {
	for _, u := range []struct{ gross, persen, mau string }{
		{"1000000", "50", "500000"},
		{"1000000", "100", "1000000"},
		// ⛔ Pembulatan LIMA desimal pada PECAHANnya, bukan pada hasilnya:
		// 33.333333% → 0.33333, lalu dikalikan utuh.
		{"1000000", "33.333333", "333330"},
		{"123.45", "12.5", "15.43125"},
	} {
		got, err := HitungClaimPaid(uangUji(t, u.gross), rasioUji(t, u.persen))
		if err != nil {
			t.Fatalf("%s × %s%%: %v", u.gross, u.persen, err)
		}
		if s := utils.FormatDecimal(got.Amount); s != u.mau {
			t.Errorf("%s × %s%% = %s, mau %s", u.gross, u.persen, s, u.mau)
		}
		if got.Currency != "IDR" {
			t.Errorf("mata uang hilang: %q", got.Currency)
		}
	}
}

// TestClaimPaidKosongBernilaiNol - `@toDecimal("")` Pega bernilai nol.
func TestClaimPaidKosongBernilaiNol(t *testing.T) {
	got, err := HitungClaimPaid(uang.Money{Currency: "IDR"}, uang.Ratio{})
	if err != nil {
		t.Fatal(err)
	}
	if s := utils.FormatDecimal(got.Amount); s != "0" {
		t.Errorf("= %s, mau 0", s)
	}
}

// TestGrossMelebihiShareDitolakKecualiEmpatKode - langkah 4 b949 + b972.
func TestGrossMelebihiShareDitolakKecualiEmpatKode(t *testing.T) {
	besar, kecil := uangUji(t, "200"), uangUji(t, "100")
	err := PeriksaClaimGrossTerhadapShare(besar, kecil, "L01")
	if !errors.Is(err, ErrClaimGrossMelebihiShare) {
		t.Fatalf("gross > share: %v, mau ErrClaimGrossMelebihiShare", err)
	}
	if err.Error() != "Claim gross tidak boleh lebih besar dari Share Nusantara Re" {
		t.Errorf("kalimat b526 berubah: %q", err.Error())
	}
	for _, kode := range []string{"L12", "L13", "L14", "L18"} {
		if err := PeriksaClaimGrossTerhadapShare(besar, kecil, kode); err != nil {
			t.Errorf("%s dikecualikan b972, tetapi ditolak: %v", kode, err)
		}
	}
	// Sama besar BUKAN pelanggaran: `>` ketat.
	if err := PeriksaClaimGrossTerhadapShare(kecil, kecil, "L01"); err != nil {
		t.Errorf("gross == share ditolak: %v", err)
	}
	if err := PeriksaClaimGrossTerhadapShare(kecil, besar, "L01"); err != nil {
		t.Errorf("gross < share ditolak: %v", err)
	}
	// ⛔ Mata uang berbeda tidak pernah dibandingkan - bahkan untuk L12.
	usd, err := uang.NewMoney("50", "USD")
	if err != nil {
		t.Fatal(err)
	}
	for _, kode := range []string{"L01", "L12"} {
		if err := PeriksaClaimGrossTerhadapShare(usd, besar, kode); !errors.Is(err, uang.ErrMataUangBerbeda) {
			t.Errorf("%s: USD lawan IDR = %v, mau ErrMataUangBerbeda", kode, err)
		}
	}
}

// TestRumusClaimPaidVERBATIMDariKorpus - rumus uang dibaca langsung dari XML.
func TestRumusClaimPaidVERBATIMDariKorpus(t *testing.T) {
	isi, err := os.ReadFile(`D:\XML\RNM_BRD\Claim Life\Activity\CountClaimAmountLife_Act.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v)", err)
	}
	teks := string(isi)
	for _, potong := range []string{
		`<PropertiesValue>@divide(@toDecimal(.PCTClaim),100,5)</PropertiesValue>`,
		`<PropertiesValue>local.ClaimGross * local.PCTClaim</PropertiesValue>`,
		`<pyStepsPreCondParamsWhen>local.ClaimGross&gt;local.ShareRNM</pyStepsPreCondParamsWhen>`,
		`"Claim gross tidak boleh lebih besar dari Share Nusantara Re"`,
		`<pyStepsTransParamsWhenTrue>6</pyStepsTransParamsWhenTrue>`,
	} {
		if !strings.Contains(teks, potong) {
			t.Errorf("korpus tidak lagi memuat %s; rumus Claim Paid harus dibaca ulang", potong)
		}
	}
}
