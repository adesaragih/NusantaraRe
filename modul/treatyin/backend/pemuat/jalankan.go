//go:build ignore

// Pemuat tabel pendaratan tab Treaty In - dijalankan dengan tangan.
//
// Pemakaian (dari akar aplikasi, dengan `.env` sudah dimuat):
//
//	go run modul/treatyin/backend/pemuat/jalankan.go -batas 10
//	go run modul/treatyin/backend/pemuat/jalankan.go -batas 10 -ikat
//	go run modul/treatyin/backend/pemuat/jalankan.go -semua -ikat
//	go run modul/treatyin/backend/pemuat/jalankan.go -kontrak 1000158,1000649
//	go run modul/treatyin/backend/pemuat/jalankan.go -cocokkan
//	go run modul/treatyin/backend/pemuat/jalankan.go -kosongkan -ikat
//
// ⛔ BAWAANNYA KERING. Tanpa `-ikat` setiap transaksi dibatalkan, dan yang
// tercetak adalah apa yang AKAN terjadi. Memuat 26.536 baris ke skema yang
// memuat 816 tabel warisan bukan hal yang pantas terjadi karena seseorang
// salah menekan panah atas.
//
// ⛔ MENGAPA BERKAS INI BERTANDA `ignore`, DAN BUKAN `cmd/`.
// `TestModulTidakMengimporModulLain` menyatakan *"cmd hanya mengimpor
// inti/...; modul dipasang lewat daftar inti/backend/daftar"*, jadi program
// di `cmd/` TIDAK BOLEH mengimpor repository modul ini. `pelanggaranLetak`
// di penjaga yang sama menyatakan kode Go modul tinggal di
// `modul/<nama>/backend/`, jadi ia juga tidak boleh berdiri di
// `modul/treatyin/alat/`. Yang tersisa adalah bentuk ini: `package main` di
// dalam modulnya sendiri, dikeluarkan dari `go build ./...` oleh tanda
// `ignore`, dijalankan dengan sengaja.
//
// ⚠️ Akibatnya ia TIDAK ikut dikompilasi `go build ./...` maupun diperiksa
// `go vet ./...`, dan karena itu ia dijaga tipis: seluruh keputusan hidup di
// `services.MuatSatuKontrak` dan `repository.PetaPendaratan`, yang ikut
// terkompilasi dan punya ujinya sendiri. Berkas ini hanya gelang dan
// pencetak.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/repository"
	"nusantarare/modul/treatyin/backend/services"
)

func main() {
	var (
		semua     = flag.Bool("semua", false, "muat SELURUH kontrak yang punya dokumen")
		batas     = flag.Int("batas", 0, "muat N kontrak pertama (urut pengenal naik)")
		kontrak   = flag.String("kontrak", "", "daftar pengenal kontrak, dipisah koma")
		ikat      = flag.Bool("ikat", false, "COMMIT tiap kontrak yang cocok; tanpa ini seluruhnya dibatalkan")
		kosongkan = flag.Bool("kosongkan", false, "kosongkan kedelapan tabel, nol pemuatan")
		cocokkan  = flag.Bool("cocokkan", false, "cetak cacah baris kedelapan tabel, nol tulisan")
	)
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("konfigurasi: %v", err)
	}
	// Satu-satunya pagar yang tidak pernah dicabut.
	if cfg.IsPegaProd {
		log.Fatal("menolak berjalan saat IS_PEGA_PROD=true (ADR-U-0005)")
	}
	if !cfg.PunyaOracle() {
		log.Fatal("ORACLE_DSN kosong")
	}
	d, err := db.Open(cfg)
	if err != nil {
		log.Fatalf("membuka oracle: %v", err)
	}
	defer func() { _ = d.Close() }()

	ctx := context.Background()
	if err := d.Ping(ctx); err != nil {
		log.Fatalf("oracle tidak terjangkau: %v", err)
	}
	g := repository.Baru(d)
	fmt.Printf("skema %s · ikat=%v\n", cfg.OracleSchema, *ikat)

	switch {
	case *cocokkan:
		cetakCacah(ctx, g)
		return
	case *kosongkan:
		jalankanPengosongan(ctx, d, g, *ikat)
		return
	}

	daftar := pilihKontrak(ctx, g, *semua, *batas, *kontrak)
	if len(daftar) == 0 {
		log.Fatal("nol kontrak dipilih; pakai -semua, -batas, atau -kontrak")
	}
	fmt.Printf("%d kontrak dipilih\n", len(daftar))
	jalankanPemuatan(ctx, d, g, daftar, *ikat)
}

func pilihKontrak(ctx context.Context, g *repository.Gudang, semua bool, batas int, daftar string) []string {
	if strings.TrimSpace(daftar) != "" {
		var out []string
		for _, s := range strings.Split(daftar, ",") {
			if s = strings.TrimSpace(s); s != "" {
				out = append(out, s)
			}
		}
		return out
	}
	n := batas
	if semua {
		n = 0
	} else if n == 0 {
		return nil
	}
	id, err := g.DaftarMasterID(ctx, n)
	if err != nil {
		log.Fatalf("mendaftar kontrak: %v", err)
	}
	return id
}

