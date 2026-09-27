package models_test

// Uji aturan murni diagnosa - butir bd.
//
// Yang dijaga di sini: gerbang `STS_REJECT`, penomoran urut, lebar kolom, dan
// daftar tahap yang membuka gridnya. Nol Oracle, nol jam.

import (
	"encoding/json"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"nusantarare/internal/models"
)

func TestDiagnosaTerkunciMengikutiKeempatGerbang(t *testing.T) {
	// ⛔ TABEL, bukan satu kasus. Kalimat yang ditiru menyebut DUA nilai
	// dengan `||`, dan yang menulis hanya satu di antaranya akan lolos uji
	// satu-kasus mana pun - lalu membiarkan diagnosa peserta yang DIAKSEP
	// tetap dapat disunting, yaitu separuh gerbangnya hilang tanpa berbunyi.
	for _, u := range []struct {
		nama   string
		kode   string
		terkun bool
	}{
		{"kosong - belum diputus", "", false},
		{"Outstanding", models.KodeOutstanding, false},
		{"Aksep", models.KodeAksep, true},
		{"Ditolak", models.KodeDitolak, true},
		{"berspasi - tetap terbaca", "  1  ", true},
		{"kode warisan 4 - TIDAK mengunci", "4", false},
	} {
		if got := models.DiagnosaTerkunci(u.kode); got != u.terkun {
			t.Errorf("%s: DiagnosaTerkunci(%q) = %v, mau %v",
				u.nama, u.kode, got, u.terkun)
		}
	}
}

func TestGerbangDiagnosaMemakaiKodeDariModels(t *testing.T) {
	// ⛔ Gerbangnya harus bergeser BERSAMA mesin status, bukan sendiri.
	// Kalau suatu hari KodeAksep bukan lagi "1", kalimat yang menulis '1'
	// telanjang akan tetap hijau dan tetap salah.
	if models.KodeAksep != "1" || models.KodeDitolak != "2" {
		t.Fatalf("kode bergeser: Aksep=%q Ditolak=%q; "+
			"gerbang `pyDisabledWhen` b4682 menyebut '1' dan '2'",
			models.KodeAksep, models.KodeDitolak)
	}
}

func TestUrutanBerikutnyaMulaiDariSatu(t *testing.T) {
	if got := models.UrutanBerikutnya(nil); got != 1 {
		t.Errorf("daftar kosong -> %d, mau 1", got)
	}
	if got := models.UrutanBerikutnya([]models.Diagnosa{}); got != 1 {
		t.Errorf("daftar nol panjang -> %d, mau 1", got)
	}
	// ⛔ TERTINGGI, bukan panjang daftar. Daftar yang pernah dihapus di
	// tengahnya lalu belum dirapatkan akan memberi nomor yang sudah dipakai
	// bila yang dihitung panjangnya.
	daftar := []models.Diagnosa{{Urutan: 1}, {Urutan: 5}, {Urutan: 3}}
	if got := models.UrutanBerikutnya(daftar); got != 6 {
		t.Errorf("tertinggi 5 -> %d, mau 6", got)
	}
}

func TestRapatkanUrutanMenutupLubang(t *testing.T) {
	rapat := models.RapatkanUrutan([]models.Diagnosa{
		{ID: 10, Urutan: 1}, {ID: 30, Urutan: 3}, {ID: 40, Urutan: 4},
	})
	for i, d := range rapat {
		if d.Urutan != i+1 {
			t.Errorf("baris %d berurut %d, mau %d", d.ID, d.Urutan, i+1)
		}
	}
	// ⚠️ Salinan, bukan tempat: masukan tidak boleh ikut berubah. Daftar yang
	// berubah di tempat akan membuat pemanggil yang masih memegang bacaan
	// lamanya melihat urutan yang sudah bergeser.
	asal := []models.Diagnosa{{ID: 1, Urutan: 7}}
	models.RapatkanUrutan(asal)
	if asal[0].Urutan != 7 {
		t.Errorf("masukan ikut berubah: urutan = %d, mau tetap 7", asal[0].Urutan)
	}
}

