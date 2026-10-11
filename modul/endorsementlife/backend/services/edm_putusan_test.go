package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/services"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

// gudangPutusan - polis NB aplikasi `UJI-PL-1` (versi 1) dan kasus `EDML-1`
// (versi 2, Type TR) yang sudah disimpan: satu peserta Delete, satu Old, satu New.
func gudangPutusan(t *testing.T) (*tiruan.Gudang, *services.Layanan, *jejakTiruan) {
	t.Helper()
	g := tiruan.Baru()
	// PRODKE new business = 0 (keputusan work owner 08-10-2026, migrasi 485).
	g.Polis["NBLF-7"] = &tiruan.Polis{ID: "NBLF-7", ProdKe: 0, Kepala: map[string]string{"TYPE": models.TypeTR}}
	g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: "UJI-NB-D1", PolisID: "NBLF-7", PLNumber: "UJI-PL-1", Nilai: map[string]string{}})
	g.Polis["EDML-1"] = &tiruan.Polis{ID: "EDML-1", OldPolicyNo: "UJI-PL-1", EdmType: models.EdmTypePerubahanData, ProdKe: 1,
		Kepala: map[string]string{"TYPE": models.TypeTR, "BUSINESS_NAME": "UJI-COB"}}
	for _, p := range []struct{ id, status string }{{"UJI-D1", models.StatusOld}, {"UJI-D2", models.StatusOld}, {"UJI-D3", models.StatusNew}} {
		g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: p.id, PolisID: "EDML-1", PLNumber: "UJI-PL-1", EdmStatus: p.status,
			Nilai: map[string]string{"CURRENCY": "IDR", "GROSS_PREMIUM_REFUND_RETRO": "10"}})
	}
	l, j := layananJejak(g)
	if _, err := l.Simpan(context.Background(), pelakuUji, "EDML-1", models.PilihanHapus{Pilih: []string{"UJI-D1"}}); err != nil {
		t.Fatal(err)
	}
	g.Panggilan, j.catatan = nil, nil
	return g, l, j
}

func TestPutuskanConfirmMeresmikanVersiDalamUrutan(t *testing.T) {
	g, l, j := gudangPutusan(t)
	ctx := context.Background()
	h, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1", Komentar: " UJI-OK "})
	if err != nil {
		t.Fatal(err)
	}
	if h.Status != models.StatusKasusSelesai || h.NoEndorsement != "UJI-PL-1/01" || h.Peserta != 3 || h.RekapWarisan != 1 {
		t.Fatalf("hasil %+v", h)
	}
	// AC 42: urutan pemanggilan adalah bagian kebenaran. JSON_POLIS sesudah Resmikan (NO_ENDORS terbit), 08-10-2026.
	if got := strings.Join(g.Panggilan, ","); got != "SisipRiwayat,AdaVersiResmi,TulisProduksiWarisan,Resmikan,TulisJSONPolis,TulisPesertaWarisan,TulisRekapWarisan" {
		t.Errorf("urutan %s", got)
	}
	p := g.Polis["EDML-1"]
	if p.NoPolis != "UJI-PL-1" || p.ProdKe != 1 || p.NoEndors != "UJI-PL-1/01" || p.PLNumberEDM != "UJI-PL-1/01" || p.Status != models.StatusKasusSelesai {
		t.Errorf("kepala %+v", p)
	}
	for _, d := range g.Peserta[1:] {
		mauLama := "0"
		if d.EdmStatus == models.StatusOld {
			mauLama = "1"
		}
		if d.Nilai["PL_NUMBER_EDM"] != "UJI-PL-1/01" || d.Nilai["STATUS"] != "1" || d.Nilai["STATUS_OLD"] != mauLama {
			t.Errorf("peserta %s %+v (STATUS 1 = TR, STATUS_OLD 1 = Old)", d.ID, d.Nilai)
		}
	}
	w := g.RekapWarisan[0]
	if w["PL_NUMBER"] != "UJI-PL-1" || w["PL_NUMBER_EDM"] != "UJI-PL-1/01" || w["IDPEGA"] != "EDML-1" || w["COB"] != "UJI-COB" ||
		w["CURRENCY"] != "IDR" || w["PREMIUM"] != "10" {
		t.Errorf("rekap warisan %+v", w)
	}
	r := g.RiwayatKasus["EDML-1"]
	if len(r) != 1 || r[0].Status != "Accept" || r[0].PIC != pelakuUji.AkunID || r[0].Komentar != "UJI-OK" || r[0].Tanggal != "2026-10-01 17:00:00" { // 10:00 UTC = 17:00 WIB (penomor.DiJakarta)
		t.Errorf("riwayat %+v", r)
	}
	if len(j.catatan) != 1 || j.catatan[0].Dari != models.TahapInputEDMLife || j.catatan[0].Ke != models.StatusKasusSelesai {
		t.Errorf("jejak %+v", j.catatan)
	}
	if k, err := l.BacaKasus(ctx, pelakuUji, "EDML-1"); err != nil || len(k.Riwayat) != 1 || k.PLNumberEDM != "UJI-PL-1/01" {
		t.Errorf("baca sesudah Confirm %+v %v", k, err)
	}
	// Kasus tertutup tidak dapat diputuskan ulang (AC 29) - nomor tidak lahir dua kali.
	if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"}); !errors.Is(err, services.ErrKasusTertutup) {
		t.Fatalf("putuskan ulang %v", err)
	}
	if v, _, _ := g.VersiBerjalan(ctx, nil, "UJI-PL-1", 0); v.ID != "EDML-1" || v.ProdKe != 1 {
		t.Errorf("versi berjalan sesudah Confirm %+v", v)
	}
}

