//go:build db

package repository_test

// Invarian Treaty In DIUJI TERHADAP ORACLE - bukan terhadap teks DDL.
//
// ⛔ Berkas ini ada karena `migrasi_invarian_test.go` TIDAK dapat membuktikan
// apa yang dibuktikan di sini. Berkas itu membaca BENTUK teks DDL; yang ini
// menyuruh Oracle MENOLAK, lalu memeriksa bahwa ia benar-benar menolak dan
// dengan constraint YANG MANA. Papan tiket menuntut keduanya: "sebuah tiket
// selesai bila perilakunya DAPAT GAGAL dan DAPAT DIPERIKSA".
//
// Tanpa ORACLE_DSN seluruh test di sini MELEWATI dengan pesan.
//
// ⭐ BERJALAN DI SKEMA APLIKASI - `POOLDATA` - dan TIDAK memakai `uji/skemauji`.
// Keputusan pemilik proses 3 Oktober 2026: skema uji tidak dipakai; yang dipakai
// POOLDATA, tempat ke-37 tabel Treaty In berdiri.
//
// ⛔ KENAPA `uji/skemauji` TIDAK BOLEH DIPAKAI DI SINI, dan ini bukan selera:
// `Bongkar` menjalankan `DROP TABLE <skema>.<nama> CASCADE CONSTRAINTS` atas
// daftar tetap yang memuat `OS_AKSEPTASI_KLAIM_LIFE`, `M_LIFE_PREMIUM_DETAIL`,
// `M_LIFE_PREMIUM_SUMMARY`, `RETROCESSIONLIFE`, `TREATYYEAR_LIFE`, dan tiruan
// Treaty Contract Out. Di skema uji itu tiruan; di POOLDATA itu TABEL WARISAN
// SUNGGUHAN berisi data. Memanggilnya di sini berarti menghapusnya.
//
// ⭐ GANTINYA: SATU TRANSAKSI PER TEST, DIAKHIRI `Rollback`. DDL auto-commit di
// Oracle, DML tidak - jadi seluruh `INSERT` di berkas ini hilang saat test
// selesai, lulus maupun gagal. Nol baris menetap di POOLDATA. Itu yang
// menggantikan "buang skemanya sesudahnya", dan ia lebih aman: ia tidak pernah
// menyentuh satu pun tabel di luar yang dimasukinya.
//
// ⚠️ Pengenal 9000xxx dan kode berawalan ZZ tetap dipakai: bila sebuah
// `Rollback` gagal, barisnya harus mudah dikenali sebagai milik uji.

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/sijms/go-ora/v2"

	"nusantarare/inti/backend/config"
)

// siapkan membuka POOLDATA lalu MEMULAI SATU TRANSAKSI; `Rollback`-nya
// didaftarkan sebagai Cleanup, dan itulah seluruh pembersihannya.
func siapkan(t *testing.T) (*sql.Tx, string, context.Context) {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("konfigurasi: %v", err)
	}
	// Produksi tetap ditolak keras - satu-satunya pagar yang tidak dicabut.
	if cfg.IsPegaProd {
		t.Fatal("menolak berjalan saat IS_PEGA_PROD=true (ADR-U-0005)")
	}
	if !cfg.PunyaOracle() {
		t.Skip("lewati: ORACLE_DSN kosong")
	}
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatalf("membuka oracle: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	ctx := context.Background()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	tx, err := sqlDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("memulai transaksi: %v", err)
	}
	// ⛔ Rollback, SELALU - tidak ada cabang yang melakukan Commit.
	t.Cleanup(func() { _ = tx.Rollback() })
	return tx, cfg.OracleSchema, ctx
}

func jalankan(ctx context.Context, db *sql.Tx, skema, q string) error {
	_, err := db.ExecContext(ctx, strings.ReplaceAll(q, "{skema}", skema))
	return err
}

