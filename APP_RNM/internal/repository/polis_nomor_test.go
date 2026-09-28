package repository

// Penyimpanan `PL_NUMBER` - tiket 03 PremiumList Life. TANPA Oracle.

import (
	"errors"
	"strings"
	"testing"
)

// TestTulisNomorPLTidakDapatMenimpaNomorYangSudahAda - gerbang kedua.
//
// ⛔ "Lahir sekali" berdiri di DUA tempat: di layanan, dan di kalimat `WHERE`
// query ini. Gerbang yang hanya ada di layanan adalah gerbang yang hilang
// ketika kelak ada pemanggil kedua - dan pemanggil kedua selalu datang.
func TestTulisNomorPLTidakDapatMenimpaNomorYangSudahAda(t *testing.T) {
	q := sqlTulisNomorPL("SKEMAUJI.T_PREMIUM_LIST_DETAIL")
	if !strings.Contains(q, "PL_NUMBER IS NULL") {
		t.Errorf("query tulis nomor tidak menyaring baris yang sudah bernomor:\n%s", q)
	}
	if !strings.Contains(q, "PREMIUM_LIST_ID = :2") {
		t.Errorf("query tulis nomor tidak dibatasi satu polis:\n%s", q)
	}
	if err := PeriksaSQL(q); err != nil {
		t.Errorf("%v\n%s", err, q)
	}
}

// TestUtuhMembedakanPolisSetengahBernomor - lubang `MIN`/`MAX` yang melewati NULL.
//
// ⛔ SEJARAHNYA, dan ia hampir lolos. Gerbang lahir-sekali membaca `MIN` dan
// `MAX` `PL_NUMBER`. Oracle - dan `sql.NullString` di sisi kami - MELEWATI
// NULL, jadi polis yang separuh barisnya bernomor menghasilkan `MIN == MAX`,
// yaitu "sudah bernomor", dan `Terbitkan` kembali lebih awal. Baris yang
// kosong TIDAK PERNAH terisi.
//
// Komentar di `sqlKeadaanNomorPL` sudah menyebut bahwa selisih `COUNT(*)`
// dengan `COUNT(PL_NUMBER)` yang menunjukkannya - tetapi tidak satu pun kode
// pernah MEMBANDINGKAN keduanya. Alat ukur yang dibaca lalu diabaikan sama
// saja dengan alat ukur yang tidak ada.
//
// ⚠️ Keadaan ini bukan kerusakan: unggahan CSV (tiket 04) menambah peserta
// SESUDAH polis bernomor, dan peserta baru memang lahir tanpa nomor.
func TestUtuhMembedakanPolisSetengahBernomor(t *testing.T) {
	for _, k := range []struct {
		apa      string
		keadaan  KeadaanNomorPL
		bernomor bool
		utuh     bool
	}{
		{"belum bernomor sama sekali",
			KeadaanNomorPL{Nomor: "", CacahPeserta: 3, CacahBernomor: 0}, false, false},
		{"seluruh baris bernomor",
			KeadaanNomorPL{Nomor: "X", CacahPeserta: 3, CacahBernomor: 3}, true, true},
		{"SEPARUH bernomor - MIN==MAX menipu",
			KeadaanNomorPL{Nomor: "X", CacahPeserta: 5, CacahBernomor: 3}, true, false},
		{"satu baris, bernomor",
			KeadaanNomorPL{Nomor: "X", CacahPeserta: 1, CacahBernomor: 1}, true, true},
		{"nol baris sama sekali",
			KeadaanNomorPL{Nomor: "", CacahPeserta: 0, CacahBernomor: 0}, false, true},
	} {
		if got := k.keadaan.Bernomor(); got != k.bernomor {
			t.Errorf("%s: Bernomor()=%v, mau %v", k.apa, got, k.bernomor)
		}
		if got := k.keadaan.Utuh(); got != k.utuh {
			t.Errorf("%s: Utuh()=%v, mau %v", k.apa, got, k.utuh)
		}
	}
}

// TestDuaSebabNolBarisDibedakan - jawaban yang benar untuk keadaan yang benar.
//
// ⛔ Ada DUA cara sampai ke "nol baris tersentuh" saat menulis nomor, dan
// keduanya menuntut jawaban berbeda:
//
//	nol baris peserta                  -> polis memang belum siap dinomori
//	ada baris, semuanya sudah bernomor -> permintaan LAIN mendahului kita
//
// Menjawab keduanya dengan "belum punya peserta" mengirim orang mencari
// peserta yang sebenarnya ada, dan menyembunyikan satu-satunya petunjuk bahwa
// dua permintaan berjalan bersamaan atas polis yang sama.
func TestDuaSebabNolBarisDibedakan(t *testing.T) {
	if ErrPolisTanpaPeserta == nil || ErrNomorPLTerbitBersamaan == nil {
		t.Fatal("salah satu galat penulisan nomor hilang")
	}
	if errors.Is(ErrNomorPLTerbitBersamaan, ErrPolisTanpaPeserta) ||
		errors.Is(ErrPolisTanpaPeserta, ErrNomorPLTerbitBersamaan) {
		t.Error("kedua galat saling membungkus; pemanggil tidak dapat " +
			"membedakan polis yang belum punya peserta dari polis yang baru " +
			"saja dinomori permintaan lain")
	}
	// Kalimatnya masing-masing menyebut apa yang harus dikerjakan.
	if !strings.Contains(ErrPolisTanpaPeserta.Error(), "unggah rincian peserta") {
		t.Errorf("pesan %q tidak menyebut jalan keluarnya", ErrPolisTanpaPeserta)
	}
	if !strings.Contains(ErrNomorPLTerbitBersamaan.Error(), "baca ulang") {
		t.Errorf("pesan %q tidak menyebut jalan keluarnya", ErrNomorPLTerbitBersamaan)
	}
}

