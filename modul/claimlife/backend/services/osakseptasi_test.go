package services

// RALAT 10-10-2026 - baris datar OS_AKSEPTASI_KLAIM_LIFE: Register tidak menyentuhnya; Save to RNM satu-satunya yang
// melahirkannya, di transaksi utama sesudah transaksi penomoran. TANPA Oracle (penjaga statik).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegisterTidakMenyentuhOSAkseptasi(t *testing.T) {
	isi, err := os.ReadFile("pendaftaran.go")
	if err != nil {
		t.Fatal(err)
	}
	kode := buangKomentarBaris(string(isi))
	for _, larang := range []string{"OS_AKSEPTASI", "TulisCerminSaveRNM", "BarisLamaDari(", "IsiTertanggungCermin(",
		"SetelCerminOutstanding("} {
		if strings.Contains(kode, larang) {
			t.Errorf("pendaftaran.go (Register) menyentuh baris datar lewat %q", larang)
		}
	}
}

func TestHanyaSaveRNMYangMelahirkanBarisDatar(t *testing.T) {
	berkas, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	pemanggil := []string{}
	for _, b := range berkas {
		if strings.HasSuffix(b, "_test.go") {
			continue
		}
		isi, err := os.ReadFile(b)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(buangKomentarBaris(string(isi)), ".TulisCerminSaveRNM(") {
			pemanggil = append(pemanggil, b)
		}
	}
	if len(pemanggil) != 1 || pemanggil[0] != "simpanrnm.go" {
		t.Errorf("pemanggil TulisCerminSaveRNM = %v, mau [simpanrnm.go]", pemanggil)
	}
}

func TestSaveRNMMenulisBarisDatarSesudahPenomoranDiTransaksiUtama(t *testing.T) {
	isi, err := os.ReadFile("simpanrnm.go")
	if err != nil {
		t.Fatal(err)
	}
	kode := buangKomentarBaris(string(isi))
	nomor := strings.Index(kode, "NomorBerikut(")
	isiNomor := strings.Index(kode, "IsiNomorKlaimKosong(")
	susun := strings.Index(kode, "klaimDatar.NomorKlaim = hasil.NomorKlaim")
	tulis := strings.Index(kode, ".TulisCerminSaveRNM(ctx, tx,")
	setel := strings.Index(kode, ".SetelCerminOutstanding(ctx, tx,")
	if nomor < 0 || isiNomor < 0 || susun < 0 || tulis < 0 || setel < 0 {
		t.Fatalf("bentuk Save to RNM tak terbaca: %d %d %d %d %d", nomor, isiNomor, susun, tulis, setel)
	}
	// Baris datar disusun dari nomor yang SUDAH lahir, ditulis di transaksi utama, lalu statusnya disetel '0'.
	if !(nomor < isiNomor && isiNomor < susun && susun < tulis && tulis < setel) {
		t.Errorf("urutan penomoran -> susun baris datar -> tulis -> setel '0' berubah: %d %d %d %d %d",
			nomor, isiNomor, susun, tulis, setel)
	}
	utama := strings.LastIndex(kode[:tulis], "x.svc.DalamTransaksi(")
	if utama < 0 || utama < isiNomor {
		t.Error("baris datar tidak ditulis di transaksi utama Save to RNM (sesudah transaksi penomoran)")
	}
	// INSERT hanya untuk baris yang belum pernah disimpan.
	if !strings.Contains(kode, "repository.CerminBolehDisisip(b.KodeStatus)") {
		t.Error("Save to RNM tidak membatasi kelahiran baris datar pada baris yang belum pernah disimpan")
	}
}

// buangKomentarBaris membuang baris komentar penuh supaya kata di komentar tidak meluluskan/menggagalkan penjaga.
func buangKomentarBaris(isi string) string {
	var b strings.Builder
	for _, baris := range strings.Split(isi, "\n") {
		if strings.HasPrefix(strings.TrimSpace(baris), "//") {
			continue
		}
		b.WriteString(baris)
		b.WriteString("\n")
	}
	return b.String()
}
