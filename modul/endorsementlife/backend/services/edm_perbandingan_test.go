package services_test

// Perbandingan summary sebelum/sesudah endorsement (permintaan work owner
// 08-10-2026) - terhadap gudang tiruan yang meniru SQL-nya.

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cockroachdb/apd/v3"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/services"
	"nusantarare/modul/endorsementlife/backend/tiruan"
)

// gudangBanding - kasus QR: 12 peserta Old IDR (GP 100, DED 10), 1 Old USD
// (GP 7.5), 3 New IDR (GP 50) hasil CSV.
func gudangBanding() *tiruan.Gudang {
	g := tiruan.Baru()
	g.Polis["EDML-1"] = &tiruan.Polis{ID: "EDML-1", OldPolicyNo: "UJI-PL-1", EdmType: models.EdmTypePerubahanData, ProdKe: 1,
		Kepala: map[string]string{"TYPE": models.TypeQR}}
	tambah := func(id, uang, gp, ded, status string) {
		g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: id, PolisID: "EDML-1", PLNumber: "UJI-PL-1", EdmStatus: status,
			Nilai: map[string]string{"CURRENCY": uang, "GROSS_PREMIUM": gp, "DEDUCTION": ded}})
	}
	for i := 1; i <= 12; i++ {
		tambah(fmt.Sprintf("UJI-O%02d", i), "IDR", "100", "10", models.StatusOld)
	}
	tambah("UJI-U1", "USD", "7.5", "", models.StatusOld)
	for i := 1; i <= 3; i++ {
		tambah(fmt.Sprintf("UJI-N%d", i), "IDR", "50", "", models.StatusNew)
	}
	return g
}

func cekSel(t *testing.T, nama string, s models.SelPerbandingan, cacah int, premi, balance string) {
	t.Helper()
	if s.Cacah != cacah || s.Nilai["PREMIUM"] != premi || s.Nilai["BALANCE"] != balance {
		t.Errorf("%s = cacah %d PREMIUM %s BALANCE %s; mau %d %s %s", nama, s.Cacah, s.Nilai["PREMIUM"], s.Nilai["BALANCE"], cacah, premi, balance)
	}
}

func TestPerbandinganSebelumDanSesudahProses(t *testing.T) {
	g := gudangBanding()
	l, _ := layananJejak(g)
	ctx := context.Background()

	// Sebelum Process Endorsement: belum ada Delete; New sudah ada.
	b, err := l.Perbandingan(ctx, pelakuUji, "EDML-1")
	if err != nil || len(b) != 2 || b[0].Currency != "IDR" || b[1].Currency != "USD" {
		t.Fatalf("%+v %v", b, err)
	}
	idr := b[0]
	// BALANCE QR = GP - DED - (...): Old 12×90 = 1080, New 3×50 = 150.
	cekSel(t, "before", idr.Before, 12, "1200", "1080")
	cekSel(t, "deleted", idr.Deleted, 0, "0", "0")
	cekSel(t, "added", idr.Added, 3, "150", "150")
	cekSel(t, "after", idr.After, 15, "1350", "1230")
	cekSel(t, "difference", idr.Difference, 3, "150", "150")
	// Mata uang terpisah: USD hanya Old.
	cekSel(t, "USD before", b[1].Before, 1, "7.5", "7.5")
	cekSel(t, "USD after", b[1].After, 1, "7.5", "7.5")
	cekSel(t, "USD difference", b[1].Difference, 0, "0", "0")

	// Process Endorsement: 2 peserta Old IDR ditandai Delete (nilainya dibalik).
	h, err := l.Simpan(ctx, pelakuUji, "EDML-1", models.PilihanHapus{Pilih: []string{"UJI-O01", "UJI-O02"}})
	if err != nil || h.Ditandai != 2 {
		t.Fatalf("simpan %+v %v", h, err)
	}
	b, err = l.Perbandingan(ctx, pelakuUji, "EDML-1")
	if err != nil || len(b) != 2 {
		t.Fatalf("%+v %v", b, err)
	}
	idr = b[0]
	cekSel(t, "before", idr.Before, 12, "1200", "1080")
	cekSel(t, "deleted", idr.Deleted, -2, "-200", "-180")
	cekSel(t, "added", idr.Added, 3, "150", "150")
	cekSel(t, "after", idr.After, 13, "1150", "1050")
	cekSel(t, "difference", idr.Difference, 1, "-50", "-30")

	// ⚠️ Rekap T_PREMIUM_LIST_SUMMARY menjumlah SELURUH baris kasus (Old sisa + Delete
	// negatif + New), jadi rekap = After + Deleted, BUKAN After (After = Old sisa + New,
	// rumus spec). Relasi itu dikunci di sini untuk SETIAP kolom rekap dan mata uang;
	// selisihnya dilaporkan ke work owner (08-10-2026).
	cekRekap(t, b, h.Rekap)
	if h.Rekap[0]["PREMIUM"] != "950" || idr.After.Nilai["PREMIUM"] != "1150" {
		t.Errorf("rekap %s, after %s", h.Rekap[0]["PREMIUM"], idr.After.Nilai["PREMIUM"])
	}
}

// cekRekap - rekap kasus = After + Deleted per kolom rekap dan mata uang.
func cekRekap(t *testing.T, b []models.BarisPerbandingan, rekap []map[string]string) {
	t.Helper()
	if len(b) != len(rekap) {
		t.Fatalf("%d mata uang perbandingan, %d rekap", len(b), len(rekap))
	}
	for i, r := range rekap {
		if b[i].Currency != r["CURRENCY"] {
			t.Fatalf("urutan mata uang %q != %q", b[i].Currency, r["CURRENCY"])
		}
		for _, k := range models.KolomRekapKasus[1:] {
			if got := tambahDesimal(t, b[i].After.Nilai[k], b[i].Deleted.Nilai[k]); got != r[k] {
				t.Errorf("%s %s: after+deleted %q != rekap %q", r["CURRENCY"], k, got, r[k])
			}
		}
	}
}

