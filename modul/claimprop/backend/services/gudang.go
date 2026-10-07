// Package services memegang aturan dagang Claim Prop di antara HTTP dan Oracle: siapa boleh apa (pemegang assignment,
// tiket 10), urutan langkah flow action Pega (pra-proses -> layar -> aksi -> pasca-proses), dan SATU transaksi per
// aksi (prompt §6 butir 5: nomor + T_CLAIM_* + OS_AKSEPTASI_KLAIM/JSON_KLAIM + log + outbox). Arah ketergantungan
// `handlers -> services -> repository`.
//
// Untuk apa berkas ini: ANTARMUKA GUDANG dan ACUAN. Layanan hanya mengenal antarmuka ini; implementasi Oracle-nya
// `repository.Gudang` / `repository.Acuan`. Uji seam HTTP memasang gudang tiruan (`tiruan`).
package services

import (
	"context"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimprop/backend/models"
	"nusantarare/modul/claimprop/backend/repository"
)

// Gudang - seluruh tulisan dan bacaan kasus.
type Gudang interface {
	Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error

	IDKasusBerikut(ctx context.Context, tx *db.Tx, awalan string) (string, error)
	SisipKasus(ctx context.Context, tx *db.Tx, id, pembuat, namaPembuat string, saat time.Time) error
	Keadaan(ctx context.Context, tx *db.Tx, id string) (models.Kasus, error)
	KunciKasus(ctx context.Context, tx *db.Tx, id, tahap string) (models.Kasus, error)
	PindahTahap(ctx context.Context, tx *db.Tx, id, lama, baru, posisi string, saat time.Time) error
	TutupKasus(ctx context.Context, tx *db.Tx, id, tahap string, saat time.Time) error
	SentuhKasus(ctx context.Context, tx *db.Tx, id string, saat time.Time) error
	DaftarKasus(ctx context.Context, s repository.SaringanKasus) ([]repository.RingkasanKasus, error)

	SimpanHalaman(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error
	SetelKomiteAdjustment(ctx context.Context, tx *db.Tx, adjID, komiteID string) error
	BacaHalaman(ctx context.Context, tx *db.Tx, id string) (*models.Halaman, error)

	NomorSementara(ctx context.Context, tx *db.Tx, tahun string) (string, error)
	NomorCFS(ctx context.Context, tx *db.Tx, kodeBis string, saat time.Time) (string, error)
	UrutNomor(ctx context.Context, tx *db.Tx, huruf string, saat time.Time) (repository.BahanNomor, error)

	SisipOS(ctx context.Context, tx *db.Tx, b models.BarisOS, saat time.Time) error
	SalinJSONKlaim(ctx context.Context, tx *db.Tx, idPega, noKlaim, noPolis string, saat time.Time) error
	CatatLogLayanan(ctx context.Context, tx *db.Tx, l repository.LogLayanan, saat time.Time) error
	SisipKatastrofe(ctx context.Context, tx *db.Tx, k models.KatastrofeBaru, saat time.Time) error
	BuatKasusKomite(ctx context.Context, tx *db.Tx, klaimID, adjID, pembuat, namaPembuat string,
		anggota []repository.AnggotaTangga, saat time.Time) (string, error)
	AntreEfek(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) (string, error)
}

// Acuan - bacaan baca-saja port activity dan pemilih layar.
type Acuan interface {
	models.Acuan
	DaftarMaster(ctx context.Context, cari string) ([]models.BarisMaster, error)
	BarisMasterDari(ctx context.Context, treatyID, grupID, cobID string) (models.BarisMaster, bool, error)
	DaftarPolis(ctx context.Context, noOffer, grup string) ([]models.BarisPolis, error)
	DaftarAdjuster(ctx context.Context, cari string) ([]models.Pilihan, error)
	DaftarMataUang(ctx context.Context) ([]models.Pilihan, error)
	DaftarJenisReas(ctx context.Context, tipe string) ([]models.Pilihan, error)
	DaftarProvinsi(ctx context.Context, cari string) ([]models.Pilihan, error)
	DaftarSebab(ctx context.Context, cari string) ([]repository.BarisSebab, error)
	BarisSebabID(ctx context.Context, id string) (repository.BarisSebab, bool, error)
	DaftarKatastrofe(ctx context.Context, cari string) ([]repository.BarisKatastrofe, error)
	BarisKatastrofeID(ctx context.Context, id string) (repository.BarisKatastrofe, bool, error)
	RingkasanOS(ctx context.Context, noClaim string) ([]repository.BarisRingkasanOS, error)
	TanggaKomite(ctx context.Context, komiteID string) ([]models.AnggotaKomite, error)
	StatusKasir(ctx context.Context, noAksep string) (string, bool, error)
	StatusKonversi(ctx context.Context, noAksep string) (string, error)
	EmailCeding(ctx context.Context, ceding string) (string, error)
	IDBankRekening(ctx context.Context, bank, cabang, akun string) (string, error)
	TingkatPelaku(ctx context.Context, operatorID string) (string, error)
	NamaPelaku(ctx context.Context, akun string) (string, error)
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

// DariDasar menyusun layanan di atas basis data bersama. Tanpa Oracle layanan menjawab `ErrTanpaOracle` (503).
func DariDasar(d *inti.Dasar) *Layanan {
	if d == nil || !d.PunyaDatabase() {
		return Baru(nil, nil, time.Now, false)
	}
	g := repository.Baru(d.DB())
	produksi := d.Lingkungan().AdalahProduksi() || d.DB().PegaProduksi()
	return Baru(penyimpanOracle{Gudang: g, dasar: d}, repository.AcuanDari(g, produksi), time.Now, produksi)
}
