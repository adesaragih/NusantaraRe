package models

// Validasi unggahan CSV - tiket 04 PremiumList Life. TANPA Oracle.

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/pkg/utils"
)

// TestUangCSVMenolakPemisahRibuan - AC tiket 04 mewajibkan kasus ini.
//
// ⛔ `1,234,567.89` dan `1.234.567,89` TIDAK boleh keduanya diterima diam-diam
// sebagai angka yang sama. Keduanya ditolak, dan pesannya menyebut sebabnya.
//
// ⛔ Pega justru MENGGANTI setiap koma dengan titik
// (`@replaceAll(.KOLOM, ",", ".")`), sehingga `1,234,567.89` menjadi
// `1.234.567.89` - bukan angka sama sekali - lalu tersimpan sebagai uang.
func TestUangCSVMenolakPemisahRibuan(t *testing.T) {
	for _, buruk := range []string{
		"1,234,567.89", // ribuan koma, desimal titik
		"1.234.567,89", // ribuan titik, desimal koma
		"1,5",          // desimal koma
		"1,000",        // ribuan koma tanpa desimal
		"1 234 567",    // ribuan spasi
	} {
		d, err := UangCSV(buruk)
		if err == nil {
			t.Errorf("%q diterima sebagai %v; mau ditolak", buruk, d)
			continue
		}
		if strings.ContainsRune(buruk, ',') && !errors.Is(err, ErrUangBerkoma) {
			t.Errorf("%q: galat %v, mau ErrUangBerkoma", buruk, err)
		}
	}
}

func TestUangCSVMenerimaDesimalTitik(t *testing.T) {
	for _, k := range []struct{ masuk, mau string }{
		{"1234567.89", "1234567.89"},
		{"0", "0"},
		{"0.00000001", "0.00000001"},
		{"-5.25", "-5.25"},
		{"  42.5  ", "42.5"},
		{"1234567890123456789.12345678", "1234567890123456789.12345678"},
	} {
		d, err := UangCSV(k.masuk)
		if err != nil {
			t.Errorf("%q ditolak: %v", k.masuk, err)
			continue
		}
		if got := utils.FormatDecimal(d); got != k.mau {
			t.Errorf("%q -> %q, mau %q", k.masuk, got, k.mau)
		}
	}
}

// TestUangCSVTidakMembulatkan - AC tiket 04.
//
// ⛔ Nilai yang lolos validasi harus IDENTIK dengan yang diunggah. Delapan
// angka di belakang koma adalah presisi kolom `NUMBER(38,8)`; satu pun yang
// hilang di sini adalah uang yang berubah tanpa ada yang meminta.
func TestUangCSVTidakMembulatkan(t *testing.T) {
	const asli = "12345678901234567890.12345678"
	d, err := UangCSV(asli)
	if err != nil {
		t.Fatalf("%q ditolak: %v", asli, err)
	}
	if got := utils.FormatDecimal(d); got != asli {
		t.Errorf("pulang sebagai %q, mau %q - ada pembulatan diam", got, asli)
	}
}

// TestUangCSVMenolakYangBukanBilangan - termasuk jebakan pustaka desimal.
//
// ⛔ `apd.NewFromString` menerima `Infinity` dan `NaN`. Keduanya BUKAN uang,
// dan keduanya akan tersimpan ke kolom `NUMBER` sebagai sesuatu yang tidak
// seorang pun maksudkan.
func TestUangCSVMenolakYangBukanBilangan(t *testing.T) {
	for _, buruk := range []string{"", "   ", "abc", "Rp1000", "1.2.3"} {
		if d, err := UangCSV(buruk); err == nil {
			t.Errorf("%q diterima sebagai %v; mau ditolak", buruk, d)
		}
	}
	// ⚠️ Kalau salah satu di bawah LOLOS, penjaganya yang harus diperketat -
	// bukan uji ini yang dilonggarkan.
	for _, aneh := range []string{"NaN", "Infinity", "-Infinity", "Inf", "1e400"} {
		d, err := UangCSV(aneh)
		if err != nil {
			continue
		}
		teks := utils.FormatDecimal(d)
		if strings.ContainsAny(teks, "nNiI") {
			t.Errorf("%q lolos menjadi %q - bukan bilangan, dan ia akan "+
				"tersimpan sebagai uang", aneh, teks)
		}
	}
}

