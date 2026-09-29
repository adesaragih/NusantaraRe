package models_test

// Uji gerbang `ProtectAccept` - tiket 01 bagian 2 PremiumList Life.
//
// ⛔ Yang dijaga: pesannya VERBATIM, gerbang posisinya benar, dan rule
// MENGUMPULKAN seluruh galat alih-alih berhenti pada yang pertama.

import (
	"os"
	"strings"
	"testing"

	"nusantarare/inti/uang"
	"nusantarare/inti/utils"
	"nusantarare/modul/premiumlist/models"
)

// uang menyusun Money dari teks desimal.
func uangUji(teks string) uang.Money {
	d, err := utils.ParseDecimal(teks)
	if err != nil {
		panic(err)
	}
	return uang.Money{Amount: d, Currency: "IDR"}
}

// penawaranLengkap adalah penawaran yang lolos seluruh pemeriksaan.
func penawaranLengkap() models.PenawaranPolis {
	return models.PenawaranPolis{
		Posisi:           models.PosisiPremium,
		TypeCeding:       "UJI-TC",
		BusinessCode:     "UJI-COB",
		Type:             "QP",
		ProductName:      "UJI-PRODUK",
		ProRateType:      "UJI-PRORATE",
		MarketingName:    "UJI-MO",
		SourceOfBusiness: "UJI-SOB",
		CacahDetail:      1,
		MataUang: []models.BarisMataUangPenawaran{{
			Premium: uangUji("100"), Balance: uangUji("50"),
		}},
	}
}

func TestPenawaranLengkapLolos(t *testing.T) {
	if pesan := models.ValidasiPenawaran(penawaranLengkap()); len(pesan) != 0 {
		t.Errorf("penawaran lengkap ditolak: %v", pesan)
	}
	if !models.PenawaranBolehMaju(penawaranLengkap()) {
		t.Error("PenawaranBolehMaju menolak penawaran lengkap")
	}
}

func TestSeluruhGalatDikumpulkanBukanYangPertamaSaja(t *testing.T) {
	// ⛔ INTI berkas ini. `ProtectAccept` menumpuk seluruh pesan ke
	// `Local.Err` lalu menampilkannya SEKALI di b4462. Pemakai yang harus
	// menekan Simpan tujuh kali untuk menemukan tujuh isian kosong akan
	// berhenti memakai layarnya.
	kosong := models.PenawaranPolis{Posisi: models.PosisiPremium}
	pesan := models.ValidasiPenawaran(kosong)
	if len(pesan) < 6 {
		t.Fatalf("galat terkumpul = %d, mau minimal enam: %v", len(pesan), pesan)
	}
	for _, mau := range []string{
		models.PesanTypeKosong, models.PesanProductNameKosong,
		models.PesanProRateTypeKosong, models.PesanMarketingKosong,
		models.PesanBelumUnggahCSV, models.PesanSOBKosong,
	} {
		if !berisi(pesan, mau) {
			t.Errorf("pesan %q tidak terkumpul", mau)
		}
	}
	// Dan gabungannya dipisah BARIS BARU, seperti `Local.Err + "\n" + ...`.
	if !strings.Contains(models.GabungPesanPenawaran(pesan), "\n") {
		t.Error("pesan tidak dipisah baris baru")
	}
}

