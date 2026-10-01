//go:build db

package handlers_test

// Seam HTTP simpan produk terhadap skema uji Oracle NYATA (paket 3).

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
	if n := u.cacah(t, "M_PRODUCT_LIFE", "ID = '100044' AND PRODUCTNAME = 'UJI PRODUK' AND BEGIN_DATE = DATE '2026-03-01' "+
		"AND JSON_VALUE(JSONDATA, '$.POLICYHODER') = 'UJI-ORG-1' AND JSON_VALUE(JSONDATA, '$.CREATEOP') = 'UJI-PELAKU'"); n != 1 {
		t.Errorf("baris baru + kolom datar + JSON: %d", n)
	}
	kode, badan = u.kirim(t, "PUT", pre+"/produk/100044", strings.Replace(badanUji, "UJI PRODUK", "UJI UBAH", 1))
	if kode != http.StatusOK {
		t.Fatalf("PUT: %d %s", kode, badan)
	}
	if n := u.cacah(t, "M_PRODUCT_LIFE", ""); n != 1 {
		t.Errorf("upsert dikunci ID: %d baris, mau 1", n)
	}
	if got := u.teks(t, `SELECT PRODUCTNAME FROM {s}.M_PRODUCT_LIFE WHERE ID = '100044'`); got != "UJI UBAH" {
		t.Errorf("kolom datar PRODUCTNAME ikut diperbarui: %q", got)
	}
}

func TestDBIdentitasTerpakaiDitolakTerang(t *testing.T) {
	u := pasangDB(t)
	u.exec(t, `INSERT INTO {s}.M_PRODUCTINWARD_LIFE (ID, JSONDATA) VALUES ('100044', '{}')`)
	u.isiMasterUji(t)
	kode, badan := u.kirim(t, "POST", pre+"/produk", badanUji)
	if kode != http.StatusInternalServerError || !strings.Contains(badan, "already used") {
		t.Errorf("ID baru yang sudah dipakai inward harus gagal terang: %d %s", kode, badan)
	}
	if n := u.cacah(t, "M_PRODUCT_LIFE", ""); n != 0 {
		t.Errorf("nol baris: %d", n)
	}
}
