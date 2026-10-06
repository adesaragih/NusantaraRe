package handlers_test

// Uji audit silang putaran 3 bab 7 (W2-W5, 7.4) lewat seam HTTP: medan dan
// daftar yang DITERIMA server dari kiriman layar mengikuti sel / wadah / grid
// yang terbuka di Section XML, dan kolom hanya-baca dihitung server. Nilai
// harapan dihitung tangan dari langkah XML yang dikutip; fixture UJI-.

import (
	"net/http"
	"strings"
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

// simpanOK - tombol Save admin; kode 200 diharapkan.
func (u *uji) simpanOK(id string, h *models.Halaman) {
	u.t.Helper()
	if kode, isi := u.simpan(id, h); kode != http.StatusOK {
		u.t.Fatalf("simpan: %d %s", kode, isi)
	}
}

// W5 - kolom `.PremiumSpreaded` / `.ClaimSpreaded` (sel `Read-only`
// `DetailPolicyTreatyIn` S30) tidak pernah diterima dari layar: ditulis
// `CountSpreading_Act` langkah 4.1 di server bila sel %Share memicunya, selain
// itu nilai tersimpan bertahan. Hitung tangan: halamanLengkap -> CountNetPremi_act
// langkah 4: NetPremium = (1000 - 0) + (0 - 0) - 0 - 0 - 0 - 0 = 1000;
// 1000 x @divide(60,100,10) = 600; klaim (0 + 0 - 0) x 60% = 0.
func TestKolomHanyaBacaSpreadingTidakDariLayar(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{
		{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60", "PremiumSpreaded": "999999", "ClaimSpreaded": "999999"},
	})
	u.simpanOK(id, h)
	sp := u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading)
	if len(sp) != 1 {
		t.Fatalf("SpreadingRiskList %v", sp)
	}
	angkaSamaTeks(t, "PremiumSpreaded (Save pertama)", sp[0]["PremiumSpreaded"], "600")
	angkaSamaTeks(t, "ClaimSpreaded (Save pertama)", sp[0]["ClaimSpreaded"], "0")

	// %Share tetap, kolom hanya-baca dipalsukan: nilai server bertahan.
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{
		{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60", "PremiumSpreaded": "1", "ClaimSpreaded": "1"},
	})
	u.simpanOK(id, h)
	sp = u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading)
	angkaSamaTeks(t, "PremiumSpreaded (Save kedua)", sp[0]["PremiumSpreaded"], "600")
	angkaSamaTeks(t, "ClaimSpreaded (Save kedua)", sp[0]["ClaimSpreaded"], "0")
}

// spreadingDasar - kasus admin dengan dua baris spreading tersimpan sesudah Save:
// A 60/60, B 40/40; NetPremium 1000 (halamanLengkap), klaim 0 -> PremiumSpreaded
// 600 / 400 (CountSpreading_Act 4.1), ClaimSpreaded 0.
func spreadingDasar(t *testing.T) (*uji, string, *models.Halaman) {
	t.Helper()
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{
		{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60"},
		{"TreatyType": "UJI-B", "SharePercentage": "40", "ClaimPercentage": "40"},
	})
	u.simpanOK(id, h)
	return u, id, h
}

// simpanSpreading - Save admin; jawabannya halaman layar (total spreading
// tanpa kolom - hanya di layar).
func (u *uji) simpanSpreading(id string, h *models.Halaman) *models.Halaman {
	u.t.Helper()
	kode, isi := u.simpan(id, h)
	if kode != http.StatusOK {
		u.t.Fatalf("simpan: %d %s", kode, isi)
	}
	return layarDari(u.t, isi).Halaman
}

