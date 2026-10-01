package repository_test

// Migrasi 059 - `T_WORK_POLIS` seragam dengan `T_WORK_CLAIM` (brief seragam
// kolom, keputusan work owner 01-10-2026). TANPA Oracle.
//
// Untuk apa berkas ini: menagih bentuk AKHIR kedua tabel kerja sesudah SELURUH
// migrasi maju setiap modul - bukan isi satu berkas - serta urutan langkah 059
// dan jalur mundurnya.
//
// ⛔ Kolom yang wajib SAMA nama DAN tipe dengan `T_WORK_CLAIM`: `ID`,
// `COVER_KEY`, `CREATE_OP`, `CREATE_OP_NAME`, `TGL_CREATE`, `TGL_UPDATE`.
// `STATUS_WORK` dan `LINI` sama NAMA, tetapi polis tetap VARCHAR2(255):
// memperkecil kolom dapat gagal pada data yang lebih panjang (brief §2).
// `POSITION` TIDAK menjadi `PY_POSITION` - dua properti Pega yang berbeda
// (brief §1).
//
// Dibaca sesudah: migrations/059_seragam_kolom_t_work_polis.sql.

import (
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"nusantarare/inti/backend/migrasi"
	"nusantarare/uji/skemauji"
)

const berkas059 = "059_seragam_kolom_t_work_polis.sql"

var (
	polaBuatTabelKerja   = regexp.MustCompile(`(?s)^CREATE TABLE \{skema\}\.(\w+) \((.*)\n\)$`)
	polaTambahKolomKerja = regexp.MustCompile(`(?s)^ALTER TABLE \{skema\}\.(\w+) ADD \((.*)\n\)$`)
	polaKolomTipeKerja   = regexp.MustCompile(`^\s*([A-Z0-9_]+)\s+([A-Z0-9_]+(?:\(\d+(?:,\d+)?\))?)`)
)

// tabelKerja - dua tabel yang dibandingkan.
var tabelKerja = map[string]bool{"T_WORK_CLAIM": true, "T_WORK_POLIS": true}

// bentukAkhirTabelKerja membaca kolom dan tipe kedua tabel kerja sesudah
// seluruh migrasi maju SETIAP modul (`skemauji.SumberMigrasi`, daftar yang sama
// dengan `-migrate`), urut nomor: CREATE, ADD, lalu DROP COLUMN - juga yang di
// dalam blok berpelindung katalog (`migrasi.KolomAlterBuang`).
func bentukAkhirTabelKerja(t *testing.T) map[string]map[string]string {
	t.Helper()
	langkah, err := migrasi.Daftar(false, skemauji.SumberMigrasi()...)
	if err != nil {
		t.Fatal(err)
	}
	hasil := map[string]map[string]string{}
	baca := func(tabel, isi string) {
		if hasil[tabel] == nil {
			hasil[tabel] = map[string]string{}
		}
		for _, b := range strings.Split(isi, "\n") {
			if m := polaKolomTipeKerja.FindStringSubmatch(b); m != nil && m[1] != "CONSTRAINT" && m[1] != "REFERENCES" {
				hasil[tabel][m[1]] = m[2]
			}
		}
	}
	for _, m := range langkah {
		for _, p := range m.Pernyataan {
			p = strings.TrimSpace(p)
			if c := polaBuatTabelKerja.FindStringSubmatch(p); c != nil && tabelKerja[c[1]] {
				baca(c[1], c[2])
				continue
			}
			if c := polaTambahKolomKerja.FindStringSubmatch(p); c != nil && tabelKerja[c[1]] {
				baca(c[1], c[2])
				continue
			}
			if tabel, kolom := migrasi.KolomAlterBuang(p); tabelKerja[tabel] {
				for _, k := range kolom {
					delete(hasil[tabel], k)
				}
			}
		}
	}
	return hasil
}

