package services_test

// Uji seam services tiket 15: identitas, himpunan tertutup, dan daftar kosong
// sebagai JAWABAN.
//
// Wewenang dan penolakan diuji DI SINI dengan `inti.Pelaku` langsung, bukan
// lewat jalur HTTP - aturannya tidak boleh bergantung pada `X-Pelaku`.

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/models"
	"nusantarare/modul/treatyin/backend/repository"
	"nusantarare/modul/treatyin/backend/services"
)

// gudangTiruan menjawab dari peta, dan MENCATAT himpunan yang diterimanya -
// supaya uji dapat membuktikan services tidak meneruskan yang tidak sah.
type gudangTiruan struct {
	isi         map[models.Himpunan][]models.Acuan
	diminta     []models.Himpunan
	galat       error
	dibuat      []models.Kontrak
	versiDibuat []models.VersiKontrak
	dibaca      []int64
	kontrak     models.KontrakDenganVersi
	serupa      []int64
	warisan     []models.Kontrak
	diperbarui  []models.Kontrak
	versiBaru   []models.VersiKontrak
	cariWarisan []string
	adaQS       bool
	dicatat     []string

	// Tab pendaratan migrasi 430, dan galatnya yang TERPISAH dari `galat`.
	periodePelaporan []models.BarisPeriodeWarisan
	portofolio       []models.BarisPortofolioWarisan
	akumulasi        []models.BarisAkumulasiWarisan
	egnpi            []models.BarisEgnpiWarisan
	retensi          []models.BarisRetensiWarisan
	angsuran         []models.BarisAngsuranWarisan
	catatan          []models.BarisCatatanWarisan
	layer            []models.BarisLayerWarisan
	limitsPohon      []map[string]any
	revisi           repository.RevisiPendaratan
	kurs             []models.BarisKursWarisan
	skalaKoasuransi  []models.BarisSkalaKoasuransiWarisan
	lampiran         []models.BarisLampiranWarisan
	katalogKategori  map[string]string
	galatTab         error

	// Panel `Existing Policy for Master ID`, galatnya terpisah lagi: ia
	// satu-satunya pembacaan yang menyentuh `TREATYINPRODUCTION`, dan
	// kegagalannya tidak boleh terbaca sebagai kegagalan tab mana pun.
	polisProduksi      []models.BarisPolisProduksi
	galatPolisProduksi error
	cedant             []models.PilihanWarisan
	jenisTreaty        []models.PilihanWarisan
	asalBisnis         []models.PilihanWarisan
	galatPilihan       error

	// Tiket 32 - apa yang SAMPAI ke gudang, supaya uji dapat membuktikan
	// masukan yang ditolak tidak pernah diteruskan.
	pemulihan    [][]models.PemulihanLimit
	layerDicatat []int64
	// Tiket 40 - nil berarti versinya yang pertama.
	versiDasar *models.VersiKontrak
	// Tiket 42.
	arsip []models.ArsipMuatanKeluar
	bukti models.BuktiArsip
	// Layar daftar kontrak.
	daftar []models.BarisDaftarKontrak
	// Layar daftar WARISAN (`TREATY_IN`).
	barisWarisan   []models.BarisDaftarWarisan
	cacahWarisan   int
	kontrakWarisan models.KontrakWarisan

	// Tombol tulis Save/Submit/Actions/Decline offer.
	kepalaTreatyIn   map[string]map[string]any
	dokumenTersimpan map[string]map[string]any
	disimpan         []models.RencanaSimpan
	pemegangPosisi   map[string][]string
	// Properti yang kolomnya dianggap BELUM terpasang (migrasi menunggu).
	belumTerpasang []string
	// Layar Adjustment — kepala `TREATY_IN_EDM`, rencana tersimpan, penghapusan.
	kepalaEDM   map[string]map[string]any
	disimpanEDM []models.RencanaPenyesuaian
	dihapusEDM  []string
	// Panel Attachment — baris yang `CatatLampiran` terima.
	lampiranBaru []models.LampiranBaru
}

// Tiket 32.
func (g *gudangTiruan) CatatPemulihanLimit(_ context.Context, idLayer int64, baris []models.PemulihanLimit) error {
	g.layerDicatat = append(g.layerDicatat, idLayer)
	g.pemulihan = append(g.pemulihan, baris)
	return g.galat
}

