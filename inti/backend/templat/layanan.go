package templat

import (
	"context"
	"errors"
	"path"
	"strconv"
	"strings"
	"unicode/utf8"
)

// ErrDitolak - berkas atau permintaan tidak memenuhi aturan slot; pesannya untuk pengguna.
var ErrDitolak = errors.New("templat: ditolak")

// GalatPeriksa membawa daftar alasan penolakan berkas.
type GalatPeriksa struct{ Hasil HasilPeriksa }

func (g GalatPeriksa) Error() string { return strings.Join(g.Hasil.Galat, "\n") }

// Is - GalatPeriksa adalah penolakan.
func (g GalatPeriksa) Is(t error) bool { return t == ErrDitolak }

// BatasCatatan - panjang CATATAN (VARCHAR2(500)).
const BatasCatatan = 500

// RingkasanSlot adalah satu baris layar Template Manager.
type RingkasanSlot struct {
	Kode        string `json:"kode"`
	Menu        string `json:"menu"`
	Grup        string `json:"grup"`
	Nama        string `json:"nama"`
	DipakaiDi   string `json:"dipakaiDi"`
	Ekstensi    string `json:"ekstensi"`
	Pemisah     string `json:"pemisah"`
	JumlahKolom int    `json:"jumlahKolom"`
	NamaUnduhan string `json:"namaUnduhan"`
	// Aktif - versi aktif; nil = berkas bawaan.
	Aktif *Versi `json:"aktif"`
}

// Layanan adalah aturan Template Manager di atas katalog dan gudang.
type Layanan struct {
	katalog *Katalog
	gudang  Gudang
}

// NewLayanan membuat layanan; gudang nil = tanpa database (hanya berkas bawaan, unggah ditolak).
func NewLayanan(k *Katalog, g Gudang) *Layanan { return &Layanan{katalog: k, gudang: g} }

// ErrTanpaDatabase - unggah dan riwayat butuh Oracle.
var ErrTanpaDatabase = errors.New("templat: Template Manager butuh database (ORACLE_DSN)")

// Slot - satu slot menurut kode.
func (l *Layanan) Slot(kode string) (Slot, error) { return l.katalog.Ambil(kode) }

// Daftar - seluruh slot beserta versi aktifnya.
func (l *Layanan) Daftar(ctx context.Context) ([]RingkasanSlot, error) {
	aktif := map[string]Versi{}
	if l.gudang != nil {
		var err error
		if aktif, err = l.gudang.Aktif(ctx); err != nil {
			return nil, err
		}
	}
	hasil := []RingkasanSlot{}
	for _, s := range l.katalog.Semua() {
		r := RingkasanSlot{Kode: s.Kode, Menu: s.Menu, Grup: s.Grup, Nama: s.Nama, DipakaiDi: s.DipakaiDi,
			Ekstensi: s.Ekstensi, JumlahKolom: s.JumlahKolom, NamaUnduhan: s.NamaUnduhan}
		if s.Pemisah != 0 {
			r.Pemisah = string(s.Pemisah)
		}
		if v, ok := aktif[s.Kode]; ok {
			r.Aktif = &v
		}
		hasil = append(hasil, r)
	}
	return hasil, nil
}

// Riwayat - seluruh versi satu slot, terbaru dulu.
func (l *Layanan) Riwayat(ctx context.Context, kode string) ([]Versi, error) {
	if _, err := l.katalog.Ambil(kode); err != nil {
		return nil, err
	}
	if l.gudang == nil {
		return []Versi{}, nil
	}
	vs, err := l.gudang.Riwayat(ctx, kode)
	if vs == nil {
		vs = []Versi{}
	}
	return vs, err
}

// isiAktif - isi versi aktif, atau berkas bawaan.
func (l *Layanan) isiAktif(ctx context.Context, s Slot) (Versi, []byte, error) {
	if l.gudang == nil {
		return Versi{}, s.Bawaan, nil
	}
	aktif, err := l.gudang.Aktif(ctx)
	if err != nil {
		return Versi{}, nil, err
	}
	v, ok := aktif[s.Kode]
	if !ok {
		return Versi{}, s.Bawaan, nil
	}
	isi, err := l.gudang.Isi(ctx, s.Kode, v.Versi)
	return v, isi, err
}

