package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/aggregate/backend/models"
	"nusantarare/modul/aggregate/backend/repository"
)

// Ukuran halaman daftar.
const (
	UkuranHalaman  = 50
	UkuranMaksimum = 200
	// BatasCariTreaty - baris view Master ID paling banyak sekali cari (Obj-Browse Pega tanpa MaxRecords: 10.000).
	BatasCariTreaty = 1000
)

// Galat layanan. Pesan untuk layar (bahasa Inggris) mengikuti sesudah `: `.
var (
	// ErrTidakAda - baris yang dicari tidak ada (404).
	ErrTidakAda = errors.New("aggregate data not found")
	// ErrMasukanTidakSah - isian ditolak (422).
	ErrMasukanTidakSah = errors.New("invalid input")
	// ErrTanpaPelaku - permintaan tanpa akun pelaku (401).
	ErrTanpaPelaku = errors.New("user account is required")
	// ErrBelumDimigrasi - migrasi 880 belum dijalankan (503).
	ErrBelumDimigrasi = repository.ErrBelumDimigrasi
)

func tolak(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrMasukanTidakSah, fmt.Sprintf(format, a...))
}

// GalatSimpan - pesan penolakan Save satu per baris layar (`SaveAggregate_Act`: `Local.Err`, dipisah `<br>`).
type GalatSimpan struct{ Pesan []string }

func (g GalatSimpan) Error() string { return strings.Join(g.Pesan, "\n") }

// Is - GalatSimpan adalah ErrMasukanTidakSah (422).
func (g GalatSimpan) Is(target error) bool { return target == ErrMasukanTidakSah }

// Transaksi menjalankan fn di dalam satu transaksi dan menutupnya.
type Transaksi func(ctx context.Context, fn func(tx *dbTx) error) error

// Gudang adalah semua yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	Daftar(ctx context.Context, kueri string, offset, ukuran int) ([]models.Kelompok, int, error)
	Rincian(ctx context.Context, k models.Kunci) ([]models.Baris, error)
	Hapus(ctx context.Context, tx *dbTx, k models.Kunci) (int64, error)
	Ringkasan(ctx context.Context, asAt string) ([]models.IrisanRingkasan, error)
	CariTreaty(ctx context.Context, kueri string, batas int) ([]models.MasterTreaty, error)
	AmbilTreaty(ctx context.Context, id string) (models.MasterTreaty, error)
	DaftarZona(ctx context.Context) (map[string]string, error)
	TahunTreaty(ctx context.Context, asAt string) (string, error)
	Kurs(ctx context.Context, tahun, mataUang string) (toUSD, toIDR string, ada bool, err error)
	Sekarang(ctx context.Context, tx *dbTx) (string, error)
	NomorBerikut(ctx context.Context, tx *dbTx) (int64, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, id, sekarang, pelaku string, b models.Baris) error
}

// Layanan memegang seluruh aturan modul ini di atas satu Gudang.
type Layanan struct {
	gudang Gudang
	tx     Transaksi
}

// BaruLayanan menyusun Layanan - dipakai uji dengan gudang tiruan dan transaksi tiruan (fn(nil)).
func BaruLayanan(g Gudang, tx Transaksi) *Layanan { return &Layanan{gudang: g, tx: tx} }

func periksaPelaku(p inti.Pelaku) error {
	if strings.TrimSpace(p.AkunID) == "" {
		return ErrTanpaPelaku
	}
	if len(p.AkunID) > models.LebarUserInput {
		return tolak("your user account is longer than %d bytes", models.LebarUserInput)
	}
	return nil
}

// tanggalSah - teks kosong atau tanggal `DD-MM-YYYY`.
func tanggalSah(s string) bool {
	if s == "" {
		return true
	}
	_, err := time.Parse("02-01-2006", s)
	return err == nil
}

func rapikanKunci(k models.Kunci) (models.Kunci, error) {
	k = models.Kunci{TanggalInput: strings.TrimSpace(k.TanggalInput), CedingCode: strings.TrimSpace(k.CedingCode),
		CedingName: strings.TrimSpace(k.CedingName), TreatyType: strings.TrimSpace(k.TreatyType),
		AsAt: strings.TrimSpace(k.AsAt), UwYear: strings.TrimSpace(k.UwYear)}
	if !tanggalSah(k.TanggalInput) || !tanggalSah(k.AsAt) {
		return k, tolak("Tanggal Input and As At must be DD-MM-YYYY")
	}
	return k, nil
}

