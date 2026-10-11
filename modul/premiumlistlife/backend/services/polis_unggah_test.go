package services

// Pengurai CSV unggahan - tiket 04 PremiumList Life. TANPA Oracle.

import (
	"errors"
	"os"
	"strings"
	"testing"

	"nusantarare/modul/premiumlistlife/backend/models"
)

// bacaBerkasLayanan membaca satu berkas layanan untuk penjaga statik.
func bacaBerkasLayanan(nama string) (string, error) {
	b, err := os.ReadFile(nama)
	return string(b), err
}

// judulLengkap menyusun baris judul yang memuat seluruh kolom wajib.
func judulLengkap(tambahan ...string) []string {
	k := []string{
		"CERTIFICATE_NO", "NAME_OF_INSURED", "PLAN", "CURRENCY", "POLICY_NO",
		"MEDICAL_STATUS", "DOB", "BEGIN_DATE", "EXPIRED_DATE", "START_DATE",
		"EFFECTIVE_DATE", "STNC", "WPC",
		"GROSS_VALUATION_BEGIN_DATE", "GROSS_VALUATION_EXPIRED_DATE",
		"RETROCESSION_VALUATION_BEGIN_DATE", "RETROCESSION_VALUATION_EXPIRED_DATE",
		"SUM_INSURED", "CEDING_RETENTION", "SUM_REASURED", "SHARE_NUSANTARA_RE",
		"GROSS_PREMIUM", "NET_PREMIUM",
		"ENTRY_AGE",
		// QR: PERIOD_MM wajib di judul (Calculate CSV, 05-10-2026).
		"PERIOD_MM",
	}
	return append(k, tambahan...)
}

// csvUji menyusun satu berkas CSV buatan.
//
// ⛔ NOL nama orang, nol nomor polis nyata - seluruh nilainya berawalan
// `UJI-`. Fixture yang menyerupai data nyata adalah data nyata yang kebetulan
// belum ketahuan.
func csvUji(judul []string, baris ...[]string) string {
	var b strings.Builder
	b.WriteString(strings.Join(judul, ","))
	b.WriteString("\n")
	for _, r := range baris {
		b.WriteString(strings.Join(r, ","))
		b.WriteString("\n")
	}
	return b.String()
}

// barisUji menyusun satu baris yang seluruh kolomnya sah.
func barisUji(sertifikat string) []string {
	return []string{
		sertifikat, "UJI PESERTA A", "UJI-PLAN", "IDR", "UJI-POL-1",
		"NM", "01/02/1990", "01/01/2026", "31/12/2026", "01/01/2026",
		"01/01/2026", "01/01/2026", "01/01/2026",
		"01/01/2026", "31/12/2026", "01/01/2026", "31/12/2026",
		"1000.50", "100.25", "900.25", "50.5", "200.75", "180.5",
		"30",
		"12",
	}
}

func TestBacaCSVUnggahMemberiNomorBarisData(t *testing.T) {
	isi := csvUji(judulLengkap(), barisUji("UJI-C1"), barisUji("UJI-C2"))
	baris, err := bacaQR(strings.NewReader(isi))
	if err != nil {
		t.Fatal(err)
	}
	if len(baris) != 2 {
		t.Fatalf("%d baris, mau 2", len(baris))
	}
	// ⛔ Nomor BARIS DATA - 1 untuk baris pertama di bawah judul. Nomor baris
	// berkas mentah akan meleset satu, dan meleset satu di berkas seribu
	// baris lebih buruk daripada tidak ada nomor sama sekali.
	if baris[0].Nomor != 1 || baris[1].Nomor != 2 {
		t.Errorf("nomor baris %d dan %d, mau 1 dan 2", baris[0].Nomor, baris[1].Nomor)
	}
	if baris[0].Nilai["CERTIFICATE_NO"] != "UJI-C1" {
		t.Errorf("nilai kolom tidak terpetakan: %v", baris[0].Nilai)
	}
}

