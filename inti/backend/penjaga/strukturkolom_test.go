package penjaga

// Pembanding daftar kolom: DDL migrasi lawan dokumen STRUKTUR - TANPA Oracle.
//
// Untuk apa berkas ini: ronde 1 tiket 14 melewatkan lima kolom yang diminta
// STRUKTUR-TABEL-CLAIM-LIFE.md tanpa satu pun test yang menyadarinya. Berkas ini
// ada supaya selisih semacam itu tidak pernah lolos lagi: ia membaca dokumen
// STRUKTUR apa adanya, membaca DDL yang ditanam ke biner, lalu membandingkan
// daftar kolomnya tabel demi tabel.
//
// Dibaca sesudah: migrasi_test.go.
//
// Istilah:
//   - DDL       : pernyataan SQL yang membentuk tabel (CREATE TABLE ...).
//   - penyimpangan : tempat DDL memang sengaja berbeda dari STRUKTUR. Setiap
//                 penyimpangan didaftar di bawah beserta alasannya, dan test
//                 ini menolak penyimpangan yang tidak terpakai - supaya daftar
//                 ini tidak menyimpan alasan untuk sesuatu yang sudah berubah.

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
)

// letakStrukturInti - dokumen STRUKTUR tabel lintas modul milik `inti`
// (900-949): M_NAV_MENU, 30-09-2026. Sejak struktur tim satu folder per modul
// di `inti/docs/` (dokumen lintas modul, milik tim inti; dulu `docs/bersama/`,
// sebelumnya `.scratch/inti/` - dipindah ke dalam APP_RNM saat repo dipangkas
// hanya berisi aplikasi, 1 Oktober 2026).
const letakStrukturInti = "../../docs/STRUKTUR-TABEL-INTI.md"

// letakStruktur menunjuk SELURUH dokumen STRUKTUR yang MENGIKAT:
// `docs/STRUKTUR-TABEL-*.md` setiap modul TERDAFTAR (punya `backend/modul.go`),
// ditambah dokumen `inti`.
//
// ⚠️ Dokumen STRUKTUR modul KERANGKA (belum dimigrasi, tanpa `backend/modul.go`)
// adalah RANCANGAN: ia boleh menggambarkan tabel dan kolom yang belum dibuat
// migrasi mana pun (mis. kolom Komite rancangan Komite Claim Prop). Ia mulai
// mengikat pada saat modulnya mendapat `backend/modul.go` - pada saat itulah
// DDL-nya wajib cocok, dan kesepakatan atas tabel bersama ikut diperiksa.
//
// ⛔ Struktur tim satu folder per modul (30-09-2026): dokumennya DITEMUKAN di
// folder setiap modul, tidak didaftar di sini - dulu "menambah modul berarti
// menambah dokumennya DI SINI", yaitu menyunting berkas milik tim inti. Yang
// tetap WAJIB: setiap tabel yang dibuat migrasi harus tercatat di dokumen
// STRUKTUR mana pun (`TestKolomDDLCocokDenganStruktur` menolak tabel DDL tanpa
// dokumen), dan setiap modul yang bermigrasi harus punya sekurangnya satu
// dokumen (diperiksa di sini).
//
// ⚠️ `T_GENERAL_KOMITE` dan `T_KOMITE_KOMITELIST` disebut DUA dokumen
// sekaligus. Itu bukan kesalahan - keduanya memang batas antara dua konteks -
// tetapi ia menuntut penjaga sendiri:
// `TestDokumenSTRUKTURSepakatAtasTabelBersama`.
func letakStruktur(t *testing.T) []string {
	t.Helper()
	cocok, err := filepath.Glob(filepath.Join(akarAplikasi, "modul", "*", "docs", "STRUKTUR-TABEL-*.md"))
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(cocok)
	var mengikat []string
	berdokumen := map[string]bool{}
	for _, c := range cocok {
		m := pemilikJalur(c)
		if _, err := os.Stat(filepath.Join(akarAplikasi, "modul", m, "backend", "modul.go")); err != nil {
			continue // modul kerangka: rancangan, belum mengikat
		}
		mengikat = append(mengikat, c)
		berdokumen[m] = true
	}
	for _, jalur := range berkasMigrasi.asal {
		if m := pemilikJalur(jalur); m != "inti" && !berdokumen[m] {
			t.Fatalf("modul %s bermigrasi tetapi tanpa docs/STRUKTUR-TABEL-*.md", m)
		}
	}
	if len(mengikat) < 4 {
		t.Fatalf("hanya %d dokumen STRUKTUR modul terdaftar ditemukan; pembacanya yang rusak", len(mengikat))
	}
	return append(mengikat, filepath.FromSlash(letakStrukturInti))
}

