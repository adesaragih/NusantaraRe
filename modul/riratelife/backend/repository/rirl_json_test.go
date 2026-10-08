package repository

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// RALAT R3: kunci yang diubah diganti, kosong = dibuang; kunci Pega lain dan bentuk angka tetap.
func TestTerapkanKunci(t *testing.T) {
	lama := `{"IDUSEDBY":101,"USEDBY":"UJI A","GENDER":"U","AGE":"30","pzInsKey":"UJI-<X>&Y","NILAI":1.50}`
	baru, err := TerapkanKunci(lama, map[string]string{"GENDER": "M", "AGE": "", "RATE": "0,5"})
	if err != nil {
		t.Fatal(err)
	}
	for _, mau := range []string{`"IDUSEDBY":101`, `"USEDBY":"UJI A"`, `"GENDER":"M"`, `"RATE":"0,5"`, `"pzInsKey":"UJI-<X>&Y"`, `"NILAI":1.50`} {
		if !strings.Contains(baru, mau) {
			t.Errorf("tanpa %s: %s", mau, baru)
		}
	}
	if strings.Contains(baru, `"AGE"`) {
		t.Errorf("AGE kosong harus dibuang: %s", baru)
	}
	if b, err := TerapkanKunci("  ", map[string]string{"USEDBY": "UJI"}); err != nil || b != `{"USEDBY":"UJI"}` {
		t.Errorf("JSONDATA kosong %q %v", b, err)
	}
	for _, rusak := range []string{"[1]", "bukan json", "null"} {
		if _, err := TerapkanKunci(rusak, map[string]string{"A": "1"}); !errors.Is(err, ErrJSONRusak) {
			t.Errorf("%q: %v", rusak, err)
		}
	}
	if TeksKunci(lama, "IDUSEDBY") != "101" || TeksKunci(`{"IDUSEDBY":" 101 "}`, "IDUSEDBY") != "101" || TeksKunci(lama, "TIDAK") != "" {
		t.Error("TeksKunci")
	}
}

// Oracle DEV menolak JSON_MERGEPATCH (ORA-00907): nol pemakaian di SQL modul ini.
func TestNolJSONMergepatchDiSQL(t *testing.T) {
	berkas, _ := filepath.Glob("*.go")
	for _, b := range berkas {
		if strings.HasSuffix(b, "_test.go") {
			continue
		}
		isi, err := os.ReadFile(b)
		if err != nil {
			t.Fatal(err)
		}
		for i, baris := range strings.Split(string(isi), "\n") {
			if strings.Contains(baris, "JSON_MERGEPATCH") && !strings.HasPrefix(strings.TrimSpace(baris), "//") {
				t.Errorf("%s:%d memakai JSON_MERGEPATCH", b, i+1)
			}
		}
	}
}
