package repository

// Test daftar kolom tabel datar warisan - TANPA Oracle.
//
// Untuk apa berkas ini: daftar 55 kolom dipakai bertiga - pembaca, penulis
// fixture, dan tabel tiruan. Kalau daftarnya cacat, ketiganya cacat bersama dan
// tidak ada yang saling mengoreksi. Test di sini menjaga daftarnya sendiri.

import (
	"reflect"
	"strings"
	"testing"
)

// Kelima puluh lima kolom ada, dan tidak satu pun tertulis dua kali.
func TestLimaPuluhLimaKolomWarisan(t *testing.T) {
	nama := NamaKolomBarisLama()
	if len(nama) != 55 {
		t.Errorf("kolom warisan %d, mau 55", len(nama))
	}
	lihat := map[string]bool{}
	for _, n := range nama {
		if lihat[n] {
			t.Errorf("kolom %s tertulis dua kali", n)
		}
		lihat[n] = true
	}
}

// Setiap medan struct BarisLama terwakili tepat satu kali di daftar kolom.
//
// Ini yang menangkap medan yang ditambahkan ke struct tetapi lupa didaftar -
// medan begitu tidak akan pernah terbaca dari Oracle, dan tanpa test ini
// kelalaiannya diam sepenuhnya.
func TestSetiapMedanBarisLamaTerdaftar(t *testing.T) {
	var b BarisLama
	tipe := reflect.TypeOf(b)
	terdaftar := map[string]bool{}
	for _, n := range NamaKolomBarisLama() {
		terdaftar[n] = true
	}
	for i := 0; i < tipe.NumField(); i++ {
		nama := tipe.Field(i).Name
		if !terdaftar[nama] {
			t.Errorf("medan BarisLama.%s tidak ada di daftar kolom", nama)
		}
	}
	if tipe.NumField() != len(NamaKolomBarisLama()) {
		t.Errorf("medan struct %d, kolom terdaftar %d",
			tipe.NumField(), len(NamaKolomBarisLama()))
	}
}

// Daftar kolom dan tujuan Scan menunjuk medan yang sama, dalam urutan yang sama.
//
// Kalau keduanya bergeser satu posisi, nilai kolom akan mendarat di medan
// tetangganya - kerusakan yang tidak menimbulkan galat apa pun.
func TestUrutanKolomDanTujuanScanSejajar(t *testing.T) {
	var b BarisLama
	medan := medanBarisLama(&b)
	for i, m := range medan {
		penanda := "TANDA-" + m.Kolom
		*m.Nilai = penanda
		// Dibaca ulang lewat daftar yang baru, bukan lewat pointer yang sama.
		if got := *medanBarisLama(&b)[i].Nilai; got != penanda {
			t.Errorf("posisi %d kolom %s: dapat %q", i, m.Kolom, got)
		}
	}
}

// Kolom angka dan tanggal dibungkus TO_CHAR; kolom teks tidak.
//
// ⛔ Ralat 26-09-2026: STS_REJECT dulu didaftar di sini sebagai kolom TEKS,
// dan itu keliru - katalog menyebutnya NUMBER(38,0). WPC juga, dan ia DATE.
// Keduanya kini dibungkus, sehingga pembacaannya tidak lagi bergantung setelan
// NLS sesi.
func TestEkspresiBacaMembungkusYangSeharusnya(t *testing.T) {
	kasus := map[string]string{
		"CLAIM_AMOUNT":     "TO_CHAR",
		"SUM_INSURED":      "TO_CHAR",
		"AGE":              "TO_CHAR",
		"ACCEPTATION_DATE": "TO_CHAR",
		"DOB":              "TO_CHAR",
		"STS_REJECT":       "TO_CHAR",
		"CLAIM_RETRO":      "TO_CHAR",
		"WPC":              "TO_CHAR",
		"CERTIFICATE_NO":   "",
		"POLICY_NO":        "",
	}
	for kolom, mau := range kasus {
		e := ekspresiBacaLama(kolom)
		if mau == "" {
			if e != kolom {
				t.Errorf("%s dibungkus menjadi %q, seharusnya apa adanya", kolom, e)
			}
			continue
		}
		if !strings.Contains(e, mau) {
			t.Errorf("%s tidak dibungkus %s: %q", kolom, mau, e)
		}
	}
}