func TestGerbangPosisiMenggerbangiSeluruhLangkah(t *testing.T) {
	// ⛔ b1207 `Position=="Offer"` menggerbangi SELURUH langkah 4 (TypeCeding,
	// COB); b2288 `Position=="Premium"` SELURUH langkah 5 (Type sampai SOB).
	// Ronde pertama hanya menggerbangi COB dan SOB (ralat sensus 28-09-2026).
	medanOffer := []string{models.PesanTypeCedingKosong, models.PesanBusinessCodeKosong}
	medanPremium := []string{
		models.PesanTypeKosong, models.PesanProductNameKosong,
		models.PesanProRateTypeKosong, models.PesanMarketingKosong,
		models.PesanBelumUnggahCSV, models.PesanSOBKosong,
	}
	kosongkan := func(posisi string) []string {
		return models.ValidasiPenawaran(models.PenawaranPolis{Posisi: posisi})
	}
	offer, premium := kosongkan(models.PosisiOffer), kosongkan(models.PosisiPremium)
	for _, m := range medanOffer {
		if !berisi(offer, m) {
			t.Errorf("%q tidak dituntut di Offer", m)
		}
		if berisi(premium, m) {
			t.Errorf("%q dituntut di Premium; ia milik langkah 4 (Offer)", m)
		}
	}
	for _, m := range medanPremium {
		if !berisi(premium, m) {
			t.Errorf("%q tidak dituntut di Premium", m)
		}
		if berisi(offer, m) {
			t.Errorf("%q dituntut di Offer; ia milik langkah 5 (Premium)", m)
		}
	}
}

func TestLangkahTerRemarkTidakDitiru(t *testing.T) {
	// ⛔ 4.1 (`Please choose no offer !`, `//` b720) dan langkah 6 (per baris,
	// `//` b2340) tidak pernah jalan di sistem lama. Kalimatnya karena itu
	// tidak boleh muncul dari gerbang ini, apa pun masukannya.
	for _, posisi := range []string{models.PosisiOffer, models.PosisiPremium} {
		for _, pesan := range models.ValidasiPenawaran(models.PenawaranPolis{Posisi: posisi}) {
			if strings.Contains(pesan, "no offer") || strings.Contains(pesan, " number ") {
				t.Errorf("%s: pesan langkah ter-remark muncul: %q", posisi, pesan)
			}
		}
	}
}

func TestDiskonNolLolosHanyaNegatifDitolak(t *testing.T) {
	// ⚠️ `DISCOUNT_PREMIUM_RETRO<0` b4190 - hanya negatif. Diskon nol sah.
	p := penawaranLengkap()
	p.Posisi = models.PosisiPremium
	p.Type = "TP"
	p.MataUang[0].DiscountPremiumRetro = uangUji("0")
	if berisi(models.ValidasiPenawaran(p), models.PesanBalanceNol) {
		t.Error("diskon NOL menyalakan gerbang TP; rule hanya menolak negatif")
	}
	p.MataUang[0].DiscountPremiumRetro = uangUji("-1")
	if !berisi(models.ValidasiPenawaran(p), models.PesanBalanceNol) {
		t.Error("diskon negatif tidak menyalakan gerbang TP")
	}
}

// TestPesanValidasiVERBATIMDariKorpus membaca rule-nya LANGSUNG.
//
// ⛔ Pesannya kalimat yang pemakai sistem lama hafal - termasuk spasi
// sebelum tanda seru dan ejaan `System Reinsurance can't null`. Memperbaiki
// ejaannya berarti layar baru berbicara dengan kalimat yang tidak pernah ada.
func TestPesanValidasiVERBATIMDariKorpus(t *testing.T) {
	const letak = `D:\XML\RNM_BRD\PremiumList Life\Activity\ProtectAccept.xml`
	isi, err := os.ReadFile(letak)
	if err != nil {
		t.Skipf("korpus tidak terjangkau di mesin ini (%v); pesan tidak terperiksa", err)
	}
	teks := string(isi)
	for _, pesan := range []string{
		models.PesanTypeCedingKosong,
		models.PesanBusinessCodeKosong, models.PesanTypeKosong,
		models.PesanProductNameKosong, models.PesanProRateTypeKosong,
		models.PesanMarketingKosong, models.PesanBelumUnggahCSV,
		models.PesanSOBKosong, models.PesanPremiumNol, models.PesanBalanceNol,
	} {
		if !strings.Contains(teks, `"`+pesan+`"`) {
			t.Errorf("pesan %q bukan kalimat mana pun di ProtectAccept.xml", pesan)
		}
	}
}

func berisi(daftar []string, cari string) bool {
	for _, x := range daftar {
		if x == cari {
			return true
		}
	}
	return false
}
