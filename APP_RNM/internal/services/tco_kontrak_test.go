package services_test

// Uji layanan kontrak treaty - TANPA Oracle (tiket 04).

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

type gudangKontrakUji struct {
	baris     map[string]models.KontrakTreaty
	urut      int
	dobel     string
	disisip   int
	diperbaru int
	jejak     []string
	anak      int64
}

func (g *gudangKontrakUji) Daftar(_ context.Context, tahunID string) ([]models.KontrakTreaty, error) {
	var hasil []models.KontrakTreaty
	for _, k := range g.baris {
		if k.IDTreatyYear == tahunID {
			hasil = append(hasil, k)
		}
	}
	return hasil, nil
}
func (g *gudangKontrakUji) Ambil(_ context.Context, tahunID, id string) (models.KontrakTreaty, error) {
	k, ada := g.baris[id]
	if !ada || k.IDTreatyYear != tahunID {
		return models.KontrakTreaty{}, repository.ErrKontrakTidakAda
	}
	return k, nil
}
func (g *gudangKontrakUji) Sisip(_ context.Context, _ *repository.Tx, k models.KontrakTreaty) (string, error) {
	g.urut++
	g.disisip++
	k.ID = "100000" + string(rune('0'+g.urut))
	g.baris[k.ID] = k
	return k.ID, nil
}
func (g *gudangKontrakUji) Perbarui(_ context.Context, _ *repository.Tx, k models.KontrakTreaty) error {
	lama, ada := g.baris[k.ID]
	if !ada || lama.IDTreatyYear != k.IDTreatyYear {
		return repository.ErrKontrakTidakAda
	}
	g.diperbaru++
	g.baris[k.ID] = k
	return nil
}
func (g *gudangKontrakUji) CariDobel(_ context.Context, _ *repository.Tx, _, _, _ string) (string, error) {
	return g.dobel, nil
}
func (g *gudangKontrakUji) JumlahAnakKombinasi(context.Context, *repository.Tx, models.KombinasiTCO, string) (int64, error) {
	return g.anak, nil
}
func (g *gudangKontrakUji) Jejak(_ context.Context, _ *repository.Tx, akun, baris, aksi, ket string, _ time.Time) error {
	g.jejak = append(g.jejak, akun+"|"+baris+"|"+aksi+"|"+ket)
	return nil
}

type tahunKontrakUji struct{}

func (tahunKontrakUji) Ambil(_ context.Context, id string) (models.TahunTreaty, error) {
	if id != "1000001" {
		return models.TahunTreaty{}, repository.ErrTahunTreatyTidakAda
	}
	return models.TahunTreaty{ID: id, TreatyYear: "2026"}, nil
}

type jenisKontrakUji []repository.JenisReasuransiTCO

func (j jenisKontrakUji) DaftarNonLife(context.Context) ([]repository.JenisReasuransiTCO, error) {
	return j, nil
}

func layananKontrak(g *gudangKontrakUji) *services.KontrakTreatyTCO {
	return services.New(nil).KontrakTreatyTCO().DenganGudang(g).DenganTahun(tahunKontrakUji{}).
		DenganJenis(jenisKontrakUji{{ID: "10003", Note: "UJI QUOTA SHARE", Tipe: "1"}}).
		DenganTransaksi(transaksiUji).DenganJam(jamUji)
}

func kontrakMasuk() services.KontrakMasuk {
	return services.KontrakMasuk{ReinsTypeID: "10003", ReinsTypeName: "KARANGAN KLIEN",
		TreatyStartDate: "2026-01-01", TreatyEndDate: "2027-01-01"}
}

func gudangKontrakKosong() *gudangKontrakUji {
	return &gudangKontrakUji{baris: map[string]models.KontrakTreaty{}}
}