func TestPutuskanDeclineMenutupTanpaVersi(t *testing.T) {
	for _, status := range []string{"2", "7"} {
		g, l, j := gudangPutusan(t)
		h, err := l.Putuskan(context.Background(), pelakuUji, "EDML-1", services.MasukanPutusan{Status: status, Komentar: "UJI-NO"})
		if err != nil || h.Status != models.StatusKasusDitolak || h.NoEndorsement != "" {
			t.Fatalf("%s: %+v %v", status, h, err)
		}
		if got := strings.Join(g.Panggilan, ","); got != "SisipRiwayat,Tolak" {
			t.Errorf("%s: urutan %s", status, got)
		}
		if p := g.Polis["EDML-1"]; p.Status != models.StatusKasusDitolak || p.NoPolis != "" || p.PLNumberEDM != "" {
			t.Errorf("%s: kepala %+v", status, p)
		}
		if r := g.RiwayatKasus["EDML-1"]; len(r) != 1 || r[0].Status != "Decline" {
			t.Errorf("%s: riwayat %+v", status, r)
		}
		if len(j.catatan) != 1 || j.catatan[0].Ke != models.StatusKasusDitolak || j.catatan[0].Komentar != "UJI-NO" {
			t.Errorf("%s: jejak %+v", status, j.catatan)
		}
		// AC 31: polis bebas di-endorse ulang.
		if k, err := l.Kelayakan(context.Background(), pelakuUji, services.MasukanKelayakan{NomorPolis: "UJI-PL-1", EdmType: "1"}); err != nil || !k.Boleh {
			t.Errorf("%s: kelayakan sesudah Decline %+v %v", status, k, err)
		}
	}
}

