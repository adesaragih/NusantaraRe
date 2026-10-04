package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/accounts/backend/models"
	"nusantarare/modul/accounts/backend/repository"
)

// Batas lebar dan ukuran.
const (
	// LebarDescription - `DESCRIPTION` VARCHAR2(4000) BYTE (migrasi 840).
	LebarDescription = 4000
	// LebarCreateOp - `CREATEOP` VARCHAR2(100) BYTE (migrasi 840).
	LebarCreateOp = 100
	// UkuranBawaan dan UkuranMaks - baris per halaman daftar.
	UkuranBawaan = 20
	UkuranMaks   = 100
	// BatasOrganisasi - pilihan "Insured Name" per pencarian (CLIENT memuat puluhan ribu organisasi).
	BatasOrganisasi = 50
	// cobaNomor - nomor sequence yang sudah terpakai dilewati paling banyak sekian kali.
	cobaNomor = 50
)

// Pesan wajib isi - VERBATIM tangkapan layar Pega ("Value cannot be blank").
const PesanWajib = "Value cannot be blank"

// Galat layanan. Pesan untuk layar (bahasa Inggris) mengikuti sesudah `: `.
var (
	// ErrMasukanTidakSah - isian ditolak (422).
	ErrMasukanTidakSah = errors.New("invalid input")
	// ErrGanda - pasangan Insured + Group Business sudah dimiliki akun yang ada (409).
	ErrGanda = errors.New("duplicate account")
	// ErrTanpaPelaku - permintaan tanpa akun pelaku (401).
	ErrTanpaPelaku = errors.New("user account is required")
	// ErrNomorHabis - setiap nomor sequence yang dicoba sudah terpakai (500 berkalimat).
	ErrNomorHabis = errors.New("SEQ_T_M_ACCOUNT keeps returning account numbers that are already used")
)

func tolak(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrMasukanTidakSah, fmt.Sprintf(format, a...))
}

// Transaksi menjalankan fn di dalam satu transaksi dan menutupnya.
type Transaksi func(ctx context.Context, fn func(tx *dbTx) error) error

// Gudang adalah semua yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	Daftar(ctx context.Context, s models.Saringan) ([]models.Account, int, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Account, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	NomorBerikut(ctx context.Context, tx *dbTx) (int64, error)
	Sisip(ctx context.Context, tx *dbTx, a models.Account) error
	PasanganLain(ctx context.Context, tx *dbTx, insuredID, groupBusinessID string) (string, error)
	CariOrganisasi(ctx context.Context, kata string, batas int) ([]models.Organisasi, error)
	AmbilOrganisasi(ctx context.Context, tx *dbTx, id string) (models.Organisasi, error)
	DaftarGroupBusiness(ctx context.Context) ([]models.GroupBusiness, error)
	AmbilGroupBusiness(ctx context.Context, tx *dbTx, id string) (models.GroupBusiness, error)
}

// Layanan memegang seluruh aturan modul ini di atas satu Gudang.
type Layanan struct {
	gudang Gudang
	tx     Transaksi
}

// BaruLayanan menyusun Layanan - dipakai uji dengan gudang tiruan dan transaksi tiruan (fn(nil)).
func BaruLayanan(g Gudang, tx Transaksi) *Layanan { return &Layanan{gudang: g, tx: tx} }

// HalamanDaftar adalah satu halaman daftar akun.
type HalamanDaftar struct {
	Daftar  []models.Account `json:"daftar"`
	Total   int              `json:"total"`
	Halaman int              `json:"halaman"`
	Ukuran  int              `json:"ukuran"`
}

// HasilOrganisasi adalah pilihan "Insured Name" untuk satu kata cari; `Lebih` = ada organisasi lain yang cocok.
type HasilOrganisasi struct {
	Daftar []models.Organisasi `json:"daftar"`
	Lebih  bool                `json:"lebih"`
}

// Pilihan adalah isi pilihan form yang kecil (dimuat sekaligus).
type Pilihan struct {
	GroupBusiness []models.GroupBusiness `json:"groupBusiness"`
}

// Daftar - satu halaman akun, terbaru lebih dulu; `kata` mencari Insured Name, Group Business, ACC-n, dan ORG-n.
func (l *Layanan) Daftar(ctx context.Context, kata string, halaman, ukuran int) (HalamanDaftar, error) {
	if halaman < 1 {
		halaman = 1
	}
	if ukuran < 1 || ukuran > UkuranMaks {
		ukuran = UkuranBawaan
	}
	d, total, err := l.gudang.Daftar(ctx, models.Saringan{Cari: strings.TrimSpace(kata), Halaman: halaman, Ukuran: ukuran})
	if err != nil {
		return HalamanDaftar{}, err
	}
	return HalamanDaftar{Daftar: d, Total: total, Halaman: halaman, Ukuran: ukuran}, nil
}

// Pilihan - daftar Group Business.
func (l *Layanan) Pilihan(ctx context.Context) (Pilihan, error) {
	gb, err := l.gudang.DaftarGroupBusiness(ctx)
	if err != nil {
		return Pilihan{}, err
	}
	return Pilihan{GroupBusiness: gb}, nil
}

