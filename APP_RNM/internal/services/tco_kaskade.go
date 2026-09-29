package services

// Kaskade hapus kontrak / reinsurer + popup konfirmasi - tiket 10 Treaty Contract Out.
//
// Untuk apa berkas ini: dua langkah yang layar jalankan -
//
//  1. DAMPAK: jumlah baris tiap jenis yang akan ikut terhapus, dan jumlah
//     klausul yang TIDAK terhapus (popup Ya/Batal, penyimpangan sadar 4);
//  2. HAPUS: dalam SATU transaksi - kunci kontrak, HITUNG ULANG dampak dengan
//     saringan yang sama, tolak (409) bila angkanya berbeda dari yang pemakai
//     konfirmasi, kaskade, pastikan yang terhapus = yang dihitung, jejak.
//
// Dengan begitu angka di popup = angka yang benar-benar terhapus, walau ada
// penulis lain di antara pratinjau dan Ya.
//
// Pega: `Delete` kontrak b11809 -> `BrowseDeleteRowTreatyInContract` ->
// `DeleteFromTREATYCONTRACT_SQL` (tanpa konfirmasi; pesan b762 "Data Berhasil
// di Hapus"); `Delete` reinsurer (`ViewDetailTreatyReinsurerGrid1.xml` b4936)
// -> `DeleteTreatyReins_Act` -> `DeleteFromTreatyReinsurer_Act` (tanpa pesan).
//
// Dibaca sesudah: tco_reinsurer.go, repository/tco_kaskade.go.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// PesanKontrakTerhapusTCO - VERBATIM `BrowseDeleteRowTreatyInContract.xml` b762.
const PesanKontrakTerhapusTCO = "Data Berhasil di Hapus"

var (
	// ErrGudangKaskadeBelumDisuntik - pelaksana kaskade belum dipasang.
	ErrGudangKaskadeBelumDisuntik = errors.New("services: pelaksana kaskade hapus belum disuntik")
	// ErrKaskadeTidakUtuh - induk yang dihapus tidak tepat satu baris (409).
	ErrKaskadeTidakUtuh = repository.ErrKaskadeTidakUtuh
	// ErrDampakBerubah - dampak sekarang berbeda dari yang dikonfirmasi (409).
	ErrDampakBerubah = errors.New("services: jumlah baris terdampak berubah sejak konfirmasi")
)

// PelaksanaKaskadeTCO menghitung dan menjalankan kaskade.
type PelaksanaKaskadeTCO interface {
	DampakKontrak(ctx context.Context, kom models.KombinasiTCO, tahunID, kontrakID string) (repository.DampakHapusTCO, error)
	HapusKontrak(ctx context.Context, tx *repository.Tx, kom models.KombinasiTCO, tahunID, kontrakID string) (repository.DampakHapusTCO, error)
	DampakReinsurer(ctx context.Context, reinsurerID string) (repository.DampakHapusTCO, error)
	HapusReinsurer(ctx context.Context, tx *repository.Tx, kom models.KombinasiTCO, reinsurerID string) (repository.DampakHapusTCO, error)
}

type kaskadeBelumDisuntik struct{}

func (kaskadeBelumDisuntik) DampakKontrak(context.Context, models.KombinasiTCO, string, string) (repository.DampakHapusTCO, error) {
	return repository.DampakHapusTCO{}, ErrGudangKaskadeBelumDisuntik
}
func (kaskadeBelumDisuntik) HapusKontrak(context.Context, *repository.Tx, models.KombinasiTCO, string, string) (repository.DampakHapusTCO, error) {
	return repository.DampakHapusTCO{}, ErrGudangKaskadeBelumDisuntik
}
func (kaskadeBelumDisuntik) DampakReinsurer(context.Context, string) (repository.DampakHapusTCO, error) {
	return repository.DampakHapusTCO{}, ErrGudangKaskadeBelumDisuntik
}
func (kaskadeBelumDisuntik) HapusReinsurer(context.Context, *repository.Tx, models.KombinasiTCO, string) (repository.DampakHapusTCO, error) {
	return repository.DampakHapusTCO{}, ErrGudangKaskadeBelumDisuntik
}

// KaskadeOracle menyusun pelaksana kaskade di atas Oracle.
func KaskadeOracle(svc *Service) PelaksanaKaskadeTCO { return repository.NewKaskadeTCO(svc.db) }

// PerekamJejakKaskadeOracle - `SisipJejakTCO` di atas Oracle.
func PerekamJejakKaskadeOracle(svc *Service) func(ctx context.Context, tx *repository.Tx, akunID, tabel, barisID, aksi,
	keterangan string, waktu time.Time) error {
	return svc.db.SisipJejakTCO
}

// DampakTampil adalah isi popup konfirmasi.
type DampakTampil struct {
	Kontrak      int64 `json:"kontrak"`
	Reinsurer    int64 `json:"reinsurer"`
	Security     int64 `json:"security"`
	Business     int64 `json:"business"`
	KlausulTetap int64 `json:"klausulTetap"`
	// Bersama - kontrak lain yang memakai kombinasi yang sama (reinsurer/security tidak ikut).
	Bersama int64 `json:"bersama"`
}

