package services_test

// Uji layanan reinsurer - TANPA Oracle (tiket 05).

import (
	"context"
	"errors"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/modul/treatycontractout/models"
	"nusantarare/modul/treatycontractout/repository"
	"nusantarare/modul/treatycontractout/services"
)

type gudangReinsurerUji struct {
	baris   map[string]models.ReinsurerTreaty
	urut    int
	dikunci int
}

func (g *gudangReinsurerUji) cocok(k models.KombinasiTCO, r models.ReinsurerTreaty) bool {
	return r.TreatyYear == k.TreatyYear && r.TreatyGroupID == k.TreatyGroupID && r.ReinsTypeID == k.ReinsTypeID
}
func (g *gudangReinsurerUji) Daftar(_ context.Context, k models.KombinasiTCO) ([]models.ReinsurerTreaty, error) {
	var hasil []models.ReinsurerTreaty
	for _, r := range g.baris {
		if g.cocok(k, r) {
			hasil = append(hasil, r)
		}
	}
	sort.Slice(hasil, func(i, j int) bool { return hasil[i].ID < hasil[j].ID })
	return hasil, nil
}
func (g *gudangReinsurerUji) Ambil(_ context.Context, k models.KombinasiTCO, id string) (models.ReinsurerTreaty, error) {
	r, ada := g.baris[id]
	if !ada || !g.cocok(k, r) {
		return models.ReinsurerTreaty{}, repository.ErrReinsurerTidakAda
	}
	return r, nil
}
func (g *gudangReinsurerUji) ShareLain(_ context.Context, _ *db.Tx, k models.KombinasiTCO, kecuali string) ([]*apd.Decimal, error) {
	var hasil []*apd.Decimal
	for id, r := range g.baris {
		if g.cocok(k, r) && id != kecuali {
			hasil = append(hasil, r.PctShare)
		}
	}
	return hasil, nil
}
func (g *gudangReinsurerUji) Sisip(_ context.Context, _ *db.Tx, r models.ReinsurerTreaty) (string, error) {
	g.urut++
	r.ID = "100000" + string(rune('0'+g.urut))
	g.baris[r.ID] = r
	return r.ID, nil
}
func (g *gudangReinsurerUji) Perbarui(_ context.Context, _ *db.Tx, r models.ReinsurerTreaty) error {
	if _, ada := g.baris[r.ID]; !ada {
		return repository.ErrReinsurerTidakAda
	}
	g.baris[r.ID] = r
	return nil
}

type kontrakPemegangUji struct{ dikunci *int }

func (k kontrakPemegangUji) Ambil(_ context.Context, tahunID, id string) (models.KontrakTreaty, error) {
	if tahunID != "1000001" || id != "1000003" {
		return models.KontrakTreaty{}, repository.ErrKontrakTidakAda
	}
	return models.KontrakTreaty{ID: id, IDTreatyYear: tahunID, ReinsTypeID: "10003", ReinsTypeName: "UJI QS"}, nil
}
func (k kontrakPemegangUji) Kunci(_ context.Context, _ *db.Tx, _, _ string) error {
	*k.dikunci++
	return nil
}

type tahunReinsurerUji struct{}

func (tahunReinsurerUji) Ambil(_ context.Context, id string) (models.TahunTreaty, error) {
	if id != "1000001" {
		return models.TahunTreaty{}, repository.ErrTahunTreatyTidakAda
	}
	return models.TahunTreaty{ID: id, TreatyYear: "2026", TreatyGroupID: "10001", TreatyGroupName: "UJI GRUP"}, nil
}

type masterReinsurerUji struct{}

