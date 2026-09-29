package models

// Bentuk `PL_NUMBER` - tiket 03 PremiumList Life. TANPA Oracle.

import (
	"errors"
	"strings"
	"testing"

	"nusantarare/inti/penomor"
)

func TestKodeTipePLEmpatCabangDariData(t *testing.T) {
	for _, tipe := range []string{"QR", "QP", "TP", "TR"} {
		kode, err := KodeTipePL(tipe)
		if err != nil {
			t.Errorf("tipe %q ditolak: %v", tipe, err)
			continue
		}
		if kode != tipe {
			t.Errorf("tipe %q -> %q, mau %q", tipe, kode, tipe)
		}
	}
	// Spasi di kedua ujung tidak membuat tipe yang sah menjadi tidak sah.
	if _, err := KodeTipePL("  TR  "); err != nil {
		t.Errorf("tipe berspasi ditolak: %v", err)
	}
}

// TestKodeTipePLMenolakTipeTanpaCabang - penyimpangan sadar dari Pega.
//
// ⛔ Di Pega `InputData.CARI20` tidak direset langkah 5, jadi tipe di luar
// keempatnya menghasilkan nomor tanpa huruf tipe - atau dengan huruf tipe
// polis SEBELUMNYA yang masih tertinggal di halaman. Keduanya nomor yang
// terlihat sah, dan keduanya salah.
func TestKodeTipePLMenolakTipeTanpaCabang(t *testing.T) {
	for _, tipe := range []string{"", "  ", "Q", "qr", "TX", "QRQP", "FAC"} {
		if _, err := KodeTipePL(tipe); !errors.Is(err, penomor.ErrTipePLTanpaCabang) {
			t.Errorf("tipe %q diterima; mau ditolak (%v)", tipe, err)
		}
	}
}

// TestSatuPenghitungUntukEmpatTipe mengunci b2717.
//
// ⛔ INI YANG PALING MUDAH DIRUSAK "PERBAIKAN". Kunci penghitungnya memakai
// teks harfiah `QR/QP/TP/TR`, BUKAN tipe yang sedang berjalan. Memecahnya
// menjadi empat penghitung menerbitkan empat deret yang masing-masing mulai
// dari 1, dan setiap nomor baru bertabrakan dengan nomor lama.
func TestSatuPenghitungUntukEmpatTipe(t *testing.T) {
	const awalan = "RNML-"
	jenis := JenisPenghitungPL(awalan)
	for _, tipe := range []string{"QR", "QP", "TP", "TR"} {
		// Apa pun tipenya, kuncinya TIDAK berubah.
		if lain := JenisPenghitungPL(awalan); lain != jenis {
			t.Fatalf("kunci penghitung berubah: %q vs %q", lain, jenis)
		}
		if strings.Contains(strings.TrimPrefix(jenis, awalan+"QR/QP/TP/TR"), tipe) {
			t.Errorf("kunci penghitung memuat tipe %q di luar teks harfiahnya", tipe)
		}
	}
	if jenis != awalan+"QR/QP/TP/TR" {
		t.Errorf("kunci penghitung = %q, mau %q", jenis, awalan+"QR/QP/TP/TR")
	}
}

// TestClassPenghitungPLTidakBergeserDiam kunci nilai yang masih `[terbuka]`.
//
// ⛔ Bila baris `GENERATE_SEQUENCE_NUMBER` di basis data nyata bertuliskan
// kelas yang LAIN, penomoran akan mulai dari satu dan setiap nomor baru
// bertabrakan dengan nomor lama. Nilainya karena itu dikunci di sini:
// menggantinya menuntut alasan, bukan sekadar satu suntingan.
func TestClassPenghitungPLTidakBergeserDiam(t *testing.T) {
	const mau = "ASM-FW-GISFW-Work-LIFE"
	if ClassPenghitungPL != mau {
		t.Errorf("ClassPenghitungPL = %q, mau %q.\n"+
			"Nilainya [terbuka - pemilik kerja / DBA]: kalau memang berubah, "+
			"ubah bersama alasannya, jangan hanya konstantanya.",
			ClassPenghitungPL, mau)
	}
}

func TestPeriodeNomorPLMemotongTahunJadiDuaAngka(t *testing.T) {
	for _, k := range []struct{ masuk, mau string }{
		{"09.2026", "09.26"},
		{"01.2000", "01.00"},
		{"12.1999", "12.99"},
		{" 03.2026 ", "03.26"},
	} {
		got, err := penomor.PeriodeNomorPL(k.masuk)
		if err != nil {
			t.Errorf("%q ditolak: %v", k.masuk, err)
			continue
		}
		if got != k.mau {
			t.Errorf("%q -> %q, mau %q", k.masuk, got, k.mau)
		}
	}
}

