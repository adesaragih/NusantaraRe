package models

// Uji audit silang putaran 3 (bab 3.2 butir 1, tombol Add/Delete grid; W2) -
// matriks fungsi murni `GabungMasukanLayar` / `SpreadingDariLayar` per posisi
// (kolom hanya-baca W5 dan anggota baris server: seam HTTP, C10 tinjauan P3):
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
	// Perintah work owner 06-10-2026: baris spreading TIDAK dapat ditambah / dihapus dan Treaty Type hanya-baca -
	// jumlah baris dan Treaty Type tetap milik server; dari layar hanya %Share / %Share Claim (grid yang terbuka).
	kiriman := func() *Halaman {
		m := HalamanBaru()
		m.SetelDaftar(DaftarSpreading, []Baris{
			{"TreatyType": "UJI-GANTI", "SharePercentage": "10"},
			{"TreatyType": "UJI-2", "SharePercentage": "20", "ClaimPercentage": "5"},
			{"TreatyType": "UJI-TAMBAH", "SharePercentage": "30"},
		})
		m.SetelDaftar(DaftarAngsuran, []Baris{{"InstallmentNo": "1", "InstallmentPercentage": "100"}})
		return m
	}
	for _, tt := range []struct {
		nama, posisi, nonProp, fakultatif, klaim string
		terima                                   bool
	}{
		{"admin Proporsional", PosisiAdmin, "", "", "", true},
		{"admin NonProp FacultativeShare 0", PosisiAdmin, "1", "0", "", true},
		{"admin NonProp FacultativeShare kosong", PosisiAdmin, "1", "", "", true},
		{"admin NonProp FacultativeShare 5", PosisiAdmin, "1", "5", "", false},
		{"admin NonProp XOL Retro (wadah S17 tersembunyi)", PosisiAdmin, "1", "0", KlaimXOLRetro, false},
		{"Sec Head Proporsional (S96 pyReadOnly)", PosisiSecHead, "", "", "", false},
		{"Dept Head Proporsional (S96 pyReadOnly)", PosisiDeptHead, "", "", "", false},
		// [keputusan work owner 07-10-2026] "HANYA ADMIN YANG BISA EDIT": S88 atasan hanya-baca
		{"Sec Head NonProp FacultativeShare 0 (S88, hanya admin)", PosisiSecHead, "1", "0", "", false},
		{"Dept Head NonProp FacultativeShare kosong (S88, hanya admin)", PosisiDeptHead, "1", "", "", false},
		{"Dept Head NonProp FacultativeShare 5", PosisiDeptHead, "1", "5", "", false},
		{"Sec Head NonProp XOL Retro (wadah S88 tersembunyi)", PosisiSecHead, "1", "0", KlaimXOLRetro, false},
	} {
		h := HalamanBaru()
		h.Setel("PolicyTreatyIn.IsNewPolicyNonProp", tt.nonProp)
		h.Setel("PolicyTreatyIn.ClaimType", tt.klaim)
		h.Setel("TreatyIn.FacultativeShare", tt.fakultatif)
		h.SetelDaftar(DaftarSpreading, []Baris{{"TreatyType": "UJI-1", "TreatyName": "UJI SATU"}, {"TreatyType": "UJI-2"}})
		if _, err := GabungMasukanLayar(h, kiriman(), tt.posisi, nil); err != nil {
			t.Fatalf("%s: %v", tt.nama, err)
		}
		d := h.AmbilDaftar(DaftarSpreading)
		if len(d) != 2 {
			t.Fatalf("%s: SpreadingRiskList %d baris, harap tetap 2 (tanpa Add/Delete)", tt.nama, len(d))
		}
		if d[0]["TreatyType"] != "UJI-1" || d[0]["TreatyName"] != "UJI SATU" {
			t.Errorf("%s: Treaty Type hanya-baca - baris 1 %+v", tt.nama, d[0])
		}
		harap := map[bool][2]string{true: {"10", "20"}, false: {"", ""}}[tt.terima]
		if d[0]["SharePercentage"] != harap[0] || d[1]["SharePercentage"] != harap[1] {
			t.Errorf("%s: %%Share = %q/%q, harap %q/%q", tt.nama, d[0]["SharePercentage"], d[1]["SharePercentage"], harap[0], harap[1])
		}
		if got := len(h.AmbilDaftar(DaftarAngsuran)); got != 0 {
			t.Errorf("%s: ListInstallment grid readOnly - %d baris diterima dari layar", tt.nama, got)
		}
		if got := SpreadingDariLayar(h, tt.posisi); got != tt.terima {
			t.Errorf("%s: SpreadingDariLayar %v (aksi CountSpreading terbuka)", tt.nama, got)
		}
	}
}
