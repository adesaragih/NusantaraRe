package repository

// Uji pengurai dokumen master XOL (tanpa Oracle): hanya medan master K8 - yang
// dibaca rule terjangkau jalur NonProp (`models.SkalarMasterXOL`,
// `models.DaftarMasterXOL`) - yang sampai ke halaman;
// angka JSON dibaca sebagai TEKS (nol float); tanggal Pega dinormalkan; baris
// kedua (M_TREATY_IN_EDM) menimpa kunci tingkat atas baris pertama, sama
// dengan `adoptJSONObject` berurutan. Fixture fiktif UJI-.

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/edmtreatyin/backend/models"
)

const dokumenUji = `{
  "pxObjClass": "UJI-KELAS",
  "ID": "UJI-T1",
  "RNMShare": 12.5,
  "FacultativeShare": "0",
  "EDMState": "",
  "ProportionType": "NonProportional",
  "RIOGR": "UJI-TIDAK-DIBACA",
  "CommentList": [{"Suggest": "UJI-TIDAK-DIBACA"}],
  "ShareReins": [{"ReinsuranceListTONP": [{"SpreadingTypeIDXOL": "UJI-KELUAR"}]}],
  "Share": [
    {"Layer": "1", "LayerType": "UJI-LT", "SpreadingTypeXOL": "", "Rahasia": "UJI-TIDAK-DIBACA",
     "GrossPremiumList": [{"Currency": "IDR", "Value": 1000000000000000000000.123456789}],
     "SpreadingListXOL": [{"ReinsTypeID": "UJI-R1", "Pct": 6, "Lain": "UJI-TIDAK-DIBACA"}]}
  ],
  "Installment": [
    {"Currency": "IDR", "AmountTotal": "2700", "PctTotal": 100,
     "InstallmentList": [
       {"Installment": 1, "PaymentDate": "20261101", "InstallmentPct": 40, "Amount": 1080, "Currency": "IDR"},
       {"Installment": 2, "PaymentDate": "20270131T170000.000 GMT", "InstallmentPct": 60, "Amount": 1620, "Currency": "IDR"}
     ]}
  ]
}`

func TestUraiMasterXOLHanyaMedanK8(t *testing.T) {
	m, err := uraiMasterXOL([]string{dokumenUji})
	if err != nil {
		t.Fatal(err)
	}
	if m.Nilai["RNMShare"] != "12.5" || m.Nilai["FacultativeShare"] != "0" || m.Nilai["ProportionType"] != "NonProportional" {
		t.Fatalf("skalar %v", m.Nilai)
	}
	for _, k := range []string{"pxObjClass", "RIOGR"} {
		if _, ada := m.Nilai[k]; ada {
			t.Errorf("%s bukan medan K8 - tidak boleh dibaca", k)
		}
	}
	// ID dibaca jalur EDM (CopyGeneralDataEDM_act 1: NoOffer = pyWorkPage.TreatyIn.ID).
	if m.Nilai["ID"] != "UJI-T1" {
		t.Errorf("ID %q", m.Nilai["ID"])
	}
	if _, ada := m.Daftar["CommentList"]; ada {
		t.Error("CommentList bukan medan K8 - tidak boleh dibaca")
	}
	// Treaty keluar jalur EDM XOL Retro (InsertToTreatyOutXOLList 3-4.7): baris dibaca, anggota di luar skema tidak.
	if r := m.Daftar["ShareReins(1).ReinsuranceListTONP"]; len(r) != 1 || len(r[0]) != 0 {
		t.Errorf("ReinsuranceListTONP %v - SpreadingTypeIDXOL di luar skema", r)
	}
	s := m.Daftar["Share"]
	if len(s) != 1 || s[0]["Layer"] != "1" || s[0]["LayerType"] != "UJI-LT" {
		t.Fatalf("Share %v", s)
	}
	if _, ada := s[0]["Rahasia"]; ada {
		t.Error("anggota di luar skema tidak boleh dibaca")
	}
	// angka JSON besar tetap utuh sebagai teks - tidak lewat float
	if v := m.Daftar["Share(1).GrossPremiumList"][0]["Value"]; v != "1000000000000000000000.123456789" {
		t.Errorf("Value %q", v)
	}
	if x := m.Daftar["Share(1).SpreadingListXOL"][0]; x["Pct"] != "6" || x["Lain"] != "" {
		t.Errorf("SpreadingListXOL %v", x)
	}
	r := m.Daftar["Installment(1).InstallmentList"]
	if len(r) != 2 || r[0]["Installment"] != "1" || r[0]["Amount"] != "1080" {
		t.Fatalf("InstallmentList %v", r)
	}
	// tanggal Pega: YYYYMMDD dan cap waktu GMT (17.00 GMT = 00.00 WIB hari berikut)
	if r[0]["PaymentDate"] != "2026-11-01" || r[1]["PaymentDate"] != "2027-02-01" {
		t.Errorf("PaymentDate %q %q", r[0]["PaymentDate"], r[1]["PaymentDate"])
	}
}