// Tiket 42.
func (g *gudangTiruan) SimpanArsipMuatanKeluar(_ context.Context, a models.ArsipMuatanKeluar) error {
	g.arsip = append(g.arsip, a)
	return g.galat
}

// Tiket 42 — nilai baliknya TIDAK punya ruas muatan, dan itu seluruh maksudnya.
func (g *gudangTiruan) BuktiArsipKontrak(_ context.Context, id int64) (models.BuktiArsip, error) {
	g.dibaca = append(g.dibaca, id)
	if g.galat != nil {
		return models.BuktiArsip{}, g.galat
	}
	return g.bukti, nil
}

// Tiket 40.
func (g *gudangTiruan) BacaVersiDasar(_ context.Context, id int64) (*models.VersiKontrak, error) {
	g.dibaca = append(g.dibaca, id)
	if g.galat != nil {
		return nil, g.galat
	}
	return g.versiDasar, nil
}

func (g *gudangTiruan) DaftarAcuan(_ context.Context, h models.Himpunan) ([]models.Acuan, error) {
	g.diminta = append(g.diminta, h)
	if g.galat != nil {
		return nil, g.galat
	}
	return g.isi[h], nil
}

// Tiket 14. `dibuat` MENCATAT apa yang sampai ke gudang, supaya uji dapat
// membuktikan masukan yang ditolak tidak pernah diteruskan.
func (g *gudangTiruan) BuatKontrakDenganVersiPertama(_ context.Context, k models.Kontrak, v models.VersiKontrak) (int64, int64, error) {
	g.dibuat = append(g.dibuat, k)
	g.versiDibuat = append(g.versiDibuat, v)
	if g.galat != nil {
		return 0, 0, g.galat
	}
	return 777, 888, nil
}

func (g *gudangTiruan) BacaKontrak(_ context.Context, id int64) (models.KontrakDenganVersi, error) {
	g.dibaca = append(g.dibaca, id)
	if g.galat != nil {
		return models.KontrakDenganVersi{}, g.galat
	}
	return g.kontrak, nil
}

var pelakuAda = inti.Pelaku{AkunID: "AKUN-UJI"}

func TestDaftarAcuanMenolakPermintaanTanpaIdentitas(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	_, err := l.DaftarAcuan(context.Background(), inti.Pelaku{}, models.HimpunanJenisPotongan)

	if !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Fatalf("mau ErrTanpaIdentitas, dapat %v", err)
	}
	// ⛔ Yang penting BUKAN hanya galatnya: gudang tidak boleh tersentuh sama
	// sekali. Penolakan yang tetap membaca basis data bukan penolakan.
	if len(g.diminta) != 0 {
		t.Errorf("gudang tersentuh %d kali walau identitas tidak ada", len(g.diminta))
	}
}

func TestDaftarAcuanMenolakHimpunanDiLuarEnam(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	for _, h := range []models.Himpunan{"", "KONTRAK", "mata_uang", "VERSI_KONTRAK", "../bahaya"} {
		_, err := l.DaftarAcuan(context.Background(), pelakuAda, h)
		if !errors.Is(err, services.ErrHimpunanTidakAda) {
			t.Errorf("himpunan %q: mau ErrHimpunanTidakAda, dapat %v", h, err)
		}
	}
	// Nama tabel tidak pernah datang dari teks permintaan: yang tidak sah
	// berhenti SEBELUM gudang.
	if len(g.diminta) != 0 {
		t.Errorf("gudang menerima %v; himpunan tak sah tidak boleh diteruskan", g.diminta)
	}
}

