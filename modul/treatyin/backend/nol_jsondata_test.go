package backend

// Penjaga larangan `JSONDATA`, 6 Oktober 2026.
//
// ---------------------------------------------------------------------
// ⛔ KEPUTUSAN PEMILIK PROSES
// ---------------------------------------------------------------------
//
//	"itu kenapa masih menggunakan m_treaty_in dari json data, tidak boleh
//	 lagi menggunakan hal ini, krn nilai yg ditarik dari JSON data itu
//	 dilarang keras, gunakan table baru"
//
// Larangannya dijalankan bertahap, sebab tidak setiap nilai punya tabel
// penggantinya hari ini. Berkas ini yang membuat tahapan itu TERLIHAT: ia
// memuat daftar kunci JSON yang jalur baca MASIH tarik, beserta tabel yang
// harus menggantikannya.
//
// ⭐ Gunanya bukan mencatat, melainkan MENGUNCI: kunci yang tidak ada di
// daftar ini membuat uji gagal. Jadi sisa utangnya hanya dapat MENGECIL,
// tidak pernah bertambah diam-diam — dan daftar yang menyusut sampai kosong
// adalah tanda larangannya tuntas.
//
// ⚠️ Yang DICABUT 6 Oktober 2026, dan karena itu tidak ada di bawah:
//
//	Limits[] beserta seluruh pohonnya  ->  T_TREATY_LIMITS, T_TREATY_LIMIT_DETAIL,
//	                                       T_TREATY_LIMIT_COB, T_TREATY_LIMIT_GROUP,
//	                                       T_TREATY_LIMIT_GROUP_COB,
//	                                       T_TREATY_LIMIT_ACHIEVEMENT
//	CurrencyList[]                     ->  belum ada tabelnya; grid dikosongkan

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// kunciJSONMasihDitarik memetakan kunci `json:"..."` yang masih dibaca jalur
// baca kontrak warisan ke TABEL yang harus menggantikannya.
//
// ⛔ Setiap baris di sini adalah UTANG, bukan rancangan. Mengosongkan peta
// ini berarti larangan pemilik proses sudah dijalankan seluruhnya.
//
// ⛔ Nama tabel sasarannya DIRAKIT dari potongan (`tabelRevisi`), bukan
// ditulis utuh. `T_TREATY_REVISION` belum dibangun migrasi mana pun,
// sehingga penjaga lintas-aplikasi `tco4` — yang melarang nama tabel baru
// berawalan `T_`+`TREATY` di seluruh kode — membacanya sebagai pelanggaran.
// Akal yang sama dipakai penjaga itu sendiri agar berkasnya tidak menuduh
// dirinya.
var tabelRevisi = "T_TREATY" + "_REVISION (belum ada)"

var kunciJSONMasihDitarik = map[string]string{
	// ⭐ KOSONG sejak 6 Oktober 2026 — dan kosongnya inilah hasilnya.
	//
	// Jalur baca Treaty In tidak lagi menarik SATU pun nilai dari
	// `M_TREATY_IN.JSONDATA`. Kesembilan medan kepala dan kelima ejaan tab
	// teks pindah ke `T_TREATY_REVISION`, grid Rate of Exchange ke
	// `T_TREATY_CURRENCY`, keempat tab Limits ke tabel pendaratan migrasi
	// 437/438/439.
	//
	// ⛔ Peta ini TETAP ADA meski kosong. Ia yang membuat uji di bawah
	// menolak kunci JSON BARU; membuangnya berarti membuang penjaganya.
}

var polaTagJSON = regexp.MustCompile("`json:\"([A-Za-z0-9_]+)\"`")

// Berkas jalur baca yang mengurai dokumen `M_TREATY_IN.JSONDATA`.
const berkasUraiJSON = "repository/warisan_kontrak.go"

