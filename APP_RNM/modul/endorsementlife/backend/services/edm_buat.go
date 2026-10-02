package services

// Tiket 01 (gerbang kelayakan) dan 02 (buat kasus + salin polis lama).
//
// Layar `EndorsmentLife_Section` memanggil `SetErrorBatalEndorsement_Act` setiap
// kali `Policy No` (b1249), `EDM Type` (b1545), atau `EDM Date` (b3347) berubah
// - di sini `Kelayakan`. Tombol `Submit` b4226 menjalankan `MappingEDMLife`
// (b4415) lalu `openAssignment` - di sini `BuatKasus`, yang menjalankan ULANG
// kelima gerbang di dalam transaksi (langkah 1 b371 juga memeriksa ulang dua di
// antaranya, prakondisi b433/b456 T=6).

import (
	"context"
	"errors"
	"fmt"
	"strings"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/jejak"
	"nusantarare/modul/endorsementlife/backend/models"
	"nusantarare/modul/endorsementlife/backend/repository"
)

// GalatKelayakan - satu atau lebih gerbang menolak; pesannya VERBATIM korpus.
type GalatKelayakan struct{ Pesan []string }

func (g GalatKelayakan) Error() string {
	return "services: endorsement rejected by eligibility gates: " + strings.Join(g.Pesan, "; ")
}

// PesanLayar - seluruh pesan gerbang, dipisah baris baru (layar menampilkannya per baris).
func (g GalatKelayakan) PesanLayar() string { return strings.Join(g.Pesan, "\n") }

// ErrKasusTerbukaGanda - pembuatan bersamaan atas polis yang sama: kasus kedua
// ditolak index unik (migrasi 482). Pesannya pesan gerbang 3.
var ErrKasusTerbukaGanda = errors.New(models.PesanEDMBelumSelesai)

// ErrArasapas - gerbang 5 tidak dapat membaca pembayaran (OQ-EDM-012).
var ErrArasapas = errors.New("services: payment status in Arasapas cannot be read; a Batal endorsement cannot be checked")

// MasukanKelayakan adalah isian yang diperiksa gerbang.
type MasukanKelayakan struct {
	NomorPolis string `json:"policyNo"`
	EdmType    string `json:"edmType"`
}

// fakta mengumpulkan jawaban basis data kelima gerbang.
func (l *Layanan) fakta(ctx context.Context, tx *db.Tx, nomorPolis, edmType string) (models.FaktaGerbang, models.Versi, error) {
	var f models.FaktaGerbang
	if nomorPolis == "" {
		return f, models.Versi{}, nil
	}
	v, ada, err := l.gudang.VersiBerjalan(ctx, tx, nomorPolis, 0)
	if err != nil {
		return f, models.Versi{}, err
	}
	f.PolisAda = ada
	f.SudahBatal = ada && v.SudahBatal()
	if f.AdaKasusTerbuka, err = l.gudang.AdaKasusTerbuka(ctx, tx, nomorPolis); err != nil {
		return f, models.Versi{}, err
	}
	// 3.10 b2319: pembayaran hanya berarti bila `EdmType=="3"` - dan skema
	// Arasapas tidak dibaca sama sekali di luar itu.
	if ada && edmType == models.EdmTypeBatal {
		dibayar, err := l.gudang.SudahDibayar(ctx, models.NomorInvoiceArasapas(nomorPolis))
		if errors.Is(err, repository.ErrArasapasTakTerbaca) {
			l.catat(fmt.Sprintf("endorsement life: gerbang 5 polis %q: %v", nomorPolis, err))
			return f, models.Versi{}, ErrArasapas
		}
		if err != nil {
			return f, models.Versi{}, err
		}
		f.SudahDibayar = dibayar
	}
	return f, v, nil
}

// Kelayakan menjalankan kelima gerbang tanpa menulis apa pun.
func (l *Layanan) Kelayakan(ctx context.Context, p inti.Pelaku, m MasukanKelayakan) (models.Kelayakan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return models.Kelayakan{}, err
	}
	np := strings.TrimSpace(m.NomorPolis)
	f, _, err := l.fakta(ctx, nil, np, m.EdmType)
	if err != nil {
		return models.Kelayakan{}, err
	}
	pesan := models.PesanGerbang(np, m.EdmType, f)
	if pesan == nil {
		pesan = []string{}
	}
	return models.Kelayakan{Pesan: pesan, Boleh: len(pesan) == 0 && models.PeriksaEdmType(m.EdmType) == nil}, nil
}

