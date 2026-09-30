package skemauji

// Penjaga skema uji - TANPA Oracle.
//
// Untuk apa berkas ini: test bertag `db` MENGHAPUS tabel di skema yang
// ditunjuk ORACLE_SCHEMA - termasuk OS_AKSEPTASI_KLAIM_LIFE, tabel datar
// warisan. Sampai ronde 4 satu-satunya penjaganya adalah IS_PEGA_PROD, dan
// Makefile bahkan memberi bawaan ORACLE_SCHEMA = POOLDATA. Menjalankan
// `go test -tags=db` di instance pengembangan yang memuat tabel warisan
// sungguhan karena itu cukup untuk menghapusnya.
//
// Test di sini menjalankan keputusan pagar itu sendiri, tanpa env dan tanpa
// Oracle, sehingga ia ikut lari di setiap `go test ./...`.
//
// Dibaca sesudah: skemauji.go.

import (
	"errors"
	"os"
	"strings"
	"testing"
)

// Pagar menuntut DUA syarat, dan menolak bila salah satunya tidak terpenuhi.
func TestPagarSkemaUjiMenuntutDuaSyarat(t *testing.T) {
	kasus := []struct {
		nama   string
		skema  string
		diakui string
		mauSah bool
	}{
		{"tanpa env sama sekali", "SKEMA_UJI_DBA", "", false},
		{"env salah eja", "SKEMA_UJI_DBA", "ya", false},
		{"env false", "SKEMA_UJI_DBA", "false", false},
		{"skema warisan meski env true", namaSkemaWarisan, "true", false},
		{"skema warisan huruf kecil", "pooldata", "true", false},
		{"skema warisan berspasi", "  POOLDATA  ", "true", false},
		{"skema warisan tanpa env", namaSkemaWarisan, "", false},
		{"turunan skema warisan", "POOLDATA_DEV", "true", false},
		{"turunan berangka", "pooldata2", "true", false},
		{"keduanya benar", "SKEMA_UJI_DBA", "true", true},
		{"keduanya benar, env huruf besar", "SKEMA_UJI_DBA", "TRUE", true},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			err := pagarSkemaUji(k.skema, k.diakui)
			if k.mauSah {
				if err != nil {
					t.Errorf("pagar menolak yang seharusnya sah: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("pagar meloloskan skema yang seharusnya ditolak")
			}
			if !errors.Is(err, ErrBukanSkemaUji) {
				t.Errorf("galatnya bukan ErrBukanSkemaUji: %v", err)
			}
		})
	}
}

// Pesan galatnya harus menyebut SEBABNYA, bukan sekadar menolak.
//
// Orang yang membacanya sedang salah menyetel env dan tidak tahu bahwa test
// ini menghapus tabel. Pesan yang hanya berbunyi "ditolak" membuat ia menyetel
// ORACLE_SKEMA_UJI=true supaya lewat - persis yang tidak boleh terjadi.
func TestPesanPagarMenyebutSebabnya(t *testing.T) {
	tanpaEnv := pagarSkemaUji("SKEMA_UJI_DBA", "").Error()
	// namaTabelLama tidak lagi disebut: pagarnya kini tinggal di inti/backend/config,
	// yang dipakai juga oleh flag -migrate-down dan tidak tahu-menahu soal tabel
	// warisan. Yang wajib tetap ada adalah nama env-nya dan kata MENGHAPUS.
	for _, mau := range []string{envSkemaUji, "MENGHAPUS"} {
		if !strings.Contains(tanpaEnv, mau) {
			t.Errorf("pesan tanpa env tidak menyebut %q: %s", mau, tanpaEnv)
		}
	}
	warisan := pagarSkemaUji(namaSkemaWarisan, "true").Error()
	for _, mau := range []string{namaSkemaWarisan, "warisan"} {
		if !strings.Contains(warisan, mau) {
			t.Errorf("pesan skema warisan tidak menyebut %q: %s", mau, warisan)
		}
	}
}

// periksaPagarSkemaUji membaca env yang benar - bukan nama lain yang mirip.
func TestPagarMembacaEnvYangBenar(t *testing.T) {
	t.Setenv(envSkemaUji, "")
	if err := periksaPagarSkemaUji("SKEMA_UJI_DBA"); err == nil {
		t.Error("env kosong seharusnya ditolak")
	}
	t.Setenv(envSkemaUji, "true")
	if err := periksaPagarSkemaUji("SKEMA_UJI_DBA"); err != nil {
		t.Errorf("env true seharusnya lolos: %v", err)
	}
	if os.Getenv(envSkemaUji) != "true" {
		t.Fatal("t.Setenv tidak berlaku; testnya yang rusak, bukan kodenya")
	}
	if err := periksaPagarSkemaUji(namaSkemaWarisan); err == nil {
		t.Error("POOLDATA seharusnya ditolak meski env true")
	}
}

// Hanya SATU galat yang layak dilewati: Oracle memang belum dikonfigurasi.
//
// Ini yang membedakan "tidak ada Oracle" - keadaan wajar di jalur B - dari
// "Oracle ada, tetapi salah yang ditunjuk". Yang pertama SKIP, yang kedua
// harus berteriak.
func TestHanyaTanpaOracleYangBolehDilewati(t *testing.T) {
	if !BolehDilewati(ErrTanpaOracle) {
		t.Error("ErrTanpaOracle seharusnya boleh dilewati")
	}
	for _, err := range []error{ErrProduksi, ErrBukanSkemaUji,
		pagarSkemaUji(namaSkemaWarisan, "true"), errors.New("galat lain")} {
		if BolehDilewati(err) {
			t.Errorf("galat ini seharusnya TIDAK dilewati: %v", err)
		}
	}
}
