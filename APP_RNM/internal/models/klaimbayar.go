package models

// Claim Paid dan batas Claim Gross - sensus Claim Life §3.1 (28-09-2026).
//
// Untuk apa berkas ini: `Activity/CountClaimAmountLife_Act.xml`, dipicu isian
// `Percent Claim (%)` `AdjustmentDetail_Section` saat berubah (b1709). Rumus
// uang, jadi dibaca ulang dari XML langkah demi langkah (pecahan sed):
//
//	1  Page-Clear-Messages
//	2  local.ClaimGross/ShareRNM/ClaimPaid/PCTClaim = 0;
//	   local.errmsg = "Claim gross tidak boleh lebih besar dari Share Nusantara Re" (b526)
//	3  local.ClaimGross = @toDecimal(.CLAIM_GROSS)                         (b634)
//	   local.ShareRNM   = @toDecimal(.SHARE_NUSANTARA_RE)                  (b681)
//	   local.PCTClaim   = @divide(@toDecimal(.PCTClaim),100,5)             (b702)
//	   local.ClaimPaid  = local.ClaimGross * local.PCTClaim                (b723)
//	   .CLAIM_PAID      = local.ClaimPaid                                  (b744)
//	4  Property-Set-Messages local.errmsg, BERPRASYARAT (urutan baris):
//	     local.ClaimGross > local.ShareRNM      WhenTrue 2 / WhenFalse 3  (b949)
//	     BusinessCode L12||L13||L14||L18        WhenTrue 3 / WhenFalse 2  (b972)
//	   lalu transisi `1==1` WhenTrue 6 = KELUAR ACTIVITY (b865)
//	5  call SpreadingClaimLife_Act                                         (b1012)
//
// Jadi: CLAIM_PAID SELALU dihitung (langkah 3 mendahului pemeriksaan), dan
// bila pesan terpasang, activity keluar SEBELUM Spreading.
//
// ⛔ `[terbuka - OQ-M3]` BELUM TERSAMBUNG, dan sebabnya bukan kode: kolom
// `PCT_CLAIM`/`CLAIM_PAID` tidak ada di `T_CLAIMLF_ADJUSTMENT` (migrasi 004),
// dan tidak ada rute yang menyunting baris adjustment. Migrasi `020`+ hanya
// dari keputusan tercatat. Aturannya ditulis di sini supaya keputusan itu
// menemukan rumus yang sudah diuji, bukan menebaknya lagi.
//
// Dibaca sesudah: money.go, businesscode.go.

import (
	"errors"
	"fmt"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/pkg/utils"
)

// ErrClaimGrossMelebihiShare - kalimat b526, VERBATIM.
var ErrClaimGrossMelebihiShare = errors.New(
	"Claim gross tidak boleh lebih besar dari Share Nusantara Re")

// desimalPersenKlaim - argumen ketiga `@divide(…,100,5)` b702.
const desimalPersenKlaim = 5

// kodeBisnisTanpaBatasShare - keempat kode b972 yang MELEWATI pemeriksaan.
var kodeBisnisTanpaBatasShare = map[string]bool{"L12": true, "L13": true, "L14": true, "L18": true}

// nolBila menjawab nol untuk nilai kosong - `@toDecimal("")` Pega bernilai nol.
func nolBila(d *apd.Decimal) *apd.Decimal {
	if d == nil {
		return apd.New(0, 0)
	}
	return d
}

// HitungClaimPaid adalah langkah 3 - MURNI, nol float.
//
// ⛔ Yang dibulatkan PECAHAN persennya (lima desimal), bukan hasilnya:
// `local.ClaimPaid = local.ClaimGross * local.PCTClaim` tanpa `@divide`
// kedua. Membulatkan hasil mengubah angka yang Pega simpan.
func HitungClaimPaid(claimGross Money, persenKlaim Ratio) (Money, error) {
	ctx := utils.DecimalContext()
	pecahan := new(apd.Decimal)
	if _, err := ctx.Quo(pecahan, nolBila(persenKlaim.Value), apd.New(100, 0)); err != nil {
		return Money{}, fmt.Errorf("models: membagi persen klaim: %w", err)
	}
	if _, err := ctx.Quantize(pecahan, pecahan, -desimalPersenKlaim); err != nil {
		return Money{}, fmt.Errorf("models: membulatkan persen klaim: %w", err)
	}
	hasil := new(apd.Decimal)
	if _, err := ctx.Mul(hasil, nolBila(claimGross.Amount), pecahan); err != nil {
		return Money{}, fmt.Errorf("models: mengalikan claim gross: %w", err)
	}
	hasil.Reduce(hasil)
	return Money{Amount: hasil, Currency: claimGross.Currency}, nil
}

// PeriksaClaimGrossTerhadapShare adalah langkah 4 - MURNI.
//
// Galat bila gross > share (ketat, b949) DAN kode bisnis BUKAN salah satu
// dari empat kode b972. Galat berarti activity keluar tanpa Spreading.
//
// ⛔ Dua uang bermata uang berbeda TIDAK dibandingkan (ADR-F-0004, sama
// dengan `Money.Add`): IDR melawan USD bukan "lebih besar", melainkan salah.
func PeriksaClaimGrossTerhadapShare(claimGross, shareNusantaraRe Money, kodeBisnis string) error {
	if claimGross.Currency != shareNusantaraRe.Currency {
		return fmt.Errorf("%w: %q lawan %q", ErrMataUangBerbeda,
			claimGross.Currency, shareNusantaraRe.Currency)
	}
	if nolBila(claimGross.Amount).Cmp(nolBila(shareNusantaraRe.Amount)) <= 0 {
		return nil
	}
	if kodeBisnisTanpaBatasShare[kodeBisnis] {
		return nil
	}
	return ErrClaimGrossMelebihiShare
}
