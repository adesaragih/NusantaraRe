package services_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	intiuang "nusantarare/inti/uang"
	"nusantarare/inti/utils"
	"nusantarare/modul/claimlife/services"
)

// Angka di berkas ini DIHITUNG TANGAN dari rumus SpreadingClaimLife_Act, dan
// turunannya ditulis di komentar tiap kasus. Test yang hanya mengulang apa
// yang kode lakukan tidak membuktikan apa pun; yang membuktikan adalah angka
// yang berasal dari luar kode.

func uang(t *testing.T, s, mata string) intiuang.Money {
	t.Helper()
	d, err := utils.ParseDecimal(s)
	if err != nil {
		t.Fatalf("uang %q: %v", s, err)
	}
	return intiuang.Money{Amount: d, Currency: mata}
}

func rasio(t *testing.T, s string) intiuang.Ratio {
	t.Helper()
	d, err := utils.ParseDecimal(s)
	if err != nil {
		t.Fatalf("rasio %q: %v", s, err)
	}
	return intiuang.Ratio{Value: d}
}

// masukanUji memberi peserta dan baris yang dipakai hampir seluruh kasus.
//
// EM_PERCENT 0,25 adalah PECAHAN LANGSUNG - SpreadingClaimLife_Act langkah 8.1
// menyalin .EM_PERCENT tanpa membaginya seratus.
func masukanUji(t *testing.T, tahunPolis int) services.MasukanSpreading {
	t.Helper()
	return services.MasukanSpreading{
		Currency:   "IDR",
		ClaimGross: uang(t, "250000000", "IDR"),
		EmPercent:  rasio(t, "0.25"),
		Umur:       "30",
		JenisKel:   "M",
		TahunPolis: tahunPolis,
	}
}

// treatyUji memberi dua treaty-year berurutan kapasitas menaik, seperti
// `order by TO_NUMBER(IDR) asc` pada GetJsonProductLife.
func treatyUji(t *testing.T) []services.TahunTreaty {
	t.Helper()
	kapasitasA := uang(t, "100000000", "IDR")
	kapasitasB := uang(t, "500000000", "IDR")
	retro := []services.BarisRetro{{
		ID:             "R1",
		ReinsurerName:  "UJI-REINSURER",
		TreatyTypeID:   "10196",
		TreatyTypeName: "QS",
		PercentShare:   rasio(t, "40"),
		Commision:      rasio(t, "20"),
		OvrComm:        rasio(t, "10"),
	}}
	return []services.TahunTreaty{
		{ID: "T-A", TreatyYearLife: "2018", IDR: &kapasitasA, Retro: retro},
		{ID: "T-B", TreatyYearLife: "2025", IDR: &kapasitasB, Retro: retro},
	}
}

func rateUji(t *testing.T) []services.BarisRate {
	t.Helper()
	// Koma sengaja: 90.436 dari 98.305 baris RATE_LIFE memakai desimal
	// Indonesia, dan Pega menguraikannya lewat replaceAll koma-ke-titik.
	return []services.BarisRate{
		{ID: "UJI-1", Umur: "30", Kontrak: "10", JenisKel: "U", Rate: rasio(t, "2.5")},
	}
}

// TestKaskadeKapasitasMengalirKeTreatyBerikut mengunci langkah 8.2.1.4-7.
//
// Klaim 250.000.000 terhadap kapasitas 100.000.000 lalu 500.000.000:
//
//	treaty A: 250.000.000 > 100.000.000  -> share A = 100.000.000, sisa 150.000.000
//	treaty B: 150.000.000 <= 500.000.000 -> share B = 150.000.000, sisa 0
func TestKaskadeKapasitasMengalirKeTreatyBerikut(t *testing.T) {
	hasil, err := services.HitungSpreading(masukanUji(t, 2), treatyUji(t), rateUji(t))
	if err != nil {
		t.Fatalf("HitungSpreading: %v", err)
	}
	if len(hasil) != 2 {
		t.Fatalf("baris spreading = %d, mau 2", len(hasil))
	}
	if got := hasil[0].RetrocadedShare.String(); got != "100000000 IDR" {
		t.Errorf("share treaty A = %q, mau 100000000 IDR", got)
	}
	if got := hasil[1].RetrocadedShare.String(); got != "150000000 IDR" {
		t.Errorf("share treaty B = %q, mau 150000000 IDR", got)
	}
}

