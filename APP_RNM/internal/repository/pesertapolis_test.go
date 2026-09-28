package repository

// Penyaring peserta dan pagar query tabel 66 juta baris - TANPA Oracle.
//
// Pemilik: tiket 02.
//
// Dibaca sesudah: pesertapolis.go.

import (
	"context"
	"os"
	"path/filepath"
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
		if _, err := r.Cari(context.Background(), kosong, "", "", 10); err == nil {
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
			// ⚠️ Jendela bacanya adalah SELURUH FUNGSI yang memuat SELECT
			// itu, bukan satu literal, dan ia dimulai dari kepala fungsinya.
			// Sejak SQL pencarian dirakit (sqlCariPeserta), penyaring dan
			// batasnya TIDAK berada di literal yang sama dengan kata SELECT:
			// penyaringnya tersusun di slice syarat SEBELUM SELECT ditulis.
			// Jendela yang berhenti di backtick pertama menuduh query yang
			// justru berpagar lengkap - dan penjaga yang menuduh hal yang benar
			// akan dilonggarkan orang, bukan dipatuhi.
			kepala := strings.LastIndex(isi[:j], "\nfunc ")
			if kepala < 0 {
				kepala = 0
			}
			ekor := strings.Index(isi[j:], "\n}")
			if ekor < 0 {
				ekor = len(isi) - j
			}
			jendela := isi[kepala : j+ekor]
			dirakit := strings.Contains(jendela, "` +")

			// ⛔ Jendela LEBAR itu hanya untuk SQL RAKITAN. Untuk SQL
			// inline, jaminannya memang berada di literal yang sama, dan
			// memperlebar jendelanya di situ justru MELONGGARKAN penjaga:
			// sebuah "PL_NUMBER =" di query LAIN dalam fungsi yang sama
			// akan menutupi query yang penyaringnya dicabut. Karena itu SQL
			// inline tetap diperiksa pada klausa WHERE-nya sendiri, seketat
			// sebelum berkas ini berubah.
			q := strings.ToUpper(tanpaKomentar(jendela))
			if !dirakit {
				lit := isi[j:]
				if b := strings.Index(lit, "`"); b >= 0 {
					lit = lit[:b]
				}
				q = strings.ToUpper(lit)
				// Klausa WHERE saja, DIBATASI ujungnya - tanpa batas itu
				// "ORDER BY CERTIFICATE_NO" ikut terbaca sebagai penyaring.
				if k := strings.Index(q, "WHERE"); k >= 0 {
					w := q[k:]
					for _, ujung := range []string{"ORDER BY", "GROUP BY", "FETCH FIRST"} {
						if u := strings.Index(w, ujung); u >= 0 {
							w = w[:u]
						}
					}
					// Batas hasil tetap dicari di SELURUH literal, bukan di WHERE.
					if !strings.Contains(q, "FETCH FIRST") &&
						!strings.Contains(q, "ROWNUM") {
						t.Errorf("%s: query inline tanpa batas hasil", nama)
					}
					q = w + " FETCH FIRST"
				}
			}
			i = j + len("SELECT")
			diperiksa++

			// ⛔ Kolomnya harus DIBANDINGKAN, bukan sekadar disebut.
			// Daftar SELECT hampir selalu memuat PL_NUMBER juga, sehingga
			// memeriksa penyebutannya saja akan meluluskan query yang
			// penyaringnya dicabut - penjaga yang memberi rasa aman palsu.
			// Yang dicari karena itu nama kolom yang diikuti operator.
			if !dibandingkan(q, "PL_NUMBER") &&
				!dibandingkan(q, "CERTIFICATE_NO") {
				t.Errorf("%s: query tanpa penyaring ber-index; tabel peserta "+
					"berisi 66,8 juta baris dan tidak ber-index pada EDMSTATUS", nama)
			}
			if !strings.Contains(q, "FETCH FIRST") &&
				!strings.Contains(q, "ROWNUM") {
				t.Errorf("%s: query ke tabel peserta tanpa batas hasil", nama)
			}

			// ⛔ BATAS PENJAGA INI, dinyatakan supaya tidak menenangkan
			// secara palsu: untuk SQL yang DIRAKIT dari potongan, pemindaian
			// teks tidak dapat membedakan penyaring WAJIB dari penyaring
			// BERSYARAT. Di sqlCariPeserta, 'CERTIFICATE_NO LIKE' yang hanya
			// terpasang bila kotaknya terisi tetap terbaca sebagai penyaring,
			// sehingga hilangnya 'PL_NUMBER = :1' yang wajib TIDAK terlihat.
			// Dibuktikan: mencabut PL_NUMBER membuat penjaga ini tetap hijau.
			//
			// Karena itu SQL rakitan wajib punya uji yang MENYEBUT nama
			// fungsinya - di sanalah bentuk akhirnya diperiksa atas seluruh
			// kombinasi masukan. Mencabut PL_NUMBER atau FETCH FIRST membuat
			// TestSQLCariPesertaSelaluBerpagar merah; keduanya sudah diuji.
			if dirakit {
				fn := namaFungsi(isi[kepala:])
				if fn == "" {
					t.Errorf("%s: SQL dirakit di luar fungsi bernama", nama)
				} else if !adaUjiMenyebut(t, fn) {
					t.Errorf("%s: %s merakit SQL tetapi tidak satu pun uji "+
						"menyebut namanya; bentuk akhirnya karena itu tidak "+
						"pernah diperiksa", nama, fn)
				}
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol query ke tabel peserta terbaca; pembacanya yang rusak, bukan kodenya")
	}
	t.Logf("%d query ke tabel peserta diperiksa", diperiksa)
}

// tanpaKomentar membuang komentar baris Go sebelum SQL dipindai.
//
// ⛔ Tanpa ini, satu kalimat komentar yang menyebut "PL_NUMBER =" sudah
// cukup memuaskan penjaga, dan query di bawahnya boleh tak berpenyaring.
func tanpaKomentar(teks string) string {
	var b strings.Builder
	for _, baris := range strings.Split(teks, "\n") {
		if k := strings.Index(baris, "//"); k >= 0 {
			baris = baris[:k]
		}
		b.WriteString(baris)
		b.WriteByte('\n')
	}
	return b.String()
}

// namaFungsi membaca nama fungsi dari potongan yang diawali "\nfunc ".
func namaFungsi(teks string) string {
	i := strings.Index(teks, "func ")
	if i < 0 {
		return ""
	}
	sisa := teks[i+len("func "):]
	if strings.HasPrefix(sisa, "(") { // metode: lewati penerimanya
		if j := strings.Index(sisa, ")"); j >= 0 {
			sisa = strings.TrimSpace(sisa[j+1:])
		}
	}
	j := strings.IndexAny(sisa, "([")
	if j < 0 {
		return ""
	}
	return strings.TrimSpace(sisa[:j])
}

// adaUjiMenyebut mencari nama fungsi itu di berkas _test.go paket ini.
func adaUjiMenyebut(t *testing.T, fn string) bool {
	t.Helper()
	berkas, err := filepath.Glob("*_test.go")
	if err != nil || len(berkas) == 0 {
		t.Fatalf("nol berkas uji terbaca: %v", err)
	}
	for _, f := range berkas {
		isi, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		if strings.Contains(string(isi), fn+"(") {
			return true
		}
	}
	return false
}

// dibandingkan menjawab apakah sebuah kolom dipakai sebagai PENYARING.
//
// Yang membedakan penyaring dari sebutan biasa adalah operatornya: kolom di
// daftar SELECT diikuti koma, kolom di WHERE diikuti '=' atau 'LIKE'. Tanpa
// pembedaan itu, query yang penyaringnya dicabut tetap lolos hanya karena
// nama kolomnya muncul di daftar SELECT.
func dibandingkan(teks, kolom string) bool {
	for i := 0; ; {
		j := strings.Index(teks[i:], kolom)
		if j < 0 {
			return false
		}
		j += i
		sisa := strings.TrimLeft(teks[j+len(kolom):], " \t")
		if strings.HasPrefix(sisa, "=") ||
			strings.HasPrefix(sisa, "LIKE") {
			return true
		}
		i = j + len(kolom)
	}
}

// AC 29: penyaringan terjadi di SATU tempat.
func TestPenyaringPesertaHanyaSatuTempat(t *testing.T) {
	// ⚠️ Yang dilarang adalah MENYARING, bukan menyebut. Skema uji membuat
	// tabel tiruan yang punya kolom EDMSTATUS dan mengisinya - itu deklarasi
	// bentuk, bukan aturan kedua. Penjaga yang tidak membedakan keduanya akan
	// memaksa tabel tiruan dibuat tanpa kolom itu, dan penyaringnya justru
	// tidak pernah teruji terhadap Oracle.
	// ⚠️ Hanya bentuk PENYARING. "EDMSTATUS)" sempat ikut dan itu keliru: ia
	// cocok dengan daftar kolom INSERT tiruan, yang bukan aturan sama sekali.
	penyaring := []string{"EDMSTATUS IS ", "EDMSTATUS NOT IN",
		"EDMSTATUS IN", "EDMSTATUS ="}
	ketemu := 0
	for nama, isi := range berkasGoSelainTest(t) {
		if !strings.Contains(isi, "EDMSTATUS") {
			continue
		}
		ketemu++
		if strings.HasSuffix(nama, "/pesertapolis.go") {
			continue
		}
		atas := strings.ToUpper(isi)
		for _, pola := range penyaring {
			if strings.Contains(atas, pola) {
				t.Errorf("%s menyaring dengan %q; penyaringnya harus satu tempat (AC 29)",
					nama, pola)
			}
		}
	}
	// Tiga berkas menyebutnya: pembacanya, skema uji yang membuat tiruannya,
	// dan - sejak tiket 05a bagian 2 (pl2) - penulis salinan warisan
	// polis_warisan.go, yang menyebutnya di DAFTAR KOLOM `INSERT`, bukan
	// sebagai penyaring. Pola penyaring di atas tetap berlaku bagi ketiganya.
	const mau = 3
	if ketemu != mau {
		t.Errorf("berkas yang menyebut EDMSTATUS = %d, mau %d", ketemu, mau)
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

// ⛔ kolomSalin dan salinKePeserta memakai POSISI, dan posisinya dikunci.
//
// salinKePeserta membaca hasil SELECT lewat indeks tetap 0-27. Satu kolom yang
// disisipkan di tengah kolomSalin akan menggeser seluruh sisanya - dan tidak
// satu pun galat muncul: nilai hanya mendarat di medan yang salah. Tanggal
// valuasi menjadi tanggal lapse, uang menjadi uang lain, dan itu baru terlihat
// jauh di hilir, kalau pernah terlihat.
//
// Yang dikunci: cacahnya 28, dan nama kolom pada tiap posisi.
//
// ⚠️ Empat kolom TERAKHIR (24-27) berbeda sifatnya dari dua puluh empat
// yang pertama: ia BAHAN, bukan isi. Tidak satu pun mendarat langsung di
// medan models.Peserta - keempatnya masuk ke dua aturan pemilihan di
// pilihpeserta.go. Karena itu posisinya dikunci TERPISAH di bawah:
// menggesernya membuat umur terbaca dari kolom share, dan sebaliknya,
// tanpa satu pun galat.
func TestUrutanKolomSalinDikunci(t *testing.T) {
	// Pemecah kasar: tiap ekspresi dipisah koma di tingkat teratas.
	var ekspresi []string
	dalam := 0
	mulai := 0
	for i, c := range kolomSalin {
		switch c {
		case '(':
			dalam++
		case ')':
			dalam--
		case ',':
			if dalam == 0 {
				ekspresi = append(ekspresi, strings.TrimSpace(kolomSalin[mulai:i]))
				mulai = i + 1
			}
		}
	}
	ekspresi = append(ekspresi, strings.TrimSpace(kolomSalin[mulai:]))

	const mau = 28
	if len(ekspresi) != mau {
		t.Fatalf("kolomSalin memuat %d ekspresi, mau %d; salinKePeserta membaca "+
			"posisi 0-%d dan akan bergeser seluruhnya", len(ekspresi), mau, mau-1)
	}

	// Urutan nama kolom, persis seperti yang dibaca salinKePeserta.
	urut := []string{
		"ID", "PL_NUMBER", "POLICY_NO", "CERTIFICATE_NO", "CURRENCY", "STNC",
		"GROSS_VALUATION_BEGIN_DATE", "GROSS_VALUATION_EXPIRED_DATE",
		"RETRO_VALUATION_BEGIN_DATE", "RETRO_VALUATION_EXPIRED_DATE",
		"WPC", "BEGIN_DATE", "EFFECTIVE_DATE", "LAPSE_DATE", "EXPIRED_DATE",
		"SUM_INSURED", "SUM_REASURED", "GROSS_PREMIUM", "NET_PREMIUM",
		"CEDING_RETENTION", "SHARE_NUSANTARA_RE", "SHARE_RETRO",
		"RETROCEDED_SHARE", "EM_PERCENT",

		// Bahan bagi pilihpeserta.go, bukan isi medan.
		"SHARE_NUSANTARA_RE_GROSS", "AGE", "ENTRY_AGE", "CURRENT_AGE",
	}
	for i, nama := range urut {
		// ⛔ Nama diambil PERSIS, bukan lewat strings.Contains. Dengan
		// Contains, ekspresi SHARE_NUSANTARA_RE_GROSS lolos sebagai
		// SHARE_NUSANTARA_RE, dan ENTRY_AGE lolos sebagai AGE - persis dua
		// pasangan yang ada di daftar ini. Penjaga yang meloloskan kolom yang
		// salah lebih buruk daripada tidak ada penjaga, sebab ia menenangkan.
		if got := namaKolom(ekspresi[i]); got != nama {
			t.Errorf("posisi %d memuat kolom %q (dari %q), mau %s",
				i, got, ekspresi[i], nama)
		}
	}
}

// namaKolom mengeluarkan nama kolom dari satu ekspresi SELECT.
//
// Bentuk yang ditemui hanya dua: nama telanjang, dan TO_CHAR(NAMA, ...).
func namaKolom(ekspresi string) string {
	e := strings.TrimSpace(ekspresi)
	if i := strings.Index(e, "("); i >= 0 {
		e = e[i+1:]
		if j := strings.IndexAny(e, ",)"); j >= 0 {
			e = e[:j]
		}
	}
	return strings.TrimSpace(e)
}
