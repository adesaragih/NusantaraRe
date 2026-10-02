//go:build db

package handlers_test

// Produk `UJI-` yang disimpan layanan dibaca kembali dari KOLOM tabel flat di Oracle sungguhan dan sama dengan
// masukannya (dulu lewat teks ketiga view DEV - view tidak dibangun ulang, K7 02-10-2026). Juga P4: gagal menulis
// baris anak membatalkan baris induk (satu transaksi).

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"nusantarare/inti/backend/db"
)

func TestDBProdukUJITersimpanSamaDenganMasukan(t *testing.T) {
	u := pasangDB(t)
	u.isiMasterUji(t)
	kode, badan := u.kirim(t, "POST", pre+"/produk", badanLengkap)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"100044"`) {
		t.Fatalf("POST: %d %s", kode, badan)
	}
	angka := func(k string) string { return fmt.Sprintf(db.FmtDesimal, k) }
	tanggal := func(k string) string { return "TO_CHAR(" + k + ", 'YYYY-MM-DD')" }
	for kolom, mau := range map[string]string{
		"PRODUCTNAME": "UJI PRODUK", "CEDING": "UJI CEDING SATU", "CEDINGID": "L0UJI1", "SOBNAME": "UJI SOB",
		"SOBID": "L0SOB", "CAUSE": "ANY CAUSE", "CAUSEID": "100004", "RIRISK": "UJI RISK", "RIRISKID": "1000117",
		angka("RICOMM"): "12.5", "INWARDNAME": "UJI PRODUK UJI PEMEGANG", "POLICYHOLDER": "UJI-ORG-1",
		"POLICYHOLDERNAME": "UJI PEMEGANG", "TREATYNUMBER": "UJI/001", "CREATEOP": "UJI-PELAKU", "UPDATEOP": "UJI-PELAKU",
		"INSURED": "UJI TERTANGGUNG", tanggal("BEGIN_DATE"): "2026-03-01", tanggal("MATURE"): "2027-02-28",
		tanggal("STNC"): "2026-03-26", angka("CEDINGLIMIT"): "150000000.12345678", angka("MINAGE"): "22",
		angka("MAXAGE"): "70", angka("MAXSUMINSURED"): "1175000000", "CURRENCY": "IDR", angka("MAXDATARECEIVE"): "90",
		angka("MAXEXPIREDCLAIM"): "180", "PAYMENT": "1", "SUBJECTTO": "UJI SYARAT", angka("BROKERAGE"): "2.5",
	} {
		if got := u.teks(t, `SELECT `+kolom+` FROM {s}.M_PRODUCTNAME_LIFE WHERE ID = '100044'`); got != mau {
			t.Errorf("M_PRODUCTNAME_LIFE.%s = %q, masukan %q", kolom, got, mau)
		}
	}
	if got := u.teks(t, `SELECT SUGGEST FROM {s}.M_PRODUCTNAME_LIFE_COMMENT WHERE PRODUCTID = '100044' AND URUT = 1`); got != "UJI komentar" {
		t.Errorf("baris komentar simpan: %q", got)
	}
}

func TestDBSimpanAtomikIndukDanAnak(t *testing.T) {
	u := pasangDB(t)
	u.isiMasterUji(t)
	// Penulis anak DIPAKSA gagal: setiap simpan menulis satu baris komentar, constraint ini menolaknya.
	u.exec(t, `ALTER TABLE {s}.M_PRODUCTNAME_LIFE_COMMENT ADD CONSTRAINT UJI_TOLAK_KOMENTAR CHECK (URUT = 0)`)
	kode, badan := u.kirim(t, "POST", pre+"/produk", badanLengkap)
	if kode == http.StatusOK {
		t.Fatalf("simpan harus gagal: %s", badan)
	}
	if n := u.cacah(t, "M_PRODUCTNAME_LIFE", ""); n != 0 {
		t.Errorf("P4: gagal menulis anak membatalkan induk - %d baris induk tersisa", n)
	}
	if n := u.cacah(t, "M_PRODUCTNAME_LIFE_DOCCLAIM", ""); n != 0 {
		t.Errorf("P4: baris anak lain ikut dibatalkan - %d tersisa", n)
	}
}
