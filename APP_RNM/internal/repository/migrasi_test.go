package repository

// Test bentuk migrasi - TANPA Oracle.
//
// Untuk apa berkas ini: memeriksa isi berkas .sql yang ditanam ke biner. Ia
// tidak menyambung ke basis data sama sekali; yang diuji adalah apa yang akan
// dikirim ke Oracle. Karena itu ia tetap berjalan di mesin tanpa instance.
//
// Banyak acceptance criteria tiket 14 berbentuk "test yang menemukan X gagal".
// Justru itu yang dikerjakan di sini.

import (
	"errors"
	"regexp"
	"strings"
	"testing"
)

func seluruhSQL(t *testing.T, mundur bool) map[string]string {
	t.Helper()
	langkah, err := daftarMigrasi(mundur)
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, m := range langkah {
		out[m.Nama] = strings.Join(m.Pernyataan, "\n")
	}
	return out
}

func gabungSemua(t *testing.T) string {
	t.Helper()
	var b strings.Builder
	for _, isi := range seluruhSQL(t, false) {
		b.WriteString(isi)
		b.WriteString("\n")
	}
	return strings.ToUpper(b.String())
}

// Delapan langkah maju, berurut, dan masing-masing punya jalur mundur.
func TestSetiapLangkahPunyaJalurMundur(t *testing.T) {
	maju, err := daftarMigrasi(false)
	if err != nil {
		t.Fatal(err)
	}
	mundur, err := daftarMigrasi(true)
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
		punyaMundur[kunciLangkah(m.Nama)] = true
	}
	for _, m := range maju {
		if !punyaMundur[kunciLangkah(m.Nama)] {
			t.Errorf("langkah %s tidak punya jalur mundur", m.Nama)
		}
	}
}

// Jalur mundur berjalan MENURUN supaya anak dibongkar sebelum induknya.
func TestJalurMundurBerurutMenurun(t *testing.T) {
	mundur, err := daftarMigrasi(true)
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
		langkah, err := daftarMigrasi(mundur)
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
				if err := PeriksaSQL(p); err != nil {
					t.Errorf("%s pernyataan %d: %v", m.Nama, i, err)
				}
			}
		}
	}
}

// Komentar murni tidak ikut menjadi pernyataan.
func TestKomentarTidakMenjadiPernyataan(t *testing.T) {
	contoh := "-- hanya komentar\n-- baris kedua\n/\nCREATE TABLE {skema}.X (A NUMBER)\n/\n"
	p := pecahPernyataan(contoh)
	if len(p) != 1 {
		t.Fatalf("dapat %d pernyataan, mau 1: %q", len(p), p)
	}
	if !strings.HasPrefix(p[0], "CREATE TABLE") {
		t.Errorf("pernyataan salah: %q", p[0])
	}
}