// MasukanKasus adalah isian `EndorsmentLife_Section` yang membuat kasus.
type MasukanKasus struct {
	NomorPolis string `json:"policyNo"`
	EdmType    string `json:"edmType"`
	EdmDate    string `json:"edmDate"`
	Deskripsi  string `json:"description"`
}

// HasilBuat - kasus baru dan cacah salinannya.
type HasilBuat struct {
	ID             string `json:"id"`
	Peserta        int    `json:"peserta"`
	Spreading      int    `json:"spreading"`
	SpreadingRetro int    `json:"spreadingRetro"`
}

// ErrMasukanTidakSah - isian kasus ditolak sebelum basis data disentuh.
var ErrMasukanTidakSah = errors.New("services: invalid endorsement input")

// BuatKasus - `Submit` b4226 → `MappingEDMLife`: gerbang ulang, kasus baru,
// salinan versi berjalan, jejak - SATU transaksi, commit sekali (commit dini
// `MappingEDMLife` 14 b3196: kasus terlihat pengguna lain begitu dibuat).
func (l *Layanan) BuatKasus(ctx context.Context, p inti.Pelaku, m MasukanKasus) (HasilBuat, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return HasilBuat{}, err
	}
	np := strings.TrimSpace(m.NomorPolis)
	if err := models.PeriksaEdmType(m.EdmType); err != nil {
		return HasilBuat{}, fmt.Errorf("%w: %v", ErrMasukanTidakSah, err)
	}
	tgl := strings.TrimSpace(m.EdmDate)
	if tgl != "" {
		if iso, ok := repository.TanggalPega(tgl); !ok || iso != tgl {
			return HasilBuat{}, fmt.Errorf("%w: EDM Date %q is not a YYYY-MM-DD date", ErrMasukanTidakSah, m.EdmDate)
		}
	}
	var hasil HasilBuat
	err := l.tx(ctx, func(tx *db.Tx) error {
		f, v, err := l.fakta(ctx, tx, np, m.EdmType)
		if err != nil {
			return err
		}
		if pesan := models.PesanGerbang(np, m.EdmType, f); len(pesan) > 0 {
			return GalatKelayakan{Pesan: pesan}
		}
		kepala, rusak, err := l.gudang.KepalaSumber(ctx, tx, v, np)
		if err != nil {
			return err
		}
		if len(rusak) > 0 {
			l.catat(fmt.Sprintf("endorsement life: tanggal kepala warisan tak terbaca, dikosongkan: %s", strings.Join(rusak, ", ")))
		}
		id, err := l.gudang.PengenalKasusBaru(ctx, tx)
		if err != nil {
			return err
		}
		if err := l.gudang.SisipKasus(ctx, tx, repository.KasusTulis{
			ID: id, NomorPolis: np, EdmType: m.EdmType, EdmDate: tgl, EdmNote: strings.TrimSpace(m.Deskripsi),
			ProdKe: v.ProdKe + 1, Pembuat: p.AkunID, Kepala: kepala,
		}); err != nil {
			if errors.Is(err, repository.ErrKasusTerbukaGanda) {
				return ErrKasusTerbukaGanda
			}
			return err
		}
		s, err := l.gudang.SalinVersi(ctx, tx, id, v, np)
		if err != nil {
			return err
		}
		hasil = HasilBuat{ID: id, Peserta: s.Peserta, Spreading: s.Spreading, SpreadingRetro: s.SpreadingRetro}
		return l.jejak.Rekam(ctx, tx, jejak.CatatanJejak{
			KlaimID: id, Dari: "", Ke: models.TahapInputEDMLife, AkunID: p.AkunID, Waktu: l.jam(),
		})
	})
	if err != nil {
		return HasilBuat{}, err
	}
	return hasil, nil
}