// cekSpreading - PremiumSpreaded per baris tersimpan dan dua total di layar
// (`CountSpreading_Act` 4.2/5: TotalSharePercentagePremium, TotalPremium).
func cekSpreading(t *testing.T, u *uji, id string, ly *models.Halaman, premi []string, totalPct, totalPremi string) {
	t.Helper()
	s := ly
	b := u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading)
	if len(b) != len(premi) {
		t.Fatalf("%d baris, harap %d: %v", len(b), len(premi), b)
	}
	for i, w := range premi {
		if w == "" {
			if b[i]["PremiumSpreaded"] != "" {
				t.Errorf("baris %d PremiumSpreaded %q, harap kosong", i+1, b[i]["PremiumSpreaded"])
			}
			continue
		}
		angkaSamaTeks(t, "PremiumSpreaded", b[i]["PremiumSpreaded"], w)
	}
	angkaSamaTeks(t, "TotalSharePercentagePremium", s.Ambil("PolicyTreatyIn.TotalSharePercentagePremium"), totalPct)
	angkaSamaTeks(t, "TotalPremium", s.Ambil("PolicyTreatyIn.TotalPremium"), totalPremi)
}

// W5 lanjutan (dipindah dari uji models, C10): sel %Share memicu
// `CountSpreading_Act` (change -> refresh) atas SELURUH baris; Add tanpa %Share
// tidak memicu (baris baru kosong sampai refresh berikutnya); Delete menggeser
// baris dan dihitung ulang. 4.1 `.ClaimPercentage = @if(.ClaimPercentage == "",
// .SharePercentage, ...)`. Hitung tangan dengan NetPremium 1000.
func TestSpreadingDipicuSelPersen(t *testing.T) {
	t.Run("Share B 50: 600 / 500, total 110 / 1100", func(t *testing.T) {
		u, id, h := spreadingDasar(t)
		h.SetelDaftar(models.DaftarSpreading, []models.Baris{
			{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60", "PremiumSpreaded": "999"},
			{"TreatyType": "UJI-B", "SharePercentage": "50", "ClaimPercentage": "40", "PremiumSpreaded": "999"},
		})
		cekSpreading(t, u, id, u.simpanSpreading(id, h), []string{"600", "500"}, "110", "1100")
	})
	t.Run("Add tanpa Share: tidak terpicu, baris baru kosong", func(t *testing.T) {
		u, id, h := spreadingDasar(t)
		h.SetelDaftar(models.DaftarSpreading, []models.Baris{
			{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60"},
			{"TreatyType": "UJI-B", "SharePercentage": "40", "ClaimPercentage": "40"},
			{"TreatyType": "UJI-C", "PremiumSpreaded": "999"},
		})
		cekSpreading(t, u, id, u.simpanSpreading(id, h), []string{"600", "400", ""}, "100", "1000")
	})
	t.Run("Add dengan Share 10: 4.1 mengisi Share Claim, 100", func(t *testing.T) {
		u, id, h := spreadingDasar(t)
		h.SetelDaftar(models.DaftarSpreading, []models.Baris{
			{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60"},
			{"TreatyType": "UJI-B", "SharePercentage": "40", "ClaimPercentage": "40"},
			{"TreatyType": "UJI-C", "SharePercentage": "10"},
		})
		cekSpreading(t, u, id, u.simpanSpreading(id, h), []string{"600", "400", "100"}, "110", "1100")
		angkaSamaTeks(t, "ClaimPercentage baris 3", u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading)[2]["ClaimPercentage"], "10")
	})
	t.Run("Delete A: B bergeser ke baris 1, 400", func(t *testing.T) {
		u, id, h := spreadingDasar(t)
		h.SetelDaftar(models.DaftarSpreading, []models.Baris{
			{"TreatyType": "UJI-B", "SharePercentage": "40", "ClaimPercentage": "40", "PremiumSpreaded": "999"},
		})
		cekSpreading(t, u, id, u.simpanSpreading(id, h), []string{"400"}, "40", "400")
	})
}

// Pengerasan tercatat (PERMINTAAN H2): anggota baris spreading yang BUKAN sel
// grid - `.TreatyName` (hanya `pyPrompt` dropdown; ditulis
// `TreatyInputPctCommSpreading`), `.Currency`, `.CurrencyID`,
// `.SplitRNMSharePct` (ditulis `TreatyNonPropSetSpreading` 3/4.1/6) - tidak
// ditulis action set sel mana pun (`CountSpreading_Act` hanya %Share dan
// kolom Spreaded), jadi TIDAK diterima dari layar: nilainya nilai baris server,
// mengikuti barisnya sendiri saat Delete; baris Add / nilai karangan kosong.
func TestAnggotaBarisSpreadingDariServer(t *testing.T) {
	anggota := []string{"TreatyName", "Currency", "CurrencyID", "SplitRNMSharePct"}
	a := map[string]string{"TreatyName": "UJI NAMA A", "Currency": "IDR", "CurrencyID": "UJI-ID-IDR", "SplitRNMSharePct": "30"}
	bb := map[string]string{"TreatyName": "UJI NAMA B", "Currency": "USD", "CurrencyID": "UJI-ID-USD", "SplitRNMSharePct": "20"}
	kosong := map[string]string{}
	// isiServer - nilai yang ditulis aktivitas pilih bisnis ke baris tersimpan.
	isiServer := func(u *uji, id string) {
		b := u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading)
		for i, isi := range []map[string]string{a, bb} {
			for kk, v := range isi {
				b[i][kk] = v
			}
		}
	}
	// salinLayar - layar menyalin baris server utuh (LayarKasus `ubahBaris`).
	salinLayar := func(u *uji, id string) []models.Baris {
		var out []models.Baris
		for _, b := range u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading) {
			x := models.Baris{}
			for kk, v := range b {
				x[kk] = v
			}
			out = append(out, x)
		}
		return out
	}
	cek := func(t *testing.T, u *uji, id string, harap []map[string]string) {
		t.Helper()
		b := u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading)
		if len(b) != len(harap) {
			t.Fatalf("%d baris, harap %d", len(b), len(harap))
		}
		for i, w := range harap {
			for _, m := range anggota {
				if b[i][m] != w[m] {
					t.Errorf("baris %d %s = %q, harap %q", i+1, m, b[i][m], w[m])
				}
			}
		}
	}

	t.Run("dibawa layar apa adanya: bertahan", func(t *testing.T) {
		u, id, h := spreadingDasar(t)
		isiServer(u, id)
		rows := salinLayar(u, id)
		rows[1]["SharePercentage"] = "50"
		h.SetelDaftar(models.DaftarSpreading, rows)
		u.simpanOK(id, h)
		cek(t, u, id, []map[string]string{a, bb})
	})
	t.Run("nilai karangan dan Add tiruan baris A: kosong", func(t *testing.T) {
		u, id, h := spreadingDasar(t)
		isiServer(u, id)
		rows := salinLayar(u, id)
		rows[1]["TreatyName"] = "UJI KARANGAN"
		rows[1]["SplitRNMSharePct"] = "99"
		tiruA := models.Baris{"TreatyType": "UJI-C"}
		for kk, v := range a {
			tiruA[kk] = v
		}
		h.SetelDaftar(models.DaftarSpreading, append(rows, tiruA))
		u.simpanOK(id, h)
		cek(t, u, id, []map[string]string{a, kosong, kosong})
	})
	t.Run("Delete A: baris B membawa anggotanya sendiri", func(t *testing.T) {
		u, id, h := spreadingDasar(t)
		isiServer(u, id)
		h.SetelDaftar(models.DaftarSpreading, salinLayar(u, id)[1:])
		u.simpanOK(id, h)
		cek(t, u, id, []map[string]string{bb})
	})
	t.Run("anggota di luar sel dan server tidak diterima", func(t *testing.T) {
		u, id, h := spreadingDasar(t)
		rows := salinLayar(u, id)
		rows[0]["UJISelundupan"] = "x"
		h.SetelDaftar(models.DaftarSpreading, rows)
		u.simpanOK(id, h)
		if _, ada := u.g.Halaman[id].AmbilDaftar(models.DaftarSpreading)[0]["UJISelundupan"]; ada {
			t.Fatal("anggota di luar sel / server diterima dari layar")
		}
	})
}

