package repository

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
	"regexp"
	"sort"
	"strings"
	"testing"
)

// letakStruktur menunjuk dokumen STRUKTUR dari folder paket ini.
const letakStruktur = "../../../.scratch/claim-life/STRUKTUR-TABEL-CLAIM-LIFE.md"

// tabelDikecualikan mendaftar tabel yang STRUKTUR sengaja tidak memuat
// kolomnya, beserta sebabnya.
// ⭐ Kosong sejak 26-09-2026. T_CLAIMLF_DOCUMENT dulu dikecualikan dengan
// alasan "daftar kolomnya tidak dapat diturunkan"; sejak butir ad kolomnya
// ADA - ia keempat belas kolom tabel warisan DOCUMENT_CLAIM yang dibaca dari
// katalog, dikurangi yang memang milik Pega. Pengecualian yang alasannya sudah
// tidak berlaku adalah lubang, bukan keringanan.
var tabelDikecualikan = map[string]string{}

// namaTabelBeda memetakan nama tabel di STRUKTUR ke nama yang dipakai DDL.
var namaTabelBeda = map[string]string{
	// Keputusan work owner 26 September 2026 (brief ronde 2 bab 2j): nama
	// STRUKTUR berukuran 36 byte, dan Oracle di bawah 12.2 menolak pengenal
	// lebih dari 30 byte dengan ORA-00972. Nama DDL berukuran 29 byte.
	"T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO": "T_CLAIMLF_ADJ_SPREADING_RETRO",
}

// namaKolomBeda memetakan nama kolom di STRUKTUR ke nama yang dipakai DDL.
var namaKolomBeda = map[string]string{
	// Keputusan work owner 26 September 2026 (brief ronde 2 bab 2j): 33 dan 35
	// byte di STRUKTUR, keduanya melewati batas 30 byte. Nama penggantinya
	// bukan karangan - STRUKTUR sendiri mencatat keduanya berasal dari korpus
	// RETRO_VALUATION_BEGIN_DATE dan RETRO_VALUATION_EXPIRED_DATE.
	"RETROCESSION_VALUATION_BEGIN_DATE":   "RETRO_VALUATION_BEGIN_DATE",
	"RETROCESSION_VALUATION_EXPIRED_DATE": "RETRO_VALUATION_EXPIRED_DATE",
}

var (
	polaJudulTabel = regexp.MustCompile(`^##\s+([A-Z][A-Z0-9_]+)\s*$`)
	polaSelKolom   = regexp.MustCompile("^`([A-Z][A-Z0-9_]*)`$")
)

// kolomMenurutStruktur membaca dokumen STRUKTUR menjadi peta tabel -> kolom.
func kolomMenurutStruktur(t *testing.T) map[string][]string {
	t.Helper()
	isi, err := os.ReadFile(filepath.FromSlash(letakStruktur))
	if err != nil {
		// Sengaja gagal, bukan melewati. Dokumen ini bagian dari repositori
		// yang sama; bila ia hilang, test inilah yang harus memberitahu.
		t.Fatalf("dokumen STRUKTUR tidak terbaca di %s: %v", letakStruktur, err)
	}
	hasil := map[string][]string{}
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
		if strings.HasPrefix(baris, "## ") {
			kini = ""
			continue
		}
		if kini == "" || !strings.HasPrefix(baris, "|") {
			continue
		}
		sel := strings.Split(strings.Trim(strings.TrimSpace(baris), "|"), "|")
		if k := polaSelKolom.FindStringSubmatch(strings.TrimSpace(sel[0])); k != nil {
			hasil[kini] = append(hasil[kini], k[1])
		}
	}
	return hasil
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
	langkah, err := daftarMigrasi(false)
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range langkah {
		for _, p := range m.Pernyataan {
			if nama, kolom := KolomCreateTable(p); nama != "" {
				hasil[nama] = kolom
				continue
			}
			// Kolom yang lahir di ALTER ikut dihitung: tanpa itu, langkah
			// migrasi lanjutan menambah kolom yang tak terlihat penjaga
			// mana pun.
			if nama, kolom := KolomAlterTambah(p); nama != "" {
				hasil[nama] = append(hasil[nama], kolom...)
			}
		}
	}
	return hasil
}

// Setiap tabel di DDL memuat persis kolom yang didaftar STRUKTUR.
//
// Inilah test yang seharusnya sudah ada di ronde 1. Ia menangkap dua arah
// sekaligus: kolom STRUKTUR yang belum dibuat, dan kolom DDL yang tidak pernah
// diminta siapa pun.
func TestKolomDDLCocokDenganStruktur(t *testing.T) {
	struktur := kolomMenurutStruktur(t)
	ddl := kolomMenurutDDL(t)

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
	case t == "date":
		return "tanggal"
	case strings.Contains(t, "desimal"), strings.Contains(t, "bulat"),
		strings.Contains(t, "angka"):
		return "angka"
	default:
		return "lain:" + t
	}
}

// tipeMenurutStruktur membaca kolom "Tipe" dokumen STRUKTUR.
func tipeMenurutStruktur(t *testing.T) map[string]map[string]string {
	t.Helper()
	isi, err := os.ReadFile(filepath.FromSlash(letakStruktur))
	if err != nil {
		t.Fatalf("dokumen STRUKTUR tidak terbaca: %v", err)
	}
	hasil := map[string]map[string]string{}
	var kini string
	for _, baris := range strings.Split(string(isi), "\n") {
		baris = strings.TrimRight(baris, "\r")
		if j := polaJudulTabel.FindStringSubmatch(baris); j != nil {
			kini = j[1]
			hasil[kini] = map[string]string{}
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
	return hasil
}

// tipeMenurutDDL membaca tipe tiap kolom dari berkas migrasi.
func tipeMenurutDDL(t *testing.T) map[string]map[string]string {
	t.Helper()
	polaKolomTipe := regexp.MustCompile(`^([A-Z][A-Z0-9_]*)\s+([A-Z0-9_]+(?:\([^)]*\))?)`)
	hasil := map[string]map[string]string{}
	for _, isi := range seluruhSQL(t, false) {
		for _, m := range polaCreateTabel.FindAllStringSubmatch(isi, -1) {
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
	return hasil
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
	"NUMBER(5)":    "AGE, umur peserta dalam tahun",
	"NUMBER(19)":   "T_CLAIMLF_DOCUMENT.ID, identitas dari sequence",
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
