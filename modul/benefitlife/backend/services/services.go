// Package services memuat aturan modul Benefit (`benefitlife`; panduan XML section `InboxBenefit`; keputusan work
// owner 08-10-2026 K1-K6, MODUL.md). Pola ririsklife, versi TANPA ringkasan / rincian / unggah / Delete (XML tidak
// memuatnya).
//
//   - Grid `BrowseBenefitLife_RD` (b4451): kolom ID (`.Number` b3858) dan Benefit (`.Benefit` b4020), 10 baris per
//     halaman (b4526), ID MENURUN bawaan (b4391/b4397), saring (b4505).
//   - Save (b1772 -> `AddToList_Act` b1791): Add bila ID kosong, selain itu Edit (Edit b4220 -> `EditList_DT` b4243
//     mengisi form dengan Number dan Benefit). Cancel (b2035 -> `NewData_DT` b2058, hanya saat `DATASHOW = 'IsEdit'`
//     b2198) murni layar.
//   - ID baru = '1' || LPAD(M_BENEFIT_LIFE_SEQ.NEXTVAL, 5, '0') (models.BentukID, K3). NEXTVAL yang menghasilkan ID
//     yang SUDAH ada = galat ErrIDTerpakai yang jelas (409), bukan 500 dan bukan lompat diam-diam.
//   - Benefit wajib dan huruf besar (models.NormalBenefit, K6).
//   - Transaksi milik Go: satu transaksi per simpan; prosedur PEGA_M_BENEFIT_LIFE tidak dipanggil (INVALID sesudah 944,
//     K3 diterima WO).
//   - Hak menu Full / View only: View only tanpa Save (403).
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/benefitlife/backend/models"
	"nusantarare/modul/benefitlife/backend/repository"
)

// Service membawa akar bersama.
type Service struct{ *inti.Dasar }

// DariDasar membungkus akar bersama.
func DariDasar(d *inti.Dasar) *Service { return &Service{Dasar: d} }

// PunyaDatabase - Oracle terpasang?
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar.PunyaDatabase() }

// LayananOracle merakit layanan di atas Oracle.
func LayananOracle(s *Service) *Layanan {
	return BaruLayanan(repository.Baru(s.DB()), s.DalamTransaksi)
}

var _ Gudang = (*repository.Gudang)(nil)

type dbTx = db.Tx

var (
	// ErrTidakAda - ID tidak ada.
	ErrTidakAda = errors.New("services: benefit tidak ada")
	// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrIDTerpakai - K3: ID dari M_BENEFIT_LIFE_SEQ sudah ada di BENEFIT_LIFE (atau nomornya melewati LPAD 5); pesannya
	// untuk pengguna / admin.
	ErrIDTerpakai = errors.New("services: ID baru dari sequence tidak dapat dipakai")
	// ErrBelumAda - tabel, kolom, atau sequence tidak ada di skema ini.
	ErrBelumAda = repository.ErrBelumAda
)

func tolak(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrMasukanTidakSah, fmt.Sprintf(format, a...))
}

func dilarang(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrDilarang, fmt.Sprintf(format, a...))
}

func idTerpakai(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrIDTerpakai, fmt.Sprintf(format, a...))
}

// Transaksi menjalankan fn di dalam satu transaksi.
type Transaksi func(ctx context.Context, fn func(tx *dbTx) error) error

// Gudang - yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	Daftar(ctx context.Context, s models.Saringan) ([]models.Benefit, int, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Benefit, error)
	NomorBaru(ctx context.Context, tx *dbTx) (string, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, b models.Benefit) error
	Ubah(ctx context.Context, tx *dbTx, b models.Benefit) error
}

// Layanan adalah aturan modul di atas Gudang.
type Layanan struct {
	gudang Gudang
	tx     Transaksi
}

// BaruLayanan membuat layanan.
func BaruLayanan(g Gudang, tx Transaksi) *Layanan { return &Layanan{gudang: g, tx: tx} }

// Aktor - pelaku permintaan: akun login dan hak menunya (Full = bukan View only).
type Aktor struct {
	AkunID string
	Penuh  bool
}