func TestPutuskanAntiDobel(t *testing.T) {
	ctx := context.Background()
	// Versi 2 resmi lahir sesudah kasus dibuat: salinan kasus basi.
	g, l, _ := gudangPutusan(t)
	g.Polis["EDML-0"] = &tiruan.Polis{ID: "EDML-0", NoPolis: "UJI-PL-1", OldPolicyNo: "UJI-PL-1", EdmType: "1", ProdKe: 1,
		PLNumberEDM: "UJI-PL-1/01", Status: models.StatusKasusSelesai}
	if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"}); !errors.Is(err, services.ErrVersiBerubah) {
		t.Fatalf("versi berubah: %v", err)
	}
	// (NO_POLIS, PROD_KE) sudah dipakai baris yang TIDAK dihitung versi berjalan -
	// mis. ditulis jalur lain di luar modul ini. Pertahanan berlapis: di jalur
	// modul ini sendiri dobel sudah dicegah UX_PL_EDM_TERBUKA + FOR UPDATE.
	g, l, _ = gudangPutusan(t)
	g.Polis["UJI-LAIN"] = &tiruan.Polis{ID: "UJI-LAIN", NoPolis: "UJI-PL-1", EdmType: "1", ProdKe: 1, Status: models.StatusKasusDitolak}
	if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"}); !errors.Is(err, services.ErrVersiBerubah) {
		t.Fatalf("dobel: %v", err)
	}
	for _, n := range g.Panggilan {
		if n == "Resmikan" || n == "TulisRekapWarisan" {
			t.Errorf("%s dipanggil sesudah penjaga menolak", n)
		}
	}
}

func TestPutuskanNomorLahirSekali(t *testing.T) {
	g, l, _ := gudangPutusan(t)
	g.Polis["EDML-1"].PLNumberEDM = "UJI-NOMOR-LAMA"
	h, err := l.Putuskan(context.Background(), pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"})
	if err != nil || h.NoEndorsement != "UJI-NOMOR-LAMA" {
		t.Fatalf("%+v %v", h, err)
	}
}

func TestPutuskanDitolak(t *testing.T) {
	ctx := context.Background()
	g := gudangKasusCSV(models.EdmTypePerubahanData) // belum Save
	l, _ := layananJejak(g)
	if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"}); !errors.Is(err, services.ErrBelumDisimpan) {
		t.Errorf("belum Save: %v", err)
	}
	if len(g.RiwayatKasus["EDML-1"]) != 0 {
		t.Error("riwayat tertulis walau ditolak")
	}
	_, l2, _ := gudangPutusan(t)
	for _, m := range []services.MasukanPutusan{{Status: "3"}, {Status: ""}, {Status: "Accept"}, {Status: "1", Komentar: strings.Repeat("X", 256)}} {
		if _, err := l2.Putuskan(ctx, pelakuUji, "EDML-1", m); !errors.Is(err, services.ErrMasukanTidakSah) {
			t.Errorf("%+v: %v", m, err)
		}
	}
	if _, err := l2.Putuskan(ctx, pelakuUji, "NBLF-7", services.MasukanPutusan{Status: "1"}); !errors.Is(err, services.ErrKasusTidakAda) {
		t.Errorf("bukan kasus EDM: %v", err)
	}
}

