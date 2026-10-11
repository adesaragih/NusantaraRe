package services_test

import (
	"testing"

	"nusantarare/modul/treatyin/backend/services"
)

func TestSetAmountConversionKaliKurs(t *testing.T) {
	b := services.EgnpiNP{Currency: "USD", Amount: "1000"}
	kurs := []services.KursNP{{Currency: "SGD", Conversion: "11000"}, {Currency: "USD", Conversion: "15000"}}
	if p := services.SetAmountConversion(&b, kurs); len(p) != 0 {
		t.Fatalf("pesan tak diduga: %v", p)
	}
	if b.AmountIDR != "15000000" {
		t.Fatalf("AmountIDR = %q, mau 15000000", b.AmountIDR)
	}
}

// ⚠️ Ekspor nol `break`: setiap kecocokan MENIMPA `local.ConvValue`, jadi
// yang menang baris TERAKHIR. Uji ini yang menahan seseorang menambahkan
// `break` "supaya lebih cepat".
func TestSetAmountConversionKecocokanTerakhirMenang(t *testing.T) {
	b := services.EgnpiNP{Currency: "USD", Amount: "2"}
	kurs := []services.KursNP{{Currency: "USD", Conversion: "10"}, {Currency: "USD", Conversion: "20"}}
	services.SetAmountConversion(&b, kurs)
	if b.AmountIDR != "40" {
		t.Fatalf("AmountIDR = %q, mau 40 (kurs terakhir)", b.AmountIDR)
	}
}

// ⛔ Kurs tak diketahui BUKAN kurs nol. Nilainya dibiarkan dan pemakai
// diberi tahu mata uang mana yang hilang.
func TestSetAmountConversionKursTakAdaTidakMenolkan(t *testing.T) {
	b := services.EgnpiNP{Currency: "JPY", Amount: "100", AmountIDR: "999"}
	pesan := services.SetAmountConversion(&b, []services.KursNP{{Currency: "USD", Conversion: "15000"}})
	if b.AmountIDR != "999" {
		t.Fatalf("AmountIDR = %q, seharusnya dibiarkan", b.AmountIDR)
	}
	if len(pesan) != 1 || pesan[0] != "Rate of Exchange JPY not found" {
		t.Fatalf("pesan = %v", pesan)
	}
}

// ⛔ SKALA 20, DIUKUR PADA KASUS YANG TIDAK HABIS DIBAGI.
//
// Tiga baris sama besar: tiap proporsi `33.333333333333333333` dan jumlahnya
// `99.999999999999999999` - BUKAN 100. Pega menghasilkan hal yang sama, dan
// angka inilah yang membuktikan skalanya benar-benar 20.
//
// ⚠️ Bentuk pertama uji ini menuntut "tepat 100" di sini dan GAGAL. Yang
// salah tuntutannya: catatan `labels.ts` mengukur 6 kontrak Pega sungguhan
// yang nilainya habis dibagi, bukan sifat yang berlaku untuk setiap masukan.
// Dicatat supaya ronde berikutnya tidak "memperbaiki" kode agar lulus.
func TestProporsiMemakaiSkalaDuaPuluh(t *testing.T) {
	rows := []services.EgnpiNP{
		{Currency: "IDR", Amount: "1", AmountIDR: "1"},
		{Currency: "IDR", Amount: "1", AmountIDR: "1"},
		{Currency: "IDR", Amount: "1", AmountIDR: "1"},
	}
	_, prop, _, pesan := services.NPSetTotalEgnpi(rows)
	if len(pesan) != 0 {
		t.Fatalf("pesan tak diduga: %v", pesan)
	}
	if rows[0].Proportion != "33.333333333333333333" {
		t.Fatalf("proporsi baris = %q, mau 20 desimal", rows[0].Proportion)
	}
	if prop != "99.999999999999999999" {
		t.Fatalf("total proporsi = %q", prop)
	}
}

// ⭐ Dan pada nilai yang HABIS dibagi, jumlahnya tepat 100 - sifat yang
// diukur pemilik proses pada 6 kontrak Pega.
func TestProporsiBerjumlahTepatSeratusBilaHabisDibagi(t *testing.T) {
	rows := []services.EgnpiNP{
		{Currency: "IDR", Amount: "1", AmountIDR: "1"},
		{Currency: "IDR", Amount: "3", AmountIDR: "3"},
	}
	_, prop, _, _ := services.NPSetTotalEgnpi(rows)
	if prop != "100" {
		t.Fatalf("total proporsi = %q, mau tepat 100", prop)
	}
}

