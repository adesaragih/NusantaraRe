package services

// Outbox keputusan Komite - tiket 06. TANPA Oracle.

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"nusantarare/internal/models"
)

func jenisEfek(a models.AkibatKeputusanKomite) []string {
	var j []string
	for _, e := range EfekKeputusanKomite(a) {
		j = append(j, e.Jenis)
	}
	return j
}

// TestEfekMenurutLangkah10Sampai12 - urutan dan gerbangnya VERBATIM.
func TestEfekMenurutLangkah10Sampai12(t *testing.T) {
	akhir := jenisEfek(models.AkibatKeputusanKomite{AkseptasiAkhir: true})
	if !reflect.DeepEqual(akhir, []string{JenisEfekKomiteArasapas, JenisEfekKomiteEmail, JenisEfekKomiteKasir}) {
		t.Errorf("Setuju akhir: %v", akhir)
	}
	// ⚠️ Ralat tiket 06: langkah 11 (email) tanpa gerbang tingkat.
	for _, a := range []models.AkibatKeputusanKomite{{Berlanjut: true}, {TolakAkhir: true}, {}} {
		if got := jenisEfek(a); !reflect.DeepEqual(got, []string{JenisEfekKomiteEmail}) {
			t.Errorf("%+v: %v, mau email saja", a, got)
		}
	}
}

// TestKunciKasirVERBATIM - dibaca dari activity-nya.
func TestKunciKasirVERBATIM(t *testing.T) {
	const letak = `D:\XML\RNM_BRD\Komite Claim Life\Activity\HitServiceToKasirKMTLife_Act.xml`
	isi, err := os.ReadFile(letak)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v)", err)
	}
	for _, s := range []string{`<Kategori_1>"Kasir"</Kategori_1>`, `<Kategori_2>"insertAllPaymentKasir"</Kategori_2>`} {
		if !strings.Contains(string(isi), s) {
			t.Errorf("HitServiceToKasirKMTLife_Act tidak memuat %s", s)
		}
	}
	if KunciKasirKomite.Kategori1 != "Kasir" || KunciKasirKomite.Kategori2 != "insertAllPaymentKasir" {
		t.Errorf("kunci Kasir %+v", KunciKasirKomite)
	}
}

// TestMuatanOutboxKomiteTanpaURL - ADR-0013.
func TestMuatanOutboxKomiteTanpaURL(t *testing.T) {
	b, err := json.Marshal(muatanOutboxKomite{KasusID: "KMTLF-UJI", Kategori1: "Kasir", Kategori2: "insertAllPaymentKasir"})
	if err != nil {
		t.Fatal(err)
	}
	// Kunci JSON-nya ditagih tepat - tidak ada medan alamat, URL, atau email.
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	for k := range m {
		switch k {
		case "kasus_id", "klaim_id", "adjustment_id", "akun_id", "waktu", "kategori_1", "kategori_2":
		default:
			t.Errorf("muatan outbox memuat medan %q", k)
		}
	}
	if !strings.Contains(string(b), `"kategori_1":"Kasir"`) {
		t.Errorf("muatan outbox tanpa kunci kategori: %s", b)
	}
}

// TestRujukanEmailPerTingkat - temuan /code-review: email tingkat 2 tidak
// boleh dianggap dobel email tingkat 1.
func TestRujukanEmailPerTingkat(t *testing.T) {
	if RujukanEfekKomite("KMTLF-000001", JenisEfekKomiteEmail, 1) ==
		RujukanEfekKomite("KMTLF-000001", JenisEfekKomiteEmail, 2) {
		t.Error("email dua tingkat berujukan sama; anti-dobel akan menelan yang kedua")
	}
	if RujukanEfekKomite("KMTLF-000001", JenisEfekKomiteKasir, 3) != "KMTLF-000001" {
		t.Error("kasir berujukan kasus berubah bentuk")
	}
	if n := len(RujukanEfekKomite("KMTLF-1234567890", JenisEfekKomiteEmail, 99)); n > 40 {
		t.Errorf("rujukan %d karakter melebihi RUJUKAN VARCHAR2(40)", n)
	}
}
