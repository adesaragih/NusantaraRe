package backend

// Uji invarian atas TEKS DDL - bukan atas Oracle (`L-3`).

import (
	"io/fs"
	"regexp"
	"strings"
	"testing"
)

func seluruhMigrasi(t *testing.T) map[string]string {
	t.Helper()
	entri, err := fs.ReadDir(berkasMigrasi, "migrations")
	if err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	for _, e := range entri {
		isi, err := fs.ReadFile(berkasMigrasi, "migrations/"+e.Name())
		if err != nil {
			t.Fatal(err)
		}
		out[e.Name()] = string(isi)
	}
	if len(out) == 0 {
		t.Fatal("nol berkas migrasi terbaca; pembacanya yang rusak")
	}
	return out
}

// tanpaKomentar membuang baris komentar `--`, menyisakan PERNYATAAN SQL saja -
// berkas migrasi modul ini MENJELASKAN keputusannya di dalam komentar, dan
// sapuan atas teks mentah akan menemukan kata-katanya di prosanya sendiri.
func tanpaKomentar(sql string) string {
	var b strings.Builder
	for _, baris := range strings.Split(sql, "\n") {
		if strings.HasPrefix(strings.TrimSpace(baris), "--") {
			continue
		}
		b.WriteString(baris)
		b.WriteByte('\n')
	}
	return b.String()
}

func majuSaja(t *testing.T) map[string]string {
	t.Helper()
	out := map[string]string{}
	for nama, isi := range seluruhMigrasi(t) {
		if !strings.HasSuffix(nama, "_down.sql") {
			out[nama] = isi
		}
	}
	return out
}

// ⛔ Modul ini TIDAK membuat satu tabel pun, dan itu keputusan - bukan keadaan
// sementara. Model datanya satu dengan Treaty In; yang dipisahkan 25-09-2026
// hanya papan tiketnya. Tabel baru di sini berarti model datanya bercabang.
func TestNolTabelBaru(t *testing.T) {
	for nama, isi := range majuSaja(t) {
		if strings.Contains(strings.ToUpper(tanpaKomentar(isi)), "CREATE TABLE") {
			t.Errorf("%s membuat tabel; modul ini hanya mengubah tabel milik `treatyin` "+
				"(model data SATU, `treaty-in/SPEC-MODEL-DATA.md`)", nama)
		}
	}
}

// Tiket 01: kolom dasar berdiri dengan kunci asingnya, dan kunci asing itu
// menunjuk VERSI_KONTRAK - tabelnya sendiri.
func TestKolomDasarBerdiriDenganKunciAsingnya(t *testing.T) {
	isi, ok := majuSaja(t)["440_versi_kontrak_dasar.sql"]
	if !ok {
		t.Fatal("migrasi 440 tidak ada")
	}
	sql := tanpaKomentar(isi)
	if !strings.Contains(sql, "ID_VERSI_KONTRAK_DASAR  NUMBER(19)") {
		t.Error("tiket 01: kolom ID_VERSI_KONTRAK_DASAR tidak ditambahkan")
	}
	if !strings.Contains(sql, "CONSTRAINT FK_VERSI_KONTRAK_DASAR") ||
		!strings.Contains(sql, "REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK)") {
		t.Error("tiket 01: kunci asing ke VERSI_KONTRAK tidak dipasang")
	}
	// Kolomnya WAJIB boleh kosong: versi pertama sebuah kontrak berdasar
	// kosong (bahan to-spec B-3), dan baris warisan boleh berdasar kosong
	// sampai tiket 10.
	if strings.Contains(sql, "ID_VERSI_KONTRAK_DASAR  NUMBER(19) NOT NULL") {
		t.Error("tiket 01: kolom dasar NOT NULL; versi pertama berdasar KOSONG")
	}
}

// Tiket 05: langkah PERLUAS - nomor urut versi menjadi boleh kosong.
func TestNomorUrutVersiDilonggarkan(t *testing.T) {
	isi, ok := majuSaja(t)["441_nomor_urut_versi_boleh_kosong.sql"]
	if !ok {
		t.Fatal("migrasi 441 tidak ada")
	}
	sql := tanpaKomentar(isi)
	if !strings.Contains(sql, "MODIFY") || !strings.Contains(sql, "NOMOR_URUT_VERSI  NUMBER(10) NULL") {
		t.Error("tiket 05: NOMOR_URUT_VERSI tidak dilonggarkan menjadi NULL")
	}
	// Lebarnya tidak boleh berubah diam-diam: MODIFY yang menyebut tipe lain
	// mengubah dua hal sekaligus, dan yang kedua tidak pernah ditinjau.
	if strings.Contains(sql, "NUMBER(9)") || strings.Contains(sql, "NUMBER(19)") {
		t.Error("tiket 05: MODIFY mengubah LEBAR kolom, bukan hanya keterisiannya")
	}
}

