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
// `[terverifikasi]` `Claim Life/Activity/CreateKMTLife_Act.xml` langkah 4
// b1188 menyalin `childPageKomite.KomiteList = .KomiteList` - tangga yang
// `GetListKomiteLife` 6.1 susun, lengkap dengan `KomiteAproval` (b1197) dan
// `KomiteEmail` (b1217). Nilainya menyeberang, tetapi TEMPATNYA di konteks
// Komite Claim Life.
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
	"modul/claimlife/backend/services/komite_statik_test.go": "berkas penjaga ini sendiri - polanya harus tertulis",
	// ⛔ DIPERSEMPIT, bukan dilonggarkan - A2, 27-09-2026.
	//
	// Premis penjaga ini di tiket 10 - *"roster dan keputusan per anggota
	// bukan milik konteks ini"* - BENAR sampai butir **af** disahkan. Dua hal
	// membatalkannya:
	//
	//   1. Brief lanjutan 4 §1 menempatkan tabelnya di rangkaian migrasi INI.
	//   2. `[terverifikasi]` `CreateKMTLife_Act.xml` - activity milik **Claim
	//      Life** - yang menulis tangganya (langkah 4 b1188, dari
	//      `GetListKomiteLife` 6.1 b1151-b1237). Penulisnya
	//      memang konteks ini; yang MEMUTUSKAN barulah Komite.
	//
	// Yang tetap dijaga: nol berkas LAIN menyentuhnya, dan nol indeks posisi
	// dipakai sebagai kunci rujukan di mana pun.
	"modul/claimlife/backend/repository/kasuskomite.go": "butir af - penulis tangga, meniru CreateKMTLife_Act (activity Claim Life)",
	// ⛔ Penjaga BENTUK SKEMA, bukan jalur simpan. Ia menyebut kedua nama
	// tabel justru untuk memeriksa bahwa KEDUA dokumen STRUKTUR - Claim Life
	// dan Komite Claim Life - sepakat atas bentuknya
	// (TestDokumenSTRUKTURSepakatAtasTabelBersama). Melarangnya menyebut nama
	// itu berarti melarang satu-satunya uji yang menjaga batas kedua konteks
	// tetap satu bentuk.
	//
	// ⚠️ Pengecualian ini TIDAK melonggarkan aturannya: berkas itu berkas
	// UJI, ia tidak menulis satu baris pun ke basis data.
	// Refactor bentuk B (30-09-2026): penjaga lintas modul, kini di inti/backend/penjaga.
	"inti/backend/penjaga/strukturkolom_test.go": "penjaga bentuk skema lintas dokumen; berkas uji, nol jalur simpan",
	// ⛔ MODUL KOMITE CLAIM LIFE sendiri - giliran 10, tiket 01 Komite.
	// Premis penjaga ini adalah batas KONTEKS: Claim Life tidak memutuskan
	// atas nama Komite. Berkas di bawah BUKAN Claim Life - brief
	// GILIRAN-3-KOMITE §0 memberi modul Komite berkas `komite_*`. Disebut
	// satu per satu, JALUR PENUH; bukan awalan `komite_`, sebab awalan itu
	// ikut membungkam `komite_test.go`/`komite_db_test.go` milik Claim Life.
	"modul/komiteclaimlife/backend/repository/komite_inbox.go":     "modul Komite Claim Life - pembaca inbox + tangga (tiket 01 Komite)",
	"modul/komiteclaimlife/backend/repository/komite_keputusan.go": "modul Komite Claim Life - PENULIS keputusan tingkat (tiket 02 Komite); inilah yang penjaga ini maksud dengan 'milik konteks Komite'",
	// ⛔ MODUL CLAIM PROP - keputusan work owner 07-10-2026 (opsi "b"): penyerahan klaim treaty ke komite
	// (`AddKomiteTreatyChild_ACT`, activity Claim Prop) menulis tangga awal, seperti `kasuskomite.go` di atas, dan grid
	// "Committe Accept Status" membaca keputusannya. Kode batas komite Claim Prop dikumpulkan di DUA berkas ini saja;
	// keputusan tingkat tetap ditulis konteks Komite.
	"modul/claimprop/backend/models/komite.go":     "modul Claim Prop - nama properti keputusan anggota (grid Committe Accept Status); keputusan work owner 07-10-2026",
	"modul/claimprop/backend/repository/komite.go": "modul Claim Prop - penulis tangga AddKomiteTreatyChild_ACT dan pembacanya; keputusan work owner 07-10-2026",
	// ⛔ MODUL KOMITE CLAIM PROP - izin work owner 08-10-2026: penulis keputusan tingkat (`KomitePostAdjustment` S6 /
	// S26.1) dan pembaca tangga / daftar kerja (`KomiteRouter` S6.1). SATU berkas; nama properti keputusan tidak dipakai.
	"modul/komiteclaimprop/backend/repository/tangga.go": "modul Komite Claim Prop - penulis keputusan tangga dan daftar kerja KomiteRouter; izin work owner 08-10-2026",
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
	err := filepath.Walk(filepath.Join("..", "..", "..", ".."), func(jalur string, info os.FileInfo, err error) error {
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
		// Jalur relatif akar APP_RNM - `internal/...`, `inti/...`, `modul/...`.
		rel := filepath.ToSlash(jalur)
		for strings.HasPrefix(rel, "../") {
			rel = strings.TrimPrefix(rel, "../")
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
		"PerbaruiStatusBaris(", "TandaiBarisOutstanding(", "CerminkanHeader(", "Transisi(",
		".KodeStatus =", ".Ubah(", ".Tolak(",
	} {
		if strings.Contains(kode, terlarang) {
			t.Errorf("komite.go memuat %q; penyerahan memindahkan KEPUTUSAN ke "+
				"Komite, ia tidak mengambil keputusan itu sendiri (ADR-U-0011)",
				terlarang)
		}
	}
}