// TestSisaNolTetapMenghasilkanBaris: XML tidak berhenti saat sisa habis -
// precondition `local.ClaimNet<=.IDR` tetap benar untuk nol, jadi barisnya ada
// dengan share nol. Menghilangkannya berarti menghilangkan treaty-year dari
// jejak perhitungan.
func TestSisaNolTetapMenghasilkanBaris(t *testing.T) {
	m := masukanUji(t, 2)
	m.ClaimGross = uang(t, "50000000", "IDR")
	hasil, err := services.HitungSpreading(m, treatyUji(t), rateUji(t))
	if err != nil {
		t.Fatalf("HitungSpreading: %v", err)
	}
	if len(hasil) != 2 {
		t.Fatalf("baris spreading = %d, mau 2", len(hasil))
	}
	if got := hasil[0].RetrocadedShare.String(); got != "50000000 IDR" {
		t.Errorf("share treaty A = %q, mau 50000000 IDR", got)
	}
	if got := hasil[1].RetrocadedShare.String(); got != "0 IDR" {
		t.Errorf("share treaty B = %q, mau 0 IDR", got)
	}
}

// TestTahunPertamaMemakaiDiscount mengunci langkah 8.2.1.9.3 ("Tahun pertama").
//
// Hitungan tangan untuk treaty A (share 100.000.000, PERCENTSHARE 40,
// OVR_COMM 10, COMMISION 20, rate mentah 2,5, EM 0,25):
//
//	RATE tersimpan = bulat(2,5 ; 4)            = 2,5
//	rate pakai     = bulat(2,5 / 1000 ; 10)    = 0,0025
//	AMOUNT         = 100.000.000 x bulat(40/100;4) = 40.000.000
//	GROSS          = 0,0025 x (1 + 0,25) x 40.000.000 = 125.000
//	Discount       = 125.000 x bulat(20/100;5) = 25.000
//	Comm           = (125.000 - 25.000) x bulat(10/100;5) = 10.000
//	NET            = 125.000 - 25.000 - 10.000 = 90.000
func TestTahunPertamaMemakaiDiscount(t *testing.T) {
	hasil, err := services.HitungSpreading(masukanUji(t, 1), treatyUji(t), rateUji(t))
	if err != nil {
		t.Fatalf("HitungSpreading: %v", err)
	}
	r := hasil[0].Retro[0]
	cek(t, "RATE", r.Rate.String(), "2.5")
	cek(t, "AMOUNT", r.Amount.String(), "40000000 IDR")
	cek(t, "GROSS", r.PremiumSpreadedGross.String(), "125000 IDR")
	cek(t, "NET", r.PremiumSpreadedNet.String(), "90000 IDR")
}

// TestBukanTahunPertamaTanpaDiscount mengunci langkah 8.2.1.9.2.
//
//	Comm = 125.000 x 0,1 = 12.500
//	NET  = 125.000 - 12.500 = 112.500
func TestBukanTahunPertamaTanpaDiscount(t *testing.T) {
	hasil, err := services.HitungSpreading(masukanUji(t, 2), treatyUji(t), rateUji(t))
	if err != nil {
		t.Fatalf("HitungSpreading: %v", err)
	}
	r := hasil[0].Retro[0]
	cek(t, "GROSS", r.PremiumSpreadedGross.String(), "125000 IDR")
	cek(t, "NET", r.PremiumSpreadedNet.String(), "112500 IDR")
}

// TestTreatyKeduaMemakaiShare-nya sendiri: share 150.000.000 x 40% = 60.000.000
// GROSS = 0,003125 x 60.000.000 = 187.500; Discount = 37.500;
// Comm = (187.500 - 37.500) x 0,1 = 15.000; NET = 135.000.
func TestTreatyKeduaMemakaiShareSendiri(t *testing.T) {
	hasil, err := services.HitungSpreading(masukanUji(t, 1), treatyUji(t), rateUji(t))
	if err != nil {
		t.Fatalf("HitungSpreading: %v", err)
	}
	r := hasil[1].Retro[0]
	cek(t, "AMOUNT", r.Amount.String(), "60000000 IDR")
	cek(t, "GROSS", r.PremiumSpreadedGross.String(), "187500 IDR")
	cek(t, "NET", r.PremiumSpreadedNet.String(), "135000 IDR")
}

