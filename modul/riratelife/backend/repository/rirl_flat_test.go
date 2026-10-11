package repository

// RALAT R7 (10-10-2026): penulis rate mengikat nilai kosong sebagai NULL (dulu "kunci JSON dibuang") lewat
// db.KosongJadiNil - pemeriksaan sumber, karena Gudang butuh Oracle nyata (uji nyata: TestDBRateKolomFlatTulis, tag db).

import (
	"os"
	"strings"
	"testing"
)

func badanFungsi(t *testing.T, isi, kepala string) string {
	t.Helper()
	i := strings.Index(isi, kepala)
	if i < 0 {
		t.Fatalf("fungsi %q tidak ditemukan", kepala)
	}
	j := strings.Index(isi[i:], "\n}\n")
	if j < 0 {
		t.Fatalf("akhir fungsi %q tidak ditemukan", kepala)
	}
	return isi[i : i+j]
}

func TestTulisRateKosongJadiNull(t *testing.T) {
	b, err := os.ReadFile("rirl.go")
	if err != nil {
		t.Fatal(err)
	}
	isi := strings.ReplaceAll(string(b), "\r\n", "\n")
	for kepala, ikat := range map[string][]string{
		"func (g *Gudang) SisipRate(":    {"k(r.IDUsedBy), k(r.UsedBy)", "k(r.Gender), k(r.Contract), k(r.Age), k(r.Rate)"},
		"func (g *Gudang) UbahRate(":     {"k(r.Gender), k(r.Contract), k(r.Age)", "k(r.Rate), r.ID"},
		"func (g *Gudang) UbahNamaRate(": {"db.KosongJadiNil(nama), idUsedBy"},
	} {
		f := badanFungsi(t, isi, kepala)
		for _, x := range ikat {
			if !strings.Contains(f, x) {
				t.Errorf("%s tidak mengikat %q (kosong harus NULL)", kepala, x)
			}
		}
		if strings.Contains(strings.ToUpper(f), "JSON") || strings.Contains(f, "ViewRate") {
			t.Errorf("%s masih memakai JSON / view", kepala)
		}
	}
	if _, err := os.Stat("rirl_json.go"); !os.IsNotExist(err) {
		t.Error("rirl_json.go (baca-ubah-tulis JSONDATA) harus sudah dihapus")
	}
}