// Uji POSITIF - dan ia yang menangkap penyaring yang terlalu ketat. Penyaring
// yang menolak salah satu dari enam lulus setiap uji negatif di atas.
func TestKelimaHimpunanDiterima(t *testing.T) {
	enam := []models.Himpunan{
		models.HimpunanJenisPotongan, models.HimpunanJenisPotongan, models.HimpunanKelasBisnis,
		models.HimpunanKelompokTreaty, models.HimpunanBahaya, models.HimpunanJenisReasuransi,
	}
	g := &gudangTiruan{isi: map[models.Himpunan][]models.Acuan{}}
	for _, h := range enam {
		g.isi[h] = []models.Acuan{{ID: 1, Kode: "K1", Nama: "N1", Aktif: "1"}}
	}
	l := services.LayananDengan(g)

	for _, h := range enam {
		baris, err := l.DaftarAcuan(context.Background(), pelakuAda, h)
		if err != nil {
			t.Errorf("himpunan %q ditolak: %v", h, err)
			continue
		}
		if len(baris) != 1 {
			t.Errorf("himpunan %q: mau 1 baris, dapat %d", h, len(baris))
		}
	}
	if len(g.diminta) != len(enam) {
		t.Errorf("gudang dipanggil %d kali, mau %d", len(g.diminta), len(enam))
	}
}

// Daftar kosong adalah JAWABAN, bukan galat: kelima tabel berdiri kosong
// sampai tiket 44 memindahkan isinya dari sistem lama.
func TestTabelAcuanKosongBukanGalat(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{isi: map[models.Himpunan][]models.Acuan{}})

	baris, err := l.DaftarAcuan(context.Background(), pelakuAda, models.HimpunanBahaya)

	if err != nil {
		t.Fatalf("tabel acuan kosong menghasilkan galat: %v", err)
	}
	if len(baris) != 0 {
		t.Fatalf("mau nol baris, dapat %d", len(baris))
	}
}

func TestGalatGudangDiteruskan(t *testing.T) {
	bocor := errors.New("oracle mati")
	l := services.LayananDengan(&gudangTiruan{galat: bocor})

	_, err := l.DaftarAcuan(context.Background(), pelakuAda, models.HimpunanJenisPotongan)

	if !errors.Is(err, bocor) {
		t.Fatalf("galat gudang tidak diteruskan; dapat %v", err)
	}
}

func TestBersusunHanyaJenisReasuransi(t *testing.T) {
	if !models.HimpunanJenisReasuransi.Bersusun() {
		t.Error("JENIS_REASURANSI bersusun (§10.6) tetapi Bersusun() false")
	}
	for _, h := range []models.Himpunan{
		models.HimpunanJenisPotongan, models.HimpunanJenisPotongan, models.HimpunanKelasBisnis,
		models.HimpunanKelompokTreaty, models.HimpunanBahaya,
	} {
		if h.Bersusun() {
			t.Errorf("%q tidak bersusun tetapi Bersusun() true", h)
		}
	}
}

// Tiket 16, 17, 18, 19.
func (g *gudangTiruan) CariKontrakSerupa(context.Context, models.Kontrak) ([]int64, error) {
	if g.galat != nil {
		return nil, g.galat
	}
	return g.serupa, nil
}

func (g *gudangTiruan) CariKontrakLewatNomorWarisan(_ context.Context, nomor string) ([]models.Kontrak, error) {
	g.cariWarisan = append(g.cariWarisan, nomor)
	if g.galat != nil {
		return nil, g.galat
	}
	return g.warisan, nil
}

func (g *gudangTiruan) PerbaruiKontrak(_ context.Context, k models.Kontrak) error {
	g.diperbarui = append(g.diperbarui, k)
	return g.galat
}

func (g *gudangTiruan) TambahVersi(_ context.Context, _ int64, v models.VersiKontrak) (int64, error) {
	g.versiBaru = append(g.versiBaru, v)
	if g.galat != nil {
		return 0, g.galat
	}
	return 999, nil
}

// Tiket 35 dan 36.
func (g *gudangTiruan) AdaQuotaSharePadaVersi(context.Context, int64) (bool, error) {
	if g.galat != nil {
		return false, g.galat
	}
	return g.adaQS, nil
}

func (g *gudangTiruan) CatatKetentuanProporsional(_ context.Context, _, _, _ int64, jenis, _ string, _ *int64) error {
	if g.galat != nil {
		return g.galat
	}
	g.dicatat = append(g.dicatat, jenis)
	return nil
}

// Layar daftar kontrak — ronde layar 1.
func (g *gudangTiruan) DaftarKontrak(context.Context) ([]models.BarisDaftarKontrak, error) {
	if g.galat != nil {
		return nil, g.galat
	}
	return g.daftar, nil
}

