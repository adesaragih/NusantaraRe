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
// Kalibrasi tafsir: `pyEditOptions` sel ditulis SEBELUM `pyValue`-nya, di
// rentang sesudah `pyValue` sebelumnya. `.CLAIM_RECEIVED_DATE` bertanda
// `Read-only` di layar Detail tetapi `Editable` di dialog
// `EditDateClaimLife_Section` (b1034) - dialog yang memang menyunting tanggal
// itu (rute `PUT …/tanggal-klaim`).

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// opsiSuntingPerMedan membaca `pyEditOptions` terakhir di sel tiap `pyValue`.
func opsiSuntingPerMedan(t *testing.T, letak string) map[string]string {
	t.Helper()
	isi, err := os.ReadFile(letak)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v); vonis tidak terperiksa", err)
	}
	nilai := regexp.MustCompile(`<pyValue>(\.[A-Za-z_][A-Za-z0-9_]*)</pyValue>`)
	opsi := regexp.MustCompile(`<pyEditOptions>([^<]+)</pyEditOptions>`)
	hasil := map[string]string{}
	terakhir := ""
	for _, b := range strings.Split(string(isi), "\n") {
		if m := opsi.FindStringSubmatch(b); m != nil {
			terakhir = m[1]
		}
		if m := nilai.FindStringSubmatch(b); m != nil {
			hasil[m[1]] = terakhir
			terakhir = ""
		}
	}
	return hasil
}

// TestKalibrasiOpsiSuntingDariDialogEditDate - tafsirnya diuji atas jawaban
// yang sudah diketahui, sebelum dipakai untuk memvonis.
func TestKalibrasiOpsiSuntingDariDialogEditDate(t *testing.T) {
	dialog := opsiSuntingPerMedan(t,
		`D:\XML\RNM_BRD\Claim Life\Section\EditDateClaimLife_Section.xml`)
	if dialog[".CLAIM_RECEIVED_DATE"] != "Editable" {
		t.Errorf("dialog Edit Date: CLAIM_RECEIVED_DATE %q, mau Editable - pembacanya rusak",
			dialog[".CLAIM_RECEIVED_DATE"])
	}
	detail := opsiSuntingPerMedan(t,
		`D:\XML\RNM_BRD\Claim Life\Section\ClaimLifeDetailGCNM.xml`)
	if detail[".CLAIM_RECEIVED_DATE"] != "Read-only" {
		t.Errorf("layar Detail: CLAIM_RECEIVED_DATE %q, mau Read-only",
			detail[".CLAIM_RECEIVED_DATE"])
	}
}

// TestGridAdjustmentNolSelDapatDisunting - vonis butir br.
func TestGridAdjustmentNolSelDapatDisunting(t *testing.T) {
	grid := opsiSuntingPerMedan(t,
		`D:\XML\RNM_BRD\Claim Life\Section\ClaimLifeDetailGCNM.xml`)
	// Keempat kolom data grid `.AdjustmentList` b17126.
	for _, medan := range []string{".STS_REJECT", ".ADJUSTMENT_DATE", ".ACCEPTEDNO", ".ACCEPTATION_DATE"} {
		if grid[medan] != "Read-only" {
			t.Errorf("grid b17126 %s: %q, mau Read-only", medan, grid[medan])
		}
	}
	panel := opsiSuntingPerMedan(t,
		`D:\XML\RNM_BRD\Claim Life\Section\AdjustmentDetail_Section.xml`)
	if len(panel) < 20 {
		t.Fatalf("hanya %d medan panel terbaca; pembacanya rusak", len(panel))
	}
	for medan, opsi := range panel {
		// Tombol (`pyTemplate*`) bukan medan data.
		if strings.HasPrefix(medan, ".pyTemplate") {
			continue
		}
		if opsi == "Editable" || opsi == "Auto" {
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
