package services_test

// Uji `LimitCalculation` — contoh disusun dari langkah Activity-nya sendiri.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

func nilai(mu, v string) services.NilaiMataUang {
	return services.NilaiMataUang{Currency: mu, CurrencyID: mu, Value: v}
}

// [3]–[4] QS 40%: Retention% = 60, Cession% = 40, lalu tiap 100% Limit
// dibagi menurut keduanya — per mata uang, URUT.
func TestQuotaShareMembagiTiapLimit(t *testing.T) {
	d := services.HitungLimit(services.MasukanLimit{
		Jenis: "qs", Otomatis: true,
		Detail: services.DetailLimit{
			QSPct:        "40",
			IOOLimitList: []services.NilaiMataUang{nilai("IDR", "1000000000"), nilai("USD", "250000")},
			// ⛔ Isi lama HARUS hilang — langkah 2 membuangnya.
			RetentionList: []services.NilaiMataUang{nilai("IDR", "999")},
		},
	})
	if d.RetentionPct != "60" || d.CessionPct != "40" {
		t.Fatalf("persen: retensi %q cession %q, mau 60/40", d.RetentionPct, d.CessionPct)
	}
	mau := [][3]string{{"IDR", "600000000", "400000000"}, {"USD", "150000", "100000"}}
	if len(d.RetentionList) != 2 || len(d.CessionList) != 2 {
		t.Fatalf("baris: retensi %d cession %d, mau 2/2", len(d.RetentionList), len(d.CessionList))
	}
	for i, m := range mau {
		if d.RetentionList[i].Currency != m[0] || d.RetentionList[i].Value != m[1] {
			t.Errorf("retensi[%d] = %+v, mau %s %s", i, d.RetentionList[i], m[0], m[1])
		}
		if d.CessionList[i].Value != m[2] {
			t.Errorf("cession[%d] = %s, mau %s", i, d.CessionList[i].Value, m[2])
		}
	}
}

// ⛔ `@divide(x,100,4)` — faktornya dibulatkan ke EMPAT tempat SEBELUM
// dikalikan. QS 33,33333 → Retention% 66,66667 → faktor 0,6667, bukan
// 0,6666667. Pada 1.000.000 bedanya 0,33 — terlihat di layar.
func TestFaktorDibulatkanEmpatTempatSebelumDikali(t *testing.T) {
	d := services.HitungLimit(services.MasukanLimit{
		Jenis: "qs", Otomatis: true,
		Detail: services.DetailLimit{QSPct: "33.33333", IOOLimitList: []services.NilaiMataUang{nilai("IDR", "1000000")}},
	})
	if d.RetentionList[0].Value != "666700" {
		t.Errorf("retensi %s, mau 666700 (faktor 0,6667)", d.RetentionList[0].Value)
	}
	if d.CessionList[0].Value != "333300" {
		t.Errorf("cession %s, mau 333300 (faktor 0,3333)", d.CessionList[0].Value)
	}
}

// [1] autocalculate == false → keluar, NOL perubahan.
func TestTanpaOtomatisTidakMengubahApaPun(t *testing.T) {
	asal := services.DetailLimit{QSPct: "40", RetentionPct: "1", IOOLimitList: []services.NilaiMataUang{nilai("IDR", "5")}}
	d := services.HitungLimit(services.MasukanLimit{Jenis: "qs", Otomatis: false, Detail: asal})
	if d.RetentionPct != "1" || len(d.RetentionList) != 0 {
		t.Errorf("berubah padahal autocalculate false: %+v", d)
	}
}

// [6]–[14] Surplus otomatis: 100% Limit dan Retention diambil dari Quota
// Share BERGRUP SAMA di pohon; Lines mengalikan.
func TestSurplusMengambilLimitQuotaShareBergrupSama(t *testing.T) {
	qs := services.DetailLimit{
		TreatyType: "QUOTA SHARE", TreatyGroup: "FIRE",
		IOOLimitList: []services.NilaiMataUang{nilai("IDR", "100"), nilai("USD", "0")},
		COBList:      []services.KelasBisnisLimit{{ClassOfBusiness: "PROPERTY", TreatyGroup: "FIRE"}},
	}
	lain := services.DetailLimit{TreatyType: "QUOTA SHARE", TreatyGroup: "MARINE",
		IOOLimitList: []services.NilaiMataUang{nilai("IDR", "7777")}}
	d := services.HitungLimit(services.MasukanLimit{
		Jenis: "surplus", Otomatis: true,
		Detail: services.DetailLimit{TreatyType: "SURPLUS", TreatyGroup: "FIRE", Surplus: "2.5"},
		Pohon:  []services.DetailLimit{lain, qs},
	})
	if d.IOOPct != "250" || d.CessionPct != "250" {
		t.Errorf("persen IOO %q cession %q, mau 250/250", d.IOOPct, d.CessionPct)
	}
	// ⛔ Baris bernilai 0 DILEWATI — langkah 12 bersyarat `.Value>0`.
	if len(d.IOOLimitList) != 1 {
		t.Fatalf("baris IOO %d, mau 1 (baris USD 0 dilewati)", len(d.IOOLimitList))
	}
	if d.IOOLimitList[0].Value != "250" || d.RetentionList[0].Value != "100" || d.CessionList[0].Value != "250" {
		t.Errorf("IOO %s retensi %s cession %s, mau 250/100/250",
			d.IOOLimitList[0].Value, d.RetentionList[0].Value, d.CessionList[0].Value)
	}
	// ⛔ Grup lain TIDAK ikut.
	for _, x := range d.RetentionList {
		if x.Value == "7777" {
			t.Error("limit grup MARINE bocor ke surplus FIRE")
		}
	}
	if len(d.COBList) != 1 || d.COBList[0].ClassOfBusiness != "PROPERTY" {
		t.Errorf("COBList %+v, mau disalin dari Quota Share", d.COBList)
	}
}

// [11] Surplus manual: Retention DIKETIK, dan IOO/Cession dihitung darinya.
func TestSurplusManualMemakaiRetensiYangDiketik(t *testing.T) {
	d := services.HitungLimit(services.MasukanLimit{
		Jenis: "surplus", Tambah: "man", Otomatis: true,
		Detail: services.DetailLimit{Surplus: "3", RetentionList: []services.NilaiMataUang{nilai("IDR", "40")}},
	})
	if len(d.RetentionList) != 1 || d.RetentionList[0].Value != "40" {
		t.Errorf("retensi %+v, mau TETAP [40] — mode manual tidak membuangnya", d.RetentionList)
	}
	if d.IOOLimitList[0].Value != "120" || d.CessionList[0].Value != "120" {
		t.Errorf("IOO %s cession %s, mau 120/120", d.IOOLimitList[0].Value, d.CessionList[0].Value)
	}
}

// ⛔ Masukan TIDAK diubah — fungsi ini dipanggil atas salinan pohon layar.
func TestMasukanTidakDiubah(t *testing.T) {
	asal := services.DetailLimit{QSPct: "40", IOOLimitList: []services.NilaiMataUang{nilai("IDR", "10")},
		RetentionList: []services.NilaiMataUang{nilai("IDR", "1")}}
	_ = services.HitungLimit(services.MasukanLimit{Jenis: "qs", Otomatis: true, Detail: asal})
	if len(asal.RetentionList) != 1 || asal.RetentionList[0].Value != "1" {
		t.Errorf("masukan berubah: %+v", asal.RetentionList)
	}
}
