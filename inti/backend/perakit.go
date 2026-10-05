package backend

// Perakit modul - struktur tim satu folder per modul, R4 (30-09-2026).
//
// Untuk apa berkas ini: sampai 30-09-2026 `modul/daftar.go` menyambung kontrak
// lintas modul DENGAN TANGAN (PembacaPolis PremiumList ke Claim Life,
// KlaimKomite Claim Life ke Komite). Berkas itu disunting setiap kali modul
// bertambah - berkas bersama yang membuat dua pengembang bertabrakan. Kini
// setiap `modul/<nama>/backend/modul.go` MENYATAKAN kontrak yang disediakan
// dan yang dibutuhkannya (`Pendaftaran`), dan perakit di sini menyambungnya
// menurut JENIS antarmuka (`inti/backend/kontrak`).
//
// ⛔ Kontrak yang dibutuhkan tanpa penyedia adalah GALAT saat menyala yang
// menyebut kontrak dan modulnya - bukan nil diam-diam yang baru meledak saat
// permintaan pertama tiba.
//
// Istilah:
//   - kontrak   : antarmuka di `inti/backend/kontrak`, tanpa implementasi.
//   - penyedia  : modul yang menyerahkan implementasi sebuah kontrak.
//   - pemakai   : modul yang menerima implementasi itu saat dirakit.

import (
	"fmt"
	"io/fs"
	"reflect"
	"sort"

	"nusantarare/inti/backend/config"
	"nusantarare/inti/backend/templat"
)

// Kontrak adalah kunci satu kontrak lintas modul: jenis antarmukanya.
type Kontrak struct{ jenis reflect.Type }

// KontrakDari menyebut antarmuka T sebagai kontrak, mis.
// `inti.KontrakDari[kontrak.PembacaPolis]()`.
func KontrakDari[T any]() Kontrak { return Kontrak{jenis: reflect.TypeFor[T]()} }

// String menyebut kontrak seperti ditulis di Go, mis. `kontrak.PembacaPolis`.
func (k Kontrak) String() string {
	if k.jenis == nil {
		return "<kontrak kosong>"
	}
	return k.jenis.String()
}

// Pendaftaran adalah yang diserahkan `backend/modul.go` setiap modul kepada
// perakit - lewat fungsi `Pendaftaran()` di paket itu, yang dipanggil daftar
// modul bangkitan (`inti/backend/daftar`).
type Pendaftaran struct {
	// Nama pengenal modul: nama folder `modul/<nama>`, MODUL_AKTIF, dan
	// GET /api/modul-aktif.
	Nama string
	// NamaLama - nama modul sebelum tabel nama modul (30-09-2026), untuk
	// MENOLAK nama itu di MODUL_AKTIF dengan kalimat yang menyebut nama barunya.
	NamaLama []string
	// Migrasi - folder `migrations/` modul ini; nil bila modul tidak
	// bermigrasi (mis. Treaty Contract Out, tco4).
	Migrasi fs.FS
	// Menyediakan - kontrak yang modul ini serahkan kepada modul lain. Setiap
	// kontrak di sini WAJIB diserahkan lewat `Sediakan` di dalam `Bangun`.
	Menyediakan []Kontrak
	// Membutuhkan - kontrak yang modul ini terima dari modul lain, dibaca
	// lewat `Ambil` di dalam `Bangun`.
	Membutuhkan []Kontrak
	// Bangun merakit modul di atas akar bersama. Dipanggil SESUDAH setiap
	// penyedia kontrak yang dibutuhkannya dibangun.
	Bangun func(p *Perakitan) (Modul, error)
	// Templat - slot berkas templat unduhan modul ini, dikelola menu Template
	// Manager (`inti/backend/templat`, keputusan work owner 04-10-2026). Slot
	// dari SEMUA modul terdaftar dirakit, juga modul yang tidak aktif.
	Templat []templat.Slot
	// HakLihat - nil = menu modul ini selalu berakses PENUH. Terisi = modul MENDAFTAR untuk akses LIHAT
	// (`M_LOGIN_GO_MENU.HAK`, keputusan work owner 04-10-2026): Kelola User menawarkan pilihan Lihat/Penuh, gerbang
	// `cmd/api` menolak rute tulis bagi pemegang LIHAT, dan layar modul menyembunyikan tombol tulisnya. Hanya modul
	// yang sudah selesai yang mendaftar - modul yang belum selesai tidak berubah perilakunya.
	HakLihat *HakLihat
}

