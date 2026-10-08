package repository

// AMBANG PEMBALIKAN §17 — Retro tidak dibangun karena JARANG.
//
// ⛔ Berkas ini pernah berjudul "penanda pekerjaan tertunda" dan artinya
// BERUBAH 4 Oktober 2026. Keputusan §14 menunda Retro; keputusan §17
// menggantikannya dengan **tidak dibangun, karena jarang**. Tidak ada lagi
// pekerjaan yang menunggu — yang ada **ambang yang diawasi**.
//
// ⛔ DAN SEBABNYA BUKAN KODE MATI. Ronde sebelumnya hendak mencabut Retro
// atas dasar penjaga `1=2` di `Section/ShareRetro.xml`, dan verifikasi
// membatalkan dasar itu: ketiga penjaga membungkus sebuah tombol, satu blok
// `BAR`, dan satu tombol kepala — **nol yang membungkus tabnya**. Tab
// non-prop berdiri tanpa syarat tampil sama sekali; tab prop bersyarat
// `TreatyIn.IsMultipleRetro`, sebuah syarat DATA. `FlowAction/ShareRetro.xml`
// berbunyi `pyRuleAvailable = Yes`. **Retro HIDUP di Pega.**
//
// ⚠️ §17 berdiri di atas DUA angka, jadi DUA penjaga mengukurnya ulang dari
// Oracle — `TestAmbangRetroJarangMasihDuaKontrak` dan
// `TestAmbangRetroMultipleMasihLima`. Satu angka yang dijaga dan satu yang
// dihafal adalah keputusan yang setengahnya dapat basi tanpa suara.
//
// ⛔ JANGAN menghapus berkas ini ketika Retro kelak dibangun. Yang dihapus
// nanti adalah seluruh berkasnya bersama kedua ujinya, dalam ronde yang
// sama dengan tabelnya.

// KontrakRetroJarang adalah kontrak yang punya `RetroList` berisi.
//
// Terukur 4 Oktober 2026 atas SELURUH 1.854 dokumen, diurai utuh sebagai
// JSON: hanya kedua pengenal ini, keduanya `NonProportional`, masing-masing
// 2 elemen — 4 elemen seluruhnya.
//
// ⛔ AMBANG PEMBALIKAN: kontrak KETIGA. Lihat §17.
var KontrakRetroJarang = []string{"1000493", "1000755"}

// LarikRetroJarang adalah kunci larik yang tidak dibangun.
const LarikRetroJarang = "RetroList"

// KunciRetroMultiple adalah penjaga tampil tab Retro cabang proporsional.
//
// ⚠️ Ia syarat DATA, bukan penjaga mati — dan itu pokok §17.
const KunciRetroMultiple = "IsMultipleRetro"

// AmbangRetroMultiple adalah cacah dokumen ber-`IsMultipleRetro = "true"`
// pada 4 Oktober 2026.
//
// ⛔ AMBANG PEMBALIKAN: kontrak KEENAM. Lihat §17.
const AmbangRetroMultiple = 5

// AlasanRetroJarang ikut tercetak bersama ambangnya, supaya yang membacanya
// tidak perlu mencari dokumen untuk tahu sebabnya.
const AlasanRetroJarang = "tab Retro TIDAK dibangun karena JARANG (KEPUTUSAN §17), bukan karena " +
	"kode mati - `pyRuleAvailable = Yes` dan nol penjaga `1=2` membungkus tabnya. " +
	"`IsMultipleRetro` true pada 5 dari 1.854; `RetroList` berisi pada 2 kontrak di bawah. " +
	"Ambang pembalikan: kontrak keenam ber-IsMultipleRetro, atau ketiga ber-RetroList."
