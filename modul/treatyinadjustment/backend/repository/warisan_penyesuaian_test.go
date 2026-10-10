package repository_test

// Uji pengurai dokumen penyesuaian - TANPA Oracle.
//
// Yang dijaga: sisi New dari halaman akar, sisi Old dari `OLDDATA`, kunci
// yang TIDAK ADA tetap tidak ada, nilai dibawa apa adanya, dan `OLDDATA`
// bersarang tidak ditelusuri.

import (
	"testing"

	"nusantarare/modul/treatyinadjustment/backend/repository"
)

const dokUji = `{
 "ID":"1000080/R02","OLDID":"1000080/R01","EDMState":"2","IsProRate":true,
 "Commencement":"20180101","TreatyContractName":"BARU",
 "ValueDifference":{"RNMShare":"-2.5","BrokeragePercent":"0","Share":[{"Layer":"1"}]},
 "ProRatePercent":45.20547945205479452,
 "CurrencyList":[{"Currency":"USD","Conversion":"14250.5","Nested":{"x":"y"}},{"Currency":"IDR","Conversion":"1"}],
 "Retention":[ ],
 "OLDDATA":{
   "Commencement":"20170101","TreatyContractName":"LAMA","ContractRefNo":null,
   "CurrencyList":[{"Currency":"EUR","Conversion":"16000"}],
   "OLDDATA":{"TreatyContractName":"LEBIH LAMA"}
 }
}`

func TestUraiMemisahNewDanOld(t *testing.T) {
	baru, lama, err := repository.UraiDokumenPenyesuaian([]byte(dokUji))
	if err != nil {
		t.Fatal(err)
	}
	if baru.Medan["TreatyContractName"] != "BARU" || lama.Medan["TreatyContractName"] != "LAMA" {
		t.Errorf("sisi tertukar: baru %q lama %q", baru.Medan["TreatyContractName"], lama.Medan["TreatyContractName"])
	}
	if baru.Larik["CurrencyList"][0]["Currency"] != "USD" || lama.Larik["CurrencyList"][0]["Currency"] != "EUR" {
		t.Error("larik CurrencyList tidak dibaca per sisi")
	}
}

// ⛔ Kunci yang TIDAK ADA tetap tidak ada - bukan string kosong. Kunci yang
// ada bernilai `null` ADA, bernilai kosong. Layar membedakan keduanya.
func TestKunciTidakAdaBerbedaDariKosong(t *testing.T) {
	baru, lama, err := repository.UraiDokumenPenyesuaian([]byte(dokUji))
	if err != nil {
		t.Fatal(err)
	}
	if _, ada := baru.Medan["ContractRefNo"]; ada {
		t.Error("ContractRefNo tidak ada di halaman akar, tetapi terbaca ada")
	}
	if v, ada := lama.Medan["ContractRefNo"]; !ada || v != "" {
		t.Errorf("ContractRefNo null di OLDDATA harus ADA dan kosong, dapat ada=%v %q", ada, v)
	}
	if _, ada := lama.Larik["Retention"]; ada {
		t.Error("Retention tidak ada di OLDDATA, tetapi terbaca ada")
	}
	if r, ada := baru.Larik["Retention"]; !ada || len(r) != 0 {
		t.Errorf("Retention kosong di akar harus ADA dan nol baris, dapat ada=%v %d", ada, len(r))
	}
}

// ⛔ Nilai APA ADANYA: angka dengan digit aslinya (bukan lewat float),
// boolean sebagai teks, urutan larik dipertahankan.
func TestNilaiApaAdanya(t *testing.T) {
	baru, _, err := repository.UraiDokumenPenyesuaian([]byte(dokUji))
	if err != nil {
		t.Fatal(err)
	}
	if got := baru.Medan["ProRatePercent"]; got != "45.20547945205479452" {
		t.Errorf("ProRatePercent %q - digitnya berubah", got)
	}
	if baru.Medan["IsProRate"] != "true" {
		t.Errorf("IsProRate %q", baru.Medan["IsProRate"])
	}
	if baru.Medan["Commencement"] != "20180101" {
		t.Errorf("tanggal diterjemahkan di repository: %q", baru.Medan["Commencement"])
	}
	if c := baru.Larik["CurrencyList"]; len(c) != 2 || c[1]["Currency"] != "IDR" {
		t.Errorf("urutan larik berubah: %v", c)
	}
	if _, ada := baru.Larik["CurrencyList"][0]["Nested"]; ada {
		t.Error("medan bersarang di dalam elemen ikut dibawa")
	}
}