// jalankanPemuatan memuat kontrak satu per satu, MASING-MASING di dalam
// transaksinya sendiri.
//
// ⛔ Satu transaksi per kontrak, bukan satu untuk 1.854. Transaksi sebesar
// itu mengunci lama, membengkakkan undo, dan - yang terburuk - membuat satu
// kontrak rusak membatalkan 1.853 yang baik. Yang gagal dicatat namanya dan
// dilewati; sisanya tetap jalan.
func jalankanPemuatan(ctx context.Context, d *db.DB, g *repository.Gudang, daftar []string, ikat bool) {
	mulai := time.Now()
	var (
		hasil          []services.HasilMuat
		gagal, tidakOK int
		totalBaris     int
	)
	for ke, id := range daftar {
		tx, err := d.Mulai(ctx)
		if err != nil {
			log.Fatalf("membuka transaksi: %v", err)
		}
		h, err := services.MuatSatuKontrak(ctx, g, tx, id)
		if err != nil {
			_ = tx.Rollback()
			gagal++
			fmt.Printf("  GAGAL  %-9s %v\n", id, err)
			continue
		}
		if !h.Cocok() {
			_ = tx.Rollback()
			tidakOK++
			fmt.Printf("  SELISIH %-9s %v\n", id, h.Selisih)
			continue
		}
		if ikat {
			if err := tx.Commit(); err != nil {
				gagal++
				fmt.Printf("  GAGAL IKAT %-9s %v\n", id, err)
				continue
			}
		} else {
			_ = tx.Rollback()
		}
		for _, n := range h.DiTabel {
			totalBaris += n
		}
		hasil = append(hasil, h)
		if (ke+1)%100 == 0 {
			fmt.Printf("  ... %d/%d kontrak, %d baris, %s\n", ke+1, len(daftar), totalBaris, time.Since(mulai).Round(time.Second))
		}
	}

	fmt.Printf("\ncocok    : %d kontrak, %d baris\n", len(hasil), totalBaris)
	fmt.Printf("selisih  : %d\n", tidakOK)
	fmt.Printf("gagal    : %d\n", gagal)
	fmt.Printf("lama     : %s\n", time.Since(mulai).Round(time.Second))
	if !ikat {
		fmt.Println("\n⚠️  KERING - seluruh transaksi dibatalkan. Tambahkan -ikat untuk menyimpan.")
	}

	if asing := services.RingkasTakTerpetakan(hasil); len(asing) != 0 {
		fmt.Println("\n⛔ KUNCI JSON TANPA KOLOM - ia mendarat sebagai ketiadaan:")
		for _, tabel := range urut(asing) {
			fmt.Printf("   %-28s %v\n", tabel, asing[tabel])
		}
	}
	if tidakOK != 0 || gagal != 0 {
		os.Exit(1)
	}
}

func jalankanPengosongan(ctx context.Context, d *db.DB, g *repository.Gudang, ikat bool) {
	id, err := g.DaftarMasterID(ctx, 0)
	if err != nil {
		log.Fatalf("mendaftar kontrak: %v", err)
	}
	tx, err := d.Mulai(ctx)
	if err != nil {
		log.Fatalf("membuka transaksi: %v", err)
	}
	total := map[string]int64{}
	for _, m := range id {
		dibuang, err := g.KosongkanKontrak(ctx, tx, m)
		if err != nil {
			_ = tx.Rollback()
			log.Fatalf("mengosongkan %s: %v", m, err)
		}
		for t, n := range dibuang {
			total[t] += n
		}
	}
	if ikat {
		if err := tx.Commit(); err != nil {
			log.Fatalf("mengikat pengosongan: %v", err)
		}
	} else {
		_ = tx.Rollback()
	}
	for _, t := range urutInt64(total) {
		fmt.Printf("  %-28s %d baris dibuang\n", t, total[t])
	}
	if !ikat {
		fmt.Println("\n⚠️  KERING - pengosongan dibatalkan. Tambahkan -ikat.")
	}
}

func cetakCacah(ctx context.Context, g *repository.Gudang) {
	cacah, err := g.CacahBarisSeluruhnya(ctx, nil)
	if err != nil {
		log.Fatalf("menghitung: %v", err)
	}
	total := 0
	for _, p := range repository.PetaPendaratan {
		n := cacah[p.Tabel]
		total += n
		tanda := "  "
		if n != p.CacahTerukur {
			tanda = "⚠️"
		}
		fmt.Printf("  %s %-28s %7d  (sapuan 3 Okt 2026: %d)\n", tanda, p.Tabel, n, p.CacahTerukur)
	}
	// ⛔ Penyebutnya DIJUMLAHKAN dari peta, bukan ditulis sebagai angka.
	// Angka yang ditulis tangan berhenti benar pada tabel kesepuluh, dan
	// berhentinya tidak terlihat - ia hanya menjadi baris yang selalu
	// bertanda peringatan.
	mau := 0
	for _, p := range repository.PetaPendaratan {
		mau += p.CacahTerukur
	}
	tandaTotal := "  "
	if total != mau {
		tandaTotal = "⚠️"
	}
	fmt.Printf("  %s %-28s %7d  (sapuan: %d)\n", tandaTotal, "TOTAL", total, mau)
	cetakPenandaRetro()
}

// cetakPenandaRetro menaruh lubang yang DISENGAJA di depan mata orang yang
// sedang merekonsiliasi.
//
// ⛔ Dicetak pada tiap `-cocokkan`, bukan disimpan di dokumen saja. Catatan
// yang hanya ada di berkas dibaca oleh yang mencarinya; yang ini lewat di
// depan yang TIDAK mencarinya, dan itulah yang diperlukan - orang yang
// menyimpulkan "27.238 cocok, selesai" justru yang harus melihatnya.
func cetakPenandaRetro() {
	fmt.Println()
	fmt.Println("⚠️ ", repository.AlasanRetroTertunda)
	for _, id := range repository.KontrakRetroTertunda {
		fmt.Printf("     kontrak %s - larik %s TIDAK dimuat%s", id, repository.LarikRetroTertunda, "\n")
	}
}

func urut(m map[string][]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func urutInt64(m map[string]int64) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
