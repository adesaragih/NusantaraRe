package services_test

// Uji layanan tahun treaty - TANPA Oracle (tiket 03). Gudang dan pelaksana
// transaksi dipalsukan; yang diuji GERBANGNYA: identitas, validasi, anti-dobel,
// pemilihan sisip/perbarui, dan jejak di transaksi yang sama.

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"nusantarare/internal/models"
	"nusantarare/internal/repository"
	"nusantarare/internal/services"
)

type gudangTahunUji struct {
	dobel       string
	disisip     []models.TahunTreaty
	diperbaru   []models.TahunTreaty
	jejak       []string
	perbaruiErr error
}

func (g *gudangTahunUji) Daftar(context.Context, int, int) (repository.HalamanTahunTreaty, error) {
	return repository.HalamanTahunTreaty{Baris: []models.TahunTreaty{{ID: "1000001", TreatyYear: "2026"}},
		Total: 1, Halaman: 1, Ukuran: 20}, nil
}
func (g *gudangTahunUji) Ambil(_ context.Context, id string) (models.TahunTreaty, error) {
	if id != "1000001" {
		return models.TahunTreaty{}, repository.ErrTahunTreatyTidakAda
	}
	return models.TahunTreaty{ID: id, TreatyYear: "2026"}, nil
}
func (g *gudangTahunUji) Sisip(_ context.Context, _ *repository.Tx, t models.TahunTreaty) (string, error) {
	g.disisip = append(g.disisip, t)
	return "1000009", nil
}
func (g *gudangTahunUji) Perbarui(_ context.Context, _ *repository.Tx, t models.TahunTreaty) error {
	if g.perbaruiErr != nil {
		return g.perbaruiErr
	}
	g.diperbaru = append(g.diperbaru, t)
	return nil
}
func (g *gudangTahunUji) CariDobel(context.Context, *repository.Tx, string, time.Time, time.Time, string) (string, error) {
	return g.dobel, nil
}
func (g *gudangTahunUji) Jejak(_ context.Context, _ *repository.Tx, akun, baris, aksi, ket string, _ time.Time) error {
	g.jejak = append(g.jejak, akun+"|"+baris+"|"+aksi+"|"+ket)
	return nil
}

// transaksiUji menjalankan fn tanpa Oracle (tx nil); gudang palsu mengabaikannya.
func transaksiUji(_ context.Context, fn func(tx *repository.Tx) error) error { return fn(nil) }

var jamUji = func() time.Time { return time.Date(2026, 9, 28, 10, 0, 0, 0, time.UTC) }

func layananTahun(g *gudangTahunUji) *services.TahunTreatyTCO {
	return services.New(nil).TahunTreatyTCO().DenganGudang(g).DenganTransaksi(transaksiUji).DenganJam(jamUji)
}

func masukWajar() services.TahunTreatyMasuk {
	return services.TahunTreatyMasuk{TreatyYear: "2026", UnderwritingYear: "2026", TreatyGroupID: "10001",
		TreatyGroupName: "UJI GRUP", Proportion: "P", StartDate: "2026-01-01", EndDate: "2026-12-31"}
}

func TestTahunTreatyTanpaIdentitasDitolak(t *testing.T) {
	l := layananTahun(&gudangTahunUji{})
	if _, err := l.Daftar(context.Background(), services.Pelaku{}, 1, 20); !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Errorf("daftar: %v", err)
	}
	if _, err := l.Simpan(context.Background(), services.Pelaku{}, masukWajar()); !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Errorf("simpan: %v", err)
	}
}

func TestTahunTreatyBawaanGagalTerang(t *testing.T) {
	_, err := services.New(nil).TahunTreatyTCO().DenganTransaksi(transaksiUji).
		Daftar(context.Background(), pelakuUjiTCO, 1, 20)
	if !errors.Is(err, services.ErrGudangTahunTreatyBelumDisuntik) {
		t.Errorf("bawaan harus gagal terang: %v", err)
	}
}

// AC 5, 6, 41: baris baru - ID dari sequence, UserID pelaku, TglUpdate jam,
// jejak di transaksi yang sama.
func TestTahunTreatySimpanBaru(t *testing.T) {
	g := &gudangTahunUji{}
	hasil, err := layananTahun(g).Simpan(context.Background(), pelakuUjiTCO, masukWajar())
	if err != nil {
		t.Fatal(err)
	}
	if hasil.ID != "1000009" || hasil.UserID != "UJI-ADMIN" || hasil.TglUpdate != "2026-09-28 10:00:00" {
		t.Errorf("hasil: %+v", hasil)
	}
	if hasil.StartDate != "2026-01-01" || hasil.EndDate != "2026-12-31" {
		t.Errorf("tanggal: %+v", hasil)
	}
	if len(g.disisip) != 1 || len(g.diperbaru) != 0 {
		t.Errorf("sisip %d, perbarui %d", len(g.disisip), len(g.diperbaru))
	}
	if len(g.jejak) != 1 || g.jejak[0] != "UJI-ADMIN|1000009|simpan|tahun treaty baru" {
		t.Errorf("jejak: %v", g.jejak)
	}
}

