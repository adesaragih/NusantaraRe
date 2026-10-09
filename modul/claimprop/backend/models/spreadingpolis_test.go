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

// Add (work owner 08-10-2026 "begitu add langsung set spreading type nya dan readonly"): Treaty Type baris baru langsung
// diisi treaty spreading polis yang BELUM ada di daftar (semua sudah ada: treaty polis pertama; tanpa spreading polis:
// kosong), mata uang estimasi pertama, tabel bawah + turunan disusun. Delete: baris hilang; tanpa baris atas, tabel
// bawah kosong.
func TestTambahHapusSpreading(t *testing.T) {
	a := acuanSpreading()
	a.JenisReas["UJI-INDUK2"] = "UJI QS INDUK DUA TRT"
	a.SpreadingPolisMap["UJI-POLIS-1"] = append(a.SpreadingPolisMap["UJI-POLIS-1"],
		models.SpreadingPolis{TreatyType: "UJI-INDUK2", SharePercentage: "50", CurrencyID: "UJI-ID-IDR", Currency: "IDR"})
	k := konteksUji(a)
	h := models.HalamanBaru()
	h.Setel(models.CD+"PolicyData.PolicyNo", "UJI-POLIS-1")
	h.SetelDaftar(models.DaftarEstimasi, []models.Baris{{"CurrencyID": "UJI-ID-IDR", "Currency": "IDR"}})
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{{"TreatyType": "UJI-INDUK", "CurrencyID": "UJI-ID-IDR", "Currency": "IDR"}})

	// 1. treaty polis berikutnya yang belum ada
	if err := models.AddSpreading(k, h, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	atas := h.AmbilDaftar(models.DaftarSpreading)
	if len(atas) != 2 || atas[1]["TreatyType"] != "UJI-INDUK2" || atas[1]["TreatyName"] != "UJI QS INDUK DUA TRT" ||
		atas[1]["SharePercentage"] != "" || atas[1]["CurrencyID"] != "UJI-ID-IDR" || atas[1]["Currency"] != "IDR" {
		t.Fatalf("Add treaty berikutnya: %v", atas)
	}
	if len(h.AmbilDaftar(models.DaftarBreakQS)) != 2 {
		t.Fatalf("tabel bawah langsung tersusun sesudah Add: %v", h.AmbilDaftar(models.DaftarBreakQS))
	}

	// 2. semua treaty polis sudah ada di mata uang ini: Add GAGAL (work owner 08-10-2026 "tidak ada spreading sama, jika
	// sama, gagal add spreadinglist, kecuali currency beda") - baris tidak bertambah, pesan tampil
	if err := models.AddSpreading(k, h, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if atas = h.AmbilDaftar(models.DaftarSpreading); len(atas) != 2 || !adaPesan(h, models.PesanSpreadingSama) {
		t.Fatalf("Add spreading sama mau gagal: %v pesan %v", atas, h.Pesan)
	}
	h.BersihkanPesan()

	// 3. mata uang lain di estimasi: treaty yang sama boleh, mata uang berbeda
	h.SetelDaftar(models.DaftarEstimasi, []models.Baris{
		{"CurrencyID": "UJI-ID-IDR", "Currency": "IDR"}, {"CurrencyID": "UJI-ID-USD", "Currency": "USD"},
	})
	if err := models.AddSpreading(k, h, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if atas = h.AmbilDaftar(models.DaftarSpreading); len(atas) != 3 || atas[2]["TreatyType"] != "UJI-INDUK" ||
		atas[2]["CurrencyID"] != "UJI-ID-USD" || adaPesan(h, models.PesanSpreadingSama) {
		t.Fatalf("Add mata uang lain: %v", atas)
	}

	// Delete sampai habis: tabel bawah ikut kosong; baris tak ada = galat
	for range 3 {
		if err := models.DeleteSpreading(k, h, 1, masterSpreading()); err != nil {
			t.Fatal(err)
		}
	}
	if n, m := len(h.AmbilDaftar(models.DaftarSpreading)), len(h.AmbilDaftar(models.DaftarBreakQS)); n != 0 || m != 0 {
		t.Fatalf("Delete semua: atas %d bawah %d, mau 0 0", n, m)
	}
	if err := models.DeleteSpreading(k, h, 1, masterSpreading()); err == nil {
		t.Fatal("Delete baris yang tidak ada mau galat")
	}

	// 4. tanpa spreading polis: baris kosong, Treaty Type tetap dapat dipilih
	h.Setel(models.CD+"PolicyData.PolicyNo", "UJI-POLIS-TANPA-PRODUKSI")
	h.SetelDaftar(models.DaftarEstimasi, []models.Baris{{"CurrencyID": "UJI-ID-IDR", "Currency": "IDR"}})
	if err := models.AddSpreading(k, h, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if atas = h.AmbilDaftar(models.DaftarSpreading); len(atas) != 1 || atas[0]["TreatyType"] != "" {
		t.Fatalf("Add tanpa spreading polis: %v", atas)
	}
	atas[0]["TreatyType"] = "UJI-INDUK"
	if err := models.SetTreatyNameSpreading(k, h, 1, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if h.AmbilDaftar(models.DaftarSpreading)[0]["TreatyName"] != "UJI QS INDUK TRT" || len(h.AmbilDaftar(models.DaftarBreakQS)) != 2 {
		t.Fatalf("pilih Treaty Type: atas %v bawah %v", h.AmbilDaftar(models.DaftarSpreading), h.AmbilDaftar(models.DaftarBreakQS))
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

// Memilih Treaty Type yang sudah ada di mata uang yang sama pada baris kosong ditolak: Treaty Type dikosongkan lagi,
// pesan tampil.
func TestPilihTreatyTypeSamaDitolak(t *testing.T) {
	k := konteksUji(acuanSpreading())
	h := models.HalamanBaru()
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{
		{"TreatyType": "UJI-INDUK", "CurrencyID": "UJI-ID-IDR", "Currency": "IDR"},
		{"TreatyType": "UJI-INDUK", "CurrencyID": "UJI-ID-IDR", "Currency": "IDR"},
	})
	if err := models.SetTreatyNameSpreading(k, h, 2, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if b := h.AmbilDaftar(models.DaftarSpreading)[1]; b["TreatyType"] != "" || !adaPesan(h, models.PesanSpreadingSama) {
		t.Fatalf("pilih treaty sama mau ditolak: %v pesan %v", b, h.Pesan)
	}
}

// Add sekali klik = satu baris per mata uang yang belum ada untuk treaty polis pertama yang masih kurang (work owner
// 08-10-2026 "jika add langsung kedetek 2 currency, langsung add 2 mengikuti currency"); urutan ikut mata uang estimasi.
func TestTambahSpreadingSatuBarisPerMataUang(t *testing.T) {
	estimasi := []models.Baris{{"CurrencyID": "UJI-ID-IDR", "Currency": "IDR"}, {"CurrencyID": "UJI-ID-USD", "Currency": "USD"}}
	k := konteksUji(acuanSpreading())
	h := models.HalamanBaru()
	h.Setel(models.CD+"PolicyData.PolicyNo", "UJI-POLIS-1")
	h.SetelDaftar(models.DaftarEstimasi, estimasi)

	// 2 mata uang, belum ada baris: 2 baris sekali klik
	if err := models.AddSpreading(k, h, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	atas := h.AmbilDaftar(models.DaftarSpreading)
	if len(atas) != 2 || atas[0]["TreatyType"] != "UJI-INDUK" || atas[0]["CurrencyID"] != "UJI-ID-IDR" ||
		atas[1]["TreatyType"] != "UJI-INDUK" || atas[1]["CurrencyID"] != "UJI-ID-USD" || atas[1]["TreatyName"] != "UJI QS INDUK TRT" {
		t.Fatalf("Add 2 mata uang: %v", atas)
	}
	if len(h.AmbilDaftar(models.DaftarBreakQS)) != 4 {
		t.Fatalf("tabel bawah 2 anak x 2 mata uang: %v", h.AmbilDaftar(models.DaftarBreakQS))
	}

	// semua sudah ada: gagal + pesan
	if err := models.AddSpreading(k, h, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if len(h.AmbilDaftar(models.DaftarSpreading)) != 2 || !adaPesan(h, models.PesanSpreadingSama) {
		t.Fatalf("Add semua sudah ada mau gagal: %v", h.AmbilDaftar(models.DaftarSpreading))
	}

	// 1 dari 2 sudah ada: hanya 1 baris
	h = models.HalamanBaru()
	h.Setel(models.CD+"PolicyData.PolicyNo", "UJI-POLIS-1")
	h.SetelDaftar(models.DaftarEstimasi, estimasi)
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{{"TreatyType": "UJI-INDUK", "CurrencyID": "UJI-ID-IDR", "Currency": "IDR"}})
	if err := models.AddSpreading(k, h, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if atas = h.AmbilDaftar(models.DaftarSpreading); len(atas) != 2 || atas[1]["CurrencyID"] != "UJI-ID-USD" {
		t.Fatalf("Add 1 dari 2 sudah ada: %v", atas)
	}

	// tanpa spreading polis: satu baris kosong per mata uang
	h = models.HalamanBaru()
	h.Setel(models.CD+"PolicyData.PolicyNo", "UJI-POLIS-TANPA-PRODUKSI")
	h.SetelDaftar(models.DaftarEstimasi, estimasi)
	if err := models.AddSpreading(k, h, masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if atas = h.AmbilDaftar(models.DaftarSpreading); len(atas) != 2 || atas[0]["TreatyType"] != "" || atas[1]["TreatyType"] != "" ||
		atas[0]["CurrencyID"] != "UJI-ID-IDR" || atas[1]["CurrencyID"] != "UJI-ID-USD" {
		t.Fatalf("Add tanpa spreading polis: %v", atas)
	}
}

// Keputusan work owner 09-10-2026 "kalo tidak ada di master, ambil dari proportionalarrg": SpreadingList master kosong
// -> tabel bawah = anak PROPORTIONALARRG setiap TreatyType tabel atas (tahun + treaty group klaim); master berisi tetap
// dipakai (PROPORTIONALARRG tidak dibaca).
func TestTabelBawahCadanganProportionalArrg(t *testing.T) {
	a := acuanSpreading()
	a.Anak["UJI-INDUK|2026|UJI-TG"] = []models.AnakSpreading{
		{ReinsTypeID: "UJI-PA1", ReinsTypeName: "UJI ANAK 1", Pct: "70.00"},
		{ReinsTypeID: "UJI-PA2", ReinsTypeName: "UJI ANAK 2", Pct: "30"},
	}
	k := konteksUji(a)
	h := models.HalamanBaru()
	h.Setel(models.CD+"TreatyYear", "2026")
	h.Setel(models.CD+"TreatyGroupID", "UJI-TG")
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-1", models.MasterTreaty{}); err != nil {
		t.Fatal(err)
	}
	bawah := h.AmbilDaftar(models.DaftarBreakQS)
	if len(bawah) != 2 || bawah[0]["TreatyType"] != "UJI-PA1" || bawah[0]["TreatyName"] != "UJI ANAK 1" ||
		bawah[0]["SharePercentage"] != "70" || bawah[1]["TreatyType"] != "UJI-PA2" || bawah[1]["CurrencyID"] != "UJI-ID-IDR" {
		t.Fatalf("master kosong: tabel bawah dari PROPORTIONALARRG: %v", bawah)
	}
	// master berisi: master yang dipakai, cadangan tidak
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-1", masterSpreading()); err != nil {
		t.Fatal(err)
	}
	if b := h.AmbilDaftar(models.DaftarBreakQS); len(b) != 2 || b[0]["TreatyType"] != "UJI-RI" {
		t.Fatalf("master berisi: tabel bawah dari master: %v", b)
	}
	// master kosong dan PROPORTIONALARRG tanpa anak: tabel bawah kosong
	h.Setel(models.CD+"TreatyGroupID", "UJI-TG-LAIN")
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-1", models.MasterTreaty{}); err != nil {
		t.Fatal(err)
	}
	if b := h.AmbilDaftar(models.DaftarBreakQS); len(b) != 0 {
		t.Fatalf("tanpa master dan tanpa anak: %v", b)
	}
}

// U/Y polis kosong (klaim lama hasil pemuat): tahun cadangan = YearofAccount (tahun master treaty).
func TestTabelBawahCadanganTanpaTreatyYear(t *testing.T) {
	a := acuanSpreading()
	a.Anak["UJI-INDUK|2018|UJI-TG"] = []models.AnakSpreading{{ReinsTypeID: "UJI-PA1", ReinsTypeName: "UJI ANAK 1", Pct: "60"}}
	k := konteksUji(a)
	h := models.HalamanBaru()
	h.Setel(models.CD+"YearofAccount", "2018")
	h.Setel(models.CD+"TreatyGroupID", "UJI-TG")
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-1", models.MasterTreaty{}); err != nil {
		t.Fatal(err)
	}
	if b := h.AmbilDaftar(models.DaftarBreakQS); len(b) != 1 || b[0]["TreatyType"] != "UJI-PA1" {
		t.Fatalf("tanpa TreatyYear: anak lewat YearofAccount: %v", b)
	}
}