// TestRingkasPolisMembacaKeduaUjungNomor - MIN dan MAX, bukan MAX saja.
//
// ⛔ Satu polis bernomor satu. Pembacaan yang hanya mengambil `MAX` mengubur
// ketidaksepakatan antar baris di balik satu nilai yang kelihatan meyakinkan.
func TestRingkasPolisMembacaKeduaUjungNomor(t *testing.T) {
	q := sqlRingkasPolis("SKEMAUJI.T_PREMIUM_LIST", "SKEMAUJI.T_PREMIUM_LIST_DETAIL")
	for _, potong := range []string{
		"MIN(d.PL_NUMBER)", "MAX(d.PL_NUMBER)",
		"COUNT(d.ID)", "COUNT(d.PL_NUMBER)",
		"LEFT JOIN", "GROUP BY p.TYPE, p.BUSINESS_CODE",
	} {
		if !strings.Contains(q, potong) {
			t.Errorf("query ringkas tidak memuat %q:\n%s", potong, q)
		}
	}
	if err := PeriksaSQL(q); err != nil {
		t.Errorf("%v\n%s", err, q)
	}
	// Satu baris, satu perjalanan: kepala polis dan nomornya tidak dibaca
	// lewat dua query yang dapat melihat dua keadaan berbeda.
	if n := strings.Count(strings.ToUpper(q), "SELECT"); n != 1 {
		t.Errorf("query ringkas memuat %d SELECT; mau 1:\n%s", n, q)
	}
	// ⛔ `COUNT(*)` atas `LEFT JOIN` tanpa pasangan mencacah SATU, bukan nol -
	// dan itu menjawab "polis ini punya satu peserta" untuk polis yang belum
	// punya satu pun, yaitu tepat gerbang yang seharusnya menolak.
	if strings.Contains(q, "COUNT(*)") {
		t.Errorf("query ringkas memakai COUNT(*) atas LEFT JOIN:\n%s", q)
	}
}

// TestKeadaanNomorPLMembacaKeduaUjungJuga - jalur di dalam transaksi.
func TestKeadaanNomorPLMembacaKeduaUjungJuga(t *testing.T) {
	q := sqlKeadaanNomorPL("SKEMAUJI.T_PREMIUM_LIST_DETAIL")
	for _, potong := range []string{"MIN(d.PL_NUMBER)", "MAX(d.PL_NUMBER)"} {
		if !strings.Contains(q, potong) {
			t.Errorf("query keadaan tidak memuat %q:\n%s", potong, q)
		}
	}
	if err := PeriksaSQL(q); err != nil {
		t.Errorf("%v\n%s", err, q)
	}
}

// TestIdentitasPolisHanyaDuaKolom - bahan nomor, bukan seluruh header.
//
// ⚠️ Membawa seluruh kolom membuat perubahan tabel di hulu mengubah bentuk
// baris kami tanpa diminta - pelajaran yang sama dengan `SELECT *` di
// `GETTanggalClosing_SQL`.
func TestIdentitasPolisHanyaDuaKolom(t *testing.T) {
	q := sqlIdentitasPolis("SKEMAUJI.T_PREMIUM_LIST")
	if strings.Contains(q, "SELECT *") {
		t.Errorf("query identitas membawa seluruh kolom:\n%s", q)
	}
	if !strings.Contains(q, "p.TYPE") || !strings.Contains(q, "p.BUSINESS_CODE") {
		t.Errorf("query identitas tidak membaca TYPE dan BUSINESS_CODE:\n%s", q)
	}
	if err := PeriksaSQL(q); err != nil {
		t.Errorf("%v\n%s", err, q)
	}
}

// TestNolCommitDiQueryNomor - ADR-U-0029, dan procedure aslinya punya satu.
//
// ⛔ `GetSequenceNumber_SQL` baris 88 memuat `COMMIT;` di dalam blok SQL-nya.
// Kami TIDAK menirunya: batas transaksi milik Go, dan `COMMIT` di tengah teks
// SQL menutup transaksi yang bukan miliknya.
func TestNolCommitDiQueryNomor(t *testing.T) {
	for nama, q := range map[string]string{
		"tulis nomor": sqlTulisNomorPL("SKEMAUJI.T_PREMIUM_LIST_DETAIL"),
		"ringkas":     sqlRingkasPolis("SKEMAUJI.T_PREMIUM_LIST", "SKEMAUJI.T_PREMIUM_LIST_DETAIL"),
		"keadaan":     sqlKeadaanNomorPL("SKEMAUJI.T_PREMIUM_LIST_DETAIL"),
		"identitas":   sqlIdentitasPolis("SKEMAUJI.T_PREMIUM_LIST"),
		"grid":        sqlGridPeserta("SKEMAUJI.T_PREMIUM_LIST_DETAIL", nil),
	} {
		if strings.Contains(strings.ToUpper(q), "COMMIT") {
			t.Errorf("query %q memuat COMMIT:\n%s", nama, q)
		}
	}
}
