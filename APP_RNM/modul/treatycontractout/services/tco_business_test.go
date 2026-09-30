package services_test

// Uji layanan business - TANPA Oracle (tiket 07).

import (
	"context"
	"errors"
	"sort"
	"strings"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/inti/backend/db"
	"nusantarare/modul/treatycontractout/models"
	"nusantarare/modul/treatycontractout/repository"
	"nusantarare/modul/treatycontractout/services"
)

type gudangBusinessUji struct {
	baris map[string]models.BusinessTreaty
	urut  int
	dobel string
}

func (g *gudangBusinessUji) cocok(k models.KombinasiTCO, b models.BusinessTreaty) bool {
	return b.TreatyYear == k.TreatyYear && b.TreatyGroupID == k.TreatyGroupID && b.ReinsTypeID == k.ReinsTypeID
}
func (g *gudangBusinessUji) Daftar(_ context.Context, k models.KombinasiTCO) ([]models.BusinessTreaty, error) {
	var hasil []models.BusinessTreaty
	for _, b := range g.baris {
		if g.cocok(k, b) {
			hasil = append(hasil, b)
		}
	}
	sort.Slice(hasil, func(i, j int) bool { return hasil[i].ID < hasil[j].ID })
	return hasil, nil
}
func (g *gudangBusinessUji) Ambil(_ context.Context, k models.KombinasiTCO, id string) (models.BusinessTreaty, error) {
	b, ada := g.baris[id]
	if !ada || !g.cocok(k, b) {
		return models.BusinessTreaty{}, repository.ErrBusinessTidakAda
	}
	return b, nil
}
func (g *gudangBusinessUji) CariDobel(_ context.Context, _ *db.Tx, _ models.KombinasiTCO, _, _ string) (string, error) {
	return g.dobel, nil
}
func (g *gudangBusinessUji) Sisip(_ context.Context, _ *db.Tx, b models.BusinessTreaty) (string, error) {
	g.urut++
	b.ID = "100000" + string(rune('0'+g.urut))
	g.baris[b.ID] = b
	return b.ID, nil
}
func (g *gudangBusinessUji) Perbarui(_ context.Context, _ *db.Tx, b models.BusinessTreaty) error {
	if _, ada := g.baris[b.ID]; !ada {
		return repository.ErrBusinessTidakAda
	}
	g.baris[b.ID] = b
	return nil
}
func (g *gudangBusinessUji) Hapus(_ context.Context, _ *db.Tx, k models.KombinasiTCO, id string) error {
	b, ada := g.baris[id]
	if !ada || !g.cocok(k, b) {
		return repository.ErrBusinessTidakAda
	}
	delete(g.baris, id)
	return nil
}

type masterBusinessUji struct{}

func (masterBusinessUji) Daftar(context.Context) ([]repository.BusinessMasterTCO, error) {
	return []repository.BusinessMasterTCO{{ID: "UJI-B1", Note: "UJI BISNIS SATU"}, {ID: "UJI-B2", Note: "UJI BISNIS DUA"}}, nil
}
func (masterBusinessUji) Ambil(_ context.Context, id string) (repository.BusinessMasterTCO, error) {
	switch id {
	case "UJI-B1":
		return repository.BusinessMasterTCO{ID: id, Note: "UJI BISNIS SATU"}, nil
	case "UJI-B2":
		return repository.BusinessMasterTCO{ID: id, Note: "UJI BISNIS DUA"}, nil
	}
	return repository.BusinessMasterTCO{}, repository.ErrBusinessMasterTidakAda
}

func layananBusiness(g *gudangBusinessUji) *services.BusinessTCO {
	dikunci := 0
	return services.New(nil).BusinessTCO().DenganGudang(g).DenganKontrak(kontrakPemegangUji{dikunci: &dikunci}).
		DenganTahun(tahunReinsurerUji{}).DenganMaster(masterBusinessUji{}).
		DenganTransaksi(transaksiUji).DenganCatat(func(string) {})
}

func gudangBusinessKosong() *gudangBusinessUji {
	return &gudangBusinessUji{baris: map[string]models.BusinessTreaty{}}
}

