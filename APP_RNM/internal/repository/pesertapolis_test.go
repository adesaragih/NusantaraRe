package repository

// Penyaring peserta dan pagar query tabel 66 juta baris - TANPA Oracle.
//
// Pemilik: tiket 02.
//
// Dibaca sesudah: pesertapolis.go.

import (
	"context"
	"regexp"
	"strings"
	"testing"
)

// AC 26-30: penyaring hidup menerima yang hidup dan menolak yang batal.
//
// ⛔ Kasus NULL dan teks kosong DIPISAH dengan sengaja. Penyaring naif
// `EDMSTATUS = ”` terlihat benar dan salah total: agregat menghitung 59,1
// juta baris NULL dan NOL baris teks kosong, jadi penyaring itu membuang
// hampir seluruh tabel tanpa satu pun galat.
func TestPesertaHidupMenyaringBatalDanDelete(t *testing.T) {
	kasus := map[string]bool{
		"":        true,  // teks kosong - nol baris begini di DEV, tetap diterima
		"Old":     true,  // AC 27
		"New":     true,  // AC 27
		"Batal":   false, // AC 26
		"Delete":  false, // AC 26, soft-delete
		" Batal ": false, // berspasi tetap batal
		"Unknown": true,  // nilai tak dikenal LOLOS - lihat komentar fungsinya
	}
	for nilai, mau := range kasus {
		if got := PesertaHidup(nilai); got != mau {
			t.Errorf("PesertaHidup(%q) = %v, mau %v", nilai, got, mau)
		}
	}
	// NULL tiba sebagai teks kosong lewat sql.NullString; diuji di atas.
}

// Pencarian tanpa nomor premium list DITOLAK, bukan dijalankan.
func TestCariMenolakTanpaNomorPremiList(t *testing.T) {
	r := &PesertaPolis{}
	for _, kosong := range []string{"", "   "} {
		if _, err := r.Cari(context.Background(), kosong, 10); err == nil {
			t.Errorf("pencarian dengan nomor %q diterima; ia memindai 66 juta baris", kosong)
		}
	}
}