// Halaman tertanam satu tingkat dibaca lewat kunci bertitik - cara Section
// mengikatnya (`TreatyIn.ValueDifference.RNMShare`).
func TestKunciBertitikMembacaHalamanTertanam(t *testing.T) {
	baru, lama, err := repository.UraiDokumenPenyesuaian([]byte(dokUji))
	if err != nil {
		t.Fatal(err)
	}
	if baru.Medan["ValueDifference.RNMShare"] != "-2.5" {
		t.Errorf("ValueDifference.RNMShare %q", baru.Medan["ValueDifference.RNMShare"])
	}
	if _, ada := lama.Medan["ValueDifference.RNMShare"]; ada {
		t.Error("OLDDATA tidak punya ValueDifference, tetapi terbaca ada")
	}
}

// ⛔ `OLDDATA` di dalam `OLDDATA` TIDAK ditelusuri - Section mengikat satu
// tingkat saja.
func TestOldDataBersarangTidakDitelusuri(t *testing.T) {
	_, lama, err := repository.UraiDokumenPenyesuaian([]byte(dokUji))
	if err != nil {
		t.Fatal(err)
	}
	if lama.Medan["TreatyContractName"] == "LEBIH LAMA" {
		t.Error("sisi Old terbaca dari OLDDATA tingkat kedua")
	}
}

func TestDokumenTanpaOldDataMemberiSisiOldKosong(t *testing.T) {
	baru, lama, err := repository.UraiDokumenPenyesuaian([]byte(`{"TreatyContractName":"X"}`))
	if err != nil {
		t.Fatal(err)
	}
	if baru.Medan["TreatyContractName"] != "X" || len(lama.Medan) != 0 || lama.Larik == nil {
		t.Errorf("baru %v lama %v", baru, lama)
	}
}

// Negatif: dokumen rusak menghasilkan galat, bukan sisi kosong diam-diam.
func TestDokumenRusakDitolak(t *testing.T) {
	for _, d := range []string{`{"ID":`, `[1,2]`, `"teks"`} {
		if _, _, err := repository.UraiDokumenPenyesuaian([]byte(d)); err == nil {
			t.Errorf("%q: mau galat", d)
		}
	}
	if _, _, err := repository.UraiDokumenPenyesuaian([]byte(`{"OLDDATA":"bukan objek"}`)); err == nil {
		t.Error("OLDDATA bukan objek: mau galat")
	}
}

// ⭐ Larik BERTITIK (`ValueDifference.Share`) dan larik BERSARANG berindeks
// (`RnmLimitListDisplay(1).Value`) — bentuk ikatan grid Value Difference
// (`TreatyInTabsNPValueDifferenceProRate.xml` @506643, sel @2232609…).
func TestLarikBertitikDanBersarangBerindeks(t *testing.T) {
	dok := `{"ValueDifference":{"Share":[{"Layer":"1",
		"RnmLimitListDisplay":[{"Currency":"IDR","Value":"100"},{"Currency":"USD","Value":"7"}],
		"DeductionList":[{"Deduction":"9"}]}]}}`
	baru, _, err := repository.UraiDokumenPenyesuaian([]byte(dok))
	if err != nil {
		t.Fatal(err)
	}
	b := baru.Larik["ValueDifference.Share"]
	if len(b) != 1 {
		t.Fatalf("ValueDifference.Share %d baris, mau 1", len(b))
	}
	for k, v := range map[string]string{"Layer": "1", "RnmLimitListDisplay(1).Currency": "IDR", "RnmLimitListDisplay(2).Value": "7"} {
		if b[0][k] != v {
			t.Errorf("%s = %q, mau %q", k, b[0][k], v)
		}
	}
	// ⛔ Larik bersarang di LUAR daftar-izin tidak dibawa.
	if _, ada := b[0]["DeductionList(1).Deduction"]; ada {
		t.Error("DeductionList ikut dibawa walau tidak diikat sel mana pun")
	}
}

// Daftar bangkitan dan daftar tangan tidak berbenturan dan tidak kosong.
func TestDaftarKunciBangkitanTerbaca(t *testing.T) {
	if len(repository.MedanKerangka) == 0 || len(repository.LarikKerangka) == 0 {
		t.Fatal("daftar bangkitan kosong — jalankan alat/ekstrak_kerangka.py")
	}
}
