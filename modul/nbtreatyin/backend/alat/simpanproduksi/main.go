// Command simpanproduksi menulis baris Utility1 `SaveJsonPolisTreatyIn_Act` - json_polis TANPA DATA_JSON,
// ACHIEVEMENT, TREATYINPRODUCTION - untuk SATU kasus NB Treaty In yang sudah Resolved-Completed sebelum penulisan
// produksi ada (keputusan work owner 06-10-2026; pemetaan models/produksi.go, MODUL.md "Alat simpanproduksi").
//
//	go run ./modul/nbtreatyin/backend/alat/simpanproduksi -kasus NB-22445             # uji-kering (bawaan)
//	go run ./modul/nbtreatyin/backend/alat/simpanproduksi -kasus NB-22445 -jalankan   # tulis
//
// Uji-kering hanya MEMBACA dan mencetak baris yang akan ditulis. `-jalankan` menulis satu transaksi; aman diulang:
// tabel yang sudah punya baris untuk IDPEGA itu dilewati. Konfigurasi dari lingkungan yang sama dengan `cmd/api`
// (`. .\muat-env.ps1`); `-jalankan` ditolak bila IS_PEGA_PROD=true (pola `pemuatlama`, ADR-U-0005).
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/services"
)

func main() {
	kasus := flag.String("kasus", "", "nomor kasus, mis. NB-22445 (wajib)")
	pengguna := flag.String("pengguna", "", "USERNAME json_polis; kosong = OPERATORID riwayat terakhir kasus")
	jalankan := flag.Bool("jalankan", false, "tulis ke json_polis / ACHIEVEMENT / TREATYINPRODUCTION (bawaan: uji-kering, baca saja)")
	flag.Parse()
	if strings.TrimSpace(*kasus) == "" {
		fmt.Fprintln(os.Stderr, "simpanproduksi: -kasus <NB-n> wajib diisi")
		flag.Usage()
		os.Exit(2)
	}
	if err := jalan(strings.TrimSpace(*kasus), strings.TrimSpace(*pengguna), *jalankan); err != nil {
		fmt.Fprintf(os.Stderr, "simpanproduksi: %v\n", err)
		os.Exit(1)
	}
}

func jalan(kasus, pengguna string, tulis bool) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("konfigurasi: %w", err)
	}
	if !cfg.PunyaOracle() {
		return errors.New("ORACLE_DSN belum dikonfigurasi - jalankan dulu `. .\\muat-env.ps1` dari folder APP_RNM")
	}
	if tulis && cfg.IsPegaProd {
		return errors.New("-jalankan ditolak: IS_PEGA_PROD=true (ADR-U-0005)")
	}
	d, err := db.Open(cfg)
	if err != nil {
		return fmt.Errorf("oracle: %w", err)
	}
	defer func() { _ = d.Close() }()
	p, err := services.PenulisProduksiDariDasar(inti.NewDasar(d))
	if err != nil {
		return err
	}
	ctx, henti := signal.NotifyContext(context.Background(), os.Interrupt)
	defer henti()
	h, err := p.Jalankan(ctx, kasus, pengguna, tulis)
	if err != nil {
		return err
	}
	cetak(h, tulis)
	return nil
}

func cetak(h services.HasilProduksi, tulis bool) {
	s := h.Simpanan
	mode := "UJI-KERING (tidak ada yang ditulis; tambahkan -jalankan untuk menulis)"
	if tulis {
		mode = "DITULIS"
	}
	fmt.Printf("Kasus   : %s\nMode    : %s\nUSERNAME: %s\n\n", s.IDPega, mode, h.Pengguna)

	nasib := func(sudah, baru int) string {
		switch {
		case sudah > 0:
			return fmt.Sprintf("dilewati - sudah ada %d baris", sudah)
		case tulis:
			return fmt.Sprintf("%d baris ditulis", baru)
		}
		return fmt.Sprintf("%d baris akan ditulis", baru)
	}
	fmt.Printf("JSON_POLIS         : %s\n", nasib(h.Ada.JSONPolis, 1))
	cetakBaris(models.KolomJSONPolis, []models.Baris{s.JSONPolis}, nil)
	fmt.Printf("ACHIEVEMENT        : %s\n", nasib(h.Ada.Capaian, len(s.Capaian)))
	cetakBaris(models.KolomCapaian, s.Capaian, []string{"NOPOLIS", "CURRENCY", "PREMIUM", "RICOMM", "BROKERAGE", "NETPREMIUM", "PROPORTIONALTYPE"})
	fmt.Printf("TREATYINPRODUCTION : %s\n", nasib(h.Ada.Produksi, len(s.Produksi)))
	cetakBaris(models.KolomProduksi, s.Produksi, []string{"JN_REAS", "CURR_ID", "PREMI_OGP", "NET_PREMIUM", "PCT_SHARE_PREMI", "DUE_TO", "LAYER", "PROD_DATE"})
	fmt.Println("(NOPOLIS ACHIEVEMENT diisi dari json_polis saat ditulis.)")
}

// cetakBaris mencetak setiap baris; `pilih` kosong = semua kolom katalog.
func cetakBaris(ks []models.Kolom, rows []models.Baris, pilih []string) {
	if len(pilih) == 0 {
		for _, k := range ks {
			pilih = append(pilih, k.Kolom)
		}
	}
	for i, b := range rows {
		var bagian []string
		for _, k := range pilih {
			bagian = append(bagian, k+"="+b[k])
		}
		fmt.Printf("  %d. %s\n", i+1, strings.Join(bagian, "  "))
	}
	fmt.Println()
}
