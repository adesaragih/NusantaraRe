package services_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/services"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

// gudangPutusan - polis NB aplikasi `UJI-PL-1` (versi 1) dan kasus `EDMLF-1`
// (versi 2, Type TR) yang sudah disimpan: satu peserta Delete, satu Old, satu New.
func gudangPutusan(t *testing.T) (*tiruan.Gudang, *services.Layanan, *jejakTiruan) {
	t.Helper()
	g := tiruan.Baru()
	g.Polis["NBLF-7"] = &tiruan.Polis{ID: "NBLF-7", ProdKe: 1, Kepala: map[string]string{"TYPE": models.TypeTR}}
	g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: "UJI-NB-D1", PolisID: "NBLF-7", PLNumber: "UJI-PL-1", Nilai: map[string]string{}})
	g.Polis["EDMLF-1"] = &tiruan.Polis{ID: "EDMLF-1", OldPolicyNo: "UJI-PL-1", EdmType: models.EdmTypePerubahanData, ProdKe: 2,
		Kepala: map[string]string{"TYPE": models.TypeTR, "BUSINESS_NAME": "UJI-COB"}}
	for _, p := range []struct{ id, status string }{{"UJI-D1", models.StatusOld}, {"UJI-D2", models.StatusOld}, {"UJI-D3", models.StatusNew}} {
		g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: p.id, PolisID: "EDMLF-1", PLNumber: "UJI-PL-1", EdmStatus: p.status,
			Nilai: map[string]string{"CURRENCY": "IDR", "GROSS_PREMIUM_REFUND_RETRO": "10"}})
	}
	l, j := layananJejak(g)
	if _, err := l.Simpan(context.Background(), pelakuUji, "EDMLF-1", models.PilihanHapus{Pilih: []string{"UJI-D1"}}); err != nil {
		t.Fatal(err)
	}
	g.Panggilan, j.catatan = nil, nil
	return g, l, j
}

func TestPutuskanConfirmMeresmikanVersiDalamUrutan(t *testing.T) {
	g, l, j := gudangPutusan(t)
	ctx := context.Background()
	h, err := l.Putuskan(ctx, pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "1", Komentar: " UJI-OK "})
	if err != nil {
		t.Fatal(err)
	}
	if h.Status != models.StatusKasusSelesai || h.NoEndorsement != "UJI-PL-1/01" || h.Peserta != 3 || h.RekapWarisan != 1 {
		t.Fatalf("hasil %+v", h)
	}
	// AC 42: urutan pemanggilan adalah bagian kebenaran.
	if got := strings.Join(g.Panggilan, ","); got != "SisipRiwayat,AdaVersiResmi,TulisProduksiWarisan,Resmikan,TulisRekapWarisan" {
		t.Errorf("urutan %s", got)
	}
	p := g.Polis["EDMLF-1"]
	if p.NoPolis != "UJI-PL-1" || p.ProdKe != 2 || p.NoEndors != "UJI-PL-1/01" || p.PLNumberEDM != "UJI-PL-1/01" || p.Status != models.StatusKasusSelesai {
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
	if w["PL_NUMBER"] != "UJI-PL-1" || w["PL_NUMBER_EDM"] != "UJI-PL-1/01" || w["IDPEGA"] != "EDMLF-1" || w["COB"] != "UJI-COB" ||
		w["CURRENCY"] != "IDR" || w["PREMIUM"] != "10" {
		t.Errorf("rekap warisan %+v", w)
	}
	r := g.RiwayatKasus["EDMLF-1"]
	if len(r) != 1 || r[0].Status != "Accept" || r[0].PIC != pelakuUji.AkunID || r[0].Komentar != "UJI-OK" || r[0].Tanggal != "2026-10-01 17:00:00" { // 10:00 UTC = 17:00 WIB (penomor.DiJakarta)
		t.Errorf("riwayat %+v", r)
	}
	if len(j.catatan) != 1 || j.catatan[0].Dari != models.TahapInputEDMLife || j.catatan[0].Ke != models.StatusKasusSelesai {
		t.Errorf("jejak %+v", j.catatan)
	}
	if k, err := l.BacaKasus(ctx, pelakuUji, "EDMLF-1"); err != nil || len(k.Riwayat) != 1 || k.PLNumberEDM != "UJI-PL-1/01" {
		t.Errorf("baca sesudah Confirm %+v %v", k, err)
	}
	// Kasus tertutup tidak dapat diputuskan ulang (AC 29) - nomor tidak lahir dua kali.
	if _, err := l.Putuskan(ctx, pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "1"}); !errors.Is(err, services.ErrKasusTertutup) {
		t.Fatalf("putuskan ulang %v", err)
	}
	if v, _, _ := g.VersiBerjalan(ctx, nil, "UJI-PL-1", 0); v.ID != "EDMLF-1" || v.ProdKe != 2 {
		t.Errorf("versi berjalan sesudah Confirm %+v", v)
	}
}