func TestTotalDanProporsiPerBaris(t *testing.T) {
	rows := []services.EgnpiNP{
		{Currency: "USD", CurrencyID: "1", Amount: "100", AmountIDR: "750"},
		{Currency: "IDR", CurrencyID: "2", Amount: "250", AmountIDR: "250"},
	}
	total, prop, perMataUang, _ := services.NPSetTotalEgnpi(rows)
	if total != "1000" {
		t.Fatalf("TotalEgnpiAmount = %q, mau 1000", total)
	}
	if rows[0].Proportion != "75" || rows[1].Proportion != "25" {
		t.Fatalf("proporsi = %q / %q", rows[0].Proportion, rows[1].Proportion)
	}
	if prop != "100" {
		t.Fatalf("total proporsi = %q", prop)
	}
	// ⛔ Grid per mata uang menjumlah `.Amount`, BUKAN `.AmountIDR`.
	if len(perMataUang) != 2 {
		t.Fatalf("baris per mata uang = %d, mau 2", len(perMataUang))
	}
	if perMataUang[0].Currency != "USD" || perMataUang[0].Value != "100" {
		t.Fatalf("baris USD = %+v, mau Value 100 (Amount, bukan AmountIDR 750)", perMataUang[0])
	}
}

func TestPerMataUangDigabungSatuBaris(t *testing.T) {
	rows := []services.EgnpiNP{
		{Currency: "USD", Amount: "100", AmountIDR: "1"},
		{Currency: "USD", Amount: "50", AmountIDR: "1"},
	}
	_, _, perMataUang, _ := services.NPSetTotalEgnpi(rows)
	if len(perMataUang) != 1 || perMataUang[0].Value != "150" {
		t.Fatalf("per mata uang = %+v, mau satu baris 150", perMataUang)
	}
}

// ⛔ Pagar langkah 6: total nol menghentikan pembagian, dan proporsi lama
// tidak dihapus.
func TestTotalNolMenahanPembagian(t *testing.T) {
	rows := []services.EgnpiNP{{Currency: "USD", Amount: "100", AmountIDR: "0", Proportion: "42"}}
	total, prop, _, pesan := services.NPSetTotalEgnpi(rows)
	if total != "0" || prop != "0" {
		t.Fatalf("total=%q proporsi=%q", total, prop)
	}
	if len(pesan) != 1 || pesan[0] != services.PesanEgnpiTotalNol {
		t.Fatalf("pesan = %v, mau %q", pesan, services.PesanEgnpiTotalNol)
	}
	if rows[0].Proportion != "42" {
		t.Fatalf("proporsi lama terhapus: %q", rows[0].Proportion)
	}
}

func TestTambahBarisMewarisiMataUangRetensiPertama(t *testing.T) {
	ret := []services.NilaiMataUang{{Currency: "USD", CurrencyID: "7"}, {Currency: "SGD", CurrencyID: "9"}}
	rows := services.TambahBarisEgnpi(nil, ret)
	if len(rows) != 1 || rows[0].Currency != "USD" || rows[0].CurrencyID != "7" {
		t.Fatalf("baris baru = %+v, mau mewarisi retensi PERTAMA", rows)
	}
	// ⚠️ Nol retensi = nol mata uang. Tidak dikarang menjadi IDR.
	kosong := services.TambahBarisEgnpi(nil, nil)
	if kosong[0].Currency != "" {
		t.Fatalf("mata uang dikarang: %q", kosong[0].Currency)
	}
}

// ⭐ `Update EGNPI Value` menyentuh SETIAP baris, bukan baris terpilih.
func TestUpdateNilaiMenyentuhSemuaBaris(t *testing.T) {
	m := services.MasukanEgnpi{
		Aksi: services.AksiEgnpiNilai,
		EGNPI: []services.EgnpiNP{
			{Currency: "USD", Amount: "1"},
			{Currency: "USD", Amount: "2"},
		},
		Kurs: []services.KursNP{{Currency: "USD", Conversion: "100"}},
	}
	h := services.HitungEgnpi(m)
	if h.EGNPI[0].AmountIDR != "100" || h.EGNPI[1].AmountIDR != "200" {
		t.Fatalf("AmountIDR = %q / %q", h.EGNPI[0].AmountIDR, h.EGNPI[1].AmountIDR)
	}
	if h.TotalEgnpiAmount != "300" {
		t.Fatalf("total = %q, mau 300", h.TotalEgnpiAmount)
	}
}

