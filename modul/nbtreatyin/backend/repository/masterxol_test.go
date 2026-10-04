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

	"nusantarare/modul/nbtreatyin/backend/models"
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
	for _, k := range []string{"pxObjClass", "ID", "RIOGR"} {
		if _, ada := m.Nilai[k]; ada {
			t.Errorf("%s bukan medan K8 - tidak boleh dibaca", k)
		}
	}
	for _, k := range []string{"CommentList", "ShareReins", "ShareReins(1).ReinsuranceListTONP"} {
		if _, ada := m.Daftar[k]; ada {
			t.Errorf("%s bukan medan K8 (treaty keluar / komentar) - tidak boleh dibaca", k)
		}
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
	baris := func(anggota []string, anak map[string][]string) map[string]any {
		r := map[string]any{"UjiAnggotaLuar": "UJI-LUAR", "UjiAnakLuar": []any{map[string]any{"A": "UJI-LUAR"}}}
		for _, a := range anggota {
			r[a] = "UJI-" + a
		}
		for nama, ang := range anak {
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
		return r
	}
	dok := map[string]any{
		"UjiSkalarLuar": "UJI-LUAR", "FlagPPH": "1", "TypeTax": "Inclusive", "pxObjClass": "UJI-KELAS",
		"UjiDaftarLuar": []any{map[string]any{"A": "UJI-LUAR"}},
		"TreatyXOLList": []any{map[string]any{"GrossPremi": "1", "ValueList": []any{map[string]any{"Layer": "1"}}}},
		"CommentList":   []any{map[string]any{"Suggest": "UJI-LUAR"}},
	}
	for _, s := range models.SkalarMasterXOL {
		dok[s] = "UJI-" + s
	}
	for nama, sk := range models.DaftarMasterXOL {
		dok[nama] = []any{baris(sk.Anggota, sk.Anak), baris(sk.Anggota, sk.Anak)}
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