// Periksa menilai berkas tanpa menyimpannya.
func (l *Layanan) Periksa(ctx context.Context, kode, nama string, isi []byte) (HasilPeriksa, error) {
	s, err := l.katalog.Ambil(kode)
	if err != nil {
		return HasilPeriksa{}, err
	}
	_, aktif, err := l.isiAktif(ctx, s)
	if err != nil {
		return HasilPeriksa{}, err
	}
	return periksaBerkas(s, nama, isi, aktif), nil
}

// Unggah menyimpan berkas sebagai versi baru yang langsung aktif.
func (l *Layanan) Unggah(ctx context.Context, kode, nama string, isi []byte, catatan, akun string) (int, error) {
	if l.gudang == nil {
		return 0, ErrTanpaDatabase
	}
	if strings.TrimSpace(akun) == "" {
		return 0, errors.New("templat: pengunggah tanpa identitas")
	}
	catatan = strings.TrimSpace(catatan)
	if utf8.RuneCountInString(catatan) > BatasCatatan || len(catatan) > BatasCatatan*4 {
		return 0, GalatPeriksa{HasilPeriksa{Galat: []string{"Note is longer than 500 characters"}}}
	}
	nama = path.Base(strings.ReplaceAll(strings.TrimSpace(nama), `\`, "/"))
	if len(nama) > 255 {
		return 0, GalatPeriksa{HasilPeriksa{Galat: []string{"File name is longer than 255 characters"}}}
	}
	h, err := l.Periksa(ctx, kode, nama, isi)
	if err != nil {
		return 0, err
	}
	if len(h.Galat) > 0 {
		return 0, GalatPeriksa{h}
	}
	return l.gudang.Sisip(ctx, kode, Versi{NamaBerkas: nama, JumlahKolom: h.JumlahKolom, Catatan: catatan, DiunggahOleh: akun}, isi)
}

// Aktifkan menjadikan satu versi aktif; 0 = kembali ke berkas bawaan.
func (l *Layanan) Aktifkan(ctx context.Context, kode string, versi int) error {
	if _, err := l.katalog.Ambil(kode); err != nil {
		return err
	}
	if l.gudang == nil {
		return ErrTanpaDatabase
	}
	if versi < 0 {
		return GalatPeriksa{HasilPeriksa{Galat: []string{"Version is not valid"}}}
	}
	return l.gudang.Aktifkan(ctx, kode, versi)
}

// Berkas adalah satu unduhan.
type Berkas struct {
	Nama string
	Mime string
	Isi  []byte
}

// mimeTetap - jenis berkas templat; TIDAK dari `mime.TypeByExtension`, yang di Windows membaca registry
// (`.csv` menjadi application/vnd.ms-excel) sehingga jawabannya berbeda antar mesin.
var mimeTetap = map[string]string{
	".csv":  "text/csv; charset=utf-8",
	".txt":  "text/plain; charset=utf-8",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".xls":  "application/vnd.ms-excel",
}

func mimeDari(nama string) string {
	if m, ok := mimeTetap[strings.ToLower(path.Ext(nama))]; ok {
		return m
	}
	return "application/octet-stream"
}

// Unduh - versi aktif (versi = -1), berkas bawaan (versi = 0), atau satu versi tertentu. Nama unduhan versi aktif
// selalu NamaUnduhan slot, apa pun nama berkas yang diunggah.
func (l *Layanan) Unduh(ctx context.Context, kode string, versi int) (Berkas, error) {
	s, err := l.katalog.Ambil(kode)
	if err != nil {
		return Berkas{}, err
	}
	switch {
	case versi < 0:
		_, isi, err := l.isiAktif(ctx, s)
		if err != nil {
			return Berkas{}, err
		}
		return Berkas{Nama: s.NamaUnduhan, Mime: mimeDari(s.NamaUnduhan), Isi: isi}, nil
	case versi == 0:
		return Berkas{Nama: s.NamaUnduhan, Mime: mimeDari(s.NamaUnduhan), Isi: s.Bawaan}, nil
	}
	if l.gudang == nil {
		return Berkas{}, ErrVersiTidakAda
	}
	isi, err := l.gudang.Isi(ctx, kode, versi)
	if err != nil {
		return Berkas{}, err
	}
	nama := strings.TrimSuffix(s.NamaUnduhan, s.Ekstensi) + " v" + strconv.Itoa(versi) + s.Ekstensi
	return Berkas{Nama: nama, Mime: mimeDari(nama), Isi: isi}, nil
}