// Halaman adalah satu halaman daftar.
type Halaman struct {
	Daftar  []models.Kelompok `json:"daftar"`
	Total   int               `json:"total"`
	Halaman int               `json:"halaman"`
	Ukuran  int               `json:"ukuran"`
}

// Daftar membaca satu halaman daftar `GridDasbordAgg`; kueri mencari Ceding Name dan Ceding Code.
func (l *Layanan) Daftar(ctx context.Context, kueri string, halaman, ukuran int) (Halaman, error) {
	if ukuran <= 0 {
		ukuran = UkuranHalaman
	}
	if ukuran > UkuranMaksimum {
		ukuran = UkuranMaksimum
	}
	if halaman < 1 {
		halaman = 1
	}
	baris, total, err := l.gudang.Daftar(ctx, strings.TrimSpace(kueri), (halaman-1)*ukuran, ukuran)
	if err != nil {
		return Halaman{}, err
	}
	return Halaman{Daftar: baris, Total: total, Halaman: halaman, Ukuran: ukuran}, nil
}

// Rincian - seluruh baris satu baris daftar (klik ganda `GridDasbordAgg`).
func (l *Layanan) Rincian(ctx context.Context, k models.Kunci) ([]models.Baris, error) {
	k, err := rapikanKunci(k)
	if err != nil {
		return nil, err
	}
	baris, err := l.gudang.Rincian(ctx, k)
	if err != nil {
		return nil, err
	}
	if len(baris) == 0 {
		return nil, ErrTidakAda
	}
	return baris, nil
}

// Hapus membuang seluruh baris satu baris daftar (Delete + ConfrimDelete `GridDasbordAgg`).
func (l *Layanan) Hapus(ctx context.Context, p inti.Pelaku, k models.Kunci) (int64, error) {
	if err := periksaPelaku(p); err != nil {
		return 0, err
	}
	k, err := rapikanKunci(k)
	if err != nil {
		return 0, err
	}
	var n int64
	err = l.tx(ctx, func(tx *dbTx) error {
		var err error
		n, err = l.gudang.Hapus(ctx, tx, k)
		if err == nil && n == 0 {
			return ErrTidakAda
		}
		return err
	})
	return n, err
}

// Ringkasan adalah isi chart bertingkat Ceding > Treaty Type > Coverage: RNM Value (USD) per daun, SELURUH As At
// dijumlah (keputusan work owner 04-10-2026: As At dibuang dari chart).
type Ringkasan struct {
	Irisan []models.IrisanRingkasan `json:"irisan"`
}

// Ringkasan - isi chart, seluruh As At.
func (l *Layanan) Ringkasan(ctx context.Context) (Ringkasan, error) {
	irisan, err := l.gudang.Ringkasan(ctx, "")
	if err != nil {
		return Ringkasan{}, err
	}
	if irisan == nil {
		irisan = []models.IrisanRingkasan{}
	}
	return Ringkasan{Irisan: irisan}, nil
}

// kunciTreaty - kunci buang-ganda `GetMasterIDAgg_Act` langkah 4.
func kunciTreaty(t models.MasterTreaty) string {
	return strings.Join([]string{t.TreatyID, t.CedingID, t.RnmShare, t.TreatyYear, t.TreatyGroup}, "#")
}

// CariTreaty - popup Master ID (`GetMasterIDAgg_Act`): kata kosong = nol baris (langkah 2 Exit-Activity); baris
// kembar menurut TREATYID, CEDINGID, RNM_SHARE, TREATYYEAR, TREATYGROUP dibuang - yang tersisa kemunculan terakhir,
// seperti langkah Java Pega yang menelusuri daftar dari belakang.
func (l *Layanan) CariTreaty(ctx context.Context, kueri string) ([]models.MasterTreaty, error) {
	kueri = strings.TrimSpace(kueri)
	if kueri == "" {
		return []models.MasterTreaty{}, nil
	}
	semua, err := l.gudang.CariTreaty(ctx, kueri, BatasCariTreaty)
	if err != nil {
		return nil, err
	}
	dilihat := map[string]bool{}
	simpan := make([]bool, len(semua))
	for i := len(semua) - 1; i >= 0; i-- {
		k := kunciTreaty(semua[i])
		if !dilihat[k] {
			dilihat[k] = true
			simpan[i] = true
		}
	}
	out := []models.MasterTreaty{}
	for i, t := range semua {
		if simpan[i] {
			out = append(out, t)
		}
	}
	return out, nil
}