// wajibTolak menuntut pernyataan GAGAL, dan gagal OLEH CONSTRAINT YANG DISEBUT.
//
// ⛔ Memeriksa namanya, bukan sekadar "ada galat": pernyataan yang ditolak
// karena alasan lain - salah ketik kolom, NOT NULL yang terlewat - adalah uji
// yang lulus secara kebetulan, dan itu lebih buruk daripada uji yang merah.
func wajibTolak(t *testing.T, ctx context.Context, db *sql.Tx, skema, constraint, q string) {
	t.Helper()
	err := jalankan(ctx, db, skema, q)
	if err == nil {
		t.Errorf("pernyataan DITERIMA, mau ditolak %s: %s", constraint, q)
		return
	}
	if !strings.Contains(strings.ToUpper(err.Error()), strings.ToUpper(constraint)) {
		t.Errorf("ditolak, tetapi BUKAN oleh %s: %v", constraint, err)
		return
	}
	t.Logf("%s menolak: %v", constraint, err)
}

// wajibTerima menuntut pernyataan BERHASIL. Uji positif ada sebab constraint
// yang menolak TERLALU BANYAK lulus setiap uji negatif yang pernah ditulis.
func wajibTerima(t *testing.T, ctx context.Context, db *sql.Tx, skema, apa, q string) {
	t.Helper()
	if err := jalankan(ctx, db, skema, q); err != nil {
		t.Errorf("%s DITOLAK, mau diterima: %v", apa, err)
	}
}

const (
	insKontrak = `INSERT INTO {skema}.KONTRAK (ID_KONTRAK,ID_CEDANT,ID_ASAL_BISNIS,SIFAT_PROPORSI,TANGGAL_MULAI,TANGGAL_BERAKHIR) VALUES (%d,1,1,'NON_PROPORSIONAL',DATE '2026-01-01',DATE '2026-12-31')`
	insVersi   = `INSERT INTO {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK,ID_KONTRAK,NOMOR_URUT_VERSI,KEADAAN_SIKLUS_HIDUP,NAMA_KONTRAK,KODE_MATA_UANG_KONTRAK,PERSEN_BAGIAN_NURE,BAGIAN_NURE_SERAGAM,MEMAKAI_BORDEREAUX,CARA_PEMBUKUAN,MEMAKAI_PRORATA,RETRO_BERGANDA) VALUES (%d,%d,%d,'DRAFT','Uji',%d,10,'0','0','X','0','0')`
)

func dasar(t *testing.T, ctx context.Context, db *sql.Tx, skema string) {
	t.Helper()
	wajibTerima(t, ctx, db, skema, "kontrak", fmt.Sprintf(insKontrak, 9000001))
}

func cacah(t *testing.T, ctx context.Context, db *sql.Tx, skema, q string) int {
	t.Helper()
	var n int
	if err := db.QueryRowContext(ctx, strings.ReplaceAll(q, "{skema}", skema)).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// Tiket 14, uji POSITIF yang papan sebut menangkap pemisahan identitas yang
// TERLALU KETAT: satu kontrak dengan TIGA versi berturut-turut diterima, dan
// ketiganya berbagi lapisan beku yang sama TANPA menyalinnya.
func TestTigaVersiPadaSatuKontrakDiterima(t *testing.T) {
	db, skema, ctx := siapkan(t)
	dasar(t, ctx, db, skema)
	for i := 1; i <= 3; i++ {
		wajibTerima(t, ctx, db, skema, fmt.Sprintf("versi ke-%d", i),
			fmt.Sprintf(insVersi, 9000000+i, 9000001, i, 9000001))
	}
	if n := cacah(t, ctx, db, skema, "SELECT COUNT(*) FROM {skema}.VERSI_KONTRAK WHERE ID_KONTRAK = 9000001"); n != 3 {
		t.Errorf("mau 3 versi, dapat %d", n)
	}
	// Lapisan beku TIDAK disalin: ia hanya ada di KONTRAK, satu baris.
	if n := cacah(t, ctx, db, skema, "SELECT COUNT(*) FROM {skema}.KONTRAK WHERE ID_KONTRAK = 9000001"); n != 1 {
		t.Errorf("lapisan beku tersalin: %d baris KONTRAK, mau 1", n)
	}
}

// INV-04 - nomor urut versi unik DI DALAM satu kontrak.
func TestNomorUrutVersiGandaDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	dasar(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "versi pertama", fmt.Sprintf(insVersi, 9000001, 9000001, 1, 9000001))
	wajibTolak(t, ctx, db, skema, "UQ_VERSI_KONTRAK", fmt.Sprintf(insVersi, 9000002, 9000001, 1, 9000001))
}

