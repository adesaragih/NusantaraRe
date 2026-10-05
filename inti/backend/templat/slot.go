// Package templat adalah Template Manager (keputusan work owner 04-10-2026): berkas templat unduhan milik
// banyak menu - Bordereaux, Aggregate, kelak NB/RNW Fac In - dikelola di SATU tempat, supaya tim IT dapat
// mengganti berkasnya tanpa deploy. Padanan Pega: Rule-File-Binary yang dicari `pyFileName` lalu diambil yang
// terbaru (`DownloadTemplate_Act` langkah 3, `pxCommitDateTime` turun).
//
// Setiap modul MENDAFTARKAN slot templatnya (`inti.Pendaftaran.Templat`); slot menentukan untuk menu apa, dipakai
// di mana, berkas bawaannya, dan berapa kolom header yang dibaca Upload CSV modul itu. Berkas yang diunggah
// disimpan berversi di `M_TEMPLATE_FILE` (migrasi inti 912); slot tanpa versi aktif memakai berkas bawaan.
//
// ⛔ Paket daun: TIDAK mengimpor `inti/backend` (paket itu mengimpor paket ini lewat `Pendaftaran`). Rute HTTP
// ada di `templat/rute`.
package templat

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"
)

// Slot adalah satu templat unduhan yang didaftarkan sebuah modul.
type Slot struct {
	// Kode pengenal tetap, huruf kecil bertitik, mis. `bordereaux.premi.fire`. Disimpan di M_TEMPLATE_FILE.KODE:
	// nilainya tidak boleh berubah sesudah ada versi yang diunggah.
	Kode string
	// Menu - KODE menu pemilik (`M_LOGIN_GO_MENU.MENU_KODE`): pemegangnya boleh mengunduh templat ini.
	Menu string
	// Grup - nama menu yang tampil di Template Manager, mis. "Bordereaux".
	Grup string
	// Nama - nama slot di dalam grupnya, mis. "PREMIUM · FIRE".
	Nama string
	// DipakaiDi - letak tombol unduhnya, mis. "Bordereaux › Input Data › Details › Template".
	DipakaiDi string
	// Ekstensi berkas yang diterima, huruf kecil berawalan titik, mis. ".csv".
	Ekstensi string
	// Pemisah kolom CSV (`;` atau `,`); 0 = berkas bukan CSV (header tidak diperiksa).
	Pemisah rune
	// JumlahKolom - jumlah kolom header yang dibaca Upload CSV modulnya; 0 = tidak diperiksa.
	JumlahKolom int
	// NamaUnduhan - nama berkas saat diunduh pengguna, apa pun nama berkas yang diunggah.
	NamaUnduhan string
	// Bawaan - isi berkas bawaan (ikut terpasang di aplikasi); dipakai selama belum ada versi aktif.
	Bawaan []byte
}

// ErrSlotTidakAda - kode slot tidak terdaftar.
var ErrSlotTidakAda = errors.New("templat: slot tidak terdaftar")

var polaKode = regexp.MustCompile(`^[a-z0-9]+(\.[a-z0-9_-]+)+$`)

// Periksa menolak slot yang bentuknya salah - dipanggil saat katalog dirakit, jadi kesalahan modul gagal TERANG
// saat proses menyala, bukan saat permintaan pertama.
func (s Slot) Periksa() error {
	switch {
	case !polaKode.MatchString(s.Kode) || len(s.Kode) > 100:
		return fmt.Errorf("templat: kode slot %q tidak sah (huruf kecil bertitik, maks. 100)", s.Kode)
	case strings.TrimSpace(s.Menu) == "":
		return fmt.Errorf("templat: slot %s tanpa menu pemilik", s.Kode)
	case strings.TrimSpace(s.Grup) == "" || strings.TrimSpace(s.Nama) == "":
		return fmt.Errorf("templat: slot %s tanpa grup atau nama", s.Kode)
	case !strings.HasPrefix(s.Ekstensi, ".") || s.Ekstensi != strings.ToLower(s.Ekstensi):
		return fmt.Errorf("templat: slot %s: ekstensi %q harus huruf kecil berawalan titik", s.Kode, s.Ekstensi)
	case s.Pemisah != 0 && s.Pemisah != ';' && s.Pemisah != ',':
		return fmt.Errorf("templat: slot %s: pemisah hanya ';' atau ','", s.Kode)
	case s.JumlahKolom < 0 || (s.JumlahKolom > 0 && s.Pemisah == 0):
		return fmt.Errorf("templat: slot %s: jumlah kolom hanya untuk slot CSV", s.Kode)
	case strings.ToLower(path.Ext(s.NamaUnduhan)) != s.Ekstensi:
		return fmt.Errorf("templat: slot %s: nama unduhan %q tidak berakhiran %s", s.Kode, s.NamaUnduhan, s.Ekstensi)
	case len(s.Bawaan) == 0:
		return fmt.Errorf("templat: slot %s tanpa berkas bawaan", s.Kode)
	}
	if s.JumlahKolom > 0 {
		kolom, err := Kepala(s.Bawaan, s.Pemisah)
		if err != nil {
			return fmt.Errorf("templat: slot %s: berkas bawaan: %w", s.Kode, err)
		}
		if len(kolom) != s.JumlahKolom {
			return fmt.Errorf("templat: slot %s: berkas bawaan %d kolom, slot %d", s.Kode, len(kolom), s.JumlahKolom)
		}
	}
	return nil
}

// Katalog adalah seluruh slot terdaftar, urut Grup lalu urutan daftar modulnya.
type Katalog struct {
	urut []Slot
	per  map[string]Slot
}

// NewKatalog merakit katalog; kode ganda atau slot cacat ditolak.
func NewKatalog(slot []Slot) (*Katalog, error) {
	k := &Katalog{per: map[string]Slot{}}
	for _, s := range slot {
		if err := s.Periksa(); err != nil {
			return nil, err
		}
		if _, ada := k.per[s.Kode]; ada {
			return nil, fmt.Errorf("templat: kode slot %s terdaftar dua kali", s.Kode)
		}
		k.per[s.Kode] = s
		k.urut = append(k.urut, s)
	}
	sort.SliceStable(k.urut, func(i, j int) bool { return k.urut[i].Grup < k.urut[j].Grup })
	return k, nil
}

// Semua - salinan slot, urut tampil.
func (k *Katalog) Semua() []Slot { return append([]Slot(nil), k.urut...) }

// Ambil - satu slot menurut kode.
func (k *Katalog) Ambil(kode string) (Slot, error) {
	s, ok := k.per[kode]
	if !ok {
		return Slot{}, ErrSlotTidakAda
	}
	return s, nil
}
