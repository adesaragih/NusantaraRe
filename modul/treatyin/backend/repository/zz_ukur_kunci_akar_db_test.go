//go:build db

package repository_test

// SEMENTARA — kunci AKAR dokumen yang belum punya kolom di T_TREATY_REVISION,
// diukur atas SELURUH korpus beserta seberapa sering ia TERISI.

import (
	"sort"
	"testing"

	"nusantarare/modul/treatyin/backend/repository"
)

func TestUkurKunciAkarTanpaKolom(t *testing.T) {
	g, ctx := gudangBaca(t)
	ids, err := g.DaftarMasterID(ctx, 0)
	if err != nil {
		t.Fatal(err)
	}
	isiSkalar := map[string]int{}
	isiLarik := map[string]int{}
	var n int
	for _, id := range ids {
		teks, err := g.BacaDokumenMentah(ctx, id)
		if err != nil {
			continue
		}
		doc, err := repository.UraiDokumen(teks)
		if err != nil {
			continue
		}
		n++
		for tabel, k := range repository.KunciTakTerpetakan(doc) {
			if tabel != "T_TREATY_REVISION" {
				continue
			}
			for _, kk := range k {
				v := doc[kk]
				switch x := v.(type) {
				case string:
					if x != "" {
						isiSkalar[kk]++
					}
				case []any:
					if len(x) > 0 {
						isiLarik[kk]++
					}
				case map[string]any:
					if len(x) > 0 {
						isiSkalar[kk]++
					}
				}
			}
		}
	}
	t.Logf("dokumen disapu: %d", n)
	cetak := func(judul string, m map[string]int) {
		type kv struct {
			k string
			n int
		}
		var l []kv
		for k, v := range m {
			l = append(l, kv{k, v})
		}
		sort.Slice(l, func(i, j int) bool { return l[i].n > l[j].n })
		t.Logf("--- %s (%d kunci)", judul, len(l))
		for _, e := range l {
			t.Logf("    %-34s terisi di %d dokumen", e.k, e.n)
		}
	}
	cetak("SKALAR akar tanpa kolom", isiSkalar)
	cetak("LARIK akar tanpa tabel", isiLarik)
}
