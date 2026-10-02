package services_test

// Rekonsiliasi eksak satu kasus endorsement kebakaran nyata terhadap Pega -
// tiket E06 AC "Satu fixture kebakaran rekonsiliasi eksak".
//
// Dibaca sesudah: beforeimage_test.go.
//
// Sumber: fixture ter-de-identifikasi NB-15 `edm-fire-1.json` (dari
// `DDL\P-5 EDM-13445 (FIRE).txt`), DIBACA DI TEMPATNYA - keputusan A01: satu
// sumber yang dijaga penjaga de-identifikasi NB (`fixture_test.go`), bukan
// salinan yang lepas dari penjaga itu. Tidak ada impor Go lintas modul; yang
// dibaca hanya berkas data. Bila berkasnya belum ada, test dilewati dengan
// alasan tertulis.
//
// Yang direkonsiliasi eksak: yang DITULIS before-image dan tidak ditimpa
// sesudahnya - lapis B (`*Old`, ditulis ulang tiap layar dibuka dari OldData
// yang sama) dan lapis C (`IsOldData`, `IsProRate`). Porsi periode: langkah 15
// `SetValueToEDMWork` atas kasus ini DITOLAK (TestPorsiPeriodeKasusNyataEDMFire,
// 150,5 hari); `ProrateEDMEnd` tersimpan direkonsiliasi terhadap penulis
// terakhirnya, `CountPaymentEdm_Act` blok 13 (TestPorsiPembayaranKasusNyata).
// Nilai porsi tersimpan di fixture bukan keluaran langkah 15 (K-048).

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti/backend/uang"
	"nusantarare/inti/backend/utils"
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/services"
	"nusantarare/modul/endorsmentfacin/backend/services/predikat"
)

var berkasFixtureFire = filepath.Join("..", "..", "..", "nbfacin", "backend", "services", "premium", "testdata", "kasus", "edm-fire-1.json")

type halaman = map[string]any

func muatFixture(t *testing.T) halaman {
	t.Helper()
	mentah, err := os.ReadFile(berkasFixtureFire)
	if os.IsNotExist(err) {
		t.Skipf("fixture NB-15 %s belum ada", berkasFixtureFire)
	}
	if err != nil {
		t.Fatal(err)
	}
	var h halaman
	if err := json.Unmarshal(mentah, &h); err != nil {
		t.Fatal(err)
	}
	return h
}

func hal(h halaman, kunci string) halaman {
	v, _ := h[kunci].(map[string]any)
	return v
}

func daftar(h halaman, kunci string) []halaman {
	mentah, _ := h[kunci].([]any)
	var hasil []halaman
	for _, m := range mentah {
		if p, ok := m.(map[string]any); ok {
			hasil = append(hasil, p)
		}
	}
	return hasil
}

func teksDari(h halaman, kunci string) string {
	s, _ := h[kunci].(string)
	return s
}

// desimalDari - medan absen atau teks kosong = nil (kosong).
func desimalDari(t *testing.T, h halaman, kunci string) *apd.Decimal {
	t.Helper()
	s := teksDari(h, kunci)
	if s == "" {
		return nil
	}
	d, err := utils.ParseDecimal(s)
	if err != nil {
		t.Fatalf("%s %q: %v", kunci, s, err)
	}
	return d
}

// uangDari - mata uang tidak diisi: rekonsiliasi membandingkan jumlah, dan
// mata uang baris belum terverifikasi dari bentuk JSON ini.
func uangDari(t *testing.T, h halaman, kunci string) uang.Money {
	return uang.Money{Amount: desimalDari(t, h, kunci)}
}

func rasioDari(t *testing.T, h halaman, kunci string) uang.Ratio {
	return uang.Ratio{Value: desimalDari(t, h, kunci)}
}

func barisMUDari(t *testing.T, h halaman) models.BarisMataUang {
	return models.BarisMataUang{
		TSI: uangDari(t, h, "TSI"), Premium: uangDari(t, h, "Premium"), Rate: rasioDari(t, h, "Rate"),
		TSIOld: uangDari(t, h, "TSIOld"), PremiumOld: uangDari(t, h, "PremiumOld"), RateOld: rasioDari(t, h, "RateOld"),
	}
}

