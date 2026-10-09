//go:build db

package repository_test

// Bukti Oracle untuk panel Attachment modul Adjustment.
//
// ⛔ BACA SAJA terhadap `M_ATTACHMENTTREATY_2`. Nol tulisan, nol `Commit`.

import (
	"testing"
)

// ⭐ 8 Oktober 2026 — penamaan kategori: kesebelas pasangan dari katalog RD
// Pega `M_KATEGORIMASTERTREATY`, termasuk keempat kode yang dulu tanpa nama.
// BACA SAJA terhadap `M_KATEGORIMASTERTREATY` dan `M_ATTACHMENTTREATY_2`.
func TestKatalogKategoriAdjustmentDariMaster(t *testing.T) {
	g, ctx := bacaSaja(t)
	k, err := g.BacaKatalogKategoriLampiran(ctx)
	if err != nil {
		t.Fatal(err)
	}
	mau := map[string]string{
		"00003": "Binding, signed share Email", "00004": "Info Pack",
		"00008": "Letter of Acknowledgment / LOA", "00009": "Claim Data",
		"00007": "Pega Proportional Calculation /Perhitungan Pega Proportional",
	}
	for kode, nama := range mau {
		if k[kode] != nama {
			t.Errorf("%s = %q, mau %q", kode, k[kode], nama)
		}
	}
	if len(k) < 11 {
		t.Errorf("katalog %d pasangan, mau >= 11: %v", len(k), k)
	}
}

// ⭐ Ke-43 baris UTUH sesudah modul ini membacanya.
func TestLampiranWarisanTidakBerubah(t *testing.T) {
	_, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)

	var n int
	if err := h.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+skema+`.M_ATTACHMENTTREATY_2`).Scan(&n); err != nil {
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
	if n < 43 {
		t.Errorf("M_ATTACHMENTTREATY_2 %d baris, terukur 43 pada 4 Oktober 2026 — "+
			"tabel WARISAN ini tidak boleh ditulis oleh modul mana pun", n)
	}
}

