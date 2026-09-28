package repository

// Penjaga statik tiket 13 - batas migrasi data.
//
// Pemilik: tiket 13.
//
// ⛔ RONDE PERTAMA JATUH SELURUHNYA, dan itu ditulis di sini supaya bentuk
// yang gagal tidak lahir kembali. Lima elakan dibangun review, dikompilasi,
// dan dijalankan; kelimanya hijau:
//
//	B1  variabel perantara: `l := baris["LAYER_2"]; strconv.Atoi(l)` - pola
//	    menuntut nama kolom BERSEBELAHAN dengan operatornya.
//	B2  `t, _ := r.db.Qualify("CURRENCY")` lalu `Sprintf("INSERT INTO %s …", t)`
//	    - dan ini yang paling telak: `Qualify` WAJIB dipakai (ADR-U-0033),
//	    sehingga nama tabel SELALU terpisah dari kata kerjanya. Pola lama
//	    karena itu hanya sanggup menangkap kode yang ADR-U-0033 sudah tolak.
//	B3  `SELECT p.JSONDATA` di baris berbeda dari `FROM m_product_life` - pola
//	    memakai `[^\n]{0,40}`, sedangkan setiap query di sini berbaris banyak.
//	B4  gelung tulis datar dihapus, penanda teksnya ditinggalkan - penjaga
//	    TERBALIK memeriksa ada-tidaknya TEKS, bukan ada-tidaknya PENULISAN.
//	B5  `type KolomLapisan struct{ LAYER_1..4 string }` disematkan ke
//	    `BarisLama` - penjaga hanya membaca sampai `\n}` pertama.
//
// Pelajarannya sama dengan tiket 12: pola yang menuntut bentuk penulisan
// tertentu selalu dapat dielakkan penulisan lain. Yang menggantikannya aturan
// atas NAMA - nama yang tidak boleh muncul tidak dapat disembunyikan di balik
// variabel perantara, struct sematan, maupun baris baru - ditambah satu
// pemeriksaan atas PENULISAN yang sesungguhnya.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// berkasSumberProduksi mengembalikan seluruh sumber yang BENAR-BENAR jalan:
// berkas `.go` bukan-test DAN berkas `.sql`.
//
// ⛔ Berkas `.sql` IKUT. Ronde pertama hanya membaca `.go`, padahal
// `migrations/*.sql` ditanam lewat `go:embed` dan dijalankan - sebuah
// `DELETE FROM CURRENCY` di sana tidak akan terlihat sama sekali.
func berkasSumberProduksi(t *testing.T) map[string]string {
	t.Helper()
	hasil := berkasGoSelainTest(t)
	err := filepath.Walk(filepath.FromSlash(akarModul),
		func(p string, info os.FileInfo, err error) error {
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
			if !strings.EqualFold(filepath.Ext(p), ".sql") {
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
	return hasil
}

// kolomTidakDibawa adalah kolom warisan yang SENGAJA tidak masuk `BarisLama`.
//
// Sensus 27-09-2026 `[terverifikasi]` - dua wadah, dan selisihnya daftar ini:
//
//	katalog `kolomWarisan` : 62 kolom
//	struct  `BarisLama`    : 55 medan
//	selisih                : 7
//
// ⚠️ Ini BUKAN "sensus dua cara" dalam arti CLAUDE.md bab 4a: dua wadah yang
// dikurangkan bukan dua cara menghitung satu hal, dan menyebutnya begitu di
// ronde pertama adalah trap 5. Yang benar-benar diperiksa dua arah adalah
// DAFTARNYA - tiap nama wajib ADA di katalog dan wajib TIDAK ADA di sumber
// produksi mana pun, keduanya di test ini.
//
// Audit: `grep -cE '\{"[A-Z_0-9]+", "[^"]*"\}' internal/repository/barislamakolom.go`
// untuk angka 62; pembacaan medan `type BarisLama struct` di `migrasidata.go`
// untuk angka 55.
//
// Alasan ketujuhnya di luar: model relasional baru tidak punya rumah bagi
// mereka, dan mengarang rumah berarti MENAFSIRKAN kolom yang perannya belum
// terverifikasi (`[data DBA]` OQ-001 untuk LAYER). Yang dilakukan migrasi:
// tidak menyentuhnya. Tabel datarnya tetap ada, jadi nilai lamanya utuh.
var kolomTidakDibawa = []string{
	"STS_KONVERSI", "TGL_KONVERSI",
	"LAYER_1", "LAYER_2", "LAYER_3", "LAYER_4",
	"INDEXLIST",
}

// berkasBolehMenyebutKolomTakDibawa berkunci akhiran jalur, dengan alasan.
var berkasBolehMenyebutKolomTakDibawa = map[string]string{
	"repository/barislamakolom.go": "katalog 62 kolom - kolomnya ADA di tabel, ia hanya tidak dibawa",
}

// migrasiDiLuarClaimLife menjawab apakah berkas itu migrasi milik modul lain.
//
// Batasnya NOMOR: Claim Life memakai 001-049, PremiumList Life mulai 050.
func migrasiDiLuarClaimLife(nama string) bool {
	i := strings.LastIndexAny(nama, "/\\")
	dasar := nama[i+1:]
	if len(dasar) < 4 || !strings.HasSuffix(nama, ".sql") {
		return false
	}
	return dasar[:4] >= "050_"
}

// TestKolomTakDibawaHanyaAdaDiKatalog menggantikan DUA penjaga ronde pertama.
//
// ⛔ Aturan atas NAMA, bukan atas bentuk penulisan. Nama yang tidak boleh
// muncul tidak dapat disembunyikan di balik variabel perantara (B1) maupun
// struct sematan (B5). Kedua elakan itu kini tertangkap oleh satu aturan yang
// lebih sederhana daripada dua pola yang digantikannya.
func TestKolomTakDibawaHanyaAdaDiKatalog(t *testing.T) {
	sumber := berkasSumberProduksi(t)
	diperiksa, dikecualikan := 0, 0
	for nama, isi := range sumber {
		// ⛔ LINGKUPNYA MODUL CLAIM LIFE. `kolomTidakDibawa` adalah selisih
		// antara katalog warisan dan struct `BarisLama` - keduanya milik
		// tabel KLAIM. Sejak modul PremiumList Life menambah migrasi 050+,
		// penjaga ini harus menyatakan lingkupnya: `LAYER_1`..`LAYER_4`
		// memang tidak punya rumah di tabel klaim, tetapi ia PUNYA rumah di
		// `T_PREMIUM_LIST` - `[keputusan work owner 2026-09-18]`,
		// STRUKTUR-TABEL-PREMIUMLIST-LIFE.md b98-101, bersumber korpus
		// `SaveLifeinProduction_SQL`.
		//
		// ⚠️ Penjaga yang membentang ke tabel modul lain akan menuduh kolom
		// yang tujuannya SUDAH diputuskan - dan penjaga yang menuduh hal yang
		// benar akan dilonggarkan orang, bukan dipatuhi.
		if migrasiDiLuarClaimLife(nama) {
			// ⚠️ TIDAK menambah `dikecualikan`: pencacah itu dikunci terhadap
			// peta `berkasBolehMenyebutKolomTakDibawa`, dan mencemarinya dengan
			// hitungan lingkup akan membuat penjaga peta itu berbunyi palsu.
			continue
		}
		diperiksa++
		if alasan, boleh := cocokAkhiran(berkasBolehMenyebutKolomTakDibawa, nama); boleh {
			dikecualikan++
			t.Logf("dikecualikan: %s (%s)", nama, alasan)
			continue
		}
		kode := buangKomentarSumber(nama, isi)
		for _, kolom := range kolomTidakDibawa {
			if strings.Contains(kode, kolom) {
				t.Errorf("%s menyebut %q. Ketujuh kolom ini SENGAJA tidak dibawa: "+
					"model relasional tidak punya rumah bagi mereka, dan mengarang "+
					"rumah berarti menafsirkan kolom yang perannya belum "+
					"terverifikasi. Bila work owner sudah memutuskan tujuannya, "+
					"hapus namanya dari kolomTidakDibawa beserta alasannya",
					nama, kolom)
			}
		}
	}
	if diperiksa < 30 {
		t.Fatalf("hanya %d berkas terbaca; pembacanya yang rusak", diperiksa)
	}
	if dikecualikan != len(berkasBolehMenyebutKolomTakDibawa) {
		t.Errorf("%d berkas dikecualikan, petanya memuat %d",
			dikecualikan, len(berkasBolehMenyebutKolomTakDibawa))
	}
	// ⛔ Arah kedua: katalognya TETAP memuat ketujuhnya - mereka ada di tabel
	// sungguhan. Menghapusnya dari katalog membuat DDL tiruan berbeda dari
	// tabelnya, dan migrasi yang membaca tabel itu akan meleset diam-diam.
	katalog, ada := cariBerkas(sumber, "repository/barislamakolom.go")
	if !ada {
		t.Fatal("barislamakolom.go tidak terbaca")
	}
	for _, kolom := range kolomTidakDibawa {
		if !strings.Contains(katalog, `"`+kolom+`"`) {
			t.Errorf("katalog kehilangan %q; kolomnya ADA di tabel sungguhan", kolom)
		}
	}
}

// masterYangTidakDisentuh adalah nama yang tidak boleh muncul di sumber mana
// pun.
//
// `[data DBA]` OQ-001: `RATE_LIFE` dan `PRODUCTINWARD_LIFE` adalah VIEW atas
// `JSONDATA`. `m_product_life` memegang JSON produk, dan membacanya menuntut
// persetujuan manusia - produk dibaca dari `product_life` RELASIONAL
// (`[keputusan work owner 2026-09-16]`).
//
// ⚠️ `CURRENCY` TIDAK DAPAT DIJAGA LEWAT NAMA, dan itu dinyatakan, bukan
// disembunyikan: ia juga nama KOLOM dan muncul puluhan kali di repository
// secara sah. Aturan atas nama akan menuduh seluruhnya. Yang menjaganya hanya
// kenyataan bahwa nol kode menyentuh master mana pun hari ini; bila itu
// berubah, `CURRENCY` harus ditinjau dengan tangan.
var masterYangTidakDisentuh = []string{
	"RATE_LIFE", "PRODUCTINWARD_LIFE", "m_product_life", "M_PRODUCT_LIFE",
}

// TestMasterViewTidakDisentuh - aturan atas NAMA, bukan atas bentuk SQL.
//
// ⛔ Menggantikan pola `INSERT INTO … RATE_LIFE` yang tidak pernah mungkin
// menggigit di basis kode ini: `Qualify` wajib dipakai (ADR-U-0033), sehingga
// nama tabel selalu terpisah dari kata kerjanya.
// berkasIzinViewProduk adalah SATU-SATUNYA berkas yang boleh menyebut
// `PRODUCTINWARD_LIFE`, beserta batas izinnya.
//
// ⛔ DIPERSEMPIT 28-09-2026 - `[DIPUTUSKAN, butir bh, veto work owner]`.
// Sebelumnya view ini terlarang mutlak, dan pembaca ambang butir `ba` sempat
// ditulis lalu DIBUANG ketika penjaga ini berbunyi. Penjaga itu benar saat
// itu, dan ia tetap benar untuk segala hal SELAIN yang dinamai di sini.
//
// Dasar pembukaannya: AC 38 melarang MENGURAI `JSONDATA` di aplikasi - bukan
// membaca kolom BERTIPE dari view milik basis data, yang Pega sendiri baca
// (`Claim Life/RDBList/GetProductName.xml`).
//
// ⛔ IZINNYA SEMPIT, DAN KESEMPITANNYA YANG MEMBUATNYA AMAN. Yang dikunci:
//
//	berkas   TEPAT satu
//	fungsi   TEPAT satu pembaca
//	kolom    TEPAT dua, bernama
//	tulisan  NOL - tidak ada INSERT/UPDATE/DELETE/MERGE
//
// Pembaca kedua, kolom ketiga, atau satu tulisan akan tetap berbunyi.
//
// ⚠️ Ronde pertama pembukaan ini MENGELABUI penjaga alih-alih
// mempersempitnya: namanya dirakit dari potongan (`"PRODUCTINWARD" +
// "_LIFE"`) sehingga pemindaian tidak menemukannya. Itu dibatalkan. Penjaga
// yang dikelabui sekali menjadi penjaga yang buta selamanya, sebab pembaca
// berikutnya akan menyalin caranya.
const berkasIzinViewProduk = "ambangproduk.go"

// izinViewProduk memeriksa bahwa berkas berizin itu tetap di dalam batasnya.
func izinViewProduk(t *testing.T, isi string) {
	t.Helper()
	if !strings.Contains(isi, "func (r *ProdukLife) Ambang(") {
		t.Errorf("%s tidak lagi memuat pembaca `ProdukLife.Ambang` - "+
			"izin butir bh terikat pada fungsi itu, bukan pada berkasnya",
			berkasIzinViewProduk)
	}
	// TEPAT dua kolom, dan keduanya bernama.
	for _, k := range KolomAmbangProduk {
		if !strings.Contains(isi, k) {
			t.Errorf("%s tidak membaca kolom %q", berkasIzinViewProduk, k)
		}
	}
	if len(KolomAmbangProduk) != 2 {
		t.Errorf("izin butir bh mencakup TEPAT dua kolom; daftarnya kini %d",
			len(KolomAmbangProduk))
	}
	if strings.Contains(isi, "SELECT *") {
		t.Errorf("%s memakai SELECT * atas view 40 kolom; izinnya hanya dua",
			berkasIzinViewProduk)
	}
	// ⛔ NOL TULISAN. Izinnya membaca, titik.
	for _, tulis := range []string{"INSERT", "UPDATE ", "DELETE", "MERGE"} {
		if strings.Contains(strings.ToUpper(isi), tulis) {
			t.Errorf("%s memuat %q - izin butir bh READ-ONLY",
				berkasIzinViewProduk, tulis)
		}
	}
	// ⛔ NOL JSONDATA. Itu batas yang AC 38 tarik, dan ia tidak bergeser.
	if strings.Contains(strings.ToUpper(isi), "JSONDATA") {
		t.Errorf("%s menyebut JSONDATA - AC 38 melarangnya", berkasIzinViewProduk)
	}
}

func TestMasterViewTidakDisentuh(t *testing.T) {
	diperiksa := 0
	for nama, isi := range berkasSumberProduksi(t) {
		diperiksa++
		// ⛔ Komentar dibuang lebih dulu. Menjelaskan MENGAPA sebuah master
		// tidak disentuh menuntut menyebut namanya - dan ronde pertama
		// menuduh justru penjelasan itu, di `spreading.go` dan `skemauji.go`.
		// Penjaga yang menuduh dokumentasinya sendiri akan dimatikan orang.
		kode := buangKomentarSumber(nama, isi)
		// Berkas berizin butir bh diperiksa dengan aturannya sendiri, yang
		// LEBIH ketat daripada sekadar "tidak menyebut namanya".
		if strings.HasSuffix(filepath.ToSlash(nama), "/"+berkasIzinViewProduk) {
			izinViewProduk(t, kode)
			continue
		}
		for _, master := range masterYangTidakDisentuh {
			if strings.Contains(kode, master) {
				t.Errorf("%s menyebut %q. Kedua master pertama VIEW atas JSONDATA "+
					"dan tidak diperlakukan sebagai tabel relasional; JSON produk "+
					"dibaca dari `product_life` relasional, dan membacanya langsung "+
					"menuntut persetujuan manusia", nama, master)
			}
		}
	}
	if diperiksa < 30 {
		t.Fatalf("hanya %d berkas terbaca; pembacanya yang rusak", diperiksa)
	}
}

// TestJalurTulisDatarWarisanWAJIBADA - penjaga TERBALIK.
//
// Tiket 13 memintanya terus terang: *"Test yang **menolak** penulisan
// `OS_AKSEPTASI_KLAIM_LIFE` justru **gagal**"* - hilir masih membaca dari
// sana, dan yang dibuang hanya JSON.
//
// ⛔ Yang dituntut PENULISANNYA, bukan teks penandanya. Ronde pertama
// memeriksa ada-tidaknya `BarisLamaDari(p)` dan `PeriksaNilaiWarisan(b)`, dan
// dielakkan dengan menghapus gelung `r.exec` sambil meninggalkan kedua
// penanda itu beserta `_ = datarT`.
func TestJalurTulisDatarWarisanWAJIBADA(t *testing.T) {
	isi, ada := cariBerkas(berkasSumberProduksi(t), "repository/pohonklaim.go")
	if !ada {
		t.Fatal("pohonklaim.go tidak terbaca; pembacanya yang rusak")
	}
	simpan := potongFungsi(buangKomentarGo(isi), "func (r *PohonKlaim) Simpan(")
	if simpan == "" {
		t.Fatal("fungsi Simpan tidak ditemukan; pembacanya yang rusak")
	}
	// Penulisan yang sesungguhnya: INSERT ke tabel datar, DI DALAM Simpan,
	// pada transaksi yang sama dengan tabel relasionalnya.
	if !strings.Contains(simpan, "INSERT INTO %s") || !strings.Contains(simpan, "datarT)") {
		t.Error("Simpan tidak lagi meng-INSERT ke tabel datar warisan. Jalur itu " +
			"WAJIB ADA: hilir - Arasapas - masih membaca OS_AKSEPTASI_KLAIM_LIFE, " +
			"dan yang dibuang hanya JSON (`[keputusan work owner 2026-09-16]`). " +
			"Dua transaksi pun tidak cukup: ia harus di dalam Simpan")
	}
	// Dan pagarnya tetap dipanggil sebelum bind (butir s1).
	if !strings.Contains(simpan, "PeriksaNilaiWarisan(b)") {
		t.Error("Simpan menulis baris datar tanpa PeriksaNilaiWarisan; nilai " +
			"bukan-angka akan sampai ke Oracle sebagai ORA-01722")
	}
}

// cocokAkhiran mencari entri peta yang akhiran jalurnya cocok.
func cocokAkhiran(peta map[string]string, nama string) (string, bool) {
	for akhiran, alasan := range peta {
		if strings.HasSuffix(nama, akhiran) {
			return alasan, true
		}
	}
	return "", false
}

// cariBerkas mencari satu berkas lewat akhiran jalurnya.
//
// ⚠️ Lewat AKHIRAN, bukan kunci persis: bentuk kunci peta bergantung akar
// telusurnya, dan penjaga yang pecah karena akarnya bergeser bukan penjaga.
func cariBerkas(peta map[string]string, akhiran string) (string, bool) {
	for nama, isi := range peta {
		if strings.HasSuffix(nama, akhiran) {
			return isi, true
		}
	}
	return "", false
}

// potongFungsi mengambil tubuh satu fungsi dari awalan tanda tangannya.
func potongFungsi(kode, awalan string) string {
	i := strings.Index(kode, awalan)
	if i < 0 {
		return ""
	}
	sisa := kode[i+len(awalan):]
	if j := strings.Index(sisa, "\nfunc "); j >= 0 {
		return sisa[:j]
	}
	return sisa
}

// buangKomentarGo menghapus baris komentar sebelum pencocokan.
//
// ⚠️ Batasnya dinyatakan: hanya komentar SEBARIS PENUH yang dibuang. Komentar
// blok dan komentar di ujung baris tetap tinggal.
func buangKomentarGo(isi string) string {
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

// buangKomentarSumber membuang komentar sesuai jenis berkasnya.
//
// ⚠️ Batasnya dinyatakan: untuk `.go` hanya komentar sebaris penuh, untuk
// `.sql` hanya baris berawalan `--`. Komentar blok dan komentar di ujung baris
// tetap tinggal, sehingga sebuah nama terlarang yang disembunyikan di sana
// tidak tertangkap. Itu celah yang diketahui, bukan yang tidak disadari.
func buangKomentarSumber(nama, isi string) string {
	awalan := "//"
	if strings.EqualFold(filepath.Ext(nama), ".sql") {
		awalan = "--"
	}
	var b strings.Builder
	for _, baris := range strings.Split(isi, "\n") {
		if strings.HasPrefix(strings.TrimSpace(baris), awalan) {
			continue
		}
		b.WriteString(baris)
		b.WriteString("\n")
	}
	return b.String()
}