func (masterReinsurerUji) Cari(_ context.Context, teks string) ([]repository.ReinsurerMasterTCO, error) {
	return []repository.ReinsurerMasterTCO{{ID: "UJI-R1", ClientName: "UJI REAS " + teks, ClientID: "UJI-C1"}}, nil
}
func (masterReinsurerUji) Ambil(_ context.Context, id string) (repository.ReinsurerMasterTCO, error) {
	switch id {
	case "UJI-R1":
		return repository.ReinsurerMasterTCO{ID: id, ClientName: "UJI REAS SATU", ClientID: "UJI-C1"}, nil
	case "UJI-R2":
		return repository.ReinsurerMasterTCO{ID: id, ClientName: "UJI REAS DUA", ClientID: "UJI-C2"}, nil
	}
	return repository.ReinsurerMasterTCO{}, repository.ErrReinsurerMasterTidakAda
}

func layananReinsurer(g *gudangReinsurerUji) *services.ReinsurerTCO {
	return services.New(nil).ReinsurerTCO().DenganGudang(g).DenganKontrak(kontrakPemegangUji{dikunci: &g.dikunci}).
		DenganTahun(tahunReinsurerUji{}).DenganMaster(masterReinsurerUji{}).
		DenganTransaksi(transaksiUji).DenganCatat(func(string) {})
}

func gudangReinsurerKosong() *gudangReinsurerUji {
	return &gudangReinsurerUji{baris: map[string]models.ReinsurerTreaty{}}
}

func reinsurerMasuk(reas, share, komisi string) services.ReinsurerMasuk {
	return services.ReinsurerMasuk{ReinsurerID: reas, PctShare: share, Ricomm: komisi, StdRating: "A"}
}

func TestReinsurerTanpaIdentitasDitolak(t *testing.T) {
	l := layananReinsurer(gudangReinsurerKosong())
	if _, err := l.Daftar(context.Background(), inti.Pelaku{}, "1000001", "1000003"); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("daftar: %v", err)
	}
	if _, err := l.Simpan(context.Background(), inti.Pelaku{}, "1000001", "1000003",
		reinsurerMasuk("UJI-R1", "10", "5")); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("simpan: %v", err)
	}
}

func TestReinsurerBawaanGagalTerang(t *testing.T) {
	_, err := services.New(nil).ReinsurerTCO().Daftar(context.Background(), pelakuUjiTCO, "1000001", "1000003")
	if err == nil || !strings.Contains(err.Error(), "is not injected") {
		t.Errorf("bawaan: %v", err)
	}
}

// AC 14, 16, 6, 41: reinsurer baru pada KOMBINASI kontrak; share/komisi
// desimal persis; nama & client dari master; jejak; kontrak dikunci.
func TestReinsurerSimpanBaruPadaKombinasi(t *testing.T) {
	g := gudangReinsurerKosong()
	h, err := layananReinsurer(g).Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		reinsurerMasuk("UJI-R1", "33,33333333", "12,5"))
	if err != nil {
		t.Fatal(err)
	}
	r := h.Reinsurer
	if r.TreatyYear != "2026" || r.TreatyGroupID != "10001" || r.ReinsTypeID != "10003" || r.ReinsTypeName != "UJI QS" ||
		r.Name != "UJI REAS SATU" || r.ClientID != "UJI-C1" || r.PctShare != "33.33333333" || r.Ricomm != "12.5" ||
		r.OperatorName != "UJI-ADMIN" || r.UserID != "" || r.TglUpdate != "" {
		t.Errorf("hasil: %+v", r)
	}
	if h.TotalShare != "33.33333333" || g.dikunci != 1 {
		t.Errorf("total %s, kunci %d", h.TotalShare, g.dikunci)
	}
}

// AC 15: total share terlihat, dan dijumlah persis.
func TestReinsurerDaftarDenganTotal(t *testing.T) {
	g := gudangReinsurerKosong()
	l := layananReinsurer(g)
	for _, s := range []string{"33.33333333", "33.33333333", "33.33333334"} {
		if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", reinsurerMasuk("UJI-R1", s, "1")); err != nil {
			t.Fatal(err)
		}
	}
	d, err := l.Daftar(context.Background(), pelakuUjiTCO, "1000001", "1000003")
	if err != nil {
		t.Fatal(err)
	}
	if d.Total != 3 || d.TotalShare != "100.00000000" || d.Kombinasi.TreatyYear != "2026" || d.Kombinasi.ReinsTypeID != "10003" {
		t.Errorf("daftar: total %d share %s kombinasi %+v", d.Total, d.TotalShare, d.Kombinasi)
	}
}

