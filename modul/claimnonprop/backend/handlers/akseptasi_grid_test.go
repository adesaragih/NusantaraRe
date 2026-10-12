package handlers_test

import (
	"net/http"
	"testing"

	"nusantarare/modul/claimnonprop/backend/models"
	"nusantarare/modul/claimnonprop/backend/services"
)

// akseptasiSiap - kasus di Input Acceptation dengan satu akseptasi (persiapan `TestAkseptasiSampaiKasusKomite`).
func (u *uji) akseptasiSiap() string {
	u.t.Helper()
	id := u.kirimKeTeknik()
	peran := models.WorkbasketAcceptation
	kode, out := u.aksi(id, teknik, peran, "SaveToOS", 0, "", nil)
	u.wajib(kode, http.StatusOK, out, "save to OS")
	h := u.g.Halaman(id)
	h.Setel(models.CD+"AppointedADJID", "UJI-ADJ")
	h.Setel(models.CD+"ConsultantID", "UJI-ADJ")
	h.Setel(models.CD+"QuotationData.BusinessOldId", "22")
	u.g.SetelHalaman(id, h)
	kode, out = u.aksi(id, teknik, peran, "AddAkseptasiCNP", 0, "", nil)
	u.wajib(kode, http.StatusOK, out, "add akseptasi")
	return id
}

// ⭐ 12-10-2026 (keputusan pemakai, penyimpangan OQ-CNP-33): grid Claim Acceptation dan Loss Allocation panel
// akseptasi bisa ditambah dan menghitung ulang Claim Amount in IDR / Cedant serta Loss Allocation.
func TestGridAkseptasiTambahDanHitungUlang(t *testing.T) {
	u := baruUji(t)
	id := u.akseptasiSiap()
	peran := models.WorkbasketAcceptation
	klaim := func(n int, p string) string {
		return models.JalurAnak(models.JalurAdj(1, models.AnakClaimAccept), n, p)
	}
	loss := func(n int, p string) string {
		return models.JalurAnak(models.JalurAdj(1, models.AnakLossAlloc), n, p)
	}

	// Angka baris klaim akseptasi 1 diubah -> IDR / Cedant dihitung ulang (Share Ceding 100, kurs 1).
	kode, out := u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "HitungKlaimAkseptasi", Indeks: 1, Baris: 1,
		Masukan: map[string]string{klaim(1, "Value"): "6668618965", klaim(1, "TPL"): "0", klaim(1, "Salvage"): "100000",
			klaim(1, "AdjusterFee"): "5000", klaim(1, "CNPOthersFee"): "0"}})
	u.wajib(kode, http.StatusOK, out, "hitung klaim akseptasi")
	h := u.g.Halaman(id)
	r := h.AmbilDaftar(models.JalurAdj(1, models.AnakClaimAccept))[0]
	harap := "6668523965" // 6.668.618.965 + 5.000 - 100.000
	if !models.SamaAngka(r["USD"], harap) || !models.SamaAngka(r["ClaimAmountCedant"], harap) {
		t.Fatalf("IDR %q / Cedant %q, harap %s (baris %v)", r["USD"], r["ClaimAmountCedant"], harap, r)
	}

	// Add Claim Acceptation -> baris kosong; Currency IDR -> kurs 1.
	kode, out = u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "TambahKlaimAkseptasi", Indeks: 1})
	u.wajib(kode, http.StatusOK, out, "tambah klaim akseptasi")
	if n := len(u.g.Halaman(id).AmbilDaftar(models.JalurAdj(1, models.AnakClaimAccept))); n != 2 {
		t.Fatalf("baris klaim akseptasi = %d, harap 2", n)
	}
	kode, out = u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "SetCurrencyAkseptasi", Indeks: 1, Baris: 2,
		Masukan: map[string]string{klaim(2, "CurrencyID"): idr}})
	u.wajib(kode, http.StatusOK, out, "mata uang klaim akseptasi")
	r2 := u.g.Halaman(id).AmbilDaftar(models.JalurAdj(1, models.AnakClaimAccept))[1]
	if r2["Currency"] != "IDR" || r2["AltValue"] != "1" {
		t.Fatalf("baris baru = %v", r2)
	}

	// Add Loss Allocation -> satu baris per mata uang klaim; Share 50 -> separuh klaim baris IDR terakhir.
	awal := len(u.g.Halaman(id).AmbilDaftar(models.JalurAdj(1, models.AnakLossAlloc)))
	kode, out = u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "TambahLossAkseptasi", Indeks: 1})
	u.wajib(kode, http.StatusOK, out, "tambah loss akseptasi")
	ls := u.g.Halaman(id).AmbilDaftar(models.JalurAdj(1, models.AnakLossAlloc))
	if len(ls) != awal+1 || ls[awal]["Currency"] != "IDR" {
		t.Fatalf("loss akseptasi = %v", ls)
	}
	kode, out = u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "HitungLossAkseptasi", Indeks: 1, Baris: awal + 1,
		Masukan: map[string]string{loss(awal+1, "ClaimPercentage"): "50"}})
	u.wajib(kode, http.StatusOK, out, "share loss akseptasi")
	l := u.g.Halaman(id).AmbilDaftar(models.JalurAdj(1, models.AnakLossAlloc))[awal]
	// Baris klaim IDR TERAKHIR = baris baru (Value kosong) -> 0.
	if !models.SamaAngka(l["ClaimAmountAdjust"], "0") {
		t.Fatalf("loss baru = %v", l)
	}
}

