package repository

// Nol kata cadangan Oracle sebagai nama kolom telanjang - pl6, 28-09-2026.
//
// ⛔ SEBAB PENJAGA INI LAHIR. Migrasi `056` memuat kolom bernama `INITIAL`.
// Ia lolos seluruh penjaga yang ada, lolos seluruh uji, lolos tinjauan - dan
// GAGAL di Oracle sungguhan saat work owner menjalankan `-migrate`:
//
//	SELECT 1 AS INITIAL   FROM DUAL  -> ORA-00923
//	SELECT 1 AS "INITIAL" FROM DUAL  -> lolos (berkutip)
//	SELECT 1 AS NO        FROM DUAL  -> lolos
//
// Sebelas langkah lain di jalan yang sama lolos, jadi bukan penanda `{skema}`
// yang salah. Yang membuatnya lolos bukan kecerobohan melainkan **ketiadaan
// pemeriksa**: nol uji di repositori ini pernah menanyakan apakah sebuah nama
// kolom boleh berdiri telanjang di Oracle.
//
// ⚠️ Daftarnya SENGAJA TIDAK LENGKAP, dan itu dinyatakan. Oracle punya
// ratusan kata cadangan; yang di bawah ini adalah yang MASUK AKAL muncul
// sebagai nama kolom di domain reasuransi - kata yang tidak akan pernah
// seseorang ketik sebagai nama kolom tidak perlu dijaga. Daftar yang
// berpura-pura lengkap lebih berbahaya daripada daftar yang mengaku parsial:
// yang pertama membuat orang berhenti berpikir.

import (
	"strings"
	"testing"
)

// kataCadanganOracle adalah kata yang TIDAK boleh berdiri telanjang sebagai
// nama kolom.
//
// Sumbernya daftar reserved word Oracle; yang disaring ke sini hanya yang
// masuk akal muncul di domain ini. Setiap entri dapat diuji sendiri dengan
// `SELECT 1 AS <kata> FROM DUAL`.
var kataCadanganOracle = map[string]bool{
	"INITIAL":    true, // yang benar-benar menggigit, 28-09-2026
	"LEVEL":      true,
	"SIZE":       true,
	"DATE":       true,
	"NUMBER":     true,
	"COMMENT":    true,
	"ORDER":      true,
	"GROUP":      true,
	"CHECK":      true,
	"DEFAULT":    true,
	"ACCESS":     true,
	"AUDIT":      true,
	"CLUSTER":    true,
	"COLUMN":     true,
	"OPTION":     true,
	"ROW":        true,
	"ROWID":      true,
	"SESSION":    true,
	"SHARE":      true,
	"START":      true,
	"SUCCESSFUL": true,
	"SYNONYM":    true,
	"TABLE":      true,
	"UID":        true,
	"USER":       true,
	"VALIDATE":   true,
	"VALUES":     true,
	"VIEW":       true,
	"MODE":       true,
	"RESOURCE":   true,
	"ONLINE":     true,
	"OFFLINE":    true,
	"INCREMENT":  true,
	"MINUS":      true,
	"PRIOR":      true,
	"PUBLIC":     true,
	"RAW":        true,
	"RENAME":     true,
	"LONG":       true,
	"FILE":       true,
	"IMMEDIATE":  true,
	"INDEX":      true,
	"EXCLUSIVE":  true,
	"COMPRESS":   true,
	"CURRENT":    true,
	"DESC":       true,
	"ASC":        true,
}

func TestNolKataCadanganOracleSebagaiKolom(t *testing.T) {
	// ⛔ Memakai pengurai PRODUKSI `KolomCreateTable`, bukan regex kedua.
	// Pengurai kedua adalah definisi kedua tentang "apa itu kolom", dan yang
	// kedua akan diam-diam berbeda - lalu penjaga ini menjaga sesuatu yang
	// bukan kolom yang sebenarnya dibuat.
	berkas := seluruhSQL(t, false)
	diperiksa := 0
	for nama, teks := range berkas {
		for _, pernyataan := range strings.Split(teks, "\n/") {
			_, kolomTabel := KolomCreateTable(pernyataan)
			for _, kolom := range kolomTabel {
				kolom = strings.ToUpper(strings.TrimSpace(kolom))
				diperiksa++
				if !kataCadanganOracle[kolom] {
					continue
				}
				t.Errorf("%s: kolom %q adalah kata cadangan Oracle dan berdiri "+
					"telanjang. Ia lolos seluruh uji di sini dan GAGAL di Oracle "+
					"sungguhan (ORA-00923), jauh dari orang yang menulisnya. "+
					"Beri nama lain - pola modul ini menambahkan konteksnya, "+
					"misalnya INITIAL -> INITIAL_SUGGEST.", nama, kolom)
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol definisi kolom terbaca; pembacanya yang rusak, bukan kodenya")
	}
	t.Logf("%d definisi kolom diperiksa terhadap %d kata cadangan",
		diperiksa, len(kataCadanganOracle))
}

// TestPenjagaKataCadanganMasihMenggigit membuktikan polanya benar-benar
// membaca kolom.
//
// ⛔ Penjaga yang polanya tidak pernah cocok adalah penjaga yang selalu
// hijau. Uji ini memberinya DDL buatan yang jelas salah, dan menuntutnya
// menemukan kolomnya - tanpa menyentuh berkas migrasi mana pun.
func TestPenjagaKataCadanganMasihMenggigit(t *testing.T) {
	const buruk = `CREATE TABLE {skema}.T_UJI (
  ID       VARCHAR2(32) NOT NULL,
  INITIAL  VARCHAR2(255),
  LEVEL    NUMBER(5),
  AMAN_    VARCHAR2(10)
)`
	var temuan []string
	_, kolom := KolomCreateTable(buruk)
	if len(kolom) == 0 {
		t.Fatal("pengurai tidak membaca satu pun kolom; polanya yang rusak")
	}
	for _, k := range kolom {
		if kataCadanganOracle[strings.ToUpper(strings.TrimSpace(k))] {
			temuan = append(temuan, k)
		}
	}
	if len(temuan) != 2 {
		t.Fatalf("temuan = %v, mau tepat dua (INITIAL dan LEVEL)", temuan)
	}
	// Dan kolom yang AMAN tidak ikut tertuduh - penjaga yang menuduh hal
	// yang benar akan dilonggarkan orang, bukan dipatuhi.
	for _, x := range temuan {
		if strings.HasPrefix(strings.ToUpper(x), "AMAN") {
			t.Errorf("kolom aman %q ikut tertuduh", x)
		}
	}
}
