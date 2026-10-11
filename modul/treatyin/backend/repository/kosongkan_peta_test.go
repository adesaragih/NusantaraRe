package repository

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// ⛔ MENGAPA PENJAGA INI ADA
//
// `alat/kosongkan-tab-treatyin.sql` ditulis ketika tabel pendaratan baru
// SEMBILAN. Migrasi 437, 438, 439, dan 446 menambah 21 lagi, dan berkas itu
// tidak ikut tumbuh - selama berbulan-bulan ia mengosongkan sembilan tabel
// dan meninggalkan 21 TERISI, tanpa satu pun galat.
//
// ⚠️ Itu jenis kegagalan yang paling mahal: ia DIAM. Pemuatan berikutnya
// menumpuk di atas sisa yang tidak terbuang, dan yang membacanya mengira
// datanya sah.
//
// Penjaga ini mengikat berkas itu ke `PetaPendaratan` - sumber kebenaran
// yang sama yang dipakai pemuat.

func bacaSkripKosongkan(t *testing.T) string {
	t.Helper()
	p := filepath.Join("..", "..", "alat", "kosongkan-tab-treatyin.sql")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("baca %s: %v", p, err)
	}
	return string(b)
}

var polaHapus = regexp.MustCompile(`(?i)DELETE\s+FROM\s+\{skema\}\.([A-Z0-9_]+)`)

func tabelDiSkrip(t *testing.T) []string {
	t.Helper()
	var out []string
	for _, m := range polaHapus.FindAllStringSubmatch(bacaSkripKosongkan(t), -1) {
		out = append(out, strings.ToUpper(m[1]))
	}
	return out
}

// ⛔ SETIAP tabel peta harus ada di skrip. Tabel yang lahir di peta tanpa
// ditambahkan di sini akan tertinggal terisi.
func TestKosongkanMemuatSeluruhPeta(t *testing.T) {
	ada := map[string]bool{}
	for _, n := range tabelDiSkrip(t) {
		ada[n] = true
	}
	var hilang []string
	for _, p := range PetaPendaratan {
		if !ada[strings.ToUpper(p.Tabel)] {
			hilang = append(hilang, p.Tabel)
		}
	}
	if len(hilang) > 0 {
		t.Fatalf("tabel peta TIDAK dikosongkan: %v", hilang)
	}
}

// ⛔ Dan sebaliknya: skrip tidak boleh menyebut tabel di luar peta. Yang
// paling berbahaya tabel WARISAN - `TREATYEXCHANGEYEARLY` dipakai bersama
// tiga modul lain, dan mengosongkannya merusak modul yang tidak ada
// hubungannya dengan Treaty In.
func TestKosongkanNolMenyentuhTabelLuar(t *testing.T) {
	sah := map[string]bool{}
	for _, p := range PetaPendaratan {
		sah[strings.ToUpper(p.Tabel)] = true
	}
	var asing []string
	for _, n := range tabelDiSkrip(t) {
		if !sah[n] {
			asing = append(asing, n)
		}
	}
	if len(asing) > 0 {
		t.Fatalf("skrip menyentuh tabel di luar peta pendaratan: %v", asing)
	}
}

// ⛔ ANAK sebelum INDUK. Urutan muat peta induk-dulu; urutan hapus harus
// kebalikannya, atau `DELETE` induk akan ditolak kunci asing.
func TestUrutanHapusAnakSebelumInduk(t *testing.T) {
	urut := map[string]int{}
	for i, n := range tabelDiSkrip(t) {
		if _, ada := urut[n]; !ada {
			urut[n] = i
		}
	}
	for _, p := range PetaPendaratan {
		if p.Induk == "" {
			continue
		}
		anak, induk := strings.ToUpper(p.Tabel), strings.ToUpper(p.Induk)
		ia, ok1 := urut[anak]
		ii, ok2 := urut[induk]
		if !ok1 || !ok2 {
			continue // ditangkap uji kelengkapan di atas
		}
		if ia > ii {
			t.Errorf("%s (anak) dihapus SESUDAH %s (induk) — kunci asing akan menolak", anak, induk)
		}
	}
}

// ⛔ `TRUNCATE` mengikat transaksinya (commit implisit), jadi pemuat yang
// gagal di tengah tidak dapat memulihkan barisnya. `DROP` membongkar
// struktur. Keduanya terlarang di berkas ini.
func TestKosongkanHanyaMemakaiDelete(t *testing.T) {
	isi := strings.ToUpper(bacaSkripKosongkan(t))
	for _, larangan := range []string{"TRUNCATE", "DROP ", "COMMIT"} {
		// Komentar boleh MENYEBUTNYA - yang dilarang pernyataannya.
		for _, baris := range strings.Split(isi, "\n") {
			b := strings.TrimSpace(baris)
			if strings.HasPrefix(b, "--") {
				continue
			}
			if strings.Contains(b, larangan) {
				t.Errorf("pernyataan terlarang %q: %s", larangan, baris)
			}
		}
	}
}
