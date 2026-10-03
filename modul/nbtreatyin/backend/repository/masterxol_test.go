package repository

// Uji pengurai dokumen master XOL (tanpa Oracle): hanya medan master K8 - yang
// dibaca rule terjangkau jalur NonProp (`models.SkalarMasterXOL`,
// `models.DaftarMasterXOL`) - yang sampai ke halaman;
// angka JSON dibaca sebagai TEKS (nol float); tanggal Pega dinormalkan; baris
// kedua (M_TREATY_IN_EDM) menimpa kunci tingkat atas baris pertama, sama
// dengan `adoptJSONObject` berurutan. Fixture fiktif UJI-.

import (
	"errors"
	"testing"
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
