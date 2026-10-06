// Package services memegang aturan modul Treaty In (`treatyin`).
//
// Arah ketergantungan: handlers -> services -> repository.
package services

import (
	"context"
	"errors"
	"fmt"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
)

// ErrHimpunanTidakAda - himpunan acuan yang diminta bukan salah satu dari enam.
var ErrHimpunanTidakAda = errors.New("himpunan acuan tidak dikenal")

// Service membawa akar layanan modul ini.
type Service struct{ *inti.Dasar }

// DariDasar merakit Service di atas akar yang perakit sediakan.
func DariDasar(d *inti.Dasar) *Service { return &Service{Dasar: d} }

// New membuat Service; db boleh nil (proses tanpa Oracle - rute menjawab 503).
func New(d *db.DB) *Service { return &Service{Dasar: inti.NewDasar(d)} }

// PunyaDatabase aman dipanggil pada nil.
func (s *Service) PunyaDatabase() bool { return s != nil && s.Dasar.PunyaDatabase() }

// Gudang adalah seluruh sentuhan basis data yang Layanan butuhkan.
type Gudang interface {
	DaftarAcuan(ctx context.Context, h models.Himpunan) ([]models.Acuan, error)

	// Tiket 14. Keduanya menulis dan membaca KONTRAK + VERSI_KONTRAK.
	// Pembuatan berjalan dalam SATU transaksi: kontrak tanpa versi pertamanya
	// bukan keadaan yang sah menurut tiket 14, dan dua pernyataan terpisah
	// dapat meninggalkannya bila yang kedua gagal.
	BuatKontrakDenganVersiPertama(ctx context.Context, k models.Kontrak, v models.VersiKontrak) (int64, int64, error)
	BacaKontrak(ctx context.Context, id int64) (models.KontrakDenganVersi, error)

	// Tiket 16, 17, 18, 19 - identitas kontrak.
	CariKontrakSerupa(ctx context.Context, k models.Kontrak) ([]int64, error)
	CariKontrakLewatNomorWarisan(ctx context.Context, nomor string) ([]models.Kontrak, error)
	PerbaruiKontrak(ctx context.Context, k models.Kontrak) error
	TambahVersi(ctx context.Context, idKontrak int64, v models.VersiKontrak) (int64, error)

	// Tiket 35 dan 36 - ketentuan proporsional.
	AdaQuotaSharePadaVersi(ctx context.Context, idVersi int64) (bool, error)
	CatatKetentuanProporsional(ctx context.Context, idVersi, idLayer, idKelompok int64, jenis, persenQS string, lines *int64) error

	// Tiket 32 - syarat berbeda tiap pemulihan limit. SELURUH daftar sekaligus:
	// nomor urut kembar hanya terlihat bila barisnya dilihat bersama.
	CatatPemulihanLimit(ctx context.Context, idLayer int64, baris []models.PemulihanLimit) error

	// Tiket 40 - versi DASAR dibaca lewat ID_VERSI_KONTRAK_DASAR, bukan dari
	// salinan. Nil tanpa galat berarti versinya yang pertama.
	BacaVersiDasar(ctx context.Context, idVersi int64) (*models.VersiKontrak, error)

	// Tiket 41 TIDAK menambah apa pun di sini: identitas bentuk lama
	// DITURUNKAN di services dari BacaKontrak. Menambah kueri tersendiri
	// berarti dua pembaca untuk kolom yang sama (INV-60).

	// Tiket 42 - arsip muatan keluar. Perhatikan bentuk nilai baliknya:
	// yang kedua `models.BuktiArsip`, yang TIDAK punya ruas muatan. Seam ini
	// karena itu tidak dapat mengembalikan isi arsip tanpa seseorang
	// menyunting struct-nya - INV-61 diwujudkan, bukan sekadar ditulis.
	SimpanArsipMuatanKeluar(ctx context.Context, a models.ArsipMuatanKeluar) error
	BuktiArsipKontrak(ctx context.Context, idKontrak int64) (models.BuktiArsip, error)

	// Layar daftar kontrak — ronde layar 1. Kesembilan kolomnya dibaca dari
	// `Section/InputTreatyInOffer.xml`; jejaknya per kolom di
	// `models.BarisDaftarKontrak`.
	DaftarKontrak(ctx context.Context) ([]models.BarisDaftarKontrak, error)

	// Layar daftar dari tabel WARISAN `TREATY_IN` — jalur TERPISAH, dan
	// sengaja: yang di atas membaca model BARU (hari ini nol baris), kedua
	// di bawah membaca sistem LAMA (1.854 baris). Menyatukannya menghapus
	// perbedaan itu tepat ketika ia mulai penting (tiket 59).
	//
	// ⛔ BACA SAJA. Nol penulisan terhadap `TREATY_IN` di seluruh seam ini.
	CacahKontrakWarisan(ctx context.Context) (int, error)
	DaftarKontrakWarisan(ctx context.Context, offset, batas int) ([]models.BarisDaftarWarisan, error)
	BacaKontrakWarisan(ctx context.Context, id string) (models.KontrakWarisan, error)

	// Tab yang sudah PINDAH dari `JSONDATA` ke tabel pendaratan migrasi 430.
	//
	// ⛔ Ketiganya seam TERPISAH dari `BacaKontrakWarisan`, dan itu
	// disengaja: kontrak yang dokumennya hilang tetap harus membuka tabnya,
	// dan tab yang tabelnya kosong tetap harus membuka kontraknya. Satu
	// panggilan yang mengembalikan semuanya mengikat kedua kegagalan itu
	// menjadi satu.
	BacaPeriodePelaporan(ctx context.Context, masterID string) ([]models.BarisPeriodeWarisan, error)
	BacaPortofolio(ctx context.Context, masterID string) ([]models.BarisPortofolioWarisan, error)
	BacaAkumulasi(ctx context.Context, masterID string) ([]models.BarisAkumulasiWarisan, error)

	// ⭐ KEEMPAT TAB — Limits · Share · Event Limits · RNM Share — dibaca
	// dari TABEL PENDARATAN sejak 6 Oktober 2026, keputusan pemilik proses:
	// nilai dari JSONDATA dilarang keras.
	//
	// Dua seam, bukan satu: bentuk PIPIH melayani empat tab sekaligus,
	// bentuk POHON melayani tab Limits proporsional. Keduanya membaca
	// tabel yang sama tetapi merangkainya berbeda, dan menyatukannya akan
	// memaksa yang satu membongkar bentuk yang lain.
	BacaLayerPendaratan(ctx context.Context, masterID string) ([]models.BarisLayerWarisan, error)
	BacaPohonLimitsPendaratan(ctx context.Context, masterID string) ([]map[string]any, error)

	// ⭐ MEDAN KEPALA dari tabel pendaratan, 6 Oktober 2026. Sesudahnya NOL
	// nilai layar Treaty In datang dari `M_TREATY_IN.JSONDATA`.
	//
	// ⛔ Grid Rate of Exchange TIDAK ikut mendarat: sumbernya
	// `TREATYEXCHANGEYEARLY`, tabel warisan yang sudah hidup — keputusan
	// pemilik proses 4 dan 6 Oktober 2026. Ia berkunci TAHUN, bukan kontrak,
	// jadi seamnya menerima `TREATYYEAR` dan bukan `MASTERID`.
	BacaRevisiPendaratan(ctx context.Context, masterID string) (repository.RevisiPendaratan, error)
	BacaKursTahunan(ctx context.Context, tahunTreaty string) ([]models.BarisKursWarisan, error)

	// Tab Co-Ins Scale - tabel pendaratan kesembilan, migrasi 432.
	BacaSkalaKoasuransi(ctx context.Context, masterID string) ([]models.BarisSkalaKoasuransiWarisan, error)

	// Panel Attachment - `M_ATTACHMENTTREATY_2`, tabel WARISAN, BACA SAJA.
	// Katalognya seam TERPISAH: ia tidak bergantung kontrak, dan satu
	// pembacaan melayani seluruh panel.
	BacaLampiranKontrak(ctx context.Context, masterID string) ([]models.BarisLampiranWarisan, error)
	BacaKatalogKategoriLampiran(ctx context.Context) (map[string]string, error)
	BacaEgnpi(ctx context.Context, masterID string) ([]models.BarisEgnpiWarisan, error)
	BacaRetensi(ctx context.Context, masterID string) ([]models.BarisRetensiWarisan, error)
	// ⭐ Panel `Existing Policy for Master ID` — `TREATYINPRODUCTION`, baca-saja.
	BacaPolisProduksi(ctx context.Context, masterID string) ([]models.BarisPolisProduksi, error)
	BacaAngsuran(ctx context.Context, masterID string) ([]models.BarisAngsuranWarisan, error)
	BacaCatatan(ctx context.Context, masterID string) ([]models.BarisCatatanWarisan, error)

	// Isi pemilih "Choose Ceding" dan "Choose Source of Business".
	//
	// ⛔ Keduanya seam TERPISAH dan tidak bergantung kontrak mana pun: satu
	// pembacaan melayani seluruh pemilih, dan menautkannya ke
	// `BacaKontrakWarisan` berarti menarik 94 baris setiap kali satu kontrak
	// dibuka.
	BacaDaftarCedant(ctx context.Context) ([]models.PilihanWarisan, error)
	BacaDaftarAsalBisnis(ctx context.Context) ([]models.PilihanWarisan, error)
	// Isi dropdown `Treaty Type` tab Limits — `BrowseReinsuranceType_RD`.
	BacaDaftarJenisTreaty(ctx context.Context) ([]models.PilihanWarisan, error)
	// `Treaty Group` dan mata uang tab Limits — `BrowseTreatyGroup_RD`,
	// `BrowseCurrencyTreatyIn_RD`.
	BacaDaftarKelompokTreaty(ctx context.Context) ([]models.PilihanWarisan, error)
	BacaDaftarMataUangLimit(ctx context.Context) ([]models.PilihanWarisan, error)
}