func TestMigrasi059WorkPolisSeragamDenganWorkClaim(t *testing.T) {
	b := bentukAkhirTabelKerja(t)
	polis, klaim := b["T_WORK_POLIS"], b["T_WORK_CLAIM"]
	// ⚠️ Instrumen diuji atas jawaban yang diketahui: klaim 001 + 016 + 017
	// (brief §0) = 14 kolom. Pengurai yang rusak akan meluluskan apa pun.
	if len(klaim) != 14 {
		t.Fatalf("T_WORK_CLAIM terbaca %d kolom, mau 14 (brief §0): %v", len(klaim), klaim)
	}
	mau := map[string]string{
		"ID": "VARCHAR2(32)", "LINI": "VARCHAR2(255)", "POSITION": "VARCHAR2(255)",
		"STATUS_WORK": "VARCHAR2(255)", "FLAG_ONGOING_POLICY": "VARCHAR2(1)",
		"COVER_KEY": "VARCHAR2(32)", "CREATE_OP": "VARCHAR2(64)", "CREATE_OP_NAME": "VARCHAR2(128)",
		"TGL_CREATE": "DATE", "TGL_UPDATE": "DATE",
	}
	if !reflect.DeepEqual(polis, mau) {
		t.Errorf("bentuk akhir T_WORK_POLIS:\n dapat %v\n mau   %v (brief §2)", polis, mau)
	}
	for _, k := range []string{"ID", "COVER_KEY", "CREATE_OP", "CREATE_OP_NAME", "TGL_CREATE", "TGL_UPDATE"} {
		if polis[k] == "" || polis[k] != klaim[k] {
			t.Errorf("%s: polis %q, klaim %q - nama dan tipe wajib sama", k, polis[k], klaim[k])
		}
	}
	for _, k := range []string{"STATUS_WORK", "LINI"} {
		if klaim[k] == "" || polis[k] == "" {
			t.Errorf("%s wajib ada di kedua tabel kerja (polis %q, klaim %q)", k, polis[k], klaim[k])
		}
	}
	if _, ada := polis["STATUS"]; ada {
		t.Error("T_WORK_POLIS.STATUS masih ada; namanya STATUS_WORK sejak 059")
	}
	if _, ada := polis["PY_POSITION"]; ada {
		t.Error("POSITION polis tidak boleh menjadi PY_POSITION - properti Pega berbeda (brief §1)")
	}
}

