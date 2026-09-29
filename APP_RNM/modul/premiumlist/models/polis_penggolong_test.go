package models_test

// Decision3 dirutekan dari bendera - GILIRAN-14 butir bq. TANPA Oracle.

import (
	"errors"
	"os"
	"strings"
	"testing"

	"nusantarare/modul/premiumlist/models"
)

// TestPenggolongOtomatisMenurutDecisionTable - `IsFlagOnGoingPolicy` utuh.
//
// `[terverifikasi]` `DecisionTable/IsFlagOnGoingPolicy.xml`: kolom
// `pyWorkPage.FlagOnGoingPolicy` operator `=`, baris "0" -> Offer (b293 ->
// b328) dan "1" -> Premium (b294 -> b329); selain itu `Decline`
// (`pyDefaultResult` b94). `Decision3` flow `InputPolicyHolder.xml` hanya punya
// DUA konektor: Premium -> ASSIGNMENT63 (b1643/b1658), Offer -> END52
// (b1793/b1807).
func TestPenggolongOtomatisMenurutDecisionTable(t *testing.T) {
	offer, err := models.PenggolongOtomatis(models.FlagPolisPenawaran)
	if err != nil {
		t.Fatal(err)
	}
	if offer.StatusWork != models.StatusPolisSelesai || offer.TahapTujuan != "" || offer.SimpanPolis {
		t.Errorf(`bendera "0" = %+v, mau tutup Resolved-Completed tanpa simpan polis`, offer)
	}
	premium, err := models.PenggolongOtomatis(models.FlagPolisPremium)
	if err != nil {
		t.Fatal(err)
	}
	if premium.TahapTujuan != models.TahapPolisDetail || premium.Ditutup() {
		t.Errorf(`bendera "1" = %+v, mau pindah ke %q`, premium, models.TahapPolisDetail)
	}
	if offer.KeDecision3 || premium.KeDecision3 {
		t.Error("hasil penggolong otomatis masih menunggu penggolong")
	}
}

// TestBenderaDiLuarTabelTidakPunyaKonektor - `Decline` tanpa jalan keluar.
//
// ⛔ Hasil `otherwise` decision table adalah `Decline`, dan `Decision3` NOL
// konektor `Decline`. Di Pega itu masalah flow; di sini galat terang - bukan
// jalur yang dikarang. Bendera kosong adalah keadaan nyata baris yang lahir
// sebelum 057 (ADR-U-0027).
func TestBenderaDiLuarTabelTidakPunyaKonektor(t *testing.T) {
	for _, flag := range []string{"", "2", " 0", "00", "true"} {
		_, err := models.PenggolongOtomatis(flag)
		if !errors.Is(err, models.ErrBenderaTanpaKonektor) {
			t.Errorf("bendera %q: galat = %v, mau ErrBenderaTanpaKonektor", flag, err)
		}
		if got := models.HasilIsFlagOnGoingPolicy(flag); got != models.HasilOtherwiseIsFlagOnGoingPolicy {
			t.Errorf("bendera %q: hasil tabel %q, mau %q (otherwise b94)",
				flag, got, models.HasilOtherwiseIsFlagOnGoingPolicy)
		}
	}
}

// TestDecisionTableDanDecision3VERBATIMDariKorpus - bentuknya dari XML.
//
// ⚠️ Melewati bila korpus tidak terjangkau - dan mengatakannya. Berkasnya
// dipecah `><` -> `>\n<` seperti bNNN dikutip.
func TestDecisionTableDanDecision3VERBATIMDariKorpus(t *testing.T) {
	const akar = `D:\XML\RNM_BRD\PremiumList Life`
	baca := func(relatif string) []string {
		isi, err := os.ReadFile(akar + `\` + relatif)
		if err != nil {
			t.Skipf("korpus tidak terjangkau di mesin ini (%v); bentuk tidak terperiksa", err)
		}
		return strings.Split(strings.ReplaceAll(string(isi), "><", ">\n<"), "\n")
	}
	tabel := baca(`DecisionTable\IsFlagOnGoingPolicy.xml`)
	for nomor, mau := range map[int]string{
		94:  `<pyDefaultResult>Decline</pyDefaultResult>`,
		293: `<rowdata REPEATINGINDEX="1">0</rowdata>`,
		294: `<rowdata REPEATINGINDEX="2">1</rowdata>`,
		328: `<rowdata REPEATINGINDEX="1">Offer</rowdata>`,
		329: `<rowdata REPEATINGINDEX="2">Premium</rowdata>`,
	} {
		if got := strings.TrimSpace(tabel[nomor-1]); got != mau {
			t.Errorf("IsFlagOnGoingPolicy b%d = %q, mau %q", nomor, got, mau)
		}
	}
	flow := strings.Join(baca(`InputPolicyHolder.xml`), "\n")
	if n := strings.Count(flow, "<pyFrom>Decision3</pyFrom>"); n != 2 {
		t.Errorf("Decision3 punya %d konektor keluar, mau 2 (Premium, Offer)", n)
	}
}
