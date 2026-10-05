// Package daftar adalah DAFTAR modul aplikasi - satu-satunya paket yang
// mengenal semua modul sekaligus.
//
// Struktur tim satu folder per modul (30-09-2026): daftar ini TIDAK ditulis
// tangan. `go generate ./inti/backend/daftar` (dari APP_RNM) menulis satu
// berkas `modul_<nama>_gen.go` per folder `modul/<nama>/backend/modul.go`, dan
// berkas itu mendaftarkan `<nama>.Pendaftaran()` ke sini. Kontrak lintas modul
// disambung perakit `inti.Rakit` menurut pernyataan setiap modul - nol baris
// per modul di berkas ini. Pembangkitnya: `bangkit/main.go`; penjaga "hasil
// bangkit = isi folder": `bangkit/main_test.go`.
//
// ⛔ Satu-satunya paket `inti` yang mengimpor modul, dan hanya paket akar
// `nusantarare/modul/<nama>/backend` (dijaga
// `inti/backend/penjaga/impor_lintas_modul_test.go`). Ia tidak dapat tinggal di
// akar `inti/backend`: setiap modul mengimpor akar itu, jadi daftar di sana
// menjadi impor melingkar.
//
// Yang memakai daftar ini: `cmd/api` dan skema uji test bertag `db`.
package daftar

// Direktif `go generate`: `pembangkit.go` (letaknya disengaja, lihat di sana).

import (
	"io/fs"
	"sort"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/templat"
)

// terdaftar diisi fungsi `init` setiap berkas bangkitan `modul_<nama>_gen.go`.
var terdaftar []inti.Pendaftaran

// daftarkan dipanggil berkas bangkitan - satu kali per modul.
func daftarkan(p inti.Pendaftaran) { terdaftar = append(terdaftar, p) }

// Terdaftar mengembalikan pendaftaran setiap modul, berurutan menurut nama.
func Terdaftar() []inti.Pendaftaran {
	hasil := append([]inti.Pendaftaran(nil), terdaftar...)
	sort.Slice(hasil, func(i, j int) bool { return hasil[i].Nama < hasil[j].Nama })
	return hasil
}

// Rakit membangun SETIAP modul terdaftar di atas satu akar bersama dan
// menyambung kontrak lintas modulnya (`inti.Rakit`).
//
// Urutannya urutan NAMA modul: urutan pendaftaran rute dan urutan
// GET /api/modul-aktif. `catat` adalah pencatat proses (log) untuk modul yang
// mencatat saat menyala. Galat = daftar yang bentuknya salah, mis. kontrak yang
// dibutuhkan tanpa penyedia; `cmd/api` menolak menyala.
func Rakit(dasar *inti.Dasar, cfg config.Config, catat func(string)) (inti.Rakitan, error) {
	return inti.Rakit(dasar, cfg, catat, Terdaftar())
}

// SumberMigrasi mengembalikan folder migrasi SETIAP modul terdaftar, ditambah
// migrasi `inti` (tabel lintas modul, 900-949).
//
// ⛔ Tidak bergantung pada modul yang aktif: skema selalu utuh, sebab tabel
// satu modul dapat dirujuk tabel modul lain (dan data warisan tidak memilih
// modul). Pelari mengurutkan langkah menurut NAMA berkas, lintas sumber.
// Modul tanpa migrasi (Treaty Contract Out, tco4) tidak menyumbang sumber.
func SumberMigrasi() []fs.FS {
	var sumber []fs.FS
	for _, p := range Terdaftar() {
		if p.Migrasi != nil {
			sumber = append(sumber, p.Migrasi)
		}
	}
	return append(sumber, inti.SumberMigrasi())
}

// SlotTemplat - slot templat unduhan dari SEMUA modul terdaftar (Template
// Manager, keputusan work owner 04-10-2026), urut nama modul.
func SlotTemplat() []templat.Slot {
	var slot []templat.Slot
	for _, p := range Terdaftar() {
		slot = append(slot, p.Templat...)
	}
	return slot
}

// HakLihat - modul terdaftar yang mendukung akses menu LIHAT (keputusan work owner 04-10-2026): nama modul (= KODE
// menunya) ke pola rute yang tetap boleh dipakai pemegang LIHAT. Modul tanpa pernyataan tidak ada di peta ini.
func HakLihat() map[string][]string {
	hasil := map[string][]string{}
	for _, p := range Terdaftar() {
		if p.HakLihat != nil {
			hasil[p.Nama] = append([]string{}, p.HakLihat.Bebas...)
		}
	}
	return hasil
}

// NamaLama memetakan nama modul SEBELUM tabel nama modul (keputusan work
// owner 30-09-2026, `PROMPT-REFACTOR-NAMA-MODUL.md`) ke namanya kini, dari
// pernyataan `NamaLama` setiap modul. Hanya dipakai untuk MENOLAK nama lama di
// MODUL_AKTIF dengan kalimat yang menyebut nama barunya - tidak pernah untuk
// menerimanya.
func NamaLama() map[string]string {
	hasil := map[string]string{}
	for _, p := range Terdaftar() {
		for _, lama := range p.NamaLama {
			hasil[lama] = p.Nama
		}
	}
	return hasil
}
