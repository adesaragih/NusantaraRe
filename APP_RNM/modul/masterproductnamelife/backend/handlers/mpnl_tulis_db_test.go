//go:build db

package handlers_test

// Seam HTTP simpan produk terhadap skema uji Oracle NYATA (paket 3) - sejak 02-10-2026 ke tabel FLAT; kedua tabel
// JSON warisan tidak pernah ditulis.

import (
	"net/http"
	"strings"
	"testing"
)

func (u *ujiDB) cacah(t *testing.T, tabel, saring string, args ...any) int {
	t.Helper()
	q := `SELECT COUNT(*) FROM {s}.` + tabel
	if saring != "" {
		q += " WHERE " + saring
	}
	var n int
	if err := u.mentah.QueryRowContext(u.ctx, strings.ReplaceAll(q, "{s}", u.skema), args...).Scan(&n); err != nil {
		t.Fatalf("mencacah %s: %v", tabel, err)
	}
	return n
}

func TestDBSimpanBaruDariSequenceLaluUpsertDikunciID(t *testing.T) {
	u := pasangDB(t)
	u.isiMasterUji(t)
	kode, badan := u.kirim(t, "POST", pre+"/produk", badanUji)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"100044"`) {
		t.Fatalf("POST: %d %s", kode, badan)
	}
	if n := u.cacah(t, "M_PRODUCTNAME_LIFE", "ID = '100044' AND PRODUCTNAME = 'UJI PRODUK' "+
		"AND POLICYHOLDER = 'UJI-ORG-1' AND CREATEOP = 'UJI-PELAKU'"); n != 1 {
		t.Errorf("baris induk baru: %d", n)
	}
	kode, badan = u.kirim(t, "PUT", pre+"/produk/100044", strings.Replace(badanUji, "UJI PRODUK", "UJI UBAH", 1))
	if kode != http.StatusOK {
		t.Fatalf("PUT: %d %s", kode, badan)
	}
	if n := u.cacah(t, "M_PRODUCTNAME_LIFE", ""); n != 1 {
		t.Errorf("upsert dikunci ID: %d baris, mau 1", n)
	}
	if got := u.teks(t, `SELECT PRODUCTNAME FROM {s}.M_PRODUCTNAME_LIFE WHERE ID = '100044'`); got != "UJI UBAH" {
		t.Errorf("PRODUCTNAME diperbarui: %q", got)
	}
	// Setiap simpan menambah satu baris komentar (OQ-MPNL-14): dua simpan = URUT 1 dan 2.
	if n := u.cacah(t, "M_PRODUCTNAME_LIFE_COMMENT", "PRODUCTID = '100044'"); n != 2 {
		t.Errorf("baris komentar: %d, mau 2", n)
	}
	if n := u.cacah(t, "M_PRODUCT_LIFE", "") + u.cacah(t, "M_PRODUCTINWARD_LIFE", ""); n != 0 {
		t.Errorf("tabel JSON warisan tidak ditulis lagi: %d baris", n)
	}
}

// Keputusan work owner 02-10-2026 ("simpan ke table flat semua"): penerbitan ID hanya memeriksa induk flat - kedua
// tabel JSON warisan tidak dibaca dan tidak disentuh aplikasi, juga tidak dibutuhkan (dibuang DBA pun simpan tetap jalan).
func TestDBIdentitasHanyaDariIndukFlat(t *testing.T) {
	u := pasangDB(t)
	u.exec(t, `INSERT INTO {s}.M_PRODUCTINWARD_LIFE (ID, JSONDATA) VALUES ('100044', '{}')`)
	u.isiMasterUji(t)
	kode, badan := u.kirim(t, "POST", pre+"/produk", badanUji)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"100044"`) {
		t.Errorf("ID yang hanya ada di tabel JSON tidak diperiksa: produk baru 100044: %d %s", kode, badan)
	}
	if n := u.cacah(t, "M_PRODUCTNAME_LIFE", ""); n != 1 {
		t.Errorf("satu baris induk: %d", n)
	}
	if got := u.teks(t, `SELECT JSONDATA FROM {s}.M_PRODUCTINWARD_LIFE WHERE ID = '100044'`); got != "{}" {
		t.Errorf("baris warisan tidak tersentuh: %q", got)
	}
	u.exec(t, `DROP TABLE {s}.M_PRODUCTINWARD_LIFE PURGE`)
	u.exec(t, `DROP TABLE {s}.M_PRODUCT_LIFE PURGE`)
	if kode, badan := u.kirim(t, "POST", pre+"/produk", badanUji); kode != http.StatusOK || !strings.Contains(badan, `"id":"100045"`) {
		t.Errorf("tanpa tabel JSON: produk baru tetap terbit: %d %s", kode, badan)
	}
}

// Kolom induk dan baris anak ditulis di Oracle sungguhan - bukan hanya di tiruan.
func TestDBKolomFlatDanBarisAnak(t *testing.T) {
	u := pasangDB(t)
	u.isiMasterUji(t)
	if kode, badan := u.kirim(t, "POST", pre+"/produk", badanLengkap); kode != http.StatusOK {
		t.Fatalf("POST: %d %s", kode, badan)
	}
	if n := u.cacah(t, "M_PRODUCTNAME_LIFE", "ID = '100044' AND RIRISKID = '1000117' AND RIRISK = 'UJI RISK'"); n != 1 {
		t.Errorf("kolom RIRISKID, RIRISK: %d", n)
	}
	if n := u.cacah(t, "M_PRODUCTNAME_LIFE_DOCCLAIM", "PRODUCTID = '100044'"); n != 2 {
		t.Errorf("dua baris DOCUMENT CLAIM: %d", n)
	}
	if got := u.teks(t, `SELECT DOCUMENT FROM {s}.M_PRODUCTNAME_LIFE_DOCCLAIM WHERE PRODUCTID = '100044' AND URUT = 2`); got != "UJI DOK B" {
		t.Errorf("urutan baris grid = URUT: %q", got)
	}
}
