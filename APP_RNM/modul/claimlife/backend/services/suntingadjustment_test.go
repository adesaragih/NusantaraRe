package services_test

// Sunting sel baris adjustment - GILIRAN-14 butir br. TANPA Oracle.
//
// ⛔ VONIS XML: NOL sel yang dapat disunting. Butir br memutuskan "HANYA
// kolom yang grid `ClaimLifeDetailGCNM.xml` b17126 tampilkan dapat disunting
// (baca `pyEditOptions`/`pyReadOnly`/`pyCondition` tiap kolom)". Keempat kolom
// grid itu `Read-only`, dan panel yang grid buka (`pyEditingMode expandPane`
// b19566 -> `pyEditAction Adjustment_Detail` b19583 ->
// `AdjustmentDetail_Section`) juga `Read-only` di SETIAP medan datanya. Rute
// sunting karena itu TIDAK dibangun: rute yang menjawab 422 untuk setiap
// kolom adalah kode mati. Uji ini mengunci vonisnya - bila korpus kelak
// berkata lain, uji ini yang merah lebih dulu.
//
// ✅ Keputusan work owner 29-09-2026 (GILIRAN-15): vonis DITERIMA - br dan OQ-N8
// ditutup, ikut XML. Uji ini tetap sebagai penjaganya.
//
// Cara membaca (sensus §4a CLAUDE.md - wadah dinamai, kunci dipasang):
//   - WADAH: grid dibatasi rentang b17126 (`pyPageListProperty
//     .AdjustmentList`) sampai b19583 (`pyEditAction` grid itu); panel adalah
//     seluruh `AdjustmentDetail_Section.xml`.
//   - KUNCI: `pyEditOptions` sebuah sel ditulis SEBELUM `pyValue`-nya, di
//     rentang sesudah `pyValue` sebelumnya; pasangan itu diikat per sel.
//   - Sel TANPA `pyEditOptions` sama sekali tidak dilewati diam-diam: ia harus
//     disebut di pengecualian bernama, dengan alasannya.
//
// Kalibrasi: ketiga tanggal dialog `EditDateClaimLife_Section` - yang memang
// disunting rute `PUT …/tanggal-klaim` - bertanda `Editable`, sedangkan
// kembarannya di layar Detail `Read-only`.

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// opsiSunting adalah hasil pembacaan satu wadah.
type opsiSunting map[string]string

// bacaOpsiSunting membaca `pyEditOptions` tiap sel ber-`pyValue` di rentang
// baris [dari, sampai] (1-based; sampai<=0 = sampai akhir berkas).
func bacaOpsiSunting(t *testing.T, letak string, dari, sampai int) opsiSunting {
	t.Helper()
	isi, err := os.ReadFile(letak)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v); vonis tidak terperiksa", err)
	}
	nilai := regexp.MustCompile(`<pyValue>(\.[A-Za-z_][A-Za-z0-9_]*)</pyValue>`)
	opsi := regexp.MustCompile(`<pyEditOptions>([^<]+)</pyEditOptions>`)
	hasil := opsiSunting{}
	terakhir := ""
	for i, b := range strings.Split(string(isi), "\n") {
		nomor := i + 1
		if nomor < dari || (sampai > 0 && nomor > sampai) {
			continue
		}
		if m := opsi.FindStringSubmatch(b); m != nil {
			terakhir = m[1]
		}
		if m := nilai.FindStringSubmatch(b); m != nil {
			// Tombol (`.pyTemplate*`) memang berulang; yang diikat per nama
			// hanya medan data.
			if _, ganda := hasil[m[1]]; ganda && !strings.HasPrefix(m[1], ".pyTemplate") {
				t.Errorf("%s muncul lebih dari sekali di wadah b%d-b%d; kunci per nama tidak cukup",
					m[1], dari, sampai)
			}
			hasil[m[1]] = terakhir
			terakhir = ""
		}
	}
	return hasil
}

const (
	letakDetail = `D:\XML\RNM_BRD\Claim Life\Section\ClaimLifeDetailGCNM.xml`
	letakPanel  = `D:\XML\RNM_BRD\Claim Life\Section\AdjustmentDetail_Section.xml`
	letakDialog = `D:\XML\RNM_BRD\Claim Life\Section\EditDateClaimLife_Section.xml`
)

