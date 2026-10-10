// Package services memuat aturan modul Plan (`planlife`; panduan XML section `InboxProductType`; keputusan work owner
// 08-10-2026 K1-K7, MODUL.md). Pola benefitlife (satu tabel, satu layar) + dua autocomplete.
//
//   - Grid `BrowseProductTypeLife_RD` (b4325): Plan Name, Business, Benefit, Edit; 10 per halaman (b4401); tanpa
//     saring (b4380); urut Plan Name / Benefit (b4247 / b4291), bawaan ID menaik (XML tanpa urutan bawaan).
//   - Save (b6403 -> `SaveProductTypeLife_Act` b6422): Add bila form tanpa ID, selain itu Edit (Edit b4027 ->
//     `EditProductTypeLife_Act` b4045 dengan ID, Business, BusinessID, Benefit, BenefitID, CoverName). New (b6667 ->
//     `NewProductTypeLife_act` b6686, hanya `DATASHOW = 'IsEdit'` b6827) murni layar.
//   - K4: teks Business / Benefit (autocomplete `pyAllowFreeFormInput` true b1125 / b1478) DICOCOKKAN ULANG ke master
//     tanpa beda huruf; cocok = nama + ID dari master, tidak cocok / ganda = 422.
//   - K5: Plan Name wajib + unik tanpa beda huruf (di luar XML); Business / Benefit wajib (tidak terbukti boleh kosong).
//   - K3: ID = '1' || LPAD(M_PRODUCT_TYPE_LIFE_SEQ.NEXTVAL, 5, '0'); ID yang sudah ada = 409 berkalimat.
//   - Hak menu Full / View only: View only tanpa Save (403).
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/planlife/backend/models"
	"nusantarare/modul/planlife/backend/repository"
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
	ErrTidakAda = errors.New("services: plan tidak ada")
	// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna.
	ErrMasukanTidakSah = errors.New("services: masukan tidak sah")
	// ErrDilarang - aktor tidak berhak; pesannya untuk pengguna.
	ErrDilarang = errors.New("services: tidak berhak")
	// ErrIDTerpakai - K3: ID dari M_PRODUCT_TYPE_LIFE_SEQ sudah ada (atau melewati LPAD 5); pesannya untuk pengguna.
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
	return fmt.Errorf("%w: new ID %s from M_PRODUCT_TYPE_LIFE_SEQ already exists; ask the administrator to move the sequence past the highest ID",
		ErrIDTerpakai, id)
}

// Transaksi menjalankan fn di dalam satu transaksi.
type Transaksi func(ctx context.Context, fn func(tx *dbTx) error) error

// Gudang - yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	Daftar(ctx context.Context, s models.Saringan) ([]models.Plan, int, error)
	Ambil(ctx context.Context, tx *dbTx, id string) (models.Plan, error)
	PemakaiNama(ctx context.Context, tx *dbTx, nama, kecualiID string) ([]models.Plan, error)
	NomorBaru(ctx context.Context, tx *dbTx) (string, error)
	AdaID(ctx context.Context, tx *dbTx, id string) (bool, error)
	Sisip(ctx context.Context, tx *dbTx, p models.Plan) error
	Ubah(ctx context.Context, tx *dbTx, p models.Plan) error
	PilihanBusiness(ctx context.Context, grup string) ([]models.PilihanBusiness, error)
	PilihanBenefit(ctx context.Context) ([]models.PilihanBenefit, error)
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
		return dilarang("Your access to the Plan menu is View only")
	}
	return nil
}

// RapikanSaringan - kolom urut di luar daftar = bawaan (ID menaik); halaman >= 1.
func RapikanSaringan(s models.Saringan) models.Saringan {
	s.Urut = strings.ToLower(strings.TrimSpace(s.Urut))
	ada := false
	for _, k := range models.KolomUrut {
		ada = ada || k == s.Urut
	}
	if !ada {
		s.Urut, s.Turun = "", false
	}
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
		d = []models.Plan{}
	}
	return models.Halaman{Daftar: d, Total: total, Halaman: s.Halaman, Ukuran: models.UkuranHalaman}, nil
}

// PilihanBusiness - autocomplete Business (K4: grup life `GROUPPANEL = '009'`).
func (l *Layanan) PilihanBusiness(ctx context.Context) ([]models.PilihanBusiness, error) {
	return l.gudang.PilihanBusiness(ctx, models.GrupBusinessLife)
}

// PilihanBenefit - autocomplete Benefit (K4: seluruh BENEFIT_LIFE).
func (l *Layanan) PilihanBenefit(ctx context.Context) ([]models.PilihanBenefit, error) {
	return l.gudang.PilihanBenefit(ctx)
}