// Tiket 14: enam tabel klaim + T_WORK_CLAIM harus dibuat.
func TestTujuhTabelDibuat(t *testing.T) {
	sql := gabungSemua(t)
	mau := []string{
		"T_WORK_CLAIM",
		"T_GENERAL_CLAIM",
		"T_CLAIMLF_PREMIUMLIST_DETAIL",
		"T_CLAIMLF_ADJUSTMENT",
		"T_CLAIMLF_ADJUSTMENT_SPREADING",
		"T_CLAIMLF_ADJ_SPREADING_RETRO",
		"T_CLAIMLF_DOCUMENT",
	}
	for _, tb := range mau {
		if !strings.Contains(sql, "CREATE TABLE {SKEMA}."+tb+" (") {
			t.Errorf("tabel %s tidak dibuat", tb)
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

// Header klaim tidak lagi memuat kolom yang pindah ke T_WORK_CLAIM,
// dan PL_NUMBER sudah berganti nama menjadi POLICY_NO.
func TestHeaderKlaimTidakMemuatKolomYangPindah(t *testing.T) {
	isi := ""
	for nama, teks := range seluruhSQL(t, false) {
		if strings.Contains(nama, "t_general_claim") {
			isi = strings.ToUpper(teks)
		}
	}
	if isi == "" {
		t.Fatal("berkas T_GENERAL_CLAIM tidak ketemu")
	}
	for _, kol := range []string{"CASEID ", "CREATE_OP", "CREATE_OP_NAME", "TGL_UPDATE", "PL_NUMBER"} {
		if strings.Contains(isi, kol) {
			t.Errorf("kolom %q masih ada di T_GENERAL_CLAIM - ia pindah atau berganti nama", kol)
		}
	}
	for _, kol := range []string{"CASEID_POLICY", "POLICY_NO", "ENDORSMENT_NO"} {
		if !strings.Contains(isi, kol) {
			t.Errorf("penunjuk polis %q tidak ada di T_GENERAL_CLAIM", kol)
		}
	}
}

// Relasi 2: shared primary key - T_GENERAL_CLAIM.ID sekaligus PK dan FK.
func TestSharedPrimaryKey(t *testing.T) {
	var isi string
	for nama, teks := range seluruhSQL(t, false) {
		if strings.Contains(nama, "t_general_claim") {
			isi = strings.ToUpper(teks)
		}
	}
	if !strings.Contains(isi, "PRIMARY KEY (ID)") {
		t.Error("T_GENERAL_CLAIM.ID bukan kunci utama")
	}
	if !strings.Contains(isi, "FOREIGN KEY (ID)") {
		t.Error("T_GENERAL_CLAIM.ID bukan kunci tamu ke T_WORK_CLAIM - shared PK tidak terbentuk")
	}
}

// Relasi 4: adjustment menggantung pada PESERTA, bukan pada header klaim.
func TestAdjustmentMenggantungPadaPeserta(t *testing.T) {
	var isi string
	for nama, teks := range seluruhSQL(t, false) {
		if strings.HasPrefix(nama, "004_") {
			isi = strings.ToUpper(teks)
		}
	}
	if !strings.Contains(isi, "REFERENCES {SKEMA}.T_CLAIMLF_PREMIUMLIST_DETAIL (ID)") {
		t.Error("FK adjustment tidak menunjuk T_CLAIMLF_PREMIUMLIST_DETAIL")
	}
	if strings.Contains(isi, "REFERENCES {SKEMA}.T_GENERAL_CLAIM") {
		t.Error("adjustment menggantung pada header klaim - seharusnya pada peserta")
	}
}

// Kaskade hanya pada relasi 3, 4, 5, 6. T_CLAIMLF_DOCUMENT (relasi 7) ditangani
// di Go, jadi kunci tamunya TANPA ON DELETE.
func TestKaskadeHanyaPadaEmpatRelasi(t *testing.T) {
	berkas := seluruhSQL(t, false)
	berkaskade := map[string]bool{"003_": true, "004_": true, "005_": true, "006_": true}
	for nama, teks := range berkas {
		isi := strings.ToUpper(teks)
		ada := strings.Contains(isi, "ON DELETE CASCADE")
		mau := false
		for awalan := range berkaskade {
			if strings.HasPrefix(nama, awalan) {
				mau = true
			}
		}
		if ada != mau {
			t.Errorf("%s: ON DELETE CASCADE ada=%v, mau=%v", nama, ada, mau)
		}
	}
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
	for _, tipe := range []string{" JSON", "CLOB", "BLOB", "JSON_KLAIM"} {
		if strings.Contains(sql, tipe) {
			t.Errorf("migrasi memuat %q - atribut klaim harus menjadi kolom bernama", tipe)
		}
	}
	if strings.Contains(sql, "FLOAT") || strings.Contains(sql, "BINARY_DOUBLE") {
		t.Error("ada kolom bertipe float - uang tidak pernah float (ADR-U-0003)")
	}
}

// Setiap kunci tamu ber-index, dan KOMITE_ID ber-index UNIK.
func TestKunciTamuBerIndex(t *testing.T) {
	// Pernyataan CREATE INDEX boleh ditulis pada beberapa baris, jadi spasinya
	// diratakan dulu. Tanpa ini, test-nya sendiri yang keliru melaporkan index
	// yang sebenarnya ada.
	sql := regexp.MustCompile(`\s+`).ReplaceAllString(gabungSemua(t), " ")
	for _, kol := range []string{
		"COVER_KEY", "CLAIM_ID", "PREMIUM_LIST_DETAIL_ID", "ADJUSTMENT_ID", "SPREADING_ID",
	} {
		if !regexp.MustCompile(`CREATE INDEX [^;]*?\(` + kol + `\)`).MatchString(sql) {
			t.Errorf("kunci tamu %s tidak ber-index", kol)
		}
	}
	if !regexp.MustCompile(`CREATE UNIQUE INDEX [^;]*?\(KOMITE_ID\)`).MatchString(sql) {
		t.Error("KOMITE_ID tidak ber-index UNIK")
	}
}

// ADR-U-0006: identitas seluruh tabel T_CLAIMLF_* dari sequence.
func TestSequenceUntukTabelYangMemakainya(t *testing.T) {
	sql := gabungSemua(t)
	for _, s := range []string{
		"SEQ_CLAIMLF_PLD", "SEQ_CLAIMLF_ADJ", "SEQ_CLAIMLF_SPR",
		"SEQ_CLAIMLF_SPR_RETRO", "SEQ_T_CLAIMLF_DOCUMENT",
	} {
		if !strings.Contains(sql, "CREATE SEQUENCE {SKEMA}."+s) {
			t.Errorf("sequence %s tidak dibuat", s)
		}
	}
}

// Identitas T_WORK_CLAIM adalah teks berformat, bukan angka sequence.
func TestIdentitasWorkClaimBerupaTeks(t *testing.T) {
	var isi string
	for nama, teks := range seluruhSQL(t, false) {
		if strings.HasPrefix(nama, "001_") {
			isi = strings.ToUpper(teks)
		}
	}
	if !regexp.MustCompile(`ID\s+VARCHAR2\(\d+\) NOT NULL`).MatchString(isi) {
		t.Error("T_WORK_CLAIM.ID bukan teks - tiket 14 menetapkan teks berformat CLM-/KMT-")
	}
	if regexp.MustCompile(`ID\s+NUMBER`).MatchString(isi) {
		t.Error("T_WORK_CLAIM.ID bertipe numerik")
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
		langkah, err := daftarMigrasi(mundur)
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

// Toleransi "objek sudah ada" hanya berlaku untuk galat yang memang berarti itu.
//
// Ini penggolong galat yang menentukan apakah migrasi meneruskan langkahnya
// atau berhenti. Menggolongkan terlalu longgar berarti menelan kerusakan
// sungguhan, jadi batasnya diuji dari kedua sisi.
func TestPenggolongGalatObjekSudahAda(t *testing.T) {
	// Tiga bentuk pembungkus yang berbeda. Yang diuji bukan kode galatnya saja
	// melainkan bahwa penggolong menemukannya di mana pun ia diletakkan driver -
	// telanjang, berawalan, dan terbungkus galat lain.
	harusYa := []string{
		"ORA-00955: name is already used by an existing object",
		"oci: ORA-00955: name is already used by an existing object",
		"repository: migrasi 001_t_work_claim.sql: go-ora: " +
			"ORA-00955: name is already used by an existing object",
	}
	harusTidak := []string{
		// ⛔ Ralat ronde 3. ORA-02264 dulu ada di daftar harusYa, dan test ini
		// justru MENGUNCI perilaku yang salah. ORA-02264 berarti nama
		// constraint terpakai, dan Oracle baru memeriksanya saat tabelnya belum
		// ada - jadi ia berarti tabelnya TIDAK terbuat, bukan sudah ada.
		"ORA-02264: name already used by an existing constraint",
		"ORA-00942: table or view does not exist",
		"ORA-01400: cannot insert NULL",
		"ORA-00972: identifier is too long",
		"sambungan terputus",
	}
	for _, p := range harusYa {
		if !sudahAda(errors.New(p)) {
			t.Errorf("sudahAda(%q) = false, seharusnya true", p)
		}
	}
	for _, p := range harusTidak {
		if sudahAda(errors.New(p)) {
			t.Errorf("sudahAda(%q) = true, seharusnya false", p)
		}
	}
	if sudahAda(nil) {
		t.Error("sudahAda(nil) = true, seharusnya false")
	}
}

// Hanya pernyataan CREATE yang boleh dilewati saat objeknya sudah ada.
func TestHanyaCreateYangBolehDilewati(t *testing.T) {
	kasus := map[string]bool{
		"CREATE TABLE {skema}.T_X (ID VARCHAR2(32))": true,
		"  create index {skema}.IX_X on ...":         true,
		"CREATE SEQUENCE {skema}.SEQ_X":              true,
		"ALTER TABLE {skema}.T_X ADD (Y DATE)":       false,
		"DROP TABLE {skema}.T_X":                     false,
		"INSERT INTO {skema}.T_MIGRASI VALUES (1)":   false,
	}
	for q, harap := range kasus {
		if pernyataanBuat(q) != harap {
			t.Errorf("pernyataanBuat(%q) = %v, seharusnya %v", q, !harap, harap)
		}
	}
}

// Ringkasan pernyataan menyebut objeknya tanpa menyalin seluruh DDL.
func TestRingkasPernyataanPendek(t *testing.T) {
	q := "CREATE TABLE {skema}.T_WORK_CLAIM (\n  ID VARCHAR2(32) NOT NULL,\n  LINI VARCHAR2(16)\n)"
	got := ringkasPernyataan(q)
	if strings.Contains(got, "VARCHAR2") {
		t.Errorf("ringkasan masih memuat badan DDL: %q", got)
	}
	if !strings.Contains(got, "T_WORK_CLAIM") {
		t.Errorf("ringkasan tidak menyebut objeknya: %q", got)
	}
}

// Nama objek terbaca dari tiap bentuk pernyataan CREATE yang dipakai migrasi.
//
// Pembacaan ini yang menentukan objek mana keberadaannya dibuktikan sesudah
// sebuah CREATE dilewati. Salah baca berarti pembuktiannya menanyakan objek
// yang keliru - dan itu sama buruknya dengan tidak membuktikan sama sekali.
func TestNamaObjekDibuatTerbaca(t *testing.T) {
	kasus := map[string]string{
		"CREATE TABLE {skema}.T_WORK_CLAIM (\n  ID VARCHAR2(32))":         "T_WORK_CLAIM",
		"CREATE INDEX {skema}.IX_PLD_CLAIM_ID ON {skema}.T_X (CLAIM_ID)":  "IX_PLD_CLAIM_ID",
		"CREATE UNIQUE INDEX {skema}.UX_ADJ_KOMITE_ID ON {skema}.T_Y (A)": "UX_ADJ_KOMITE_ID",
		"CREATE SEQUENCE {skema}.SEQ_CLAIMLF_PLD START WITH 1":            "SEQ_CLAIMLF_PLD",
		"create table {skema}.t_kecil (id number(19))":                    "T_KECIL",
		"INSERT INTO {skema}.T_MIGRASI (NAMA) VALUES (:1)":                "",
		"DROP TABLE {skema}.T_WORK_CLAIM CASCADE CONSTRAINTS":             "",
	}
	for q, mau := range kasus {
		if got := namaObjekDibuat(q); got != mau {
			t.Errorf("namaObjekDibuat(%.50s) = %q, mau %q", q, got, mau)
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
	langkah, err := daftarMigrasi(false)
	if err != nil {
		t.Fatal(err)
	}
	diperiksa := 0
	for _, m := range langkah {
		for _, p := range m.Pernyataan {
			if !pernyataanBuat(p) {
				continue
			}
			diperiksa++
			if namaObjekDibuat(p) == "" {
				t.Errorf("%s: nama objek tidak terbaca dari %q", m.Nama, ringkasPernyataan(p))
			}
		}
	}
	// Angkanya dikunci: kalau pemisah pernyataan rusak lagi, cacahnya anjlok
	// dan test ini gagal alih-alih diam-diam memeriksa lebih sedikit.
	const mau = 19
	if diperiksa != mau {
		t.Errorf("pernyataan CREATE diperiksa %d, mau %d", diperiksa, mau)
	}
}

// Pembanding bentuk tabel menyebut kedua arah selisihnya.
//
// Ini bagian MURNI dari pra-terbang butir x: ia tidak menyentuh Oracle sama
// sekali, sehingga perilakunya terkunci di setiap `go test` biasa. Yang
// dibandingkan hanya NAMA kolom - tipe sengaja tidak, sebab selisih tipe belum
// tentu salah dan akan menghasilkan penolakan palsu.
func TestSelisihKolomMenyebutKeduaArah(t *testing.T) {
	kasus := []struct {
		nama          string
		ddl, katalog  []string
		kurang, lebih []string
	}{
		{"sama persis",
			[]string{"ID", "NAMA"}, []string{"ID", "NAMA"}, nil, nil},
		{"urutan berbeda tetap sama",
			[]string{"ID", "NAMA"}, []string{"NAMA", "ID"}, nil, nil},
		{"huruf kecil di katalog tetap sama",
			[]string{"ID", "NAMA"}, []string{"id", "nama"}, nil, nil},
		{"katalog kurang satu kolom",
			[]string{"ID", "NAMA", "TGL"}, []string{"ID", "NAMA"},
			[]string{"TGL"}, nil},
		{"katalog punya kolom yang tidak diminta",
			[]string{"ID"}, []string{"ID", "IDPEGA", "NAMAFILE"},
			nil, []string{"IDPEGA", "NAMAFILE"}},
		{"berselisih di kedua arah",
			[]string{"ID", "TGL"}, []string{"ID", "IDPEGA"},
			[]string{"TGL"}, []string{"IDPEGA"}},
		{"tabel katalog kosong",
			[]string{"ID"}, nil, []string{"ID"}, nil},
	}
	for _, k := range kasus {
		t.Run(k.nama, func(t *testing.T) {
			kurang, lebih := SelisihKolom(k.ddl, k.katalog)
			if !samaDaftar(kurang, k.kurang) {
				t.Errorf("kurang = %v, mau %v", kurang, k.kurang)
			}
			if !samaDaftar(lebih, k.lebih) {
				t.Errorf("lebih = %v, mau %v", lebih, k.lebih)
			}
		})
	}
}

func samaDaftar(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// KolomCreateTable membaca nama dan kolom dari setiap CREATE TABLE migrasi.
//
// Cacahnya dikunci: delapan CREATE TABLE. Pernyataan yang BUKAN CREATE TABLE -
// CREATE INDEX dan CREATE SEQUENCE - harus mengembalikan nama kosong, kalau
// tidak pra-terbang akan mencari "bentuk" sebuah sequence.
func TestKolomCreateTableMembacaSeluruhTabel(t *testing.T) {
	langkah, err := daftarMigrasi(false)
	if err != nil {
		t.Fatal(err)
	}
	tabel, bukanTabel := 0, 0
	for _, m := range langkah {
		for _, p := range m.Pernyataan {
			nama, kolom := KolomCreateTable(p)
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
	// 7 tabel + 5 sequence + 7 index = 19 pernyataan CREATE, cocok dengan
	// cacah yang dikunci TestSeluruhCreateDapatDibacaNamanya.
	const mauTabel = 7
	if tabel != mauTabel {
		t.Errorf("CREATE TABLE terbaca %d, mau %d", tabel, mauTabel)
	}
	if bukanTabel == 0 {
		t.Error("nol pernyataan bukan-tabel; CREATE INDEX dan SEQUENCE seharusnya ada")
	}
}
