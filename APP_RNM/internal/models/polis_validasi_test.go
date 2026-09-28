package models_test

// Uji gerbang `ProtectAccept` - tiket 01 bagian 2 PremiumList Life.
//
// ⛔ Yang dijaga: pesannya VERBATIM, gerbang posisinya benar, dan rule
// MENGUMPULKAN seluruh galat alih-alih berhenti pada yang pertama.

import (
	"os"
	"strings"
	"testing"

	"nusantarare/internal/models"
	"nusantarare/pkg/utils"
)

// uang menyusun Money dari teks desimal.
func uang(teks string) models.Money {
	d, err := utils.ParseDecimal(teks)
	if err != nil {
		panic(err)
	}
	return models.Money{Amount: d, Currency: "IDR"}
}

// nisbah menyusun Ratio dari teks desimal.
func nisbah(teks string) models.Ratio {
	d, err := utils.ParseDecimal(teks)
	if err != nil {
		panic(err)
	}
	return models.Ratio{Value: d}
}

// penawaranLengkap adalah penawaran yang lolos seluruh pemeriksaan.
func penawaranLengkap() models.PenawaranPolis {
	return models.PenawaranPolis{
		Posisi:           models.PosisiOffer,
		NoOffer:          "UJI-OFR-1",
		TypeCeding:       "UJI-TC",
		BusinessCode:     "UJI-COB",
		Type:             "QP",
		ProductName:      "UJI-PRODUK",
		ProRateType:      "UJI-PRORATE",
		MarketingName:    "UJI-MO",
		SourceOfBusiness: "UJI-SOB",
		Detail: []models.BarisDetailPenawaran{{
			SumInsured:   uang("1000"),
			SumReasured:  uang("900"),
			Rate:         nisbah("0.5"),
			GrossPremium: uang("100"),
			NetPremium:   uang("90"),
		}},
		MataUang: []models.BarisMataUangPenawaran{{
			Premium: uang("100"), Balance: uang("50"),
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
	kosong := models.PenawaranPolis{Posisi: models.PosisiOffer}
	pesan := models.ValidasiPenawaran(kosong)
	if len(pesan) < 8 {
		t.Fatalf("galat terkumpul = %d, mau minimal delapan: %v", len(pesan), pesan)
	}
	for _, mau := range []string{
		models.PesanNoOfferKosong, models.PesanTypeCedingKosong,
		models.PesanBusinessCodeKosong, models.PesanTypeKosong,
		models.PesanProductNameKosong, models.PesanProRateTypeKosong,
		models.PesanMarketingKosong, models.PesanBelumUnggahCSV,
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

func TestGerbangPosisiMemilihPemeriksaan(t *testing.T) {
	// ⛔ `COB can't null` HANYA di posisi Offer (b1207); `SOB can't null`
	// HANYA di posisi Premium (b2288). Menukarnya membuat layar menuntut
	// isian yang layarnya sendiri tidak punya.
	offer := penawaranLengkap()
	offer.Posisi = models.PosisiOffer
	offer.BusinessCode = ""
	offer.SourceOfBusiness = ""
	pesan := models.ValidasiPenawaran(offer)
	if !berisi(pesan, models.PesanBusinessCodeKosong) {
		t.Error("COB tidak dituntut di posisi Offer")
	}
	if berisi(pesan, models.PesanSOBKosong) {
		t.Error("SOB dituntut di posisi Offer; ia milik posisi Premium")
	}

	premium := penawaranLengkap()
	premium.Posisi = models.PosisiPremium
	premium.BusinessCode = ""
	premium.SourceOfBusiness = ""
	pesan = models.ValidasiPenawaran(premium)
	if !berisi(pesan, models.PesanSOBKosong) {
		t.Error("SOB tidak dituntut di posisi Premium")
	}
	if berisi(pesan, models.PesanBusinessCodeKosong) {
		t.Error("COB dituntut di posisi Premium; ia milik posisi Offer")
	}
}

func TestPesanBarisMemakaiNomorUrutSatuBased(t *testing.T) {
	// `local.idx` NOMOR URUT, 1-based - padanan `.pxListSubscript`. Nomor
	// urut dipertahankan supaya pesannya dapat dibandingkan dengan sistem
	// lama.
	p := penawaranLengkap()
	p.Detail = append(p.Detail, models.BarisDetailPenawaran{}) // baris KEDUA kosong
	pesan := models.ValidasiPenawaran(p)
	if !berisi(pesan, "Sum insured number 2 is 0") {
		t.Errorf("pesan baris kedua tidak bernomor 2: %v", pesan)
	}
	if berisi(pesan, "Sum insured number 1 is 0") {
		t.Errorf("baris pertama yang terisi ikut tertuduh: %v", pesan)
	}
	if berisi(pesan, "Sum insured number 0 is 0") {
		t.Error("nomor urut 0-based; rule aslinya 1-based")
	}
}

func TestSumReasuredDanRateHanyaMenolakNolBukanNegatif(t *testing.T) {
	// ⚠️ PENYIMPANGAN YANG DIPERTAHANKAN. Rule menguji `SUM_REASURED==0`
	// b2705 dan `RATE==0` b2867 - TANPA "atau kurang" - sedangkan empat
	// pemeriksaan lain memakai `==0 || <0`. "Memperbaikinya" berarti menolak
	// baris yang sistem lama terima.
	p := penawaranLengkap()
	p.Detail[0].SumReasured = uang("-5")
	p.Detail[0].Rate = nisbah("-0.5")
	pesan := models.ValidasiPenawaran(p)
	if berisi(pesan, "Sum reasured number 1 is 0") {
		t.Error("sum reasured negatif ditolak; rule aslinya hanya menguji == 0")
	}
	if berisi(pesan, "Rate number 1 is 0") {
		t.Error("rate negatif ditolak; rule aslinya hanya menguji == 0")
	}
	// Sedangkan yang NOL tetap ditolak.
	p.Detail[0].SumReasured = uang("0")
	if !berisi(models.ValidasiPenawaran(p), "Sum reasured number 1 is 0") {
		t.Error("sum reasured nol tidak ditolak")
	}
}

func TestKosongDiperlakukanSebagaiNol(t *testing.T) {
	// ⛔ PENYIMPANGAN SADAR, disempitkan ke gerbang ini. ADR-U-0027 berkata
	// kosong bukan nol, dan itu tetap berlaku di seluruh model. Tetapi Pega
	// menilai `.SUM_INSURED==0` BENAR untuk properti yang belum diisi;
	// menolak menyamakannya berarti isian yang dikosongkan pemakai lolos
	// gerbang yang di sistem lama menahannya.
	p := penawaranLengkap()
	p.Detail[0].SumInsured = models.Money{} // kosong, bukan nol
	if !berisi(models.ValidasiPenawaran(p), "Sum insured number 1 is 0") {
		t.Error("sum insured KOSONG lolos; Pega menahannya")
	}
}

func TestDiskonNolLolosHanyaNegatifDitolak(t *testing.T) {
	// ⚠️ `DISCOUNT_PREMIUM_RETRO<0` b4190 - hanya negatif. Diskon nol sah.
	p := penawaranLengkap()
	p.Posisi = models.PosisiPremium
	p.Type = "TP"
	p.MataUang[0].DiscountPremiumRetro = uang("0")
	if berisi(models.ValidasiPenawaran(p), models.PesanBalanceNol) {
		t.Error("diskon NOL menyalakan gerbang TP; rule hanya menolak negatif")
	}
	p.MataUang[0].DiscountPremiumRetro = uang("-1")
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
		models.PesanNoOfferKosong, models.PesanTypeCedingKosong,
		models.PesanBusinessCodeKosong, models.PesanTypeKosong,
		models.PesanProductNameKosong, models.PesanProRateTypeKosong,
		models.PesanMarketingKosong, models.PesanBelumUnggahCSV,
		models.PesanSOBKosong, models.PesanPremiumNol, models.PesanBalanceNol,
	} {
		if !strings.Contains(teks, `"`+pesan+`"`) {
			t.Errorf("pesan %q bukan kalimat mana pun di ProtectAccept.xml", pesan)
		}
	}
	// Dan kelima pesan berbaris - dengan `local.idx` di tengahnya.
	for _, awalan := range []string{
		"Sum insured number ", "Sum reasured number ", "Rate number ",
		"Gross premium number ", "Net premium number ",
	} {
		if !strings.Contains(teks, `"`+awalan+`" +local.idx+ " is 0"`) {
			t.Errorf("pola pesan baris %q tidak cocok dengan rule", awalan)
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
