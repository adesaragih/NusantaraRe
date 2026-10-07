//go:build db

package repository_test

// Uji seam repository lawan Oracle SUNGGUHAN yang menegaskan NILAI KOLOM
// LANGSUNG (spec-penyimpanan §7 "Pulang-pergi saja tidak cukup"): NOURUT
// sesudah baris dihapus (AC 9, tiket 17), kode bernol-depan (AC 12-13, tiket
// 18), `""` lawan `"0"` pada penanda (AC 15 RALAT P11, tiket 18), dan
// kegagalan menulis tabel ANAK yang membatalkan induk (AC 45, tiket 20).
//
// ⛔ Belum pernah dijalankan: skema uji K11 kosong (PERMINTAAN C4). Seluruh
// tabel di sini dibuat migrasi modul ini (320-327) - tidak ada tiruan.
// Pemanggil `skemauji.Buka()` tetap satu (`pasang`). Fixture berawalan UJI-.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	intidb "nusantarare/inti/backend/db"
	"nusantarare/inti/backend/galat"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

// kolomTeks membaca satu kolom sebagai teks apa adanya ("<NULL>" bila NULL).
func kolomTeks(t *testing.T, ctx context.Context, sqlDB *sql.DB, q string, args ...any) string {
	t.Helper()
	var v sql.NullString
	if err := sqlDB.QueryRowContext(ctx, q, args...).Scan(&v); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
	if !v.Valid {
		return "<NULL>"
	}
	return v.String
}

// barisAngsuran membaca (NOURUT, INSTALLMENT_NO) T_POLIS_INSTALMENT satu
// polis, berurut NOURUT.
func barisAngsuran(t *testing.T, ctx context.Context, sqlDB *sql.DB, skema, id string) []string {
	t.Helper()
	rows, err := sqlDB.QueryContext(ctx, fmt.Sprintf(`SELECT TO_CHAR(NOURUT), TO_CHAR(INSTALLMENT_NO)
		FROM %s.T_POLIS_INSTALMENT WHERE POLIS_ID = :1 ORDER BY NOURUT`, skema), id)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var no, inst sql.NullString
		if err := rows.Scan(&no, &inst); err != nil {
			t.Fatal(err)
		}
		out = append(out, no.String+":"+inst.String)
	}
	return out
}

func simpanKasusBaru(t *testing.T, ctx context.Context, d *intidb.DB, g *repository.Gudang, id string, h *models.Halaman) {
	t.Helper()
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error {
		if err := g.SisipKasus(ctx, tx, id, "UJI-AKUN", "UJI NAMA"); err != nil {
			return err
		}
		return g.SimpanHalaman(ctx, tx, id, h)
	}); err != nil {
		t.Fatal(err)
	}
}

func angsuranUji(no ...string) []models.Baris {
	var b []models.Baris
	for _, n := range no {
		b = append(b, models.Baris{"InstallmentNo": n, "DueDate": "2026-11-01", "Premium": "1.5"})
	}
	return b
}

// Spec-penyimpanan AC 9: "Pada satu polis NB dengan tiga baris angsuran,
// menghapus baris kedua menghasilkan NOURUT 1 dan 2. Test yang menemukan 1 dan
// 3 gagal." (ID-12) - kolom NOURUT dibaca LANGSUNG.
func TestNourutDinomoriUlangSesudahBarisKeduaDihapus(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-NOURUT"
	h := models.HalamanBaru()
	h.Setel("PositionNote", models.PosisiAdmin)
	h.SetelDaftar(models.DaftarAngsuran, angsuranUji("1", "2", "3"))
	simpanKasusBaru(t, ctx, d, g, id, h)
	if got := fmt.Sprint(barisAngsuran(t, ctx, sqlDB, skema, id)); got != "[1:1 2:2 3:3]" {
		t.Fatalf("tiga angsuran: %s", got)
	}
	// baris kedua (InstallmentNo 2) dihapus di layar, lalu disimpan
	h.SetelDaftar(models.DaftarAngsuran, angsuranUji("1", "3"))
	if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SimpanHalaman(ctx, tx, id, h) }); err != nil {
		t.Fatal(err)
	}
	// NOURUT 1, 2 (bukan 1, 3); baris ber-InstallmentNo 3 kini NOURUT 2.
	if got := fmt.Sprint(barisAngsuran(t, ctx, sqlDB, skema, id)); got != "[1:1 2:3]" {
		t.Fatalf("sesudah hapus baris kedua: %s, harap [1:1 2:3]", got)
	}
}

// Spec-penyimpanan AC 12: "`GroupPanel` bernilai "006" tersimpan dan terbaca
// kembali sebagai "006". Test yang menemukan "6" atau 6 gagal." (ID-16), dan
// AC 13 (`BusinessOldId` "01") - dibaca dari KOLOM, beserta tipenya.
func TestKodeBernolDepanUtuhDiKolom(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-KODE"
	h := models.HalamanBaru()
	h.Setel("PositionNote", models.PosisiAdmin)
	h.Setel("Quotation.GroupPanel", "006")
	h.Setel("Quotation.BusinessOldId", "01")
	simpanKasusBaru(t, ctx, d, g, id, h)
	for kolom, harap := range map[string]string{"GROUP_PANEL": "006", "BUSINESS_OLD_ID": "01"} {
		if got := kolomTeks(t, ctx, sqlDB, fmt.Sprintf(`SELECT %s FROM %s.T_POLIS_QUOTATION WHERE POLIS_ID = :1`, kolom, skema), id); got != harap {
			t.Errorf("kolom %s = %q, harap %q", kolom, got, harap)
		}
		if tipe := kolomTeks(t, ctx, sqlDB, `SELECT DATA_TYPE FROM ALL_TAB_COLUMNS WHERE OWNER = UPPER(:1)
			AND TABLE_NAME = 'T_POLIS_QUOTATION' AND COLUMN_NAME = :2`, skema, kolom); tipe != "VARCHAR2" {
			t.Errorf("tipe %s = %s, harap VARCHAR2 (kode tetap teks)", kolom, tipe)
		}
	}
	b, err := g.BacaHalaman(ctx, nil, id)
	if err != nil {
		t.Fatal(err)
	}
	if b.Ambil("Quotation.GroupPanel") != "006" || b.Ambil("PolicyTreatyIn.QuotationData.GroupPanel") != "006" {
		t.Errorf("GroupPanel terbaca %q / %q", b.Ambil("Quotation.GroupPanel"), b.Ambil("PolicyTreatyIn.QuotationData.GroupPanel"))
	}
}