func TestTanggalCSVKetatDDMMYYYY(t *testing.T) {
	if _, err := TanggalCSV("31/01/2026"); err != nil {
		t.Errorf("31/01/2026 ditolak: %v", err)
	}
	for _, buruk := range []string{
		"2026-01-31", "1/2/2026", "31-01-2026", "31/13/2026", "32/01/2026", "",
	} {
		if _, err := TanggalCSV(buruk); err == nil {
			t.Errorf("%q diterima; mau ditolak", buruk)
		}
	}
}

// TestTanggalEpochDiterima - pengecualian `02/01/1970` di Pega.
//
// ⛔ Pega mengecualikannya (`@if(.DOB=="02/01/1970",true,...)`) sebab
// `@toDate` di sana mengembalikan NOL untuk tanggal itu, dan pemeriksanya
// membandingkan dengan nol. Pengurai kami tidak punya cacat itu, jadi tanggal
// lahir 2 Januari 1970 diterima TANPA pengecualian - dan uji ini menjaga
// supaya tidak ada yang "merapikannya" menjadi aturan yang menolaknya.
func TestTanggalEpochDiterima(t *testing.T) {
	got, err := TanggalCSV(TanggalEpochCSV)
	if err != nil {
		t.Fatalf("%s ditolak: %v", TanggalEpochCSV, err)
	}
	if got.Day() != 2 || got.Month() != 1 || got.Year() != 1970 {
		t.Errorf("terurai menjadi %v", got)
	}
}

// TestKolomUangUnggahTigaPuluhDua - sensus, dan ralat tiket.
func TestKolomUangUnggahTigaPuluhDua(t *testing.T) {
	const mau = 32
	if n := len(KolomUangUnggah); n != mau {
		t.Errorf("%d kolom uang, mau %d.\n"+
			"Jendela: precondition @PropertyHasValue di langkah 2 "+
			"ValidasiUploadPL_act.xml. Tiket 04 menulis 33; yang benar 32 - "+
			"cara kedua (deskripsi sub-langkah) menjawab 31 karena "+
			"FLEET_DISCOUNT deskripsinya kosong.", n, mau)
	}
	lihat := map[string]bool{}
	for _, k := range KolomUangUnggah {
		if lihat[k] {
			t.Errorf("kolom %q muncul dua kali", k)
		}
		lihat[k] = true
	}
	// Selisih kedua cara itu bernama, dan ia HARUS ada di daftar.
	if !lihat["FLEET_DISCOUNT"] {
		t.Error("FLEET_DISCOUNT hilang - itu persis kolom yang cara kedua lewatkan")
	}
	// Keenam yang wajib adalah bagian dari ketiga puluh dua.
	for k := range kolomUangWajib {
		if !lihat[k] {
			t.Errorf("kolom wajib %q tidak ada di daftar kolom uang", k)
		}
	}
}