// ⭐ 12-10-2026 (laporan pemakai): Claim Amount Loss Allocation akseptasi berubah -> XOL Allocation akseptasi disusun
// ulang per layer master (UR = deductible 2,4 M; XL 1st sampai limit 3,1 M; sisanya XL 2nd), RNM Share 30 %.
func TestLossAkseptasiMenyusunUlangXOL(t *testing.T) {
	u := baruUji(t)
	id := u.akseptasiSiap()
	peran := models.WorkbasketAcceptation
	loss := func(p string) string { return models.JalurAnak(models.JalurAdj(1, models.AnakLossAlloc), 1, p) }
	kode, out := u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "HitungLossJumlahAkseptasi", Indeks: 1, Baris: 1,
		Masukan: map[string]string{loss("ClaimAmountAdjust"): "4000000000", loss("CNPFlagXOL"): "true"}})
	u.wajib(kode, http.StatusOK, out, "claim amount loss akseptasi")
	xol := u.g.Halaman(id).AmbilDaftar(models.JalurAdj(1, models.AnakXOL))
	if len(xol) != 3 {
		t.Fatalf("XOL akseptasi = %d baris: %v", len(xol), xol)
	}
	harap := [][2]string{{"2400000000", "0"}, {"1600000000", "480000000"}, {"0", "0"}}
	for i, w := range harap {
		x := xol[i]
		if !models.SamaAngka(x["ClaimAmountAdjust"], w[0]) || !models.SamaAngka(x["ClaimSpreaded"], w[1]) {
			t.Errorf("layer %d (%s) = alokasi %s / RNM %s, harap %v", i+1, x["TreatyName"], x["ClaimAmountAdjust"],
				x["ClaimSpreaded"], w)
		}
		// Gross Value OS = ClaimEstimation + AdjClaimValue = alokasi baru.
		var kal models.Kalkulator
		if g := models.Teks(kal.Tambah(kal.B(x, "ClaimEstimation"), kal.B(x, "AdjClaimValue"))); !models.SamaAngka(g, w[0]) {
			t.Errorf("layer %d estimasi + penyesuaian = %s, harap %s", i+1, g, w[0])
		}
	}
	// Spreading In ikut: Σ ClaimSpreaded layer non-UR × share spreading.
	if in := u.g.Halaman(id).AmbilDaftar(models.JalurAdj(1, models.AnakSpreadIn)); len(in) == 0 ||
		!models.SamaAngka(in[0]["ClaimSpreaded"], "480000000") {
		t.Errorf("Spreading In = %v", in)
	}
}

// ⭐ 12-10-2026 (tangkapan Pega akseptasi): layer XOL akseptasi memakai rumus AdjClaimCNP - Claim Amount layer NETO,
// Claim Amount RNM dari nilai BRUTO (fee kembali ke 100 %). Master uji: UR 2,4 M, XL 1st limit 3,1 M, RNM 30 %,
// MDP 1.000, reinstatement 100 %. Klaim 3 M, Salvage 256 jt -> neto 2,744 M: UR 2,4 M, XL 1st 344 jt; bruto
// 344 + 256 = 600 jt; RNM 180 jt; Salvage RNM 76,8 jt; Total Claim RNM 103,2 jt.
func TestXOLAkseptasiRumusAdjClaimCNP(t *testing.T) {
	u := baruUji(t)
	id := u.akseptasiSiap()
	peran := models.WorkbasketAcceptation
	klaim := func(p string) string { return models.JalurAnak(models.JalurAdj(1, models.AnakClaimAccept), 1, p) }
	loss := func(p string) string { return models.JalurAnak(models.JalurAdj(1, models.AnakLossAlloc), 1, p) }
	kode, out := u.aksiP(id, teknik, peran, services.PermintaanAksi{Aksi: "HitungKlaimAkseptasi", Indeks: 1, Baris: 1,
		Masukan: map[string]string{klaim("Value"): "3000000000", klaim("TPL"): "0", klaim("AdjusterFee"): "0",
			klaim("Salvage"): "256000000", klaim("CNPOthersFee"): "0", loss("ClaimPercentage"): "100",
			loss("CNPFlagXOL"): "true"}})
	u.wajib(kode, http.StatusOK, out, "hitung klaim akseptasi")
	xol := u.g.Halaman(id).AmbilDaftar(models.JalurAdj(1, models.AnakXOL))
	if len(xol) < 2 {
		t.Fatalf("XOL = %v", xol)
	}
	cek := func(i int, p, harap string) {
		t.Helper()
		if !models.SamaAngka(xol[i][p], harap) {
			t.Errorf("layer %d %s = %q, harap %s", i+1, p, xol[i][p], harap)
		}
	}
	cek(0, "TotalClaim", "2400000000")
	cek(1, "TotalClaim", "344000000")
	cek(1, "ClaimAmountAdjust", "600000000")
	cek(1, "ClaimSpreaded", "180000000")
	cek(1, "Salvage", "76800000")
	cek(1, "TotalClaimRNM", "103200000")
}