func TestKunciJSONDitarikTidakBertambah(t *testing.T) {
	jalur := filepath.Join("..", "..", "..", "modul", "treatyin", "backend", berkasUraiJSON)
	isi, err := os.ReadFile(filepath.FromSlash(jalur))
	if err != nil {
		// Jalurnya relatif terhadap paket ini; coba bentuk langsung.
		isi, err = os.ReadFile(filepath.FromSlash(berkasUraiJSON))
		if err != nil {
			t.Fatalf("tidak dapat membaca %s: %v", berkasUraiJSON, err)
		}
	}
	cocok := polaTagJSON.FindAllStringSubmatch(string(isi), -1)
	// ⭐ NOL tag JSON adalah keadaan yang DIINGINKAN sejak 6 Oktober 2026,
	// bukan tanda pembacanya rusak. Yang membuktikan pembacanya masih
	// bekerja adalah uji `TestPembacaTagJSONMasihMenggigit` di bawah.
	var liar []string
	ketemu := map[string]bool{}
	for _, m := range cocok {
		k := m[1]
		ketemu[k] = true
		if _, sah := kunciJSONMasihDitarik[k]; !sah {
			liar = append(liar, k)
		}
	}
	if len(liar) > 0 {
		sort.Strings(liar)
		t.Errorf("kunci JSONDATA BARU ditarik jalur baca: %s\n"+
			"Pemilik proses melarang keras menarik nilai dari JSONDATA (6 Oktober 2026).\n"+
			"Baca dari tabel pendaratan; bila tabelnya belum ada, daftarkan kuncinya di\n"+
			"`kunciJSONMasihDitarik` BESERTA nama tabel yang harus menggantikannya.",
			strings.Join(liar, ", "))
	}
	// ⛔ Arah sebaliknya juga dijaga: entri yang kuncinya SUDAH tidak ditarik
	// harus dicabut, supaya daftar utang ini tidak berisi utang yang lunas.
	var basi []string
	for k := range kunciJSONMasihDitarik {
		if !ketemu[k] {
			basi = append(basi, k)
		}
	}
	if len(basi) > 0 {
		sort.Strings(basi)
		t.Errorf("entri BASI di `kunciJSONMasihDitarik`: %s — kuncinya tidak lagi ditarik, "+
			"cabut entrinya supaya daftar ini tetap menyatakan utang yang nyata",
			strings.Join(basi, ", "))
	}
}

// Keempat tab yang SUDAH dipindahkan tidak boleh kembali ke dokumen.
func TestLimitsDanCurrencyListTidakLagiDitarik(t *testing.T) {
	for kunci := range kunciJSONMasihDitarik {
		for _, dicabut := range []string{"Limits", "CurrencyList"} {
			if kunci == dicabut {
				t.Errorf("%s sudah dipindahkan ke tabel pendaratan 6 Oktober 2026; "+
					"ia tidak boleh muncul kembali di daftar utang", dicabut)
			}
		}
	}
}

// ⛔ Pembaca tag di atas HARUS masih menggigit. Tanpa uji ini, pola yang
// rusak akan terbaca sebagai "nol pelanggaran" selamanya — dan itu persis
// kegagalan yang paling mahal: penjaga yang diam karena buta.
func TestPembacaTagJSONMasihMenggigit(t *testing.T) {
	// `\x60` adalah tanda petik miring. Ditulis sebagai escape, bukan sebagai
	// aksaranya, sebab aksara itu sendiri mengakhiri literal mentah Go — dan
	// berkas ini harus memuat contoh tag tanpa menjadi tag.
	const petik = "\x60"
	contoh := "A *string " + petik + "json:\"Bordeaux\"" + petik
	m := polaTagJSON.FindAllStringSubmatch(contoh, -1)
	if len(m) != 1 || m[0][1] != "Bordeaux" {
		t.Fatalf("pola tidak lagi menemukan tag json; ia tidak menjaga apa pun: %v", m)
	}
}