func TestPotongDiagnosaMelaporkanBukanMemotong(t *testing.T) {
	if kolom, muat := models.PotongDiagnosa("A00", "Diabetes", "Metabolik"); !muat {
		t.Errorf("nilai wajar ditolak di kolom %s", kolom)
	}
	for _, u := range []struct {
		nama             string
		icd, isi, grup   string
		kolomYangDiharap string
	}{
		{"ICD kepanjangan", strings.Repeat("A", 101), "x", "", "ICD_CODE"},
		{"nama kepanjangan", "A00", strings.Repeat("B", 1001), "", "DISEASE"},
		{"grup kepanjangan", "A00", "x", strings.Repeat("C", 256), "GROUP_DIAGNOSE"},
	} {
		kolom, muat := models.PotongDiagnosa(u.icd, u.isi, u.grup)
		if muat {
			t.Errorf("%s: diterima, mau ditolak", u.nama)
		}
		if kolom != u.kolomYangDiharap {
			t.Errorf("%s: kolom = %q, mau %q", u.nama, kolom, u.kolomYangDiharap)
		}
	}
	// ⛔ Batas TEPAT masih muat. Uji yang hanya mencoba "jauh kepanjangan"
	// tidak dapat membedakan `>` dari `>=`, dan yang kedua menolak nilai
	// terpanjang yang sah - kolomnya selebar itu justru supaya muat.
	if _, muat := models.PotongDiagnosa(strings.Repeat("A", 100),
		strings.Repeat("B", 1000), strings.Repeat("C", 255)); !muat {
		t.Error("nilai sepanjang lebar kolomnya ditolak")
	}
}

func TestTahapPembukaGridPesertaAdaTiga(t *testing.T) {
	// ⛔ TIGA, dan Medical Check TERMASUK. Ia tidak MEMUAT
	// `ClaimLifeDetailGCNM`; ia MEMBUKANYA lewat `ViewClaimDetailLifeGCNM`
	// b17416. Berhenti pada daftar pemuat berarti menutup layar diagnosa
	// bagi Medical Advisor - orang yang paling masuk akal mengisinya.
	for _, u := range []struct {
		tahap models.Tahap
		mau   bool
	}{
		{models.TahapOutstanding, true},
		{models.TahapMedicalCheck, true},
		{models.TahapClaimAnalis, true},
		{models.TahapInputRegister, false},
		{models.TahapTidakDikenal, false},
	} {
		if got := models.TahapBergridPeserta(u.tahap); got != u.mau {
			t.Errorf("TahapBergridPeserta(%v) = %v, mau %v", u.tahap, got, u.mau)
		}
	}
	if n := len(models.TahapPembukaGridPeserta()); n != 3 {
		t.Errorf("TahapPembukaGridPeserta() = %d tahap, mau 3", n)
	}
}

// TestGerbangDiagnosaVERBATIMDariKorpus membaca kalimatnya LANGSUNG dari XML.
//
// ⛔ Penjaga yang membandingkan salinan dengan salinan tidak menjaga apa pun.
// Uji ini membuka berkas korpus dan menghitung berapa kali kalimat gerbangnya
// muncul - kalau rule-nya berubah, di sinilah ketahuannya, bukan di rapat.
//
// ⚠️ Korpus DIBACA, tidak pernah ditulis. Bila ia tidak ada (mesin lain),
// uji ini MELEWATI, bukan gagal - dan mengatakannya, supaya "hijau" tidak
// berarti "terperiksa".
func TestGerbangDiagnosaVERBATIMDariKorpus(t *testing.T) {
	const letak = `D:\XML\RNM_BRD\Claim Life\Section\ClaimLifeDetailGCNM.xml`
	isi, err := os.ReadFile(letak)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v); gerbang tidak terperiksa", err)
	}
	const kalimat = `<pyDisabledWhen>.STS_REJECT=='1' || .STS_REJECT=='2'</pyDisabledWhen>`
	// ⛔ TUJUH, bukan empat - dan uji inilah yang mengoreksinya.
	// Ronde pertama uji ini menuntut 4 sebab bacaan kami hanya menyisir
	// rentang grid (b4100-b6260). Kalimat yang sama menjaga TIGA medan
	// catatan di layar yang sama:
	//
	//	b2628  `.ADMIN_NOTES`     (label `ADMIN NOTES` b2611)
	//	b4682  `Add`              grid diagnosa
	//	b5059  `Find Disease`     grid diagnosa
	//	b5870  `.GROUPDIAGNOSE`   grid diagnosa
	//	b6152  `Delete`           grid diagnosa
	//	b7335  `.RECOMMENDATION`  (label `RECOMMENDATION` b7316)
	//	b15234 `.NOTES`           (label `CLAIM ANALYST` b15217)
	//
	// ⭐ Artinya gerbangnya lebih luas daripada diagnosa: peserta yang
	// sudah diputus membekukan SELURUH isian layar Detail, bukan hanya
	// daftar diagnosanya. Ketiga medan catatan itu belum dibangun; ketika
	// dibangun, gerbang yang sama berlaku - dan angka 7 di sini yang akan
	// mengingatkannya.
	if n := strings.Count(string(isi), kalimat); n != 7 {
		t.Errorf("gerbang `%s` muncul %d kali, mau 7 "+
			"(4 di grid diagnosa: b4682, b5059, b5870, b6152; "+
			"3 di medan catatan: b2628, b7335, b15234)", kalimat, n)
	}
	// Dan ketiga medan catatan itu masih ada - supaya "7" tidak dapat
	// dipenuhi oleh tujuh kemunculan yang berbeda sama sekali.
	for _, label := range []string{
		"<pyValue>.ADMIN_NOTES</pyValue>",
		"<pyValue>.RECOMMENDATION</pyValue>",
		"<pyLabelPreview>CLAIM ANALYST</pyLabelPreview>",
	} {
		if !strings.Contains(string(isi), label) {
			t.Errorf("korpus tidak lagi memuat medan bergerbang %s", label)
		}
	}
	// Dan keempat kontrolnya masih bernama seperti yang tiket 08 catat.
	for _, label := range []string{
		`<pyLabel>Add</pyLabel>`,
		`<pyLabel>Find Disease</pyLabel>`,
		`<pyLabel>Delete</pyLabel>`,
		`<pyAction>addRow</pyAction>`,
		`<pyAction>deleteRow</pyAction>`,
		`<pyPageListProperty>.DiagnoseList</pyPageListProperty>`,
	} {
		if !strings.Contains(string(isi), label) {
			t.Errorf("korpus tidak lagi memuat %s", label)
		}
	}
}

