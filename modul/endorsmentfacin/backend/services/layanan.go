package services

// Layanan modul di atas Oracle - irisan I2 (periksa sebelum Create). Urutan query mengikuti
// `SetErrorBatalEndorsement_Act` langkah 5 dan 7, lalu `CheckEDMPolisDate` langkah 3; query gerbang yang dilompati
// korpus (RI slip, Internal Retro bukan 1) tidak dijalankan.
//
// Dibaca sesudah: periksa.go.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/endorsmentfacin/backend/models"
	"nusantarare/modul/endorsmentfacin/backend/repository"
)

// Service membawa akar bersama dan kontrak lintas modul yang dirakit.
type Service struct {
	*inti.Dasar
	// jenisBisnis - kontrak.JenisBisnisFacIn dari nbfacin (BusinessType_DeT); nil = belum tersambung.
	jenisBisnis kontrak.JenisBisnisFacIn
}

// DariDasar membungkus akar bersama.
func DariDasar(d *inti.Dasar) *Service { return &Service{Dasar: d} }

// DenganJenisBisnis memasang kontrak JenisBisnisFacIn (modul.go, dari perakit).
func (s *Service) DenganJenisBisnis(j kontrak.JenisBisnisFacIn) *Service {
	s.jenisBisnis = j
	return s
}

// PunyaDatabase - Oracle terpasang?
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar.PunyaDatabase() }

// LayananOracle merakit layanan di atas Oracle.
func LayananOracle(s *Service) *Layanan {
	return BaruLayanan(repository.Baru(s.DB())).DenganTransaksi(s.DalamTransaksi).DenganJenisBisnis(s.jenisBisnis)
}

var _ Gudang = (*repository.Gudang)(nil)

// JenisPengecualianTanggal - `M_EDM_PENGECUALIAN_POLIS.JENIS` pengganti literal `CheckEDMPolisDate` langkah 2.
const JenisPengecualianTanggal = "TANGGAL"

// ErrMasukanTidakSah - isian ditolak; pesannya untuk pengguna (Inggris, seperti layar).
var ErrMasukanTidakSah = errors.New("services: masukan tidak sah")

func tolak(format string, a ...any) error {
	return fmt.Errorf("%w: %s", ErrMasukanTidakSah, fmt.Sprintf(format, a...))
}

// Gudang - yang layanan butuhkan dari Oracle (`repository.Gudang`; tiruan di uji).
type Gudang interface {
	StatusEDM(ctx context.Context, polis string) (string, error)
	JenisBisnis(ctx context.Context, polis string) (string, error)
	Klaim(ctx context.Context, polis string) ([]string, error)
	AdaPembayaran(ctx context.Context, polis string) (bool, error)
	AdaRenewal(ctx context.Context, polis string) (bool, error)
	AdaFacout(ctx context.Context, polis string) (bool, error)
	AdaProteksi(ctx context.Context, polis string, tipe int) (bool, error)
	PeriodeProduksi(ctx context.Context, polis string) (*repository.PeriodeTeks, error)
	Pengecualian(ctx context.Context, jenis string) ([]string, error)
	EDMBelumSelesai(ctx context.Context, polis string) (string, bool, error)
	PolisSistemBaru(ctx context.Context, polis string) (bool, error)
	EdmTypeSelesaiProduksi(ctx context.Context, polis string) (string, bool, error)
	BisnisProduksi(ctx context.Context, polis string) ([]repository.BarisBisnis, error)

	// Tulisan I3 - di dalam transaksi pemanggil.
	NomorKasusBerikut(ctx context.Context, tx *db.Tx) (string, error)
	SisipWorkPolis(ctx context.Context, tx *db.Tx, id, akun string) error
	SisipKasus(ctx context.Context, tx *db.Tx, k models.KepalaKasus) error
}

// Transaksi menjalankan fn di dalam satu transaksi (commit bila nil, rollback bila galat).
type Transaksi func(ctx context.Context, fn func(tx *db.Tx) error) error

// Layanan adalah aturan modul di atas Gudang.
type Layanan struct {
	gudang Gudang
	tx     Transaksi
	// jenisBisnis - kontrak JenisBisnisFacIn; nil = lini polis sistem baru belum dapat dihitung (A52).
	jenisBisnis kontrak.JenisBisnisFacIn
}

// DenganJenisBisnis memasang kontrak JenisBisnisFacIn (nil diperbolehkan).
func (l *Layanan) DenganJenisBisnis(j kontrak.JenisBisnisFacIn) *Layanan {
	l.jenisBisnis = j
	return l
}