// HakLihat - pernyataan satu modul untuk akses menu LIHAT.
type HakLihat struct {
	// Bebas - pola rute modul ini (persis seperti didaftarkan ke mux, mis. `POST /api/aggregate/pratinjau`) yang TETAP
	// boleh dipakai pemegang LIHAT: POST yang hanya membaca, atau keputusan alur kerja yang dijaga modulnya sendiri
	// (mis. Approve/Reject Checker). Selain GET/HEAD dan pola ini, pemegang LIHAT dijawab 403.
	Bebas []string
}

// Perakitan adalah yang diterima `Bangun` satu modul: akar bersama,
// konfigurasi, pencatat proses, dan kontrak yang dibutuhkannya.
type Perakitan struct {
	dasar           *Dasar
	cfg             config.Config
	catat           func(string)
	modul           string
	tersedia        map[Kontrak]any
	bolehDisediakan map[Kontrak]bool
	disediakan      map[Kontrak]any
	galat           []error
}

// Dasar - akar bersama (koneksi, lingkungan, folder unggahan).
func (p *Perakitan) Dasar() *Dasar { return p.dasar }

// Config - konfigurasi proses.
func (p *Perakitan) Config() config.Config { return p.cfg }

// Catat - pencatat proses (log), untuk modul yang mencatat saat menyala.
func (p *Perakitan) Catat(s string) { p.catat(s) }

// Ambil membaca kontrak T yang dibutuhkan modul yang sedang dirakit.
func Ambil[T any](p *Perakitan) T {
	k := KontrakDari[T]()
	var kosong T
	nilai, ada := p.tersedia[k]
	if !ada {
		p.galat = append(p.galat, fmt.Errorf("modul %s mengambil kontrak %s yang tidak ada di Membutuhkan-nya", p.modul, k))
		return kosong
	}
	t, _ := nilai.(T) // nilai selalu T: Sediakan[T] yang mengisinya
	return t
}

// Sediakan menyerahkan implementasi kontrak T dari modul yang sedang dirakit.
func Sediakan[T any](p *Perakitan, nilai T) {
	k := KontrakDari[T]()
	if !p.bolehDisediakan[k] {
		p.galat = append(p.galat, fmt.Errorf("modul %s menyediakan kontrak %s yang tidak ada di Menyediakan-nya", p.modul, k))
		return
	}
	// ⛔ nil ANTARMUKA maupun penunjuk nil di dalamnya: keduanya meledak saat
	// dipanggil pemakai, jauh dari sebabnya.
	if v := reflect.ValueOf(nilai); !v.IsValid() || (v.Kind() == reflect.Pointer && v.IsNil()) {
		p.galat = append(p.galat, fmt.Errorf("modul %s menyediakan kontrak %s bernilai nil", p.modul, k))
		return
	}
	p.disediakan[k] = nilai
}

// Sambungan adalah satu kontrak yang disambung perakit: dari penyedia ke pemakai.
type Sambungan struct {
	Kontrak  string
	Penyedia string
	Pemakai  string
}

// Rakitan adalah hasil perakitan.
type Rakitan struct {
	// Modul - setiap modul terdaftar, berurutan menurut nama.
	Modul []Modul
	// Sambungan - setiap kontrak yang disambung, berurutan menurut kontrak
	// lalu pemakai.
	Sambungan []Sambungan
}

// Rakit membangun SETIAP modul terdaftar di atas satu akar bersama, dan
// menyambung kontraknya menurut jenis antarmuka.
//
// Modul dibangun sesudah penyedia setiap kontrak yang dibutuhkannya; urutan
// hasilnya urutan NAMA modul, bukan urutan bangun dan bukan urutan daftar.
func Rakit(dasar *Dasar, cfg config.Config, catat func(string), daftar []Pendaftaran) (Rakitan, error) {
	penyedia, err := periksaDaftar(daftar)
	if err != nil {
		return Rakitan{}, err
	}
	urut, err := urutanBangun(daftar, penyedia)
	if err != nil {
		return Rakitan{}, err
	}
	nilai := map[Kontrak]any{}
	var hasil Rakitan
	for _, d := range urut {
		p := &Perakitan{dasar: dasar, cfg: cfg, catat: catat, modul: d.Nama,
			tersedia: map[Kontrak]any{}, bolehDisediakan: map[Kontrak]bool{}, disediakan: map[Kontrak]any{}}
		for _, k := range d.Membutuhkan {
			p.tersedia[k] = nilai[k]
			hasil.Sambungan = append(hasil.Sambungan, Sambungan{Kontrak: k.String(), Penyedia: penyedia[k], Pemakai: d.Nama})
		}
		for _, k := range d.Menyediakan {
			p.bolehDisediakan[k] = true
		}
		m, err := d.Bangun(p)
		if err != nil {
			return Rakitan{}, fmt.Errorf("modul %s: %w", d.Nama, err)
		}
		if len(p.galat) > 0 {
			return Rakitan{}, p.galat[0]
		}
		if m == nil || m.Nama() != d.Nama {
			nama := "<nil>"
			if m != nil {
				nama = m.Nama()
			}
			return Rakitan{}, fmt.Errorf("modul terdaftar %s membangun modul bernama %s; keduanya harus sama "+
				"(nama folder modul/<nama>, MODUL_AKTIF)", d.Nama, nama)
		}
		for _, k := range d.Menyediakan {
			if _, ada := p.disediakan[k]; !ada {
				return Rakitan{}, fmt.Errorf("modul %s menyatakan Menyediakan kontrak %s, tetapi Bangun-nya tidak "+
					"menyerahkannya lewat Sediakan", d.Nama, k)
			}
		}
		for k, v := range p.disediakan {
			nilai[k] = v
		}
		hasil.Modul = append(hasil.Modul, m)
	}
	sort.Slice(hasil.Modul, func(i, j int) bool { return hasil.Modul[i].Nama() < hasil.Modul[j].Nama() })
	sort.Slice(hasil.Sambungan, func(i, j int) bool {
		a, b := hasil.Sambungan[i], hasil.Sambungan[j]
		if a.Kontrak != b.Kontrak {
			return a.Kontrak < b.Kontrak
		}
		return a.Pemakai < b.Pemakai
	})
	return hasil, nil
}

