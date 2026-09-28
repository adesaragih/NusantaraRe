package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	kurs *repository.MataUang
	// Butir bk: ambang `MAXEXPIREDCLAIM` penanda `MAX CLAIM RECEIVED`.
	polis  *repository.RingkasPolisLife
	produk *repository.ProdukLife
}

// KlaimLife mengembalikan layanan klaim Life, atau nil bila tanpa database.
func (s *Service) KlaimLife() *KlaimLife {
	if !s.PunyaDatabase() {
		return nil
	}
	return &KlaimLife{
		repo:   repository.NewKlaimLife(s.db),
		diag:   repository.NewDiagnosa(s.db),
		kurs:   repository.NewMataUang(s.db),
		polis:  repository.NewRingkasPolisLife(s.db),
		produk: repository.NewProdukLife(s.db),
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
		// ⛔ `SetCurrencyID_Act` - celah sensus 28-09-2026, kini ditutup.
		//
		// Rule itu adalah `pyPreDataTransform` `AdjustmentDetail_Section.xml`
		// b1206: ia menerjemahkan `.CURRENCY` menjadi `.CURRENCYID` lewat
		// `GetCurrencyID` b419 SETIAP KALI layar adjustment dimuat. Di sini
		// pun ia berjalan di jalur BACA, bukan jalur tulis - meniru letaknya,
		// bukan hanya hasilnya.
		//
		// ⚠️ Sebelum ini `CURRENCYID` hanya DIBAWA (`WarisiKolom` menyalinnya
		// dari baris sebelumnya) dan tidak pernah DITERBITKAN, sehingga baris
		// PERTAMA sebuah peserta lahir tanpa pengenal mata uang - dan
		// `HitungTotalPeserta` menolaknya dengan kalimat "mata uang beragam",
		// yaitu kalimat yang benar tentang hal yang salah.
		if err := k.lengkapiPengenalMataUang(ctx, &peserta[i]); err != nil {
			return nil, err
		}
		total, err := models.HitungTotalPeserta(peserta[i].Baris, peserta[i].MataUang)
		if err != nil {
			return nil, fmt.Errorf("peserta %s: %w", peserta[i].NomorSertifikat, err)
		}
		peserta[i].Total = total
	}
	// ⭐ BUTIR bk - `.MAXCLAIM_RECEIVED` DIHITUNG saat baca, di jalur yang
	// sama dengan total dan pengenal mata uang: di Pega ia hanya Property-Set
	// halaman (ValidasiClaimReceived_Act b582), nol penulis tabel.
	ambang, alasan, err := k.ambangTerimaKlaim(ctx, klaim.NomorPolis)
	if err != nil {
		return nil, err
	}
	isiPenandaTerimaKlaim(peserta, ambang, alasan)
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

// lengkapiPengenalMataUang mengisi `CURRENCYID` baris yang belum punya.
//
// ⛔ Hanya yang KOSONG. Baris yang sudah berpengenal tidak disentuh: nilai
// yang tersimpan adalah nilai yang pernah benar, dan menimpanya dengan hasil
// pencarian hari ini berarti mengubah data lama karena tabel rujukan berubah.
//
// ⚠️ Satu pencarian per KODE, bukan per baris. Peserta berbaris puluhan
// yang seluruhnya `IDR` tidak perlu puluhan query.
func (k *KlaimLife) lengkapiPengenalMataUang(ctx context.Context,
	p *models.Peserta) error {

	if k == nil || k.kurs == nil || p == nil {
		return nil
	}
	sudah := map[string]string{}
	for i := range p.Baris {
		b := &p.Baris[i]
		if strings.TrimSpace(b.CurrencyID) != "" {
			continue
		}
		kode := strings.TrimSpace(b.JumlahKlaim.Currency)
		if kode == "" {
			kode = strings.TrimSpace(p.MataUang)
		}
		if kode == "" {
			// Mata uang yang memang belum diisi bukan mata uang yang salah
			// (ADR-U-0027). Dibiarkan, dan total peserta yang menjumlahkannya
			// akan berkata apa adanya.
			continue
		}
		id, ada := sudah[kode]
		if !ada {
			var err error
			id, err = k.kurs.Pengenal(ctx, kode)
			if err != nil {
				return err
			}
			sudah[kode] = id
		}
		b.CurrencyID = id
	}
	return nil
}

// ambangTerimaKlaim membaca `MAXEXPIREDCLAIM` produk polis sebuah klaim.
//
// Polis atau produk yang TIDAK DITEMUKAN bukan kerusakan layar Detail: ia
// menjadi ALASAN (penanda tidak dihitung, dan itu dinyatakan). Galat basis
// data lain tetap menggagalkan pembacaan.
func (k *KlaimLife) ambangTerimaKlaim(ctx context.Context, nomorPolis string) (
	ambang, alasan string, err error) {

	if strings.TrimSpace(nomorPolis) == "" {
		return "", "klaim tanpa nomor polis; ambang MAXEXPIREDCLAIM tidak dapat dibaca", nil
	}
	polis, err := k.polis.Ringkas(ctx, nomorPolis)
	if errors.Is(err, repository.ErrPolisNomorTakDitemukan) {
		return "", "polis belum ada di PremiumList Life; ambang MAXEXPIREDCLAIM tidak dapat dibaca", nil
	}
	if err != nil {
		return "", "", err
	}
	if strings.TrimSpace(polis.ProductNameID) == "" {
		return "", "polis tidak menyebut produk; ambang MAXEXPIREDCLAIM tidak dapat dibaca", nil
	}
	a, err := k.produk.Ambang(ctx, polis.ProductNameID)
	if errors.Is(err, repository.ErrAmbangProdukTakDitemukan) {
		return "", "produk polis tidak ada di view produk; ambang MAXEXPIREDCLAIM tidak dapat dibaca", nil
	}
	if err != nil {
		return "", "", err
	}
	return a.MaxExpiredClaim, "", nil
}

// isiPenandaTerimaKlaim menulis `.MAXCLAIM_RECEIVED` tiap peserta - MURNI.
//
// `alasan` tak kosong berarti ambangnya tak terbaca: penanda TIDAK dihitung,
// dan alasannya ditulis - kosong tanpa alasan akan terbaca "sah".
func isiPenandaTerimaKlaim(peserta []models.Peserta, maxExpiredClaim, alasan string) {
	for i := range peserta {
		if alasan != "" {
			peserta[i].PenandaTerimaKlaimAlasan = alasan
			continue
		}
		penanda, err := models.PenandaTerimaKlaim(peserta[i], maxExpiredClaim)
		if err != nil {
			peserta[i].PenandaTerimaKlaimAlasan = err.Error()
			continue
		}
		peserta[i].PenandaTerimaKlaim = penanda
	}
}