// TestJudulDinaikkanHurufBesar - berkas dari mesin berbeda berbeda hurufnya.
func TestBacaCSVUnggahMenaikkanHurufJudul(t *testing.T) {
	j := judulLengkap()
	j[0] = strings.ToLower(j[0])
	isi := csvUji(j, barisUji("UJI-C1"))
	baris, err := bacaQR(strings.NewReader(isi))
	if err != nil {
		t.Fatalf("judul huruf kecil ditolak: %v", err)
	}
	if baris[0].Nilai["CERTIFICATE_NO"] != "UJI-C1" {
		t.Error("judul huruf kecil tidak dinaikkan menjadi huruf besar")
	}
}

// TestBOMExcelDibuang - jebakan yang menolak SETIAP berkas yang benar.
//
// ⛔ Excel menempelkan BOM UTF-8 di depan judul PERTAMA. Tanpa membuangnya,
// kolom pertama dianggap HILANG - dan berkas yang jelas-jelas memuatnya
// ditolak seluruhnya, dengan pesan yang menuduh kolom yang terlihat ada.
func TestBOMExcelDibuang(t *testing.T) {
	isi := string(rune(0xFEFF)) + csvUji(judulLengkap(), barisUji("UJI-C1"))
	baris, err := bacaQR(strings.NewReader(isi))
	if err != nil {
		t.Fatalf("berkas ber-BOM ditolak: %v", err)
	}
	if baris[0].Nilai["CERTIFICATE_NO"] != "UJI-C1" {
		t.Error("BOM membuat kolom pertama tidak terbaca")
	}
}

func TestBacaCSVUnggahMenolakYangTidakTerbaca(t *testing.T) {
	for _, k := range []struct {
		apa   string
		isi   string
		mauIs error
	}{
		{"berkas kosong", "", ErrCSVKosong},
		{"hanya judul", csvUji(judulLengkap()), ErrCSVKosong},
		{"judul ganda", csvUji(append(judulLengkap(), "DOB"), barisUji("UJI-C1")), ErrCSVJudulGanda},
		{"kolom wajib hilang", "CERTIFICATE_NO\nUJI-C1\n", ErrCSVKolomKurang},
	} {
		_, err := bacaQR(strings.NewReader(k.isi))
		if !errors.Is(err, k.mauIs) {
			t.Errorf("%s: galat %v, mau %v", k.apa, err, k.mauIs)
		}
	}
}

// TestKolomHilangDilaporkanSekaliDanMenyebutNamanya.
//
// ⛔ Kolom yang hilang SAMA SEKALI menghasilkan satu penolakan "kolom kosong"
// untuk setiap baris kalau diperiksa per baris - seribu kalimat yang
// mengatakan satu hal.
func TestKolomHilangDilaporkanSekaliDanMenyebutNamanya(t *testing.T) {
	j := judulLengkap()
	var tanpa []string
	for _, k := range j {
		if k != "SUM_INSURED" {
			tanpa = append(tanpa, k)
		}
	}
	_, err := bacaQR(strings.NewReader(csvUji(tanpa, barisUji("UJI-C1")[:len(tanpa)])))
	if !errors.Is(err, ErrCSVKolomKurang) {
		t.Fatalf("galat %v, mau ErrCSVKolomKurang", err)
	}
	if !strings.Contains(err.Error(), "SUM_INSURED") {
		t.Errorf("pesan %q tidak menyebut kolom yang hilang", err)
	}
}

// TestBarisKosongDiUjungDilewati - penyunting menambahkannya sendiri.
func TestBarisKosongDiUjungDilewati(t *testing.T) {
	isi := csvUji(judulLengkap(), barisUji("UJI-C1")) + "\n\n"
	baris, err := bacaQR(strings.NewReader(isi))
	if err != nil {
		t.Fatalf("baris kosong di ujung menolak berkas: %v", err)
	}
	if len(baris) != 1 {
		t.Errorf("%d baris, mau 1 - baris kosong ikut terbaca", len(baris))
	}
}

