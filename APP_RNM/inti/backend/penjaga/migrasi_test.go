package penjaga

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/migrasi"
)

// Delapan langkah maju, berurut, dan masing-masing punya jalur mundur.
func TestSetiapLangkahPunyaJalurMundur(t *testing.T) {
	maju, err := migrasi.Daftar(false, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	mundur, err := migrasi.Daftar(true, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	if len(maju) == 0 {
		t.Fatal("tidak ada langkah maju sama sekali")
	}
	if len(maju) != len(mundur) {
		t.Fatalf("maju %d langkah, mundur %d langkah", len(maju), len(mundur))
	}

	punyaMundur := map[string]bool{}
	for _, m := range mundur {
		punyaMundur[migrasi.KunciLangkah(m.Nama)] = true
	}
	for _, m := range maju {
		if !punyaMundur[migrasi.KunciLangkah(m.Nama)] {
			t.Errorf("langkah %s tidak punya jalur mundur", m.Nama)
		}
	}
}

// Jalur mundur berjalan MENURUN supaya anak dibongkar sebelum induknya.
func TestJalurMundurBerurutMenurun(t *testing.T) {
	mundur, err := migrasi.Daftar(true, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(mundur); i++ {
		if mundur[i-1].Nama < mundur[i].Nama {
			t.Errorf("urutan mundur naik di %s lalu %s", mundur[i-1].Nama, mundur[i].Nama)
		}
	}
}

// Setiap pernyataan harus berisi, menyebut skema, dan lolos PeriksaSQL.
func TestSetiapPernyataanSahDanBerskema(t *testing.T) {
	for _, mundur := range []bool{false, true} {
		langkah, err := migrasi.Daftar(mundur, berkasMigrasi)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range langkah {
			if len(m.Pernyataan) == 0 {
				t.Errorf("%s: nol pernyataan", m.Nama)
			}
			for i, p := range m.Pernyataan {
				if strings.TrimSpace(p) == "" {
					t.Errorf("%s pernyataan %d kosong", m.Nama, i)
				}
				if !strings.Contains(p, "{skema}") {
					t.Errorf("%s pernyataan %d tidak menyebut skema (ADR-U-0033): %.60s",
						m.Nama, i, p)
				}
				// ADR-U-0029: nol COMMIT di teks SQL.
				if err := db.PeriksaSQL(p); err != nil {
					t.Errorf("%s pernyataan %d: %v", m.Nama, i, err)
				}
			}
		}
	}
}

// ⛔ Nama-nama yang DIBUANG. Tiket 14 menyatakan test yang menemukannya gagal.
func TestNamaYangDibuangTidakAda(t *testing.T) {
	sql := gabungSemua(t)
	// Struktur tim satu folder per modul (30-09-2026): namanya dinyatakan bab
	// "Nama terlarang di migrasi" MODUL.md modul yang membuangnya.
	terlarang := namaTerlarangDiMigrasi(t)
	if len(terlarang) == 0 {
		t.Fatal("nol nama terlarang dinyatakan; pembacanya yang rusak")
	}
	for n, sebab := range terlarang {
		if strings.Contains(sql, n) {
			t.Errorf("nama terlarang %q muncul di migrasi - %s", n, sebab)
		}
	}
}

// Kaskade pada relasi 3, 4, 5, 6 - dan relasi 9, roster komite.
// T_CLAIMLF_DOCUMENT (relasi 7) ditangani di Go, jadi kunci tamunya TANPA
// ON DELETE.
//
// ⚠️ LINGKUPNYA MODUL CLAIM LIFE - migrasi 001-049. Sejak modul
// PremiumList Life menambah migrasi 050+, penjaga ini harus menyatakan
// lingkupnya: kebijakan kaskade kedua modul BERBEDA dan keduanya disengaja.
// Claim Life memilih kaskade pada empat relasi (kini lima); PremiumList Life
// memilih kaskade pada SELURUH FK (tiket 00 AC 46, dijaga
// TestSeluruhFKPohonPolisBerkaskade). Penjaga yang membentang ke modul lain
// akan memaksa salah satunya mengalah tanpa alasan.
//
// ⛔ RALAT 27-09-2026, dan penjaga ini sempat MENEGAKKAN cacatnya sendiri.
// Daftar di bawah ditulis ketika hanya migrasi Claim Life ada, dan ia
// menuntut `013_tabel_komite.sql` TIDAK berkaskade. Tetapi
// `STRUKTUR-TABEL-KOMITE-CLAIM-LIFE.md` menyebut relasi
// `T_GENERAL_KOMITE` -> `T_KOMITE_KOMITELIST` sebagai `ON DELETE CASCADE` di
// TIGA tempat (baris 159, 238, 253), dan 013 membuatnya tanpa `ON DELETE`.
// Jadi penjaga ini bukan hanya melewatkan cacat - ia menahan perbaikannya.
//
// ⚠️ Akibat cacat itu nyata: menghapus satu kasus komite DITOLAK Oracle
// (ORA-02292) selama masih ada baris roster yang menunjuknya, dan itu jalur
// yang tiket 05 perlukan. Migrasi `030` memperbaikinya lewat ALTER.
//
// ⛔ KEDUA PELAJARAN ITU SATU, dan disatukan di sini 27-09-2026: daftar yang
// ditulis sebelum modul kedua lahir berhenti menjadi penjaga dan mulai
// menjadi pagar. Yang satu menyempitkan LINGKUPnya, yang lain memperbaiki
// ISInya; keduanya perlu.
// ⛔ NAMANYA BERUBAH 27-09-2026, dan sebabnya adalah namanya sendiri.
// Ia lahir sebagai `TestKaskadeHanyaPadaEmpatRelasi` ketika relasi berkaskade
// memang empat. Kini enam berkas terdaftar, dan nama yang menyebut ANGKA
// berbohong setiap kali relasi ketujuh lahir - sementara nama yang menyebut
// ATURANNYA tidak pernah berbohong. Nama lamanya ditulis di sini supaya
// pencarian atasnya tetap sampai ke tempat ini.
//
// ⚠️ Yang dijaga bukan jumlahnya melainkan kesengajaannya: kaskade ada
// HANYA pada relasi yang terdaftar di bawah, dan mendaftarkan yang baru
// menuntut bukti - bukan kemudahan.
func TestKaskadeHanyaPadaRelasiTerdaftar(t *testing.T) {
	berkas := seluruhSQL(t, false)
	// Refactor bentuk B paket 8: daftarnya PER MODUL, menurut folder migrasi
	// berkasnya - dulu satu daftar untuk rentang nomor 001-049. Modul yang
	// tidak menyatakan kebijakannya menjaga kaskadenya sendiri (PremiumList:
	// `TestSeluruhFKPohonPolisBerkaskade`); Treaty Contract Out tanpa migrasi.
	// Struktur tim satu folder per modul (30-09-2026): daftarnya dinyatakan bab
	// "Kaskade ON DELETE CASCADE" MODUL.md setiap modul (`kaskadePerModul`).
	berkaskade := kaskadePerModul(t)
	if len(berkaskade) == 0 {
		t.Fatal("nol kebijakan kaskade dinyatakan; pembacanya yang rusak")
	}
	diperiksa := map[string]int{}
	for nama, teks := range berkas {
		modul := berkasMigrasi.modul(nama)
		terdaftar, diatur := berkaskade[modul]
		if !diatur {
			continue
		}
		diperiksa[modul]++
		isi := strings.ToUpper(teks)
		ada := strings.Contains(isi, "ON DELETE CASCADE")
		mau := false
		for awalan := range terdaftar {
			if strings.HasPrefix(nama, awalan) {
				mau = true
			}
		}
		if ada != mau {
			t.Errorf("%s: ON DELETE CASCADE ada=%v, mau=%v", nama, ada, mau)
		}
	}
	// ⛔ Kebijakan yang dinyatakan tanpa satu pun berkas tidak boleh mematikan
	// penjaganya diam-diam: setiap modul yang menyatakannya harus punya berkas.
	for modul := range berkaskade {
		if diperiksa[modul] == 0 {
			t.Errorf("modul %q menyatakan kebijakan kaskade tetapi tidak punya satu pun berkas migrasi", modul)
		}
	}
	t.Logf("berkas diperiksa per modul: %v", diperiksa)
}

// AC: seluruh kolom uang bertipe desimal, tidak ada yang berupa teks, dan
// tidak ada kolom JSON yang menyimpan atribut klaim.
func TestKolomUangDesimalDanNolJSON(t *testing.T) {
	sql := gabungSemua(t)
	for _, kol := range []string{"CLAIM_AMOUNT", "SUM_INSURED", "SUM_REASURED"} {
		pola := regexp.MustCompile(kol + `\s+NUMBER\(38,8\)`)
		if !pola.MatchString(sql) {
			t.Errorf("kolom uang %s tidak bertipe NUMBER(38,8)", kol)
		}
		if regexp.MustCompile(kol + `\s+(VARCHAR2|CHAR|CLOB)`).MatchString(sql) {
			t.Errorf("kolom uang %s bertipe teks", kol)
		}
	}
	// ⛔ DIPERSEMPIT, bukan dilonggarkan - A2, 27-09-2026.
	//
	// Yang dilarang adalah ATRIBUT KLAIM yang bersembunyi di dalam dokumen -
	// itulah inti keputusan "yang dibuang hanya JSON". `T_LOG_SERVICE_RNM.MUATAN`
	// bukan atribut klaim: ia BADAN PESAN antrean, yang memang berbentuk
	// dokumen dan memang tidak boleh dipecah menjadi kolom - setiap jenis efek
	// punya bentuk muatannya sendiri.
	//
	// Cakupannya dinyatakan: pengecualian berlaku untuk SATU kolom di SATU
	// tabel, dan penjaga terpisah memastikan muatan itu tidak memuat nama
	// orang, kredensial, maupun alamat.
	// ⛔ DIPERIKSA PER BERKAS, bukan dengan memotong teks gabungan.
	//
	// Ronde pertama memotong blok `CREATE TABLE` dari teks yang sudah
	// digabung - dan potongannya MELESET: 6.382 dari 10.037 karakter ikut
	// terbuang, sehingga penjaganya berhenti memeriksa sebagian besar
	// migrasi tanpa ada yang tahu. Ketahuan hanya karena panjangnya dicetak.
	//
	// Pemeriksaan per berkas tidak dapat salah potong: satu berkas
	// dikecualikan dengan namanya, sisanya utuh.
	const berkasOutbox = "015_t_log_service_rnm"
	dokumenDiOutbox := 0
	for nama, isi := range seluruhSQL(t, false) {
		atas := strings.ToUpper(isi)
		if strings.Contains(nama, berkasOutbox) {
			dokumenDiOutbox += strings.Count(atas, "CLOB")
			continue
		}
		for _, tipe := range []string{" JSON", "CLOB", "BLOB", "JSON_KLAIM"} {
			if strings.Contains(atas, tipe) {
				t.Errorf("%s memuat %q - atribut klaim harus menjadi kolom bernama",
					nama, tipe)
			}
		}
	}
	// Dan outbox-nya memang hanya punya SATU kolom dokumen.
	if dokumenDiOutbox != 1 {
		t.Errorf("T_LOG_SERVICE_RNM memuat %d kolom CLOB, mau tepat 1 (MUATAN)",
			dokumenDiOutbox)
	}
	if strings.Contains(sql, "FLOAT") || strings.Contains(sql, "BINARY_DOUBLE") {
		t.Error("ada kolom bertipe float - uang tidak pernah float (ADR-U-0003)")
	}
}

// Berkas migrasi tidak boleh berawalan byte order mark.
//
// Ini bukan kerewelan gaya. BOM membuat baris komentar pertama lolos menjadi
// bagian pernyataan SQL, dan Oracle menolaknya - sementara seluruh test tanpa
// basis data tetap hijau. Sekali terjadi di berkas 004; sejak itu dijaga.
func TestBerkasMigrasiTanpaBOM(t *testing.T) {
	entri, err := berkasMigrasi.ReadDir("migrations")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entri {
		isi, err := berkasMigrasi.ReadFile("migrations/" + e.Name())
		if err != nil {
			t.Fatal(err)
		}
		if len(isi) >= 3 && isi[0] == 0xEF && isi[1] == 0xBB && isi[2] == 0xBF {
			t.Errorf("%s berawalan BOM UTF-8", e.Name())
		}
	}
}

// Berkas migrasi tidak boleh memuat byte carriage return.
//
// Sepupu dekat jebakan BOM di atas. Alat Windows - PowerShell, penyunting yang
// disetel salah, git tanpa .gitattributes - menulis akhiran baris CRLF. CR yang
// terbawa masuk ke teks pernyataan yang dikirim ke Oracle, dan sekali lagi
// seluruh test tanpa basis data tetap hijau sementara instance menolaknya.
// Berkas .gitattributes di akar repositori menjaga sisi git; test ini menjaga
// sisi berkas.
func TestBerkasMigrasiTanpaCR(t *testing.T) {
	entri, err := berkasMigrasi.ReadDir("migrations")
	if err != nil {
		t.Fatal(err)
	}
	if len(entri) == 0 {
		t.Fatal("nol berkas migrasi terbaca; pembacanya yang rusak")
	}
	for _, e := range entri {
		isi, err := berkasMigrasi.ReadFile("migrations/" + e.Name())
		if err != nil {
			t.Fatal(err)
		}
		if n := strings.Count(string(isi), "\r"); n > 0 {
			t.Errorf("%s memuat %d byte CR; akhiran barisnya harus LF", e.Name(), n)
		}
	}
}

// Pernyataan pertama setiap berkas harus benar-benar mulai dengan kata perintah
// SQL - bukan dengan sisa komentar.
func TestPernyataanMulaiDenganPerintah(t *testing.T) {
	for _, mundur := range []bool{false, true} {
		langkah, err := migrasi.Daftar(mundur, berkasMigrasi)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range langkah {
			for i, p := range m.Pernyataan {
				kata := strings.ToUpper(strings.Fields(p)[0])
				switch kata {
				case "CREATE", "ALTER", "DROP", "INSERT", "UPDATE", "DELETE", "SELECT":
				default:
					t.Errorf("%s pernyataan %d mulai dengan %q, bukan perintah SQL",
						m.Nama, i, kata)
				}
			}
		}
	}
}

// Setiap pernyataan CREATE di berkas migrasi harus dapat dibaca namanya.
//
// Kalau ada satu saja yang tidak terbaca, jalur "dilewati lalu dibuktikan"
// menolak dengan galat - dan lebih baik test ini yang menemukannya lebih dulu,
// di mesin tanpa Oracle.
//
// ⚠️ Versi pertama test ini memecah ulang teks yang sudah disambung seluruhSQL,
// sehingga pemisahnya tidak pernah memisah apa pun dan yang diperiksa hanya
// SATU pernyataan per berkas - 8 dari 19. Penjaga yang lebih lemah dari
// namanya. Sekarang pernyataannya diambil dari daftarMigrasi apa adanya.
func TestSeluruhCreateDapatDibacaNamanya(t *testing.T) {
	langkah, err := migrasi.Daftar(false, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	diperiksa := 0
	for _, m := range langkah {
		for _, p := range m.Pernyataan {
			if !migrasi.PernyataanBuat(p) {
				continue
			}
			diperiksa++
			if migrasi.NamaObjekDibuat(p) == "" {
				t.Errorf("%s: nama objek tidak terbaca dari %q", m.Nama, migrasi.RingkasPernyataan(p))
			}
		}
	}
	// ⛔ Cacahnya dihitung DUA CARA (CLAUDE.md §4a): dari pernyataan hasil
	// pemisah pelari, dan dari baris `CREATE` di teks mentah berkas migrasi.
	// Kalau pemisah pernyataan rusak lagi, keduanya berselisih dan test ini
	// gagal alih-alih diam-diam memeriksa lebih sedikit.
	//
	// Struktur tim satu folder per modul (30-09-2026): dulu angkanya DIKUNCI
	// (`const mau = 56`, riwayat setiap kenaikannya di git log berkas ini), dan
	// setiap modul yang menambah tabel menyunting berkas milik tim inti ini -
	// dua modul yang melakukannya bersamaan bertabrakan di satu baris. Dua cara
	// hitung menjaga hal yang sama tanpa angka yang harus disunting. Instrumennya
	// diuji atas angka lama: 56 baris CREATE, 21 CREATE TABLE (30-09-2026).
	if mentah := cacahBarisMentah(t, polaBarisCreate); diperiksa != mentah {
		t.Errorf("pernyataan CREATE diperiksa %d, teks mentah memuat %d baris CREATE - pemisah pernyataan rusak?",
			diperiksa, mentah)
	}
	if diperiksa < 20 {
		t.Fatalf("hanya %d pernyataan CREATE terbaca; pembacanya yang rusak", diperiksa)
	}
}

// KolomCreateTable membaca nama dan kolom dari setiap CREATE TABLE migrasi.
//
// Cacahnya dihitung dua cara (lihat di bawah). Pernyataan yang BUKAN CREATE TABLE -
// CREATE INDEX dan CREATE SEQUENCE - harus mengembalikan nama kosong, kalau
// tidak pra-terbang akan mencari "bentuk" sebuah sequence.
func TestKolomCreateTableMembacaSeluruhTabel(t *testing.T) {
	langkah, err := migrasi.Daftar(false, berkasMigrasi)
	if err != nil {
		t.Fatal(err)
	}
	tabel, bukanTabel := 0, 0
	for _, m := range langkah {
		for _, p := range m.Pernyataan {
			nama, kolom := migrasi.KolomCreateTable(p)
			if nama == "" {
				bukanTabel++
				continue
			}
			tabel++
			if len(kolom) == 0 {
				t.Errorf("%s: tabel %s terbaca tanpa satu pun kolom", m.Nama, nama)
			}
			for _, k := range kolom {
				if strings.HasPrefix(k, "CONSTRAINT") || strings.HasPrefix(k, "REFERENCES") {
					t.Errorf("%s: %s menganggap %q sebagai kolom", m.Nama, nama, k)
				}
			}
		}
	}
	// T_MIGRASI dibuat siapkanTabelMigrasi, di luar berkas migrasi - tidak
	// terhitung di sini. Dua cara hitung, seperti TestSeluruhCreateDapatDibacaNamanya:
	// dulu `const mauTabel = 21` yang disunting setiap modul yang menambah tabel.
	if mentah := cacahBarisMentah(t, polaBarisCreateTable); tabel != mentah {
		t.Errorf("CREATE TABLE terbaca %d, teks mentah memuat %d baris CREATE TABLE", tabel, mentah)
	}
	if tabel < 10 {
		t.Fatalf("hanya %d CREATE TABLE terbaca; pembacanya yang rusak", tabel)
	}
	if bukanTabel == 0 {
		t.Error("nol pernyataan bukan-tabel; CREATE INDEX dan SEQUENCE seharusnya ada")
	}
}
