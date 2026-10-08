//go:build db

package repository_test

// Bukti Oracle untuk tiket 76 (`NILAI_SELISIH`) dan 77
// (`NILAI_SEBELUM_PRO_RATE`).
//
// ⛔ Satu transaksi per test, diakhiri `Rollback` lewat `t.Cleanup` di
// `siapkan(t)`. Nol `Commit`, nol `DROP`, nol teardown — kedua tabel tetap
// NOL BARIS sesudah berkas ini berjalan, dan
// `TestZNilaiTetapNolBarisSesudahUji` membuktikannya alih-alih
// mengandaikannya.
//
// =====================================================================
// APA YANG DIBANGUN TIKET INI, DAN APA YANG TIDAK
// =====================================================================
//   Tiket 76 berbunyi *"PK membangun TABELNYA — entitas beserta kolom,
//   kunci alami, kunci asing, dan perilaku hapusnya"*, dan menyatakan
//   terang-terangan: *"Tidak termasuk: padanan kunci bisnis saat simpan —
//   tiket 06; TABEL INI TIDAK MENGISI DIRINYA SENDIRI."*
//
//   Jadi yang dibuktikan di sini BENTUKNYA: kunci alami yang menolak
//   kembar, mata uang DI DALAM kunci itu, dan kaskade yang benar-benar
//   menghapus. Jalur tulis aplikasi milik tiket 06.
//
// ⛔ TIKET 77 TIDAK PUNYA RUMUS, dan itu temuan yang diukur. Activity
// keempat belas `TreatyEDMProRateCalculation.xml` — satu-satunya rule yang
// menyentuh `ValueBeforeProrate` — langkahnya berbunyi persis
// *"Copy value from ValueDifference to ValueBeforeProrate"*, dengan
// `PropertiesName=TreatyIn.ValueBeforeProrate` dan
// `PropertiesValue=TreatyIn.ValueDifference`. Ia SALINAN, bukan hitungan.
// `GRL-15` menyatakan mesin pro rata sengaja tidak dibangun kembali, jadi
// angkanya masuk sebagai FAKTA — dan tiket 77 menulis jalur gagalnya
// sendiri: *"nilai di tabel ini DIHITUNG ULANG dari nilai sesudah pro rata
// -> itu melanggar maksudnya."*

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

// Pernyataan penyisip. Pengenal 9100xxx: bila sebuah `Rollback` gagal,
// barisnya harus mudah dikenali sebagai milik uji.
const (
	insKontrak = `INSERT INTO {skema}.KONTRAK (ID_KONTRAK,ID_CEDANT,ID_ASAL_BISNIS,SIFAT_PROPORSI,TANGGAL_MULAI,TANGGAL_BERAKHIR)
	              VALUES (%d,1,1,'NON_PROPORSIONAL',DATE '2026-01-01',DATE '2026-12-31')`
	// ⚠️ `MEMAKAI_PRORATA` adalah ruas kedua dari belakang — uji positif
	// kedua tiket 77 menyetelnya '0'.
	insVersi = `INSERT INTO {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK,ID_KONTRAK,NOMOR_URUT_VERSI,KEADAAN_SIKLUS_HIDUP,NAMA_KONTRAK,KODE_MATA_UANG_KONTRAK,PERSEN_BAGIAN_NURE,BAGIAN_NURE_SERAGAM,MEMAKAI_BORDEREAUX,CARA_PEMBUKUAN,MEMAKAI_PRORATA,RETRO_BERGANDA)
	            VALUES (%d,%d,%d,'DRAFT','Uji',1,10,'0','0','X','%s','0')`
	insSelisih = `INSERT INTO {skema}.NILAI_SELISIH (ID_NILAI_SELISIH,ID_VERSI_KONTRAK,KODE_BESARAN,KODE_MATA_UANG,NILAI_LAMA,NILAI_BARU)
	              VALUES (%d,%d,'%s','%s',100,200)`
	insSebelum = `INSERT INTO {skema}.NILAI_SEBELUM_PRO_RATE (ID_NILAI_SEBELUM_PRO_RATE,ID_VERSI_KONTRAK,KODE_BESARAN,KODE_MATA_UANG,NILAI)
	              VALUES (%d,%d,'%s','%s',150)`
)

func jalankan(ctx context.Context, tx *db.Tx, skema, q string) error {
	_, err := tx.ExecContext(ctx, strings.ReplaceAll(q, "{skema}", skema))
	return err
}

