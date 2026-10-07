// Package services memegang aturan dagang EDM Treaty In di antara HTTP dan Oracle: siapa boleh apa (keanggotaan
// antrean), urutan langkah flow action Pega (pra-proses -> layar -> pasca DT -> pasca aktivitas -> connector), dan
// SATU transaksi per tindakan. Arah ketergantungan `handlers -> services -> repository`. Pola: salinan
// `modul/nbtreatyin/backend/services` (06-10-2026), isinya menurut korpus EDM.
//
// Untuk apa berkas ini: ANTARMUKA GUDANG. Layanan hanya mengenal antarmuka ini; implementasi Oracle-nya
// `repository.Gudang` dibungkus `penyimpanOracle` (transaksi dari `inti.Dasar`). Uji seam HTTP memasang gudang
// tiruan (`backend/tiruan`).
package services

import (
	"context"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/edmtreatyin/backend/models"
	"nusantarare/modul/edmtreatyin/backend/repository"
)

// Gudang adalah seluruh kebutuhan penyimpanan layanan.
type Gudang interface {
	// Transaksi menjalankan fn di satu transaksi; galat membatalkan semuanya.
	Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error

	IDKasusBerikut(ctx context.Context, tx *db.Tx) (string, error)
	SisipKasus(ctx context.Context, tx *db.Tx, id, pembuat, namaPembuat string, gb repository.GenerasiBaru) error
	Keadaan(ctx context.Context, tx *db.Tx, id string) (models.Kasus, error)
	KunciKasus(ctx context.Context, tx *db.Tx, id, statusHarap string) error
	PindahPosisi(ctx context.Context, tx *db.Tx, id, statusLama, posisiBaru, position string) error
	TutupKasus(ctx context.Context, tx *db.Tx, id, statusLama, statusAkhir string) error
	DaftarKasus(ctx context.Context, s models.SaringanKasus) ([]models.RingkasanKasus, error)

	// rantai generasi (repository/generasi.go)
	CacahGenerasiJSONPolis(ctx context.Context, nopolis string) (int, error)
	NoMasterDariNoPolis(ctx context.Context, nopolis string) (string, error)
	AdaEDMBerjalan(ctx context.Context, nopolis string) (bool, error)
	GenerasiTerakhir(ctx context.Context, tx *db.Tx, nopolis string) (repository.GenerasiPolis, error)
	SetelNomorPolisSelesai(ctx context.Context, tx *db.Tx, id, nopolis string) error
	LepasGenerasi(ctx context.Context, tx *db.Tx, id string) error
	BuangSelisihGenerasi(ctx context.Context, tx *db.Tx, polisID string) error

	// halaman (repository/polis.go, halaman_edm.go, selisih.go)
	SimpanHalaman(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error
	BacaHalaman(ctx context.Context, tx *db.Tx, id string) (*models.Halaman, error)
	BacaGenerasi(ctx context.Context, tx *db.Tx, id string) (*models.Halaman, error)
	SimpanSelisih(ctx context.Context, tx *db.Tx, polisID string, h *models.Halaman, k repository.KunciSelisih) error
	HariClosing(ctx context.Context, tx *db.Tx) (int, error)

	// acuan (repository/acuan.go, bisnis_edm.go)
	DaftarBisnisEDM(ctx context.Context, s repository.SaringanPopupEDM) ([]models.Baris, error)
	IDMataUangDariNama(ctx context.Context, nama string) (string, error)
	NamaMataUang(ctx context.Context, id string) (string, error)
	StsPKPAgen(ctx context.Context, sobID string) (string, error)
	BisnisDariKunci(ctx context.Context, kunci string) (models.BarisBisnis, error)
	MO(ctx context.Context, id string) (models.BarisMO, error)
	DaftarMataUang(ctx context.Context) ([]models.Pilihan, error)
	DaftarMO(ctx context.Context) ([]models.Pilihan, error)
	DaftarJenisReas(ctx context.Context) ([]models.Pilihan, error)

	// PembacaMaster - master kontrak treaty (JSONDATA M_TREATY_IN / M_TREATY_IN_EDM / M_TREATY_OUT, pengecualian
	// K8 NB - repository/masterxol.go).
	PembacaMaster

	// SimpanPolisProduksi - Utility1 `SaveJsonPolisTreatyInEDM_Act`: json_polis (tanpa DATA_JSON), ACHIEVEMENT,
	// TREATYINPRODUCTION (models/produksi.go).
	SimpanPolisProduksi(ctx context.Context, tx *db.Tx, s models.SimpananPolis) error
	CatatRiwayat(ctx context.Context, tx *db.Tx, r models.Riwayat) error
	// CatatUsulan menulis catatan SuggestList ke POOLDATA.HISTORYAKSEPTASIPRODUCTION (ketetapan NB K4).
	CatatUsulan(ctx context.Context, tx *db.Tx, idPega string, baris []models.UsulanProduksi) error
	NamaTampilan(ctx context.Context, loginID string) (string, error)
	// PemegangKotakMasuk - pemegang aktif workbasket tujuan Submit, untuk nama di NBStatus.
	PemegangKotakMasuk(ctx context.Context, workbasket string) (models.PemegangKotakMasuk, error)
	// HitungKotakMasuk - cacah berkas yang menunggu akun per posisi (kotak masuk Beranda).
	HitungKotakMasuk(ctx context.Context, akun string, admin bool, atasan []string) (map[string]int, error)
}

// PembacaMaster adalah sumber master kontrak treaty.
type PembacaMaster interface {
	// MasterXOL membaca master satu kontrak menurut ID (`SetTreatyIn_Act` / RDB BrowseTreatyIn: M_TREATY_IN ∪
	// M_TREATY_IN_EDM) - pra-proses `InputPolicyTreatyInPre_Act` langkah 10 (TreatyRealizationCheckXOLList).
	MasterXOL(ctx context.Context, noKontrak string) (models.MasterXOL, error)
	// MasterEDM - pembaca rantai `EDMChooseBusiness_Act` (OLDID / ID / EDM ID / M_TREATY_OUT, mata uang).
	MasterEDM(ctx context.Context) models.PembacaMasterEDM
}

// penyimpanOracle - `repository.Gudang` + transaksi `inti.Dasar`.
type penyimpanOracle struct {
	*repository.Gudang
	dasar *inti.Dasar
}

// Transaksi = `inti.Dasar.DalamTransaksi`.
func (p penyimpanOracle) Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error {
	return p.dasar.DalamTransaksi(ctx, fn)
}

// MasterXOL = `repository.MasterXOLDariJSON` (K8).
func (p penyimpanOracle) MasterXOL(ctx context.Context, noKontrak string) (models.MasterXOL, error) {
	return p.Gudang.MasterXOLDariJSON(ctx, noKontrak)
}

// MasterEDM = `repository.PembacaMasterEDM`.
func (p penyimpanOracle) MasterEDM(ctx context.Context) models.PembacaMasterEDM {
	return p.Gudang.PembacaMasterEDM(ctx)
}

// DariDasar menyusun layanan di atas basis data bersama. Tanpa Oracle (`PunyaDatabase` palsu) layanan menjawab
// `ErrTanpaOracle` - 503.
func DariDasar(d *inti.Dasar) *Layanan {
	if d == nil || !d.PunyaDatabase() {
		return Baru(nil, time.Now)
	}
	l := Baru(penyimpanOracle{Gudang: repository.Baru(d.DB()), dasar: d}, time.Now)
	return l.DenganKonversi(nil, d.Lingkungan().AdalahProduksi())
}
