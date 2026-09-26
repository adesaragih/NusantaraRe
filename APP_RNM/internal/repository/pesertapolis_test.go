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
	// Dua berkas menyebutnya: pembacanya, dan skema uji yang membuat tiruannya.
	const mau = 2
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
// salinKePeserta membaca hasil SELECT lewat indeks tetap 0-23. Satu kolom yang
// disisipkan di tengah kolomSalin akan menggeser seluruh sisanya - dan tidak
// satu pun galat muncul: nilai hanya mendarat di medan yang salah. Tanggal
// valuasi menjadi tanggal lapse, uang menjadi uang lain, dan itu baru terlihat
// jauh di hilir, kalau pernah terlihat.
//
// Yang dikunci: cacahnya 24, dan nama kolom pada tiap posisi.
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

	const mau = 24
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
	}
	for i, nama := range urut {
		if !strings.Contains(ekspresi[i], nama) {
			t.Errorf("posisi %d memuat %q, mau kolom %s", i, ekspresi[i], nama)
		}
	}
}
