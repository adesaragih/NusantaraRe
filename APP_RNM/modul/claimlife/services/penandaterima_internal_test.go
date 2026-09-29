package services

import (
	"strings"
	"testing"

	"nusantarare/modul/claimlife/models"
)

// TestIsiPenandaTerimaKlaim - butir bk: dihitung saat baca; tak terhitung DINYATAKAN.
func TestIsiPenandaTerimaKlaim(t *testing.T) {
	peserta := []models.Peserta{
		{TanggalKejadian: "2026-01-01 00:00:00", TanggalKlaimTeks: models.TanggalKlaimTeks{TanggalTerimaKlaim: "2026-01-31 00:00:00"}},
		{TanggalKejadian: "2026-01-01 00:00:00", TanggalKlaimTeks: models.TanggalKlaimTeks{TanggalTerimaKlaim: "2026-01-10 00:00:00"}},
	}
	isiPenandaTerimaKlaim(peserta, "20", "")
	if peserta[0].PenandaTerimaKlaim != "31/01/2026" || peserta[1].PenandaTerimaKlaim != "" {
		t.Errorf("penanda = %q, %q", peserta[0].PenandaTerimaKlaim, peserta[1].PenandaTerimaKlaim)
	}
	for _, p := range peserta {
		if p.PenandaTerimaKlaimAlasan != "" {
			t.Errorf("alasan terisi pada penanda terhitung: %q", p.PenandaTerimaKlaimAlasan)
		}
	}

	tak := []models.Peserta{{TanggalKejadian: "2026-01-01 00:00:00", TanggalKlaimTeks: models.TanggalKlaimTeks{TanggalTerimaKlaim: "2026-01-31 00:00:00"}}}
	isiPenandaTerimaKlaim(tak, "", "polis belum ada")
	if tak[0].PenandaTerimaKlaim != "" || tak[0].PenandaTerimaKlaimAlasan != "polis belum ada" {
		t.Errorf("tak terhitung: %+v", tak[0])
	}

	kosong := []models.Peserta{{TanggalKejadian: "2026-01-01 00:00:00", TanggalKlaimTeks: models.TanggalKlaimTeks{TanggalTerimaKlaim: "2026-01-31 00:00:00"}}}
	isiPenandaTerimaKlaim(kosong, "", "")
	if !strings.Contains(kosong[0].PenandaTerimaKlaimAlasan, "ambang") {
		t.Errorf("ambang kosong tidak dinyatakan: %q", kosong[0].PenandaTerimaKlaimAlasan)
	}
}
