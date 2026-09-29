package services_test

// Kontrak penyerahan ke Komite - TANPA Oracle.
//
// Pemilik: tiket 10. Dibaca sesudah: komite.go.

import (
	"context"
	"errors"
	"testing"

	"nusantarare/inti"
	"nusantarare/inti/galat"
	"nusantarare/inti/kontrak"
	"nusantarare/inti/utils"
	"nusantarare/modul/claimlife/models"
	"nusantarare/modul/claimlife/services"
)

func barisLengkap(t *testing.T) models.BarisAdjustment {
	t.Helper()
	return models.BarisAdjustment{
		ID: "UJI-ADJ-1", KodeStatus: kontrak.KodeOutstanding,
		JumlahKlaim: uang(t, "1500000.00", "IDR"), CurrencyID: "IDR",
		NamaBank: "UJI-BANK", IDBank: "UJI-006", NomorRekening: "0012345",
	}
}

// TestGerbangRekeningMenolakTiapMedanKosong mengunci prasyarat XML.
//
// `[terverifikasi]` `Activity/GetListKomiteLife.xml` pecahan baris 655:
// `.NameOfBank==""||.NoAccount==""||.IDOfBank==""`, dengan
// `pyStepsPreCondParamsWhenTrue=2` (lanjut ke langkah galat) dan
// `WhenFalse=3` (lewati). Ketiganya setara - satu kosong sudah cukup.
func TestGerbangRekeningMenolakTiapMedanKosong(t *testing.T) {
	for _, k := range []struct {
		apa  string
		ubah func(*models.BarisAdjustment)
	}{
		{"nama bank kosong", func(b *models.BarisAdjustment) { b.NamaBank = "" }},
		{"id bank kosong", func(b *models.BarisAdjustment) { b.IDBank = "" }},
		{"nomor rekening kosong", func(b *models.BarisAdjustment) { b.NomorRekening = "" }},
		// ⛔ Spasi bukan isi. Korpus membandingkan dengan "" apa adanya;
		// membiarkan spasi lolos berarti gerbangnya dapat ditembus spasi.
		{"nama bank hanya spasi", func(b *models.BarisAdjustment) { b.NamaBank = "   " }},
	} {
		b := barisLengkap(t)
		k.ubah(&b)
		if err := services.PeriksaRekening(b); !errors.Is(err, services.ErrRekeningBelumLengkap) {
			t.Errorf("%s: galat = %v, mau ErrRekeningBelumLengkap", k.apa, err)
		}
	}
	if err := services.PeriksaRekening(barisLengkap(t)); err != nil {
		t.Errorf("baris lengkap ditolak: %v", err)
	}
}

// TestPesanRekeningPersisSepertiXML - teksnya kontrak dengan pengguna lama.
//
// `[terverifikasi]` `GetListKomiteLife.xml` pecahan baris 449-450:
// `local.errmsg = "Name of bank cannot be empty"`, dipasang ke `.NameOfBank`
// sebagai `pyMessageLabel`.
func TestPesanRekeningPersisSepertiXML(t *testing.T) {
	const mau = "Name of bank cannot be empty"
	if got := services.ErrRekeningBelumLengkap.Error(); got != mau {
		t.Errorf("pesan = %q, mau %q", got, mau)
	}
}

// TestAmbangRosterMemakaiNilaiMutlak mengunci `@if(IsADj<0,IsADj* -1,IsADj)`.
//
// `[terverifikasi]` `GetListKomiteLife.xml` pecahan baris 790.
func TestAmbangRosterMemakaiNilaiMutlak(t *testing.T) {
	for _, k := range []struct{ masuk, mau string }{
		{"1500000.00", "1500000.00"},
		{"-1500000.00", "1500000.00"},
		{"0", "0"},
		{"-0.00000001", "0.00000001"},
	} {
		got, err := services.AmbangRoster(uang(t, k.masuk, "IDR"))
		if err != nil {
			t.Fatalf("%s: %v", k.masuk, err)
		}
		// ⛔ Dibandingkan lewat FormatDecimal, BUKAN Amount.String(). Keduanya
		// berbeda: `String()` merender -0.00000001 sebagai "1E-8", sedangkan
		// `Text('f')` - yang dipakai MarshalJSON dan karena itu yang benar-benar
		// menyeberang - merendernya "0.00000001". Menguji lewat String() berarti
		// menguji bentuk yang tidak pernah dilihat Komite.
		if got := utils.FormatDecimal(got.Amount); got != k.mau {
			t.Errorf("AmbangRoster(%s) = %s, mau %s", k.masuk, got, k.mau)
		}
		if got.Currency != "IDR" {
			t.Errorf("AmbangRoster(%s) kehilangan mata uang: %q", k.masuk, got.Currency)
		}
	}
}

// TestStatusKlaimRosterAdalahLIFE - `[terverifikasi]` `GetListKomiteLife.xml`
// pecahan baris 810-811: `Param.STS_KLAIM = "LIFE"`.
func TestStatusKlaimRosterAdalahLIFE(t *testing.T) {
	if services.StatusKlaimRoster != "LIFE" {
		t.Errorf("StatusKlaimRoster = %q, mau \"LIFE\"", services.StatusKlaimRoster)
	}
}