// PRODKE SATU hitungan mulai 0 (keputusan work owner 08-10-2026): NB (warisan PRODKE kosong = 0) → endorsement
// pertama PRODKE 1 `<polis>/01`, kedua PRODKE 2 `/02`. Rantai Pega yang sudah punya endorsement warisan PRODKE 1
// berlanjut ke PRODKE 2 `/02`.
func TestNomorEndorsementPertamaDanKeduaPolisNBWarisan(t *testing.T) {
	ctx := context.Background()
	g := tiruan.Baru()
	g.PolisWarisan = append(g.PolisWarisan, &tiruan.PolisWarisan{IDPega: "UJI-IDPEGA-NB", NoPolis: "UJI-PL-W", ProdKe: 0})
	l, _ := layananJejak(g)
	kasus := func(id string, prodKe int) {
		g.Polis[id] = &tiruan.Polis{ID: id, OldPolicyNo: "UJI-PL-W", EdmType: models.EdmTypePerubahanData, ProdKe: prodKe,
			Kepala: map[string]string{"TYPE": models.TypeTR, "BUSINESS_NAME": "UJI-COB"}}
		g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: "UJI-P-" + id, PolisID: id, PLNumber: "UJI-PL-W", EdmStatus: models.StatusNew,
			Nilai: map[string]string{"CURRENCY": "IDR"}})
		if _, err := l.Simpan(ctx, pelakuUji, id, models.PilihanHapus{}); err != nil {
			t.Fatalf("simpan %s: %v", id, err)
		}
	}
	kasus("EDML-11", 1)
	h, err := l.Putuskan(ctx, pelakuUji, "EDML-11", services.MasukanPutusan{Status: "1"})
	if err != nil || h.NoEndorsement != "UJI-PL-W/01" || g.Polis["EDML-11"].ProdKe != 1 {
		t.Fatalf("endorsement pertama NB warisan: %+v %v (PRODKE %d)", h, err, g.Polis["EDML-11"].ProdKe)
	}
	kasus("EDML-12", 2)
	h, err = l.Putuskan(ctx, pelakuUji, "EDML-12", services.MasukanPutusan{Status: "1"})
	if err != nil || h.NoEndorsement != "UJI-PL-W/02" || g.Polis["EDML-12"].ProdKe != 2 {
		t.Fatalf("endorsement kedua: %+v %v (PRODKE %d)", h, err, g.Polis["EDML-12"].ProdKe)
	}
	// JSON_POLIS sekali per kasus, PRODKE baru + NOENDORS (08-10-2026).
	if j := g.JSONPolisTertulis["EDML-12"]; j["PRODKE"] != "2" || j["NOENDORS"] != "UJI-PL-W/02" || j["NOPOLIS"] != "UJI-PL-W" {
		t.Errorf("JSON_POLIS kasus kedua = %+v", j)
	}

	// Rantai Pega: NB (kosong = 0) + endorsement warisan PRODKE 1 → endorsement sistem baru PRODKE 2 `/02`.
	g = tiruan.Baru()
	g.PolisWarisan = append(g.PolisWarisan, &tiruan.PolisWarisan{IDPega: "UJI-IDPEGA-A", NoPolis: "UJI-PL-W", ProdKe: 0,
		TglInput: "2026-01-01 08:00:00"},
		&tiruan.PolisWarisan{IDPega: "UJI-IDPEGA-Z", NoPolis: "UJI-PL-W", ProdKe: 1, EdmType: "1", TglInput: "2026-03-01 08:00:00"})
	l, _ = layananJejak(g)
	kasus("EDML-13", 2)
	if v, _, _ := g.VersiBerjalan(ctx, nil, "UJI-PL-W", 0); v.ID != "UJI-IDPEGA-Z" || v.ProdKe != 1 {
		t.Fatalf("versi berjalan rantai Pega: %+v", v)
	}
	if h, err := l.Putuskan(ctx, pelakuUji, "EDML-13", services.MasukanPutusan{Status: "1"}); err != nil || h.NoEndorsement != "UJI-PL-W/02" {
		t.Errorf("sesudah endorsement Pega PRODKE 1: %+v %v", h, err)
	}
}