// periksaDaftar menolak daftar yang bentuknya salah SEBELUM satu modul pun
// dibangun, dan mengembalikan penyedia setiap kontrak.
func periksaDaftar(daftar []Pendaftaran) (map[Kontrak]string, error) {
	nama := map[string]bool{}
	penyedia := map[Kontrak]string{}
	for _, d := range daftar {
		if d.Nama == "" {
			return nil, fmt.Errorf("sebuah pendaftaran modul tanpa Nama")
		}
		if nama[d.Nama] {
			return nil, fmt.Errorf("modul %s terdaftar dua kali", d.Nama)
		}
		nama[d.Nama] = true
		if d.Bangun == nil {
			return nil, fmt.Errorf("modul %s terdaftar tanpa Bangun", d.Nama)
		}
		for _, k := range append(append([]Kontrak(nil), d.Menyediakan...), d.Membutuhkan...) {
			if k.jenis == nil || k.jenis.Kind() != reflect.Interface {
				return nil, fmt.Errorf("modul %s menyebut %s sebagai kontrak; kontrak harus antarmuka "+
					"(inti/backend/kontrak)", d.Nama, k)
			}
		}
		for _, k := range d.Menyediakan {
			if lain, ada := penyedia[k]; ada {
				return nil, fmt.Errorf("kontrak %s disediakan dua modul: %s dan %s; satu kontrak satu penyedia",
					k, lain, d.Nama)
			}
			penyedia[k] = d.Nama
		}
	}
	for _, d := range daftar {
		for _, k := range d.Membutuhkan {
			if _, ada := penyedia[k]; !ada {
				return nil, fmt.Errorf("modul %s membutuhkan kontrak %s, tetapi tidak ada modul terdaftar "+
					"yang menyediakannya (Menyediakan di backend/modul.go penyedianya)", d.Nama, k)
			}
		}
	}
	return penyedia, nil
}

// urutanBangun menyusun modul sehingga setiap penyedia dibangun sebelum
// pemakainya; di antara yang sama-sama siap, menurut nama. Ketergantungan
// melingkar DITOLAK - tanpa penjagaan ini perakit berputar tanpa henti.
func urutanBangun(daftar []Pendaftaran, penyedia map[Kontrak]string) ([]Pendaftaran, error) {
	sisa := append([]Pendaftaran(nil), daftar...)
	sort.Slice(sisa, func(i, j int) bool { return sisa[i].Nama < sisa[j].Nama })
	selesai := map[string]bool{}
	var urut []Pendaftaran
	for len(sisa) > 0 {
		maju := false
		for i, d := range sisa {
			siap := true
			for _, k := range d.Membutuhkan {
				if !selesai[penyedia[k]] {
					siap = false
				}
			}
			if siap {
				urut = append(urut, d)
				selesai[d.Nama] = true
				sisa = append(sisa[:i], sisa[i+1:]...)
				maju = true
				break
			}
		}
		if !maju {
			var tersangkut []string
			for _, d := range sisa {
				tersangkut = append(tersangkut, d.Nama)
			}
			return nil, fmt.Errorf("ketergantungan kontrak melingkar di antara modul %v; "+
				"setiap modul menunggu kontrak modul lain", tersangkut)
		}
	}
	return urut, nil
}
