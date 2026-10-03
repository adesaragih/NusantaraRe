package repository

// Penomoran yang selamat melewati cutover - tiket 13.
//
// Untuk apa berkas ini: memastikan penghitung nomor yang hidup sesudah migrasi
// TIDAK BERADA DI BAWAH nomor tertinggi yang sudah terpakai, per kombinasi
// kunci.
//
// ⚠️ "Tidak di bawah", bukan "di atas" - dan selisihnya penting. `[data DBA]`
// `SUMBER-PENOMORAN-DBA.md` baris 72-86: prosedurnya membaca `no_seq`, lalu
// `+1`, lalu `UPDATE`. Jadi `NO_SEQ` adalah nomor TERAKHIR YANG SUDAH
// DITERBITKAN, bukan nomor berikutnya - penghitung yang SAMA dengan tanda air
// sudah benar. Ronde pertama menulis "di atas" di prosa sementara kodenya
// menerima kesamaan; prosa yang lebih ketat daripada kodenya mengundang orang
// "memperbaiki" kode yang sudah benar.
//
// Dibaca sesudah: migrasidata.go.
//
// ⛔ Ia TIDAK menjalankan migrasi dan TIDAK menyentuh tabel produksi.
// Menjalankan migrasi data menuntut persetujuan manusia; yang ada di sini
// aturannya, dapat diuji tanpa Oracle.
//
// Istilah:
//   - tanda air : nomor tertinggi yang sudah terpakai pada satu kombinasi.

import (
	"errors"
	"fmt"
	"sort"
)

// ErrSequenceMundur - penghitung hidup berada di bawah nomor yang terpakai.
var ErrSequenceMundur = errors.New(
	"repository: penghitung nomor berada di bawah nomor yang sudah terpakai")

// KunciSequence adalah kunci primer komposit `GENERATE_SEQUENCE_NUMBER`.
//
// ⛔ RALAT LABEL 27-09-2026. Ronde pertama menulis `[terverifikasi]` OQ-002.
// Keduanya keliru:
//
//   - Labelnya DINAIKKAN. Sumbernya
//     `modul/claimlife/docs/SUMBER-PENOMORAN-DBA.md` baris 116, 125, 136-137,
//     dan berkas itu melabeli dirinya sendiri di baris 3: `[data DBA - dibaca
//     sendiri 26 September 2026 dari katalog instance pengembangan, BELUM
//     DIKONFIRMASI DBA]`. Menaikkannya menjadi `[terverifikasi]` melanggar
//     disiplin label: `[terverifikasi]` hanya untuk yang dibaca sesi ini dari
//     korpus, dengan path dan nomor baris.
//   - Nomor OQ-nya salah. OQ-002 tentang BODY stored procedure; fakta DDL
//     seperti tipe kolom milik OQ-001, yang masih terbuka.
//
// `[data DBA - belum dikonfirmasi DBA]` `(CLASS, JENIS, TAHUN)` dengan
// `TAHUN VARCHAR2(5)` dan `NO_SEQ NUMBER`.
//
// ⛔ Ketiganya TEKS, termasuk tahun. `TAHUN VARCHAR2(5)` - mengubahnya menjadi
// bilangan membuang bentuk aslinya dan membuat dua tahun yang tersimpan
// berbeda tampak sama (ADR-U-0022).
type KunciSequence struct {
	Class string
	Jenis string
	Tahun string
}

// String menulis kunci sebagai satu kata yang terbaca di pesan galat.
func (k KunciSequence) String() string {
	return fmt.Sprintf("(%s, %s, %s)", k.Class, k.Jenis, k.Tahun)
}

// NomorTerpakai adalah satu nomor yang sudah dipegang data lama.
type NomorTerpakai struct {
	Kunci KunciSequence
	// Urut adalah `NO_SEQ`.
	//
	// ⚠️ `int` dipilih sadar, dan batasnya dinyatakan: `NO_SEQ NUMBER` di
	// Oracle tak berbatas, tetapi keluaran prosedurnya `LPAD(no_seq, 5, '0')`
	// `[data DBA]` `SUMBER-PENOMORAN-DBA.md` baris 122 - jadi bentuk
	// nomornya sendiri pecah di 99999, jauh sebelum `int` pecah.
	Urut int
}

// TandaAirSequence mencari nomor tertinggi per kombinasi kunci.
//
// ⛔ PER KOMBINASI, bukan satu penghitung global. Satu penghitung global akan
// membuat nomor MELOMPAT pada kombinasi yang jarang dipakai - ia mewarisi
// tinggi kombinasi lain - dan MENGULANG pada kombinasi yang ramai. Keduanya
// baru terlihat berbulan-bulan sesudah cutover, ketika memperbaikinya berarti
// menyentuh nomor yang sudah dicetak di dokumen.
func TandaAirSequence(terpakai []NomorTerpakai) map[KunciSequence]int {
	air := make(map[KunciSequence]int, len(terpakai))
	for _, n := range terpakai {
		if lama, ada := air[n.Kunci]; !ada || n.Urut > lama {
			air[n.Kunci] = n.Urut
		}
	}
	return air
}

// PeriksaTandaAir menolak penghitung yang belum melewati nomor terpakai.
//
// ⛔ Kunci yang ADA di tanda air tetapi HILANG dari penghitung juga ditolak.
// Penghitung yang tidak mengenal sebuah kombinasi akan mulai dari nol, dan
// nomor pertamanya sesudah cutover akan bertabrakan dengan nomor yang sudah
// dipegang klaim lama - kegagalan yang persis sama, hanya lebih senyap.
//
// ⚠️ Galatnya menyebut KUNCI mana yang bermasalah. Galat migrasi yang hanya
// berkata "penomoran salah" memaksa orang menebak di antara ribuan kombinasi.
func PeriksaTandaAir(penghitung, tandaAir map[KunciSequence]int) error {
	var kurang []string
	for kunci, air := range tandaAir {
		nilai, ada := penghitung[kunci]
		switch {
		case !ada:
			kurang = append(kurang, fmt.Sprintf("%s tidak dikenal penghitung "+
				"(terpakai sampai %d)", kunci, air))
		case nilai < air:
			kurang = append(kurang, fmt.Sprintf("%s penghitung %d < terpakai %d",
				kunci, nilai, air))
		}
	}
	if len(kurang) == 0 {
		return nil
	}
	// Diurutkan supaya pesannya sama di setiap jalannya - laporan migrasi yang
	// berubah urutannya tiap dijalankan tidak dapat dibandingkan.
	sort.Strings(kurang)
	// ⛔ Dibatasi. Migrasi yang salah seluruhnya akan menghasilkan ribuan
	// kombinasi, dan galat sepanjang itu tidak terbaca siapa pun - ia juga
	// berakhir di log. Cacah penuhnya tetap disebut.
	const tampil = 10
	dipotong := kurang
	ekor := ""
	if len(dipotong) > tampil {
		dipotong = dipotong[:tampil]
		ekor = fmt.Sprintf(" (dan %d lainnya)", len(kurang)-tampil)
	}
	return fmt.Errorf("%w: %d kombinasi bermasalah: %v%s",
		ErrSequenceMundur, len(kurang), dipotong, ekor)
}