func daftarMUDari(t *testing.T, hs []halaman) []models.BarisMataUang {
	var h []models.BarisMataUang
	for _, b := range hs {
		h = append(h, barisMUDari(t, b))
	}
	return h
}

// offerFireDari - halaman OfferFacIn JSON → model, sebatas yang dibaca lapis
// B dan C lini kebakaran.
func offerFireDari(t *testing.T, h halaman) models.OfferFacIn {
	t.Helper()
	o := models.OfferFacIn{CurrencyList: daftarMUDari(t, daftar(h, "CurrencyList"))}
	for _, c := range daftar(h, "CedingCedantList") {
		o.CedingCedantList = append(o.CedingCedantList, models.Cedant{CurrencyList: daftarMUDari(t, daftar(c, "CurrencyList"))})
	}
	for _, l := range daftar(h, "LocationList") {
		p := hal(l, "Property")
		lok := models.Lokasi{IsOldData: teksDari(l, "IsOldData"), Property: models.PropertiLokasi{
			TotalTSIList:           daftarMUDari(t, daftar(p, "TotalTSIList")),
			TotalTSIPremiGrossList: daftarMUDari(t, daftar(p, "TotalTSIPremiGrossList")),
		}}
		for _, it := range daftar(p, "PropertyItemList") {
			item := models.ItemProperti{
				IsOldData:     teksDari(it, "IsOldData"),
				TSIObjectItem: uangDari(t, it, "TSIObjectItem"), TotalGrossPremi: uangDari(t, it, "TotalGrossPremi"),
				TotalPremiumNusantaraRe: uangDari(t, it, "TotalPremiumNusantaraRe"),
				TSIObjectItemOld:        uangDari(t, it, "TSIObjectItemOld"), TotalGrossPremiOld: uangDari(t, it, "TotalGrossPremiOld"),
				TotalPremiumNusantaraReOld: uangDari(t, it, "TotalPremiumNusantaraReOld"),
			}
			for _, c := range daftar(it, "CoverageList") {
				item.CoverageList = append(item.CoverageList, models.Coverage{IsOldData: teksDari(c, "IsOldData")})
			}
			lok.Property.PropertyItemList = append(lok.Property.PropertyItemList, item)
		}
		o.LocationList = append(o.LocationList, lok)
	}
	return o
}

func waktuPega(t *testing.T, s string) time.Time {
	t.Helper()
	w, err := time.Parse("20060102T150405.000 MST", s)
	if err != nil {
		t.Fatalf("DateTime Pega %q: %v", s, err)
	}
	return w
}

// samaJumlah - rekonsiliasi: kosong hanya sama dengan kosong; selain itu
// jumlahnya sama persis (`Cmp`, jadi "0" dan "0.0" sama).
func samaJumlah(t *testing.T, label string, kita uang.Money, pega uang.Money) bool {
	t.Helper()
	if kita.Kosong() != pega.Kosong() || (!kita.Kosong() && kita.Amount.Cmp(pega.Amount) != 0) {
		t.Errorf("%s: kita %q, Pega %q", label, kita.String(), pega.String())
		return false
	}
	return true
}

func samaRasioJumlah(t *testing.T, label string, kita uang.Ratio, pega uang.Ratio) {
	t.Helper()
	if kita.Kosong() != pega.Kosong() || (!kita.Kosong() && kita.Value.Cmp(pega.Value) != 0) {
		t.Errorf("%s: kita %s, Pega %s", label, kita, pega)
	}
}

