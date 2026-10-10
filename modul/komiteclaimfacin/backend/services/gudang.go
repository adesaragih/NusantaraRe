// Package services memegang aturan dagang Komite Claim Fac In di antara HTTP dan Oracle: siapa boleh memutuskan
// (pemegang tingkat berjalan, KCF-01), urutan langkah `KomitePost_Adjustment` / `KomitePost_Reject` /
// `KomitePost_CloseClaim`, dan SATU transaksi per Submit (keputusan + tangga + perluasan + nomor + tulis balik ke klaim
// induk lewat kontrak + tabel warisan + outbox). Arah ketergantungan `handlers -> services -> repository`; klaim induk
// lewat `kontrak.KlaimFacInKomite`, nol impor `modul/claimfacin`. Pola disalin dari Komite Claim Prop / Non Prop (bukan
// impor).
//
// Untuk apa berkas ini: ANTARMUKA GUDANG dan ACUAN. Implementasi Oracle-nya `repository.Gudang` / `repository.Acuan`;
// uji memasang `tiruan`.
package services

import (
	"context"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/inti/backend/kontrak"
	"nusantarare/modul/komiteclaimfacin/backend/models"
	"nusantarare/modul/komiteclaimfacin/backend/repository"
)

// Gudang - tulisan dan bacaan kasus komite serta tabel warisan.
type Gudang interface {
	Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error

	BacaKasus(ctx context.Context, tx *db.Tx, id string, kunci bool) (models.Kasus, error)
	BacaTangga(ctx context.Context, tx *db.Tx, id string) ([]models.Anggota, error)
	TambahAnggota(ctx context.Context, tx *db.Tx, id string, baru []models.Anggota) ([]models.Anggota, error)
	DaftarKerja(ctx context.Context, akun string, peran []string) ([]models.BarisKerja, error)
	TulisAnggota(ctx context.Context, tx *db.Tx, id string, u models.UbahAnggota) error
	SimpanKepala(ctx context.Context, tx *db.Tx, id string, countLama int, kp models.Kepala) error
	TutupKasus(ctx context.Context, tx *db.Tx, id string, selesai bool, posisi string, saat time.Time) error

	UrutNomorAkseptasi(ctx context.Context, tx *db.Tx, saat time.Time) (models.BahanNomor, error)
	SisipOS(ctx context.Context, tx *db.Tx, b models.BarisOS, saat time.Time) error
	SalinJSONKlaim(ctx context.Context, tx *db.Tx, idPega, noKlaim, noPolis string, saat time.Time) error
	CatatLogLayanan(ctx context.Context, tx *db.Tx, l models.LogLayanan, saat time.Time) error
	CatatRiwayatAkseptasi(ctx context.Context, tx *db.Tx, r models.RiwayatAkseptasi, saat time.Time) error
	UbahSubProgres(ctx context.Context, tx *db.Tx, komiteID, posisi string) error
	SisipKlaimDitolak(ctx context.Context, tx *db.Tx, k models.KlaimDitolak) error
	AntreEfek(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) (string, error)
}

// Acuan - bacaan baca-saja RDB komite.
type Acuan interface {
	RosterKomite(ctx context.Context) ([]models.AnggotaRoster, error)
	KursStandar(ctx context.Context, cur string) (string, error)
	NamaJenisReas(ctx context.Context) (map[string]string, error)
	IDBankRekening(ctx context.Context, bank, cabang, akun string) (string, error)
	StatusKonversi(ctx context.Context, noAksep string) (string, error)
	EmailCeding(ctx context.Context, ceding string) (string, error)
	NamaPelaku(ctx context.Context, akun string) (string, error)
	EmailPelaku(ctx context.Context, akun string) (string, error)
	EmailAnggotaWorkbasket(ctx context.Context, workbasket string) ([]string, error)
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

// DariDasar menyusun layanan di atas basis data bersama dan kontrak Claim Fac In. Tanpa Oracle layanan menjawab
// `ErrTanpaOracle` (503).
func DariDasar(d *inti.Dasar, k kontrak.KlaimFacInKomite) *Layanan {
	if d == nil || !d.PunyaDatabase() {
		return Baru(nil, nil, k, time.Now, false)
	}
	g := repository.Baru(d.DB())
	produksi := d.Lingkungan().AdalahProduksi() || d.DB().PegaProduksi()
	return Baru(penyimpanOracle{Gudang: g, dasar: d}, repository.AcuanDari(g, produksi), k, time.Now, produksi)
}