// TestRateDibagiSeribuSekaliSaja menjaga jebakan per-mil: RATE yang TERSIMPAN
// adalah rate mentah, sedangkan yang dipakai berhitung sudah dibagi 1000.
// Menyimpan rate yang sudah dibagi berarti kehilangan angka aslinya.
func TestRateDibagiSeribuSekaliSaja(t *testing.T) {
	hasil, err := services.HitungSpreading(masukanUji(t, 2), treatyUji(t), rateUji(t))
	if err != nil {
		t.Fatalf("HitungSpreading: %v", err)
	}
	if got := hasil[0].Retro[0].Rate.String(); got != "2.5" {
		t.Fatalf("RATE tersimpan = %q, mau 2.5 (rate mentah, bukan hasil bagi seribu)", got)
	}
}

// TestRateUnisexMengabaikanJenisKelamin mengunci precondition 8.2.1.9.1:
// GENDER memuat "U" -> cocokkan UMUR saja.
func TestRateUnisexMengabaikanJenisKelamin(t *testing.T) {
	r, err := services.PilihRate(rateUji(t), "30", "F")
	if err != nil {
		t.Fatalf("PilihRate: %v", err)
	}
	if got := r.String(); got != "2.5" {
		t.Errorf("rate = %q, mau 2.5", got)
	}
}

// TestRateBerjenisKelaminMenuntutKecocokan mengunci precondition 8.2.1.9.2:
// GENDER tidak memuat "U" -> UMUR dan JENIS KELAMIN harus cocok keduanya.
func TestRateBerjenisKelaminMenuntutKecocokan(t *testing.T) {
	baris := []services.BarisRate{
		{ID: "UJI-2", Umur: "30", JenisKel: "F", Rate: rasio(t, "3.1")},
	}
	if _, err := services.PilihRate(baris, "30", "M"); !errors.Is(err, services.ErrRateTidakDitemukan) {
		t.Fatalf("galat = %v, mau ErrRateTidakDitemukan", err)
	}
	r, err := services.PilihRate(baris, "30", "F")
	if err != nil {
		t.Fatalf("PilihRate cocok: %v", err)
	}
	if got := r.String(); got != "3.1" {
		t.Errorf("rate = %q, mau 3.1", got)
	}
}

// TestRateGandaDilaporkanBukanDipilih: GetRateRetro tidak memakai ORDER BY dan
// putaran Pega tidak berhenti pada yang pertama cocok, jadi "yang menang" di
// sistem lama tidak tertentu. Go melaporkan, tidak memilih diam-diam.
func TestRateGandaDilaporkanBukanDipilih(t *testing.T) {
	baris := []services.BarisRate{
		{ID: "UJI-3", Umur: "30", JenisKel: "U", Rate: rasio(t, "2.5")},
		{ID: "UJI-4", Umur: "30", JenisKel: "U", Rate: rasio(t, "9.9")},
	}
	_, err := services.PilihRate(baris, "30", "M")
	if !errors.Is(err, services.ErrRateBerganda) {
		t.Fatalf("galat = %v, mau ErrRateBerganda", err)
	}
	// Pesannya harus menyebut ID barisnya, supaya DBA dapat menelusurinya.
	if !strings.Contains(err.Error(), "UJI-3") || !strings.Contains(err.Error(), "UJI-4") {
		t.Errorf("pesan tidak menyebut kedua ID: %v", err)
	}
}

