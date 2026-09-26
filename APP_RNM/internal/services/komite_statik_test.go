package services_test

// Penjaga statik tiket 10 - batas konteks Komite.
//
// Pemilik: tiket 10. Dibaca sesudah: komite.go.
//
// Dua aturan dijaga di sini, dan keduanya diminta tiket dalam bentuk test:
//
//  1. Roster dan keputusan per anggota BUKAN milik Claim Life. Menyalinnya
//     ke sini berarti dua sumber kebenaran atas satu keputusan.
//  2. Rujukan ke Komite memakai ID STABIL, bukan indeks posisi. Indeks posisi
//     berubah ketika baris di atasnya dihapus, dan rujukan yang berubah
//     sendiri menunjuk keputusan milik baris lain.

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// polaMilikKomite mencocokkan penyimpan keputusan komite per anggota.
//
// `[terverifikasi]` `Claim Life/Activity/CreateKMTLife_Act.xml` menulis
// `childPageKomite.KomiteList(<LAST>).KomiteAproval` dan `.KomiteEmail` -
// nilainya menyeberang, tetapi TEMPATNYA di konteks Komite Claim Life.
var polaMilikKomite = regexp.MustCompile(
	`KomiteAproval|KomiteComment|DateApprove|T_KOMITE_KOMITELIST`)

// polaIndeksPosisi mencocokkan indeks posisi Pega yang dipakai sebagai kunci.
//
// `[terverifikasi]` `CreateKMTLife_Act.xml` menaut lewat `.pxListSubscript`
// beserta `IndexAdjustment` / `IndexPremiumList`. Itu bentuk LAMA; penyimpangan
// sadarnya adalah menautkan lewat `KOMITE_ID` (AC 62, tiket 14).
var polaIndeksPosisi = regexp.MustCompile(
	`IndexAdjustment|IndexPremiumList|pxListSubscript`)

// berkasKomiteBolehMenyebut adalah tempat nama-nama itu boleh muncul sebagai
// PROSA - penjelasan mengapa sesuatu TIDAK dilakukan di sini.
//
// ⚠️ Daftar ini sengaja pendek dan beralasan. Daftar pengecualian yang tumbuh
// tanpa alasan adalah cara paling umum sebuah penjaga mati diam-diam.
// ⛔ BERKUNCI JALUR PENUH, bukan nama berkas. Ronde pertama memakai
// `filepath.Base`, sehingga berkas bernama `komite_statik_test.go` DI MANA PUN
// ikut dikecualikan - satu berkas baru bernama sama di paket lain sudah cukup
// untuk membungkam penjaga ini. Bypass itu dibangun dan dijalankan, dan ia
// hijau.
var berkasKomiteBolehMenyebut = map[string]string{
	"internal/services/komite_statik_test.go": "berkas penjaga ini sendiri - polanya harus tertulis",
}

// TestNolPenyimpanKeputusanKomiteDiKonteksIni menegakkan batas konteks.
//
// ⛔ Yang dijaga adalah KODE, bukan komentar: komentar dibuang lebih dulu,
// sebab menjelaskan mengapa sesuatu bukan milik kita justru menuntut
// menyebut namanya.
func TestNolPenyimpanKeputusanKomiteDiKonteksIni(t *testing.T) {
	diperiksa, dikecualikan := 0, 0
	// ⛔ Akar telusurnya AKAR MODUL, bukan `internal/`. Ronde pertama memakai
	// ".." - yang dari paket ini berarti `internal/` saja, sehingga `cmd/` dan
	// `pkg/` tidak pernah dibaca.
	err := filepath.Walk(filepath.Join("..", ".."), func(jalur string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(jalur, ".go") {
			return nil
		}
		// ⛔ Dihitung LEBIH DULU. Ronde pertama menaikkan pencacah sesudah
		// pengecualian, sehingga berkas yang dikecualikan tidak pernah terhitung
		// dan swa-periksa `diperiksa < 20` tidak mungkin menangkap pengecualian
		// yang terlalu lebar.
		diperiksa++
		rel := filepath.ToSlash(jalur)
		if i := strings.Index(rel, "internal/"); i >= 0 {
			rel = rel[i:]
		}
		if alasan, boleh := berkasKomiteBolehMenyebut[rel]; boleh {
			t.Logf("dikecualikan: %s (%s)", rel, alasan)
			dikecualikan++
			return nil
		}
		isi, err := os.ReadFile(jalur)
		if err != nil {
			return err
		}
		kode := buangKomentar(string(isi))
		if m := polaMilikKomite.FindString(kode); m != "" {
			t.Errorf("%s: menyebut %q di luar komentar; roster dan keputusan per "+
				"anggota milik konteks Komite Claim Life, bukan Claim Life",
				filepath.ToSlash(jalur), m)
		}
		if m := polaIndeksPosisi.FindString(kode); m != "" {
			t.Errorf("%s: memakai %q; rujukan ke Komite memakai KOMITE_ID - ID "+
				"stabil, bukan indeks posisi (AC 62)", filepath.ToSlash(jalur), m)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	// ⛔ Pembacanya sendiri diperiksa. Penjaga yang tidak membaca berkas apa
	// pun hijau selamanya, dan hijau selamanya tidak dapat dibedakan dari
	// aman selamanya.
	if diperiksa < 40 {
		t.Fatalf("hanya %d berkas terbaca; pembacanya yang rusak, bukan kodenya",
			diperiksa)
	}
	// ⛔ Pengecualian yang tumbuh adalah cara penjaga ini mati diam-diam.
	// Jumlahnya dikunci ke jumlah baris petanya.
	if dikecualikan != len(berkasKomiteBolehMenyebut) {
		t.Errorf("%d berkas dikecualikan, sedangkan petanya memuat %d; "+
			"pengecualian yang tidak terpakai atau berlipat sama-sama cacat",
			dikecualikan, len(berkasKomiteBolehMenyebut))
	}
}

// TestPenyerahanTidakMemutuskanStatusBaris - tahap dan status tetap terpisah.
//
// ⛔ Menyerahkan ke Komite TIDAK mengubah status baris. Barisnya tetap
// Outstanding sampai Komite memutuskan, dan keputusan itu masuk lewat tiket 11.
// Bila penyerahan ikut menulis status, baris akan tampak sudah diputus padahal
// belum ada yang memutuskannya.
func TestPenyerahanTidakMemutuskanStatusBaris(t *testing.T) {
	isi, err := os.ReadFile(filepath.Join("..", "services", "komite.go"))
	if err != nil {
		t.Fatal(err)
	}
	kode := buangKomentar(string(isi))
	// ⚠️ Yang dicari PENULISAN, bukan pembacaan. Ronde pertama juga melarang
	// `KodeStatus:` - dan langsung menangkap `MuatanKomite{KodeStatus: ...}`,
	// yang justru AC 24: status baris SAAT penyerahan ikut menyeberang.
	// Membaca status untuk dilaporkan bukan memutuskannya.
	for _, terlarang := range []string{
		"PerbaruiStatusBaris(", "CerminkanHeader(", "Transisi(",
		".KodeStatus =", ".Ubah(", ".Tolak(",
	} {
		if strings.Contains(kode, terlarang) {
			t.Errorf("komite.go memuat %q; penyerahan memindahkan KEPUTUSAN ke "+
				"Komite, ia tidak mengambil keputusan itu sendiri (ADR-U-0011)",
				terlarang)
		}
	}
}
