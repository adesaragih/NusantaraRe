package backend

// Antarmuka modul untuk perakitan aplikasi - refactor bentuk B.
//
// Untuk apa berkas ini: `cmd/api` tidak mengenal isi modul mana pun. Setiap
// modul menyerahkan dirinya lewat `Modul`, dan `cmd/api` hanya memanggil
// daftar modul yang AKTIF (MODUL_AKTIF): rutenya didaftarkan, pekerjanya
// dinyalakan. Migrasi TIDAK lewat sini - ia selalu dari SEMUA modul terdaftar
// (`daftar.SumberMigrasi`), supaya skema selalu utuh.

import (
	"context"
	"net/http"
)

// Modul adalah yang diserahkan setiap modul kepada `cmd/api`.
type Modul interface {
	// Nama pengenal modul di MODUL_AKTIF dan di GET /api/modul-aktif.
	Nama() string
	// DaftarkanRute mendaftarkan seluruh rute HTTP modul ke mux bersama.
	DaftarkanRute(mux *http.ServeMux)
	// JalankanPekerja menyalakan pekerja latar modul, bila ada. Ia berhenti
	// bersama ctx.
	JalankanPekerja(ctx context.Context) Pekerja
}

// Pekerja adalah pekerja latar sebuah modul yang sedang berjalan.
type Pekerja struct {
	// Selesai tertutup saat pekerja modul BENAR-BENAR berhenti; langsung
	// tertutup bila modul tidak punya pekerja atau pekerjanya tidak dinyalakan.
	Selesai <-chan struct{}
	// PesanTerlambat dicetak `cmd/api` bila pekerja belum berhenti ketika batas
	// waktu penutupan habis. Kosong = tidak ada yang perlu dicetak.
	PesanTerlambat string
}

// TanpaPekerja - untuk modul yang tidak punya pekerja latar.
func TanpaPekerja() Pekerja {
	selesai := make(chan struct{})
	close(selesai)
	return Pekerja{Selesai: selesai}
}
