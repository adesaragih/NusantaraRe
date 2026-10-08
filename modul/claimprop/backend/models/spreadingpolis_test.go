package models_test

// Spreading klaim (keputusan work owner 08-10-2026): tabel atas terisi otomatis dari polis saat polis dipilih; tabel
// bawah = SpreadingList master treaty (SetTreatyNameSpreading_Act langkah 11, sumber XML - PROPORTIONALARRG dibatalkan
// sesudah uji data); Add / Delete aktif. AddSpreading_Act / DeleteSpreading_Act tidak diekspor.

import (
	"testing"

	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/tiruan"
)

func acuanSpreading() *tiruan.Acuan {
	a := tiruan.AcuanBaru()
	a.JenisReas["UJI-INDUK"] = "UJI QS INDUK TRT"
	a.SpreadingPolisMap["UJI-POLIS-1"] = []models.SpreadingPolis{
		{TreatyType: "UJI-INDUK", SharePercentage: "100", CurrencyID: "UJI-ID-IDR", Currency: "IDR"},
	}
	return a
}

func masterSpreading() models.MasterTreaty {
	return models.MasterTreaty{Limits: []models.LimitMaster{{Detail: []models.DetailLimit{{SpreadingList: []models.SpreadingMaster{
		{ReinsTypeID: "UJI-RI", ReinsTypeName: "UJI QS RI", Pct: "60"},
		{ReinsTypeID: "UJI-OR", ReinsTypeName: "UJI QS OR", Pct: "40.00"},
	}}}}}}
}

func TestIsiSpreadingDariPolis(t *testing.T) {
	k := konteksUji(acuanSpreading())
	h := models.HalamanBaru()
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-1", masterSpreading()); err != nil {
		t.Fatal(err)
	}
	atas := h.AmbilDaftar(models.DaftarSpreading)
	if len(atas) != 1 {
		t.Fatalf("spreading atas %d baris: %v", len(atas), atas)
	}
	for kunci, mau := range map[string]string{"TreatyType": "UJI-INDUK", "TreatyName": "UJI QS INDUK TRT",
		"SharePercentage": "100", "CurrencyID": "UJI-ID-IDR", "Currency": "IDR"} {
		if atas[0][kunci] != mau {
			t.Fatalf("atas %s = %q, mau %q", kunci, atas[0][kunci], mau)
		}
	}
	bawah := h.AmbilDaftar(models.DaftarBreakQS)
	if len(bawah) != 2 || bawah[0]["TreatyType"] != "UJI-RI" || bawah[0]["TreatyName"] != "UJI QS RI" ||
		bawah[0]["SharePercentage"] != "60" || bawah[1]["TreatyType"] != "UJI-OR" || bawah[1]["SharePercentage"] != "40" ||
		bawah[0]["CurrencyID"] != "UJI-ID-IDR" || bawah[1]["Currency"] != "IDR" {
		t.Fatalf("tabel bawah dari SpreadingList master: %v", bawah)
	}

	// polis lain tanpa produksi: spreading diganti (kosong), bukan ditumpuk; tabel bawah ikut kosong
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-TANPA-PRODUKSI", masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if n, m := len(h.AmbilDaftar(models.DaftarSpreading)), len(h.AmbilDaftar(models.DaftarBreakQS)); n != 0 || m != 0 {
		t.Fatalf("pilih polis lain: atas %d bawah %d, mau 0 0", n, m)
	}
}

// Tabel bawah per mata uang estimasi (SetTreatyNameSpreading_Act langkah 12-13), bila estimasi sudah ada.
func TestTabelBawahPerMataUangEstimasi(t *testing.T) {
	k := konteksUji(acuanSpreading())
	h := models.HalamanBaru()
	h.SetelDaftar(models.DaftarEstimasi, []models.Baris{
		{"CurrencyID": "UJI-ID-IDR", "Currency": "IDR"}, {"CurrencyID": "UJI-ID-USD", "Currency": "USD"},
	})
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-1", masterSpreading()); err != nil {
		t.Fatal(err)
	}
	bawah := h.AmbilDaftar(models.DaftarBreakQS)
	if len(bawah) != 4 || bawah[2]["CurrencyID"] != "UJI-ID-USD" || bawah[2]["TreatyType"] != "UJI-RI" {
		t.Fatalf("2 anak x 2 mata uang estimasi: %v", bawah)
	}
}

// Add: baris kosong bermata uang estimasi pertama; pilih Treaty Type (SetTreatyNameSpreading) mengisi nama dan
// menyusun ulang tabel bawah. Delete: baris hilang; tanpa baris atas, tabel bawah kosong.
func TestTambahPilihHapusSpreading(t *testing.T) {
	k := konteksUji(acuanSpreading())
	h := models.HalamanBaru()
	h.SetelDaftar(models.DaftarEstimasi, []models.Baris{{"CurrencyID": "UJI-ID-IDR", "Currency": "IDR"}})
	if err := models.AddSpreading(k, h, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	atas := h.AmbilDaftar(models.DaftarSpreading)
	if len(atas) != 1 || atas[0]["TreatyType"] != "" || atas[0]["CurrencyID"] != "UJI-ID-IDR" || atas[0]["Currency"] != "IDR" {
		t.Fatalf("Add: %v", atas)
	}
	atas[0]["TreatyType"] = "UJI-INDUK"
	if err := models.SetTreatyNameSpreading(k, h, 1, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if h.AmbilDaftar(models.DaftarSpreading)[0]["TreatyName"] != "UJI QS INDUK TRT" || len(h.AmbilDaftar(models.DaftarBreakQS)) != 2 {
		t.Fatalf("pilih Treaty Type: atas %v bawah %v", h.AmbilDaftar(models.DaftarSpreading), h.AmbilDaftar(models.DaftarBreakQS))
	}
	if err := models.DeleteSpreading(k, h, 1, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if n, m := len(h.AmbilDaftar(models.DaftarSpreading)), len(h.AmbilDaftar(models.DaftarBreakQS)); n != 0 || m != 0 {
		t.Fatalf("Delete baris terakhir: atas %d bawah %d, mau 0 0", n, m)
	}
	if err := models.DeleteSpreading(k, h, 1, masterSpreading()); err == nil {
		t.Fatal("Delete baris yang tidak ada mau galat")
	}
}

// Kasus baru dengan spreading terisi dari polis lolos ProteksiData langkah 5; semua baris dihapus = ditolak lagi.
func TestSpreadingDariPolisLolosProteksiData(t *testing.T) {
	k := konteksUji(acuanSpreading())
	h := models.HalamanBaru()
	models.ProteksiData(h)
	if !adaPesan(h, models.PesanIsiSpreading) {
		t.Fatal("tanpa spreading: pesan ProteksiData 5 harus muncul")
	}
	h = models.HalamanBaru()
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-1", masterSpreading()); err != nil {
		t.Fatal(err)
	}
	models.ProteksiData(h)
	if adaPesan(h, models.PesanIsiSpreading) {
		t.Fatal("spreading dari polis: pesan ProteksiData 5 tidak boleh muncul")
	}
	h.BersihkanPesan()
	if err := models.DeleteSpreading(k, h, 1, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	models.ProteksiData(h)
	if !adaPesan(h, models.PesanIsiSpreading) {
		t.Fatal("semua baris dihapus: pesan ProteksiData 5 harus muncul lagi")
	}
}

func adaPesan(h *models.Halaman, pesan string) bool {
	for _, daftar := range h.Pesan {
		for _, p := range daftar {
			if p == pesan {
				return true
			}
		}
	}
	return false
}