// KonfirmasiHapus adalah jumlah yang pemakai lihat di popup lalu setujui.
type KonfirmasiHapus struct {
	Reinsurer, Security, Business int64
}

func tampilDampak(d repository.DampakHapusTCO) DampakTampil {
	return DampakTampil{Kontrak: d.Kontrak, Reinsurer: d.Reinsurer, Security: d.Security, Business: d.Business,
		KlausulTetap: d.KlausulTetap, Bersama: d.Bersama}
}

// KaskadeTCO melayani popup + kaskade hapus.
type KaskadeTCO struct {
	svc       *Service
	kaskade   PelaksanaKaskadeTCO
	kontrak   PemegangKontrakTCO
	tahun     PemeriksaTahunTCO
	reinsurer PembacaReinsurerTCO
	jejak     func(ctx context.Context, tx *repository.Tx, akunID, tabel, barisID, aksi, keterangan string, waktu time.Time) error
	jam       func() time.Time
	transaksi func(ctx context.Context, fn func(tx *repository.Tx) error) error
}

// KaskadeTCO menyusun layanannya; bawaannya gagal terang.
func (s *Service) KaskadeTCO() *KaskadeTCO {
	return &KaskadeTCO{svc: s, kaskade: kaskadeBelumDisuntik{}, kontrak: pemegangKontrakBelumDisuntik{},
		tahun: gudangTahunTreatyBelumDisuntik{}, reinsurer: reinsurerBelumDisuntik{},
		jejak: func(context.Context, *repository.Tx, string, string, string, string, string, time.Time) error {
			return ErrGudangKaskadeBelumDisuntik
		}, jam: time.Now, transaksi: s.DalamTransaksi}
}

func (l *KaskadeTCO) salin() *KaskadeTCO { s := *l; return &s }

// DenganKaskade memasang pelaksana kaskade.
func (l *KaskadeTCO) DenganKaskade(k PelaksanaKaskadeTCO) *KaskadeTCO {
	s := l.salin()
	s.kaskade = k
	return s
}

// DenganKontrak memasang pemegang kontrak.
func (l *KaskadeTCO) DenganKontrak(k PemegangKontrakTCO) *KaskadeTCO {
	s := l.salin()
	s.kontrak = k
	return s
}

// DenganTahun memasang pemeriksa tahun treaty.
func (l *KaskadeTCO) DenganTahun(t PemeriksaTahunTCO) *KaskadeTCO {
	s := l.salin()
	s.tahun = t
	return s
}

// DenganReinsurer memasang pembaca reinsurer.
func (l *KaskadeTCO) DenganReinsurer(r PembacaReinsurerTCO) *KaskadeTCO {
	s := l.salin()
	s.reinsurer = r
	return s
}

// DenganJejak memasang perekam jejak.
func (l *KaskadeTCO) DenganJejak(j func(ctx context.Context, tx *repository.Tx, akunID, tabel, barisID, aksi,
	keterangan string, waktu time.Time) error) *KaskadeTCO {
	s := l.salin()
	s.jejak = j
	return s
}

// DenganJam mengganti sumber waktu - dipakai uji.
func (l *KaskadeTCO) DenganJam(j func() time.Time) *KaskadeTCO { s := l.salin(); s.jam = j; return s }

// DenganTransaksi mengganti pelaksana transaksi - dipakai uji.
func (l *KaskadeTCO) DenganTransaksi(f func(ctx context.Context, fn func(tx *repository.Tx) error) error) *KaskadeTCO {
	s := l.salin()
	s.transaksi = f
	return s
}

func (l *KaskadeTCO) kombinasi(ctx context.Context, tahunID, kontrakID string) (models.KombinasiTCO, error) {
	t, err := l.tahun.Ambil(ctx, tahunID)
	if err != nil {
		return models.KombinasiTCO{}, err
	}
	k, err := l.kontrak.Ambil(ctx, tahunID, kontrakID)
	if err != nil {
		return models.KombinasiTCO{}, err
	}
	return models.KombinasiDari(t, k), nil
}

// DampakHapusKontrak - isi popup: yang ikut terhapus, dan klausul yang tetap.
func (l *KaskadeTCO) DampakHapusKontrak(ctx context.Context, pelaku Pelaku, tahunID, kontrakID string) (DampakTampil, error) {
	if err := WajibIdentitas(pelaku); err != nil {
		return DampakTampil{}, err
	}
	kom, err := l.kombinasi(ctx, tahunID, kontrakID)
	if err != nil {
		return DampakTampil{}, err
	}
	d, err := l.kaskade.DampakKontrak(ctx, kom, tahunID, kontrakID)
	if err != nil {
		return DampakTampil{}, err
	}
	return tampilDampak(d), nil
}

