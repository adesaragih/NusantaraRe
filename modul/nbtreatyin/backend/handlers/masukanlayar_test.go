package handlers_test

// Uji audit silang putaran 3 bab 7 (W2-W5, 7.4) lewat seam HTTP: medan dan
// daftar yang DITERIMA server dari kiriman layar mengikuti sel / wadah / grid
// yang terbuka di Section XML, dan kolom hanya-baca dihitung server. Nilai
// harapan dihitung tangan dari langkah XML yang dikutip; fixture UJI-.

import (
	"net/http"
	"testing"

	"nusantarare/modul/nbtreatyin/backend/models"
)

// npSampaiSecHead - kasus NonProp (kontrak `kontrakNP`, master FacultativeShare
// `fakultatif`) dipilih dan disetujui admin; kembali id-nya, menunggu Sec Head.
func npSampaiSecHead(u *uji, fakultatif string) string {
	u.t.Helper()
	kontrakNP(u)
	m := masterNP()
	m.Nilai["FacultativeShare"] = fakultatif
	u.g.Master["UJI-T-NP"] = m
	id := u.buat()
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D-NP"}); kode != http.StatusOK {
		u.t.Fatalf("pilih bisnis NonProp: %d %s", kode, isi)
	}
	_, isi := u.panggil("GET", "/kasus/"+id, admin, nil)
	h := u.layar(isi).Halaman
	h.Setel("PolicyTreatyIn.QuotationData.MOID", "UJI-MO-1")
	h.Setel("PolicyTreatyIn.IsApproved", "1")
	h.Setel("PolicyTreatyIn.Suggest", "UJI-catatan")
	if kode, isi := u.kirim(id, admin, h); kode != http.StatusOK {
		u.t.Fatalf("admin submit NonProp: %d %s", kode, isi)
	}
	return id
}

// W2 - `Section/DetailDeptHeadTreatyIn_UW` S88 (`.IsNewPolicyNonProp = 1 &&
// .ClaimType != 'XOL Retro'`) menyertakan `DetailPoliciesNonProportional`
// (SUB_SECTION `pyEditOptions=Auto`) -> `DetailPolicyTreatyInNonProportional`
// S73 -> `Section/SpreadingRiskList`: tombol Add/Delete tampil bila
// `pyWorkPage.TreatyIn.FacultativeShare = 0 || = ”`; `.TreatyType`,
// `.SharePercentage`, `.ClaimPercentage` `pyReadOnlyCondition FacultativeShare >0`;
// %Share: change -> refresh `CountSpreading_Act(Index=.pxListSubscript)`.
// Syaratnya SAMA dengan layar admin - atasan dapat menyunting grid itu.
//
// CountSpreading_Act langkah 4.1: `.PremiumSpreaded = Primary.NetPremium *
// @divide(.SharePercentage,100,10)`; NetPremium polis XOL = 2700
// (InputPolicyTreatyInDetail_NonProp, hitung tangan `TestNonPropPilihBisnis...`):
// 2700 x 60% = 1620, 2700 x 40% = 1080.
func TestAtasanMenyuntingSpreadingNonProp(t *testing.T) {
	u := baru(t)
	id := npSampaiSecHead(u, "0")
	_, isi := u.panggil("GET", "/kasus/"+id, secHead, nil)
	h := u.layar(isi).Halaman
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{
		{"TreatyType": "UJI-SPR-ID", "SharePercentage": "60"},
		{"TreatyType": "UJI-SPR-2", "SharePercentage": "40"},
	})
	h.Setel("PolicyTreatyIn.IsApproved", "1")
	h.Setel("PolicyTreatyIn.Suggest", "UJI-sec")
	kode, isi := u.panggil("POST", "/kasus/"+id+"/hitung", secHead,
		map[string]any{"urutan": []map[string]string{{"aksi": "CountSpreading"}}, "indeks": 2, "halaman": h})
	if kode != http.StatusOK {
		t.Fatalf("Sec Head CountSpreading (sel %%Share terbuka, FacultativeShare 0): %d %s", kode, isi)
	}
	hasil := u.layar(isi).Halaman
	if kode, isi := u.kirim(id, secHead, hasil); kode != http.StatusOK {
		t.Fatalf("Sec Head submit: %d %s", kode, isi)
	}
	sp := u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading)
	if len(sp) != 2 || sp[1]["TreatyType"] != "UJI-SPR-2" {
		t.Fatalf("SpreadingRiskList suntingan Sec Head tersimpan: %v", sp)
	}
	angkaSamaTeks(t, "Spreading(1).PremiumSpreaded", sp[0]["PremiumSpreaded"], "1620")
	angkaSamaTeks(t, "Spreading(2).PremiumSpreaded", sp[1]["PremiumSpreaded"], "1080")
}

// W2 - FacultativeShare > 0: Add/Delete tersembunyi, sel terkunci - kiriman
// diabaikan dan CountSpreading ditolak (409). Grid Proporsional atasan
// (`DetailDeptHeadTreatyIn_UW` S96) `pyReadOnly=true` tanpa syarat - tetap
// hanya-baca.
func TestAtasanSpreadingTerkunciBilaFakultatifAtauProporsional(t *testing.T) {
	u := baru(t)
	id := npSampaiSecHead(u, "5")
	lama := u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading)
	p := putusan("1")
	p.SetelDaftar(models.DaftarSpreading, []models.Baris{{"TreatyType": "UJI-KARANGAN", "SharePercentage": "1"}})
	if kode, _ := u.panggil("POST", "/kasus/"+id+"/hitung", secHead,
		map[string]any{"urutan": []map[string]string{{"aksi": "CountSpreading"}}, "indeks": 1, "halaman": p}); kode != http.StatusConflict {
		t.Fatalf("FacultativeShare 5: CountSpreading atasan %d, harap 409", kode)
	}
	if kode, isi := u.kirim(id, secHead, p); kode != http.StatusOK {
		t.Fatalf("Sec Head submit: %d %s", kode, isi)
	}
	if sp := u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading); len(sp) != len(lama) || (len(sp) > 0 && sp[0]["TreatyType"] == "UJI-KARANGAN") {
		t.Fatalf("FacultativeShare 5: grid terkunci, kiriman diabaikan: %v", sp)
	}

	// Proporsional: grid atasan hanya-baca.
	id2 := u.buat()
	if kode, isi := u.kirim(id2, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("admin: %d %s", kode, isi)
	}
	if kode, _ := u.panggil("POST", "/kasus/"+id2+"/hitung", secHead,
		map[string]any{"urutan": []map[string]string{{"aksi": "CountSpreading"}}, "indeks": 1, "halaman": p}); kode != http.StatusConflict {
		t.Fatalf("Proporsional: CountSpreading atasan %d, harap 409", kode)
	}
	if kode, isi := u.kirim(id2, secHead, p); kode != http.StatusOK {
		t.Fatalf("Sec Head submit: %d %s", kode, isi)
	}
	if sp := u.g.Halaman[id2].AmbilDaftar(models.DaftarSpreading); len(sp) != 0 {
		t.Fatalf("Proporsional: grid atasan pyReadOnly, kiriman diabaikan: %v", sp)
	}
}
