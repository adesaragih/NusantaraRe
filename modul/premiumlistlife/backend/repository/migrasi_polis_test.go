package repository

import (
	"regexp"
	"strings"
	"testing"
)

// milikPremiumList menjawab apakah berkas migrasi itu milik PremiumList Life.
//
// Rentangnya 050-099 sejak struktur tim satu folder per modul (R2, MODUL.md);
// semula 050-079 (PROMPT-EKSEKUSI-HULU-HILIR.md §4). Lahir 28-09-2026
// ketika modul ketiga (Treaty Contract Out, 300-319) menambah migrasi dan
// "bukan Claim Life" tidak lagi berarti "PremiumList".
func milikPremiumList(nama string) bool {
	return nama >= "050_" && nama < "100_"
}

// Seluruh FK pohon polis BERKASKADE - tiket 00 PremiumList Life AC 46.
//
// ⛔ Kebijakan yang BERBEDA dari Claim Life, dan sengaja. Pohon polis empat
// tingkat (polis -> peserta -> spreading -> spreading retro) dan hapus polis
// harus membersihkan seluruh turunannya; menangani kaskade di Go untuk pohon
// sedalam itu berarti empat perjalanan pulang-pergi dan satu kesempatan
// gagal di tengah.
//
// ⚠️ Migrasi 050 dan 051 TIDAK punya FK sama sekali - T_WORK_POLIS berdiri
// sendiri, dan T_PREMIUM_LIST berbagi PK dengannya tanpa constraint. Uji ini
// menuntut keduanya TANPA kaskade, supaya constraint yang diam-diam
// ditambahkan di antara keduanya berbunyi.
func TestSeluruhFKPohonPolisBerkaskade(t *testing.T) {
	tanpaFK := map[string]bool{"050_": true, "051_": true}
	diperiksa := 0
	for nama, teks := range seluruhSQL(t, false) {
		// ⛔ DIPERSEMPIT ke rentang PremiumList Life 28-09-2026 (sesi Treaty
		// Contract Out). Sebelumnya "bukan Claim Life" berarti "PremiumList",
		// sebab hanya dua modul yang bermigrasi. Sejak migrasi 300-319 ada,
		// tabel yang menggantung pada KUNCI GABUNGAN tanpa FK (spec Treaty
		// Contract Out §2) akan dituduh "tidak punya FOREIGN KEY". Penjaga
		// yang menuduh hal yang benar akan dilonggarkan orang; ia karena itu
		// dipersempit ke rentang modul yang kebijakannya ia jaga - tidak
		// dilonggarkan. Kebijakan FK Treaty Contract Out dijaga
		// tco_migrasi_test.go.
		if !milikPremiumList(nama) || strings.Contains(nama, "_down") {
			continue
		}
		isi := strings.ToUpper(teks)
		// Penjaga ini tentang TABEL anak. Migrasi yang tidak membuat tabel -
		// 057 hanya menambah sequence dan satu kolom (butir bn) - tidak punya
		// induk untuk ditunjuk.
		if !strings.Contains(isi, "CREATE TABLE") {
			continue
		}
		diperiksa++
		punyaFK := strings.Contains(isi, "FOREIGN KEY")
		berkaskade := strings.Contains(isi, "ON DELETE CASCADE")

		bebas := false
		for awalan := range tanpaFK {
			if strings.HasPrefix(nama, awalan) {
				bebas = true
			}
		}
		if bebas {
			if punyaFK {
				t.Errorf("%s punya FOREIGN KEY; 050 dan 051 seharusnya tanpa FK "+
					"- hubungan keduanya SHARED PK tanpa kolom penyambung", nama)
			}
			continue
		}
		if !punyaFK {
			t.Errorf("%s tidak punya FOREIGN KEY; seluruh tabel anak pohon polis "+
				"menunjuk induknya", nama)
			continue
		}
		if !berkaskade {
			t.Errorf("%s punya FK TANPA ON DELETE CASCADE (tiket 00 AC 46)", nama)
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol migrasi PremiumList ditelusuri - penjaga ini tidak menjaga apa pun")
	}
}

// Setiap FK pohon polis PUNYA INDEX - tiket 00 AC 49.
//
// ⛔ Skala jutaan baris. FK tanpa index membuat setiap hapus induk memindai
// seluruh tabel anak, dan pada T_PREMIUM_LIST_DETAIL itu berarti memindai
// jutaan baris untuk menghapus satu polis.
func TestSetiapFKPohonPolisBerindex(t *testing.T) {
	diperiksa := 0
	polaFK := regexp.MustCompile(`FOREIGN KEY \(([A-Z_]+)\)`)
	for nama, teks := range seluruhSQL(t, false) {
		// Dipersempit ke rentang PremiumList Life (lihat penjaga di atas);
		// index FK Treaty Contract Out dijaga tco_migrasi_test.go.
		if !milikPremiumList(nama) || strings.Contains(nama, "_down") {
			continue
		}
		isi := strings.ToUpper(teks)
		for _, m := range polaFK.FindAllStringSubmatch(isi, -1) {
			kolom := m[1]
			diperiksa++
			// Index-nya harus ada DAN menyebut kolom FK itu, bukan kolom lain.
			// Index pada kolom lain menenangkan tanpa menjaga: hapus induk
			// tetap memindai seluruh tabel anak.
			if !strings.Contains(isi, "CREATE INDEX") {
				t.Errorf("%s: FK pada %s, dan berkas itu tidak membuat index apa pun",
					nama, kolom)
				continue
			}
			if !strings.Contains(isi, "("+kolom+")"+"\n") {
				t.Errorf("%s: ada CREATE INDEX tetapi tidak pada kolom FK %s", nama, kolom)
			}
		}
	}
	if diperiksa == 0 {
		t.Fatal("nol FK pohon polis ditemukan - penjaga ini tidak menjaga apa pun")
	}
}
