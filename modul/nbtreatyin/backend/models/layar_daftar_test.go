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

import (
	"testing"
	"time"
)

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

// W5 - kolom `.PremiumSpreaded` / `.ClaimSpreaded` grid spreading `Read-only`
// (DetailPolicyTreatyIn S30 C[2.3]/C[2.5], SpreadingRiskList, DetailDeptHeadTreatyIn_UW
// S96): nilainya TIDAK PERNAH dari layar. Ia hanya ditulis `CountSpreading_Act`
// langkah 4.1, yang terpicu sel `.SharePercentage` / `.ClaimPercentage`
// (change -> refresh) - selain itu nilai server (tersimpan) yang bertahan.
//
//	4.1  .SharePercentage = @if(.SharePercentage=="",(100/@LengthOfPageList(..)),.SharePercentage)
//	     .ClaimPercentage = @if(.ClaimPercentage == "",.SharePercentage,.ClaimPercentage)
//	     .PremiumSpreaded = Primary.NetPremium * @divide(.SharePercentage,100,10)
//	     .ClaimSpreaded   = (Primary.ExcessLoss + Primary.Claim - Primary.SalvageValue)* @divide(.ClaimPercentage,100,10)
//	4.2/5 total = jumlah keempat kolom
//
// Hitung tangan: NetPremium 1000, Claim 200 (ExcessLoss, SalvageValue 0).
func TestKolomHanyaBacaSpreadingDariServer(t *testing.T) {
	tersimpan := func() *Halaman {
		h := HalamanBaru()
		h.Setel("PolicyTreatyIn.NetPremium", "1000")
		h.Setel("PolicyTreatyIn.Claim", "200")
		h.SetelDaftar(DaftarSpreading, []Baris{
			{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60", "PremiumSpreaded": "600", "ClaimSpreaded": "120"},
			{"TreatyType": "UJI-B", "SharePercentage": "40", "ClaimPercentage": "40", "PremiumSpreaded": "400", "ClaimSpreaded": "80"},
		})
		return h
	}
	jalankan := func(t *testing.T, baris []Baris) *Halaman {
		t.Helper()
		h := tersimpan()
		m := HalamanBaru()
		m.SetelDaftar(DaftarSpreading, baris)
		p, err := GabungMasukanLayar(h, m, PosisiAdmin, nil)
		if err != nil {
			t.Fatal(err)
		}
		if err := TerapkanPemicu(h, p, time.Date(2026, 10, 3, 0, 0, 0, 0, time.UTC)); err != nil {
			t.Fatal(err)
		}
		return h
	}
	cek := func(t *testing.T, h *Halaman, harap [][2]string, total [4]string) {
		t.Helper()
		b := h.AmbilDaftar(DaftarSpreading)
		if len(b) != len(harap) {
			t.Fatalf("%d baris, harap %d: %v", len(b), len(harap), b)
		}
		for i, w := range harap {
			if !samaNilai(b[i]["PremiumSpreaded"], w[0]) || !samaNilai(b[i]["ClaimSpreaded"], w[1]) {
				t.Errorf("baris %d: Premium/ClaimSpreaded %q/%q, harap %q/%q", i+1, b[i]["PremiumSpreaded"], b[i]["ClaimSpreaded"], w[0], w[1])
			}
		}
		for i, j := range []string{"TotalSharePercentagePremium", "TotalPremium", "TotalSharePercentageClaim", "TotalClaim"} {
			if got := h.Ambil("PolicyTreatyIn." + j); !samaNilai(got, total[i]) {
				t.Errorf("%s = %q, harap %s", j, got, total[i])
			}
		}
	}
	t.Run("kolom hanya-baca palsu, %Share tetap: nilai tersimpan", func(t *testing.T) {
		h := jalankan(t, []Baris{
			{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60", "PremiumSpreaded": "999", "ClaimSpreaded": "999"},
			{"TreatyType": "UJI-C", "SharePercentage": "40", "ClaimPercentage": "40", "PremiumSpreaded": "999"},
		})
		cek(t, h, [][2]string{{"600", "120"}, {"400", "80"}}, [4]string{"100", "1000", "100", "200"})
		if got := h.AmbilDaftar(DaftarSpreading)[1]["TreatyType"]; got != "UJI-C" {
			t.Errorf("sel .TreatyType terbuka: %q", got)
		}
	})
	t.Run("%Share berubah: CountSpreading_Act 4.1 seluruh baris", func(t *testing.T) {
		h := jalankan(t, []Baris{
			{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60", "PremiumSpreaded": "999"},
			{"TreatyType": "UJI-B", "SharePercentage": "50", "ClaimPercentage": "40", "PremiumSpreaded": "999"},
		})
		cek(t, h, [][2]string{{"600", "120"}, {"500", "80"}}, [4]string{"110", "1100", "100", "200"})
	})
	t.Run("Add tanpa %Share: tidak ada refresh, baris baru kosong", func(t *testing.T) {
		h := jalankan(t, []Baris{
			{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60"},
			{"TreatyType": "UJI-B", "SharePercentage": "40", "ClaimPercentage": "40"},
			{"TreatyType": "UJI-C", "PremiumSpreaded": "999"},
		})
		cek(t, h, [][2]string{{"600", "120"}, {"400", "80"}, {"", ""}}, [4]string{"100", "1000", "100", "200"})
	})
	t.Run("Add dengan %Share: 4.1 mengisi %Share Claim kosong", func(t *testing.T) {
		h := jalankan(t, []Baris{
			{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60"},
			{"TreatyType": "UJI-B", "SharePercentage": "40", "ClaimPercentage": "40"},
			{"TreatyType": "UJI-C", "SharePercentage": "10"},
		})
		cek(t, h, [][2]string{{"600", "120"}, {"400", "80"}, {"100", "20"}}, [4]string{"110", "1100", "110", "220"})
	})
	t.Run("Delete: baris bergeser, dihitung ulang", func(t *testing.T) {
		h := jalankan(t, []Baris{{"TreatyType": "UJI-B", "SharePercentage": "40", "ClaimPercentage": "40", "PremiumSpreaded": "999"}})
		cek(t, h, [][2]string{{"400", "80"}}, [4]string{"40", "400", "40", "80"})
	})
}

// samaNilai - dua teks angka sama nilainya ("" hanya sama dengan "").
func samaNilai(a, b string) bool {
	if a == "" || b == "" {
		return a == b
	}
	x, err1 := AngkaTeks("a", a)
	y, err2 := AngkaTeks("b", b)
	return err1 == nil && err2 == nil && x.Cmp(y) == 0
}
