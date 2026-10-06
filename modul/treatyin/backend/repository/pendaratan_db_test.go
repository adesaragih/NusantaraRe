//go:build db

package repository_test

// Bukti Oracle untuk pemuat tabel pendaratan - SEPULUH KONTRAK NYATA.
//
// ⛔ Satu transaksi per test, diakhiri `Rollback`. Nol `Commit`, nol `DROP`,
// nol teardown, nol `DELETE` di luar transaksi. Sesudah berkas ini selesai,
// kedelapan tabel harus kosong persis seperti sebelumnya - dan
// `TestNolBarisTersisaSesudahSeluruhUji` membuktikannya, bukan
// mengandaikannya.
//
// ⛔ `M_TREATY_IN` dan `TREATY_IN` DIBACA saja di sini.
//
// Kesepuluh kontraknya TIDAK dikarang: ia dipilih dari `TREATY_IN` menurut
// `PROPORTIONTYPE`, lima dan lima, dan dipilih yang dokumennya paling besar -
// dokumen besar memuat lebih banyak larik, dan larik yang kosong tidak
// membuktikan apa pun tentang pemuat.

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/repository"
	"nusantarare/modul/treatyin/backend/services"
)

// gudangUji membuka Gudang sungguhan beserta satu transaksi yang SELALU
// dibatalkan.
//
// ⚠️ Berdiri sendiri, tidak memakai `siapkan(t)`: yang itu mengembalikan
// `*sql.Tx` telanjang, sementara pemuat bekerja lewat `*db.Tx`. Aturannya
// sama persis - satu transaksi per test, `Rollback` lewat `t.Cleanup`.
func gudangUji(t *testing.T) (*repository.Gudang, *db.Tx, context.Context) {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("konfigurasi: %v", err)
	}
	if cfg.IsPegaProd {
		t.Fatal("menolak berjalan saat IS_PEGA_PROD=true (ADR-U-0005)")
	}
	if !cfg.PunyaOracle() {
		t.Skip("lewati: ORACLE_DSN kosong")
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatalf("membuka oracle: %v", err)
	}
	t.Cleanup(func() { _ = d.Close() })
	ctx := context.Background()
	if err := d.Ping(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}
	tx, err := d.Mulai(ctx)
	if err != nil {
		t.Fatalf("memulai transaksi: %v", err)
	}
	// ⛔ Rollback, SELALU.
	t.Cleanup(func() { _ = tx.Rollback() })
	return repository.Baru(d), tx, ctx
}

// kontrakUji memilih lima kontrak dengan `PROPORTIONTYPE` yang diminta,
// yang dokumennya PALING BESAR.
func kontrakUji(t *testing.T, ctx context.Context, sifat string) []string {
	t.Helper()
	cfg, _ := config.Load()
	sqlDB, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = sqlDB.Close() }()

	q := `SELECT m.ID FROM ` + cfg.OracleSchema + `.M_TREATY_IN m
	        JOIN ` + cfg.OracleSchema + `.TREATY_IN t ON t.ID = m.ID
	       WHERE t.PROPORTIONTYPE = :1
	       ORDER BY DBMS_LOB.GETLENGTH(m.JSONDATA) DESC
	       FETCH FIRST 5 ROWS ONLY`
	baris, err := sqlDB.QueryContext(ctx, q, sifat)
	if err != nil {
		t.Fatalf("memilih kontrak %s: %v", sifat, err)
	}
	defer func() { _ = baris.Close() }()
	var out []string
	for baris.Next() {
		var id string
		if err := baris.Scan(&id); err != nil {
			t.Fatal(err)
		}
		out = append(out, id)
	}
	if len(out) != 5 {
		t.Fatalf("%s: dapat %d kontrak, mau 5", sifat, len(out))
	}
	return out
}

