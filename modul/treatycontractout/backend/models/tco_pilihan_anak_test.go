package models

import (
	"html"
	"os"
	"regexp"
	"strings"
	"testing"
)

// akarKorpusTCO - korpus READ-ONLY; uji yang membacanya DILEWATI bila tak terjangkau.
const akarKorpusTCO = `D:\XML\RNM_BRD\Treaty Contract Out`

func ids(d []PilihanReinsAnak) string {
	var s []string
	for _, p := range d {
		s = append(s, p.ID+"="+p.Nama)
	}
	return strings.Join(s, ",")
}

func TestPilihanReinsAnakDari(t *testing.T) {
	for nama, mau := range map[string]string{
		"2019 QS 101M TRT": "10028=QS (OR),10004=QS (R/I),10007=ORS",
		"UJI SPL 2020":     "10248=SPL (OR),10249=SPL (RI),10007=ORS",
		"UJI XOL LAYER":    "10028=QS (OR),10004=QS (R/I),10217=XL",
		"ORS":              "10007=ORS",
		"UJI ORS 2020":     "10007=ORS",
		"UJI SURPLUS":      "",
		"QS 2020":          "",                         // tanpa spasi di depan - " QS " tidak cocok
		"UJI qs 2020":      "",                         // @contains peka huruf besar-kecil
		"UJI QS ORS SHARE": "10007=ORS,10004=QS (R/I)", // 2.4 menimpa pxResults(1); ORS di (3) tampil sekali
		"UJI SPL XOL BOTH": "10028=QS (OR),10004=QS (R/I),10217=XL",
		"":                 "",
	} {
		if g := ids(PilihanReinsAnakDari(nama)); g != mau {
			t.Errorf("%q: %s, mau %s", nama, g, mau)
		}
	}
}

// TestPilihanReinsAnakDariKorpus - tabel langkah diturunkan ULANG dari
// `Activity/TreatyContractSetReinsTypeList.xml`: kata `@contains`, urutan, ID
// `.CARI1`, dan nama `.CARI2` per indeks.
func TestPilihanReinsAnakDariKorpus(t *testing.T) {
	b, err := os.ReadFile(akarKorpusTCO + `\Activity\TreatyContractSetReinsTypeList.xml`)
	if err != nil {
		t.Skipf("korpus tidak terjangkau (%v)", err)
	}
	isi := string(b)
	kata := regexp.MustCompile(`@contains\(InputData\.CARIDESCFACIN,(&quot;|")([^"&]*)(&quot;|")\)`)
	pasangan := regexp.MustCompile(`<PropertiesName>ReinsTypeList\.pxResults\((\d)\)\.(CARI[12])</PropertiesName>\s*<PropertiesValue>([^<]*)</PropertiesValue>`)
	var kataXML []string
	for _, m := range kata.FindAllStringSubmatch(isi, -1) {
		kataXML = append(kataXML, m[2])
	}
	if strings.Join(kataXML, "|") != " QS | SPL | XOL |ORS" {
		t.Fatalf("kata @contains korpus %q", kataXML)
	}
	// Potong per langkah 2.1-2.4 pada deskripsinya, lalu baca pasangan CARI1/CARI2.
	desk := []string{"JIKA QUOTA SHARE", "JIKA SURPLUS", "JIKA XOL", "JIKA ORS"}
	for i, d := range desk {
		awal := strings.Index(isi, "<pyStepsDescription>"+d+"</pyStepsDescription>")
		if awal < 0 {
			t.Fatalf("langkah %q tidak ada", d)
		}
		potong := isi[awal:]
		if i+1 < len(desk) {
			potong = potong[:strings.Index(potong, "<pyStepsDescription>"+desk[i+1]+"</pyStepsDescription>")]
		}
		baris := map[string]*PilihanReinsAnak{}
		var urut []string
		for _, m := range pasangan.FindAllStringSubmatch(potong, -1) {
			if baris[m[1]] == nil {
				baris[m[1]] = &PilihanReinsAnak{}
				urut = append(urut, m[1])
			}
			v := strings.Trim(html.UnescapeString(m[3]), `"`)
			if m[2] == "CARI1" {
				baris[m[1]].ID = v
			} else {
				baris[m[1]].Nama = v
			}
		}
		var korpus []PilihanReinsAnak
		for _, n := range urut {
			korpus = append(korpus, *baris[n])
		}
		if ids(korpus) != ids(langkahPilihanReinsAnak[i].isi) {
			t.Errorf("langkah %s: kode %s, korpus %s", d, ids(langkahPilihanReinsAnak[i].isi), ids(korpus))
		}
		if langkahPilihanReinsAnak[i].kata != kataXML[i] {
			t.Errorf("langkah %s: kata kode %q, korpus %q", d, langkahPilihanReinsAnak[i].kata, kataXML[i])
		}
	}
}