// AC 21, 6, 41: kode dan nama dari master; TreatyYearID = tahun induk; jejak.
func TestBusinessSimpanBaru(t *testing.T) {
	g := gudangBusinessKosong()
	h, err := layananBusiness(g).Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		services.BusinessMasuk{BizCode: "UJI-B1", BizName: "KARANGAN", IsActive: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if h.BizName != "UJI BISNIS SATU" || h.TreatyYearID != "1000001" || h.TreatyYear != "2026" ||
		h.ReinsTypeID != "10003" || h.UserID != "" || h.TglUpdate != "" || !h.Aktif {
		t.Errorf("hasil: %+v", h)
	}
}

// AC 22: dinonaktifkan tanpa dihapus, tetap terbaca sebagai nonaktif.
// AC 23: SELURUH medan yang dikirim tersimpan.
func TestBusinessNonaktifkanDanPerbaruiSeluruhMedan(t *testing.T) {
	g := gudangBusinessKosong()
	l := layananBusiness(g)
	a, _ := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", services.BusinessMasuk{BizCode: "UJI-B1", IsActive: "1"})
	_, err := l.Simpan(context.Background(), inti.Pelaku{AkunID: "UJI-PENGUBAH"}, "1000001", "1000003",
		services.BusinessMasuk{ID: a.ID, BizCode: "UJI-B2", IsActive: "0"})
	if err != nil {
		t.Fatal(err)
	}
	b := g.baris[a.ID]
	// OQ-TCO-25: USERID/TGLUPDATE tidak diisi, seperti Pega (data DEV 2/4.621 terisi).
	if b.IsActive != "0" || b.BizCode != "UJI-B2" || b.BizName != "UJI BISNIS DUA" || b.UserID != "" || !b.TglUpdate.IsZero() ||
		b.TreatyYearID != "1000001" || b.TreatyGroupName != "UJI GRUP" || b.ReinsTypeName != "UJI QS" {
		t.Errorf("perbarui tidak menulis seluruh medan: %+v", b)
	}
	d, err := l.Daftar(context.Background(), pelakuUjiTCO, "1000001", "1000003")
	if err != nil || d.Total != 1 || d.Daftar[0].Aktif || d.Daftar[0].IsActive != "0" {
		t.Errorf("baris nonaktif tidak terbaca sebagai nonaktif: %+v %v", d, err)
	}
}

func TestBusinessGerbang(t *testing.T) {
	kasus := []struct {
		masuk services.BusinessMasuk
		mau   error
	}{
		{services.BusinessMasuk{BizCode: "", IsActive: "1"}, models.ErrBusinessKodeKosong},
		{services.BusinessMasuk{BizCode: "UJI-B1", IsActive: ""}, models.ErrBusinessAktifTakSah},
		{services.BusinessMasuk{BizCode: "UJI-B9", IsActive: "1"}, services.ErrBusinessDiLuarMaster},
	}
	for _, k := range kasus {
		g := gudangBusinessKosong()
		_, err := layananBusiness(g).Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", k.masuk)
		if !errors.Is(err, k.mau) {
			t.Errorf("mau %v, dapat %v", k.mau, err)
		}
		if len(g.baris) != 0 {
			t.Errorf("%v: tetap menulis", k.mau)
		}
	}
	g := gudangBusinessKosong()
	g.dobel = "1000004"
	_, err := layananBusiness(g).Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		services.BusinessMasuk{BizCode: "UJI-B1", IsActive: "1"})
	if !errors.Is(err, services.ErrBusinessDobel) || !strings.Contains(err.Error(), "Data has already been entered") {
		t.Errorf("dobel: %v", err)
	}
}

// AC 63/64: hapus satu baris, satu tabel; pesan VERBATIM DeleteRowBusiness langkah 3.
func TestBusinessHapus(t *testing.T) {
	g := gudangBusinessKosong()
	l := layananBusiness(g)
	a, _ := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", services.BusinessMasuk{BizCode: "UJI-B1", IsActive: "1"})
	pesan, err := l.Hapus(context.Background(), pelakuUjiTCO, "1000001", "1000003", a.ID)
	if err != nil || pesan != "Data with ID "+a.ID+" successfully deleted" {
		t.Errorf("hapus: %q %v", pesan, err)
	}
	if len(g.baris) != 0 {
		t.Errorf("baris %d", len(g.baris))
	}
	if _, err := l.Hapus(context.Background(), pelakuUjiTCO, "1000001", "1000003", a.ID); !errors.Is(err, services.ErrBusinessTidakAda) {
		t.Errorf("hapus kedua: %v", err)
	}
}

func TestBusinessMaster(t *testing.T) {
	d, err := layananBusiness(gudangBusinessKosong()).Master(context.Background(), pelakuUjiTCO)
	if err != nil || len(d) != 2 || d[0].ID != "UJI-B1" {
		t.Errorf("master: %+v %v", d, err)
	}
}

// tco4 (temuan /code-review): UPDATE procedure hanya lima kolom - jawaban simpan
// memuat kolom yang benar-benar tersimpan, bukan nilai kombinasi saat ini.
func TestBusinessPerbaruiMenjawabYangTersimpan(t *testing.T) {
	g := gudangBusinessKosong()
	l := layananBusiness(g)
	a, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", services.BusinessMasuk{BizCode: "UJI-B1", IsActive: "1"})
	if err != nil {
		t.Fatal(err)
	}
	lama := g.baris[a.ID]
	lama.TreatyYearID, lama.TreatyGroupName = "", "UJI GRUP LAMA" // baris warisan: TREATYYEARID NULL
	g.baris[a.ID] = lama
	b, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		services.BusinessMasuk{ID: a.ID, BizCode: "UJI-B1", IsActive: "0"})
	if err != nil {
		t.Fatal(err)
	}
	if b.TreatyYearID != "" || b.TreatyGroupName != "UJI GRUP LAMA" || b.Aktif {
		t.Errorf("jawaban bukan yang tersimpan: %+v", b)
	}
}

// OQ-TCO-25: pelaku tidak ditulis ke kolom warisan, tetapi tercatat di log aplikasi.
func TestBusinessSimpanMencatatPelakuDiLog(t *testing.T) {
	var log []string
	l := layananBusiness(gudangBusinessKosong()).DenganCatat(func(s string) { log = append(log, s) })
	h, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", services.BusinessMasuk{BizCode: "UJI-B1", IsActive: "1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(log) != 1 || !strings.Contains(log[0], "business "+h.ID) || !strings.Contains(log[0], "UJI-ADMIN") {
		t.Errorf("log: %q", log)
	}
}
