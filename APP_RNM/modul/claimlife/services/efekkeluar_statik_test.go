package services_test

// Penjaga statik tiket 12 - alamat layanan dan payload keluar.
//
// Pemilik: tiket 12. Dibaca sesudah: efekkeluar.go.
//
// Dua aturan dijaga, keduanya diminta tiket dalam bentuk test:
//
//  1. ADR-U-0013: nol URL sebagai literal, konstanta, MAUPUN env var. Yang
//     boleh jadi konstanta hanyalah kunci kategori. Pemisahan dev-prod
//     terjadi lewat isi tabel per-database, bukan percabangan di kode.
//  2. `[keputusan work owner 2026-09-16]`: konversi ke produksi tidak lagi
//     mengirim payload JSON; hilir membaca langsung dari tabel klaim.
//
// ⛔ Ketiga kelemahan penjaga tiket 10 ditutup SEJAK AWAL di sini:
// pengecualian berkunci jalur penuh (bukan nama berkas), akar telusurnya akar
// modul (bukan `internal/`), dan pencacahnya naik SEBELUM pengecualian
// sehingga swa-periksa tidak dapat dikelabui pengecualian yang terlalu lebar.

// Refactor bentuk B paket 8 (30-09-2026): TestNolAlamatLayananDiKode (alamat layanan, ADR-U-0013) berlaku untuk
// SELURUH aplikasi, jadi pindah apa adanya ke
// `inti/backend/penjaga/lintasaplikasi_test.go`.
import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/outbox"
	"nusantarare/modul/claimlife/services"
)

// TestFlagLingkunganTidakMenggerbangiPenyimpanan - penyimpangan sadar.
//
// ⛔ Flag lingkungan hanya menggerbangi EFEK KELUAR. Di Pega `IsPEGAPROD`
// menggerbangi simpan utama juga; bila kita menirunya, lingkungan
// non-produksi tidak dapat dipakai menguji sama sekali - klaim tidak akan
// pernah tersimpan di sana.
//
// ⛔ PENJAGA TEKSTUAL DICABUT - ia terbukti dapat dielakkan.
//
// Ronde pertama memeriksa teks sumber: `polaCabangLingkungan` menuntut
// `Produksi()` BERSEBELAHAN dengan `if`. Dua elakan dibangun, dikompilasi,
// dan dijalankan - keduanya hijau:
//
//	prod := p.penyalur.lingkungan.AdalahProduksi()
//	if !prod { return nil }              // menyimpan dilewati di non-produksi
//
//	var _ = "if !p.lingkungan.AdalahProduksi()" // memuaskan strings.Contains
//
// Elakan kedua paling telak: ia memuaskan pemeriksaan "gerbangnya masih ada"
// dengan sebuah LITERAL TEKS, sementara gerbang aslinya dihapus - efek keluar
// menyala di setiap lingkungan, email nyata dari lingkungan uji.
//
// Pelajarannya bukan "perkuat polanya". Pola atas teks sumber selalu dapat
// dielakkan oleh penulisan ulang yang setara. Yang menggantikannya adalah
// penjaga PERILAKU di `efekkeluar_test.go` dan di bawah: ia memanggil kodenya
// dan memeriksa apa yang TERJADI, sehingga penulisan ulang apa pun yang
// mengubah perilaku akan tertangkap apa pun bentuknya.
func TestFlagLingkunganTidakMenggerbangiPenyimpanan(t *testing.T) {
	// ⛔ PERILAKU, bukan teks. Layanan disusun di lingkungan BUKAN produksi -
	// bawaan `svc.Komite()` - lalu jalur penyimpanannya dipanggil TANPA Oracle.
	//
	// Yang benar: ia berjalan sampai menyentuh basis data, lalu berhenti di
	// `ErrTanpaOracle`. Bila seseorang menyisipkan gerbang lingkungan di depan
	// penyimpanan - dalam bentuk APA PUN, termasuk yang dinaikkan ke variabel
	// atau ditulis sebagai `switch` - jalur itu akan pulang lebih awal dan
	// mengembalikan nil, dan test ini merah. Penulisan ulang yang setara tidak
	// menolongnya: yang diperiksa akibatnya, bukan bentuknya.
	svc := services.New(nil)
	err := svc.Komite().Serahkan(context.Background(),
		inti.Pelaku{AkunID: "UJI-AKUN", Peran: []string{inti.PeranAdmin}},
		"CLM-1", "P-1", "A-1", saatUji)
	if err == nil {
		t.Fatal("penyimpanan pulang tanpa galat di lingkungan bukan produksi; " +
			"flag lingkungan menggerbangi SIMPAN, bukan hanya efek keluar - " +
			"tanpa Oracle jalur ini seharusnya berhenti di ErrTanpaOracle")
	}
	if !errors.Is(err, db.ErrTanpaOracle) {
		t.Errorf("galat = %v, mau ErrTanpaOracle; jalur penyimpanan tidak "+
			"sampai menyentuh basis data", err)
	}
}

// TestGerbangLingkunganBenarBenarMenggerbangi - pasangan test di atas.
//
// ⛔ Test di atas memastikan gerbangnya TIDAK ada di penyimpanan. Test ini
// memastikan ia ADA di efek keluar. Tanpa keduanya, menghapus gerbang itu
// seluruhnya akan tetap hijau - dan itulah elakan ketiga yang terbukti
// berhasil terhadap ronde pertama.
func TestGerbangLingkunganBenarBenarMenggerbangi(t *testing.T) {
	efek := &efekUji{nama: "UJI-GERBANG"}
	antre := &antreanUji{}

	outbox.NewPenyalur(inti.BukanProduksi, antre, efek).
		Salurkan(context.Background(), outbox.MuatanEfek{KlaimID: "CLM-1"})
	if efek.dipanggil != 0 {
		t.Fatalf("efek berjalan %d kali di BUKAN produksi; gerbangnya hilang - "+
			"email nyata akan terkirim dari lingkungan uji", efek.dipanggil)
	}

	outbox.NewPenyalur(inti.Produksi, antre, efek).
		Salurkan(context.Background(), outbox.MuatanEfek{KlaimID: "CLM-1"})
	if efek.dipanggil != 1 {
		t.Errorf("efek berjalan %d kali di produksi, mau 1; gerbangnya menutup "+
			"terlalu rapat dan efek keluar tidak pernah berjalan di mana pun",
			efek.dipanggil)
	}
}