// TestNamaJSONDiagnosaDikunci mengunci himpunan kunci JSON diagnosa.
//
// ⛔ SISI GO DARI KONTRAK DUA SISI. Pasangannya `Diagnosa` di
// `frontend/src/services/api.ts` dan `GridDiagnosa.test.ts`.
//
// ⛔ CACAT YANG NYARIS TERJADI, dan sebab uji ini lahir. Ronde pertama
// `models.Diagnosa` ditulis TANPA tag JSON sama sekali, sementara `api.ts`
// sudah mendeklarasikan `kodeIcd`, `pesertaId`, dan `groupDiagnose`. Go akan
// mengirim `KodeICD`, `PesertaID`, `GroupDiagnose`; React membaca
// `undefined`; grid tampil dengan tiga kolom kosong dan nol galat di kedua
// sisi. Itu bentuk cacat lintas-lapis KEENAM di modul ini - dan yang kelima
// lahir di giliran yang sama ketika keempat pendahulunya didaftarkan.
//
// ⛔ Dikunci sebagai HIMPUNAN UTUH: daftar nama terlarang selalu kalah dari
// nama yang belum terpikirkan.
func TestNamaJSONDiagnosaDikunci(t *testing.T) {
	b, err := json.Marshal(models.Diagnosa{
		ID: 1, PesertaID: "UJI-PES-1", Urutan: 1, KodeICD: "E11",
	})
	if err != nil {
		t.Fatal(err)
	}
	var peta map[string]json.RawMessage
	if err := json.Unmarshal(b, &peta); err != nil {
		t.Fatal(err)
	}
	kunci := make([]string, 0, len(peta))
	for k := range peta {
		kunci = append(kunci, k)
	}
	sort.Strings(kunci)
	mau := []string{"groupDiagnose", "id", "kodeIcd", "kodeStatus", "nama",
		"pesertaId", "urutan"}
	if !reflect.DeepEqual(kunci, mau) {
		t.Errorf("kunci JSON diagnosa = %v, mau %v.\n"+
			"Ubah JUGA tipe Diagnosa di api.ts - keduanya satu kontrak.", kunci, mau)
	}
	// ⛔ Dan tidak satu pun kunci di atas dapat memuat nama orang: diagnosa
	// adalah nama PENYAKIT, bukan nama tertanggung.
	for _, k := range kunci {
		if strings.Contains(strings.ToLower(k), "insured") ||
			strings.Contains(strings.ToLower(k), "holder") {
			t.Errorf("kunci %q menyerempet nama orang", k)
		}
	}
}