// ⛔ Setiap query ke tabel peserta menyaring ber-index DAN berbatas hasil.
//
// Kenapa ini test statik dan bukan test perilaku: yang dijaga bukan hasilnya
// melainkan BENTUK query-nya, dan akibat melanggarnya - pemindaian penuh atas
// 66,8 juta baris - tidak akan terlihat pada skema uji yang isinya tiga baris.
// Ia baru terlihat di produksi, saat sudah terlambat.
func TestQueryTabelPesertaSelaluBerindexDanBerbatas(t *testing.T) {
	diperiksa := 0
	for nama, isi := range berkasGoSelainTest(t) {
		// Hanya berkas yang benar-benar menyentuh tabel itu. Nama tabelnya
		// tidak muncul di teks query - ia datang lewat Qualify sebagai %s -
		// jadi yang dicari adalah konstanta pemiliknya.
		if !strings.Contains(isi, "namaTabelPeserta") {
			continue
		}
		for i := 0; ; {
			j := strings.Index(isi[i:], "SELECT")
			if j < 0 {
				break
			}
			j += i
			panjang := strings.Index(isi[j:], "`")
			if panjang < 0 {
				panjang = len(isi) - j
			}
			q := strings.ToUpper(isi[j : j+panjang])
			i = j + len("SELECT")
			diperiksa++

			// ⛔ Penyaringnya dicari SESUDAH WHERE, bukan di mana pun. Kolom
			// yang diminta SELECT hampir selalu memuat PL_NUMBER juga, jadi
			// memeriksa seluruh teks akan meluluskan query yang penyaringnya
			// dicabut - penjaga yang memberi rasa aman palsu. Ditemukan saat
			// penjaga ini diuji pada kasus buruknya.
			// Klausa WHERE saja, DIBATASI ujungnya. Tanpa batas itu,
			// "ORDER BY CERTIFICATE_NO" ikut terbaca sebagai penyaring -
			// lubang kedua yang ditemukan saat penjaga ini diuji.
			where := ""
			if k := strings.Index(q, "WHERE"); k >= 0 {
				where = q[k:]
				for _, ujung := range []string{"ORDER BY", "GROUP BY", "FETCH FIRST"} {
					if u := strings.Index(where, ujung); u >= 0 {
						where = where[:u]
					}
				}
			}
			if !strings.Contains(where, "PL_NUMBER") && !strings.Contains(where, "CERTIFICATE_NO") {
				t.Errorf("%s: query tanpa penyaring ber-index; tabel peserta berisi "+
					"66,8 juta baris dan tidak ber-index pada EDMSTATUS", nama)
			}
			if !strings.Contains(q, "FETCH FIRST") && !strings.Contains(q, "ROWNUM") {
				t.Errorf("%s: query ke tabel peserta tanpa batas hasil", nama)
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol query ke tabel peserta terbaca; pembacanya yang rusak, bukan kodenya")
	}
	t.Logf("%d query ke tabel peserta diperiksa", diperiksa)
}

// AC 29: penyaringan terjadi di SATU tempat.
func TestPenyaringPesertaHanyaSatuTempat(t *testing.T) {
	ketemu := 0
	for nama, isi := range berkasGoSelainTest(t) {
		if strings.Contains(isi, "EDMSTATUS") {
			ketemu++
			if !strings.HasSuffix(nama, "/pesertapolis.go") {
				t.Errorf("%s menyebut EDMSTATUS; penyaringnya harus satu tempat (AC 29)", nama)
			}
		}
	}
	if ketemu != 1 {
		t.Errorf("berkas yang menyebut EDMSTATUS = %d, mau 1", ketemu)
	}
}

// Pengenal work object berbentuk tetap: awalan + enam digit.
func TestRakitPengenalWork(t *testing.T) {
	kasus := map[string]string{
		"1":       AwalanKlaim + "000001",
		"42":      AwalanKlaim + "000042",
		"999999":  AwalanKlaim + "999999",
		"1000000": AwalanKlaim + "1000000", // ⛔ TIDAK dipotong: lihat komentarnya
	}
	for urut, mau := range kasus {
		if got := RakitPengenalWork(AwalanKlaim, urut); got != mau {
			t.Errorf("RakitPengenalWork(%q) = %q, mau %q", urut, got, mau)
		}
	}
	if got := RakitPengenalWork(AwalanKomite, "7"); got != AwalanKomite+"000007" {
		t.Errorf("awalan komite = %q", got)
	}
}

// AC 25: STATUS dan STATUSOLD tidak pernah dipakai sebagai penanda hidup/mati.
//
// Keduanya ada di tabel peserta dan terlihat menggoda - namanya persis seperti
// penanda keadaan. Yang menentukan peserta masih hidup bagi Claim Life adalah
// EDMSTATUS (verdict V14), dan memakai kolom lain akan menyaring dengan aturan
// yang tidak pernah ditetapkan siapa pun.
func TestStatusDanStatusOldBukanPenandaHidup(t *testing.T) {
	// \b = batas kata. Tanpa itu, "EDMSTATUS" ikut tercocok sebagai "STATUS"
	// dan penjaga ini menuduh penyaring yang justru benar.
	//
	// ⛔ Polanya sempat ditulis \\b di dalam raw string, dan di sana dua
	// backslash berarti backslash HARFIAH - bukan batas kata. Pola itu tidak
	// pernah cocok dengan apa pun, sehingga test lulus hampa sementara bab
	// Implementasi justru membanggakannya sebagai penjaga yang diperbaiki.
	// Ditemukan /code-review, dan dibuktikan dengan menjalankan kedua pola
	// berdampingan atas "WHERE STATUS = 1".
	terlarang := regexp.MustCompile(`(?i)\b(STATUSOLD|STATUS)\s*=`)
	diperiksa := 0
	for nama, isi := range berkasGoSelainTest(t) {
		if !strings.Contains(isi, "namaTabelPeserta") {
			continue
		}
		diperiksa++
		// Komentar dibuang: yang dijaga adalah query, bukan penjelasannya.
		var kode []string
		for _, b := range strings.Split(isi, "\n") {
			if strings.HasPrefix(strings.TrimSpace(b), "//") {
				continue
			}
			kode = append(kode, b)
		}
		for _, m := range terlarang.FindAllString(strings.Join(kode, "\n"), -1) {
			if strings.Contains(strings.ToUpper(m), "EDMSTATUS") {
				continue
			}
			t.Errorf("%s memakai %q sebagai penyaring peserta; yang menentukan "+
				"hidup/mati hanya EDMSTATUS (AC 25, verdict V14)", nama, strings.TrimSpace(m))
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol berkas pembaca peserta terbaca; pembacanya yang rusak")
	}
}
