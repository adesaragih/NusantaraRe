// Program jalankan menerapkan paket deidentifikasi pada berkas kasus JSON dan
// memverifikasi hasilnya (tiket NB-15).
//
//	go run ./modul/nbfacin/backend/alat/deidentifikasi/jalankan \
//	    -keluar <folder-di-luar-repositori> [-kandidat] <berkas-kasus.json> ...
//
// Berkas yang KEBOCORAN-nya tidak nol TIDAK ditulis, dan program keluar dengan
// kode 1: nilai yang dibuang muncul lagi di medan lain, jadi daftar BUANG belum
// lengkap. `-kandidat` mencetak NAMA kunci bernilai teks untuk ditinjau manusia.
//
// ⛔ Berkas mentah TIDAK pernah masuk repositori. Hasilnya baru boleh disimpan di
// repositori setelah daftar kebocoran dan kandidat identitas yang dicetak program
// ini ditinjau manusia. Program hanya mencetak JALUR dan NAMA KUNCI, tidak pernah
// nilainya.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"nusantarare/modul/nbfacin/backend/alat/deidentifikasi"
)

func main() {
	keluar := flag.String("keluar", "", "folder hasil (di luar repositori)")
	cetakKandidat := flag.Bool("kandidat", false, "cetak nama kunci kandidat identitas")
	flag.Parse()
	if *keluar == "" || flag.NArg() == 0 {
		log.Fatal("jalankan: -keluar dan minimal satu berkas wajib")
	}
	if err := os.MkdirAll(*keluar, 0o755); err != nil {
		log.Fatal(err)
	}
	gagal := false
	for _, f := range flag.Args() {
		asli, err := os.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		asli = bytes.TrimPrefix(asli, []byte{0xEF, 0xBB, 0xBF}) // BOM UTF-8
		bersih, lap, err := deidentifikasi.Bersihkan(asli)
		if err != nil {
			log.Fatalf("%s: %v", filepath.Base(f), err)
		}
		bocor, err := deidentifikasi.Periksa(asli, bersih)
		if err != nil {
			fmt.Printf("%s: VERIFIKASI GAGAL: %v\n", filepath.Base(f), err)
			gagal = true
			continue
		}
		kandidat, err := deidentifikasi.KandidatIdentitas(bersih)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: dikosongkan %v | kebocoran %d | kandidat kunci teks %d\n", filepath.Base(f), lap.Dibuang, len(bocor), len(kandidat))
		for _, b := range bocor {
			fmt.Printf("   bocor di %s\n", b)
		}
		if *cetakKandidat {
			fmt.Printf("   kandidat: %s\n", strings.Join(kandidat, ", "))
		}
		if len(bocor) > 0 {
			fmt.Printf("   TIDAK DITULIS: %d kebocoran\n", len(bocor))
			gagal = true
			continue
		}
		tujuan := filepath.Join(*keluar, strings.TrimSuffix(filepath.Base(f), filepath.Ext(f))+".bersih.json")
		if err := os.WriteFile(tujuan, bersih, 0o644); err != nil {
			log.Fatal(err)
		}
	}
	if gagal {
		os.Exit(1)
	}
}