// Layanan memegang aturan modul ini di atas satu Gudang.
type Layanan struct{ gudang Gudang }

// LayananOracle merakit Layanan di atas Oracle.
func LayananOracle(s *Service) *Layanan { return &Layanan{gudang: repository.Baru(s.DB())} }

// LayananDengan merakit Layanan di atas Gudang mana pun - dipakai uji.
func LayananDengan(g Gudang) *Layanan { return &Layanan{gudang: g} }

// himpunanSah adalah KELIMA himpunan acuan, DISEBUT satu per satu.
// `mata-uang` dicabut 4 Oktober 2026 — migrasi 434.
var himpunanSah = map[models.Himpunan]bool{
	models.HimpunanJenisPotongan:   true,
	models.HimpunanKelasBisnis:     true,
	models.HimpunanKelompokTreaty:  true,
	models.HimpunanBahaya:          true,
	models.HimpunanJenisReasuransi: true,
}

// DaftarAcuan membaca isi satu himpunan acuan.
//
// Identitas pelaku diperiksa lebih dulu: tidak ada jalur baca tanpa identitas,
// termasuk untuk data acuan.
func (l *Layanan) DaftarAcuan(ctx context.Context, p inti.Pelaku, h models.Himpunan) ([]models.Acuan, error) {
	if err := inti.WajibIdentitas(p); err != nil {
		return nil, err
	}
	if !himpunanSah[h] {
		return nil, fmt.Errorf("%w: %q", ErrHimpunanTidakAda, h)
	}
	return l.gudang.DaftarAcuan(ctx, h)
}

// Pesan mengembalikan kalimat yang layak sampai ke layar.
func Pesan(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
