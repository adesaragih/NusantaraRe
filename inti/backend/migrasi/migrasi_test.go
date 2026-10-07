package migrasi

import (
	"errors"
	"strings"
	"testing"
)

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
		if PernyataanBuat(q) != harap {
			t.Errorf("pernyataanBuat(%q) = %v, seharusnya %v", q, !harap, harap)
		}
	}
}

// Ringkasan pernyataan menyebut objeknya tanpa menyalin seluruh DDL.
func TestRingkasPernyataanPendek(t *testing.T) {
	q := "CREATE TABLE {skema}.T_WORK_CLAIM (\n  ID VARCHAR2(32) NOT NULL,\n  LINI VARCHAR2(16)\n)"
	got := RingkasPernyataan(q)
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
		if got := NamaObjekDibuat(q); got != mau {
			t.Errorf("namaObjekDibuat(%.50s) = %q, mau %q", q, got, mau)
		}
	}
}

// View yang diganti tabel bernama sama dikenali hanya bila DROP VIEW-nya
// datang SEBELUM CREATE TABLE di langkah yang sama (masterprovince 880, dulu masterdata 760).
func TestViewDibongkarDulu(t *testing.T) {
	p := pecahPernyataan("CREATE TABLE {skema}.P_SALIN (\n  ID VARCHAR2(10)\n)\n/\n" +
		"-- komentar\nDROP VIEW {skema}.Province\n/\nCREATE TABLE {skema}.PROVINCE (\n  ID VARCHAR2(10)\n)\n/\n" +
		"DROP VIEW {skema}.CITY\n/\n")
	if len(p) != 4 {
		t.Fatalf("%d pernyataan: %q", len(p), p)
	}
	kasus := []struct {
		i    int
		nama string
		mau  bool
	}{
		{2, "PROVINCE", true},  // DROP VIEW sebelumnya, huruf kecil tetap sama
		{0, "P_SALIN", false},  // tidak ada yang dibongkar sebelumnya
		{2, "P_SALIN", false},  // view lain
		{1, "PROVINCE", false}, // DROP VIEW itu sendiri, bukan sebelumnya
		{3, "CITY", false},     // DROP VIEW SESUDAH tidak dihitung
		{4, "CITY", true},
	}
	for _, k := range kasus {
		if got := ViewDibongkarDulu(p, k.i, k.nama); got != k.mau {
			t.Errorf("ViewDibongkarDulu(%d, %s) = %v, mau %v", k.i, k.nama, got, k.mau)
		}
	}
	if ViewDibongkarDulu([]string{"DROP VIEW {skema}.X CASCADE CONSTRAINTS"}, 1, "X") {
		t.Error("DROP VIEW berekor tidak boleh dikenali")
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