// ⭐ Pembacaannya menyaring per kontrak dan cacahnya cocok dengan Oracle.
func TestBacaLampiranKontrakAdjustment(t *testing.T) {
	g, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)

	var id string
	if err := h.QueryRowContext(ctx,
		`SELECT TREATYID FROM `+skema+`.M_ATTACHMENTTREATY_2
		  WHERE TREATYID IS NOT NULL GROUP BY TREATYID
		  ORDER BY COUNT(*) DESC, TREATYID FETCH FIRST 1 ROWS ONLY`).Scan(&id); err != nil {
		t.Skipf("lewati: nol lampiran berpengenal kontrak: %v", err)
	}

	baris, err := g.BacaLampiranKontrak(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	var n int
	if err := h.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM `+skema+`.M_ATTACHMENTTREATY_2 WHERE TREATYID = :1`, id).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if len(baris) != n {
		t.Errorf("%s: dibaca %d, Oracle punya %d", id, len(baris), n)
	}
	if len(baris) == 0 {
		t.Fatalf("%s: nol lampiran padahal ia yang terbanyak", id)
	}
	for i, b := range baris {
		if b.ID == "" || b.NamaBerkas == "" {
			t.Errorf("%s baris %d: medan pokok kosong (%+v)", id, i, b)
		}
	}
	t.Logf("kontrak %s: %d lampiran", id, len(baris))

	// ⛔ Kuerinya BENAR-BENAR menyaring.
	kosong, err := g.BacaLampiranKontrak(ctx, "ZZ-TIDAK-ADA")
	if err != nil {
		t.Fatal(err)
	}
	if len(kosong) != 0 {
		t.Errorf("pengenal tak dikenal mengembalikan %d lampiran", len(kosong))
	}
	if kosong == nil {
		t.Error("nil, mau irisan kosong")
	}
}

// ⭐ Harness `siapkan(t)` benar-benar MEMBUKA transaksi dan benar-benar
// MEMBATALKANNYA.
//
// ⛔ Uji pertama yang memakai harness baru harus membuktikan harnessnya,
// bukan hanya memakainya. Harness yang transaksinya tidak pernah terbuka
// membuat setiap uji tulis sesudahnya menulis LANGSUNG ke POOLDATA.
func TestHarnessMembukaTransaksiDanMembatalkannya(t *testing.T) {
	_, tx, ctx := siapkan(t)

	if !tx.Terisi() {
		t.Fatal("transaksi kosong; uji tulis mana pun sesudah ini akan menulis langsung ke skema")
	}
	// Satu pernyataan yang hanya hidup di dalam transaksi.
	var satu int
	if err := tx.QueryRowContext(ctx, "SELECT 1 FROM DUAL").Scan(&satu); err != nil {
		t.Fatalf("transaksi tidak dapat dipakai: %v", err)
	}
	if satu != 1 {
		t.Fatalf("SELECT 1 mengembalikan %d", satu)
	}
	// ⚠️ Pembatalannya dikerjakan `t.Cleanup`, dan itu sengaja TIDAK
	// dipanggil di sini: yang hendak dibuktikan adalah jalur yang dipakai
	// uji lain, bukan jalur yang ditulis khusus untuk uji ini.
}

// ---------------------------------------------------------------------
// Panel History — `T_VIEW_COMMENT`, BACA SAJA.
// ---------------------------------------------------------------------

// ⭐ Modul ini membaca tabel itu SENDIRI, dan isinya cocok dengan Oracle.
//
// ⚠️ Tabelnya dibuat dan dimuat modul `treatyin` (migrasi 430, 11.365
// baris). Modul ini bukan pemiliknya — yang dibuktikan di sini hanya bahwa
// pembacaannya benar dan tidak menulis apa pun.
func TestBacaRiwayatKontrakAdjustment(t *testing.T) {
	g, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)

	var id string
	if err := h.QueryRowContext(ctx,
		`SELECT MASTERID FROM `+skema+`.T_VIEW_COMMENT
		  GROUP BY MASTERID ORDER BY COUNT(*) DESC, MASTERID FETCH FIRST 1 ROWS ONLY`).Scan(&id); err != nil {
		t.Skipf("lewati: nol riwayat: %v", err)
	}

	baris, err := g.BacaRiwayatKontrak(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	var n int
	if err := h.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM `+skema+`.T_VIEW_COMMENT WHERE MASTERID = :1`, id).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if len(baris) != n {
		t.Errorf("%s: dibaca %d baris, Oracle punya %d", id, len(baris), n)
	}
	if len(baris) == 0 {
		t.Fatalf("%s: nol riwayat padahal ia yang terbanyak", id)
	}
	// Keempat kolomnya berpadanan satu-ke-satu; yang pokok tidak boleh kosong.
	if baris[0].Tanggal == "" || baris[0].Operator == "" {
		t.Errorf("%s baris pertama: medan pokok kosong (%+v)", id, baris[0])
	}
	t.Logf("kontrak %s: %d baris riwayat, pertama %q oleh %q",
		id, len(baris), baris[0].Tanggal, baris[0].Operator)

	// ⛔ Kuerinya BENAR-BENAR menyaring.
	kosong, err := g.BacaRiwayatKontrak(ctx, "ZZ-TIDAK-ADA")
	if err != nil {
		t.Fatal(err)
	}
	if len(kosong) != 0 {
		t.Errorf("pengenal tak dikenal mengembalikan %d baris", len(kosong))
	}
	if kosong == nil {
		t.Error("nil, mau irisan kosong")
	}
}

// ⭐ URUTANNYA sesuai `URUTAN`, bukan urutan acak Oracle.
func TestRiwayatBerurutMenurutUrutan(t *testing.T) {
	g, ctx := bacaSaja(t)
	h, skema := sqlMentah(t)

	var id string
	if err := h.QueryRowContext(ctx,
		`SELECT MASTERID FROM `+skema+`.T_VIEW_COMMENT
		  GROUP BY MASTERID HAVING COUNT(*) > 3 ORDER BY COUNT(*) DESC, MASTERID
		  FETCH FIRST 1 ROWS ONLY`).Scan(&id); err != nil {
		t.Skipf("lewati: nol kontrak berriwayat banyak: %v", err)
	}

	baris, err := g.BacaRiwayatKontrak(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	// Dibandingkan dengan urutan yang Oracle berikan untuk ORDER BY URUTAN.
	rows, err := h.QueryContext(ctx,
		`SELECT OPERATORNAME FROM `+skema+`.T_VIEW_COMMENT WHERE MASTERID = :1 ORDER BY URUTAN`, id)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rows.Close() }()
	var mau []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			t.Fatal(err)
		}
		mau = append(mau, s)
	}
	if len(baris) != len(mau) {
		t.Fatalf("%d lawan %d baris", len(baris), len(mau))
	}
	for i := range mau {
		if baris[i].Operator != mau[i] {
			t.Errorf("urutan ke-%d: %q, mau %q", i, baris[i].Operator, mau[i])
		}
	}
	t.Logf("urutan riwayat terbukti pada kontrak %s, %d baris", id, len(baris))
}