// CariOrganisasi - pilihan "Insured Name" untuk satu kata cari (kosong = urutan nama teratas).
func (l *Layanan) CariOrganisasi(ctx context.Context, kata string) (HasilOrganisasi, error) {
	d, err := l.gudang.CariOrganisasi(ctx, strings.TrimSpace(kata), BatasOrganisasi+1)
	if err != nil {
		return HasilOrganisasi{}, err
	}
	h := HasilOrganisasi{Daftar: d}
	if len(d) > BatasOrganisasi {
		h.Daftar, h.Lebih = d[:BatasOrganisasi], true
	}
	return h, nil
}

// periksaPelaku - akun pelaku ditulis ke `CREATEOP` VARCHAR2(100).
func periksaPelaku(p inti.Pelaku) error {
	if strings.TrimSpace(p.AkunID) == "" {
		return ErrTanpaPelaku
	}
	if len(p.AkunID) > LebarCreateOp {
		return tolak("your user account is longer than %d bytes", LebarCreateOp)
	}
	return nil
}

// isi - isian yang sudah diperiksa: organisasi dan Group Business dari tabelnya, Description dirapikan.
func (l *Layanan) isi(ctx context.Context, tx *dbTx, isian models.Isian) (models.Account, error) {
	insID, gbID := strings.TrimSpace(isian.InsuredID), strings.TrimSpace(isian.GroupBusinessID)
	desk := strings.TrimSpace(isian.Description)
	var kosong []string
	if insID == "" {
		kosong = append(kosong, "Insured Name: "+PesanWajib)
	}
	if gbID == "" {
		kosong = append(kosong, "Group Business: "+PesanWajib)
	}
	if len(kosong) > 0 {
		return models.Account{}, tolak("%s", strings.Join(kosong, "; "))
	}
	if len(desk) > LebarDescription {
		return models.Account{}, tolak("Description is longer than %d bytes", LebarDescription)
	}
	org, err := l.gudang.AmbilOrganisasi(ctx, tx, insID)
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Account{}, tolak("Insured Name: the organization is not in the organization list")
	}
	if err != nil {
		return models.Account{}, err
	}
	gb, err := l.gudang.AmbilGroupBusiness(ctx, tx, gbID)
	if errors.Is(err, repository.ErrTidakAda) {
		return models.Account{}, tolak("Group Business: not in the business group list")
	}
	if err != nil {
		return models.Account{}, err
	}
	return models.Account{InsuredID: org.ID, OrgID: models.BagianAkhir(org.ID, models.AwalanOrg),
		InsuredName: org.Nama, GroupBusinessID: gb.ID, GroupBusiness: gb.Note, Description: desk}, nil
}

// tolakGanda - pasangan Insured + Group Business yang sudah dimiliki akun yang ada.
func (l *Layanan) tolakGanda(ctx context.Context, tx *dbTx, a models.Account) error {
	lain, err := l.gudang.PasanganLain(ctx, tx, a.InsuredID, a.GroupBusinessID)
	if err != nil || lain == "" {
		return err
	}
	return fmt.Errorf("%w: %s already has an account in %s (%s)", ErrGanda, a.InsuredName, a.GroupBusiness,
		models.BagianAkhir(lain, models.AwalanID))
}

// idBaru - `ACC-<SEQ_T_M_ACCOUNT>`; nomor yang sudah terpakai (mis. dibuat sebelum migrasi 842) dilewati.
func (l *Layanan) idBaru(ctx context.Context, tx *dbTx) (string, error) {
	for i := 0; i < cobaNomor; i++ {
		n, err := l.gudang.NomorBerikut(ctx, tx)
		if err != nil {
			return "", err
		}
		id := models.AwalanID + models.AwalanIDView + strconv.FormatInt(n, 10)
		ada, err := l.gudang.AdaID(ctx, tx, id)
		if err != nil {
			return "", err
		}
		if !ada {
			return id, nil
		}
	}
	return "", ErrNomorHabis
}

// Tambah - akun baru: Owner = akun pelaku, Create Date = SYSDATE. Satu-satunya tulisan modul ini: akun hanya
// diinput sekali - nol ubah, nol hapus (keputusan work owner 04-10-2026).
func (l *Layanan) Tambah(ctx context.Context, p inti.Pelaku, isian models.Isian) (models.Account, error) {
	if err := periksaPelaku(p); err != nil {
		return models.Account{}, err
	}
	var id string
	err := l.tx(ctx, func(tx *dbTx) error {
		a, err := l.isi(ctx, tx, isian)
		if err != nil {
			return err
		}
		if err := l.tolakGanda(ctx, tx, a); err != nil {
			return err
		}
		if a.ID, err = l.idBaru(ctx, tx); err != nil {
			return err
		}
		a.CreateOp = strings.TrimSpace(p.AkunID)
		id = a.ID
		return l.gudang.Sisip(ctx, tx, a)
	})
	if err != nil {
		return models.Account{}, err
	}
	return l.gudang.Ambil(ctx, nil, id)
}
