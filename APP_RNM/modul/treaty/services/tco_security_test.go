package services_test

// Uji layanan security reinsurer - TANPA Oracle (tiket 06).

import (
	"context"
	"errors"
	"sort"
	"strings"
	"testing"

	"github.com/cockroachdb/apd/v3"

	"nusantarare/inti"
	"nusantarare/inti/db"
	"nusantarare/modul/treaty/models"
	"nusantarare/modul/treaty/repository"
	"nusantarare/modul/treaty/services"
)

// gudangSecurityUji meniru MTREATYSECURITY warisan (tco4): tanpa identitas,
// baris dikunci REAS_SECURITY terpangkas (per reinsurer - uji memakai satu).
type gudangSecurityUji struct {
	baris map[string]models.SecurityReinsurer
}

func (g *gudangSecurityUji) Daftar(_ context.Context, reasID, thn string) ([]repository.SecurityTCO, error) {
	var hasil []repository.SecurityTCO
	for _, s := range g.baris {
		if s.ReasID == reasID && s.ThnTreaty == thn {
			hasil = append(hasil, repository.SecurityTCO{SecurityReinsurer: s, ClientName: "NAMA " + s.ReasSecurity})
		}
	}
	sort.Slice(hasil, func(i, j int) bool { return hasil[i].ID < hasil[j].ID })
	return hasil, nil
}
func (g *gudangSecurityUji) Ambil(_ context.Context, reasID, id string) (repository.SecurityTCO, error) {
	s, ada := g.baris[id]
	if !ada || s.ReasID != reasID {
		return repository.SecurityTCO{}, repository.ErrSecurityTidakAda
	}
	return repository.SecurityTCO{SecurityReinsurer: s}, nil
}
func (g *gudangSecurityUji) CariDobel(_ context.Context, _ *db.Tx, reasID, sec, kecuali string) (string, error) {
	for id, s := range g.baris {
		if s.ReasID == reasID && s.ReasSecurity == sec && id != kecuali {
			return id, nil
		}
	}
	return "", nil
}
func (g *gudangSecurityUji) Sisip(_ context.Context, _ *db.Tx, s models.SecurityReinsurer) (string, error) {
	s.ID = strings.TrimSpace(s.ReasSecurity)
	g.baris[s.ID] = s
	return s.ID, nil
}

// Perbarui - UpdateMTreatySecurity: berkunci nama LAMA (`s.ID`), nama baru ditulis.
func (g *gudangSecurityUji) Perbarui(_ context.Context, _ *db.Tx, s models.SecurityReinsurer) error {
	lama, ada := g.baris[s.ID]
	if !ada || lama.ReasID != s.ReasID {
		return repository.ErrSecurityTidakAda
	}
	delete(g.baris, s.ID)
	s.ID = strings.TrimSpace(s.ReasSecurity)
	g.baris[s.ID] = s
	return nil
}
func (g *gudangSecurityUji) Hapus(_ context.Context, _ *db.Tx, reasID, id string) error {
	s, ada := g.baris[id]
	if !ada || s.ReasID != reasID {
		return repository.ErrSecurityTidakAda
	}
	delete(g.baris, id)
	return nil
}

// reinsurerIndukSec - satu reinsurer pada kombinasi 2026/10001/10003.
func reinsurerIndukSec() *gudangReinsurerUji {
	return &gudangReinsurerUji{baris: map[string]models.ReinsurerTreaty{
		"1000007": {ID: "1000007", TreatyYear: "2026", TreatyGroupID: "10001", ReinsTypeID: "10003",
			ReinsurerID: "UJI-R1", Name: "UJI REAS SATU", PctShare: apd.New(40, 0)},
	}}
}

func layananSecurity(g *gudangSecurityUji, r *gudangReinsurerUji) *services.SecurityTCO {
	dikunci := 0
	return services.New(nil).SecurityTCO().DenganGudang(g).DenganReinsurer(r).
		DenganKontrak(kontrakPemegangUji{dikunci: &dikunci}).DenganTahun(tahunReinsurerUji{}).
		DenganMaster(masterReinsurerUji{}).DenganTransaksi(transaksiUji)
}