// Tabel yang LEBIH DARI SATU dokumen gambarkan kini DITURUNKAN dari dokumennya
// (`TestDokumenSTRUKTURSepakatAtasTabelBersama`), tidak didaftar: daftar
// bernama tabel satu modul adalah daftar yang harus disunting modul itu.
// tabelDikecualikan mendaftar tabel yang STRUKTUR sengaja tidak memuat
// kolomnya, beserta sebabnya.
// ⭐ Kosong sejak 26-09-2026. T_CLAIMLF_DOCUMENT dulu dikecualikan dengan
// alasan "daftar kolomnya tidak dapat diturunkan"; sejak butir ad kolomnya
// ADA - ia keempat belas kolom tabel warisan DOCUMENT_CLAIM yang dibaca dari
// katalog, dikurangi yang memang milik Pega. Pengecualian yang alasannya sudah
// tidak berlaku adalah lubang, bukan keringanan.
var tabelDikecualikan = map[string]string{}

// tabelBukanMilikKita (`modulmd_test.go`) mendaftar tabel yang STRUKTUR
// gambarkan tetapi yang SENGAJA tidak kita buat, beserta alasannya - dari bab
// "Tabel warisan: dibaca, tidak dibuat" `MODUL.md` setiap modul.
//
// ⛔ MEKANISME YANG BERBEDA dari tabelDikecualikan, dan bedanya penting.
// `tabelDikecualikan` berarti "tabelnya kita buat, kolomnya saja yang tidak
// dibandingkan" - dan `TestTabelDikecualikanTetapDibuat` menegakkannya.
// Pernyataan tabel warisan berarti "tabelnya BUKAN milik kita": ia lahir di
// sistem lama, kita hanya membacanya. Memakai mekanisme pertama untuk maksud
// kedua akan membuat penjaga itu menuntut kita membuat tabel orang lain.
//
// ⚠️ Arah sebaliknya dijaga pula: bila DDL kelak MEMBUAT salah satunya,
// TestTabelBukanMilikKitaTidakDibuat berbunyi - sebab yang berubah saat itu
// adalah kepemilikannya, dan itu keputusan work owner.
// Nama tabel dan kolom STRUKTUR yang SENGAJA ditulis lain di DDL dinyatakan di
// bab "Nama STRUKTUR yang berbeda di DDL" `MODUL.md` modul pemiliknya
// (`namaStrukturBeda`, `modulmd_test.go`) - mis. batas 30 byte pengenal Oracle
// di bawah 12.2 (ORA-00972).

var (
	polaJudulTabel = regexp.MustCompile(`^##\s+([A-Z][A-Z0-9_]+)\s*$`)
	polaSelKolom   = regexp.MustCompile("^`([A-Z][A-Z0-9_]*)`$")
)

// kolomMenurutStruktur membaca dokumen STRUKTUR menjadi peta tabel -> kolom.
func kolomMenurutStruktur(t *testing.T) map[string][]string {
	t.Helper()
	hasil := map[string][]string{}
	for _, letak := range letakStruktur(t) {
		bacaSatuStruktur(t, letak, hasil)
	}
	return hasil
}

// bacaSatuStruktur menambahkan kolom satu dokumen STRUKTUR ke peta bersama.
//
// ⚠️ Tabel yang disebut LEBIH DARI SATU dokumen tidak digabung diam-diam:
// nama kolom yang sudah ada dilewati, dan kesepakatan antardokumen diperiksa
// `TestDokumenSTRUKTURSepakatAtasTabelBersama`. Menggabungkan kolomnya akan
// menyembunyikan dokumen mana yang sudah usang.
func bacaSatuStruktur(t *testing.T, letak string, hasil map[string][]string) {
	t.Helper()
	isi, err := os.ReadFile(filepath.FromSlash(letak))
	if err != nil {
		// Sengaja gagal, bukan melewati. Dokumen ini bagian dari repositori
		// yang sama; bila ia hilang, test inilah yang harus memberitahu.
		t.Fatalf("dokumen STRUKTUR tidak terbaca di %s: %v", letak, err)
	}
	var kini string
	for _, baris := range strings.Split(string(isi), "\n") {
		baris = strings.TrimRight(baris, "\r")
		if j := polaJudulTabel.FindStringSubmatch(baris); j != nil {
			kini = j[1]
			hasil[kini] = nil
			continue
		}
		// ⛔ Judul BUKAN-tabel mengakhiri bab tabel sebelumnya.
		//
		// Tanpa ini, tiap blok ralat bertanggal yang ditambahkan di akhir
		// dokumen - dan seluruhnya berjudul bukan-nama-tabel - barisnya
		// diatribusikan ke tabel terakhir yang kebetulan disebut. Ronde 5 dan
		// 6 menambahkan tiga blok begitu, dan kolom di dalamnya diam-diam
		// dihitung sebagai kolom tabel lain. Baru terlihat saat satu
		// pengecualian dicabut.
		// ⛔ JUDUL TINGKAT MANA PUN mengakhiri bab tabel - bukan hanya "## ".
		//
		// Ronde sebelumnya hanya memeriksa "## ", dan itu menambal separuh
		// lubang: dokumen STRUKTUR Komite memuat dua bab ber-"### " yang
		// TABELNYA berisi nama kolom - satu daftar ganti-nama dan satu daftar
		// penutup. Keduanya akan diatribusikan ke tabel terakhir yang
		// kebetulan disebut, sehingga `T_KOMITE_KOMITELIST` tampak punya 26
		// kolom padahal sembilan. Diperiksa 27-09-2026 saat dokumen ketiga
		// hendak ditambahkan.
		if strings.HasPrefix(baris, "#") {
			kini = ""
			continue
		}
		if kini == "" || !strings.HasPrefix(baris, "|") {
			continue
		}
		sel := strings.Split(strings.Trim(strings.TrimSpace(baris), "|"), "|")
		if k := polaSelKolom.FindStringSubmatch(strings.TrimSpace(sel[0])); k != nil {
			// ⛔ Tabel yang disebut DUA dokumen tidak digabung diam-diam:
			// nama kolom yang sudah ada dilewati, dan kesepakatan kedua
			// dokumen diperiksa TestDokumenSTRUKTURSepakatAtasTabelBersama.
			sudah := false
			for _, n := range hasil[kini] {
				if n == k[1] {
					sudah = true
				}
			}
			if !sudah {
				hasil[kini] = append(hasil[kini], k[1])
			}
		}
	}
}