// bacaPerintah059 membaca pernyataan satu arah 059 dari folder modul ini.
func bacaPerintah059(t *testing.T, nama string) []string {
	t.Helper()
	p, err := migrasi.PernyataanLangkah(os.DirFS(".."), nama)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// Urutan maju: satu-satunya pernyataan TANPA pelindung adalah penambah kolom
// di urutan pertama - penjaga `pelanggaranBlokPLSQL` menolak kolom baru lewat
// blok, sebab blok tidak terlihat pembanding STRUKTUR.
func TestMigrasi059LangkahMaju(t *testing.T) {
	p := bacaPerintah059(t, berkas059)
	if len(p) != 6 {
		t.Fatalf("059 terbaca %d pernyataan, mau 6", len(p))
	}
	tabel, kolom := migrasi.KolomAlterTambah(p[0])
	if tabel != "T_WORK_POLIS" || !reflect.DeepEqual(kolom,
		[]string{"STATUS_WORK", "COVER_KEY", "CREATE_OP", "CREATE_OP_NAME", "TGL_CREATE", "TGL_UPDATE"}) {
		t.Errorf("pernyataan 1 menambah %s %v", tabel, kolom)
	}
	blok := func(i int) migrasi.PerintahKatalog {
		t.Helper()
		pk, ok := migrasi.BacaPerintahKatalog(p[i])
		if !ok {
			t.Fatalf("pernyataan %d bukan blok berpelindung katalog:\n%s", i+1, p[i])
		}
		return pk
	}
	salin, buang, fk := blok(1), blok(2), blok(3)
	if salin.Objek != "STATUS" || !salin.BilaAda || salin.Perintah != "UPDATE {skema}.T_WORK_POLIS SET STATUS_WORK = STATUS" {
		t.Errorf("salin STATUS -> STATUS_WORK: %+v", salin)
	}
	// ⛔ Salin DULU, baru buang: urutan sebaliknya membuang isi status kerja.
	if buang.Objek != "STATUS" || !buang.BilaAda || buang.Perintah != "ALTER TABLE {skema}.T_WORK_POLIS DROP COLUMN STATUS" {
		t.Errorf("buang STATUS: %+v", buang)
	}
	if fk.Katalog != "ALL_CONSTRAINTS" || fk.Objek != "FK_WORK_POLIS_COVER_KEY" || fk.BilaAda ||
		fk.Perintah != "ALTER TABLE {skema}.T_WORK_POLIS ADD CONSTRAINT FK_WORK_POLIS_COVER_KEY FOREIGN KEY (COVER_KEY) REFERENCES {skema}.T_WORK_POLIS (ID)" {
		t.Errorf("FK COVER_KEY (tanpa ON DELETE, sama dengan klaim butir d): %+v", fk)
	}
	if p[4] != "CREATE INDEX {skema}.IX_WORK_POLIS_COVER_KEY ON {skema}.T_WORK_POLIS (COVER_KEY)" {
		t.Errorf("index COVER_KEY: %q", p[4])
	}
	isi := p[5]
	for _, w := range []string{
		"UPDATE {skema}.T_WORK_POLIS w",
		"p.ID_PEGA = w.ID",                // penghubung - polis_kasus.go sqlSisipPremiumListKosong
		"= 1",                             // hanya tepat satu pasangan
		"p.TGL_INPUT",                     // TGL_CREATE <- TGL_INPUT
		"LENGTH(p.CREATE_OP_NAME) <= 128", // 255 -> 128 tanpa memotong
		"w.TGL_CREATE IS NULL",            // baris baru tidak ditimpa
	} {
		if !strings.Contains(isi, w) {
			t.Errorf("pengisian baris lama tanpa %q:\n%s", w, isi)
		}
	}
	// ⛔ CREATE_OP baris lama TETAP kosong (ADR-U-0027): T_PREMIUM_LIST tidak
	// menyimpan akun pembuat. TGL_UPDATE pun tidak: kolom itu tidak ada di sana.
	for _, w := range []string{"CREATE_OP,", "CREATE_OP =", "TGL_UPDATE"} {
		if strings.Contains(isi, w) {
			t.Errorf("pengisian baris lama tidak boleh menyentuh %q:\n%s", w, isi)
		}
	}
	semua := strings.Join(p, "\n")
	for _, w := range []string{"RENAME", "ON DELETE", "COMMIT"} {
		if strings.Contains(semua, w) {
			t.Errorf("059 maju memuat %q", w)
		}
	}
}

// Jalur mundur: seluruhnya berpelindung - index dan FK dulu, lima kolom baru,
// lalu STATUS_WORK kembali menjadi STATUS (isinya ikut kembali).
func TestMigrasi059LangkahMundur(t *testing.T) {
	p := bacaPerintah059(t, strings.TrimSuffix(berkas059, ".sql")+"_down.sql")
	mau := []struct{ katalog, objek, perintah string }{
		{"ALL_INDEXES", "IX_WORK_POLIS_COVER_KEY", "DROP INDEX {skema}.IX_WORK_POLIS_COVER_KEY"},
		{"ALL_CONSTRAINTS", "FK_WORK_POLIS_COVER_KEY", "ALTER TABLE {skema}.T_WORK_POLIS DROP CONSTRAINT FK_WORK_POLIS_COVER_KEY"},
		{"ALL_TAB_COLUMNS", "TGL_UPDATE", "ALTER TABLE {skema}.T_WORK_POLIS DROP COLUMN TGL_UPDATE"},
		{"ALL_TAB_COLUMNS", "TGL_CREATE", "ALTER TABLE {skema}.T_WORK_POLIS DROP COLUMN TGL_CREATE"},
		{"ALL_TAB_COLUMNS", "CREATE_OP_NAME", "ALTER TABLE {skema}.T_WORK_POLIS DROP COLUMN CREATE_OP_NAME"},
		{"ALL_TAB_COLUMNS", "CREATE_OP", "ALTER TABLE {skema}.T_WORK_POLIS DROP COLUMN CREATE_OP"},
		{"ALL_TAB_COLUMNS", "COVER_KEY", "ALTER TABLE {skema}.T_WORK_POLIS DROP COLUMN COVER_KEY"},
		{"ALL_TAB_COLUMNS", "STATUS_WORK", "ALTER TABLE {skema}.T_WORK_POLIS RENAME COLUMN STATUS_WORK TO STATUS"},
	}
	if len(p) != len(mau) {
		t.Fatalf("059_down terbaca %d pernyataan, mau %d", len(p), len(mau))
	}
	for i, m := range mau {
		pk, ok := migrasi.BacaPerintahKatalog(p[i])
		if !ok || pk.Katalog != m.katalog || pk.Objek != m.objek || !pk.BilaAda || pk.Perintah != m.perintah {
			t.Errorf("pernyataan mundur %d: %+v (ok=%v), mau %+v", i+1, pk, ok, m)
		}
	}
}
