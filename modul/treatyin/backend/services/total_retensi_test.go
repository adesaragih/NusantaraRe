package services_test

// Bukti rumus panel `Total Retention Amount`.
//
// ⛔ Rumusnya DIBACA, bukan ditebak dari nama panelnya —
// `Activity/TreatyInNPSetTotal.xml` (`pyRuleAvailable = Yes`, nol penjaga
// `1=2`), cabang `param.type=retention`. Uji ini menjaga ketiga hal yang
// Activity itu nyatakan: pengelompokan per mata uang, penjumlahan ke baris
// yang sudah ada, dan urutan kemunculan pertama.

import (
	"testing"

	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/services"
)

func retensi(pasang ...string) []models.BarisRetensiWarisan {
	var b []models.BarisRetensiWarisan
	for i := 0; i+1 < len(pasang); i += 2 {
		b = append(b, models.BarisRetensiWarisan{MataUang: pasang[i], Jumlah: pasang[i+1]})
	}
	return b
}

// ⭐ Nol baris -> irisan KOSONG, bukan nil. Panelnya yang menyatakan
// "No items"; services tidak boleh mengirim `null` ke layar.
func TestTotalRetensiNolBaris(t *testing.T) {
	got := services.TotalRetensiPerMataUang(nil)
	if got == nil {
		t.Fatal("nil, mau irisan kosong")
	}
	if len(got) != 0 {
		t.Errorf("%d baris dari nol masukan", len(got))
	}
}

// ⭐ Mata uang yang SAMA dijumlahkan ke SATU baris — `Appendflag` Activity.
func TestTotalRetensiMenjumlahkanMataUangSama(t *testing.T) {
	got := services.TotalRetensiPerMataUang(retensi(
		"IDR", "1000", "IDR", "2500", "IDR", "0.5",
	))
	if len(got) != 1 {
		t.Fatalf("%d baris, mau 1: %+v", len(got), got)
	}
	if got[0].MataUang != "IDR" || got[0].Nilai != "3500.5" {
		t.Errorf("%+v, mau {IDR 3500.5}", got[0])
	}
}

// ⭐ URUTANNYA kemunculan pertama, bukan abjad.
//
// ⛔ Mengurutkan abjad di sini akan membuat layar ini berbeda dari layar
// lama tanpa satu pun sumber yang memintanya — `<APPEND>` Activity menaruh
// mata uang baru di ujung.
func TestTotalRetensiBerurutKemunculan(t *testing.T) {
	got := services.TotalRetensiPerMataUang(retensi(
		"USD", "10", "IDR", "20", "EUR", "30", "IDR", "5",
	))
	mau := []models.BarisTotalRetensiWarisan{
		{MataUang: "USD", Nilai: "10"},
		{MataUang: "IDR", Nilai: "25"},
		{MataUang: "EUR", Nilai: "30"},
	}
	if len(got) != len(mau) {
		t.Fatalf("%d baris, mau %d: %+v", len(got), len(mau), got)
	}
	for i := range mau {
		if got[i] != mau[i] {
			t.Errorf("baris %d: %+v, mau %+v", i, got[i], mau[i])
		}
	}
}

// ⭐ Baris KOTOR dihitung nol dan TIDAK menggugurkan panelnya.
//
// ⚠️ Activity-nya punya pesan "Error Amount is empty" untuk keadaan ini, dan
// pesan itu milik jalur TULIS yang belum dibangun. Layar baca-saja tidak
// boleh menolak menampilkan apa pun karena satu baris tidak dapat diurai.
func TestTotalRetensiBarisKotorDihitungNol(t *testing.T) {
	got := services.TotalRetensiPerMataUang(retensi(
		"IDR", "100", "IDR", "", "IDR", "n/a", "IDR", "25",
	))
	if len(got) != 1 {
		t.Fatalf("%d baris, mau 1: %+v", len(got), got)
	}
	if got[0].Nilai != "125" {
		t.Errorf("nilai %q, mau 125", got[0].Nilai)
	}
}

// ⭐ Mata uang KOSONG tetap satu barisnya sendiri.
//
// ⛔ Menggabungkannya dengan baris bermata-uang akan menyembunyikan data
// kotor di balik angka yang terlihat benar.
func TestTotalRetensiMataUangKosongBerdiriSendiri(t *testing.T) {
	got := services.TotalRetensiPerMataUang(retensi("IDR", "100", "", "7", "  ", "3"))
	if len(got) != 2 {
		t.Fatalf("%d baris, mau 2: %+v", len(got), got)
	}
	// `  ` dipangkas menjadi `` — keduanya satu kelompok.
	if got[1].MataUang != "" || got[1].Nilai != "10" {
		t.Errorf("%+v, mau {\"\" 10}", got[1])
	}
}

// ⭐ UANG, bukan float biner.
//
// ⛔ `0.1 + 0.2` dengan `float64` menghasilkan `0,30000000000000004`, dan
// angka itu lalu tampil di layar sebagai nilai yang tidak pernah diketik
// siapa pun. Uji ini gagal pada hari seseorang menyederhanakannya.
func TestTotalRetensiTanpaSisaBiner(t *testing.T) {
	got := services.TotalRetensiPerMataUang(retensi("IDR", "0.1", "IDR", "0.2"))
	if got[0].Nilai != "0.3" {
		t.Errorf("nilai %q, mau 0.3 — sisa biner bocor ke layar", got[0].Nilai)
	}
}

// ⭐ Nilai BESAR tidak kehilangan digit.
//
// Retensi tersimpan sebagai teks dan dapat melampaui lebar `float64` yang
// tepat (2^53). Uji ini memakai nilai yang justru melewati batas itu.
func TestTotalRetensiNilaiBesarUtuh(t *testing.T) {
	got := services.TotalRetensiPerMataUang(retensi(
		"IDR", "9007199254740993", "IDR", "1",
	))
	if got[0].Nilai != "9007199254740994" {
		t.Errorf("nilai %q, mau 9007199254740994", got[0].Nilai)
	}
}