func tambahDesimal(t *testing.T, a, b string) string {
	t.Helper()
	x, err := utils.ParseDecimal(a)
	if err != nil {
		t.Fatal(err)
	}
	y, err := utils.ParseDecimal(b)
	if err != nil {
		t.Fatal(err)
	}
	var z, r apd.Decimal
	if _, err := utils.DecimalContext().Add(&z, x, y); err != nil {
		t.Fatal(err)
	}
	r.Reduce(&z)
	return utils.FormatDecimal(&r)
}

func TestPerbandinganBatalDanGerbang(t *testing.T) {
	g := gudangBanding()
	g.Polis["EDML-1"].EdmType = models.EdmTypeBatal
	g.Peserta = g.Peserta[:13] // tanpa New
	l, _ := layananJejak(g)
	ctx := context.Background()
	if _, err := l.Simpan(ctx, pelakuUji, "EDML-1", models.PilihanHapus{}); err != nil {
		t.Fatal(err)
	}
	b, err := l.Perbandingan(ctx, pelakuUji, "EDML-1")
	if err != nil {
		t.Fatal(err)
	}
	cekSel(t, "Batal before", b[0].Before, 12, "1200", "1080")
	cekSel(t, "Batal deleted", b[0].Deleted, -12, "-1200", "-1080")
	cekSel(t, "Batal after", b[0].After, 0, "0", "0")
	cekSel(t, "Batal difference", b[0].Difference, -12, "-1200", "-1080")
	r, err := g.RekapKasus(ctx, nil, "EDML-1")
	if err != nil {
		t.Fatal(err)
	}
	cekRekap(t, b, r)

	if _, err := l.Perbandingan(ctx, inti.Pelaku{}, "EDML-1"); err == nil {
		t.Error("tanpa identitas diterima")
	}
	for _, id := range []string{"NBLF-1", "EDML-9"} {
		if _, err := l.Perbandingan(ctx, pelakuUji, id); !errors.Is(err, services.ErrKasusTidakAda) {
			t.Errorf("%s: %v", id, err)
		}
	}
}

// TestRekapWarisanSelisih - PREMIUM_SELISIH / BALANCE_SELISIH M_LIFE_PREMIUM_SUMMARY (09-10-2026) = baris
// Difference Summary Comparison per mata uang; PREMIUM/BALANCE rekap warisan tetap seperti sebelumnya.
func TestRekapWarisanSelisih(t *testing.T) {
	ctx := context.Background()
	for _, k := range []struct {
		nama    string
		edmType string
		pilih   []string
	}{
		{"Perubahan Data", models.EdmTypePerubahanData, []string{"UJI-O01", "UJI-O02"}},
		{"Batal", models.EdmTypeBatal, nil},
	} {
		g := gudangBanding()
		g.Polis["EDML-1"].EdmType = k.edmType
		if k.edmType == models.EdmTypeBatal {
			g.Peserta = g.Peserta[:13] // Batal tanpa New
		}
		g.Polis["NBLF-1"] = &tiruan.Polis{ID: "NBLF-1", ProdKe: 0, Kepala: map[string]string{"TYPE": models.TypeQR}}
		g.Peserta = append(g.Peserta, &tiruan.Peserta{ID: "UJI-NB-1", PolisID: "NBLF-1", PLNumber: "UJI-PL-1", Nilai: map[string]string{}})
		l, _ := layananJejak(g)
		if _, err := l.Simpan(ctx, pelakuUji, "EDML-1", models.PilihanHapus{Pilih: k.pilih}); err != nil {
			t.Fatalf("%s simpan: %v", k.nama, err)
		}
		banding, err := l.Perbandingan(ctx, pelakuUji, "EDML-1")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := l.Putuskan(ctx, pelakuUji, "EDML-1", services.MasukanPutusan{Status: "1"}); err != nil {
			t.Fatalf("%s putuskan: %v", k.nama, err)
		}
		if len(g.RekapWarisan) != len(banding) {
			t.Fatalf("%s: %d baris rekap warisan, %d mata uang", k.nama, len(g.RekapWarisan), len(banding))
		}
		for _, b := range banding {
			var w map[string]string
			for _, r := range g.RekapWarisan {
				if r["CURRENCY"] == b.Currency {
					w = r
				}
			}
			if w == nil || w["IDPEGA"] != "EDML-1" {
				t.Fatalf("%s %s: baris rekap warisan %v", k.nama, b.Currency, w)
			}
			if w["PREMIUM_SELISIH"] != b.Difference.Nilai["PREMIUM"] || w["BALANCE_SELISIH"] != b.Difference.Nilai["BALANCE"] {
				t.Errorf("%s %s: selisih %s/%s, Difference %s/%s", k.nama, b.Currency, w["PREMIUM_SELISIH"], w["BALANCE_SELISIH"],
					b.Difference.Nilai["PREMIUM"], b.Difference.Nilai["BALANCE"])
			}
		}
		idr := banding[0]
		switch k.edmType {
		case models.EdmTypePerubahanData:
			if idr.Difference.Nilai["PREMIUM"] != "-50" || idr.Difference.Nilai["BALANCE"] != "-30" {
				t.Errorf("Difference IDR %+v", idr.Difference)
			}
		case models.EdmTypeBatal:
			// Batal: selisih = −Before.
			if idr.Difference.Nilai["PREMIUM"] != "-1200" || idr.Before.Nilai["PREMIUM"] != "1200" {
				t.Errorf("Batal IDR %+v", idr)
			}
		}
	}
}
