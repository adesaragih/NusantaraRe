// Package menu menyusun menu aplikasi dari tabel `M_NAV_MENU` untuk
// `GET /api/menu`: GROUPMENU -> modul, SATU tingkat.
//
// Permintaan work owner 30-09-2026 (`PROMPT-MENU-DARI-TABEL-M_NAV_MENU.md`):
// menu dibuat dari tabel supaya kelak dapat disaring per akun sesudah login
// ada. Keputusan work owner 30-09-2026 berikutnya
// (`PROMPT-MENU-DATAR-PER-GROUPMENU.md`): "menu jangan ada model seperti
// child ... 1 modul 1 menu" - tidak ada butir di bawah modul; klik tombol modul
// membuka halaman awalnya (frontend `HALAMAN_AWAL_<X>`). Paket ini milik
// `inti` karena menu memuat SEMUA modul - termasuk yang belum dimigrasi - dan
// tidak pernah mengimpor modul mana pun.
//
// Tiga berkas, tiga lapis: `menu.go` (bentuk dan aturan penyusunan, tanpa I/O),
// `pembaca.go` (SQL), `rute.go` (HTTP).
package menu

import (
	"sort"

	inti "nusantarare/inti/backend"
)

// Golongan adalah isi CHECK `GROUPMENU`, dalam urutan tampil sidebar - urutan
// yang work owner tulis ("TREATY, FACULTATIVE, KLAIM, MASTER"), bukan urutan
// abjad `ORDER BY GROUPMENU`. `inti/backend/penjaga/menu_test.go` membandingkannya
// dengan CHECK di migrasi 900.
var Golongan = []string{"TREATY", "FACULTATIVE", "KLAIM", "MASTER"}

// Baris adalah satu baris MODUL `M_NAV_MENU` yang aktif (`KODE = MODUL`).
type Baris struct {
	ID        int64
	Kode      string
	Label     string
	Golongan  string
	Modul     string
	Urutan    int
	Dimigrasi bool
}

// Modul adalah satu tombol menu - satu modul.
type Modul struct {
	// Kode = nama modul backend (tabel nama modul), sama dengan Modul.
	Kode string `json:"kode"`
	// Label = nama folder korpus VERBATIM - label tombol.
	Label string `json:"label"`
	Modul string `json:"modul"`
	// Urutan di dalam golongannya.
	Urutan int `json:"urutan"`
	// Dimigrasi false = modul belum punya layar (tombol nonaktif,
	// "belum dimigrasi").
	Dimigrasi bool `json:"dimigrasi"`
}

// BagianGolongan adalah satu kepala bagian sidebar (TREATY, ...).
type BagianGolongan struct {
	Kode  string  `json:"kode"`
	Modul []Modul `json:"modul"`
}

// Menu adalah badan `GET /api/menu`.
type Menu struct {
	Golongan []BagianGolongan `json:"golongan"`
}

// Susun membangun menu dari baris aktif `M_NAV_MENU`.
//
// Aturannya:
//   - golongan menurut `Golongan`; golongan tanpa modul yang dikirim tidak
//     dikirim
//   - modul menurut URUTAN, lalu ID - tidak bergantung pada ORDER BY
//     pembacanya
//   - modul yang SUDAH dimigrasi tetapi tidak ada di `modulAktif`
//     (MODUL_AKTIF) tidak dikirim: halamannya tidak ada di proses ini. Modul
//     yang BELUM dimigrasi tetap dikirim - frontend menampilkannya nonaktif
//   - baris yang bukan baris modul (`KODE <> MODUL`: lima butir anak 900
//     selama 901 belum berjalan) tidak dikirim, juga bila pembaca
//     meloloskannya
//
// ⚠️ GROUPMENU di luar `Golongan` tidak mungkin - CHECK menolaknya di Oracle -
// jadi modul begitu tidak punya tempat dan tidak dikirim.
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

	modulDari := map[string][]Modul{}
	for _, b := range urut {
		if b.Kode != b.Modul || (b.Dimigrasi && !aktif[b.Modul]) {
			continue
		}
		modulDari[b.Golongan] = append(modulDari[b.Golongan],
			Modul{Kode: b.Kode, Label: b.Label, Modul: b.Modul, Urutan: b.Urutan, Dimigrasi: b.Dimigrasi})
	}
	m := Menu{Golongan: []BagianGolongan{}}
	for _, g := range Golongan {
		if len(modulDari[g]) > 0 {
			m.Golongan = append(m.Golongan, BagianGolongan{Kode: g, Modul: modulDari[g]})
		}
	}
	return m
}

// GolonganAdmin adalah kepala bagian menu APLIKASI - di bawah golongan tabel.
//
// ⛔ BUKAN isi CHECK `GROUPMENU`: menu aplikasi bukan baris `M_NAV_MENU`, yang
// dijaga tepat dua puluh baris, satu per folder modul korpus (keputusan work
// owner 30-09-2026, "1 modul 1 menu").
const GolonganAdmin = "ADMIN"

// KodeKelolaUser adalah KODE menu Kelola User (keputusan work owner
// 01-10-2026) - nilai `M_LOGIN_GO_MENU.MENU_KODE` yang membuka layar dan API
// kelola pengguna. Pemegangnya adalah admin.
const KodeKelolaUser = "kelolauser"

// MenuAplikasi adalah menu milik aplikasi, bukan modul korpus: hidup di kode
// seperti Beranda, tampil di golongan `GolonganAdmin` hanya bagi akun yang
// memegang KODE-nya.
var MenuAplikasi = []Modul{
	{Kode: KodeKelolaUser, Label: "Kelola User", Modul: KodeKelolaUser, Urutan: 1, Dimigrasi: true},
}

// SaringMenuUntukAkun menyisakan menu yang KODE-nya dimiliki akun
// (`M_LOGIN_GO_MENU`, Kelola User 01-10-2026): modul lain tidak dikirim,
// golongan yang kosong hilang, lalu menu aplikasi yang dimiliki menyusul di
// golongan `GolonganAdmin`. Menu asalnya tidak diubah.
//
// ⛔ Ini saringan TAMPILAN. Penegaknya gerbang 403 di `cmd/api` (rute modul)
// dan rute Kelola User sendiri - menu yang disembunyikan di sini tetap ditolak
// di sana bila jalurnya diketik langsung.
func SaringMenuUntukAkun(m Menu, kode []string) Menu {
	hasil := Menu{Golongan: []BagianGolongan{}}
	for _, g := range m.Golongan {
		var modul []Modul
		for _, x := range g.Modul {
			if inti.PunyaMenu(kode, x.Kode) {
				modul = append(modul, x)
			}
		}
		if len(modul) > 0 {
			hasil.Golongan = append(hasil.Golongan, BagianGolongan{Kode: g.Kode, Modul: modul})
		}
	}
	var aplikasi []Modul
	for _, x := range MenuAplikasi {
		if inti.PunyaMenu(kode, x.Kode) {
			aplikasi = append(aplikasi, x)
		}
	}
	if len(aplikasi) > 0 {
		hasil.Golongan = append(hasil.Golongan, BagianGolongan{Kode: GolonganAdmin, Modul: aplikasi})
	}
	return hasil
}