func TestPutuskanDeclineMenutupTanpaVersi(t *testing.T) {
	for _, status := range []string{"2", "7"} {
		g, l, j := gudangPutusan(t)
		h, err := l.Putuskan(context.Background(), pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: status, Komentar: "UJI-NO"})
		if err != nil || h.Status != models.StatusKasusDitolak || h.NoEndorsement != "" {
			t.Fatalf("%s: %+v %v", status, h, err)
		}
		if got := strings.Join(g.Panggilan, ","); got != "SisipRiwayat,Tolak" {
			t.Errorf("%s: urutan %s", status, got)
		}
		if p := g.Polis["EDMLF-1"]; p.Status != models.StatusKasusDitolak || p.NoPolis != "" || p.PLNumberEDM != "" {
			t.Errorf("%s: kepala %+v", status, p)
		}
		if r := g.RiwayatKasus["EDMLF-1"]; len(r) != 1 || r[0].Status != "Decline" {
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
	g.Polis["EDMLF-0"] = &tiruan.Polis{ID: "EDMLF-0", NoPolis: "UJI-PL-1", OldPolicyNo: "UJI-PL-1", EdmType: "1", ProdKe: 2,
		PLNumberEDM: "UJI-PL-1/01", Status: models.StatusKasusSelesai}
	if _, err := l.Putuskan(ctx, pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "1"}); !errors.Is(err, services.ErrVersiBerubah) {
		t.Fatalf("versi berubah: %v", err)
	}
	// (NO_POLIS, PROD_KE) sudah dipakai baris yang TIDAK dihitung versi berjalan -
	// mis. ditulis jalur lain di luar modul ini. Pertahanan berlapis: di jalur
	// modul ini sendiri dobel sudah dicegah UX_PL_EDM_TERBUKA + FOR UPDATE.
	g, l, _ = gudangPutusan(t)
	g.Polis["UJI-LAIN"] = &tiruan.Polis{ID: "UJI-LAIN", NoPolis: "UJI-PL-1", EdmType: "1", ProdKe: 2, Status: models.StatusKasusDitolak}
	if _, err := l.Putuskan(ctx, pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "1"}); !errors.Is(err, services.ErrVersiBerubah) {
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
	g.Polis["EDMLF-1"].PLNumberEDM = "UJI-NOMOR-LAMA"
	h, err := l.Putuskan(context.Background(), pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "1"})
	if err != nil || h.NoEndorsement != "UJI-NOMOR-LAMA" {
		t.Fatalf("%+v %v", h, err)
	}
}

func TestPutuskanDitolak(t *testing.T) {
	ctx := context.Background()
	g := gudangKasusCSV(models.EdmTypePerubahanData) // belum Save
	l, _ := layananJejak(g)
	if _, err := l.Putuskan(ctx, pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "1"}); !errors.Is(err, services.ErrBelumDisimpan) {
		t.Errorf("belum Save: %v", err)
	}
	if len(g.RiwayatKasus["EDMLF-1"]) != 0 {
		t.Error("riwayat tertulis walau ditolak")
	}
	_, l2, _ := gudangPutusan(t)
	for _, m := range []services.MasukanPutusan{{Status: "3"}, {Status: ""}, {Status: "Accept"}, {Status: "1", Komentar: strings.Repeat("X", 256)}} {
		if _, err := l2.Putuskan(ctx, pelakuUji, "EDMLF-1", m); !errors.Is(err, services.ErrMasukanTidakSah) {
			t.Errorf("%+v: %v", m, err)
		}
	}
	if _, err := l2.Putuskan(ctx, pelakuUji, "NBLF-7", services.MasukanPutusan{Status: "1"}); !errors.Is(err, services.ErrKasusTidakAda) {
		t.Errorf("bukan kasus EDM: %v", err)
	}
}

