package repository

import (
	"database/sql"
	"reflect"
	"testing"
)

// barisTiruan - barisKueri tanpa basis data: tiap baris [JABATAN, JABATAN_ATASAN] (nil = NULL).
type barisTiruan struct {
	isi [][2]*string
	i   int
}

func (b *barisTiruan) Next() bool { b.i++; return b.i <= len(b.isi) }
func (b *barisTiruan) Err() error { return nil }
func (b *barisTiruan) Scan(dest ...any) error {
	for c, d := range dest {
		v := b.isi[b.i-1][c]
		*d.(*sql.NullString) = sql.NullString{}
		if v != nil {
			*d.(*sql.NullString) = sql.NullString{String: *v, Valid: true}
		}
	}
	return nil
}

// TestOpsiBandingBarisPertama - bug DEV 10-10-2026 ("Banding to hanya KADIVTEKNIK, harusnya juga atasannya"): GetAksepBanding
// langkah 5 = [JABATAN, JABATAN_ATASAN] BARIS PERTAMA; baris ke-2 dst. diabaikan; NULL -> ""; nol baris -> kosong.
func TestOpsiBandingBarisPertama(t *testing.T) {
	s := func(x string) *string { return &x }
	for nama, k := range map[string]struct {
		isi [][2]*string
		mau []string
	}{
		"dua baris":   {[][2]*string{{s("KADIVTEKNIK"), s("DIREKTURTEKNIK")}, {s("SENIORUW"), s("KADIVTEKNIK")}}, []string{"KADIVTEKNIK", "DIREKTURTEKNIK"}},
		"atasan NULL": {[][2]*string{{s("KADIVTEKNIK"), nil}}, []string{"KADIVTEKNIK", ""}},
		"nol baris":   {nil, []string{}},
	} {
		got, err := opsiBandingBarisPertama(&barisTiruan{isi: k.isi})
		if err != nil || !reflect.DeepEqual(got, k.mau) {
			t.Errorf("%s: %v %v", nama, got, err)
		}
	}
}
