//go:build db

package handlers_test

// Uji wajib brief bab 4 (paket 4): produk `UJI-` yang disimpan layanan dibaca
// kembali lewat TEKS SQL KETIGA VIEW DEV di skema uji tiruan dan memberi nilai
// sama dengan masukannya. Teks view dibaca dari `docs/dba-view-produk-life.md`
// (definisi `ALL_VIEWS` DEV, 30-09-2026) - bukan disalin ke sini - supaya uji
// ini selalu memakai teks yang terdokumentasi.
//
// Juga P4: gagal menulis sisi inward membatalkan sisi umum (satu transaksi).

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
)

// teksView - `=== VIEW <NAMA>` → teks SQL-nya, dari dokumen DBA.
func teksView(t *testing.T) map[string]string {
	t.Helper()
	isi, err := os.ReadFile("../../docs/dba-view-produk-life.md")
	if err != nil {
		t.Fatalf("dokumen view: %v", err)
	}
	blok := regexp.MustCompile("(?s)```sql\n(.*?)```").FindStringSubmatch(string(isi))
	if blok == nil {
		t.Fatal("blok sql tidak ditemukan di dba-view-produk-life.md")
	}
	hasil := map[string]string{}
	bagian := regexp.MustCompile(`(?m)^=== VIEW (\w+)\s*$`).Split(blok[1], -1)
	nama := regexp.MustCompile(`(?m)^=== VIEW (\w+)\s*$`).FindAllStringSubmatch(blok[1], -1)
	for i, n := range nama {
		hasil[n[1]] = strings.TrimSpace(bagian[i+1])
	}
	for _, v := range []string{"DOCUMENTCLAIM_LIFE", "PRODUCTINWARD_LIFE", "PRODUCT_LIFE"} {
		if hasil[v] == "" {
			t.Fatalf("teks view %s tidak terbaca", v)
		}
	}
	return hasil
}

// pasangView membuat ketiga view di skema uji; nama tabel diarahkan ke skema uji.
func (u *ujiDB) pasangView(t *testing.T) {
	t.Helper()
	ganti := strings.NewReplacer(
		"pooldata.M_PRODUCT_LIFE", u.skema+".M_PRODUCT_LIFE",
		"FROM m_productinward_life a", "FROM "+u.skema+".M_PRODUCTINWARD_LIFE a",
		"FROM M_PRODUCT_LIFE a", "FROM "+u.skema+".M_PRODUCT_LIFE a")
	for nama, teks := range teksView(t) {
		u.exec(t, fmt.Sprintf(`CREATE OR REPLACE VIEW %s.%s AS %s`, u.skema, nama, ganti.Replace(teks)))
		nama := nama
		t.Cleanup(func() { _, _ = u.mentah.ExecContext(u.ctx, fmt.Sprintf(`DROP VIEW %s.%s`, u.skema, nama)) })
	}
}

func (u *ujiDB) isiMasterUji(t *testing.T) {
	t.Helper()
	for _, q := range []string{
		`INSERT INTO {s}.AGENT VALUES ('L0UJI1', 'UJI CEDING SATU', '1')`,
		`INSERT INTO {s}.AGENT VALUES ('L0SOB', 'UJI SOB', '1')`,
		`INSERT INTO {s}.CLIENT VALUES ('UJI-ORG-1', 'UJI PEMEGANG', 'LIFE')`,
		`INSERT INTO {s}.CURRENCY VALUES ('1', 'IDR')`,
		`INSERT INTO {s}.RIRISK_LIFE_SUMMARY VALUES ('1000117', 'UJI RISK')`,
		`INSERT INTO {s}.CAUSEOFLOSS_LIFE VALUES ('100004', 'ANY CAUSE')`,
	} {
		u.exec(t, q)
	}
}

// barisView membaca kolom bernama satu baris view sebagai teks.
func (u *ujiDB) barisView(t *testing.T, view, id string, kolom []string) map[string]string {
	t.Helper()
	q := fmt.Sprintf(`SELECT %s FROM %s.%s WHERE ID = :1`, strings.Join(kolom, ", "), u.skema, view)
	sel := make([]sql.NullString, len(kolom))
	ptr := make([]any, len(kolom))
	for i := range sel {
		ptr[i] = &sel[i]
	}
	if err := u.mentah.QueryRowContext(u.ctx, q, id).Scan(ptr...); err != nil {
		t.Fatalf("membaca view %s: %v", view, err)
	}
	hasil := map[string]string{}
	for i, k := range kolom {
		hasil[k] = sel[i].String
	}
	return hasil
}

