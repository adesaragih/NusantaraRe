// Package services memegang aturan dagang Komite Claim Non Prop di antara HTTP dan Oracle: siapa boleh memutuskan
// (pemegang assignment KomiteRouter), urutan langkah `KomitePostAdjustment`, dan SATU transaksi per Submit (keputusan +
// tangga + nomor + tulis balik ke klaim induk lewat kontrak + tabel warisan + outbox). Arah ketergantungan
// `handlers -> services -> repository`; klaim induk lewat `kontrak.KlaimTreatyNonPropKomite`, nol impor
// `modul/claimnonprop`. Pola disalin dari Komite Claim Prop (bukan impor).
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
	"nusantarare/modul/komiteclaimnonprop/backend/models"
	"nusantarare/modul/komiteclaimnonprop/backend/repository"
)

// Gudang - tulisan dan bacaan kasus komite serta tabel warisan.
type Gudang interface {
	Transaksi(ctx context.Context, fn func(tx *db.Tx) error) error

	BacaKasus(ctx context.Context, tx *db.Tx, id string, kunci bool) (models.Kasus, error)
	BacaTangga(ctx context.Context, tx *db.Tx, id string) ([]models.Anggota, error)
	DaftarKerja(ctx context.Context, akun string, peran []string) ([]models.BarisKerja, error)
	TulisAnggota(ctx context.Context, tx *db.Tx, id string, u models.UbahAnggota) error
	SimpanKepala(ctx context.Context, tx *db.Tx, id string, countLama int, kp models.Kepala) error
	KomentarAwal(ctx context.Context, klaimID, adjID, id string) (string, error)
	TutupKasus(ctx context.Context, tx *db.Tx, id string, selesai bool, posisi string, saat time.Time) error

	UrutNomorAkseptasi(ctx context.Context, tx *db.Tx, saat time.Time) (models.BahanNomor, error)
	SisipOS(ctx context.Context, tx *db.Tx, b models.BarisOSAkseptasi, saat time.Time) error
	SisipXOL2(ctx context.Context, tx *db.Tx, b models.BarisXOL2, saat time.Time) error
	SimpanOSSubjectivity(ctx context.Context, tx *db.Tx, b models.BarisOSSubjectivity, saat time.Time) error
	SalinJSONKlaim(ctx context.Context, tx *db.Tx, idPega, noKlaim, noPolis string, saat time.Time) error
	CatatRiwayatAkseptasi(ctx context.Context, tx *db.Tx, r models.RiwayatAkseptasi, saat time.Time) error
	AntreEfek(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) (string, error)
}

// Acuan - bacaan baca-saja RDB komite.
type Acuan interface {
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

// DariDasar menyusun layanan di atas basis data bersama dan kontrak Claim Non Prop. Tanpa Oracle layanan menjawab
// `ErrTanpaOracle` (503).
func DariDasar(d *inti.Dasar, k kontrak.KlaimTreatyNonPropKomite) *Layanan {
	if d == nil || !d.PunyaDatabase() {
		return Baru(nil, nil, k, time.Now, false)
	}
	g := repository.Baru(d.DB())
	produksi := d.Lingkungan().AdalahProduksi() || d.DB().PegaProduksi()
	return Baru(penyimpanOracle{Gudang: g, dasar: d}, repository.AcuanDari(g, produksi), k, time.Now, produksi)
}
