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
type Service struct {
	*inti.Dasar
	// garamToken - `STORAGE_TOKEN_SALT`, bahan token penyimpanan BARU (token
	// berlaku di `GCP_IMAGE` dipakai ulang tanpa garam). ⛔ Tidak pernah
	// dicetak, dicatat, atau masuk pesan galat.
	garamToken string
}

// DenganGaramToken memasang garam token penyimpanan lampiran — dipanggil
// sekali dari `modul.go` dengan `STORAGE_TOKEN_SALT`; boleh kosong.
func (s *Service) DenganGaramToken(garam string) *Service {
	salin := *s
	salin.garamToken = garam
	return &salin
}

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
	// Larik AKAR tab Limits Non-Prop (`Summary of Limit`, `Total All Layers`).
	BacaLimitsAkarPendaratan(ctx context.Context, masterID string) (models.LimitsAkar, error)
	// Sub-tab Achievement — RDB `GetAchievement` dan `GetCurrencyToIDR_SQL`.
	BacaAchievement(ctx context.Context, idKontrak string) ([]models.BarisAchievement, error)
	BacaKursKeIDR(ctx context.Context, idMataUang []string) (map[string]string, error)

	// ⭐ MEDAN KEPALA dari tabel pendaratan, 6 Oktober 2026. Sesudahnya NOL
	// nilai layar Treaty In datang dari `M_TREATY_IN.JSONDATA`.
	//
	// ⛔ Grid Rate of Exchange TIDAK ikut mendarat: sumbernya
	// `TREATYEXCHANGEYEARLY`, tabel warisan yang sudah hidup — keputusan
	// pemilik proses 4 dan 6 Oktober 2026. Ia berkunci TAHUN, bukan kontrak,
	// jadi seamnya menerima `TREATYYEAR` dan bukan `MASTERID`.
	BacaRevisiPendaratan(ctx context.Context, masterID string) (repository.RevisiPendaratan, error)
	// Larik total yang tab pegang di penampung halaman (`T_TREATY_TOTAL`,
	// `repository.LarikTotalPenampung`) — supaya isian yang di-Save tampil
	// kembali tanpa Refresh.
	BacaTotalPenampung(ctx context.Context, masterID string) (map[string][]map[string]any, error)
	// Medan akar `T_TREATY_HAZARD_LIMIT` — isi tab Event Limits (Non-Prop)
	// dan kedua batas Co-Ins. Tabelnya berdiri sejak migrasi `446` dan Save
	// menulisnya; pembacanya baru ada 8 Oktober 2026.
	BacaBatasBahaya(ctx context.Context, masterID string) (map[string]string, error)
	// Tombol `Submit` sub-tab Achievement — `InsertToLogAchievement` ke
	// `LOG_ACHIEVEMENT`, satu transaksi (keputusan pemakai 8 Oktober 2026).
	CatatLogAchievement(ctx context.Context, masterID, operator string, baris []repository.BarisLogSiap) error
	BacaKursTahunan(ctx context.Context, tahunTreaty string) ([]models.BarisKursWarisan, error)
	// BacaKursKontrak — kurs MILIK kontrak (`T_TREATY_KURS`, migrasi 455);
	// tanpa catatan → `BacaKursTahunan`.
	BacaKursKontrak(ctx context.Context, masterID, tahunTreaty string) ([]models.BarisKursWarisan, error)

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
	// `TreatyInCheckCedingBlacklist` — `AGENT.STATUSACTIVE` per pengenal agen
	// (`agen_daftar_negatif.go`). Baca saja.
	BacaStatusAktifAgen(ctx context.Context, ids []string) (map[string]string, error)
	// Isi dropdown `Treaty Type` tab Limits — `BrowseReinsuranceType_RD`.
	BacaDaftarJenisTreaty(ctx context.Context) ([]models.PilihanWarisan, error)
	// `Treaty Group` dan mata uang tab Limits — `BrowseTreatyGroup_RD`,
	// `BrowseCurrencyTreatyIn_RD`.
	BacaDaftarKelompokTreaty(ctx context.Context) ([]models.PilihanWarisan, error)
	BacaDaftarMataUangLimit(ctx context.Context) ([]models.PilihanWarisan, error)
	// Autocomplete `Class of Business` per Treaty Group —
	// `BrowseTreatyBusinessWOType_RD`.
	BacaDaftarKelasBisnisTreaty(ctx context.Context, treatyGroupID string) ([]models.PilihanWarisan, error)

	// ⭐ Tab Share Non-Prop — pohon `T_TREATY_SHARE*`, `T_TREATY_RETRO_SHARE`,
	// `T_TREATY_FAC_*` (pendaratan), dan kedua RD spreading atas
	// `PROPORTIONALARRG` + `TREATYYEAR` (master, BACA SAJA).
	BacaSharePendaratan(ctx context.Context, masterID string) (models.SharePendaratan, error)
	BacaIndukSpreading(ctx context.Context, treatyGroupID, treatyDescID, tanggalMulai, kecuali string) ([]models.SusunanSpreading, error)
	BacaAnakSpreading(ctx context.Context, parentReinsTypeID, treatyYearID string) ([]models.SusunanSpreading, error)
	// RD `Limit_MstTrt_RD` cabang PROP — kelima filternya, parameter kosong dilewati.
	BacaAnakSpreadingProp(ctx context.Context, treatyYear, treatyGroupID, treatyDescID, parent, treatyYearID string) ([]models.SusunanSpreading, error)
	// Autocomplete `Reinsurer Name` / `Facultative Reinsurers` —
	// `BrowseAgentNusaRe_RD` TANPA saringan `ChildCount`.
	BacaDaftarReasuradurShare(ctx context.Context) ([]models.PilihanWarisan, error)
	// Skalar akar Share: kolom `T_TREATY_REVISION` (migrasi 445; toleran bila
	// belum dijalankan) dan cadangan `TREATYINDETAIL` (warisan, baca saja).
	BacaShareAkarRevisi(ctx context.Context, masterID string) (map[string]string, error)
	BacaShareDetailWarisan(ctx context.Context, masterID string) (models.ShareDetailWarisan, error)

	// ⭐ Tombol tulis Save/Submit/Actions/Decline offer — keputusan pemilik
	// proses 6–7 Oktober 2026. Satu-satunya jalur TULIS form Treaty In:
	// `T_TREATY_*` (peta pendaratan), kepala `TREATY_IN`, kurs
	// `TREATYEXCHANGEYEARLY` (tambah & ubah), dalam satu transaksi.
	BacaKepalaTreatyIn(ctx context.Context, id string) (map[string]any, bool, error)
	BacaDokumenPendaratan(ctx context.Context, masterID string) (map[string]any, error)
	SimpanKontrak(ctx context.Context, r models.RencanaSimpan) (string, error)
	// Pemegang workbasket (Kelola User) — `PositionUsername` jalur naik.
	PemegangPosisi(ctx context.Context, workbasket string) ([]string, error)
	// Divisi akun (`M_LOGIN_GO.DIVISION_CODE`) — `IT` memegang Force Edit.
	DivisiAkun(ctx context.Context, akunID string) (string, error)
	// Properti dokumen yang kolom/tabelnya BELUM terpasang (migrasi menunggu)
	// — dilaporkan Save, tidak ditelan.
	KunciBelumTerpasang(ctx context.Context, doc map[string]any) ([]string, error)

	// ⭐ Tombol tulis layar ADJUSTMENT (EDM) — 7 Oktober 2026: kepala
	// `TREATY_IN_EDM` dan kedua sisi dokumennya di `T_TREATY_*`.
	BacaKepalaPenyesuaian(ctx context.Context, id string) (map[string]any, bool, error)
	SimpanPenyesuaian(ctx context.Context, r models.RencanaPenyesuaian) error
	HapusPenyesuaian(ctx context.Context, id string) error

	// ⭐ Unggah panel Attachment — 8 Oktober 2026: objek `T_STORAGE_IMAGE`
	// dan baris `M_ATTACHMENTTREATY_2`, satu transaksi.
	NamaAplikasiSimpanan(ctx context.Context) (string, error)
	CatatLampiran(ctx context.Context, l models.LampiranBaru) (string, error)
	// Modal View File — unduh, hapus, ganti kategori.
	BacaObjekSimpanan(ctx context.Context, imageID string) (models.ObjekSimpanan, bool, error)
	PerbaruiObjekSimpanan(ctx context.Context, o models.ObjekSimpanan, tanggal string) error
	HapusLampiran(ctx context.Context, idKontrak, idLampiran, imageID string) error
	UbahKategoriLampiran(ctx context.Context, idKontrak string, ubah []models.PerubahanKategori) error
}

// Layanan memegang aturan modul ini di atas satu Gudang.
type Layanan struct {
	gudang Gudang
	// simpanan - `ServiceGoogle` panel Attachment; nil = belum disambung.
	simpanan PengirimSimpanan
}

// LayananOracle merakit Layanan di atas Oracle.
func LayananOracle(s *Service) *Layanan {
	return &Layanan{gudang: repository.Baru(s.DB()), simpanan: s.pengirimSimpanan()}
}

// LayananDengan merakit Layanan di atas Gudang mana pun - dipakai uji.
func LayananDengan(g Gudang) *Layanan { return &Layanan{gudang: g} }

// LayananDenganSimpanan - seperti `LayananDengan`, dengan pengirim
// penyimpanan tiruan. Dipakai uji.
func LayananDenganSimpanan(g Gudang, s PengirimSimpanan) *Layanan {
	return &Layanan{gudang: g, simpanan: s}
}

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