// ⛔ Masukan tidak boleh berubah - rute `/hitung/` murni.
func TestHitungEgnpiTidakMengubahMasukan(t *testing.T) {
	asal := []services.EgnpiNP{{Currency: "USD", Amount: "1", Proportion: "x"}}
	services.HitungEgnpi(services.MasukanEgnpi{
		Aksi:  services.AksiEgnpiNilai,
		EGNPI: asal,
		Kurs:  []services.KursNP{{Currency: "USD", Conversion: "100"}},
	})
	if asal[0].AmountIDR != "" || asal[0].Proportion != "x" {
		t.Fatalf("masukan berubah: %+v", asal[0])
	}
}

func TestHapusBaris(t *testing.T) {
	h := services.HitungEgnpi(services.MasukanEgnpi{
		Aksi:   services.AksiEgnpiHapus,
		Indeks: 1,
		EGNPI: []services.EgnpiNP{
			{Currency: "A", AmountIDR: "1"},
			{Currency: "B", AmountIDR: "1"},
			{Currency: "C", AmountIDR: "1"},
		},
	})
	if len(h.EGNPI) != 2 || h.EGNPI[1].Currency != "C" {
		t.Fatalf("sesudah hapus = %+v", h.EGNPI)
	}
}

// ⭐ Larik hasil nol pernah `nil` - `null` di JSON membuat layar jatuh pada
// `.length`. Pelajaran yang sudah dibayar sekali di modul ini.
func TestHasilEgnpiNolLarikNil(t *testing.T) {
	h := services.HitungEgnpi(services.MasukanEgnpi{Aksi: "entah"})
	if h.EGNPI == nil || h.TotalEgnpiAmountNP == nil || h.Pesan == nil {
		t.Fatalf("ada larik nil: %+v", h)
	}
}

// ⛔ CACAT YANG DILAPORKAN PEMILIK PROSES 7 Oktober 2026.
//
// Menekan `Add` melahirkan baris KOSONG, total menjadi nol, dan pagar
// bagi-nol meneriakkan `Error Divide by Zero` di layar yang pemakainya
// belum sempat mengetik apa pun.
//
// ⭐ Di Pega pagar itu MILIK SATU TOMBOL: `TreatyInNPSetTotal` langkah 6
// hanya berjalan lewat `Update Total`. `TreatyInNonAddItem` dan
// `TreatyInEGNPIListValue` nol memanggilnya.
func TestPesanBagiNolHanyaDariUpdateTotal(t *testing.T) {
	sepi := []string{services.AksiEgnpiTambah, services.AksiEgnpiHapus,
		services.AksiEgnpiNilai, services.AksiEgnpiKonversi}
	for _, aksi := range sepi {
		h := services.HitungEgnpi(services.MasukanEgnpi{Aksi: aksi})
		for _, p := range h.Pesan {
			if p == services.PesanEgnpiTotalNol {
				t.Fatalf("aksi %q meneriakkan %q", aksi, services.PesanEgnpiTotalNol)
			}
		}
	}
	// ⭐ Dan `Update Total` TETAP meneriakkannya — pagar yang diam bukan
	// pagar.
	h := services.HitungEgnpi(services.MasukanEgnpi{Aksi: services.AksiEgnpiTotal})
	if len(h.Pesan) != 1 || h.Pesan[0] != services.PesanEgnpiTotalNol {
		t.Fatalf("Update Total nol meneriakkan pagar bagi-nol: %v", h.Pesan)
	}
}

// ⚠️ Totalnya TETAP dihitung ulang pada aksi lain — yang disaring hanya
// pesannya. Menghentikan hitungannya akan membuat total basi di sebelah
// baris yang baru berubah.
func TestTotalTetapDihitungWalauPesanDisaring(t *testing.T) {
	h := services.HitungEgnpi(services.MasukanEgnpi{
		Aksi:  services.AksiEgnpiNilai,
		EGNPI: []services.EgnpiNP{{Currency: "USD", Amount: "2"}},
		Kurs:  []services.KursNP{{Currency: "USD", Conversion: "100"}},
	})
	if h.TotalEgnpiAmount != "200" {
		t.Fatalf("total = %q, mau 200", h.TotalEgnpiAmount)
	}
	if len(h.Pesan) != 0 {
		t.Fatalf("pesan tak diduga: %v", h.Pesan)
	}
}
