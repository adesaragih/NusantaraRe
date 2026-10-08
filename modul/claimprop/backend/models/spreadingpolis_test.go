package models_test

// Spreading klaim dari polis (keputusan work owner 08-10-2026: "kamu bisa ambil spreading dari polisnya? untuk table
// bawahnya itu child dari spreading atas, bisa di select dari PROPORTIONALARRG?" - dipilih: terisi otomatis saat polis
// dipilih). AddSpreading_Act / DeleteSpreading_Act tidak diekspor; Add/Delete manual tetap nonaktif.

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
	a.AnakSpreadingMap["UJI-INDUK|2024|UJI-GRUP"] = []models.AnakSpreading{
		{ReinsTypeID: "UJI-RI", ReinsTypeName: "QS (R/I)", Pct: "60"},
		{ReinsTypeID: "UJI-OR", ReinsTypeName: "QS (OR)", Pct: "40.00"},
	}
	return a
}

func halamanSpreading() *models.Halaman {
	h := models.HalamanBaru()
	h.Setel(models.CD+"TreatyYear", "2024")
	h.Setel(models.CD+"TreatyGroupID", "UJI-GRUP")
	return h
}

func TestIsiSpreadingDariPolis(t *testing.T) {
	k := konteksUji(acuanSpreading())
	h := halamanSpreading()
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-1"); err != nil {
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
	if len(bawah) != 2 || bawah[0]["TreatyType"] != "UJI-RI" || bawah[0]["TreatyName"] != "QS (R/I)" ||
		bawah[0]["SharePercentage"] != "60" || bawah[1]["TreatyType"] != "UJI-OR" || bawah[1]["SharePercentage"] != "40" ||
		bawah[0]["CurrencyID"] != "UJI-ID-IDR" || bawah[1]["Currency"] != "IDR" {
		t.Fatalf("tabel bawah dari PROPORTIONALARRG: %v", bawah)
	}

	// polis lain tanpa produksi: spreading diganti (kosong), bukan ditumpuk
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-TANPA-PRODUKSI"); err != nil {
		t.Fatal(err)
	}
	if n, m := len(h.AmbilDaftar(models.DaftarSpreading)), len(h.AmbilDaftar(models.DaftarBreakQS)); n != 0 || m != 0 {
		t.Fatalf("pilih polis lain: atas %d bawah %d, mau 0 0", n, m)
	}
}

// Anak PROPORTIONALARRG dicari menurut induk + treaty group klaim pada tahun arrangement TERBARU yang <= tahun treaty
// klaim (`[data DEV 08-10-2026]` klaim 2025 memakai arrangement 2024, klaim 2020 arrangement 2019, klaim 2018
// arrangement 2017); grup lain atau tahun sebelum arrangement pertama = tabel bawah kosong.
func TestTabelBawahHanyaAnakIndukTahunGrup(t *testing.T) {
	for _, c := range []struct {
		tahun, grup string
		bawah       int
	}{
		{"2024", "UJI-GRUP", 2},
		{"2025", "UJI-GRUP", 2}, // arrangement 2024 masih berlaku
		{"2023", "UJI-GRUP", 0}, // sebelum arrangement pertama
		{"2024", "UJI-GRUP-LAIN", 0},
	} {
		k := konteksUji(acuanSpreading())
		h := halamanSpreading()
		h.Setel(models.CD+"TreatyYear", c.tahun)
		h.Setel(models.CD+"TreatyGroupID", c.grup)
		if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-1"); err != nil {
			t.Fatal(err)
		}
		if len(h.AmbilDaftar(models.DaftarSpreading)) != 1 || len(h.AmbilDaftar(models.DaftarBreakQS)) != c.bawah {
			t.Fatalf("tahun %s grup %s: atas %v bawah %v, mau bawah %d", c.tahun, c.grup,
				h.AmbilDaftar(models.DaftarSpreading), h.AmbilDaftar(models.DaftarBreakQS), c.bawah)
		}
	}
}

// Tabel bawah per mata uang estimasi (SetTreatyNameSpreading_Act langkah 12-13), bila estimasi sudah ada.
func TestTabelBawahPerMataUangEstimasi(t *testing.T) {
	k := konteksUji(acuanSpreading())
	h := halamanSpreading()
	h.SetelDaftar(models.DaftarEstimasi, []models.Baris{
		{"CurrencyID": "UJI-ID-IDR", "Currency": "IDR"}, {"CurrencyID": "UJI-ID-USD", "Currency": "USD"},
	})
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-1"); err != nil {
		t.Fatal(err)
	}
	bawah := h.AmbilDaftar(models.DaftarBreakQS)
	if len(bawah) != 4 || bawah[2]["CurrencyID"] != "UJI-ID-USD" || bawah[2]["TreatyType"] != "UJI-RI" {
		t.Fatalf("2 anak x 2 mata uang estimasi: %v", bawah)
	}
}

// Ganti Treaty Type baris atas (SetTreatyNameSpreading) menyusun ulang tabel bawah dari PROPORTIONALARRG.
func TestGantiTreatyTypeMenyusunAnak(t *testing.T) {
	k := konteksUji(acuanSpreading())
	h := halamanSpreading()
	h.SetelDaftar(models.DaftarSpreading, []models.Baris{{"TreatyType": "UJI-INDUK", "SharePercentage": "100",
		"CurrencyID": "UJI-ID-IDR", "Currency": "IDR"}})
	if err := models.SetTreatyNameSpreading(k, h, 1); err != nil {
		t.Fatal(err)
	}
	if h.AmbilDaftar(models.DaftarSpreading)[0]["TreatyName"] != "UJI QS INDUK TRT" {
		t.Fatalf("nama treaty: %v", h.AmbilDaftar(models.DaftarSpreading))
	}
	if len(h.AmbilDaftar(models.DaftarBreakQS)) != 2 {
		t.Fatalf("anak: %v", h.AmbilDaftar(models.DaftarBreakQS))
	}
}

// Kasus baru dengan spreading terisi dari polis tidak lagi ditolak ProteksiData langkah 5 ("please Fill SpreadingList").
func TestSpreadingDariPolisLolosProteksiData(t *testing.T) {
	k := konteksUji(acuanSpreading())
	h := halamanSpreading()
	models.ProteksiData(h)
	if !adaPesan(h, models.PesanIsiSpreading) {
		t.Fatal("tanpa spreading: pesan ProteksiData 5 harus muncul")
	}
	h = halamanSpreading()
	if err := models.IsiSpreadingPolis(k, h, "UJI-POLIS-1"); err != nil {
		t.Fatal(err)
	}
	models.ProteksiData(h)
	if adaPesan(h, models.PesanIsiSpreading) {
		t.Fatal("spreading dari polis: pesan ProteksiData 5 tidak boleh muncul")
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
