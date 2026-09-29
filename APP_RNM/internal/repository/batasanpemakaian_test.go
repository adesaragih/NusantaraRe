package repository

// Penjaga batas pemakaian - TANPA Oracle.
//
// Untuk apa berkas ini: beberapa aturan proyek ini tidak dapat dijaga oleh
// kompilator maupun oleh test perilaku, sebab yang dilarang bukan hasilnya
// melainkan SIAPA yang memanggil. Test di sini membaca berkas sumber sebagai
// teks dan memeriksa hal itu.
//
// Dibaca sesudah: pohonklaim.go.
//
// Istilah:
//   - test statik : test yang membaca kode sumber, bukan menjalankannya.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// akarModul menunjuk folder APP_RNM dari folder paket ini.
const akarModul = "../.."

// berkasGoSelainTest mengumpulkan seluruh berkas .go yang BUKAN test.
func berkasGoSelainTest(t *testing.T) map[string]string {
	t.Helper()
	hasil := map[string]string{}
	err := filepath.Walk(filepath.FromSlash(akarModul), func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// frontend dan hasil bangun tidak memuat kode Go yang relevan.
			switch info.Name() {
			case "frontend", "node_modules", "bin", ".git":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(p, ".go") || strings.HasSuffix(p, "_test.go") {
			return nil
		}
		isi, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		hasil[filepath.ToSlash(p)] = string(isi)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hasil) == 0 {
		t.Fatal("nol berkas .go terbaca; pembacanya yang rusak, bukan kodenya")
	}
	return hasil
}

// PohonKlaim.HapusFisik tidak boleh dipanggil dari kode yang bukan test.
//
// ADR-U-0031 `[keputusan work owner]`: penghapusan klaim di jalur pengguna
// adalah PENANDA ditambah nilai balik, dan Akibat 1-nya berbunyi "nol perintah
// hapus fisik pada jalur pengguna di lapisan mana pun". Satu-satunya pemakai
// yang sah adalah test yang membuktikan kaskade `ON DELETE CASCADE` (AC 38) -
// itu pembuktian bentuk basis data, bukan jalur pengguna.
//
// ⚠️ Dua ronde sebelumnya keduanya salah, dan cara salahnya berlawanan.
// Pola pertama mencocokkan `.Hapus(` apa pun, sehingga `services.Penghapusan
// .Hapus` - yang justru MENOLAK menghapus - ikut tertuduh. Pola kedua
// mencocokkan `.Hapus(ctx, tx` dan karena itu LOLOS oleh nama variabel lain
// (`trx`, `h.tx`) maupun pemanggilan yang terpotong baris, padahal
// komentarnya mengaku menangkap "setiap pemanggilan". Sebab keduanya sama:
// dua hal berbeda memakai satu nama.
//
// Yang diperbaiki adalah NAMANYA, bukan polanya. Hapus fisik kini bernama
// `HapusFisik`, sehingga penjaga ini mencocokkan nama - bukan menebak bentuk
// pemanggilan - dan tidak dapat lolos oleh nama variabel maupun baris yang
// terpotong.
func TestHapusFisikTidakDipanggilDiLuarTest(t *testing.T) {
	for nama, isi := range berkasGoSelainTest(t) {
		for _, baris := range strings.Split(isi, "\n") {
			// Definisi metodenya sendiri jelas bukan pemanggilan.
			if strings.Contains(baris, "func (r *PohonKlaim) HapusFisik(") {
				continue
			}
			if strings.Contains(baris, "HapusFisik(") {
				t.Errorf("%s memanggil hapus fisik: %s", nama, strings.TrimSpace(baris))
			}
		}
	}
}

// Lapisan handlers tidak boleh menyentuh repository langsung.
//
// Arah ketergantungan proyek ini satu arah: handlers -> services -> repository.
// Memotongnya membuat aturan dagang tersebar ke lapisan yang tugasnya hanya
// menerima permintaan HTTP.
func TestHandlersTidakMengimporRepository(t *testing.T) {
	diperiksa := 0
	// Refactor bentuk B (30-09-2026): handlers SETIAP modul, dan repository
	// SETIAP modul. `inti/db` ikut dilarang: isinya dulu kepala paket
	// repository (koneksi dan transaksi), jadi larangan lama tetap utuh.
	polaRepository := regexp.MustCompile(`"nusantarare/(internal/repository|modul/[^/"]+/repository|inti/db)"`)
	for nama, isi := range berkasGoSelainTest(t) {
		if !strings.Contains(nama, "/handlers/") {
			continue
		}
		diperiksa++
		if polaRepository.MatchString(isi) {
			t.Errorf("%s mengimpor repository; seharusnya lewat services", nama)
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol berkas handlers terbaca; pembacanya yang rusak")
	}
}

// Baris datar warisan ditulis 18 dari 55 kolom - dan angkanya dikunci di sini.
//
// `Simpan` menulis sebagian kolom saja; 37 sisanya tinggal NULL dan didaftar
// namanya di komentar fungsi itu. Komentar dapat basi tanpa ada yang tahu, jadi
// pembagiannya dikunci test. Kalau kelak sebuah kolom mulai ditulis - misalnya
// tiket 02 atau 03 membawa atribut polisnya - test ini gagal, dan komentar itu
// ikut diperbarui alih-alih menyesatkan diam-diam.
func TestCacahKolomWarisanYangDitulisSimpan(t *testing.T) {
	isi, ada := berkasGoSelainTest(t)["../../internal/repository/pohonklaim.go"]
	if !ada {
		t.Fatal("pohonklaim.go tidak terbaca; pembacanya yang rusak")
	}
	awal := strings.Index(isi, "INSERT INTO %s\n\t\t\t(ID, CASEID")
	if awal < 0 {
		t.Fatal("pernyataan INSERT baris datar warisan tidak ketemu")
	}
	akhir := strings.Index(isi[awal:], ")\n\t\t\tVALUES")
	if akhir < 0 {
		t.Fatal("akhir daftar kolom tidak ketemu")
	}
	daftar := isi[awal : awal+akhir]
	daftar = daftar[strings.Index(daftar, "(")+1:]

	ditulis := map[string]bool{}
	for _, n := range strings.Split(daftar, ",") {
		if n = strings.TrimSpace(n); n != "" {
			ditulis[n] = true
		}
	}
	const mau = 18
	if len(ditulis) != mau {
		t.Errorf("Simpan menulis %d kolom warisan, komentar fungsinya menyebut %d",
			len(ditulis), mau)
	}
	seluruh := NamaKolomBarisLama()
	var kosong int
	for _, n := range seluruh {
		if !ditulis[n] {
			kosong++
		}
		delete(ditulis, n)
	}
	if len(ditulis) != 0 {
		t.Errorf("Simpan menulis kolom yang tidak ada di daftar 55: %v", ditulis)
	}
	if kosong != len(seluruh)-mau {
		t.Errorf("kolom yang tinggal NULL = %d, mau %d", kosong, len(seluruh)-mau)
	}
}

// ADR-U-0033: nol nama tabel telanjang di dalam teks query.
//
// "Telanjang" berarti tanpa awalan skema. Query begitu benar hanya selama sesi
// kebetulan menunjuk skema yang tepat, dan salahnya baru muncul saat pindah
// lingkungan - jauh dari orang yang menulisnya. ADR-U-0033 Akibat 3 menuntut
// test yang menemukannya gagal; sampai ronde 3 test itu tidak pernah ada, dan
// SYS.ALL_OBJECTS sempat lolos sebagai ALL_OBJECTS telanjang.
//
// Yang dianggap SAH sesudah FROM / INTO / UPDATE / JOIN:
//   - "%s"            nama yang sudah dilewatkan Qualify
//   - "A.B"           sudah berawalan skema, termasuk SYS.
//   - "{skema}.B"     penanda di berkas migrasi
//   - "DUAL"          tabel semu milik Oracle, tidak punya skema
//
// Dan `FOR UPDATE` DILEWATI: ia klausa penguncian baris, bukan pernyataan
// UPDATE, sehingga kata sesudahnya (`SKIP`, `NOWAIT`, atau tidak ada) bukan
// nama tabel. Penyempitan ini dibuktikan MASIH MENGGIGIT sebelum dipakai.
// ⛔ KOMENTAR DIBUANG SEBELUM PENCOCOKAN, sejak 27-09-2026. Tanpa itu
// prosa yang MENERANGKAN sebuah query - hal biasa di repositori ini -
// dituduh sebagai query-nya. Yang menyalakannya satu kalimat di
// `diagnosa.go`: "menirunya dengan N UPDATE berarti daftar berlubang", dan
// penjaga membaca `UPDATE berarti` sebagai `UPDATE <nama tabel>`.
//
// ⚠️ Penjaga yang menuduh hal yang BENAR akan dilonggarkan orang,
// bukan dipatuhi. Jadi ia dipersempit sekarang - dan dibuktikan MASIH
// MENGGIGIT sebelum dipakai, sebagaimana penyempitan `FOR UPDATE` di atas.
//
// Kembaran yang sudah lebih dulu melakukan hal yang sama, dengan sebab
// yang sama persis: `polaKomentar` di `models/kodestatus_test.go`.
var polaKomentarBaris = regexp.MustCompile(`(?m)^\s*//.*$`)

// polaTabelTelanjang mencari kata sesudah FROM / INTO / UPDATE / JOIN.
var polaTabelTelanjang = regexp.MustCompile(
	`(?i)(\bFOR\s+)?\b(FROM|INTO|UPDATE|JOIN)\s+([A-Za-z_{%][\w{}%.]*)`)

func TestNolNamaTabelTelanjangDiQuery(t *testing.T) {
	diperiksa := 0
	for nama, isi := range berkasGoSelainTest(t) {
		// Refactor bentuk B (30-09-2026): SQL kini juga tinggal di
		// `modul/*/repository` dan `inti/*`. Dulu hanya `/internal/repository/`
		// - saringan itu diam-diam menyempit begitu kode pindah (cacah log
		// dasar 275 rujukan; sesudah paket 1-2 tanpa perbaikan ini, 168).
		if !strings.Contains(nama, "/repository/") && !strings.Contains(nama, "/inti/") {
			continue
		}
		isi = polaKomentarBaris.ReplaceAllString(isi, "")
		for _, m := range polaTabelTelanjang.FindAllStringSubmatch(isi, -1) {
			if m[1] != "" { // klausa `FOR UPDATE`, bukan pernyataan
				continue
			}
			objek := m[3]
			diperiksa++
			switch {
			case strings.Contains(objek, "."): // berawalan skema atau {skema}
			case strings.Contains(objek, "%s"): // datang dari Qualify
			case strings.EqualFold(objek, "DUAL"): // tabel semu Oracle
			default:
				t.Errorf("%s: %s %s - nama tabel telanjang (ADR-U-0033)",
					nama, strings.ToUpper(m[2]), objek)
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol rujukan tabel terbaca; pembacanya yang rusak")
	}
	t.Logf("%d rujukan tabel diperiksa", diperiksa)
}

// Setiap pemanggil skemauji.Buka memeriksa BolehDilewati sebelum melewat.
//
// Kenapa ini dijaga test dan bukan diserahkan ke kehati-hatian: pola enam baris
// itu disalin ke ENAM tempat di tiga paket. Test db ketujuh yang menyalin
// bentuk lama - `if err != nil { t.Skipf(...) }` - akan membuka kembali lubang
// yang ditutup ronde 4, dan tidak ada yang menyadarinya sebab test yang
// MELEWAT tetap terlihat hijau.
//
// Yang dijaga: sesudah setiap `skemauji.Buka()`, dalam sepuluh baris
// berikutnya, harus ada `BolehDilewati`.
func TestSetiapPemanggilBukaMemeriksaBolehDilewati(t *testing.T) {
	var berkas []string
	err := filepath.Walk(filepath.FromSlash(akarModul), func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			switch info.Name() {
			case "frontend", "node_modules", "bin", ".git":
				return filepath.SkipDir
			}
			return nil
		}
		// Berkas ini sendiri dikecualikan: ia menyebut nama pemanggilnya di
		// komentar dan di literal polanya, sehingga akan menghitung dirinya
		// sendiri sebagai empat pemanggil. Alat ukur tidak boleh masuk ke
		// dalam benda yang diukurnya.
		if strings.HasSuffix(p, "_test.go") &&
			!strings.HasSuffix(filepath.ToSlash(p), "/batasanpemakaian_test.go") {
			berkas = append(berkas, p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	diperiksa := 0
	for _, nama := range berkas {
		isi, err := os.ReadFile(nama)
		if err != nil {
			t.Fatal(err)
		}
		baris := strings.Split(string(isi), "\n")
		for i, b := range baris {
			// Hanya pemanggilan sungguhan - yang hasilnya ditampung - bukan
			// penyebutan namanya di komentar.
			if !strings.Contains(b, ":= skemauji.Buka()") {
				continue
			}
			diperiksa++
			jendela := strings.Join(baris[i:min(len(baris), i+10)], "\n")
			if !strings.Contains(jendela, "BolehDilewati") {
				t.Errorf("%s:%d memanggil skemauji.Buka() tanpa memeriksa BolehDilewati; "+
					"galat pagar skema uji akan dilewati diam-diam", filepath.ToSlash(nama), i+1)
			}
		}
	}
	// Delapan sejak 26-09-2026 malam: TestBentukTabelBerbedaMenggagalkanMigrasi
	// (butir x) dan services/pendaftaran_db_test.go (tiket 02) masing-masing
	// membuka koneksi uji sendiri.
	//
	// SEMBILAN sejak 28-09-2026: imageid_db_test.go mengadu rumus IMAGEID
	// dengan STANDARD_HASH Oracle, dan ia membuka koneksinya sendiri sebab
	// yang diuji bukan tabel mana pun melainkan sebuah fungsi.
	//
	// SEPULUH sejak 28-09-2026 sore: tco_pindah_db_test.go (tiket 01 Treaty
	// Contract Out) membuka koneksinya sendiri untuk mengisi tiruan warisan
	// dan membaca tabel T_* sesudah commit.
	//
	// SEBELAS sejak tiket 02 Treaty Contract Out: handlers/tco_db_test.go
	// membuka koneksinya sendiri untuk mengisi tiruan master REINSURANCETYPE
	// dan menyalakan stub identitas pada Router-nya.
	//
	// SEPULUH lagi sejak tco4 (29-09-2026): tco_pindah_db_test.go dibuang
	// bersama migrasi data ke tabel T_* - tidak ada tabel tujuan.
	//
	// SEBELAS sejak GILIRAN-17 (OQ-N2): cermintertanggung_db_test.go membuka
	// koneksinya sendiri untuk mengisi baris sumber peserta tiruan dan
	// menghitung kecocokan nama/DOB DI DALAM Oracle (nilai tidak dibaca ke Go).
	const mau = 11
	if diperiksa != mau {
		t.Errorf("pemanggil skemauji.Buka() ditemukan %d, mau %d; "+
			"bila memang bertambah, perbarui angkanya di sini", diperiksa, mau)
	}
}

// Nama tabel lama hanya boleh muncul sebagai rujukan ke yang WARISAN.
//
// `[keputusan work owner 26-09-2026, butir v1]` mengganti nama tabel dokumen
// yang BARU menjadi T_CLAIMLF_DOCUMENT, sebab POOLDATA.DOCUMENT_CLAIM sudah ada
// dengan empat belas kolom milik kelas Pega dan 295 baris. Rename itu tuntas di
// SQL dan kode, tetapi tertinggal di lima komentar - dan komentar yang menyebut
// nama tabel yang tidak ada akan menyesatkan pembaca berikutnya, yang justru
// pembaca yang paling membutuhkan komentar itu.
//
// Dua bentuk yang SAH, dan hanya dua:
//   - "Int-DOCUMENT_CLAIM"  nama kelas Pega, bukan nama tabel; tidak berubah
//   - satu baris dengan kata "warisan"  rujukan sadar ke tabel warisan
func TestNamaTabelDokumenLamaHanyaUntukWarisan(t *testing.T) {
	var berkas []string
	// Refactor bentuk B (30-09-2026): dulu hanya `internal/repository`; kini
	// juga `inti/` dan `modul/` - nama tabel warisan tidak boleh lahir kembali
	// di modul mana pun.
	for _, akar := range []string{
		filepath.FromSlash(akarModul + "/internal/repository"),
		filepath.FromSlash(akarModul + "/inti"),
		filepath.FromSlash(akarModul + "/modul"),
	} {
		err := filepath.Walk(akar, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			n := filepath.ToSlash(p)
			// Berkas ini sendiri menyebut nama itu di komentar dan literalnya, jadi
			// ia akan menghitung dirinya sendiri. Alat ukur tidak boleh masuk ke
			// dalam benda yang diukurnya.
			//
			// ⚠️ Pengecualian ini TIDAK dipakai menyembunyikan pelanggaran: setiap
			// baris berkas ini yang menyebut nama lama ditulis agar lulus aturannya
			// sendiri - lewat frasa Int-, atau bersama kata "warisan".
			if strings.HasSuffix(n, "/batasanpemakaian_test.go") {
				return nil
			}
			if strings.HasSuffix(n, ".sql") || strings.HasSuffix(n, ".go") {
				berkas = append(berkas, p)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(berkas) == 0 {
		t.Fatal("nol berkas terbaca; pembacanya yang rusak")
	}

	const lama = "DOCUMENT_CLAIM"
	diperiksa := 0
	for _, nama := range berkas {
		isi, err := os.ReadFile(nama)
		if err != nil {
			t.Fatal(err)
		}
		for i, baris := range strings.Split(string(isi), "\n") {
			if !strings.Contains(baris, lama) {
				continue
			}
			diperiksa++
			// Hanya nama kelas Pega yang dibuang sebelum memeriksa. Nama
			// tabel BARU tidak perlu dibuang: "T_CLAIMLF_DOCUMENT" tidak
			// memuat "DOCUMENT_CLAIM" sebagai potongan teks.
			bersih := strings.ReplaceAll(baris, "Int-"+lama, "")
			if !strings.Contains(bersih, lama) {
				continue
			}
			if strings.Contains(strings.ToLower(baris), "warisan") {
				continue
			}
			t.Errorf("%s:%d menyebut %s sebagai tabel BARU; namanya kini "+
				"T_CLAIMLF_DOCUMENT (butir v1): %s",
				filepath.ToSlash(nama), i+1, lama, strings.TrimSpace(baris))
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol baris memuat nama itu; pembacanya yang rusak, bukan kodenya")
	}
	t.Logf("%d baris memuat nama lama, seluruhnya sah", diperiksa)
}

// Jalur tulis baris datar warisan memanggil pagar nilainya.
//
// ⛔ Kenapa ini dijaga terpisah dari pagar itu sendiri: menguji
// PeriksaNilaiWarisan membuktikan fungsinya benar, BUKAN bahwa ia dipasang.
// Saat pagar itu ditulis (butir s1), mencabut pemanggilannya dari Simpan tidak
// menggagalkan satu test pun - fungsinya tetap lulus sendirian, dan jalur tulis
// diam-diam kembali mengirim teks ke kolom NUMBER. Test ini menutup selisih
// antara "ada" dan "dipakai".
// ⛔ "Satu jalur" ternyata tidak cukup: tinjauan menemukan jalur KEDUA,
// skemauji.IsiBarisLama, yang mem-bind NilaiBarisLama mentah ke tabel tiruan -
// persis jalur yang melahirkan ORA-01722 pada ronde 5. Test ini karena itu
// memeriksa SELURUH penulis dan mengunci cacahnya, supaya jalur ketiga tidak
// dapat muncul tanpa pagar.
//
// ⚠️ Batasnya, dikatakan terus terang: ia membaca TEKS, jadi menulis
// `_ = PeriksaNilaiWarisan(b)` - memanggil lalu membuang galatnya - akan
// meluluskannya. Yang dijaga "dipanggil atau tidak", bukan "dipakai benar".
func TestSetiapJalurTulisWarisanDipagari(t *testing.T) {
	// Dua penanda penulis: bind medan satu per satu (Simpan), dan bind lewat
	// daftar nilai (IsiBarisLama).
	penanda := []string{"BarisLamaDari(p)", "NilaiBarisLama(b)"}
	ditemukan := 0
	for nama, isi := range berkasGoSelainTest(t) {
		for _, tanda := range penanda {
			if !strings.Contains(isi, tanda) {
				continue
			}
			ditemukan++
			if !strings.Contains(isi, "PeriksaNilaiWarisan(b)") {
				t.Errorf("%s menulis baris datar warisan lewat %s tanpa memanggil "+
					"PeriksaNilaiWarisan; nilai bukan-angka akan sampai ke Oracle "+
					"sebagai ORA-01722 (butir s1)", nama, tanda)
			}
		}
	}
	const mau = 2
	if ditemukan != mau {
		t.Errorf("penulis baris datar warisan ditemukan %d, mau %d; bila memang "+
			"bertambah, pagarnya dipasang dulu lalu angka ini diperbarui", ditemukan, mau)
	}
}

// Jalur penghapusan tidak boleh menyentuh tabel peserta polis sumbernya.
//
// Tiket 15 AC: "Menghapus klaim tidak menyentuh M_LIFE_PREMIUM_DETAIL - ia
// hanya dibaca sebagai sumber snapshot peserta", dan "tidak menghapus peserta
// di premium list sumbernya". Tanpa penjaga ini kedua AC itu benar secara
// kebetulan, dan kebetulan tidak bertahan.
func TestJalurHapusTidakMenyentuhTabelSumber(t *testing.T) {
	const sumber = "M_LIFE_PREMIUM_DETAIL"
	berkasHapus := []string{"pohonklaim.go", "hapus.go"}
	diperiksa := 0
	for nama, isi := range berkasGoSelainTest(t) {
		dasar := filepath.Base(nama)
		cocok := false
		for _, b := range berkasHapus {
			if dasar == b {
				cocok = true
			}
		}
		if !cocok {
			continue
		}
		diperiksa++
		for _, baris := range strings.Split(isi, "\n") {
			if !strings.Contains(baris, sumber) {
				continue
			}
			// Menyebutnya di komentar adalah penjelasan, bukan sentuhan.
			if strings.HasPrefix(strings.TrimSpace(baris), "//") {
				continue
			}
			t.Errorf("%s menyentuh %s di jalur hapus: %s",
				nama, sumber, strings.TrimSpace(baris))
		}
	}
	if diperiksa < 2 {
		t.Fatalf("hanya %d berkas jalur hapus terbaca; pembacanya yang rusak", diperiksa)
	}
}
