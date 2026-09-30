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

	"nusantarare/inti/backend/migrasi"
)

func seluruhSQL(t *testing.T, mundur bool) map[string]string {
	t.Helper()
	langkah, err := migrasi.Daftar(mundur, berkasMigrasi)
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

// milikClaimLife menjawab apakah berkas migrasi itu milik modul Claim Life.
//
// ⛔ Batasnya NOMOR, dan itu keputusan yang tercatat: Claim Life memakai
// 001-049, PremiumList Life mulai 050. Memisahkan lewat nama tabel akan
// gagal pada tabel yang namanya tidak menyebut modulnya.
func milikClaimLife(nama string) bool {
	return nama < "050_"
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