// TestRateGandaBernilaiSamaBukanAmbiguitas: dua baris yang nilainya identik
// tidak membuat hasilnya tidak tentu, jadi ia bukan galat.
func TestRateGandaBernilaiSamaBukanAmbiguitas(t *testing.T) {
	baris := []services.BarisRate{
		{ID: "UJI-5", Umur: "30", JenisKel: "U", Rate: rasio(t, "2.5")},
		{ID: "UJI-6", Umur: "30", JenisKel: "U", Rate: rasio(t, "2.50")},
	}
	r, err := services.PilihRate(baris, "30", "M")
	if err != nil {
		t.Fatalf("PilihRate: %v", err)
	}
	if got := r.String(); got != "2.5" {
		t.Errorf("rate = %q, mau 2.5", got)
	}
}

// TestTanpaRateCocokGagalTerang - Pega memakai rate baris sebelumnya yang
// SUDAH dibagi seribu lalu membaginya seribu lagi. Itu cacat, bukan aturan.
func TestTanpaRateCocokGagalTerang(t *testing.T) {
	baris := []services.BarisRate{
		{ID: "UJI-7", Umur: "45", JenisKel: "U", Rate: rasio(t, "2.5")},
	}
	if _, err := services.PilihRate(baris, "30", "M"); !errors.Is(err, services.ErrRateTidakDitemukan) {
		t.Fatalf("galat = %v, mau ErrRateTidakDitemukan", err)
	}
}

// TestKapasitasKosongMenggagalkanSpreading: TREATYYEAR_LIFE di POOLDATA tidak
// punya kolom IDR/USD. Menganggap kapasitas kosong sebagai nol akan membuat
// seluruh share nol tanpa seorang pun tahu.
func TestKapasitasKosongMenggagalkanSpreading(t *testing.T) {
	treaty := treatyUji(t)
	treaty[0].IDR = nil
	_, err := services.HitungSpreading(masukanUji(t, 1), treaty, rateUji(t))
	if !errors.Is(err, services.ErrKapasitasBelumDiketahui) {
		t.Fatalf("galat = %v, mau ErrKapasitasBelumDiketahui", err)
	}
}

// TestMataUangTanpaCabangDitolak: XML hanya punya cabang IDR dan USD.
func TestMataUangTanpaCabangDitolak(t *testing.T) {
	m := masukanUji(t, 1)
	m.Currency = "SGD"
	m.ClaimGross = uang(t, "1000", "SGD")
	_, err := services.HitungSpreading(m, treatyUji(t), rateUji(t))
	if !errors.Is(err, services.ErrMataUangTanpaCabang) {
		t.Fatalf("galat = %v, mau ErrMataUangTanpaCabang", err)
	}
}

// TestTanpaTreatyNolBaris - klaim yang seluruhnya ditahan sendiri adalah
// keadaan biasa, bukan kegagalan.
func TestTanpaTreatyNolBaris(t *testing.T) {
	hasil, err := services.HitungSpreading(masukanUji(t, 1), nil, rateUji(t))
	if err != nil {
		t.Fatalf("HitungSpreading: %v", err)
	}
	if len(hasil) != 0 {
		t.Fatalf("baris = %d, mau 0", len(hasil))
	}
}

// TestTreatyTanpaRetroTetapBerbaris: baris spreading lahir dari daftar
// treaty-year, bukan dari daftar reinsurer. Nol reinsurer berarti nol anak,
// bukan nol baris.
func TestTreatyTanpaRetroTetapBerbaris(t *testing.T) {
	treaty := treatyUji(t)
	treaty[0].Retro = nil
	treaty[1].Retro = nil
	hasil, err := services.HitungSpreading(masukanUji(t, 1), treaty, rateUji(t))
	if err != nil {
		t.Fatalf("HitungSpreading: %v", err)
	}
	if len(hasil) != 2 {
		t.Fatalf("baris = %d, mau 2", len(hasil))
	}
	if len(hasil[0].Retro) != 0 {
		t.Fatalf("anak = %d, mau 0", len(hasil[0].Retro))
	}
}

