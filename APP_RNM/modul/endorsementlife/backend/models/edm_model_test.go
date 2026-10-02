package models

import (
	"html"
	"os"
	"regexp"
	"strings"
	"testing"
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
	if err != nil || id != "EDMLF-7" {
		t.Fatalf("RakitPengenalKasus(7) = %q, %v", id, err)
	}
	for _, rusak := range []string{"", "x", "-1", "1.5"} {
		if _, err := RakitPengenalKasus(rusak); err == nil {
			t.Errorf("%q diterima", rusak)
		}
	}
	if !KasusEDM("EDMLF-1") || KasusEDM("NBLF-1") || KasusEDM("EDMLF-") {
		t.Error("KasusEDM salah menggolongkan pengenal")
	}
}

// K3 keputusan work owner 01-10-2026 (OQ-EDM-008): NOMOR = rumus Pega `GenerateNoEDM_Life` 3–4
// (`Local.Prodke` int ← `PRODKE` baris terbaru b898/b899, `CARI4 = Prodke+1` b946, pad dua digit b967);
// VERSI tetap E1 (`NVL(PRODKE, 1)` + 1). Polis NB warisan (`PRODKE` kosong) → `/01`, lalu `/02`.
func TestNomorEndorsementRumusPega(t *testing.T) {
	kasus := []struct {
		nama   string
		polis  string
		versi  Versi
		nomor  string
		prodKe int
	}{
		{"NB warisan PRODKE kosong", "UJI-PL-1", Versi{ID: "UJI-IDPEGA-1", ProdKe: 1, UrutanPega: 0}, "UJI-PL-1/01", 2},
		{"endorsement pertama sistem baru", "UJI-PL-1", Versi{ID: "EDMLF-1", ProdKe: 2, UrutanPega: 1}, "UJI-PL-1/02", 3},
		{"endorsement warisan Pega /01", "UJI-PL-1", Versi{ID: "UJI-IDPEGA-2", ProdKe: 1, UrutanPega: 1}, "UJI-PL-1/02", 2},
		{"delapan", "UJI-PL-1", Versi{ID: "x", ProdKe: 9, UrutanPega: 8}, "UJI-PL-1/09", 10},
		{"dua digit tanpa pad tambahan b967", "UJI-PL-1", Versi{ID: "x", ProdKe: 10, UrutanPega: 9}, "UJI-PL-1/10", 11},
		{"tiga digit", " UJI-PL-1 ", Versi{ID: "x", ProdKe: 100, UrutanPega: 99}, "UJI-PL-1/100", 101},
	}
	for _, k := range kasus {
		n, b, err := NomorEndorsement(k.polis, k.versi)
		if err != nil || n != k.nomor || b != k.prodKe {
			t.Errorf("%s: NomorEndorsement = %q, %d, %v; mau %q, %d", k.nama, n, b, err, k.nomor, k.prodKe)
		}
	}
	if _, _, err := NomorEndorsement("", Versi{ProdKe: 1}); err == nil {
		t.Error("nomor polis kosong diterima")
	}
	if _, _, err := NomorEndorsement("UJI-PL-1", Versi{ProdKe: 0}); err == nil {
		t.Error("PRODKE 0 diterima sebagai versi")
	}
	if _, _, err := NomorEndorsement("UJI-PL-1", Versi{ProdKe: 1, UrutanPega: -1}); err == nil {
		t.Error("urutan Pega negatif diterima")
	}
}

// Urutan Pega endorsement sistem baru = akhiran nomornya: Pega menulis `JSON_POLIS.PRODKE = CARI4`
// (`InsertJsonPolisEDM` b101), angka yang sama dengan akhiran `NOPOLIS||'/'||CARI14`.
func TestUrutanDariNomor(t *testing.T) {
	for nomor, mau := range map[string]int{"UJI-PL-1/01": 1, "UJI-PL-1/10": 10, "UJI/PL/1/100": 100, " UJI-PL-1/02 ": 2} {
		if got, err := UrutanDariNomor(nomor); err != nil || got != mau {
			t.Errorf("UrutanDariNomor(%q) = %d, %v; mau %d", nomor, got, err, mau)
		}
	}
	for _, rusak := range []string{"", "UJI-PL-1", "UJI-PL-1/", "UJI-PL-1/0x", "UJI-PL-1/-1"} {
		if _, err := UrutanDariNomor(rusak); err == nil {
			t.Errorf("%q diterima", rusak)
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
	edm := Versi{Jenis: SumberAplikasi, ID: "EDMLF-1", ProdKe: 2, EdmType: "1"}
	if LebihBaru(nb, edm).ID != "EDMLF-1" || LebihBaru(edm, nb).ID != "EDMLF-1" {
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
