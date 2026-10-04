package predikat

// Test registry predikat EDM (tiket E01). Mesin evaluatornya salinan nbfacin
// (lihat predikat.go); yang diuji di sini adalah registry VARIAN EDM dan
// sumber datanya. Nilai sintetis; tidak ada nomor polis atau nama.

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type kasusUji map[string]string

func (k kasusUji) Nilai(jalur string) (string, bool) {
	v, ada := k[jalur]
	return v, ada
}

func lini(jenis string) KasusEDM { return KasusEDM{JenisBisnisLama: &jenis} }

func harus(t *testing.T, label, nama string, k Kasus, mau bool) {
	t.Helper()
	got, err := Eval(nama, k)
	if err != nil || got != mau {
		t.Errorf("%s: %s = %v (%v), mau %v", label, nama, got, err, mau)
	}
}

func panikMemuat(t *testing.T, label, alasan string, f func()) {
	t.Helper()
	defer func() {
		t.Helper()
		r := recover()
		if r == nil {
			t.Errorf("%s: tidak panic", label)
		} else if s := fmt.Sprint(r); !strings.Contains(s, alasan) {
			t.Errorf("%s: panic %q, mau memuat %q", label, s, alasan)
		}
	}()
	f()
}

// TestRegistryUtuh - seluruh registry: asal, logika terurai, label berbaris
// kondisi, rujukan ada.
//
// [terverifikasi] 202 berkas `Endorsment Fac In\When`, 201 identitas, dihitung
// dua cara 01-10-2026:
//
//	cara 1: ls *.xml | wc -l → 202; grep -ho '<pxInsName>[^<]*</pxInsName>' *.xml | sort -u | wc -l → 201
//	cara 2: pengurai XML Python, nama sesudah '!' huruf besar → 201, kembar 1 (ISFLAGOLDDATA)
func TestRegistryUtuh(t *testing.T) {
	if len(registry) != 201 {
		t.Errorf("registry %d predikat, mau 201", len(registry))
	}
	for nama, pr := range registry {
		if nama != strings.ToUpper(nama) {
			t.Errorf("%s: kunci wajib huruf besar", nama)
		}
		if len(pr.asal) == 0 || !strings.HasPrefix(pr.asal[0], `Endorsment Fac In\When\`) || !strings.Contains(pr.asal[0], "!"+nama+")") {
			t.Errorf("%s: asal %v bukan berkas EDM beridentitas", nama, pr.asal)
		}
		if pr.panik != "" {
			continue
		}
		label := pohonLogika[nama].Label()
		if len(label) != len(pr.kondisi) {
			t.Errorf("%s: %d label, %d baris kondisi", nama, len(label), len(pr.kondisi))
		}
		for _, l := range label {
			kd, ada := pr.kondisi[l]
			if !ada {
				t.Errorf("%s: label %s tanpa kondisi", nama, l)
			} else if kd.jenis == rujukWhen {
				if _, ada := registry[kd.rujukan]; !ada {
					t.Errorf("%s: merujuk %s yang tidak ada", nama, kd.rujukan)
				}
			}
		}
	}
	if asal := registry["ISFLAGOLDDATA"].asal; len(asal) != 2 {
		t.Errorf("ISFLAGOLDDATA asal %v, mau dua berkas", asal)
	}
}

// TestSumberJenisBisnisLamaEnamPredikat - sumber data per-rule: tepat enam
// predikat membaca hasil query (cara 2 atas registry; cara 1 grep korpus, lihat
// JalurJenisBisnisLama), dan sisanya properti kasus.
func TestSumberJenisBisnisLamaEnamPredikat(t *testing.T) {
	var dapat []string
	for nama := range registry {
		for _, j := range Sumber(nama) {
			if j == JalurJenisBisnisLama {
				dapat = append(dapat, nama)
			}
		}
	}
	mau := map[string]bool{"ISANEKA": true, "ISFIRE": true, "ISGOLFINSURANCE": true, "ISMBU": true, "ISMARINECARGO": true, "ISPA": true}
	if len(dapat) != len(mau) {
		t.Errorf("pembaca %s: %v", JalurJenisBisnisLama, dapat)
	}
	for _, n := range dapat {
		if !mau[n] {
			t.Errorf("%s membaca jenis bisnis lama, tidak diharapkan", n)
		}
	}
	if s := Sumber("IsEDM"); !reflect.DeepEqual(s, []string{"pyWorkPage.OfferFacIn.QuotationData.StatusBusiness"}) {
		t.Errorf("Sumber IsEDM %v", s)
	}
}

// TestTigaPuluhTujuhCabangLini - 5 · 25 · 1 · 4 · 1 · 1 = 37 baris kondisi
// (E01), bukan 203. Dihitung dari registry; cara 2: grep `pyConditionValue1`
// per berkas korpus pada saat tiket ditulis.
func TestTigaPuluhTujuhCabangLini(t *testing.T) {
	jumlah := 0
	for nama, mau := range map[string]int{"ISFIRE": 5, "ISANEKA": 25, "ISPA": 1, "ISMBU": 4, "ISMARINECARGO": 1, "ISGOLFINSURANCE": 1} {
		if n := len(registry[nama].kondisi); n != mau {
			t.Errorf("%s: %d cabang, mau %d", nama, n, mau)
		}
		jumlah += len(registry[nama].kondisi)
	}
	if jumlah != 37 {
		t.Errorf("total %d, mau 37", jumlah)
	}
}

// TestPredikatLiniMembacaPolisLama - varian EDM membaca hasil query, BUKAN
// `pyWorkPage.Quotation.BusinessType` seperti varian NB.
func TestPredikatLiniMembacaPolisLama(t *testing.T) {
	harus(t, "fire", "IsFire", lini("Fire"), true)
	harus(t, "fire style", "IsFire", lini("FireStyle2"), true)
	harus(t, "bukan fire", "IsFire", lini("PA"), false)
	harus(t, "pa", "IsPA", lini("PA"), true)
	harus(t, "life dari langkah 8 bukan lini mana pun", "IsFire", lini("Life"), false)
	// Properti kasus ala NB tidak berpengaruh.
	harus(t, "properti NB diabaikan", "IsFire", KasusEDM{
		Properti:        map[string]string{"pyWorkPage.Quotation.BusinessType": "Fire"},
		JenisBisnisLama: ptr("Aneka"),
	}, false)
}

func ptr(s string) *string { return &s }

// TestUrutanQueryMengikat - predikat lini dinilai sebelum query → panic, bukan
// salah diam-diam.
func TestUrutanQueryMengikat(t *testing.T) {
	panikMemuat(t, "IsFire tanpa query", "sebelum query", func() { _, _ = Eval("IsFire", KasusEDM{}) })
	// Predikat non-lini tidak menyentuh query.
	harus(t, "IsEDM tanpa query", "IsEDM", KasusEDM{Properti: map[string]string{"pyWorkPage.OfferFacIn.QuotationData.StatusBusiness": "3"}}, true)
}

// TestK046_IsMBU_CabangGanda_HasilTidakBerubah - MBUCar dan MBUMotorCycle
// masing-masing muncul dua kali; redundan.
func TestK046_IsMBU_CabangGanda_HasilTidakBerubah(t *testing.T) {
	hitung := map[string]int{}
	for _, kd := range registry["ISMBU"].kondisi {
		hitung[kd.kanan.teks]++
	}
	if hitung["MBUCar"] != 2 || hitung["MBUMotorCycle"] != 2 {
		t.Errorf("cabang ISMBU %v", hitung)
	}
	harus(t, "mobil", "IsMBU", lini("MBUCar"), true)
	harus(t, "motor", "IsMBU", lini("MBUMotorCycle"), true)
	harus(t, "lain", "IsMBU", lini("Fire"), false)
}

// TestK046_IsAneka_LabelKodeBisnis_Diabaikan_IkutiValue1 - teks tampilan
// `pyConditionString` IsAneka berbunyi `Kode Bisnis = "24"` dsb.
// ([terverifikasi] `grep -o "<pyConditionString>[^<]*" IsAneka.xml`), tetapi
// yang mengikat ekspresi tersimpan atas jenis bisnis.
func TestK046_IsAneka_LabelKodeBisnis_Diabaikan_IkutiValue1(t *testing.T) {
	for _, kd := range registry["ISANEKA"].kondisi {
		if kd.kiri != JalurJenisBisnisLama {
			t.Errorf("ISANEKA membaca %s", kd.kiri)
		}
	}
	harus(t, "kode bisnis tampilan tidak berlaku", "IsAneka", KasusEDM{
		Properti: map[string]string{"pyWorkPage.Quotation.BusinessCode": "24"}, JenisBisnisLama: ptr(""),
	}, false)
	harus(t, "jenis bisnis aneka", "IsAneka", lini("Liability"), true)
}

// TestK046_COB_EDM_TanpaDelegasiSubRule - keenam predikat lini EDM
// membandingkan literal; nol `evaluateWhen` ([terverifikasi] `grep -c
// evaluateWhen` keenam berkas → 0).
func TestK046_COB_EDM_TanpaDelegasiSubRule(t *testing.T) {
	for _, nama := range []string{"ISFIRE", "ISANEKA", "ISPA", "ISMBU", "ISMARINECARGO", "ISGOLFINSURANCE"} {
		for l, kd := range registry[nama].kondisi {
			if kd.jenis != banding {
				t.Errorf("%s baris %s bukan perbandingan literal", nama, l)
			}
		}
	}
}

// TestIsEDMBukanKebalikanIsNotEDM - IsEDM membaca OfferFacIn.QuotationData,
// IsNotEDM membaca Quotation: keduanya bisa benar sekaligus.
func TestIsEDMBukanKebalikanIsNotEDM(t *testing.T) {
	k := kasusUji{"pyWorkPage.OfferFacIn.QuotationData.StatusBusiness": "3", "pyWorkPage.Quotation.StatusBusiness": "1"}
	harus(t, "IsEDM", "IsEDM", k, true)
	harus(t, "IsNotEDM", "IsNotEDM", k, true)
}

// TestPredikatPembayaran - gerbang porsi pembayaran (porsipembayaran.go).
func TestPredikatPembayaran(t *testing.T) {
	k := func(edmType, typ string) kasusUji {
		return kasusUji{
			"pyWorkPage.OfferFacIn.QuotationData.StatusBusiness": "3",
			"pyWorkPage.OfferFacIn.QuotationData.EdmType":        edmType,
			"pyWorkPage.OfferFacIn.QuotationData.Type":           typ,
		}
	}
	harus(t, "extend", "IsEdmExtendPeriod", k("4", "1"), true)
	harus(t, "extend type lain", "IsEdmExtendPeriod", k("4", "3"), false)
	harus(t, "adj TSI", "IsEdmAdjTSI", k("4", "2"), true)
	harus(t, "adj rate lewat adj TSI", "IsEdmAdjRate", k("4", "2"), true)
}

// TestNilaiRahasiaDisaring - predikat yang membandingkan nomor polis atau nama
// marketing dengan literal menjadi panic tanpa nilai (generator + saring).
func TestNilaiRahasiaDisaring(t *testing.T) {
	panikMemuat(t, "IsErrorSpreading", "nomor polis", func() { _, _ = Eval("IsErrorSpreading", kasusUji{}) })
	panikMemuat(t, "IsTBonding", "nama", func() { _, _ = Eval("IsTBonding", kasusUji{}) })
}

// TestRegistryTanpaLiteralIdentitas - penjaga tanpa korpus (CLAUDE.md §4.10,
// §10): tidak satu pun baris kondisi membandingkan medan identitas dengan
// literal TIDAK-KOSONG. Literal kosong (cek isian) boleh.
func TestRegistryTanpaLiteralIdentitas(t *testing.T) {
	identitas := map[string]bool{"OldPolicyNo": true, "PolicyNo": true, "MarketingName": true, "pyTelephone": true, "pyUserIdentifier": true}
	for nama, pr := range registry {
		for l, kd := range pr.kondisi {
			seg := strings.Split(kd.kiri, ".")
			if kd.jenis == banding && identitas[seg[len(seg)-1]] && kd.kanan.teks != "" {
				t.Errorf("%s baris %s: %s dibandingkan dengan literal (nilai tidak dicetak)", nama, l, kd.kiri)
			}
		}
	}
}

// TestTafsirBerbedaDanNamaTakDikenal - sikap mesin (butir 20) berlaku di EDM.
func TestTafsirBerbedaDanNamaTakDikenal(t *testing.T) {
	_, err := Eval("IsEDM", kasusUji{"pyWorkPage.OfferFacIn.QuotationData.StatusBusiness": "3.0"})
	if !errors.Is(err, ErrTafsirBerbeda) {
		t.Errorf("3.0 vs angka 3: %v", err)
	}
	panikMemuat(t, "nama tak dikenal", "tidak ada di registry", func() { _, _ = Eval("IsTidakAda", kasusUji{}) })
}

// --- E02: tiga predikat yang perilakunya berbeda di endorsement ---
//
// Varian NB dikutip dari korpus `NB FacIn/When/<nama>.xml` (bukan dari kode
// nbfacin): IsUW A OR B OR C (Director, GroupLeader, Underwriting); IsClaim
// `pyWorkPage.pyWorkIDPrefix = "CLM-"`; IsTravel satu sumber
// `pyWorkPage.Quotation.BusinessType`.

const workbasket = "pxRequestor.OperatorID.pyWorkBasketList(1).pyWorkBasketName"

// TestK046_IsUW_MarketingSebagaiUnderwriter_EDM - EDM mengakui empat
// workbasket; ReasFacInMarketing tambahannya.
func TestK046_IsUW_MarketingSebagaiUnderwriter_EDM(t *testing.T) {
	if n := len(registry["ISUW"].kondisi); n != 4 {
		t.Errorf("ISUW %d workbasket, mau 4 (NB 3)", n)
	}
	for _, wb := range []string{"ReasFacInDirector", "ReasFacInMarketing", "ReasFacInGroupLeader", "ReasFacInUnderwriting"} {
		harus(t, wb, "IsUW", kasusUji{workbasket: wb}, true)
	}
	harus(t, "tanpa workbasket", "IsUW", kasusUji{workbasket: ""}, false)
}

// TestK046_IsClaim_UjiPrefiks_vs_UjiHalaman - EDM menguji keberadaan halaman
// `pyWorkPage.ClaimData`, bukan prefiks ID kasus.
func TestK046_IsClaim_UjiPrefiks_vs_UjiHalaman(t *testing.T) {
	ada := KasusEDM{HalamanBernilai: map[string]bool{JalurHalamanKlaim: true}}
	tidak := KasusEDM{
		Properti:        map[string]string{"pyWorkPage.pyWorkIDPrefix": "CLM-"},
		HalamanBernilai: map[string]bool{JalurHalamanKlaim: false},
	}
	harus(t, "halaman klaim ada", "IsClaim", ada, true)
	harus(t, "prefiks CLM- tanpa halaman", "IsClaim", tidak, false)
	if _, err := Eval("IsClaim", KasusEDM{}); err == nil {
		t.Error("PropertyHasValue tak terjawab: tanpa galat")
	}
	panikMemuat(t, "Kasus tanpa PemeriksaNilai", "PemeriksaNilai", func() { _, _ = Eval("IsClaim", kasusUji{}) })
}

// TestK046_IsTravel_LabelTidakDieksekusi - dua sumber (relatif `.Quotation` dan
// `pyWorkPage.OfferFacIn.QuotationData`); label tampilan `Kode Bisnis = "77"`
// tidak dieksekusi.
func TestK046_IsTravel_LabelTidakDieksekusi(t *testing.T) {
	harus(t, "sumber relatif", "IsTravel", kasusUji{".Quotation.BusinessType": "Travel"}, true)
	harus(t, "sumber OfferFacIn", "IsTravel", kasusUji{"pyWorkPage.OfferFacIn.QuotationData.BusinessType": "Travel"}, true)
	harus(t, "sumber NB tidak dibaca", "IsTravel", kasusUji{"pyWorkPage.Quotation.BusinessType": "Travel"}, false)
	harus(t, "label kode bisnis", "IsTravel", kasusUji{"pyWorkPage.Quotation.BusinessCode": "77", ".Quotation.BusinessCode": "77"}, false)
}

// TestIsNotEDMSamaSamaSalah - kedua halaman sengaja tidak sinkron ke arah lain:
// OfferFacIn bukan EDM, Quotation EDM → keduanya salah.
func TestIsNotEDMSamaSamaSalah(t *testing.T) {
	k := kasusUji{"pyWorkPage.OfferFacIn.QuotationData.StatusBusiness": "1", "pyWorkPage.Quotation.StatusBusiness": "3"}
	harus(t, "IsEDM", "IsEDM", k, false)
	harus(t, "IsNotEDM", "IsNotEDM", k, false)
}
