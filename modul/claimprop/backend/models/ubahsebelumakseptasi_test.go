package models_test

import (
	"testing"

	"nusantarare/modul/claimprop/backend/models"
)

// medanTampil - medan pertama berjalur `j` yang tampil di pohon tata (nil = tersembunyi).
func medanTampil(ts []models.Tata, j string) *models.Tata {
	if m := medanBerjalur(ts, j); len(m) > 0 {
		return &m[0]
	}
	return nil
}

// Keputusan work owner 09-10-2026 "buat boleh di ubah kalau belum ada akseptasi, yang isanyacceptation itu": di layar
// Input Acceptation ketujuh medan informasi klaim terbuka selama IsAnyAcceptation bukan 1 (XML menguncinya dengan
// IsOutstanding, yang selalu 1 sesudah Save Outstanding); sesudah ada akseptasi terkunci. Layar Outstanding tetap.
func TestMedanKlaimBolehDiubahSebelumAkseptasi(t *testing.T) {
	medan := []string{models.CD + "ConsultantID", models.CD + "AppointedADJID", models.CD + "ReportDescription",
		models.CD + "Location", models.CD + "Occupation", models.CD + "Province", models.CD + "PostalCode"}
	h := models.HalamanBaru()
	h.Setel("IsOutstanding", "1")
	h.Setel("IsAnyAcceptation", "0")
	ts := models.Evaluasi(h, models.LayarAkseptasi(), false)
	for _, j := range medan {
		m := medanTampil(ts, j)
		if m == nil || m.HanyaBaca || m.Nonaktif {
			t.Errorf("belum ada akseptasi: %s harus tampil dan dapat diubah, dapat %+v", j, m)
		}
	}
	for _, id := range []string{"AdjusterConsultantBaru1", "AdjusterConsultantBaru2"} {
		if b := cariTata(ts, id); b == nil || b.Nonaktif {
			t.Errorf("belum ada akseptasi: tombol + %s harus aktif, dapat %+v", id, b)
		}
	}

	h.Setel("IsAnyAcceptation", "1")
	ts = models.Evaluasi(h, models.LayarAkseptasi(), false)
	for _, j := range medan {
		if m := medanTampil(ts, j); m != nil && !m.HanyaBaca {
			t.Errorf("sudah ada akseptasi: %s harus terkunci", j)
		}
	}

	ts = models.Evaluasi(h, models.LayarOutstanding(), false)
	for _, j := range []string{models.CD + "ReportDescription", models.CD + "Location", models.CD + "Occupation"} {
		if m := medanTampil(ts, j); m == nil || !m.HanyaBaca {
			t.Errorf("layar Outstanding ber-IsOutstanding 1: %s tetap terkunci seperti XML", j)
		}
	}
}