// SaveTreatyReinsurerDetail1_Act b1357/b1452: total BARU > 100 ditolak; mengubah
// baris sendiri tidak menghitung share lamanya dua kali.
func TestReinsurerTotalLebihDari100Ditolak(t *testing.T) {
	g := gudangReinsurerKosong()
	l := layananReinsurer(g)
	a, _ := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", reinsurerMasuk("UJI-R1", "60", "1"))
	_, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", reinsurerMasuk("UJI-R2", "40.00000001", "1"))
	if !errors.Is(err, models.ErrTotalShareMelebihi100) || !strings.Contains(err.Error(), "Percentage cannot be more than 100!") {
		t.Fatalf("mau total > 100 ditolak, dapat %v", err)
	}
	if len(g.baris) != 1 {
		t.Error("penolakan tetap menulis")
	}
	// Ubah baris 60 menjadi 100: share lamanya TIDAK ikut dijumlah.
	m := reinsurerMasuk("UJI-R1", "100", "1")
	m.ID = a.Reinsurer.ID
	if h, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", m); err != nil || h.TotalShare != "100" {
		t.Errorf("ubah ke 100: %+v %v", h, err)
	}
}

func TestReinsurerGerbang(t *testing.T) {
	kasus := []struct {
		tahun, kontrak string
		masuk          services.ReinsurerMasuk
		mau            error
	}{
		{"1000001", "1000003", reinsurerMasuk("", "10", "1"), models.ErrReinsurerKosong},
		{"1000001", "1000003", reinsurerMasuk("UJI-R9", "10", "1"), services.ErrReinsurerDiLuarMaster},
		{"1000001", "1000003", reinsurerMasuk("UJI-R1", "", "1"), models.ErrPersenKosong},
		{"1000001", "1000003", reinsurerMasuk("UJI-R1", "10", "101"), models.ErrPersenDiLuarRentang},
		{"1000001", "1000003", reinsurerMasuk("UJI-R1", "1.000,5", "1"), models.ErrPersenBukanDesimal},
		{"1000001", "1999999", reinsurerMasuk("UJI-R1", "10", "1"), services.ErrKontrakTidakAda},
		{"1999999", "1000003", reinsurerMasuk("UJI-R1", "10", "1"), services.ErrTahunTreatyTidakAda},
	}
	for _, k := range kasus {
		g := gudangReinsurerKosong()
		_, err := layananReinsurer(g).Simpan(context.Background(), pelakuUjiTCO, k.tahun, k.kontrak, k.masuk)
		if !errors.Is(err, k.mau) {
			t.Errorf("mau %v, dapat %v", k.mau, err)
		}
		if len(g.baris) != 0 {
			t.Errorf("%v: tetap menulis", k.mau)
		}
	}
}