// TestTinjauTidakMenyentuhApaPun - AC tiket 04.
//
// ⛔ Penjaga STATIK, dan sengaja: uji perilaku menuntut Oracle, dan uji yang
// menuntut Oracle akan SKIP di mesin yang tidak punya - yaitu diam, yaitu
// hijau. Yang dijaga di sini bentuk kodenya: `Tinjau` tidak boleh membuka
// transaksi sama sekali.
func TestTinjauTidakMenyentuhApaPun(t *testing.T) {
	isi, err := bacaBerkasLayanan("polis_unggah.go")
	if err != nil {
		t.Fatal(err)
	}
	a := strings.Index(isi, "func (u *UnggahPremiumList) Tinjau(")
	b := strings.Index(isi, "func periksaBerkas(")
	if a < 0 || b < 0 || b < a {
		t.Fatal("fungsi Tinjau atau periksaBerkas tidak ditemukan")
	}
	badan := isi[a:b]
	for _, jejak := range []string{"DalamTransaksi", "ExecContext", "INSERT", "DELETE", "UPDATE"} {
		if strings.Contains(badan, jejak) {
			t.Errorf("Tinjau memuat %q - ia harus MEMBACA saja", jejak)
		}
	}
	// Tinjau dan Simpan SATU jalan; Validate CSV TANPA hitung QR - hanya bentuk
	// dan batas usia / sum insured (keputusan work owner 05-10-2026). Jalan itu
	// pun tidak menulis.
	if !strings.Contains(badan, "u.periksaDanHitung(ctx, polisID, berkas, tipe, false)") {
		t.Error("Tinjau tidak memakai periksaDanHitung tanpa hitung QR")
	}
	c := strings.Index(isi, "func (u *UnggahPremiumList) periksaDanHitung(")
	d := strings.Index(isi, "// HasilSimpanUnggah adalah jawaban penyimpanan.")
	if c < 0 || d < c {
		t.Fatal("periksaDanHitung tidak ditemukan")
	}
	for _, jejak := range []string{"DalamTransaksi", "ExecContext", "INSERT", "DELETE", "UPDATE"} {
		if strings.Contains(isi[c:d], jejak) {
			t.Errorf("periksaDanHitung memuat %q - ia harus MEMBACA saja", jejak)
		}
	}
}

// TestSimpanMemvalidasiUlang - tinjauan yang lolos bukan izin menyimpan.
func TestSimpanMemvalidasiUlang(t *testing.T) {
	isi, err := bacaBerkasLayanan("polis_unggah.go")
	if err != nil {
		t.Fatal(err)
	}
	a := strings.Index(isi, "func (u *UnggahPremiumList) Simpan(")
	if a < 0 {
		t.Fatal("fungsi Simpan tidak ditemukan")
	}
	badan := isi[a:]
	// Simpan → periksaDanHitung → periksaBerkas (keputusan work owner 05-10-2026).
	if !strings.Contains(badan, "u.periksaDanHitung(ctx, polisID, berkas, tipe, true)") ||
		!strings.Contains(isi, "baris, hasil, err := periksaBerkas(berkas, tipe, dataPolis.Terapkan)") {
		t.Error("Simpan tidak memvalidasi ulang berkasnya; klien yang dapat " +
			"melewatkan tinjauan dapat menyimpan apa saja")
	}
	if !strings.Contains(badan, "HapusPesertaPolis") ||
		!strings.Contains(badan, "SisipPeserta") {
		t.Error("Simpan tidak mengganti isi (hapus lalu sisip)")
	}
	// ⛔ Keduanya di dalam SATU transaksi: hapus yang berhasil lalu sisip
	// yang gagal meninggalkan polis TANPA peserta sama sekali.
	iTx := strings.Index(badan, "DalamTransaksi(ctx")
	iHapus := strings.Index(badan, "HapusPesertaPolis")
	iSisip := strings.Index(badan, "SisipPeserta")
	if iTx < 0 || iTx > iHapus || iTx > iSisip {
		t.Error("hapus dan sisip tidak keduanya di dalam satu transaksi")
	}
}