func wajibPenuh(a Aktor) error {
	if a.AkunID == "" {
		return dilarang("Log in to save")
	}
	if !a.Penuh {
		return dilarang("Your access to the Benefit menu is View only")
	}
	return nil
}

// RapikanSaringan - pangkas saringan, halaman >= 1.
func RapikanSaringan(s models.Saringan) models.Saringan {
	s.ID, s.Benefit = strings.TrimSpace(s.ID), strings.TrimSpace(s.Benefit)
	if s.Halaman < 1 {
		s.Halaman = 1
	}
	return s
}

// Daftar - satu halaman grid.
func (l *Layanan) Daftar(ctx context.Context, s models.Saringan) (models.Halaman, error) {
	s = RapikanSaringan(s)
	d, total, err := l.gudang.Daftar(ctx, s)
	if err != nil {
		return models.Halaman{}, err
	}
	if d == nil {
		d = []models.Benefit{}
	}
	return models.Halaman{Daftar: d, Total: total, Halaman: s.Halaman, Ukuran: models.UkuranHalaman}, nil
}

func petaGalat(err error) error {
	if errors.Is(err, repository.ErrTidakAda) {
		return ErrTidakAda
	}
	return err
}

// idBaru - K3: '1' || LPAD(NEXTVAL, 5, '0'). ID yang sudah ada (mis. ditulis prosedur Pega dengan sequence yang
// tertinggal) = ErrIDTerpakai yang menyebut ID dan sequence-nya - TIDAK dilompati diam-diam (keputusan work owner).
func (l *Layanan) idBaru(ctx context.Context, tx *dbTx) (string, error) {
	nomor, err := l.gudang.NomorBaru(ctx, tx)
	if err != nil {
		return "", err
	}
	id, err := models.BentukID(nomor)
	if err != nil {
		return "", idTerpakai("M_BENEFIT_LIFE_SEQ gave number %s, which does not fit the ID formula '1' || LPAD(number, 5); ask the administrator", strings.TrimSpace(nomor))
	}
	ada, err := l.gudang.AdaID(ctx, tx, id)
	if err != nil {
		return "", err
	}
	if ada {
		return "", idTerpakai("new ID %s from M_BENEFIT_LIFE_SEQ already exists; ask the administrator to move the sequence past the highest ID", id)
	}
	return id, nil
}

// Simpan - Save (`AddToList_Act` b1791): Add bila id kosong (ID dari sequence, K3), selain itu Edit (`EditList_DT`
// b4243: Number = ID dari jalur, Benefit). Satu transaksi.
func (l *Layanan) Simpan(ctx context.Context, a Aktor, id string, isi models.Isian) (models.Benefit, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Benefit{}, err
	}
	id = strings.TrimSpace(id)
	benefit, err := models.NormalBenefit(isi.Benefit)
	if err != nil {
		return models.Benefit{}, tolak("%s", err.Error())
	}
	err = l.tx(ctx, func(tx *dbTx) error {
		if id != "" {
			if _, err := l.gudang.Ambil(ctx, tx, id); err != nil {
				return err
			}
			return l.gudang.Ubah(ctx, tx, models.Benefit{ID: id, Benefit: benefit})
		}
		baru, err := l.idBaru(ctx, tx)
		if err != nil {
			return err
		}
		if err := l.gudang.Sisip(ctx, tx, models.Benefit{ID: baru, Benefit: benefit}); err != nil {
			if errors.Is(err, repository.ErrKembar) {
				// Balapan dengan penulis lain sesudah AdaID (PK SYS_C009031 menolak): galat yang sama jelasnya.
				return idTerpakai("new ID %s from M_BENEFIT_LIFE_SEQ already exists; ask the administrator to move the sequence past the highest ID", baru)
			}
			return err
		}
		id = baru
		return nil
	})
	if err != nil {
		return models.Benefit{}, petaGalat(err)
	}
	b, err := l.gudang.Ambil(ctx, nil, id)
	return b, petaGalat(err)
}
