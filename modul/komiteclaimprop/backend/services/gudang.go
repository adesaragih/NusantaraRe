// Package services memegang aturan dagang Komite Claim Prop di antara HTTP dan Oracle: siapa boleh memutuskan
// (pemegang assignment KomiteRouter, keputusan 30 / ADR-0014), urutan langkah `KomitePostAdjustment`, dan SATU
// transaksi per Submit (keputusan + tangga + nomor + tulis balik ke klaim induk lewat kontrak + tabel warisan + log +
// outbox). Arah ketergantungan `handlers -> services -> repository`; klaim induk lewat `kontrak.KlaimTreatyKomite`,
// nol impor `modul/claimprop`.
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
	"nusantarare/inti/backend/penyimpanan"
	"nusantarare/modul/komiteclaimprop/backend/models"
	"nusantarare/modul/komiteclaimprop/backend/repository"
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
	SalinJSONKlaim(ctx context.Context, tx *db.Tx, idPega, noKlaim, noPolis string, saat time.Time) error
	CatatLogLayanan(ctx context.Context, tx *db.Tx, l models.LogLayanan, saat time.Time) error
	CatatRiwayatAkseptasi(ctx context.Context, tx *db.Tx, r models.RiwayatAkseptasi, saat time.Time) error
	AntreEfek(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) (string, error)
	SisipDokumenKlaim(ctx context.Context, tx *db.Tx, d models.BarisDokumenKlaim) error
}

// PenyimpananBerkas - bagian `inti/backend/penyimpanan` yang dipakai modul ini (InsertGoogleStorage_Act +
// Insert_T_Storage_SQL; tiruan di uji).
type PenyimpananBerkas interface {
	Unggah(ctx context.Context, m penyimpanan.MasukUnggah) (penyimpanan.Objek, error)
	Catat(ctx context.Context, tx *db.Tx, o penyimpanan.Objek) error
}

// Acuan - bacaan baca-saja RDB komite.
type Acuan interface {
	TahunTreaty(ctx context.Context, grup, ymd string) (string, error)
	LimitPLA(ctx context.Context, tahun, grup, reins string) (string, error)
	DaftarRetro(ctx context.Context, tahun, grup, reins string) ([]models.BarisRetro, error)
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

// DariDasar menyusun layanan di atas basis data bersama dan kontrak Claim Prop. Tanpa Oracle layanan menjawab
// `ErrTanpaOracle` (503).
func DariDasar(d *inti.Dasar, k kontrak.KlaimTreatyKomite) *Layanan {
	if d == nil || !d.PunyaDatabase() {
		return Baru(nil, nil, k, time.Now, false)
	}
	g := repository.Baru(d.DB())
	produksi := d.Lingkungan().AdalahProduksi() || d.DB().PegaProduksi()
	return Baru(penyimpanOracle{Gudang: g, dasar: d}, repository.AcuanDari(g, produksi), k, time.Now, produksi)
}