// TestPeriodeNomorPLMenolakBentukLain - `@substring` yang membabi buta.
func TestPeriodeNomorPLMenolakBentukLain(t *testing.T) {
	for _, buruk := range []string{
		"", "09.26", "9.2026", "092026", "09-2026", "09.20267", "ab.2026", "09.20a6",
	} {
		if _, err := penomor.PeriodeNomorPL(buruk); !errors.Is(err, penomor.ErrPeriodeNomorPLTakBerbentuk) {
			t.Errorf("periode %q diterima; mau ditolak", buruk)
		}
	}
}

func TestRakitNomorPLBerbentukSepertiB3126(t *testing.T) {
	got := penomor.RakitNomorPL("RNML-", "QR", "LF", "09.26", 7)
	const mau = "RNML-QRLF.09.26.00007"
	if got != mau {
		t.Errorf("nomor = %q, mau %q", got, mau)
	}
}

// TestUrutNomorPLLimaAngkaDanTidakDipotong - `LPAD(v_seq,5,'0')`.
//
// ⚠️ Urut yang sudah lebih panjang dari lima TIDAK dipotong. Memotongnya
// menerbitkan nomor yang bertabrakan - persis yang `LPAD` sendiri tidak
// lakukan.
func TestUrutNomorPLLimaAngkaDanTidakDipotong(t *testing.T) {
	for _, k := range []struct {
		urut int
		mau  string
	}{
		{1, "00001"},
		{99, "00099"},
		{99999, "99999"},
		{100000, "100000"},
		{1234567, "1234567"},
	} {
		got := penomor.RakitNomorPL("P", "QR", "B", "09.26", k.urut)
		bagian := strings.Split(got, ".")
		akhir := bagian[len(bagian)-1]
		if akhir != k.mau {
			t.Errorf("urut %d -> %q, mau %q", k.urut, akhir, k.mau)
		}
	}
}

// TestTahunNomorPLDuaAngka menjaga perbedaan dari nomor klaim.
//
// ⛔ Procedure yang SAMA mengembalikan `MM.YYYY`; nomor klaim memakainya apa
// adanya, nomor PL memotongnya. Menyeragamkan keduanya membuat salah satunya
// tidak lagi cocok dengan nomor yang sudah beredar.
func TestTahunNomorPLDuaAngka(t *testing.T) {
	nomor, err := NomorPL(BahanNomorPL{
		Awalan: "RNML-", Tipe: "TR", KodeBisnis: "LF",
		PeriodeMMYYYY: "09.2026", Urut: 12,
	})
	if err != nil {
		t.Fatal(err)
	}
	if nomor != "RNML-TRLF.09.26.00012" {
		t.Fatalf("nomor = %q", nomor)
	}
	bagian := strings.Split(nomor, ".")
	if len(bagian) != 4 {
		t.Fatalf("nomor %q tidak berbentuk empat bagian ber-titik", nomor)
	}
	if len(bagian[2]) != 2 {
		t.Errorf("bagian tahun %q berpanjang %d, mau 2 - bukan empat seperti nomor klaim",
			bagian[2], len(bagian[2]))
	}
	if strings.Contains(nomor, "2026") {
		t.Errorf("nomor %q memuat tahun empat angka", nomor)
	}
}

// TestNomorPLMenolakBahanKosong - nol bahan kosong yang dibiarkan lewat.
func TestNomorPLMenolakBahanKosong(t *testing.T) {
	utuh := BahanNomorPL{
		Awalan: "RNML-", Tipe: "QR", KodeBisnis: "LF",
		PeriodeMMYYYY: "09.2026", Urut: 1,
	}
	if _, err := NomorPL(utuh); err != nil {
		t.Fatalf("bahan utuh ditolak: %v", err)
	}
	for _, k := range []struct {
		apa   string
		ubah  func(b BahanNomorPL) BahanNomorPL
		mauIs error
	}{
		{"awalan kosong", func(b BahanNomorPL) BahanNomorPL {
			b.Awalan = "  "
			return b
		}, penomor.ErrAwalanProduksiKosong},
		{"tipe kosong", func(b BahanNomorPL) BahanNomorPL {
			b.Tipe = ""
			return b
		}, penomor.ErrTipePLTanpaCabang},
		{"kode bisnis kosong", func(b BahanNomorPL) BahanNomorPL {
			b.KodeBisnis = " "
			return b
		}, penomor.ErrKodeBisnisKosong},
		{"periode tak berbentuk", func(b BahanNomorPL) BahanNomorPL {
			b.PeriodeMMYYYY = "09.26"
			return b
		}, penomor.ErrPeriodeNomorPLTakBerbentuk},
	} {
		if _, err := NomorPL(k.ubah(utuh)); !errors.Is(err, k.mauIs) {
			t.Errorf("%s: galat %v, mau %v", k.apa, err, k.mauIs)
		}
	}
	// Urut nol atau negatif berarti penghitungnya tidak dinaikkan.
	for _, urut := range []int{0, -1} {
		b := utuh
		b.Urut = urut
		if _, err := NomorPL(b); err == nil {
			t.Errorf("urut %d diterima; mau ditolak", urut)
		}
	}
}