// TestRekonsiliasiEksakFireEDM - E06: lapis A dari OldData fixture, lalu lapis
// B dan C dibandingkan nilai demi nilai dengan yang disimpan Pega.
func TestRekonsiliasiEksakFireEDM(t *testing.T) {
	fx := muatFixture(t)
	pega := offerFireDari(t, fx)
	halLama := hal(fx, "OldData")
	lama := offerFireDari(t, halLama)
	periodeLama := hal(halLama, "PolicyData")
	lama.PolicyData.StartDateTime = waktuPega(t, teksDari(periodeLama, "StartDateTime"))
	lama.PolicyData.EndDateTime = waktuPega(t, teksDari(periodeLama, "EndDateTime"))

	q := hal(fx, "QuotationData")
	kasus := models.KasusEndorsement{OfferFacIn: models.OfferFacIn{QuotationData: models.Quotation{
		StatusBusiness: teksDari(q, "StatusBusiness"),
		Type:           models.JenisPenyesuaian(teksDari(q, "Type")),
		EdmDate:        waktuPega(t, teksDari(q, "EdmDate")),
	}}}
	h, err := services.PrepareBeforeImage(services.MasukanBeforeImage{
		Kasus: kasus, PolisLama: &lama, Predikat: services.Predikat{IsFire: true, IsEDM: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	kita := h.Kasus.OfferFacIn

	if kita.IsProRate != teksDari(fx, "IsProRate") {
		t.Errorf("IsProRate: kita %q, Pega %q", kita.IsProRate, teksDari(fx, "IsProRate"))
	}
	// Pasangan baris menurut NOMOR URUT memang aturan lapis B (langkah 4
	// `SetOldData` memakai `.pxListSubscript`); panjang daftar dicek lebih dulu
	// supaya selisih panjang dilaporkan, bukan panic.
	samaPanjang := func(label string, kita, pega int) {
		t.Helper()
		if kita != pega {
			t.Fatalf("%s: kita %d baris, Pega %d", label, kita, pega)
		}
	}
	samaPanjang("LocationList", len(kita.LocationList), len(pega.LocationList))
	samaPanjang("CedingCedantList", len(kita.CedingCedantList), len(pega.CedingCedantList))
	dibanding := 0
	for i, lp := range pega.LocationList {
		lk := kita.LocationList[i]
		samaPanjang("PropertyItemList", len(lk.Property.PropertyItemList), len(lp.Property.PropertyItemList))
		samaPanjang("TotalTSIList", len(lk.Property.TotalTSIList), len(lp.Property.TotalTSIList))
		samaPanjang("TotalTSIPremiGrossList", len(lk.Property.TotalTSIPremiGrossList), len(lp.Property.TotalTSIPremiGrossList))
		if lk.IsOldData != lp.IsOldData {
			t.Errorf("lokasi %d IsOldData: kita %q, Pega %q", i+1, lk.IsOldData, lp.IsOldData)
		}
		for j, ip := range lp.Property.PropertyItemList {
			ik := lk.Property.PropertyItemList[j]
			samaPanjang("CoverageList", len(ik.CoverageList), len(ip.CoverageList))
			if ik.IsOldData != ip.IsOldData {
				t.Errorf("item %d IsOldData: kita %q, Pega %q", j+1, ik.IsOldData, ip.IsOldData)
			}
			samaJumlah(t, "TSIObjectItemOld", ik.TSIObjectItemOld, ip.TSIObjectItemOld)
			samaJumlah(t, "TotalGrossPremiOld", ik.TotalGrossPremiOld, ip.TotalGrossPremiOld)
			samaJumlah(t, "TotalPremiumNusantaraReOld", ik.TotalPremiumNusantaraReOld, ip.TotalPremiumNusantaraReOld)
			dibanding += 3
			for c, cp := range ip.CoverageList {
				if ik.CoverageList[c].IsOldData != cp.IsOldData {
					t.Errorf("item %d coverage %d IsOldData: kita %q, Pega %q", j+1, c+1, ik.CoverageList[c].IsOldData, cp.IsOldData)
				}
				dibanding++
			}
		}
		for j, bp := range lp.Property.TotalTSIList {
			samaJumlah(t, "TotalTSIList.TSIOld", lk.Property.TotalTSIList[j].TSIOld, bp.TSIOld)
			dibanding++
		}
		for j, bp := range lp.Property.TotalTSIPremiGrossList {
			bk := lk.Property.TotalTSIPremiGrossList[j]
			samaJumlah(t, "TotalTSIPremiGross.TSIOld", bk.TSIOld, bp.TSIOld)
			// K-046: polis lama asal NB tanpa PremiumOld/RateOld → 0.
			samaJumlah(t, "TotalTSIPremiGross.PremiumOld (K-046)", bk.PremiumOld, bp.PremiumOld)
			samaRasioJumlah(t, "TotalTSIPremiGross.RateOld (K-046)", bk.RateOld, bp.RateOld)
			dibanding += 3
		}
	}
	for i, cp := range pega.CedingCedantList {
		samaPanjang("cedant CurrencyList", len(kita.CedingCedantList[i].CurrencyList), len(cp.CurrencyList))
		for j, bp := range cp.CurrencyList {
			bk := kita.CedingCedantList[i].CurrencyList[j]
			samaJumlah(t, "cedant TSIOld", bk.TSIOld, bp.TSIOld)
			samaJumlah(t, "cedant PremiumOld", bk.PremiumOld, bp.PremiumOld)
			dibanding += 2
		}
	}
	// Uji instrumen atas jawaban yang SUDAH diketahui: 9 item × 3 + 45
	// coverage + 1 TotalTSIList + 1 TotalTSIPremiGross × 3 + 1 cedant × 2 = 78,
	// dihitung terpisah dengan Python atas fixture yang sama:
	//   py -c "import json;d=json.load(open('edm-fire-1.json'));…"  (laporan §5)
	// Pembaca fixture yang rusak memberi angka lain, bukan lulus diam-diam.
	if dibanding != 78 {
		t.Fatalf("%d nilai dibandingkan, mau 78; pembaca fixture yang rusak", dibanding)
	}
	t.Logf("%d nilai lapis B/C direkonsiliasi eksak terhadap Pega", dibanding)
}

// TestPorsiPeriodeKasusNyataEDMFire - langkah 15 atas tanggal KASUS NYATA
// fixture NB-15 (dibaca dari berkasnya, tidak disalin): EdmToStart = 150,5 hari,
// EdmToEnd = −0,5 hari, TotalPeriod = 150. Port ini menolak selisih yang bukan
// hari bulat (satuan selisih DateTime Pega dan pemotongannya ke `int` belum
// terverifikasi) - jadi kasus nyata satu-satunya ini DITOLAK. Ditulis apa
// adanya.
//
// ⚠️ Nilai tersimpan fixture (`ProrateStartEDM` 1.00666666666666666667,
// `ProrateEDMEnd` 1) BUKAN pembanding langkah 15: −0,5/150 tidak mungkin
// menjadi 1, jadi nilai itu ditulis (atau ditimpa) rule lain -
// `CountPaymentEdm_Act`, lihat laporan.
func TestPorsiPeriodeKasusNyataEDMFire(t *testing.T) {
	fx := muatFixture(t)
	halLama := hal(fx, "OldData")
	periodeLama := hal(halLama, "PolicyData")
	lama := offerFireDari(t, halLama)
	lama.PolicyData.StartDateTime = waktuPega(t, teksDari(periodeLama, "StartDateTime"))
	lama.PolicyData.EndDateTime = waktuPega(t, teksDari(periodeLama, "EndDateTime"))
	q := hal(fx, "QuotationData")
	h, err := services.PrepareBeforeImage(services.MasukanBeforeImage{
		Kasus: models.KasusEndorsement{OfferFacIn: models.OfferFacIn{QuotationData: models.Quotation{
			StatusBusiness: teksDari(q, "StatusBusiness"), EdmDate: waktuPega(t, teksDari(q, "EdmDate")),
		}}},
		PolisLama: &lama, Predikat: services.Predikat{IsFire: true, IsEDM: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !errors.Is(h.GalatPorsiPeriode, services.ErrSatuanSelisihWaktuBelumTerverifikasi) {
		t.Fatalf("GalatPorsiPeriode %v, mau ErrSatuanSelisihWaktuBelumTerverifikasi (kasus nyata ditolak)", h.GalatPorsiPeriode)
	}
	if !h.Kasus.OfferFacIn.ProrateStartEDM.Kosong() || !h.Kasus.OfferFacIn.ProrateEDMEnd.Kosong() {
		t.Error("rasio terisi padahal ditolak")
	}
}

// TestPorsiPembayaranKasusNyata - kasus nyata NB-15 (StatusBusiness 3, EdmType
// 4, Type 1): seluruh gerbang dinilai lewat REGISTRY predikat EDM atas
// properti QuotationData fixture, lalu `CountPaymentEdm_Act` blok 13 memberi
// `ProrateEDMEnd = 1` - SAMA dengan nilai tersimpan fixture. ⚠️ Lemah: blok 13
// menulis konstanta.
//
// ⚠️ `[dugaan]` Hasil query jenis bisnis polis lama (`OutData.pxResults(1).CARI2`)
// tidak ada di fixture; diisi `OldData.QuotationData.BusinessType`. Yang
// bergantung padanya hanya IsMarineCargo (langkah 9 / 14.3 / 1.3 TSIObj), dan
// blok 13 tidak membacanya.
func TestPorsiPembayaranKasusNyata(t *testing.T) {
	fx := muatFixture(t)
	q := hal(fx, "QuotationData")
	properti := map[string]string{}
	for _, medan := range []string{"StatusBusiness", "EdmType", "Type"} {
		properti["pyWorkPage.OfferFacIn.QuotationData."+medan] = teksDari(q, medan)
	}
	jenisLama := teksDari(hal(hal(fx, "OldData"), "QuotationData"), "BusinessType")
	p, err := services.PredikatPembayaranDari(predikat.KasusEDM{Properti: properti, JenisBisnisLama: &jenisLama})
	if err != nil {
		t.Fatal(err)
	}
	mau := services.PredikatPembayaran{IsEDM: true, IsEdmExtendPeriod: true}
	if p != mau {
		t.Fatalf("predikat %+v, mau %+v", p, mau)
	}
	pd := hal(fx, "PolicyData")
	h, err := services.HitungPorsiPembayaran(services.MasukanPorsiPembayaran{
		Predikat:      p,
		EDMDay:        teksDari(q, "EDMDay"),
		StartDateTime: waktuPega(t, teksDari(pd, "StartDateTime")), EndDateTime: waktuPega(t, teksDari(pd, "EndDateTime")),
		EdmDate: waktuPega(t, teksDari(q, "EdmDate")), CacahMataUang: len(daftar(fx, "CurrencyList")),
		EdmType: models.JenisEndorsemen(teksDari(q, "EdmType")),
	})
	if err != nil {
		t.Fatal(err)
	}
	samaRasioJumlah(t, "ProrateEDMEnd", h.ProrateEDMEnd, rasioDari(t, fx, "ProrateEDMEnd"))
	if !h.ProrateStartEDM.Kosong() {
		t.Errorf("ProrateStartEDM ditulis blok 13: %s", h.ProrateStartEDM)
	}
	// Langkah 4: 17:00 GMT = 00:00 WIB → tanggal Jakarta pukul 05:00 GMT.
	for nama, w := range map[string]time.Time{"StartDateTime": h.StartDateTime, "EndDateTime": h.EndDateTime} {
		if w.Hour() != 5 || w.Minute() != 0 || w.Location() != time.UTC {
			t.Errorf("%s ternormalisasi %v", nama, w)
		}
	}
}

// TestPorsiPembayaranTanggalNyataBulat - BUKAN rekonsiliasi (predikat
// sintetis): tanggal nyata NB-15 sesudah normalisasi langkah 4 (polis →
// tanggal Jakarta 05:00 GMT) berselisih HARI BULAT dengan EdmDate tersimpan
// (`T050000.000 GMT`), sehingga blok 14 tidak menolak - berbeda dengan langkah
// 15 `SetValueToEDMWork` yang membaca OldData tanpa normalisasi.
func TestPorsiPembayaranTanggalNyataBulat(t *testing.T) {
	fx := muatFixture(t)
	q, pd := hal(fx, "QuotationData"), hal(fx, "PolicyData")
	h, err := services.HitungPorsiPembayaran(services.MasukanPorsiPembayaran{
		Predikat:      services.PredikatPembayaran{IsEDM: true, IsEdmAdjRate: true},
		EDMDay:        teksDari(q, "EDMDay"),
		StartDateTime: waktuPega(t, teksDari(pd, "StartDateTime")), EndDateTime: waktuPega(t, teksDari(pd, "EndDateTime")),
		EdmDate: waktuPega(t, teksDari(q, "EdmDate")), CacahMataUang: 1,
	})
	if err != nil {
		t.Fatalf("tanggal nyata ditolak: %v", err)
	}
	if h.Day <= 0 || h.ProrateEDMEnd.Kosong() {
		t.Errorf("%+v", h)
	}
}
