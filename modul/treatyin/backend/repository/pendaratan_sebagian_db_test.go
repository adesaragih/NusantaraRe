//go:build db

package repository_test

// Bukti Oracle untuk perbaikan cacat KEHILANGAN DATA tombol Save
// (7 Oktober 2026).
//
// ⛔ Satu transaksi, diakhiri `Rollback`. Nol `Commit`, nol teardown.
//
// ⛔ NOL SENTUHAN `M_TREATY_IN` DAN NOL BACA JSON. Berkas ini TIDAK memakai
// `kontrakUji`, yang memilih kontrak dengan membaca `M_TREATY_IN.JSONDATA`
// — pemilik proses menutup akses itu TOTAL, termasuk untuk uji. Dokumennya
// disusun di sini sebagai data Go, dan pengenal kontraknya karangan: yang
// hendak dibuktikan perilaku PEMUAT, bukan isi satu kontrak nyata.

import (
	"testing"

	"nusantarare/inti/backend/config"
)

// skemaUji - skema Oracle yang sedang dipakai.
func skemaUji(t *testing.T) string {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("konfigurasi: %v", err)
	}
	return cfg.OracleSchema
}

// Pengenal karangan — tidak pernah mengikat, dan tidak boleh bertabrakan
// dengan kontrak nyata mana pun.
const idUjiSebagian = "UJI-SAVE-SEBAGIAN"

// Dokumen UTUH: satu baris di beberapa tabel, berikut skalar akar.
func dokumenUtuh() map[string]any {
	return map[string]any{
		"ID":                 idUjiSebagian,
		"TreatyContractName": "KONTRAK UJI",
		// ⭐ Skalar akar milik tab yang TIDAK akan dibuka pada Save kedua.
		"Exclusions":  "SELURUH PENGECUALIAN",
		"Information": "KETERANGAN",
		"Limits": []any{
			map[string]any{"TreatyType": "QUOTA SHARE", "Detail": []any{
				map[string]any{"TreatyGroup": "PROPERTY"},
			}},
		},
		"Portfolio": []any{map[string]any{"Type": "A", "Description": "B"}},
	}
}

// ⛔ UJI UTAMA: Save yang TIDAK membawa `Limits` tidak boleh menghapusnya.
//
// Inilah cacat yang pemakai laporkan: *"saat melakukan inputan tiba tiba
// datanya hilang dan semua di reset malah tidak tersimpan"*.
func TestSaveTanpaLimitsTidakMenghapusLimits(t *testing.T) {
	g, tx, ctx := gudangUji(t)

	if _, err := g.MuatKontrak(ctx, tx, idUjiSebagian, dokumenUtuh()); err != nil {
		t.Fatalf("memuat dokumen utuh: %v", err)
	}
	sebelum, err := g.CacahBarisKontrak(ctx, tx, idUjiSebagian)
	if err != nil {
		t.Fatal(err)
	}
	if sebelum["T_TREATY_LIMITS"] == 0 {
		t.Fatalf("persiapan gagal: nol baris Limits — uji ini kehilangan sasarannya (%v)", sebelum)
	}

	// Save berikutnya: pemakai hanya menyunting KEPALA. Nol larik terkirim,
	// persis bentuk `susunDokumen` ketika nol tab dibuka.
	doc := map[string]any{"ID": idUjiSebagian, "TreatyContractName": "DIUBAH"}
	if _, err := g.MuatKontrakSebagian(ctx, tx, idUjiSebagian, doc); err != nil {
		t.Fatalf("Save sebagian: %v", err)
	}
	sesudah, err := g.CacahBarisKontrak(ctx, tx, idUjiSebagian)
	if err != nil {
		t.Fatal(err)
	}
	for _, tabel := range []string{"T_TREATY_LIMITS", "T_TREATY_LIMIT_DETAIL", "T_TREATY_PORTFOLIO"} {
		if sesudah[tabel] != sebelum[tabel] {
			t.Errorf("%s: %d baris sebelum Save, %d sesudah — Save menghapus tab yang tidak dikirim",
				tabel, sebelum[tabel], sesudah[tabel])
		}
	}
}

