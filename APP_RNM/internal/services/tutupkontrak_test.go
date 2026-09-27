package services

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Penjaga kontrak butir bb: SETIAP rute pengubah menolak kasus yang tertutup.
//
// ⛔ Kenapa penjaga statik dan bukan tujuh uji perilaku. Aturan yang ditulis
// ulang di tujuh berkas adalah aturan yang suatu hari hanya ada di enam, dan
// yang ketujuh tidak akan berbunyi: tiap berkas hijau sendirian. Itu bentuk
// cacat yang sudah terjadi tiga kali di modul ini (envelope `galat`, rute
// tanpa pemanggil, penanda `IsCheck`). Penjaga ini menagih pemanggilannya
// dari daftar, sehingga layanan pengubah KEDELAPAN pun akan tertagih - asal
// namanya ditambahkan ke sini, dan menambahkannya adalah pekerjaan sadar.

// layananPengubah memetakan berkas layanan ke fungsi masuknya.
//
// ⚠️ Daftar ini DITULIS TANGAN dan itu disengaja. Menurunkannya otomatis dari
// "berkas yang memanggil DalamTransaksi" akan membuat penjaga ini menjaga
// apa pun yang kebetulan ada, bukan apa yang kita putuskan harus dijaga -
// dan daftar yang menyesuaikan diri tidak pernah gagal.
var layananPengubah = map[string]string{
	"akseptasi.go":   "SimpanAdjustment",
	"dol.go":         "Set",
	"hapus.go":       "Hapus",
	"hasilkomite.go": "Tambah",
	"komite.go":      "Serahkan",
	"tahap.go":       "Pindah",
	"statusbaris.go": "ubah",
}

func TestSetiapLayananPengubahMemeriksaKasusTerbuka(t *testing.T) {
	const panggilan = "PastikanKasusTerbuka(ctx, klaimID)"
	for berkas, fungsi := range layananPengubah {
		isi, err := os.ReadFile(berkas)
		if err != nil {
			t.Errorf("membaca %s: %v", berkas, err)
			continue
		}
		teks := string(isi)
		if !strings.Contains(teks, "func ") || !strings.Contains(teks, fungsi+"(ctx context.Context") {
			t.Errorf("%s: fungsi masuk %q tidak ditemukan lagi - daftar "+
				"layananPengubah harus disesuaikan, bukan penjaga ini dimatikan",
				berkas, fungsi)
			continue
		}
		if !strings.Contains(teks, panggilan) {
			t.Errorf("%s (%s) tidak memanggil %s.\n"+
				"Butir bb: sesudah kasus ditutup, SETIAP rute pengubah harus "+
				"menolak. Rute yang lupa memeriksanya tidak akan membuat satu "+
				"pun uji lain merah - ia hanya akan mengubah kasus yang sudah "+
				"selesai.", berkas, fungsi, panggilan)
		}
	}
}

// TestDaftarLayananPengubahMencakupSeluruhRutePengubah menagih arah
// sebaliknya: rute pengubah yang BARU harus masuk daftar di atas.
//
// ⛔ Tanpa uji ini, daftar itu hanya menjaga dirinya sendiri. Rute kedelapan
// dapat lahir lengkap dengan uji perilakunya yang hijau, dan penjaga di atas
// akan tetap diam - sebab ia hanya memeriksa nama yang sudah tertulis.
func TestDaftarLayananPengubahMencakupSeluruhRutePengubah(t *testing.T) {
	// Berkas layanan yang MENULIS ke basis data adalah yang memanggil
	// DalamTransaksi. Setiap satu di antaranya harus ada di daftar, atau
	// dinyatakan alasannya di pengecualian bernama di bawah.
	dikecualikan := map[string]string{
		// Pendaftaran MELAHIRKAN kasus; belum ada kasus untuk ditutup.
		"pendaftaran.go": "melahirkan kasus, bukan mengubah kasus yang ada",
		// Outbox dan pelaksana efek bekerja atas baris antrean, bukan atas
		// kasus - dan efek yang sudah terlanjur diantre tetap harus selesai
		// walau kasusnya kemudian ditutup.
		"antrean.go": "menjalankan antrean efek, bukan mengubah kasus",
		// Tutup itu sendiri: ia yang MENUTUP, dan ia memeriksa dengan
		// pembacaan status kerjanya sendiri.
		"tutup.go": "layanan penutupnya sendiri",
		// Penerbit nomor akseptasi dipanggil DARI akseptasi.go, yang sudah
		// diperiksa di pintunya.
		"spreading.go": "dipanggil dari akseptasi.go yang sudah memeriksa",
		// Tolak MENYALURKAN ke statusbaris.go `ubah`, yang memeriksanya.
		// Penjaga di dua pintu menuju satu ruang adalah dua tempat untuk lupa.
		"tolak.go": "menyalurkan ke statusbaris.go ubah yang sudah memeriksa",
		// services.go adalah tempat DalamTransaksi DIDEFINISIKAN, bukan
		// pemakainya. Ia cocok dengan pencariannya karena namanya sendiri.
		"services.go": "mendefinisikan DalamTransaksi, bukan memakainya",
	}
	masuk, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("membaca direktori layanan: %v", err)
	}
	diperiksa := 0
	for _, e := range masuk {
		nama := e.Name()
		if e.IsDir() || !strings.HasSuffix(nama, ".go") ||
			strings.HasSuffix(nama, "_test.go") {
			continue
		}
		isi, err := os.ReadFile(filepath.Join(".", nama))
		if err != nil {
			t.Fatalf("membaca %s: %v", nama, err)
		}
		if !strings.Contains(string(isi), "DalamTransaksi(ctx") {
			continue
		}
		diperiksa++
		if _, ada := layananPengubah[nama]; ada {
			continue
		}
		if alasan, ada := dikecualikan[nama]; ada {
			if alasan == "" {
				t.Errorf("%s dikecualikan tanpa alasan tertulis", nama)
			}
			continue
		}
		t.Errorf("%s menulis ke basis data tetapi tidak ada di layananPengubah "+
			"maupun di daftar pengecualian.\n"+
			"Butir bb: tiap layanan pengubah harus menolak kasus yang sudah "+
			"ditutup, atau menyatakan alasannya dengan nama.", nama)
	}
	if diperiksa == 0 {
		t.Fatal("nol layanan bertransaksi ditemukan - penjaga ini tidak menjaga apa pun")
	}
}