// barisSah menyusun satu baris yang seluruh kolomnya benar.
//
// ⛔ NOL NAMA ORANG, nol nomor polis nyata. `UJI-*` adalah penanda fixture
// buatan, dan ia sengaja tidak menyerupai data siapa pun.
func barisSah(nomor int) BarisUnggah {
	n := map[string]string{
		"CERTIFICATE_NO":  "UJI-CERT-" + strings.Repeat("0", 3) + itoa(nomor),
		"NAME_OF_INSURED": "UJI PESERTA " + itoa(nomor),
		"POLICY_NO":       "UJI-POL-1",
		"PLAN":            "UJI-PLAN",
		"CURRENCY":        "IDR",
		"MEDICAL_STATUS":  "NM",
	}
	for _, k := range []string{
		"DOB", "BEGIN_DATE", "EXPIRED_DATE", "START_DATE", "EFFECTIVE_DATE",
		"STNC", "WPC", "GROSS_VALUATION_BEGIN_DATE",
		"GROSS_VALUATION_EXPIRED_DATE", "RETROCESSION_VALUATION_BEGIN_DATE",
		"RETROCESSION_VALUATION_EXPIRED_DATE",
	} {
		n[k] = "01/02/2026"
	}
	for k := range kolomUangWajib {
		n[k] = "1000.50"
	}
	return BarisUnggah{Nomor: nomor, Nilai: n}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

func TestValidasiUnggahMelewatkanBerkasYangBenar(t *testing.T) {
	h := ValidasiUnggah([]BarisUnggah{barisSah(1), barisSah(2)})
	if !h.Lolos() {
		for _, p := range h.Ditolak {
			t.Errorf("baris %d kolom %s: %s (%s)", p.Baris, p.Kolom, p.Pesan, p.Sebab)
		}
	}
	if h.CacahBaris != 2 {
		t.Errorf("cacah baris %d, mau 2", h.CacahBaris)
	}
}

// TestSetiapPenolakanMenyebutKolomDanBaris - AC tiket 04.
func TestSetiapPenolakanMenyebutKolomDanBaris(t *testing.T) {
	b := barisSah(7)
	b.Nilai["PLAN"] = ""
	b.Nilai["SUM_INSURED"] = "1,000"
	b.Nilai["DOB"] = "1970-01-02"
	h := ValidasiUnggah([]BarisUnggah{b})
	if h.Lolos() {
		t.Fatal("berkas rusak diterima")
	}
	for _, p := range h.Ditolak {
		if p.Baris != 7 {
			t.Errorf("penolakan %q menyebut baris %d, mau 7", p.Kolom, p.Baris)
		}
		if p.Kolom == "" {
			t.Errorf("penolakan tanpa nama kolom: %+v", p)
		}
		if p.Pesan == "" || p.Sebab == "" {
			t.Errorf("penolakan tanpa pesan atau sebab: %+v", p)
		}
	}
	// Ketiganya tertangkap, bukan hanya yang pertama.
	kena := map[string]string{}
	for _, p := range h.Ditolak {
		kena[p.Kolom] = p.Sebab
	}
	for _, k := range []string{"PLAN", "SUM_INSURED", "DOB"} {
		if _, ada := kena[k]; !ada {
			t.Errorf("kolom %q tidak dilaporkan", k)
		}
	}
	// ⛔ Sebab koma MENYEBUT komanya - pesan verbatim "SUM INSURED HARUS ADA"
	// sendirian akan berbohong tentang penyebabnya.
	if s := kena["SUM_INSURED"]; !strings.Contains(s, "koma") {
		t.Errorf("sebab SUM_INSURED = %q; mau menyebut koma", s)
	}
}

// TestNetPremiumHanyaDiperiksaKeberadaannya - OQ-069 TERBUKA.
//
// ⛔ Pesannya menjanjikan "LEBIH BESAR DARI GROSS PREMIUM", tetapi korpus
// TIDAK pernah membandingkannya. Mengarang perbandingan itu akan menolak
// berkas yang di sistem lama diterima - dan arah yang salah (net lazimnya
// lebih KECIL dari gross) akan menolak setiap berkas yang benar.
func TestNetPremiumHanyaDiperiksaKeberadaannya(t *testing.T) {
	b := barisSah(1)
	b.Nilai["NET_PREMIUM"] = "1.00"
	b.Nilai["GROSS_PREMIUM"] = "999999.00"
	h := ValidasiUnggah([]BarisUnggah{b})
	if !h.Lolos() {
		for _, p := range h.Ditolak {
			t.Errorf("net < gross ditolak di %s: %s (%s)", p.Kolom, p.Pesan, p.Sebab)
		}
		t.Error("OQ-069 masih terbuka: JANGAN mengarang perbandingan net vs gross")
	}
	// Pesannya tetap disalin verbatim, apa adanya.
	if !strings.Contains(PesanNetPremium, "LEBIH BESAR DARI GROSS PREMIUM") {
		t.Error("pesan verbatim OQ-069 berubah")
	}
}

// TestKolomUangTakWajibTetapDiperiksaBentuknya.
//
// ⛔ Aturan "separator titik" dinyatakan enam kali di nama langkah korpus dan
// ditegakkan SEKALI (`@contains(.GROSS_PREMIUM,",")`). Kolom uang yang tidak
// diperiksa adalah kolom uang yang menerima apa saja.
func TestKolomUangTakWajibTetapDiperiksaBentuknya(t *testing.T) {
	b := barisSah(1)
	b.Nilai["TAX"] = "1,5"
	h := ValidasiUnggah([]BarisUnggah{b})
	found := false
	for _, p := range h.Ditolak {
		if p.Kolom == "TAX" {
			found = true
			if !strings.Contains(p.Sebab, "koma") {
				t.Errorf("sebab TAX = %q, mau menyebut koma", p.Sebab)
			}
		}
	}
	if !found {
		t.Error("kolom uang tak wajib yang berkoma TIDAK ditolak")
	}
	// Tetapi KOSONG pada kolom tak wajib bukan galat.
	c := barisSah(2)
	c.Nilai["TAX"] = ""
	if h2 := ValidasiUnggah([]BarisUnggah{c}); !h2.Lolos() {
		for _, p := range h2.Ditolak {
			t.Errorf("kolom tak wajib kosong ditolak: %+v", p)
		}
	}
}

func TestSeluruhBarisDiperiksaBukanHanyaYangPertama(t *testing.T) {
	a, b := barisSah(1), barisSah(2)
	a.Nilai["PLAN"] = ""
	b.Nilai["CURRENCY"] = ""
	h := ValidasiUnggah([]BarisUnggah{a, b})
	if n := len(h.BarisDitolak()); n != 2 {
		t.Errorf("%d baris dilaporkan, mau 2 - berhenti di baris pertama "+
			"memaksa orang mengunggah ulang sebanyak jumlah kesalahannya", n)
	}
}

// TestPenolakanLangkahTerRemarkTidakDitegakkan - sensus remark 28-09-2026.
//
// ⛔ Sertifikat ganda (9.5/10/13), sertifikat terpakai (9.2-9.4/12),
// `MEDICAL_STATUS` (9.11/19), dan lampiran (11/33) ter-remark di
// `ValidasiUploadPL_act`; duplikat nama+DOB tidak pernah menjadi penolakan
// unggah. Berkas yang hanya "melanggar" kelimanya LOLOS - seperti di Pega.
func TestPenolakanLangkahTerRemarkTidakDitegakkan(t *testing.T) {
	a, b := barisSah(1), barisSah(2)
	b.Nilai["CERTIFICATE_NO"] = a.Nilai["CERTIFICATE_NO"]
	b.Nilai["NAME_OF_INSURED"] = strings.ToLower(a.Nilai["NAME_OF_INSURED"])
	b.Nilai["DOB"] = a.Nilai["DOB"]
	a.Nilai["MEDICAL_STATUS"] = "X"
	delete(b.Nilai, "MEDICAL_STATUS")
	if h := ValidasiUnggah([]BarisUnggah{a, b}); !h.Lolos() {
		t.Errorf("penolakan langkah ter-remark muncul: %+v", h.Ditolak)
	}
	for _, k := range KolomWajibUnggah() {
		if k == "MEDICAL_STATUS" {
			t.Error("MEDICAL_STATUS masih kolom judul wajib; pemeriksanya ter-remark")
		}
	}
}