// K3 keputusan work owner 01-10-2026 (OQ-EDM-008): polis NB WARISAN (`JSON_POLIS.PRODKE` kosong) →
// endorsement pertama `<polis>/01` (versi 2, E1), kedua `<polis>/02` (versi 3) - rumus `GenerateNoEDM_Life`
// 3–4 (`Local.Prodke` int, kosong = 0; `CARI4 = Prodke+1`; pad dua digit). Rantai Pega yang sudah punya
// endorsement `/01` warisan (`PRODKE` 1) berlanjut ke `/02`.
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
	kasus("EDMLF-11", 2)
	h, err := l.Putuskan(ctx, pelakuUji, "EDMLF-11", services.MasukanPutusan{Status: "1"})
	if err != nil || h.NoEndorsement != "UJI-PL-W/01" || g.Polis["EDMLF-11"].ProdKe != 2 {
		t.Fatalf("endorsement pertama NB warisan: %+v %v (versi %d)", h, err, g.Polis["EDMLF-11"].ProdKe)
	}
	kasus("EDMLF-12", 3)
	h, err = l.Putuskan(ctx, pelakuUji, "EDMLF-12", services.MasukanPutusan{Status: "1"})
	if err != nil || h.NoEndorsement != "UJI-PL-W/02" || g.Polis["EDMLF-12"].ProdKe != 3 {
		t.Fatalf("endorsement kedua: %+v %v (versi %d)", h, err, g.Polis["EDMLF-12"].ProdKe)
	}

	// Rantai Pega: NB (kosong) + endorsement warisan `/01` (`PRODKE` 1) → endorsement sistem baru `/02`, versi 2.
	g = tiruan.Baru()
	g.PolisWarisan = append(g.PolisWarisan, &tiruan.PolisWarisan{IDPega: "UJI-IDPEGA-NB", NoPolis: "UJI-PL-W", ProdKe: 0},
		&tiruan.PolisWarisan{IDPega: "UJI-IDPEGA-E1", NoPolis: "UJI-PL-W", ProdKe: 1, EdmType: "1"})
	l, _ = layananJejak(g)
	kasus("EDMLF-13", 2)
	if v, _, _ := g.VersiBerjalan(ctx, nil, "UJI-PL-W", 0); v.ID != "UJI-IDPEGA-E1" || v.UrutanPega != 1 {
		t.Fatalf("versi berjalan rantai Pega: %+v", v)
	}
	if h, err := l.Putuskan(ctx, pelakuUji, "EDMLF-13", services.MasukanPutusan{Status: "1"}); err != nil || h.NoEndorsement != "UJI-PL-W/02" {
		t.Errorf("sesudah endorsement Pega /01: %+v %v", h, err)
	}
}

// K4 keputusan work owner 01-10-2026 (OQ-EDM-010): Confirm menulis SATU baris `LIFEINPRODUCTION` di transaksi
// yang sama, SEBELUM peserta diresmikan - urutan `InsertJsonPolisLife_Act` 10 b2663 (tanpa prakondisi, selalu
// jalan) lalu 11/12; Decline tidak menulis; gagal menulis = transaksi batal.
func TestPutuskanMenulisProduksiWarisan(t *testing.T) {
	g, l, _ := gudangPutusan(t)
	ctx := context.Background()
	if _, err := l.Putuskan(ctx, pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "1"}); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(g.Panggilan, ","); got != "SisipRiwayat,AdaVersiResmi,TulisProduksiWarisan,Resmikan,TulisRekapWarisan" {
		t.Errorf("urutan %s", got)
	}
	if len(g.ProduksiWarisan) != 1 {
		t.Fatalf("baris LIFEINPRODUCTION: %d", len(g.ProduksiWarisan))
	}
	b := g.ProduksiWarisan[0]
	if b["IDPEGA"] != "EDMLF-1" || b["NOPOLIS"] != "UJI-PL-1" || b["NOENDORS"] != "UJI-PL-1/01" || b["TYPE"] != models.TypeTR ||
		b["BUSINESSNAME"] != "UJI-COB" {
		t.Errorf("baris %+v", b)
	}

	g, l, _ = gudangPutusan(t)
	if _, err := l.Putuskan(ctx, pelakuUji, "EDMLF-1", services.MasukanPutusan{Status: "2"}); err != nil || len(g.ProduksiWarisan) != 0 {
		t.Errorf("Decline tidak menulis LIFEINPRODUCTION: %v %d", err, len(g.ProduksiWarisan))
	}
}