func TestDBProdukUJIDibacaTigaViewSamaDenganMasukan(t *testing.T) {
	u := pasangDB(t)
	u.pasangView(t)
	u.isiMasterUji(t)
	kode, badan := u.kirim(t, "POST", pre+"/produk", badanLengkap)
	if kode != http.StatusOK || !strings.Contains(badan, `"id":"100044"`) {
		t.Fatalf("POST: %d %s", kode, badan)
	}
	in := u.barisView(t, "PRODUCTINWARD_LIFE", "100044", []string{"PRODUCTID", "INSURED", "BEGIN", "MATURE", "STNC",
		"CEDINGLIMIT", "MINAGE", "MAXAGE", "MAXSUMINSURED", "POLICYHODER", "POLICYHODERNAME", "CURRENCY",
		"MAXDATARECEIVE", "MAXEXPIREDCLAIM", "PAYMENT", "SUBJECTTO", "BROKERAGE"})
	mau := map[string]string{"PRODUCTID": "100044", "INSURED": "UJI TERTANGGUNG", "BEGIN": "01/03/2026",
		"MATURE": "28/02/2027", "STNC": "26/03/2026", "CEDINGLIMIT": "150000000.123456789", "MINAGE": "22",
		"MAXAGE": "70", "MAXSUMINSURED": "1175000000", "POLICYHODER": "UJI-ORG-1", "POLICYHODERNAME": "UJI PEMEGANG",
		"CURRENCY": "IDR", "MAXDATARECEIVE": "90", "MAXEXPIREDCLAIM": "180", "PAYMENT": "1", "SUBJECTTO": "UJI SYARAT",
		"BROKERAGE": "2.5"}
	for k, v := range mau {
		if in[k] != v {
			t.Errorf("PRODUCTINWARD_LIFE.%s = %q, masukan %q", k, in[k], v)
		}
	}
	um := u.barisView(t, "PRODUCT_LIFE", "100044", []string{"CEDING", "CEDINGID", "SOBNAME", "SOBID", "CAUSE",
		"CAUSEID", "PRODUCTNAME", "RIRISK", "RIRISKID", "RICOMM", "INWARDNAME", "POLICYHODER", "POLICYHODERNAME",
		"TREATYNUMBER", "CREATEOP", "UPDATEOP"})
	mauUm := map[string]string{"CEDING": "UJI CEDING SATU", "CEDINGID": "L0UJI1", "SOBNAME": "UJI SOB", "SOBID": "L0SOB",
		"CAUSE": "ANY CAUSE", "CAUSEID": "100004", "PRODUCTNAME": "UJI PRODUK", "RIRISK": "UJI RISK", "RIRISKID": "1000117",
		"RICOMM": "12.5", "INWARDNAME": "UJI PRODUK UJI PEMEGANG", "POLICYHODER": "UJI-ORG-1",
		"POLICYHODERNAME": "UJI PEMEGANG", "TREATYNUMBER": "UJI/001", "CREATEOP": "UJI-PELAKU", "UPDATEOP": "UJI-PELAKU"}
	for k, v := range mauUm {
		if um[k] != v {
			t.Errorf("PRODUCT_LIFE.%s = %q, masukan %q", k, um[k], v)
		}
	}
	// Claim Life membaca view inward dengan kunci ID produk (`GetProductName.xml`).
	if got := u.teks(t, `SELECT MAXEXPIREDCLAIM FROM {s}.PRODUCTINWARD_LIFE WHERE ID = '100044'`); got != "180" {
		t.Errorf("pembaca Claim Life (WHERE ID = produk): %q", got)
	}
}

func TestDBSimpanAtomikDuaTabel(t *testing.T) {
	u := pasangDB(t)
	u.isiMasterUji(t)
	// Sisi inward DIPAKSA gagal: constraint yang menolak setiap baris.
	u.exec(t, `ALTER TABLE {s}.M_PRODUCTINWARD_LIFE ADD CONSTRAINT UJI_TOLAK_INWARD CHECK (ID = 'X')`)
	kode, badan := u.kirim(t, "POST", pre+"/produk", badanLengkap)
	if kode == http.StatusOK {
		t.Fatalf("simpan harus gagal: %s", badan)
	}
	if n := u.cacah(t, "M_PRODUCT_LIFE", ""); n != 0 {
		t.Errorf("P4: gagal sisi inward membatalkan sisi umum - %d baris M_PRODUCT_LIFE tersisa", n)
	}
}