// TestTahunPolisDihitungDariTahunSaja mengunci langkah 8.1:
// local.Year = tahun(GROSS_VALUATION_BEGIN_DATE) - tahun(BEGIN_DATE) + 1.
func TestTahunPolisDihitungDariTahunSaja(t *testing.T) {
	kasus := []struct {
		mulai, valuasi string
		mau            int
	}{
		{"2020-05-01", "2020-08-01", 1},
		{"2020-12-31", "2021-01-01", 2},
		{"2018-01-01", "2025-06-30", 8},
	}
	for _, k := range kasus {
		mulai, err := time.Parse("2006-01-02", k.mulai)
		if err != nil {
			t.Fatalf("mulai: %v", err)
		}
		valuasi, err := time.Parse("2006-01-02", k.valuasi)
		if err != nil {
			t.Fatalf("valuasi: %v", err)
		}
		if got := services.TahunPolis(mulai, valuasi); got != k.mau {
			t.Errorf("TahunPolis(%s, %s) = %d, mau %d", k.mulai, k.valuasi, got, k.mau)
		}
	}
}

func cek(t *testing.T, nama, got, mau string) {
	t.Helper()
	if got != mau {
		t.Errorf("%s = %q, mau %q", nama, got, mau)
	}
}

// TestUangDibulatkanDelapanDesimal menutup ADR-U-0016 Akibat 3: lapisan
// services membulatkan pada desimal kedelapan SAAT MEMUAT.
//
// Tanpa pembulatan itu, Oracle-lah yang membulatkan diam-diam saat menulis ke
// NUMBER(38,8), sehingga angka yang dihitung dan angka yang tersimpan berbeda
// tanpa satu pun galat.
//
// Hitungan tangan, dipilih supaya hasilnya PASTI melewati delapan desimal:
//
//	rate mentah 2,4681  -> rate pakai = bulat(0,0024681 ; 10) = 0,0024681
//	AMOUNT  = 1 x bulat(100/100 ; 4) = 1
//	GROSS   = 0,0024681 x 1,25 x 1 = 0,003085125   <- SEMBILAN desimal
//	dibulatkan delapan                = 0,00308513
func TestUangDibulatkanDelapanDesimal(t *testing.T) {
	kapasitas := uang(t, "1", "IDR")
	treaty := []services.TahunTreaty{{
		ID: "T-A", TreatyYearLife: "2025", IDR: &kapasitas,
		Retro: []services.BarisRetro{{
			ID: "R1", ReinsurerName: "UJI-REINSURER",
			PercentShare: rasio(t, "100"),
			Commision:    rasio(t, "0"),
			OvrComm:      rasio(t, "0"),
		}},
	}}
	m := services.MasukanSpreading{
		Currency:   "IDR",
		ClaimGross: uang(t, "1", "IDR"),
		EmPercent:  rasio(t, "0.25"),
		Umur:       "30",
		JenisKel:   "M",
		TahunPolis: 2,
	}
	rate := []services.BarisRate{
		{ID: "UJI-8", Umur: "30", JenisKel: "U", Rate: rasio(t, "2.4681")},
	}
	hasil, err := services.HitungSpreading(m, treaty, rate)
	if err != nil {
		t.Fatalf("HitungSpreading: %v", err)
	}
	r := hasil[0].Retro[0]
	cek(t, "GROSS", r.PremiumSpreadedGross.String(), "0.00308513 IDR")
	// NET diturunkan dari GROSS yang SUDAH dibulatkan, supaya angka yang
	// tersimpan konsisten satu sama lain: NET = GROSS - Comm berlaku pada
	// nilai yang benar-benar masuk basis data, bukan hanya pada nilai antara.
	cek(t, "NET", r.PremiumSpreadedNet.String(), "0.00308513 IDR")
}

// TestTreatyTidakTerurutDitolak - urutan kaskade ditentukan XML
// (`order by TO_NUMBER(IDR) asc`). Daftar yang tiba tanpa urutan itu membagi
// klaim ke treaty-year yang salah, dan angkanya tetap tampak masuk akal.
func TestTreatyTidakTerurutDitolak(t *testing.T) {
	terurut := treatyUji(t)
	terbalik := []services.TahunTreaty{terurut[1], terurut[0]}
	_, err := services.HitungSpreading(masukanUji(t, 1), terbalik, rateUji(t))
	if !errors.Is(err, services.ErrTreatyTidakTerurut) {
		t.Fatalf("galat = %v, mau ErrTreatyTidakTerurut", err)
	}
}