// ⭐ Larik KOSONG tetap berwenang — grid yang pemakai kosongkan memang
// terhapus. Tanpa ini, perbaikan di atas akan membuat grid mustahil
// dikosongkan, dan itu cacat yang lain lagi.
func TestSaveDenganLarikKosongMenghapusBarisnya(t *testing.T) {
	g, tx, ctx := gudangUji(t)

	if _, err := g.MuatKontrak(ctx, tx, idUjiSebagian, dokumenUtuh()); err != nil {
		t.Fatalf("memuat dokumen utuh: %v", err)
	}
	doc := map[string]any{"ID": idUjiSebagian, "Limits": []any{}}
	if _, err := g.MuatKontrakSebagian(ctx, tx, idUjiSebagian, doc); err != nil {
		t.Fatalf("Save sebagian: %v", err)
	}
	cacah, err := g.CacahBarisKontrak(ctx, tx, idUjiSebagian)
	if err != nil {
		t.Fatal(err)
	}
	if cacah["T_TREATY_LIMITS"] != 0 || cacah["T_TREATY_LIMIT_DETAIL"] != 0 {
		t.Errorf("`Limits: []` tidak mengosongkan grid: %d induk, %d detail",
			cacah["T_TREATY_LIMITS"], cacah["T_TREATY_LIMIT_DETAIL"])
	}
	// Tab lain tidak ikut terbawa.
	if cacah["T_TREATY_PORTFOLIO"] == 0 {
		t.Error("Portfolio ikut terhapus padahal dokumen tidak membawanya")
	}
}

// ⛔ KOLOM tabel AKAR juga harus bertahan — cacat yang sama berpindah dari
// baris ke kolom bila tabel akar sekadar dihapus lalu disisipkan ulang.
func TestSaveSebagianMempertahankanKolomAkar(t *testing.T) {
	g, tx, ctx := gudangUji(t)

	if _, err := g.MuatKontrak(ctx, tx, idUjiSebagian, dokumenUtuh()); err != nil {
		t.Fatalf("memuat dokumen utuh: %v", err)
	}
	doc := map[string]any{"ID": idUjiSebagian, "TreatyContractName": "DIUBAH"}
	if _, err := g.MuatKontrakSebagian(ctx, tx, idUjiSebagian, doc); err != nil {
		t.Fatalf("Save sebagian: %v", err)
	}

	var nama, exclusions, information string
	q := `SELECT TREATYCONTRACTNAME, EXCLUSIONS, INFORMATION FROM ` +
		skemaUji(t) + `.T_TREATY_REVISION WHERE MASTERID = :1`
	if err := tx.QueryRowContext(ctx, q, idUjiSebagian).Scan(&nama, &exclusions, &information); err != nil {
		t.Fatalf("membaca baris akar: %v", err)
	}
	// Yang dokumen BAWA menang.
	if nama != "DIUBAH" {
		t.Errorf("TREATYCONTRACTNAME = %q, mau DIUBAH", nama)
	}
	// Yang dokumen TIDAK bawa bertahan.
	if exclusions != "SELURUH PENGECUALIAN" {
		t.Errorf("EXCLUSIONS = %q — kolom tab yang tidak dibuka hilang", exclusions)
	}
	if information != "KETERANGAN" {
		t.Errorf("INFORMATION = %q — kolom tab yang tidak dibuka hilang", information)
	}
}

// ⭐ Pemuat MASSAL tidak berubah: dokumen utuh dua kali menghasilkan cacah
// yang sama persis. Itulah idempoten yang perbaikan ini tidak boleh cabut.
func TestPemuatMassalTetapIdempotenDiOracle(t *testing.T) {
	g, tx, ctx := gudangUji(t)

	if _, err := g.MuatKontrak(ctx, tx, idUjiSebagian, dokumenUtuh()); err != nil {
		t.Fatal(err)
	}
	pertama, err := g.CacahBarisKontrak(ctx, tx, idUjiSebagian)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.MuatKontrak(ctx, tx, idUjiSebagian, dokumenUtuh()); err != nil {
		t.Fatal(err)
	}
	kedua, err := g.CacahBarisKontrak(ctx, tx, idUjiSebagian)
	if err != nil {
		t.Fatal(err)
	}
	for tabel, n := range pertama {
		if kedua[tabel] != n {
			t.Errorf("%s: %d lalu %d — pemuat massal tidak lagi idempoten", tabel, n, kedua[tabel])
		}
	}
}
