// Package services memegang aturan dagang NB Treaty In di antara HTTP dan
// Oracle: siapa boleh apa (keanggotaan antrean, AC 11-14), urutan langkah
// flow action Pega (pra-proses -> layar -> pasca DT -> pasca aktivitas ->
// connector), dan SATU transaksi per tindakan (AC 29, 83). Arah
// ketergantungan `handlers -> services -> repository` (AC 60).
//
// Untuk apa berkas ini: ANTARMUKA GUDANG. Layanan hanya mengenal antarmuka
// ini; implementasi Oracle-nya `repository.Gudang` dibungkus `penyimpanOracle`
// (transaksi dari `inti.Dasar`). Ia BUKAN jahitan uji tersendiri (spec §6.2:
// tiga jahitan saja): uji seam 1 (HTTP handler) memasang gudang tiruan di
// baliknya, uji seam 2 (`repository`, `-tags db`) memakai Oracle sungguhan.
package services

import (
	"context"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/nbtreatyin/backend/models"
	"nusantarare/modul/nbtreatyin/backend/repository"
)

// Gudang adalah seluruh kebutuhan penyimpanan layanan.
type Gudang interface {
	// Transaksi menjalankan fn di satu transaksi; galat membatalkan semuanya.
	Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error

	IDKasusBerikut(ctx context.Context, tx *db.Tx) (string, error)
	SisipKasus(ctx context.Context, tx *db.Tx, id, pembuat, namaPembuat string) error
	Keadaan(ctx context.Context, tx *db.Tx, id string) (models.Kasus, error)
	KunciKasus(ctx context.Context, tx *db.Tx, id, statusHarap string) error
	PindahPosisi(ctx context.Context, tx *db.Tx, id, statusLama, posisiBaru string) error
	TutupKasus(ctx context.Context, tx *db.Tx, id, statusLama, statusAkhir string) error
	DaftarKasus(ctx context.Context, s models.SaringanKasus) ([]models.RingkasanKasus, error)

	SimpanHalaman(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error
	BacaHalaman(ctx context.Context, tx *db.Tx, id string) (*models.Halaman, error)
	SetelNomorPolis(ctx context.Context, tx *db.Tx, id, nopol string) error
	TerbitkanNomorPolis(ctx context.Context, tx *db.Tx, h *models.Halaman, sekarang time.Time) (repository.BahanNomor, error)
	HariClosing(ctx context.Context, tx *db.Tx) (int, error)

	DetailKontrak(ctx context.Context, id string) (models.BarisKontrak, error)
	// DetailKontrakTreaty - baris view pertama ber-TREATYID itu (master berkas salinan dokumen Pega lama).
	DetailKontrakTreaty(ctx context.Context, treatyID string) (models.BarisKontrak, error)
	KomisiKontrak(ctx context.Context, treatyID string) ([]models.BarisKontrak, error)
	DaftarBisnis(ctx context.Context, s models.SaringanBisnis) ([]models.BarisKontrak, error)
	IDMataUangDariNama(ctx context.Context, nama string) (string, error)
	NamaMataUang(ctx context.Context, id string) (string, error)
	OJKGrupTreaty(ctx context.Context, grupID string) (string, error)
	OldIDGrupTreaty(ctx context.Context, grupID string) (string, error)
	KlienDariNama(ctx context.Context, nama string) (string, error)
	StsPKPAgen(ctx context.Context, sobID string) (string, error)
	BisnisDariKunci(ctx context.Context, kunci string) (models.BarisBisnis, error)
	MO(ctx context.Context, id string) (models.BarisMO, error)
	DaftarMataUang(ctx context.Context) ([]models.Pilihan, error)
	DaftarMO(ctx context.Context) ([]models.Pilihan, error)
	DaftarJenisSpreading(ctx context.Context) ([]models.Pilihan, error)
	DaftarJenisReas(ctx context.Context) ([]models.Pilihan, error)
	PolisSerupa(ctx context.Context, h *models.Halaman) ([]string, error)
	DaftarAgenHierarki(ctx context.Context) ([]models.BarisAgen, error)
	AgenHierarki(ctx context.Context, id string) (models.BarisAgen, bool, error)

	// PembacaMasterTreaty - master kontrak jalur XOL (K8).
	PembacaMasterTreaty

	// SimpanPolisProduksi - Utility1 `SaveJsonPolisTreatyIn_Act` sesudah realisasi selesai: json_polis (tanpa
	// DATA_JSON), ACHIEVEMENT, TREATYINPRODUCTION (`[keputusan work owner 06-10-2026]`, models/produksi.go).
	SimpanPolisProduksi(ctx context.Context, tx *db.Tx, s models.SimpananPolis) error
	CatatRiwayat(ctx context.Context, tx *db.Tx, r models.Riwayat) error
	// CatatUsulan menulis catatan SuggestList ke POOLDATA.HISTORYAKSEPTASIPRODUCTION
	// (SaveViewSuggest -> InsertViewSuggest_SQL), NOURUT berikutnya per IDPEGA.
	CatatUsulan(ctx context.Context, tx *db.Tx, idPega string, baris []models.UsulanProduksi) error
	DaftarRiwayat(ctx context.Context, idPega string) ([]models.Riwayat, error)
	NamaTampilan(ctx context.Context, loginID string) (string, error)
	// PemegangKotakMasuk - pemegang aktif workbasket tujuan Submit, untuk nama di NBStatus.
	PemegangKotakMasuk(ctx context.Context, workbasket string) (models.PemegangKotakMasuk, error)
	// HitungKotakMasuk - cacah berkas yang menunggu akun per posisi (kotak masuk Beranda).
	HitungKotakMasuk(ctx context.Context, akun string, admin bool, atasan []string) (map[string]int, error)
}

// PembacaMasterTreaty adalah sumber master kontrak treaty untuk jalur NB
// NonProporsional / XOL (`[keputusan work owner]` K8): `TreatyIn.Share()`,
// `Installment()`, `FacultativeShareList()`, `Limits()`, ... menurut
// `models.SkalarMasterXOL` / `models.DaftarMasterXOL`.
//
// Kini dipenuhi `repository.MasterXOLDariJSON` - baca-saja `JSONDATA`
// `M_TREATY_IN` / `M_TREATY_IN_EDM`, pengecualian sempit atas P29. Tanda
// tangannya sengaja tidak menyebut JSON: begitu tabel master modul `treatyin`
// (`KONTRAK`, `LAYER`, `BAGIAN`, `PEMULIHAN_LIMIT`, `TERMIN`, `POTONGAN`) terisi,
// antarmuka ini dipenuhi kontrak modul itu (PERMINTAAN-TIM-INTI bagian E).
type PembacaMasterTreaty interface {
	// MasterXOL membaca master satu kontrak menurut nomornya
	// (`PolicyTreatyIn.NoOffer`). Nol master = `ErrMasterXOLTidakAda`; dokumen
	// tak terurai = `ErrMasterXOLRusak`.
	MasterXOL(ctx context.Context, noKontrak string) (models.MasterXOL, error)
}

// penyimpanOracle - `repository.Gudang` + transaksi `inti.Dasar`.
type penyimpanOracle struct {
	*repository.Gudang
	dasar *inti.Dasar
}

// Transaksi = `inti.Dasar.DalamTransaksi` (ADR-U-0029).
func (p penyimpanOracle) Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error {
	return p.dasar.DalamTransaksi(ctx, fn)
}

// MasterXOL = `repository.MasterXOLDariJSON` (K8) - satu-satunya titik ganti
// ke kontrak `PembacaMasterTreaty` modul treatyin.
func (p penyimpanOracle) MasterXOL(ctx context.Context, noKontrak string) (models.MasterXOL, error) {
	return p.Gudang.MasterXOLDariJSON(ctx, noKontrak)
}

// DariDasar menyusun layanan di atas basis data bersama. Tanpa Oracle
// (`PunyaDatabase` palsu) layanan menjawab `ErrTanpaOracle` - 503.
func DariDasar(d *inti.Dasar) *Layanan {
	if d == nil || !d.PunyaDatabase() {
		return Baru(nil, time.Now)
	}
	l := Baru(penyimpanOracle{Gudang: repository.Baru(d.DB()), dasar: d}, time.Now)
	return l.DenganKonversi(nil, d.Lingkungan().AdalahProduksi())
}
