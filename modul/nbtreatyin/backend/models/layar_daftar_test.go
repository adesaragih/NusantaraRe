package models

// Uji audit silang putaran 3 (bab 3.2 butir 1, tombol Add/Delete grid):
// daftar baris yang diterima dari layar mengikuti grid yang dapat disunting
// di XML.
//
//	DetailPolicyTreatyIn (Proporsional, wadah `.IsNewPolicyNonProp != 1`):
//	  grid `.SpreadingRiskList` bertombol Add/Delete; grid `.ListInstallment`
//	  bersel `.InstallmentPercentage` / `.Premium` terbuka
//	Section/SpreadingRiskList (di DetailPolicyTreatyInNonProportional):
//	  Add/Delete dan sel terkunci bila `pyWorkPage.TreatyIn.FacultativeShare > 0`
//	InstallmentList NonProp: masterDetail ber-pyReadOnly - tidak diterima
//	DetailDeptHeadTreatyIn_UW: seluruh grid hanya-baca

import "testing"

func TestGabungMasukanDaftarMenurutGridXML(t *testing.T) {
	kiriman := func() *Halaman {
		m := HalamanBaru()
		m.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "UJI-1"}, {"TreatyType": "UJI-2"}})
		m.SetelDaftar(DaftarAngsuran, []Baris{{"InstallmentNo": "1", "InstallmentPercentage": "100"}})
		return m
	}
	for _, tt := range []struct {
		nama, posisi, nonProp, fakultatif string
		spreading, angsuran               int
	}{
		{"admin Proporsional", PosisiAdmin, "", "", 2, 1},
		{"admin NonProp FacultativeShare 0", PosisiAdmin, "1", "0", 2, 0},
		{"admin NonProp FacultativeShare kosong", PosisiAdmin, "1", "", 2, 0},
		{"admin NonProp FacultativeShare 5", PosisiAdmin, "1", "5", 0, 0},
		{"Sec Head", PosisiSecHead, "", "", 0, 0},
		{"Dept Head", PosisiDeptHead, "", "", 0, 0},
	} {
		h := HalamanBaru()
		h.Setel("PolicyTreatyIn.IsNewPolicyNonProp", tt.nonProp)
		h.Setel("TreatyIn.FacultativeShare", tt.fakultatif)
		GabungMasukanLayar(h, kiriman(), tt.posisi, nil)
		if got := len(h.AmbilDaftar(DaftarSpreading)); got != tt.spreading {
			t.Errorf("%s: SpreadingRiskList %d baris, harap %d", tt.nama, got, tt.spreading)
		}
		if got := len(h.AmbilDaftar(DaftarAngsuran)); got != tt.angsuran {
			t.Errorf("%s: ListInstallment %d baris, harap %d", tt.nama, got, tt.angsuran)
		}
	}
}