// W4 - polis NonProp baru: medan uang berada di wadah S19 `.IsNewPolicyNonProp
// != 1 && .IsNewPolicyListFormat != 1` yang TERSEMBUNYI; nilainya milik
// InputPolicyTreatyInDetail_NonProp 18-19 (PremiOgp 3000, Deduction1 306.6,
// Deduction2 0 - hitung tangan `TestNonPropPilihBisnisHitungSimpanBacaKembali`).
// Kiriman layar Save untuk medan itu diabaikan.
func TestKirimanLayarTidakMenimpaUangMasterNonProp(t *testing.T) {
	u := baru(t)
	kontrakNP(u)
	id := u.buat()
	if kode, isi := u.panggil("POST", "/kasus/"+id+"/pilih-bisnis", admin, map[string]any{"idDetail": "UJI-D-NP"}); kode != http.StatusOK {
		t.Fatalf("pilih bisnis NonProp: %d %s", kode, isi)
	}
	_, isi := u.panggil("GET", "/kasus/"+id, admin, nil)
	h := u.layar(isi).Halaman
	h.Setel("PolicyTreatyIn.IsApproved", "1")
	h.Setel("PolicyTreatyIn.Suggest", "UJI-catatan")
	h.Setel("PolicyTreatyIn.QuotationData.MOID", "UJI-MO-1")
	h.Setel("PolicyTreatyIn.PremiOgp", "1")
	h.Setel("PolicyTreatyIn.Deduction1", "2")
	h.Setel("PolicyTreatyIn.Deduction2", "3")
	h.Setel("PolicyTreatyIn.Installment", "4")
	u.simpanOK(id, h)
	s := u.g.Halaman[id]
	angkaSamaTeks(t, "PremiOgp", s.Ambil("PolicyTreatyIn.PremiOgp"), "3000")
	angkaSamaTeks(t, "Deduction1", s.Ambil("PolicyTreatyIn.Deduction1"), "306.6")
	if got := s.Ambil("PolicyTreatyIn.Deduction2"); got != "" && got != "0" {
		t.Errorf("Deduction2 = %q, harap 0 (InputPolicyTreatyInDetail_NonProp)", got)
	}
	if got := s.Ambil("PolicyTreatyIn.Installment"); got == "4" {
		t.Error(".Installment di wadah S19 tersembunyi: kiriman diabaikan")
	}
	if got := s.Ambil("PolicyTreatyIn.QuotationData.MOID"); got != "UJI-MO-1" {
		t.Errorf("medan tampil tetap diterima: MOID %q", got)
	}
}