// OQ-PL-12 (GILIRAN-17): SATU jalan tinjau/simpan mengisi 0 pada uang kosong
// SEBELUM validasi - urutannya yang membuat "HARUS ADA" tidak berbunyi.
func TestPeriksaBerkasMengisiNolSebelumValidasi(t *testing.T) {
	isi, err := os.ReadFile("polis_unggah.go")
	if err != nil {
		t.Fatal(err)
	}
	s := string(isi)
	i := strings.Index(s, "func periksaBerkas(")
	if i < 0 {
		t.Fatal("periksaBerkas tidak ada")
	}
	badan := s[i:]
	if j := strings.Index(badan[1:], "\nfunc "); j >= 0 {
		badan = badan[:j+1]
	}
	nol, val := strings.Index(badan, "models.IsiNolUangKosong(baris)"), strings.Index(badan, "models.ValidasiUnggah(tipe, baris)")
	if nol < 0 || val < 0 || nol > val {
		t.Errorf("periksaBerkas harus memanggil IsiNolUangKosong SEBELUM ValidasiUnggah (nol=%d, validasi=%d)", nol, val)
	}
}

// TestRingkasPenolakan - 409 Calculate CSV menyebut baris, kolom, dan pesan
// (05-10-2026), dibatasi supaya tidak menjadi seribu kalimat.
func TestRingkasPenolakan(t *testing.T) {
	d := []models.Penolakan{{Baris: 1, Kolom: "RATE", Pesan: "No rate found for age 2, contract 1 in R/I Rate UJI-1"}}
	got := RingkasPenolakan(d)
	if got != "Calculate CSV rejected 1 row(s); nothing was saved. Row 1 RATE: No rate found for age 2, contract 1 in R/I Rate UJI-1." {
		t.Errorf("ringkasan = %q", got)
	}
	var banyak []models.Penolakan
	for i := 1; i <= 12; i++ {
		banyak = append(banyak, models.Penolakan{Baris: i, Kolom: "RATE", Pesan: "UJI"})
	}
	if got := RingkasPenolakan(banyak); !strings.HasSuffix(got, " ... and 2 more.") {
		t.Errorf("ringkasan tidak dibatasi: %q", got)
	}
	if RingkasPenolakan(nil) != "" {
		t.Error("tanpa penolakan tetap berkalimat")
	}
}

// TestPolicyHolderDanPlanDariPolis - POLICY_HOLDER / PLAN peserta ditimpa data
// polis SEBELUM validasi per baris, di jalan bersama Validate / Calculate CSV
// (keputusan work owner 05-10-2026).
func TestPolicyHolderDanPlanDariPolis(t *testing.T) {
	isi, err := bacaBerkasLayanan("polis_unggah.go")
	if err != nil {
		t.Fatal(err)
	}
	jalan := isi[strings.Index(isi, "func (u *UnggahPremiumList) periksaDanHitung("):]
	iPeriksa := strings.Index(jalan, "dataPolis.Periksa()")
	iBerkas := strings.Index(jalan, "periksaBerkas(berkas, tipe, dataPolis.Terapkan)")
	if iPeriksa < 0 || iBerkas < iPeriksa {
		t.Error("data polis tidak diperiksa sebelum berkas diproses")
	}
	b := isi[strings.Index(isi, "func periksaBerkas("):]
	nol, kait, val := strings.Index(b, "models.IsiNolUangKosong(baris)"), strings.Index(b, "f(baris)"),
		strings.Index(b, "models.ValidasiUnggah(tipe, baris)")
	if !(nol < kait && kait < val) {
		t.Error("kait data polis tidak berjalan sebelum validasi per baris")
	}
}