func gudangSecurityKosong() *gudangSecurityUji {
	return &gudangSecurityUji{baris: map[string]models.SecurityReinsurer{}}
}

// AC 17, 41: security di bawah reinsurer beserta porsinya; tahun dari reinsurer;
// nama dari master; jejak.
func TestSecuritySimpanBaru(t *testing.T) {
	g := gudangSecurityKosong()
	h, err := layananSecurity(g, reinsurerIndukSec()).Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003",
		"1000007", services.SecurityMasuk{ReasSecurity: " UJI-R2 ", ClientName: "KARANGAN", PctShare: "33,33333333"})
	if err != nil {
		t.Fatal(err)
	}
	if h.ReasID != "1000007" || h.ThnTreaty != "2026" || h.ReasSecurity != "UJI-R2" || h.ClientName != "UJI REAS DUA" ||
		h.PctShare != "33.33333333" {
		t.Errorf("hasil: %+v", h)
	}
	b := g.baris[h.ID]
	if b.TopID != "" || b.TpTreaty != "" || b.UserID != "" {
		t.Errorf("kolom [terbuka] tidak dibiarkan kosong seperti warisan: %+v", b)
	}
}

// User story 14: mengganti nama security menimpa baris yang SAMA - seperti
// UpdateMTreatySecurity b89 (berkunci nama LAMA). tco4: tabel tanpa identitas,
// maka kunci baris itu kini nama baru.
func TestSecurityGantiNamaTetapBarisYangSama(t *testing.T) {
	g := gudangSecurityKosong()
	l := layananSecurity(g, reinsurerIndukSec())
	a, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007",
		services.SecurityMasuk{ReasSecurity: "UJI-R1", PctShare: "10"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007",
		services.SecurityMasuk{ID: a.ID, ReasSecurity: "UJI-R2", PctShare: "15"})
	if err != nil {
		t.Fatal(err)
	}
	if a.ID != "UJI-R1" || b.ID != "UJI-R2" || len(g.baris) != 1 || g.baris["UJI-R2"].ReasID != "1000007" {
		t.Errorf("ganti nama: %+v %+v", b, g.baris)
	}
}

func TestSecurityGerbang(t *testing.T) {
	kasus := []struct {
		nama        string
		reas, masuk string
		m           services.SecurityMasuk
		mau         error
	}{
		{"security kosong", "1000007", "", services.SecurityMasuk{PctShare: "1"}, models.ErrSecurityKosong},
		{"share kosong", "1000007", "", services.SecurityMasuk{ReasSecurity: "UJI-R1"}, models.ErrPersenKosong},
		{"share > 100", "1000007", "", services.SecurityMasuk{ReasSecurity: "UJI-R1", PctShare: "100.5"}, models.ErrPersenDiLuarRentang},
		{"di luar master", "1000007", "", services.SecurityMasuk{ReasSecurity: "UJI-X", PctShare: "1"}, services.ErrReinsurerDiLuarMaster},
		{"reinsurer bukan milik kombinasi", "1000099", "", services.SecurityMasuk{ReasSecurity: "UJI-R1", PctShare: "1"}, services.ErrReinsurerTidakAda},
		{"ubah baris yang tidak ada", "1000007", "1000999", services.SecurityMasuk{ReasSecurity: "UJI-R1", PctShare: "1"}, services.ErrSecurityTidakAda},
	}
	for _, k := range kasus {
		g := gudangSecurityKosong()
		k.m.ID = k.masuk
		_, err := layananSecurity(g, reinsurerIndukSec()).Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003",
			k.reas, k.m)
		if !errors.Is(err, k.mau) {
			t.Errorf("%s: %v, mau %v", k.nama, err, k.mau)
		}
		if len(g.baris) != 0 {
			t.Errorf("%s: gagal tetapi menulis", k.nama)
		}
	}
	if _, err := layananSecurity(gudangSecurityKosong(), reinsurerIndukSec()).Simpan(context.Background(), inti.Pelaku{},
		"1000001", "1000003", "1000007", services.SecurityMasuk{ReasSecurity: "UJI-R1", PctShare: "1"}); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
}

