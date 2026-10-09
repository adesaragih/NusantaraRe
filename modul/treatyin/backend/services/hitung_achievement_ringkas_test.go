package services_test

import (
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

// Pembantu: satu Limits berisi Detail bernama, dan hasil per-Detail-nya.
func susunRingkas(grup, jenis, prem, net, inc, lr string) (services.LimitAchievement, services.HasilDetailAchievement) {
	return services.LimitAchievement{
			TreatyType: jenis,
			Detail:     []services.DetailAchievement{{TreatyGroup: grup}},
		},
		services.HasilDetailAchievement{
			TotalAchPremium: prem, TotalAchNetPremium: net,
			TotalAchIncured: inc, LossRatio: lr,
		}
}

func ringkasUji(pasang ...[2]any) services.RingkasanAchievement {
	m := services.MasukanAchievement{}
	h := services.HasilAchievement{}
	for _, p := range pasang {
		m.Limits = append(m.Limits, p[0].(services.LimitAchievement))
		h.Limits = append(h.Limits, []services.HasilDetailAchievement{p[1].(services.HasilDetailAchievement)})
	}
	return services.RingkasAchievement(m, h)
}

func detailRingkas(grup, jenis, prem, net, inc, lr string) [2]any {
	l, d := susunRingkas(grup, jenis, prem, net, inc, lr)
	return [2]any{l, d}
}

// ⭐ Keenam kolom ekspor terisi, dan `Reins Type` datang dari Limits induk
// sementara `Treaty Group` dari Detail.
func TestRingkasAchievementKeenamKolom(t *testing.T) {
	r := ringkasUji(detailRingkas("PROPERTY", "SURPLUS", "1000", "800", "400", "40"))
	if len(r.Baris) != 2 {
		t.Fatalf("baris = %d, mau 2 (satu data + satu Total In IDR)", len(r.Baris))
	}
	b := r.Baris[0]
	if b.TreatyGroup != "PROPERTY" || b.TreatyType != "SURPLUS" {
		t.Fatalf("kolom 1/2 = %q/%q", b.TreatyGroup, b.TreatyType)
	}
	if b.TotalAchPremium != "1000" || b.TotalAchNetPremium != "800" ||
		b.TotalAchIncured != "400" || b.LossRatio != "40" {
		t.Fatalf("kolom 3-6 = %+v", b)
	}
}

// ⛔ Langkah [12.3]: Detail ber-`TreatyGroup` KOSONG dibuang dari grid.
func TestDetailTanpaTreatyGroupDibuang(t *testing.T) {
	r := ringkasUji(
		detailRingkas("", "SURPLUS", "999", "999", "999", "99"),
		detailRingkas("PROPERTY", "SURPLUS", "100", "80", "40", "50"),
	)
	if len(r.Baris) != 2 {
		t.Fatalf("baris = %d, mau 2 (satu data + Total)", len(r.Baris))
	}
	if r.Baris[0].TreatyGroup != "PROPERTY" {
		t.Fatalf("baris bergrup kosong tidak dibuang: %+v", r.Baris[0])
	}
	// Dan nilainya TIDAK ikut ke total.
	if r.SumTotalAchievNetPremium != "80" {
		t.Fatalf("total net = %q, mau 80 — baris buangan ikut terhitung", r.SumTotalAchievNetPremium)
	}
}

// ⛔ Langkah [12.4] menimpa `.SumLossRatio` dengan PEMBAGIAN, membuang
// penjumlahan yang baru saja ditumpuk [12.2.7]. Jadi total Loss Ratio BUKAN
// jumlah Loss Ratio per baris.
func TestTotalLossRatioDihitungUlangBukanDijumlah(t *testing.T) {
	r := ringkasUji(
		detailRingkas("A", "X", "100", "100", "50", "50"),
		detailRingkas("B", "X", "100", "100", "10", "10"),
	)
	// Jumlah naif: 50 + 10 = 60. Yang benar: 60/200 x 100 = 30.
	if r.SumLossRatio == "60" {
		t.Fatal("SumLossRatio dijumlah dari baris — [12.4] menimpanya dengan pembagian")
	}
	if r.SumLossRatio != "30" {
		t.Fatalf("SumLossRatio = %q, mau 30 (60/200 x 100)", r.SumLossRatio)
	}
}

func TestTotalNetDanIncuredDijumlah(t *testing.T) {
	r := ringkasUji(
		detailRingkas("A", "X", "10", "8", "4", "50"),
		detailRingkas("B", "X", "20", "12", "6", "50"),
	)
	if r.SumTotalAchievNetPremium != "20" {
		t.Fatalf("total net = %q, mau 20", r.SumTotalAchievNetPremium)
	}
	if r.SumTotalAchievIncured != "10" {
		t.Fatalf("total incured = %q, mau 10", r.SumTotalAchievIncured)
	}
}

// ⭐ Baris " Total In IDR" memakai kolom `Reins Type`, BUKAN `Treaty Group` —
// `.Detail(<LAST>).TreatyType` di ekspor.
func TestBarisTotalMemakaiKolomReinsType(t *testing.T) {
	r := ringkasUji(detailRingkas("A", "X", "10", "8", "4", "50"))
	akhir := r.Baris[len(r.Baris)-1]
	if akhir.TreatyType != " Total In IDR" {
		t.Fatalf("kolom Reins Type baris total = %q", akhir.TreatyType)
	}
	if akhir.TreatyGroup != "" {
		t.Fatalf("Treaty Group baris total terisi: %q", akhir.TreatyGroup)
	}
	// Premium baris total = SIGMA TotalAchPremium, walau nol sel kaki
	// menampilkannya (kaki hanya punya tiga sel).
	if akhir.TotalAchPremium != "10" {
		t.Fatalf("premium baris total = %q, mau 10", akhir.TotalAchPremium)
	}
}

// ⚠️ Nol baris data => nol baris total. Baris " Total In IDR" sendirian di
// grid kosong akan terbaca sebagai data.
func TestGridKosongNolBarisTotal(t *testing.T) {
	r := ringkasUji()
	if len(r.Baris) != 0 {
		t.Fatalf("baris = %+v, mau kosong", r.Baris)
	}
	if r.SumLossRatio != "0" {
		t.Fatalf("SumLossRatio = %q, mau 0", r.SumLossRatio)
	}
}

// ⛔ Net nol => Loss Ratio nol, bukan pembagian nol.
func TestNetNolTidakMembagiNol(t *testing.T) {
	r := ringkasUji(detailRingkas("A", "X", "10", "0", "4", "0"))
	if r.SumLossRatio != "0" {
		t.Fatalf("SumLossRatio = %q, mau 0", r.SumLossRatio)
	}
}

// ⭐ Nol larik nil — `null` di JSON membuat layar jatuh pada `.length`.
func TestRingkasanNolLarikNil(t *testing.T) {
	if ringkasUji().Baris == nil {
		t.Fatal("Baris nil")
	}
}
