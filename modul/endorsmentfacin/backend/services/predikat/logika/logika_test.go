package logika

import (
	"reflect"
	"testing"
)

// TestUraiMenolakYangBelumTerverifikasi - AND dan OR tanpa kurung pada satu
// tingkat (urutan prioritasnya di Pega belum terverifikasi), kata NOT (tidak ada
// di korpus NB), dan logika yang tidak utuh.
func TestUraiMenolakYangBelumTerverifikasi(t *testing.T) {
	for _, s := range []string{"A AND B OR C", "A && B || C", "NOT", "NOT A", "A AND NOT", "A OR", "(A OR B", "A B", ""} {
		if _, err := Urai(s); err == nil {
			t.Errorf("%q: terurai, mau ditolak", s)
		}
	}
}

// TestNilai - tabel kebenaran kecil; tiap baris membedakan satu kesalahan
// pengurai (AND/OR tertukar, negasi hilang, kurung diabaikan).
func TestNilai(t *testing.T) {
	for _, u := range []struct {
		logika string
		hasil  map[string]bool
		mau    bool
	}{
		{"A OR B", map[string]bool{"A": false, "B": true}, true},
		{"A AND B", map[string]bool{"A": false, "B": true}, false},
		{"A && B", map[string]bool{"A": true, "B": true}, true},
		{"A Or B", map[string]bool{"A": false, "B": false}, false},
		{"!A AND !B", map[string]bool{"A": false, "B": false}, true},
		{"!A AND !B", map[string]bool{"A": true, "B": false}, false},
		{"!(A OR B)", map[string]bool{"A": false, "B": true}, false},
		{"A AND (B OR C)", map[string]bool{"A": true, "B": false, "C": true}, true},
		{"A AND (B OR C)", map[string]bool{"A": false, "B": true, "C": true}, false},
	} {
		s, err := Urai(u.logika)
		if err != nil {
			t.Fatalf("%q: %v", u.logika, err)
		}
		if got := s.Nilai(u.hasil); got != u.mau {
			t.Errorf("%q atas %v = %v, mau %v", u.logika, u.hasil, got, u.mau)
		}
	}
}

// TestLabelUrutPega - A, B, …, Z, lalu AA: pendek dulu, lalu abjad.
func TestLabelUrutPega(t *testing.T) {
	s, err := Urai("AA OR B OR A OR Z")
	if err != nil {
		t.Fatal(err)
	}
	if got, mau := s.Label(), []string{"A", "B", "Z", "AA"}; !reflect.DeepEqual(got, mau) {
		t.Errorf("label %v, mau %v", got, mau)
	}
}
