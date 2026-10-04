package models

// Uji audit silang putaran 3 (bab 3.2 butir 1, tombol Add/Delete grid; W2, W5):
// daftar baris yang diterima dari layar mengikuti grid yang dapat disunting
// di XML.
//
//	DetailPolicyTreatyIn (Proporsional, wadah S19 `.IsNewPolicyNonProp != 1 &&
//	  .IsNewPolicyListFormat != 1`): grid S30 `.SpreadingRiskList` bertombol
//	  Add/Delete, sel TreatyType/%Share/%Share Claim terbuka
//	Section/SpreadingRiskList (di DetailPolicyTreatyInNonProportional, disertakan
//	  `DetailPoliciesNonProportional` - wadah `.IsNewPolicyNonProp = 1 &&
//	  .ClaimType != 'XOL Retro'` di KEDUA layar: admin S17, atasan S88, sel
//	  SUB_SECTION `pyEditOptions=Auto`): Add/Delete tampil bila
//	  `pyWorkPage.TreatyIn.FacultativeShare = 0 || = ''`; TreatyType/%Share/%Share
//	  Claim `pyReadOnlyCondition FacultativeShare > 0`
//	DetailDeptHeadTreatyIn_UW grid S96 (Proporsional): sel `pyReadOnly=true` tanpa
//	  syarat, tanpa Add/Delete - tidak diterima
//	grid `.ListInstallment` (S45 / S107): `pyEditingMode` / `pyRowEditing` `readOnly`
//	  di kedua layar - TIDAK PERNAH diterima (W5; baris dari FillPaymentInstallment)

import "testing"

func TestGabungMasukanDaftarMenurutGridXML(t *testing.T) {
	kiriman := func() *Halaman {
		m := HalamanBaru()
		m.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "UJI-1"}, {"TreatyType": "UJI-2"}})
		m.SetelDaftar(DaftarAngsuran, []Baris{{"InstallmentNo": "1", "InstallmentPercentage": "100"}})
		return m
	}
	for _, tt := range []struct {
		nama, posisi, nonProp, fakultatif, klaim string
		spreading                                int
	}{
		{"admin Proporsional", PosisiAdmin, "", "", "", 2},
		{"admin NonProp FacultativeShare 0", PosisiAdmin, "1", "0", "", 2},
		{"admin NonProp FacultativeShare kosong", PosisiAdmin, "1", "", "", 2},
		{"admin NonProp FacultativeShare 5", PosisiAdmin, "1", "5", "", 0},
		{"admin NonProp XOL Retro (wadah S17 tersembunyi)", PosisiAdmin, "1", "0", KlaimXOLRetro, 0},
		{"Sec Head Proporsional (S96 pyReadOnly)", PosisiSecHead, "", "", "", 0},
		{"Dept Head Proporsional (S96 pyReadOnly)", PosisiDeptHead, "", "", "", 0},
		{"Sec Head NonProp FacultativeShare 0 (S88)", PosisiSecHead, "1", "0", "", 2},
		{"Dept Head NonProp FacultativeShare kosong (S88)", PosisiDeptHead, "1", "", "", 2},
		{"Dept Head NonProp FacultativeShare 5", PosisiDeptHead, "1", "5", "", 0},
		{"Sec Head NonProp XOL Retro (wadah S88 tersembunyi)", PosisiSecHead, "1", "0", KlaimXOLRetro, 0},
	} {
		h := HalamanBaru()
		h.Setel("PolicyTreatyIn.IsNewPolicyNonProp", tt.nonProp)
		h.Setel("PolicyTreatyIn.ClaimType", tt.klaim)
		h.Setel("TreatyIn.FacultativeShare", tt.fakultatif)
		if _, err := GabungMasukanLayar(h, kiriman(), tt.posisi, nil); err != nil {
			t.Fatalf("%s: %v", tt.nama, err)
		}
		if got := len(h.AmbilDaftar(DaftarSpreading)); got != tt.spreading {
			t.Errorf("%s: SpreadingRiskList %d baris, harap %d", tt.nama, got, tt.spreading)
		}
		if got := len(h.AmbilDaftar(DaftarAngsuran)); got != 0 {
			t.Errorf("%s: ListInstallment grid readOnly - %d baris diterima dari layar", tt.nama, got)
		}
		if got := SpreadingDariLayar(h, tt.posisi); got != (tt.spreading > 0) {
			t.Errorf("%s: SpreadingDariLayar %v (aksi CountSpreading terbuka)", tt.nama, got)
		}
	}
}