// TestVersiAbaikanBarisJSONPolisSistemBaru - JSON_POLIS kini juga memuat baris PremiumList (`NBLF-`, PRODKE '0')
// dan endorsement modul ini (`EDML-`): versi dipilih dari jalur APLIKASI (peserta di tabel aplikasi), bukan
// dari baris warisan yang hanya rekam.
func TestVersiAbaikanBarisJSONPolisSistemBaru(t *testing.T) {
	ctx := context.Background()
	g := tiruan.Baru()
	g.Polis["NBLF-5"] = &tiruan.Polis{ID: "NBLF-5", ProdKe: 0, Kepala: map[string]string{"TYPE": "QR"}}
	g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: "UJI-NB5-D1", PolisID: "NBLF-5", PLNumber: "UJI-PL-X"})
	g.PolisWarisan = append(g.PolisWarisan,
		&tiruan.PolisWarisan{IDPega: "NBLF-5", NoPolis: "UJI-PL-X", ProdKe: 0},  // ditulis Confirm PremiumList
		&tiruan.PolisWarisan{IDPega: "EDML-20", NoPolis: "UJI-PL-X", ProdKe: 1}) // ditulis Confirm endorsement
	v, ada, err := g.VersiBerjalan(ctx, nil, "UJI-PL-X", 0)
	if err != nil || !ada || v.ID != "NBLF-5" || v.Jenis != models.SumberAplikasi || v.ProdKe != 0 {
		t.Fatalf("NB aplikasi: %+v %v %v", v, ada, err)
	}
	g.Polis["EDML-20"] = &tiruan.Polis{ID: "EDML-20", NoPolis: "UJI-PL-X", OldPolicyNo: "UJI-PL-X", EdmType: "1", ProdKe: 1,
		PLNumberEDM: "UJI-PL-X/01", Status: models.StatusKasusSelesai}
	if v, _, _ := g.VersiBerjalan(ctx, nil, "UJI-PL-X", 0); v.ID != "EDML-20" || v.Jenis != models.SumberAplikasi || v.ProdKe != 1 {
		t.Errorf("endorsement aplikasi PRODKE 1: %+v", v)
	}
}

// TestPutuskanMenulisJSONPolisTglProd - TGL_PROD (08-10-2026): bulan Confirm = periode → tanggal Confirm;
// periode bergeser (lewat TANGGAL_CLOSING) → tanggal 1 bulan periode. TANGGAL_CLOSING kosong = Confirm batal.
func TestPutuskanMenulisJSONPolisTglProd(t *testing.T) {
	ctx := context.Background()
	for _, k := range []struct {
		jam  time.Time
		hari int
		mau  string
	}{
		{time.Date(2026, 10, 8, 3, 0, 0, 0, time.UTC), 25, "2026-10-08"},  // 08-10 10:00 WIB, ≤ 25
		{time.Date(2026, 10, 27, 3, 0, 0, 0, time.UTC), 25, "2026-11-01"}, // 27-10 > 25 → periode 11.2026
	} {
		g, _, j := gudangPutusan(t)
		g.HariTutupBuku = k.hari
		l := services.BaruLayanan(g, tiruan.Transaksi, nil).DenganJejak(j).DenganJam(func() time.Time { return k.jam })
		h, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"})
		if err != nil || h.NoEndorsement != "UJI-PL-1/01" {
			t.Fatalf("%v: %+v %v", k.jam, h, err)
		}
		b := g.JSONPolisTertulis["EDML-1"]
		if b["TGL_PROD"] != k.mau || b["PRODKE"] != "1" || b["NOENDORS"] != "UJI-PL-1/01" || b["IDPEGA"] != "EDML-1" ||
			b["USERNAME"] != g.Polis["EDML-1"].Pembuat {
			t.Errorf("%v: JSON_POLIS = %+v, mau TGL_PROD %s", k.jam, b, k.mau)
		}
	}
	g, l, _ := gudangPutusan(t)
	g.HariTutupBuku = 0
	if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"}); err == nil ||
		g.Polis["EDML-1"].Status == models.StatusKasusSelesai && len(g.JSONPolisTertulis) != 0 {
		t.Errorf("TANGGAL_CLOSING kosong tidak membatalkan Confirm: %v", err)
	}
}

