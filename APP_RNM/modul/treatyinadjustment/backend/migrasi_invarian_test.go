package backend

// Uji invarian atas TEKS DDL - bukan atas Oracle (`L-3`).

import (
	"io/fs"
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

// INV-18: nol kaskade hapus - bawaan Oracle MENOLAK, dan menolak yang dikehendaki.
func TestNolKaskadeHapus(t *testing.T) {
	for nama, isi := range majuSaja(t) {
		if strings.Contains(strings.ToUpper(tanpaKomentar(isi)), "ON DELETE") {
			t.Errorf("INV-18: %s memuat ON DELETE; kebijakan modul ini MENOLAK, tanpa klausa", nama)
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