// Layar daftar WARISAN — jalur baca `TREATY_IN`.
func (g *gudangTiruan) CacahKontrakWarisan(context.Context) (int, error) {
	if g.galat != nil {
		return 0, g.galat
	}
	return g.cacahWarisan, nil
}

func (g *gudangTiruan) DaftarKontrakWarisan(_ context.Context, offset, batas int) ([]models.BarisDaftarWarisan, error) {
	if g.galat != nil {
		return nil, g.galat
	}
	// Memotong seperti basis data memotong, supaya uji penomoran halaman
	// menguji penomorannya - bukan tiruan yang selalu mengembalikan semua.
	if offset >= len(g.barisWarisan) {
		return []models.BarisDaftarWarisan{}, nil
	}
	akhir := offset + batas
	if akhir > len(g.barisWarisan) {
		akhir = len(g.barisWarisan)
	}
	return append([]models.BarisDaftarWarisan{}, g.barisWarisan[offset:akhir]...), nil
}

// Satu kontrak WARISAN.
func (g *gudangTiruan) BacaKontrakWarisan(_ context.Context, id string) (models.KontrakWarisan, error) {
	if g.galat != nil {
		return models.KontrakWarisan{}, g.galat
	}
	return g.kontrakWarisan, nil
}

// Tab yang PINDAH ke tabel pendaratan migrasi 430 — tiruannya mengembalikan
// apa yang ditaruh di medan di bawah, dan `galatTab` membuat ketiganya
// gagal tanpa membuat `BacaKontrakWarisan` ikut gagal.
//
// ⛔ Galatnya TERPISAH dari `galat`, dan itu pokoknya: kontrak yang
// dokumennya hilang tetap harus membuka tabnya, dan tab yang tabelnya
// gagal dibaca tidak boleh terbaca sebagai "kontraknya tidak ada".
func (g *gudangTiruan) BacaPeriodePelaporan(_ context.Context, _ string) ([]models.BarisPeriodeWarisan, error) {
	return g.periodePelaporan, g.galatTab
}

func (g *gudangTiruan) BacaPortofolio(_ context.Context, _ string) ([]models.BarisPortofolioWarisan, error) {
	return g.portofolio, g.galatTab
}

func (g *gudangTiruan) BacaAkumulasi(_ context.Context, _ string) ([]models.BarisAkumulasiWarisan, error) {
	return g.akumulasi, g.galatTab
}

// Empat tab pendaratan berikutnya. `galatTab` dipakai ulang supaya uji
// "satu tab gagal, seluruh pembacaan gagal" berlaku untuk kedelapan.
func (g *gudangTiruan) BacaEgnpi(_ context.Context, _ string) ([]models.BarisEgnpiWarisan, error) {
	return g.egnpi, g.galatTab
}

// ⭐ Panel `Existing Policy for Master ID`. Nol baris adalah keadaan yang
// SAH — gambar 01 dokumen desain memperlihatkannya berbunyi "No items".
func (g *gudangTiruan) BacaPolisProduksi(_ context.Context, _ string) ([]models.BarisPolisProduksi, error) {
	return g.polisProduksi, g.galatPolisProduksi
}

func (g *gudangTiruan) BacaRetensi(_ context.Context, _ string) ([]models.BarisRetensiWarisan, error) {
	return g.retensi, g.galatTab
}
func (g *gudangTiruan) BacaAngsuran(_ context.Context, _ string) ([]models.BarisAngsuranWarisan, error) {
	return g.angsuran, g.galatTab
}
func (g *gudangTiruan) BacaCatatan(_ context.Context, _ string) ([]models.BarisCatatanWarisan, error) {
	return g.catatan, g.galatTab
}

// Isi kedua pemilih "Choose …". Galatnya BERBAGI `galatPilihan`, bukan
// `galatTab`: keduanya katalog yang tidak bergantung kontrak, dan uji yang
// mematikan satu tab tidak boleh ikut mematikan pemilihnya.
func (g *gudangTiruan) BacaDaftarCedant(_ context.Context) ([]models.PilihanWarisan, error) {
	return g.cedant, g.galatPilihan
}
func (g *gudangTiruan) BacaDaftarAsalBisnis(_ context.Context) ([]models.PilihanWarisan, error) {
	return g.asalBisnis, g.galatPilihan
}
func (g *gudangTiruan) BacaDaftarJenisTreaty(_ context.Context) ([]models.PilihanWarisan, error) {
	return g.jenisTreaty, g.galatPilihan
}
func (g *gudangTiruan) BacaDaftarKelompokTreaty(_ context.Context) ([]models.PilihanWarisan, error) {
	return g.jenisTreaty, g.galatPilihan
}
func (g *gudangTiruan) BacaDaftarMataUangLimit(_ context.Context) ([]models.PilihanWarisan, error) {
	return g.jenisTreaty, g.galatPilihan
}
func (g *gudangTiruan) BacaDaftarKelasBisnisTreaty(_ context.Context, _ string) ([]models.PilihanWarisan, error) {
	return g.jenisTreaty, g.galatPilihan
}

