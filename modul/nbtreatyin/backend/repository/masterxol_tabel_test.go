package repository

// Pembaca master Non-Prop dari tabel `T_TREATY_*` (`MasterXOLDariTabel`) -
// laporan pemakai 11-10-2026: master baru (mis. 1002318) tidak punya dokumen
// JSON, dan Choose Business Non-Prop gagal "kontrak ini tidak ditemukan".

import (
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

// Setiap PageList master punya tabel sumber, dan setiap kolom yang dipetakan
// adalah anggota DAFTAR TERTUTUP (`models.DaftarMasterXOL`) - pembaca tabel tidak
// boleh membawa medan di luar daftar itu.
func TestSumberMasterTabelMeliputiDaftarMaster(t *testing.T) {
	for nama, skema := range models.DaftarMasterXOL {
		s, ada := SumberDaftarMasterXOL[nama]
		if !ada {
			t.Errorf("PageList %q tidak punya tabel sumber", nama)
			continue
		}
		anggota := map[string]bool{}
		for _, a := range skema.Anggota {
			anggota[a] = true
		}
		for a := range s.kolom {
			if !anggota[a] {
				t.Errorf("%s.%s dipetakan tetapi bukan anggota daftar tertutup", nama, a)
			}
		}
		for anak, ang := range skema.Anak {
			sa, ada := s.anak[anak]
			if !ada {
				t.Errorf("%s.%s tidak punya tabel sumber", nama, anak)
				continue
			}
			boleh := map[string]bool{}
			for _, a := range ang {
				boleh[a] = true
			}
			for a := range sa.kolom {
				if !boleh[a] {
					t.Errorf("%s.%s.%s dipetakan tetapi bukan anggota daftar tertutup", nama, anak, a)
				}
			}
		}
	}
	for nama := range SumberDaftarMasterXOL {
		if _, ada := models.DaftarMasterXOL[nama]; !ada {
			t.Errorf("sumber %q bukan PageList master", nama)
		}
	}
}

// Setiap skalar master punya kolom `T_TREATY_REVISION`.
func TestSkalarMasterTabelLengkap(t *testing.T) {
	for _, k := range models.SkalarMasterXOL {
		if KolomSkalarMasterXOL[k] == "" {
			t.Errorf("skalar %q tanpa kolom revisi", k)
		}
	}
}
