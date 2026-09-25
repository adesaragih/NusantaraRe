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
		"T_CLAIMLF_ADJUSTMENT_SPREADING_RETRO",
		"DOCUMENT_CLAIM",
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

// Kaskade hanya pada relasi 3, 4, 5, 6. DOCUMENT_CLAIM (relasi 7) ditangani
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

// ADR-U-0006: identitas T_CLAIMLF_* dan DOCUMENT_CLAIM dari sequence.
func TestSequenceUntukTabelYangMemakainya(t *testing.T) {
	sql := gabungSemua(t)
	for _, s := range []string{
		"SEQ_CLAIMLF_PLD", "SEQ_CLAIMLF_ADJ", "SEQ_CLAIMLF_SPR",
		"SEQ_CLAIMLF_SPR_RETRO", "SEQ_DOCUMENT_CLAIM",
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