// ⭐ UJI UTAMA RONDE INI: sepuluh kontrak nyata, lima prop dan lima
// non-prop, dimuat lalu DICOCOKKAN - cacah baris di tabel harus sama persis
// dengan cacah elemen di dokumennya.
func TestSepuluhKontrakNyataMuatDanCocok(t *testing.T) {
	g, tx, ctx := gudangUji(t)

	var hasil []services.HasilMuat
	for _, sifat := range []string{"Proportional", "NonProportional"} {
		for _, id := range kontrakUji(t, ctx, sifat) {
			h, err := services.MuatSatuKontrak(ctx, g, tx, id)
			if err != nil {
				t.Fatalf("%s %s: %v", sifat, id, err)
			}
			if !h.Cocok() {
				t.Errorf("%s %s TIDAK cocok: %v", sifat, id, h.Selisih)
			}
			hasil = append(hasil, h)
			t.Logf("%-15s %-8s dokumen=%v", sifat, id, ringkasCacah(h.DiDokumen))
		}
	}
	if len(hasil) != 10 {
		t.Fatalf("%d kontrak dimuat, mau 10", len(hasil))
	}

	// ⛔ Pemuatan yang mendaratkan NOL baris pada kesepuluh kontrak akan
	// lulus seluruh pemeriksaan di atas. Dituntut ada yang benar-benar
	// mendarat.
	total := 0
	for _, h := range hasil {
		for _, n := range h.DiTabel {
			total += n
		}
	}
	if total == 0 {
		t.Fatal("nol baris mendarat dari sepuluh kontrak; pemuatnya tidak melakukan apa pun")
	}
	t.Logf("TOTAL %d baris mendarat dari 10 kontrak", total)

	// ⛔ `pxObjClass` SENGAJA tidak mendarat. Ia nama kelas Pega, dan
	// kolomnya dicabut migrasi 436 atas keputusan pemilik proses 5 Oktober
	// 2026: menyimpan nama kelas sistem lama di dalam tabel model baru
	// mengundang orang menurunkan perilaku darinya.
	//
	// ⚠ Ia DISEBUT di sini, bukan disaring di `RingkasTakTerpetakan`.
	// Penjaga yang dilonggarkan di sumbernya berhenti melaporkan kunci
	// BERIKUTNYA yang hilang; daftar sengaja-diabaikan di tempat
	// pemeriksaannya tetap berbunyi untuk yang lain.
	// ⛔ Kunci yang SENGAJA tidak menjadi kolom, dan ketiganya beralasan
	// berbeda:
	//
	//   pxObjClass / pxListSubscript / pyTemplate*  perabot Pega, dicabut
	//       migrasi 436 atas keputusan pemilik proses.
	//   KunciAnak                                   BUKAN kolom melainkan
	//       TABEL ANAKNYA SENDIRI. Melaporkannya sebagai hilang membuat
	//       penjaga ini merah pada struktur yang justru benar.
	sengaja := map[string]bool{"pxObjClass": true, "pxListSubscript": true}
	for _, p := range repository.PetaPendaratan {
		if p.KunciAnak != "" {
			sengaja[p.KunciAnak] = true
		}
	}
	// ⚠️ TERTUNDA, bukan terlewat. Ketujuh belas larik di bawah milik TUJUH
	// tabel yang `Diagram-Skema-Tabel-TreatyIn-dan-EDM-v2.xlsx` sebut tetapi
	// jalur sumbernya TERPOTONG di dalam selnya sendiri — berbunyi `| Tr`,
	// `| TreatyIn.Limi`, putus di tengah kata. Ketujuhnya menggabungkan
	// beberapa larik, dan daftarnya tidak lengkap di mana pun:
	//
	//   T_TREATY_LIMIT_AMOUNT · T_TREATY_LIMIT_MEASURE · T_TREATY_SHARE_AMOUNT
	//   T_TREATY_TOTAL · T_TREATY_LIMIT_SUMMARY · T_TREATY_FAC_SHARE_AMOUNT
	//   T_TREATY_SHARE_SPREAD_AMOUNT
	//
	// ⛔ Didaftar di sini SUPAYA penjaga ini tetap berbunyi untuk larik
	// BERIKUTNYA yang kehilangan rumahnya. Menghapus pemeriksaannya akan
	// menyembunyikan yang kedelapan.
	for _, k := range []string{
		"MDPList", "PremiumEarnedList", "EgnpiTotalList", "Reinstatement_List",
		"LayerList", "GrossPremiumList", "NetPremiumList", "DeductionTotalList",
		"ClassofBusinessList", "RnmGrossPremiDisplay", "RnmLimitList", "RnmLimitListDisplay",
		"RNMSpreadedListXOL", "RNMSpreadedListRIXOL", "RNMSpreadedListNetXOL",
		"RNMSpreadedListNetRIXOL", "RNMSpreadedListGrossXOL", "RNMSpreadedListGrossRIXOL",
		"RNMSpreadedListDeductXOL", "RNMSpreadedListDeductRIXOL",
	} {
		sengaja[k] = true
	}
	for tabel, kunci := range services.RingkasTakTerpetakan(hasil) {
		var nyata []string
		for _, k := range kunci {
			if !sengaja[k] && !strings.HasPrefix(k, "pyTemplate") {
				nyata = append(nyata, k)
			}
		}
		if len(nyata) != 0 {
			t.Errorf("%s: kunci JSON tanpa kolom %v; kunci tanpa kolom mendarat "+
				"sebagai ketiadaan, dan cacah barisnya tetap cocok", tabel, nyata)
		}
	}
}

