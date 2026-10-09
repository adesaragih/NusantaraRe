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
	// Template Manager (keputusan work owner 04-10-2026): ISI adalah BERKAS templat utuh (CSV/XLSX yang diunduh
	// pengguna), bukan atribut klaim dan bukan dokumen JSON. Pengecualiannya SATU kolom BLOB di SATU tabel, dan
	// berkas itu tidak boleh memuat CLOB atau JSON.
	const berkasTemplat = "912_m_template_file"
	// R/I Rate Life (keputusan work owner 07-10-2026): 928 MEMBUANG kolom JSON warisan M_RATE_LIFE_SUMMARY.JSONDATA
	// (ringkasan satu tabel) dan 930 MEMBUANG M_RATE_LIFE.JSONDATA (rincian tabel flat) - searah dengan penjaga ini.
	// R/I Comm Life (keputusan work owner 08-10-2026, pola sama): 932 membuang M_RICOMM_LIFE_SUMMARY.JSONDATA dan 934
	// membuang M_RICOMM_LIFE.JSONDATA - tetap untuk tinjauan tim inti (modul/ricommlife/docs/PR-RICOMMLIFE.md).
	// R/I Risk (keputusan work owner 08-10-2026): 937 / 940 membuang JSONDATA tabel RIRISK_LIFE_SUMMARY / RIRISK_LIFE
	// (nama baru M_RIRISK_LIFE*) - tinjauan tim inti (modul/ririsklife/docs/PR-RIRISKLIFE.md).
	// Benefit (keputusan work owner 08-10-2026 K1/K2, pola R/I Risk): 944 membuang JSONDATA tabel BENEFIT_LIFE (nama baru
	// M_BENEFIT_LIFE) - tinjauan tim inti (modul/benefitlife/docs/PR-BENEFITLIFE.md).
	// Plan (keputusan work owner 08-10-2026 K1, pola Benefit): 948 membuang JSONDATA tabel PRODUCT_TYPE_LIFE (nama baru
	// M_PRODUCT_TYPE_LIFE) - tinjauan tim inti (modul/planlife/docs/PR-PLANLIFE.md).
	// Cause Of Loss Life (keputusan work owner 08-10-2026 K1, pola Benefit, migrasi MODUL K0): 092 membuang JSONDATA tabel
	// CAUSEOFLOSS_LIFE (nama baru M_CAUSEOFLOSS_LIFE) - tinjauan tim inti (modul/causeoflosslife/docs/PR-CAUSEOFLOSSLIFE.md).
	// Cover Life (keputusan work owner 08-10-2026 C1, pola Cause Of Loss TANPA RENAME, migrasi MODUL K0): 086 membuang
	// JSONDATA tabel M_COVER_LIFE (nama tetap) - tinjauan tim inti (modul/coverlife/docs/PR-COVERLIFE.md).
	// Pengecualiannya SATU perintah per berkas yang dinamai, persis.
	const perintahBuangJSON = "DROP COLUMN JSONDATA CASCADE CONSTRAINTS"
	buangJSON := map[string]int{"928_m_rate_life_summary_satu_tabel": 0, "930_m_rate_life_satu_tabel": 0,
		"932_m_ricomm_life_summary_satu_tabel": 0, "934_m_ricomm_life_satu_tabel": 0,
		"937_ririsk_life_summary_satu_tabel": 0, "940_ririsk_life_satu_tabel": 0, "944_benefit_life_satu_tabel": 0,
		"948_product_type_life_satu_tabel": 0, "092_causeofloss_life_satu_tabel": 0, "086_cover_life_satu_tabel": 0}
	dokumenDiOutbox, blobTemplat := 0, 0
	for nama, isi := range seluruhSQL(t, false) {
		atas := strings.ToUpper(isi)
		if strings.Contains(nama, berkasOutbox) {
			dokumenDiOutbox += strings.Count(atas, "CLOB")
			continue
		}
		if strings.Contains(nama, berkasTemplat) && !strings.Contains(nama, "_down") {
			blobTemplat += len(regexp.MustCompile(`(?m)^\s*ISI\s+BLOB\s+NOT NULL,$`).FindAllString(atas, -1))
			atas = regexp.MustCompile(`(?m)^\s*ISI\s+BLOB\s+NOT NULL,$`).ReplaceAllString(atas, "")
		}
		if kunci := strings.TrimSuffix(nama, ".sql"); !strings.Contains(nama, "_down") {
			if _, ada := buangJSON[kunci]; ada {
				buangJSON[kunci] += strings.Count(atas, perintahBuangJSON)
				atas = strings.ReplaceAll(atas, perintahBuangJSON, "")
			}
		}
		for _, tipe := range []string{" JSON", "CLOB", "BLOB", "JSON_KLAIM"} {
			if strings.Contains(atas, tipe) {
				t.Errorf("%s memuat %q - atribut klaim harus menjadi kolom bernama",
					nama, tipe)
			}
		}
	}
	for berkas, n := range buangJSON {
		if n != 1 {
			t.Errorf("%s memuat %d perintah %q, mau tepat 1", berkas, n, perintahBuangJSON)
		}
	}
	if blobTemplat != 1 {
		t.Errorf("M_TEMPLATE_FILE memuat %d kolom ISI BLOB, mau tepat 1", blobTemplat)
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
				// Blok PL/SQL hanya dalam DUA bentuk: berpelindung katalog
				// (`migrasi.BacaPerintahKatalog`, 901 menu datar) dan sequence
				// dari kueri (`migrasi.BacaSequenceDariKueri`, 810 Company
				// Detail) - dan perintah di dalamnya sendiri harus perintah SQL.
				if kata == "DECLARE" {
					perintah, alasan := perintahBlokPLSQL(p, mundur)
					if alasan != "" {
						t.Errorf("%s pernyataan %d: %s", m.Nama, i, alasan)
						continue
					}
					kata = strings.ToUpper(strings.Fields(perintah)[0])
				}
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

// polaTambahKolomSebaris - `ALTER TABLE ... ADD (kolom ...)` sebaris, bentuk
// yang ditulis di dalam EXECUTE IMMEDIATE (bukan ADD CONSTRAINT).
var polaTambahKolomSebaris = regexp.MustCompile(`(?i)^ALTER\s+TABLE\s+\S+\s+ADD\s*\(`)

// polaGantiNamaTabel - satu-satunya bentuk RENAME tabel di blok berpelindung (R/I Risk 935/938): nama baru TANPA skema
// (syarat Oracle `RENAME TO`).
var polaGantiNamaTabel = regexp.MustCompile(`^ALTER TABLE \{skema\}\.(\w+) RENAME TO (\w+)$`)

// perintahBlokPLSQL menjawab perintah di dalam blok PL/SQL, atau mengapa blok
// itu tidak sah (kosong = sah). Dua bentuk diterima: sequence dari kueri
// (`migrasi.BacaSequenceDariKueri`, 810 Company Detail - jalur maju saja,
// mundurnya cukup DROP SEQUENCE biasa) dan berpelindung katalog
// (`pelanggaranBlokPLSQL`).
func perintahBlokPLSQL(p string, mundur bool) (string, string) {
	if s, ok := migrasi.BacaSequenceDariKueri(p); ok {
		if mundur {
			return "", "sequence dari kueri hanya untuk jalur maju; jalur mundur cukup DROP SEQUENCE " + s.Nama
		}
		return "CREATE SEQUENCE {skema}." + s.Nama, ""
	}
	return pelanggaranBlokPLSQL(p, mundur)
}

// pelanggaranBlokPLSQL menjawab perintah di dalam blok PL/SQL berpelindung
// katalog (`migrasi.BacaPerintahKatalog`, 901 menu datar), atau mengapa blok
// itu tidak sah (kosong = sah).
func pelanggaranBlokPLSQL(p string, mundur bool) (string, string) {
	pk, ok := migrasi.BacaPerintahKatalog(p)
	if !ok || len(strings.Fields(pk.Perintah)) == 0 {
		return "", "blok PL/SQL di luar bentuk berpelindung katalog"
	}
	// Perintahnya atas tabel DAN objek yang DITANYAKAN katalog (DROP INDEX
	// menyebut indeksnya, bukan tabelnya) - blok yang memeriksa satu hal lalu
	// berbuat pada hal lain ditolak.
	atasTabel := strings.Contains(pk.Perintah, "{skema}."+pk.Tabel) ||
		(pk.Katalog == "ALL_INDEXES" && strings.Contains(pk.Perintah, "{skema}."+pk.Objek))
	// Objek dicocokkan sebagai KATA UTUH: kolom `C` tidak boleh "ditemukan"
	// di dalam kata COLUMN.
	objekUtuh := regexp.MustCompile(`(^|[^A-Za-z0-9_])` + regexp.QuoteMeta(pk.Objek) + `([^A-Za-z0-9_]|$)`)
	// RENAME tabel (R/I Risk 935/938, keputusan work owner 08-10-2026 K1): blok menanyakan kolom tabel SUMBER
	// (ALL_TAB_COLUMNS, n > 0 = sumber masih ada) dan perintahnya PERSIS `ALTER TABLE {skema}.<sumber> RENAME TO <baru>`
	// - kolom yang ditanyakan hanya bukti tabel sumber masih ada, jadi ia tidak disebut perintah. Target yang sudah ada
	// sebagai tabel membuat RENAME gagal keras (ORA-00955), bukan dilewati. ⚠️ Tinjauan tim inti.
	if m := polaGantiNamaTabel.FindStringSubmatch(pk.Perintah); m != nil {
		if pk.Katalog != "ALL_TAB_COLUMNS" || !pk.BilaAda || m[1] != pk.Tabel {
			return "", "RENAME harus berpelindung ALL_TAB_COLUMNS tabel sumbernya (n > 0): " + pk.Perintah
		}
		return pk.Perintah, ""
	}
	if !atasTabel || !objekUtuh.MatchString(pk.Perintah) {
		return "", "blok memeriksa " + pk.Tabel + "." + pk.Objek + " tetapi perintahnya " + pk.Perintah
	}
	// Kolom BARU lewat blok tidak terlihat `KolomAlterTambah` (pembanding
	// STRUKTUR, penjaga tipe): di jalur maju kolom ditambah `ALTER ... ADD (`
	// biasa. Jalur mundur (901_down) boleh - penjaga bentuk membaca jalur maju.
	if !mundur && polaTambahKolomSebaris.MatchString(pk.Perintah) {
		return "", "kolom baru lewat blok berpelindung tidak terlihat penjaga STRUKTUR: " + pk.Perintah
	}
	return pk.Perintah, ""
}

// Aturan blok PL/SQL MENGGIGIT (temuan /code-review): blok yang memeriksa satu
// objek lalu berbuat pada objek lain, perintah kosong, dan kolom baru lewat
// blok di jalur maju - semuanya ditolak.
func TestAturanBlokPLSQLMenggigit(t *testing.T) {
	blok := func(katalog, tabel, kolom, objek, banding, perintah string) string {
		return "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS." + katalog + "\n" +
			"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = '" + tabel + "' AND " + kolom + " = '" + objek + "';\n" +
			"  IF n " + banding + " 0 THEN\n    EXECUTE IMMEDIATE '" + perintah + "';\n  END IF;\nEND;"
	}
	for _, k := range []struct {
		nama   string
		p      string
		mundur bool
		sah    bool
	}{
		{"buang kolom yang ditanyakan", blok("ALL_TAB_COLUMNS", "T_A", "COLUMN_NAME", "C", ">", "ALTER TABLE {skema}.T_A DROP COLUMN C"), false, true},
		{"buang indeks yang ditanyakan", blok("ALL_INDEXES", "T_A", "INDEX_NAME", "IX_A", ">", "DROP INDEX {skema}.IX_A"), false, true},
		{"memeriksa T_A, membuang tabel lain", blok("ALL_TAB_COLUMNS", "T_A", "COLUMN_NAME", "C", ">", "DROP TABLE {skema}.T_B"), false, false},
		{"memeriksa kolom C, membuang kolom D", blok("ALL_TAB_COLUMNS", "T_A", "COLUMN_NAME", "C", ">", "ALTER TABLE {skema}.T_A DROP COLUMN D"), false, false},
		{"perintah kosong", blok("ALL_INDEXES", "T_A", "INDEX_NAME", "IX_A", ">", ""), false, false},
		{"kolom baru lewat blok, jalur maju", blok("ALL_TAB_COLUMNS", "T_A", "COLUMN_NAME", "C", "=", "ALTER TABLE {skema}.T_A ADD (C VARCHAR2(10))"), false, false},
		{"kolom kembali lewat blok, jalur mundur", blok("ALL_TAB_COLUMNS", "T_A", "COLUMN_NAME", "C", "=", "ALTER TABLE {skema}.T_A ADD (C VARCHAR2(10))"), true, true},
		{"tanpa pemeriksaan katalog", "BEGIN\n  EXECUTE IMMEDIATE 'DROP TABLE {skema}.T_A';\nEND;", false, false},
		{"RENAME tabel sumber yang ditanyakan", blok("ALL_TAB_COLUMNS", "M_T_A", "COLUMN_NAME", "ID", ">", "ALTER TABLE {skema}.M_T_A RENAME TO T_A"), false, true},
		{"RENAME tabel lain", blok("ALL_TAB_COLUMNS", "M_T_A", "COLUMN_NAME", "ID", ">", "ALTER TABLE {skema}.M_T_B RENAME TO T_A"), false, false},
		{"RENAME bila sumber TIDAK ada", blok("ALL_TAB_COLUMNS", "M_T_A", "COLUMN_NAME", "ID", "=", "ALTER TABLE {skema}.M_T_A RENAME TO T_A"), false, false},
		{"RENAME berpelindung indeks", blok("ALL_INDEXES", "M_T_A", "INDEX_NAME", "IX_A", ">", "ALTER TABLE {skema}.M_T_A RENAME TO T_A"), false, false},
		{"DROP VIEW berpelindung ALL_VIEWS", "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_VIEWS\n" +
			"   WHERE OWNER = UPPER('{skema}') AND VIEW_NAME = 'V_A';\n  IF n > 0 THEN\n    EXECUTE IMMEDIATE 'DROP VIEW {skema}.V_A';\n  END IF;\nEND;", false, true},
		{"ALL_VIEWS membuang view lain", "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_VIEWS\n" +
			"   WHERE OWNER = UPPER('{skema}') AND VIEW_NAME = 'V_A';\n  IF n > 0 THEN\n    EXECUTE IMMEDIATE 'DROP VIEW {skema}.V_B';\n  END IF;\nEND;", false, false},
		{"indeks bila belum ada indeks berkolom pertama itu", "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_IND_COLUMNS\n" +
			"   WHERE TABLE_OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_A' AND COLUMN_NAME = 'C' AND COLUMN_POSITION = 1;\n" +
			"  IF n = 0 THEN\n    EXECUTE IMMEDIATE 'CREATE INDEX {skema}.IX_T_A_C ON {skema}.T_A (C)';\n  END IF;\nEND;", false, true},
		{"ALL_IND_COLUMNS mengindeks kolom lain", "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_IND_COLUMNS\n" +
			"   WHERE TABLE_OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_A' AND COLUMN_NAME = 'C' AND COLUMN_POSITION = 1;\n" +
			"  IF n = 0 THEN\n    EXECUTE IMMEDIATE 'CREATE INDEX {skema}.IX_T_A_D ON {skema}.T_A (D)';\n  END IF;\nEND;", false, false},
	} {
		if _, alasan := pelanggaranBlokPLSQL(k.p, k.mundur); (alasan == "") != k.sah {
			t.Errorf("%s: sah=%v, mau %v (%s)", k.nama, alasan == "", k.sah, alasan)
		}
	}
}

// Bentuk kedua - sequence yang nilai awalnya dihitung dari data (810 Company Detail) - diterima di jalur maju dan
// terbaca sebagai CREATE SEQUENCE; di jalur mundur DITOLAK (mundur cukup DROP SEQUENCE biasa).
func TestBlokSequenceDariKueriHanyaJalurMaju(t *testing.T) {
	blok := "DECLARE\n  n    NUMBER;\n  awal NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_SEQUENCES\n" +
		"   WHERE SEQUENCE_OWNER = UPPER('{skema}') AND SEQUENCE_NAME = 'SEQ_A';\n  IF n = 0 THEN\n" +
		"    SELECT NVL(MAX(N), 0) + 1 INTO awal FROM (SELECT 1 N FROM {skema}.T_A);\n" +
		"    EXECUTE IMMEDIATE 'CREATE SEQUENCE {skema}.SEQ_A START WITH ' || awal || ' INCREMENT BY 1 NOCACHE';\n" +
		"  END IF;\nEND;"
	perintah, alasan := perintahBlokPLSQL(blok, false)
	if alasan != "" || perintah != "CREATE SEQUENCE {skema}.SEQ_A" {
		t.Errorf("jalur maju: perintah %q, alasan %q", perintah, alasan)
	}
	if _, alasan := perintahBlokPLSQL(blok, true); alasan == "" {
		t.Error("jalur mundur menerima sequence dari kueri")
	}
	// Bentuk pertama tetap lewat jalan yang sama.
	katalog := "DECLARE\n  n NUMBER;\nBEGIN\n  SELECT COUNT(*) INTO n FROM SYS.ALL_INDEXES\n" +
		"   WHERE OWNER = UPPER('{skema}') AND TABLE_NAME = 'T_A' AND INDEX_NAME = 'IX_A';\n" +
		"  IF n > 0 THEN\n    EXECUTE IMMEDIATE 'DROP INDEX {skema}.IX_A';\n  END IF;\nEND;"
	if perintah, alasan := perintahBlokPLSQL(katalog, false); alasan != "" || perintah != "DROP INDEX {skema}.IX_A" {
		t.Errorf("blok katalog: perintah %q, alasan %q", perintah, alasan)
	}
}
