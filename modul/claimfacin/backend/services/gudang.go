// Package services memegang aturan dagang Claim Fac In di antara HTTP dan Oracle: siapa boleh apa (pemegang
// assignment), urutan langkah flow action Pega (pra-proses -> layar -> aksi -> pasca-proses), dan SATU transaksi per
// aksi (prompt §6 butir 3: nomor + tabel datar + OS_AKSEPTASI_KLAIM / JSON_KLAIM / T_VIEW_SUGGEST + outbox). Arah
// ketergantungan `handlers -> services -> repository`. Pola `modul/claimnonprop/backend/services` (disalin, bukan
// diimpor).
//
// Untuk apa berkas ini: ANTARMUKA GUDANG dan ACUAN. Layanan hanya mengenal antarmuka ini; implementasi Oracle-nya
// `repository.Gudang` / `repository.Acuan`. Uji seam HTTP memasang gudang tiruan (`tiruan`).
package services

import (
	"context"
	"time"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/claimfacin/backend/models"
	"nusantarare/modul/claimfacin/backend/repository"
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
	TutupKasusStatus(ctx context.Context, tx *db.Tx, id, tahap, status string, saat time.Time) error
	SentuhKasus(ctx context.Context, tx *db.Tx, id string, saat time.Time) error
	DaftarKasus(ctx context.Context, s repository.SaringanKasus) ([]repository.RingkasanKasus, error)

	SimpanHalaman(ctx context.Context, tx *db.Tx, id string, h *models.Halaman) error
	BacaHalaman(ctx context.Context, tx *db.Tx, id string) (*models.Halaman, error)

	UrutNomor(ctx context.Context, tx *db.Tx, huruf string, saat time.Time) (repository.BahanNomor, error)

	SisipOS(ctx context.Context, tx *db.Tx, b models.BarisOS, saat time.Time) error
	BacaOS(ctx context.Context, kasusID, noKlaim string) ([]models.BarisOSTersimpan, error)
	SalinJSONKlaim(ctx context.Context, tx *db.Tx, idPega, noKlaim, noPolis string, saat time.Time) error
	SisipKatastrofe(ctx context.Context, tx *db.Tx, k models.BarisKatastrofe, saat time.Time) error
	CatatLogLayanan(ctx context.Context, tx *db.Tx, l repository.LogLayanan, saat time.Time) error
	TulisProgres(ctx context.Context, tx *db.Tx, p models.Progres, s models.SubProgres) error
	AntreEfek(ctx context.Context, tx *db.Tx, jenis, rujukan, muatan string, saat time.Time) (string, error)
	TandaiDLAOS(ctx context.Context, tx *db.Tx, noAkseptasi, noDLA string, saat time.Time) error

	KategoriLampiran(ctx context.Context, id string) ([]models.KategoriLampiran, error)
	DaftarLampiran(ctx context.Context, id string) ([]models.Lampiran, error)
	SisipDokumenKlaim(ctx context.Context, tx *db.Tx, d models.BarisDokumenKlaim) error
	HapusDokumenKlaim(ctx context.Context, tx *db.Tx, id, lid string) error
	PindahKategoriDokumen(ctx context.Context, tx *db.Tx, id, lid, kategori string) error

	BuatKasusKomite(ctx context.Context, tx *db.Tx, klaimID, adjID, transfer string, teks repository.TeksKomite,
		pembuat, namaPembuat string, anggota []repository.AnggotaTangga, saat time.Time) (string, error)
	// TutupKomiteAnak = `CloseAllSubCases=true` (pxForceCaseClose / ASMForceCaseClose): kasus komite KMT- klaim yang
	// masih terbuka ditutup berstatus `status`, kecuali `kecuali` (kosong = semua).
	TutupKomiteAnak(ctx context.Context, tx *db.Tx, klaimID, kecuali, status string, saat time.Time) error
	AdaKomiteTutupTerbuka(ctx context.Context, tx *db.Tx, klaimID string) (bool, error)
	SetelKomiteAdjustment(ctx context.Context, tx *db.Tx, adjID, komiteID, komiteLama string) error
	UbahAdjustmentKomite(ctx context.Context, tx *db.Tx, klaimID, adjID string, nilai map[string]string) error
}

// Acuan - bacaan baca-saja port activity dan pemilih layar.
type Acuan interface {
	models.Acuan
	AdaPolis(ctx context.Context, nopolis, prodke string) (bool, error)
	DaftarMataUang(ctx context.Context) ([]models.Pilihan, error)
	DaftarJenisReas(ctx context.Context) ([]models.Pilihan, error)
	DaftarAdjuster(ctx context.Context, cari string) ([]models.Pilihan, error)
	DaftarNegara(ctx context.Context, cari string) ([]models.Pilihan, error)
	DaftarProvinsi(ctx context.Context, negara, cari string) ([]models.Pilihan, error)
	DaftarKota(ctx context.Context, provinsi, cari string) ([]models.Pilihan, error)
	DaftarDistrik(ctx context.Context, kotaID, cari string) ([]models.Pilihan, error)
	DaftarRW(ctx context.Context, distrikID, cari string) ([]models.Pilihan, error)
	DaftarSebab(ctx context.Context, cari string) ([]repository.BarisSebab, error)
	BarisSebabID(ctx context.Context, id string) (repository.BarisSebab, bool, error)
	DaftarKatastrofe(ctx context.Context, cari string) ([]repository.BarisKatastrofe, error)
	BarisKatastrofeID(ctx context.Context, id string) (repository.BarisKatastrofe, bool, error)
	ProgresKlaim(ctx context.Context, kasusID string) ([]models.Baris, map[int][]models.Baris, error)
	IDBankRekening(ctx context.Context, bank, cabang, akun string) (string, error)
	StatusKasir(ctx context.Context, noAksep string) (string, bool, error)
	StatusKonversi(ctx context.Context, noAksep string) (string, error)
	EmailCeding(ctx context.Context, ceding string) (string, error)
	LimitDirekturUtama(ctx context.Context) (string, bool, error)
	SaldoPremi(ctx context.Context, invoice, cur string) (string, error)
	AdaProteksiPremi(ctx context.Context, nopolis string) (bool, error)
	TanggaKomite(ctx context.Context, komiteID string) ([]models.AnggotaKomite, error)
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