// ADR-0056 (K-4) dan ADR-U-0029: nol trigger, nol procedure, nol COMMIT.
//
// Tiket 03 MEMINTA trigger INV-54; keputusannya ada di
// `docs/KEPUTUSAN-TIKET-02-03.md` §2 - ditegakkan di services, sejalan dengan
// INV-53 yang sebentuk. Uji ini menjaga keputusan itu.
func TestNolTriggerProcedureDanCommit(t *testing.T) {
	for nama, isi := range seluruhMigrasi(t) {
		atas := strings.ToUpper(tanpaKomentar(isi))
		for _, terlarang := range []string{"CREATE TRIGGER", "CREATE OR REPLACE TRIGGER", "CREATE PROCEDURE", "CREATE FUNCTION", "\nCOMMIT"} {
			if strings.Contains(atas, terlarang) {
				t.Errorf("%s memuat %q; aturan bisnis tidak turun ke basis data (ADR-0056)", nama, terlarang)
			}
		}
	}
}

// gabungan menyatukan seluruh langkah maju menjadi satu teks.
func gabungan(t *testing.T) string {
	t.Helper()
	var b strings.Builder
	for _, isi := range majuSaja(t) {
		b.WriteString(isi)
	}
	return b.String()
}

// INV-18 - perilaku hapus DITETAPKAN SADAR, diadu dengan ERD.md §2.
//
// ⚠️ RALAT 2 Oktober 2026. Uji ini pernah bernama TestNolKaskadeHapus dan
// menuntut nol `ON DELETE` dengan alasan "kebijakan modul ini MENOLAK". Itu
// membaca INV-18 terbalik — invariannya menuntut keputusan yang DINYATAKAN,
// bukan bawaan yang kebetulan cocok, dan sumbernya `ERD.md` §2.
//
// Modul ini punya SATU kunci asing sendiri, dan §2.2 menyatakannya
// [hapus: tolak]: "menghapus versi dasar akan membuat seluruh baris selisih
// kehilangan artinya." TOLAK diwujudkan dengan tidak menulis klausa ON DELETE.
func TestPerilakuHapusSesuaiERD(t *testing.T) {
	pola := regexp.MustCompile(`(?is)CONSTRAINT\s+(FK_\w+)\s+FOREIGN KEY\s*\([^)]*\)\s*REFERENCES\s+\{skema\}\.\w+\s*\([^)]*\)([^,\n]*)`)
	cocok := pola.FindAllStringSubmatch(tanpaKomentar(gabungan(t)), -1)
	if len(cocok) != 1 {
		t.Fatalf("mau tepat 1 kunci asing di modul ini, dapat %d", len(cocok))
	}
	nama, ekor := cocok[0][1], strings.ToUpper(strings.TrimSpace(cocok[0][2]))
	if nama != "FK_VERSI_KONTRAK_DASAR" {
		t.Errorf("kunci asing tak terduga: %s", nama)
	}
	if ekor != "" {
		t.Errorf("%s berperilaku %q; ERD.md §2.2 menuntut TOLAK, yaitu tanpa klausa ON DELETE", nama, ekor)
	}
}

// Nol kaskade, dan `MODUL.md` menyatakannya dengan tabel kosong — sehingga
// `TestKaskadeHanyaPadaRelasiTerdaftar` menolak yang pertama menambahkannya.
func TestNolKaskadeDiModulIni(t *testing.T) {
	for nama, isi := range majuSaja(t) {
		if strings.Contains(strings.ToUpper(tanpaKomentar(isi)), "ON DELETE") {
			t.Errorf("%s memuat ON DELETE; modul ini menyatakan NOL kaskade di MODUL.md", nama)
		}
	}
}

// Yang TIDAK diuji di sini, sebab penjaga inti sudah menegakkannya atas SELURUH
// modul - menyalinnya ke sini berarti dua tempat yang harus sepakat:
//
//   jalur mundur tiap langkah   -> penjaga `TestSetiapLangkahPunyaJalurMundur`
//   nomor di rentang modulnya   -> penjaga `TestSetiapMigrasiDiRentangAtauSlotModulnya`
//   slot menu hanya menyentuh
//   baris menu modulnya sendiri -> penjaga `TestSlotMenuHanyaMenyentuhMenuModulnya`
//
// Yang DI SINI adalah yang khusus modul ini: invarian bernomor dan pernyataan
// keputusan yang penjaga umum tidak dapat mengetahuinya.