// BaruLayanan membuat layanan; tulisan (Buat) butuh DenganTransaksi.
func BaruLayanan(g Gudang) *Layanan { return &Layanan{gudang: g} }

// DenganTransaksi memasang pembuka transaksi.
func (l *Layanan) DenganTransaksi(tx Transaksi) *Layanan {
	l.tx = tx
	return l
}

// UraiMasukanPeriksa - validasi isian (422 bila salah). policyNo dipakai APA ADANYA (teks bebas `pxTextInput`), hanya
// kekosongannya yang diuji.
func UraiMasukanPeriksa(m models.MasukanPeriksa) (DataPeriksa, error) {
	if strings.TrimSpace(m.PolicyNo) == "" {
		return DataPeriksa{}, tolak("Policy No is required")
	}
	switch models.JenisEndorsemen(m.EdmType) {
	case "", models.EdmBatalSejakSemula, models.EdmBatalProrata, models.EdmPerubahan:
	default:
		return DataPeriksa{}, tolak("EDM Type must be 1, 2 or 4")
	}
	d := DataPeriksa{
		PolicyNo: m.PolicyNo, EdmType: models.JenisEndorsemen(m.EdmType), EdmTypeNew: m.EdmTypeNew,
		EndorsementInternalRetro: m.EndorsementInternalRetro,
	}
	if s := strings.TrimSpace(m.EndorsementDate); s != "" {
		// Tanggal layar = tanggal kalender pengguna (Asia/Jakarta); CheckEDMPolisDate memformatnya di zona itu.
		w, err := time.ParseInLocation(formatKabel, s, jakarta)
		if err != nil {
			return DataPeriksa{}, tolak("Endorsement Date must be DD-MM-YYYY")
		}
		d.TanggalEndorsement = &w
	}
	return d, nil
}

// Periksa - `POST /periksa`.
func (l *Layanan) Periksa(ctx context.Context, m models.MasukanPeriksa) (models.HasilPeriksa, error) {
	d, err := UraiMasukanPeriksa(m)
	if err != nil {
		return models.HasilPeriksa{}, err
	}
	if err := l.bacaData(ctx, &d); err != nil {
		return models.HasilPeriksa{}, err
	}
	return SusunHasilPeriksa(d)
}

// bacaData - seluruh query pemeriksaan atas masukan yang sudah terurai.
func (l *Layanan) bacaData(ctx context.Context, d *DataPeriksa) error {
	if err := l.bacaGerbang(ctx, d); err != nil {
		return err
	}
	// Langkah 7 - jenis bisnis polis lama.
	var err error
	if d.JenisBisnis, err = l.gudang.JenisBisnis(ctx, d.PolicyNo); err != nil {
		return err
	}
	// A51 (keputusan work owner 12-10-2026 opsi 1a + urutan sumber): DATA_JSON dulu; HANYA bila kosong dan polisnya
	// sistem baru, turunkan dari tabel - FACINPRODUCTION.BUSINESSCODE → BUSINESS (GROUPPANEL, OLDID) → kontrak
	// JenisBisnisFacIn (BusinessType_DeT milik nbfacin). Bukan tepat satu baris → ("", "") → "UNKNOWN" (pola NB).
	if d.SistemBaru && strings.TrimSpace(d.JenisBisnis) == "" && l.jenisBisnis != nil {
		baris, err := l.gudang.BisnisProduksi(ctx, d.PolicyNo)
		if err != nil {
			return err
		}
		var b repository.BarisBisnis
		if len(baris) == 1 {
			b = baris[0]
		}
		d.JenisBisnis, d.JenisBisnisDariTabel = l.jenisBisnis.BusinessType(b.GroupPanel, b.OldID), true
	}
	if d.TanggalEndorsement != nil {
		return l.bacaTanggal(ctx, d)
	}
	return nil
}

