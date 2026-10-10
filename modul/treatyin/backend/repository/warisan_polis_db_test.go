//go:build db

package repository_test

// Bukti Oracle untuk panel `Existing Policy for Master ID`.
//
// ⛔ BACA SAJA terhadap `TREATYINPRODUCTION`. Nol transaksi, nol tulisan,
// nol `Commit`, nol `DROP`.
//
// ⚠️ Yang dijaga di sini bukan "kueri jalan" melainkan bahwa ia memberi
// PERSIS apa yang gambar 01 dan 26 dokumen desain perlihatkan. Uji yang
// hanya memanggil fungsinya akan tetap hijau pada hari kuerinya salah.

import (
	"database/sql"
	"testing"

	"nusantarare/inti/backend/config"
	"nusantarare/modul/treatyin/backend/repository"
)

// ⭐ Kontrak gambar 26 memberi SATU baris, dan isinya cocok sampai ke
// pemotongan 18 aksara `IDPEGA`.
func TestPolisProduksiCocokDenganGambar26(t *testing.T) {
	g, ctx := gudangBaca(t)

	baris, err := g.BacaPolisProduksi(ctx, "1001841")
	if err != nil {
		t.Fatal(err)
	}
	if len(baris) != 1 {
		t.Fatalf("%d baris, gambar 26 memperlihatkan 1: %+v", len(baris), baris)
	}
	b := baris[0]
	if b.NomorPolis != "RNM-QR.T02.05.2025.11987" {
		t.Errorf("Policy No %q, gambar 26 berbunyi RNM-QR.T02.05.2025.11987", b.NomorPolis)
	}
	// ⛔ INI yang membuktikan `@substring(.CARI2,18)` dijalankan: nilai
	// mentahnya `ASM-FW-GISFW-WORK NB-147044`, dan layar lama berbunyi
	// `NB-147044`.
	if b.PegaID != "NB-147044" {
		t.Errorf("Pega ID %q, gambar 26 berbunyi NB-147044 — "+
			"pemotongan 18 aksara tidak berjalan", b.PegaID)
	}
	// Kedua kolom kuartal KOSONG di gambarnya, dan kosong itu sah.
	t.Logf("kuartal %q tahun %q", b.Kuartal, b.TahunKuartal)
}

// ⭐ Kontrak gambar 01 memberi NOL baris — dan gambarnya berbunyi `No items`.
//
// ⛔ Irisan KOSONG, bukan nil: panel yang menerima `null` merender beda.
func TestPolisProduksiKosongCocokDenganGambar01(t *testing.T) {
	g, ctx := gudangBaca(t)

	baris, err := g.BacaPolisProduksi(ctx, "1001846")
	if err != nil {
		t.Fatal(err)
	}
	if baris == nil {
		t.Fatal("nil, mau irisan kosong")
	}
	if len(baris) != 0 {
		t.Errorf("%d baris, gambar 01 berbunyi \"No items\": %+v", len(baris), baris)
	}
}

// ⭐ Pencocokannya atas TUJUH aksara pertama, bukan kesamaan penuh.
//
// ⚠️ Itu yang SQL rule-nya lakukan, dan ia berbeda hasilnya: pengenal yang
// tujuh aksara pertamanya sama tetapi ekornya berbeda tetap ikut terbaca.
// Uji ini merah pada hari seseorang "merapikan" `SUBSTR` menjadi `=`.
func TestPolisProduksiMencocokkanTujuhAksaraPertama(t *testing.T) {
	g, ctx := gudangBaca(t)
	cfg, _ := config.Load()
	h, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = h.Close() }()

	var lewatRepo, lewatSQL int
	baris, err := g.BacaPolisProduksi(ctx, "1001841")
	if err != nil {
		t.Fatal(err)
	}
	lewatRepo = len(baris)

	if err := h.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM (SELECT DISTINCT NOPOLIS, IDPEGA, QUARTER, QUARTER_YEAR
		   FROM `+cfg.OracleSchema+`.`+repository.TabelWarisanProduksi+`
		  WHERE SUBSTR(NOOFFER,1,7) = SUBSTR('1001841',1,7))`).Scan(&lewatSQL); err != nil {
		t.Fatal(err)
	}
	if lewatRepo != lewatSQL {
		t.Errorf("repository %d baris, SQL rule-nya %d", lewatRepo, lewatSQL)
	}

	// ⚠️ Dan kesamaan PENUH diukur di sampingnya — TANPA dijadikan syarat.
	//
	// Pada kontrak ini keduanya kebetulan memberi 1 baris, jadi angka itu
	// tidak membedakan apa pun hari ini. Yang menjaga `SUBSTR` tetap ada
	// adalah SQL rule-nya, bukan uji ini; angkanya dicatat supaya hari
	// keduanya berbeda terlihat di log, bukan supaya uji ini merah.
	var lewatSamaPersis int
	if err := h.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM (SELECT DISTINCT NOPOLIS, IDPEGA, QUARTER, QUARTER_YEAR
		   FROM `+cfg.OracleSchema+`.`+repository.TabelWarisanProduksi+`
		  WHERE NOOFFER = '1001841')`).Scan(&lewatSamaPersis); err != nil {
		t.Fatal(err)
	}
	t.Logf("tujuh aksara: %d baris; kesamaan penuh: %d baris", lewatSQL, lewatSamaPersis)
}

// ⭐ Tabelnya UTUH — modul ini tidak pernah menulisinya.
func TestTreatyInProductionTidakBerubah(t *testing.T) {
	cfg, _ := config.Load()
	h, err := sql.Open("oracle", cfg.OracleDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = h.Close() }()
	_, ctx := gudangBaca(t)

	var n int
	if err := h.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM `+cfg.OracleSchema+`.`+repository.TabelWarisanProduksi).Scan(&n); err != nil {
		t.Fatal(err)
	}
	// ⛔ LANTAI, bukan angka PERSIS — ralat 6 Oktober 2026.
	//
	// Angka persis membuat gerbang ini merah setiap kali bisnis berjalan
	// normal: Pega menambah baris ke tabel warisan ini, dan penambahan itu
	// bukan regresi. Terukur hari itu — `TREATY_IN` dan `M_TREATY_IN` sudah
	// 1.855 dari 1.854 (kontrak `1002059`), `M_TREATY_IN_DETAIL` 27.618 dari
	// 27.617.
	//
	// ⚠️ YANG HILANG karenanya, dan itu dinyatakan: penambahan baris oleh
	// MODUL INI tidak lagi tertangkap cacahnya. Yang menangkapnya
	// `TestWarisanHanyaDibaca`, yang menyapu naskah SQL dan menolak
	// `INSERT`/`UPDATE`/`DELETE` terhadap tabel ini — penjaga yang
	// sesungguhnya, dan yang 6 Oktober 2026 diperluas supaya mencakupnya.
	if n < 41936 {
		t.Errorf("%s %d baris, terukur 41.936 pada 5 Oktober 2026 — "+
			"tabel WARISAN ini tidak boleh ditulis oleh modul mana pun",
			repository.TabelWarisanProduksi, n)
	}
}