// cocokBusiness - K4: teks = NOTE master (tanpa beda huruf / spasi tepi). Kosong / >1 cocok = pesan.
func cocokBusiness(daftar []models.PilihanBusiness, teks string) (models.PilihanBusiness, string) {
	var hasil []models.PilihanBusiness
	for _, b := range daftar {
		if models.KunciTeks(b.Note) == models.KunciTeks(teks) {
			hasil = append(hasil, b)
		}
	}
	switch len(hasil) {
	case 1:
		return hasil[0], ""
	case 0:
		return models.PilihanBusiness{}, fmt.Sprintf("Business %q is not in the Business list; choose one from the list", teks)
	}
	return models.PilihanBusiness{}, fmt.Sprintf("Business %q matches more than one Business (IDs %s and %s); ask the administrator", teks, hasil[0].ID, hasil[1].ID)
}

// cocokBenefit - K4: teks = BENEFIT master (tanpa beda huruf / spasi tepi).
func cocokBenefit(daftar []models.PilihanBenefit, teks string) (models.PilihanBenefit, string) {
	var hasil []models.PilihanBenefit
	for _, b := range daftar {
		if models.KunciTeks(b.Benefit) == models.KunciTeks(teks) {
			hasil = append(hasil, b)
		}
	}
	switch len(hasil) {
	case 1:
		return hasil[0], ""
	case 0:
		return models.PilihanBenefit{}, fmt.Sprintf("Benefit %q is not in the Benefit list; choose one from the list", teks)
	}
	return models.PilihanBenefit{}, fmt.Sprintf("Benefit %q matches more than one Benefit (IDs %s and %s); ask the administrator", teks, hasil[0].ID, hasil[1].ID)
}

// rakit - isian -> baris (nama dan ID dari master). Galat = SEMUA masalah sekaligus (422).
func (l *Layanan) rakit(ctx context.Context, isi models.Isian) (models.Plan, error) {
	rapi, err := models.PeriksaIsian(isi)
	if err != nil {
		return models.Plan{}, tolak("%s", err.Error())
	}
	bizz, err := l.gudang.PilihanBusiness(ctx, models.GrupBusinessLife)
	if err != nil {
		return models.Plan{}, err
	}
	bens, err := l.gudang.PilihanBenefit(ctx)
	if err != nil {
		return models.Plan{}, err
	}
	var pesan []string
	b, salah := cocokBusiness(bizz, rapi.Business)
	if salah != "" {
		pesan = append(pesan, salah)
	}
	e, salah := cocokBenefit(bens, rapi.Benefit)
	if salah != "" {
		pesan = append(pesan, salah)
	}
	if len(pesan) > 0 {
		return models.Plan{}, tolak("%s", strings.Join(pesan, "; "))
	}
	return models.Plan{CoverName: rapi.CoverName, Business: strings.TrimSpace(b.Note), BusinessID: b.ID,
		Benefit: strings.TrimSpace(e.Benefit), BenefitID: e.ID}, nil
}

func petaGalat(err error) error {
	if errors.Is(err, repository.ErrTidakAda) {
		return ErrTidakAda
	}
	return err
}

// idBaru - K3: '1' || LPAD(NEXTVAL, 5, '0'); ID yang sudah ada = ErrIDTerpakai (tidak dilompati diam-diam).
func (l *Layanan) idBaru(ctx context.Context, tx *dbTx) (string, error) {
	nomor, err := l.gudang.NomorBaru(ctx, tx)
	if err != nil {
		return "", err
	}
	id, err := models.BentukID(nomor)
	if err != nil {
		return "", fmt.Errorf("%w: M_PRODUCT_TYPE_LIFE_SEQ gave number %s, which does not fit the ID formula '1' || LPAD(number, 5); ask the administrator",
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

// Simpan - Save (`SaveProductTypeLife_Act` b6422): Add bila id kosong (ID dari sequence, K3), selain itu Edit
// (`EditProductTypeLife_Act` b4045). Satu transaksi.
func (l *Layanan) Simpan(ctx context.Context, a Aktor, id string, isi models.Isian) (models.Plan, error) {
	if err := wajibPenuh(a); err != nil {
		return models.Plan{}, err
	}
	id = strings.TrimSpace(id)
	p, err := l.rakit(ctx, isi)
	if err != nil {
		return models.Plan{}, err
	}
	err = l.tx(ctx, func(tx *dbTx) error {
		if id != "" {
			if _, err := l.gudang.Ambil(ctx, tx, id); err != nil {
				return err
			}
		}
		ganda, err := l.gudang.PemakaiNama(ctx, tx, p.CoverName, id)
		if err != nil {
			return err
		}
		if len(ganda) > 0 {
			return tolak("Plan Name %s is already used by ID %s", p.CoverName, ganda[0].ID)
		}
		if id != "" {
			p.ID = id
			return l.gudang.Ubah(ctx, tx, p)
		}
		baru, err := l.idBaru(ctx, tx)
		if err != nil {
			return err
		}
		p.ID = baru
		if err := l.gudang.Sisip(ctx, tx, p); err != nil {
			if errors.Is(err, repository.ErrKembar) {
				return idTerpakai(baru) // balapan dengan penulis lain sesudah AdaID (PK menolak)
			}
			return err
		}
		id = baru
		return nil
	})
	if err != nil {
		return models.Plan{}, petaGalat(err)
	}
	hasil, err := l.gudang.Ambil(ctx, nil, id)
	return hasil, petaGalat(err)
}