// AC 8: ID terisi = pembaruan, bukan baris baru.
func TestTahunTreatySimpanPerbarui(t *testing.T) {
	g := &gudangTahunUji{}
	m := masukWajar()
	m.ID = "1000001"
	hasil, err := layananTahun(g).Simpan(context.Background(), pelakuUjiTCO, m)
	if err != nil {
		t.Fatal(err)
	}
	if hasil.ID != "1000001" || len(g.disisip) != 0 || len(g.diperbaru) != 1 {
		t.Errorf("hasil %+v sisip %d perbarui %d", hasil, len(g.disisip), len(g.diperbaru))
	}
	if g.jejak[0] != "UJI-ADMIN|1000001|simpan|tahun treaty diperbarui" {
		t.Errorf("jejak: %v", g.jejak)
	}
	// Baris yang tidak ada: galat repository diteruskan, nol jejak.
	g2 := &gudangTahunUji{perbaruiErr: repository.ErrTahunTreatyTidakAda}
	if _, err := layananTahun(g2).Simpan(context.Background(), pelakuUjiTCO, m); !errors.Is(err, services.ErrTahunTreatyTidakAda) {
		t.Errorf("tidak ada: %v", err)
	}
	if len(g2.jejak) != 0 {
		t.Error("jejak ditulis padahal pembaruan gagal")
	}
}

// AC 73: dobel ditolak dengan pesan yang menyebut tahun treaty mana; nol tulisan.
func TestTahunTreatyAntiDobel(t *testing.T) {
	g := &gudangTahunUji{dobel: "1000005"}
	_, err := layananTahun(g).Simpan(context.Background(), pelakuUjiTCO, masukWajar())
	if !errors.Is(err, services.ErrTahunTreatyDobel) {
		t.Fatalf("mau ErrTahunTreatyDobel, dapat %v", err)
	}
	for _, mau := range []string{"1000005", "2026-01-01", "2026-12-31", "10001"} {
		if !strings.Contains(err.Error(), mau) {
			t.Errorf("pesan tanpa %q: %s", mau, err)
		}
	}
	if len(g.disisip)+len(g.diperbaru)+len(g.jejak) != 0 {
		t.Error("dobel tetap menulis")
	}
}

// AC 9 dan prasyarat b388/b411/b434 dijalankan SEBELUM transaksi.
func TestTahunTreatyValidasiSebelumTransaksi(t *testing.T) {
	kasus := []struct {
		ubah func(*services.TahunTreatyMasuk)
		mau  error
	}{
		{func(m *services.TahunTreatyMasuk) { m.EndDate = "2025-12-31" }, models.ErrPeriodeTerbalik},
		{func(m *services.TahunTreatyMasuk) { m.TreatyGroupID = "" }, models.ErrTahunTreatyGrupKosong},
		{func(m *services.TahunTreatyMasuk) { m.TreatyYear = "dua ribu" }, models.ErrTahunTreatyBukanAngka},
		{func(m *services.TahunTreatyMasuk) { m.StartDate = "01/01/2026" }, services.ErrPermintaanTidakSah},
	}
	for _, k := range kasus {
		g := &gudangTahunUji{}
		m := masukWajar()
		k.ubah(&m)
		_, err := layananTahun(g).Simpan(context.Background(), pelakuUjiTCO, m)
		if !errors.Is(err, k.mau) {
			t.Errorf("mau %v, dapat %v", k.mau, err)
		}
		if len(g.disisip)+len(g.jejak) != 0 {
			t.Errorf("%v: tetap menulis", k.mau)
		}
	}
}

func TestTahunTreatyDaftarDanAmbil(t *testing.T) {
	l := layananTahun(&gudangTahunUji{})
	hal, err := l.Daftar(context.Background(), pelakuUjiTCO, 0, 0)
	if err != nil || hal.Total != 1 || len(hal.Baris) != 1 || hal.Baris[0].ID != "1000001" {
		t.Errorf("daftar: %+v %v", hal, err)
	}
	if _, err := l.Ambil(context.Background(), pelakuUjiTCO, "9"); !errors.Is(err, services.ErrTahunTreatyTidakAda) {
		t.Errorf("ambil tidak ada: %v", err)
	}
}

type pembacaGrupUji struct{ baris []repository.GrupTreatyTCO }

func (p pembacaGrupUji) Daftar(context.Context) ([]repository.GrupTreatyTCO, error) {
	return p.baris, nil
}

func TestGrupTreatyMasterKosongAdalahKegagalan(t *testing.T) {
	l := services.New(nil).GrupTreaty().DenganPembaca(pembacaGrupUji{})
	if _, err := l.Daftar(context.Background(), pelakuUjiTCO); !errors.Is(err, services.ErrMasterGrupTreatyKosong) {
		t.Errorf("master kosong: %v", err)
	}
	l = l.DenganPembaca(pembacaGrupUji{baris: []repository.GrupTreatyTCO{{ID: "10001", TreatyGroupName: "UJI GRUP"}}})
	d, err := l.Daftar(context.Background(), pelakuUjiTCO)
	if err != nil || len(d) != 1 || d[0].TreatyGroupName != "UJI GRUP" {
		t.Errorf("daftar: %+v %v", d, err)
	}
}
