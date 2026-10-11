// Package services memuat aturan modul Cover Life (`coverlife`; panduan XML section `InboxCoverLife`; keputusan work
// owner 08-10-2026 C1-C4 dan K5, MODUL.md). Pola causeoflosslife (satu tabel, satu layar), TANPA saring / urut pilihan
// / Delete / Upload (XML tidak memuatnya).
//
//   - Grid `BrowseCoverLife_RD` (b4026): kolom ID (`.ID` b3449) dan Cover (`.Cover` b3596); Note TIDAK tampil; 50 baris
//     per halaman (b4101), ID MENAIK tetap (b3967 / b3972), tanpa saring (b4080).
//   - Save (b5407 -> `AddToList_Act` b5426): Add bila ID kosong, selain itu Edit (Edit b3784 -> `EditList_DT` b3807
//     mengisi form dengan ID, Cover, Note). Cancel (b5671 -> `NewData_DT` b5694, hanya saat `DATASHOW = 'IsEdit'`
//     b5831) murni layar.
//   - ID baru = '1' || LPAD(M_COVER_LIFE_SEQ.NEXTVAL, 5, '0') (models.BentukID, C2). NEXTVAL yang menghasilkan ID yang
//     SUDAH ada = galat ErrIDTerpakai yang jelas (409), bukan 500 dan bukan lompat diam-diam.
//   - C3: Cover wajib (XML), Note opsional; huruf TIDAK diubah (XML tanpa pengubah huruf); Cover dipangkas dan tidak
//     boleh kembar tanpa beda huruf = di luar XML (PARITAS, WO boleh menolak).
//   - Transaksi milik Go: satu transaksi per simpan; prosedur PEGA_M_COVER_LIFE tidak dipanggil (INVALID sesudah 086,
//     C2 diterima WO).
//   - Hak menu Full / View only: View only tanpa Save (403).
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/coverlife/backend/models"
	"nusantarare/modul/coverlife/backend/repository"
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
	ErrTidakAda = errors.New("services: cover tidak ada")
	// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrIDTerpakai - C2: ID dari M_COVER_LIFE_SEQ sudah ada di M_COVER_LIFE (atau nomornya melewati LPAD 5); pesannya
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

func idTerpakai(id string) error {
	return fmt.Errorf("%w: new ID %s from M_COVER_LIFE_SEQ already exists; ask the administrator to move the sequence past the highest ID",
		ErrIDTerpakai, id)
}

// Transaksi menjalankan fn di dalam satu transaksi.
type Transaksi func(ctx context.Context, fn func(tx *dbTx) error) error

// Gudang - yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	Daftar(ctx context.Context, s models.Saringan) ([]models.Cover, int, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Cover, error)
	PemakaiNama(ctx context.Context, tx *dbTx, nama, kecualiID string) ([]models.Cover, error)
	NomorBaru(ctx context.Context, tx *dbTx) (string, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, c models.Cover) error
	Ubah(ctx context.Context, tx *dbTx, c models.Cover) error
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
		return dilarang("Your access to the Cover Life menu is View only")
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
		d = []models.Cover{}
	}
	return models.Halaman{Daftar: d, Total: total, Halaman: s.Halaman, Ukuran: models.UkuranHalaman}, nil
}

func petaGalat(err error) error {
	if errors.Is(err, repository.ErrTidakAda) {
		return ErrTidakAda
	}
	return err
}

// idBaru - C2: '1' || LPAD(NEXTVAL, 5, '0'). ID yang sudah ada = ErrIDTerpakai yang menyebut ID dan sequence-nya -
// TIDAK dilompati diam-diam.
func (l *Layanan) idBaru(ctx context.Context, tx *dbTx) (string, error) {
	nomor, err := l.gudang.NomorBaru(ctx, tx)
	if err != nil {
		return "", err
	}
	id, err := models.BentukID(nomor)
	if err != nil {
		return "", fmt.Errorf("%w: M_COVER_LIFE_SEQ gave number %s, which does not fit the ID formula '1' || LPAD(number, 5); ask the administrator",
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

// Simpan - Save (`AddToList_Act` b5426): Add bila id kosong (ID dari sequence, C2), selain itu Edit (`EditList_DT`
// b3807: ID dari jalur, Cover, Note). Satu transaksi.
func (l *Layanan) Simpan(ctx context.Context, a Aktor, id string, isi models.Isian) (models.Cover, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Cover{}, err
	}
	id = strings.TrimSpace(id)
	cover, err := models.NormalCover(isi.Cover)
	if err != nil {
		return models.Cover{}, tolak("%s", err.Error())
	}
	note, err := models.NormalNote(isi.Note)
	if err != nil {
		return models.Cover{}, tolak("%s", err.Error())
	}
	err = l.tx(ctx, func(tx *dbTx) error {
		if id != "" {
			if _, err := l.gudang.Ambil(ctx, tx, id); err != nil {
				return err
			}
		}
		ganda, err := l.gudang.PemakaiNama(ctx, tx, cover, id)
		if err != nil {
			return err
		}
		if len(ganda) > 0 {
			return tolak("Cover %s is already used by ID %s", cover, ganda[0].ID)
		}
		if id != "" {
			return l.gudang.Ubah(ctx, tx, models.Cover{ID: id, Cover: cover, Note: note})
		}
		baru, err := l.idBaru(ctx, tx)
		if err != nil {
			return err
		}
		if err := l.gudang.Sisip(ctx, tx, models.Cover{ID: baru, Cover: cover, Note: note}); err != nil {
			if errors.Is(err, repository.ErrKembar) {
				return idTerpakai(baru) // balapan dengan penulis lain sesudah AdaID (PK SYS_C009203 menolak)
			}
			return err
		}
		id = baru
		return nil
	})
	if err != nil {
		return models.Cover{}, petaGalat(err)
	}
	c, err := l.gudang.Ambil(ctx, nil, id)
	return c, petaGalat(err)
}