// ⭐ IDEMPOTEN: memuat kontrak yang SAMA dua kali menghasilkan keadaan yang
// sama persis - bukan dua kali lipat barisnya.
//
// ⛔ Ini sifat yang paling mudah dikira benar tanpa diuji. Pemuat yang hanya
// menyisipkan lulus setiap uji "apakah datanya masuk"; yang menangkapnya
// hanya pemuatan kedua.
func TestMuatDuaKaliTidakMenggandakan(t *testing.T) {
	g, tx, ctx := gudangUji(t)
	id := kontrakUji(t, ctx, "Proportional")[0]

	pertama, err := services.MuatSatuKontrak(ctx, g, tx, id)
	if err != nil {
		t.Fatal(err)
	}
	kedua, err := services.MuatSatuKontrak(ctx, g, tx, id)
	if err != nil {
		t.Fatal(err)
	}
	if !kedua.Cocok() {
		t.Errorf("muatan kedua tidak cocok: %v", kedua.Selisih)
	}
	for tabel, n := range pertama.DiTabel {
		if kedua.DiTabel[tabel] != n {
			t.Errorf("%s: muatan pertama %d baris, kedua %d — pemuatnya menggandakan",
				tabel, n, kedua.DiTabel[tabel])
		}
	}
}

// ⭐ TERBALIKKAN: `KosongkanKontrak` mengembalikan kedelapan tabel ke nol
// untuk kontrak itu, dan tidak menyentuh kontrak lain.
func TestKosongkanMengembalikanKeNolTanpaMenyentuhKontrakLain(t *testing.T) {
	g, tx, ctx := gudangUji(t)
	id := kontrakUji(t, ctx, "Proportional")
	sasaran, tetangga := id[0], id[1]

	if _, err := services.MuatSatuKontrak(ctx, g, tx, sasaran); err != nil {
		t.Fatal(err)
	}
	hTetangga, err := services.MuatSatuKontrak(ctx, g, tx, tetangga)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := g.KosongkanKontrak(ctx, tx, sasaran); err != nil {
		t.Fatal(err)
	}
	sesudah, err := g.CacahBarisKontrak(ctx, tx, sasaran)
	if err != nil {
		t.Fatal(err)
	}
	for tabel, n := range sesudah {
		if n != 0 {
			t.Errorf("%s masih memuat %d baris kontrak %s sesudah dikosongkan", tabel, n, sasaran)
		}
	}
	// ⛔ Dan TETANGGANYA utuh. Pengosongan yang membuang segalanya lulus
	// pemeriksaan di atas dengan sempurna.
	masih, err := g.CacahBarisKontrak(ctx, tx, tetangga)
	if err != nil {
		t.Fatal(err)
	}
	for tabel, n := range hTetangga.DiTabel {
		if masih[tabel] != n {
			t.Errorf("%s: kontrak tetangga %s kehilangan baris, %d -> %d",
				tabel, tetangga, n, masih[tabel])
		}
	}
}

// ⭐ Butir angsuran menunjuk INDUK yang benar - bukan sekadar induk mana pun.
//
// `IDINDUK` yang tertukar antar induk menghasilkan cacah baris yang cocok
// sempurna dan angka angsuran yang menempel pada termin yang salah.
func TestButirAngsuranMenunjukIndukYangMelahirkannya(t *testing.T) {
	g, tx, ctx := gudangUji(t)
	cfg, _ := config.Load()

	var diuji int
	for _, sifat := range []string{"Proportional", "NonProportional"} {
		for _, id := range kontrakUji(t, ctx, sifat) {
			h, err := services.MuatSatuKontrak(ctx, g, tx, id)
			if err != nil {
				t.Fatal(err)
			}
			if h.DiTabel["T_TREATY_INSTALLMENT_ITEM"] == 0 {
				continue
			}
			diuji++
			var yatim int
			q := `SELECT COUNT(*) FROM ` + cfg.OracleSchema + `.T_TREATY_INSTALLMENT_ITEM b
			       WHERE b.MASTERID = :1
			         AND NOT EXISTS (SELECT 1 FROM ` + cfg.OracleSchema + `.T_TREATY_INSTALLMENT i
			                          WHERE i.ID = b.IDINDUK AND i.MASTERID = b.MASTERID)`
			if err := tx.QueryRowContext(ctx, q, id).Scan(&yatim); err != nil {
				t.Fatal(err)
			}
			if yatim != 0 {
				t.Errorf("kontrak %s: %d butir angsuran menunjuk induk milik kontrak lain atau induk yang tidak ada", id, yatim)
			}
		}
	}
	if diuji == 0 {
		t.Skip("lewati: kesepuluh kontrak terbesar tidak punya butir angsuran")
	}
	t.Logf("%d kontrak punya butir angsuran dan seluruhnya menunjuk induk yang benar", diuji)
}

