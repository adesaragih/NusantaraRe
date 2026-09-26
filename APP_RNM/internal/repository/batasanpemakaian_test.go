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

// Hapus tidak boleh dipanggil dari kode yang bukan test.
//
// Hapus adalah DELETE fisik. ADR-U-0031 menetapkan penghapusan klaim di jalur
// pengguna berupa penanda ditambah nilai balik, dan pelaksananya tiket 15.
// Selama tiket itu belum dikerjakan, satu-satunya pemakai yang sah adalah test
// yang membuktikan kaskade (AC 38).
func TestHapusTidakDipanggilDiLuarTest(t *testing.T) {
	for nama, isi := range berkasGoSelainTest(t) {
		// Definisi metodenya sendiri jelas bukan pemanggilan.
		for _, baris := range strings.Split(isi, "\n") {
			if strings.Contains(baris, "func (r *PohonKlaim) Hapus(") {
				continue
			}
			if strings.Contains(baris, ".Hapus(") {
				t.Errorf("%s memanggil .Hapus(): %s", nama, strings.TrimSpace(baris))
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
	for nama, isi := range berkasGoSelainTest(t) {
		if !strings.Contains(nama, "/internal/handlers/") {
			continue
		}
		diperiksa++
		if strings.Contains(isi, `"nusantarare/internal/repository"`) {
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
func TestNolNamaTabelTelanjangDiQuery(t *testing.T) {
	pola := regexp.MustCompile(`(?i)\b(FROM|INTO|UPDATE|JOIN)\s+([A-Za-z_{%][\w{}%.]*)`)
	diperiksa := 0
	for nama, isi := range berkasGoSelainTest(t) {
		if !strings.Contains(nama, "/internal/repository/") {
			continue
		}
		for _, m := range pola.FindAllStringSubmatch(isi, -1) {
			objek := m[2]
			diperiksa++
			switch {
			case strings.Contains(objek, "."): // berawalan skema atau {skema}
			case strings.Contains(objek, "%s"): // datang dari Qualify
			case strings.EqualFold(objek, "DUAL"): // tabel semu Oracle
			default:
				t.Errorf("%s: %s %s - nama tabel telanjang (ADR-U-0033)",
					nama, strings.ToUpper(m[1]), objek)
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
	const mau = 6
	if diperiksa != mau {
		t.Errorf("pemanggil skemauji.Buka() ditemukan %d, mau %d; "+
			"bila memang bertambah, perbarui angkanya di sini", diperiksa, mau)
	}
}
