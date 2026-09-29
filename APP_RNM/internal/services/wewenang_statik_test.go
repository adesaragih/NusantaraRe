package services_test

// Penjaga statik tiket 07 - penegakan peran.
//
// Pemilik: tiket 07. Dibaca sesudah: wewenang.go.
//
// ⛔ Letaknya di paket `services`, bukan di `models`. Ronde pertama menaruhnya
// di `models_test` yang membaca `../services/*.go` - arah ketergantungan
// proyek ini satu arah, dan penjaga yang membalikkannya mengajari kebiasaan
// yang justru dilarangnya.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"nusantarare/modul/premiumlist/models"
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
}

// TestSetiapPenulisStatusBergerbangPeran menelusuri SELURUH `internal/`.
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

// polaNamaOrangTetap mencocokkan pemberian TEKS TETAP ke medan yang namanya
// menandakan nama orang - termasuk teks berkutip-balik.
//
// ⛔ STRUKTURAL, bukan berdaftar-nama. Percobaan pertama memuat daftar nama
// operator nyata dari korpus, dan itu sendiri melanggar pagar keamanan brief
// induk ("nilai nama orang tidak disalin ke artefak mana pun").
//
// ⚠️ Ia HEURISTIK, dan batasnya dinyatakan di tiket: nama yang masuk lewat
// konstanta perantara, lewat medan yang tidak terdaftar, atau lewat
// perbandingan `AkunID == "..."` tidak tertangkap. Penjaga yang menyebut
// batasnya lebih berguna daripada penjaga yang mengaku sempurna.
var polaNamaOrangTetap = regexp.MustCompile(
	"(?i)(OpName|NamaOrang|PolicyHolder|NameOfInsured|Tertanggung)" +
		"\\s*[:=]+\\s*[\"`][^\"`]+[\"`]")

// pesanVerbatimYangSah adalah NILAI yang cocok dengan pola di atas tetapi
// BUKAN nama orang - masing-masing beserta alasannya.
//
// ⛔ DIPERSEMPIT 28-09-2026 dengan daftar bernama, BUKAN dengan melonggarkan
// polanya. Penjaga yang menuduh hal yang benar akan dilonggarkan orang, bukan
// dipatuhi - pelajaran yang sudah dibayar dua kali di repo ini (nama tabel
// telanjang, ambang tutup buku).
//
// Yang menuduh di sini: konstanta PESAN VALIDASI yang disalin VERBATIM dari
// `ValidasiUploadPL_act.xml`. Namanya memuat "PolicyHolder"/"Tertanggung"
// sebab itulah KOLOM yang divalidasi, dan nilainya kalimat galat huruf besar
// - bukan nama siapa pun.
//
// ⛔ KUNCINYA DIAMBIL DARI `models`, TIDAK DIKETIK ULANG. Menuliskan
// kalimatnya harfiah di sini akan membuat penjaga ini menuduh DIRINYA
// SENDIRI - dan itu persis yang terjadi pada ronde pertama penyempitan ini.
// Mengambilnya dari konstantanya juga berarti daftar ini ikut basi begitu
// pesannya berubah, alih-alih diam-diam tetap mengecualikan teks lama.
// ⛔ ALASANNYA DI KOMENTAR, BUKAN DI NILAI. Menuliskannya sebagai nilai teks
// membuat BARIS DAFTAR INI SENDIRI cocok dengan polanya - `...Tertanggung:
// "alasan"` - dan penjaga ini menuduh dirinya sendiri. Sudah terjadi, dua
// kali, saat penyempitan ini ditulis.
var pesanVerbatimYangSah = []string{
	// ValidasiUploadPL_act `local.err3` - pesan kolom NAME_OF_INSURED.
	models.PesanNamaTertanggung,
	// ValidasiUploadPL_act `local.err17` - pesan rujukan master POLICY HOLDER.
	models.PesanPolicyHolder,
}

// pesanVerbatimDiterima menjawab apakah sebuah nilai ada di daftar itu.
func pesanVerbatimDiterima(nilai string) bool {
	for _, p := range pesanVerbatimYangSah {
		if p == nilai {
			return true
		}
	}
	return false
}

// polaNilaiTerkutip mengambil nilai di dalam tanda kutip sebuah baris cocok.
var polaNilaiTerkutip = regexp.MustCompile("[\"`]([^\"`]+)[\"`]")

func TestNolNamaOrangDiKode(t *testing.T) {
	diperiksa := 0
	err := filepath.Walk(akarAplikasiPindai, func(jalur string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && lewatiFolderPindai(info.Name()) {
			return filepath.SkipDir
		}
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(jalur, ".go") {
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		diperiksa++
		for _, baris := range strings.Split(string(isi), "\n") {
			if strings.HasPrefix(strings.TrimSpace(baris), "//") {
				continue
			}
			cocok := polaNamaOrangTetap.FindString(baris)
			if cocok == "" {
				continue
			}
			// Nilai sintetis BERAWALAN UJI- memang bentuk yang pagar keamanan
			// brief tuntut untuk fixture. Diperiksa di AWAL teksnya, bukan di
			// mana saja - "Budi UJI-1" bukan nilai sintetis.
			if awalanUjiSintetis(cocok) {
				continue
			}
			// Pesan validasi VERBATIM - didaftar satu per satu beserta
			// alasannya, lihat `pesanVerbatimYangSah`. Yang dicocokkan
			// NILAINYA, bukan seluruh barisnya: baris yang sama dapat ditulis
			// dengan spasi berbeda, dan daftar yang mencocokkan spasi akan
			// lolos begitu seseorang menjalankan gofmt.
			if m := polaNilaiTerkutip.FindStringSubmatch(cocok); m != nil &&
				pesanVerbatimDiterima(m[1]) {
				continue
			}
			t.Errorf("%s memberi nilai tetap ke medan bernama-orang: %s",
				filepath.ToSlash(jalur), strings.TrimSpace(baris))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if diperiksa < 10 {
		t.Fatalf("hanya %d berkas terbaca; pembacanya yang rusak", diperiksa)
	}
}

// awalanUjiSintetis menyatakan teks yang dikutip dimulai dengan UJI-.
func awalanUjiSintetis(cocok string) bool {
	for _, kutip := range []string{"\"", "`"} {
		i := strings.Index(cocok, kutip)
		if i < 0 {
			continue
		}
		return strings.HasPrefix(cocok[i+1:], "UJI-")
	}
	return false
}