// W3 - `.QuotationData.ProportionalType` (`pyReadOnly=true`) dan
// `.IsNewPolicyNonProp` hanya diubah tombol "Enable / Disable Input Type"
// (pyVisible `.TreatyType='XOL'`) -> DataTransform `TreatyEnableDisableInput`:
// langkah 1 ProportionalType "NonProportional", 2-4 IsNewPolicyNonProp "0".
func TestEnableDisableHanyaDariTombolXOL(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.QuotationData.ProportionalType", models.JenisNonProporsional)
	u.simpanOK(id, h)
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.QuotationData.ProportionalType"); got != "" {
		t.Fatalf("TreatyType bukan XOL (tombol tak tampil): medan terkunci, tersimpan %q", got)
	}

	u.g.Halaman[id].Setel("PolicyTreatyIn.TreatyType", "XOL")
	// klik tombol: refresh tanpa simpan, hasil DT dipegang layar
	ly := u.hitung(id, map[string]any{"urutan": []map[string]string{{"aksi": "TreatyEnableDisableInput"}}, "halaman": h})
	if ly.Halaman.Ambil("PolicyTreatyIn.QuotationData.ProportionalType") != models.JenisNonProporsional ||
		ly.Halaman.Ambil("PolicyTreatyIn.IsNewPolicyNonProp") != "0" {
		t.Fatalf("hasil DT: %q / %q", ly.Halaman.Ambil("PolicyTreatyIn.QuotationData.ProportionalType"),
			ly.Halaman.Ambil("PolicyTreatyIn.IsNewPolicyNonProp"))
	}
	u.simpanOK(id, ly.Halaman)
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.QuotationData.ProportionalType"); got != models.JenisNonProporsional {
		t.Fatalf("XOL: hasil TreatyEnableDisableInput yang dipegang layar disimpan; tersimpan %q", got)
	}

	palsu := ly.Halaman.Salin()
	palsu.Setel("PolicyTreatyIn.QuotationData.ProportionalType", models.JenisProporsional)
	kode, isi := u.simpan(id, palsu)
	if kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "TreatyEnableDisableInput") {
		t.Fatalf("nilai yang bukan hasil tombol: %d %s, harap 422", kode, isi)
	}
}