func (_ *gudangTiruan) BacaSharePendaratan(_ context.Context, _ string) (models.SharePendaratan, error) {
	return models.SharePendaratan{}, nil
}

func (_ *gudangTiruan) BacaIndukSpreading(_ context.Context, _, _, _ string) ([]models.SusunanSpreading, error) {
	return []models.SusunanSpreading{}, nil
}

func (_ *gudangTiruan) BacaAnakSpreading(_ context.Context, _, _ string) ([]models.SusunanSpreading, error) {
	return []models.SusunanSpreading{}, nil
}

func (_ *gudangTiruan) BacaAnakSpreadingProp(_ context.Context, _, _, _, _, _ string) ([]models.SusunanSpreading, error) {
	return []models.SusunanSpreading{}, nil
}

func (_ *gudangTiruan) BacaDaftarReasuradurShare(_ context.Context) ([]models.PilihanWarisan, error) {
	return []models.PilihanWarisan{}, nil
}

func (_ *gudangTiruan) BacaShareAkarRevisi(_ context.Context, _ string) (map[string]string, error) {
	return map[string]string{}, nil
}

func (_ *gudangTiruan) BacaShareDetailWarisan(_ context.Context, _ string) (models.ShareDetailWarisan, error) {
	return models.ShareDetailWarisan{}, nil
}

func (g *gudangTiruan) BacaSkalaKoasuransi(_ context.Context, _ string) ([]models.BarisSkalaKoasuransiWarisan, error) {
	return g.skalaKoasuransi, g.galatTab
}

// ⭐ Keempat tab dari TABEL PENDARATAN — keputusan 6 Oktober 2026.
func (g *gudangTiruan) BacaLayerPendaratan(_ context.Context, _ string) ([]models.BarisLayerWarisan, error) {
	return g.layer, g.galatTab
}

func (g *gudangTiruan) BacaAchievement(_ context.Context, _ string) ([]models.BarisAchievement, error) {
	return []models.BarisAchievement{}, nil
}
func (g *gudangTiruan) BacaKursKeIDR(_ context.Context, _ []string) (map[string]string, error) {
	return map[string]string{}, nil
}

func (g *gudangTiruan) BacaLimitsAkarPendaratan(_ context.Context, _ string) (models.LimitsAkar, error) {
	return models.LimitsAkar{}, nil
}

func (g *gudangTiruan) BacaTotalPenampung(_ context.Context, _ string) (map[string][]map[string]any, error) {
	return map[string][]map[string]any{}, nil
}

func (g *gudangTiruan) BacaPohonLimitsPendaratan(_ context.Context, _ string) ([]map[string]any, error) {
	return g.limitsPohon, g.galatTab
}

// ⭐ Medan kepala dan grid Rate of Exchange — juga dari pendaratan.
func (g *gudangTiruan) BacaRevisiPendaratan(_ context.Context, _ string) (repository.RevisiPendaratan, error) {
	return g.revisi, g.galatTab
}

func (g *gudangTiruan) BacaKursTahunan(_ context.Context, _ string) ([]models.BarisKursWarisan, error) {
	return g.kurs, g.galatTab
}

// Panel Attachment - `M_ATTACHMENTTREATY_2`, tabel warisan.
func (g *gudangTiruan) BacaLampiranKontrak(_ context.Context, _ string) ([]models.BarisLampiranWarisan, error) {
	return g.lampiran, g.galatTab
}

func (g *gudangTiruan) BacaKatalogKategoriLampiran(_ context.Context) (map[string]string, error) {
	return g.katalogKategori, g.galatTab
}
