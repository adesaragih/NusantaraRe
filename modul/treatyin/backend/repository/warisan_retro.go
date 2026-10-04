package repository

// PENANDA — kontrak yang satu cabang datanya SENGAJA ditinggalkan.
//
// ⛔ Berkas ini ada supaya sebuah lubang tidak menjadi diam. Keputusan
// pemilik proses 4 Oktober 2026 menunda tab Retro
// (`docs/KEPUTUSAN-PENYELARASAN-REPO.md` §14): dua kontrak dari 1.854 tidak
// cukup untuk merancang tiga tabel bersarang. Keputusan itu benar, dan ia
// punya ongkos — dua kontrak akan terlihat "selesai dipindahkan" pada
// rekonsiliasi tiket `44` padahal `RetroList`-nya tidak ikut.
//
// ⚠️ Penandaannya WAJIB dapat disapu mesin, dan tidak boleh bergantung pada
// ingatan siapa pun. Tiga hal mewujudkannya, dan ketiganya diperlukan:
//
//	1. daftar di bawah — dapat di-`grep`, dan menyebut sebabnya di tempat;
//	2. `TestKontrakRetroTertundaMasihDuaItu` (`-tags db`) MENGUKUR ULANG dari
//	   Oracle dan merah begitu kontrak ketiga muncul atau salah satu hilang;
//	3. pemuat mencetaknya pada tiap `-cocokkan`, sehingga ia lewat di depan
//	   mata orang yang sedang merekonsiliasi — bukan hanya tersimpan.
//
// ⛔ JANGAN menghapus daftar ini ketika tab Retro dibangun. Yang dihapus
// nanti adalah seluruh berkasnya, bersama ujinya, dalam ronde yang sama
// dengan tabelnya — supaya tidak ada tenggang waktu saat penandanya hilang
// tetapi datanya belum pindah.

// KontrakRetroTertunda adalah kontrak yang punya `RetroList` berisi.
//
// Terukur 4 Oktober 2026 atas SELURUH 1.854 dokumen, diurai utuh sebagai
// JSON: hanya kedua pengenal ini, keduanya `NonProportional`, masing-masing
// 2 elemen — 4 elemen seluruhnya.
var KontrakRetroTertunda = []string{"1000493", "1000755"}

// LarikRetroTertunda adalah kunci larik yang ditinggalkan di dokumen.
const LarikRetroTertunda = "RetroList"

// AlasanRetroTertunda ikut tercetak bersama penandanya, supaya yang membaca
// peringatannya tidak perlu mencari dokumen untuk tahu sebabnya.
const AlasanRetroTertunda = "tab Retro ditunda (KEPUTUSAN §14): 2 kontrak dari 1.854 tidak cukup " +
	"merancang 3 tabel bersarang; 8 dari 17 medannya turunan (INV-58). `RetroList` kontrak di " +
	"bawah TIDAK ikut dimuat, jadi rekonsiliasi tiket 44 tidak boleh menghitungnya selesai."