// W5 - grid `.ListInstallment` S45 `readOnly`: baris dari action set server.
//
//	.Installment 2 -> FillPaymentInstallment: 3.2 pct = @Math.divide(100,2,4) = 50;
//	   3.5 Premium = @Math.divide(BalanceDueTo x 50, 100, 4); DueDate = @CurrentDateTime()
//	   (jam uji 2026-10-03). halamanLengkap: BalanceDueTo = NetPremium 1000 -> 500.
//	PremiOgp 2000 -> CountOGPONP_Act 10 SetValidateInstallment_Act 3.1:
//	   Premium = 2000 x @divide(50,100,4) = 1000; 9 CountSpreading_Act: 2000 x 60% = 1200.
func TestAngsuranDariActionSetServer(t *testing.T) {
	u := baru(t)
	id := u.buat()
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.Installment", "2")
	h.SetelDaftar(models.DaftarAngsuran, []models.Baris{{"InstallmentNo": "1", "InstallmentPercentage": "100", "Premium": "1", "PaymentTotal": "1", "DueDate": "2000-01-01"}})
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{{"TreatyType": "UJI-A", "SharePercentage": "60", "ClaimPercentage": "60"}})
	u.simpanOK(id, h)
	ang := u.g.Halaman[id].AmbilDaftar(models.DaftarAngsuran)
	if len(ang) != 2 {
		t.Fatalf("FillPaymentInstallment(2): %v", ang)
	}
	for i, b := range ang {
		angkaSamaTeks(t, "Angsuran.InstallmentPercentage", b["InstallmentPercentage"], "50")
		angkaSamaTeks(t, "Angsuran.Premium", b["Premium"], "500")
		angkaSamaTeks(t, "Angsuran.PaymentTotal", b["PaymentTotal"], "500")
		if b["DueDate"] != "2026-10-03" {
			t.Errorf("baris %d DueDate %q, harap @CurrentDateTime() 2026-10-03", i+1, b["DueDate"])
		}
	}

	// baris palsu tanpa pemicu: baris server bertahan
	u.simpanOK(id, h)
	if ang := u.g.Halaman[id].AmbilDaftar(models.DaftarAngsuran); len(ang) != 2 {
		t.Fatalf("baris kiriman layar tidak diterima: %v", ang)
	}

	h.Setel("PolicyTreatyIn.PremiOgp", "2000")
	u.simpanOK(id, h)
	s := u.g.Halaman[id]
	for _, b := range s.AmbilDaftar(models.DaftarAngsuran) {
		angkaSamaTeks(t, "SetValidateInstallment Premium", b["Premium"], "1000")
		angkaSamaTeks(t, "SetValidateInstallment PaymentTotal", b["PaymentTotal"], "1000")
	}
	angkaSamaTeks(t, "CountOGPONP_Act 9 PremiumSpreaded", s.AmbilDaftar(models.DaftarSpreading)[0]["PremiumSpreaded"], "1200")
}