// K4 keputusan work owner 01-10-2026 (OQ-EDM-010): Confirm menulis SATU baris `LIFEINPRODUCTION` di transaksi
// yang sama, SEBELUM peserta diresmikan - urutan `InsertJsonPolisLife_Act` 10 b2663 (tanpa prakondisi, selalu
// jalan) lalu 11/12; Decline tidak menulis; gagal menulis = transaksi batal.
func TestPutuskanMenulisProduksiWarisan(t *testing.T) {
	g, l, _ := gudangPutusan(t)
	ctx := context.Background()
	g.Polis["EDML-1"].Kepala["WPC"] = "2026-10-01"
	g.Polis["EDML-1"].Kepala["PRODUCT_NAME_ID"] = "UJI-PROD-ID"
	g.Polis["EDML-1"].Kepala["PRODUCT_NAME"] = "UJI-PRODUK"
	if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1", Komentar: "  UJI-KOMENTAR  "}); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(g.Panggilan, ","); got != "SisipRiwayat,AdaVersiResmi,TulisProduksiWarisan,Resmikan,TulisJSONPolis,TulisPesertaWarisan,TulisRekapWarisan" {
		t.Errorf("urutan %s", got)
	}
	if len(g.ProduksiWarisan) != 1 {
		t.Fatalf("baris LIFEINPRODUCTION: %d", len(g.ProduksiWarisan))
	}
	b := g.ProduksiWarisan[0]
	if b["IDPEGA"] != "EDML-1" || b["NOPOLIS"] != "UJI-PL-1" || b["NOENDORS"] != "UJI-PL-1/01" || b["TYPE"] != models.TypeTR ||
		b["BUSINESSNAME"] != "UJI-COB" {
		t.Errorf("baris %+v", b)
	}
	// DESCRIPTION ← Comment keputusan, di-trim (keputusan work owner 09-10-2026).
	if b["DESCRIPTION"] != "UJI-KOMENTAR" {
		t.Errorf("DESCRIPTION %q", b["DESCRIPTION"])
	}
	// WPC / PRODUCTNAMEID / PRODUCTNAME ← kepala kasus (keputusan work owner 09-10-2026).
	if b["WPC"] != "2026-10-01" || b["PRODUCTNAMEID"] != "UJI-PROD-ID" || b["PRODUCTNAME"] != "UJI-PRODUK" {
		t.Errorf("WPC/PRODUCT %q %q %q", b["WPC"], b["PRODUCTNAMEID"], b["PRODUCTNAME"])
	}
	g, l, _ = gudangPutusan(t)
	if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"}); err != nil || g.ProduksiWarisan[0]["DESCRIPTION"] != "" {
		t.Errorf("Comment kosong → DESCRIPTION kosong (NULL): %v %+v", err, g.ProduksiWarisan)
	}

	g, l, _ = gudangPutusan(t)
	if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "2"}); err != nil || len(g.ProduksiWarisan) != 0 {
		t.Errorf("Decline tidak menulis LIFEINPRODUCTION: %v %d", err, len(g.ProduksiWarisan))
	}
}