// bacaGerbang - query langkah 5 yang benar-benar dijalankan korpus.
func (l *Layanan) bacaGerbang(ctx context.Context, d *DataPeriksa) error {
	g := l.gudang
	// 5.6-5.9 gerbang 2 - berjalan juga untuk RI slip.
	id, ada, err := g.EDMBelumSelesai(ctx, d.PolicyNo)
	switch {
	case errors.Is(err, repository.ErrBelumAda):
		d.G2TakTerbaca = true // tabel 262 belum dibuat DBA
	case err != nil:
		return err
	default:
		d.G2Diperiksa = true
		if ada {
			d.EDMBelumSelesai = []string{id}
		}
	}
	// A51: penanda polis sistem baru - dibaca sebelum cabang RI slip karena gerbang bisnis (langkah 7-18) juga
	// memakainya.
	if d.SistemBaru, err = g.PolisSistemBaru(ctx, d.PolicyNo); err != nil {
		return err
	}
	// 5.1 / 5.10 RI slip: 5.2-5.5 dan 5.11-5.23 dilompati.
	if d.EdmType == models.EdmPerubahan && d.EdmTypeNew == models.EdmTypeNewRISlip {
		return nil
	}
	if err := l.bacaStatusEDM(ctx, d); err != nil { // 5.2
		return err
	}
	if d.Klaim, err = g.Klaim(ctx, d.PolicyNo); err != nil { // 5.11
		return err
	}
	// 5.15 dijalankan Pega untuk setiap EdmType, tetapi hasilnya hanya dibaca bila EdmType 1/2 - lintas skema, jadi
	// dibaca hanya bila menentukan.
	if d.EdmType == models.EdmBatalSejakSemula || d.EdmType == models.EdmBatalProrata {
		ada, err := g.AdaPembayaran(ctx, d.PolicyNo)
		switch {
		case errors.Is(err, repository.ErrSkemaLainTakTerbaca):
			// Pembayaran tetap nil → gerbang 4 "belum-diperiksa".
		case err != nil:
			return err
		default:
			d.Pembayaran = &ada
		}
	}
	if d.Renewal, err = g.AdaRenewal(ctx, d.PolicyNo); err != nil { // 5.19
		return err
	}
	retro, err := samaKode(d.EndorsementInternalRetro, "1") // 5.23
	if err != nil {
		return err
	}
	if retro {
		if d.Facout, err = g.AdaFacout(ctx, d.PolicyNo); err != nil {
			return err
		}
	}
	return l.bacaKlep(ctx, d)
}

// bacaStatusEDM - CARI20 gerbang 1 (dan StatusEDMPolis OpenCase). Polis era Pega: GetEDMStatus_SQL apa adanya. Polis
// sistem baru (A51, persetujuan work owner 12-10-2026): EDM_TYPE kasus EDM sistem baru yang berproduksi; tabel 262
// belum ada → G1TakTerbaca.
func (l *Layanan) bacaStatusEDM(ctx context.Context, d *DataPeriksa) error {
	var err error
	if !d.SistemBaru {
		if d.StatusEDM, err = l.gudang.StatusEDM(ctx, d.PolicyNo); err != nil {
			return err
		}
		d.StatusEDMDibaca = true
		return nil
	}
	v, _, err := l.gudang.EdmTypeSelesaiProduksi(ctx, d.PolicyNo)
	switch {
	case errors.Is(err, repository.ErrBelumAda):
		d.G1TakTerbaca = true
	case err != nil:
		return err
	default:
		d.StatusEDM = v
	}
	d.StatusEDMDibaca = true
	return nil
}

// bacaKlep - `BrowseOpenProteksiEdm_RD` langkah 5.3/5.12/5.16/5.20, Param.Type 1-4. Pega menjalankan keempatnya
// tanpa melihat syarat tolak gerbangnya (bacaGerbang sudah keluar untuk RI slip, yang melompati keempatnya).
func (l *Layanan) bacaKlep(ctx context.Context, d *DataPeriksa) error {
	d.Klep = map[KlepPembatal]bool{}
	for _, k := range []KlepPembatal{KlepBatal, KlepKlaim, KlepPembayaran, KlepRenewal} {
		ada, err := l.gudang.AdaProteksi(ctx, d.PolicyNo, int(k))
		switch {
		case errors.Is(err, repository.ErrBelumAda):
			d.Klep, d.KlepTakTerbaca = nil, true // tabel tidak ada: seluruh klep kosong
			return nil
		case err != nil:
			return err
		}
		d.Klep[k] = ada
	}
	return nil
}

// bacaTanggal - pengecualian (pengganti literal langkah 2) dan periode `GetStartDate` (langkah 3).
func (l *Layanan) bacaTanggal(ctx context.Context, d *DataPeriksa) error {
	daftar, err := l.gudang.Pengecualian(ctx, JenisPengecualianTanggal)
	switch {
	case errors.Is(err, repository.ErrBelumAda):
		d.PengecualianBelumAda = true // keputusan work owner 11-10-2026: tabel belum ada = daftar kosong
	case err != nil:
		return err
	default:
		d.PengecualianTanggal = daftar
	}
	p, err := l.gudang.PeriodeProduksi(ctx, d.PolicyNo)
	if err != nil {
		return err
	}
	if p != nil {
		d.Periode = &PeriodeMentah{Mulai: p.Mulai, Akhir: p.Akhir}
	}
	return nil
}