func TestKontrakTanpaIdentitasDitolak(t *testing.T) {
	l := layananKontrak(gudangKontrakKosong())
	if _, err := l.Daftar(context.Background(), services.Pelaku{}, "1000001"); !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Errorf("daftar: %v", err)
	}
	if _, err := l.Simpan(context.Background(), services.Pelaku{}, "1000001", kontrakMasuk()); !errors.Is(err, services.ErrTanpaIdentitas) {
		t.Errorf("simpan: %v", err)
	}
}

func TestKontrakBawaanGagalTerang(t *testing.T) {
	_, err := services.New(nil).KontrakTreatyTCO().Daftar(context.Background(), pelakuUjiTCO, "1000001")
	if !errors.Is(err, services.ErrGudangKontrakBelumDisuntik) {
		t.Errorf("bawaan: %v", err)
	}
}

// AC 7, 5, 6, 10, 41: kontrak baru di dalam tahun; nama jenis dari MASTER,
// bukan dari klien; UserID pelaku; jejak di transaksi yang sama.
func TestKontrakSimpanBaru(t *testing.T) {
	g := gudangKontrakKosong()
	k, err := layananKontrak(g).Simpan(context.Background(), pelakuUjiTCO, "1000001", kontrakMasuk())
	if err != nil {
		t.Fatal(err)
	}
	if k.ID != "1000001" || k.IDTreatyYear != "1000001" || k.ReinsTypeName != "UJI QUOTA SHARE" ||
		k.UserID != "UJI-ADMIN" || k.TglUpdate != "2026-09-28 10:00:00" || k.TreatyStartDate != "2026-01-01" {
		t.Errorf("hasil: %+v", k)
	}
	if g.disisip != 1 || len(g.jejak) != 1 || !strings.Contains(g.jejak[0], "|simpan|kontrak baru") {
		t.Errorf("sisip %d jejak %v", g.disisip, g.jejak)
	}
}

// AC 8: ID terisi = pembaruan; kontrak tahun lain tidak dapat diperbarui lewat tahun ini.
func TestKontrakSimpanPerbarui(t *testing.T) {
	g := gudangKontrakKosong()
	g.baris["1000009"] = models.KontrakTreaty{ID: "1000009", IDTreatyYear: "1000001", ReinsTypeID: "10003"}
	g.baris["1000008"] = models.KontrakTreaty{ID: "1000008", IDTreatyYear: "1000002", ReinsTypeID: "10003"}
	m := kontrakMasuk()
	m.ID = "1000009"
	if _, err := layananKontrak(g).Simpan(context.Background(), pelakuUjiTCO, "1000001", m); err != nil {
		t.Fatal(err)
	}
	if g.disisip != 0 || g.diperbaru != 1 || !strings.Contains(g.jejak[0], "kontrak diperbarui") {
		t.Errorf("sisip %d perbarui %d jejak %v", g.disisip, g.diperbaru, g.jejak)
	}
	m.ID = "1000008"
	if _, err := layananKontrak(g).Simpan(context.Background(), pelakuUjiTCO, "1000001", m); !errors.Is(err, services.ErrKontrakTidakAda) {
		t.Errorf("kontrak tahun lain: %v", err)
	}
}

// Gerbang sebelum transaksi: nol tulisan untuk setiap penolakan.
func TestKontrakGerbang(t *testing.T) {
	kasus := []struct {
		tahun string
		ubah  func(*services.KontrakMasuk)
		mau   error
	}{
		{"9999999", func(*services.KontrakMasuk) {}, services.ErrTahunTreatyTidakAda},
		{"1000001", func(m *services.KontrakMasuk) { m.ReinsTypeID = "" }, models.ErrKontrakJenisReasuransiKosong},
		{"1000001", func(m *services.KontrakMasuk) { m.ReinsTypeID = "99999" }, services.ErrJenisReasuransiDiLuarDaftar},
		{"1000001", func(m *services.KontrakMasuk) { m.TreatyEndDate = "2025-12-31" }, models.ErrPeriodeTerbalik},
		{"1000001", func(m *services.KontrakMasuk) {
			m.TreatyStartDate, m.TreatyEndDate = "2027-01-01", "2028-01-01"
		}, models.ErrKontrakTahunMulaiBeda},
		{"1000001", func(m *services.KontrakMasuk) { m.TreatyStartDate = "01/01/2026" }, services.ErrPermintaanTidakSah},
		{"1000001", func(m *services.KontrakMasuk) { m.TreatyEndDate = "" }, models.ErrKontrakAkhirKosong},
	}
	for _, k := range kasus {
		g := gudangKontrakKosong()
		m := kontrakMasuk()
		k.ubah(&m)
		_, err := layananKontrak(g).Simpan(context.Background(), pelakuUjiTCO, k.tahun, m)
		if !errors.Is(err, k.mau) {
			t.Errorf("mau %v, dapat %v", k.mau, err)
		}
		if g.disisip+g.diperbaru+len(g.jejak) != 0 {
			t.Errorf("%v: tetap menulis", k.mau)
		}
	}
}