// K5 keputusan work owner 01-10-2026 (OQ-EDM-016): Confirm menulis SETIAP peserta kasus - Old, New, Delete -
// ke `M_LIFE_PREMIUM_DETAIL` (`SaveMasterLPDet` 11.6 b5098, PRE=false b5118 → selalu jalan) sesudah 11.2–11.5
// (STATUS/STATUSOLD) dan sebelum rekap 12; baris Delete/Batal tersimpan tetapi TIDAK lolos penyaring Claim Life.
func TestPutuskanMenulisPesertaWarisan(t *testing.T) {
	g, l, _ := gudangPutusan(t)
	ctx := context.Background()
	h, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(g.Panggilan, ","); got != "SisipRiwayat,AdaVersiResmi,TulisProduksiWarisan,Resmikan,TulisJSONPolis,TulisPesertaWarisan,TulisRekapWarisan" {
		t.Errorf("urutan %s", got)
	}
	if len(g.PesertaWarisanTertulis) != 3 || h.PesertaWarisan != 3 {
		t.Fatalf("baris M_LIFE_PREMIUM_DETAIL: %d (hasil %d)", len(g.PesertaWarisanTertulis), h.PesertaWarisan)
	}
	lolos := map[string]bool{}
	for _, b := range g.PesertaWarisanTertulis {
		if b["PL_NUMBER"] != "UJI-PL-1" || b["PL_NUMBER_EDM"] != "UJI-PL-1/01" || b["IDPEGA"] != "EDML-1" || b["STATUS"] != "1" {
			t.Errorf("baris %+v", b)
		}
		mauLama := "0"
		if b["EDMSTATUS"] == models.StatusOld {
			mauLama = "1"
		}
		if b["STATUSOLD"] != mauLama {
			t.Errorf("STATUSOLD: %+v", b)
		}
		// Penyaring Claim Life DITIRU lewat `models.StatusHidup` - kesetaraannya dengan teks `penyaringHidup`
		// Claim Life dijaga `models.TestStatusPesertaSejalanPenyaringClaimLife` (dibaca dari sumbernya, tidak diimpor).
		lolos[b["EDMSTATUS"]] = models.StatusHidup(b["EDMSTATUS"])
	}
	if !lolos[models.StatusOld] || !lolos[models.StatusNew] || lolos[models.StatusDelete] {
		t.Errorf("penyaring Claim Life atas baris tertulis: %v", lolos)
	}
	for _, v := range []string{models.StatusBatal, " Batal ", models.StatusDelete} {
		if models.StatusHidup(v) {
			t.Errorf("%q lolos penyaring Claim Life", v)
		}
	}
	if !models.StatusHidup("") {
		t.Error("peserta new business (EDMSTATUS NULL) harus lolos")
	}
	// Decline: nol baris warisan.
	g, l, _ = gudangPutusan(t)
	if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "2"}); err != nil || len(g.PesertaWarisanTertulis) != 0 {
		t.Errorf("Decline: %v %d", err, len(g.PesertaWarisanTertulis))
	}
}

// TestPutuskanAnnuityRefundDanCadangan - LIFEINPRODUCTION ANNUITYINTEREST / PREMIUMREFUNDFACTOR (09-10-2026):
// dari kepala kasus; kasus lama tanpa nilai memakai kepala versi yang di-endorse.
func TestPutuskanAnnuityRefundDanCadangan(t *testing.T) {
	ctx := context.Background()
	g, l, _ := gudangPutusan(t)
	g.Polis["EDML-1"].Kepala["ANNUITY_INTEREST"] = "0.05"
	g.Polis["EDML-1"].Kepala["PREMIUM_REFUND_FACTOR"] = "1.25"
	g.Polis["NBLF-7"].Kepala["ANNUITY_INTEREST"] = "9"
	if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"}); err != nil {
		t.Fatal(err)
	}
	if b := g.ProduksiWarisan[0]; b["ANNUITYINTEREST"] != "0.05" || b["PREMIUMREFUNDFACTOR"] != "1.25" {
		t.Errorf("dari kepala kasus: %q %q", b["ANNUITYINTEREST"], b["PREMIUMREFUNDFACTOR"])
	}

	// Kasus lama (kepala kosong): cadangan dari versi NBLF-7 yang di-endorse.
	g, l, _ = gudangPutusan(t)
	g.Polis["NBLF-7"].Kepala["ANNUITY_INTEREST"] = "0.04"
	g.Polis["NBLF-7"].Kepala["PREMIUM_REFUND_FACTOR"] = "2"
	if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"}); err != nil {
		t.Fatal(err)
	}
	if b := g.ProduksiWarisan[0]; b["ANNUITYINTEREST"] != "0.04" || b["PREMIUMREFUNDFACTOR"] != "2" {
		t.Errorf("cadangan versi: %q %q", b["ANNUITYINTEREST"], b["PREMIUMREFUNDFACTOR"])
	}
	// Kasus lama tidak diubah (tanpa isi mundur).
	if g.Polis["EDML-1"].Kepala["ANNUITY_INTEREST"] != "" {
		t.Error("kepala kasus lama ikut ditulis")
	}
}
