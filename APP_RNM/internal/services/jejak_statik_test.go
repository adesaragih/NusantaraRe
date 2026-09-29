package services_test

// Penjaga statik tiket 09 - jejak audit setiap transisi.
//
// Pemilik: tiket 09. Dibaca sesudah: statusbaris.go, tahap.go.
//
// ADR-U-0007: sistem baru merekam SIAPA dan KAPAN untuk setiap transisi status
// klaim DAN setiap jalur balik. Tabelnya belum ada (butir am masih `[USULAN]`),
// tetapi bentuknya sudah dapat dijaga: setiap penulis transisi wajib merekam
// jejaknya, dan merekamnya DI DALAM fungsi yang sama.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// polaPenulisTransisi mencocokkan pemanggilan repository yang MENGUBAH keadaan
// tangga: status baris, pencerminan, dan perpindahan tahap.
//
// ⚠️ CAKUPANNYA DINYATAKAN, bukan dilebihkan. Ia TIDAK memuat `KodeStatus =`,
// sehingga jalur `Pendaftaran.Daftar` -> `TandaiOutstanding` ->
// `PohonKlaim.Simpan` - yang menulis Outstanding saat INSERT, bukan lewat
// `PerbaruiStatusBaris` - berada di luar jangkauannya. Celah itu nyata,
// dinyatakan di tiket, dan menunggu butir am bersama tabelnya. Pesan galat di
// bawah karena itu menyebut apa yang benar-benar diperiksa, bukan "setiap
// transisi".
var polaPenulisTransisi = regexp.MustCompile(
	`PerbaruiStatusBaris\(|TandaiBarisOutstanding\(|CerminkanHeader\(|PerbaruiTahap\(`)

// TestSetiapPenulisTransisiMerekamJejak menelusuri seluruh lapisan layanan.
//
// ⛔ Diperiksa PER FUNGSI, bukan per berkas. Ronde pertama memeriksa berkasnya
// - dan komentarnya sendiri mengaku menjaga "di dalam transaksi yang sama",
// padahal satu berkas dengan dua fungsi bebas (satu menulis transisi, satu
// lagi memanggil Rekam tanpa pernah dijalankan) akan lolos. Pengakuan yang
// lebih kuat daripada yang diperiksanya adalah cacat yang berulang di sesi
// ini; ia diperbaiki di tempat ia muncul.
// berkasPenerusKontrak - penerus dan penolak kontrak Claim Life untuk
// Komite; lihat pengecualian di dalam uji di bawah.
var berkasPenerusKontrak = map[string]string{
	"klaimuntukkomite.go":      "penerus kontrak (Claim Life): satu baris penerus per metode",
	"klaim_belum_disambung.go": "penolak kontrak (Komite): setiap metode menjawab galat, nol tulisan",
}

func TestSetiapPenulisTransisiMerekamJejak(t *testing.T) {
	diperiksa := 0
	penerusKontrakTerlihat := map[string]bool{}
	err := filepath.Walk(akarAplikasiPindai, func(jalur string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && lewatiFolderPindai(info.Name()) {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(jalur, ".go") ||
			strings.HasSuffix(jalur, "_test.go") {
			return nil
		}
		// Repository MENJALANKAN SQL-nya; jejaknya direkam layanan yang
		// memanggilnya. Yang dijaga di sini lapisan layanan.
		if strings.Contains(filepath.ToSlash(jalur), "/repository/") {
			return nil
		}
		// Refactor bentuk B (30-09-2026): penerus kontrak Claim Life untuk
		// Komite diperlakukan sama dengan repository - satu baris penerus per
		// metode. Jejaknya direkam pemanggil Komite (komite_akseptasi.go) di
		// fungsi yang sama dengan panggilan kontraknya, dan fungsi itu tetap
		// diperiksa di sini.
		if _, ada := berkasPenerusKontrak[filepath.Base(jalur)]; ada {
			penerusKontrakTerlihat[filepath.Base(jalur)] = true
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		for _, fn := range pecahPerFungsi(string(isi)) {
			// ⛔ Komentar dibuang di KEDUA sisi. Ronde sebelumnya hanya
			// membuangnya di sisi penulis, sehingga satu baris prosa yang
			// MENYEBUT `jejak.Rekam(ctx, tx,` - misalnya komentar kepala
			// fungsi berikutnya, yang ikut terpotong ke tubuh sebelumnya -
			// cukup untuk membungkam penjaga ini.
			tubuh := buangKomentar(fn.tubuh)
			if !polaPenulisTransisi.MatchString(tubuh) {
				continue
			}
			diperiksa++
			if !strings.Contains(tubuh, "jejak.Rekam(ctx, tx,") {
				t.Errorf("%s: %s memanggil penulis transisi repository tetapi tidak "+
					"merekam jejaknya di fungsi yang sama; ADR-U-0007 menuntut siapa "+
					"dan kapan untuk setiap transisi", filepath.ToSlash(jalur), fn.nama)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if diperiksa < 2 {
		t.Fatalf("hanya %d fungsi penulis transisi ditemukan; pembacanya yang rusak, "+
			"bukan kodenya", diperiksa)
	}
	// Pengecualian yang tidak terpakai adalah pengecualian yang basi.
	for nama := range berkasPenerusKontrak {
		if !penerusKontrakTerlihat[nama] {
			t.Errorf("%s tidak ditemukan; cabut pengecualiannya", nama)
		}
	}
}

// TestJejakBawaanGagalTerang - yang belum diputuskan terlihat sebagai satu
// galat yang menyebut apa yang ditunggu, bukan sebagai transisi yang diam-diam
// tak tercatat.
func TestJejakBawaanGagalTerang(t *testing.T) {
	for _, berkas := range []string{"statusbaris.go", "tahap.go"} {
		isi, err := os.ReadFile(filepath.Join("..", "services", berkas))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(isi), "jejak: jejak.JejakBelumDiputuskan{}") {
			t.Errorf("%s tidak memakai JejakBelumDiputuskan sebagai bawaan; "+
				"transisi akan berjalan tanpa jejak, dan tidak ada yang tahu", berkas)
		}
	}
}

type fungsiGo struct {
	nama  string
	tubuh string
}

// polaAwalFungsi mencocokkan baris pembuka sebuah fungsi atau method.
var polaAwalFungsi = regexp.MustCompile(`(?m)^func (?:\([^)]*\) )?([A-Za-z_][A-Za-z0-9_]*)`)

// pecahPerFungsi memotong sebuah berkas Go menjadi tubuh fungsi-fungsinya.
//
// Pemotongan berdasarkan baris `func` di kolom nol - cukup untuk penjaga
// statik, dan tidak menuntut pengurai Go yang utuh. Batasnya dinyatakan:
// fungsi bersarang (closure) ikut ke tubuh induknya, yang justru yang
// diinginkan di sini - sisipan di dalam transaksi memang closure.
func pecahPerFungsi(isi string) []fungsiGo {
	letak := polaAwalFungsi.FindAllStringSubmatchIndex(isi, -1)
	out := make([]fungsiGo, 0, len(letak))
	for i, l := range letak {
		akhir := len(isi)
		if i+1 < len(letak) {
			akhir = letak[i+1][0]
		}
		out = append(out, fungsiGo{nama: isi[l[2]:l[3]], tubuh: isi[l[0]:akhir]})
	}
	return out
}

// buangKomentar menghapus baris komentar sebelum pencocokan, supaya prosa yang
// MENYEBUT sebuah pemanggilan tidak dituduh sebagai pemanggilannya.
func buangKomentar(isi string) string {
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
