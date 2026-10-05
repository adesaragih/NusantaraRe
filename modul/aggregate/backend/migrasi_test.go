package backend

// Migrasi 880 dan slot menu 996 - TANPA Oracle: bentuk SQL-nya.

import (
	"reflect"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

func langkah(t *testing.T, nama string) []string {
	t.Helper()
	p, err := migrasi.PernyataanLangkah(berkasMigrasi, nama)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// 880 adalah blok sequence-dari-kueri yang membaca nomor AGG terbesar AGGREGATE; mundurnya membuang sequence.
func TestSequenceAggregate(t *testing.T) {
	p := langkah(t, "880_seq_aggregate.sql")
	if len(p) != 1 {
		t.Fatalf("880: %d pernyataan", len(p))
	}
	s, ok := migrasi.BacaSequenceDariKueri(p[0])
	if !ok || s.Nama != "SEQ_AGGREGATE" || !strings.Contains(s.Kueri, "'^AGG-([0-9]+)$'") ||
		!strings.HasSuffix(s.Kueri, "FROM {skema}.AGGREGATE") {
		t.Errorf("880 terbaca %+v %v", s, ok)
	}
	if p := langkah(t, "880_seq_aggregate_down.sql"); !reflect.DeepEqual(p, []string{"DROP SEQUENCE {skema}.SEQ_AGGREGATE"}) {
		t.Errorf("880_down: %q", p)
	}
}

// Slot menu 996: satu UPDATE DIMIGRASI baris modul ini.
func TestSlotMenu(t *testing.T) {
	if p := langkah(t, "996_menu_aggregate.sql"); len(p) != 1 || !strings.Contains(p[0], "SET DIMIGRASI = '1'") ||
		!strings.HasSuffix(p[0], "WHERE KODE = 'aggregate'") {
		t.Errorf("996: %q", p)
	}
}
