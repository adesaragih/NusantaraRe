package services_test

// Uji seam services tiket 01 dan 05: identitas, pengenal tak sah, dan
// perbedaan antara "kontrak tidak ada" dengan "rantai kosong".

import (
	"context"
	"errors"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyinadjustment/backend/models"
	"nusantarare/modul/treatyinadjustment/backend/services"
)

type gudangTiruan struct {
	kontrak  []models.Kontrak
	versi    map[int64][]models.Versi
	disentuh int
	galat    error

	// Panel Attachment - galatnya TERPISAH dari `galat`, sebab kontrak yang
	// rantai versinya gagal dibaca tetap harus membuka lampirannya.
	lampiran        []models.BarisLampiranWarisan
	katalogKategori map[string]string
	riwayat         []models.BarisRiwayatWarisan
	galatLampiran   error

	// Layar Adjustment - galatnya TERPISAH lagi: daftar penyesuaian dan
	// rantai versi model baru dibaca dari tabel yang berbeda.
	daftarPenyesuaian []models.BarisPenyesuaian
	penyesuaian       map[string]models.Penyesuaian
	galatPenyesuaian  error

	// Picker Add dan tombol `Choose`.
	master        []models.BarisMasterPilihan
	hanyaNonProp  *bool
	dokumenMaster map[string]models.SisiPenyesuaian
	adaRevisi     map[string]bool
}

func (g *gudangTiruan) DaftarKontrak(context.Context) ([]models.Kontrak, error) {
	g.disentuh++
	return g.kontrak, g.galat
}

func (g *gudangTiruan) RantaiVersi(_ context.Context, id int64) ([]models.Versi, error) {
	g.disentuh++
	if g.galat != nil {
		return nil, g.galat
	}
	return g.versi[id], nil
}

var pelakuAda = inti.Pelaku{AkunID: "AKUN-UJI"}

func nomor(n int64) *int64 { return &n }

func TestKeduaJalurMenolakPermintaanTanpaIdentitas(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	if _, err := l.DaftarKontrak(context.Background(), inti.Pelaku{}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("DaftarKontrak: mau ErrTanpaIdentitas, dapat %v", err)
	}
	if _, err := l.RantaiVersi(context.Background(), inti.Pelaku{}, 1); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("RantaiVersi: mau ErrTanpaIdentitas, dapat %v", err)
	}
	// Penolakan yang tetap membaca basis data bukan penolakan.
	if g.disentuh != 0 {
		t.Errorf("gudang tersentuh %d kali walau identitas tidak ada", g.disentuh)
	}
}

func TestPengenalKontrakTidakSahDitolakSebelumGudang(t *testing.T) {
	g := &gudangTiruan{}
	l := services.LayananDengan(g)

	for _, id := range []int64{0, -1, -99} {
		if _, err := l.RantaiVersi(context.Background(), pelakuAda, id); !errors.Is(err, services.ErrIDTidakSah) {
			t.Errorf("id %d: mau ErrIDTidakSah, dapat %v", id, err)
		}
	}
	if g.disentuh != 0 {
		t.Errorf("gudang tersentuh %d kali untuk pengenal tak sah", g.disentuh)
	}
}

// ⛔ Rantai kosong dan kontrak tidak ada BERBEDA artinya. Sebuah kontrak selalu
// punya sekurangnya versi pertamanya (tiket 14), jadi daftar kosong yang
// dikembalikan apa adanya akan terbaca di layar sebagai "kontrak ini tidak
// punya versi" - padahal kontraknya sendiri yang tidak ada.
func TestRantaiKosongBerartiKontrakTidakAda(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{versi: map[int64][]models.Versi{}})

	_, err := l.RantaiVersi(context.Background(), pelakuAda, 42)

	if !errors.Is(err, services.ErrKontrakTidakAda) {
		t.Fatalf("mau ErrKontrakTidakAda, dapat %v", err)
	}
}

// Uji POSITIF - dan ia yang menangkap lingkup yang terlalu sempit: kontrak
// dengan beberapa versi diterima, termasuk versi pertama yang DASARNYA KOSONG
// dan baris warisan yang NOMOR URUTNYA KOSONG (tiket 05).
func TestRantaiBerisiDiterimaTermasukYangKosongDasarDanNomornya(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{versi: map[int64][]models.Versi{
		7: {
			{ID: 1, IDKontrak: 7, NomorUrutVersi: nomor(1), IDVersiDasar: nil},
			{ID: 2, IDKontrak: 7, NomorUrutVersi: nomor(2), IDVersiDasar: nomor(1)},
			{ID: 3, IDKontrak: 7, NomorUrutVersi: nil, IDVersiDasar: nomor(2)},
		},
	}})

	versi, err := l.RantaiVersi(context.Background(), pelakuAda, 7)

	if err != nil {
		t.Fatalf("rantai sah ditolak: %v", err)
	}
	if len(versi) != 3 {
		t.Fatalf("mau 3 versi, dapat %d", len(versi))
	}
	if versi[0].IDVersiDasar != nil {
		t.Error("versi pertama harus berdasar KOSONG (bahan to-spec B-3)")
	}
	if versi[2].NomorUrutVersi != nil {
		t.Error("baris warisan harus boleh bernomor urut KOSONG (tiket 05)")
	}
}

func TestDaftarKontrakKosongBukanGalat(t *testing.T) {
	l := services.LayananDengan(&gudangTiruan{kontrak: []models.Kontrak{}})

	k, err := l.DaftarKontrak(context.Background(), pelakuAda)

	if err != nil {
		t.Fatalf("daftar kontrak kosong menghasilkan galat: %v", err)
	}
	if len(k) != 0 {
		t.Fatalf("mau nol kontrak, dapat %d", len(k))
	}
}

func TestGalatGudangDiteruskan(t *testing.T) {
	bocor := errors.New("oracle mati")
	l := services.LayananDengan(&gudangTiruan{galat: bocor})

	if _, err := l.DaftarKontrak(context.Background(), pelakuAda); !errors.Is(err, bocor) {
		t.Errorf("DaftarKontrak tidak meneruskan galat gudang; dapat %v", err)
	}
	if _, err := l.RantaiVersi(context.Background(), pelakuAda, 1); !errors.Is(err, bocor) {
		t.Errorf("RantaiVersi tidak meneruskan galat gudang; dapat %v", err)
	}
}

// Panel Attachment - `M_ATTACHMENTTREATY_2`, tabel warisan.
func (g *gudangTiruan) BacaLampiranKontrak(_ context.Context, _ string) ([]models.BarisLampiranWarisan, error) {
	return g.lampiran, g.galatLampiran
}

func (g *gudangTiruan) BacaKatalogKategoriLampiran(_ context.Context) (map[string]string, error) {
	return g.katalogKategori, g.galatLampiran
}

// Panel History - `T_VIEW_COMMENT`, tabel warisan.
// Panel `Existing Policy for Master ID` - `TREATYINPRODUCTION`, tabel warisan.
func (g *gudangTiruan) BacaPolisMaster(_ context.Context, _ string) ([]models.BarisPolisMaster, error) {
	return []models.BarisPolisMaster{}, g.galatLampiran
}

func (g *gudangTiruan) BacaRiwayatKontrak(_ context.Context, _ string) ([]models.BarisRiwayatWarisan, error) {
	return g.riwayat, g.galatLampiran
}