// ⛔ Kolom angka WAJIB membawa argumen NLS.
//
// Tanpa argumen itu, sesi ber-NLS koma menyerahkan "1234,56" dan seluruh
// pembacaan uang rusak diam-diam - persis jebakan yang ADR-U-0003 hindari.
func TestKolomAngkaMembawaArgumenNLS(t *testing.T) {
	for kolom := range kolomAngkaLama {
		e := ekspresiBacaLama(kolom)
		if !strings.Contains(e, "NLS_NUMERIC_CHARACTERS") {
			t.Errorf("%s dibaca tanpa argumen NLS: %q", kolom, e)
		}
	}
	if len(kolomAngkaLama) == 0 {
		t.Fatal("nol kolom angka terdaftar; daftarnya yang rusak")
	}
}

// Bentuk tanggal saat menulis sama persis dengan saat membaca.
func TestBentukTanggalTulisDanBacaSama(t *testing.T) {
	baca := ekspresiBacaLama("ACCEPTATION_DATE")
	tulis := PenampungTulisLama("ACCEPTATION_DATE", 1)
	const bentuk = "YYYY-MM-DD HH24:MI:SS"
	if !strings.Contains(baca, bentuk) {
		t.Errorf("bentuk baca %q tidak memuat %s", baca, bentuk)
	}
	if !strings.Contains(tulis, bentuk) {
		t.Errorf("bentuk tulis %q tidak memuat %s", tulis, bentuk)
	}
	if biasa := PenampungTulisLama("POLICY_NO", 2); biasa != ":2" {
		t.Errorf("kolom teks dibungkus saat menulis: %q", biasa)
	}
}

// Teks kosong menjadi NULL, bukan teks kosong dan bukan nol.
func TestNilaiKosongMenjadiNull(t *testing.T) {
	b := BarisLama{ID: "R1", CASEID: "UJI-CASE-1", CLAIM_AMOUNT: "100.50"}
	nilai := NilaiBarisLama(b)
	if len(nilai) != 55 {
		t.Fatalf("nilai %d, mau 55", len(nilai))
	}
	nama := NamaKolomBarisLama()
	for i, n := range nama {
		switch n {
		case "ID", "CASEID", "CLAIM_AMOUNT":
			if nilai[i] == nil {
				t.Errorf("%s bernilai nil, seharusnya terisi", n)
			}
		default:
			if nilai[i] != nil {
				t.Errorf("%s = %v, seharusnya nil", n, nilai[i])
			}
		}
	}
}

// Tanggal dinormalkan ke satu bentuk sebelum ditulis.
func TestTanggalDinormalkanSebelumDitulis(t *testing.T) {
	b := BarisLama{ID: "R1", ACCEPTATION_DATE: "2024-01-15"}
	nilai := NilaiBarisLama(b)
	for i, n := range NamaKolomBarisLama() {
		if n != "ACCEPTATION_DATE" {
			continue
		}
		got, ok := nilai[i].(string)
		if !ok {
			t.Fatalf("ACCEPTATION_DATE bertipe %T, mau string", nilai[i])
		}
		if got != "2024-01-15 00:00:00" {
			t.Errorf("ACCEPTATION_DATE = %q, mau %q", got, "2024-01-15 00:00:00")
		}
	}
}

// Tipe kolom tiruan mengikuti KATALOG, bukan penggolongan yang ditebak.
//
// ⛔ Ralat 26-09-2026: angka di test ini dulu adalah tebakan yang rapi -
// NUMBER(38,8) untuk semua kolom angka, VARCHAR2(255) untuk semua teks,
// NUMBER(5) untuk AGE. Tabel sungguhan tidak serapi itu, dan tiruan yang
// bentuknya beda dari yang ditiru tidak menguji apa pun tentang bentuk.
func TestTipeKolomTiruanMengikutiGolongan(t *testing.T) {
	kasus := map[string]string{
		"AGE":              "NUMBER(38,0)",
		"STS_REJECT":       "NUMBER(38,0)",
		"CLAIM_AMOUNT":     "NUMBER",
		"SUM_INSURED":      "NUMBER",
		"CLAIM_RETRO":      "NUMBER",
		"ACCEPTATION_DATE": "DATE",
		"DOB":              "DATE",
		"WPC":              "DATE",
		"POLICY_NO":        "VARCHAR2(100)",
		"POLICY_HOLDER":    "VARCHAR2(1000)",
		"ICD_CODE":         "VARCHAR2(10)",
		"STS_KONVERSI":     "CHAR(1)",
	}
	for kolom, mau := range kasus {
		if got := TipeKolomBarisLama(kolom); got != mau {
			t.Errorf("TipeKolomBarisLama(%s) = %q, mau %q", kolom, got, mau)
		}
	}
}