func periksaKonfirmasi(d repository.DampakHapusTCO, k KonfirmasiHapus) error {
	if d.Reinsurer != k.Reinsurer || d.Security != k.Security || d.Business != k.Business {
		return fmt.Errorf("%w: dikonfirmasi %d reinsurer, %d security, %d business; sekarang %d, %d, %d - tinjau ulang",
			ErrDampakBerubah, k.Reinsurer, k.Security, k.Business, d.Reinsurer, d.Security, d.Business)
	}
	return nil
}

// HapusKontrak - kaskade kontrak -> business, security, reinsurer dalam SATU
// transaksi (AC 42, 45); klausul tidak disentuh (AC 44).
func (l *KaskadeTCO) HapusKontrak(ctx context.Context, pelaku Pelaku, tahunID, kontrakID string, k KonfirmasiHapus) (string, error) {
	if err := WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	var kom models.KombinasiTCO
	err := l.transaksi(ctx, func(tx *repository.Tx) error {
		c := repository.DenganBacaTxTCO(ctx, tx)
		if err := l.kontrak.Kunci(c, tx, tahunID, kontrakID); err != nil {
			return err
		}
		// Temuan /code-review: kombinasi dibaca SESUDAH kunci, lewat transaksi -
		// bukan salinan basi dari sebelum penulis lain mengubah kontraknya.
		var err error
		if kom, err = l.kombinasi(c, tahunID, kontrakID); err != nil {
			return err
		}
		sekarang, err := l.kaskade.DampakKontrak(c, kom, tahunID, kontrakID)
		if err != nil {
			return err
		}
		if err := periksaKonfirmasi(sekarang, k); err != nil {
			return err
		}
		terhapus, err := l.kaskade.HapusKontrak(c, tx, kom, tahunID, kontrakID)
		if err != nil {
			return err
		}
		if err := periksaKonfirmasi(terhapus, k); err != nil {
			return err
		}
		return l.jejak(c, tx, pelaku.AkunID, repository.TabelKontrakTCO, kontrakID, repository.AksiJejakHapus,
			fmt.Sprintf("kontrak dihapus kombinasi %s/%s/%s beserta %d reinsurer, %d security, %d business; %d klausul tidak disentuh",
				kom.TreatyYear, kom.TreatyGroupID, kom.ReinsTypeID, terhapus.Reinsurer, terhapus.Security, terhapus.Business,
				sekarang.KlausulTetap), l.jam())
	})
	if err != nil {
		return "", err
	}
	return PesanKontrakTerhapusTCO, nil
}

// DampakHapusReinsurer - popup hapus reinsurer: security yang ikut terhapus.
func (l *KaskadeTCO) DampakHapusReinsurer(ctx context.Context, pelaku Pelaku, tahunID, kontrakID, reinsurerID string) (
	DampakTampil, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return DampakTampil{}, err
	}
	kom, err := l.kombinasi(ctx, tahunID, kontrakID)
	if err != nil {
		return DampakTampil{}, err
	}
	if _, err := l.reinsurer.Ambil(ctx, kom, reinsurerID); err != nil {
		return DampakTampil{}, err
	}
	d, err := l.kaskade.DampakReinsurer(ctx, reinsurerID)
	if err != nil {
		return DampakTampil{}, err
	}
	return tampilDampak(d), nil
}

// HapusReinsurer - `DeleteFromTreatyReinsurer_Act`: security lalu reinsurer,
// satu transaksi, angka dikonfirmasi.
func (l *KaskadeTCO) HapusReinsurer(ctx context.Context, pelaku Pelaku, tahunID, kontrakID, reinsurerID string,
	k KonfirmasiHapus) (string, error) {

	if err := WajibIdentitas(pelaku); err != nil {
		return "", err
	}
	err := l.transaksi(ctx, func(tx *repository.Tx) error {
		c := repository.DenganBacaTxTCO(ctx, tx)
		if err := l.kontrak.Kunci(c, tx, tahunID, kontrakID); err != nil {
			return err
		}
		kom, err := l.kombinasi(c, tahunID, kontrakID)
		if err != nil {
			return err
		}
		if _, err := l.reinsurer.Ambil(c, kom, reinsurerID); err != nil {
			return err
		}
		sekarang, err := l.kaskade.DampakReinsurer(c, reinsurerID)
		if err != nil {
			return err
		}
		if err := periksaKonfirmasi(sekarang, k); err != nil {
			return err
		}
		terhapus, err := l.kaskade.HapusReinsurer(c, tx, kom, reinsurerID)
		if err != nil {
			return err
		}
		if err := periksaKonfirmasi(terhapus, k); err != nil {
			return err
		}
		return l.jejak(c, tx, pelaku.AkunID, repository.TabelReinsurerTCO, reinsurerID, repository.AksiJejakHapus,
			fmt.Sprintf("reinsurer dihapus beserta %d security", terhapus.Security), l.jam())
	})
	if err != nil {
		return "", err
	}
	// `DeleteTreatyReins_Act` tidak menampilkan pesan - `[tidak ada di korpus]`.
	return fmt.Sprintf("Reinsurer dengan ID %s dihapus beserta %d security", reinsurerID, k.Security), nil
}
