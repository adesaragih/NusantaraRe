// Package services memuat aturan modul Disease Life (`diseaselife`; panduan XML section `InboxDisease`; keputusan
// work owner 08-10-2026 D1-D4 dan K5, MODUL.md). Pola causeoflosslife / benefitlife (satu tabel, satu layar Edit / Save
// / Cancel), TANPA Delete / Upload (XML tidak memuatnya).
//
//   - Grid `BrowseDiseaseLife_RD` (b5094): kolom ID (`.Number` b4321), ICD Code (`.ICD_Code` b4484), Disease
//     (`.Disease` b4629), 10 baris per halaman (b5169), bawaan ID MENURUN (b5012 / b5017), urut pilihan ID / ICD Code
//     (b5014 / b5036), saring (b5148) - SEMUA DI SERVER: 97.586 baris DEV tidak pernah dimuat sekaligus.
//   - Save (b2102 -> `AddToList_Act` b2121): Add bila ID kosong, selain itu Edit (Edit b4829 -> `EditList_DT` b4852
//     mengisi form dengan Number, Disease, ICD_Code). Cancel (b2365 -> `NewData_DT` b2388, hanya saat
//     `DATASHOW = 'IsEdit'` b2525) murni layar.
//   - ID baru = TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL) (D1.1); ID yang SUDAH ada (mis. ditulis prosedur Pega) = galat
//     ErrIDTerpakai yang jelas (409), bukan 500 dan bukan lompat diam-diam.
//   - D3: ICD Code dan Disease wajib; huruf besar pada keduanya (`SetUpperCase_DT` dipanggil kedua medan); ICD Code unik
//     tanpa beda huruf + pangkas = di luar XML (PARITAS).
//   - Transaksi milik Go: satu transaksi per simpan; prosedur PEGA_DISEASE_LIFE tidak dipanggil.
//   - Hak menu Full / View only: View only tanpa Save (403).
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/diseaselife/backend/models"
	"nusantarare/modul/diseaselife/backend/repository"
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
	ErrTidakAda = errors.New("services: disease tidak ada")
	// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrIDTerpakai - ID dari SEQ_DISEASE_LIFE sudah ada di DISEASE_LIFE; pesannya untuk pengguna / admin.
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
	return fmt.Errorf("%w: new ID %s from SEQ_DISEASE_LIFE already exists; ask the administrator to move the sequence past the highest ID",
		ErrIDTerpakai, id)
}

// Transaksi menjalankan fn di dalam satu transaksi.
type Transaksi func(ctx context.Context, fn func(tx *dbTx) error) error

// Gudang - yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	Daftar(ctx context.Context, s models.Saringan) ([]models.Penyakit, int, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Penyakit, error)
	PemakaiICD(ctx context.Context, tx *dbTx, icd, kecualiID string) ([]models.Penyakit, error)
	NomorBaru(ctx context.Context, tx *dbTx) (string, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, p models.Penyakit) error
	Ubah(ctx context.Context, tx *dbTx, p models.Penyakit) error
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
		return dilarang("Your access to the Disease Life menu is View only")
	}
	return nil
}

// Daftar - satu halaman grid (halaman >= 1, urut ID / ICD Code, saring dirapikan).
func (l *Layanan) Daftar(ctx context.Context, s models.Saringan) (models.Halaman, error) {
	if s.Halaman < 1 {
		s.Halaman = 1
	}
	if s.Urut != models.UrutICD {
		s.Urut = models.UrutID
	}
	s.ICDCode, s.Disease = models.RapikanSaring(s.ICDCode), models.RapikanSaring(s.Disease)
	d, total, err := l.gudang.Daftar(ctx, s)
	if err != nil {
		return models.Halaman{}, err
	}
	if d == nil {
		d = []models.Penyakit{}
	}
	return models.Halaman{Daftar: d, Total: total, Halaman: s.Halaman, Ukuran: models.UkuranHalaman}, nil
}

func petaGalat(err error) error {
	if errors.Is(err, repository.ErrTidakAda) {
		return ErrTidakAda
	}
	return err
}

// idBaru - D1.1: TO_CHAR(SEQ_DISEASE_LIFE.NEXTVAL). ID yang sudah ada = ErrIDTerpakai yang menyebut ID dan
// sequence-nya - TIDAK dilompati diam-diam.
func (l *Layanan) idBaru(ctx context.Context, tx *dbTx) (string, error) {
	nomor, err := l.gudang.NomorBaru(ctx, tx)
	if err != nil {
		return "", err
	}
	id, err := models.BentukID(nomor)
	if err != nil {
		return "", fmt.Errorf("%w: SEQ_DISEASE_LIFE gave %q, which is not a valid ID; ask the administrator",
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

// Simpan - Save (`AddToList_Act` b2121): Add bila id kosong (ID dari SEQ_DISEASE_LIFE), selain itu Edit
// (`EditList_DT` b4852: ID dari jalur, ICD Code, Disease). Satu transaksi.
func (l *Layanan) Simpan(ctx context.Context, a Aktor, id string, isi models.Isian) (models.Penyakit, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Penyakit{}, err
	}
	id = strings.TrimSpace(id)
	icd, err := models.NormalICDCode(isi.ICDCode)
	if err != nil {
		return models.Penyakit{}, tolak("%s", err.Error())
	}
	nama, err := models.NormalDisease(isi.Disease)
	if err != nil {
		return models.Penyakit{}, tolak("%s", err.Error())
	}
	err = l.tx(ctx, func(tx *dbTx) error {
		if id != "" {
			if _, err := l.gudang.Ambil(ctx, tx, id); err != nil {
				return err
			}
		}
		ganda, err := l.gudang.PemakaiICD(ctx, tx, icd, id)
		if err != nil {
			return err
		}
		if len(ganda) > 0 {
			return tolak("ICD Code %s is already used by ID %s", icd, ganda[0].ID)
		}
		if id != "" {
			return l.gudang.Ubah(ctx, tx, models.Penyakit{ID: id, ICDCode: icd, Disease: nama})
		}
		baru, err := l.idBaru(ctx, tx)
		if err != nil {
			return err
		}
		if err := l.gudang.Sisip(ctx, tx, models.Penyakit{ID: baru, ICDCode: icd, Disease: nama}); err != nil {
			if errors.Is(err, repository.ErrKembar) {
				return idTerpakai(baru) // balapan dengan penulis lain sesudah AdaID (PK_DISEASE_LIFE menolak)
			}
			return err
		}
		id = baru
		return nil
	})
	if err != nil {
		return models.Penyakit{}, petaGalat(err)
	}
	p, err := l.gudang.Ambil(ctx, nil, id)
	return p, petaGalat(err)
}
