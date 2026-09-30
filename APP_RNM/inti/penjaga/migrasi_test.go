package penjaga

import (
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/db"
	"nusantarare/inti/migrasi"
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
	terlarang := map[string]string{
		"T_CLAIMLF_POLICY":            "tabel dihapus 2026-09-18, bukan diganti nama",
		"T_CLAIMLF_MARKETING":         "tabel dihapus 2026-09-18, bukan diganti nama",
		"T_CLAIM_POLICY":              "tabel dihapus 2026-09-18",
		"T_CLAIM_MARKETING":           "tabel dihapus 2026-09-18",
		"WORK_CLAIM_ID":               "dibuang; hubungannya shared primary key",
		"KMT_NO":                      "dibuang",
		"T_CLAIMLF_ADJUSTMENT_KOMITE": "roster komite tidak disimpan di Claim Life",
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
	// tidak disebut di sini menjaga kaskadenya sendiri (PremiumList:
	// `TestSeluruhFKPohonPolisBerkaskade`); Treaty Contract Out tanpa migrasi.
	berkaskade := map[string]map[string]bool{
		"claimlife": {
			"003_": true, "004_": true, "005_": true, "006_": true,
			// Relasi 10: diagnosa per peserta (butir bd). Buktinya bukan
			// selera: `.DiagnoseList` hidup DI DALAM halaman peserta -
			// `SetDisease.xml` b389 menutup dengan `Obj-Save pyWorkPage`,
			// bukan menyimpan halaman diagnosa sendiri. Menghapus peserta
			// karena itu menghapus daftarnya.
			"018_": true,
		},
		"komiteclaimlife": {
			// Relasi 9: roster komite. 013 membuatnya TANPA kaskade (cacat),
			// 030 memasangnya lewat ALTER. Keduanya terdaftar: yang pertama
			// karena kelak diperbaiki di tempatnya, yang kedua karena ia
			// perbaikannya.
			"030_": true,
		},
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
	// ⛔ Nama modul yang salah ketik di daftar tidak boleh mematikan
	// penjaganya diam-diam: setiap modul yang disebut harus punya berkas.
	for modul := range berkaskade {
		if diperiksa[modul] == 0 {
			t.Errorf("modul %q di daftar kaskade tidak punya satu pun berkas migrasi", modul)
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
	// Angkanya dikunci: kalau pemisah pernyataan rusak lagi, cacahnya anjlok
	// dan test ini gagal alih-alih diam-diam memeriksa lebih sedikit.
	// Dua puluh sejak 26-09-2026: langkah 009 menambah SEQ_WORK_CLAIM
	// (butir aa, tiket 02).
	//
	// ⛔ Angkanya DIPERBARUI, bukan dilonggarkan - A1 27-09-2026 menambah
	// butir am: 1 tabel (T_CLAIMLF_JEJAK) + 1 sequence (SEQ_CLAIMLF_JEJAK) +
	// 2 index = 4 pernyataan CREATE baru.
	//
	// ⛔ Diperbarui LAGI - A1 menambah butir af (2 tabel + 1 sequence +
	// 2 index) dan temuan audit A0 (BUSINESS_CODE, ALTER - tidak dihitung).
	//
	// ⛔ Diperbarui LAGI - A2 menambah butir aq: 1 tabel (T_LOG_SERVICE_RNM) +
	// 1 sequence + 2 index = 4 pernyataan CREATE baru.
	//
	// ⛔ Diperbarui LAGI - tiket 00 PremiumList Life menambah TUJUH tabel
	// (T_WORK_POLIS, T_PREMIUM_LIST, _DETAIL, _SPREADING, _SPREADING_RETRO,
	// _SUMMARY, T_VIEW_SUGGEST) dan ENAM index FK (dua pada _DETAIL, satu
	// pada masing-masing tabel anak lainnya). Nol sequence: pengenalnya
	// dirakit di repository, pola PengenalWorkBerikut (butir pl3).
	//
	// ⛔ Diperbarui LAGI - butir bd menambah T_CLAIMLF_DIAGNOSE:
	// 1 tabel + 1 sequence (SEQ_CLAIMLF_DIAGNOSE) + 1 index (FK peserta).
	// ALTER pada DISEASE tidak dihitung - ia bukan CREATE.
	//
	// 11+7+1 = 19 tabel + 10 sequence + 13+6+1 = 20 index = 49.
	// +1 tabel dari 019 (butir be, kartu berkas unggahan) = 50.
	//
	// ⛔ Diperbarui LAGI - butir bn (GILIRAN-13): +1 sequence SEQ_WORK_POLIS
	// (057). RALAT atas catatan tiket 00 di atas: "nol sequence" untuk polis
	// keliru - pl3 memutuskan SEQ_WORK_POLIS, dan ia terlewat. ALTER kolom
	// FLAG_ONGOING_POLICY tidak dihitung.
	// ⛔ Diperbarui LAGI - tiket 01 Treaty Contract Out (300-306) menambah
	// TUJUH tabel (T_TREATYYEAR, T_TREATYCONTRACT, T_TREATYREINSURER,
	// T_MTREATYSECURITY, T_TREATYBUSINESS, T_PROPORTIONALARRG,
	// T_TREATYCO_JEJAK) + 7 sequence + 7 index = 21 pernyataan CREATE = 71.
	// +3 dari 307 (tiket 12 Treaty Contract Out, 29-09-2026): tabel
	// T_TREATYYEAR_LAMPIRAN + index + sequence = 74.
	// Penyatuan 29-09-2026: 50 (dasar) + 1 (057, butir bn) + 24 (300-307 Treaty Contract Out) = 75.
	// ⛔ tco4 (keputusan work owner 29-09-2026): Treaty Contract Out NOL tabel
	// baru - 300-307 dibuang, kembali ke 50 + 1 = 51.
	// ⛔ OQ-PL-15 (GILIRAN-15): 058 membuat ULANG SEQ_WORK_POLIS - DROP (tidak
	// dihitung) lalu CREATE SEQUENCE ... START WITH 22374 (+1) = 52.
	const mau = 52
	if diperiksa != mau {
		t.Errorf("pernyataan CREATE diperiksa %d, mau %d", diperiksa, mau)
	}
}

// KolomCreateTable membaca nama dan kolom dari setiap CREATE TABLE migrasi.
//
// Cacahnya dikunci: delapan CREATE TABLE. Pernyataan yang BUKAN CREATE TABLE -
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
	// Tujuh, bukan delapan: T_MIGRASI dibuat siapkanTabelMigrasi, di luar
	// berkas migrasi. Sesudah migrasi, katalog memang memuat delapan tabel.
	// 11 tabel Claim Life + 7 tabel PremiumList Life (tiket 00) = 18.
	// 20 tabel + 10 sequence + 20 index = 50 pernyataan CREATE, cocok dengan
	// cacah yang dikunci TestSeluruhCreateDapatDibacaNamanya.
	// +1 tabel, +1 sequence, +1 index dari 018 (butir bd, diagnosa).
	// +1 tabel dari 019 (butir be) - TANPA sequence dan TANPA index:
	// identitasnya cap waktu `models.IDDokumenBaru`, bukan nomor kita,
	// dan PK-nya sudah berindeks sendiri.
	// +7 tabel dari 300-306 (tiket 01 Treaty Contract Out) = 27.
	// +1 tabel dari 307 (tiket 12 Treaty Contract Out, lampiran) = 28.
	// ⛔ tco4 (29-09-2026): 300-307 dibuang - kembali ke 20.
	const mauTabel = 20
	if tabel != mauTabel {
		t.Errorf("CREATE TABLE terbaca %d, mau %d", tabel, mauTabel)
	}
	if bukanTabel == 0 {
		t.Error("nol pernyataan bukan-tabel; CREATE INDEX dan SEQUENCE seharusnya ada")
	}
}
