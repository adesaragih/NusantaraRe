package models

// Dampak penghapusan klaim - tiket 15.
//
// Untuk apa berkas ini: pengguna berhak tahu BERAPA baris yang akan ikut
// terhapus sebelum ia menekan Ya. Angka di peringatan itu dihitung dari data,
// bukan diperkirakan.
//
// Dibaca sesudah: pohonklaim.go.
//
// Istilah:
//   - kaskade : baris anak yang ikut terhapus bersama induknya.

import "fmt"

// DampakHapus adalah cacah baris yang akan ikut terhapus, per jenis.
//
// ⛔ Setiap tingkat punya medannya sendiri, bukan satu angka total. Satu angka
// menyembunyikan tingkat mana yang ternyata lebih besar dari dugaan - dan
// justru itulah yang membuat pengguna menekan Batal.
type DampakHapus struct {
	// Header adalah baris T_GENERAL_CLAIM itu sendiri - 0 atau 1.
	Header         int
	Peserta        int
	Adjustment     int
	Spreading      int
	SpreadingRetro int
	Dokumen        int
	// WorkClaim adalah baris tangga kerja klaim ini - 0 atau 1.
	WorkClaim int
	// BarisDatarWarisan adalah baris OS_AKSEPTASI_KLAIM_LIFE ber-CASEID sama.
	//
	// ⚠️ `[terbuka - work owner]` NASIB baris ini saat klaim dihapus belum
	// diputuskan siapa pun. Kode penghapusan yang ada (tiket 14) MENGHAPUSNYA;
	// tiket 15 menandainya sebagai pertanyaan yang belum dijawab. Angkanya
	// ditampilkan justru supaya keputusan yang belum diambil itu TERLIHAT oleh
	// yang menekan tombol, bukan tersembunyi di dalam kaskade.
	BarisDatarWarisan int
}

// Total menjumlahkan seluruh baris yang akan terhapus.
func (d DampakHapus) Total() int {
	return d.Header + d.Peserta + d.Adjustment + d.Spreading + d.SpreadingRetro +
		d.Dokumen + d.WorkClaim
}

// TotalTermasukWarisan menambahkan baris datar warisan ke jumlahnya.
//
// ⛔ Dipisah dari Total dengan sengaja. Nasib baris datar warisan saat klaim
// dihapus **belum diputuskan** (tiket 15 AC 11, `[terbuka - work owner]`);
// menjumlahkannya ke dalam "yang akan ikut terhapus" berarti menjawab
// pertanyaan itu diam-diam, lewat sebuah angka.
func (d DampakHapus) TotalTermasukWarisan() int {
	return d.Total() + d.BarisDatarWarisan
}

// Kosong menyatakan tidak ada apa pun yang akan terhapus.
func (d DampakHapus) Kosong() bool { return d.Total() == 0 }

// String menulis rincian dampak untuk pesan dan jejak audit.
func (d DampakHapus) String() string {
	return fmt.Sprintf(
		"header %d, peserta %d, adjustment %d, spreading %d, spreading retro %d, "+
			"dokumen %d, baris work %d, baris datar warisan %d",
		d.Header, d.Peserta, d.Adjustment, d.Spreading, d.SpreadingRetro,
		d.Dokumen, d.WorkClaim, d.BarisDatarWarisan)
}
