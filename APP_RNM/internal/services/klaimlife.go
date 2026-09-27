package services

import (
	"context"
	"errors"
	"fmt"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
)

// ErrKlaimTidakAda dikembalikan bila klaim yang diminta tidak ada.
var ErrKlaimTidakAda = errors.New("services: klaim tidak ada")

// KlaimLife merakit agregat klaim Life dari bacaan repository.
//
// Perakitan ada DI SINI, bukan di repository dan bukan di handlers.
type KlaimLife struct {
	repo *repository.KlaimLife
	diag *repository.Diagnosa
}

// KlaimLife mengembalikan layanan klaim Life, atau nil bila tanpa database.
func (s *Service) KlaimLife() *KlaimLife {
	if !s.PunyaDatabase() {
		return nil
	}
	return &KlaimLife{
		repo: repository.NewKlaimLife(s.db),
		diag: repository.NewDiagnosa(s.db),
	}
}

// Ambil merakit satu klaim beserta seluruh peserta dan seluruh baris
// adjustment-nya.
//
// Tiket 01 AC-1: seluruh baris `AdjustmentList`, bukan hanya yang terakhir.
// Setiap baris tetap melekat pada pesertanya; meratakannya ke header menghapus
// informasi peserta pemilik dan mematahkan mesin status (ADR-U-0011).
func (k *KlaimLife) Ambil(ctx context.Context, id string) (*models.Klaim, error) {
	if k == nil || k.repo == nil {
		return nil, repository.ErrTanpaOracle
	}
	if id == "" {
		return nil, fmt.Errorf("%w: pengenal klaim kosong", ErrWajibIsi)
	}

	klaim, err := k.repo.AmbilHeader(ctx, id)
	if err != nil {
		return nil, err
	}
	if klaim == nil {
		return nil, ErrKlaimTidakAda
	}

	peserta, err := k.repo.AmbilPeserta(ctx, id)
	if err != nil {
		return nil, err
	}
	baris, err := k.repo.AmbilBaris(ctx, id)
	if err != nil {
		return nil, err
	}
	// Dokumen pendukung per peserta - `LoadDocumentLife_ACT` (lihat
	// repository.AmbilDokumen untuk pohon langkahnya dan dua penyimpangan
	// sadarnya).
	dokumen, err := k.repo.AmbilDokumen(ctx, id)
	if err != nil {
		return nil, err
	}
	// Diagnosa per peserta - butir bd. Satu query untuk seluruh klaim,
	// sebentuk dengan dokumen di atasnya.
	diagnosa, err := k.diag.AmbilDiagnosa(ctx, id)
	if err != nil {
		return nil, err
	}

	for i := range peserta {
		if b, ada := baris[peserta[i].ID]; ada {
			peserta[i].Baris = b
		}
		if d, ada := dokumen[peserta[i].ID]; ada {
			peserta[i].Dokumen = d
		}
		if g, ada := diagnosa[peserta[i].ID]; ada {
			peserta[i].Diagnosa = g
		}
		// Keenam total uang peserta dihitung DI SINI, saat dibaca, dan
		// tidak disimpan (models.HitungTotalPeserta punya bukti XML-nya).
		//
		// Di luar `if ada` dengan sengaja: peserta TANPA baris adjustment
		// tetap bertotal nol, sebab penampung Pega mulai dari literal 0
		// (`SavePesertaClaim.xml` b4027..b4159) dan langkah 8.2 b4592
		// menuliskannya apa adanya. Menaruhnya di dalam `if` akan membuat
		// peserta itu bertotal KOSONG, dan layar akan berkata "belum ada
		// datanya" untuk peserta yang datanya lengkap dan berjumlah nol.
		total, err := models.HitungTotalPeserta(peserta[i].Baris, peserta[i].MataUang)
		if err != nil {
			return nil, fmt.Errorf("peserta %s: %w", peserta[i].NomorSertifikat, err)
		}
		peserta[i].Total = total
	}
	klaim.Peserta = peserta

	// ⭐ BUTIR bb: tahap dan status kerja ikut menyeberang, sebab layar Detail
	// memerlukan keduanya untuk memutuskan apakah `Close Claim` pantas
	// ditawarkan. Keduanya hidup di T_WORK_CLAIM, bukan di header klaim.
	//
	// ⚠️ Baris work yang TIDAK ADA bukan galat. Aplikasi ini tidak pernah
	// menyisipkan ke T_WORK_CLAIM - baris itu lahir di sistem lama - jadi
	// klaim tanpa baris work adalah keadaan nyata. Ia menjadi tahap kosong,
	// dan tahap kosong tidak menawarkan tombol apa pun (gagal TERTUTUP).
	// Galat LAIN tetap menggagalkan pembacaan: menelan galat basis data di
	// sini akan membuat layar diam-diam menyembunyikan tombol yang
	// seharusnya ada, dan tidak ada yang tahu kenapa.
	tahap, _, err := k.repo.TahapDanPeran(ctx, id)
	if err != nil && !errors.Is(err, repository.ErrWorkTidakAda) {
		return nil, err
	}
	if err == nil {
		klaim.Tahap = tahap
		status, err := k.repo.StatusWorkKlaim(ctx, id)
		if err != nil && !errors.Is(err, repository.ErrWorkTidakAda) {
			return nil, err
		}
		klaim.StatusWork = status
	}
	return klaim, nil
}

// JalankanMigrasi membentuk tabel di basis data yang dikonfigurasi.
//
// Ia dipanggil oleh `go run ./cmd/api -migrate` (target `make migrate`).
// Pelarinya menolak berjalan bila lingkungan menunjuk produksi Pega
// (ADR-U-0005), dan aman dijalankan berulang kali.
func (s *Service) JalankanMigrasi(ctx context.Context) (repository.LaporanMigrasi, error) {
	if !s.PunyaDatabase() {
		return repository.LaporanMigrasi{}, repository.ErrTanpaOracle
	}
	return s.db.JalankanMigrasi(ctx)
}

// BongkarMigrasi menjalankan jalur mundur tiap langkah yang TERCATAT selesai.
//
// ⛔ Ia MENGHAPUS tabel. Sampai 26-09-2026 satu-satunya pemanggilnya adalah
// skema uji, sehingga orang yang ingin membongkar skema uji sendiri terpaksa
// menyalin isi berkas *_down.sql ke sqlplus - dan itu melewati pengaman
// T_MIGRASI, yang hanya membongkar langkah yang benar-benar tercatat.
//
// Pemanggilnya WAJIB memagari lebih dulu lewat Config.PastikanSkemaUji.
// Lapisan ini tidak membaca environment sendiri.
func (s *Service) BongkarMigrasi(ctx context.Context) (repository.LaporanMigrasi, error) {
	if !s.PunyaDatabase() {
		return repository.LaporanMigrasi{}, repository.ErrTanpaOracle
	}
	return s.db.BongkarMigrasi(ctx)
}