// Medan tersembunyi (IUDate, StatusOn, UserID) DIPERTAHANKAN saat diubah
// (`SetUbahTreatyReinsurerList_Act` b1097-b1177); OperatorName = pengubah
// terakhir; TglUpdate tidak diisi (OQ-TCO-25). StartDate/EndDate TIDAK dibawa:
// repository selalu mengikat NULL (OQ-TCO-01), jadi jawaban pun kosong
// (temuan /code-review lanjutan 4).
func TestReinsurerPerbaruiMempertahankanMedanTersembunyi(t *testing.T) {
	g := gudangReinsurerKosong()
	mulai := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	g.baris["1000009"] = models.ReinsurerTreaty{ID: "1000009", TreatyYear: "2026", TreatyGroupID: "10001",
		ReinsTypeID: "10003", ReinsurerID: "UJI-R1", IUDate: "UJI-IU", StartDate: mulai, StatusOn: "1",
		UserID: "UJI-PEMBUAT", PctShare: apd.New(10, 0)}
	m := reinsurerMasuk("UJI-R2", "20", "2")
	m.ID = "1000009"
	h, err := layananReinsurer(g).Simpan(context.Background(), inti.Pelaku{AkunID: "UJI-PENGUBAH"}, "1000001", "1000003", m)
	if err != nil {
		t.Fatal(err)
	}
	r := g.baris["1000009"]
	if r.IUDate != "UJI-IU" || !r.StartDate.IsZero() || r.StatusOn != "1" || r.UserID != "UJI-PEMBUAT" ||
		r.OperatorName != "UJI-PENGUBAH" || r.Name != "UJI REAS DUA" || h.TotalShare != "20" || !r.TglUpdate.IsZero() {
		t.Errorf("perbarui: %+v total %s", r, h.TotalShare)
	}
	// Reinsurer kombinasi LAIN tidak dapat diubah lewat kontrak ini.
	g.baris["1000008"] = models.ReinsurerTreaty{ID: "1000008", TreatyYear: "2025", TreatyGroupID: "10001", ReinsTypeID: "10003"}
	m.ID = "1000008"
	if _, err := layananReinsurer(g).Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", m); !errors.Is(err, services.ErrReinsurerTidakAda) {
		t.Errorf("kombinasi lain: %v", err)
	}
}

func TestReinsurerCariMaster(t *testing.T) {
	d, err := layananReinsurer(gudangReinsurerKosong()).CariMaster(context.Background(), pelakuUjiTCO, "ab")
	if err != nil || len(d) != 1 || d[0].ClientName != "UJI REAS ab" || d[0].ClientID != "UJI-C1" {
		t.Errorf("cari: %+v %v", d, err)
	}
}

// OQ-TCO-25 (lanjutan 4): USERID/TGLUPDATE reinsurer TIDAK diisi layanan, seperti
// Pega (`NewTreatyReinsurerDetail_Act` b917 mengosongkan UserId; tidak ada langkah
// yang mengisi TglUpdate; data DEV 0/430 terisi). Baris lama yang kosong tetap kosong.
func TestReinsurerKolomPelakuKosongSepertiPega(t *testing.T) {
	g := gudangReinsurerKosong()
	g.baris["1000009"] = models.ReinsurerTreaty{ID: "1000009", TreatyYear: "2026", TreatyGroupID: "10001",
		ReinsTypeID: "10003", ReinsurerID: "UJI-R1", PctShare: apd.New(10, 0)}
	m := reinsurerMasuk("UJI-R1", "20", "2")
	m.ID = "1000009"
	if _, err := layananReinsurer(g).Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", m); err != nil {
		t.Fatal(err)
	}
	if r := g.baris["1000009"]; r.UserID != "" || !r.TglUpdate.IsZero() || r.OperatorName != "UJI-ADMIN" {
		t.Errorf("kolom pelaku: UserID %q TglUpdate %v OperatorName %q", r.UserID, r.TglUpdate, r.OperatorName)
	}
}

// OQ-TCO-25: pelaku tidak ditulis ke kolom warisan, tetapi tercatat di log aplikasi.
func TestReinsurerSimpanMencatatPelakuDiLog(t *testing.T) {
	var log []string
	l := layananReinsurer(gudangReinsurerKosong()).DenganCatat(func(s string) { log = append(log, s) })
	h, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", reinsurerMasuk("UJI-R1", "10", "1"))
	if err != nil {
		t.Fatal(err)
	}
	if len(log) != 1 || !strings.Contains(log[0], "reinsurer "+h.Reinsurer.ID) || !strings.Contains(log[0], "UJI-ADMIN") {
		t.Errorf("log: %q", log)
	}
}