// W5 - `pyDefaultValue` sel terbuka `Section/DetailPolicyTreatyIn`:
// `.QuotationData.IsSurveyReport` "No" (sel tampil bila Quotation bukan
// NonProportional) saat layar dirender; `.TypeTax` "Inclusive" bila FlagPPH
// dicentang dan TypeTax kosong.
func TestNilaiBawaanSelLayarAdmin(t *testing.T) {
	u := baru(t)
	id := u.buat()
	_, isi := u.panggil("GET", "/kasus/"+id, admin, nil)
	if got := u.layar(isi).Halaman.Ambil("PolicyTreatyIn.QuotationData.IsSurveyReport"); got != "No" {
		t.Fatalf("IsSurveyReport bawaan sel %q, harap No", got)
	}
	h := halamanLengkap("1")
	h.Setel("PolicyTreatyIn.FlagPPH", "true")
	u.simpanOK(id, h)
	if got := u.g.Halaman[id].Ambil("PolicyTreatyIn.TypeTax"); got != models.TypeTaxInclusive {
		t.Fatalf("TypeTax bawaan sel %q, harap Inclusive", got)
	}
}

// 7.4 audit silang P3 - tombol Submit Dept Head (IsApproved 1) ->
// `GeneratePolicyNoTreaty_Act` langkah 11 (`TreatyGroupOldID==""`) -> Call
// `FetchTreatyGroupOldID`: langkah 2 `Param.Errmsg = "Cannot fetch Treaty Group
// ID, Contact IT"`, langkah 3 `Page-Set-Messages pyWorkPage` bila
// `pyWorkPage.PolicyTreatyIn.TreatyGroupID==""`. Nomor tetap dibentuk dan disimpan
// (langkah 30 `Obj-Save WithErrors=true`), tetapi halaman berpesan tidak dapat
// di-submit: OK (finishAssignment) ditolak 422 dan berkas tetap di Dept Head.
func TestGrupTreatyKosongMenahanSubmitDeptHead(t *testing.T) {
	u := baru(t)
	u.g.OldIDGrup = "" // RDB FetchTreatyGroupOLDID atas ID kosong: nol baris
	id := u.buat()
	if kode, isi := u.kirim(id, admin, halamanLengkap("1")); kode != http.StatusOK {
		t.Fatalf("admin: %d %s", kode, isi)
	}
	if kode, isi := u.kirim(id, secHead, putusan("1")); kode != http.StatusOK {
		t.Fatalf("Sec Head: %d %s", kode, isi)
	}
	kode, isi := u.panggil("POST", "/kasus/"+id+"/nomor-polis", deptHead, map[string]any{"halaman": putusan("1")})
	if kode != http.StatusOK {
		t.Fatalf("nomor polis tetap dibentuk (Obj-Save WithErrors): %d %s", kode, isi)
	}
	kode, isi = u.kirim(id, deptHead, putusan("1"))
	if kode != http.StatusUnprocessableEntity || !strings.Contains(isi, "Cannot fetch Treaty Group ID, Contact IT") {
		t.Fatalf("submit Dept Head TreatyGroupID kosong: %d %s", kode, isi)
	}
	if k := u.g.Kasus[id]; k.PositionNote != models.PosisiDeptHead || k.StatusWork == models.StatusSelesai {
		t.Fatalf("berkas tertahan di Dept Head: %+v", k)
	}

	// TreatyGroupID terisi: langkah 3 dilewati, submit berjalan.
	u.g.Halaman[id].Setel("PolicyTreatyIn.TreatyGroupID", "UJI-GRUP")
	if kode, isi := u.kirim(id, deptHead, putusan("1")); kode != http.StatusOK || u.g.Kasus[id].StatusWork != models.StatusSelesai {
		t.Fatalf("TreatyGroupID terisi: %d %s", kode, isi)
	}
}