// ⭐ UNIQUE (MASTERID, URUTAN) BENAR-BENAR menolak - jaring kedua di balik
// idempotensi. Tanpa uji ini, constraint yang salah kolom tetap terlihat
// seperti ada.
func TestUrutanGandaDitolakConstraint(t *testing.T) {
	g, tx, ctx := gudangUji(t)
	cfg, _ := config.Load()
	id := kontrakUji(t, ctx, "Proportional")[0]

	if _, err := services.MuatSatuKontrak(ctx, g, tx, id); err != nil {
		t.Fatal(err)
	}
	q := `INSERT INTO ` + cfg.OracleSchema + `.T_VIEW_COMMENT (ID, MASTERID, URUTAN)
	      SELECT ` + cfg.OracleSchema + `.SEQ_TV_COMMENT.NEXTVAL, MASTERID, URUTAN
	        FROM ` + cfg.OracleSchema + `.T_VIEW_COMMENT WHERE MASTERID = :1 AND ROWNUM = 1`
	_, err := tx.ExecContext(ctx, q, id)
	if err == nil {
		t.Fatal("urutan ganda DITERIMA; UQ_MTI_COMMENT tidak menjaga apa pun")
	}
	if !strings.Contains(err.Error(), "UQ_MTI_COMMENT") {
		t.Errorf("ditolak, tetapi bukan oleh UQ_MTI_COMMENT: %v", err)
	}
}

// ⛔ Penutup: kedelapan tabel UTUH - cacahnya persis seperti sebelum
// berkas ini berjalan.
//
// ⚠️ Pernah berbunyi "kedelapan tabel KOSONG", dan itu benar sampai
// 3 Oktober 2026 pukul 16:40 - saat pemuatan sungguhan mengikat 26.536
// baris. Sejak itu "kosong" adalah ukuran yang SALAH: ia akan merah pada
// basis data yang sehat, dan hijau pada basis data yang baru saja
// kehilangan seluruh muatannya. Yang dijaga BUKAN nol, melainkan TIDAK
// BERUBAH - dan itu persis yang `Rollback` janjikan.
//
// Cacah pembandingnya diambil `TestMain` SEBELUM satu pun uji berjalan,
// sehingga ia bekerja di basis data yang sudah dimuat maupun yang belum.
func TestZCacahTidakBerubahSesudahSeluruhUji(t *testing.T) {
	if cacahAwal == nil {
		t.Skip("lewati: cacah awal tidak terambil (oracle tidak terjangkau)")
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("konfigurasi: %v", err)
	}
	d, err := db.Open(cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = d.Close() }()
	ctx := context.Background()
	if err := d.Ping(ctx); err != nil {
		t.Skipf("lewati: oracle tidak terjangkau: %v", err)
	}

	sesudah, err := repository.Baru(d).CacahBarisSeluruhnya(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	for tabel, n := range sesudah {
		if n != cacahAwal[tabel] {
			t.Errorf("%s: %d baris sebelum uji, %d sesudah — sebuah Rollback tidak terjadi, "+
				"atau sebuah Commit masuk ke berkas ini", tabel, cacahAwal[tabel], n)
		}
	}
	t.Logf("kedelapan tabel utuh: %v", sesudah)
}

// cacahAwal - isi kedelapan tabel SEBELUM uji mana pun berjalan. nil bila
// Oracle tidak terjangkau.
var cacahAwal map[string]int

// TestMain mengambil cacah awal, lalu menjalankan ujinya.
//
// ⛔ Ia TIDAK menyiapkan maupun membersihkan apa pun - nol `DROP`, nol
// teardown, nol `DELETE`. Satu-satunya tugasnya membaca.
func TestMain(m *testing.M) {
	cfg, err := config.Load()
	if err == nil && cfg.PunyaOracle() && !cfg.IsPegaProd {
		if d, err := db.Open(cfg); err == nil {
			ctx := context.Background()
			if d.Ping(ctx) == nil {
				if c, err := repository.Baru(d).CacahBarisSeluruhnya(ctx, nil); err == nil {
					cacahAwal = c
				}
			}
			_ = d.Close()
		}
	}
	os.Exit(m.Run())
}

func ringkasCacah(m map[string]int) map[string]int {
	out := map[string]int{}
	for k, v := range m {
		if v != 0 {
			// ⛔ Potong menurut AWALAN yang benar-benar ada, bukan pada
			// posisi tetap. Nama tabel berganti dari `M_TREATYIN_*` menjadi
			// `T_TREATY_*` (migrasi 436), dan potongan berposisi tetap
			// menghasilkan kunci cacat seperti `PORTING_PERIOD` — yang tetap
			// terbaca sebagai nama tabel sampai seseorang membacanya pelan.
			pendek := k
			for _, awalan := range []string{"M_TREATYIN_", "T_TREATY_", "T_VIEW_"} {
				if strings.HasPrefix(pendek, awalan) {
					pendek = strings.TrimPrefix(pendek, awalan)
					break
				}
			}
			out[pendek] = v
		}
	}
	return out
}