// wajibTolak menuntut pernyataan GAGAL, dan gagal OLEH CONSTRAINT YANG
// DISEBUT.
//
// ⛔ Memeriksa NAMANYA, bukan sekadar "ada galat": pernyataan yang ditolak
// karena alasan lain — salah ketik kolom, NOT NULL yang terlewat — adalah
// uji yang lulus secara kebetulan, dan itu lebih buruk daripada uji merah.
func wajibTolak(t *testing.T, ctx context.Context, tx *db.Tx, skema, constraint, q string) {
	t.Helper()
	err := jalankan(ctx, tx, skema, q)
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

// wajibTerima menuntut pernyataan BERHASIL.
//
// ⛔ Uji positif ada sebab constraint yang menolak TERLALU BANYAK lulus
// setiap uji negatif yang pernah ditulis untuknya.
func wajibTerima(t *testing.T, ctx context.Context, tx *db.Tx, skema, apa, q string) {
	t.Helper()
	if err := jalankan(ctx, tx, skema, q); err != nil {
		t.Errorf("%s DITOLAK, mau diterima: %v", apa, err)
	}
}

func cacah(t *testing.T, ctx context.Context, tx *db.Tx, skema, q string) int {
	t.Helper()
	var n int
	if err := tx.QueryRowContext(ctx, strings.ReplaceAll(q, "{skema}", skema)).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

// pondasiVersi menyiapkan satu kontrak dan satu versi.
func pondasiVersi(t *testing.T, ctx context.Context, tx *db.Tx, skema, prorata string) {
	t.Helper()
	wajibTerima(t, ctx, tx, skema, "kontrak", fmt.Sprintf(insKontrak, 9100001))
	wajibTerima(t, ctx, tx, skema, "versi", fmt.Sprintf(insVersi, 9100001, 9100001, 1, prorata))
}

// ============================ TIKET 76 ============================

// ⛔ NEGATIF tiket 76 — baris kembar pada sumbu PENUH ditolak kunci
// alaminya. Jalur gagal pertama tiket: *"dua baris selisih berbesaran dan
// bermata-uang sama pada satu versi -> ditolak kunci alaminya."*
func TestTiket76SelisihKembarDitolak(t *testing.T) {
	_, tx, ctx := siapkan(t)
	skema := skemaUji(t)
	pondasiVersi(t, ctx, tx, skema, "0")

	wajibTerima(t, ctx, tx, skema, "selisih pertama",
		fmt.Sprintf(insSelisih, 9100001, 9100001, "BrokeragePercent", "IDR"))
	wajibTolak(t, ctx, tx, skema, "UQ_NILAI_SELISIH",
		fmt.Sprintf(insSelisih, 9100002, 9100001, "BrokeragePercent", "IDR"))
}

// ⛔ NEGATIF tiket 76 — nilai tanpa mata uang ditolak (`INV-36`).
func TestTiket76SelisihTanpaMataUangDitolak(t *testing.T) {
	_, tx, ctx := siapkan(t)
	skema := skemaUji(t)
	pondasiVersi(t, ctx, tx, skema, "0")

	q := `INSERT INTO {skema}.NILAI_SELISIH (ID_NILAI_SELISIH,ID_VERSI_KONTRAK,KODE_BESARAN,NILAI_LAMA,NILAI_BARU)
	      VALUES (9100003,9100001,'BrokeragePercent',100,200)`
	// Ditegakkan NOT NULL pada kolomnya; namanya dibangkitkan Oracle, jadi
	// yang diadu kodenya.
	err := jalankan(ctx, tx, skema, q)
	if err == nil {
		t.Fatal("nilai selisih tanpa mata uang DITERIMA; INV-36 menuntut mata uang")
	}
	if !strings.Contains(err.Error(), "ORA-01400") {
		t.Errorf("ditolak, tetapi bukan oleh NOT NULL: %v", err)
	}
	t.Logf("tanpa mata uang ditolak: %v", err)
}

// ⛔ NEGATIF tiket 76 — kolom rujukan versi lama TIDAK ADA.
//
// Jalur gagal keempat tiket: *"kolom `ID_VERSI_KONTRAK_LAMA` ditambahkan ->
// ditolak tinjauan skema, `ERD.md` §2.6."* Sisi lama dibaca lewat
// `VERSI_KONTRAK.ID_VERSI_KONTRAK_DASAR` pada induknya.
func TestTiket76NolKolomRujukanVersiLama(t *testing.T) {
	_, tx, ctx := siapkan(t)
	skema := skemaUji(t)

	n := cacah(t, ctx, tx, skema,
		`SELECT COUNT(*) FROM all_tab_columns WHERE owner='`+skema+`'
		   AND table_name='NILAI_SELISIH' AND column_name LIKE '%LAMA%'
		   AND column_name <> 'NILAI_LAMA'`)
	if n != 0 {
		t.Errorf("%d kolom rujukan versi lama ada; ERD.md §2.6 menolaknya — "+
			"sisi lama dibaca lewat VERSI_KONTRAK.ID_VERSI_KONTRAK_DASAR", n)
	}
}

// ⭐ POSITIF tiket 76 — dan ia yang membuktikan `ADR-0048` butir 3 sungguh
// diterapkan: besaran yang SAMA dalam IDR dan USD adalah DUA fakta yang
// sah, sebab mata uang ada di dalam kunci.
//
// ⛔ `UQ` tanpa mata uang akan menolak yang kedua, dan menolaknya diam-diam
// benar menurut setiap uji negatif di atas.
func TestTiket76BesaranSamaDuaMataUangDiterima(t *testing.T) {
	_, tx, ctx := siapkan(t)
	skema := skemaUji(t)
	pondasiVersi(t, ctx, tx, skema, "0")

	wajibTerima(t, ctx, tx, skema, "selisih IDR",
		fmt.Sprintf(insSelisih, 9100010, 9100001, "BrokeragePercent", "IDR"))
	wajibTerima(t, ctx, tx, skema, "selisih USD",
		fmt.Sprintf(insSelisih, 9100011, 9100001, "BrokeragePercent", "USD"))

	n := cacah(t, ctx, tx, skema,
		`SELECT COUNT(*) FROM {skema}.NILAI_SELISIH WHERE ID_VERSI_KONTRAK = 9100001`)
	if n != 2 {
		t.Errorf("%d baris selisih, mau 2 — mata uang tidak ada di dalam kunci", n)
	}
}

// ⭐ POSITIF KEDUA tiket 76 — menghapus versi induk menghapus baris
// selisihnya, DIJALANKAN terhadap Oracle dan bukan diargumentasikan.
func TestTiket76HapusVersiIkutMenghapusSelisih(t *testing.T) {
	_, tx, ctx := siapkan(t)
	skema := skemaUji(t)
	pondasiVersi(t, ctx, tx, skema, "0")
	wajibTerima(t, ctx, tx, skema, "selisih",
		fmt.Sprintf(insSelisih, 9100020, 9100001, "BrokeragePercent", "IDR"))

	if n := cacah(t, ctx, tx, skema,
		`SELECT COUNT(*) FROM {skema}.NILAI_SELISIH WHERE ID_VERSI_KONTRAK = 9100001`); n != 1 {
		t.Fatalf("%d baris sebelum hapus, mau 1", n)
	}
	wajibTerima(t, ctx, tx, skema, "hapus versi",
		`DELETE FROM {skema}.VERSI_KONTRAK WHERE ID_VERSI_KONTRAK = 9100001`)
	if n := cacah(t, ctx, tx, skema,
		`SELECT COUNT(*) FROM {skema}.NILAI_SELISIH WHERE ID_VERSI_KONTRAK = 9100001`); n != 0 {
		t.Errorf("%d baris selisih tersisa sesudah versinya dihapus; ERD.md §2.6 "+
			"menuntut ikut hapus", n)
	}
}

// ============================ TIKET 77 ============================

// ⛔ NEGATIF tiket 77 — baris kembar ditolak kunci alaminya.
func TestTiket77SebelumProRataKembarDitolak(t *testing.T) {
	_, tx, ctx := siapkan(t)
	skema := skemaUji(t)
	pondasiVersi(t, ctx, tx, skema, "1")

	wajibTerima(t, ctx, tx, skema, "sebelum pro rata pertama",
		fmt.Sprintf(insSebelum, 9100030, 9100001, "RNMShare", "IDR"))
	wajibTolak(t, ctx, tx, skema, "UQ_NILAI_SEBELUM_PRO_RATE",
		fmt.Sprintf(insSebelum, 9100031, 9100001, "RNMShare", "IDR"))
}

// ⛔ NEGATIF tiket 77 — nilai tanpa mata uang ditolak (`INV-36`).
func TestTiket77SebelumProRataTanpaMataUangDitolak(t *testing.T) {
	_, tx, ctx := siapkan(t)
	skema := skemaUji(t)
	pondasiVersi(t, ctx, tx, skema, "1")

	q := `INSERT INTO {skema}.NILAI_SEBELUM_PRO_RATE (ID_NILAI_SEBELUM_PRO_RATE,ID_VERSI_KONTRAK,KODE_BESARAN,NILAI)
	      VALUES (9100032,9100001,'RNMShare',150)`
	err := jalankan(ctx, tx, skema, q)
	if err == nil {
		t.Fatal("nilai sebelum pro rata tanpa mata uang DITERIMA; INV-36 menuntut mata uang")
	}
	if !strings.Contains(err.Error(), "ORA-01400") {
		t.Errorf("ditolak, tetapi bukan oleh NOT NULL: %v", err)
	}
}

// ⭐ POSITIF tiket 77 — besaran yang sama dalam dua mata uang diterima.
func TestTiket77BesaranSamaDuaMataUangDiterima(t *testing.T) {
	_, tx, ctx := siapkan(t)
	skema := skemaUji(t)
	pondasiVersi(t, ctx, tx, skema, "1")

	wajibTerima(t, ctx, tx, skema, "sebelum pro rata IDR",
		fmt.Sprintf(insSebelum, 9100040, 9100001, "RNMShare", "IDR"))
	wajibTerima(t, ctx, tx, skema, "sebelum pro rata USD",
		fmt.Sprintf(insSebelum, 9100041, 9100001, "RNMShare", "USD"))

	if n := cacah(t, ctx, tx, skema,
		`SELECT COUNT(*) FROM {skema}.NILAI_SEBELUM_PRO_RATE WHERE ID_VERSI_KONTRAK = 9100001`); n != 2 {
		t.Errorf("%d baris, mau 2", n)
	}
}

// ⭐⭐ POSITIF KEDUA tiket 77 — DAN IA YANG MEMBUKTIKAN TABEL INI BUKAN
// TURUNAN.
//
// Tiket menulisnya sendiri: *"versi yang `MEMAKAI_PRORATA`-nya TIDAK
// disetel -> baris sebelum-pro-rata TETAP BOLEH ADA, dan nilainya sama
// dengan nilai berjalan. Tabel yang menolak baris itu mengaku dirinya hasil
// hitungan."*
//
// ⛔ Inilah uji yang gagal bila seseorang kelak memasang CHECK
// `MEMAKAI_PRORATA = '1'` pada tabel ini, atau menghitung isinya alih-alih
// menerimanya sebagai fakta. `GRL-15`: mesin pro rata sengaja tidak
// dibangun kembali, jadi tidak ada yang dapat menghitungnya ulang.
func TestTiket77BarisTetapSahWalauProRataTidakDipakai(t *testing.T) {
	_, tx, ctx := siapkan(t)
	skema := skemaUji(t)
	// ⛔ MEMAKAI_PRORATA = '0' — pro rata TIDAK dipakai pada versi ini.
	pondasiVersi(t, ctx, tx, skema, "0")

	wajibTerima(t, ctx, tx, skema, "sebelum pro rata pada versi tanpa pro rata",
		fmt.Sprintf(insSebelum, 9100050, 9100001, "RNMShare", "IDR"))

	if n := cacah(t, ctx, tx, skema,
		`SELECT COUNT(*) FROM {skema}.NILAI_SEBELUM_PRO_RATE WHERE ID_VERSI_KONTRAK = 9100001`); n != 1 {
		t.Errorf("%d baris; tabel ini MENOLAK baris pada versi tanpa pro rata, "+
			"yang berarti ia mengaku dirinya hasil hitungan — tiket 77 menolak itu", n)
	}
}

// ⭐ POSITIF tiket 77 — kaskade dari versinya, dijalankan terhadap Oracle.
func TestTiket77HapusVersiIkutMenghapusSebelumProRata(t *testing.T) {
	_, tx, ctx := siapkan(t)
	skema := skemaUji(t)
	pondasiVersi(t, ctx, tx, skema, "1")
	wajibTerima(t, ctx, tx, skema, "sebelum pro rata",
		fmt.Sprintf(insSebelum, 9100060, 9100001, "RNMShare", "IDR"))

	wajibTerima(t, ctx, tx, skema, "hapus versi",
		`DELETE FROM {skema}.VERSI_KONTRAK WHERE ID_VERSI_KONTRAK = 9100001`)
	if n := cacah(t, ctx, tx, skema,
		`SELECT COUNT(*) FROM {skema}.NILAI_SEBELUM_PRO_RATE WHERE ID_VERSI_KONTRAK = 9100001`); n != 0 {
		t.Errorf("%d baris tersisa sesudah versinya dihapus", n)
	}
}

// ⛔ Penutup: kedua tabel KOSONG. Dijalankan di luar transaksi mana pun,
// jadi ia melihat apa yang orang lain lihat.
//
// Namanya berawalan `TestZ` supaya `go test` menjalankannya terakhir —
// urutan dalam satu berkas adalah urutan deklarasi.
func TestZNilaiTetapNolBarisSesudahUji(t *testing.T) {
	_, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)

	for _, tabel := range []string{"NILAI_SELISIH", "NILAI_SEBELUM_PRO_RATE"} {
		var n int
		if err := h.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+skema+`.`+tabel).Scan(&n); err != nil {
			t.Fatal(err)
		}
		if n != 0 {
			t.Errorf("%s memuat %d baris sesudah uji; sebuah Rollback tidak terjadi, "+
				"atau sebuah Commit masuk ke berkas ini", tabel, n)
		}
	}
}