// TestRosterBawaanGagalTerang - butir af masih `[USULAN]`, dan tabel
// `POOLDATA.EMAILKOMITE` tabel produksi. Yang belum diputuskan terlihat
// sebagai satu galat yang menyebut apa yang ditunggu.
func TestRosterBawaanGagalTerang(t *testing.T) {
	_, err := services.RosterBelumDiputuskan{}.AmbilAnggota(
		context.Background(), uang(t, "1", "IDR"), services.StatusKlaimRoster)
	if !errors.Is(err, services.ErrRosterBelumDiputuskan) {
		t.Fatalf("galat = %v, mau ErrRosterBelumDiputuskan", err)
	}
	if !errors.Is(err, services.ErrRosterBelumDiputuskan) {
		t.Error("galat tidak dapat dikenali pemanggil")
	}
}

// TestRosterKosongGagalTerangBukanTanggaNolTingkat - `[keputusan work owner]`.
//
// ⚠️ Ini penjaga defensif, dan ia BERBEDA dari Pega: `[terverifikasi]`
// `CreateKMTLife_Act` men-set `KomiteCount = 1` meski `KomiteLoop = 0`, jadi
// sistem lama diam-diam membuat tangga satu tingkat dari roster kosong.
func TestRosterKosongGagalTerangBukanTanggaNolTingkat(t *testing.T) {
	if err := services.PeriksaTingkatKomite(0); !errors.Is(
		err, services.ErrRosterKomiteKosong) {
		t.Errorf("roster kosong: galat = %v, mau ErrRosterKomiteKosong", err)
	}
	if err := services.PeriksaTingkatKomite(1); err != nil {
		t.Errorf("roster satu tingkat ditolak: %v", err)
	}
}

// TestMataUangSatuKlaimWajibSeragam - invariant OQ-060.
func TestMataUangSatuKlaimWajibSeragam(t *testing.T) {
	seragam := []models.Peserta{{Baris: []models.BarisAdjustment{
		{CurrencyID: "IDR"}, {CurrencyID: "IDR"},
	}}}
	if err := services.PeriksaSatuMataUang(seragam); err != nil {
		t.Errorf("klaim seragam ditolak: %v", err)
	}
	campur := []models.Peserta{{Baris: []models.BarisAdjustment{
		{CurrencyID: "IDR"}, {CurrencyID: "USD"},
	}}}
	if err := services.PeriksaSatuMataUang(campur); !errors.Is(
		err, services.ErrMataUangKlaimCampur) {
		t.Errorf("klaim campur: galat = %v, mau ErrMataUangKlaimCampur", err)
	}

	// ⛔ DUA kolom mata uang, dan yang dijaga harus yang MENYEBERANG.
	// `CurrencyID` datang dari `CURRENCY_ID`; `JumlahKlaim.Currency` datang
	// dari `CURRENCY`. Yang ikut ke Komite lewat AmbangRoster adalah yang
	// kedua. Ronde pertama hanya menjaga yang pertama, sehingga klaim seperti
	// di bawah - `CURRENCY_ID` seragam, `CURRENCY` campur - lolos utuh.
	campurYangMenyeberang := []models.Peserta{{
		Baris: []models.BarisAdjustment{
			{ID: "A-1", CurrencyID: "IDR", JumlahKlaim: uang(t, "10", "IDR")},
			{ID: "A-2", CurrencyID: "IDR", JumlahKlaim: uang(t, "20", "USD")},
		}}}
	if err := services.PeriksaSatuMataUang(campurYangMenyeberang); !errors.Is(
		err, services.ErrMataUangKlaimCampur) {
		t.Errorf("CURRENCY campur di balik CURRENCY_ID seragam: galat = %v, "+
			"mau ErrMataUangKlaimCampur", err)
	}
}

// TestSerahkanMenjagaPagarnya - identitas, pengenal, Oracle.
func TestSerahkanMenjagaPagarnya(t *testing.T) {
	svc := services.New(nil)
	ctx := context.Background()
	if err := svc.Komite().Serahkan(ctx, inti.Pelaku{}, "CLM-1", "P-1", "A-1",
		saatUji); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: galat = %v, mau ErrTanpaIdentitas", err)
	}
	if err := svc.Komite().Serahkan(ctx, pelakuBerperan(inti.PeranSPV),
		"  ", "P-1", "A-1", saatUji); !errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Errorf("pengenal kosong: galat = %v, mau ErrPermintaanTidakSah", err)
	}
	if err := svc.Komite().Serahkan(ctx, pelakuBerperan(inti.PeranSPV),
		"CLM-1", "P-1", "A-1", saatUji); err == nil {
		t.Error("tanpa Oracle: lolos tanpa galat")
	}
}

// TestSerahkanHanyaDariOutstanding - baris final tidak diserahkan lagi.
func TestSerahkanHanyaDariOutstanding(t *testing.T) {
	for _, kode := range []string{kontrak.KodeAksep, kontrak.KodeDitolak} {
		b := barisLengkap(t)
		b.KodeStatus = kode
		if err := services.PeriksaBolehDiserahkan(b); !errors.Is(
			err, services.ErrBarisBukanOutstanding) {
			t.Errorf("kode %q: galat = %v, mau ErrBarisBukanOutstanding", kode, err)
		}
	}
	if err := services.PeriksaBolehDiserahkan(barisLengkap(t)); err != nil {
		t.Errorf("baris Outstanding ditolak: %v", err)
	}
	// ⛔ Baris yang SUDAH punya KomiteID tidak diserahkan dua kali.
	b := barisLengkap(t)
	b.KomiteID = "KMT-000001"
	if err := services.PeriksaBolehDiserahkan(b); !errors.Is(
		err, services.ErrBarisSudahDiserahkan) {
		t.Errorf("baris sudah diserahkan: galat = %v, mau ErrBarisSudahDiserahkan", err)
	}
}