// Kunci asing VERSI_KONTRAK -> KONTRAK: versi yatim ditolak.
func TestVersiYatimDitolak(t *testing.T) {
	db, skema, ctx := siapkan(t)
	dasar(t, ctx, db, skema)
	wajibTolak(t, ctx, db, skema, "FK_VERSI_KONTRAK_1", fmt.Sprintf(insVersi, 9000003, 8888888, 1, 9000001))
}

// INV-68 - KODE unik di KELIMA tabel acuan (MATA_UANG dicabut, migrasi 434), diperiksa satu per satu, bukan
// disimpulkan dari satu.
func TestKodeGandaDitolakDiKelimaTabelAcuan(t *testing.T) {
	db, skema, ctx := siapkan(t)
	acuan := []struct{ tabel, kunci, constraint string }{
		{"JENIS_POTONGAN", "ID_JENIS_POTONGAN", "UQ_JENIS_POTONGAN"},
		{"KELAS_BISNIS", "ID_KELAS_BISNIS", "UQ_KELAS_BISNIS"},
		{"KELOMPOK_TREATY", "ID_KELOMPOK_TREATY", "UQ_KELOMPOK_TREATY"},
		{"BAHAYA", "ID_BAHAYA", "UQ_BAHAYA"},
		{"JENIS_REASURANSI", "ID_JENIS_REASURANSI", "UQ_JENIS_REASURANSI"},
	}
	for i, a := range acuan {
		ins := fmt.Sprintf("INSERT INTO {skema}.%s (%s,KODE,NAMA,AKTIF) VALUES (%%d,'ZZ%d','Uji','1')", a.tabel, a.kunci, i)
		wajibTerima(t, ctx, db, skema, a.tabel+" baris pertama", fmt.Sprintf(ins, 9000001+i*10))
		wajibTolak(t, ctx, db, skema, a.constraint, fmt.Sprintf(ins, 9000002+i*10))
	}
}

// INV-62 uji POSITIF, dan ia POKOK tiket 15: bahaya KESEMBILAN ditambahkan,
// lalu batas tanggungan dicatat untuknya - keduanya TANPA perubahan skema.
//
// Sistem lama menaruh empat bahaya sebagai empat pasang KOLOM di kepala
// kontrak; di sana bahaya kesembilan menuntut kolom kesembilan.
func TestBahayaKesembilanLangsungDapatDipakai(t *testing.T) {
	db, skema, ctx := siapkan(t)
	dasar(t, ctx, db, skema)
	wajibTerima(t, ctx, db, skema, "versi", fmt.Sprintf(insVersi, 9000001, 9000001, 1, 9000001))
	for i := 1; i <= 9; i++ {
		wajibTerima(t, ctx, db, skema, fmt.Sprintf("bahaya ke-%d", i),
			fmt.Sprintf("INSERT INTO {skema}.BAHAYA (ID_BAHAYA,KODE,NAMA,AKTIF) VALUES (%d,'ZB%d','Uji','1')", 9000000+i, i))
	}
	wajibTerima(t, ctx, db, skema, "batas untuk bahaya KESEMBILAN",
		"INSERT INTO {skema}.BATAS_PER_BAHAYA (ID_BATAS_PER_BAHAYA,ID_VERSI_KONTRAK,ID_BAHAYA,NILAI_BATAS) VALUES (9000001,9000001,9000009,1000)")
}