// Security yang sama dua kali di bawah reinsurer yang sama ditolak 409 - warisan
// menyisipkannya lagi (`Local.IsUpdate` dihitung lalu tidak dipakai).
func TestSecurityDobelDitolak(t *testing.T) {
	g := gudangSecurityKosong()
	l := layananSecurity(g, reinsurerIndukSec())
	a, _ := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007",
		services.SecurityMasuk{ReasSecurity: "UJI-R1", PctShare: "10"})
	_, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007",
		services.SecurityMasuk{ReasSecurity: "UJI-R1", PctShare: "20"})
	var dobel services.GalatSecurityDobel
	if !errors.As(err, &dobel) || dobel.IDLain != a.ID || !errors.Is(err, services.ErrSecurityDobel) {
		t.Errorf("dobel: %v", err)
	}
	// Menyimpan ulang baris itu sendiri bukan dobel.
	if _, err := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007",
		services.SecurityMasuk{ID: a.ID, ReasSecurity: "UJI-R1", PctShare: "20"}); err != nil {
		t.Errorf("simpan ulang baris sendiri: %v", err)
	}
}

// Hapus SATU baris menurut ID; reinsurer induk dan security lain tidak tersentuh.
func TestSecurityHapusSatuBaris(t *testing.T) {
	g := gudangSecurityKosong()
	r := reinsurerIndukSec()
	l := layananSecurity(g, r)
	a, _ := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007",
		services.SecurityMasuk{ReasSecurity: "UJI-R1", PctShare: "10"})
	b, _ := l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007",
		services.SecurityMasuk{ReasSecurity: "UJI-R2", PctShare: "20"})
	pesan, err := l.Hapus(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007", a.ID)
	if err != nil || !strings.Contains(pesan, a.ID) {
		t.Fatalf("hapus: %q %v", pesan, err)
	}
	if _, ada := g.baris[b.ID]; !ada || len(g.baris) != 1 || len(r.baris) != 1 {
		t.Errorf("hapus menyentuh baris lain: %+v %+v", g.baris, r.baris)
	}
	if _, err := l.Hapus(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007", a.ID); !errors.Is(err, services.ErrSecurityTidakAda) {
		t.Errorf("hapus dua kali: %v", err)
	}
}

func TestSecurityDaftarPerReinsurer(t *testing.T) {
	g := gudangSecurityKosong()
	l := layananSecurity(g, reinsurerIndukSec())
	_, _ = l.Simpan(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007",
		services.SecurityMasuk{ReasSecurity: "UJI-R1", PctShare: "10"})
	g.baris["1000999"] = models.SecurityReinsurer{ID: "1000999", ReasID: "1000008", ThnTreaty: "2026", ReasSecurity: "UJI-R2"}
	d, err := l.Daftar(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007")
	if err != nil || d.Total != 1 || d.Daftar[0].ReasSecurity != "UJI-R1" || d.Reinsurer.ID != "1000007" {
		t.Errorf("daftar: %+v %v", d, err)
	}
}

func TestSecurityBelumDisuntikGagalTerang(t *testing.T) {
	_, err := services.New(nil).SecurityTCO().DenganTahun(tahunReinsurerUji{}).
		DenganKontrak(kontrakPemegangUji{dikunci: new(int)}).DenganReinsurer(reinsurerIndukSec()).
		Daftar(context.Background(), pelakuUjiTCO, "1000001", "1000003", "1000007")
	if !errors.Is(err, services.ErrGudangSecurityBelumDisuntik) {
		t.Errorf("bawaan: %v", err)
	}
}