// Spec-penyimpanan AC 15 (bunyi baru RALAT P11): Oracle menyamakan teks
// kosong dengan NULL - `IsApproved` "" tersimpan NULL dan terbaca kembali "";
// "0" tersimpan '0' dan terbaca "0". Keduanya TIDAK pernah menyatu.
func TestIsApprovedKosongDanNolTetapBerbeda(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-PENANDA"
	h := models.HalamanBaru()
	h.Setel("PositionNote", models.PosisiAdmin)
	h.Setel("PolicyTreatyIn.IsApproved", "")
	simpanKasusBaru(t, ctx, d, g, id, h)
	q := fmt.Sprintf(`SELECT IS_APPROVED FROM %s.T_GENERAL_POLIS_TREATY WHERE ID = :1`, skema)
	for _, c := range []struct{ masuk, kolom string }{{"", "<NULL>"}, {"0", "0"}, {"", "<NULL>"}} {
		h.Setel("PolicyTreatyIn.IsApproved", c.masuk)
		if err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SimpanHalaman(ctx, tx, id, h) }); err != nil {
			t.Fatal(err)
		}
		if got := kolomTeks(t, ctx, sqlDB, q, id); got != c.kolom {
			t.Errorf("IsApproved %q: kolom %s, harap %s", c.masuk, got, c.kolom)
		}
		b, err := g.BacaHalaman(ctx, nil, id)
		if err != nil {
			t.Fatal(err)
		}
		if got := b.Ambil("PolicyTreatyIn.IsApproved"); got != c.masuk {
			t.Errorf("IsApproved %q terbaca kembali %q", c.masuk, got)
		}
	}
}

// Spec-penyimpanan AC 45: "Kegagalan menulis tabel anak mana pun membatalkan
// seluruh penyimpanan. Test yang menemukan baris induk tersisa gagal." (ID-32)
//
// Kegagalan lahir di ORACLE pada baris anak kedua: `InstallmentNo`
// 12345678901 lolos pemeriksaan Go (bilangan bulat) tetapi melampaui
// `T_POLIS_INSTALMENT.INSTALLMENT_NO NUMBER(10)` (ORA-01438) - SESUDAH induk
// T_GENERAL_POLIS_TREATY diperbarui, anak lama dihapus, T_POLIS_QUOTATION ditulis
// ulang, dan baris anak pertama disisipkan, di transaksi yang sama.
func TestGagalTulisTabelAnakMembatalkanInduk(t *testing.T) {
	sqlDB, skema, ctx, d := pasang(t)
	g := repository.Baru(d)
	const id = "UJI-NB-ANAK-GAGAL"
	h := models.HalamanBaru()
	h.Setel("PositionNote", models.PosisiAdmin)
	h.Setel("PolicyTreatyIn.PremiOgp", "100")
	h.Setel("Quotation.GroupPanel", "006")
	h.SetelDaftar(models.DaftarAngsuran, angsuranUji("1", "2"))
	simpanKasusBaru(t, ctx, d, g, id, h)

	ubah := models.HalamanBaru()
	ubah.Setel("PositionNote", models.PosisiAdmin)
	ubah.Setel("PolicyTreatyIn.PremiOgp", "999")
	ubah.Setel("Quotation.GroupPanel", "007")
	ubah.SetelDaftar(models.DaftarAngsuran, angsuranUji("1", "12345678901"))
	err := dalamTx(t, ctx, d, func(tx *intidb.Tx) error { return g.SimpanHalaman(ctx, tx, id, ubah) })
	if err == nil {
		t.Fatal("baris anak yang ditolak Oracle harus menggagalkan penyimpanan")
	}
	if errors.Is(err, galat.ErrPermintaanTidakSah) {
		t.Fatalf("galat lahir di Go, bukan di tabel anak: %v", err)
	}
	premi := kolomTeks(t, ctx, sqlDB, fmt.Sprintf(`SELECT %s FROM %s.T_GENERAL_POLIS_TREATY WHERE ID = :1`,
		fmt.Sprintf(intidb.FmtDesimal, "PREMI_OGP"), skema), id)
	if premi != "100" {
		t.Errorf("induk T_GENERAL_POLIS_TREATY.PREMI_OGP = %s sesudah gagal, harap 100", premi)
	}
	if gp := kolomTeks(t, ctx, sqlDB, fmt.Sprintf(`SELECT GROUP_PANEL FROM %s.T_POLIS_QUOTATION WHERE POLIS_ID = :1`, skema), id); gp != "006" {
		t.Errorf("T_POLIS_QUOTATION.GROUP_PANEL = %s sesudah gagal, harap 006", gp)
	}
	if got := fmt.Sprint(barisAngsuran(t, ctx, sqlDB, skema, id)); got != "[1:1 2:2]" {
		t.Errorf("anak lama T_POLIS_INSTALMENT sesudah gagal: %s, harap [1:1 2:2]", got)
	}
}
