package services_test

// Penjaga statik tiket 07 - penegakan peran.
//
// Pemilik: tiket 07. Dibaca sesudah: wewenang.go.
//
// ⛔ Letaknya di paket `services`, bukan di `models`. Ronde pertama menaruhnya
// di `models_test` yang membaca `../services/*.go` - arah ketergantungan
// proyek ini satu arah, dan penjaga yang membalikkannya mengajari kebiasaan
// yang justru dilarangnya.

// Refactor bentuk B paket 8 (30-09-2026): TestNolNamaOrangDiKode (nama orang) berlaku untuk
// SELURUH aplikasi, jadi pindah apa adanya ke
// `inti/penjaga/lintasaplikasi_test.go`.
import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// polaPenulisStatus mencocokkan baris yang MENULIS status baris adjustment.
//
// ⛔ `=` yang TIDAK diikuti `=` - yaitu penugasan, bukan perbandingan.
// Pola lamanya `KodeStatus\s*=` menangkap `b.KodeStatus == KodeAksep` pula,
// dan itu PEMBACAAN. Ia tidak pernah berbunyi selama satu-satunya pembanding
// di basis kode menyelipkan `)` di antaranya (`TrimSpace(b.KodeStatus) ==`);
// begitu tanda kurung itu hilang, penjaga ini menuduh gerbang tutup yang
// hanya membaca status. Penjaga yang menuduh hal yang benar akan
// dilonggarkan orang, bukan dipatuhi.
var polaPenulisStatus = regexp.MustCompile(
	`KodeStatus\s*=([^=]|$)|PerbaruiStatusBaris\(|TandaiBarisOutstanding\(|CerminkanHeader\(|CabutPenandaDipilih\(`)

// berkasYangBolehMenulisStatus adalah berkas layanan yang memang penulisnya,
// masing-masing dengan gerbang yang wajib ia panggil.
//
// Menambah baris ke sini adalah tindakan SADAR - dan itulah gunanya: setiap
// penulis status baru harus dipertanggungjawabkan, bukan diam-diam ikut.
var berkasYangBolehMenulisStatus = map[string]string{
	// Lapisan layanan - di sinilah gerbangnya, dan tiap berkas menyebut
	// gerbang yang wajib ia panggil.
	"statusbaris.go": "WajibPeranPengubahStatus(pelaku, ke)",
	"tolak.go":       "WajibPeran(pelaku, PeranRejectOutstanding)",
	"pendaftaran.go": "WajibPeran(pelaku, PeranInputRegister)",
	"adjustment.go":  "", // TambahBaris/WarisiKolom - murni, tanpa pelaku
	// Tiket 11: ia menulis Outstanding pada baris BARU, bukan mengubah
	// keputusan baris lama. Keputusan atas baris lama milik Komite Claim
	// Life `[keputusan work owner 2026-09-15]`.
	"hasilkomite.go": "WajibPeran(pelaku, PeranSimpanOutstanding)",
	// Audit A0: jalur akseptasi Claim Life sendiri (SaveAdjustment_Act).
	// Gerbangnya PEMEGANG TAHAP, bukan daftar peran datar - pohon XML
	// membuktikan tombolnya tidak bergerbang peran.
	"akseptasi.go": "WajibPemegangTahap(pelaku, tahap)",
	// GILIRAN-11 paket 1: Save to RNM menulis Outstanding (`0`) HANYA pada
	// baris yang belum berstatus (langkah 22.1.3.2). Gerbangnya pemegang tahap
	// Outstanding - tombolnya di layar itu, dan layar itu dipegang Admin.
	"simpanrnm.go": "WajibPemegangTahap(pelaku, tahap)",
	// Komite Claim Life tiket 04a: akseptasi tingkat akhir komite. Gerbangnya
	// ANGGOTA BERJALAN tangga (ADR-0014), diulang di penulis statusnya.
	"komite_akseptasi.go": "periksaGiliran(kasus, pelaku.AkunID)",

	// Lapisan repository - ia MENJALANKAN SQL-nya, dan memang tidak memegang
	// pelaku. Gerbangnya ada di layanan yang memanggilnya; yang dijaga di
	// sini adalah bahwa daftarnya tidak bertambah diam-diam.
	"klaimlife.go":    "",
	"kolompeserta.go": "",
	"pohonklaim.go":   "",
	"migrasidata.go":  "",
	// Refactor bentuk B (30-09-2026): penerus kontrak Claim Life untuk Komite
	// (`inti/kontrak.KlaimKomite`). Seperti repository, ia tidak memegang
	// pelaku; gerbangnya `periksaGiliran` di komite_akseptasi.go - pemanggil
	// yang terdaftar di atas - dan deklarasi antarmukanya (klaim.go) nol
	// pernyataan.
	"klaimuntukkomite.go": "",
	"klaim.go":            "",
	// Penolak kontrak di sisi Komite: menjawab galat, nol tulisan.
	"klaim_belum_disambung.go": "",
}

// TestSetiapPenulisStatusBergerbangPeran menelusuri SELURUH aplikasi (dulu `internal/`).
//
// Ronde pertama hanya membaca satu berkas dan mencari satu potongan teks,
// sambil mengaku memeriksa "setiap fungsi layanan yang mengubah status" -
// pengakuan yang lebih luas daripada yang diperiksanya. Dua penulis memang
// lolos darinya.
func TestSetiapPenulisStatusBergerbangPeran(t *testing.T) {
	diperiksa := 0
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
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		teks := string(isi)
		var menulis bool
		for _, baris := range strings.Split(teks, "\n") {
			if strings.HasPrefix(strings.TrimSpace(baris), "//") {
				continue
			}
			if polaPenulisStatus.MatchString(baris) {
				menulis = true
				break
			}
		}
		if !menulis {
			return nil
		}
		diperiksa++
		dasar := filepath.Base(jalur)
		gerbang, sah := berkasYangBolehMenulisStatus[dasar]
		if !sah {
			t.Errorf("%s menulis status baris tetapi bukan penulis yang terdaftar; "+
				"daftarkan di berkasYangBolehMenulisStatus beserta gerbangnya, atau "+
				"jangan menulis status di sana", filepath.ToSlash(jalur))
			return nil
		}
		if gerbang != "" && !strings.Contains(teks, gerbang) {
			t.Errorf("%s menulis status tetapi tidak memanggil %s",
				filepath.ToSlash(jalur), gerbang)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if diperiksa < 3 {
		t.Fatalf("hanya %d penulis status ditemukan; pembacanya yang rusak, bukan kodenya",
			diperiksa)
	}
}
