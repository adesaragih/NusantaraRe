package services

import (
	"errors"
	"strings"
	"testing"
)

// START DATE sesudah END DATE ditolak; sama hari diterima (04-10-2026).
func TestTahunMenolakTanggalMundur(t *testing.T) {
	m := TahunMasuk{UnderwritingYear: "2026", TreatyYear: "2026", StartDate: "2026-12-31", EndDate: "2026-01-01"}
	if _, err := m.keModel(); !errors.Is(err, ErrMasukanTidakSah) || !strings.Contains(err.Error(), PesanTanggalMundur) {
		t.Errorf("tanggal mundur: galat %v", err)
	}
	for _, akhir := range []string{"2026-12-31", "2026-01-01"} {
		m.StartDate, m.EndDate = "2026-01-01", akhir
		if _, err := m.keModel(); err != nil {
			t.Errorf("start 2026-01-01 end %s ditolak: %v", akhir, err)
		}
	}
}
