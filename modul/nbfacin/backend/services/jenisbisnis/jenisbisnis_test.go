package jenisbisnis

import "testing"

// TestTabelFirstMatch - 36 baris; untuk setiap baris dan setiap OldID-nya, Dari mengembalikan hasil baris pertama yang cocok
// (first match: hasil baris itu sendiri atau baris lebih awal ber-GroupPanel sama); baris tanpa OldID cocok untuk OldID apa pun.
func TestTabelFirstMatch(t *testing.T) {
	if JumlahBaris() != 36 {
		t.Fatalf("baris %d", JumlahBaris())
	}
	pertama := func(gp, old string) string {
		for _, b := range tabel {
			if b.groupPanel != gp {
				continue
			}
			if len(b.oldID) == 0 {
				return b.hasil
			}
			for _, o := range b.oldID {
				if o == old {
					return b.hasil
				}
			}
		}
		return Bawaan
	}
	n := 0
	for _, b := range tabel {
		kunci := b.oldID
		if len(kunci) == 0 {
			kunci = []string{"", "zz"}
		}
		for _, o := range kunci {
			n++
			if got := Dari(b.groupPanel, o); got != pertama(b.groupPanel, o) {
				t.Errorf("%s/%s: %q", b.groupPanel, o, got)
			}
		}
	}
	if n < 36 {
		t.Errorf("kunci diperiksa %d", n)
	}
}