// TestKalibrasiOpsiSuntingDariDialogEditDate - tafsirnya diuji atas jawaban
// yang sudah diketahui, sebelum dipakai untuk memvonis.
func TestKalibrasiOpsiSuntingDariDialogEditDate(t *testing.T) {
	dialog := bacaOpsiSunting(t, letakDialog, 1, 0)
	detail := bacaOpsiSunting(t, letakDetail, 1, 0)
	for _, medan := range []string{".CLAIM_RECEIVED_DATE", ".COMPLETE_DATE", ".CONFIRMATION_DATE"} {
		if dialog[medan] != "Editable" {
			t.Errorf("dialog Edit Date %s: %q, mau Editable - pembacanya rusak", medan, dialog[medan])
		}
	}
	if detail[".CLAIM_RECEIVED_DATE"] != "Read-only" {
		t.Errorf("layar Detail CLAIM_RECEIVED_DATE: %q, mau Read-only", detail[".CLAIM_RECEIVED_DATE"])
	}
}

// TestGridAdjustmentNolSelDapatDisunting - vonis butir br.
func TestGridAdjustmentNolSelDapatDisunting(t *testing.T) {
	// Wadah GRID: b17126-b19583, bukan seluruh berkas.
	grid := bacaOpsiSunting(t, letakDetail, 17126, 19583)
	mauGrid := []string{".STS_REJECT", ".ADJUSTMENT_DATE", ".ACCEPTEDNO", ".ACCEPTATION_DATE"}
	for _, medan := range mauGrid {
		if grid[medan] != "Read-only" {
			t.Errorf("grid b17126 %s: %q, mau Read-only", medan, grid[medan])
		}
	}
	for medan := range grid {
		if !strings.HasPrefix(medan, ".pyTemplate") && !berisi(mauGrid, medan) {
			t.Errorf("grid b17126 memuat medan %s yang tidak diperiksa", medan)
		}
	}

	panel := bacaOpsiSunting(t, letakPanel, 1, 0)
	if len(panel) < 20 {
		t.Fatalf("hanya %d medan panel terbaca; pembacanya rusak", len(panel))
	}
	// ⛔ Pengecualian bernama: sel tanpa `pyEditOptions`.
	tanpaOpsi := map[string]string{
		// b6245: medan daftar-sumber (`pyPromptClass`
		// ASM-FW-GISFW-Int-BANKACCOUNT) milik kontrol autocomplete
		// `.NameOfBank` b6122 - bukan sel tersendiri.
		".NAMEOFBANK": "sumber autocomplete .NameOfBank",
	}
	for medan, opsi := range panel {
		switch {
		case strings.HasPrefix(medan, ".pyTemplate"): // tombol, bukan medan data
		case opsi == "":
			if _, ada := tanpaOpsi[medan]; !ada {
				t.Errorf("panel %s: tanpa pyEditOptions dan tidak dikecualikan", medan)
			}
		case opsi != "Read-only":
			t.Errorf("panel Adjustment_Detail %s: %q - sel dapat disunting; butir br "+
				"menuntut rute sunting untuknya", medan, opsi)
		}
	}
	// Medan uang yang paling mungkin disangka dapat diisi - dinyatakan satu
	// per satu supaya vonisnya terbaca, bukan hanya tersirat.
	for _, medan := range []string{".CLAIM_GROSS", ".CLAIM_PAID", ".PCTClaim", ".CURRENCY",
		".SHARE_NUSANTARA_RE", ".CEDING_RETENTION", ".SUM_INSURED", ".SUM_REASURED",
		".SHARE_RETRO", ".RETROCEDED_SHARE", ".CLAIM_RETRO"} {
		if panel[medan] != "Read-only" {
			t.Errorf("panel %s: %q, mau Read-only", medan, panel[medan])
		}
	}
}

func berisi(daftar []string, x string) bool {
	for _, d := range daftar {
		if d == x {
			return true
		}
	}
	return false
}