// Satu jenis reasuransi satu kontrak per tahun: 409 "Data sudah pernah di Input".
func TestKontrakDobelDitolak(t *testing.T) {
	g := gudangKontrakKosong()
	g.dobel = "1000004"
	_, err := layananKontrak(g).Simpan(context.Background(), pelakuUjiTCO, "1000001", kontrakMasuk())
	if !errors.Is(err, services.ErrKontrakDobel) {
		t.Fatalf("mau ErrKontrakDobel, dapat %v", err)
	}
	for _, mau := range []string{"Data sudah pernah di Input", "1000004", "10003"} {
		if !strings.Contains(err.Error(), mau) {
			t.Errorf("pesan tanpa %q: %s", mau, err)
		}
	}
	if g.disisip+len(g.jejak) != 0 {
		t.Error("dobel tetap menulis")
	}
}

// SetTanggalTreatyContract: tanggal akhir bawaan dihitung dari tahun induknya.
func TestKontrakAkhirBawaan(t *testing.T) {
	l := layananKontrak(gudangKontrakKosong())
	akhir, err := l.AkhirBawaan(context.Background(), pelakuUjiTCO, "1000001", "2026-03-01")
	if err != nil || akhir != "2027-03-01" {
		t.Errorf("akhir bawaan %q %v", akhir, err)
	}
	if _, err := l.AkhirBawaan(context.Background(), pelakuUjiTCO, "1000001", "kemarin"); !errors.Is(err, services.ErrPermintaanTidakSah) {
		t.Errorf("mulai bukan tanggal: %v", err)
	}
	if _, err := l.AkhirBawaan(context.Background(), pelakuUjiTCO, "9999999", "2026-03-01"); !errors.Is(err, services.ErrTahunTreatyTidakAda) {
		t.Errorf("tahun tidak ada: %v", err)
	}
}

// Temuan /code-review: jenis reasuransi kontrak = kunci kombinasi anak-anaknya.
func TestKontrakGantiJenisBeranakDitolak(t *testing.T) {
	g := gudangKontrakKosong()
	g.baris["1000003"] = models.KontrakTreaty{ID: "1000003", IDTreatyYear: "1000001", ReinsTypeID: "10003"}
	g.anak = 2
	l := services.New(nil).KontrakTreatyTCO().DenganGudang(g).DenganTahun(tahunKontrakUji{}).
		DenganJenis(jenisKontrakUji{{ID: "10003", Note: "UJI QS", Tipe: "1"}, {ID: "10005", Note: "UJI SURPLUS", Tipe: "2"}}).
		DenganTransaksi(transaksiUji).DenganJam(jamUji)
	m := kontrakMasuk()
	m.ID, m.ReinsTypeID = "1000003", "10005"
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", m); !errors.Is(err, services.ErrKontrakBeranak) || g.diperbaru != 0 {
		t.Errorf("ganti jenis beranak: %v (diperbarui %d)", err, g.diperbaru)
	}
	g.anak = 0
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", m); err != nil || g.diperbaru != 1 {
		t.Errorf("ganti jenis tanpa anak: %v", err)
	}
}
