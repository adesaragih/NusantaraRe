// Package services memuat aturan modul Cause Of Loss Life (`causeoflosslife`; panduan XML section
// `InboxCauseofLossLife`; keputusan work owner 08-10-2026 K0-K6, MODUL.md). Pola benefitlife (satu tabel, satu layar),
// TANPA saring / urut pilihan / Delete / Upload (XML tidak memuatnya).
//
//   - Grid `BrowseCauseofLossLife_RD` (b3981): kolom ID (`.ID` b3415) dan Cause of Loss (`.CauseofLoss` b3562), 10
//     baris per halaman (b4057), ID MENAIK tetap (b3923 / b3929), tanpa saring (b4036).
//   - Save (b5366 -> `AddToList_Act` b5385): Add bila ID kosong, selain itu Edit (Edit b3747 -> `EditList_DT` b3770
//     mengisi form dengan ID dan CauseofLoss). Cancel (b5624 -> `NewData_DT` b5647, hanya saat `DATASHOW = 'IsEdit'`
//     b5787) murni layar.
//   - ID baru = '1' || LPAD(M_CAUSEOFLOSS_LIFE_SEQ.NEXTVAL, 5, '0') (models.BentukID, K3). NEXTVAL yang menghasilkan
//     ID yang SUDAH ada = galat ErrIDTerpakai yang jelas (409), bukan 500 dan bukan lompat diam-diam.
//   - K4: Cause of Loss wajib (XML); huruf TIDAK diubah (XML tanpa pengubah huruf); dipangkas dan tidak boleh kembar
//     tanpa beda huruf = di luar XML, pola modul MASTER TREATY (PARITAS, WO boleh menolak).
//   - Transaksi milik Go: satu transaksi per simpan; prosedur PEGA_M_CAUSEOFLOSS_LIFE tidak dipanggil (INVALID sesudah
//     092, K3 diterima WO).
//   - Hak menu Full / View only: View only tanpa Save (403).
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/causeoflosslife/backend/models"
	"nusantarare/modul/causeoflosslife/backend/repository"
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
	ErrTidakAda = errors.New("services: cause of loss tidak ada")
	// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrIDTerpakai - K3: ID dari M_CAUSEOFLOSS_LIFE_SEQ sudah ada di CAUSEOFLOSS_LIFE (atau nomornya melewati LPAD 5);
	// pesannya untuk pengguna / admin.
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

func idTerpakai(id string) error {
	return fmt.Errorf("%w: new ID %s from M_CAUSEOFLOSS_LIFE_SEQ already exists; ask the administrator to move the sequence past the highest ID",
		ErrIDTerpakai, id)
}

// Transaksi menjalankan fn di dalam satu transaksi.
type Transaksi func(ctx context.Context, fn func(tx *dbTx) error) error

// Gudang - yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	Daftar(ctx context.Context, s models.Saringan) ([]models.CauseOfLoss, int, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.CauseOfLoss, error)
	PemakaiNama(ctx context.Context, tx *dbTx, nama, kecualiID string) ([]models.CauseOfLoss, error)
	NomorBaru(ctx context.Context, tx *dbTx) (string, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, c models.CauseOfLoss) error
	Ubah(ctx context.Context, tx *dbTx, c models.CauseOfLoss) error
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
		return dilarang("Your access to the Cause Of Loss Life menu is View only")
	}
	return nil
}

// Daftar - satu halaman grid (halaman >= 1).
func (l *Layanan) Daftar(ctx context.Context, s models.Saringan) (models.Halaman, error) {
	if s.Halaman < 1 {
		s.Halaman = 1
	}
	d, total, err := l.gudang.Daftar(ctx, s)
	if err != nil {
		return models.Halaman{}, err
	}
	if d == nil {
		d = []models.CauseOfLoss{}
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
		return "", fmt.Errorf("%w: M_CAUSEOFLOSS_LIFE_SEQ gave number %s, which does not fit the ID formula '1' || LPAD(number, 5); ask the administrator",
			ErrIDTerpakai, strings.TrimSpace(nomor))
	}
	ada, err := l.gudang.AdaID(ctx, tx, id)
	if err != nil {
		return "", err
	}
	if ada {
		return "", idTerpakai(id)
	}
	return id, nil
}

// Simpan - Save (`AddToList_Act` b5385): Add bila id kosong (ID dari sequence, K3), selain itu Edit (`EditList_DT`
// b3770: ID dari jalur, CauseofLoss). Satu transaksi.
func (l *Layanan) Simpan(ctx context.Context, a Aktor, id string, isi models.Isian) (models.CauseOfLoss, error) {
	if err := wajibPenuh(a); err != nil {
		return models.CauseOfLoss{}, err
	}
	id = strings.TrimSpace(id)
	nama, err := models.NormalCauseOfLoss(isi.CauseOfLoss)
	if err != nil {
		return models.CauseOfLoss{}, tolak("%s", err.Error())
	}
	err = l.tx(ctx, func(tx *dbTx) error {
		if id != "" {
			if _, err := l.gudang.Ambil(ctx, tx, id); err != nil {
				return err
			}
		}
		ganda, err := l.gudang.PemakaiNama(ctx, tx, nama, id)
		if err != nil {
			return err
		}
		if len(ganda) > 0 {
			return tolak("Cause of Loss %s is already used by ID %s", nama, ganda[0].ID)
		}
		if id != "" {
			return l.gudang.Ubah(ctx, tx, models.CauseOfLoss{ID: id, CauseOfLoss: nama})
		}
		baru, err := l.idBaru(ctx, tx)
		if err != nil {
			return err
		}
		if err := l.gudang.Sisip(ctx, tx, models.CauseOfLoss{ID: baru, CauseOfLoss: nama}); err != nil {
			if errors.Is(err, repository.ErrKembar) {
				return idTerpakai(baru) // balapan dengan penulis lain sesudah AdaID (PK SYS_C008825 menolak)
			}
			return err
		}
		id = baru
		return nil
	})
	if err != nil {
		return models.CauseOfLoss{}, petaGalat(err)
	}
	c, err := l.gudang.Ambil(ctx, nil, id)
	return c, petaGalat(err)
}
