//go:build db

package repository_test

// Untuk apa berkas ini: SATU PINTU pembuat TIRUAN UJI (keputusan WO U1,
// PROMPT putaran 3 bab 2, 04-10-2026).
//
// Modul ini membaca/menulis beberapa tabel dan view WARISAN `POOLDATA` yang
// TIDAK dibuat migrasinya (320-327 hanya delapan tabel diagram). Skema uji
// yang dipasang `skemauji.Pasang` hanya memuat hasil migrasi, sehingga uji db
// yang menyentuh objek warisan membuat TIRUAN sementara - hanya di skema uji,
// bukan migrasi, bukan tabel aplikasi. Syarat U1, ditegakkan di `buatTiruan`
// dan diperiksa statik oleh `TestTiruanUjiLewatSatuPintuBerpagar`
// (tiruan_penjaga_test.go, tanpa Oracle):
//
//  1. menolak `POOLDATA` (dan setiap skema bernama warisan) dengan penjaga yang
//     SAMA dengan `uji/skemauji` - `config.PagarSkemaUji`, fungsi yang dipanggil
//     `skemauji.pagarSkemaUji` dan diuji `uji/skemauji/pagar_test.go` - SEBELUM
//     satu pun DDL dijalankan;
//  2. dibuat HANYA bila objeknya belum ada (objek yang disediakan DBA di skema
//     uji dipakai apa adanya, tidak pernah dibuang);
//  3. dibuang sesudah uji (`t.Cleanup` ... `DROP TABLE ... PURGE`);
//  4. setiap nama tiruan adalah konstanta berkomentar "tiruan uji, bukan tabel
//     aplikasi".

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"nusantarare/inti/backend/config"
)

// buatTiruan membuat tabel tiruan `skema.nama` berkolom `kolom` bila objek
// bernama itu belum ada, dan mendaftarkan pembuangannya. Mengembalikan true
// bila tiruan dibuat uji ini; false bila objek sungguhan sudah ada.
func buatTiruan(t *testing.T, ctx context.Context, sqlDB *sql.DB, skema, nama string, kolom []string) bool {
	t.Helper()
	// U1: penjaga yang SAMA dengan uji/skemauji - menolak POOLDATA meski
	// ORACLE_SKEMA_UJI=true, dan menolak skema yang belum diakui skema uji.
	if err := config.PagarSkemaUji(skema, os.Getenv(config.EnvSkemaUji)); err != nil {
		t.Fatalf("tiruan uji %s ditolak: %v", nama, err)
	}
	var ada int
	if err := sqlDB.QueryRowContext(ctx, `SELECT COUNT(*) FROM ALL_OBJECTS WHERE OWNER = UPPER(:1) AND OBJECT_NAME = :2`,
		skema, nama).Scan(&ada); err != nil {
		t.Fatal(err)
	}
	if ada > 0 {
		return false
	}
	if _, err := sqlDB.ExecContext(ctx, fmt.Sprintf(`CREATE TABLE %s.%s (%s)`, skema, nama, strings.Join(kolom, ", "))); err != nil {
		t.Fatalf("tiruan uji %s: %v", nama, err)
	}
	t.Cleanup(func() { _, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`DROP TABLE %s.%s PURGE`, skema, nama)) })
	return true
}

// tiruan uji, bukan tabel aplikasi: riwayat catatan usulan WARISAN yang ditulis
// `CatatUsulan` dan DIBACA `BacaHalaman` (K4) - tanpa objek ini setiap uji db
// yang membuka halaman gagal ORA-00942, bukan menguji klaimnya.
const tabelRiwayatProduksi = "HISTORYAKSEPTASIPRODUCTION"

// siapkanRiwayatProduksi - 15 kolom = INSERT `RDBList/InsertViewSuggest_SQL`
// (IDPEGA ... PERCENT_RNM). `TGL_INP` DATE (`To_date(...)` di rule);
// `KETERANGAN` 4000 (`substr({CARI10},0,3990)`); tipe kolom lain belum dicek
// ke katalog (PERMINTAAN C8) ⇒ VARCHAR2(1000) seperti kolom view warisan.
// `NOURUT` sengaja TEKS: pembacaan `ORDER BY TO_NUMBER(NOURUT)` harus benar
// juga bila kolomnya VARCHAR2 ("10" < "2" bila diurut sebagai teks).
func siapkanRiwayatProduksi(t *testing.T, ctx context.Context, sqlDB *sql.DB, skema string) {
	t.Helper()
	kolom := []string{"IDPEGA VARCHAR2(1000)", "TYPE_POLIS VARCHAR2(1000)", "NOURUT VARCHAR2(1000)",
		"POSISI VARCHAR2(1000)", "PIC VARCHAR2(1000)", "TGL_INP DATE", "DIV VARCHAR2(1000)", "TYPE VARCHAR2(1000)",
		"PUTARAN VARCHAR2(1000)", "APPROVAL VARCHAR2(1000)", "KETERANGAN VARCHAR2(4000)", "AKSES_LOGIN VARCHAR2(1000)",
		"B2B VARCHAR2(1000)", "BUSINESS_CODE VARCHAR2(1000)", "PERCENT_RNM VARCHAR2(1000)"}
	if !buatTiruan(t, ctx, sqlDB, skema, tabelRiwayatProduksi, kolom) {
		// tabel DBA: baris UJI- dibuang sesudah uji, isi lain tidak disentuh.
		t.Cleanup(func() {
			_, _ = sqlDB.ExecContext(ctx, fmt.Sprintf(`DELETE FROM %s.%s WHERE IDPEGA LIKE '%%UJI-NB-%%'`, skema, tabelRiwayatProduksi))
		})
	}
}

// tiruan uji, bukan tabel aplikasi: tabel DASAR `T_GENERAL_POLIS` bersama
// FacIn + Treaty In (keputusan WO 04-10-2026). Ia dibuat migrasi `nbfacin`
// `182_t_general_polis`, yang belum ada di repo; migrasi 320 modul ini hanya
// `ALTER TABLE ... ADD (` kolom Treaty atasnya. Tanpa tabel dasar,
// `skemauji.Pasang` gagal di 320 (ORA-00942). Bila skema uji sudah memuatnya
// (182 sudah di repo/dijalankan DBA), tabel itu dipakai apa adanya.
const tabelGeneralPolisDasar = "T_GENERAL_POLIS"

// siapkanGeneralPolisDasar - tujuh kolom dasar 182 (ID VARCHAR2(32) PK, IDPEGA
// VARCHAR2(50) - keputusan WO 04-10-2026); tipe lima kolom FacIn lainnya dari
// draf `nbfacin` (`services/loader/skema_gen.go`). DIPANGGIL SEBELUM
// `skemauji.Pasang`: pembuangannya (t.Cleanup) berjalan SESUDAH
// `skemauji.Bongkar` membuang tujuh tabel anak dan kolom Treaty (320_down).
func siapkanGeneralPolisDasar(t *testing.T, ctx context.Context, sqlDB *sql.DB, skema string) {
	t.Helper()
	buatTiruan(t, ctx, sqlDB, skema, tabelGeneralPolisDasar, []string{
		"ID VARCHAR2(32) NOT NULL PRIMARY KEY", "IDPEGA VARCHAR2(50)", "COB_GROUP VARCHAR2(20)",
		"START_DATE_TIME VARCHAR2(30)", "OFFERING_DATE VARCHAR2(30)", "END_DATE_TIME VARCHAR2(30)", "FOLLOWING VARCHAR2(50)"})
}
