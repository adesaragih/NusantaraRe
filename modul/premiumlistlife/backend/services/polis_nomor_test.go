package services

// Penjaga statik penerbitan `PL_NUMBER` - tiket 03 PremiumList Life.
//
// ⛔ KENAPA PENJAGA STATIK DAN BUKAN UJI PERILAKU. Jalur yang dijaga di sini
// hanya berjalan dengan Oracle, dan uji yang menuntut Oracle akan SKIP di
// mesin mana pun yang tidak punya - yaitu diam, yaitu hijau. Penjaga statik
// berbunyi di setiap mesin.

import (
	"os"
	"strings"
	"testing"
)

const berkasNomorPL = "polis_nomor.go"

func sumberNomorPL(t *testing.T) string {
	t.Helper()
	isi, err := os.ReadFile(berkasNomorPL)
	if err != nil {
		t.Fatalf("membaca %s: %v", berkasNomorPL, err)
	}
	return string(isi)
}

// TestGerbangLahirSekaliMenyembuhkanBarisYangTertinggal.
//
// ⛔ LUBANG YANG INI TUTUP. `MIN` dan `MAX` melewati NULL, jadi polis yang
// separuh barisnya bernomor terbaca "sudah bernomor" oleh gerbang
// lahir-sekali. Kembali begitu saja meninggalkan baris yang kosong TIDAK
// PERNAH terisi - dan keadaan itu akan lahir sendiri begitu unggahan CSV
// (tiket 04) menambah peserta SESUDAH polis bernomor.
//
// Yang benar bukan nomor baru dan bukan galat: baris baru diberi nomor yang
// SUDAH ada, dan penghitung tetap tidak tersentuh.
func TestGerbangLahirSekaliMenyembuhkanBarisYangTertinggal(t *testing.T) {
	teks := sumberNomorPL(t)
	if !strings.Contains(teks, "keadaan.Utuh()") {
		t.Errorf("%s tidak memeriksa keadaan.Utuh().\n"+
			"Tanpa itu, polis yang separuh barisnya bernomor kembali lebih "+
			"awal dan sisa barisnya tidak pernah terisi - tanpa satu pun galat.",
			berkasNomorPL)
	}
	if !strings.Contains(teks, "nomorPolis.TulisNomor(ctx, tx, polisID, keadaan.Nomor)") {
		t.Errorf("%s tidak menuliskan nomor yang SUDAH ada ke baris yang "+
			"tertinggal.\nNomor BARU untuk baris yang tertinggal akan "+
			"memecah satu polis menjadi dua nomor.", berkasNomorPL)
	}
}

// TestPenghitungTidakTersentuhBilaSudahBernomor.
//
// ⛔ "Lahir sekali" berarti PENGHITUNGNYA tidak bergerak, bukan sekadar
// nomornya tidak berubah. Penghitung yang naik tanpa nomor tersimpan
// menerbitkan lubang di deret yang tidak akan pernah terisi.
func TestPenghitungTidakTersentuhBilaSudahBernomor(t *testing.T) {
	teks := sumberNomorPL(t)
	iGerbang := strings.Index(teks, "if keadaan.Bernomor()")
	iUrut := strings.Index(teks, "penghitung.UrutNomorBerikut(")
	if iGerbang < 0 || iUrut < 0 {
		t.Fatal("gerbang lahir-sekali atau pemanggilan penghitung tidak ditemukan")
	}
	if iGerbang > iUrut {
		t.Error("gerbang lahir-sekali berdiri SESUDAH penghitung dinaikkan; " +
			"submit kedua akan membakar satu nomor")
	}
}

// TestGerbangKasusTertutupDiLuarTransaksiPenomoran.
//
// ⚠️ Menolak lebih dahulu berarti kasus tertutup tidak pernah sempat
// menyentuh baris penghitung - dan tidak pernah memegang kuncinya.
func TestGerbangKasusTertutupDiLuarTransaksiPenomoran(t *testing.T) {
	teks := sumberNomorPL(t)
	iTutup := strings.Index(teks, "models.KasusPolisTertutup(")
	iTx := strings.Index(teks, "DalamTransaksi(ctx")
	if iTutup < 0 || iTx < 0 {
		t.Fatal("gerbang kasus tertutup atau transaksinya tidak ditemukan")
	}
	if iTutup > iTx {
		t.Error("gerbang kasus tertutup berdiri DI DALAM transaksi penomoran; " +
			"kasus yang sudah ditutup ikut mengantre di kunci penghitung")
	}
}

// TestNolPembentukBentukNomorDiLapisanLayanan - AC tiket 03.
//
// ⛔ Yang merakit bentuk nomor `models.NomorPL`. Bentuk nomor yang tersebar
// di lapisan layanan adalah bentuk yang berubah tanpa satu pun uji murni
// menyadarinya.
func TestNolPembentukBentukNomorDiLapisanLayanan(t *testing.T) {
	teks := sumberNomorPL(t)
	// Komentar dibuang: prosa yang MENERANGKAN bentuk nomor tidak boleh
	// dituduh merakitnya.
	var bersih []string
	for _, b := range strings.Split(teks, "\n") {
		if t := strings.TrimSpace(b); strings.HasPrefix(t, "//") {
			continue
		}
		bersih = append(bersih, b)
	}
	kode := strings.Join(bersih, "\n")
	for _, jejak := range []string{"QR/QP/TP/TR", "PanjangUrutNomorPL", `+ "." +`} {
		if strings.Contains(kode, jejak) {
			t.Errorf("%s merakit bentuk nomor sendiri (%q); itu milik models",
				berkasNomorPL, jejak)
		}
	}
	if !strings.Contains(kode, "models.NomorPL(") {
		t.Errorf("%s tidak memanggil models.NomorPL", berkasNomorPL)
	}
}