// kolomMenurutDDL membaca berkas migrasi maju menjadi peta tabel -> kolom.
//
// Pemecahnya sendiri TIDAK lagi tinggal di sini. Sejak 26-09-2026 ia adalah
// KolomCreateTable di migrasi.go, sebab pra-terbang bentuk (butir x)
// membutuhkan pemecah yang sama saat aplikasi berjalan sungguhan. Satu pemecah
// untuk keduanya berarti test ini menguji alat yang benar-benar dipakai, bukan
// kembarannya.
func kolomMenurutDDL(t *testing.T) map[string][]string {
	t.Helper()
	hasil := map[string][]string{}
	langkah, err := migrasi.Daftar(false, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range langkah {
		for _, p := range m.Pernyataan {
			if nama, kolom := migrasi.KolomCreateTable(p); nama != "" {
				hasil[nama] = kolom
				continue
			}
			// Kolom yang lahir di ALTER ikut dihitung: tanpa itu, langkah
			// migrasi lanjutan menambah kolom yang tak terlihat penjaga
			// mana pun.
			if nama, kolom := migrasi.KolomAlterTambah(p); nama != "" {
				hasil[nama] = append(hasil[nama], kolom...)
				continue
			}
			// Kolom yang DIBUANG langkah lanjutan (901 menu datar: PARENT_ID)
			// ikut dikurangi - tanpa ini STRUKTUR yang sudah benar dituduh
			// kekurangan kolom yang memang tidak ada lagi.
			if nama, kolom := migrasi.KolomAlterBuang(p); nama != "" {
				hasil[nama] = kurangiKolom(hasil[nama], kolom)
			}
		}
	}
	return hasil
}

// kurangiKolom membuang `buang` dari `kolom`, urutan sisanya tetap.
func kurangiKolom(kolom, buang []string) []string {
	var sisa []string
	for _, k := range kolom {
		dibuang := false
		for _, b := range buang {
			dibuang = dibuang || strings.EqualFold(k, b)
		}
		if !dibuang {
			sisa = append(sisa, k)
		}
	}
	return sisa
}

// Setiap tabel di DDL memuat persis kolom yang didaftar STRUKTUR.
//
// Inilah test yang seharusnya sudah ada di ronde 1. Ia menangkap dua arah
// sekaligus: kolom STRUKTUR yang belum dibuat, dan kolom DDL yang tidak pernah
// diminta siapa pun.
func TestKolomDDLCocokDenganStruktur(t *testing.T) {
	struktur := kolomMenurutStruktur(t)
	ddl := kolomMenurutDDL(t)
	namaTabelBeda := namaStrukturBeda(t, "tabel")
	namaKolomBeda := namaStrukturBeda(t, "kolom")
	tabelBukanMilikKita := tabelBukanMilikKita(t)

	if len(ddl) == 0 {
		t.Fatal("tidak satu pun CREATE TABLE terbaca dari berkas migrasi")
	}

	terpakaiTabel := map[string]bool{}
	terpakaiKolom := map[string]bool{}
	diperiksa := 0

	for tabelStruktur, diharap := range struktur {
		if len(diharap) == 0 {
			continue // judul tanpa tabel kolom, misalnya bab pohon relasi
		}
		namaDDL := tabelStruktur
		if ganti, ada := namaTabelBeda[tabelStruktur]; ada {
			namaDDL = ganti
			terpakaiTabel[tabelStruktur] = true
		}
		if sebab, dikecualikan := tabelDikecualikan[namaDDL]; dikecualikan {
			t.Logf("%s dikecualikan: %s", namaDDL, sebab)
			continue
		}
		if sebab, bukanMilik := tabelBukanMilikKita[namaDDL]; bukanMilik {
			t.Logf("%s bukan milik kita: %s", namaDDL, sebab)
			continue
		}
		ada, punya := ddl[namaDDL]
		if !punya {
			t.Errorf("STRUKTUR memuat tabel %s, DDL tidak membuatnya", namaDDL)
			continue
		}
		diperiksa++

		punyaKolom := map[string]bool{}
		for _, k := range ada {
			punyaKolom[k] = true
		}
		diminta := map[string]bool{}
		for _, k := range diharap {
			nk := k
			if ganti, beda := namaKolomBeda[k]; beda {
				nk = ganti
				terpakaiKolom[k] = true
			}
			diminta[nk] = true
			if !punyaKolom[nk] {
				t.Errorf("%s: STRUKTUR meminta kolom %s, DDL tidak memuatnya", namaDDL, nk)
			}
		}
		for _, k := range ada {
			if !diminta[k] {
				t.Errorf("%s: DDL memuat kolom %s yang tidak diminta STRUKTUR", namaDDL, k)
			}
		}
	}

	if diperiksa == 0 {
		t.Fatal("tidak satu pun tabel terbanding; pembacanya yang rusak, bukan DDL-nya")
	}

	// Penyimpangan yang tidak terpakai berarti alasannya sudah basi.
	for n := range namaTabelBeda {
		if !terpakaiTabel[n] {
			t.Errorf("penyimpangan nama tabel %s tidak terpakai; cabut dari daftar", n)
		}
	}
	for n := range namaKolomBeda {
		if !terpakaiKolom[n] {
			t.Errorf("penyimpangan nama kolom %s tidak terpakai; cabut dari daftar", n)
		}
	}
}

// Dokumen STRUKTUR yang menggambarkan tabel yang SAMA wajib SEPAKAT atasnya.
//
// ⛔ Contoh hari ini: `T_GENERAL_KOMITE` dan `T_KOMITE_KOMITELIST` adalah BATAS antara Claim
// Life dan Komite Claim Life: yang satu menyerahkan kasus, yang lain
// memutuskan. Dokumen yang berbeda isinya berarti salah satu konteks bekerja
// dari bentuk yang sudah usang - dan bedanya baru terlihat ketika satu sisi
// menulis kolom yang sisi lain tidak baca.
//
// ⚠️ Diperiksa DARI DOKUMEN, bukan dari DDL. DDL hanya satu, jadi
// membandingkannya dengan dirinya sendiri tidak membuktikan apa pun tentang
// kesepakatan kedua konteks.
func TestDokumenSTRUKTURSepakatAtasTabelBersama(t *testing.T) {
	letak := letakStruktur(t)
	if len(letak) < 2 {
		t.Skip("hanya satu dokumen STRUKTUR; tidak ada yang dibandingkan")
	}
	perDokumen := make([]map[string][]string, 0, len(letak))
	for _, l := range letak {
		satu := map[string][]string{}
		bacaSatuStruktur(t, l, satu)
		perDokumen = append(perDokumen, satu)
	}
	// Tabel bersama DITURUNKAN: yang kolomnya digambarkan lebih dari satu dokumen.
	sebutan := map[string]int{}
	for _, d := range perDokumen {
		for tabel, kolom := range d {
			if len(kolom) > 0 {
				sebutan[tabel]++
			}
		}
	}
	var tabelBersama []string
	for tabel, n := range sebutan {
		if n > 1 {
			tabelBersama = append(tabelBersama, tabel)
		}
	}
	sort.Strings(tabelBersama)
	t.Logf("tabel yang digambarkan lebih dari satu dokumen STRUKTUR: %v", tabelBersama)
	// Prasyarat POSITIF: dulu dua tabel bernama wajib ada di kedua dokumen.
	// Sejak diturunkan, penjaga yang menemukan nol tabel bersama lulus atas
	// apa pun - padahal batas Claim Life dan Komite Claim Life ada hari ini.
	if len(tabelBersama) == 0 {
		t.Fatalf("nol tabel digambarkan lebih dari satu dokumen STRUKTUR (%d dokumen); pembacanya yang rusak", len(letak))
	}
	for _, tabel := range tabelBersama {
		var acuan []string
		var acuanDari string
		for i, d := range perDokumen {
			kolom := d[tabel]
			if len(kolom) == 0 {
				continue
			}
			if acuan == nil {
				acuan, acuanDari = kolom, letak[i]
				continue
			}
			if !reflect.DeepEqual(acuan, kolom) {
				t.Errorf("%s digambarkan BERBEDA oleh dua dokumen STRUKTUR.\n"+
					"  %s: %v\n  %s: %v\n"+
					"Keduanya batas antara dua konteks; bentuk yang berbeda berarti "+
					"salah satunya sudah usang.", tabel, acuanDari, acuan,
					letak[i], kolom)
			}
		}
	}
}

// Tabel yang dikecualikan harus memang ada di DDL - pengecualian bukan alasan
// untuk lupa membuatnya.
func TestTabelDikecualikanTetapDibuat(t *testing.T) {
	ddl := kolomMenurutDDL(t)
	for nama := range tabelDikecualikan {
		if _, ada := ddl[nama]; !ada {
			t.Errorf("%s dikecualikan dari perbandingan kolom, tetapi DDL juga tidak membuatnya", nama)
		}
	}
}

// Tabel yang BUKAN milik kita tidak boleh diam-diam dibuat migrasi.
//
// ⛔ Arah sebaliknya dari tabelBukanMilikKita. Bila DDL kelak membuat salah
// satunya, yang berubah adalah KEPEMILIKAN tabel warisan - dan itu keputusan
// work owner, bukan keputusan yang boleh menyelinap lewat satu berkas migrasi.
func TestTabelBukanMilikKitaTidakDibuat(t *testing.T) {
	ddl := kolomMenurutDDL(t)
	tabelBukanMilikKita := tabelBukanMilikKita(t)
	if len(tabelBukanMilikKita) == 0 {
		t.Skip("daftarnya kosong; tidak ada yang dijaga")
	}
	for nama, sebab := range tabelBukanMilikKita {
		if _, ada := ddl[nama]; ada {
			t.Errorf("%s dibuat migrasi, padahal ia dinyatakan bukan milik kita: %s.\n"+
				"Bila kepemilikannya memang berpindah, cabut barisnya dari bab "+
				"\"Tabel warisan\" MODUL.md modulnya - itu keputusan work owner.",
				nama, sebab)
		}
	}
}

// Tidak satu pun pengenal yang dibuat migrasi melewati 30 byte.
//
// Oracle di bawah 12.2 menolak pengenal lebih dari 30 byte dengan ORA-00972.
// Instance sasaran belum dipastikan DBA, jadi batas amanlah yang dipakai.
func TestPengenalTidakLebihDari30Byte(t *testing.T) {
	pola := []*regexp.Regexp{
		regexp.MustCompile(`(?i)CREATE\s+TABLE\s+\{skema\}\.(\w+)`),
		regexp.MustCompile(`(?i)CREATE\s+SEQUENCE\s+\{skema\}\.(\w+)`),
		regexp.MustCompile(`(?i)CREATE\s+(?:UNIQUE\s+)?INDEX\s+\{skema\}\.(\w+)`),
		regexp.MustCompile(`(?i)CONSTRAINT\s+(\w+)`),
	}
	diperiksa := 0
	for nama, isi := range seluruhSQL(t, false) {
		for _, p := range pola {
			for _, m := range p.FindAllStringSubmatch(isi, -1) {
				diperiksa++
				if n := len(m[1]); n > 30 {
					t.Errorf("%s: pengenal %s berukuran %d byte, lebih dari 30 (ORA-00972)",
						nama, m[1], n)
				}
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol pengenal terbaca; pembacanya yang rusak")
	}
}

// golonganDDL menyederhanakan tipe Oracle menjadi satu kata golongan.
func golonganDDL(tipe string) string {
	t := strings.ToUpper(tipe)
	switch {
	case strings.HasPrefix(t, "VARCHAR2"), strings.HasPrefix(t, "CHAR"),
		strings.HasPrefix(t, "CLOB"):
		return "teks"
	case strings.HasPrefix(t, "DATE"), strings.HasPrefix(t, "TIMESTAMP"):
		return "tanggal"
	case strings.HasPrefix(t, "NUMBER"), strings.HasPrefix(t, "INTEGER"):
		return "angka"
	case t == "BLOB":
		// Satu-satunya: M_TEMPLATE_FILE.ISI (Template Manager, 04-10-2026) - dijaga TestKolomUangDesimalDanNolJSON.
		return "berkas"
	default:
		return "lain:" + t
	}
}

// golonganStruktur menyederhanakan tipe yang ditulis dokumen STRUKTUR.
func golonganStruktur(tipe string) string {
	t := strings.ToLower(strings.TrimSpace(tipe))
	switch {
	case t == "teks":
		return "teks"
	case t == "date", t == "timestamp":
		// ⛔ `timestamp` ditambahkan A1 27-09-2026, bukan untuk meloloskan
		// sesuatu melainkan karena kosakatanya memang kurang: `golonganDDL`
		// sudah menggolongkan TIMESTAMP sebagai tanggal sejak awal, sedangkan
		// sisi STRUKTUR hanya mengenal DATE. Jejak audit memakai TIMESTAMP
		// dengan alasan: dua transisi dapat terjadi dalam detik yang sama, dan
		// DATE Oracle tidak dapat membedakannya.
		return "tanggal"
	case strings.Contains(t, "desimal"), strings.Contains(t, "bulat"),
		strings.Contains(t, "angka"):
		return "angka"
	case t == "berkas":
		return "berkas"
	default:
		return "lain:" + t
	}
}

// tipeMenurutStruktur membaca kolom "Tipe" dokumen STRUKTUR.
func tipeMenurutStruktur(t *testing.T) map[string]map[string]string {
	t.Helper()
	hasil := map[string]map[string]string{}
	for _, letak := range letakStruktur(t) {
		bacaSatuTipe(t, letak, hasil)
	}
	return hasil
}

// bacaSatuTipe menambahkan tipe kolom satu dokumen STRUKTUR ke peta bersama.
func bacaSatuTipe(t *testing.T, letak string, hasil map[string]map[string]string) {
	t.Helper()
	isi, err := os.ReadFile(filepath.FromSlash(letak))
	if err != nil {
		t.Fatalf("dokumen STRUKTUR tidak terbaca di %s: %v", letak, err)
	}
	var kini string
	for _, baris := range strings.Split(string(isi), "\n") {
		baris = strings.TrimRight(baris, "\r")
		if j := polaJudulTabel.FindStringSubmatch(baris); j != nil {
			kini = j[1]
			if hasil[kini] == nil {
				hasil[kini] = map[string]string{}
			}
			continue
		}
		// ⛔ JUDUL TINGKAT MANA PUN mengakhiri bab tabel. Pembaca kolom sudah
		// punya penjaga ini; pembaca TIPE tidak pernah punya, dan itu baru
		// terlihat ketika dokumen kedua ditambahkan: bab "### Daftar penutup"
		// di STRUKTUR Komite bertabel DUA kolom (nama + keterangan), sehingga
		// KETERANGAN terbaca sebagai TIPE - "penyetuju ke berapa" menjadi
		// golongan tipe. Penjaga yang membandingkan tipe lalu menuduh DDL
		// yang benar.
		if strings.HasPrefix(baris, "#") {
			kini = ""
			continue
		}
		if kini == "" || !strings.HasPrefix(baris, "|") {
			continue
		}
		sel := strings.Split(strings.Trim(strings.TrimSpace(baris), "|"), "|")
		if len(sel) < 2 {
			continue
		}
		if k := polaSelKolom.FindStringSubmatch(strings.TrimSpace(sel[0])); k != nil {
			hasil[kini][k[1]] = strings.TrimSpace(sel[1])
		}
	}
}

// tipeMenurutDDL membaca tipe tiap kolom dari berkas migrasi.
func tipeMenurutDDL(t *testing.T) map[string]map[string]string {
	t.Helper()
	polaKolomTipe := regexp.MustCompile(`^([A-Z][A-Z0-9_]*)\s+([A-Z0-9_]+(?:\([^)]*\))?)`)
	hasil := map[string]map[string]string{}
	for _, isi := range seluruhSQL(t, false) {
		for _, m := range migrasi.PolaCreateTabel.FindAllStringSubmatch(isi, -1) {
			nama, badan := strings.ToUpper(m[1]), m[2]
			kolom := map[string]string{}
			for _, b := range strings.Split(badan, "\n") {
				b = strings.TrimSpace(b)
				atas := strings.ToUpper(b)
				if b == "" || strings.HasPrefix(atas, "CONSTRAINT") ||
					strings.HasPrefix(atas, "REFERENCES") {
					continue
				}
				if k := polaKolomTipe.FindStringSubmatch(atas); k != nil {
					kolom[k[1]] = k[2]
				}
			}
			hasil[nama] = kolom
		}
	}
	// ⛔ `ALTER … MODIFY` ikut diterapkan, SESUDAH seluruh CREATE terbaca.
	//
	// Tanpa ini penjaga lebar kolom membaca bentuk saat tabel LAHIR, bukan
	// bentuknya sesudah seluruh migrasi berjalan - dan ia akan menuduh skema
	// yang justru sudah diperbaiki migrasi berikutnya. Penjaga yang menuduh
	// hal yang benar akan dilonggarkan orang, bukan dipatuhi.
	//
	// ⚠️ Urutannya dijaga `daftarMigrasi` lewat `seluruhSQL`; MODIFY yang
	// mendahului CREATE-nya sendiri tidak mungkin, sebab nomor migrasi naik.
	terapkanAlterModify(t, hasil)
	// Kolom yang DIBUANG (901 menu datar) keluar dari peta tipe - lihat
	// `kolomMenurutDDL`.
	for _, isi := range seluruhSQL(t, false) {
		for _, p := range strings.Split(isi, "\n") {
			if nama, kolom := migrasi.KolomAlterBuang(p); nama != "" {
				for _, k := range kolom {
					delete(hasil[nama], k)
				}
			}
		}
	}
	return hasil
}

// polaAwalModify mengenali baris pembuka `ALTER TABLE {skema}.X MODIFY (`.
//
// ⚠️ BERBASIS BARIS, bukan satu regex atas seluruh blok. Regex bertanda
// kurung bersarang yang mencoba menangkap `VARCHAR2(32)` di dalam `MODIFY(…)`
// berhenti di kurung tutup PERTAMA, sehingga tipenya terbaca `VARCHAR2`
// tanpa lebar - dan penjaga lebar lalu membandingkan dua nilai yang
// dua-duanya salah. Pemecah baris memakai aturan yang sama dengan pembaca
// CREATE di atas, jadi keduanya tidak dapat berbeda tafsir.
var polaAwalModify = regexp.MustCompile(
	`(?i)^ALTER\s+TABLE\s+\{skema\}\.([A-Z][A-Z0-9_]*)\s+MODIFY\s*\($`)

// modifyTabelWarisan - tabel WARISAN (tidak dibuat CREATE migrasi mana pun) yang boleh di-MODIFY, beserta sebabnya.
// Bentuk lengkapnya tidak diketahui DDL repo, jadi blok MODIFY-nya dilewati, bukan diterapkan.
var modifyTabelWarisan = map[string]string{
	"BORDEREAUX_CLAIM_AVIATION": "892: kolom spread NUMBER(10,4) warisan tak muat nominal, dilebarkan ke NUMBER(38,8) (04-10-2026)",
	"BORDEREAUX_PREMI_AVIATION": "892: kolom spread NUMBER(10,4) warisan tak muat nominal, dilebarkan ke NUMBER(38,8) (04-10-2026)",
}

// terapkanAlterModify menimpa tipe kolom dengan bentuk sesudah ALTER MODIFY.
func terapkanAlterModify(t *testing.T, hasil map[string]map[string]string) {
	t.Helper()
	polaKolomTipe := regexp.MustCompile(`^([A-Z][A-Z0-9_]*)\s+([A-Z0-9_]+(?:\([^)]*\))?)`)
	for nama, isi := range seluruhSQL(t, false) {
		if strings.Contains(nama, "_down") {
			continue // jalur mundur bukan bentuk akhir
		}
		tabel := ""
		for _, baris := range strings.Split(isi, "\n") {
			b := strings.TrimSpace(baris)
			if m := polaAwalModify.FindStringSubmatch(b); m != nil {
				tabel = strings.ToUpper(m[1])
				if hasil[tabel] == nil {
					if _, boleh := modifyTabelWarisan[tabel]; !boleh {
						t.Errorf("%s: MODIFY atas tabel %s yang tidak pernah dibuat CREATE",
							nama, tabel)
					}
					tabel = ""
				}
				continue
			}
			if tabel == "" {
				continue
			}
			if b == ")" || b == "/" {
				tabel = ""
				continue
			}
			atas := strings.ToUpper(strings.TrimSuffix(b, ","))
			if k := polaKolomTipe.FindStringSubmatch(atas); k != nil {
				hasil[tabel][k[1]] = k[2]
			}
		}
	}
}

// Golongan tipe tiap kolom DDL cocok dengan yang ditulis STRUKTUR.
//
// Pembandingnya sengaja KASAR - teks, tanggal, angka - sebab STRUKTUR memang
// tidak menetapkan presisi fisik (`[data DBA]`). Yang dijaga di sini adalah
// kesalahan yang berakibat nyata: kolom tanggal yang dibuat sebagai teks, atau
// kolom uang yang dibuat sebagai VARCHAR2.
func TestGolonganTipeDDLCocokDenganStruktur(t *testing.T) {
	struktur := tipeMenurutStruktur(t)
	ddl := tipeMenurutDDL(t)
	namaTabelBeda := namaStrukturBeda(t, "tabel")
	namaKolomBeda := namaStrukturBeda(t, "kolom")
	diperiksa := 0

	for tabelStruktur, kolom := range struktur {
		if len(kolom) == 0 {
			continue
		}
		namaDDL := tabelStruktur
		if ganti, ada := namaTabelBeda[tabelStruktur]; ada {
			namaDDL = ganti
		}
		if _, dikecualikan := tabelDikecualikan[namaDDL]; dikecualikan {
			continue
		}
		punya, ada := ddl[namaDDL]
		if !ada {
			continue // sudah dilaporkan TestKolomDDLCocokDenganStruktur
		}
		for k, tipeDoc := range kolom {
			nk := k
			if ganti, beda := namaKolomBeda[k]; beda {
				nk = ganti
			}
			tipeDDL, punyaKolom := punya[nk]
			if !punyaKolom {
				continue // sudah dilaporkan test nama kolom
			}
			diperiksa++
			mau, dapat := golonganStruktur(tipeDoc), golonganDDL(tipeDDL)
			if mau != dapat {
				t.Errorf("%s.%s: STRUKTUR menulis %q (golongan %s), DDL %q (golongan %s)",
					namaDDL, nk, tipeDoc, mau, tipeDDL, dapat)
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol kolom terbanding tipenya; pembacanya yang rusak")
	}
	t.Logf("golongan tipe dibandingkan untuk %d kolom", diperiksa)
}

// presisiSah mendaftar satu-satunya bentuk NUMBER yang boleh muncul di DDL,
// beserta sebabnya.
var presisiSah = map[string]string{
	// Persen dan rate ikut di sini atas KETETAPAN MODUL Claim Life, yang
	// ADR-U-0016 Akibat 2 serahkan kepada modul - bukan penyimpangan darinya.
	"NUMBER(38,8)": "uang, share, persen, dan rate - keputusan work owner c, 26 September 2026",
	"NUMBER(5)":    "AGE, umur peserta dalam tahun; M_NAV_MENU.URUTAN, urutan di dalam GROUPMENU",
	"NUMBER(10)":   "M_NAV_MENU.ID, identitas dari sequence - brief menu 30-09-2026 (PARENT_ID dibuang 901)",
	"NUMBER(19)":   "T_CLAIMLF_DOCUMENT.ID, identitas dari sequence",

	// Modul nbtreatyin (migrasi 320-327): diagram grilling sheet NB Treaty In
	// Prop F20 "uang · persen -> skala MINIMAL 9 desimal (P29)" dan J69 "NB:
	// 100 / jumlah baris presisi 10" - perintah work owner 04-10-2026.
	"NUMBER(38,10)": "uang dan persen NB Treaty In - diagram NB Treaty In Prop F20/J69, perintah work owner 04-10-2026",
}

// ⛔ Tidak satu pun kolom bertipe NUMBER tanpa presisi.
//
// NUMBER polos di Oracle berarti presisi arbitrer - Oracle menyimpan apa pun
// yang diberikan. Itu terdengar aman dan justru tidak: dua kolom yang menyimpan
// hal sama menjadi bertipe berbeda tanpa ada yang menyadarinya. Persis itu yang
// terjadi ronde 2 - CEDING_RETENTION bertipe NUMBER(38,8) di berkas 003 tetapi
// NUMBER polos di 004, dan tidak satu pun test menangkapnya: pembanding golongan
// sengaja kasar, dan test AC 41 hanya menuntut kata "NUMBER" ada.
func TestNolNumberTanpaPresisi(t *testing.T) {
	pola := regexp.MustCompile(`(?m)^\s+([A-Z][A-Z0-9_]*)\s+(NUMBER(?:\([0-9, ]*\))?)`)
	diperiksa := 0
	for nama, isi := range seluruhSQL(t, false) {
		for _, m := range pola.FindAllStringSubmatch(isi, -1) {
			kolom, tipe := m[1], strings.ReplaceAll(m[2], " ", "")
			diperiksa++
			if _, sah := presisiSah[tipe]; !sah {
				t.Errorf("%s kolom %s bertipe %s; yang sah hanya %v",
					nama, kolom, tipe, daftarPresisiSah())
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol kolom NUMBER terbaca; pembacanya yang rusak, bukan DDL-nya")
	}
	t.Logf("%d kolom NUMBER diperiksa presisinya", diperiksa)
}

func daftarPresisiSah() []string {
	out := make([]string, 0, len(presisiSah))
	for k := range presisiSah {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// pasanganLebarKolom memasangkan kolom penunjuk dengan kolom yang ditunjuknya.
//
// ⛔ Lebar yang berbeda TIDAK pernah gagal di Oracle - ia hanya berbohong.
// Kolom 40 karakter yang menunjuk kolom 32 karakter menjanjikan ruang yang
// tidak dapat dipakai: nilai ke-33 sampai ke-40 tidak akan pernah punya
// induk. Shared PK yang lebarnya berbeda dari induknya bukan shared PK,
// melainkan kebetulan yang sedang cocok.
var pasanganLebarKolom = []struct{ anakTabel, anakKolom, indukTabel, indukKolom string }{
	{"T_GENERAL_KOMITE", "ID", "T_WORK_CLAIM", "ID"},
	{"T_GENERAL_KOMITE", "ADJUSTMENT_ID", "T_CLAIMLF_ADJUSTMENT", "ID"},
	{"T_KOMITE_KOMITELIST", "DATA_KOMITE_ID", "T_GENERAL_KOMITE", "ID"},
	{"T_CLAIMLF_ADJUSTMENT", "KOMITE_ID", "T_WORK_CLAIM", "ID"},
	{"T_WORK_CLAIM", "COVER_KEY", "T_WORK_CLAIM", "ID"},
	// Menu per akun (Kelola User, 01-10-2026): akun dan KODE menu yang ditunjuknya.
	{"M_LOGIN_GO_MENU", "LOGIN_ID", "M_LOGIN_GO", "LOGIN_ID"},
	{"M_LOGIN_GO_MENU", "MENU_KODE", "M_NAV_MENU", "KODE"},
	{"M_LOGIN_GO_WORKBASKET", "LOGIN_ID", "M_LOGIN_GO", "LOGIN_ID"},
}

func TestLebarKolomPenunjukSamaDenganIndukNya(t *testing.T) {
	tipe := tipeMenurutDDL(t)
	diperiksa := 0
	for _, p := range pasanganLebarKolom {
		anak, ada := tipe[p.anakTabel]
		if !ada {
			t.Errorf("tabel %s tidak ada di DDL", p.anakTabel)
			continue
		}
		induk, ada := tipe[p.indukTabel]
		if !ada {
			t.Errorf("tabel %s tidak ada di DDL", p.indukTabel)
			continue
		}
		ta, tb := anak[p.anakKolom], induk[p.indukKolom]
		if ta == "" || tb == "" {
			t.Errorf("%s.%s (%q) atau %s.%s (%q) tidak terbaca dari DDL",
				p.anakTabel, p.anakKolom, ta, p.indukTabel, p.indukKolom, tb)
			continue
		}
		diperiksa++
		if ta != tb {
			t.Errorf("%s.%s bertipe %s sedangkan induknya %s.%s bertipe %s.\n"+
				"Kunci tamu antarlebar berbeda diterima Oracle tetapi menjanjikan "+
				"ruang yang tidak dapat dipakai.",
				p.anakTabel, p.anakKolom, ta, p.indukTabel, p.indukKolom, tb)
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol pasangan diperiksa - penjaga ini tidak menjaga apa pun")
	}
}
