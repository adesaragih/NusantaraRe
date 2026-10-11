package models

import (
	"errors"
	"html"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

// akarKorpus - korpus READ-ONLY; uji yang membacanya DILEWATI bila tak terjangkau.
const akarKorpus = `D:\XML\RNM_BRD\Endorsement Life`

func bacaKorpus(t *testing.T, jalur string) string {
	t.Helper()
	b, err := os.ReadFile(akarKorpus + `\` + jalur)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	return html.UnescapeString(string(b))
}

// langkah memotong satu langkah activity: dari `<pyStepsActivityName>` …
// sampai langkah berikutnya - cukup untuk membaca pasangan Properties-nya.
func pasanganProperti(isi, awal, akhir string) [][2]string {
	i := strings.Index(isi, awal)
	if i < 0 {
		return nil
	}
	potong := isi[i:]
	if j := strings.Index(potong[len(awal):], akhir); j >= 0 {
		potong = potong[:len(awal)+j]
	}
	pola := regexp.MustCompile(`<PropertiesName>([^<]*)</PropertiesName>\s*<PropertiesValue>([^<]*)</PropertiesValue>`)
	var hasil [][2]string
	for _, m := range pola.FindAllStringSubmatch(potong, -1) {
		hasil = append(hasil, [2]string{m[1], m[2]})
	}
	return hasil
}

func TestPeriksaEdmTypeHanyaSatuDanTiga(t *testing.T) {
	for _, v := range []string{"1", "3"} {
		if err := PeriksaEdmType(v); err != nil {
			t.Errorf("%q ditolak: %v", v, err)
		}
	}
	for _, v := range []string{"", "2", "4", " 1", "1 ", "Batal"} {
		if err := PeriksaEdmType(v); err == nil {
			t.Errorf("%q diterima, padahal enum efektif {1, 3}", v)
		}
	}
}

func TestPengenalKasusTanpaPadding(t *testing.T) {
	id, err := RakitPengenalKasus("7")
	if err != nil || id != "EDML-7" {
		t.Fatalf("RakitPengenalKasus(7) = %q, %v", id, err)
	}
	for _, rusak := range []string{"", "x", "-1", "1.5"} {
		if _, err := RakitPengenalKasus(rusak); err == nil {
			t.Errorf("%q diterima", rusak)
		}
	}
	if !KasusEDM("EDML-1") || KasusEDM("NBLF-1") || KasusEDM("EDML-") {
		t.Error("KasusEDM salah menggolongkan pengenal")
	}
}

// TestPengenalEDMLSamaDenganPega - `EDML-<n>` (keputusan work owner 07-10-2026, migrasi 484 mulai 828).
func TestPengenalEDMLSamaDenganPega(t *testing.T) {
	if AwalanKasus != "EDML-" {
		t.Fatalf("AwalanKasus = %q", AwalanKasus)
	}
	if id, err := RakitPengenalKasus("828"); err != nil || id != "EDML-828" {
		t.Errorf("RakitPengenalKasus(828) = %q, %v", id, err)
	}
	for _, sah := range []string{"EDML-828", "EDML-1", "EDML-1234"} {
		if !KasusEDM(sah) {
			t.Errorf("%q ditolak", sah)
		}
	}
	// Awalan lama, NB, kosong, dan akhiran bukan angka bukan kasus modul ini.
	for _, bukan := range []string{"EDMLF-1", "EDMLF-828", "NBLF-828", "EDML-", "EDML", "EDML-12a", "EDML--1", "edml-1", ""} {
		if KasusEDM(bukan) {
			t.Errorf("%q diterima", bukan)
		}
	}
}

// PRODKE SATU hitungan mulai 0 (keputusan work owner 08-10-2026): NO_ENDORS = PL + "/" + PRODKE baru dua digit.
func TestNomorEndorsementDariPRODKE(t *testing.T) {
	for _, k := range []struct {
		polis  string
		prodKe int
		nomor  string
	}{
		{"RNML-QRL2.10.26.00042", 1, "RNML-QRL2.10.26.00042/01"}, // endorsement pertama (NB PRODKE 0 → 1)
		{"RNML-QRL2.10.26.00042", 2, "RNML-QRL2.10.26.00042/02"}, // dari versi 1
		{"UJI-PL-1", 9, "UJI-PL-1/09"},
		{"UJI-PL-1", 10, "UJI-PL-1/10"},
		{" UJI-PL-1 ", 100, "UJI-PL-1/100"}, // > 99: angka apa adanya (pola Pega b966)
	} {
		if n, err := NomorEndorsement(k.polis, k.prodKe); err != nil || n != k.nomor {
			t.Errorf("NomorEndorsement(%q, %d) = %q, %v; mau %q", k.polis, k.prodKe, n, err, k.nomor)
		}
	}
	if _, err := NomorEndorsement("", 1); err == nil {
		t.Error("nomor polis kosong diterima")
	}
	if _, err := NomorEndorsement("UJI-PL-1", 0); err == nil {
		t.Error("PRODKE 0 (new business) diterima sebagai nomor endorsement")
	}
}

// TestTanggalProduksiEndorsement - TGL_PROD JSON_POLIS (08-10-2026): bulan sama → tanggal Confirm; bergeser → tgl 1.
func TestTanggalProduksiEndorsement(t *testing.T) {
	wib := time.FixedZone("WIB", 7*60*60)
	for _, k := range []struct {
		saat    time.Time
		periode string
		mau     string
	}{
		{time.Date(2026, 10, 8, 15, 30, 0, 0, wib), "10.2026", "2026-10-08"},
		{time.Date(2026, 10, 28, 9, 0, 0, 0, wib), "11.2026", "2026-11-01"},
		{time.Date(2026, 12, 30, 9, 0, 0, 0, wib), "01.2027", "2027-01-01"},
		{time.Date(2026, 10, 31, 20, 0, 0, 0, time.UTC), "11.2026", "2026-11-01"}, // 01-11 03:00 WIB
	} {
		got, err := TanggalProduksi(k.saat, k.periode)
		if err != nil || got.Format("2006-01-02") != k.mau || got.Hour() != 0 {
			t.Errorf("%v %s = %v, %v; mau %s", k.saat, k.periode, got, err, k.mau)
		}
	}
	for _, salah := range []string{"", "2026-10", "13.2026"} {
		if _, err := TanggalProduksi(time.Now(), salah); !errors.Is(err, ErrTanggalProduksi) {
			t.Errorf("periode %q: %v", salah, err)
		}
	}
}

func TestStatusJenisDanStatusLamaVERBATIM(t *testing.T) {
	for tipe, mau := range map[string]string{"QR": "0", "QP": "0", "TR": "1", "TP": "1", "": "", "XX": ""} {
		if g := StatusJenis(tipe); g != mau {
			t.Errorf("StatusJenis(%q) = %q, mau %q", tipe, g, mau)
		}
	}
	for s, mau := range map[string]string{"Old": "1", "New": "0", "Delete": "0", "Batal": "0", "": "0"} {
		if g := StatusLama(s); g != mau {
			t.Errorf("StatusLama(%q) = %q, mau %q", s, g, mau)
		}
	}
}

func TestStatusHidupMenjagaPesertaNB(t *testing.T) {
	// ⛔ AC 48a: kosong/NULL HIDUP - penyaring naif NOT IN membuang seluruh NB.
	for s, mau := range map[string]bool{"": true, "Old": true, "New": true, "Delete": false, "Batal": false} {
		if g := StatusHidup(s); g != mau {
			t.Errorf("StatusHidup(%q) = %v, mau %v", s, g, mau)
		}
	}
}

func TestKeputusanIsLifeAccepted(t *testing.T) {
	if !Diterima("1") || Diterima("2") || Diterima("7") {
		t.Error("IsLifeAccepted: hanya 1 → Confirm")
	}
	if LabelKeputusanRiwayat("1") != "Accept" || LabelKeputusanRiwayat("2") != "Decline" || LabelKeputusanRiwayat("7") != "Decline" {
		t.Error("AddHistorySuggest b353 dilanggar")
	}
	for _, v := range []string{"1", "2", "7"} {
		if err := PeriksaKeputusan(v); err != nil {
			t.Errorf("%q ditolak", v)
		}
	}
	for _, v := range []string{"", "0", "3", "Accept"} {
		if err := PeriksaKeputusan(v); err == nil {
			t.Errorf("%q diterima", v)
		}
	}
}

func TestLebihBaruMemilihPRODKETerbesar(t *testing.T) {
	nb := Versi{Jenis: SumberWarisan, ID: "UJI-IDPEGA-1", ProdKe: 1}
	edm := Versi{Jenis: SumberAplikasi, ID: "EDML-1", ProdKe: 2, EdmType: "1"}
	if LebihBaru(nb, edm).ID != "EDML-1" || LebihBaru(edm, nb).ID != "EDML-1" {
		t.Error("PRODKE terbesar tidak menang")
	}
	seri := Versi{Jenis: SumberAplikasi, ID: "NBLF-1", ProdKe: 1}
	if LebihBaru(nb, seri).Jenis != SumberAplikasi {
		t.Error("seri PRODKE: sumber aplikasi harus menang")
	}
	if LebihBaru(Versi{}, nb).ID != nb.ID {
		t.Error("versi kosong menang")
	}
	if !(Versi{EdmType: "3"}).SudahBatal() || (Versi{EdmType: "1"}).SudahBatal() {
		t.Error("gerbang 4 @contains(...,\"3\") dilanggar")
	}
}

// TestKolomJurnalBalikDariKorpus - ke-32 kolom diturunkan ulang dari
// `SetPremi_EDM` 2.1 (`.X <= .X * -1`), urutan dan isinya.
func TestKolomJurnalBalikDariKorpus(t *testing.T) {
	isi := bacaKorpus(t, `Activity\SetPremi_EDM.xml`)
	pasangan := pasanganProperti(isi, "<pyStepsDescription>Set 0 jika EDM Batal</pyStepsDescription>", "<pyStepsDescription>flag pengurangan</pyStepsDescription>")
	var dariKorpus []string
	for _, p := range pasangan {
		kolom := strings.TrimPrefix(p[0], ".")
		if p[1] != "."+kolom+" * -1" {
			t.Errorf("langkah 2.1 memuat %q <= %q, bukan pembalikan tanda", p[0], p[1])
			continue
		}
		dariKorpus = append(dariKorpus, kolom)
	}
	if len(dariKorpus) != 32 {
		t.Fatalf("korpus memberi %d kolom, mau 32; pembacanya yang rusak?", len(dariKorpus))
	}
	if strings.Join(dariKorpus, ",") != strings.Join(KolomJurnalBalik, ",") {
		t.Errorf("KolomJurnalBalik tidak VERBATIM korpus:\n kode  %v\n korpus %v", KolomJurnalBalik, dariKorpus)
	}
}

// TestKolomKepalaSalinDariKorpus - 24 properti kepala `MappingEDMLife` 9
// (`.X <= WorkLife.X`), ditambah `PL_NUMBER ← TempWork.PolicyNo`.
func TestKolomKepalaSalinDariKorpus(t *testing.T) {
	isi := bacaKorpus(t, `Activity\MappingEDMLife.xml`)
	pasangan := pasanganProperti(isi, "<pyStepsDescription>Mapping detail dari Life</pyStepsDescription>", "<pyStepsDescription>Copy page dari Life ke EDM</pyStepsDescription>")
	var properti []string
	nomorPolis := false
	for _, p := range pasangan {
		nama := strings.TrimPrefix(strings.TrimPrefix(p[0], "pyWorkPage"), ".")
		switch {
		case p[1] == "TempWork.PolicyNo" && nama == "PremiumListSummary.PL_NUMBER":
			nomorPolis = true
		case p[1] == "WorkLife."+nama:
			properti = append(properti, nama)
		default:
			t.Errorf("pasangan langkah 9 tak dikenal: %q <= %q", p[0], p[1])
		}
	}
	if !nomorPolis {
		t.Error("langkah 9 tidak lagi memetakan PL_NUMBER dari Policy No (RALAT R01)")
	}
	var kode []string
	for _, k := range KolomKepalaSalin {
		kode = append(kode, k.Properti)
	}
	if strings.Join(properti, ",") != strings.Join(kode, ",") {
		t.Errorf("KolomKepalaSalin tidak VERBATIM korpus:\n kode   %v\n korpus %v", kode, properti)
	}
	for _, k := range KolomKepalaSalin {
		if !KolomSah(k.Kolom) {
			t.Errorf("kolom %q tidak sah dirakit ke SQL", k.Kolom)
		}
	}
}

func TestPesanVERBATIMNomorInvoice(t *testing.T) {
	if NomorInvoiceArasapas("UJI.PL.01") != "UJIPL01" {
		t.Error("b464 @replaceAll(PolicyNo,\".\",\"\") dilanggar")
	}
}
