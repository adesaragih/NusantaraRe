package repository

import (
	"context"
	"strings"
	"testing"

	"nusantarare/modul/nbfacin/backend/models"
)

// TestSisipAnekaBerkelompokDari - tiket 73 P6: SEQ_NO lokasi = awal+1.. (unggah CSV berkelompok menyambung lokasi yang ada);
// awal 0 = perilaku GantiObjekAneka (1..n).
func TestSisipAnekaBerkelompokDari(t *testing.T) {
	lokasi := []models.LokasiAneka{{ObjectNo: "8"}, {ObjectNo: "9", Okupasi: []models.OkupasiAneka{{Objek: []models.ObjekAneka{{ObjectName: "A"}}}}}}
	for awal, mau := range map[int][]any{0: {1, 2}, 7: {8, 9}} {
		e := &eksekutorTiruan{u: &urutTiruan{n: map[string]int{}}, seragam: true}
		baru, _, err := sisipAnekaBerkelompokDari(context.Background(), e, tabelAnekaUji(), "K1", lokasi, uidUrut(), awal)
		if err != nil {
			t.Fatal(err)
		}
		var seq []any
		for _, b := range e.baris {
			if strings.Fields(b.q)[2] == "S.LOC" {
				seq = append(seq, b.arg[2])
			}
		}
		if len(seq) != 2 || seq[0] != mau[0] || seq[1] != mau[1] {
			t.Errorf("awal %d: SEQ_NO %v", awal, seq)
		}
		if _, ada := baru[kunciAneka(mau[1].(int), 1, 1)]; !ada {
			t.Errorf("awal %d: posisi objek %v", awal, baru)
		}
	}
	if q := sqlSeqLokasiTerakhir("S.LOC"); bindTerbesar(t, q) != 1 {
		t.Error(q)
	}
}