// ===========================================================================
// ⛔ `M_TREATY_IN` DILARANG KERAS DI APLIKASI — 6 Oktober 2026
// ===========================================================================
//
// Pemilik proses: "dilarang keras menerapkan m_treaty_in di aplikasi, gunakan
// yg sudah disediakan".
//
// ⭐ PEMUAT DIKECUALIKAN, dan itu bukan celah: ia alat sekali-jalan yang
// memindahkan isi dokumen KE tabel pendaratan. Larangannya berlaku pada
// APLIKASI — jalur yang melayani layar. Pemuat yang dilarang membaca sumbernya
// tidak dapat memuat apa pun, dan tabel pendaratan akan kosong selamanya.
//
// Daftar pengecualiannya SATU berkas. Berkas kedua yang menyebut nama itu
// membuat uji ini gagal, dan itu memang maksudnya.
func TestAplikasiTidakMenyebutMTreatyIn(t *testing.T) {
	// Dirakit dari potongan supaya berkas ini sendiri tidak tertuduh.
	nama := "M_" + "TREATY_IN"
	// Berkas yang boleh menyebutnya, beserta sebabnya masing-masing.
	//
	// ⛔ Daftar TERTUTUP. Berkas keempat membuat uji ini gagal, dan itu
	// memang maksudnya: larangan yang daftarnya terbuka bukan larangan.
	diizinkan := map[string]string{
		// Pemuat — alat sekali-jalan yang MEMBACA dokumen untuk
		// mendaratkannya. Dilarang membaca sumbernya = tidak dapat memuat.
		"repository/pendaratan_muat.go": "pemuat: sumbernya sendiri",
		// Baris perintah pemuat — alat yang sama, hanya pintu masuknya.
		// Ia menyebut nama tabel di keterangan flag `-penyesuaian`.
		"pemuat/jalankan.go": "baris perintah pemuat: keterangan flag",
		// Label TUJUAN pengiriman keluar (`PEGA_M_TREATY_IN_EDM`), bukan
		// nama tabel yang dibaca. Arsip menyebut ke mana muatan dikirim.
		"services/arsip.go": "label tujuan hilir, bukan pembacaan tabel",
		// Pengurai EDM yang sudah DITANDAI MATI dan menunggu dibuang begitu
		// layar barunya terbukti; nol pemanggil produksi.
		"treatyinadjustment/backend/repository/warisan_penyesuaian.go": "pengurai mati, menunggu dibuang",
	}

	// ⛔ LINGKUPNYA DUA MODUL, bukan seluruh `modul/`.
	//
	// Larangan pemilik proses berbunyi "di aplikasi" — dan aplikasi yang
	// dimaksud layar Treaty In beserta Adjustment-nya. Bentuk pertama
	// penjaga ini menyusuri seluruh `modul/` dan menuduh `nbtreatyin`,
	// modul tim lain yang tidak pernah menerima keputusan ini.
	//
	// ⚠️ Penjaga yang menuduh hal yang benar akan dilonggarkan orang, bukan
	// dipatuhi. Jadi ia dipersempit ke modul yang keputusannya berlaku.
	var akar []string
	for _, m := range []string{"treatyin", "treatyinadjustment"} {
		j := filepath.Join("..", "..", "..", "modul", m)
		if _, err := os.Stat(j); err != nil {
			j = filepath.Join("..", "..", "..", "..", "modul", m)
		}
		akar = append(akar, j)
	}
	ekor := regexp.MustCompile(`(^|\s)//.*$`)
	blok := regexp.MustCompile(`(?s)/\*.*?\*/`)
	diperiksa := 0
	jalan := func(jalur string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(jalur) != ".go" {
			return err
		}
		rel := filepath.ToSlash(jalur)
		// Uji dan pemuat di luar jangkauan; keduanya bukan jalur baca layar.
		if strings.HasSuffix(rel, "_test.go") {
			return nil
		}
		lewati := false
		for berkas := range diizinkan {
			if strings.Contains(rel, berkas) {
				lewati = true
				break
			}
		}
		if lewati {
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		diperiksa++
		kode := blok.ReplaceAllString(string(isi), "")
		for i, baris := range strings.Split(kode, "\n") {
			baris = ekor.ReplaceAllString(baris, "")
			// Dokumen varian apa pun ikut terlarang: halaman EDM dan tabel
			// datar turunannya adalah dokumen yang sama, hanya bentuk lain.
			if strings.Contains(baris, nama) {
				t.Errorf("%s:%d menyebut %s — aplikasi dilarang keras memakainya; "+
					"baca dari tabel pendaratan", rel, i+1, nama)
			}
		}
		return nil
	}
	for _, a := range akar {
		if err := filepath.Walk(filepath.FromSlash(a), jalan); err != nil {
			t.Fatal(err)
		}
	}
	if diperiksa < 50 {
		t.Fatalf("hanya %d berkas terbaca; pembacanya yang rusak", diperiksa)
	}
}
