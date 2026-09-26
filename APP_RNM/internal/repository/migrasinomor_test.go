package repository

// Penomoran sesudah migrasi - TANPA Oracle.
//
// Pemilik: tiket 13. Dibaca sesudah: migrasinomor.go.

import (
	"errors"
	"strings"
	"testing"
)

// TestTandaAirDihitungPerKombinasiKunci - AC tiket 13.
//
// ⛔ `NO_SEQ` terjaga per kombinasi `(CLASS, JENIS, TAHUN)`, BUKAN sebagai
// penghitung global.
//
// `[data DBA - belum dikonfirmasi DBA]` `SUMBER-PENOMORAN-DBA.md` baris 116,
// 125, 136-137: `GENERATE_SEQUENCE_NUMBER` berkunci primer komposit
// `(CLASS, JENIS, TAHUN)` dengan `NO_SEQ NUMBER`. Ronde pertama menulis
// `[terverifikasi]` OQ-002 di sini - label dinaikkan, dan nomor OQ-nya pun
// salah; lihat ralat di `migrasinomor.go`.
//
// Satu penghitung global akan membuat nomor MELOMPAT di kombinasi yang jarang
// dipakai dan MENGULANG di kombinasi yang ramai - dua kegagalan yang AC ini
// larang, dan keduanya baru terlihat berbulan-bulan sesudah cutover.
func TestTandaAirDihitungPerKombinasiKunci(t *testing.T) {
	terpakai := []NomorTerpakai{
		{KunciSequence{"CLM", "A", "2026"}, 7},
		{KunciSequence{"CLM", "A", "2026"}, 3},
		{KunciSequence{"CLM", "A", "2025"}, 91},
		{KunciSequence{"CLM", "AR", "2026"}, 2},
		{KunciSequence{"KMT", "A", "2026"}, 40},
	}
	air := TandaAirSequence(terpakai)

	if len(air) != 4 {
		t.Fatalf("kombinasi kunci = %d, mau 4; kunci digabung atau dipecah salah", len(air))
	}
	for kunci, mau := range map[KunciSequence]int{
		{"CLM", "A", "2026"}:  7,
		{"CLM", "A", "2025"}:  91,
		{"CLM", "AR", "2026"}: 2,
		{"KMT", "A", "2026"}:  40,
	} {
		if got := air[kunci]; got != mau {
			t.Errorf("tanda air %+v = %d, mau %d", kunci, got, mau)
		}
	}
	// ⛔ Bukan penghitung global: tahun berbeda TIDAK saling menindih. Bila
	// ia global, kombinasi 2026 akan mewarisi 91 dari 2025 dan seluruh nomor
	// 2026 melompat.
	if air[KunciSequence{"CLM", "A", "2026"}] == air[KunciSequence{"CLM", "A", "2025"}] {
		t.Error("dua tahun berbagi tanda air; penghitungnya global, bukan per kunci")
	}
}

// TestTandaAirKosongBukanNol - kombinasi yang belum pernah dipakai.
func TestTandaAirKosongBukanNol(t *testing.T) {
	air := TandaAirSequence(nil)
	if len(air) != 0 {
		t.Errorf("tanda air dari nol nomor berisi %d kunci", len(air))
	}
	// ⚠️ Kunci yang tidak ada berbeda dari kunci bernilai 0. Peta Go
	// mengembalikan 0 untuk keduanya, jadi pemanggil WAJIB memeriksa
	// keberadaannya - dan `PeriksaTandaAir` melakukannya.
	if _, ada := air[KunciSequence{"CLM", "A", "2026"}]; ada {
		t.Error("kunci yang belum pernah dipakai dilaporkan ada")
	}
}

// TestPenomoranTidakMundurSesudahMigrasi - AC 42 spec.
//
// Sesudah migrasi, penghitung yang hidup wajib TIDAK BERADA DI BAWAH nomor
// tertinggi yang sudah terpakai; bila ia di bawahnya, nomor berikutnya
// MENGULANG nomor yang sudah dipegang klaim lain.
//
// ⚠️ SAMA DENGANNYA sudah benar - `NO_SEQ` adalah nomor terakhir yang
// DITERBITKAN, bukan nomor berikutnya.
func TestPenomoranTidakMundurSesudahMigrasi(t *testing.T) {
	air := map[KunciSequence]int{
		{"CLM", "A", "2026"}: 7,
		{"CLM", "A", "2025"}: 91,
	}
	// Penghitung hidup di atas atau sama dengan tanda air: aman.
	if err := PeriksaTandaAir(map[KunciSequence]int{
		{"CLM", "A", "2026"}: 7,
		{"CLM", "A", "2025"}: 100,
	}, air); err != nil {
		t.Errorf("penghitung yang cukup ditolak: %v", err)
	}
	// ⛔ Satu kunci mundur sudah cukup untuk menolak seluruhnya.
	err := PeriksaTandaAir(map[KunciSequence]int{
		{"CLM", "A", "2026"}: 6,
		{"CLM", "A", "2025"}: 100,
	}, air)
	if !errors.Is(err, ErrSequenceMundur) {
		t.Fatalf("galat = %v, mau ErrSequenceMundur", err)
	}
	// Pesannya menyebut KUNCI mana yang bermasalah - galat migrasi yang tidak
	// menyebut kuncinya memaksa orang menebak di antara ribuan kombinasi.
	if got := err.Error(); !strings.Contains(got, "CLM") || !strings.Contains(got, "2026") {
		t.Errorf("pesan tidak menyebut kuncinya: %q", got)
	}
	// ⛔ Kunci yang ADA di tanda air tetapi HILANG dari penghitung juga
	// ditolak: penghitung yang tidak tahu kombinasi itu akan mulai dari nol.
	if err := PeriksaTandaAir(map[KunciSequence]int{
		{"CLM", "A", "2026"}: 9,
	}, air); !errors.Is(err, ErrSequenceMundur) {
		t.Errorf("kunci hilang dari penghitung: galat = %v, mau ErrSequenceMundur", err)
	}
}