// TestSpreadingDiJalanCalculate - keputusan work owner 06-10-2026:
//   - spreading hanya Type QR, sesudah hitung premi (cabang hitung QR);
//   - Validate CSV (hitung=false) tidak menjalankannya dan tidak menulis;
//   - Calculate CSV menulis peserta + spreading di SATU transaksi, spreading lama dihapus eksplisit.
func TestSpreadingDiJalanCalculate(t *testing.T) {
	isi, err := bacaBerkasLayanan("polis_unggah.go")
	if err != nil {
		t.Fatal(err)
	}
	jalan := isi[strings.Index(isi, "func (u *UnggahPremiumList) periksaDanHitung("):]
	jalan = jalan[:strings.Index(jalan, "\n}\n")]
	iQR := strings.Index(jalan, "if hitung && models.HitungPesertaPerTipe(tipe) {")
	iHitung := strings.Index(jalan, "models.HitungPesertaQR(m, baris, lewati)")
	iSpread := strings.Index(jalan, "u.spreadingQR(ctx, rujukan, kepala, baris, lewati)")
	if iQR < 0 || iHitung < iQR || iSpread < iHitung {
		t.Error("spreading tidak di cabang hitung QR sesudah hitung premi")
	}
	simpan := isi[strings.Index(isi, "func (u *UnggahPremiumList) Simpan("):]
	iTx := strings.Index(simpan, "DalamTransaksi(ctx")
	iHapus := strings.Index(simpan, "spreading.HapusSpreadingPolis(ctx, tx, polisID)")
	iSisipP := strings.Index(simpan, "peserta.SisipPeserta(ctx, tx, polisID, baris)")
	iSisipS := strings.Index(simpan, "spreading.SisipSpreading(ctx, tx, polisID, pelaku.AkunID, idPeserta, baris)")
	if iTx < 0 || iHapus < iTx || iSisipP < iHapus || iSisipS < iSisipP {
		t.Error("peserta dan spreading tidak ditulis di satu transaksi (hapus spreading → sisip peserta → sisip spreading)")
	}
	tinjau := isi[strings.Index(isi, "func (u *UnggahPremiumList) Tinjau("):strings.Index(isi, "func periksaBerkas(")]
	if strings.Contains(tinjau, "SisipSpreading") || strings.Contains(tinjau, "DalamTransaksi") {
		t.Error("Validate CSV menulis spreading")
	}
}

// TestConfirmDigerbangiKesiapan - keputusan work owner 06-10-2026: Confirm menjalankan
// pemeriksaan kesiapan SEBELUM transaksi yang menerbitkan nomor dan menutup kasus.
func TestConfirmDigerbangiKesiapan(t *testing.T) {
	isi, err := bacaBerkasLayanan("polis_penawaran.go")
	if err != nil {
		t.Fatal(err)
	}
	badan := isi[strings.Index(isi, "func (p *Penawaran) PutuskanDenganNomor("):]
	iSiap := strings.Index(badan, "p.periksaSiapConfirm(ctx, pelaku, keadaan)")
	iTerapkan := strings.Index(badan, "p.terapkan(ctx, pelaku, keadaan, akibat, sebab, saat)")
	if iSiap < 0 || iTerapkan < iSiap {
		t.Error("Confirm tidak memeriksa kesiapan sebelum menerbitkan nomor")
	}
	k, err := bacaBerkasLayanan("polis_kesiapan.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, jejak := range []string{"models.SusunDataPolis(d.IsianDataPolis)", "PesanSiapWPC", "PesanSiapPeserta",
		"PesanSiapSummary", "PesanSiapSpreading", "PesanSiapPeriode"} {
		if !strings.Contains(k, jejak) {
			t.Errorf("kesiapan tanpa %s", jejak)
		}
	}
	for _, jejak := range []string{"DalamTransaksi", "ExecContext"} {
		if strings.Contains(k, jejak) {
			t.Errorf("kesiapan menulis (%s)", jejak)
		}
	}
}
