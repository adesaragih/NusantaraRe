// Package menu menyusun menu aplikasi dari tabel `M_NAV_MENU` untuk
// `GET /api/menu`: pohon GROUPMENU -> kelompok modul -> butir menu.
//
// Permintaan work owner 30-09-2026 (`PROMPT-MENU-DARI-TABEL-M_NAV_MENU.md`):
// menu dibuat dari tabel supaya kelak dapat disaring per akun sesudah login
// ada. Paket ini milik `inti` karena menu memuat SEMUA modul - termasuk yang
// belum dimigrasi - dan tidak pernah mengimpor modul mana pun.
//
// Tiga berkas, tiga lapis: `menu.go` (bentuk dan aturan penyusunan, tanpa I/O),
// `pembaca.go` (SQL), `rute.go` (HTTP).
package menu

import (
	"sort"

	"nusantarare/inti"
)

// Golongan adalah isi CHECK `GROUPMENU`, dalam urutan tampil sidebar - urutan
// yang work owner tulis ("TREATY, FACULTATIVE, KLAIM, MASTER"), bukan urutan
// abjad `ORDER BY GROUPMENU`. `inti/penjaga/menu_test.go` membandingkannya
// dengan CHECK di migrasi 900.
var Golongan = []string{"TREATY", "FACULTATIVE", "KLAIM", "MASTER"}

// Baris adalah satu baris `M_NAV_MENU` yang aktif.
type Baris struct {
	ID int64
	// IndukID 0 = baris KELOMPOK modul (`PARENT_ID` kosong).
	IndukID   int64
	Kode      string
	Label     string
	Golongan  string
	Modul     string
	Urutan    int
	Dimigrasi bool
}

// Butir adalah satu butir menu - halaman frontend yang dapat dibuka.
type Butir struct {
	// Kode = kunci halaman frontend (`modul/daftar.ts`).
	Kode  string `json:"kode"`
	Label string `json:"label"`
	Modul string `json:"modul"`
}

// Kelompok adalah satu kelompok modul beserta butirnya.
type Kelompok struct {
	// Kode = nama modul backend (tabel nama modul).
	Kode  string `json:"kode"`
	Label string `json:"label"`
	Modul string `json:"modul"`
	// Dimigrasi false = modul belum punya layar ("belum dimigrasi").
	Dimigrasi bool    `json:"dimigrasi"`
	Butir     []Butir `json:"butir"`
}

// BagianGolongan adalah satu kepala bagian sidebar (TREATY, ...).
type BagianGolongan struct {
	Kode     string     `json:"kode"`
	Kelompok []Kelompok `json:"kelompok"`
}

// Menu adalah badan `GET /api/menu`.
type Menu struct {
	Golongan []BagianGolongan `json:"golongan"`
}

// Susun membangun pohon dari baris aktif `M_NAV_MENU`.
//
// Aturannya:
//   - golongan menurut `Golongan`; golongan tanpa kelompok tidak dikirim
//   - kelompok dan butir menurut URUTAN, lalu ID - tidak bergantung pada
//     ORDER BY pembacanya
//   - butir milik modul yang TIDAK ada di `modulAktif` (MODUL_AKTIF) tidak
//     dikirim; kelompoknya tetap, dengan butir kosong
//   - butir yang induknya tidak terbaca (induk nonaktif) dan butir di bawah
//     butir (tingkat ketiga) tidak dikirim: tabel ini dua tingkat
//
// ⚠️ GROUPMENU di luar `Golongan` tidak mungkin - CHECK menolaknya di Oracle -
// jadi kelompok begitu tidak punya tempat dan tidak dikirim.
func Susun(baris []Baris, modulAktif []string) Menu {
	aktif := map[string]bool{}
	for _, n := range modulAktif {
		aktif[n] = true
	}
	urut := append([]Baris(nil), baris...)
	sort.SliceStable(urut, func(i, j int) bool {
		if urut[i].Urutan != urut[j].Urutan {
			return urut[i].Urutan < urut[j].Urutan
		}
		return urut[i].ID < urut[j].ID
	})

	butirDari := map[int64][]Butir{}
	for _, b := range urut {
		if b.IndukID != 0 && aktif[b.Modul] {
			butirDari[b.IndukID] = append(butirDari[b.IndukID], Butir{Kode: b.Kode, Label: b.Label, Modul: b.Modul})
		}
	}
	kelompokDari := map[string][]Kelompok{}
	for _, k := range urut {
		if k.IndukID != 0 {
			continue
		}
		butir := butirDari[k.ID]
		if butir == nil {
			butir = []Butir{}
		}
		kelompokDari[k.Golongan] = append(kelompokDari[k.Golongan],
			Kelompok{Kode: k.Kode, Label: k.Label, Modul: k.Modul, Dimigrasi: k.Dimigrasi, Butir: butir})
	}
	m := Menu{Golongan: []BagianGolongan{}}
	for _, g := range Golongan {
		if len(kelompokDari[g]) > 0 {
			m.Golongan = append(m.Golongan, BagianGolongan{Kode: g, Kelompok: kelompokDari[g]})
		}
	}
	return m
}

// SaringMenuUntukPelaku adalah TITIK SAMBUNG saringan menu per akun.
//
// ⛔ Akses per akun menyusul: hari ini ia meneruskan SEMUA, untuk pelaku siapa
// pun. Tabel aksesnya (mis. `M_NAV_MENU_AKSES`: akun atau peran -> MENU_ID)
// dan login berada di luar lingkup brief menu 30-09-2026 - dicatat, tidak
// dibangun. Sesudah keduanya ada, saringannya disambung di sini, dan
// `GET /api/menu` sudah memanggilnya.
func SaringMenuUntukPelaku(_ inti.Pelaku, m Menu) Menu {
	return m
}
