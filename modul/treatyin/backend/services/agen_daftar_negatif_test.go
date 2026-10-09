package services_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	inti "nusantarare/inti/backend"
	"nusantarare/modul/treatyin/backend/services"
)

// Gudang tiruan bersama: nol agen tercatat.
func (g *gudangTiruan) BacaStatusAktifAgen(context.Context, []string) (map[string]string, error) {
	return map[string]string{}, nil
}

// gudangStatusAgen - `AGENT.STATUSACTIVE` tiruan, dan pengenal yang diminta.
type gudangStatusAgen struct {
	*gudangTiruan
	status map[string]string
	galat  error
	minta  []string
}

func (g *gudangStatusAgen) BacaStatusAktifAgen(_ context.Context, ids []string) (map[string]string, error) {
	g.minta = ids
	return g.status, g.galat
}

func TestDaftarNegatifHanyaUntukStatusInactiveApaAdanya(t *testing.T) {
	g := &gudangStatusAgen{gudangTiruan: &gudangTiruan{}, status: map[string]string{
		"G1": "inactive", "G2": "1", "G3": "0", "G4": "Inactive",
	}}
	l := services.LayananDengan(g)
	kasus := []struct {
		id      string
		negatif bool
	}{
		{"G1", true},
		// ⛔ `'0'` BUKAN "inactive" — menafsirkannya keputusan pemilik proses.
		{"G3", false},
		{"G2", false},
		// `==` Pega persis: huruf besar tidak dilipat.
		{"G4", false},
		// Tidak ada di AGENT.
		{"G9", false},
	}
	for _, k := range kasus {
		h, err := l.PeriksaDaftarNegatifAgen(context.Background(), pelakuAda, k.id, "")
		if err != nil {
			t.Fatalf("%s: %v", k.id, err)
		}
		if h.Cedant.DaftarNegatif != k.negatif {
			t.Errorf("%s: daftarNegatif %v, mau %v", k.id, h.Cedant.DaftarNegatif, k.negatif)
		}
		mauPesan := ""
		if k.negatif {
			mauPesan = "This name is on Agent Negative List"
		}
		if h.Cedant.Pesan != mauPesan {
			t.Errorf("%s: pesan %q, mau %q", k.id, h.Cedant.Pesan, mauPesan)
		}
		if h.AsalBisnis.DaftarNegatif || h.AsalBisnis.Pesan != "" {
			t.Errorf("%s: asal bisnis kosong tidak boleh berpesan: %+v", k.id, h.AsalBisnis)
		}
	}
}

func TestDaftarNegatifMemeriksaKeduaMedanSekaligus(t *testing.T) {
	g := &gudangStatusAgen{gudangTiruan: &gudangTiruan{}, status: map[string]string{"C1": "1", "S1": "inactive"}}
	h, err := services.LayananDengan(g).PeriksaDaftarNegatifAgen(context.Background(), pelakuAda, " C1 ", "S1")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(g.minta, []string{"C1", "S1"}) {
		t.Errorf("pengenal diminta %q", g.minta)
	}
	if h.Cedant.DaftarNegatif || !h.Cedant.Ditemukan || h.Cedant.StatusAktif != "1" {
		t.Errorf("cedant %+v", h.Cedant)
	}
	if !h.AsalBisnis.DaftarNegatif || h.AsalBisnis.Pesan != services.PesanDaftarNegatifAgen {
		t.Errorf("asal bisnis %+v", h.AsalBisnis)
	}
}

func TestDaftarNegatifMeneruskanGalatDanMenolakTanpaIdentitas(t *testing.T) {
	bocor := errors.New("oracle mati")
	g := &gudangStatusAgen{gudangTiruan: &gudangTiruan{}, galat: bocor}
	if _, err := services.LayananDengan(g).PeriksaDaftarNegatifAgen(context.Background(), pelakuAda, "C1", ""); !errors.Is(err, bocor) {
		t.Errorf("galat gudang: %v", err)
	}
	if _, err := services.LayananDengan(g).PeriksaDaftarNegatifAgen(context.Background(), inti.Pelaku{}, "C1", ""); !errors.Is(err, inti.ErrTanpaIdentitas) {
		t.Errorf("tanpa identitas: %v", err)
	}
}