func TestUraiMasterXOLBarisBerikutMenimpa(t *testing.T) {
	edm := `{"RNMShare": "20", "Installment": [{"Currency": "USD"}]}`
	m, err := uraiMasterXOL([]string{dokumenUji, edm})
	if err != nil {
		t.Fatal(err)
	}
	if m.Nilai["RNMShare"] != "20" || m.Nilai["ProportionType"] != "NonProportional" {
		t.Fatalf("kunci baris kedua menimpa, kunci lain tetap: %v", m.Nilai)
	}
	if d := m.Daftar["Installment"]; len(d) != 1 || d[0]["Currency"] != "USD" {
		t.Fatalf("daftar diganti utuh: %v", d)
	}
	if _, ada := m.Daftar["Installment(1).InstallmentList"]; ada {
		t.Fatal("daftar bersarang baris lama ikut terbuang")
	}
	if len(m.Daftar["Share"]) != 1 {
		t.Fatal("daftar yang tidak disebut baris kedua tetap")
	}
}

func TestUraiMasterXOLRusak(t *testing.T) {
	if _, err := uraiMasterXOL([]string{`{"RNMShare": `}); !errors.Is(err, ErrMasterXOLRusak) {
		t.Fatalf("dokumen rusak: %v", err)
	}
}

// dokumenMaksimal - dokumen master UJI- yang memuat SETIAP medan daftar
// tertutup (`models.MedanMasterXOL`) beserta medan di luar daftar di setiap
// tingkat: skalar dan daftar akar lain, keluaran `TreatyXOLList` (F1: tidak
// dibaca), medan polis (`FlagPPH`, `TypeTax`), anggota baris lain, dan daftar
// bersarang lain.
func dokumenMaksimal(t *testing.T) string {
	t.Helper()
	var baris func(sk models.SkemaDaftarMaster) map[string]any
	baris = func(sk models.SkemaDaftarMaster) map[string]any {
		r := map[string]any{"UjiAnggotaLuar": "UJI-LUAR", "UjiAnakLuar": []any{map[string]any{"A": "UJI-LUAR"}}}
		for _, a := range sk.Anggota {
			r[a] = "UJI-" + a
		}
		for nama, ang := range sk.Anak {
			var sub []any
			for i := 0; i < 2; i++ {
				rr := map[string]any{"UjiAnggotaLuar": "UJI-LUAR"}
				for _, a := range ang {
					rr[a] = "UJI-" + a
				}
				sub = append(sub, rr)
			}
			r[nama] = sub
		}
		for nama, s := range sk.Sub { // daftar bersarang bertingkat (treaty keluar)
			r[nama] = []any{baris(s), baris(s)}
		}
		return r
	}
	dok := map[string]any{
		"UjiSkalarLuar": "UJI-LUAR", "FlagPPH": "1", "TypeTax": "Inclusive", "pxObjClass": "UJI-KELAS",
		"UjiDaftarLuar": []any{map[string]any{"A": "UJI-LUAR"}},
		"TreatyXOLList": []any{map[string]any{"GrossPremi": "1", "ValueList": []any{map[string]any{"Layer": "1"}}}},
		"CommentList":   []any{map[string]any{"Suggest": "UJI-LUAR"}},
	}
	// taruh menaruh nilai di jalur titik (`ValueDifference.Share` = objek tersemat), beserta medan luar di objeknya.
	taruh := func(jalur string, v any) {
		cur := dok
		bagian := strings.Split(jalur, ".")
		for _, b := range bagian[:len(bagian)-1] {
			o, ok := cur[b].(map[string]any)
			if !ok {
				o = map[string]any{"UjiObjekLuar": "UJI-LUAR", "UjiDaftarObjekLuar": []any{map[string]any{"A": "UJI-LUAR"}}}
				cur[b] = o
			}
			cur = o
		}
		cur[bagian[len(bagian)-1]] = v
	}
	for _, s := range models.SkalarMasterXOL {
		taruh(s, "UJI-"+s)
	}
	for nama, sk := range models.DaftarMasterXOL {
		taruh(nama, []any{baris(sk), baris(sk)})
	}
	b, err := json.Marshal(dok)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

var polaSubskrip = regexp.MustCompile(`\(\d+\)`)

// F1: keluaran pengurai - dan karenanya keluaran `MasterXOLDariJSON`
// (`TestMasterXOLDariJSONHanyaMengembalikanHasilUrai`) - HANYA memuat medan
// daftar tertutup `models.MedanMasterXOL`. Uji ini GAGAL bila satu saja medan
// di luar daftar sampai ke halaman.
func TestUraiMasterXOLHanyaMedanDaftarTertutup(t *testing.T) {
	m, err := uraiMasterXOL([]string{dokumenMaksimal(t)})
	if err != nil {
		t.Fatal(err)
	}
	boleh := map[string]bool{}
	daftarBoleh := map[string]bool{}
	for _, j := range models.MedanMasterXOL() {
		boleh[j] = true
		for i := strings.LastIndex(j, "."); i > 0; i = strings.LastIndex(j[:i], ".") {
			daftarBoleh[j[:i]] = true
		}
	}
	dapat := map[string]bool{}
	for k := range m.Nilai {
		dapat[k] = true
		if !boleh[k] {
			t.Errorf("skalar %s di luar daftar tertutup", k)
		}
	}
	for k, bs := range m.Daftar {
		jalur := polaSubskrip.ReplaceAllString(k, "")
		if !daftarBoleh[jalur] {
			t.Errorf("daftar %s di luar daftar tertutup", k)
		}
		for _, b := range bs {
			for a := range b {
				dapat[jalur+"."+a] = true
				if !boleh[jalur+"."+a] {
					t.Errorf("medan %s.%s di luar daftar tertutup", k, a)
				}
			}
		}
	}
	for j := range boleh {
		if !dapat[j] {
			t.Errorf("medan daftar %s tidak terbaca dari dokumen yang memuatnya", j)
		}
	}
}

// Pengurai umum: daftar di objek tersemat (`ValueDifference.Share`, `ActualValue.Share`) dan daftar bersarang
// bertingkat (`ShareReins(i).ReinsuranceListTONP(j).GrossPremiumList`, `SpreadingTONP(1).ReinsuranceListTONP(1)
// .LayerList(k).DeductionList`) sampai ke kunci halaman berjalur `models.JalurAnak`; tanggal di objek tersemat
// dinormalkan; objek tersemat yang tidak memuat daftar itu dilewati.
func TestUraiMasterXOLObjekTersematDanBertingkat(t *testing.T) {
	dok := `{
  "ValueDifference": {"Lain": "UJI-LUAR",
    "Share": [{"GrossPremiumList": [{"Currency": "IDR", "Value": 12.5}], "NetPremiumList": [{"Currency": "IDR", "Value": "7"}]}],
    "Installment": [{"Currency": "IDR", "AmountTotal": 300,
      "InstallmentList": [{"Installment": 1, "DueDate": "20261101", "PaymentDate": "20261201", "Amount": 100, "InstallmentPct": 33.3333}]}]},
  "ActualValue": {"FacultativeShareList": [{"DeductionList": [{"Currency": "USD", "Deduction": 3}]}]},
  "ShareReins": [
    {"ReinsuranceListTONP": [
      {"ReinsName": "UJI-SOB", "Layer": "1", "Rahasia": "UJI-TIDAK-DIBACA", "GrossPremiumList": [{"Currency": "IDR", "Value": 40}]},
      {"ReinsName": "UJI-LAIN", "Layer": "2", "DeductionList": [{"Currency": "IDR", "Deduction": 4}]}]}],
  "SpreadingTONP": [{"ReinsuranceListTONP": [{"LayerList": [{}, {"DeductionList": [{"Currency": "IDR", "Deduction": 9}]}]}]}],
  "Limits": [{"TreatyGroupList": [{"TreatyGroup": "UJI-GRUP"}]}],
  "IsProRate": true, "ProRatePercent": 50, "InstallmentNo": "3"
}`
	m, err := uraiMasterXOL([]string{dok})
	if err != nil {
		t.Fatal(err)
	}
	if v := m.Daftar["ValueDifference.Share(1).GrossPremiumList"]; len(v) != 1 || v[0]["Value"] != "12.5" {
		t.Errorf("ValueDifference.Share GrossPremiumList %v", v)
	}
	r := m.Daftar["ValueDifference.Installment(1).InstallmentList"]
	if len(r) != 1 || r[0]["DueDate"] != "2026-11-01" || r[0]["PaymentDate"] != "2026-12-01" || r[0]["InstallmentPct"] != "33.3333" {
		t.Errorf("ValueDifference.Installment InstallmentList %v", r)
	}
	if v := m.Daftar["ValueDifference.Installment"]; len(v) != 1 || v[0]["AmountTotal"] != "300" {
		t.Errorf("ValueDifference.Installment %v", v)
	}
	if _, ada := m.Daftar["ActualValue.Share"]; ada {
		t.Error("ActualValue.Share tidak ada di dokumen - tidak boleh muncul")
	}
	if v := m.Daftar["ActualValue.FacultativeShareList(1).DeductionList"]; len(v) != 1 || v[0]["Deduction"] != "3" {
		t.Errorf("ActualValue.FacultativeShareList %v", v)
	}
	tonp := m.Daftar["ShareReins(1).ReinsuranceListTONP"]
	if len(tonp) != 2 || tonp[0]["ReinsName"] != "UJI-SOB" || tonp[1]["Layer"] != "2" || tonp[0]["Rahasia"] != "" {
		t.Errorf("ReinsuranceListTONP %v", tonp)
	}
	if v := m.Daftar["ShareReins(1).ReinsuranceListTONP(1).GrossPremiumList"]; len(v) != 1 || v[0]["Value"] != "40" {
		t.Errorf("TONP(1) GrossPremiumList %v", v)
	}
	if v := m.Daftar["ShareReins(1).ReinsuranceListTONP(2).DeductionList"]; len(v) != 1 || v[0]["Deduction"] != "4" {
		t.Errorf("TONP(2) DeductionList %v", v)
	}
	if v := m.Daftar["SpreadingTONP(1).ReinsuranceListTONP(1).LayerList"]; len(v) != 2 {
		t.Errorf("LayerList %v", v)
	}
	if v := m.Daftar["SpreadingTONP(1).ReinsuranceListTONP(1).LayerList(2).DeductionList"]; len(v) != 1 || v[0]["Deduction"] != "9" {
		t.Errorf("LayerList(2) DeductionList %v", v)
	}
	if v := m.Daftar["Limits(1).TreatyGroupList"]; len(v) != 1 || v[0]["TreatyGroup"] != "UJI-GRUP" {
		t.Errorf("TreatyGroupList %v", v)
	}
	if m.Nilai["IsProRate"] != "true" || m.Nilai["ProRatePercent"] != "50" || m.Nilai["InstallmentNo"] != "3" {
		t.Errorf("skalar EDM %v", m.Nilai)
	}
}

// uraiMasterXOLPerBaris: satu master per baris hasil, berurutan, tanpa saling menimpa (adopsi berurutan di models).
func TestUraiMasterXOLPerBaris(t *testing.T) {
	ms, err := uraiMasterXOLPerBaris([]string{`{"ID": "UJI-A", "RNMShare": "10"}`, `{"ID": "UJI-B"}`})
	if err != nil {
		t.Fatal(err)
	}
	if len(ms) != 2 || ms[0].Nilai["ID"] != "UJI-A" || ms[0].Nilai["RNMShare"] != "10" || ms[1].Nilai["ID"] != "UJI-B" {
		t.Fatalf("per baris %v", ms)
	}
	if _, ada := ms[1].Nilai["RNMShare"]; ada {
		t.Error("baris kedua tidak mewarisi kunci baris pertama")
	}
	if _, err := uraiMasterXOLPerBaris([]string{`{}`, `{"ID": `}); !errors.Is(err, ErrMasterXOLRusak) {
		t.Fatalf("baris rusak: %v", err)
	}
	if ms, err := uraiMasterXOLPerBaris(nil); err != nil || len(ms) != 0 {
		t.Fatalf("nol baris: %v %v", ms, err)
	}
}
